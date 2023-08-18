package server

import (
	"context"
	"net/http"
	"slices"
)

func (s *Server) redirectIfLoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if userID := getUserIDFromSessionCookie(r); userID != -1 {
			ctx := context.WithValue(r.Context(), "userID", userID)
			http.Redirect(w, r.WithContext(ctx), "/", http.StatusSeeOther)
			return
		}

		if userID := getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken); userID != -1 {
			ctx := context.WithValue(r.Context(), "userID", userID)
			w.Header().Set("HX-Redirect", "/")
			http.Redirect(w, r.WithContext(ctx), "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) mustBeLoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		excludedURIs := []string{"/user-initials", "/settings/tabs/recipes", "/settings/tabs/profile"}
		if !slices.Contains(excludedURIs, r.RequestURI) {
			http.SetCookie(w, NewRedirectCookie(r.RequestURI))
		}

		if userID := getUserIDFromSessionCookie(r); userID != -1 {
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if userID := getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken); userID != -1 {
			ctx := context.WithValue(r.Context(), "userID", userID)
			w.Header().Set("HX-Redirect", "/")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", "/auth/login")
			w.WriteHeader(http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		}
	})
}
