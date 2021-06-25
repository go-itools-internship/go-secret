/*
	Package io contain custom implementation of io Reader.
	Using for read data many times without changes.
*/
package crypto

// A lootReader implements the io.Reader.
type loopReader struct {
	data []byte // data to read
}

// LoopReader returns a new Reader reading from p.
// 	Regardless of the number of iterations, returns the same data in buffer size.
func LoopReader(p []byte) *loopReader {
	return &loopReader{data: p}
}

// Read reads up to len(p) bytes into p. It returns the number of bytes.
// Not return EOF error for use with io.ReadFull method.
func (r *loopReader) Read(p []byte) (n int, err error) {
	n = copy(p, r.data)
	if len(r.data) < len(p) { // if data to read less than len of buffer
		return len(p), nil
	}
	return n, nil
}
