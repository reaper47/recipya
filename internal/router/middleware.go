package router

import (
	"net/http"

	"github.com/reaper47/recipya/internal/repository"
)

type authenticationMiddleware struct {
	tokenUsers map[string]string
}

func (m *authenticationMiddleware) Populate() {

}

func (m *authenticationMiddleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !repository.IsAuthenticated(w, req) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, req)
	}
}
