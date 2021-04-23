/*
Package storage provides functions for storing and retrieving data.
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
	file    *os.File
}

func NewFileVault(path string) (*fileVault, error) {
	var storage = make(map[string][]byte)

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("unable to open or create file: %w", err)
		}
	} else {
		// decode data from json format from file and put it in storage
		err = json.NewDecoder(file).Decode(&storage)
		if err != nil {
			return nil, fmt.Errorf("unable to write data from file to map: %w", err)
		}
	}
	return &fileVault{storage: storage, path: path, file: file}, err
}

func (f *fileVault) SaveData(key, encodedValue []byte) error {

	// open the file for writing only
	file, err := os.OpenFile(f.path, os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}

	// set coming encodedValue to our storage
	f.storage[string(key)] = encodedValue

	// encode the data in storage and put it to file in json format
	err = json.NewEncoder(file).Encode(f.storage)
	if err != nil {
		return fmt.Errorf("unable to write data from map to file: %w", err)
	}
	return err
}

func (f *fileVault) ReadData(key []byte) ([]byte, error) {
	file, err := os.OpenFile(f.path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}

	// decode data from json format from file and put it in storage
	err = json.NewDecoder(file).Decode(&f.storage)
	if err != nil {
		return nil, fmt.Errorf("unable to write data from file to map: %w", err)
	}

	data := f.storage[string(key)]

	return data, err
}

func (f *fileVault) Close() (err error) {

	//encode the data in storage and put it to file in json format
	//err = json.NewEncoder(f.file).Encode(f.storage)
	//if err != nil {
	//	return fmt.Errorf("unable to write data from map to file: %w", err)
	//}

	defer func() {
		cerr := f.file.Close()
		if err == nil {
			err = cerr
		}
	}()
	return err
}
