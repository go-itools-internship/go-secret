package provider

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProvider_SetData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		key := []byte{1, 1, 1}
		value := []byte{0, 1, 3}
		encodedValue := []byte{0, 1, 3, 5, 34}
		encodedKey := []byte{0, 1}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockCr.On("Encode", value).Return(encodedValue, nil)
		mockCr.On("Encode", key).Return(encodedKey, nil)
		mockDs.On("SaveData", encodedKey, encodedValue).Return(nil)

		p := NewProvider(mockCr, mockDs)
		err := p.SetData(key, value)
		require.NoError(t, err)

		mockCr.AssertExpectations(t)
		mockDs.AssertExpectations(t)
	})

	t.Run("encode method error", func(t *testing.T) {
		key := []byte{1, 1, 1}
		value := []byte{0, 1, 3}
		encodedValue := []byte{0, 1, 3, 5, 34}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockCr.On("Encode", value).Return(encodedValue, fmt.Errorf("test"))

		p := NewProvider(mockCr, mockDs)
		err := p.SetData(key, value)
		require.Error(t, err, "error assert failed")
		require.EqualValues(t, "provider, SetData method: encode value error: test", err.Error())

		mockCr.AssertExpectations(t)
	})

	t.Run("setData method error", func(t *testing.T) {
		key := []byte{1, 1, 1}
		value := []byte{0, 1, 3}
		encodedValue := []byte{0, 1, 3, 5, 34}
		encodedKey := []byte{0, 1, 3, 5}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockCr.On("Encode", value).Return(encodedValue, nil)
		mockCr.On("Encode", key).Return(encodedKey, nil)
		mockDs.On("SaveData", encodedKey, encodedValue).Return(fmt.Errorf("test"))

		p := NewProvider(mockCr, mockDs)
		err := p.SetData(key, value)
		require.Error(t, err, "error assert failed")
		require.EqualValues(t, "provider, SetData method: save error: test", err.Error())

		mockDs.AssertExpectations(t)
	})
}

func TestProvider_GetData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		key := []byte{1, 1, 1}
		value := []byte{0, 1, 3}
		encodedValue := []byte{0, 1, 3, 5, 34}
		encodedKey := []byte{1, 1, 1}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockDs.On("ReadData", key).Return(encodedValue, nil)
		mockCr.On("Encode", key).Return(encodedKey, nil)
		mockCr.On("Decode", encodedValue).Return(value, nil)

		p := NewProvider(mockCr, mockDs)
		data, err := p.GetData(encodedKey)
		require.NoError(t, err)
		assert.Equal(t, value, data)

		mockDs.AssertExpectations(t)
		mockCr.AssertExpectations(t)
	})

	t.Run("read data error", func(t *testing.T) {
		key := []byte{1, 1, 1}
		encodedKey := []byte{1, 1, 9}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockCr.On("Encode", encodedKey).Return(key, nil)
		mockDs.On("ReadData", key).Return(nil, fmt.Errorf("test"))

		p := NewProvider(mockCr, mockDs)
		data, err := p.GetData(encodedKey)
		require.Error(t, err, "error assert failed")
		require.EqualValues(t, "provider, GetData method: read data error: test", err.Error())
		assert.Equal(t, "", string(data))
		mockDs.AssertExpectations(t)
	})

	t.Run("decode error", func(t *testing.T) {
		key := []byte{1, 1, 1}
		value := []byte{0, 1, 3}
		encodedValue := []byte{0, 1, 3, 4}
		encodedKey := []byte{1, 1, 1}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockCr.On("Encode", key).Return(encodedKey, nil)
		mockDs.On("ReadData", key).Return(encodedValue, nil)
		mockCr.On("Decode", encodedValue).Return(value, fmt.Errorf("test"))

		p := NewProvider(mockCr, mockDs)
		data, err := p.GetData(encodedKey)
		require.Error(t, err, "error assert failed")
		require.EqualValues(t, "provider, GetData method: decode error: test", err.Error())
		assert.Equal(t, "", string(data))
		mockCr.AssertExpectations(t)
	})
}
