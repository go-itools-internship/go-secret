// Package secret contains interfaces for API secrets, storage and encoder.
package secret

// Provider organizes a gateway for managing data by setting/getting by a key.
type Provider interface {
	// SetData set by key and put value and return error
	SetData(key, value []byte) error
	// GetData get data by key and return array of bytes and error
	GetData(key []byte) ([]byte, error)
}

// Cryptographer describes the behavior for encrypting and decrypting data
type Cryptographer interface {
	// Encode takes value and return encode value and error
	Encode(value []byte) ([]byte, error)
	// Decode takes encode value and return decode value and error
	Decode(encodedValue []byte) ([]byte, error)
}

// DataSaver describes the behavior of storing and reading data in the storage
// For implementation, we can use any type of storage (for example: cloud, file, local memory)
type DataSaver interface {
	// SaveData save encoded value by key and
	// It return any write error encountered.
	SaveData(key, encodedValue []byte) error
	// ReadData get encoded data by key
	// It returns the array of bytes and any write error encountered.
	ReadData(key []byte) ([]byte, error)
}
