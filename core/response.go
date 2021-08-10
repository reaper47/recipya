package core

import (
	"encoding/json"
	"net/http"

	"github.com/reaper47/recipya/api"
)

func writeSuccessJson(object interface{}, w http.ResponseWriter) {
	addHeadersJson(w)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(object)
}

func writeCreatedJson(object interface{}, w http.ResponseWriter) {
	addHeadersJson(w)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(object)
}

func writeErrorJson(code int, message string, w http.ResponseWriter) {
	payload := api.ErrorJson{
		Objects: api.Error{
			Code:    code,
			Message: message,
			Status:  http.StatusText(code),
		}}

	addHeadersJson(w)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func addHeadersJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
