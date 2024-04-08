package server_test

import (
	"encoding/json"
	"fmt"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestHandlers_Auth_ChangePassword(t *testing.T) {
	srv := newServerTest()

	const uri = "/auth/change-password"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("mistake in form", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("password-current=test1&password-new=test2&password-confirm=test"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Passwords do not match.\",\"title\":\"Form Error\"}"}`)
	})

	t.Run("Wrong current password", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("password-current=cat&password-new=test2&password-confirm=test2"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Current password is incorrect.\",\"title\":\"Form Error\"}"}`)
	})

	t.Run("new password same as current", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("password-current=test2&password-new=test2&password-confirm=test2"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"New password is same as current.\",\"title\":\"Form Error\"}"}`)
	})

	t.Run("cannot change password if autologin", func(t *testing.T) {
		app.Config.Server.IsAutologin = true
		defer func() {
			app.Config.Server.IsAutologin = false
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("password-current=test2&password-new=test2&password-confirm=test2"))

		assertStatus(t, rr.Code, http.StatusForbidden)
	})

	t.Run("valid form", func(t *testing.T) {
		repo := &mockRepository{
			UsersRegistered: []models.User{
				{ID: 1, Email: "test@example.com"},
			},
		}
		srv.Repository = repo

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("password-current=test1&password-new=test2&password-confirm=test2"))

		assertStatus(t, rr.Code, http.StatusNoContent)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-info\",\"message\":\"\",\"title\":\"Password updated.\"}"}`)
		if !slices.Contains(repo.UsersUpdated, 1) {
			t.Fatal("must have called UpdatePassword on the user")
		}
	})
}

func TestHandlers_Auth_Confirm(t *testing.T) {
	srv := newServerTest()
	srv.Repository = &mockRepository{
		UsersRegistered: []models.User{
			{ID: 1, Email: "test@example.com"},
		},
	}

	const uri = "/auth/confirm"

	t.Run("missing token", func(t *testing.T) {
		rr := sendRequest(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "Location", "/")
	})

	t.Run("invalid token for existing user", func(t *testing.T) {
		token, _ := auth.CreateToken(map[string]any{"userID": 2}, 1*time.Nanosecond)

		rr := sendRequest(srv, http.MethodGet, uri+"?token="+token, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		want := []string{
			`<title hx-swap-oob="true">Token Expired | Recipya</title>`,
			"The token associated with the URL expired.",
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("user does not exist", func(t *testing.T) {
		token, _ := auth.CreateToken(map[string]any{"userID": 2}, 1*time.Hour)

		rr := sendRequest(srv, http.MethodGet, uri+"?token="+token, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusNotFound)
		want := []string{
			`<title hx-swap-oob="true">Confirm Error | Recipya</title>`,
			"An error occurred when you requested to confirm your account.",
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("valid confirmation token for existing user", func(t *testing.T) {
		token, _ := auth.CreateToken(map[string]any{"userID": 1}, 1*time.Hour)

		rr := sendRequest(srv, http.MethodGet, uri+"?token="+token, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Success | Recipya</title>`,
			"Your account has been confirmed.",
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Auth_DeleteUser(t *testing.T) {
	srv := newServerTest()
	originalRepo := &mockRepository{
		UsersRegistered: []models.User{
			{ID: 1},
			{ID: 2},
		},
	}
	srv.Repository = originalRepo

	uri := "/auth/user"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodDelete, uri)
	})

	t.Run("demo version cannot be deleted", func(t *testing.T) {
		repo := &mockRepository{
			UsersRegistered: []models.User{
				{ID: 1, Email: "demo@demo.com"},
				{ID: 2},
			},
		}
		srv.Repository = repo
		app.Config.Server.IsDemo = true
		defer func() {
			app.Config.Server.IsDemo = false
			srv.Repository = originalRepo
		}()

		rr := sendRequestAsLoggedIn(srv, http.MethodDelete, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusTeapot)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Your savings account has been deleted.\",\"title\":\"General Error\"}"}`)
	})

	t.Run("cannot delete user if autologin", func(t *testing.T) {
		app.Config.Server.IsAutologin = true
		defer func() {
			app.Config.Server.IsAutologin = false
		}()

		rr := sendRequestAsLoggedIn(srv, http.MethodDelete, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusForbidden)
	})

	t.Run("user cannot delete other user", func(t *testing.T) {
		repo := originalRepo
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInOther(srv, http.MethodDelete, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", "/")
		if !slices.ContainsFunc(repo.UsersRegistered, func(user models.User) bool { return user.ID == 1 }) {
			t.Fatal("user 1 should not have been deleted")
		}
	})

	t.Run("valid request", func(t *testing.T) {
		clear(server.SessionData.Data)
		originalNumSessions := len(server.SessionData.Data) + 1
		repo := originalRepo
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", "/")
		if slices.ContainsFunc(repo.UsersRegistered, func(user models.User) bool { return user.ID == 1 }) {
			t.Fatal("user 1 should have been deleted")
		}
		if len(server.SessionData.Data) == originalNumSessions {
			t.Fatalf("expected one less number of sessions")
		}
	})

	t.Run("delete a deleted user does nothing", func(t *testing.T) {
		repo := originalRepo
		defer func() {
			srv.Repository = originalRepo
		}()
		rr := httptest.NewRecorder()
		r := prepareRequest(http.MethodDelete, uri, noHeader, nil)
		srv.Router.ServeHTTP(rr, r)

		rr = httptest.NewRecorder()
		srv.Router.ServeHTTP(rr, r)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		if slices.ContainsFunc(repo.UsersRegistered, func(user models.User) bool { return user.ID == 1 }) {
			t.Fatal("user 1 should have been deleted")
		}
	})
}

func TestHandlers_Auth_ForgotPassword(t *testing.T) {
	emailMock := &mockEmail{}
	repo := &mockRepository{
		UsersRegistered: []models.User{{ID: 1, Email: "test@example.com"}},
	}
	srv := newServerTest()
	srv.Email = emailMock
	srv.Repository = repo

	uri := "/auth/forgot-password"

	t.Run("anonymous user accesses forgot password", func(t *testing.T) {
		rr := sendRequest(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Forgot Password | Recipya</title>`,
			`<input required type="email" placeholder="Enter your email address" class="input input-bordered w-full" name="email">`,
			`<button class="btn btn-primary btn-block btn-sm">Reset password</button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("logged-in user cannot access forgot password form", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", "/settings")
	})

	t.Run("logged in user cannot post", func(t *testing.T) {
		numHits := emailMock.hitCount

		rr := sendRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("email=not@exist.com"))

		assertStatus(t, rr.Code, http.StatusForbidden)
		if emailMock.hitCount == numHits+1 {
			t.Fatal("an email should not have been sent")
		}
	})

	testcases := []struct {
		name        string
		email       string
		isEmailSent bool
		want        []string
	}{
		{
			name:  "user does not exist",
			email: "email=not@exist.com",
		},
		{
			name:        "user exists",
			email:       "test@example.com",
			isEmailSent: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			numHits := emailMock.hitCount

			rr := sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader("email="+tc.email))

			assertStatus(t, rr.Code, http.StatusOK)
			if tc.isEmailSent && emailMock.hitCount != numHits+1 {
				t.Fatal("an email should have been sent")
			} else if !tc.isEmailSent && emailMock.hitCount == numHits+1 {
				t.Fatal("an email should not have been sent")
			}
			want := []string{
				`<h2 class="card-title underline self-center">Password Reset Requested</h2>`,
				`An email with instructions on how to reset your password has been sent to you. Please check your inbox and follow the provided steps to regain access to your account.`,
				`<a href="/" class="btn btn-primary btn-block btn-sm">Back Home</a>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}

	testcases2 := []struct {
		name        string
		tokenExpire int
	}{
		{name: "reset password link has expired token", tokenExpire: -1},
		{name: "reset password link has invalid token"},
		{name: "reset password link has no token"},
	}
	for _, tc := range testcases2 {
		t.Run(tc.name, func(t *testing.T) {
			var token string
			if tc.tokenExpire != 0 {
				token, _ = auth.CreateToken(map[string]any{"userID": 1}, time.Duration(tc.tokenExpire)*time.Second)
			} else if strings.Contains(tc.name, "invalid") {
				token = "hello"
			}

			rr := sendRequest(srv, http.MethodGet, fmt.Sprintf("%s/reset?token=%s", uri, token), noHeader, nil)

			assertStatus(t, rr.Code, http.StatusBadRequest)
			want := []string{
				`<title hx-swap-oob="true">Token Expired | Recipya</title>`,
				"The token associated with the URL expired.",
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}

	t.Run("reset password link is valid", func(t *testing.T) {
		token, _ := auth.CreateToken(map[string]any{"userID": 1}, 1*time.Second)

		rr := sendRequest(srv, http.MethodGet, fmt.Sprintf("%s/reset?token=%s", uri, token), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Reset Password | Recipya</title>`,
			`<input name="user-id" type="hidden" value="1">`,
			`<input required type="password" placeholder="Enter your new password" class="input input-bordered w-full" name="password">`,
			`<input required type="password" placeholder="Retype your password" class="input input-bordered w-full" name="password-confirm">`,
			`<button class="btn btn-primary btn-block btn-sm">Change</button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("reset password invalid", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+"/reset", formHeader, strings.NewReader("user-id=1&password=test&password-confirm=test2"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Password is invalid.\",\"title\":\"Form Error\"}"}`)
	})

	t.Run("reset password valid", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+"/reset", formHeader, strings.NewReader("user-id=1&password=test&password-confirm=test"))

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-info\",\"message\":\"Password updated.\",\"title\":\"\"}"}`)
		assertHeader(t, rr, "HX-Redirect", "/auth/login")
		if !slices.Contains(repo.UsersUpdated, 1) {
			t.Fatal("must have called UpdatePassword on the user")
		}
	})
}

func TestHandlers_Auth_Login(t *testing.T) {
	srv := newServerTest()
	repo := &mockRepository{
		UsersRegistered: []models.User{
			{ID: 1, Email: "test@example.com"},
		},
	}
	srv.Repository = repo

	const uri = "/auth/login"

	testcases := []struct {
		name string
		form string
	}{
		{name: "email invalid", form: "email=hello@test&password=123"},
		{name: "email invalid", form: "email=hello@test.com"},
		{name: "user not found", form: "email=hello@test.com&password=123"},
	}
	for _, tc := range testcases {
		t.Run("invalid credentials when "+tc.name, func(t *testing.T) {
			rr := sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader(tc.form))

			assertStatus(t, rr.Code, http.StatusNoContent)
			var got map[string]string
			_ = json.Unmarshal([]byte(rr.Header().Get("HX-Trigger")), &got)
			wantHeader := `{"action":"","background":"alert-error","message":"Credentials are invalid.","title":"Form Error"}`
			if got["showToast"] != wantHeader {
				t.Fatalf("got\n%q\nbut want\n%q", got["showToast"], wantHeader)
			}
		})
	}

	t.Run("redirect to home when logged in", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "Location", "/")
	})

	t.Run("redirect to accessed uri after logged in", func(t *testing.T) {
		otherURI := "/recipes/add"
		r := httptest.NewRequest(http.MethodPost, uri, strings.NewReader("email=test@example.com&password=123&remember-me=false"))
		r.Header.Set("Content-Type", string(formHeader))
		r.AddCookie(server.NewRedirectCookie(otherURI))

		rr := httptest.NewRecorder()
		srv.Router.ServeHTTP(rr, r)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", otherURI)
	})

	t.Run("redirect to index if autologin enabled", func(t *testing.T) {
		app.Config.Server.IsAutologin = true
		defer func() {
			app.Config.Server.IsAutologin = false
		}()

		rr := sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader("email=test@example.com&password=123&remember-me=false"))

		assertStatus(t, rr.Code, http.StatusSeeOther)
	})

	t.Run("hide signup button when registrations disabled", func(t *testing.T) {
		app.Config.Server.IsNoSignups = true
		defer func() {
			app.Config.Server.IsNoSignups = false
		}()

		rr := sendRequest(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		body := getBodyHTML(rr)
		assertStringsInHTML(t, body, []string{
			`<input required type="email" placeholder="Enter your email address" class="input input-bordered w-full" name="email">`,
			`<input required type="password" placeholder="Enter your password" class="input input-bordered w-full" name="password">`,
			`<input type="checkbox" class="checkbox checkbox-primary" name="remember-me" value="yes">`,
			`<button class="btn btn-primary btn-block btn-sm">Log In</button>`,
			`<a class="btn btn-sm btn-ghost" href="/auth/forgot-password">Forgot your password?</a>`,
		})
		assertStringsNotInHTML(t, body, []string{
			`<p class="text-center">Don't have an account?</p>`,
			`<a class="btn btn-sm btn-block btn-outline" href="/auth/register">Sign Up</a>`,
		})
	})

	t.Run("login  successful", func(t *testing.T) {
		clear(server.SessionData.Data)

		rr := sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader("email=test@example.com&password=123&remember-me=false"))

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", "/")
		var (
			isUserInSession   bool
			isCookieStoresSID bool
		)
		for sid, userID := range server.SessionData.Data {
			if userID == 1 {
				isUserInSession = true
				isCookieStoresSID = slices.ContainsFunc(rr.Result().Cookies(), func(cookie *http.Cookie) bool {
					return cookie.Name == "session" && cookie.Value == sid.String()
				})
				break
			}
		}
		if !isUserInSession {
			t.Fatal("expected user to be in the server's session data")
		}
		if !isCookieStoresSID {
			t.Fatal("expected session ID to be stored in a cookie named 'session'")
		}
	})

	t.Run("user checked remember me", func(t *testing.T) {
		numAuthTokensBefore := len(repo.AuthTokens)

		rr := sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader("email=test@example.com&password=123&remember-me=yes"))

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", "/")

		cookies := rr.Result().Cookies()
		index := slices.IndexFunc(cookies, func(cookie *http.Cookie) bool { return cookie.Name == "remember_me" })
		if index == -1 {
			t.Fatal("there must be a session cookie")
		}
		if cookies[index].Expires.Before(time.Now().Add(30 * 24 * time.Hour).Add(-1 * time.Minute)) {
			t.Fatalf("got expiration %v but want an expiration of 1 month", cookies[index].Expires)
		}

		if len(repo.AuthTokens) != numAuthTokensBefore+1 {
			t.Fatal("expected an authentication token to be added to the database")
		}
	})

	t.Run("user checked remember me and accesses login page", func(t *testing.T) {
		rr := sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader("email=test@example.com&password=123&remember-me=yes"))
		r := httptest.NewRequest(http.MethodGet, uri, nil)
		for _, c := range rr.Result().Cookies() {
			if c.Name == "remember_me" {
				r.AddCookie(c)
			}
		}

		rr = httptest.NewRecorder()
		srv.Router.ServeHTTP(rr, r)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", "/")
	})
}

func TestHandlers_Auth_Logout(t *testing.T) {
	srv := newServerTest()
	repo := &mockRepository{}
	srv.Repository = repo

	const uri = "/auth/logout"

	t.Run("cannot log out a user who is already logged out", func(t *testing.T) {
		originalNumSessions := len(server.SessionData.Data)

		rr := sendRequest(srv, http.MethodPost, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusNoContent)
		if originalNumSessions != len(server.SessionData.Data) {
			t.Fatalf("expected same number of sessions")
		}
	})

	t.Run("valid logout for a logged-in user", func(t *testing.T) {
		clear(server.SessionData.Data)
		originalNumSessions := len(server.SessionData.Data) + 1

		rr := sendRequestAsLoggedIn(srv, http.MethodPost, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		if len(server.SessionData.Data) != originalNumSessions-1 {
			t.Fatalf("expected one less number of sessions")
		}
		var isCookieInvalid bool
		for _, c := range rr.Result().Cookies() {
			if c.Name == "session" {
				isCookieInvalid = c.MaxAge == -1
			}
		}
		if !isCookieInvalid {
			t.Fatal("expected the session cookie to be invalidated")
		}
	})

	t.Run("cannot logout when autologin enabled", func(t *testing.T) {
		app.Config.Server.IsAutologin = true
		defer func() {
			app.Config.Server.IsAutologin = false
		}()

		rr := repo.sendRequestAsLoggedInRememberMe(srv, http.MethodPost, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusForbidden)
	})

	t.Run("remember me user has its token deleted on logout", func(t *testing.T) {
		originalNumAuthTokens := len(repo.AuthTokens) + 1

		rr := repo.sendRequestAsLoggedInRememberMe(srv, http.MethodPost, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		var isCookieInvalid bool
		for _, c := range rr.Result().Cookies() {
			if c.Name == "remember_me" {
				isCookieInvalid = c.MaxAge == -1
			}
		}
		if !isCookieInvalid {
			t.Fatal("expected the remember me cookie to be invalidated")
		}
		if len(repo.AuthTokens) != originalNumAuthTokens-1 {
			t.Fatal("expected one less auth token in the database")
		}
	})
}

func TestHandlers_Auth_Register(t *testing.T) {
	srv := newServerTest()

	const uri = "/auth/register"

	t.Run("redirect to home when logged in", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
	})

	t.Run("invalid registration user already exists", func(t *testing.T) {
		email := "invalid@register.com"
		_ = sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader("email="+email+"&password=test123&password-confirm=test123"))
		originalNumUsers := len(srv.Repository.Users())

		rr := sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader("email="+email+"&password=test123&password-confirm=test123"))

		assertStatus(t, rr.Code, http.StatusUnprocessableEntity)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"User might be registered or password invalid.\",\"title\":\"Database Error\"}"}`)
		if len(srv.Repository.Users()) != originalNumUsers {
			t.Fatalf("expected no user to be added to the db of %d users", originalNumUsers)
		}
	})

	t.Run("redirect to index autologin enabled", func(t *testing.T) {
		app.Config.Server.IsAutologin = true
		defer func() {
			app.Config.Server.IsAutologin = false
		}()

		rrGet := sendRequest(srv, http.MethodGet, uri, noHeader, nil)
		rrPost := sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader("email=test@test.com&password=test123&password-confirm=test123"))

		assertStatus(t, rrGet.Code, http.StatusSeeOther)
		assertStatus(t, rrPost.Code, http.StatusSeeOther)
	})

	t.Run("cannot register when no signups disabled", func(t *testing.T) {
		originalNumUsers := len(srv.Repository.Users())
		app.Config.Server.IsNoSignups = true
		defer func() {
			app.Config.Server.IsNoSignups = false
		}()

		rrGet := sendRequest(srv, http.MethodGet, uri, noHeader, nil)
		rrPost := sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader("email=test@test.com&password=test123&password-confirm=test123"))

		assertStatus(t, rrGet.Code, http.StatusSeeOther)
		assertHeader(t, rrGet, "Location", "/auth/login")
		assertStatus(t, rrPost.Code, http.StatusSeeOther)
		if len(srv.Repository.Users()) != originalNumUsers {
			t.Fatal("expected no users to be registered")
		}
	})

	t.Run("valid registration for new user", func(t *testing.T) {
		email := "regsiter@valid.com"
		originalNumUsers := len(srv.Repository.Users())

		rr := sendRequest(srv, http.MethodPost, uri, formHeader, strings.NewReader("email="+email+"&password=test123&password-confirm=test123"))

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", "/auth/login")
		users := srv.Repository.Users()
		if len(users) != originalNumUsers+1 {
			t.Fatalf("expected %d users but got %d", originalNumUsers+1, len(users))
		}
		if users[len(users)-1].Email != email {
			t.Fatalf("got user %+v but want %s", users[len(users)-1], email)
		}
	})
}
