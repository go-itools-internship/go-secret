package main

import (
	"fmt"
	pkg "github.com/go-itools-internship/go-secret/pkg/io/storage"
)

var (
	revision = "unknown"
)

func main() {
	fmt.Printf("secret %s\n", revision)
	fmt.Println("Hi from go-secret!")

	storage := make(map[string][]byte)
	fv := pkg.NewFileVault(storage, "internal/fileStorage/")
	fv.SaveData([]byte("f1"), []byte("123"))
	fv.SaveData([]byte("f2"), []byte("12345"))
	fv.SaveData([]byte("f3"), []byte("123"))

	fv.ReadData([]byte("f1"))
	fv.ReadData([]byte("f2"))
	fv.ReadData([]byte("f3"))
}
