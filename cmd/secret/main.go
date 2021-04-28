package main

import (
	"context"
	"fmt"

	"github.com/go-itools-internship/go-secret/cmd/secret/cmd"
)

func main() {
	p := cmd.NewRoot()
	err := p.RootExecute(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
