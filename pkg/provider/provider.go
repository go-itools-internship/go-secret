// Package provider implements crypto and storage
package provider

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/go-itools-internship/go-secret/pkg/secret"
)

type provider struct {
	cryptographer secret.Cryptographer
	dataSaver     secret.DataSaver
	logger        *zap.Logger
}

func NewProvider(cryptographer secret.Cryptographer, dataSaver secret.DataSaver, logger *zap.SugaredLogger) *provider {
	return &provider{cryptographer: cryptographer, dataSaver: dataSaver}
}

func (p *provider) SetData(key, value []byte) error {
	encodedValue, err := p.cryptographer.Encode(value)
	if err != nil {
		return fmt.Errorf("provider, SetData method: encode error: %w", err)
	}
	saveError := p.dataSaver.SaveData(key, encodedValue)
	if saveError != nil {
		return fmt.Errorf("provider, SetData method: save error: %w", saveError)
	}
	return nil
}

func (p *provider) GetData(key []byte) ([]byte, error) {
	data, err := p.dataSaver.ReadData(key)
	if err != nil {
		return nil, fmt.Errorf("provider, GetData method: read data error: %w", err)
	}
	decode, err := p.cryptographer.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("provider, GetData method: decode error: %w", err)
	}
	return decode, nil
}
