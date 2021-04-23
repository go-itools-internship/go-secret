package http

import (
	"encoding/json"
	"github.com/go-itools-internship/go-secret/pkg/secret"
	"log"
	"net/http"
)

const (
	ParamActionKey = "action"
	ParamCipherKey = "cipher"
	ParamGetterKey = "key"
)

type ActionFactoryFunc func(cipher string) (secret.Provider, func())

type actions struct {
	ss map[string]ActionFactoryFunc
}

func NewActions(ss map[string]ActionFactoryFunc) *actions {
	return &actions{
		ss: ss,
	}
}

func (a *actions) GetByKey(w http.ResponseWriter, r *http.Request) {
	getterKey := r.URL.Query().Get(ParamGetterKey)
	if getterKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	actionType := r.URL.Query().Get(ParamActionKey)
	if _, ok := a.ss[actionType]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cipherKey := r.Header.Get(ParamCipherKey)
	p, tearDownFn := a.ss[actionType](cipherKey)
	if tearDownFn != nil {
		defer tearDownFn()
	}

	result, err := p.GetData([]byte(getterKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var responseBody struct {
		Value string `json:"value"`
	}
	responseBody.Value = string(result)
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

/*
{
	"getter": "cloud-key",
	"action": "memory",
	"value": "123-456"
}
*/
func (a *actions) SetByKey(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		GetterKey  string `json:"getter"`
		ActionType string `json:"action"`
		Value      string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("cannot close request body: ", err.Error())
		}
	}()

	if requestBody.GetterKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, ok := a.ss[requestBody.ActionType]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"error": "cannot find provided action type"}`)); err != nil {
			log.Println("cannot write response body: ", err.Error())
		}
		return
	}

	cipherKey := r.Header.Get(ParamCipherKey)
	p, tearDownFn := a.ss[requestBody.ActionType](cipherKey)
	if tearDownFn != nil {
		defer tearDownFn()
	}

	err := p.SetData([]byte(requestBody.GetterKey), []byte(requestBody.Value))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
