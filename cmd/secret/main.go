package main

import (
	"fmt"
)

var (
	revision = "unknown"
)

func main() {
	fmt.Printf("secret %s\n", revision)
	fmt.Println("Hi from go-secret!")
}
