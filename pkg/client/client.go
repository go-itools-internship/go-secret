// Package client work with server api and set/get keys to it
package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	api "github.com/go-itools-internship/go-secret/pkg/http"
)

type client struct {
	client http.Client
	url    string
}

func New(url string) *client {
	client1 := http.Client{Timeout: time.Second}
	newClient := &client{client: client1, url: url}
	return newClient
}

func (c *client) GetByKey(cipherKey string, key string, ctx context.Context) (data string, err error) {
	req, err := http.NewRequest(http.MethodGet, c.url, nil)
	if err != nil {
		return "", fmt.Errorf("can't create request %w", err)
	}
	req.RequestURI = ""
	req.Header.Set(api.ParamCipherKey, cipherKey)
	query := req.URL.Query()
	query.Set(api.ParamGetterKey, key)
	query.Set(api.ParamMethodKey, "local")
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	defer resp.Body.Close()
	fmt.Println(resp)
	//I must get and return data
	//data :=

	return "", nil
}

func (c *client) SetByKey(cipherKey string, key string, value string) error {
	//req := httptest.NewRequest(http.MethodPost, "http://localhost:"+port, body)

	body := bytes.NewBufferString(`{"getter":"key-value","method":"local","value":""}`)
	req, err := http.NewRequest(http.MethodPost, c.url, body)
	if err != nil {
		return fmt.Errorf("can't create request %w", err)
	}
	req.Header.Set(api.ParamCipherKey, cipherKey)
	req.Header.Set(api.ParamMethodKey, "")
	req.Header.Set(api.ParamGetterKey, key)
	req.Header.Set("value", value)
	req.RequestURI = ""

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("can't create responce %w", err)
	}
	fmt.Println(resp)
	return nil
}
