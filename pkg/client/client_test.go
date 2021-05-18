package client

import (
	"fmt"
	"testing"
)

func TestClient_SetByKey(t *testing.T) {

	//srv := &http.Server{Addr: ":" + port, Handler: func() {}} // handler for server
	//srvHttp := httptest.NewServer()
	cipherKey := "c-key"
	key := "key"
	value := "test-value"

	client := New("key")
	err := client.SetByKey(cipherKey, key, value)
	if err != nil {
		fmt.Println(err)
	}

}

func TestClient_GetByKey(t *testing.T) {

}
