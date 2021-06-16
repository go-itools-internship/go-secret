/*
Package http provides handlers for different methods for REST API.
*/
package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-itools-internship/go-secret/pkg/secret"
)

const (
	ParamMethodKey = "method"
	ParamCipherKey = "cipher"
	ParamGetterKey = "key"
)

// MethodFactoryFunc type specifies signature to create/fetch a provider
// based on cipher key.
// Returns also a tear down function that should be called provider's work is done.
// Tear down function could be nil so need to check for nil before invoking.
type MethodFactoryFunc func(cipher string) (secret.Provider, func())

type methods struct {
	ss     map[string]MethodFactoryFunc
	logger *zap.SugaredLogger
}

// NewMethods initializes a structure that provides HTTP handler functions
// to organize REST API access to different type of provides based on "method" type.
//
// Accepts `ss` map with a set of method-provider pair.
func NewMethods(ss map[string]MethodFactoryFunc, logger *zap.SugaredLogger) *methods {
	return &methods{
		ss:     ss,
		logger: logger,
	}
}

// GetByKey method fetches a value specified by getter key.
// Uses cipher key to access encrypted data.
// Requires to provide getter key, cipher key (as a header) and method type to access as.
func (a *methods) GetByKey(w http.ResponseWriter, r *http.Request) {
	getterKey := r.URL.Query().Get(ParamGetterKey)
	if getterKey == "" {
		a.writeErrorResponse(w, http.StatusBadRequest, errors.New("cannot find getter key: empty"))
		return
	}

	actionType := r.URL.Query().Get(ParamMethodKey)
	if _, ok := a.ss[actionType]; !ok {
		a.writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("cannot find provided method type %s", actionType))
		return
	}

	cipherKey := r.Header.Get(ParamCipherKey)
	p, tearDownFn := a.ss[actionType](cipherKey)
	if tearDownFn != nil {
		defer tearDownFn()
	}

	result, err := p.GetData([]byte(getterKey))
	if err != nil {
		a.writeErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("cannot get data by key: %w", err))
		return
	}

	var responseBody struct {
		Value string `json:"value"`
	}
	responseBody.Value = string(result)
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		a.writeErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("cannot write response: %w", err))
	}
}

// SetByKey method sets a new value for specified getter key.
// Value is encrypted using cipher key (provided in header).
// Requires to provide getter key, cipher key (as a header) and method type to access as.
//
// Example of request body:
//
//    {
//        "getter": "cloud-key",
//        "method": "memory",
//        "value": "123-456"
//    }
func (a *methods) SetByKey(w http.ResponseWriter, r *http.Request) {
	logger := a.logger.Named("set-by-key")
	var requestBody struct {
		GetterKey  string `json:"getter"`
		MethodType string `json:"method"`
		Value      string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		a.writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("cannot decode body: %w", err))
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Warnf("cannot close request body: %s", err.Error())
		}
	}()

	if requestBody.GetterKey == "" {
		a.writeErrorResponse(w, http.StatusBadRequest, errors.New("cannot find getter key: empty"))
		return
	}
	if _, ok := a.ss[requestBody.MethodType]; !ok {
		a.writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("cannot find provided method type %s", requestBody.MethodType))
		return
	}

	cipherKey := r.Header.Get(ParamCipherKey)
	p, tearDownFn := a.ss[requestBody.MethodType](cipherKey)
	if tearDownFn != nil {
		defer tearDownFn()
	}

	err := p.SetData([]byte(requestBody.GetterKey), []byte(requestBody.Value))
	if err != nil {
		a.writeErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("cannot set data: %w", err))
		return
	}
	logger.Infof(requestBody.Value)
	w.WriteHeader(http.StatusNoContent)
}

func (a *methods) writeErrorResponse(w http.ResponseWriter, status int, response error) {
	logger := a.logger.Named("write-error-response")
	w.WriteHeader(status)
	if response != nil {
		if _, err := fmt.Fprintf(w, `{"error":"%s"}`, response.Error()); err != nil {
			logger.Warnf("cannot write response body: %s", err.Error())
		}
	}
}
