package server_test

import (
	"errors"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"slices"
	"strings"
	"testing"
)

func TestHandlers_Cookbooks(t *testing.T) {
	srv := newServerTest()
	originalRepo := &mockRepository{
		CookbooksRegistered: map[int64][]models.Cookbook{1: {}},
	}
	srv.Repository = originalRepo

	uri := "/cookbooks"

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
			`<p>Why not start by creating a cookbook by clicking the <a class="underline font-semibold cursor-pointer" hx-post="/cookbooks" hx-prompt="Enter the name of your cookbook" hx-target="#content" hx-push-url="true">Add cookbook</a> button at the top? </p>`,
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
		repo := originalRepo
		repo.CookbooksRegistered[1] = []models.Cookbook{
			{ID: 1, Title: "Lovely Canada"},
			{ID: 2, Title: "Lovely America"},
			{ID: 3, Title: "Lovely Ukraine"},
		}
		repo.UserSettingsRegistered = map[int64]*models.UserSettings{1: {CookbooksViewMode: models.GridViewMode}}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, 1, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.GridViewMode})
		assertCookbooksViewMode(t, models.GridViewMode, getBodyHTML(rr))
	})

	t.Run("have cookbooks list preferred mode", func(t *testing.T) {
		repo := originalRepo
		repo.CookbooksRegistered[1] = []models.Cookbook{
			{ID: 1, Title: "Lovely Canada"},
			{ID: 2, Title: "Lovely America"},
			{ID: 3, Title: "Lovely Ukraine"},
		}
		repo.UserSettingsRegistered = map[int64]*models.UserSettings{1: {CookbooksViewMode: models.ListViewMode}}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, 1, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.ListViewMode})
		assertCookbooksViewMode(t, models.ListViewMode, getBodyHTML(rr))
	})

	t.Run("have cookbooks grid view select list", func(t *testing.T) {
		repo := originalRepo
		repo.CookbooksRegistered[1] = []models.Cookbook{
			{ID: 1, Title: "Lovely Canada"},
			{ID: 2, Title: "Lovely America"},
			{ID: 3, Title: "Lovely Ukraine"},
		}
		repo.UserSettingsRegistered = map[int64]*models.UserSettings{1: {CookbooksViewMode: models.GridViewMode}}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?view=list", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertUserSettings(t, 1, repo.UserSettingsRegistered[1], &models.UserSettings{CookbooksViewMode: models.ListViewMode})
		assertCookbooksViewMode(t, models.ListViewMode, getBodyHTML(rr))
	})

	t.Run("have cookbooks list view select grid", func(t *testing.T) {
		repo := originalRepo
		repo.CookbooksRegistered[1] = []models.Cookbook{
			{ID: 1, Title: "Lovely Canada"},
			{ID: 2, Title: "Lovely America"},
			{ID: 3, Title: "Lovely Ukraine"},
		}
		repo.UserSettingsRegistered = map[int64]*models.UserSettings{1: {CookbooksViewMode: models.ListViewMode}}
		defer func() {
			srv.Repository = originalRepo
		}()

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
		defer func() {
			srv.Repository = originalRepo
		}()
		title := "Lovely America"

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader(title))

		assertStatus(t, rr.Code, http.StatusCreated)
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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Cookbook created.\",\"backgroundColor\":\"bg-blue-500\"}"}`)
		want := []string{
			`<div class="col-span-1 bg-white rounded-lg shadow-md dark:bg-neutral-700">`,
			`<img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-[180px] px-4 pt-4" src="" alt="Cookbook image">`,
			`<p class="font-semibold">Lovely America</p>`,
			`<span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"> 0 </span>`,
			`<button class="w-full border-2 border-gray-800 rounded-lg center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/cookbooks/1" hx-target="#content" hx-push-url="true"> Open </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}
