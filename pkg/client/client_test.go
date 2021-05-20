package client

import (
	"context"
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
	t.Run("set by key", func(t *testing.T) {
		t.Run("set by key success", func(t *testing.T) {
			ctx := context.Background()
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
				require.EqualValues(t, r.Header.Get(api.ParamCipherKey), cipherKey)
			}))
			defer s.Close()

			path := s.URL
			client := New(path)
			err := client.SetByKey(ctx, getter, value, method, cipherKey)
			require.NoError(t, err)
			require.EqualValues(t, path, s.URL)
		})
		t.Run("context error,when set by key with closed context", func(t *testing.T) {
			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, 10*time.Millisecond)
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(100 * time.Millisecond)
				w.WriteHeader(http.StatusNoContent)
				require.EqualValues(t, r.Header.Get(api.ParamCipherKey), cipherKey)
			}))
			defer s.Close()

			path := s.URL
			client := New(path)
			err := client.SetByKey(ctx, getter, value, method, cipherKey)
			require.Error(t, err)
		})
		t.Run("expected error if set by key with wrong url", func(t *testing.T) {
			wrongCipherKey := "wrong"
			ctx := context.Background()
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
				r.Header.Set(api.ParamCipherKey, wrongCipherKey)
				require.NotEqualValues(t, r.Header.Get(api.ParamCipherKey), cipherKey)
			}))
			defer s.Close()
			path := s.URL
			client := New(path + path)
			err := client.SetByKey(ctx, getter, value, method, cipherKey)
			require.Error(t, err)
		})
	})
}

func TestClient_GetByKey(t *testing.T) {
	cipherKey := "c-key"
	getter := "key"
	requestData := "Test value"
	method := "cloud"
	t.Run("get by key", func(t *testing.T) {
		t.Run("get by key success", func(t *testing.T) {
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
			defer s.Close()

			path := s.URL
			client := New(path)
			data, err := client.GetByKey(ctx, getter, cipherKey, method)
			require.NoError(t, err)

			require.EqualValues(t, requestData, data)
			require.EqualValues(t, s.URL, client.url)
		})
		t.Run("expected error if get by key with wrong url", func(t *testing.T) {
			ctx := context.Background()

			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json := `{"value":"Test value"}`
				_, err := w.Write([]byte(json))
				require.NoError(t, err)
			}))
			defer s.Close()

			path := s.URL
			client := New(path + "wrong-url")
			data, err := client.GetByKey(ctx, getter, cipherKey, method)
			require.NotEqualValues(t, requestData, data)
			require.NotEqualValues(t, s.URL, client.url)
			require.Error(t, err)
		})
	})
}
