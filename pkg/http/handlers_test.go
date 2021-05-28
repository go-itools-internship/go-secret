package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/go-itools-internship/go-secret/pkg/secret"
	"github.com/stretchr/testify/require"
)

const jsonTerminator = "\n"

func TestHTTPHandlers(t *testing.T) {
	t.Run("set by key", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			mockProvider := new(MockProvider)
			defer mockProvider.AssertExpectations(t)
			mockProvider.On("SetData", []byte("test-getter-1"), []byte("test-value-1")).Return(nil).Once()

			expectedSipherKey := "1234-5678"

			tearDownFnCounter := int64(1)
			a := NewMethods(map[string]MethodFactoryFunc{
				"test-method": func(cipher string) (secret.Provider, func()) {
					require.EqualValues(t, expectedSipherKey, cipher)
					return mockProvider, func() {
						atomic.AddInt64(&tearDownFnCounter, -1)
					}
				},
			}, createSugarLogger())

			s := httptest.NewServer(http.HandlerFunc(a.SetByKey))
			defer s.Close()

			body := bytes.NewBufferString(`{"getter":"test-getter-1","method":"test-method","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, s.URL, body)
			req.Header.Set(ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusNoContent, resp.StatusCode)

			require.Eventually(t, func() bool {
				return atomic.LoadInt64(&tearDownFnCounter) == 0
			}, 300*time.Millisecond, 100*time.Millisecond)
		})

		t.Run("error when method is not found", func(t *testing.T) {
			a := NewMethods(map[string]MethodFactoryFunc{
				"test-method": func(cipher string) (secret.Provider, func()) {
					t.Error("test method should not be invoked")
					return nil, func() {
						t.Errorf("tear down should not be invoked")
					}
				},
			}, createSugarLogger())

			s := httptest.NewServer(http.HandlerFunc(a.SetByKey))
			defer s.Close()

			body := bytes.NewBufferString(`{"getter":"test-getter-1","method":"test-method-1","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, s.URL, body)
			req.Header.Set(ParamCipherKey, "")
			req.RequestURI = ""

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.EqualValues(t, `{"error":"cannot find provided method type test-method-1"}`, respBody)
		})

		t.Run("error when provider returned error", func(t *testing.T) {
			mockProvider := new(MockProvider)
			defer mockProvider.AssertExpectations(t)
			mockProvider.On("SetData", []byte("test-getter-1"), []byte("test-value-1")).Return(fmt.Errorf("test error")).Once()

			expectedSipherKey := "1234-5678"

			a := NewMethods(map[string]MethodFactoryFunc{
				"test-method": func(cipher string) (secret.Provider, func()) {
					require.EqualValues(t, expectedSipherKey, cipher)
					return mockProvider, nil
				},
			}, createSugarLogger())

			s := httptest.NewServer(http.HandlerFunc(a.SetByKey))
			defer s.Close()

			body := bytes.NewBufferString(`{"getter":"test-getter-1","method":"test-method","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, s.URL, body)
			req.Header.Set(ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.EqualValues(t, `{"error":"cannot set data: test error"}`, respBody)
		})

		t.Run("error when cannot decode request body", func(t *testing.T) {
			mockProvider := new(MockProvider)
			defer mockProvider.AssertExpectations(t)

			a := NewMethods(map[string]MethodFactoryFunc{}, createSugarLogger())

			s := httptest.NewServer(http.HandlerFunc(a.SetByKey))
			defer s.Close()

			body := bytes.NewBufferString(`abcd`)
			req := httptest.NewRequest(http.MethodPost, s.URL, body)
			req.RequestURI = ""

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Contains(t, string(respBody), `{"error":"cannot decode body:`)
		})
	})

	t.Run("get by key", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			mockProvider := new(MockProvider)
			defer mockProvider.AssertExpectations(t)
			mockProvider.On("GetData", []byte("test-getter-1")).Return([]byte("test-value-1"), nil).Once()

			expectedSipherKey := "1234-5678"

			tearDownFnCounter := int64(1)
			a := NewMethods(map[string]MethodFactoryFunc{
				"test-method": func(cipher string) (secret.Provider, func()) {
					require.EqualValues(t, expectedSipherKey, cipher)
					return mockProvider, func() {
						atomic.AddInt64(&tearDownFnCounter, -1)
					}
				},
			}, createSugarLogger())

			s := httptest.NewServer(http.HandlerFunc(a.GetByKey))
			defer s.Close()

			req := httptest.NewRequest(http.MethodGet, s.URL, nil)
			req.RequestURI = ""
			req.Header.Set(ParamCipherKey, expectedSipherKey)
			query := req.URL.Query()
			query.Set(ParamGetterKey, "test-getter-1")
			query.Set(ParamMethodKey, "test-method")
			req.URL.RawQuery = query.Encode()

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusOK, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.EqualValues(t, `{"value":"test-value-1"}`+jsonTerminator, string(respBody))

			require.Eventually(t, func() bool {
				return atomic.LoadInt64(&tearDownFnCounter) == 0
			}, 300*time.Millisecond, 100*time.Millisecond)
		})

		t.Run("error when no getter provided", func(t *testing.T) {
			mockProvider := new(MockProvider)
			defer mockProvider.AssertExpectations(t)

			a := NewMethods(map[string]MethodFactoryFunc{}, createSugarLogger())

			s := httptest.NewServer(http.HandlerFunc(a.GetByKey))
			defer s.Close()

			req := httptest.NewRequest(http.MethodGet, s.URL, nil)
			req.RequestURI = ""
			query := req.URL.Query()
			query.Set(ParamGetterKey, "")

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.EqualValues(t, `{"error":"cannot find getter key: empty"}`, respBody)
		})

		t.Run("error when no method type exists", func(t *testing.T) {
			mockProvider := new(MockProvider)
			defer mockProvider.AssertExpectations(t)

			tearDownFnCounter := int64(1)
			a := NewMethods(map[string]MethodFactoryFunc{
				"test-method": func(cipher string) (secret.Provider, func()) {
					t.Errorf("test method should not be invoked")
					return mockProvider, func() {
						atomic.AddInt64(&tearDownFnCounter, -1)
					}
				},
			}, createSugarLogger())

			s := httptest.NewServer(http.HandlerFunc(a.GetByKey))
			defer s.Close()

			req := httptest.NewRequest(http.MethodGet, s.URL, nil)
			req.RequestURI = ""
			query := req.URL.Query()
			query.Set(ParamGetterKey, "test-getter-1")
			query.Set(ParamMethodKey, "test-method-1")
			req.URL.RawQuery = query.Encode()

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.EqualValues(t, `{"error":"cannot find provided method type test-method-1"}`, respBody)

			require.Eventually(t, func() bool {
				return atomic.LoadInt64(&tearDownFnCounter) == 1 // NOTE: do not expect to be changed
			}, 300*time.Millisecond, 100*time.Millisecond)
		})

		t.Run("error when GetData returned error", func(t *testing.T) {
			mockProvider := new(MockProvider)
			defer mockProvider.AssertExpectations(t)
			mockProvider.On("GetData", []byte("test-getter-1")).Return(nil, fmt.Errorf("test error")).Once()

			expectedSipherKey := "1234-5678"

			a := NewMethods(map[string]MethodFactoryFunc{
				"test-method": func(cipher string) (secret.Provider, func()) {
					require.EqualValues(t, expectedSipherKey, cipher)
					return mockProvider, nil
				},
			}, createSugarLogger())

			s := httptest.NewServer(http.HandlerFunc(a.GetByKey))
			defer s.Close()

			req := httptest.NewRequest(http.MethodGet, s.URL, nil)
			req.RequestURI = ""
			req.Header.Set(ParamCipherKey, expectedSipherKey)
			query := req.URL.Query()
			query.Set(ParamGetterKey, "test-getter-1")
			query.Set(ParamMethodKey, "test-method")
			req.URL.RawQuery = query.Encode()

			resp, err := s.Client().Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Contains(t, string(respBody), `{"error":"cannot get data by key:`)
		})
	})
}

func createSugarLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	sugar := logger.Sugar()
	return sugar
}
