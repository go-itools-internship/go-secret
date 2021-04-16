// Package contains interfaces for API secrets, storage and encoder.
package secret

//Provider defines what we can do with the data in our application
//Describes the behavior of commands from the terminal
type Provider interface {
	// set by key and put value
	SetData(key, value []byte) error
	//get data by key
	GetData(key []byte) ([]byte, error)
}

// The cryptographer describes the behavior for encrypting and decrypting data
type Cryptographer interface {
	//takes value and return encode value
	Encode(value []byte) ([]byte, error)
	//takes encode value and return decode value
	Decode(encodedValue []byte) ([]byte, error)
}

// The DataSaver describes the behavior of storing and reading data in the storage
// For implementation, we can use any type of storage (for example: cloud, file, local memory)
type DataSaver interface {
	//save encoded value by key
	SaveData(key, encodedValue []byte) error
	//get encoded data by key
	ReadData(key []byte) ([]byte, error)
}
