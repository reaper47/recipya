package server

import (
	"errors"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"strings"
	"time"
)

const (
	cookieNameRedirect   = "redirect"
	cookieNameRememberMe = "remember_me"
	cookieNameSession    = "session"
)

// SessionData maps a UUID to a user id. It's used to track who is logged in session-wise.
var SessionData map[uuid.UUID]int64

// NewRedirectCookie creates a URL redirection cookie for an anonymous user.
func NewRedirectCookie(uri string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieNameRedirect,
		Value:    uri,
		Path:     "/",
		Secure:   app.Config.Server.IsProduction,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

// NewSessionCookie creates a session cookie for a logged-in user.
func NewSessionCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieNameSession,
		Value:    value,
		Path:     "/",
		Secure:   app.Config.Server.IsProduction,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func getUserIDFromSessionCookie(r *http.Request) int64 {
	c, err := r.Cookie(cookieNameSession)
	if err != nil {
		return -1
	}

	if c.MaxAge == -1 {
		return -1
	}

	sid, err := uuid.Parse(c.Value)
	if err != nil {
		return -1
	}

	if userID, ok := SessionData[sid]; ok {
		return userID
	}
	return -1
}

// NewRememberMeCookie creates a cookie for when the user checks remember on login.
func NewRememberMeCookie(selector, validator string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieNameRememberMe,
		Value:    selector + ":" + validator,
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Secure:   app.Config.Server.IsProduction,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func getUserIDFromRememberMeCookie(r *http.Request, getAuthToken func(selector, validator string) (models.AuthToken, error)) int64 {
	c, err := r.Cookie(cookieNameRememberMe)
	if errors.Is(err, http.ErrNoCookie) || c == nil {
		return -1
	}

	parts := strings.Split(c.Value, ":")
	if len(parts) != 2 {
		return -1
	}

	token, err := getAuthToken(parts[0], parts[1])
	if err != nil {
		return -1
	}
	return token.UserID
}
