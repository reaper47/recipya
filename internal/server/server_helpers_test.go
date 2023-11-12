package server_test

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/server"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
)

type header string

const (
	formHeader   header = "application/x-www-form-urlencoded"
	noHeader     header = ""
	promptHeader header = "prompt"
)

var (
	regexRmSpaces = regexp.MustCompile(`>\s+<`)
	regexOneLine  = regexp.MustCompile(`\s{2,}|\n`)
)

func createMultipartForm(fields map[string]string) (contentType string, body string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for name, value := range fields {
		if strings.HasSuffix(value, ".jpg") {
			field, _ := writer.CreateFormFile(name, value)
			_, _ = field.Write([]byte("not a real file"))
		} else {
			field, _ := writer.CreateFormField(name)
			_, _ = field.Write([]byte(value))
		}
	}
	_ = writer.Close()

	return writer.FormDataContentType(), buf.String()
}

func sendRequest(srv *server.Server, method, target string, contentType header, body *strings.Reader) *httptest.ResponseRecorder {
	if body == nil {
		body = strings.NewReader("")
	}
	r := httptest.NewRequest(method, target, body)

	if contentType != noHeader {
		r.Header.Set("Content-Type", string(contentType))
	}

	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, r)
	return rr
}

func sendRequestAsLoggedIn(srv *server.Server, method, target string, contentType header, body *strings.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, prepareRequest(method, target, contentType, body))
	return rr
}

func sendRequestAsLoggedInOther(srv *server.Server, method, target string, contentType header, body *strings.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, prepareRequestOther(method, target, contentType, body))
	return rr
}

func sendHxRequestAsLoggedIn(srv *server.Server, method, target string, contentType header, body *strings.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r := prepareRequest(method, target, contentType, body)
	if contentType == promptHeader {
		b, _ := io.ReadAll(body)
		r.Header.Set("HX-Prompt", string(b))
	}
	r.Header.Set("HX-Request", "true")
	srv.Router.ServeHTTP(rr, r)
	return rr
}

func sendHxRequestAsLoggedInOther(srv *server.Server, method, target string, contentType header, body *strings.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r := prepareRequestOther(method, target, contentType, body)
	if contentType == promptHeader {
		b, _ := io.ReadAll(body)
		r.Header.Set("HX-Prompt", string(b))
	}
	r.Header.Set("HX-Request", "true")
	srv.Router.ServeHTTP(rr, r)
	return rr
}

func (m *mockRepository) sendRequestAsLoggedInRememberMe(srv *server.Server, method, target string, contentType header, body *strings.Reader) *httptest.ResponseRecorder {
	r := prepareRequest(method, target, contentType, body)
	selector, validator := auth.GenerateSelectorAndValidator()
	_ = m.AddAuthToken(selector, validator, 1)
	r.AddCookie(server.NewRememberMeCookie(selector, validator))
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, r)
	return rr
}

func prepareRequest(method, target string, contentType header, body *strings.Reader) *http.Request {
	if body == nil {
		body = strings.NewReader("")
	}
	r := httptest.NewRequest(method, target, body)
	sid := uuid.New()
	server.SessionData[sid] = 1
	r.AddCookie(server.NewSessionCookie(sid.String()))
	r = r.WithContext(context.WithValue(r.Context(), server.UserIDKey, int64(1)))

	if contentType != noHeader {
		r.Header.Set("Content-Type", string(contentType))
	}
	return r
}

func prepareRequestOther(method, target string, contentType header, body *strings.Reader) *http.Request {
	if body == nil {
		body = strings.NewReader("")
	}
	r := httptest.NewRequest(method, target, body)
	sid := uuid.New()
	server.SessionData[sid] = 2
	r.AddCookie(server.NewSessionCookie(sid.String()))
	r = r.WithContext(context.WithValue(r.Context(), server.UserIDKey, int64(2)))

	if contentType != noHeader {
		r.Header.Set("Content-Type", string(contentType))
	}
	return r
}

func getBodyHTML(rr *httptest.ResponseRecorder) string {
	body, _ := io.ReadAll(rr.Body)
	bodyStr := strings.ReplaceAll(string(body), "\n", "")
	bodyStr = regexOneLine.ReplaceAllString(bodyStr, " ")
	bodyStr = regexRmSpaces.ReplaceAllString(bodyStr, "><")
	return bodyStr
}
