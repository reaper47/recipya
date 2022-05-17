package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
	"github.com/reaper47/recipya/internal/router/handlers"
)

func TestIndex(t *testing.T) {
	t.Run("unauthenticated", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		rr := sendRequest(req, handlers.Index)

		assertStatusCode(t, rr.Code, http.StatusOK)
		approvals.VerifyWithExtension(t, rr.Body, ".html")
	})

	t.Run("authenticated redirect to /recipes", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(testCookie)

		rr := sendRequest(req, handlers.Index)

		assertStatusCode(t, rr.Code, http.StatusSeeOther)
		assertLocation(t, rr, "/recipes")
	})
}

func sendRequest(req *http.Request, f func(http.ResponseWriter, *http.Request)) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(rr, req)
	return rr
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %d but want %d", got, want)
	}
}

func assertLocation(t *testing.T, rr *httptest.ResponseRecorder, want string) {
	t.Helper()

	got, _ := rr.Result().Location()
	if got.String() != want {
		t.Fatalf("got Location %q but want %q", got, want)
	}
}
