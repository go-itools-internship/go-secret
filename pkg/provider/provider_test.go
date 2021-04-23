package provider

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProvider_SetData(t *testing.T) {
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
}

func TestProvider_GetData(t *testing.T) {
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
}
