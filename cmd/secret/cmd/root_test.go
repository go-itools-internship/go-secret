package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	api "github.com/go-itools-internship/go-secret/pkg/http"
	"github.com/phayes/freeport"

	"github.com/stretchr/testify/require"
)

func TestRoot_Set(t *testing.T) {
	t.Run("expect one keys", func(t *testing.T) {
		key := "key value"
		path := "testFile.txt"
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		r := New()
		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--path", path})
		err := r.Execute(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.Remove(path))
		}()

		testFile, err := os.Open(path)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, testFile.Close())
		}()

		fileData := make(map[string]string)
		require.NoError(t, json.NewDecoder(testFile).Decode(&fileData))

		var got string
		require.Len(t, fileData, 1)
		for key := range fileData {
			got = key
			break // we iterate one time to get first key
		}
		require.EqualValues(t, key, got)
	})

	t.Run("expect two keys", func(t *testing.T) {
		firstKey := "first key"
		secondKey := "second key"
		path := "testFile.txt"
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		r := New()
		r.cmd.SetArgs([]string{"set", "--key", firstKey, "--value", "test value", "--cipher-key", "ck", "--path", path})
		err := r.Execute(ctx)
		require.NoError(t, err)

		r2 := New()
		r2.cmd.SetArgs([]string{"set", "--key", secondKey, "--value", "test value", "--cipher-key", "ck", "--path", path})
		err = r2.Execute(ctx)
		require.NoError(t, err)

		require.NoError(t, err)
		defer func() {
			require.NoError(t, os.Remove(path))
		}()

		testFile, err := os.Open(path)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, testFile.Close())
		}()

		fileData := make(map[string]string)
		require.NoError(t, json.NewDecoder(testFile).Decode(&fileData))

		require.Len(t, fileData, 2)
		require.Contains(t, fileData, firstKey)
		require.Contains(t, fileData, secondKey)
	})
}

func TestRoot_Get(t *testing.T) {
	key := "key value"
	value := "60OBdPOOkSOu6kn8ZuMuXtAPVrUEFkPREydDwY6+ip/LrAFaHSc="
	path := "testFile.txt"
	file, err := os.Create(path)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(path))
	}()

	defer func() {
		require.NoError(t, file.Close())
	}()

	fileTestData := make(map[string]string)
	fileTestData[key] = value
	require.NoError(t, json.NewEncoder(file).Encode(&fileTestData))

	r := New()
	r.cmd.SetArgs([]string{"get", "--key", key, "--cipher-key", "ck", "--path", path})
	executeErr := r.Execute(ctx)
	require.NoError(t, executeErr)

	testFile, err := os.Open(path)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, testFile.Close())
	}()

	fileData := make(map[string]string)
	require.NoError(t, json.NewDecoder(testFile).Decode(&fileData))
	var got string
	require.Len(t, fileData, 1)
	for _, value := range fileData {
		got = value
		break // we iterate one time to get first value
	}
	require.EqualValues(t, value, got)
}

func TestRoot_Server(t *testing.T) {
	key := "key-value"
	path := "testFile.txt"
	t.Run("set by key", func(t *testing.T) {
		t.Run("expect set method success", func(t *testing.T) {
			expectedSipherKey := "key value"
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			// get free port after cli creating
			port := createAndExecuteCliCommand(ctx, t)

			client := http.Client{Timeout: time.Second}
			body := bytes.NewBufferString(`{"getter":"key-value","method":"local","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:"+port, body)
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusNoContent, resp.StatusCode)
			require.NoError(t, resp.Body.Close())

			defer func() {
				require.NoError(t, os.Remove(path))
			}()
		})
		t.Run("expect url not found error", func(t *testing.T) {
			expectedSipherKey := "key value"
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			port := createAndExecuteCliCommand(ctx, t)

			client := http.Client{Timeout: time.Second}
			body := bytes.NewBufferString(`{"getter":"key-value","method":"local","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:"+port+"/error", body)
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusNotFound, resp.StatusCode)
			require.NoError(t, resp.Body.Close())

			defer func() {
				require.NoError(t, os.Remove(path))
			}()
		})
	})
	t.Run("get by key", func(t *testing.T) {
		t.Run("get method success", func(t *testing.T) {
			expectedSipherKey := "key value"
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			port := createAndExecuteCliCommand(ctx, t)

			client := http.Client{Timeout: time.Second}

			body := bytes.NewBufferString(`{"getter":"key-value","method":"local","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:"+port, body)
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""
			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusNoContent, resp.StatusCode)
			require.NoError(t, resp.Body.Close())

			req = httptest.NewRequest(http.MethodGet, "http://localhost:"+port, nil)
			req.RequestURI = ""
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			query := req.URL.Query()
			query.Set(api.ParamGetterKey, key)
			query.Set(api.ParamMethodKey, "local")
			req.URL.RawQuery = query.Encode()

			resp, err = client.Do(req)
			require.NoError(t, err)
			respBody, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			fmt.Println(string(respBody))
			require.EqualValues(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, resp.Body.Close())

			defer func() {
				require.NoError(t, os.Remove(path))
			}()
		})
		t.Run("get method with error, when url not found error", func(t *testing.T) {
			invalidKey := "invalid-key"
			require.NotEqual(t, key, invalidKey)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			expectedSipherKey := "key value"

			port := createAndExecuteCliCommand(ctx, t)

			client := http.Client{Timeout: time.Second}
			req := httptest.NewRequest(http.MethodGet, "http://localhost:"+port+"/errorUrl", nil)
			req.RequestURI = ""
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			query := req.URL.Query()
			query.Set(api.ParamGetterKey, invalidKey)
			query.Set(api.ParamMethodKey, "local")
			req.URL.RawQuery = query.Encode()

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusNotFound, resp.StatusCode)
			require.NoError(t, resp.Body.Close())

			defer func() {
				require.NoError(t, os.Remove(path))
			}()
		})
		t.Run("get method with error, when key does not exist on server", func(t *testing.T) {
			invalidKey := "invalid-key"
			require.NotEqual(t, key, invalidKey)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			expectedSipherKey := "key value"

			port := createAndExecuteCliCommand(ctx, t)

			client := http.Client{Timeout: time.Second}
			req := httptest.NewRequest(http.MethodGet, "http://localhost:"+port, nil)
			req.RequestURI = ""
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			query := req.URL.Query()
			query.Set(api.ParamGetterKey, invalidKey)
			query.Set(api.ParamMethodKey, "local")
			req.URL.RawQuery = query.Encode()

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)
			require.NoError(t, resp.Body.Close())

			defer func() {
				require.NoError(t, os.Remove(path))
			}()
		})
	})
}

func createAndExecuteCliCommand(ctx context.Context, t *testing.T) (freePort string) {
	key := "key-value"
	path := "testFile.txt"
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	r := New()
	r.cmd.SetArgs([]string{"server", "--cipher-key", key, "--path", path, "--port", strconv.Itoa(port)})
	go func() {
		err := r.Execute(ctx)
		require.Error(t, err)
	}()
	time.Sleep(2 * time.Second)
	return strconv.Itoa(port)
}
