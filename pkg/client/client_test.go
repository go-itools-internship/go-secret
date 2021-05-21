package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/go-itools-internship/go-secret/pkg/http"

	"github.com/stretchr/testify/require"
)

func TestClient_SetByKey(t *testing.T) {
	cipherKey := "c-key"
	getter := "key"
	method := "cloud"
	value := "test-value"
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.EqualValues(t, "/", r.URL.Path)
			require.EqualValues(t, r.Header.Get(api.ParamCipherKey), cipherKey)
			w.WriteHeader(http.StatusNoContent)

			var requestBody struct {
				GetterKey  string `json:"getter"`
				MethodType string `json:"method"`
				Value      string `json:"value"`
			}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&requestBody))
			require.EqualValues(t, requestBody.Value, value)
			require.EqualValues(t, requestBody.MethodType, method)
			require.EqualValues(t, requestBody.GetterKey, getter)
		}))
		defer s.Close()

		c := New(s.URL)
		require.EqualValues(t, s.URL, c.url)
		err := c.SetByKey(ctx, getter, value, method, cipherKey)
		require.NoError(t, err)

	})
	t.Run("context error if closed context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer s.Close()

		c := New(s.URL)
		err := c.SetByKey(ctx, getter, value, method, cipherKey)
		require.Error(t, err)
	})
	t.Run("expected error if server does not respond", func(t *testing.T) {
		ctx := context.Background()
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer s.Close()
		wrongURL := "http://127.0.0.1:8888"
		c := New(wrongURL)
		require.NotEqual(t, s.URL, wrongURL)
		err := c.SetByKey(ctx, getter, value, method, cipherKey)
		require.Error(t, err)
		require.Contains(t, err.Error(), "secret client: can't do request")
	})
	t.Run("expected error if wrong status code", func(t *testing.T) {
		ctx := context.Background()
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("test bad request body"))
			require.NoError(t, err)
		}))
		defer s.Close()

		c := New(s.URL)
		err := c.SetByKey(ctx, getter, value, method, cipherKey)
		require.Error(t, err)
		require.EqualValues(t, "secret client: can't set data: body: \"test bad request body\", status code: 400", err.Error())
	})
}

func TestClient_GetByKey(t *testing.T) {
	cipherKey := "c-key"
	getter := "key"
	requestData := "Test value"
	method := "cloud"
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.EqualValues(t, "/", r.URL.Path)
			require.EqualValues(t, r.Header.Get(api.ParamCipherKey), cipherKey)
			require.EqualValues(t, r.URL.Query().Get(api.ParamGetterKey), getter)
			require.EqualValues(t, r.URL.Query().Get(api.ParamMethodKey), method)
			testJSON := `{"value":"Test value"}`
			_, err := w.Write([]byte(testJSON))
			require.NoError(t, err)
		}))
		defer s.Close()

		c := New(s.URL)
		require.EqualValues(t, s.URL, c.url)
		data, err := c.GetByKey(ctx, getter, method, cipherKey)
		require.NoError(t, err)
		require.EqualValues(t, requestData, data)
	})
	t.Run("expected error if server does not respond", func(t *testing.T) {
		ctx := context.Background()

		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			j := `{"value":"Test value"}`
			_, err := w.Write([]byte(j))
			require.NoError(t, err)
		}))
		defer s.Close()

		wrongURL := "http://127.0.0.1:8888"
		c := New(wrongURL)
		require.NotEqual(t, s.URL, wrongURL)
		data, err := c.GetByKey(ctx, getter, method, cipherKey)
		require.Error(t, err)
		require.Contains(t, err.Error(), "secret client: can't do response")
		require.Empty(t, data)
	})
	t.Run("expected error if no json", func(t *testing.T) {
		ctx := context.Background()
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))
		defer s.Close()

		c := New(s.URL)
		data, err := c.GetByKey(ctx, getter, method, cipherKey)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot decode body")
		require.Empty(t, data)
	})
	t.Run("expected error if wrong status code", func(t *testing.T) {
		ctx := context.Background()
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			j := `{"value":"test bad request body"}`
			_, err := w.Write([]byte(j))
			require.NoError(t, err)

		}))
		defer s.Close()

		c := New(s.URL)
		_, err := c.GetByKey(ctx, getter, method, cipherKey)
		require.Error(t, err)
		require.EqualValues(t, "secret client: can't get data: body: \"test bad request body\", status code: 400", err.Error())
	})
}
