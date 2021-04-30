// Package cmd provides functions set and get with cobra library
package cmd

import (
	"context"
	"fmt"

	"github.com/go-itools-internship/go-secret/pkg/crypto"
	"github.com/go-itools-internship/go-secret/pkg/io/storage"
	"github.com/go-itools-internship/go-secret/pkg/provider"
	"github.com/spf13/cobra"
)

type root struct {
	options   options
	key       *string
	cipherKey *string
	value     *string
	path      *string
	cmd       *cobra.Command
}

type options struct {
	version string
}

var defaultOption = options{
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

// New function create and set flags and commands to the cobra CLI
// RootOptions adds additional features to the cobra CLI
func New(opts ...RootOptions) *root {
	options := defaultOption
	for _, opt := range opts {
		opt(&options)
	}

	var secret = &cobra.Command{
		Use:     "secret",
		Short:   "Contains commands to set and get encrypt data to storage",
		Long:    "Create CLI to set and get secrets via the command line",
		Version: options.version,
	}
	v := secret.PersistentFlags().StringP("value", "v", "", "value to be encrypted")
	k := secret.PersistentFlags().StringP("key", "k", "", "key for pair key-value")
	ck := secret.PersistentFlags().StringP("cipher-key", "c", "", "cipher key for data encryption and decryption")
	p := secret.PersistentFlags().StringP("path", "p", "file.txt", "the place where the key/value will be stored/got")
	rootData := &root{cipherKey: ck, key: k, value: v, path: p, cmd: secret, options: options}

	secret.AddCommand(rootData.getCmd())
	secret.AddCommand(rootData.setCmd())

	return rootData
}

func (r *root) setCmd() *cobra.Command {
	var setCmd = &cobra.Command{
		Use:   "set",
		Short: "Saves data to the specified path in encrypted form",
		Long:  "it takes keys and a value and path from user and saves value in encrypted manner in specified storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			var cr = crypto.NewCryptographer([]byte(*r.cipherKey))
			ds, err := storage.NewFileVault(*r.path)
			if err != nil {
				return fmt.Errorf("can't create storage by path: %w", err)
			}
			pr := provider.NewProvider(cr, ds)
			err = pr.SetData([]byte(*r.key), []byte(*r.value))
			if err != nil {
				return fmt.Errorf("can't set data %w", err)
			}
			return nil
		},
	}
	return setCmd
}

func (r *root) getCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get data from specified path in decrypted form",
		Long:  "it takes keys and path from user and get value in decrypted manner from specified storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			var cr = crypto.NewCryptographer([]byte(*r.cipherKey))
			ds, err := storage.NewFileVault(*r.path)
			if err != nil {
				return fmt.Errorf("can't get storage by path: %w", err)
			}
			pr := provider.NewProvider(cr, ds)
			data, err := pr.GetData([]byte(*r.key))
			if err != nil {
				return fmt.Errorf("can't get data by key: %w", err)
			}
			fmt.Println(string(data))
			return nil
		},
	}
	return getCmd
}
