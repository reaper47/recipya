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
			`<div class="three-dots-container cursor-pointer h-fit justify-self-end hover:text-red-600" _="on click if menuOpen add .hidden to menuOpen end then set $menuOpen to #cookbook-menu-container then set $menuOpen.style.left to (event.pageX - 104) +'px' then set $menuOpen.style.top to (event.pageY + 24)+'px' then set $li to closest <li/> to event.target then if not $li set $li to closest <section/> end then set $deleteOption to #cookbook-menu-delete then set $addButton to #add-cookbook then js const page = $addButton.getAttribute('hx-post').split('?')[1]; $deleteOption.setAttribute('hx-delete', ` + "`/cookbooks/${$li.id.split('-')[1]}?${page}`" + `); $deleteOption.setAttribute('hx-target', ` + "`#${$li.id}`" + `); end then call htmx.process($menuOpen) then toggle .hidden on #cookbook-menu-container">`,
			`<a id="cookbook-menu-share" class="flex p-1" hx-get="/cookbooks/1/share">`,
			`<a id="cookbook-menu-download" class="flex p-1" hx-get="/cookbooks/1/download">`,
			`<a id="cookbook-menu-delete" class="flex p-1" hx-delete="/cookbooks/1" hx-swap="outerHTML" hx-target="closest .cookbook" hx-confirm="Are you sure you want to delete 'Cookbooks'? Its recipes will not be deleted."><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg><div class="pl-1 align-bottom">Delete</div></a>`,
			`<span class="w-fit h-fit text-xs text-center font-medium select-none py-1 px-2 bg-indigo-700 text-white self-end justify-self-end"> 0 </span>`,
			`<button class="w-full border-t center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600 hover:rounded-b-md" hx-get="/cookbooks/2" hx-target="#content" hx-push-url="true"> Open </button>`,
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
			`<button class="w-full border-2 border-gray-800 rounded-lg center hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/cookbooks/1" hx-target="#content" hx-push-url="true"> Open </button>`,
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

func prepareCookbook(srv *server.Server) (*mockFiles, *mockRepository, func()) {
	originalFiles := srv.Files
	originalRepo := srv.Repository

	recipes := models.Recipes{{ID: 1, Name: "Chicken"}}

	files := &mockFiles{}
	repo := &mockRepository{
		CookbooksRegistered: map[int64][]models.Cookbook{
			1: {
				models.Cookbook{ID: 1, Title: "Lovely Canada", Recipes: recipes},
				models.Cookbook{ID: 2, Title: "Lovely Ukraine"},
				models.Cookbook{ID: 3, Title: "Lovely America"},
			},
		},
		RecipesRegistered: map[int64]models.Recipes{1: recipes},
	}
	srv.Files = files
	srv.Repository = repo

	return files, repo, func() {
		srv.Files = originalFiles
		srv.Repository = originalRepo
	}
}
