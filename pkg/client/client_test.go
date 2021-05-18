package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_SetByKey(t *testing.T) {
	cipherKey := "c-key"
	getter := "key"
	method := "cloud"
	value := "test-value"

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	client := New(s.URL)
	err, resp := client.SetByKey(cipherKey, getter, value, method)
	if err != nil {
		fmt.Println(err)
	}
	require.EqualValues(t, http.StatusOK, resp.StatusCode)
}

func TestClient_GetByKey(t *testing.T) {
	cipherKey := "c-key"
	getter := "key"
	requestData := "data"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// put data
		_, err := w.Write([]byte(requestData))
		require.NoError(t, err)
	}))

	client := New(s.URL)
	data, err := client.GetByKey(cipherKey, getter)
	if err != nil {
		fmt.Println(err)
	}
	require.EqualValues(t, requestData, data)

}
