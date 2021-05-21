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

	api "github.com/go-itools-internship/go-secret/pkg/http"
)

type client struct {
	client http.Client // TODO client options
	url    string      // address where client will be work with server
}

// New function creates client with timeout
func New(url string) *client {
	client1 := http.Client{Timeout: 20 * time.Second}
	newClient := &client{client: client1, url: url}
	return newClient
}

// GetByKey get data from server by key method and cipher key
func (c *client) GetByKey(ctx context.Context, key, method, cipherKey string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return "", fmt.Errorf("secret client: can't create request %w", err)
	}
	req.Header.Set(api.ParamCipherKey, cipherKey)

	query := req.URL.Query()
	query.Set(api.ParamGetterKey, key)
	query.Set(api.ParamMethodKey, method)
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("secret client: can't do response %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("cannot close request body: ", err.Error())
		}
	}()

	var responseBody struct {
		Value string `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", fmt.Errorf("cannot decode body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("secret client: can't get data: body: %q, status code: %d", responseBody.Value, resp.StatusCode)
	}
	return responseBody.Value, nil
}

// SetByKey set data to server by  getterKey, value, method, cipherKey
func (c *client) SetByKey(ctx context.Context, getterKey, value, method, cipherKey string) error {
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

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("secret client: can't do request %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("cannot close request body: ", err.Error())
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
