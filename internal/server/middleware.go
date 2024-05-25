package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/utils/regex"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

// Key is a type alias for a context key.
type Key string

const (
	SearchOptsKey Key = "opts"   // SearchOptsKey is the key to identify a SearchOptionsRecipes struct.
	UserIDKey     Key = "userID" // UserIDKey is the key to identify a user ID.
)

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
	"/ws":                                {},
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()

		now := time.Now()
		next.ServeHTTP(rec, r)
		elapsed := time.Since(now)

		result := rec.Result()

		var b bytes.Buffer
		r.Body = io.NopCloser(io.TeeReader(r.Body, &b))
		headers := map[string]string{
			"HX-Trigger":  rec.Header().Get("HX-Trigger"),
			"HX-Redirect": rec.Header().Get("HX-Redirect"),
		}

		s.Logger.LogAttrs(
			context.Background(),
			slog.LevelInfo,
			"HTTP roundtrip",
			slog.Int64("userID", getUserID(r)),
			slog.Group(
				"request",
				slog.String("method", r.Method),
				slog.String("uri", r.URL.EscapedPath()),
				slog.String("referrer", r.Header.Get("Referer")),
				slog.String("userAgent", r.Header.Get("User-Agent")),
				slog.String("ipAddress", getRemoteAddress(r)),
				slog.String("body", b.String()),
				slog.String("form", r.Form.Encode()),
			),
			slog.Group(
				"response",
				slog.Int("code", rec.Code),
				slog.String("duration", elapsed.String()),
				slog.String("body", rec.Body.String()),
				slog.Any("headers", headers),
			),
		)

		maps.Copy(w.Header(), result.Header)
		w.WriteHeader(rec.Code)
		_, _ = rec.Body.WriteTo(w)
	})
}

func getRemoteAddress(r *http.Request) string {
	realIP := r.Header.Get("X-Real-Ip")
	forwarded := r.Header.Get("X-Forwarded-For")
	if realIP == "" && forwarded == "" {
		i := strings.LastIndex(r.RemoteAddr, ":")
		if i == -1 {
			return r.RemoteAddr
		}
		return r.RemoteAddr[:i]
	}

	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	return realIP
}

func (s *Server) onlyAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromSessionCookie(r)
		if userID == -1 {
			userID = getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken)
		}

		if userID != 1 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Access denied: You are not an admin.")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) redirectIfLoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Config.Server.IsAutologin {
			ctx := context.WithValue(r.Context(), UserIDKey, int64(1))
			http.Redirect(w, r.WithContext(ctx), "/recipes", http.StatusSeeOther)
			return
		}

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

func redirectIfNoSignupsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Config.Server.IsNoSignups {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) mustBeLoggedInMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Config.Server.IsAutologin {
			ctx := context.WithValue(r.Context(), UserIDKey, int64(1))

			if SessionData.Data == nil {
				SessionData.Data = make(map[uuid.UUID]int64)
			}

			var isFound bool
			for _, v := range SessionData.Data {
				if v == 1 {
					isFound = true
					break
				}
			}

			if !isFound {
				sid := uuid.New()
				SessionData.Set(sid, 1)
				http.SetCookie(w, NewSessionCookie(sid.String()))
			}

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		uri := r.URL.Path
		if uri != "/ws" {
			_, found := excludedURIs[uri]
			if found || regex.WildcardURL.MatchString(uri) {
				uri = "/"
			}
			http.SetCookie(w, NewRedirectCookie(uri))
		}

		userID := getUserIDFromSessionCookie(r)
		if userID != -1 {
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID = getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken)
		if userID != -1 {
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
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
