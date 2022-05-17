package router

import (
	"context"
	"net/http"

	"github.com/reaper47/recipya/internal/constants"
	"github.com/reaper47/recipya/internal/repository"
)

type authenticationMiddleware struct{}

// Middleware is the authentication middleware.
func (a *authenticationMiddleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		s, isAuthenticated := repository.IsAuthenticated(req)
		if !isAuthenticated || s.UserID == -1 {
			http.Redirect(w, req, "/auth/signin", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(req.Context(), constants.UserID, s)
		req = req.WithContext(ctx)

		next.ServeHTTP(w, req)
	}
}
