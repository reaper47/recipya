package server_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
)

type header string

const (
	formData     header = "multipart/form-data"
	formHeader   header = "application/x-www-form-urlencoded"
	noHeader     header = ""
	promptHeader header = "prompt"
)

func createWSServer() (*server.Server, *httptest.Server, *websocket.Conn) {
	srv := newServerTest()
	repo := &mockRepository{
		categories:        map[int64][]string{1: {"breakfast", "lunch", "dinner"}},
		RecipesRegistered: make(map[int64]models.Recipes),
		UsersRegistered: []models.User{
			{ID: 1, Email: "test@example.com"},
		},
		UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
	}
	srv.Repository = repo

	if server.SessionData.Data == nil {
		server.SessionData.Data = make(map[uuid.UUID]int64)
	}

	sid := uuid.New()
	server.SessionData.Set(sid, 1)

	h := http.Header{}
	h.Add("Cookie", server.NewSessionCookie(sid.String()).String())

	ts := httptest.NewServer(srv.Router)
	u := strings.Replace(ts.URL, "http", "ws", 1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, _, err := websocket.Dial(ctx, u+"/ws", &websocket.DialOptions{HTTPHeader: h})
	if err != nil {
		panic(err)
	}
	return srv, ts, c
}

func createMultipartForm(fields map[string][]string) (contentType string, body string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for name, values := range fields {
		for _, value := range values {
			if strings.HasSuffix(value, ".jpg") {
				field, _ := writer.CreateFormFile(name, value)
				if field == nil {
					log.Println("createMultipartForm.CreateFormField: field is nil after writer.CreateFormField")
					return "", ""
				}
				_, _ = field.Write([]byte("not a real file"))
			} else {
				field, _ := writer.CreateFormField(name)
				if field == nil {
					log.Println("createMultipartForm.CreateFormField: field is nil after writer.CreateFormField")
					return "", ""
				}
				_, _ = field.Write([]byte(value))
			}
		}
	}
	_ = writer.Close()

	return writer.FormDataContentType(), buf.String()
}

func sendRequestNoBody(srv *server.Server, method, target string) *httptest.ResponseRecorder {
	return sendRequest(srv, method, target, noHeader, nil)
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

func sendRequestAsLoggedInNoBody(srv *server.Server, method, target string) *httptest.ResponseRecorder {
	return sendRequestAsLoggedIn(srv, method, target, noHeader, nil)
}

func sendRequestAsLoggedIn(srv *server.Server, method, target string, contentType header, body *strings.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, prepareRequest(method, target, contentType, body))
	return rr
}

func sendRequestAsLoggedInOtherNoBody(srv *server.Server, method, target string) *httptest.ResponseRecorder {
	return sendRequestAsLoggedInOther(srv, method, target, noHeader, nil)
}

func sendRequestAsLoggedInOther(srv *server.Server, method, target string, contentType header, body *strings.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	srv.Router.ServeHTTP(rr, prepareRequestOther(method, target, contentType, body))
	return rr
}

func sendHxRequestAsLoggedInNoBody(srv *server.Server, method, target string) *httptest.ResponseRecorder {
	return sendHxRequestAsLoggedIn(srv, method, target, noHeader, nil)
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

func sendHxRequestAsLoggedInOtherNoBody(srv *server.Server, method, target string) *httptest.ResponseRecorder {
	return sendHxRequestAsLoggedInOther(srv, method, target, noHeader, nil)
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

	sid := uuid.New()
	if server.SessionData.Data == nil {
		server.SessionData.Data = make(map[uuid.UUID]int64)
	}
	server.SessionData.Set(sid, 1)

	r := httptest.NewRequest(method, target, body)
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

	sid := uuid.New()
	if server.SessionData.Data == nil {
		server.SessionData.Data = make(map[uuid.UUID]int64)
	}
	server.SessionData.Set(sid, 2)

	r := httptest.NewRequest(method, target, body)
	r.AddCookie(server.NewSessionCookie(sid.String()))
	r = r.WithContext(context.WithValue(r.Context(), server.UserIDKey, int64(2)))

	if contentType != noHeader {
		r.Header.Set("Content-Type", string(contentType))
	}
	return r
}

func getBodyHTML(rr *httptest.ResponseRecorder) string {
	body, _ := io.ReadAll(rr.Body)

	cases := []models.Replace{
		{Old: "\r\n", New: ""},
		{Old: "\n", New: ""},
		{Old: "\r", New: ""},
		{Old: "&#39;", New: "'"},
		{Old: "&#34;", New: `"`},
		{Old: "&lt;", New: "<"},
		{Old: "&gt;", New: ">"},
	}
	for _, c := range cases {
		body = bytes.ReplaceAll(body, []byte(c.Old), []byte(c.New))
	}

	body = bytes.Join(bytes.Fields(body), []byte(" "))
	return string(body)
}

func readMessage(tb testing.TB, c *websocket.Conn, number int) (websocket.MessageType, []byte) {
	tb.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var (
		mt   websocket.MessageType
		data []byte
		err  error
	)

	for i := 0; i < number; i++ {
		time.Sleep(200 * time.Millisecond)
		mt, data, err = c.Read(ctx)
		if err != nil {
			tb.Fatalf("failed to read message: %v", err)
		}
	}

	return mt, data
}
