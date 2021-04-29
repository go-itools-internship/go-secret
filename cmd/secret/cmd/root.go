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
	key       *string
	cipherKey *string
	value     *string
	path      *string
	cmd       *cobra.Command
}

func (r *root) Execute(ctx context.Context) error {
	return r.cmd.ExecuteContext(ctx)
}

func New(version string) *root {
	var secret = &cobra.Command{
		Use:     "secret",
		Short:   "Root includes cli commands",
		Long:    "Create CLI to set and get secrets via the command line",
		Version: version,
	}
	v := secret.PersistentFlags().StringP("value", "v", "", "value to be encrypted")
	k := secret.PersistentFlags().StringP("key", "k", "", "key to encrypt value")
	ck := secret.PersistentFlags().StringP("cipherKey", "c", "", "cipher key to cryptographer")
	p := secret.PersistentFlags().StringP("path", "p", "file.txt", "path to file")
	rootData := &root{cipherKey: ck, key: k, value: v, path: p, cmd: secret}
	secret.AddCommand(rootData.getCmd())
	secret.AddCommand(rootData.setCmd())
	return rootData
}

func (r *root) setCmd() *cobra.Command {
	var setCmd = &cobra.Command{
		Use:   "set",
		Short: "Set data in file by key",
		Run: func(cmd *cobra.Command, args []string) {
			var cr = crypto.NewCryptographer([]byte(*r.cipherKey))
			ds, err := storage.NewFileVault(*r.path)
			if err != nil {
				fmt.Println(err)
			}
			pr := provider.NewProvider(cr, ds)
			err = pr.SetData([]byte(*r.key), []byte(*r.value))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("data set in file success")
		},
	}
	return setCmd
}

func (r *root) getCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get data by key",
		Run: func(cmd *cobra.Command, args []string) {
			var cr = crypto.NewCryptographer([]byte(*r.cipherKey))
			ds, err := storage.NewFileVault(*r.path)
			if err != nil {
				fmt.Println(err)
			}
			pr := provider.NewProvider(cr, ds)
			data, err := pr.GetData([]byte(*r.key))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(data))
		},
	}
	return getCmd
}
