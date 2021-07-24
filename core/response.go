package core

import (
	"encoding/json"
	"net/http"

	"github.com/reaper47/recipya/api"
)

func writeSuccessJson(object interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(object)
}

func writeErrorJson(code int, message string, w http.ResponseWriter) {
	payload := api.ErrorJson{
		Objects: api.Error{
			Code:    code,
			Message: message,
			Status:  http.StatusText(code),
		}}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(payload)
}
