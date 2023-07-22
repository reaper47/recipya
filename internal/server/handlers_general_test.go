package server_test

import (
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"strings"
	"testing"
)

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
			`<div hx-get="/user-initials" hx-trigger="load" hx-target="#user-initials"><button id="avatar-button" class="items-center h-10 px-4 mr-1 text-center text-gray-800 bg-blue-200 border rounded-full" _="on click toggle .hidden on next <div/>"><span id="user-initials">A</span></button><div id="avatar-dropdown-container" class="absolute w-24 right-3 top-[3.6rem] bg-white rounded-lg hidden"><ul id="avatar-menu"><li class="hover:bg-blue-100 border-2 rounded-t-lg"><a class="flex" href="/settings" hx-get="/settings" hx-target="#content" hx-push-url="/settings"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg><div class="pl-1 align-bottom">Settings</div></a></li><li class="hover:bg-blue-100 border-2"><a href="#" class="flex"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg><span class="pl-1 align-bottom">About</span></a></li><li class="hover:bg-blue-100 border-2 rounded-b-lg"><a hx-post="/auth/logout" class="flex" href="#"><svg xmlns="http://www.w3.org/2000/svg" style="margin-left: 0;" class="w-5 h-5 ml-0 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/></svg><span class="pl-1 align-bottom">Log out</span></a></li></ul></div>`,
			`<script defer> window.addEventListener("DOMContentLoaded", () => { const recipesPattern = new RegExp("^/recipes/\\d+$"); const pathsHideAddRecipeButton = [ "/settings", "/recipes/add", "/recipes/add/manual", "/recipes/" ]; const toggleAddRecipeButton = () => { if (pathsHideAddRecipeButton.some(path => path === location.pathname || recipesPattern.test(location.pathname))) { document.querySelector("#add-recipe")?.classList.add("hidden"); } else { document.querySelector("#add-recipe")?.classList.remove("hidden"); } } toggleAddRecipeButton(); document.addEventListener("htmx:pushedIntoHistory", toggleAddRecipeButton); }); </script>`,
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

func TestHandlers_General_Settings(t *testing.T) {
	srv := newServerTest()

	uri := "/settings"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
	})

	t.Run("tab profile", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"/tabs/profile", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<p class="grid justify-end font-semibold">Change password:</p>`,
			`<form class="h-fit w-fit border p-2" hx-post="/auth/change-password" enctype="multipart/form-data" hx-indicator="#fullscreen-loader" hx-swap="none">`,
			`<input class="w-full rounded-lg bg-gray-100 px-4 py-2" id="password-current" name="password-current" placeholder="Current password..." required type="password"/>`,
			`<input class="w-full rounded-lg bg-gray-100 px-4 py-2" id="password-new" name="password-new" placeholder="New password..." required type="password"/>`,
			`<input class="w-full rounded-lg bg-gray-100 px-4 py-2" id="password-confirm" name="password-confirm" placeholder="Retype new password..." required type="password"/>`,
			`<button type="submit" class="w-full p-2 font-semibold text-white duration-300 bg-blue-500 rounded-lg hover:bg-blue-800" > Update </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("tab recipes", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"/tabs/recipes", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<p class="font-semibold grid justify-end">Export data:</p>`,
			`<button type="button" hx-get="/recipes/export" hx-swap="none" class="h-fit w-fit bg-white border border-gray-300 focus:outline-none hover:bg-gray-100 focus:ring-4 focus:ring-gray-200 rounded-lg px-5 py-2.5 mr-2 mb-2 dark:bg-gray-800 dark:text-white dark:border-gray-600 dark:hover:bg-gray-700 dark:hover:border-gray-600 dark:focus:ring-gray-700">`,
			`<svg xmlns="http://www.w3.org/2000/svg" style="margin-left: 0;" class="w-5 h-5 ml-0 self-center" fill="black" viewBox="0 0 24 24" stroke="currentColor"><path d="M16 11v5H2v-5H0v5a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-5z"/><path d="m9 14 5-6h-4V0H8v8H4z"/></svg>`,
			`<label for="systems" class="grid justify-end mb-2 font-semibold dark:text-white">Measurement System:</label>`,
			`<select id="systems" name="system" hx-post="/settings/measurement-system" hx-swap="none" class="h-fit w-fit bg-gray-50 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"><option value="imperial" selected>Imperial</option><option value="metric">Metric</option></select>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
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
