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
			`<p class="pb-2">Your cookbooks collection looks a bit empty at the moment.</p>`,
			`<p>Why not start by creating a cookbook by clicking the <a class="underline font-semibold cursor-pointer" hx-post="/cookbooks" hx-prompt="Enter the name of your cookbook" hx-target=".cookbooks-display" hx-swap="beforeend">Add cookbook</a> button at the top?</p>`,
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Error getting cookbooks.\",\"title\":\"Database Error\"}"}`)
	})

	t.Run("have cookbooks grid preferred mode", func(t *testing.T) {
		_, repo, revertFunc := prepare(srv, models.GridViewMode)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.GridViewMode})
		assertCookbooksViewMode(t, models.GridViewMode, getBodyHTML(rr))
	})

	t.Run("have cookbooks list preferred mode", func(t *testing.T) {
		_, repo, revertFunc := prepare(srv, models.ListViewMode)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.ListViewMode})
		body := getBodyHTML(rr)
		assertCookbooksViewMode(t, models.ListViewMode, body)
		want := []string{
			`<div class="p-2 hover:bg-red-600 hover:text-white" title="Display as grid" hx-get="/cookbooks?view=grid" hx-target="#content" hx-trigger="mousedown">`,
			`<h2 class="card-title text-base w-[20ch] sm:w-[29ch] break-words md:text-xl">Lovely Canada</h2>`,
			`<h2 class="card-title text-base w-[20ch] sm:w-[29ch] break-words md:text-xl">Lovely America</h2>`,
			`<h2 class="card-title text-base w-[20ch] sm:w-[29ch] break-words md:text-xl">Lovely Ukraine</h2>`,
			`<form id="cookbook-image-form-1" enctype="multipart/form-data" hx-swap="none" hx-put="/cookbooks/1/image" hx-trigger="change from:#cookbook-image-1"><input id="cookbook-image-1" type="file" accept="image/*" name="image" required class="hidden" _="on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from next <button/>"></form>`,
			`<form id="cookbook-image-form-2" enctype="multipart/form-data" hx-swap="none" hx-put="/cookbooks/2/image" hx-trigger="change from:#cookbook-image-2">`,
			`<form id="cookbook-image-form-3" enctype="multipart/form-data" hx-swap="none" hx-put="/cookbooks/3/image" hx-trigger="change from:#cookbook-image-3">`,
			`<span class="three-dots-container indicator-item indicator-end badge badge-neutral rounded-md p-1 select-none cursor-pointer hover:bg-secondary" _="on mousedown openCookbookOptionsMenu(event)"><svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-three-dots-vertical" viewBox="0 0 16 16"><path d="M9.5 13a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0zm0-5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0zm0-5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0z"></path></svg></span>`,
			`<a id="cookbook_menu_share" hx-post="/cookbooks/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call share_dialog.showModal()">`,
			`<a id="cookbook_menu_download" hx-get="/cookbooks/1/download">`,
			`<a id="cookbook_menu_delete" hx-delete="/cookbooks/1" hx-swap="outerHTML" hx-target="closest .cookbook" hx-confirm="Are you sure you want to delete this cookbook? Its recipes will not be deleted.">`,
			`<button class="btn btn-outline btn-sm" hx-get="/cookbooks/1?page=1" hx-target="#content" hx-trigger="mousedown" hx-push-url="/cookbooks/1">Open</button>`,
			`<footer id="pagination" class="footer footer-center bg-base-200 pb-12 p-2 md:pb-2 text-base-content gap-2" onload="__templ_updateAddCookbookURL`,
			`<div class="join gap-0"><button class="join-item btn btn-disabled">«</button><!-- Left Section --><button aria-current="page" class="join-item btn btn-active">1</button><!-- Middle Section --><!-- Right Section --><button class="join-item btn btn-disabled">»</button></div><div class="text-center"><p class="text-sm">Showing <span class="font-medium">1</span> to <span class="font-medium">3</span> of <span id="search-count" class="font-medium">3</span> results</p></div></footer>`,
		}
		assertStringsInHTML(t, body, want)
	})

	t.Run("have cookbooks grid view select list", func(t *testing.T) {
		_, repo, revertFunc := prepare(srv, models.GridViewMode)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?view=list", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.ListViewMode})
		assertCookbooksViewMode(t, models.ListViewMode, getBodyHTML(rr))
	})

	t.Run("have cookbooks list view select grid", func(t *testing.T) {
		_, repo, revertFunc := prepare(srv, models.ListViewMode)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?view=grid", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.GridViewMode})
		assertCookbooksViewMode(t, models.GridViewMode, getBodyHTML(rr))
	})

	t.Run("title must not be empty", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader(""))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Title must not be empty.\",\"title\":\"Request Error\"}"}`)
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
			`<section id="cookbook-1" class="cookbook card card-compact bg-base-100 shadow-lg indicator">`,
			`<img class="rounded-t-lg w-full border-b h-32 text-center object-cover max-w-48 md:h-48 hover:bg-gray-100 hover:opacity-80" src="/static/img/cookbooks-new/placeholder.webp" onClick="__templ_cookbookImageClick`,
			`<form id="cookbook-image-form-1" enctype="multipart/form-data" hx-swap="none" hx-put="/cookbooks/1/image" hx-trigger="change from:#cookbook-image-1"><input id="cookbook-image-1" type="file" accept="image/*" name="image" required class="hidden" _="on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from next <button/>"></form>`,
			`<p class="font-semibold w-[18ch] break-words">Lovely America</p>`,
			`<button class="btn btn-block btn-sm btn-outline" hx-get="/cookbooks/1?page=1" hx-target="#content" hx-trigger="mousedown" hx-push-url="/cookbooks/1">Open</button>`,
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
		t.Helper()
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Error deleting cookbook.\",\"title\":\"Database Error\"}"}`)
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

	t.Run("cookbook does not exist", func(t *testing.T) {
		_, _, revertFunc := prepare(srv)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(20), noHeader, nil)

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
			`<search><form class="w-72 flex md:w-96" hx-get="/cookbooks/4/recipes/search" hx-vals="{"page": 1}" hx-target="#search-results" hx-push-url="true" hx-trigger="submit, change target:.sort-option"><div class="relative w-full"><label><input type="search" id="search-recipes" name="q" class="input input-bordered input-sm w-full z-20" placeholder="Search for recipes..." value="" _="on keyup if event.target.value !== '' then remove .md:block from #search-shortcut else add .md:block to #search-shortcut end then if (event.key === 'Backspace' or event.key === 'Delete') and event.target.value === '' then send submit to closest <form/> end"></label> <kbd id="search-shortcut" class="hidden absolute top-1 right-12 font-sans font-semibold select-none dark:text-slate-500 md:block"><abbr title="Control" class="no-underline text-slate-300 dark:text-slate-500">Ctrl </abbr> /</kbd> <button type="submit" class="absolute top-0 right-0 px-2 btn btn-sm btn-primary"><svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"></path></svg><span class="sr-only">Search</span></button></div>`,
			`<details class="dropdown dropdown-left"><summary class="btn btn-sm ml-1"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75a4.5 4.5 0 0 1-4.884 4.484c-1.076-.091-2.264.071-2.95.904l-7.152 8.684a2.548 2.548 0 1 1-3.586-3.586l8.684-7.152c.833-.686.995-1.874.904-2.95a4.5 4.5 0 0 1 6.336-4.486l-3.276 3.276a3.004 3.004 0 0 0 2.25 2.25l3.276-3.276c.256.565.398 1.192.398 1.852Z"></path> <path stroke-linecap="round" stroke-linejoin="round" d="M4.867 19.125h.008v.008h-.008v-.008Z"></path></svg></summary><div tabindex="0" class="dropdown-content z-10 menu menu-sm p-2 shadow bg-base-200 w-52 sm:menu-md prose" _="on click remove @open from closest <details/>"><h4>Search Method</h4><div class="form-control"><label class="label cursor-pointer"><span class="label-text">By name</span> <input type="radio" name="mode" class="radio radio-sm" value="name"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Full search</span> <input type="radio" name="mode" class="radio radio-sm" value="full" checked></label></div></div></details>`,
			`<details class="dropdown dropdown-left"><summary class="btn btn-sm ml-1"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12"></path></svg></summary><div tabindex="0" class="dropdown-content z-10 menu menu-sm p-2 shadow bg-base-200 w-52 sm:menu-md prose" _="on click remove @open from closest <details/>"><h4>Sort</h4><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Default</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="default" checked></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Name:<br>A to Z</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="a-z"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Name:<br>Z to A</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="z-a"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Date created:<br>Newest to oldest</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="new-old"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Date created:<br>Oldest to newest</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="old-new"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Random</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="random"></label></div></div></details>`,
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
					`<title hx-swap-oob="true">Lovely Canada | Recipya</title>`,
					`<div id="content-title" hx-swap-oob="innerHTML">Lovely Canada</div>`,
					`<script defer> function initReorder()`,
					`<form class="w-72 flex md:w-96" hx-get="/cookbooks/1/recipes/search" hx-vals="{"page": 1}" hx-target="#search-results" hx-push-url="true" hx-trigger="submit, change target:.sort-option"><div class="relative w-full"><label><input type="search" id="search-recipes" name="q" class="input input-bordered input-sm w-full z-20" placeholder="Search for recipes..." value="" _="on keyup if event.target.value !== '' then remove .md:block from #search-shortcut else add .md:block to #search-shortcut end then if (event.key === 'Backspace' or event.key === 'Delete') and event.target.value === '' then send submit to closest <form/> end"></label> <kbd id="search-shortcut" class="hidden absolute top-1 right-12 font-sans font-semibold select-none dark:text-slate-500 md:block"><abbr title="Control" class="no-underline text-slate-300 dark:text-slate-500">Ctrl </abbr> /</kbd> <button type="submit" class="absolute top-0 right-0 px-2 btn btn-sm btn-primary"><svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"></path></svg><span class="sr-only">Search</span></button></div><details class="dropdown dropdown-left"><summary class="btn btn-sm ml-1"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75a4.5 4.5 0 0 1-4.884 4.484c-1.076-.091-2.264.071-2.95.904l-7.152 8.684a2.548 2.548 0 1 1-3.586-3.586l8.684-7.152c.833-.686.995-1.874.904-2.95a4.5 4.5 0 0 1 6.336-4.486l-3.276 3.276a3.004 3.004 0 0 0 2.25 2.25l3.276-3.276c.256.565.398 1.192.398 1.852Z"></path> <path stroke-linecap="round" stroke-linejoin="round" d="M4.867 19.125h.008v.008h-.008v-.008Z"></path></svg></summary><div tabindex="0" class="dropdown-content z-10 menu menu-sm p-2 shadow bg-base-200 w-52 sm:menu-md prose" _="on click remove @open from closest <details/>"><h4>Search Method</h4><div class="form-control"><label class="label cursor-pointer"><span class="label-text">By name</span> <input type="radio" name="mode" class="radio radio-sm" value="name"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Full search</span> <input type="radio" name="mode" class="radio radio-sm" value="full" checked></label></div></div></details> <details class="dropdown dropdown-left"><summary class="btn btn-sm ml-1"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12"></path></svg></summary><div tabindex="0" class="dropdown-content z-10 menu menu-sm p-2 shadow bg-base-200 w-52 sm:menu-md prose" _="on click remove @open from closest <details/>"><h4>Sort</h4><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Default</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="default" checked></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Name:<br>A to Z</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="a-z"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Name:<br>Z to A</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="z-a"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Date created:<br>Newest to oldest</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="new-old"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Date created:<br>Oldest to newest</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="old-new"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Random</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="random"></label></div></div></details></form>`,
					`<div class="card card-side card-bordered card-compact bg-base-100 shadow-lg sm:w-[30rem]">`,
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

	t.Run("recipeID missing in form", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodPost, uri(1), formHeader, strings.NewReader("cookbookId=1"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Missing 'recipeId' in body.\",\"title\":\"Form Error\"}"}`)
	})

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

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri(2), formHeader, strings.NewReader("recipeId=2"))

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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Could not fetch cookbook.\",\"title\":\"Database Error\"}"}`)
	})

	t.Run("other user cannot download someone else's cookbook", func(t *testing.T) {
		_, _, revert := prepareCookbook(srv)
		files := &mockFiles{}
		srv.Files = files
		defer revert()

		rr := sendHxRequestAsLoggedInOther(srv, http.MethodGet, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Could not fetch cookbook.\",\"title\":\"Database Error\"}"}`)
		if files.exportHitCount > 0 {
			t.Fatal("export function must not have been called")
		}
	})

	t.Run("cannot download a cookbook with no recipes", func(t *testing.T) {
		_, _, revert := prepareCookbook(srv)
		files := &mockFiles{}
		srv.Files = files
		defer revert()

		rr := sendHxRequestAsLoggedInOther(srv, http.MethodGet, uri(5), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Cookbook is empty.\",\"title\":\"Request Error\"}"}`)
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
		t.Helper()
		assertStatus(t, gotStatusCode, wantStatusCode)
		assertUploadImageHitCount(t, files.uploadImageHitCount, wantImageHitCount)
		assertImage(t, repo.CookbooksRegistered[1][0].Image, wantImage)
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPut, uri(1))
	})

	t.Run("empty image in form", func(t *testing.T) {
		files, repo, revertFunc := prepareCookbook(srv)
		defer revertFunc()

		rr := sendReq("")

		assert(t, files, repo, rr.Code, http.StatusBadRequest, uuid.Nil, 0)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Could not retrieve the image from the form.\",\"title\":\"Form Error\"}"}`)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Error uploading image.\",\"title\":\"Files Error\"}"}`)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Error updating image.\",\"title\":\"Database Error\"}"}`)
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

	uri := func(id int64) string {
		return fmt.Sprintf("/cookbooks/%d/recipes/search", id)
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri(1))
	})

	t.Run("id must be positive", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(0), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertStringsInHTML(t, getBodyHTML(rr), []string{"Cookbook ID in path must be > 0."})
	})

	t.Run("page query parameter is mandatory", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertStringsInHTML(t, getBodyHTML(rr), []string{"Missing page number in query."})
	})

	t.Run("display cookbook recipes if search is empty", func(t *testing.T) {
		_, _, revertFunc := prepareCookbook(srv)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(1)+"?page=1", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Redirect", "/cookbooks/1")
		assertStringsInHTML(t, getBodyHTML(rr), []string{"Query parameter must not be 'q' empty."})
	})

	t.Run("no results", func(t *testing.T) {
		_, _, revertFunc := prepareCookbook(srv)
		defer revertFunc()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(1)+"?page=1&q=hello&mode=name", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<div class="grid place-content-center text-sm text-center h-3/5 md:text-base"><p>No results found.</p></div>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	searchMethods := []struct {
		name string
		mode string
		want []string
	}{
		{
			name: "empty defaults to by name",
			mode: "",
			want: []string{
				`<li class="indicator recipe cookbook"><div class="card card-side card-compact bg-base-100 shadow-md max-w-[30rem] border indicator dark:border-slate-600"><figure class="relative"><img src="/static/img/recipes/placeholder.webp" alt="Recipe image" class="w-28 h-full object-cover"></figure><div class="card-body"><h2 class="card-title text-base w-[20ch] sm:w-[29ch] break-words md:text-xl">Best Hot Dogs</h2><p></p><div class="text-sm sm:text-base"><p class="text-sm pb-1">Category:</p><div class="badge badge-primary badge-">American</div></div><div class="card-actions justify-end"><button class="btn btn-outline btn-sm" hx-post="/cookbooks/1" hx-vals="{"recipeId": 3}" hx-swap="outerHTML" hx-target="closest .recipe" _="on click put (#search-count.textContent as Number) - 1 into #search-count">Add</button></div></div></div></li>`,
			},
		},
		{
			name: "name",
			mode: "name",
			want: []string{
				`<li class="indicator recipe cookbook"><div class="card card-side card-compact bg-base-100 shadow-md max-w-[30rem] border indicator dark:border-slate-600"><figure class="relative"><img src="/static/img/recipes/placeholder.webp" alt="Recipe image" class="w-28 h-full object-cover"></figure><div class="card-body"><h2 class="card-title text-base w-[20ch] sm:w-[29ch] break-words md:text-xl">Best Hot Dogs</h2><p></p><div class="text-sm sm:text-base"><p class="text-sm pb-1">Category:</p><div class="badge badge-primary badge-">American</div></div><div class="card-actions justify-end"><button class="btn btn-outline btn-sm" hx-post="/cookbooks/1" hx-vals="{"recipeId": 3}" hx-swap="outerHTML" hx-target="closest .recipe" _="on click put (#search-count.textContent as Number) - 1 into #search-count">Add</button></div></div></div></li>`,
			},
		},
		{
			name: "full",
			mode: "full",
			want: []string{
				`<li class="indicator recipe cookbook"><div class="card card-side card-compact bg-base-100 shadow-md max-w-[30rem] border indicator dark:border-slate-600"><figure class="relative"><img src="/static/img/recipes/placeholder.webp" alt="Recipe image" class="w-28 h-full object-cover"></figure><div class="card-body"><h2 class="card-title text-base w-[20ch] sm:w-[29ch] break-words md:text-xl">Best Hot Dogs</h2><p></p><div class="text-sm sm:text-base"><p class="text-sm pb-1">Category:</p><div class="badge badge-primary badge-">American</div></div><div class="card-actions justify-end"><button class="btn btn-outline btn-sm" hx-post="/cookbooks/1" hx-vals="{"recipeId": 3}" hx-swap="outerHTML" hx-target="closest .recipe" _="on click put (#search-count.textContent as Number) - 1 into #search-count">Add</button></div></div></div></li>`,
				`<li class="indicator recipe cookbook"><div class="card card-side card-compact bg-base-100 shadow-md max-w-[30rem] border indicator dark:border-slate-600"><figure class="relative"><img src="/static/img/recipes/placeholder.webp" alt="Recipe image" class="w-28 h-full object-cover"></figure><div class="card-body"><h2 class="card-title text-base w-[20ch] sm:w-[29ch] break-words md:text-xl">Grandma's Shepard's Pie</h2><p></p><div class="text-sm sm:text-base"><p class="text-sm pb-1">Category:</p><div class="badge badge-primary badge-">Chinese</div></div><div class="card-actions justify-end"><button class="btn btn-outline btn-sm" hx-post="/cookbooks/1" hx-vals="{"recipeId": 4}" hx-swap="outerHTML" hx-target="closest .recipe" _="on click put (#search-count.textContent as Number) - 1 into #search-count">Add</button></div></div></div></li>`,
			},
		},
	}
	for _, tc := range searchMethods {
		t.Run(tc.name, func(t *testing.T) {
			_, _, revertFunc := prepareCookbook(srv)
			defer revertFunc()

			rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(1)+"?page=1&q=best&mode="+tc.mode, noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			assertStringsInHTML(t, getBodyHTML(rr), tc.want)
		})
	}

	searches := []struct {
		query string
		want  []string
	}{
		{
			query: "lov",
			want: []string{
				`<article class="grid gap-8 p-4 text-sm justify-center md:p-0"><ul class="cookbooks-display grid gap-2 p-2 md:p-0"><li class="indicator recipe cookbook"><div class="card card-side card-compact bg-base-100 shadow-md max-w-[30rem] border indicator dark:border-slate-600"><figure class="relative"><img src="/static/img/recipes/placeholder.webp" alt="Recipe image" class="w-28 h-full object-cover"></figure><div class="card-body"><h2 class="card-title text-base w-[20ch] sm:w-[29ch] break-words md:text-xl">Lovely Xenophobia</h2><p></p><div class="text-sm sm:text-base"><p class="text-sm pb-1">Category:</p><div class="badge badge-primary badge-"></div></div><div class="card-actions justify-end"><button class="btn btn-outline btn-sm" hx-post="/cookbooks/1" hx-vals="{"recipeId": 2}" hx-swap="outerHTML" hx-target="closest .recipe" _="on click put (#search-count.textContent as Number) - 1 into #search-count">Add</button></div></div></div></li></ul></article>`,
			},
		},
		{
			query: "chi",
			want: []string{
				`<article class="grid gap-8 p-4 text-sm justify-center md:p-0"><ul class="cookbooks-display grid gap-2 p-2 md:p-0"><li class="indicator recipe cookbook"><div class="card card-side card-compact bg-base-100 shadow-md max-w-[30rem] border indicator dark:border-slate-600"><figure class="relative"><img src="/static/img/recipes/placeholder.webp" alt="Recipe image" class="w-28 h-full object-cover"></figure><div class="card-body"><h2 class="card-title text-base w-[20ch] sm:w-[29ch] break-words md:text-xl">Chicken</h2><p></p><div class="text-sm sm:text-base"><p class="text-sm pb-1">Category:</p><div class="badge badge-primary badge-"></div></div><div class="card-actions justify-end"><button class="btn btn-outline btn-sm" hx-post="/cookbooks/1" hx-vals="{"recipeId": 1}" hx-swap="outerHTML" hx-target="closest .recipe" _="on click put (#search-count.textContent as Number) - 1 into #search-count">Add</button></div></div></div></li></ul></article>`,
			},
		},
		{
			query: "lovely+xenophobia",
			want: []string{
				`<article class="grid gap-8 p-4 text-sm justify-center md:p-0"><ul class="cookbooks-display grid gap-2 p-2 md:p-0"><li class="indicator recipe cookbook"><div class="card card-side card-compact bg-base-100 shadow-md max-w-[30rem] border indicator dark:border-slate-600"><figure class="relative"><img src="/static/img/recipes/placeholder.webp" alt="Recipe image" class="w-28 h-full object-cover"></figure><div class="card-body"><h2 class="card-title text-base w-[20ch] sm:w-[29ch] break-words md:text-xl">Lovely Xenophobia</h2><p></p><div class="text-sm sm:text-base"><p class="text-sm pb-1">Category:</p><div class="badge badge-primary badge-"></div></div><div class="card-actions justify-end"><button class="btn btn-outline btn-sm" hx-post="/cookbooks/1" hx-vals="{"recipeId": 2}" hx-swap="outerHTML" hx-target="closest .recipe" _="on click put (#search-count.textContent as Number) - 1 into #search-count">Add</button></div></div></div></li></ul></article>`,
			},
		},
	}
	for _, tc := range searches {
		t.Run("results search titles only "+tc.query, func(t *testing.T) {
			_, _, revertFunc := prepareCookbook(srv)
			defer revertFunc()

			rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri(1)+"?page=1&q="+tc.query, noHeader, nil)

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
			name:      "missing recipe IDs",
			form:      "cookbook-id=1",
			wantToast: "Missing recipe IDs in body.",
		},
		{
			name:      "invalid recipe IDs",
			form:      "recipe-id=8&recipe-id=3&recipe-id=0&recipe-id=-1",
			wantToast: `Recipe ID \\\"-1\\\" is invalid.`,
		},
	}
	for _, tc := range missingBodyPartsTestcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, uri, formHeader, strings.NewReader(tc.form))

			assertStatus(t, rr.Code, http.StatusBadRequest)
			assertHeader(t, rr, "HX-Trigger", fmt.Sprintf(`{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"%s\",\"title\":\"Form Error\"}"}`, tc.wantToast))
		})
	}

	t.Run("valid request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, uri, formHeader, strings.NewReader("recipe-id=1"))

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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to create share link.\",\"title\":\"Database Error\"}"}`)
	})

	t.Run("create valid share link", func(t *testing.T) {
		_, _, revert := prepareCookbook(srv)
		defer revert()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<div class="grid grid-flow-col gap-2"><label><input type="url" value="example.com/c/33320755-82f9-47e5-bb0a-d1b55cbd3f7b" class="input input-bordered w-full" readonly="readonly"></label> <script type="text/javascript">function __templ_copyToClipboard`,
			`{if (window.navigator.clipboard) { navigator.clipboard.writeText(text); copy_button.textContent = "Copied!"; copy_button.setAttribute("disabled", true); copy_button.classList.toggle(".btn-disabled"); } else { alert('Your browser does not support the clipboard feature. Please copy the link manually.'); }}</script><button id="copy_button" class="btn btn-neutral" title="Copy to clipboard" onClick="__templ_copyToClipboard`,
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
			`<a class="btn btn-ghost text-lg" href="/">Recipya</a>`,
			`<a href="/auth/login" class="btn btn-ghost">Log In</a>`,
			`<a href="/auth/register" class="btn btn-ghost">Sign Up</a>`,
			`<section class="grid justify-center p-2 sm:p-4"><p class="grid justify-center font-semibold underline mt-4 md:hidden">Lovely Ukraine</p></section>`,
			`<p>The user has not added recipes to this cookbook yet.</p></div>`,
		}
		notWant := []string{
			`id="share-dialog"`,
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
			`<a href="/auth/login" class="btn btn-ghost">Log In</a>`,
			`<a href="/auth/register" class="btn btn-ghost">Sign Up</a>`,
			`<section class="grid gap-4 text-sm justify-center md:p-4 md:text-base"><div class="flex flex-col h-full"><section class="grid justify-center p-2 sm:p-4 sm:pb-0"><p class="grid justify-center font-semibold underline mt-4 md:mt-0 md:text-xl">Lovely Canada</p></section>`,
			`<form hx-put="/cookbooks/1/reorder" hx-trigger="end" hx-swap="none"><input type="hidden" name="cookbook-id" value="1"><ul class="cookbooks-display grid gap-8 p-2 place-items-center text-sm md:p-0 md:text-base"><li class="indicator recipe cookbook"><input type="hidden" name="recipe-id" value="3"><div class="indicator-item indicator-bottom badge badge-secondary cursor-none">1</div><div class="card card-side card-bordered card-compact bg-base-100 shadow-lg sm:w-[30rem]"><figure class="w-28 sm:w-32"><img src="/static/img/recipes/placeholder.webp" alt="Recipe image" class="object-cover"></figure><div class="card-body"><h2 class="card-title text-base w-[20ch] sm:w-full break-words md:text-xl">Gotcha</h2><p></p><div><p class="text-sm pb-1">Category:</p><div class="badge badge-primary badge-">American</div></div><div class="card-actions justify-end"><button class="btn btn-outline btn-sm" hx-get="/r/3?cookbook=1" hx-target="#content" hx-trigger="mousedown" hx-swap="innerHTML" hx-push-url="true">View</button></div></div></div></li></ul></form>`,
		}
		body := getBodyHTML(rr)
		assertStringsInHTML(t, body, want)
		notWant := []string{`id="share-dialog"`}
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
				`<section class="grid gap-4 text-sm justify-center md:p-4 md:text-base"><div class="flex flex-col h-full"><section class="grid justify-center p-2 sm:p-4 sm:pb-0"><p class="grid justify-center font-semibold underline mt-4 md:mt-0 md:text-xl">Lovely Canada</p></section>`,
				`<form hx-put="/cookbooks/1/reorder" hx-trigger="end" hx-swap="none"><input type="hidden" name="cookbook-id" value="1"><ul class="cookbooks-display grid gap-8 p-2 place-items-center text-sm md:p-0 md:text-base"><li class="indicator recipe cookbook"><input type="hidden" name="recipe-id" value="3"><div class="indicator-item indicator-bottom badge badge-secondary cursor-none">1</div><div class="card card-side card-bordered card-compact bg-base-100 shadow-lg sm:w-[30rem]"><figure class="w-28 sm:w-32"><img src="/static/img/recipes/placeholder.webp" alt="Recipe image" class="object-cover"></figure><div class="card-body"><h2 class="card-title text-base w-[20ch] sm:w-full break-words md:text-xl">Gotcha</h2><p></p><div><p class="text-sm pb-1">Category:</p><div class="badge badge-primary badge-">American</div></div><div class="card-actions justify-end"><button class="btn btn-outline btn-sm" hx-get="/r/3?cookbook=1" hx-target="#content" hx-trigger="mousedown" hx-swap="innerHTML" hx-push-url="true">View</button></div></div></div></li></ul></form>`,
			}
			body := getBodyHTML(rr)
			assertStringsInHTML(t, body, want)
			notWant := []string{`id="share-dialog"`, `title="Share recipe"`}
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
				`<section class="grid gap-4 text-sm justify-center md:p-4 md:text-base"><div class="flex flex-col h-full"><section class="grid justify-center p-2 sm:p-4 sm:pb-0">`,
				`<search><form class="w-72 flex md:w-96" hx-get="/cookbooks/2/recipes/search" hx-vals="{"page": 1}" hx-target="#search-results" hx-push-url="true" hx-trigger="submit, change target:.sort-option"><div class="relative w-full"><label><input type="search" id="search-recipes" name="q" class="input input-bordered input-sm w-full z-20" placeholder="Search for recipes..." value="" _="on keyup if event.target.value !== '' then remove .md:block from #search-shortcut else add .md:block to #search-shortcut end then if (event.key === 'Backspace' or event.key === 'Delete') and event.target.value === '' then send submit to closest <form/> end"></label> <kbd id="search-shortcut" class="hidden absolute top-1 right-12 font-sans font-semibold select-none dark:text-slate-500 md:block"><abbr title="Control" class="no-underline text-slate-300 dark:text-slate-500">Ctrl </abbr> /</kbd> <button type="submit" class="absolute top-0 right-0 px-2 btn btn-sm btn-primary"><svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"></path></svg><span class="sr-only">Search</span></button></div><details class="dropdown dropdown-left"><summary class="btn btn-sm ml-1"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75a4.5 4.5 0 0 1-4.884 4.484c-1.076-.091-2.264.071-2.95.904l-7.152 8.684a2.548 2.548 0 1 1-3.586-3.586l8.684-7.152c.833-.686.995-1.874.904-2.95a4.5 4.5 0 0 1 6.336-4.486l-3.276 3.276a3.004 3.004 0 0 0 2.25 2.25l3.276-3.276c.256.565.398 1.192.398 1.852Z"></path> <path stroke-linecap="round" stroke-linejoin="round" d="M4.867 19.125h.008v.008h-.008v-.008Z"></path></svg></summary><div tabindex="0" class="dropdown-content z-10 menu menu-sm p-2 shadow bg-base-200 w-52 sm:menu-md prose" _="on click remove @open from closest <details/>"><h4>Search Method</h4><div class="form-control"><label class="label cursor-pointer"><span class="label-text">By name</span> <input type="radio" name="mode" class="radio radio-sm" value="name"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Full search</span> <input type="radio" name="mode" class="radio radio-sm" value="full"></label></div></div></details> <details class="dropdown dropdown-left"><summary class="btn btn-sm ml-1"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12"></path></svg></summary><div tabindex="0" class="dropdown-content z-10 menu menu-sm p-2 shadow bg-base-200 w-52 sm:menu-md prose" _="on click remove @open from closest <details/>"><h4>Sort</h4><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Default</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="default"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Name:<br>A to Z</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="a-z"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Name:<br>Z to A</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="z-a"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Date created:<br>Newest to oldest</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="new-old"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Date created:<br>Oldest to newest</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="old-new"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Random</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="random"></label></div></div></details></form></search>`,
				`<p class="grid justify-center font-semibold underline mt-4 md:mt-0 md:text-xl md:hidden">Lovely Canada</p></section></div><div id="search-results" class="md:min-h-[79vh]"><form hx-put="/cookbooks/1/reorder" hx-trigger="end" hx-swap="none"><input type="hidden" name="cookbook-id" value="1"><ul class="cookbooks-display grid gap-8 p-2 place-items-center text-sm md:p-0 md:text-base"><li class="indicator recipe cookbook"><input type="hidden" name="recipe-id" value="3"><div class="indicator-item indicator-bottom badge badge-secondary cursor-move handle">1</div><div class="indicator-item badge badge-neutral h-6 w-8"><button title="Remove recipe from cookbook" class="btn btn-ghost btn-xs p-0" hx-delete="/cookbooks/1/recipes/3" hx-swap="outerHTML" hx-target="closest .recipe" hx-confirm="Are you sure you want to remove this recipe from the cookbook?" hx-indicator="#fullscreen-loader"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg></button></div><div class="card card-side card-bordered card-compact bg-base-100 shadow-lg sm:w-[30rem]"><figure class="w-28 sm:w-32"><img src="/static/img/recipes/placeholder.webp" alt="Recipe image" class="object-cover"></figure><div class="card-body"><h2 class="card-title text-base w-[20ch] sm:w-full break-words md:text-xl">Gotcha</h2><p></p><div><p class="text-sm pb-1">Category:</p><div class="badge badge-primary badge-">American</div></div><div class="card-actions justify-end"><button class="btn btn-outline btn-sm" hx-get="/recipes/3" hx-target="#content" hx-swap="innerHTML" hx-push-url="true">View</button></div></div></div></li></ul></form></div>`,
			}
			body := getBodyHTML(rr)
			assertStringsInHTML(t, body, want)
			notWant := []string{`id="share-dialog"`, `title="Share recipe"`}
			assertStringsNotInHTML(t, body, notWant)
		})
	}
}

func prepareCookbook(srv *server.Server) (*mockFiles, *mockRepository, func()) {
	originalFiles := srv.Files
	originalRepo := srv.Repository

	recipes := models.Recipes{
		{ID: 1, Name: "Chicken"},
		{ID: 2, Name: "Lovely Xenophobia"},
		{ID: 3, Name: "Best Hot Dogs", Category: "American"},
		{ID: 4, Name: "Grandma's Shepard's Pie", Category: "Chinese", Description: "The best chinese pie ever originally not from China."},
		{ID: 5, Name: "Double angus hamburger", Category: "American", Description: "The most delicious angus hamburger in the world."},
	}
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
				models.Cookbook{ID: 4, Title: "Lovely Canada 2", Recipes: models.Recipes{recipe3}},
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

func assertCookbooks(tb testing.TB, got, want []models.Cookbook) {
	tb.Helper()
	if !slices.EqualFunc(got, want, func(c1 models.Cookbook, c2 models.Cookbook) bool {
		return c1.ID == c2.ID && slices.EqualFunc(c1.Recipes, c2.Recipes, func(r1 models.Recipe, r2 models.Recipe) bool {
			return r1.ID == r2.ID
		})
	}) {
		tb.Fatalf("got\n%+v\nbut want\n%+v for user 1", got, want)
	}
}
