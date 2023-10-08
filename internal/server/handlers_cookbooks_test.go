package server_test

import (
	"net/http"
	"testing"
)

func TestHandlers_Cookbooks(t *testing.T) {
	srv := newServerTest()

	uri := "/cookbooks"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("no cookbooks", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Cookbooks | Recipya</title>`,
			`<li id="recipes-sidebar-recipes" class="recipes-sidebar-not-selected" hx-get="/recipes" hx-target="#content" hx-push-url="true" hx-swap-oob="true"><img src="/static/img/cherries.svg" alt=""><span class="hidden md:block ml-1">Recipes</span></li>`,
			`<li id="recipes-sidebar-cookbooks" class="recipes-sidebar-selected" hx-get="/cookbooks" hx-target="#content" hx-push-url="true" hx-swap-oob="true"><img src="/static/img/cookbook.svg" alt=""><span class="hidden md:block ml-1">Cookbooks</span></li>`,
			`<div class="grid place-content-center text-sm h-full text-center md:text-base"><p>Your cookbooks collection looks a bit empty at the moment.</p><p> Why not start by creating a cookbook by clicking the <a class="underline font-semibold cursor-pointer" hx-get="/cookbooks/new" hx-target="#content" hx-push-url="true">Add cookbook</a> button at the top? </p></div>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}
