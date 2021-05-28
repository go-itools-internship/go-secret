// Package client work with server api and set/get keys to it
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"

	api "github.com/go-itools-internship/go-secret/pkg/http"
)

type client struct {
	options options
	url     string // address where client will be work with server
	logger  *zap.SugaredLogger
}

// options client
type options struct {
	c *http.Client
}

var defaultOptions = options{
	c: &http.Client{Timeout: 20 * time.Second},
}

type Option func(o *options)

// Client provide any http client
func Client(c *http.Client) Option {
	return func(options *options) {
		options.c = c
	}
}

// New function initializes a structure that provides client accessing functions.
//
// Accepts url where client will be work with server and client options.
func New(url string, logger *zap.SugaredLogger, opts ...Option) *client {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	newClient := &client{options: options, url: url, logger: logger}
	return newClient
}

// GetByKey get data from server by key method and cipher key.
// 	Key for pair key-value.
// 	Cipher key for data encryption and decryption.
// 	Method to choose different providers.
func (c *client) GetByKey(ctx context.Context, key, method, cipherKey string) (string, error) {
	logger := c.logger.Named("get-by-key")
	var responseBody struct {
		Value string `json:"value"`
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return "", fmt.Errorf("secret client: can't create request %w", err)
	}
	req.Header.Set(api.ParamCipherKey, cipherKey)
	query := req.URL.Query()
	query.Set(api.ParamGetterKey, key)
	query.Set(api.ParamMethodKey, method)
	req.URL.RawQuery = query.Encode()

	resp, err := c.options.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("secret client: can't do response %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Warnf("secret client: cannot close request body: %s", err.Error())
		}
	}()
	if resp.StatusCode != http.StatusOK {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("secret client: can't get response body %w", err)
		}
		return "", fmt.Errorf("secret client: can't get data: body: %q, status code: %d", responseBody, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", fmt.Errorf("secret client: cannot decode body: %w", err)
	}
	return responseBody.Value, nil
}

// SetByKey set data to server by getterKey, value, method, cipherKey.
// 	GetterKey for pair key-value.
// 	Cipher key to set data encryption.
// 	Method to choose different providers.
func (c *client) SetByKey(ctx context.Context, getterKey, value, method, cipherKey string) error {
	logger := c.logger.Named("set-by-key")
	postBody, err := json.Marshal(map[string]string{
		"getter": getterKey,
		"method": method,
		"value":  value,
	})
	if err != nil {
		return fmt.Errorf("secret client: can't marshal body %w", err)
	}
	body := bytes.NewBuffer(postBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, body)
	if err != nil {
		return fmt.Errorf("secret client: can't create request %w", err)
	}
	req.Header.Set(api.ParamCipherKey, cipherKey)

	resp, err := c.options.c.Do(req)
	if err != nil {
		return fmt.Errorf("secret client: can't do request %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Warnf("secret client: cannot close request body: %s", err.Error())
		}
	}()
	if resp.StatusCode != http.StatusNoContent {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("secret client: can't get response body %w", err)
		}
		return fmt.Errorf("secret client: can't set data: body: %q, status code: %d", responseBody, resp.StatusCode)
	}
	return nil
}
