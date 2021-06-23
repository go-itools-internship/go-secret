/*
Package storage provides functions for storing and retrieving data.
*/
package storage

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type fileVault struct {
	storage map[string][]byte
	path    string
}

func NewFileVault(path string) (*fileVault, error) {
	storage := make(map[string][]byte)

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("filevault: unable to open file: %w", err)
		}
		if f, err = os.Create(path); err != nil {
			return nil, fmt.Errorf("filevault: unable to create file: %w", err)
		}
		if _, err := fmt.Fprintf(f, `{}`); err != nil {
			return nil, fmt.Errorf("filevault: unable to init file: %w", err)
		}
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = fmt.Errorf("filevault: %w", cerr)
		}
	}()

	if err := json.NewDecoder(f).Decode(&storage); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("filevault: unable to decode data: %w", err)
		}
	}

	return &fileVault{storage: storage, path: path}, nil
}

func (f *fileVault) SaveData(key, encodedValue []byte) error {
	file, err := os.OpenFile(f.path, os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("filevault: unable to open file while saving: %w", err)
	}
	defer func() {
		if cerr := file.Close(); err == nil {
			err = fmt.Errorf("filevault: %w", cerr)
		}
	}()

	f.storage[hex.EncodeToString(key)] = encodedValue
	if err = json.NewEncoder(file).Encode(f.storage); err != nil {
		return fmt.Errorf("filevault: unable to encode data while saving: %w", err)
	}
	return nil
}

func (f *fileVault) ReadData(key []byte) ([]byte, error) {
	file, err := os.OpenFile(f.path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("filevault: unable to open file while reading: %w", err)
	}
	defer func() {
		if cerr := file.Close(); err == nil {
			err = fmt.Errorf("filevault: %w", cerr)
		}
	}()

	if err := json.NewDecoder(file).Decode(&f.storage); err != nil {
		return nil, fmt.Errorf("filevault: unable to decode while reading: %w", err)
	}

	data, ok := f.storage[hex.EncodeToString(key)]
	if !ok {
		return nil, fmt.Errorf("filevault: cannot read data: not found")
	}

	return data, nil
}
