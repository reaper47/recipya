package server_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"github.com/reaper47/recipya/internal/services"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHandlers_Recipes(t *testing.T) {
	srv := newServerTest()

	const uri = "/recipes"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("user has no recipes", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		want := []string{
			`<title hx-swap-oob="true">Recipes | Recipya</title>`,
			`<p class="pb-2">Your recipe collection looks a bit empty at the moment.</p><p>Why not start adding recipes by clicking the <a class="underline font-semibold cursor-pointer" hx-get="/recipes/add" hx-target="#content" hx-push-url="true">Add recipe</a> button at the top?</p>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("user has recipes", func(t *testing.T) {
		srv.Repository = &mockRepository{
			RecipesRegistered: map[int64]models.Recipes{
				1: {
					{ID: 1, Name: "One", Category: "Lunch", Description: "Recipe one"},
					{ID: 2, Name: "Two", Category: "Soup", Description: "Recipe two"},
					{ID: 3, Name: "Three", Category: "Dinner", Description: "Recipe three"},
				},
			},
		}
		defer func() {
			srv.Repository = &mockRepository{}
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		got := getBodyHTML(rr)
		assertStringsInHTML(t, got, []string{
			`<title hx-swap-oob="true">Recipes | Recipya</title>`,
			`<form class="w-72 flex md:w-96" hx-get="/recipes/search" hx-vals="{"page": 1}" hx-target="#list-recipes" hx-push-url="true" hx-trigger="submit, change target:.sort-option"><div class="w-full"><label class="input input-bordered input-sm flex justify-between px-0 gap-2 z-20"><button type="button" id="search_shortcut" class="pl-2" popovertarget="search_help" _="on click toggle .hidden on #search_help"><svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 self-center" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg></button> <input id="search_recipes" class="w-full" type="search" name="q" placeholder="Search for recipes..." value="" _="on keyup if event.target.value !== '' then remove .md:block from #search_shortcut else add .md:block to #search_shortcut then if (event.key is not 'Delete' and not event.key.startsWith('Arrow')) then send submit to closest <form/> then end end"> <button type="submit" class="px-2 btn btn-sm btn-primary"><svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"></path></svg><span class="sr-only">Search</span></button></label></div><div class="dropdown dropdown-left ml-1"><div tabindex="0" role="button" class="btn btn-sm p-1"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12"></path></svg></div><div tabindex="0" class="dropdown-content z-10 menu menu-sm p-2 shadow bg-base-200 w-52 sm:menu-md prose"><h4>Sort</h4><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Default</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="default" checked></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Name:<br>A to Z</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="a-z"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Name:<br>Z to A</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="z-a"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Date created:<br>Newest to oldest</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="new-old"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Date created:<br>Oldest to newest</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="old-new"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Random</span> <input type="radio" name="sort" class="radio radio-sm sort-option" value="random"></label></div></div></div></form>`,
			`<div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex">`,
			`<img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the One recipe">`,
			`<img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Two recipe">`,
			`<img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Three recipe">`,
		})
		assertStringsNotInHTML(t, got, []string{`Your recipe collection looks a bit empty at the moment.`})
	})
}

func TestHandlers_Recipes_New(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/add"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, uri)

			want := []string{
				`<title hx-swap-oob="true">Add Recipe | Recipya</title>`,
				`<img class="object-cover w-full h-40 rounded-t-xl" src="/static/img/recipes/new/import.webp" alt="Writing on a piece of paper with a traditional pen.">`,
				`<button class="underline" hx-get="/recipes/supported-websites" hx-target="#search-results" onclick="supported_websites_dialog.showModal()">supported</button>`,
				`<dialog id="websites_dialog" class="modal"><div class="modal-box"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="font-bold text-lg">Fetch recipes from websites</h3><form class="py-4" hx-post="/recipes/add/website" hx-swap="none" _="on submit call websites_dialog.close() then set me.querySelector('textarea').value to ''"><div class="grid mb-4"><label class="form-control"><div class="label"><span class="label-text">Enter one or more URLs, each on a new line.</span></div><textarea class="textarea textarea-bordered whitespace-pre-line" placeholder="URL 1URL 2URL 3URL 4etc..." name="urls" rows="5"></textarea></label></div><button class="btn btn-block btn-primary btn-sm">Submit</button></form></div></dialog>`,
				`<dialog id="supported_websites_dialog" class="modal"><div class="modal-box h-2/3"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="mb-1"><label><input type="search" placeholder="Search a website" class="input input-bordered input-sm w-11/12" _="on input show <tbody>tr/> in next <table/> when its textContent.toLowerCase() contains my value.toLowerCase()"></label></h3><div class="overflow-x-auto"><table class="table table-zebra table-sm"><thead><tr class="text-center"><th class="py-1">Number</th><th class="py-1">Website</th></tr></thead> <tbody id="search-results"></tbody></table></div></div></dialog>`,
				`<dialog id="supported_apps_import_dialog" class="modal"><div class="modal-box h-2/3"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="mb-1"><label><input type="search" placeholder="Search an application" class="input input-bordered input-sm w-11/12" _="on input show <tbody>tr/> in next <table/> when its textContent.toLowerCase() contains my value.toLowerCase()"></label></h3><div class="overflow-x-auto"><table class="table table-zebra table-sm"><thead><tr class="text-center"><th class="py-1">Number</th><th class="py-1">Application</th></tr></thead> <tbody id="application-results"></tbody></table></div></div></dialog>`,
				`<dialog id="add_ocr_dialog" class="modal"><div class="modal-box"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="font-bold text-lg">Scan Recipe</h3><form class="py-4" hx-post="/recipes/add/ocr" hx-encoding="multipart/form-data" hx-indicator="#fullscreen-loader" hx-swap="none" _="on submit add_ocr_dialog.close()"><div class="grid mb-4"><label for="add-ocr-files-input" class="text-sm font-medium mb-1">Select your recipe's images ordered by page or a recipe document in the PDF format.</label> <input id="add-ocr-files-input" type="file" name="files" accept=".jpg, .jpeg, .png, .bmp, .tiff, .heif, .pdf" multiple class="p-2 border border-gray-300 rounded-lg shadow focus:ring-2 focus:ring-purple-600 dark:bg-gray-900 dark:border-none"></div><button class="btn btn-block btn-primary btn-sm">Submit</button></form></div></dialog>`,
				`<dialog id="import_recipes_dialog" class="modal"><div class="modal-box"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="font-bold text-lg">Import Recipes</h3><form class="py-4" hx-post="/recipes/add/import" enctype="multipart/form-data" hx-indicator="#fullscreen-loader" hx-swap="none"><div class="grid mb-4"><label for="import-dialog-file" class="text-sm font-semibold mb-1">Choose files in the .json, .txt, .zip or other application format.</label> <input id="import-dialog-file" type="file" name="files" accept=".cml,.crumb,.json,.mxp,.paprikarecipes,.txt,.zip" multiple class="p-2 border border-gray-300 rounded-lg shadow focus:ring-2 focus:ring-purple-600 dark:bg-gray-900 dark:border-none"></div><button type="submit" class="btn btn-block btn-primary btn-sm" onclick="import_recipes_dialog.close()">Submit</button></form></div></dialog>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}

func TestHandlers_Recipes_AddImport(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.Close()

	uri := ts.URL + "/recipes/add/import"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("payload too big", func(t *testing.T) {
		b := bytes.NewBuffer(make([]byte, 130<<20))
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formData, strings.NewReader(b.String()))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 3, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Could not parse the uploaded files.","title":"Form Error"}}`)
	})

	t.Run("error parsing files", func(t *testing.T) {
		contentType, body := createMultipartForm(map[string][]string{"files": {"file1"}})
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 3, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Could not retrieve the files or the directory from the form.","title":"Form Error"}}`)
	})

	t.Run("valid request", func(t *testing.T) {
		contentType, body := createMultipartForm(map[string][]string{"files": {"file1.jpg"}})
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusAccepted)
	})
}

func TestHandlers_Recipes_AddManual(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.Close()

	uri := ts.URL + "/recipes/add/manual"

	repo := &mockRepository{
		categories:        map[int64][]string{1: {"breakfast", "lunch", "dinner"}},
		RecipesRegistered: make(map[int64]models.Recipes),
		UsersRegistered: []models.User{
			{ID: 1, Email: "test@example.com"},
		},
		UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
	}

	resetRepo := func() int {
		repo = &mockRepository{
			categories:        map[int64][]string{1: {"breakfast", "lunch", "dinner"}},
			RecipesRegistered: make(map[int64]models.Recipes),
			UsersRegistered: []models.User{
				{ID: 1, Email: "test@example.com"},
			},
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		return len(repo.RecipesRegistered)
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, uri)

			want := []string{
				`<title hx-swap-oob="true">Add Manual | Recipya</title>`,
				`<form class="card-body" style="padding: 0" enctype="multipart/form-data" hx-post="/recipes/add/manual" hx-indicator="#fullscreen-loader">`,
				`<input required type="text" name="title" placeholder="Title of the recipe*" autocomplete="off" class="input w-full btn-ghost text-center">`,
				`<img src="" alt="Image preview of the recipe." class="object-cover mb-2 w-full max-h-[39rem]"> <span class="grid gap-1 max-w-sm" style="margin: auto auto 0.25rem;"><div class="mr-1"><input type="file" accept="image/*" name="images" class="file-input file-input-sm file-input-bordered w-full max-w-sm" _="on dragover or dragenter halt the event then set the target's style.background to 'lightgray' on dragleave or drop set the target's style.background to '' on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from me.parentElement.parentElement.querySelectorAll('button') then add .hidden to the parentElement of me">`,
				`<input type="number" min="1" name="yield" value="1" class="input input-bordered input-sm w-24 md:w-20 lg:w-24">`,
				`<input type="text" list="categories" name="category" class="input input-bordered input-sm w-48 md:w-36 lg:w-48" placeholder="Breakfast" autocomplete="off"> <datalist id="categories"><option>breakfast</option><option>lunch</option><option>dinner</option></datalist>`,
				`<textarea name="description" placeholder="This Thai curry chicken will make you drool." class="textarea w-full h-full resize-none"></textarea>`,
				`<table class="table table-zebra table-xs md:h-fit"><thead><tr><th>Time</th><th>h:m:s</th></tr></thead> <tbody><tr><td>Prep</td><td><label><input type="text" name="time-preparation" value="00:15:00" class="input input-bordered input-xs max-w-24 html-duration-picker"></label></td></tr><tr><td>Cooking</td><td><label><input type="text" name="time-cooking" value="00:30:00" class="input input-bordered input-xs max-w-24 html-duration-picker"></label></td></tr></tbody></table>`,
				`<table class="table table-zebra table-xs"><thead><tr><th>Nutrition<br>(per 100g)</th><th>Amount</th></tr></thead> <tbody><tr><td>Calories</td><td><label><input type="text" name="calories" autocomplete="off" placeholder="368kcal" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Total carbs</td><td><label><input type="text" name="total-carbohydrates" autocomplete="off" placeholder="35g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Sugars</td><td><label><input type="text" name="sugars" autocomplete="off" placeholder="3g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Protein</td><td><label><input type="text" name="protein" autocomplete="off" placeholder="21g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Total fat</td><td><label><input type="text" name="total-fat" autocomplete="off" placeholder="15g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Saturated fat</td><td><label><input type="text" name="saturated-fat" autocomplete="off" placeholder="1.8g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Cholesterol</td><td><label><input type="text" name="cholesterol" autocomplete="off" placeholder="1.1mg" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Sodium</td><td><label><input type="text" name="sodium" autocomplete="off" placeholder="100mg" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Fiber</td><td><label><input type="text" name="fiber" autocomplete="off" placeholder="8g" class="input input-bordered input-xs max-w-24"></label></td></tr></tbody></table>`,
				`<ol id="tools-list" class="pl-4 list-decimal"><li class="pb-2"><div class="grid grid-flow-col items-center"><label><input type="text" name="tools" placeholder="1 frying pan" class="input input-bordered input-sm w-full" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)"></label><div class="ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: Enter" onclick="addItem(event)">+</button> <button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" _="on click if (closest <ol/>).childElementCount > 1 remove closest <li/> else set input to (closest <li/>).querySelector('input') then set input.value to '' then input.focus()">-</button><div class="inline-block h-4 cursor-move handle ml-2"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li></ol>`,
				`<ol id="ingredients-list" class="pl-4 list-decimal"><li class="pb-2"><div class="grid grid-flow-col items-center"><label><input required type="text" name="ingredients" value="" placeholder="1 cup of chopped onions" class="input input-bordered input-sm w-full" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)"></label><div class="ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: Enter" onclick="addItem(event)">+</button> <button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" _="on click if (closest <ol/>).childElementCount > 1 remove closest <li/> else set input to (closest <li/>).querySelector('input') then set input.value to '' then input.focus()">-</button><div class="inline-block h-4 cursor-move handle ml-2"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li></ol></div><div class="col-span-6 px-6 py-2 border-gray-700 md:rounded-bl-none md:col-span-4"><h2 class="font-semibold text-center pb-2"><span class="underline">Instructions</span> <sup class="text-red-600">*</sup></h2><ol id="instructions-list" class="grid list-decimal"><li class="pt-2 md:pl-0"><div class="flex"><label class="w-11/12"><textarea required name="instructions" rows="3" class="textarea textarea-bordered w-full" placeholder="Mix all ingredients together" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)"></textarea></label><div class="grid ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: CTRL + Enter" onclick="addItem(event)">+</button> <button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" _="on click if (closest <ol/>).childElementCount > 1 remove closest <li/> else set input to (closest <li/>).querySelector('textarea') then set input.value to '' then input.focus()">-</button><div class="h-4 cursor-move handle grid place-content-center"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li></ol>`,
				`<div class="col-span-6 px-6 py-2 border-gray-700 md:rounded-bl-none md:col-span-4"><h2 class="font-semibold text-center pb-2"><span class="underline">Instructions</span> <sup class="text-red-600">*</sup></h2>`,
				`<ol id="instructions-list" class="grid list-decimal"><li class="pt-2 md:pl-0"><div class="flex"><label class="w-11/12"><textarea required name="instructions" rows="3" class="textarea textarea-bordered w-full" placeholder="Mix all ingredients together" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)"></textarea></label><div class="grid ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: CTRL + Enter" onclick="addItem(event)">+</button> <button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" _="on click if (closest <ol/>).childElementCount > 1 remove closest <li/> else set input to (closest <li/>).querySelector('textarea') then set input.value to '' then input.focus()">-</button><div class="h-4 cursor-move handle grid place-content-center"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li></ol>`,
				`<button class="btn btn-primary btn-block btn-sm">Submit</button>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}

	t.Run("missing source defaults to unknown", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		if repo.RecipesRegistered[1][0].URL != "Unknown" {
			t.Fatalf("got source %q; want 'unknown'", repo.RecipesRegistered[1][0].URL)
		}
	})

	t.Run("missing category defaults to uncategorized", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		if repo.RecipesRegistered[1][0].Category != "uncategorized" {
			t.Fatalf("got category %q; want 'uncategorized'", repo.RecipesRegistered[1][0].Category)
		}
	})

	t.Run("missing yield defaults to 1", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		if repo.RecipesRegistered[1][0].Yield != 1 {
			t.Fatalf("got yield %d; want 1", repo.RecipesRegistered[1][0].Yield)
		}
	})

	t.Run("can only be one category", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"category":     {"dinner,lunch"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		if repo.RecipesRegistered[1][0].Category != "dinner" {
			t.Fatalf("got category %s; want dinner", repo.RecipesRegistered[1][0].Category)
		}
	})

	t.Run("subcategories are possible", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"category":     {"beverages:cocktails:vodka"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		if repo.RecipesRegistered[1][0].Category != "beverages:cocktails:vodka" {
			t.Fatalf("got category %s; want beverages:cocktails:vodka", repo.RecipesRegistered[1][0].Category)
		}
	})

	t.Run("submit recipe", func(t *testing.T) {
		_ = resetRepo()

		contentType, body := createMultipartForm(map[string][]string{
			"title":               {"Salsa"},
			"images":              {"eggs.jpg"},
			"category":            {"appetizers"},
			"source":              {"Mommy"},
			"description":         {"The best"},
			"calories":            {"666"},
			"total-carbohydrates": {"31g"},
			"sugars":              {"0.1mg"},
			"protein":             {"5g"},
			"total-fat":           {"0g"},
			"saturated-fat":       {"0g"},
			"cholesterol":         {"256mg"},
			"sodium":              {"777mg"},
			"fiber":               {"2g"},
			"time-preparation":    {"00:15:30"},
			"time-cooking":        {"00:30:15"},
			"tools":               {"wok", "3 pan"},
			"ingredients":         {"ing1", "ing2"},
			"instructions":        {"ins1", "ins2"},
			"yield":               {"4"},
		})
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		id := int64(len(repo.RecipesRegistered))
		gotRecipe := repo.RecipesRegistered[1][id-1]
		want := models.Recipe{
			Category:     "appetizers",
			Description:  "The best",
			ID:           1,
			Images:       gotRecipe.Images,
			Ingredients:  []string{"ing1", "ing2"},
			Instructions: []string{"ins1", "ins2"},
			Keywords:     make([]string, 0),
			Name:         "Salsa",
			Nutrition: models.Nutrition{
				Calories:           "666",
				Cholesterol:        "256mg",
				Fiber:              "2g",
				Protein:            "5g",
				SaturatedFat:       "0g",
				Sodium:             "777mg",
				Sugars:             "0.1mg",
				TotalCarbohydrates: "31g",
				TotalFat:           "0g",
				UnsaturatedFat:     "",
			},
			Times: models.Times{
				Prep:  15*time.Minute + 30*time.Second,
				Cook:  30*time.Minute + 15*time.Second,
				Total: 45*time.Minute + 45*time.Second,
			},
			Tools: []models.Tool{{Name: "wok", Quantity: 1}, {Name: "pan", Quantity: 3}},
			URL:   "Mommy",
			Yield: 4,
		}
		if len(gotRecipe.Images) == 0 {
			t.Fatal("got no images when want images")
		}
		if !cmp.Equal(want, gotRecipe) {
			t.Log(cmp.Diff(want, gotRecipe))
			t.Fail()
		}
		assertHeader(t, rr, "HX-Redirect", "/recipes/"+strconv.FormatInt(id, 10))
	})
}

func TestHandlers_Recipes_AddOCR(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.Close()

	originalIntegrations := srv.Integrations
	originalRepo := srv.Repository
	app.Config.Integrations.AzureDI.Key = "chicken"
	app.Config.Integrations.AzureDI.Endpoint = "kyiv"

	uri := ts.URL + "/recipes/add/ocr"

	sendReq := func(files string) *httptest.ResponseRecorder {
		contentType, body := createMultipartForm(map[string][]string{"files": {files}})
		return sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("warn user when feature disabled", func(t *testing.T) {
		original := app.Config
		app.Config.Integrations.AzureDI.Key = ""
		app.Config.Integrations.AzureDI.Endpoint = ""
		defer func() {
			app.Config = original
		}()

		rr := sendReq("hello.jpg")

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-warning","message":"Please consult the docs to enable OCR.","title":"Feature Disabled"}}`)
	})

	t.Run("files must not be empty", func(t *testing.T) {
		rr := sendReq("")

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Could not retrieve the image from the form.","title":"Form Error"}}`)
	})

	t.Run("processing OCR failed", func(t *testing.T) {
		srv.Integrations = &mockIntegrations{
			processImageOCRFunc: func(_ []io.Reader) (models.Recipes, error) {
				return models.Recipes{}, errors.New("error")
			},
		}
		defer func() {
			srv.Integrations = originalIntegrations
		}()

		rr := sendReq("hello.jpg")

		assertStatus(t, rr.Code, http.StatusAccepted)
		assertWebsocket(t, c, 3, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Could not process OCR.","title":"Integrations Error"}}`)
	})

	t.Run("adding processed recipe failed", func(t *testing.T) {
		srv.Repository = &mockRepository{
			AddRecipesFunc: func(_ models.Recipes, _ int64, _ chan models.Progress) ([]int64, []models.ReportLog, error) {
				return nil, nil, errors.New("oops")
			},
			UsersRegistered:        originalRepo.Users(),
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendReq("hello.jpg")

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Recipes could not be added.","title":"Database Error"}}`
		assertWebsocket(t, c, 3, want)
	})

	t.Run("valid request", func(t *testing.T) {
		repo := &mockRepository{
			RecipesRegistered:      make(map[int64]models.Recipes),
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendReq("hello.jpg")

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"View /recipes/1","background":"alert-info","message":"Recipe scanned and uploaded.","title":"Operation Successful"}}`
		assertWebsocket(t, c, 3, want)
		if len(repo.RecipesRegistered[1]) != 1 && repo.RecipesRegistered[1][0].ID != 1 {
			t.Fatal("expected the recipe to be added")
		}
	})
}

func TestHandlers_Recipes_AddWebsite(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.Close()

	originalRepo := srv.Repository
	originalScraper := srv.Scraper
	originalFiles := srv.Files

	uri := ts.URL + "/recipes/add/website"

	prepare := func() *mockRepository {
		repo := &mockRepository{
			RecipesRegistered:      map[int64]models.Recipes{1: make(models.Recipes, 0)},
			Reports:                map[int64][]models.Report{1: make([]models.Report, 0)},
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		return repo
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("no input", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("websites="))

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("no valid URLs", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=I am a pig\noink oink"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"No valid URLs found.","title":"Request Error"}}`)
	})

	t.Run("add one valid URL from unsupported websites", func(t *testing.T) {
		repo := prepare()
		srv.Files = &mockFiles{}
		srv.Scraper = &mockScraper{
			scraperFunc: func(_ string, _ services.FilesService) (models.RecipeSchema, error) {
				return models.RecipeSchema{}, errors.New("unsupported website")
			},
		}
		defer func() {
			srv.Repository = originalRepo
			srv.Scraper = originalScraper
			srv.Files = originalFiles
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=https://www.example.com"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"View /reports?view=latest","background":"alert-error","message":"Fetching the recipe failed.","title":"Operation Failed"}}`
		assertWebsocket(t, c, 4, want)
		if len(repo.Reports[1]) != 1 {
			t.Fatalf("got reports %v but want one report added", repo.Reports[1])
		}
	})

	t.Run("add one valid URL from supported websites", func(t *testing.T) {
		repo := prepare()
		defer func() {
			srv.Repository = repo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=https://www.example.com"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"View /recipes/1","background":"alert-info","message":"Recipe has been added to your collection.","title":"Operation Successful"}}`
		assertWebsocket(t, c, 4, want)
		if len(repo.Reports[1]) != 1 {
			t.Fatalf("got reports %v but want one report added", repo.Reports[1])
		}
		if len(repo.RecipesRegistered[1]) != 1 {
			t.Fatal("expected 3 recipes")
		}
	})

	t.Run("add duplicates", func(t *testing.T) {
		repo := prepare()
		defer func() {
			srv.Repository = repo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=https://www.example.com\nhttps://www.example.com"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"View /recipes/1","background":"alert-info","message":"Recipe has been added to your collection.","title":"Operation Successful"}}`
		assertWebsocket(t, c, 4, want)
		if len(repo.Reports[1]) != 1 {
			t.Fatalf("got reports %v but want one report added", repo.Reports[1])
		}
		if len(repo.RecipesRegistered[1]) != 1 {
			t.Fatal("expected 3 recipes")
		}
	})

	t.Run("add many valid URLs from supported websites", func(t *testing.T) {
		repo := prepare()
		defer func() {
			srv.Repository = repo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=https://www.example.com\nhttps://www.hello.com\nhttp://helloiam.bob.com\njesus.com"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		assertWebsocket(t, c, 6, `{"type":"toast","fileName":"","data":"","toast":{"action":"View /reports?view=latest","background":"alert-info","message":"Fetched 3 recipes. 0 skipped","title":"Operation Successful"}}`)
		if len(repo.Reports[1]) != 1 {
			t.Fatalf("got reports %v but want one report added", repo.Reports[1])
		}
		if len(repo.RecipesRegistered[1]) != 3 {
			t.Fatal("expected 3 recipes")
		}
	})
}

func TestHandlers_Recipes_Delete(t *testing.T) {
	repo := &mockRepository{
		RecipesRegistered: map[int64]models.Recipes{1: make(models.Recipes, 0)},
		UsersRegistered:   []models.User{{ID: 1, Email: "test@example.com"}},
	}
	srv := newServerTest()
	srv.Repository = repo

	uri := "/recipes"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodDelete, uri+"/5")
	})

	t.Run("cannot delete recipe that does not exist", func(t *testing.T) {
		numRecipesBefore := len(repo.RecipesRegistered)
		defer func() {
			srv.Repository = repo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodDelete, uri+"/5")

		assertStatus(t, rr.Code, http.StatusNoContent)
		if numRecipesBefore != len(repo.RecipesRegistered) {
			t.Fail()
		}
	})

	t.Run("can delete user's recipe", func(t *testing.T) {
		_, _, _ = srv.Repository.AddRecipes(models.Recipes{{ID: 1, Name: "Chicken"}}, 1, nil)

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodDelete, uri+"/1")

		assertStatus(t, rr.Code, http.StatusNoContent)
		assertHeader(t, rr, "HX-Redirect", "/")
	})
}

func TestHandlers_Recipes_Edit(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.Close()

	baseRecipe := models.Recipe{
		Category:     "american",
		CreatedAt:    time.Now(),
		Cuisine:      "indonesian",
		Description:  "A delicious recipe!",
		ID:           1,
		Images:       []uuid.UUID{uuid.New()},
		Ingredients:  []string{"ing1", "ing2", "ing3"},
		Instructions: []string{"ins1", "ins2", "ins3"},
		Keywords:     []string{"chicken", "big", "marinade"},
		Name:         "Chicken Jersey",
		Nutrition: models.Nutrition{
			Calories:           "354",
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
			Prep:  30 * time.Minute,
			Cook:  1 * time.Hour,
			Total: 1*time.Hour + 30*time.Minute,
		},
		Tools:     []models.Tool{{Name: "spoons"}, {Name: "drill"}},
		UpdatedAt: time.Now(),
		URL:       "https://example.com/recipes/yummy",
		Yield:     12,
	}

	xc, _ := srv.Repository.Categories(1)
	repo := &mockRepository{
		categories:        map[int64][]string{1: xc},
		RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}},
	}
	srv.Repository = repo
	originalRepo := srv.Repository

	resetRepo := func() int {
		repo = &mockRepository{
			categories:             map[int64][]string{1: xc},
			RecipesRegistered:      map[int64]models.Recipes{1: {baseRecipe}},
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		return len(repo.RecipesRegistered)
	}

	uri := ts.URL + "/recipes/%d/edit"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, fmt.Sprintf(uri, 5))
		assertMustBeLoggedIn(t, srv, http.MethodPut, fmt.Sprintf(uri, 5))
	})

	t.Run("error fetching recipe", func(t *testing.T) {
		srv.Repository = &mockRepository{
			RecipeFunc: func(id, userID int64) (*models.Recipe, error) {
				return nil, errors.New("oops")
			},
			RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, fmt.Sprintf(uri, 1))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to retrieve recipe.","title":"Database Error"}}`)
	})

	t.Run("successful request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, fmt.Sprintf(uri, 1))

		want := []string{
			`<title hx-swap-oob="true">Edit Chicken Jersey | Recipya</title>`,
			`<input required type="text" name="title" placeholder="Title of the recipe*" autocomplete="off" class="input w-full btn-ghost text-center" value="Chicken Jersey">`,
			`<img alt="Image preview of the recipe." class="object-cover mb-2 w-full max-h-[39rem]" src="/data/images/` + baseRecipe.Images[0].String() + `.jpg"> <span class="grid gap-1 max-w-sm" style="margin: auto auto 0.25rem;"><div class="mr-1"><input type="file" accept="image/*" name="images" class="file-input file-input-sm file-input-bordered w-full max-w-sm" value="/data/images/` + baseRecipe.Images[0].String() + `.jpg" _="on dragover or dragenter halt the event then set the target's style.background to 'lightgray' on dragleave or drop set the target's style.background to '' on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from me.parentElement.parentElement.querySelectorAll('button') then add .hidden to the parentElement of me">`,
			`<input type="text" list="categories" name="category" class="input input-bordered input-sm w-48 md:w-36 lg:w-48" placeholder="Breakfast" autocomplete="off" value="american"> <datalist id="categories"><option>breakfast</option><option>lunch</option><option>dinner</option></datalist>`,
			`<input type="number" min="1" name="yield" class="input input-bordered input-sm w-24 md:w-20 lg:w-24" value="12">`,
			`<input type="text" placeholder="Source" name="source" class="input input-bordered input-sm md:w-28 lg:w-40 xl:w-44" value="https://example.com/recipes/yummy"`,
			`<textarea name="description" placeholder="This Thai curry chicken will make you drool." class="textarea w-full h-full resize-none">A delicious recipe!</textarea>`,
			`<tbody><tr><td>Prep</td><td><label><input type="text" name="time-preparation" class="input input-bordered input-xs max-w-24 html-duration-picker" value="00:30:00"></label></td></tr><tr><td>Cooking</td><td><label><input type="text" name="time-cooking" class="input input-bordered input-xs max-w-24 html-duration-picker" value="01:00:00"></label></td></tr></tbody>`,
			`<tbody><tr><td>Calories</td><td><label><input type="text" name="calories" autocomplete="off" placeholder="368kcal" class="input input-bordered input-xs max-w-24" value="354"></label></td></tr><tr><td>Total carbs</td><td><label><input type="text" name="total-carbohydrates" autocomplete="off" placeholder="35g" class="input input-bordered input-xs max-w-24" value="7g"></label></td></tr><tr><td>Sugars</td><td><label><input type="text" name="sugars" autocomplete="off" placeholder="3g" class="input input-bordered input-xs max-w-24" value="6g"></label></td></tr><tr><td>Protein</td><td><label><input type="text" name="protein" autocomplete="off" placeholder="21g" class="input input-bordered input-xs max-w-24" value="3g"></label></td></tr><tr><td>Total fat</td><td><label><input type="text" name="total-fat" autocomplete="off" placeholder="15g" class="input input-bordered input-xs max-w-24" value="8g"></label></td></tr><tr><td>Saturated fat</td><td><label><input type="text" name="saturated-fat" autocomplete="off" placeholder="1.8g" class="input input-bordered input-xs max-w-24" value="4g"></label></td></tr><tr><td>Cholesterol</td><td><label><input type="text" name="cholesterol" autocomplete="off" placeholder="1.1mg" class="input input-bordered input-xs max-w-24" value="1g"></label></td></tr><tr><td>Sodium</td><td><label><input type="text" name="sodium" autocomplete="off" placeholder="100mg" class="input input-bordered input-xs max-w-24" value="5g"></label></td></tr><tr><td>Fiber</td><td><label><input type="text" name="fiber" autocomplete="off" placeholder="8g" class="input input-bordered input-xs max-w-24" value="2g"></label></td></tr></tbody>`,
			`<input type="text" name="tools" placeholder="1 frying pan" class="input input-bordered input-sm w-full" value="0 spoons" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)">`,
			`<input type="text" name="time-preparation" class="input input-bordered input-xs max-w-24 html-duration-picker" value="00:30:00">`,
			`<input required type="text" name="ingredients" value="ing1" placeholder="1 cup of chopped onions" class="input input-bordered input-sm w-full" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)">`,
			`<textarea required name="instructions" rows="3" class="textarea textarea-bordered w-full" placeholder="Mix all ingredients together" _="on keydown if event.key is 'Enter' halt the event then call addItem(event)">ins1</textarea>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("no updated image", func(t *testing.T) {
		files := &mockFiles{}
		srv.Files = files
		srv.Repository = &mockRepository{RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}}}
		contentType, body := createMultipartForm(map[string][]string{})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		got, _ := srv.Repository.Recipe(baseRecipe.ID, 1)
		if !slices.Equal(got.Images, baseRecipe.Images) {
			t.Fatal("image should not have been updated")
		}
		if files.uploadImageHitCount > 0 {
			t.Fatal("must not have uploaded image")
		}
	})

	t.Run("updated image", func(t *testing.T) {
		files := &mockFiles{}
		srv.Files = files
		srv.Repository = &mockRepository{RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}}}
		contentType, body := createMultipartForm(map[string][]string{"images": {"eggs.jpg"}})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		got, _ := srv.Repository.Recipe(baseRecipe.ID, 1)
		if slices.Equal(got.Images, baseRecipe.Images) {
			t.Fatal("image should have been updated")
		}
		if files.uploadImageHitCount == 0 {
			t.Fatal("must have uploaded image")
		}
	})

	t.Run("update every field", func(t *testing.T) {
		files := &mockFiles{}
		srv.Files = files
		srv.Repository = &mockRepository{RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}}}
		fields := map[string][]string{
			"title":               {"Salsa"},
			"image":               {"jesus.jpg"},
			"category":            {"appetizers"},
			"source":              {"Mommy"},
			"description":         {"The best"},
			"calories":            {"666"},
			"total-carbohydrates": {"31g"},
			"sugars":              {"0.1mg"},
			"protein":             {"5g"},
			"total-fat":           {"24g"},
			"saturated-fat":       {"58g"},
			"cholesterol":         {"256mg"},
			"sodium":              {"777mg"},
			"fiber":               {"2g"},
			"time-preparation":    {"00:15:30"},
			"time-cooking":        {"00:30:15"},
			"ingredients":         {"cheese", "avocado"},
			"instructions":        {"mix", "eat"},
			"yield":               {"4"},
		}
		contentType, body := createMultipartForm(fields)

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		assertHeader(t, rr, "HX-Redirect", "/recipes/1")
		got, _ := srv.Repository.Recipe(baseRecipe.ID, 1)
		want := models.Recipe{
			Category:     "appetizers",
			CreatedAt:    baseRecipe.CreatedAt,
			Description:  "The best",
			ID:           baseRecipe.ID,
			Images:       got.Images,
			Ingredients:  []string{"cheese", "avocado"},
			Instructions: []string{"mix", "eat"},
			Keywords:     []string{"chicken", "big", "marinade"},
			Name:         "Salsa",
			Nutrition: models.Nutrition{
				Calories:           "666",
				Cholesterol:        "256mg",
				Fiber:              "2g",
				Protein:            "5g",
				SaturatedFat:       "58g",
				Sodium:             "777mg",
				Sugars:             "0.1mg",
				TotalCarbohydrates: "31g",
				TotalFat:           "24g",
				UnsaturatedFat:     "",
			},
			Times: models.Times{
				Prep:  15*time.Minute + 30*time.Second,
				Cook:  30*time.Minute + 15*time.Second,
				Total: 45*time.Minute + 45*time.Second,
			},
			Tools:     baseRecipe.Tools,
			UpdatedAt: baseRecipe.UpdatedAt,
			URL:       "Mommy",
			Yield:     4,
		}
		if !cmp.Equal(*got, want) {
			t.Log(cmp.Diff(*got, want))
			t.Fail()
		}
	})

	t.Run("missing source defaults to unknown", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		if repo.RecipesRegistered[1][0].URL != "Unknown" {
			t.Fatalf("got source %q; want 'unknown'", repo.RecipesRegistered[1][0].URL)
		}
	})

	t.Run("missing category defaults to uncategorized", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		if repo.RecipesRegistered[1][0].Category != "uncategorized" {
			t.Fatalf("got category %q; want 'uncategorized'", repo.RecipesRegistered[1][0].Category)
		}
	})

	t.Run("missing yield defaults to 1", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		if repo.RecipesRegistered[1][0].Yield != 1 {
			t.Fatalf("got yield %d; want 1", repo.RecipesRegistered[1][0].Yield)
		}
	})

	t.Run("can only be one category", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"category":     {"dinner,lunch"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		if repo.RecipesRegistered[1][0].Category != "dinner" {
			t.Fatalf("got category %s; want dinner", repo.RecipesRegistered[1][0].Category)
		}
	})

	t.Run("subcategories are possible", func(t *testing.T) {
		_ = resetRepo()
		contentType, body := createMultipartForm(map[string][]string{
			"title":        {"title"},
			"category":     {"beverages:cocktails:vodka"},
			"source":       {"Mommy"},
			"ingredients":  {"ing1"},
			"instructions": {"ins1"},
		})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		if repo.RecipesRegistered[1][0].Category != "beverages:cocktails:vodka" {
			t.Fatalf("got category %s; want beverages:cocktails:vodka", repo.RecipesRegistered[1][0].Category)
		}
	})
}

func TestHandlers_Recipes_Scale(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.Close()

	originalRepo := srv.Repository

	uri := ts.URL + "/recipes/1/scale"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	yieldTestcases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "yield query parameter must be present",
			in:   "",
			want: "No yield in the query.",
		},
		{
			name: "yield query parameter must be greater than zero",
			in:   "-1",
			want: "Yield must be greater than zero.",
		},
		{
			name: "yield query parameter must be greater than zero",
			in:   "0",
			want: "Yield must be greater than zero.",
		},
	}
	for _, tc := range yieldTestcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?yield="+tc.in)

			assertStatus(t, rr.Code, http.StatusBadRequest)
			assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"`+tc.want+`","title":"General Error"}}`)
		})
	}

	t.Run("cannot find recipe in database", func(t *testing.T) {
		srv.Repository = &mockRepository{
			RecipesRegistered: map[int64]models.Recipes{1: make(models.Recipes, 0)},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?yield=6")

		assertStatus(t, rr.Code, http.StatusNotFound)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Recipe not found.","title":"General Error"}}`)
	})

	t.Run("valid request double yield", func(t *testing.T) {
		srv.Repository = &mockRepository{
			RecipesRegistered: map[int64]models.Recipes{
				1: []models.Recipe{
					{
						ID:   1,
						Name: "American Jerky",
						Ingredients: []string{
							"2lb chicken",
							"1/2 cup bread loaf",
							"½ tbsp beef broth",
							"7 1/2 cups flour",
							"2 big apples",
							"Lots of big apples",
							"2.5 slices of bacon",
							"2 1/3 cans of bamboo sticks",
							"1½can of tomato paste",
							"6 ¾ peanut butter jars",
							"7.5mL of whiskey",
							"2 tsp lemon juice",
							"Ground ginger",
							"3 Large or 4 medium ripe Hass avocados",
							"1/4-1/2 teaspoon salt plus more for seasoning",
							"1/2 fresh pineapple, cored and cut into 1 1/2-inch pieces",
							"Un sac de chips de 1kg",
							"Two 15-ounce can Goya beans",
						},
						Instructions: []string{},
						Yield:        4,
					},
				},
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?yield=8")

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">4 lb chicken</span></label>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">1 cup bread loaf</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">1 tbsp beef broth</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">15 cups flour</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">4 big apples</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Lots of big apples</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">5 slices of bacon</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">4 2/3 cans of bamboo sticks</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">3 can of tomato paste</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">13 1/2 peanut butter jars</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">15 mL of whiskey</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">1 1/3 tbsp lemon juice</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Ground ginger</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">6 Large or 8 medium ripe Hass avocados</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">1/2 tsp salt plus more for seasoning</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">1 fresh pineapple, cored and cut into 3-inch pieces</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Un sac de chips de 1kg</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">4 15-ounce can Goya beans</span>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_Categories(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.Close()

	baseRecipes := map[int64]models.Recipes{
		1: {
			{ID: 1, Name: "Chinese Firmware", Category: "breakfast"},
			{ID: 2, Name: "Lovely Canada", Category: "lunch"},
			{ID: 3, Name: "Lovely Ukraine", Category: "dinner"},
			{ID: 4, Name: "Space Disco", Category: "snack"},
			{ID: 4, Name: "Maple Pancakes", Category: "breakfast"},
		},
	}
	srv.Repository = &mockRepository{
		categories:        map[int64][]string{1: {"uncategorized", "chicken"}},
		RecipesRegistered: baseRecipes,
	}
	originalRepo := srv.Repository

	uri := ts.URL + "/recipes/categories"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("delete category updates recipes", func(t *testing.T) {
		d := url.Values{}
		d.Add("category", "breakfast")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri+"?"+d.Encode(), formHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		for _, r := range srv.Repository.RecipesAll(1) {
			if r.Category == "breakfast" {
				t.Fatal("expected recipes with category 'breakfast' to be updated")
			}
		}
	})

	t.Run("failed to delete category", func(t *testing.T) {
		srv.Repository = &mockRepository{
			DeleteCategoryFunc: func(_ string, _ int64) error {
				return errors.New("that's some bad hat harry")
			},
			RecipesRegistered: baseRecipes,
		}
		defer func() {
			srv.Repository = originalRepo
		}()
		d := url.Values{}
		d.Add("category", "breakfast")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri+"?"+d.Encode(), formHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to delete category.","title":"General Error"}}`)
	})

	t.Run("cannot delete uncategorized category", func(t *testing.T) {
		d := url.Values{}
		d.Add("category", "uncategorized")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri+"?"+d.Encode(), formHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to delete category.","title":"General Error"}}`)
	})

	t.Run("add new category", func(t *testing.T) {
		categoriesBefore, _ := srv.Repository.Categories(1)
		d := url.Values{}
		d.Add("category", "uber chicken")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(d.Encode()))

		assertStatus(t, rr.Code, http.StatusCreated)
		categoriesAfter, _ := srv.Repository.Categories(1)
		if len(categoriesBefore) == len(categoriesAfter) {
			t.Fatal("categories should have been updated")
		}
		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<div class="badge badge-outline p-3 pr-0"><form class="inline-flex" hx-delete="/recipes/categories" hx-target="closest <div/>" hx-swap="delete"><input type="hidden" name="category" value="uber chicken"> <span class="select-none">uber chicken</span> <button type="submit" class="btn btn-xs btn-ghost">X</button></form></div>`,
			`<div class="badge badge-outline p-3 pr-0"><form class="inline-flex" hx-post="/recipes/categories" hx-target="closest <div/>" hx-swap="outerHTML"><label class="form-control"><input required type="text" placeholder="New category" class="input input-ghost input-xs w-[16ch] focus:outline-none" name="category" autocomplete="off"></label> <button class="btn btn-xs btn-ghost">&#10003;</button></form></div>`,
		})
	})

	t.Run("cannot add existing category", func(t *testing.T) {
		categoriesBefore, _ := srv.Repository.Categories(1)
		d := url.Values{}
		d.Add("category", "uncategorized")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(d.Encode()))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		categoriesAfter, _ := srv.Repository.Categories(1)
		if len(categoriesBefore) != len(categoriesAfter) {
			t.Fatal("categories should not have been updated")
		}
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to add category.","title":"General Error"}}`)
	})

	t.Run("failed to add category", func(t *testing.T) {
		srv.Repository = &mockRepository{
			AddRecipeCategoryFunc: func(_ string, _ int64) error {
				return errors.New("that's very bad")
			},
		}
		defer func() {
			srv.Repository = originalRepo
		}()
		categoriesBefore, _ := srv.Repository.Categories(1)
		d := url.Values{}
		d.Add("category", "uber chicken")

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(d.Encode()))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		categoriesAfter, _ := srv.Repository.Categories(1)
		if len(categoriesBefore) != len(categoriesAfter) {
			t.Fatal("categories should not have been updated")
		}
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to add category.","title":"General Error"}}`)
	})
}

func TestHandlers_Recipes_Search(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.Close()

	srv.Repository = &mockRepository{
		RecipesRegistered: map[int64]models.Recipes{
			1: {
				{ID: 1, Name: "Chinese Firmware"},
				{ID: 2, Name: "Lovely Canada"},
				{ID: 3, Name: "Lovely Ukraine"},
			},
		},
	}

	uri := ts.URL + "/recipes/search"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("page below 1 returns recipes from first page", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?q=lovely&page=0&mode=name")

		assertStatus(t, rr.Code, http.StatusOK)
		assertStringsInHTML(t, getBodyHTML(rr), []string{
			`<button class="join-item btn btn-disabled">«</button>`,
			`<!-- Left Section --><button aria-current="page" class="join-item btn btn-active">1</button>`,
			`<!-- Right Section --><button class="join-item btn btn-disabled">»</button></div><div class="text-center"><p class="text-sm">Showing <span class="font-medium">1</span> to <span class="font-medium">2</span> of <span id="search-count" class="font-medium">2</span> results</p></div>`,
		})
	})

	t.Run("empty query redirects to recipes index", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?q=&page=0&sort=a-z")

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "HX-Retarget", "#content")
	})

	t.Run("no results", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?q=kool-aid&mode=name")

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
				`<section class="card-side sm:card card-compact card-bordered bg-base-100 shadow-lg indicator w-full"><span class="hidden sm:block"><span class="badge badge-primary select-none cursor-pointer hover:bg-neutral indicator-item indicator-center" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put 'cat:' into #search_recipes.value"></span></span><figure class="relative cursor-pointer" hx-get="/recipes/2" hx-target="#content" hx-push-url="true" hx-trigger="mousedown" hx-swap="innerHTML show:window:top transition:true"><img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Canada recipe"><div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex"><p class="p-2 text-sm"></p></div></figure><div class="card-body justify-between" lang="en"><h2 class="sm:font-semibold sm:w-[25ch] sm:break-words sm:min-h-20">Lovely Canada</h2><div class="sm:hidden"><span class="badge badge-primary select-none cursor-pointer hover:bg-neutral" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put 'cat:' into #search_recipes.value"></span></div><div class="card-actions flex-col-reverse h-fit"><button class="btn btn-block btn-xs btn-outline sm:btn-sm " hx-get="/recipes/2" hx-target="#content" hx-trigger="mousedown" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true">View</button></div></div></section>`,
				`<section class="card-side sm:card card-compact card-bordered bg-base-100 shadow-lg indicator w-full"><span class="hidden sm:block"><span class="badge badge-primary select-none cursor-pointer hover:bg-neutral indicator-item indicator-center" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put 'cat:' into #search_recipes.value"></span></span><figure class="relative cursor-pointer" hx-get="/recipes/3" hx-target="#content" hx-push-url="true" hx-trigger="mousedown" hx-swap="innerHTML show:window:top transition:true"><img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Ukraine recipe"><div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex"><p class="p-2 text-sm"></p></div></figure><div class="card-body justify-between" lang="en"><h2 class="sm:font-semibold sm:w-[25ch] sm:break-words sm:min-h-20">Lovely Ukraine</h2><div class="sm:hidden"><span class="badge badge-primary select-none cursor-pointer hover:bg-neutral" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put 'cat:' into #search_recipes.value"></span></div><div class="card-actions flex-col-reverse h-fit"><button class="btn btn-block btn-xs btn-outline sm:btn-sm " hx-get="/recipes/3" hx-target="#content" hx-trigger="mousedown" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true">View</button></div></div></section>`,
			},
		},
		{
			query: "chi",
			want: []string{
				`<section class="card-side sm:card card-compact card-bordered bg-base-100 shadow-lg indicator w-full"><span class="hidden sm:block"><span class="badge badge-primary select-none cursor-pointer hover:bg-neutral indicator-item indicator-center" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put 'cat:' into #search_recipes.value"></span></span><figure class="relative cursor-pointer" hx-get="/recipes/1" hx-target="#content" hx-push-url="true" hx-trigger="mousedown" hx-swap="innerHTML show:window:top transition:true"><img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Chinese Firmware recipe"><div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex"><p class="p-2 text-sm"></p></div></figure><div class="card-body justify-between" lang="en"><h2 class="sm:font-semibold sm:w-[25ch] sm:break-words sm:min-h-20">Chinese Firmware</h2><div class="sm:hidden"><span class="badge badge-primary select-none cursor-pointer hover:bg-neutral" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put 'cat:' into #search_recipes.value"></span></div><div class="card-actions flex-col-reverse h-fit"><button class="btn btn-block btn-xs btn-outline sm:btn-sm " hx-get="/recipes/1" hx-target="#content" hx-trigger="mousedown" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true">View</button></div></div></section>`,
			},
		},
		{
			query: "lovely",
			want: []string{
				`<section class="card-side sm:card card-compact card-bordered bg-base-100 shadow-lg indicator w-full"><span class="hidden sm:block"><span class="badge badge-primary select-none cursor-pointer hover:bg-neutral indicator-item indicator-center" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put 'cat:' into #search_recipes.value"></span></span><figure class="relative cursor-pointer" hx-get="/recipes/2" hx-target="#content" hx-push-url="true" hx-trigger="mousedown" hx-swap="innerHTML show:window:top transition:true"><img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Canada recipe"><div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex"><p class="p-2 text-sm"></p></div></figure><div class="card-body justify-between" lang="en"><h2 class="sm:font-semibold sm:w-[25ch] sm:break-words sm:min-h-20">Lovely Canada</h2><div class="sm:hidden"><span class="badge badge-primary select-none cursor-pointer hover:bg-neutral" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put 'cat:' into #search_recipes.value"></span></div><div class="card-actions flex-col-reverse h-fit"><button class="btn btn-block btn-xs btn-outline sm:btn-sm " hx-get="/recipes/2" hx-target="#content" hx-trigger="mousedown" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true">View</button></div></div></section>`,
				`<section class="card-side sm:card card-compact card-bordered bg-base-100 shadow-lg indicator w-full"><span class="hidden sm:block"><span class="badge badge-primary select-none cursor-pointer hover:bg-neutral indicator-item indicator-center" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put 'cat:' into #search_recipes.value"></span></span><figure class="relative cursor-pointer" hx-get="/recipes/3" hx-target="#content" hx-push-url="true" hx-trigger="mousedown" hx-swap="innerHTML show:window:top transition:true"><img class="h-28 w-24 object-cover rounded-t-lg sm:h-40 sm:min-w-full sm:w-full" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Ukraine recipe"><div class="hidden absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 items-center justify-center text-white select-none rounded-t-lg sm:flex"><p class="p-2 text-sm"></p></div></figure><div class="card-body justify-between" lang="en"><h2 class="sm:font-semibold sm:w-[25ch] sm:break-words sm:min-h-20">Lovely Ukraine</h2><div class="sm:hidden"><span class="badge badge-primary select-none cursor-pointer hover:bg-neutral" hx-get="/recipes/search" hx-target="#list-recipes" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true" hx-vals="{"q": "cat:"}" _="on click put 'cat:' into #search_recipes.value"></span></div><div class="card-actions flex-col-reverse h-fit"><button class="btn btn-block btn-xs btn-outline sm:btn-sm " hx-get="/recipes/3" hx-target="#content" hx-trigger="mousedown" hx-push-url="true" hx-swap="innerHTML show:window:top transition:true">View</button></div></div></section>`,
			},
		},
	}
	for _, tc := range searches {
		t.Run("results for "+tc.query, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?mode=name&q="+tc.query, formHeader, strings.NewReader("id=1&page=1&q="+tc.query))

			assertStatus(t, rr.Code, http.StatusOK)
			assertStringsInHTML(t, getBodyHTML(rr), tc.want)
		})
	}
}

func TestHandlers_Recipes_Share(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.Close()

	uri := func(id int64) string {
		return fmt.Sprintf("%s/recipes/%d/share", ts.URL, id)
	}

	app.Config.Server.URL = "https://www.recipya.com"
	recipe := models.Recipe{
		Category:     "American",
		Description:  "This is the most delicious recipe!",
		ID:           1,
		Images:       []uuid.UUID{uuid.New()},
		Ingredients:  []string{"Ing1", "Ing2", "Ing3"},
		Instructions: []string{"Ins1", "Ins2", "Ins3"},
		Name:         "Chicken Jersey",
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
	_, _, _ = srv.Repository.AddRecipes(models.Recipes{recipe}, 1, nil)
	link, _ := srv.Repository.AddShareLink(models.Share{RecipeID: 1, CookbookID: -1, UserID: 1})

	t.Run("create valid share link", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodPost, uri(1))

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<input type="url" value="` + strings.TrimPrefix(ts.URL, "http://") + `/r/33320755-82f9-47e5-bb0a-d1b55cbd3f7b" class="input input-bordered w-full" readonly="readonly"></label>`,
			`<button id="copy_button" class="btn btn-neutral" title="Copy to clipboard" onClick="__templ_copyToClipboard`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("create invalid share link", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodPost, uri(10))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to create share link.","title":"General Error"}}`)
	})

	t.Run("access share link anonymous", func(t *testing.T) {
		rr := sendRequestNoBody(srv, http.MethodGet, link)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
			`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d=`,
			`<a href="/auth/login" class="btn btn-ghost">Log In</a> <a href="/auth/register" class="btn btn-ghost">Sign Up</a>`,
			`<span class="text-center pb-2 print:w-full">Chicken Jersey</span>`,
			`<button class="mr-2" title="Print recipe" _="on click print()">`,
			`<img id="output" style="object-fit: cover" alt="Image of the recipe" class="w-full max-h-80 md:max-h-[34rem]" src="/data/images/` + recipe.Images[0].String() + `.jpg">`,
			`<div class="badge badge-primary badge-outline">American</div>`,
			`<p class="text-sm text-center">2 servings</p>`,
			`<a class="btn btn-sm btn-outline no-underline print:hidden" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank">Source</a><p class="hidden print:block print:whitespace-nowrap print:overflow-hidden print:text-ellipsis print:max-w-xs">https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/</p>`,
			`<textarea class="textarea w-full h-full resize-none" readonly>This is the most delicious recipe!</textarea>`,
			`<p class="text-xs">Per 100g: calories 500; total carbohydrates 7g; sugar 6g; protein 3g; total fat 8g; saturated fat 4g; cholesterol 1g; fiber 2g</p>`,
			`<table class="table table-zebra table-xs md:h-fit"><thead><tr><th>Time</th><th>h:m:s</th></tr></thead> <tbody><tr><td>Prep:</td><td><time datetime="PT05M">5m</time></td></tr><tr><td>Cooking:</td><td><time datetime="PT1H05M">1h05m</time></td></tr><tr><td>Total:</td><td><time datetime="PT1H10M">1h10m</time></td></tr></tbody></table>`,
			`<table class="table table-zebra table-xs print:hidden"><thead><tr><th>Nutrition (per 100g)</th><th>Amount</th></tr></thead> <tbody><tr><td>Calories:</td><td>500</td></tr><tr><td>Total carbs:</td><td>7g</td></tr><tr><td>Sugars:</td><td>6g</td></tr><tr><td>Protein:</td><td>3g</td></tr><tr><td>Total fat:</td><td>8g</td></tr><tr><td>Saturated fat:</td><td>4g</td></tr><tr><td>Cholesterol:</td><td>1g</td></tr><tr><td>Sodium:</td><td>5g</td></tr><tr><td>Fiber:</td><td>2g</td></tr></tbody></table>`,
			`<div id="ingredients-instructions-container" class="grid text-sm md:grid-flow-col md:col-span-6"><div class="col-span-6 border-gray-700 px-4 py-2 border-y md:col-span-2 md:border-r md:border-y-0 print:hidden"><h2 class="font-semibold text-center underline">Ingredients</h2><ul><li class="form-control hover:bg-gray-100 dark:hover:bg-gray-700"><label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Ing1</span></label></li><li class="form-control hover:bg-gray-100 dark:hover:bg-gray-700"><label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Ing2</span></label></li><li class="form-control hover:bg-gray-100 dark:hover:bg-gray-700"><label class="label justify-start"><input type="checkbox" class="checkbox"> <span class="label-text pl-2">Ing3</span></label></li></ul></div><div class="col-span-6 px-8 py-2 border-gray-700 md:rounded-bl-none md:col-span-4 print:hidden"><h2 class="font-semibold text-center underline md:pb-2">Instructions</h2><ol class="grid list-decimal"><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins1</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins2</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins3</span></li></ol></div></div><div class="hidden print:grid col-span-6 ml-2 my-1"><h1 class="text-sm print:mb-1"><b>Ingredients</b></h1><ol class="col-span-6 w-full print:mb-2" style="column-count: 1"><li class="text-sm"><label><input type="checkbox"></label> <span class="pl-2">Ing1</span></li><li class="text-sm"><label><input type="checkbox"></label> <span class="pl-2">Ing2</span></li><li class="text-sm"><label><input type="checkbox"></label> <span class="pl-2">Ing3</span></li></ol></div><div class="hidden col-span-5 overflow-visible print:inline"><h1 class="text-sm print:ml-2 print:mb-1"><b>Instructions</b></h1><ol class="col-span-6 list-decimal w-full ml-6"><li class="print:mr-4"><span class="text-sm whitespace-pre-line">Ins1</span></li><li class="print:mr-4"><span class="text-sm whitespace-pre-line">Ins2</span></li><li class="print:mr-4"><span class="text-sm whitespace-pre-line">Ins3</span></li></ol></div></div>`,
			`<h2 class="font-semibold text-center underline md:pb-2">Instructions</h2><ol class="grid list-decimal"><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins1</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins2</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins3</span></li></ol>`,
		}
		notWant := []string{
			`id="share-dialog"`,
		}
		body := getBodyHTML(rr)
		assertStringsInHTML(t, body, want)
		assertStringsNotInHTML(t, body, notWant)
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
	}{
		{name: "other user Hx-Request", sendFunc: sendHxRequestAsLoggedInOtherNoBody},
		{name: "other user no Hx-Request", sendFunc: sendRequestAsLoggedInOtherNoBody},
	}
	for _, tc := range testcases {
		t.Run("access share link logged in "+tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, link)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
				`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d=`,
				`<button class="mr-2" title="Add recipe to collection" hx-get="/recipes/1/share/add" hx-push-url="true">`,
			}
			notWant := []string{
				`id="share-dialog"`,
				`title="Share recipe"`,
			}
			body := getBodyHTML(rr)
			assertStringsInHTML(t, body, want)
			assertStringsNotInHTML(t, body, notWant)
		})
	}

	testcases2 := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
	}{
		{name: "host user Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody},
		{name: "host user no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody},
	}
	for _, tc := range testcases2 {
		t.Run("access share link logged "+tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, link)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
				`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d=`,
				`<button class="mr-2" title="Print recipe" _="on click print()">`,
				`<button class="mr-2" hx-delete="/recipes/1" hx-swap="none" title="Delete recipe" hx-confirm="Are you sure you wish to delete this recipe?" hx-indicator="#fullscreen-loader">`,
			}
			notWant := []string{
				`id="share-dialog-result`,
				`title="Add recipe to collection`,
				`title="Share recipe`,
				`title="Add recipe to collection`,
			}
			body := getBodyHTML(rr)
			assertStringsInHTML(t, body, want)
			assertStringsNotInHTML(t, body, notWant)
		})
	}
}

func TestHandlers_Recipes_ShareAdd(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.Close()

	uri := func(id int64) string {
		return fmt.Sprintf("%s/recipes/%d/share/add", ts.URL, id)
	}

	originalRepo := srv.Repository

	app.Config.Server.URL = "https://www.recipya.com"
	recipe := models.Recipe{
		Category:     "American",
		Description:  "This is the most delicious recipe!",
		ID:           1,
		Images:       []uuid.UUID{uuid.New()},
		Ingredients:  []string{"Ing1", "Ing2", "Ing3"},
		Instructions: []string{"Ins1", "Ins2", "Ins3"},
		Name:         "Chicken Jersey",
		Times: models.Times{
			Prep:  5 * time.Minute,
			Cook:  1*time.Hour + 5*time.Minute,
			Total: 1*time.Hour + 10*time.Minute,
		},
		URL:   "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/",
		Yield: 2,
	}
	_, _, _ = srv.Repository.AddRecipes(models.Recipes{recipe}, 1, nil)

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri(1))
	})

	t.Run("recipe ID in path must be positive", func(t *testing.T) {
		rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri(-1))

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
		isHxReq  bool
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody, isHxReq: true},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody, isHxReq: false},
	}
	for _, tc := range testcases {
		t.Run("failed to add shared recipe "+tc.name, func(t *testing.T) {
			srv.Repository = &mockRepository{
				AddShareRecipeFunc: func(_, _ int64) (int64, error) {
					return 0, errors.New("failed to add shared recipe")
				},
			}
			defer func() {
				srv.Repository = originalRepo
			}()

			rr := tc.sendFunc(srv, http.MethodGet, uri(1))

			if tc.isHxReq {
				assertStatus(t, rr.Code, http.StatusInternalServerError)
				assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to add shared recipe to user's collection.","title":"General Error"}}`)
			} else {
				assertStatus(t, rr.Code, http.StatusSeeOther)
				assertHeader(t, rr, "Location", "/recipes")
			}
		})
	}

	testcases2 := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
		isHxReq  bool
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody, isHxReq: true},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody, isHxReq: false},
	}
	for _, tc := range testcases2 {
		t.Run("recipe exists in collection"+tc.name, func(t *testing.T) {
			_, _, _ = srv.Repository.AddRecipes(models.Recipes{recipe}, 1, nil)
			defer func() {
				srv.Repository = originalRepo
			}()

			rr := tc.sendFunc(srv, http.MethodGet, uri(1))

			if tc.isHxReq {
				assertStatus(t, rr.Code, http.StatusInternalServerError)
				assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to add shared recipe to user's collection.","title":"General Error"}}`)
			} else {
				assertStatus(t, rr.Code, http.StatusSeeOther)
				assertHeader(t, rr, "Location", "/recipes")
			}
		})
	}

	testcases3 := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
		isHxReq  bool
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody, isHxReq: true},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody, isHxReq: false},
	}
	for _, tc := range testcases3 {
		t.Run("valid request"+tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, uri(8))

			if tc.isHxReq {
				assertStatus(t, rr.Code, http.StatusOK)
				assertHeader(t, rr, "HX-Redirect", "/recipes/2")
			} else {
				assertStatus(t, rr.Code, http.StatusSeeOther)
				assertHeader(t, rr, "Location", "/recipes/2")
			}
		})
	}
}

func TestHandlers_Recipes_SupportedApplications(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/supported-applications"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("returns list of applications to logged in user", func(t *testing.T) {
		var want []string
		applications := [][]string{
			{"AccuChef", "https://www.accuchef.com"},
			{"ChefTap", "https://cheftap.com"},
			{"Crouton", "https://crouton.app"},
			{"Easy Recipe Deluxe", "https://easy-recipe-deluxe.software.informer.com"},
			{"Kalorio", "https://www.kalorio.de"},
			{"MasterCook", "https://www.mastercook.com"},
			{"Paprika", "https://www.paprikaapp.com"},
			{"Recipe Keeper", "https://recipekeeperonline.com"},
			{"RecipeSage", "https://recipesage.com"},
			{"Saffron", "https://www.mysaffronapp.com"},
		}
		for i, a := range applications {
			want = append(want, `<tr class="border text-center"><td class="border dark:border-gray-800">`+strconv.Itoa(i+1)+`</td><td class="border py-1 dark:border-gray-800"><a class="underline" href="`+a[1]+`" target="_blank">`+a[0]+`</a></td></tr>`)
		}

		rr := sendRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "Content-Type", "text/html")
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_SupportedWebsites(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/supported-websites"

	website1 := `<tr class="border text-center"><td class="border dark:border-gray-800">1</td><td class="border py-1 dark:border-gray-800"><a class="underline" href="https://101cookbooks.com" target="_blank">101cookbooks.com</a></td></tr>`
	website2 := `<tr class="border text-center"><td class="border dark:border-gray-800">2</td><td class="border py-1 dark:border-gray-800"><a class="underline" href="https://www.afghankitchenrecipes.com" target="_blank">afghankitchenrecipes.com</a></td></tr>`

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("returns list of websites to logged in user", func(t *testing.T) {
		rr := sendRequestAsLoggedInNoBody(srv, http.MethodGet, uri)

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "Content-Type", "text/html")
		want := []string{website1, website2}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_View(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri+"/999")
	})

	t.Run("recipe is not in user collection", func(t *testing.T) {
		rr := sendRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"/999")

		assertStatus(t, rr.Code, http.StatusNotFound)
		want := []string{
			`<title hx-swap-oob="true">Page Not Found | Recipya</title>`,
			"Page Not Found",
			"The page you requested to view is not found. Please go back to the main page.",
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string) *httptest.ResponseRecorder
	}{
		{name: "logged in Hx-Request", sendFunc: sendHxRequestAsLoggedInNoBody},
		{name: "logged in no Hx-Request", sendFunc: sendRequestAsLoggedInNoBody},
	}
	for _, tc := range testcases {
		t.Run("recipe is in user's collection when "+tc.name, func(t *testing.T) {
			image, _ := uuid.Parse("e81ba735-a4af-4c66-8c17-2f2ccc1b1a95")
			r := models.Recipe{
				Category:     "American",
				Description:  "This is the most delicious recipe!",
				ID:           1,
				Images:       []uuid.UUID{image},
				Ingredients:  []string{"Ing1", "Ing2", "Ing3"},
				Instructions: []string{"Ins1", "Ins2", "Ins3"},
				Name:         "Chicken Jersey",
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
			_, _, _ = srv.Repository.AddRecipes(models.Recipes{r}, 1, nil)

			rr := tc.sendFunc(srv, http.MethodGet, uri+"/1")

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + r.Name + " | Recipya</title>",
				`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d='M 12.276 18.55 v -0.748 a 4.79 4.79 0 0 1 1.463 -3.458 a 5.763 5.763 0 0 0 1.804 -4.21 a 5.821 5.821 0 0 0 -6.475 -5.778 c -2.779 0.307 -4.99 2.65 -5.146 5.448 a 5.82 5.82 0 0 0 1.757 4.503 a 4.906 4.906 0 0 1 1.5 3.495 v 0.747 a 1.44 1.44 0 0 0 1.44 1.439 h 2.218 a 1.44 1.44 0 0 0 1.44 -1.439 z m -1.058 0 c 0 0.209 -0.17 0.38 -0.38 0.38 h -2.22 c -0.21 0 -0.38 -0.171 -0.38 -0.38 v -0.748 c 0 -1.58 -0.664 -3.13 -1.822 -4.254 A 4.762 4.762 0 0 1 4.98 9.863 c 0.127 -2.289 1.935 -4.204 4.205 -4.455 a 4.762 4.762 0 0 1 5.3 4.727 a 4.714 4.714 0 0 1 -1.474 3.443 a 5.853 5.853 0 0 0 -1.791 4.225 v 0.746 z M 11.45 20.51 H 8.006 a 0.397 0.397 0 1 0 0 0.795 h 3.444 a 0.397 0.397 0 1 0 0 -0.794 z M 11.847 22.162 a 0.397 0.397 0 0 0 -0.397 -0.397 H 8.006 a 0.397 0.397 0 1 0 0 0.794 h 3.444 c 0.22 0 0.397 -0.178 0.397 -0.397 z z z z z z z z M 10.986 23.416 H 8.867 a 0.397 0.397 0 1 0 0 0.794 h 1.722 c 0.22 0 0.397 -0.178 0.397 -0.397 z' to #icon-bulb else call initWakeLock() then add @d='M12.276 18.55v-.748a4.79 4.79 0 0 1 1.463-3.458 5.763 5.763 0 0 0 1.804-4.21 5.821 5.821 0 0 0-6.475-5.778c-2.779.307-4.99 2.65-5.146 5.448a5.82 5.82 0 0 0 1.757 4.503 4.906 4.906 0 0 1 1.5 3.495v.747a1.44 1.44 0 0 0 1.44 1.439h2.218a1.44 1.44 0 0 0 1.44-1.439zm-1.058 0c0 .209-.17.38-.38.38h-2.22c-.21 0-.38-.171-.38-.38v-.748c0-1.58-.664-3.13-1.822-4.254A4.762 4.762 0 0 1 4.98 9.863c.127-2.289 1.935-4.204 4.205-4.455a4.762 4.762 0 0 1 5.3 4.727 4.714 4.714 0 0 1-1.474 3.443 5.853 5.853 0 0 0-1.791 4.225v.746zM11.45 20.51H8.006a.397.397 0 1 0 0 .795h3.444a.397.397 0 1 0 0-.794zM11.847 22.162a.397.397 0 0 0-.397-.397H8.006a.397.397 0 1 0 0 .794h3.444c.22 0 .397-.178.397-.397zM.397 10.125h2.287a.397.397 0 1 0 0-.794H.397a.397.397 0 1 0 0 .794zM19.456 9.728a.397.397 0 0 0-.397-.397h-2.287a.397.397 0 1 0 0 .794h2.287c.22 0 .397-.178.397-.397zM9.331.397v2.287a.397.397 0 1 0 .794 0V.397a.397.397 0 1 0-.794 0zM16.045 2.85 14.43 4.465a.397.397 0 1 0 .561.561l1.617-1.617a.397.397 0 1 0-.562-.56zM5.027 14.429a.397.397 0 0 0-.56 0l-1.618 1.616a.397.397 0 1 0 .562.562l1.617-1.617a.397.397 0 0 0 0-.561zM4.466 5.027a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.56L3.41 2.848a.397.397 0 1 0-.561.561zM16.045 16.607a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.562L14.99 14.43a.397.397 0 1 0-.561.56zM10.986 23.416a.397.397 0 0 0-.397-.397H8.867a.397.397 0 1 0 0 .794h1.722c.22 0 .397-.178.397-.397z' to #icon-bulb end"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor"><path id="icon-bulb" d="M12.276 18.55v-.748a4.79 4.79 0 0 1 1.463-3.458 5.763 5.763 0 0 0 1.804-4.21 5.821 5.821 0 0 0-6.475-5.778c-2.779.307-4.99 2.65-5.146 5.448a5.82 5.82 0 0 0 1.757 4.503 4.906 4.906 0 0 1 1.5 3.495v.747a1.44 1.44 0 0 0 1.44 1.439h2.218a1.44 1.44 0 0 0 1.44-1.439zm-1.058 0c0 .209-.17.38-.38.38h-2.22c-.21 0-.38-.171-.38-.38v-.748c0-1.58-.664-3.13-1.822-4.254A4.762 4.762 0 0 1 4.98 9.863c.127-2.289 1.935-4.204 4.205-4.455a4.762 4.762 0 0 1 5.3 4.727 4.714 4.714 0 0 1-1.474 3.443 5.853 5.853 0 0 0-1.791 4.225v.746zM11.45 20.51H8.006a.397.397 0 1 0 0 .795h3.444a.397.397 0 1 0 0-.794zM11.847 22.162a.397.397 0 0 0-.397-.397H8.006a.397.397 0 1 0 0 .794h3.444c.22 0 .397-.178.397-.397zM.397 10.125h2.287a.397.397 0 1 0 0-.794H.397a.397.397 0 1 0 0 .794zM19.456 9.728a.397.397 0 0 0-.397-.397h-2.287a.397.397 0 1 0 0 .794h2.287c.22 0 .397-.178.397-.397zM9.331.397v2.287a.397.397 0 1 0 .794 0V.397a.397.397 0 1 0-.794 0zM16.045 2.85 14.43 4.465a.397.397 0 1 0 .561.561l1.617-1.617a.397.397 0 1 0-.562-.56zM5.027 14.429a.397.397 0 0 0-.56 0l-1.618 1.616a.397.397 0 1 0 .562.562l1.617-1.617a.397.397 0 0 0 0-.561zM4.466 5.027a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.56L3.41 2.848a.397.397 0 1 0-.561.561zM16.045 16.607a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.562L14.99 14.43a.397.397 0 1 0-.561.56zM10.986 23.416a.397.397 0 0 0-.397-.397H8.867a.397.397 0 1 0 0 .794h1.722c.22 0 .397-.178.397-.397z"></path></svg></button>`,
				`<button class="ml-2" title="Edit recipe" hx-get="/recipes/1/edit" hx-push-url="true" hx-target="#content"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg></button>`,
				`<span class="text-center pb-2 print:w-full">Chicken Jersey</span>`,
				`<button class="mr-2" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call share_dialog.showModal()">`,
				`<button class="mr-2" title="Print recipe" _="on click print()">`,
				`<button class="mr-2" hx-delete="/recipes/1" hx-swap="none" title="Delete recipe" hx-confirm="Are you sure you wish to delete this recipe?" hx-indicator="#fullscreen-loader">`,
				`<img id="output" style="object-fit: cover" alt="Image of the recipe" class="w-full max-h-80 md:max-h-[34rem]" src="/data/images/e81ba735-a4af-4c66-8c17-2f2ccc1b1a95.jpg">`,
				`<div class="badge badge-primary badge-outline">American</div>`,
				`<button class="mr-2" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call share_dialog.showModal()"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor">`,
				`<form autocomplete="off" _="on submit halt the event" class="print:hidden"><label class="form-control w-full"><div class="label p-0"><span class="label-text">Servings</span></div><input id="yield" type="number" min="1" name="yield" value="2" class="input input-bordered input-sm w-24" hx-get="/recipes/1/scale" hx-trigger="input" hx-target="#ingredients-instructions-container"></label></form>`,
				`<a class="btn btn-sm btn-outline no-underline print:hidden" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank">Source</a>`,
				`<textarea class="textarea w-full h-full resize-none" readonly>This is the most delicious recipe!</textarea>`,
				`<p class="text-xs">Per 100g: calories 500; total carbohydrates 7g; sugar 6g; protein 3g; total fat 8g; saturated fat 4g; cholesterol 1g; fiber 2g</p>`,
				`<table class="table table-zebra table-xs md:h-fit"><thead><tr><th>Time</th><th>h:m:s</th></tr></thead> <tbody><tr><td>Prep:</td><td><time datetime="PT05M">5m</time></td></tr><tr><td>Cooking:</td><td><time datetime="PT1H05M">1h05m</time></td></tr><tr><td>Total:</td><td><time datetime="PT1H10M">1h10m</time></td></tr></tbody></table>`,
				`<table class="table table-zebra table-xs print:hidden"><thead><tr><th>Nutrition (per 100g)</th><th>Amount</th></tr></thead> <tbody><tr><td>Calories:</td><td>500</td></tr><tr><td>Total carbs:</td><td>7g</td></tr><tr><td>Sugars:</td><td>6g</td></tr><tr><td>Protein:</td><td>3g</td></tr><tr><td>Total fat:</td><td>8g</td></tr><tr><td>Saturated fat:</td><td>4g</td></tr><tr><td>Cholesterol:</td><td>1g</td></tr><tr><td>Sodium:</td><td>5g</td></tr><tr><td>Fiber:</td><td>2g</td></tr></tbody></table>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}
