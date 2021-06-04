package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-itools-internship/go-secret/pkg/io/storage"
	"github.com/go-redis/redis/v8"

	api "github.com/go-itools-internship/go-secret/pkg/http"
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
	t.Run("expect set data only redis storage", func(t *testing.T) {
		key := "key value"
		path := ""
		redisURL := "localhost:6379"
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		r := New()
		r.cmd.SetArgs([]string{"set", "--key", key, "--value", "test value", "--cipher-key", "ck", "--redis-url", redisURL})
		err := r.Execute(ctx)
		require.NoError(t, err)

		_, err = os.Open(path)
		require.Error(t, err)

		key = "key"
		encodedValue := "value"

		rdb := redis.NewClient(&redis.Options{Addr: redisURL, Password: "", DB: 0})
		s := storage.NewRedisVault(rdb)
		err = s.SaveData([]byte(key), []byte(encodedValue))
		require.NoError(t, err)
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
	key := "key-value"
	value := "60OBdPOOkSOu6kn8ZuMuXtAPVrUEFkPREydDwY6+ip/LrAFaHSc="
	path := "testFile.txt"
	t.Run("success", func(t *testing.T) {
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
	})
	t.Run("success set and read data from redis storage", func(t *testing.T) {
		path := ""
		redisURL := "localhost:6379"
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		r := New()
		r.cmd.SetArgs([]string{"set", "--key", key, "--value", value, "--cipher-key", "ck", "--redis-url", redisURL})
		err := r.Execute(ctx)
		require.NoError(t, err)

		_, err = os.Open(path)
		require.Error(t, err)

		rdb := redis.NewClient(&redis.Options{Addr: redisURL, Password: "", DB: 0})
		s := storage.NewRedisVault(rdb)
		err = s.SaveData([]byte(key), []byte(value))
		require.NoError(t, err)

		val, err := s.ReadData([]byte(key))
		require.NoError(t, err)
		require.EqualValues(t, value, val)
	})
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
			port := createAndExecuteCliCommand(ctx)
			defer func() {
				require.NoError(t, os.Remove(path))
			}()

			client := http.Client{Timeout: time.Second}
			body := bytes.NewBufferString(`{"getter":"key-value","method":"local","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:"+port, body)
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusNoContent, resp.StatusCode)
			require.NoError(t, resp.Body.Close())
		})
		t.Run("expect redis set method success", func(t *testing.T) {
			expectedSipherKey := "key value"
			redisURL := "localhost:6379"
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		t.Run("expect url not found error", func(t *testing.T) {
			expectedSipherKey := "key value"
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			port := createAndExecuteCliCommand(ctx)

			client := http.Client{Timeout: time.Second}
			body := bytes.NewBufferString(`{"getter":"key-value","method":"local","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:"+port+"/error", body)
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusNotFound, resp.StatusCode)
			require.NoError(t, resp.Body.Close())
		})
	})
	t.Run("get by key", func(t *testing.T) {
		t.Run("get method success", func(t *testing.T) {
			expectedSipherKey := "key value"
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			port := createAndExecuteCliCommand(ctx)
			defer func() {
				require.NoError(t, os.Remove(path))
			}()

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
			_, err = ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusOK, resp.StatusCode)
			require.NoError(t, resp.Body.Close())
		})
		t.Run("expect bad request status if set local method and try get by remote method", func(t *testing.T) {
			expectedSipherKey := "key value"
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			port, err := GetFreePort()
			require.NoError(t, err)
			r := New()
			r.cmd.SetArgs([]string{"server", "--port", strconv.Itoa(port), "--path", path})
			go func() {
				err := r.Execute(ctx)
				if err != nil {
					fmt.Println(err)
				}
			}()
			time.Sleep(2 * time.Second)

			client := http.Client{Timeout: time.Second}
			body := bytes.NewBufferString(`{"getter":"key-value","method":"remote","value":"test-value-1"}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:"+strconv.Itoa(port), body)
			req.Header.Set(api.ParamCipherKey, expectedSipherKey)
			req.RequestURI = ""

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
			require.NoError(t, resp.Body.Close())
		})
		t.Run("expect bad request status if set remote method and try get by local method", func(t *testing.T) {
			expectedSipherKey := "key value"
			redisURL := "localhost:6379"
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
			expectedSipherKey := "key value"
			redisURL := "localhost:6379"
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
		t.Run("middleware check success", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			port := createAndExecuteCliCommand(ctx)

			client := http.Client{Timeout: time.Second}
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/ping", port), nil)
			req.RequestURI = ""

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusOK, resp.StatusCode)
			require.EqualValues(t, "text/plain", resp.Header.Get("Content-Type"))
			require.NoError(t, resp.Body.Close())
		})
		t.Run("error when used wrong cipher key", func(t *testing.T) {
			wrongSipherKey := "wrong key"
			expectedSipherKey := "key value"
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			port := createAndExecuteCliCommand(ctx)
			defer func() {
				require.NoError(t, os.Remove(path))
			}()

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
			req.Header.Set(api.ParamCipherKey, wrongSipherKey)
			query := req.URL.Query()
			query.Set(api.ParamGetterKey, key)
			query.Set(api.ParamMethodKey, "local")
			req.URL.RawQuery = query.Encode()

			resp, err = client.Do(req)
			require.NoError(t, err)
			_, err = ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)
			require.NoError(t, resp.Body.Close())
		})
		t.Run("get method with error, when url not found error", func(t *testing.T) {
			invalidKey := "invalid-key"
			require.NotEqual(t, key, invalidKey)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			expectedSipherKey := "key value"

			port := createAndExecuteCliCommand(ctx)

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
		})
		t.Run("get method with error, when key does not exist on server", func(t *testing.T) {
			invalidKey := "invalid-key"
			require.NotEqual(t, key, invalidKey)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			expectedSipherKey := "key value"

			port := createAndExecuteCliCommand(ctx)
			defer func() {
				require.NoError(t, os.Remove(path))
			}()

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
		})
	})
}

func TestRoot_ServerPing(t *testing.T) {
	route := "/ping"
	t.Run("success", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.EqualValues(t, http.MethodGet, r.Method)
			require.EqualValues(t, route, r.URL.Path)
		}))
		defer s.Close()

		// Parse server url for wright flags format
		sURL, h, p := ParseURL(s.URL)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		r := New()
		r.cmd.SetArgs([]string{"server", "ping", "--url", fmt.Sprintf("%s://%s", sURL.Scheme, h), "--port", p, "--route", route})
		require.NoError(t, r.Execute(ctx))
	})
	t.Run("error when response status is not found", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte("test request body"))
			require.NoError(t, err)
		}))
		defer s.Close()
		// Parse server url for wright flags format
		sURL, h, p := ParseURL(s.URL)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		r := New()
		r.cmd.SetArgs([]string{"server", "ping", "--url", fmt.Sprintf("%s://%s", sURL.Scheme, h), "--port", p, "--route", route})
		err := r.Execute(ctx)
		require.Error(t, err)
		require.EqualValues(t, "server response is not expected: body \"test request body\", wrong status code 404", err.Error())
	})
	t.Run("error when server connection refused", func(t *testing.T) {
		testURL := "http://localhost"
		port := "8880"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		r := New()
		r.cmd.SetArgs([]string{"server", "ping", "--url", testURL, "--port", port, "--route", route})
		err := r.Execute(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("server response error: Get %q:", fmt.Sprintf("%s:%s%s", testURL, port, route)))
		require.Contains(t, err.Error(), fmt.Sprintf("%s: connect: connection refused", port))
	})
	t.Run("error when request time to make a request is over", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.EqualValues(t, http.MethodGet, r.Method)
			require.EqualValues(t, route, r.URL.Path)
			time.Sleep(3 * time.Second)
		}))
		defer s.Close()
		// Parse server url for wright flags format
		sURL, h, p := ParseURL(s.URL)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		r := New()
		r.cmd.SetArgs([]string{"server", "ping", "--url", fmt.Sprintf("%s://%s", sURL.Scheme, h), "--port", p, "--route", route, "--timeout", "2s"})
		err := r.Execute(ctx)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Client.Timeout exceeded while awaiting headers")
	})
}

func ParseURL(s string) (*url.URL, string, string) {
	serverURL, err := url.Parse(s)
	if err != nil {
		fmt.Println(err)
	}
	h, p, err := net.SplitHostPort(serverURL.Host)
	if err != nil {
		fmt.Println(err)
	}
	return serverURL, h, p
}

func createAndExecuteCliCommand(ctx context.Context) (freePort string) {
	path := "testFile.txt"
	port, err := GetFreePort()
	if err != nil {
		fmt.Println(err)
	}
	r := New()
	r.cmd.SetArgs([]string{"server", "--path", path, "--port", strconv.Itoa(port)})
	go func() {
		err := r.Execute(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(2 * time.Second)
	return strconv.Itoa(port)
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func(l *net.TCPListener) {
		err := l.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(l)
	return l.Addr().(*net.TCPAddr).Port, nil
}
