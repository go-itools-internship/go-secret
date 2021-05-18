// Package client work with server api and set/get keys to it
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (c *client) GetByKey(cipherKey string, key string) (data string, err error) {
	req, err := http.NewRequest(http.MethodGet, c.url, nil)
	if err != nil {
		return "", fmt.Errorf("can't create request %w", err)
	}
	req.RequestURI = ""
	req.Header.Set(api.ParamCipherKey, cipherKey)
	query := req.URL.Query()
	query.Set(api.ParamGetterKey, key)
	query.Set(api.ParamMethodKey, "cloud")
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	//defer resp.Body.Close()
	fmt.Println(resp)

	//type Data struct{
	//	getter string
	//	method string
	//	value string
	//}
	//data1 := json.NewDecoder(resp.Body)
	//data1.Decode()
	body, err := ioutil.ReadAll(resp.Body)
	//I must get and return data
	//data :=
	return string(body), nil
}

func (c *client) SetByKey(cipherKey string, getterKey string, value string, method string) (error, *http.Response) {
	postBody, _ := json.Marshal(map[string]string{
		"getter": getterKey,
		"method": method,
		"value":  value,
	})
	body := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(http.MethodPost, c.url, body)
	if err != nil {
		return fmt.Errorf("can't create request %w", err), nil
	}
	req.Header.Set(api.ParamCipherKey, cipherKey)
	req.RequestURI = ""

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("can't create responce %w", err), nil
	}
	fmt.Println(resp)
	return nil, resp
}
