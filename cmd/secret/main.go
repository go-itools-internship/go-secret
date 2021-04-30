package main

import (
	"context"
	"fmt"

	"github.com/go-itools-internship/go-secret/cmd/secret/cmd"
)

var (
	revision = "unknown"
)

func main() {
	p := cmd.New(cmd.Version(revision))
	err := p.Execute(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
