package server_test

import (
	"errors"
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

		got := getBodyHTML(rr)
		want := []string{
			`<title hx-swap-oob="true">Home | Recipya</title>`,
			`<h1 class="mb-4 text-2xl font-bold leading-tight text-white md:text-5xl"> A powerful recipe manager that will blow your kitchen away </h1>`,
		}
		assertStringsInHTML(t, got, want)
		notWant := []string{
			`<span id="user-initials">A</span>`,
			"Add recipe",
			`<li id="recipes-sidebar-recipes" class="p-2 hover:bg-red-600 hover:text-white" hx-get="/recipes" hx-target="#content" hx-push-url="true" hx-swap-oob="true"> Recipes </li>`,
		}
		assertStringsNotInHTML(t, got, notWant)
	})

	t.Run("logged in basic access", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)
		got := getBodyHTML(rr)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Home | Recipya</title>`,
			`<div hx-get="/user-initials" hx-trigger="load" hx-target="#user-initials"><button id="avatar-button" class="items-center h-10 px-4 mr-1 text-center text-gray-800 bg-blue-200 border rounded-full" _="on click toggle .hidden on next <div/>"><span id="user-initials">A</span></button><div id="avatar-dropdown-container" class="absolute w-24 right-3 top-[3.6rem] bg-white rounded-lg hidden z-10 dark:bg-gray-900"><ul id="avatar-menu" class="before:content-[''] before:absolute before:right-2 before:top-[-9px] before:border-x-[10px] before:border-x-transparent before:border-b-[10px] before:border-b-[#333] dark:before:border-b-[gray]"><li class="border-2 rounded-t-lg hover:bg-blue-100 dark:border-gray-600 dark:hover:bg-blue-600"><a class="flex" href="/settings" hx-get="/settings" hx-target="#content" hx-push-url="true"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg><div class="pl-1 align-bottom">Settings</div></a></li><li class="border-2 hover:bg-blue-100 dark:border-gray-500 dark:hover:bg-blue-600"><a href="#" class="flex"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg><span class="pl-1 align-bottom">About</span></a></li><li class="border-2 rounded-b-lg hover:bg-blue-100 dark:border-gray-500 dark:hover:bg-blue-600"><a hx-post="/auth/logout" class="flex" href="#"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 ml-0 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/></svg><span class="pl-1 align-bottom">Log out</span></a></li></ul></div>`,
			`<script defer> var timeoutToast = timeoutToast || null; htmx.on('showToast', function (event) { const toastContainer = document.getElementById('toast-container'); const {message, backgroundColor} = JSON.parse(event.detail.value); const toast = document.createElement('div'); toast.classList.add('text-white', 'px-6', 'py-4', 'rounded', 'shadow-md', backgroundColor); toast.textContent = message; if (toastContainer.firstChild) { toastContainer.replaceChild(toast, toastContainer.firstChild); } else { toastContainer.appendChild(toast); } toastContainer.classList.remove('hidden'); timeoutToast = setTimeout(function () { toastContainer.classList.add('hidden'); toastContainer?.removeChild(toast); }, 5000); }); </script>`,
			"Add recipe",
			`<li id="recipes-sidebar-recipes" class="recipes-sidebar-selected cursor-pointer" hx-get="/recipes" hx-target="#content" hx-push-url="true" hx-swap-oob="true" _="on load if location.pathname is not '/recipes' then remove .recipes-sidebar-selected then add .recipes-sidebar-not-selected"><img src="/static/img/cherries.svg" alt=""><span class="hidden md:block ml-1">Recipes</span></li>`,
			`<li id="recipes-sidebar-cookbooks" class="recipes-sidebar-not-selected cursor-pointer" hx-get="/cookbooks" hx-target="#content" hx-push-url="true" hx-swap-oob="true" _="on load if location.pathname is '/cookbooks' then remove .recipes-sidebar-not-selected add .recipes-sidebar-selected"><img src="/static/img/cookbook.svg" alt=""><span class="hidden md:block ml-1">Cookbooks</span></li>`,
			`<button id="add-recipe" class="app-bar-center-button" hx-get="/recipes/add" hx-target="#content" hx-push-url="true"> Add recipe </button>`,
			`<button id="add-cookbook" class="app-bar-center-button" hx-post="/cookbooks" hx-prompt="Enter the name of your cookbook" hx-target=".cookbooks-display" hx-swap="beforeend"> Add cookbook </button>`,
		}
		assertStringsInHTML(t, got, want)
		notWant := []string{
			`<h1 class="mb-4 text-2xl font-bold leading-tight text-white md:text-5xl"> A powerful recipe manager that will blow your kitchen away </h1>`,
			`<a href="/auth/login" class="mr-4 rounded-lg px-2 py-1 text-white hover:bg-green-600">Log In</a>`,
			`<section class="flex h-screen w-full items-center justify-center bg-indigo-100 dark:bg-gray-800">`,
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
