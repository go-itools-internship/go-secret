package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	api "github.com/go-itools-internship/go-secret/pkg/http"
	"github.com/stretchr/testify/require"
)

func TestRoot_Server_Redis(t *testing.T) {
	t.Run("set by key", func(t *testing.T) {
		t.Run("expect redis set method success", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			port, err := GetFreePort()
			require.NoError(t, err)
			r := New()
			r.cmd.SetArgs([]string{"server", "--port", strconv.Itoa(port), "--redis-url", redisURL})
			go func() {
				err := r.Execute(ctx)
				if err != nil {
					fmt.Println(err)
				}
			}()
			time.Sleep(2 * time.Second)
			defer func() {
				require.NoError(t, os.Remove("file.txt"))
			}()

			client := http.Client{Timeout: time.Second}
			body := bytes.NewBufferString(`{"getter":"key-value","method":"remote","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:"+strconv.Itoa(port), body)
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusNoContent, resp.StatusCode)
			require.NoError(t, resp.Body.Close())
		})
	})
	t.Run("get by key", func(t *testing.T) {
		t.Run("expect bad request status if set local method and try get by remote redis method", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			port, err := GetFreePort()
			require.NoError(t, err)
			r := New()
			r.cmd.SetArgs([]string{"server", "--port", strconv.Itoa(port), "--redis-url", redisURL})
			go func() {
				err := r.Execute(ctx)
				if err != nil {
					fmt.Println(err)
				}
			}()
			time.Sleep(2 * time.Second)

			defer func() {
				require.NoError(t, os.Remove("file.txt"))
			}()

			client := http.Client{Timeout: time.Second}
			body := bytes.NewBufferString(`{"getter":"key-value","method":"local","value":"test-value-2"}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:"+strconv.Itoa(port), body)
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""
			resp, err := client.Do(req)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())

			req = httptest.NewRequest(http.MethodGet, "http://localhost:"+strconv.Itoa(port), nil)
			req.RequestURI = ""
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			query := req.URL.Query()
			query.Set(api.ParamGetterKey, key)
			query.Set(api.ParamMethodKey, "remote")
			req.URL.RawQuery = query.Encode()

			resp, err = client.Do(req)
			require.NoError(t, err)
			var requestBody struct {
				Value string `json:"value"`
			}
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&requestBody))
			require.EqualValues(t, "test-value-1", requestBody.Value)
			require.NoError(t, resp.Body.Close())
		})
		t.Run("expect bad request status if set remote redis method and try get by local method", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			port, err := GetFreePort()
			require.NoError(t, err)
			r := New()
			r.cmd.SetArgs([]string{"server", "--port", strconv.Itoa(port), "--redis-url", redisURL})
			go func() {
				err := r.Execute(ctx)
				if err != nil {
					fmt.Println(err)
				}
			}()
			time.Sleep(2 * time.Second)
			defer func() {
				require.NoError(t, os.Remove("file.txt"))
			}()

			client := http.Client{Timeout: time.Second}
			body := bytes.NewBufferString(`{"getter":"key-value","method":"remote","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:"+strconv.Itoa(port), body)
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""
			resp, err := client.Do(req)
			require.NoError(t, err)
			require.NoError(t, resp.Body.Close())

			req = httptest.NewRequest(http.MethodGet, "http://localhost:"+strconv.Itoa(port), nil)
			req.RequestURI = ""
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			query := req.URL.Query()
			query.Set(api.ParamGetterKey, key)
			query.Set(api.ParamMethodKey, "local")
			req.URL.RawQuery = query.Encode()

			resp, err = client.Do(req)
			require.NoError(t, err)
			data, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Contains(t, string(data), "cannot get data by key")
			require.NoError(t, resp.Body.Close())
		})
		t.Run("expect redis get method success", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			port, err := GetFreePort()
			require.NoError(t, err)
			r := New()
			r.cmd.SetArgs([]string{"server", "--port", strconv.Itoa(port), "--redis-url", redisURL})
			go func() {
				err := r.Execute(ctx)
				if err != nil {
					fmt.Println(err)
				}
			}()
			time.Sleep(2 * time.Second)
			defer func() {
				require.NoError(t, os.Remove("file.txt"))
			}()

			client := http.Client{Timeout: time.Second}
			body := bytes.NewBufferString(`{"getter":"key-value","method":"remote","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:"+strconv.Itoa(port), body)
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusNoContent, resp.StatusCode)
			require.NoError(t, resp.Body.Close())

			req = httptest.NewRequest(http.MethodGet, "http://localhost:"+strconv.Itoa(port), nil)
			req.RequestURI = ""
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			query := req.URL.Query()
			query.Set(api.ParamGetterKey, key)
			query.Set(api.ParamMethodKey, "remote")
			req.URL.RawQuery = query.Encode()

			resp, err = client.Do(req)
			require.NoError(t, err)
			_, err = ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, resp.Body.Close())
		})
	})
}
