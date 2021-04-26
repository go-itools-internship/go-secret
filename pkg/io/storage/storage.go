/*
Package storage provides functions for storing and retrieving data.
*/
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type fileVault struct {
	storage map[string][]byte
	path    string
}

func NewFileVault(path string) (*fileVault, error) {
	storage := make(map[string][]byte)

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(path)
			if err != nil {
				return nil, fmt.Errorf("unable to create file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unable to open file: %w", err)
		}
	} else {
		err = json.NewDecoder(file).Decode(&storage)
		if err != nil {
			return nil, fmt.Errorf("unable to write data from file to map: %w", err)
		}
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	return &fileVault{storage: storage, path: path}, err
}

func (f *fileVault) SaveData(key, encodedValue []byte) error {

	file, err := os.OpenFile(f.path, os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	f.storage[string(key)] = encodedValue

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
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	err = json.NewDecoder(file).Decode(&f.storage)
	if err != nil {
		return nil, fmt.Errorf("unable to write data from file to map: %w", err)
	}

	data := f.storage[string(key)]

	return data, err
}
