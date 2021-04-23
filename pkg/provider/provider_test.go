package provider

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProvider_SetData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		key := []byte{1, 1, 1}
		value := []byte{0, 1, 3}
		encodedValue := []byte{0, 1, 3}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockCr.On("Encode", value).Return(encodedValue, nil)
		mockDs.On("SaveData", key, encodedValue).Return(nil)

		p := NewProvider(mockCr, mockDs)
		err := p.SetData(key, value)
		require.NoError(t, err)

		mockCr.AssertExpectations(t)
		mockDs.AssertExpectations(t)
	})

	t.Run("encode method error", func(t *testing.T) {
		key := []byte{1, 1, 1}
		value := []byte{0, 1, 3}
		encodedValue := []byte{0, 1, 3}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockCr.On("Encode", value).Return(encodedValue, fmt.Errorf("test"))

		p := NewProvider(mockCr, mockDs)
		err := p.SetData(key, value)
		require.Error(t, err, "error assert failed")
		require.EqualValues(t, "provider, SetData method: encode error: test", err.Error())

		mockCr.AssertExpectations(t)
	})

	t.Run("setData method error", func(t *testing.T) {
		key := []byte{1, 1, 1}
		value := []byte{0, 1, 3}
		encodedValue := []byte{0, 1, 3}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockCr.On("Encode", value).Return(encodedValue, nil)
		mockDs.On("SaveData", key, encodedValue).Return(fmt.Errorf("test"))

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
		encodedValue := []byte{0, 1, 3}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockDs.On("ReadData", key).Return(encodedValue, nil)
		mockCr.On("Decode", encodedValue).Return(value, nil)

		p := NewProvider(mockCr, mockDs)
		data, err := p.GetData(key)
		require.NoError(t, err)

		mockDs.AssertExpectations(t)
		mockCr.AssertExpectations(t)
		assert.Equal(t, value, data)
	})

	t.Run("read data error", func(t *testing.T) {
		key := []byte{1, 1, 1}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockDs.On("ReadData", key).Return(nil, fmt.Errorf("test"))

		p := NewProvider(mockCr, mockDs)
		data, err := p.GetData(key)
		require.Error(t, err, "error assert failed")
		require.EqualValues(t, "provider, GetData method: read data error: test", err.Error())
		assert.Equal(t, "", string(data))
		mockDs.AssertExpectations(t)
	})

	t.Run("decode error", func(t *testing.T) {
		key := []byte{1, 1, 1}
		value := []byte{0, 1, 3}
		encodedValue := []byte{0, 1, 3}
		mockCr := new(MockCryptographer)
		mockDs := new(MockDataSaver)

		mockDs.On("ReadData", key).Return(encodedValue, nil)
		mockCr.On("Decode", encodedValue).Return(value, fmt.Errorf("test"))

		p := NewProvider(mockCr, mockDs)
		data, err := p.GetData(key)
		require.Error(t, err, "error assert failed")
		require.EqualValues(t, "provider, GetData method: decode error: test", err.Error())
		assert.Equal(t, "", string(data))
		mockCr.AssertExpectations(t)
	})
}
