package server_test

import (
	"github.com/reaper47/recipya/internal/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func assertHeader(t testing.TB, rr *httptest.ResponseRecorder, key, value string) {
	t.Helper()
	got := rr.Result().Header.Get(key)
	if got != value {
		t.Fatalf("expected header %s to be %s but got %s", key, value, got)
	}
}

func assertMustBeLoggedIn(t *testing.T, srv *server.Server, method string, uri string) {
	t.Helper()
	rr := sendRequest(srv, method, uri, noHeader, nil)

	assertStatus(t, rr.Code, http.StatusSeeOther)
	assertHeader(t, rr, "Location", "/auth/login")
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("expected status %d but got %d", want, got)
	}
}

func assertStringsInHTML(t testing.TB, bodyHTML string, wants []string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(bodyHTML, want) {
			t.Fatalf("string not found in HTML:\n%s", want)
		}
	}
}

func assertStringsNotInHTML(t *testing.T, bodyHTML string, wants []string) {
	for _, want := range wants {
		if strings.Contains(bodyHTML, want) {
			t.Fatalf("string found in HTML when it should not:\n%s", want)
		}
	}
}
