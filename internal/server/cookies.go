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
	cookieNameRememberMe = "remember_me"
	cookieNameSession    = "session"
)

// NewSessionCookie creates a session cookie for a logged-in user.
func NewSessionCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieNameSession,
		Value:    value,
		Path:     "/",
		Secure:   app.Config.IsCookieSecure(),
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

	userID, ok := SessionData.Get(sid)
	if ok {
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
		Secure:   app.Config.IsCookieSecure(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
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
	if err != nil || token.IsExpired() {
		return -1
	}

	for _, id := range SessionData.Data {
		if id == token.UserID {
			return token.UserID
		}
	}
	return -1
}
