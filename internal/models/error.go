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

// NewErrorJSON returns the error message as JSON.
func NewErrorJSON(code int, message string) ([]byte, error) {
	return json.Marshal(errorWrapper{
		Error: errorDetails{
			Code:    code,
			Message: message,
			Status:  http.StatusText(code),
		},
	})
}
