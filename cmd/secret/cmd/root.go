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
	rootCmd   *cobra.Command
}

func (r *root) RootExecute(ctx context.Context) error { //todo take context
	err := r.rootCmd.ExecuteContext(ctx)
	return err
}

func NewRoot() *root {
	var rootCmd = &cobra.Command{
		Use:   "root",
		Short: "root includes cli commands",
	}
	v := rootCmd.PersistentFlags().StringP("value", "v", "", "value to provider")
	k := rootCmd.PersistentFlags().StringP("key", "k", "", "key to provider")
	ck := rootCmd.PersistentFlags().StringP("cipherKey", "c", "", "cipher key to provider")
	p := rootCmd.PersistentFlags().StringP("path", "p", "file.txt", "path to provider")
	rootData := &root{cipherKey: ck, key: k, value: v, path: p, rootCmd: rootCmd}
	rootCmd.AddCommand(rootData.getCmd())
	rootCmd.AddCommand(rootData.setCmd())
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
			fmt.Println(string(data))
		},
	}
	return getCmd
}
