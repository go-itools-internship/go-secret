/*
Package  storage provides functions for storing and retrieving data.
from json file
*/
package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type fileVault struct {
	storage map[string][]byte
	path    string
}

func NewFileVault(path string) (*fileVault, error) {
	var storage = make(map[string][]byte)

	// open the file at the specified path, create it if the file is not found
	file, err := os.OpenFile(path, os.O_CREATE, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to open or create file: %w", err)
	}
	defer file.Close()

	return &fileVault{storage: storage, path: path}, nil
}

func (f *fileVault) SaveData(key, encodedValue []byte) error {

	// open the file for writing only
	file, err := os.OpenFile(f.path, os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	// set coming encodedValue to our storage
	f.storage[string(key)] = encodedValue

	// encode the data in storage and put it to file in json format
	err = json.NewEncoder(file).Encode(f.storage)
	if err != nil {
		return fmt.Errorf("unable to write data from map to file: %w", err)
	}
	return nil
}

func (f *fileVault) ReadData(key []byte) ([]byte, error) {
	file, err := os.OpenFile(f.path, os.O_RDONLY, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	// decode data from json format from file and put it in storage
	err = json.NewDecoder(file).Decode(&f.storage)
	if err != nil {
		return nil, fmt.Errorf("unable to write data from file to map: %w", err)
	}

	data := f.storage[string(key)]

	return data, nil
}
