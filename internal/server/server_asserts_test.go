package server_test

import (
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func assertCookbooksViewMode(t testing.TB, mode models.ViewMode, got string) {
	t.Helper()

	var (
		want    []string
		notWant []string
	)

	if mode == models.GridViewMode {
		want = []string{
			`<title hx-swap-oob="true">Cookbooks | Recipya</title>`,
			`<li id="recipes-sidebar-recipes" class="recipes-sidebar-not-selected" hx-get="/recipes" hx-target="#content" hx-push-url="true" hx-swap-oob="true"><img src="/static/img/cherries.svg" alt=""><span class="hidden md:block ml-1">Recipes</span></li>`,
			`<li id="recipes-sidebar-cookbooks" class="recipes-sidebar-selected" hx-get="/cookbooks" hx-target="#content" hx-push-url="true" hx-swap-oob="true"><img src="/static/img/cookbook.svg" alt=""><span class="hidden md:block ml-1">Cookbooks</span></li>`,
			`<div class="grid grid-flow-col place-content-end p-1">`,
			`<div class="p-2 bg-blue-600 text-white" title="Display as grid"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="20" height="20" fill="currentColor"><path d="M1 2.5A1.5 1.5 0 0 1 2.5 1h3A1.5 1.5 0 0 1 7 2.5v3A1.5 1.5 0 0 1 5.5 7h-3A1.5 1.5 0 0 1 1 5.5v-3zM2.5 2a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5h-3zm6.5.5A1.5 1.5 0 0 1 10.5 1h3A1.5 1.5 0 0 1 15 2.5v3A1.5 1.5 0 0 1 13.5 7h-3A1.5 1.5 0 0 1 9 5.5v-3zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5h-3zM1 10.5A1.5 1.5 0 0 1 2.5 9h3A1.5 1.5 0 0 1 7 10.5v3A1.5 1.5 0 0 1 5.5 15h-3A1.5 1.5 0 0 1 1 13.5v-3zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5h-3zm6.5.5A1.5 1.5 0 0 1 10.5 9h3a1.5 1.5 0 0 1 1.5 1.5v3a1.5 1.5 0 0 1-1.5 1.5h-3A1.5 1.5 0 0 1 9 13.5v-3zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5h-3z"/></svg></div>`,
			`<div class="p-2 hover:bg-red-600 hover:text-white" title="Display as list" hx-get="/cookbooks?view=list" hx-target="#content"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="20" height="20" fill="currentColor"><path fill-rule="evenodd" d="M5 11.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm-3 1a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 4a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 4a1 1 0 1 0 0-2 1 1 0 0 0 0 2z"/></svg></div>`,
			`<div id="cookbooks-container" class="grid justify-center"><div class="grid grid-cols-5 gap-2">`,
			`<button class="w-full border-2 border-gray-800 rounded-lg center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/cookbooks/1" hx-target="#content" hx-push-url="true"> Open </button>`,
			`<button class="w-full border-2 border-gray-800 rounded-lg center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/cookbooks/2" hx-target="#content" hx-push-url="true"> Open </button>`,
			`<button class="w-full border-2 border-gray-800 rounded-lg center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/cookbooks/3" hx-target="#content" hx-push-url="true"> Open </button>`,
		}

		notWant = []string{
			`<div class="grid place-content-center text-sm h-full text-center md:text-base">`,
			`<p>Your cookbooks collection looks a bit empty at the moment.</p>`,
			`<p>Why not start by creating a cookbook by clicking the <a class="underline font-semibold cursor-pointer" hx-post="/cookbooks" hx-prompt="Enter the name of your cookbook" hx-target="#content" hx-push-url="true">Add cookbook</a> button at the top? </p>`,
		}
	} else {
		want = []string{
			`<title hx-swap-oob="true">Cookbooks | Recipya</title>`,
			`<li id="recipes-sidebar-recipes" class="recipes-sidebar-not-selected" hx-get="/recipes" hx-target="#content" hx-push-url="true" hx-swap-oob="true"><img src="/static/img/cherries.svg" alt=""><span class="hidden md:block ml-1">Recipes</span></li>`,
			`<li id="recipes-sidebar-cookbooks" class="recipes-sidebar-selected" hx-get="/cookbooks" hx-target="#content" hx-push-url="true" hx-swap-oob="true"><img src="/static/img/cookbook.svg" alt=""><span class="hidden md:block ml-1">Cookbooks</span></li>`,
			`<div class="grid grid-flow-col place-content-end p-1">`,
			`<div class="p-2 hover:bg-red-600 hover:text-white" title="Display as grid" hx-get="/cookbooks?view=grid" hx-target="#content"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="20" height="20" fill="currentColor"><path d="M1 2.5A1.5 1.5 0 0 1 2.5 1h3A1.5 1.5 0 0 1 7 2.5v3A1.5 1.5 0 0 1 5.5 7h-3A1.5 1.5 0 0 1 1 5.5v-3zM2.5 2a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5h-3zm6.5.5A1.5 1.5 0 0 1 10.5 1h3A1.5 1.5 0 0 1 15 2.5v3A1.5 1.5 0 0 1 13.5 7h-3A1.5 1.5 0 0 1 9 5.5v-3zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5h-3zM1 10.5A1.5 1.5 0 0 1 2.5 9h3A1.5 1.5 0 0 1 7 10.5v3A1.5 1.5 0 0 1 5.5 15h-3A1.5 1.5 0 0 1 1 13.5v-3zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5h-3zm6.5.5A1.5 1.5 0 0 1 10.5 9h3a1.5 1.5 0 0 1 1.5 1.5v3a1.5 1.5 0 0 1-1.5 1.5h-3A1.5 1.5 0 0 1 9 13.5v-3zm1.5-.5a.5.5 0 0 0-.5.5v3a.5.5 0 0 0 .5.5h3a.5.5 0 0 0 .5-.5v-3a.5.5 0 0 0-.5-.5h-3z"/></svg></div>`,
			`<div class="p-2 bg-blue-600 text-white" title="Display as list"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="20" height="20" fill="currentColor"><path fill-rule="evenodd" d="M5 11.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm-3 1a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 4a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 4a1 1 0 1 0 0-2 1 1 0 0 0 0 2z"/></svg></div>`,
			`<div id="cookbooks-container" class="grid justify-center">`,
			`<button class="w-full border-t center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/cookbooks/1" hx-target="#content" hx-push-url="true"> Open </button>`,
			`<button class="w-full border-t center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/cookbooks/2" hx-target="#content" hx-push-url="true"> Open </button>`,
			`<button class="w-full border-t center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/cookbooks/3" hx-target="#content" hx-push-url="true"> Open </button>`,
		}
		notWant = []string{
			`<div class="grid place-content-center text-sm h-full text-center md:text-base">`,
			`<p>Your cookbooks collection looks a bit empty at the moment.</p>`,
			`<p>Why not start by creating a cookbook by clicking the <a class="underline font-semibold cursor-pointer" hx-post="/cookbooks" hx-prompt="Enter the name of your cookbook" hx-target="#content" hx-push-url="true">Add cookbook</a> button at the top? </p>`,
		}
	}

	assertStringsInHTML(t, got, want)
	assertStringsNotInHTML(t, got, notWant)
}

func assertHeader(t testing.TB, rr *httptest.ResponseRecorder, key, value string) {
	t.Helper()
	got := rr.Result().Header.Get(key)
	if got != value {
		t.Fatalf("expected header %s to be %s but got %s", key, value, got)
	}
}

func assertMustBeLoggedIn(t testing.TB, srv *server.Server, method string, uri string) {
	t.Helper()
	rr := sendRequest(srv, method, uri, noHeader, nil)

	assertStatus(t, rr.Code, http.StatusSeeOther)
	assertHeader(t, rr, "Location", "/auth/login")
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("expected status %d but got %d", want, got)
	}
}

func assertStringsInHTML(t testing.TB, bodyHTML string, wants []string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(bodyHTML, want) {
			t.Fatalf("string not found in HTML:\n%s", want)
		}
	}
}

func assertStringsNotInHTML(t testing.TB, bodyHTML string, wants []string) {
	t.Helper()
	for _, want := range wants {
		if strings.Contains(bodyHTML, want) {
			t.Fatalf("string found in HTML when it should not:\n%s", want)
		}
	}
}

func assertUserSettings(t testing.TB, id int64, got, want *models.UserSettings) {
	t.Helper()
	if got.ConvertAutomatically != want.ConvertAutomatically {
		t.Fatalf("settings ConvertAutomatically got %t but want %t", got.ConvertAutomatically, want.ConvertAutomatically)
	}

	if got.CookbooksViewMode != want.CookbooksViewMode {
		t.Fatalf("settings CookbooksViewMode got %d but want %d", got.CookbooksViewMode, want.CookbooksViewMode)
	}

	if got.MeasurementSystem != want.MeasurementSystem {
		t.Fatalf("settings MeasurementSystem got %q but want %q", got.MeasurementSystem, want.MeasurementSystem)
	}
}
