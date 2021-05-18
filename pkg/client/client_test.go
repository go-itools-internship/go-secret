package client

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestClient_SetByKey(t *testing.T) {

	//srv := &http.Server{Addr: ":" + port, Handler: func() {}} // handler for server
	srvHttp := httptest.NewServer()

	client := New("key")
	err := client.SetByKey()
	if err != nil {
		fmt.Println(err)
	}

}

func TestClient_GetByKey(t *testing.T) {

}
