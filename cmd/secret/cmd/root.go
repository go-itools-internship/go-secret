// Package cmd provides functions set and get with cobra library
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	api "github.com/go-itools-internship/go-secret/pkg/http"
	secretApi "github.com/go-itools-internship/go-secret/pkg/secret"

	"github.com/go-chi/chi/v5"
	"github.com/go-itools-internship/go-secret/pkg/crypto"
	"github.com/go-itools-internship/go-secret/pkg/io/storage"
	"github.com/go-itools-internship/go-secret/pkg/provider"
	"github.com/spf13/cobra"
)

type root struct {
	options options
	cmd     *cobra.Command
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

	rootData := &root{cmd: secret, options: options}

	secret.AddCommand(rootData.setCmd())
	secret.AddCommand(rootData.getCmd())
	secret.AddCommand(rootData.serverCmd())

	return rootData
}

func (r *root) setCmd() *cobra.Command {
	var key string
	var cipherKey string
	var value string
	var path string

	var setCmd = &cobra.Command{
		Use:   "set",
		Short: "Saves data to the specified path in encrypted form",
		Long:  "it takes keys and a value and path from user and saves value in encrypted manner in specified storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			var cr = crypto.NewCryptographer([]byte(cipherKey))
			ds, err := storage.NewFileVault(path)
			if err != nil {
				return fmt.Errorf("can't create storage by path: %w", err)
			}
			pr := provider.NewProvider(cr, ds)
			err = pr.SetData([]byte(key), []byte(value))
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

	return setCmd
}

func (r *root) getCmd() *cobra.Command {
	var key string
	var cipherKey string
	var path string
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get data from specified path in decrypted form",
		Long:  "it takes keys and path from user and get value in decrypted manner from specified storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			var cr = crypto.NewCryptographer([]byte(cipherKey))
			ds, err := storage.NewFileVault(path)
			if err != nil {
				return fmt.Errorf("can't get storage by path: %w", err)
			}
			pr := provider.NewProvider(cr, ds)
			data, err := pr.GetData([]byte(key))
			if err != nil {
				return fmt.Errorf("can't get data by key: %w", err)
			}
			fmt.Println(string(data))
			return nil
		},
	}
	getCmd.Flags().StringVarP(&key, "key", "k", key, "key for pair key-value")
	getCmd.Flags().StringVarP(&cipherKey, "cipher-key", "c", cipherKey, "cipher key for data encryption and decryption")
	getCmd.Flags().StringVarP(&path, "path", "p", "file.txt", "the place where the value will be got")

	return getCmd
}

func (r *root) serverCmd() *cobra.Command {
	var path string
	var port string
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run server runner mode to start the app as a daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			ds, err := storage.NewFileVault(path)
			if err != nil {
				return fmt.Errorf("can't get storage by path: %w", err)
			}
			store := make(map[string]api.MethodFactoryFunc)
			store["local"] = func(cipher string) (secretApi.Provider, func()) {
				cr := crypto.NewCryptographer([]byte(cipher))
				return provider.NewProvider(cr, ds), nil
			}

			handler := api.NewMethods(store)
			router := chi.NewRouter()

			srv := &http.Server{Addr: ":" + port, Handler: router}
			//router.With(middleware.Heartbeat("/ping")).Route("/ping", func(r chi.Router) {
			//	router.Get("/", handler.GetByKey)
			//	router.Post("/", handler.SetByKey)
			//})
			router.Use(middleware.Heartbeat("/ping"))
			router.Get("/", handler.GetByKey)
			router.Post("/", handler.SetByKey)
			done := make(chan os.Signal, 1)
			shutdownCh := make(chan struct{})
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				err := srv.ListenAndServe()
				if err != nil {
					fmt.Println("connection error: %w", err)
				}
			}()
			fmt.Println("Server started")

			select {
			case <-done:
				fmt.Println("Server stopped")
			case <-cmd.Context().Done():
				fmt.Println("Server stopped with context")
			}

			go func(ctx context.Context) {
				defer close(shutdownCh)
				ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				defer cancel()
				err = srv.Shutdown(ctx)
				if err != nil {
					fmt.Println("Server Shutdown Failed", err)
				}
			}(context.Background())
			<-shutdownCh
			fmt.Println("Server exit")

			return nil
		},
	}
	serverCmd.Flags().StringVarP(&path, "path", "p", "file.txt", "the place where the key/value will be stored/got")
	serverCmd.Flags().StringVarP(&port, "port", "t", "8888", "localhost address")

	return serverCmd
}
