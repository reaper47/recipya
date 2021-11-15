package model

import "net/http"

// HttpResponse is a general struct that stores
// information related to the HTTP request.
type HttpResponse struct {
	Url      string
	Response *http.Response
	Err      error
}
