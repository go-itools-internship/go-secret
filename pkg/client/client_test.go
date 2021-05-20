package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/go-itools-internship/go-secret/pkg/http"

	"github.com/stretchr/testify/require"
)

func TestClient_SetByKey(t *testing.T) {
	cipherKey := "c-key"
	getter := "key"
	method := "cloud"
	value := "test-value"
	ctx := context.Background()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		require.EqualValues(t, r.Header.Get(api.ParamCipherKey), cipherKey)
	}))

	client := New(s.URL)
	err := client.SetByKey(ctx, getter, value, method, cipherKey)
	if err != nil {
		fmt.Println(err)
	}
}

func TestClient_GetByKey(t *testing.T) {
	cipherKey := "c-key"
	getter := "key"
	requestData := "Test value"
	method := "cloud"
	ctx := context.Background()

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json := `{"value":"Test value"}`
		_, err := w.Write([]byte(json))
		require.NotNil(t, w)
		require.EqualValues(t, r.URL.Query().Get(api.ParamGetterKey), getter)
		require.EqualValues(t, r.URL.Query().Get(api.ParamMethodKey), method)
		require.EqualValues(t, r.Header.Get(api.ParamCipherKey), cipherKey)
		require.NoError(t, err)
	}))

	client := New(s.URL)
	data, err := client.GetByKey(ctx, getter, cipherKey, method)
	if err != nil {
		fmt.Println(err)
	}
	require.EqualValues(t, requestData, data)

}
