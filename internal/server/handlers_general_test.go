package server_test

import (
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"strings"
	"testing"
)

func TestHandlers_General_AvatarDropdown(t *testing.T) {
	srv := newServerTest()

	uri := "/avatar-dropdown"

	t.Run("anonymous access forbidden", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
	})

	t.Run("logged in", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		want := []string{
			`<a href="/settings" class="flex">`, `<div class="pl-1 align-bottom">Settings</div>`,
			`<a href="#" class="flex">`, `<span class="pl-1 align-bottom">About</span>`,
			`<a hx-post="/auth/logout" class="flex" href="#">`, `<span class="pl-1 align-bottom">Log out</span>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_General_Index(t *testing.T) {
	srv := newServerTest()
	srv.Repository = &mockRepository{
		UsersRegistered: []models.User{
			{ID: 1, Email: "test@example.com"},
		},
	}

	const uri = "/"

	t.Run("anonymous access", func(t *testing.T) {
		rr := sendRequest(srv, http.MethodGet, uri, noHeader, nil)

		got := getBodyHTML(rr)
		want := []string{
			`<title hx-swap-oob="true">Home | Recipya</title>`,
			`<h1 class="mb-4 text-2xl font-bold leading-tight text-white md:text-5xl"> A powerful recipe manager that will blow your kitchen away </h1>`,
		}
		assertStringsInHTML(t, got, want)
		notWant := []string{
			`<span id="user-initials">A</span>`,
			"Add recipe"}
		assertStringsNotInHTML(t, got, notWant)
	})

	t.Run("logged in basic access", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)
		got := getBodyHTML(rr)

		want := []string{
			`<title hx-swap-oob="true">Home | Recipya</title>`,
			`<span id="user-initials">A</span>`,
			"Add recipe",
		}
		assertStringsInHTML(t, got, want)
		notWant := []string{
			`<h1 class="mb-4 text-2xl font-bold leading-tight text-white md:text-5xl"> A powerful recipe manager that will blow your kitchen away </h1>`,
			`<a href="/auth/login" class="mr-4 rounded-lg px-2 py-1 text-white hover:bg-green-600">Log In</a>`,
			`<section class="flex h-screen w-full items-center justify-center bg-indigo-100">`,
		}
		assertStringsNotInHTML(t, got, notWant)
	})
}

func TestHandlers_General_NotFound(t *testing.T) {
	srv := newServerTest()

	rr := sendRequestAsLoggedIn(srv, http.MethodGet, "/i-dont-exist-haha", noHeader, nil)

	assertStatus(t, rr.Code, http.StatusNotFound)
	want := []string{
		`<title hx-swap-oob="true">Not Found | Recipya</title>`,
		"Page Not Found",
		"The page you requested to view is not found. Please go back to the main page.",
	}
	assertStringsInHTML(t, getBodyHTML(rr), want)
}

func TestHandlers_General_UserInitials(t *testing.T) {
	srv := newServerTest()
	srv.Repository = &mockRepository{
		UsersRegistered: []models.User{
			{ID: 1, Email: "test@example.com"},
		},
	}

	const uri = "/user-initials"

	t.Run("anonymous user doesn't have initials", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
	})

	t.Run("logged in user has initials", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		body := getBodyHTML(rr)
		want := string(strings.ToUpper(srv.Repository.Users()[0].Email)[0])
		if body != want {
			t.Fatalf("got %s but want %s", body, want)
		}
	})
}
