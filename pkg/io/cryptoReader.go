/*
	Package io contain custom implementation of io Reader.
	Using for read data many times without changes.
*/
package io

// A RootReader implements the io.Reader
type rootReader struct {
	data []byte // data to read
}

// NewRootReader returns a new Reader reading from p.
func NewRootReader(p []byte) *rootReader {
	return &rootReader{data: p}
}

// Read reads up to len(p) bytes into p. It returns the number of bytes
// Not return EOF error for use with io.ReadFull method
func (r *rootReader) Read(p []byte) (n int, err error) {
	n = copy(p, r.data)
	if len(r.data) < len(p) { // if data to read less than len of buffer
		return len(p), nil
	}
	return n, nil
}
