package main

import (
	"fmt"
	"secret"
)

var (
	revision = "unknown"
)

type Storage struct {
	fileStorage map[string]string
	path        string
}

type File struct {
	Storage
}

func (f File) Set(key, value string) (err error) {
	return err
}

func (f File) Get(key string) (value string, err error) {
	return "", err
}

type Memory struct {
	Storage
}

func (m Memory) Set(key, value string) (err error) {
	return err
}

func (m Memory) Get(key string) (value string, err error) {
	return "", err
}

type Cloud struct {
	Storage
}

func (c Cloud) Set(key, value string) (err error) {
	return err
}

func (c Cloud) Get(key string) (value string, err error) {
	return "", err
}

type Provider interface { // for Api
	Set()
	Get()
}

type Storager interface { // for storage
	pickTypeOfStorage()
}

func pickTypeOfStorage(pickType int) func() {
	switch pickType {
	case 1:
		return FileVault()
	case 2:
		return MemoryVault()
	case 3:
		return CloudVault()
	default:
		fmt.Println("Invalid value")
		return pickTypeOfStorage(pickType)
	}
}

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
