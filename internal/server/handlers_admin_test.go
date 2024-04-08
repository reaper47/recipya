package server_test

import (
	"fmt"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"strings"
	"testing"
)

func TestHandlers_Admin(t *testing.T) {
	srv := newServerTest()

	const uri = "/admin"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("admin is user with ID 1", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		body := getBodyHTML(rr)
		assertStringsNotInHTML(t, body, []string{"Access denied: You are not an admin."})
		assertStringsInHTML(t, body, []string{
			`<div class="card card-compact card-bordered mt-4"><div class="card-body">`,
			`<h2 class="card-title">Users</h2>`,
			`<table class="table table-zebra"><thead><tr><th>Name</th><th>Password</th><th></th></tr></thead> <tbody></tbody></table></div>`,
		})
	})

	t.Run("other users cannot access", func(t *testing.T) {
		rr := sendRequestAsLoggedInOther(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusForbidden)
		assertStringsInHTML(t, getBodyHTML(rr), []string{"Access denied: You are not an admin."})
	})
}

func TestHandlers_Admin_AddUser(t *testing.T) {
	srv := newServerTest()
	originalRepo := srv.Repository

	const uri = "/admin/users"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("other users cannot access", func(t *testing.T) {
		rr := sendRequestAsLoggedInOther(srv, http.MethodPost, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusForbidden)
		assertStringsInHTML(t, getBodyHTML(rr), []string{"Access denied: You are not an admin."})
	})

	t.Run("demo cannot add users", func(t *testing.T) {
		srv.Repository = &mockRepository{
			UsersRegistered: []models.User{
				{ID: 1, Email: "demo@demo.com"},
			},
		}
		app.Config.Server.IsDemo = true
		defer func() {
			app.Config.Server.IsDemo = false
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("email=chiccken@power.com&password=123"))

		assertStatus(t, rr.Code, http.StatusTeapot)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"OK\",\"background\":\"alert-error\",\"message\":\"\",\"title\":\"Every day is Christmas.\"}"}`)
		if len(srv.Repository.Users()) != 1 {
			t.Fail()
		}
	})

	t.Run("cannot add existing user", func(t *testing.T) {
		srv.Repository = &mockRepository{
			UsersRegistered: []models.User{
				{ID: 1, Email: "demo@demo.com"},
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("email=demo@demo.com&password=123"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		if len(srv.Repository.Users()) != 1 {
			t.Fail()
		}
	})

	t.Run("valid request", func(t *testing.T) {
		srv.Repository = &mockRepository{
			UsersRegistered: []models.User{
				{ID: 1, Email: "demo@demo.com"},
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("email=bob@gmail.com&password=bob123"))

		assertStatus(t, rr.Code, http.StatusCreated)
		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<tr><td>bob@gmail.com</td><td>*****</td><th><button class="btn btn-ghost btn-xs" title="Delete user" hx-delete="/admin/users/bob@gmail.com" hx-target="closest tr" hx-swap="outerHTML" hx-confirm="Are you sure you wish to delete this user?" hx-indicator="#fullscreen-loader"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button></th></tr>`,
			`<tr><td><input type="text" name="email" placeholder="Enter new email" class="input input-sm input-bordered w-full"></td><td><input type="password" name="password" placeholder="Enter new password" class="input input-sm input-bordered w-full"></td><th><button class="btn btn-ghost btn-xs" hx-post="/admin/users" hx-include="[name='email'], [name='password']" hx-target="closest tr" hx-swap="outerHTML" hx-indicator="#fullscreen-loader"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle> <line x1="12" y1="8" x2="12" y2="16"></line> <line x1="8" y1="12" x2="16" y2="12"></line></svg></button></th></tr>`,
		})
		if len(srv.Repository.Users()) != 2 {
			t.Fail()
		}
	})
}

func TestHandlers_Admin_DeleteUser(t *testing.T) {
	srv := newServerTest()
	originalRepo := srv.Repository

	const uri = "/admin/users/%s"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodDelete, fmt.Sprintf(uri, "example@example.com"))
	})

	t.Run("other users cannot access", func(t *testing.T) {
		rr := sendRequestAsLoggedInOther(srv, http.MethodDelete, fmt.Sprintf(uri, "example@example.com"), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusForbidden)
		assertStringsInHTML(t, getBodyHTML(rr), []string{"Access denied: You are not an admin."})
	})

	t.Run("cannot delete accounts when demo", func(t *testing.T) {
		srv.Repository = &mockRepository{
			UsersRegistered: []models.User{
				{ID: 1, Email: "demo@demo.com"},
				{ID: 2, Email: "chicken@power.com"},
			},
		}
		app.Config.Server.IsDemo = true
		defer func() {
			app.Config.Server.IsDemo = false
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, fmt.Sprintf(uri, "chicken@power.com"), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusTeapot)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Who do you think you are, eh?\",\"title\":\"General Error\"}"}`)
		if len(srv.Repository.Users()) != 2 {
			t.Fail()
		}
	})

	t.Run("delete invalid user succeeds", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, fmt.Sprintf(uri, "hello@bye.com"), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
	})

	t.Run("cannot delete admin", func(t *testing.T) {
		srv.Repository = &mockRepository{
			UsersRegistered: []models.User{
				{ID: 1, Email: "admin@admin.com"},
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, fmt.Sprintf(uri, "admin@admin.com"), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Cannot delete admin.\",\"title\":\"General Error\"}"}`)
		if len(srv.Repository.Users()) != 1 {
			t.Fail()
		}
	})

	t.Run("can delete other user when admin", func(t *testing.T) {
		srv.Repository = &mockRepository{
			UsersRegistered: []models.User{
				{ID: 1, Email: "admin@admin.com"},
				{ID: 2, Email: "yay@nay.com"},
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, fmt.Sprintf(uri, "yay@nay.com"), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		if len(srv.Repository.Users()) != 1 {
			t.Fail()
		}
	})
}
