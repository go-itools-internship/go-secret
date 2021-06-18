// Package cmd provides functions set and get with cobra library
package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"

	"github.com/go-redis/redis/v8"

	"github.com/go-chi/chi/v5/middleware"

	api "github.com/go-itools-internship/go-secret/pkg/http"
	secretApi "github.com/go-itools-internship/go-secret/pkg/secret"

	"github.com/go-chi/chi/v5"
	"github.com/go-itools-internship/go-secret/pkg/crypto"
	"github.com/go-itools-internship/go-secret/pkg/io/storage"
	"github.com/go-itools-internship/go-secret/pkg/provider"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type root struct {
	options options
	cmd     *cobra.Command
	logger  *zap.SugaredLogger
}

type chiLogger struct {
	sugar *zap.SugaredLogger
}

func (c *chiLogger) Print(v ...interface{}) {
	c.sugar.Info(v...)
}

type options struct {
	version string
}

var defaultOptions = options{
	version: "undefined",
}

type RootOptions func(o *options)

// Version is optional function can add version flag to root
func Version(ver string) RootOptions {
	return func(o *options) {
		o.version = ver
	}
}

// Execute executes the secret commands.
func (r *root) Execute(ctx context.Context) error {
	return r.cmd.ExecuteContext(ctx)
}

// New function create commands to the cobra CLI
// RootOptions adds additional features to the cobra CLI
func New(opts ...RootOptions) *root {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	var secret = &cobra.Command{
		Use:     "secret",
		Short:   "Contains commands to set and get encrypt data to storage",
		Long:    "Create CLI to set and get secrets via the command line",
		Version: options.version,
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.DisableCaller = true

	logg, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("logger can't initilize %s", err))
	}

	rootData := &root{cmd: secret, options: options, logger: logg.Sugar()}

	secret.AddCommand(rootData.setCmd())
	secret.AddCommand(rootData.getCmd())
	secret.AddCommand(rootData.serverCmd())
	secret.SilenceUsage = true // write false if you want to see options when an error occurs

	return rootData
}

func (r *root) setCmd() *cobra.Command {
	var key string
	var cipherKey string
	var value string
	var path string
	var redisURL string
	var postgresURL string
	var migration string
	var setCmd = &cobra.Command{
		Use:   "set",
		Short: "Saves data to the specified storage in encrypted form",
		Long:  "it takes keys and a value from user and saves value in encrypted manner in specified storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			var ds secretApi.DataSaver
			var cr = crypto.NewCryptographer([]byte(cipherKey))
			logger := r.logger.Named("set-cmd")
			logger.Info("Start")
			switch {
			case redisURL != "":
				rdb := redis.NewClient(&redis.Options{Addr: redisURL, Password: "", DB: 0})
				defer disconnectRDB(rdb, logger)
				err := rdb.Ping(r.cmd.Context()).Err()
				if err != nil {
					return fmt.Errorf("redis db is not reachable:  %w", err)
				}
				ds = storage.NewRedisVault(rdb)
			case postgresURL != "":
				err := migrateUp(postgresURL, migration, logger)
				if err != nil {
					if errors.Is(err, migrate.ErrNoChange) {
						logger.Infof("can't migrate db:  %s", err)
					} else {
						return fmt.Errorf("migrate error :  %w", err)
					}
				}
				pdb, err := sqlx.ConnectContext(r.cmd.Context(), "postgres", postgresURL)
				if err != nil {
					return fmt.Errorf("postgres url is not reachable:  %w", err)
				}
				logger.Infof("pdb after connection %v", pdb)
				defer disconnectPDB(pdb, logger)
				ds = storage.NewPostgreVault(pdb)
			case path != "":
				var err error
				ds, err = storage.NewFileVault(path)
				if err != nil {
					return fmt.Errorf("can't create storage by path: %w", err)
				}
			}

			pr := provider.NewProvider(cr, ds)
			logger.Info("prepare get data by key: ", key)
			err := pr.SetData([]byte(key), []byte(value))
			logger.Info("ready get data by key: ", key)
			if err != nil {
				return fmt.Errorf("can't set data %w", err)
			}
			return nil
		},
	}
	setCmd.Flags().StringVarP(&value, "value", "v", value, "value to be encrypted")
	setCmd.Flags().StringVarP(&key, "key", "k", key, "key for pair key-value")
	setCmd.Flags().StringVarP(&cipherKey, "cipher-key", "c", cipherKey, "cipher key for data encryption and decryption")
	setCmd.Flags().StringVarP(&path, "path", "p", "file.txt", "the place where the key/value will be stored/got")
	setCmd.Flags().StringVarP(&redisURL, "redis-url", "r", "", "redis url address. Example: localhost:6379")
	setCmd.Flags().StringVarP(&postgresURL, "postgres-url", "s", "", "postgres url address. Example: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	setCmd.Flags().StringVarP(&migration, "migration", "m", "", "migration up route to scripts/migrations folder. Example: file://../../../scripts/migrations")

	return setCmd
}

func (r *root) getCmd() *cobra.Command {
	var key string
	var cipherKey string
	var path string
	var redisURL string
	var postgresURL string
	var migration string
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get data from specified storage in decrypted form",
		Long:  "it takes keys from user and get value in decrypted manner from specified storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := r.logger.Named("get-cmd")
			logger.Info("Start")
			var ds secretApi.DataSaver
			var cr = crypto.NewCryptographer([]byte(cipherKey))

			switch {
			case redisURL != "":
				rdb := redis.NewClient(&redis.Options{Addr: redisURL, Password: "", DB: 0})
				defer disconnectRDB(rdb, logger)
				err := rdb.Ping(r.cmd.Context()).Err()
				if err != nil {
					return fmt.Errorf("redis db is not reachable:  %w", err)
				}
				ds = storage.NewRedisVault(rdb)
			case postgresURL != "":
				err := migrateUp(postgresURL, migration, logger)
				if err != nil {
					if errors.Is(err, migrate.ErrNoChange) {
						logger.Infof("can't migrate db:  %s", err)
					} else {
						return fmt.Errorf("migrate error :  %w", err)
					}
				}
				pdb, err := sqlx.ConnectContext(r.cmd.Context(), "postgres", postgresURL)
				if err != nil {
					return fmt.Errorf("postgres url is not reachable:  %w", err)
				}
				logger.Infof("pdb after connection %v", pdb)
				defer disconnectPDB(pdb, logger)
				ds = storage.NewPostgreVault(pdb)
			case path != "":
				var err error
				ds, err = storage.NewFileVault(path)
				if err != nil {
					return fmt.Errorf("can't get storage by path: %w", err)
				}
			}

			pr := provider.NewProvider(cr, ds)
			logger.Info("prepare by get data by key: ", key)
			data, err := pr.GetData([]byte(key))
			if err != nil {
				return fmt.Errorf("can't get data by key: %w", err)
			}
			logger.Info("ready get data by key: ", key)
			cmd.Println(string(data))
			return nil
		},
	}
	getCmd.Flags().StringVarP(&key, "key", "k", key, "key for pair key-value")
	getCmd.Flags().StringVarP(&cipherKey, "cipher-key", "c", cipherKey, "cipher key for data encryption and decryption")
	getCmd.Flags().StringVarP(&path, "path", "p", "file.txt", "the place where the value will be got")
	getCmd.Flags().StringVarP(&redisURL, "redis-url", "r", "", "redis url address. Example: localhost:6379")
	getCmd.Flags().StringVarP(&postgresURL, "postgres-url", "s", "", "postgres url address. Example: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	getCmd.Flags().StringVarP(&migration, "migration", "m", "", "migration up route to scripts/migrations folder. Example: file://../../../scripts/migrations")

	return getCmd
}

func (r *root) serverCmd() *cobra.Command {
	var path string
	var port string
	var redisURL string
	var postgresURL string
	var migration string
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run server runner mode to start the app as a daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := r.logger.Named("server")
			store := make(map[string]api.MethodFactoryFunc)
			logger.Info("Start")
			switch {
			case redisURL != "":
				rdb := redis.NewClient(&redis.Options{Addr: redisURL, Password: "", DB: 0})
				defer disconnectRDB(rdb, logger)
				err := rdb.Ping(r.cmd.Context()).Err()
				if err != nil {
					return fmt.Errorf("redis db is not reachable:  %w", err)
				}
				dataRedis := storage.NewRedisVault(rdb)
				// remote method set handler for redis storage
				store["remote"] = func(cipher string) (secretApi.Provider, func()) {
					cr := crypto.NewCryptographer([]byte(cipher))
					return provider.NewProvider(cr, dataRedis), nil
				}
			case postgresURL != "":
				err := migrateUp(postgresURL, migration, logger)
				if err != nil {
					if errors.Is(err, migrate.ErrNoChange) {
						logger.Infof("can't migrate db:  %s", err)
					} else {
						return fmt.Errorf("migrate error :  %w", err)
					}
				}
				pdb, err := sqlx.ConnectContext(r.cmd.Context(), "postgres", postgresURL)
				if err != nil {
					return fmt.Errorf("postgres url is not reachable:  %w", err)
				}
				logger.Infof("pdb after connection %v", pdb)
				defer disconnectPDB(pdb, logger)
				dataPostgres := storage.NewPostgreVault(pdb)
				// remote method set handler for postgres storage
				store["remote"] = func(cipher string) (secretApi.Provider, func()) {
					cr := crypto.NewCryptographer([]byte(cipher))
					return provider.NewProvider(cr, dataPostgres), nil
				}
			}
			if path != "" {
				ds, err := storage.NewFileVault(path)
				if err != nil {
					return fmt.Errorf("can't get storage by path: %s", err)
				}
				store["local"] = func(cipher string) (secretApi.Provider, func()) {
					cr := crypto.NewCryptographer([]byte(cipher))
					return provider.NewProvider(cr, ds), nil
				}
			}

			handler := api.NewMethods(store, logger.Named("handler"))
			router := chi.NewRouter()
			srv := &http.Server{Addr: ":" + port, Handler: router}
			router.Use(middleware.Heartbeat("/ping"), middleware.RequestLogger(&middleware.DefaultLogFormatter{
				Logger: &chiLogger{logger.Named("api")},
			}))
			router.Post("/", handler.SetByKey)
			router.Get("/", handler.GetByKey)

			done := make(chan os.Signal, 1)
			shutdownCh := make(chan struct{})
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				logger.Infof("listening ")
				err := srv.ListenAndServe()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Errorf("connection error: %s", err)
				}
			}()
			logger.Info("listener started")

			select {
			case <-done:
				logger.Info("listener stopped")
			case <-cmd.Context().Done():
				logger.Info("listener stopped with context")
			}

			go func(ctx context.Context) {
				defer close(shutdownCh)
				ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				logger.Infof("ctx-shutdown-pdb after connection %v", ctx)
				defer cancel()
				err := srv.Shutdown(ctx)
				if err != nil {
					logger.Error("listener shutdown failed", err)
				}
			}(context.Background())
			<-shutdownCh
			logger.Info("exit")

			return nil
		},
	}
	serverCmd.Flags().StringVarP(&path, "path", "p", "file.txt", "the place where the key/value will be stored/got")
	serverCmd.Flags().StringVarP(&port, "port", "t", "8888", "localhost address")
	serverCmd.Flags().StringVarP(&redisURL, "redis-url", "r", "", "redis url address. Example: localhost:6379")
	serverCmd.Flags().StringVarP(&postgresURL, "postgres-url", "s", "", "postgres url address. Example: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable\"")
	serverCmd.Flags().StringVarP(&migration, "migration", "m", "", "migration up route to scripts/migrations folder. Example: file://../../../scripts/migrations")
	serverCmd.AddCommand(r.serverPingCmd())
	return serverCmd
}

func (r *root) serverPingCmd() *cobra.Command {
	var url string
	var port string
	var route string
	var timeout time.Duration
	var serverPingCmd = &cobra.Command{
		Use:   "ping",
		Short: "Check a health check route endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := http.Client{
				Timeout: timeout,
			}
			resp, err := client.Get(fmt.Sprintf("%s:%s%s", url, port, route))
			if err != nil {
				return fmt.Errorf("server response error: %w", err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					fmt.Println("server: can't close request body: ", err.Error())
				}
			}()

			if resp.StatusCode != http.StatusOK {
				responseBody, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("server: can't get response body %w", err)
				}
				return fmt.Errorf("server response is not expected: body %q, wrong status code %d", responseBody, resp.StatusCode)
			}
			return nil
		},
	}
	serverPingCmd.Flags().StringVarP(&port, "port", "p", "8880", "port to connect. Address shouldn't have port. Default: '8880")
	serverPingCmd.Flags().StringVarP(&route, "route", "r", "/ping", "health check route. Default: '/ping'")
	serverPingCmd.Flags().StringVarP(&url, "url", "u", "http://localhost", "url for server checking. Url shouldn't contain port. Default: 'http://localhost'")
	serverPingCmd.Flags().DurationVarP(&timeout, "timeout", "t", 15*time.Second, "max request time to make a request. Default: '15 seconds'")
	return serverPingCmd
}

func migrateUp(postgres, source string, logger *zap.SugaredLogger) error {
	logger = logger.Named("migration")
	logger.Info("starting from source=%s", source)
	m, err := migrate.New(
		source,
		postgres)
	if err != nil {
		return err
	}
	logger.Info("created")
	err = m.Up()
	if err != nil {
		return err
	}
	logger.Info("finished")
	return nil
}

func disconnectPDB(pdb *sqlx.DB, logger *zap.SugaredLogger) {
	logger = logger.Named("disconnect")
	err := pdb.Close()
	if err != nil {
		logger.Warnf("can't disconnect postgres db, error=%v", err)
		return
	}
	logger.Info("pdb disconnect")
}

func disconnectRDB(rdb *redis.Client, logger *zap.SugaredLogger) {
	logger = logger.Named("disconnect")
	err := rdb.Close()
	if err != nil {
		logger.Warnf("can't disconnect redis db, error=%v", err)
		return
	}
	logger.Info("rdb disconnected")
}
