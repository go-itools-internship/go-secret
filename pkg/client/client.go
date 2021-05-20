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
	client http.Client
	url    string // address where client will be work with server
}

// New function creates client with timeout
func New(url string) *client {
	client1 := http.Client{Timeout: time.Second}
	newClient := &client{client: client1, url: url}
	return newClient
}

func (c *client) GetByKey(ctx context.Context, key string, cipherKey string, method string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return "", fmt.Errorf("secret client: can't create request %w", err)
	}
	req.RequestURI = ""
	req.Header.Set(api.ParamCipherKey, cipherKey)

	query := req.URL.Query()
	query.Set(api.ParamGetterKey, key)
	query.Set(api.ParamMethodKey, method)
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't read responce data %w", err)
	}
	//
	var responseBody struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(data, &responseBody); err != nil {
		return "", fmt.Errorf("cannot write response: %w", err)
	}
	return responseBody.Value, nil
}

func (c *client) SetByKey(ctx context.Context, getterKey string, value string, method string, cipherKey string) error {
	postBody, _ := json.Marshal(map[string]string{
		"getter": getterKey,
		"method": method,
		"value":  value,
	})
	body := bytes.NewBuffer(postBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, body)
	if err != nil {
		return fmt.Errorf("can't create request %w", err)
	}
	req.Header.Set(api.ParamCipherKey, cipherKey)
	req.RequestURI = ""

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("can't create responce %w", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("can't set data %w", err)
	}
	fmt.Println(resp)
	return nil
}
