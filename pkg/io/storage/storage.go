package pkg

import (
	"fmt"
	"io"
	"os"
)

type FileVault struct {
	Storage map[string][]byte
	Path    string
}

func NewFileVault(storage map[string][]byte, path string) *FileVault {
	return &FileVault{Storage: storage, Path: path}
}

func (f *FileVault) SaveData(key, encodedValue []byte) error {
	file, err := os.Create(f.Path + string(key))
	if err != nil {
		fmt.Println("Unable to Create file:", err)
		os.Exit(1)
	}
	f.Storage[string(key)] = encodedValue

	defer file.Close()
	_, err = file.Write(encodedValue)
	if err != nil {
		fmt.Println("Unable to write encodedValue:", err)
		os.Exit(1)
	}

	fmt.Println("Done.")
	return err
}

func (f *FileVault) ReadData(key []byte) ([]byte, error) {
	file, err := os.Open(f.Path + string(key))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data := make([]byte, 64)

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		fmt.Print(string(data[:n]))
	}
	fmt.Println()
	return data, err
}
