// Copyright 2016 Nevio Vesic
// Please check out LICENSE file for more information about limitations
// MIT License

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	disposable "github.com/0x19/disposable/protos"
	uuid "github.com/satori/go.uuid"
)

// GetExternalIP - Will check and get current machine external IP address
func GetExternalIP() (string, error) {
	rsp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf)), nil
}

// StringInSlice - Will check if string in list. This is equivalent to python if x in []
func StringInSlice(str string, list []string) bool {
	for _, value := range list {
		if value == str {
			return true
		}
	}
	return false
}

// DecodeJSONBody -
func DecodeJSONBody(model interface{}, rc io.ReadCloser) error {
	decoder := json.NewDecoder(rc)
	if err := decoder.Decode(model); err != nil {
		return err
	}

	return nil
}

func DecodeRequestBody(i interface{}, body io.Reader) *disposable.DisposableResponse {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(i); err != nil {
		return &disposable.DisposableResponse{
			Status:    false,
			RequestId: uuid.NewV4().String(),
			Error: &disposable.Error{
				Message: ErrorJSONParseError,
				Type:    TypeJSONParseError,
				Info: map[string]string{
					"error": err.Error(),
				},
			},
		}
	}

	return nil
}
