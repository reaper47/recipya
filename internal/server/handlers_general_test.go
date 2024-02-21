package server_test

import (
	"errors"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"strings"
	"testing"
)

func TestHandlers_General_Download(t *testing.T) {
	srv := newServerTest()

	uri := "/download"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri+"/file.zip")
	})

	t.Run("file does not exist", func(t *testing.T) {
		srv.Files = &mockFiles{
			ReadTempFileFunc: func(name string) ([]byte, error) {
				return nil, errors.New("file does not exist")
			},
		}
		defer func() {
			srv.Files = &mockFiles{}
		}()

		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri+"/does-not-exists", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusNotFound)
	})

	t.Run("file exists", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri+"/exists", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "Content-Type", "text/plain; charset=utf-8")
		assertHeader(t, rr, "Content-Disposition", `attachment; filename="exists"`)
		assertHeader(t, rr, "Content-Length", "6")
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

		assertStatus(t, rr.Code, http.StatusSeeOther)
		got := getBodyHTML(rr)
		want := []string{
			`<a href="/guide">See Other</a>`,
		}
		assertStringsInHTML(t, got, want)
		notWant := []string{
			`<span id="user-initials">A</span>`,
			"Add recipe",
			`<li id="recipes-sidebar-recipes" class="p-2 hover:bg-red-600 hover:text-white" hx-get="/recipes" hx-target="#content" hx-push-url="true" hx-swap-oob="true"> Recipes </li>`,
		}
		assertStringsNotInHTML(t, got, notWant)
	})

	t.Run("hide elements on main page when autologin enabled", func(t *testing.T) {
		app.Config.Server.IsAutologin = true
		defer func() {
			app.Config.Server.IsAutologin = false
		}()
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)
		got := getBodyHTML(rr)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<li class="border-2 rounded-b-lg hover:bg-blue-100 dark:border-gray-500 dark:hover:bg-blue-600"><a hx-post="/auth/logout" class="flex" href="#"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 ml-0 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/></svg><span class="pl-1 align-bottom">Log out</span></a></li>`,
		}
		assertStringsNotInHTML(t, got, want)
	})

	t.Run("logged in basic access", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)
		got := getBodyHTML(rr)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Recipes | Recipya</title>`,
			`<div class="dropdown dropdown-end" hx-get="/user-initials" hx-trigger="load" hx-target="#user-initials">`,
			`<div class="bg-neutral text-neutral-content w-10 rounded-full"><span id="user-initials">A</span></div>`,
			`<ul tabindex="0" id="avatar_dropdown" class="menu dropdown-content mt-3 z-10 p-2 shadow bg-base-100 rounded-box before:content-[''] before:absolute before:right-2 before:top-[-9px] before:border-x-[15px] before:border-x-transparent before:border-b-[8px] before:border-b-[#333] dark:before:border-b-[gray]">`,
			`<li onclick="document.activeElement?.blur()"><a href="/admin" hx-get="/admin" hx-target="#content" hx-push-url="true"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M12 21v-8.25M15.75 21v-8.25M8.25 21v-8.25M3 9l9-6 9 6m-1.5 12V10.332A48.36 48.36 0 0 0 12 9.75c-2.551 0-5.056.2-7.5.582V21M3 21h18M12 6.75h.008v.008H12V6.75Z"></path></svg>Admin</a></li>`,
			`<li onclick="document.activeElement?.blur()"><a href="/reports" hx-get="/reports" hx-target="#content" hx-push-url="true"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M3 3v1.5M3 21v-6m0 0 2.77-.693a9 9 0 0 1 6.208.682l.108.054a9 9 0 0 0 6.086.71l3.114-.732a48.524 48.524 0 0 1-.005-10.499l-3.11.732a9 9 0 0 1-6.085-.711l-.108-.054a9 9 0 0 0-6.208-.682L3 4.5M3 15V4.5"></path></svg>Reports</a></li><div class="divider m-0"></div>`,
			`<li onclick="document.activeElement?.blur()"><a href="/settings" hx-get="/settings" hx-target="#content" hx-push-url="true"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>Settings</a></li>`,
			`<li onclick="about_dialog.showModal()"><a><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>About</a></li><li><a hx-post="/auth/logout"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 ml-0 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg>Log out</a></li></ul></div>`,
			`<div id="ws-notification-container" class="z-20 fixed bottom-0 right-0 p-6 cursor-default hidden"><div class="bg-blue-500 text-white px-4 py-2 rounded shadow-md"><p class="font-medium text-center pb-1"></p></div></div>`,
			"Add recipe",
			`<a class="tooltip tooltip-right active" data-tip="Recipes">`,
			`<a class="tooltip tooltip-right" data-tip="Cookbooks"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">`,
			`<button id="add_recipe" class="btn btn-primary btn-sm hover:btn-accent" hx-get="/recipes/add" hx-target="#content" hx-push-url="true">Add recipe</button>`,
			`<button id="add_cookbook" class="btn btn-primary btn-sm hover:btn-accent" hx-post="/cookbooks" hx-prompt="Enter the name of your cookbook" hx-target="#cookbooks-display" hx-swap="beforeend">Add cookbook</button>`,
		}
		assertStringsInHTML(t, got, want)
		notWant := []string{`A powerful recipe manager that will blow your kitchen away`, `href="/auth/login"`}
		assertStringsNotInHTML(t, got, notWant)
	})
}

func TestHandlers_General_NotFound(t *testing.T) {
	srv := newServerTest()

	rr := sendRequestAsLoggedIn(srv, http.MethodGet, "/i-dont-exist-haha", noHeader, nil)

	assertStatus(t, rr.Code, http.StatusNotFound)
	want := []string{
		`<title hx-swap-oob="true">Page Not Found | Recipya</title>`,
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
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
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
