package crypto

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootReader_Read(t *testing.T) {
	t.Run("success if buffer more len of buffer", func(t *testing.T) {
		reader := LoopReader([]byte{1, 3, 4, 5})
		p := make([]byte, 8)
		n, err := io.ReadFull(reader, p)
		require.NoError(t, err)
		require.EqualValues(t, []byte{1, 3, 4, 5, 0, 0, 0, 0}, p)
		require.EqualValues(t, len(p), n)
	})
	t.Run("no error if reader read more data than len of buffer", func(t *testing.T) {
		reader := LoopReader([]byte{1, 3, 4, 5, 6})
		p := make([]byte, 4)
		n, err := io.ReadFull(reader, p)
		require.NoError(t, err)
		require.EqualValues(t, []byte{1, 3, 4, 5}, p)
		require.EqualValues(t, len(p), n)
	})
	t.Run("no error reader nil", func(t *testing.T) {
		reader := LoopReader(nil)
		p := make([]byte, 4)
		n, err := io.ReadFull(reader, p)
		require.NoError(t, err)
		require.EqualValues(t, []byte{0, 0, 0, 0}, p)
		require.EqualValues(t, len(p), n)
	})
	t.Run("when reader nil and buffer empty", func(t *testing.T) {
		reader := LoopReader(nil)
		p := make([]byte, 0)
		n, err := io.ReadFull(reader, p)
		require.NoError(t, err)
		require.EqualValues(t, []byte{}, p)
		require.EqualValues(t, len(p), n)
	})
	t.Run("no error if buffer empty", func(t *testing.T) {
		reader := LoopReader([]byte{1, 3, 4})
		p := make([]byte, 0)
		n, err := io.ReadFull(reader, p)
		require.NoError(t, err)
		require.EqualValues(t, []byte{}, p)
		require.EqualValues(t, len(p), n)
	})
	t.Run("return same result", func(t *testing.T) {
		reader := LoopReader([]byte{1, 3, 4, 5, 6})
		i := 0
		for i < 100 {
			p := make([]byte, 8)
			n, err := io.ReadFull(reader, p)
			require.NoError(t, err)
			require.EqualValues(t, []byte{1, 3, 4, 5, 6, 0, 0, 0}, p)
			require.EqualValues(t, len(p), n)
			i++
		}
	})
	t.Run("success if read slice of slice buffers", func(t *testing.T) {
		reader := LoopReader([]byte{1, 3, 4, 5, 6})
		m := [][]byte{make([]byte, 3), make([]byte, 6), make([]byte, 10), make([]byte, 1)}
		expected := [][]byte{{1, 3, 4}, {1, 3, 4, 5, 6, 0}, {1, 3, 4, 5, 6, 0, 0, 0, 0, 0}, {1}}
		for i := range m {
			n, err := io.ReadFull(reader, m[i])
			require.NoError(t, err)
			require.EqualValues(t, expected[i], m[i])
			require.EqualValues(t, len(m[i]), n)
		}
	})
}
