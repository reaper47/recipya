package models

import (
	"encoding/json"
	"net/http"
)

type errorWrapper struct {
	Error errorDetails `json:"error"`
}

type errorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NewErrorJson returns the error message as JSON.
func NewErrorJson(code int, message string) ([]byte, error) {
	e := errorDetails{
		Code:    code,
		Message: message,
		Status:  http.StatusText(code),
	}
	xb, err := json.Marshal(errorWrapper{Error: e})
	if err != nil {
		return nil, err
	}
	return xb, nil
}
