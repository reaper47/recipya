package server_test

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestHandlers_Cookbooks(t *testing.T) {
	srv := newServerTest()
	originalRepo := &mockRepository{
		CookbooksRegistered: map[int64][]models.Cookbook{1: {}},
	}
	srv.Repository = originalRepo

	uri := "/cookbooks"

	prepare := func(srv *server.Server, viewMode models.ViewMode) ([]models.Cookbook, *mockRepository, func()) {
		_, repo, revertFunc := prepareCookbook(srv)

		repo.UserSettingsRegistered = map[int64]*models.UserSettings{
			1: {CookbooksViewMode: viewMode},
		}

		originalCookbooks := make([]models.Cookbook, len(repo.CookbooksRegistered[1]))
		copy(originalCookbooks, repo.CookbooksRegistered[1])

		return originalCookbooks, repo, revertFunc
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("no cookbooks", func(t *testing.T) {
		srv.Repository = &mockRepository{
			CookbooksRegistered:    map[int64][]models.Cookbook{1: {}},
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Cookbooks | Recipya</title>`,
			`<li id="recipes-sidebar-recipes" class="recipes-sidebar-not-selected" hx-get="/recipes" hx-target="#content" hx-push-url="true" hx-swap-oob="true"><img src="/static/img/cherries.svg" alt=""><span class="hidden md:block ml-1">Recipes</span></li>`,
			`<li id="recipes-sidebar-cookbooks" class="recipes-sidebar-selected" hx-get="/cookbooks" hx-target="#content" hx-push-url="true" hx-swap-oob="true"><img src="/static/img/cookbook.svg" alt=""><span class="hidden md:block ml-1">Cookbooks</span></li>`,
			`<div class="grid place-content-center text-sm h-full text-center md:text-base">`,
			`<p>Your cookbooks collection looks a bit empty at the moment.</p>`,
			`<p>Why not start by creating a cookbook by clicking the <a class="underline font-semibold cursor-pointer" hx-post="/cookbooks" hx-prompt="Enter the name of your cookbook" hx-target="#content">Add cookbook</a> button at the top? </p>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("get cookbooks error", func(t *testing.T) {
		srv.Repository = &mockRepository{
			CookbooksFunc: func(userID int64) ([]models.Cookbook, error) {
				return nil, errors.New("error fetching cookbooks")
			},
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Error getting cookbooks.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("have cookbooks grid preferred mode", func(t *testing.T) {
		_, repo, revertFunc := prepare(srv, models.GridViewMode)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, 1, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.GridViewMode})
		assertCookbooksViewMode(t, models.GridViewMode, getBodyHTML(rr))
	})

	t.Run("have cookbooks list preferred mode", func(t *testing.T) {
		_, repo, revertFunc := prepare(srv, models.ListViewMode)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, 1, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.ListViewMode})
		body := getBodyHTML(rr)
		assertCookbooksViewMode(t, models.ListViewMode, body)
		want := []string{
			`<li id="cookbook-1" class="cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700">`,
			`<p class="font-semibold"> Lovely Canada </p>`,
			`<p class="font-semibold"> Lovely America </p>`,
			`<p class="font-semibold"> Lovely Ukraine </p>`,
			`<form id="cookbook-image-form-1" enctype="multipart/form-data" hx-put="/cookbooks/1/image" hx-trigger="change from:#cookbook-image-1" hx-swap="none"><input id="cookbook-image-1" type="file" accept="image/*" name="image" required class="hidden" _="on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from next <button/>"></form>`,
			`<form id="cookbook-image-form-2" enctype="multipart/form-data" hx-put="/cookbooks/2/image"`,
			`<form id="cookbook-image-form-3" enctype="multipart/form-data" hx-put="/cookbooks/3/image"`,
			`<div class="three-dots-container cursor-pointer h-fit justify-self-end hover:text-red-600" _="on click if menuOpen add .hidden to menuOpen end then set $menuOpen to #cookbook-menu-container`,
			`<a id="cookbook-menu-share" class="flex p-1" hx-post="/cookbooks/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call #share-dialog.showModal()">`,
			`<a id="cookbook-menu-download" class="flex p-1" hx-get="/cookbooks/1/download">`,
			`<a id="cookbook-menu-delete" class="flex p-1" hx-delete="/cookbooks/1" hx-swap="outerHTML" hx-target="closest .cookbook" hx-confirm="Are you sure you want to delete 'Cookbooks'? Its recipes will not be deleted."><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg><div class="pl-1 align-bottom">Delete</div></a>`,
			`<span class="w-fit h-fit text-xs text-center font-medium select-none py-1 px-2 bg-indigo-700 text-white self-end justify-self-end"> 0 </span>`,
			`<button class="w-full border-t center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/cookbooks/2?page=1" hx-target="#content" hx-push-url="/cookbooks/2"> Open </button>`,
			`<footer id="pagination" class="grid justify-center pb-2 " ><div class="grid p-2 border-t border-gray-200 sm:px-6 dark:border-t-gray-800"><div class="text-sm font-light sm:hidden"><button disabled class="pagination-selected-edge-button-mobile"> Previous </button><button disabled class="pagination-selected-edge-button-mobile"> Next </button></div><div class="hidden select-none sm:flex"><div class="col-span-8"><nav class="inline-flex -space-x-px rounded shadow-sm" aria-label="Pagination"><button disabled class="pagination-square p-2 cursor-not-allowed rounded-l"><span class="sr-only">Previous</span><svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true"><path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd"/></svg></button><button aria-current="page" class="pagination-selected cursor-default bg-indigo-50"> 1 </button><button disabled class="pagination-square p-2 cursor-not-allowed rounded-r disabled"><span class="sr-only">Next</span><svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true"><path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"/></svg></button></nav></div></div></div><div class="text-center"><p class="text-sm"> Showing <span class="font-medium">1</span> to <span class="font-medium">3</span> of <span class="font-medium">3</span> results </p></div></footer>`,
			`<script defer> document.addEventListener('click', (event) => { const cookbookContainers = document.querySelectorAll(".cookbook-menu"); cookbookContainers.forEach(c => { if (c && !c.classList.contains("hidden") && !event.target.classList.contains("three-dots-container") && !["svg", "path"].includes(event.target.tagName)) { c.classList.add("hidden"); } }); }); </script>`,
		}
		assertStringsInHTML(t, body, want)
	})

	t.Run("have cookbooks grid view select list", func(t *testing.T) {
		_, repo, revertFunc := prepare(srv, models.GridViewMode)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?view=list", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, 1, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.ListViewMode})
		assertCookbooksViewMode(t, models.ListViewMode, getBodyHTML(rr))
	})

	t.Run("have cookbooks list view select grid", func(t *testing.T) {
		_, repo, revertFunc := prepare(srv, models.ListViewMode)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?view=grid", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, 1, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.GridViewMode})
		assertCookbooksViewMode(t, models.GridViewMode, getBodyHTML(rr))
	})

	t.Run("title must not be empty", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader(""))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Title must not be empty.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("create cookbook", func(t *testing.T) {
		repo := originalRepo
		repo.CookbooksRegistered[1] = make([]models.Cookbook, 0)
		repo.UserSettingsRegistered = map[int64]*models.UserSettings{1: {CookbooksViewMode: models.GridViewMode}}
		defer func() {
			srv.Repository = originalRepo
		}()
		title := "Lovely America"

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader(title))

		assertStatus(t, rr.Code, http.StatusOK)
		cookbooks, ok := repo.CookbooksRegistered[1]
		if !ok {
			t.Fatal("user is not registered in the cookbooks map")
		}
		isFound := slices.ContainsFunc(cookbooks, func(cookbook models.Cookbook) bool {
			return cookbook.Title == title
		})
		if !isFound {
			t.Fatal("cookbook must have been added to the user's collection")
		}
		want := []string{
			`<section id="cookbook-1" class="cookbook relative col-span-1 bg-white rounded-lg shadow-md dark:bg-neutral-700">`,
			`<img class="rounded-t-lg w-full border-b dark:border-b-gray-800 h-32 md:h-48 text-center hover:bg-gray-100 dark:hover:bg-blue-100 hover:opacity-80" src="/static/img/cookbooks-new/placeholder.png" alt="Cookbook image"><form id="cookbook-image-form-0" enctype="multipart/form-data" hx-put="/cookbooks/0/image" hx-trigger="change from:#cookbook-image-0" hx-swap="none">`,
			`<form id="cookbook-image-form-0" enctype="multipart/form-data" hx-put="/cookbooks/0/image" hx-trigger="change from:#cookbook-image-0" hx-swap="none"><input id="cookbook-image-0" type="file" accept="image/*" name="image" required class="hidden" _="on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from next <button/>"></form>`,
			`<div class="relative" onclick="document.querySelector('#cookbook-image-0').click()">`,
			`<p class="font-semibold">Lovely America</p>`,
			`<span class="grid justify-self-end text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg px-2 h-fit"> 0 </span>`,
			`<button class="w-full border-2 border-gray-800 rounded-lg center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/cookbooks/1?page=1" hx-target="#content" hx-push-url="/cookbooks/1"> Open </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Cookbooks_Cookbook(t *testing.T) {
	srv := newServerTest()

	uri := func(id int) string {
		return fmt.Sprintf("/cookbooks/%d", id)
	}

	assertCookbooksEqual := func(t *testing.T, originalCookbooks, cookbooks []models.Cookbook) {
		isCookbooksEqual := slices.EqualFunc(originalCookbooks, cookbooks, func(c1 models.Cookbook, c2 models.Cookbook) bool {
			return c1.ID == c2.ID
		})
		if !isCookbooksEqual {
			t.Fatal("did not expect a cookbook to be deleted")
		}
	}

	prepare := func(srv *server.Server) ([]models.Cookbook, *mockRepository, func()) {
		_, repo, revertFunc := prepareCookbook(srv)
		originalCookbooks := make([]models.Cookbook, len(repo.CookbooksRegistered[1]))
		copy(originalCookbooks, repo.CookbooksRegistered[1])
		return originalCookbooks, repo, revertFunc
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodDelete, uri(1))
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri(1))
	})

	t.Run("cannot delete cookbooks from other user", func(t *testing.T) {
		originalCookbooks, repo, revertFunc := prepare(srv)
		repo.CookbooksRegistered[2] = []models.Cookbook{{ID: 1}}
		defer revertFunc()

		rr := sendHxRequestAsLoggedInOther(srv, http.MethodDelete, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertCookbooksEqual(t, originalCookbooks, repo.CookbooksRegistered[1])
	})

	t.Run("error deleting cookbook", func(t *testing.T) {
		originalCookbooks, repo, revertFunc := prepare(srv)
		defer revertFunc()
		repo.DeleteCookbookFunc = func(_, _ int64) error {
			return errors.New("error deleting")
		}

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Error deleting cookbook.\",\"backgroundColor\":\"bg-red-500\"}"}`)
		assertCookbooksEqual(t, originalCookbooks, repo.CookbooksRegistered[1])
	})

	testcases := []struct{ name string }{
		{"success even when cookbook does not exist"},
		{"success when cookbook exists"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			originalCookbooks, repo, revertFunc := prepare(srv)
			defer revertFunc()

			rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri(4), noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			assertCookbooksEqual(t, originalCookbooks, repo.CookbooksRegistered[1])
		})
	}

	t.Run("deleting a cookbook does not delete recipes", func(t *testing.T) {
		_, repo, revertFunc := prepare(srv)
		defer revertFunc()
		originalRecipes := repo.CookbooksRegistered[1][0].Recipes
		repo.RecipesRegistered[1] = make(models.Recipes, len(originalRecipes))
		repo.RecipesRegistered[1] = originalRecipes

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri(4), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		if len(originalRecipes) != len(repo.CookbooksRegistered[1][0].Recipes) {
			t.Fatal("no recipe should have been deleted")
		}
	})

	t.Run("open cookbook page must be specified hx request", func(t *testing.T) {
		_, _, revertFunc := prepare(srv)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Missing page query parameter.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("open cookbook page must be specified normal request", func(t *testing.T) {
		_, _, revertFunc := prepare(srv)
		defer revertFunc()

		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri(20), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "Location", "/cookbooks")
	})

	t.Run("cookbook does not exist hx request", func(t *testing.T) {
		_, _, revertFunc := prepare(srv)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(20)+"?page=1", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusNotFound)
		want := []string{"404 page not found"}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("cookbook does not belong to the user", func(t *testing.T) {
		_, _, revertFunc := prepare(srv)
		defer revertFunc()

		rr := sendHxRequestAsLoggedInOther(srv, http.MethodGet, uri(1)+"?page=1", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusNotFound)
		want := []string{"404 page not found"}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("cookbook belongs to user and is empty", func(t *testing.T) {
		cookbooks, repo, revertFunc := prepare(srv)
		defer revertFunc()
		id := len(cookbooks) + 1
		repo.CookbooksRegistered[1] = append(cookbooks, models.Cookbook{ID: int64(id), Title: "Ensiferum"})
		srv.Repository = repo

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(id)+"?page=1", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Ensiferum | Recipya</title>`,
			`<div id="content-title" hx-swap-oob="innerHTML">Ensiferum</div>`,
			`<section class="grid justify-center p-4"><div class="relative">`,
			`<form class="w-72 md:w-96" hx-post="/cookbooks/recipes/search" hx-target="#search-results" hx-vals='{"id": 4, "page": 1}'><div class="flex"><div class="relative w-full"><label><input type="search" id="search-recipes" name="q" class="block z-20 p-2.5 w-full text-sm bg-gray-50 rounded-r-lg border border-gray-300 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400" placeholder="Search for recipes to add..." _="on keyup if event.target.value !== '' then remove .md:block from #search-shortcut else add .md:block to #search-shortcut end then if event.key === 'Backspace' and event.target.value === '' then send submit to closest <form/> end"></label><kbd id="search-shortcut" class="hidden absolute top-2 right-12 font-sans font-semibold select-none dark:text-slate-500 md:block"><abbr title="Control" class="no-underline text-slate-300 dark:text-slate-500">Ctrl </abbr> / </kbd><button type="submit" class="absolute top-0 right-0 p-2.5 text-sm font-medium h-full text-white bg-blue-700 rounded-r-lg border border-blue-700 hover:bg-blue-800 dark:bg-blue-600 dark:hover:bg-blue-700"><svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"/></svg><span class="sr-only">Search</span></button></div></div></form>`,
			`<p class="grid justify-center font-semibold underline mt-2 md:mt-0 md:text-xl md:hidden"> Ensiferum </p>`,
			`<section id="search-results" class="justify-center grid"><div class="grid place-content-center text-sm text-center md:text-base" style="height: 50vh"><p>Your cookbook looks a bit empty at the moment.</p><p>Why not add recipes to your cookbook by searching for recipes in the search box above?</p></div></section>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	requestTypes := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string, contentType header, body *strings.Reader) *httptest.ResponseRecorder
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedIn},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedIn},
	}
	for _, tc := range requestTypes {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("cookbook has recipes", func(t *testing.T) {
				_, _, revertFunc := prepare(srv)
				defer revertFunc()

				rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(1)+"?page=1", noHeader, nil)

				assertStatus(t, rr.Code, http.StatusOK)
				want := []string{
					` <title hx-swap-oob="true">Lovely Canada | Recipya</title>`,
					`<div id="content-title" hx-swap-oob="innerHTML">Lovely Canada</div>`,
					`<script defer> function initReorder() { document.querySelectorAll("#search-results").forEach(sortable => { const sortableInstance = new Sortable(sortable, { animation: 150, ghostClass: 'blue-background-class', handle: '.handle', onEnd: function (event) { Array.from(document.querySelector('#search-results').children).forEach((c, i) => { const p = c.querySelector('.handle'); p.innerText = i + 1; }); }, }); sortable.addEventListener("htmx:afterSwap", function () { sortableInstance.option("disabled", false); }); }); } document.addEventListener("keydown", (event) => { if (event.ctrlKey && event.key === "/") { event.preventDefault(); document.querySelector("#search-recipes").focus(); } }); loadSortableJS().then(initReorder); </script>`,
					`<section class="grid justify-center p-4"><div class="relative">`,
					`<form class="sortable" hx-put="/cookbooks/1/reorder" hx-trigger="end" hx-swap="none"><input type='hidden' name='cookbook-id' value='1'/><ul id="search-results" class="cookbooks-display grid gap-2 p-2 md:p-0 text-sm md:text-base">`,
					`<li class="recipe cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700"><input type='hidden' name='recipe-id' value='3'/><div class="grid grid-cols-4"><img class="w-fit col-span-1 border-r h-[90px]" src="/data/images/00000000-0000-0000-0000-000000000000.jpg" alt="Recipe image"><div class="grid col-span-3 gap-1 p-2"><div class="grid grid-flow-col"><p class="font-semibold">Gotcha</p><div class="grid justify-end"><div class="pb-4 pl-4 text-xs"><span title="Remove recipe from cookbook" hx-delete="/cookbooks/1/recipes/3" hx-swap="outerHTML" hx-target="closest .recipe" hx-confirm="Are you sure you want to remove this recipe from the cookbook?" hx-indicator="#fullscreen-loader"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg></span></div></div></div><div class="text-xs"><span class="w-fit leading-none flex items-center justify-center p-2 text-blue-700 bg-blue-100 border border-blue-300 rounded-full dark:border-gray-800"> American </span></div></div></div><div class="flex border-t dark:border-gray-800"><button class="w-full center hover:bg-gray-800 hover:text-white hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/recipes/3" hx-target="#content" hx-swap="innerHTML" hx-push-url="true"> View </button><p class="justify-self-end px-4 select-none cursor-move handle"> 1 </p></div></li>`,
				}
				assertStringsInHTML(t, getBodyHTML(rr), want)
			})
		})
	}
}

func TestHandlers_Cookbooks_AddRecipe(t *testing.T) {
	srv := newServerTest()
	originalRepo := srv.Repository

	uri := func(cookbookID int64) string {
		return fmt.Sprintf("/cookbooks/%d", cookbookID)
	}

	revert := func() {
		srv.Repository = originalRepo
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri(1))
	})

	testcases := []struct {
		name      string
		form      string
		wantToast string
	}{
		{
			name:      "cookbookID missing in form",
			form:      "recipeId=1",
			wantToast: "Missing 'cookbookId' in body.",
		},
		{
			name:      "recipeID missing in form",
			form:      "cookbookId=1",
			wantToast: "Missing 'recipeId' in body.",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendRequestAsLoggedIn(srv, http.MethodPost, uri(1), formHeader, strings.NewReader(tc.form))

			assertStatus(t, rr.Code, http.StatusBadRequest)
			assertHeader(t, rr, "HX-Trigger", fmt.Sprintf(`{"showToast":"{\"message\":\"%s\",\"backgroundColor\":\"bg-red-500\"}"}`, tc.wantToast))
		})
	}

	t.Run("valid request", func(t *testing.T) {
		recipes := models.Recipes{
			{ID: 1, Name: "Cheese Toasts"},
			{ID: 2, Name: "Maple Syrup Waffles"},
			{ID: 3, Name: "Chicken Jerky"},
		}
		repo := &mockRepository{
			CookbooksRegistered: map[int64][]models.Cookbook{
				1: {
					{ID: 1, Recipes: []models.Recipe{recipes[0]}},
					{ID: 2, Recipes: []models.Recipe{recipes[1]}},
				},
			},
			RecipesRegistered: map[int64]models.Recipes{1: slices.Clone(recipes)},
		}
		srv.Repository = repo
		defer revert()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri(1), formHeader, strings.NewReader("cookbookId=2&recipeId=2"))

		assertStatus(t, rr.Code, http.StatusCreated)
		want := []models.Cookbook{
			{ID: 1, Recipes: []models.Recipe{recipes[0]}},
			{ID: 2, Recipes: []models.Recipe{recipes[1], recipes[2]}},
		}
		assertCookbooks(t, repo.CookbooksRegistered[1], want)
	})
}

func TestHandlers_Cookbooks_DeleteCookbookRecipe(t *testing.T) {
	srv := newServerTest()
	originalRepo := srv.Repository

	uri := func(cookbookID, recipeID int64) string {
		return fmt.Sprintf("/cookbooks/%d/recipes/%d", cookbookID, recipeID)
	}

	restore := func() {
		srv.Repository = originalRepo
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodDelete, uri(1, 2))
	})

	t.Run("no cookbook matches cookbookID", func(t *testing.T) {
		repo := &mockRepository{CookbooksRegistered: map[int64][]models.Cookbook{
			1: {{ID: 1}},
		}}
		srv.Repository = repo
		defer restore()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri(10, 10), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		if len(repo.CookbooksRegistered[1]) != 1 {
			t.Fatal("expected the cookbooks to remain untouched")
		}
	})

	t.Run("the recipe is not found in the cookbook", func(t *testing.T) {
		repo := &mockRepository{CookbooksRegistered: map[int64][]models.Cookbook{
			1: {{ID: 1}, {ID: 2}},
		}}
		srv.Repository = repo
		defer restore()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri(1, 10), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		if len(repo.CookbooksRegistered[1]) != 2 {
			t.Fatal("expected the cookbooks to remain untouched")
		}
	})

	t.Run("other user cannot delete other user's recipes", func(t *testing.T) {
		cookbooks1 := []models.Cookbook{{ID: 1, Recipes: []models.Recipe{{ID: 1}, {ID: 2}}}}
		cookbooks2 := []models.Cookbook{{ID: 3}, {ID: 4}}
		repo := &mockRepository{CookbooksRegistered: map[int64][]models.Cookbook{
			1: slices.Clone(cookbooks1),
			2: slices.Clone(cookbooks2),
		}}
		srv.Repository = repo
		defer restore()

		rr := sendHxRequestAsLoggedInOther(srv, http.MethodDelete, uri(1, 1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertCookbooks(t, repo.CookbooksRegistered[1], cookbooks1)
		assertCookbooks(t, repo.CookbooksRegistered[2], cookbooks2)
	})

	t.Run("valid request", func(t *testing.T) {
		repo := &mockRepository{CookbooksRegistered: map[int64][]models.Cookbook{
			1: {
				{ID: 1, Recipes: []models.Recipe{{ID: 1}, {ID: 2}}},
				{ID: 2, Recipes: models.Recipes{{ID: 1}}},
			},
		}}
		srv.Repository = repo
		defer restore()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri(1, 1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertCookbooks(t, repo.CookbooksRegistered[1], []models.Cookbook{
			{ID: 1, Recipes: []models.Recipe{{ID: 2}}},
			{ID: 2, Recipes: models.Recipes{{ID: 1}}},
		})
	})

	t.Run("valid request delete all recipes in cookbook", func(t *testing.T) {
		repo := &mockRepository{CookbooksRegistered: map[int64][]models.Cookbook{
			1: {{ID: 1, Recipes: []models.Recipe{{ID: 1}}}},
		}}
		srv.Repository = repo
		defer restore()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri(1, 1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertCookbooks(t, repo.CookbooksRegistered[1], []models.Cookbook{
			{ID: 1, Recipes: []models.Recipe{}},
		})
		want := []string{
			`<div class="grid place-content-center text-sm text-center md:text-base" style="height: 50vh">`,
			`<p>Your cookbook looks a bit empty at the moment.</p>`,
			`<p>Why not add recipes to your cookbook by searching for recipes in the search box above?</p>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Cookbooks_DownloadCookbook(t *testing.T) {
	srv := newServerTest()

	uri := func(id int64) string {
		return fmt.Sprintf("/cookbooks/%d/download", id)
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri(1))
	})

	t.Run("cookbook not exist", func(t *testing.T) {
		_, _, revert := prepareCookbook(srv)
		defer revert()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(10), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Could not fetch cookbook.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("other user cannot download someone else's cookbook", func(t *testing.T) {
		_, _, revert := prepareCookbook(srv)
		files := &mockFiles{}
		srv.Files = files
		defer revert()

		rr := sendHxRequestAsLoggedInOther(srv, http.MethodGet, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Could not fetch cookbook.\",\"backgroundColor\":\"bg-red-500\"}"}`)
		if files.exportHitCount > 0 {
			t.Fatal("export function must not have been called")
		}
	})

	t.Run("valid request", func(t *testing.T) {
		_, _, revert := prepareCookbook(srv)
		files := &mockFiles{}
		srv.Files = files
		defer revert()

		rr := sendHxRequestAsLoggedInOther(srv, http.MethodGet, uri(4), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", "/download/Lovely Canada 2.pdf")
		if files.exportHitCount != 1 {
			t.Fatal("export function must have been called")
		}
	})

}

func TestHandlers_Cookbooks_Image(t *testing.T) {
	srv := newServerTest()

	uri := func(id int) string {
		return fmt.Sprintf("/cookbooks/%d/image", id)
	}

	sendReq := func(image string) *httptest.ResponseRecorder {
		fields := map[string]string{"image": image}
		contentType, body := createMultipartForm(fields)
		return sendHxRequestAsLoggedIn(srv, http.MethodPut, uri(1), header(contentType), strings.NewReader(body))
	}

	assert := func(t *testing.T, files *mockFiles, repo *mockRepository, gotStatusCode, wantStatusCode int, wantImage uuid.UUID, wantImageHitCount int) {
		assertStatus(t, gotStatusCode, wantStatusCode)
		assertUploadImageHitCount(t, files.uploadImageHitCount, wantImageHitCount)
		assertImage(t, repo.CookbooksRegistered[1][0].Image, wantImage)
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodDelete, uri(1))
	})

	t.Run("empty image in form", func(t *testing.T) {
		files, repo, revertFunc := prepareCookbook(srv)
		defer revertFunc()

		rr := sendReq("")

		assert(t, files, repo, rr.Code, http.StatusBadRequest, uuid.Nil, 0)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Could not retrieve the image from the form.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("upload image failed", func(t *testing.T) {
		files, repo, revertFunc := prepareCookbook(srv)
		files.uploadImageFunc = func(_ io.ReadCloser) (uuid.UUID, error) {
			return uuid.Nil, errors.New("error uploading")
		}
		srv.Files = files
		defer revertFunc()

		rr := sendReq("eggs.jpg")

		assert(t, files, repo, rr.Code, http.StatusInternalServerError, uuid.Nil, 0)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Error uploading image.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("updating image failed", func(t *testing.T) {
		files, repo, revertFunc := prepareCookbook(srv)
		repo.UpdateCookbookImageFunc = func(id int64, image uuid.UUID, userID int64) error {
			return errors.New("error")
		}
		srv.Repository = repo
		defer revertFunc()

		rr := sendReq("eggs.jpg")

		assert(t, files, repo, http.StatusInternalServerError, rr.Code, uuid.Nil, 1)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Error updating the cookbook's image.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("upload image", func(t *testing.T) {
		files, repo, revertFunc := prepareCookbook(srv)
		defer revertFunc()

		rr := sendReq("eggs.jpg")

		assertStatus(t, rr.Code, http.StatusCreated)
		assertUploadImageHitCount(t, files.uploadImageHitCount, 1)
		assertImageNotNil(t, repo.CookbooksRegistered[1][0].Image)
	})
}

func TestHandlers_Cookbooks_RecipesSearch(t *testing.T) {
	srv := newServerTest()
	originalRepo := srv.Repository

	uri := "/cookbooks/recipes/search"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	revert := func() {
		srv.Repository = originalRepo
	}

	missingInBodyTestcases := []struct {
		name      string
		body      string
		wantToast string
	}{
		{
			name:      "missing cookbook ID in body",
			body:      "page=1",
			wantToast: "Missing cookbook ID in body.",
		},
		{
			name:      "missing page number in body",
			body:      "id=1",
			wantToast: "Missing page number in body.",
		},
	}
	for _, tc := range missingInBodyTestcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(tc.body))

			assertStatus(t, rr.Code, http.StatusBadRequest)
			assertHeader(t, rr, "HX-Trigger", fmt.Sprintf(`{"showToast":"{\"message\":\"%s\",\"backgroundColor\":\"bg-red-500\"}"}`, tc.wantToast))
		})
	}

	t.Run("display cookbook recipes if search is empty", func(t *testing.T) {
		srv.Repository = &mockRepository{
			CookbooksRegistered: map[int64][]models.Cookbook{
				1: {
					{
						ID: 1,
						Recipes: models.Recipes{
							{ID: 1, Name: "Americano", Category: "breakfast"},
							{ID: 2, Name: "Chicken Latte", Category: "lunch"},
						},
					},
				},
			},
			RecipesRegistered: nil,
		}
		defer revert()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("id=1&page=1"))

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<li class="recipe cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700"><input type='hidden' name='recipe-id' value='1'/><div class="grid grid-cols-4"><img class="w-fit col-span-1 border-r h-[90px]" src="/data/images/00000000-0000-0000-0000-000000000000.jpg" alt="Recipe image"><div class="grid col-span-3 gap-1 p-2"><div class="grid grid-flow-col"><p class="font-semibold">Americano</p><div class="grid justify-end"><div class="pb-4 pl-4 text-xs"><span title="Remove recipe from cookbook" hx-delete="/cookbooks/1/recipes/1" hx-swap="outerHTML" hx-target="closest .recipe" hx-confirm="Are you sure you want to remove this recipe from the cookbook?" hx-indicator="#fullscreen-loader"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg></span></div></div></div><div class="text-xs"><span class="w-fit leading-none flex items-center justify-center p-2 text-blue-700 bg-blue-100 border border-blue-300 rounded-full dark:border-gray-800"> breakfast </span></div></div></div><div class="flex border-t dark:border-gray-800"><button class="w-full center hover:bg-gray-800 hover:text-white hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/recipes/1" hx-target="#content" hx-swap="innerHTML" hx-push-url="true"> View </button><p class="justify-self-end px-4 select-none cursor-move handle"> 1 </p></div></li>`,
			`<li class="recipe cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700"><input type='hidden' name='recipe-id' value='2'/><div class="grid grid-cols-4"><img class="w-fit col-span-1 border-r h-[90px]" src="/data/images/00000000-0000-0000-0000-000000000000.jpg" alt="Recipe image"><div class="grid col-span-3 gap-1 p-2"><div class="grid grid-flow-col"><p class="font-semibold">Chicken Latte</p><div class="grid justify-end"><div class="pb-4 pl-4 text-xs"><span title="Remove recipe from cookbook" hx-delete="/cookbooks/1/recipes/2" hx-swap="outerHTML" hx-target="closest .recipe" hx-confirm="Are you sure you want to remove this recipe from the cookbook?" hx-indicator="#fullscreen-loader"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg></span></div></div></div><div class="text-xs"><span class="w-fit leading-none flex items-center justify-center p-2 text-blue-700 bg-blue-100 border border-blue-300 rounded-full dark:border-gray-800"> lunch </span></div></div></div><div class="flex border-t dark:border-gray-800"><button class="w-full center hover:bg-gray-800 hover:text-white hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/recipes/2" hx-target="#content" hx-swap="innerHTML" hx-push-url="true"> View </button><p class="justify-self-end px-4 select-none cursor-move handle"> 2 </p></div></li>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("no results", func(t *testing.T) {
		_, _, revertFunc := prepareCookbook(srv)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("id=1&page=1&q=hello"))

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<div class="grid place-content-center text-sm text-center h-3/5 md:text-base"><p>No results found.</p></div>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	searches := []struct {
		query string
		want  []string
	}{
		{
			query: "lov",
			want: []string{
				`<article class="grid gap-4 p-4 text-sm justify-center md:p-0 md:text-base">`,
				`<ul class="cookbooks-display grid gap-2 p-2 md:p-0 text-sm md:text-base"><li class="recipe cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700"><div class="grid grid-cols-4"><img class="col-span-1 border-r" src="/data/images/00000000-0000-0000-0000-000000000000.jpg" alt="Recipe image"><div class="grid col-span-3 gap-1 p-2"><div class="grid grid-flow-col"><p class="font-semibold">Lovely Xenophobia</p><div class="grid justify-end"><div class="pb-4 pl-4 text-xs"><span class="w-fit leading-none flex items-center justify-center p-2 text-blue-700 bg-blue-100 border border-blue-300 rounded-full dark:border-gray-800"></span></div></div></div></div></div><button class="w-full border-t center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600 hover:rounded-b-md" hx-post="/cookbooks/1" hx-vals='{"cookbookId": 1, "recipeId": 2}' hx-swap="outerHTML" hx-target="closest .recipe"> Add </button></li>`,
			},
		},
		{
			query: "chi",
			want: []string{
				`<article class="grid gap-4 p-4 text-sm justify-center md:p-0 md:text-base">`,
				`<ul class="cookbooks-display grid gap-2 p-2 md:p-0 text-sm md:text-base"><li class="recipe cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700"><div class="grid grid-cols-4"><img class="col-span-1 border-r" src="/data/images/00000000-0000-0000-0000-000000000000.jpg" alt="Recipe image"><div class="grid col-span-3 gap-1 p-2"><div class="grid grid-flow-col"><p class="font-semibold">Chicken</p><div class="grid justify-end"><div class="pb-4 pl-4 text-xs"><span class="w-fit leading-none flex items-center justify-center p-2 text-blue-700 bg-blue-100 border border-blue-300 rounded-full dark:border-gray-800"></span></div></div></div></div></div><button class="w-full border-t center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600 hover:rounded-b-md" hx-post="/cookbooks/1" hx-vals='{"cookbookId": 1, "recipeId": 1}' hx-swap="outerHTML" hx-target="closest .recipe"> Add </button></li>`,
			},
		},
		{
			query: "lovely xenophobia",
			want: []string{
				`<article class="grid gap-4 p-4 text-sm justify-center md:p-0 md:text-base">`,
				`<ul class="cookbooks-display grid gap-2 p-2 md:p-0 text-sm md:text-base"><li class="recipe cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700"><div class="grid grid-cols-4"><img class="col-span-1 border-r" src="/data/images/00000000-0000-0000-0000-000000000000.jpg" alt="Recipe image"><div class="grid col-span-3 gap-1 p-2"><div class="grid grid-flow-col"><p class="font-semibold">Lovely Xenophobia</p><div class="grid justify-end"><div class="pb-4 pl-4 text-xs"><span class="w-fit leading-none flex items-center justify-center p-2 text-blue-700 bg-blue-100 border border-blue-300 rounded-full dark:border-gray-800"></span></div></div></div></div></div><button class="w-full border-t center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600 hover:rounded-b-md" hx-post="/cookbooks/1" hx-vals='{"cookbookId": 1, "recipeId": 2}' hx-swap="outerHTML" hx-target="closest .recipe"> Add </button></li>`,
			},
		},
	}
	for _, tc := range searches {
		t.Run("results search titles only "+tc.query, func(t *testing.T) {
			_, _, revertFunc := prepareCookbook(srv)
			defer revertFunc()

			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("id=1&page=1&q="+tc.query))

			assertStatus(t, rr.Code, http.StatusOK)
			assertStringsInHTML(t, getBodyHTML(rr), tc.want)
		})
	}
}

func TestHandlers_Cookbooks_ReorderRecipes(t *testing.T) {
	srv := newServerTest()

	uri := "/cookbooks/1/reorder"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPut, uri)
	})

	missingBodyPartsTestcases := []struct {
		name      string
		form      string
		wantToast string
	}{
		{
			name:      "missing cookbook ID",
			form:      "recipe-id=8&recipe-id=3",
			wantToast: "Missing cookbook ID in body.",
		},
		{
			name:      "missing recipe IDs",
			form:      "cookbook-id=1",
			wantToast: "Missing recipe IDs in body.",
		},
		{
			name:      "invalid recipe IDs",
			form:      "cookbook-id=1&recipe-id=8&recipe-id=3&recipe-id=0&recipe-id=-1",
			wantToast: `Recipe ID \"-1\" is invalid.`,
		},
	}
	for _, tc := range missingBodyPartsTestcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, uri, formHeader, strings.NewReader(tc.form))

			assertStatus(t, rr.Code, http.StatusBadRequest)
			assertHeader(t, rr, "HX-Trigger", fmt.Sprintf(`{"showToast":"{\"message\":\"%s\",\"backgroundColor\":\"bg-red-500\"}"}`, tc.wantToast))
		})
	}

	t.Run("valid request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, uri, formHeader, strings.NewReader("cookbook-id=1&recipe-id=1"))

		assertStatus(t, rr.Code, http.StatusNoContent)
	})
}

func TestHandlers_Cookbooks_Share(t *testing.T) {
	srv := newServerTest()

	uri := func(id int64) string {
		return fmt.Sprintf("/cookbooks/%d/share", id)
	}

	link := "/c/33320755-82f9-47e5-bb0a-d1b55cbd3f7b"
	link2 := "/c/43320755-82f9-47e5-bb0a-d1b55cbd3f72"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri(1))
	})

	t.Run("cookbook not registered", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodPost, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Failed to create share link.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("create valid share link", func(t *testing.T) {
		_, _, revert := prepareCookbook(srv)
		defer revert()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<input type="url" value="example.com` + link + `" class="w-full rounded-lg bg-gray-100 px-4 py-2" readonly="readonly">`,
			`<button class="w-24 font-semibold p-2 bg-gray-300 rounded-lg hover:bg-gray-400" title="Copy to clipboard" _="on click if window.navigator.clipboard then call navigator.clipboard.writeText('example.com/c/33320755-82f9-47e5-bb0a-d1b55cbd3f7b') then put 'Copied!' into me then add @title='Copied to clipboard!' then toggle @disabled on me then toggle .cursor-not-allowed .bg-green-600 .text-white .hover:bg-gray-400 on me else call alert('Your browser does not support the clipboard feature. Please copy the link manually.') end"> Copy </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("access share link anonymous no recipes", func(t *testing.T) {
		_, _, revert := prepareCookbook(srv)
		defer revert()

		rr := sendRequest(srv, http.MethodGet, link2, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Lovely Ukraine | Recipya</title>`,
			`<a href="/auth/login" class="flex-auto mr-2 rounded-lg px-2 py-1 hover:text-white hover:bg-green-600"> Log In </a>`,
			`<section class="grid justify-center p-4"><p class="grid justify-center font-semibold underline mt-2 md:mt-0 md:text-xl md:hidden"> Lovely Ukraine </p></section>`,
			`<div class="grid place-content-center text-sm text-center md:text-base" style="height: 50vh"><p>The user has not added recipes to this cookbook yet.</p></div>`,
		}
		notWant := []string{
			`<dialog id="share-dialog" class="p-4 border-4 border-black min-w-[15rem]"><div id="share-dialog-result" class="pb-4"></div>`,
		}
		body := getBodyHTML(rr)
		assertStringsInHTML(t, body, want)
		assertStringsNotInHTML(t, body, notWant)
	})

	t.Run("access share link anonymous have recipes", func(t *testing.T) {
		_, _, revert := prepareCookbook(srv)
		defer revert()

		rr := sendRequest(srv, http.MethodGet, link, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">Lovely Canada | Recipya</title>`,
			`<a href="/auth/login" class="flex-auto mr-2 rounded-lg px-2 py-1 hover:text-white hover:bg-green-600"> Log In </a>`,
			`<a href="/auth/register" class="flex-auto mr-4 rounded-lg px-2 py-1 bg-amber-300 dark:bg-orange-600 hover:text-white hover:bg-red-600"> Sign Up </a>`,
			`<section class="grid justify-center p-4"><p class="grid justify-center font-semibold underline mt-2 md:mt-0 md:text-xl "> Lovely Canada </p></section>`,
			`<form class="sortable" hx-put="/cookbooks/2/reorder" hx-trigger="end" hx-swap="none"><input type='hidden' name='cookbook-id' value='1'/><ul id="search-results" class="cookbooks-display grid gap-2 p-2 md:p-0 text-sm md:text-base">`,
			`<li class="recipe cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700"><input type='hidden' name='recipe-id' value='3'/><div class="grid grid-cols-4"><img class="w-fit col-span-1 border-r h-[90px]" src="/data/images/00000000-0000-0000-0000-000000000000.jpg" alt="Recipe image"><div class="grid col-span-3 gap-1 p-2"><div class="grid grid-flow-col"><p class="font-semibold">Gotcha</p><div class="grid justify-end"><div class="pb-4 pl-4 text-xs"></div></div></div><div class="text-xs"><span class="w-fit leading-none flex items-center justify-center p-2 text-blue-700 bg-blue-100 border border-blue-300 rounded-full dark:border-gray-800"> American </span></div></div></div><div class="flex border-t dark:border-gray-800"><button class="w-full center hover:bg-gray-800 hover:text-white hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/r/3?cookbook=1" hx-target="#content" hx-swap="innerHTML" hx-push-url="true"> View </button><p class="justify-self-end px-4 select-none cursor-none"> 1 </p></div></li>`,
		}
		notWant := []string{
			`<dialog id="share-dialog" class="p-4 border-4 border-black min-w-[15rem]"><div id="share-dialog-result" class="pb-4"></div>`,
		}
		body := getBodyHTML(rr)
		assertStringsInHTML(t, body, want)
		assertStringsNotInHTML(t, body, notWant)
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string, contentType header, body *strings.Reader) *httptest.ResponseRecorder
	}{
		{name: "other user Hx-Request", sendFunc: sendHxRequestAsLoggedInOther},
		{name: "other user no Hx-Request", sendFunc: sendRequestAsLoggedInOther},
	}
	for _, tc := range testcases {
		t.Run("access share link logged in "+tc.name, func(t *testing.T) {
			_, _, revert := prepareCookbook(srv)
			defer revert()

			rr := tc.sendFunc(srv, http.MethodGet, link, noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">Lovely Canada | Recipya</title>`,
				`<section class="grid gap-4 p-4 text-sm justify-center md:text-base">`,
				`<div class="flex flex-col h-full"><section class="grid justify-center p-4"><p class="grid justify-center font-semibold underline mt-2 md:mt-0 md:text-xl "> Lovely Canada </p></section></div>`,
				`<form class="sortable" hx-put="/cookbooks/2/reorder" hx-trigger="end" hx-swap="none"><input type='hidden' name='cookbook-id' value='1'/><ul id="search-results" class="cookbooks-display grid gap-2 p-2 md:p-0 text-sm md:text-base">`,
				`<li class="recipe cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700">`,
				`<input type='hidden' name='recipe-id' value='3'/><div class="grid grid-cols-4">`,
				`<img class="w-fit col-span-1 border-r h-[90px]" src="/data/images/00000000-0000-0000-0000-000000000000.jpg" alt="Recipe image">`,
				`<div class="grid col-span-3 gap-1 p-2"><div class="grid grid-flow-col"><p class="font-semibold">Gotcha</p><div class="grid justify-end"><div class="pb-4 pl-4 text-xs"></div></div></div><div class="text-xs"><span class="w-fit leading-none flex items-center justify-center p-2 text-blue-700 bg-blue-100 border border-blue-300 rounded-full dark:border-gray-800"> American </span></div></div></div><div class="flex border-t dark:border-gray-800"><button class="w-full center hover:bg-gray-800 hover:text-white hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/r/3?cookbook=1" hx-target="#content" hx-swap="innerHTML" hx-push-url="true"> View </button><p class="justify-self-end px-4 select-none cursor-none"> 1 </p></div></li>`,
			}
			notWant := []string{
				`<dialog id="share-dialog" class="p-4 border-4 border-black min-w-[15rem]"><div id="share-dialog-result" class="pb-4"></div>`,
				`<button class="mr-2" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call #share-dialog.showModal()">`,
			}
			body := getBodyHTML(rr)
			assertStringsInHTML(t, body, want)
			assertStringsNotInHTML(t, body, notWant)
		})
	}

	testcases2 := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string, contentType header, body *strings.Reader) *httptest.ResponseRecorder
	}{
		{name: "host user Hx-Request", sendFunc: sendHxRequestAsLoggedIn},
		{name: "host user no Hx-Request", sendFunc: sendRequestAsLoggedIn},
	}
	for _, tc := range testcases2 {
		t.Run("access share link logged in "+tc.name, func(t *testing.T) {
			_, _, revert := prepareCookbook(srv)
			defer revert()

			rr := tc.sendFunc(srv, http.MethodGet, link, noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">Lovely Canada | Recipya</title>`,
				`<section class="grid gap-4 p-4 text-sm justify-center md:text-base">`,
				`<form class="sortable" hx-put="/cookbooks/2/reorder" hx-trigger="end" hx-swap="none"><input type='hidden' name='cookbook-id' value='1'/><ul id="search-results" class="cookbooks-display grid gap-2 p-2 md:p-0 text-sm md:text-base">`,
				`<li class="recipe cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700">`,
				`<input type='hidden' name='recipe-id' value='3'/><div class="grid grid-cols-4">`,
				`<img class="w-fit col-span-1 border-r h-[90px]" src="/data/images/00000000-0000-0000-0000-000000000000.jpg" alt="Recipe image">`,
				`<li class="recipe cookbook relative grid max-w-[30rem] border bg-white rounded-md shadow-md md:min-w-[30rem] dark:bg-neutral-700"><input type='hidden' name='recipe-id' value='3'/><div class="grid grid-cols-4"><img class="w-fit col-span-1 border-r h-[90px]" src="/data/images/00000000-0000-0000-0000-000000000000.jpg" alt="Recipe image"><div class="grid col-span-3 gap-1 p-2"><div class="grid grid-flow-col"><p class="font-semibold">Gotcha</p><div class="grid justify-end"><div class="pb-4 pl-4 text-xs"><span title="Remove recipe from cookbook" hx-delete="/cookbooks/1/recipes/3" hx-swap="outerHTML" hx-target="closest .recipe" hx-confirm="Are you sure you want to remove this recipe from the cookbook?" hx-indicator="#fullscreen-loader"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg></span></div></div></div><div class="text-xs"><span class="w-fit leading-none flex items-center justify-center p-2 text-blue-700 bg-blue-100 border border-blue-300 rounded-full dark:border-gray-800"> American </span></div></div></div><div class="flex border-t dark:border-gray-800"><button class="w-full center hover:bg-gray-800 hover:text-white hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/recipes/3" hx-target="#content" hx-swap="innerHTML" hx-push-url="true"> View </button><p class="justify-self-end px-4 select-none cursor-move handle"> 1 </p></div></li>`,
			}
			notWant := []string{
				`<dialog id="share-dialog" class="p-4 border-4 border-black min-w-[15rem]"><div id="share-dialog-result" class="pb-4"></div>`,
				`<button class="mr-2" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call #share-dialog.showModal()">`,
			}
			body := getBodyHTML(rr)
			assertStringsInHTML(t, body, want)
			assertStringsNotInHTML(t, body, notWant)
		})
	}
}

func prepareCookbook(srv *server.Server) (*mockFiles, *mockRepository, func()) {
	originalFiles := srv.Files
	originalRepo := srv.Repository

	recipes := models.Recipes{{ID: 1, Name: "Chicken"}, {ID: 2, Name: "Lovely Xenophobia"}}
	recipe3 := models.Recipe{
		Category:     "American",
		Description:  "This is the most delicious recipe!",
		ID:           3,
		Image:        uuid.Nil,
		Ingredients:  []string{"Ing1", "Ing2", "Ing3"},
		Instructions: []string{"Ins1", "Ins2", "Ins3"},
		Name:         "Gotcha",
		Nutrition: models.Nutrition{
			Calories:           "500",
			Cholesterol:        "1g",
			Fiber:              "2g",
			Protein:            "3g",
			SaturatedFat:       "4g",
			Sodium:             "5g",
			Sugars:             "6g",
			TotalCarbohydrates: "7g",
			TotalFat:           "8g",
			UnsaturatedFat:     "9g",
		},
		Times: models.Times{
			Prep:  5 * time.Minute,
			Cook:  1*time.Hour + 5*time.Minute,
			Total: 1*time.Hour + 10*time.Minute,
		},
		URL:   "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/",
		Yield: 2,
	}

	files := &mockFiles{}
	repo := &mockRepository{
		CookbooksRegistered: map[int64][]models.Cookbook{
			1: {
				models.Cookbook{ID: 1, Title: "Lovely Canada", Recipes: models.Recipes{recipe3}},
				models.Cookbook{ID: 2, Title: "Lovely Ukraine"},
				models.Cookbook{ID: 3, Title: "Lovely America"},
			},
			2: {
				models.Cookbook{ID: 4, Title: "Lovely Canada 2"},
				models.Cookbook{ID: 5, Title: "Lovely Ukraine 2"},
				models.Cookbook{ID: 6, Title: "Lovely America 2"},
			},
		},
		RecipesRegistered: map[int64]models.Recipes{1: recipes},
		ShareLinks: map[string]models.Share{
			"/c/33320755-82f9-47e5-bb0a-d1b55cbd3f7b": {CookbookID: 1, RecipeID: -1, UserID: 1},
			"/c/43320755-82f9-47e5-bb0a-d1b55cbd3f72": {CookbookID: 2, RecipeID: -1, UserID: 1},
		},
	}
	srv.Files = files
	srv.Repository = repo

	return files, repo, func() {
		srv.Files = originalFiles
		srv.Repository = originalRepo
	}
}

func assertCookbooks(t testing.TB, got, want []models.Cookbook) {
	t.Helper()
	if !slices.EqualFunc(got, want, func(c1 models.Cookbook, c2 models.Cookbook) bool {
		return c1.ID == c2.ID && slices.EqualFunc(c1.Recipes, c2.Recipes, func(r1 models.Recipe, r2 models.Recipe) bool {
			return r1.ID == r2.ID
		})
	}) {
		t.Fatalf("got\n%+v\nbut want\n%+v for user 1", got, want)
	}
}
