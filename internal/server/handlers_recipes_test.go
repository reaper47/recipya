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
	"maps"
	"net/http"
	"net/http/httptest"
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
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

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

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		got := getBodyHTML(rr)
		want := []string{
			`<title hx-swap-oob="true">Recipes | Recipya</title>`,
			`<form class="w-72 flex md:w-96" hx-post="/recipes/search" hx-vals="{"page": 1}" hx-target="#list-recipes" _="on submit if #search-recipes.value is not '' then add .hidden to #pagination else remove .hidden from #pagination end"><div class="relative w-full"><label><input type="search" id="search-recipes" name="q" class="input input-bordered input-sm w-full z-20" placeholder="Search for recipes..." _="on keyup if event.target.value !== '' then remove .md:block from #search-shortcut else add .md:block to #search-shortcut end then if event.key === 'Backspace' and event.target.value === '' then send submit to closest <form/> end"></label><kbd id="search-shortcut" class="hidden absolute top-1 right-12 font-sans font-semibold select-none dark:text-slate-500 md:block"><abbr title="Control" class="no-underline text-slate-300 dark:text-slate-500">Ctrl </abbr> /</kbd><button type="submit" class="absolute top-0 right-0 px-2 btn btn-sm btn-primary"><svg class="w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20"><path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"></path></svg><span class="sr-only">Search</span></button></div><details class="dropdown dropdown-left"><summary class="btn btn-sm ml-1"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75a4.5 4.5 0 0 1-4.884 4.484c-1.076-.091-2.264.071-2.95.904l-7.152 8.684a2.548 2.548 0 1 1-3.586-3.586l8.684-7.152c.833-.686.995-1.874.904-2.95a4.5 4.5 0 0 1 6.336-4.486l-3.276 3.276a3.004 3.004 0 0 0 2.25 2.25l3.276-3.276c.256.565.398 1.192.398 1.852Z"></path><path stroke-linecap="round" stroke-linejoin="round" d="M4.867 19.125h.008v.008h-.008v-.008Z"></path></svg></summary><div tabindex="0" class="dropdown-content z-10 menu menu-sm p-2 shadow bg-base-200 w-52 sm:menu-md prose"><h4>Search Method</h4><div class="form-control"><label class="label cursor-pointer"><span class="label-text">By name</span><input type="radio" name="search-method" class="radio radio-sm" checked value="name"></label></div><div class="form-control"><label class="label cursor-pointer"><span class="label-text">Full search</span><input type="radio" name="search-method" class="radio radio-sm" value="full"></label></div></div></details></form>`,
			`<img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the One recipe">`,
			`<div class="absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 flex items-center justify-center text-white select-none rounded-t-lg">`,
			`<img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the Two recipe">`,
			`<img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the Three recipe">`,
		}
		assertStringsInHTML(t, got, want)
		notWant := []string{`Your recipe collection looks a bit empty at the moment.`}
		assertStringsNotInHTML(t, got, notWant)
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
		sendFunc func(server *server.Server, target, uri string, contentType header, body *strings.Reader) *httptest.ResponseRecorder
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedIn},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedIn},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, uri, noHeader, nil)

			want := []string{
				`<title hx-swap-oob="true">Add Recipe | Recipya</title>`,
				`<img class="object-cover w-full h-40 rounded-t-xl" src="/static/img/recipes/new/import.webp" alt="Writing on a piece of paper with a traditional pen.">`,
				`<button class="underline" hx-get="/recipes/supported-websites" hx-target="#search-results" onclick="supported_websites_dialog.showModal()">supported</button>`,
				`<dialog id="websites_dialog" class="modal"><div class="modal-box"><form method="dialog"><button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button></form><h3 class="font-bold text-lg">Fetch recipes from websites</h3><form class="py-4" hx-post="/recipes/add/website" hx-indicator="#fullscreen-loader" hx-swap="none"><div class="grid mb-4"><label class="form-control"><div class="label"><span class="label-text">Enter one or more URLs, each on a new line.</span></div><textarea class="textarea textarea-bordered h-24 whitespace-pre-line" placeholder="URL 1URL 2URL 3URL 4etc..." name="urls" rows="10"></textarea></label></div><button type="submit" class="w-full p-2 font-semibold text-white bg-blue-500 border rounded-lg hover:bg-blue-800 dark:border-gray-800" onclick="websites_dialog.close()">Submit</button></form></div></dialog>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}

func TestHandlers_Recipes_AddImport(t *testing.T) {
	srv, ts, c := createWSServer()
	defer func() {
		_ = c.Close()
	}()

	originalBrokers := maps.Clone(srv.Brokers)

	uri := ts.URL + "/recipes/add/import"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("no ws connection", func(t *testing.T) {
		srv.Brokers = map[int64]*models.Broker{}
		defer func() {
			srv.Brokers = originalBrokers
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-warning\",\"message\":\"Connection lost. Please reload page.\",\"title\":\"Websocket\"}"}`)
	})

	t.Run("payload too big", func(t *testing.T) {
		b := bytes.NewBuffer(make([]byte, 130<<20))
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formData, strings.NewReader(b.String()))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Could not parse the uploaded files.\",\"title\":\"Form Error\"}"}`)
		want := `<div id="ws-notification-container" class="z-20 fixed bottom-0 right-0 p-6 cursor-default hidden"> <div class="bg-blue-500 text-white px-4 py-2 rounded shadow-md"> <p class="font-medium text-center pb-1"></p> <div id="export-progress"><progress max="100" value="100.000000"></progress></div> </div> </div>`
		assertWebsocket(t, c, 2, want)
	})

	t.Run("error parsing files", func(t *testing.T) {
		contentType, body := createMultipartForm(map[string]string{"files": "file1"})
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Could not retrieve the files or the directory from the form.\",\"title\":\"Form Error\"}"}`)
		want := `<div id="ws-notification-container" class="z-20 fixed bottom-0 right-0 p-6 cursor-default hidden"> <div class="bg-blue-500 text-white px-4 py-2 rounded shadow-md"> <p class="font-medium text-center pb-1"></p> <div id="export-progress"><progress max="100" value="100.000000"></progress></div> </div> </div>`
		assertWebsocket(t, c, 2, want)
	})

	t.Run("valid request", func(t *testing.T) {
		contentType, body := createMultipartForm(map[string]string{"files": "file1.jpg"})
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusAccepted)
	})
}

func TestHandlers_Recipes_AddManual(t *testing.T) {
	srv := newServerTest()
	repo := &mockRepository{}
	srv.Repository = repo

	uri := "/recipes/add/manual"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	testcases := []struct {
		name     string
		sendFunc func(server *server.Server, target, uri string, contentType header, body *strings.Reader) *httptest.ResponseRecorder
	}{
		{name: "is logged in Hx-Request", sendFunc: sendHxRequestAsLoggedIn},
		{name: "is logged in no Hx-Request", sendFunc: sendRequestAsLoggedIn},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, uri, noHeader, nil)

			want := []string{
				`<title hx-swap-oob="true">Add Manual | Recipya</title>`,
				`<form class="card-body" style="padding: 0" enctype="multipart/form-data" hx-post="/recipes/add/manual" hx-indicator="#fullscreen-loader">`,
				`<input required type="text" name="title" placeholder="Title of the recipe*" autocomplete="off" class="input w-full btn-ghost text-center">`,
				`<img src="" alt="Image preview of the recipe." class="object-cover"><span><input type="file" accept="image/*" name="image" class="file-input file-input-sm w-full max-w-xs" _="on dragover or dragenter halt the event then set the target's style.background to 'lightgray' on dragleave or drop set the target's style.background to '' on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from next <button/>">`,
				`<input type="number" min="1" name="yield" value="1" class="input input-bordered input-sm w-24 md:w-20 lg:w-24">`,
				`<input required type="text" list="categories" name="category" class="input input-bordered input-sm w-48 md:w-36 lg:w-48" placeholder="Breakfast"><datalist id="categories"><option>breakfast</option><option>lunch</option><option>dinner</option></datalist>`,
				`<textarea required name="description" placeholder="This Thai curry chicken will make you drool." class="textarea w-full h-full resize-none"></textarea>`,
				`<table class="table table-zebra table-xs md:h-fit"><thead><tr><th>Time</th><th>h:m:s</th></tr></thead><tbody><tr><td>Prep</td><td><label><input type="text" name="time-preparation" value="00:15:00" class="input input-bordered input-xs max-w-24 html-duration-picker"></label></td></tr><tr><td>Cooking</td><td><label><input type="text" name="time-cooking" value="00:30:00" class="input input-bordered input-xs max-w-24 html-duration-picker"></label></td></tr></tbody></table></div><table class="table table-zebra table-xs"><thead><tr><th>Nutrition<br>(per 100g)</th><th>Amount</th></tr></thead><tbody><tr><td>Calories</td><td><label><input type="text" name="calories" autocomplete="off" placeholder="368kcal" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Total carbs</td><td><label><input type="text" name="total-carbohydrates" autocomplete="off" placeholder="35g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Sugars</td><td><label><input type="text" name="sugars" autocomplete="off" placeholder="3g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Protein</td><td><label><input type="text" name="protein" autocomplete="off" placeholder="21g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Total fat</td><td><label><input type="text" name="total-fat" autocomplete="off" placeholder="15g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Saturated fat</td><td><label><input type="text" name="saturated-fat" autocomplete="off" placeholder="1.8g" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Cholesterol</td><td><label><input type="text" name="cholesterol" autocomplete="off" placeholder="1.1mg" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Sodium</td><td><label><input type="text" name="sodium" autocomplete="off" placeholder="100mg" class="input input-bordered input-xs max-w-24"></label></td></tr><tr><td>Fiber</td><td><label><input type="text" name="fiber" autocomplete="off" placeholder="8g" class="input input-bordered input-xs max-w-24"></label></td></tr></tbody></table>`,
				`<table class="table table-zebra table-xs md:h-fit"><thead><tr><th>Time</th><th>h:m:s</th></tr></thead><tbody><tr><td>Prep</td><td><label><input type="text" name="time-preparation" value="00:15:00" class="input input-bordered input-xs max-w-24 html-duration-picker"></label></td></tr><tr><td>Cooking</td><td><label><input type="text" name="time-cooking" value="00:30:00" class="input input-bordered input-xs max-w-24 html-duration-picker"></label></td></tr></tbody></table>`,
				`<ol id="ingredients-list" class="pl-4 list-decimal"><li class="pb-2"><div class="grid grid-flow-col items-center"><label><input required type="text" name="ingredient-1" placeholder="Ingredient #1" class="input input-bordered input-sm w-full" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')"></label><div class="ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: Enter" hx-post="/recipes/add/manual/ingredient" hx-target="#ingredients-list" hx-swap="beforeend" hx-include="[name^='ingredient']">+</button><button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/1" hx-include="[name^='ingredient']">-</button><div class="inline-block h-4 cursor-move handle"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li></ol></div>`,
				`<div class="col-span-6 px-6 py-2 border-gray-700 md:rounded-bl-none md:col-span-4"><h2 class="font-semibold text-center pb-2"><span class="underline">Instructions</span><sup class="text-red-600">*</sup></h2>`,
				`<ol id="instructions-list" class="grid list-decimal"><li class="pt-2 md:pl-0"><div class="flex"><label class="w-11/12"><textarea required name="instruction-1" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #1" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')"></textarea></label><div class="grid ml-2"><button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: CTRL + Enter" hx-post="/recipes/add/manual/instruction" hx-target="#instructions-list" hx-swap="beforeend" hx-include="[name^='instruction']">+</button><button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/1" hx-include="[name^='instruction']">-</button><div class="h-4 cursor-move handle grid place-content-center"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path></svg></div></div></div></li></ol>`,
				`<button class="btn btn-primary btn-block btn-sm">Submit</button>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}

	t.Run("submit recipe", func(t *testing.T) {
		repo = &mockRepository{
			RecipesRegistered:      make(map[int64]models.Recipes),
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		originalNumRecipes := len(repo.RecipesRegistered)

		fields := map[string]string{
			"title":               "Salsa",
			"image":               "eggs.jpg",
			"category":            "appetizers",
			"source":              "Mommy",
			"description":         "The best",
			"calories":            "666",
			"total-carbohydrates": "31g",
			"sugars":              "0.1mg",
			"protein":             "5g",
			"total-fat":           "0g",
			"saturated-fat":       "0g",
			"cholesterol":         "256mg",
			"sodium":              "777mg",
			"fiber":               "2g",
			"time-preparation":    "00:15:30",
			"time-cooking":        "00:30:15",
			"ingredient-1":        "ing1",
			"ingredient-2":        "ing2",
			"instruction-1":       "ins1",
			"instruction-2":       "ins2",
			"yield":               "4",
		}
		contentType, body := createMultipartForm(fields)
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusCreated)
		id := int64(len(repo.RecipesRegistered))
		if len(repo.RecipesRegistered) != originalNumRecipes+1 {
			t.Fatal("expected one more recipe to be added to the database")
		}
		gotRecipe := repo.RecipesRegistered[1][id-1]
		want := models.Recipe{
			Category:     "appetizers",
			CreatedAt:    time.Time{},
			Cuisine:      "",
			Description:  "The best",
			Image:        gotRecipe.Image,
			Ingredients:  []string{"ing1", "ing2"},
			Instructions: []string{"ins1", "ins2"},
			Keywords:     nil,
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
			Tools:     nil,
			UpdatedAt: time.Time{},
			URL:       "Mommy",
			Yield:     4,
		}
		if gotRecipe.Image == uuid.Nil {
			t.Fatal("got nil UUID image when want something other than nil")
		}
		if !cmp.Equal(want, gotRecipe) {
			t.Log(cmp.Diff(want, gotRecipe))
			t.Fail()
		}
		assertHeader(t, rr, "HX-Redirect", "/recipes/"+strconv.FormatInt(id, 10))
	})
}

func TestHandlers_Recipes_AddManualIngredient(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/add/manual/ingredient"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("does not yield input when previous input empty", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("ingredient-1="))

		assertStatus(t, rr.Code, http.StatusUnprocessableEntity)
	})

	t.Run("yields new ingredient input", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("ingredient-1=ingredient1"))

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<input required type="text" name="ingredient-2" placeholder="Ingredient #2" autofocus class="input input-bordered input-sm w-full" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')"></label><div class="ml-2">`,
			`<button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: Enter" hx-post="/recipes/add/manual/ingredient" hx-target="#ingredients-list" hx-swap="beforeend" hx-include="[name^='ingredient']">+</button>`,
			`<button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/2" hx-include="[name^='ingredient']">-</button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_AddManualIngredientDelete(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/add/manual/ingredient"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("does not yield input when only one input left", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+"/1", formHeader, strings.NewReader("ingredient-1="))

		assertStatus(t, rr.Code, http.StatusUnprocessableEntity)
	})

	testcases := []struct {
		name  string
		entry int
		want  []string
	}{
		{
			name:  "delete last entry",
			entry: 4,
			want: []string{
				`<input required type="text" name="ingredient-1" placeholder="Ingredient #1" class="input input-bordered input-sm w-full" value="one" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<input required type="text" name="ingredient-2" placeholder="Ingredient #2" class="input input-bordered input-sm w-full" value="two" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<input required type="text" name="ingredient-3" placeholder="Ingredient #3" class="input input-bordered input-sm w-full" value="three" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
			},
		},
		{
			name:  "delete first entry",
			entry: 1,
			want: []string{
				`<input required type="text" name="ingredient-1" placeholder="Ingredient #1" class="input input-bordered input-sm w-full" value="two" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<input required type="text" name="ingredient-2" placeholder="Ingredient #2" class="input input-bordered input-sm w-full" value="three" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<input required type="text" name="ingredient-3" placeholder="Ingredient #3" class="input input-bordered input-sm w-full" value="''" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
			},
		},
		{
			name:  "delete middle entry",
			entry: 3,
			want: []string{
				`<input required type="text" name="ingredient-1" placeholder="Ingredient #1" class="input input-bordered input-sm w-full" value="one" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<input required type="text" name="ingredient-2" placeholder="Ingredient #2" class="input input-bordered input-sm w-full" value="two" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')"`,
				`<input required type="text" name="ingredient-3" placeholder="Ingredient #3" class="input input-bordered input-sm w-full" value="''" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+"/"+strconv.Itoa(tc.entry), formHeader, strings.NewReader("ingredient-1=one&ingredient-2=two&ingredient-3=three&ingredient-4=''"))

			assertStatus(t, rr.Code, http.StatusOK)
			tc.want = append(tc.want, []string{
				`<button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/1" hx-include="[name^='ingredient']">-</button>`,
				`<button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/2" hx-include="[name^='ingredient']">-</button>`,
				`<button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/3" hx-include="[name^='ingredient']">-</button>`,
			}...)
			assertStringsInHTML(t, getBodyHTML(rr), tc.want)
		})
	}
}

func TestHandlers_Recipes_AddManualInstruction(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/add/manual/instruction"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("does not yield input when previous input empty", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("instruction-1="))

		assertStatus(t, rr.Code, http.StatusUnprocessableEntity)
	})

	t.Run("yields new instruction input", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("instruction-1=one"))

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<textarea required name="instruction-2" rows="3" autofocus class="textarea textarea-bordered w-full" placeholder="Instruction #2" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')"></textarea>`,
			`<button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/2" hx-include="[name^='instruction']">-</button>`,
			`<button type="button" class="btn btn-square btn-sm btn-outline btn-success" title="Shortcut: CTRL + Enter" hx-post="/recipes/add/manual/instruction" hx-target="#instructions-list" hx-swap="beforeend" hx-include="[name^='instruction']">+</button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_AddManualInstructionDelete(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/add/manual/instruction"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("does not yield input when only one input left", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+"/1", formHeader, strings.NewReader("instruction-1="))

		assertStatus(t, rr.Code, http.StatusUnprocessableEntity)
	})

	testcases := []struct {
		name  string
		entry int
		want  []string
	}{
		{
			name:  "delete last entry",
			entry: 4,
			want: []string{
				`<textarea required name="instruction-1" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #1" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">One</textarea>`,
				`<textarea required name="instruction-2" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #2" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">Two</textarea>`,
				`<textarea required name="instruction-3" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #3" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">Three</textarea>`,
			},
		},
		{
			name:  "delete first entry",
			entry: 1,
			want: []string{
				`<textarea required name="instruction-1" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #1" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">Two</textarea>`,
				`<textarea required name="instruction-2" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #2" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">Three</textarea>`,
				`<textarea required name="instruction-3" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #3" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">''</textarea>`,
			},
		},
		{
			name:  "delete middle entry",
			entry: 3,
			want: []string{
				`<textarea required name="instruction-1" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #1" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">One</textarea>`,
				`<textarea required name="instruction-2" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #2" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">Two</textarea>`,
				`<textarea required name="instruction-3" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #3" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">''</textarea>`,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+"/"+strconv.Itoa(tc.entry), formHeader, strings.NewReader("instruction-1=One&instruction-2=Two&instruction-3=Three&instruction-4=''"))

			assertStatus(t, rr.Code, http.StatusOK)
			tc.want = append(tc.want, []string{
				`<button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/1" hx-include="[name^='instruction']">-</button>`,
				`<button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/2" hx-include="[name^='instruction']">-</button>`,
				`<button type="button" class="delete-button btn btn-square btn-sm btn-outline btn-error" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/3" hx-include="[name^='instruction']">-</button>`,
			}...)
			assertStringsInHTML(t, getBodyHTML(rr), tc.want)
		})
	}
}

func TestHandlers_Recipes_AddOCR(t *testing.T) {
	srv := newServerTest()
	originalIntegrations := srv.Integrations
	originalRepo := srv.Repository

	uri := "/recipes/add/ocr"

	sendReq := func(image string) *httptest.ResponseRecorder {
		fields := map[string]string{"image": image}
		contentType, body := createMultipartForm(fields)
		return sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, header(contentType), strings.NewReader(body))
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("image file must not be empty", func(t *testing.T) {
		rr := sendReq("")

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Could not retrieve the image from the form.\",\"title\":\"Form Error\"}"}`)
	})

	t.Run("processing OCR failed", func(t *testing.T) {
		srv.Integrations = &mockIntegrations{
			ProcessImageOCRFunc: func(_ io.Reader) (models.Recipe, error) {
				return models.Recipe{}, errors.New("error")
			},
		}
		defer func() {
			srv.Integrations = originalIntegrations
		}()

		rr := sendReq("hello.jpg")

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Could not process OCR.\",\"title\":\"Integrations Error\"}"}`)
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

		assertStatus(t, rr.Code, http.StatusCreated)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-info\",\"message\":\"\",\"title\":\"Recipe scanned and uploaded.\"}"}`)
		if len(repo.RecipesRegistered[1]) != 1 && repo.RecipesRegistered[1][0].ID != 1 {
			t.Fatal("expected the recipe to be added")
		}
	})
}

func TestHandlers_Recipes_AddRequestWebsite(t *testing.T) {
	srv := newServerTest()
	emailMock := &mockEmail{}
	srv.Email = emailMock

	uri := "/recipes/add/request-website"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("request website successful", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(`website=https://www.eatingbirdfood.com/cinnamon-rolls`))

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", "/recipes/add")
	})
}

func TestHandlers_Recipes_AddWebsite(t *testing.T) {
	srv, ts, c := createWSServer()
	defer func() {
		_ = c.Close()
	}()

	originalBrokers := maps.Clone(srv.Brokers)
	originalRepo := srv.Repository
	originalScraper := srv.Scraper

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

	t.Run("no ws connection", func(t *testing.T) {
		srv.Brokers = map[int64]*models.Broker{}
		defer func() {
			srv.Brokers = originalBrokers
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-warning\",\"message\":\"Connection lost. Please reload page.\",\"title\":\"Websocket\"}"}`)
	})

	t.Run("no input", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("websites="))

		assertStatus(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("no valid URLs", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("urls=I am a pig\noink oink"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"No valid URLs found.\",\"title\":\"Request Error\"}"}`)
	})

	t.Run("add one valid URL from unsupported websites", func(t *testing.T) {
		repo := prepare()
		srv.Scraper = &mockScraper{
			scraperFunc: func(_ string, _ services.FilesService) (models.RecipeSchema, error) {
				return models.RecipeSchema{}, errors.New("unsupported website")
			},
		}
		defer func() {
			srv.Repository = originalRepo
			srv.Scraper = originalScraper
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
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"View /reports?view=latest","background":"alert-info","message":"Fetched 3 recipes. 0 skipped","title":"Operation Successful"}}`
		assertWebsocket(t, c, 6, want)
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

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri+"/5", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusNoContent)
		if numRecipesBefore != len(repo.RecipesRegistered) {
			t.Fail()
		}
	})

	t.Run("can delete user's recipe", func(t *testing.T) {
		r := &models.Recipe{ID: 1, Name: "Chicken"}
		_, _ = srv.Repository.AddRecipe(r, 1, models.UserSettings{})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri+"/1", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusNoContent)
		assertHeader(t, rr, "HX-Redirect", "/")
	})
}

func TestHandlers_Recipes_Edit(t *testing.T) {
	srv := newServerTest()

	baseRecipe := models.Recipe{
		Category:     "american",
		CreatedAt:    time.Now(),
		Cuisine:      "indonesian",
		Description:  "A delicious recipe!",
		ID:           1,
		Image:        uuid.New(),
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
		Tools:     []string{"spoons", "drill"},
		UpdatedAt: time.Now(),
		URL:       "https://example.com/recipes/yummy",
		Yield:     12,
	}
	srv.Repository = &mockRepository{
		RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}},
	}

	uri := "/recipes/%d/edit"

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
			srv.Repository = &mockRepository{RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}}}
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, fmt.Sprintf(uri, 1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to retrieve recipe.\",\"title\":\"Database Error\"}"}`)
	})

	t.Run("successful request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, fmt.Sprintf(uri, 1), noHeader, nil)

		want := []string{
			`<title hx-swap-oob="true">Edit Chicken Jersey | Recipya</title>`,
			`<input required type="text" name="title" placeholder="Title of the recipe*" autocomplete="off" class="input w-full btn-ghost text-center" value="Chicken Jersey">`,
			`<img src="/data/images/` + baseRecipe.Image.String() + `.jpg" alt="Image preview of the recipe." class="object-cover"><span><input type="file" accept="image/*" name="image" class="file-input file-input-sm w-full max-w-xs" value="/data/images/` + baseRecipe.Image.String() + `.jpg" _="on dragover or dragenter halt the event then set the target's style.background to 'lightgray' on dragleave or drop set the target's style.background to '' on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from next <button/>">`,
			`<input required type="text" list="categories" name="category" class="input input-bordered input-sm w-48 md:w-36 lg:w-48" placeholder="Breakfast" value="american"><datalist id="categories"></datalist>`,
			`<input type="number" min="1" name="yield" class="input input-bordered input-sm w-24 md:w-20 lg:w-24" value="12">`,
			`<input required type="text" placeholder="Source" name="source" class="input input-bordered input-sm md:w-28 lg:w-40 xl:w-44" value="https://example.com/recipes/yummy"`,
			`<textarea required name="description" placeholder="This Thai curry chicken will make you drool." class="textarea w-full h-full resize-none">A delicious recipe!</textarea>`,
			`<tbody><tr><td>Prep</td><td><label><input type="text" name="time-preparation" class="input input-bordered input-xs max-w-24 html-duration-picker" value="00:30:00"></label></td></tr><tr><td>Cooking</td><td><label><input type="text" name="time-cooking" class="input input-bordered input-xs max-w-24 html-duration-picker" value="01:00:00"></label></td></tr></tbody>`,
			`<tbody><tr><td>Calories</td><td><label><input type="text" name="calories" autocomplete="off" placeholder="368kcal" class="input input-bordered input-xs max-w-24" value="354"></label></td></tr><tr><td>Total carbs</td><td><label><input type="text" name="total-carbohydrates" autocomplete="off" placeholder="35g" class="input input-bordered input-xs max-w-24" value="7g"></label></td></tr><tr><td>Sugars</td><td><label><input type="text" name="sugars" autocomplete="off" placeholder="3g" class="input input-bordered input-xs max-w-24" value="6g"></label></td></tr><tr><td>Protein</td><td><label><input type="text" name="protein" autocomplete="off" placeholder="21g" class="input input-bordered input-xs max-w-24" value="3g"></label></td></tr><tr><td>Total fat</td><td><label><input type="text" name="total-fat" autocomplete="off" placeholder="15g" class="input input-bordered input-xs max-w-24" value="8g"></label></td></tr><tr><td>Saturated fat</td><td><label><input type="text" name="saturated-fat" autocomplete="off" placeholder="1.8g" class="input input-bordered input-xs max-w-24" value="4g"></label></td></tr><tr><td>Cholesterol</td><td><label><input type="text" name="cholesterol" autocomplete="off" placeholder="1.1mg" class="input input-bordered input-xs max-w-24" value="1g"></label></td></tr><tr><td>Sodium</td><td><label><input type="text" name="sodium" autocomplete="off" placeholder="100mg" class="input input-bordered input-xs max-w-24" value="5g"></label></td></tr><tr><td>Fiber</td><td><label><input type="text" name="fiber" autocomplete="off" placeholder="8g" class="input input-bordered input-xs max-w-24" value="2g"></label></td></tr></tbody>`,
			`<input type="text" name="time-preparation" class="input input-bordered input-xs max-w-24 html-duration-picker" value="00:30:00">`,
			`<input required type="text" name="ingredient-1" placeholder="Ingredient #1" class="input input-bordered input-sm w-full" value="ing1" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
			`<textarea required name="instruction-1" rows="3" class="textarea textarea-bordered w-full" placeholder="Instruction #1" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">ins1</textarea>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("no updated image", func(t *testing.T) {
		files := &mockFiles{}
		srv.Files = files
		srv.Repository = &mockRepository{RecipesRegistered: map[int64]models.Recipes{1: {baseRecipe}}}
		contentType, body := createMultipartForm(map[string]string{"image": ""})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		got, _ := srv.Repository.Recipe(baseRecipe.ID, 1)
		if got.Image != baseRecipe.Image {
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
		contentType, body := createMultipartForm(map[string]string{"image": "eggs.jpg"})

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPut, fmt.Sprintf(uri, 1), header(contentType), strings.NewReader(body))

		assertStatus(t, rr.Code, http.StatusNoContent)
		got, _ := srv.Repository.Recipe(baseRecipe.ID, 1)
		if got.Image == baseRecipe.Image {
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
		fields := map[string]string{
			"title":               "Salsa",
			"image":               "jesus.jpg",
			"category":            "appetizers",
			"source":              "Mommy",
			"description":         "The best",
			"calories":            "666",
			"total-carbohydrates": "31g",
			"sugars":              "0.1mg",
			"protein":             "5g",
			"total-fat":           "24g",
			"saturated-fat":       "58g",
			"cholesterol":         "256mg",
			"sodium":              "777mg",
			"fiber":               "2g",
			"time-preparation":    "00:15:30",
			"time-cooking":        "00:30:15",
			"ingredient-1":        "cheese",
			"ingredient-2":        "avocado",
			"instruction-1":       "mix",
			"instruction-2":       "eat",
			"yield":               "4",
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
			Image:        got.Image,
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
}

func TestHandlers_Recipes_Scale(t *testing.T) {
	srv := newServerTest()
	originalRepo := srv.Repository

	uri := "/recipes/1/scale"

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
			rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?yield="+tc.in, noHeader, nil)

			assertStatus(t, rr.Code, http.StatusBadRequest)
			assertHeader(t, rr, "HX-Trigger", fmt.Sprintf(`{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"%s\",\"title\":\"\"}"}`, tc.want))

		})
	}

	t.Run("cannot find recipe in database", func(t *testing.T) {
		srv.Repository = &mockRepository{
			RecipesRegistered: map[int64]models.Recipes{1: make(models.Recipes, 0)},
		}
		defer func() {
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?yield=6", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusNotFound)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Recipe not found.\",\"title\":\"\"}"}`)
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

		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri+"?yield=8", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">4 lb chicken</span></label>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">1 cup bread loaf</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">1 tbsp beef broth</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">15 cups flour</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">4 big apples</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">Lots of big apples</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">5 slices of bacon</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">4 2/3 cans of bamboo sticks</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">3 can of tomato paste</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">13 1/2 peanut butter jars</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">15 mL of whiskey</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">1 1/3 tbsp lemon juice</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">Ground ginger</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">6 Large or 8 medium ripe Hass avocados</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">1/2 tsp salt plus more for seasoning</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">1 fresh pineapple, cored and cut into 3-inch pieces</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">Un sac de chips de 1kg</span>`,
			`<label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">4 15-ounce can Goya beans</span>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_Search(t *testing.T) {
	srv := newServerTest()
	srv.Repository = &mockRepository{
		RecipesRegistered: map[int64]models.Recipes{
			1: {
				{ID: 1, Name: "Chinese Firmware"},
				{ID: 2, Name: "Lovely Canada"},
				{ID: 3, Name: "Lovely Ukraine"},
			},
		},
	}

	uri := "/recipes/search"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("display all recipes if search is empty", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(""))

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<section class="card card-compact bg-base-100 shadow-lg indicator"><span class="indicator-item indicator-center badge badge-primary select-none"></span><figure class="relative cursor-pointer" hx-get="/recipes/1" hx-target="#content" hx-push-url="true"><img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the Chinese Firmware recipe"><div class="absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 flex items-center justify-center text-white select-none rounded-t-lg"><p class="p-2 text-sm"></p></div></figure><div class="card-body" lang="en"><h2 class="font-semibold w-[25ch] break-words">Chinese Firmware</h2><div class="card-actions h-full flex-col-reverse"><button class="btn btn-block btn-sm btn-outline" hx-get="/recipes/1" hx-target="#content" hx-push-url="true">View</button></div></div></section>`,
			`<section class="card card-compact bg-base-100 shadow-lg indicator"><span class="indicator-item indicator-center badge badge-primary select-none"></span><figure class="relative cursor-pointer" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"><img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Canada recipe"><div class="absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 flex items-center justify-center text-white select-none rounded-t-lg"><p class="p-2 text-sm"></p></div></figure><div class="card-body" lang="en"><h2 class="font-semibold w-[25ch] break-words">Lovely Canada</h2><div class="card-actions h-full flex-col-reverse"><button class="btn btn-block btn-sm btn-outline" hx-get="/recipes/2" hx-target="#content" hx-push-url="true">View</button></div></div></section>`,
			`<section class="card card-compact bg-base-100 shadow-lg indicator"><span class="indicator-item indicator-center badge badge-primary select-none"></span><figure class="relative cursor-pointer" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"><img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Ukraine recipe"><div class="absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 flex items-center justify-center text-white select-none rounded-t-lg"><p class="p-2 text-sm"></p></div></figure><div class="card-body" lang="en"><h2 class="font-semibold w-[25ch] break-words">Lovely Ukraine</h2><div class="card-actions h-full flex-col-reverse"><button class="btn btn-block btn-sm btn-outline" hx-get="/recipes/3" hx-target="#content" hx-push-url="true">View</button></div></div></section>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("no results", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("q=kool-aid"))

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
				`<section class="card card-compact bg-base-100 shadow-lg indicator"><span class="indicator-item indicator-center badge badge-primary select-none"></span><figure class="relative cursor-pointer" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"><img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Canada recipe"><div class="absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 flex items-center justify-center text-white select-none rounded-t-lg"><p class="p-2 text-sm"></p></div></figure><div class="card-body" lang="en"><h2 class="font-semibold w-[25ch] break-words">Lovely Canada</h2><div class="card-actions h-full flex-col-reverse"><button class="btn btn-block btn-sm btn-outline" hx-get="/recipes/2" hx-target="#content" hx-push-url="true">View</button></div></div></section>`,
				`<section class="card card-compact bg-base-100 shadow-lg indicator"><span class="indicator-item indicator-center badge badge-primary select-none"></span><figure class="relative cursor-pointer" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"><img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Ukraine recipe"><div class="absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 flex items-center justify-center text-white select-none rounded-t-lg"><p class="p-2 text-sm"></p></div></figure><div class="card-body" lang="en"><h2 class="font-semibold w-[25ch] break-words">Lovely Ukraine</h2><div class="card-actions h-full flex-col-reverse"><button class="btn btn-block btn-sm btn-outline" hx-get="/recipes/3" hx-target="#content" hx-push-url="true">View</button></div></div></section>`,
			},
		},
		{
			query: "chi",
			want: []string{
				`<section class="card card-compact bg-base-100 shadow-lg indicator"><span class="indicator-item indicator-center badge badge-primary select-none"></span><figure class="relative cursor-pointer" hx-get="/recipes/1" hx-target="#content" hx-push-url="true"><img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the Chinese Firmware recipe"><div class="absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 flex items-center justify-center text-white select-none rounded-t-lg"><p class="p-2 text-sm"></p></div></figure><div class="card-body" lang="en"><h2 class="font-semibold w-[25ch] break-words">Chinese Firmware</h2><div class="card-actions h-full flex-col-reverse"><button class="btn btn-block btn-sm btn-outline" hx-get="/recipes/1" hx-target="#content" hx-push-url="true">View</button></div></div></section>`,
			},
		},
		{
			query: "lovely",
			want: []string{
				`<section class="card card-compact bg-base-100 shadow-lg indicator"><span class="indicator-item indicator-center badge badge-primary select-none"></span><figure class="relative cursor-pointer" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"><img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Canada recipe"><div class="absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 flex items-center justify-center text-white select-none rounded-t-lg"><p class="p-2 text-sm"></p></div></figure><div class="card-body" lang="en"><h2 class="font-semibold w-[25ch] break-words">Lovely Canada</h2><div class="card-actions h-full flex-col-reverse"><button class="btn btn-block btn-sm btn-outline" hx-get="/recipes/2" hx-target="#content" hx-push-url="true">View</button></div></div></section>`,
				`<section class="card card-compact bg-base-100 shadow-lg indicator"><span class="indicator-item indicator-center badge badge-primary select-none"></span><figure class="relative cursor-pointer" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"><img class="h-48 max-w-48 w-full object-cover rounded-t-lg" src="/static/img/recipes/placeholder.webp" alt="Image for the Lovely Ukraine recipe"><div class="absolute inset-0 bg-black opacity-0 hover:opacity-80 transition-opacity duration-300 flex items-center justify-center text-white select-none rounded-t-lg"><p class="p-2 text-sm"></p></div></figure><div class="card-body" lang="en"><h2 class="font-semibold w-[25ch] break-words">Lovely Ukraine</h2><div class="card-actions h-full flex-col-reverse"><button class="btn btn-block btn-sm btn-outline" hx-get="/recipes/3" hx-target="#content" hx-push-url="true">View</button></div></div></section>`,
			},
		},
	}
	for _, tc := range searches {
		t.Run("results for "+tc.query, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("id=1&page=1&q="+tc.query))

			assertStatus(t, rr.Code, http.StatusOK)
			assertStringsInHTML(t, getBodyHTML(rr), tc.want)
		})
	}
}

func TestHandlers_Recipes_Share(t *testing.T) {
	srv := newServerTest()

	uri := func(id int64) string {
		return fmt.Sprintf("/recipes/%d/share", id)
	}

	app.Config.Server.URL = "https://www.recipya.com"
	recipe := &models.Recipe{
		Category:     "American",
		Description:  "This is the most delicious recipe!",
		ID:           1,
		Image:        uuid.New(),
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
	_, _ = srv.Repository.AddRecipe(recipe, 1, models.UserSettings{})
	link, _ := srv.Repository.AddShareLink(models.Share{RecipeID: 1, CookbookID: -1, UserID: 1})

	t.Run("create valid share link", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<input type="url" value="example.com/r/33320755-82f9-47e5-bb0a-d1b55cbd3f7b" class="input input-bordered w-full" readonly="readonly"></label>`,
			`<button id="copy_button" class="btn btn-neutral" title="Copy to clipboard" onClick="__templ_copyToClipboard`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("create invalid share link", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri(10), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"action\":\"\",\"background\":\"alert-error\",\"message\":\"Failed to create share link.\",\"title\":\"\"}"}`)
	})

	t.Run("access share link anonymous", func(t *testing.T) {
		rr := sendRequest(srv, http.MethodGet, link, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
			`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d=`,
			`<a href="/auth/login" class="btn btn-ghost">Log In</a><a href="/auth/register" class="btn btn-ghost">Sign Up</a>`,
			`<span class="text-center pb-2">Chicken Jersey</span>`,
			`<button class="mr-2" title="Print recipe" _="on click print()">`,
			`<img id="output" style="object-fit: cover" alt="Image of the recipe" class="w-full" src="/data/images/` + recipe.Image.String() + `.jpg">`,
			`<div class="badge badge-primary badge-outline">American</div>`,
			`<p class="text-sm text-center">2 servings</p>`,
			`<a class="btn btn-sm btn-outline no-underline print:hidden" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank">Source</a><p class="hidden print:block print:whitespace-nowrap print:overflow-hidden print:text-ellipsis print:max-w-xs">https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/</p>`,
			`<textarea class="textarea w-full h-full resize-none" readonly>This is the most delicious recipe!</textarea>`,
			`<p class="text-xs">Per 100g: 500 calories; total carbohydrates 7g; sugars 6g; protein 3g; total fat 8g; saturated fat 4g; cholesterol 1g; sodium 5g; fiber 2g.</p>`,
			`<table class="table table-zebra table-xs md:h-fit print:hidden"><thead><tr><th>Time</th><th>h:m:s</th></tr></thead><tbody><tr><td>Prep:</td><td><time datetime="PT05M">5m</time></td></tr><tr><td>Cooking:</td><td><time datetime="PT1H05M">1h05m</time></td></tr><tr><td>Total:</td><td><time datetime="PT1H10M">1h10m</time></td></tr></tbody></table>`,
			`<table class="table table-zebra table-xs print:hidden"><thead><tr><th>Nutrition (per 100g</th><th>Amount</th></tr></thead><tbody><tr><td>Calories:</td><td>500</td></tr><tr><td>Total carbs:</td><td>7g</td></tr><tr><td>Sugars:</td><td>6g</td></tr><tr><td>Protein:</td><td>3g</td></tr><tr><td>Total fat:</td><td>8g</td></tr><tr><td>Saturated fat:</td><td>4g</td></tr><tr><td>Cholesterol:</td><td>1g</td></tr><tr><td>Sodium:</td><td>5g</td></tr><tr><td>Fiber:</td><td>2g</td></tr></tbody></table>`,
			`<div id="ingredients-instructions-container" class="grid text-sm md:grid-flow-col md:col-span-6"><div class="col-span-6 border-gray-700 px-4 py-2 border-y md:col-span-2 md:border-r md:border-y-0 print:hidden"><h2 class="font-semibold text-center underline md:pb-2">Ingredients</h2><ul><li class="form-control hover:bg-gray-100 dark:hover:bg-gray-700"><label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">Ing1</span></label></li><li class="form-control hover:bg-gray-100 dark:hover:bg-gray-700"><label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">Ing2</span></label></li><li class="form-control hover:bg-gray-100 dark:hover:bg-gray-700"><label class="label justify-start"><input type="checkbox" class="checkbox"><span class="label-text pl-2">Ing3</span></label></li></ul></div><div class="col-span-6 px-8 py-2 border-gray-700 md:rounded-bl-none md:col-span-4 print:hidden"><h2 class="font-semibold text-center underline md:pb-2">Instructions</h2><ol class="grid list-decimal"><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins1</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins2</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins3</span></li></ol></div></div><div class="hidden print:grid col-span-6 ml-2 mt-2 mb-4"><h1 class="text-sm"><b>Ingredients</b></h1><ol class="col-span-6 w-full" style="column-count: 1"><li class="text-sm"><label><input type="checkbox"></label><span class="pl-2">Ing1</span></li><li class="text-sm"><label><input type="checkbox"></label><span class="pl-2">Ing2</span></li><li class="text-sm"><label><input type="checkbox"></label><span class="pl-2">Ing3</span></li></ol></div><div class="hidden print:grid col-span-5"><h1 class="text-sm"><b>Instructions</b></h1><ol class="col-span-6 list-decimal w-full ml-6" style="column-count: 2; column-gap: 2.5rem"><li><span class="text-sm whitespace-pre-line">Ins1</span></li><li><span class="text-sm whitespace-pre-line">Ins2</span></li><li><span class="text-sm whitespace-pre-line">Ins3</span></li></ol></div></div>`,
			`<h2 class="font-semibold text-center underline md:pb-2">Instructions</h2><ol class="grid list-decimal"><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins1</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins2</span></li><li class="min-w-full py-2 select-none hover:bg-gray-100 dark:hover:bg-gray-700" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins3</span></li></ol></div></div><div class="hidden print:grid col-span-6 ml-2 mt-2 mb-4"><h1 class="text-sm"><b>Ingredients</b></h1><ol class="col-span-6 w-full" style="column-count: 1"><li class="text-sm"><label><input type="checkbox"></label><span class="pl-2">Ing1</span></li><li class="text-sm"><label><input type="checkbox"></label><span class="pl-2">Ing2</span></li><li class="text-sm"><label><input type="checkbox"></label><span class="pl-2">Ing3</span></li></ol></div><div class="hidden print:grid col-span-5"><h1 class="text-sm"><b>Instructions</b></h1><ol class="col-span-6 list-decimal w-full ml-6" style="column-count: 2; column-gap: 2.5rem"><li><span class="text-sm whitespace-pre-line">Ins1</span></li><li><span class="text-sm whitespace-pre-line">Ins2</span></li><li><span class="text-sm whitespace-pre-line">Ins3</span></li></ol></div>`,
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
		sendFunc func(server *server.Server, target, uri string, contentType header, body *strings.Reader) *httptest.ResponseRecorder
	}{
		{name: "other user Hx-Request", sendFunc: sendHxRequestAsLoggedInOther},
		{name: "other user no Hx-Request", sendFunc: sendRequestAsLoggedInOther},
	}
	for _, tc := range testcases {
		t.Run("access share link logged in "+tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, link, noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
				`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d=`,
				`<button class="mr-2" title="Add recipe to collection" hx-post="/recipes/1/share/add">`,
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
		sendFunc func(server *server.Server, target, uri string, contentType header, body *strings.Reader) *httptest.ResponseRecorder
	}{
		{name: "host user Hx-Request", sendFunc: sendHxRequestAsLoggedIn},
		{name: "host user no Hx-Request", sendFunc: sendRequestAsLoggedIn},
	}
	for _, tc := range testcases2 {
		t.Run("access share link logged "+tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, link, noHeader, nil)

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

func TestHandlers_Recipes_SupportedWebsites(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/supported-websites"

	website1 := `<tr class="border text-center"><td class="border dark:border-gray-800">1</td><td class="border py-1 dark:border-gray-800"><a class="underline" href="https://101cookbooks.com" target="_blank">101cookbooks.com</a></td></tr>`
	website2 := `<tr class="border text-center"><td class="border dark:border-gray-800">2</td><td class="border py-1 dark:border-gray-800"><a class="underline" href="https://www.afghankitchenrecipes.com" target="_blank">afghankitchenrecipes.com</a></td></tr>`

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("returns list of websites to logged in user", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

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
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri+"/999", noHeader, nil)

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
		sendFunc func(server *server.Server, target, uri string, contentType header, body *strings.Reader) *httptest.ResponseRecorder
	}{
		{name: "logged in Hx-Request", sendFunc: sendHxRequestAsLoggedIn},
		{name: "logged in no Hx-Request", sendFunc: sendRequestAsLoggedIn},
	}
	for _, tc := range testcases {
		t.Run("recipe is in user's collection when "+tc.name, func(t *testing.T) {
			image, _ := uuid.Parse("e81ba735-a4af-4c66-8c17-2f2ccc1b1a95")
			r := &models.Recipe{
				Category:     "American",
				Description:  "This is the most delicious recipe!",
				ID:           1,
				Image:        image,
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
			_, _ = srv.Repository.AddRecipe(r, 1, models.UserSettings{})

			rr := tc.sendFunc(srv, http.MethodGet, uri+"/1", noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + r.Name + " | Recipya</title>",
				`<button title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d='M 12.276 18.55 v -0.748 a 4.79 4.79 0 0 1 1.463 -3.458 a 5.763 5.763 0 0 0 1.804 -4.21 a 5.821 5.821 0 0 0 -6.475 -5.778 c -2.779 0.307 -4.99 2.65 -5.146 5.448 a 5.82 5.82 0 0 0 1.757 4.503 a 4.906 4.906 0 0 1 1.5 3.495 v 0.747 a 1.44 1.44 0 0 0 1.44 1.439 h 2.218 a 1.44 1.44 0 0 0 1.44 -1.439 z m -1.058 0 c 0 0.209 -0.17 0.38 -0.38 0.38 h -2.22 c -0.21 0 -0.38 -0.171 -0.38 -0.38 v -0.748 c 0 -1.58 -0.664 -3.13 -1.822 -4.254 A 4.762 4.762 0 0 1 4.98 9.863 c 0.127 -2.289 1.935 -4.204 4.205 -4.455 a 4.762 4.762 0 0 1 5.3 4.727 a 4.714 4.714 0 0 1 -1.474 3.443 a 5.853 5.853 0 0 0 -1.791 4.225 v 0.746 z M 11.45 20.51 H 8.006 a 0.397 0.397 0 1 0 0 0.795 h 3.444 a 0.397 0.397 0 1 0 0 -0.794 z M 11.847 22.162 a 0.397 0.397 0 0 0 -0.397 -0.397 H 8.006 a 0.397 0.397 0 1 0 0 0.794 h 3.444 c 0.22 0 0.397 -0.178 0.397 -0.397 z z z z z z z z M 10.986 23.416 H 8.867 a 0.397 0.397 0 1 0 0 0.794 h 1.722 c 0.22 0 0.397 -0.178 0.397 -0.397 z' to #icon-bulb else call initWakeLock() then add @d='M12.276 18.55v-.748a4.79 4.79 0 0 1 1.463-3.458 5.763 5.763 0 0 0 1.804-4.21 5.821 5.821 0 0 0-6.475-5.778c-2.779.307-4.99 2.65-5.146 5.448a5.82 5.82 0 0 0 1.757 4.503 4.906 4.906 0 0 1 1.5 3.495v.747a1.44 1.44 0 0 0 1.44 1.439h2.218a1.44 1.44 0 0 0 1.44-1.439zm-1.058 0c0 .209-.17.38-.38.38h-2.22c-.21 0-.38-.171-.38-.38v-.748c0-1.58-.664-3.13-1.822-4.254A4.762 4.762 0 0 1 4.98 9.863c.127-2.289 1.935-4.204 4.205-4.455a4.762 4.762 0 0 1 5.3 4.727 4.714 4.714 0 0 1-1.474 3.443 5.853 5.853 0 0 0-1.791 4.225v.746zM11.45 20.51H8.006a.397.397 0 1 0 0 .795h3.444a.397.397 0 1 0 0-.794zM11.847 22.162a.397.397 0 0 0-.397-.397H8.006a.397.397 0 1 0 0 .794h3.444c.22 0 .397-.178.397-.397zM.397 10.125h2.287a.397.397 0 1 0 0-.794H.397a.397.397 0 1 0 0 .794zM19.456 9.728a.397.397 0 0 0-.397-.397h-2.287a.397.397 0 1 0 0 .794h2.287c.22 0 .397-.178.397-.397zM9.331.397v2.287a.397.397 0 1 0 .794 0V.397a.397.397 0 1 0-.794 0zM16.045 2.85 14.43 4.465a.397.397 0 1 0 .561.561l1.617-1.617a.397.397 0 1 0-.562-.56zM5.027 14.429a.397.397 0 0 0-.56 0l-1.618 1.616a.397.397 0 1 0 .562.562l1.617-1.617a.397.397 0 0 0 0-.561zM4.466 5.027a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.56L3.41 2.848a.397.397 0 1 0-.561.561zM16.045 16.607a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.562L14.99 14.43a.397.397 0 1 0-.561.56zM10.986 23.416a.397.397 0 0 0-.397-.397H8.867a.397.397 0 1 0 0 .794h1.722c.22 0 .397-.178.397-.397z' to #icon-bulb end"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor"><path id="icon-bulb" d="M12.276 18.55v-.748a4.79 4.79 0 0 1 1.463-3.458 5.763 5.763 0 0 0 1.804-4.21 5.821 5.821 0 0 0-6.475-5.778c-2.779.307-4.99 2.65-5.146 5.448a5.82 5.82 0 0 0 1.757 4.503 4.906 4.906 0 0 1 1.5 3.495v.747a1.44 1.44 0 0 0 1.44 1.439h2.218a1.44 1.44 0 0 0 1.44-1.439zm-1.058 0c0 .209-.17.38-.38.38h-2.22c-.21 0-.38-.171-.38-.38v-.748c0-1.58-.664-3.13-1.822-4.254A4.762 4.762 0 0 1 4.98 9.863c.127-2.289 1.935-4.204 4.205-4.455a4.762 4.762 0 0 1 5.3 4.727 4.714 4.714 0 0 1-1.474 3.443 5.853 5.853 0 0 0-1.791 4.225v.746zM11.45 20.51H8.006a.397.397 0 1 0 0 .795h3.444a.397.397 0 1 0 0-.794zM11.847 22.162a.397.397 0 0 0-.397-.397H8.006a.397.397 0 1 0 0 .794h3.444c.22 0 .397-.178.397-.397zM.397 10.125h2.287a.397.397 0 1 0 0-.794H.397a.397.397 0 1 0 0 .794zM19.456 9.728a.397.397 0 0 0-.397-.397h-2.287a.397.397 0 1 0 0 .794h2.287c.22 0 .397-.178.397-.397zM9.331.397v2.287a.397.397 0 1 0 .794 0V.397a.397.397 0 1 0-.794 0zM16.045 2.85 14.43 4.465a.397.397 0 1 0 .561.561l1.617-1.617a.397.397 0 1 0-.562-.56zM5.027 14.429a.397.397 0 0 0-.56 0l-1.618 1.616a.397.397 0 1 0 .562.562l1.617-1.617a.397.397 0 0 0 0-.561zM4.466 5.027a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.56L3.41 2.848a.397.397 0 1 0-.561.561zM16.045 16.607a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.562L14.99 14.43a.397.397 0 1 0-.561.56zM10.986 23.416a.397.397 0 0 0-.397-.397H8.867a.397.397 0 1 0 0 .794h1.722c.22 0 .397-.178.397-.397z"></path></svg></button>`,
				`<button class="ml-2" title="Edit recipe" hx-get="/recipes/1/edit" hx-push-url="true" hx-target="#content"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg></button>`,
				`<span class="text-center pb-2">Chicken Jersey</span>`,
				`<button class="mr-2" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call share_dialog.showModal()">`,
				`<button class="mr-2" title="Print recipe" _="on click print()">`,
				`<button class="mr-2" hx-delete="/recipes/1" hx-swap="none" title="Delete recipe" hx-confirm="Are you sure you wish to delete this recipe?" hx-indicator="#fullscreen-loader">`,
				`<img id="output" style="object-fit: cover" alt="Image of the recipe" class="w-full" src="/data/images/e81ba735-a4af-4c66-8c17-2f2ccc1b1a95.jpg"></div><div class="grid grid-cols-3 col-span-3 md:grid-flow-row md:grid-rows-4 print:grid-rows-2">`,
				`<div class="badge badge-primary badge-outline">American</div>`,
				`<button class="mr-2" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call share_dialog.showModal()"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor">`,
				`<form autocomplete="off" _="on submit halt the event" class="print:hidden"><label class="form-control w-full"><div class="label p-0"><span class="label-text">Servings</span></div><input id="yield" type="number" min="1" name="yield" value="2" class="input input-bordered input-sm w-24" hx-get="/recipes/1/scale" hx-trigger="input" hx-target="#ingredients-instructions-container"></label></form>`,
				`<a class="btn btn-sm btn-outline no-underline print:hidden" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank">Source</a>`,
				`<textarea class="textarea w-full h-full resize-none" readonly>This is the most delicious recipe!</textarea>`,
				`<p class="text-xs">Per 100g: 500 calories; total carbohydrates 7g; sugars 6g; protein 3g; total fat 8g; saturated fat 4g; cholesterol 1g; sodium 5g; fiber 2g.</p>`,
				`<table class="table table-zebra table-xs md:h-fit print:hidden"><thead><tr><th>Time</th><th>h:m:s</th></tr></thead><tbody><tr><td>Prep:</td><td><time datetime="PT05M">5m</time></td></tr><tr><td>Cooking:</td><td><time datetime="PT1H05M">1h05m</time></td></tr><tr><td>Total:</td><td><time datetime="PT1H10M">1h10m</time></td></tr></tbody></table>`,
				`<table class="table table-zebra table-xs print:hidden"><thead><tr><th>Nutrition (per 100g</th><th>Amount</th></tr></thead><tbody><tr><td>Calories:</td><td>500</td></tr><tr><td>Total carbs:</td><td>7g</td></tr><tr><td>Sugars:</td><td>6g</td></tr><tr><td>Protein:</td><td>3g</td></tr><tr><td>Total fat:</td><td>8g</td></tr><tr><td>Saturated fat:</td><td>4g</td></tr><tr><td>Cholesterol:</td><td>1g</td></tr><tr><td>Sodium:</td><td>5g</td></tr><tr><td>Fiber:</td><td>2g</td></tr></tbody></table>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}
