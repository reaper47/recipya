package server

import (
	"context"
	"github.com/reaper47/recipya/internal/utils/regex"
	"net/http"
	"strings"
)

// Key is a type alias for a context key.
type Key string

// UserIDKey is the key to identify a user ID.
const UserIDKey Key = "userID"

var excludedURIs = map[string]struct{}{
	"/auth/change-password":              {},
	"/auth/confirm":                      {},
	"/auth/register/validate-password":   {},
	"/auth/user":                         {},
	"/cookbooks/recipes/search":          {},
	"/integrations/import/nextcloud":     {},
	"/recipes/add/import":                {},
	"/recipes/add/manual/ingredient":     {},
	"/recipes/add/manual/instruction":    {},
	"/recipes/add/manual/ingredient/*":   {},
	"/recipes/add/manual/instruction/*":  {},
	"/recipes/add/ocr":                   {},
	"/recipes/add/request-website":       {},
	"/recipes/add/website":               {},
	"/recipes/search":                    {},
	"/recipes/supported-websites":        {},
	"/settings/export/recipes?type=json": {},
	"/settings/calculate-nutrition":      {},
	"/settings/convert-automatically":    {},
	"/settings/measurement-system":       {},
	"/settings/tabs/profile":             {},
	"/settings/tabs/recipes":             {},
	"/user-initials":                     {},
}

func (s *Server) redirectIfLoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromSessionCookie(r)
		if userID != -1 {
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			http.Redirect(w, r.WithContext(ctx), "/", http.StatusSeeOther)
			return
		}

		userID = getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken)
		if userID != -1 {
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			w.Header().Set("HX-Redirect", "/")
			http.Redirect(w, r.WithContext(ctx), "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) mustBeLoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		_, found := excludedURIs[r.RequestURI]
		if found || regex.WildcardURL.MatchString(r.RequestURI) {
			if strings.HasPrefix(r.RequestURI, "/settings") {
				uri = "/settings"
			} else {
				uri = "/"
			}
		}
		http.SetCookie(w, NewRedirectCookie(uri))

		userID := getUserIDFromSessionCookie(r)
		if userID != -1 {
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID = getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken)
		if userID != -1 {
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
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
