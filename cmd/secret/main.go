package main

import (
	"fmt"
	"secret"
)

var (
	revision = "unknown"
)

func main() {
	fmt.Printf("secret %s\n", revision)
	fmt.Println("Hi from go-secret!")

	v := secret.FileVault("encoding-key", "path/to/file")
	err := v.Set("key-name", "key-value")
	value, err := v.Get("key-name")
	fmt.Println(value) // "key-value"
}

//Pick type of storage?
//1-File
//2-Memory
//3-Cloud
//Pick set or get or back to main menu
//1-set -> encode
//2-get -> decode
//3-back

	//if 1-set
	//input key and data
		// -> print set OK

	//if 2-get
	//input key to get data
		//print data

	//if 3-back
		//backToMainMenu() auto-method
// exit or back main menu

//Provider (? -> set/get) ->  Crypter(encode) -> Storage(Cloud/File/Memory)
//						  <-		 (decode) <-

//Сryptographer - internal
//1-encode
//2-decode
