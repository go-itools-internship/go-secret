package http

import (
	"bytes"
	"fmt"
	"github.com/go-itools-internship/go-secret/pkg/secret"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPAction(t *testing.T) {
	t.Run("set by key", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			mockProvider := new(MockProvider)
			mockProvider.On("SetData", []byte("test-getter-1"), []byte("test-value-1")).Return(nil).Once()

			expectedSipherKey := "1234-5678"

			a := NewActions(map[string]ActionFactoryFunc{
				"test-action": func(cipher string) (secret.Provider, func()) {
					require.EqualValues(t, expectedSipherKey, cipher)
					return mockProvider, func() {}
				},
			})

			s := httptest.NewServer(http.HandlerFunc(a.SetByKey))
			defer s.Close()

			body := bytes.NewBufferString(`{"getter":"test-getter-1","action":"test-action","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, s.URL, body)
			req.Header.Set(ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusNoContent, resp.StatusCode)

			mockProvider.AssertExpectations(t)
		})

		t.Run("error when action is not found", func(t *testing.T) {
			a := NewActions(map[string]ActionFactoryFunc{
				"test-action": func(cipher string) (secret.Provider, func()) {
					t.Error("test action should not be invoked")
					return nil, func() {}
				},
			})

			s := httptest.NewServer(http.HandlerFunc(a.SetByKey))
			defer s.Close()

			body := bytes.NewBufferString(`{"getter":"test-getter-1","action":"test-action-1","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, s.URL, body)
			req.Header.Set(ParamCipherKey, "")
			req.RequestURI = ""

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.EqualValues(t, `{"error": "cannot find provided action type"}`, respBody)
		})

		t.Run("error when provider returned error", func(t *testing.T) {
			mockProvider := new(MockProvider)
			mockProvider.On("SetData", []byte("test-getter-1"), []byte("test-value-1")).Return(fmt.Errorf("test error")).Once()

			expectedSipherKey := "1234-5678"

			a := NewActions(map[string]ActionFactoryFunc{
				"test-action": func(cipher string) (secret.Provider, func()) {
					require.EqualValues(t, expectedSipherKey, cipher)
					return mockProvider, func() {}
				},
			})

			s := httptest.NewServer(http.HandlerFunc(a.SetByKey))
			defer s.Close()

			body := bytes.NewBufferString(`{"getter":"test-getter-1","action":"test-action","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, s.URL, body)
			req.Header.Set(ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)

			mockProvider.AssertExpectations(t)
		})
	})

	t.Run("get by key", func(t *testing.T) {

	})
}
