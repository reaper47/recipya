package server_test

import (
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHandlers_Recipes(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	t.Run("user has no recipes", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		want := []string{
			`<title hx-swap-oob="true">Home | Recipya</title>`,
			`<div class="grid place-content-center text-sm h-full text-center md:text-base"><p>Your recipe collection looks a bit empty at the moment.</p><p> Why not start adding some of your favorite recipes by clicking the <a class="underline font-semibold cursor-pointer" hx-get="/recipes/add" hx-target="#content" hx-push-url="true">Add recipe</a> button at the top? </p></div>`,
			`<li id="recipes-sidebar-recipes" class="recipes-sidebar-selected" hx-get="/recipes" hx-target="#content" hx-push-url="true" hx-swap-oob="true"><img src="/static/img/cherries.svg" alt=""><span class="hidden md:block ml-1">Recipes</span></li>`,
			`<li id="recipes-sidebar-cookbooks" class="recipes-sidebar-not-selected" hx-get="/cookbooks" hx-target="#content" hx-push-url="true" hx-swap-oob="true"><img src="/static/img/cookbook.svg" alt=""><span class="hidden md:block ml-1">Cookbooks</span></li>`,
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
			`<title hx-swap-oob="true">Home | Recipya</title>`,
			`<form class="w-72 md:w-96" hx-post="/recipes/search" hx-vals='{"page": 1}' hx-target="#list-recipes" _="on submit if #search-recipes.value is not '' then add .hidden to #pagination else remove .hidden from #pagination end">`,
			`<article class="grid gap-4 p-4 text-sm grid-cols-4 md:m-auto md:max-w-7xl md:text-base md:grid-cols-10">`,
			`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the One recipe" hx-get="/recipes/1" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">One</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"> Lunch </span></div><p class="pb-8 text-justify text-sm">Recipe one</p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/1" hx-target="#content" hx-push-url="true"> View </button></div></div>`,
			`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the Two recipe" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">Two</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"> Soup </span></div><p class="pb-8 text-justify text-sm">Recipe two</p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"> View </button></div></div>`,
			`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the Three recipe" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">Three</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"> Dinner </span></div><p class="pb-8 text-justify text-sm">Recipe three</p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"> View </button></div></div>`,
		}
		assertStringsInHTML(t, got, want)
		notWant := []string{
			`<div class="grid place-content-center text-sm h-full text-center md:text-base"><p>Your recipe collection looks a bit empty at the moment.</p><p> Why not start adding some of your favorite recipes by clicking the <a class="underline font-semibold cursor-pointer" hx-get="/recipes/add" hx-target="#content" hx-push-url="true"> Add recipe </a> button at the top? </p></div>`,
		}
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
				`<img class="object-cover w-full h-40 rounded-t-xl" src="/static/img/recipes-new/import.webp" alt="Writing on a piece of paper with a traditional pen."/>`,
				`<button id="search-button" class="underline" hx-get="/recipes/supported-websites" hx-target="#search-results" onclick="document.querySelector('#supported-websites-dialog').showModal()" > supported </button>`,
				`<button class="w-full border-2 border-gray-800 rounded-lg hover:bg-gray-800 hover:text-white center" hx-target="#content" hx-push-url="/recipes/add/unsupported-website" hx-prompt="Enter the recipe's URL" hx-post="/recipes/add/website" hx-indicator="#fullscreen-loader"> Fetch </button>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}

func TestHandlers_Recipes_AddManual(t *testing.T) {
	repo := &mockRepository{}
	srv := server.NewServer(repo, &mockEmail{}, &mockFiles{}, &mockIntegrations{})

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
				`<title hx-swap-oob="true">Manual | Recipya</title>`,
				`<form hx-post="/recipes/add/manual" enctype="multipart/form-data" class="grid max-w-6xl grid-cols-6 m-auto mt-4 pb-4 dark:bg-gray-600 dark:rounded-t-lg">`,
				`<input type="text" name="title" id="title" placeholder="Title of the recipe*" required autocomplete="off" class="w-full py-2 bg-gray-50 font-semibold text-center text-gray-600 placeholder-gray-400 rounded-t-lg dark:bg-gray-900 dark:text-gray-200">`,
				`<label class="grid col-span-6 place-content-center w-full h-96 border-b border-x border-black md:border-r-0 md:col-span-4 text-sm dark:border-gray-800"><img src="" alt="Image preview of the recipe." class="h-full"><span><input type="file" accept="image/*" name="image" required _="on dragover or dragenter halt the event then set the target's style.background to 'lightgray' on dragleave or drop set the target's style.background to '' on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from next <button/>"><button type="button" class="px-2 bg-red-300 border border-gray-800 rounded-lg hidden hover:bg-red-600 hover:text-white dark:border-gray-800 dark:bg-red-600 dark:hover:bg-red-700" _="on click set {value: ''} on previous <input/> then set {src: ''} on previous <img/> then add .hidden"> Delete </button></span></label>`,
				`<div class="grid col-span-2 py-2 border-black place-items-center text-sm md:col-span-3 md:border-t dark:border-gray-800"><div><label for="yield">Yields</label><input id="yield" type="number" min="1" name="yield" value="4" class="w-24 rounded bg-gray-100 p-2 dark:bg-gray-900"><span>servings</span></div>`,
				`<label class="grid place-content-center p-2 font-medium text-sm text-gray-900 bg-blue-100 border border-blue-300 rounded-full w-fit"><input type="text" list="categories" name="category" class="bg-transparent text-center focus:outline-none" placeholder="Breakfast*" required><datalist id="categories"><option>breakfast</option><option>lunch</option><option>dinner</option></datalist></label></div><div class="grid col-span-2 py-2 border-black place-items-center text-sm md:col-span-3 md:border-t dark:border-gray-800"><div><label for="yield">Yields</label><input id="yield" type="number" min="1" name="yield" value="4" class="w-24 rounded bg-gray-100 p-2 dark:bg-gray-900"><span>servings</span></div>`,
				`<textarea id="description" name="description" rows="10" placeholder="This Thai curry chicken will make you drool." required class="p-2 border border-gray-300 rounded-t-lg dark:bg-gray-800 dark:border-none" ></textarea>`,
				`<th scope="col" class="text-right md:text-center"><p>Nutrition<br>(per 100g)</p></th><th scope="col" class="text-center"><p>Amount<br>(optional)</p></th>`,
				`<th scope="col" class="text-right">Time</th><th scope="col" class="text-center">h:m:s</th>`,
				`<ol id="ingredients-list" class="pl-8 list-decimal">`,
				`<ol id="instructions-list" class="pl-4 list-decimal">`,
				`<input type="text" name="ingredient-1" placeholder="Ingredient #1" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<button type="submit" class="col-span-6 p-2 font-semibold text-white bg-blue-500 hover:bg-blue-800"> Submit </button>`,
				`loadScript("https://cdn.jsdelivr.net/npm/html-duration-picker@latest/dist/html-duration-picker.min.js") .then(() => HtmlDurationPicker.init()) loadScript("https://cdn.jsdelivr.net/npm/sortablejs@latest/Sortable.min.js")`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}

	t.Run("submit recipe", func(t *testing.T) {
		repo = &mockRepository{
			RecipesRegistered: make(map[int64]models.Recipes),
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
	srv := server.NewServer(&mockRepository{}, &mockEmail{}, &mockFiles{}, &mockIntegrations{})

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
			`<input type="text" name="ingredient-2" placeholder="Ingredient #2" required autofocus class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
			`<button type="button" class="delete-button w-10 h-10 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/2" hx-include="[name^='ingredient']"> - </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_AddManualIngredientDelete(t *testing.T) {
	srv := server.NewServer(&mockRepository{}, &mockEmail{}, &mockFiles{}, &mockIntegrations{})

	uri := "/recipes/add/manual/ingredient/"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("does not yield input when only one input left", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+"1", formHeader, strings.NewReader("ingredient-1="))

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
				`<input type="text" name="ingredient-1" placeholder="Ingredient #1" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="one" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<input type="text" name="ingredient-2" placeholder="Ingredient #2" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="two" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<input type="text" name="ingredient-3" placeholder="Ingredient #3" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="three" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
			},
		},
		{
			name:  "delete first entry",
			entry: 1,
			want: []string{
				`<input type="text" name="ingredient-1" placeholder="Ingredient #1" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="two" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<input type="text" name="ingredient-2" placeholder="Ingredient #2" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="three" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<input type="text" name="ingredient-3" placeholder="Ingredient #3" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="''" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
			},
		},
		{
			name:  "delete middle entry",
			entry: 3,
			want: []string{
				`<input type="text" name="ingredient-1" placeholder="Ingredient #1" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="one" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')"><`,
				`<input type="text" name="ingredient-2" placeholder="Ingredient #2" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="two" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
				`<input type="text" name="ingredient-3" placeholder="Ingredient #3" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="''" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')">`,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+strconv.Itoa(tc.entry), formHeader, strings.NewReader("ingredient-1=one&ingredient-2=two&ingredient-3=three&ingredient-4=''"))

			assertStatus(t, rr.Code, http.StatusOK)
			want := append(tc.want, []string{
				`&nbsp;<button type="button" class="delete-button w-10 h-10 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/1" hx-include="[name^='ingredient']">-</button>`,
				`&nbsp;<button type="button" class="delete-button w-10 h-10 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/2" hx-include="[name^='ingredient']">-</button>`,
				`&nbsp;<button type="button" class="delete-button w-10 h-10 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/3" hx-include="[name^='ingredient']">-</button>`,
			}...)
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}

func TestHandlers_Recipes_AddManualInstruction(t *testing.T) {
	srv := server.NewServer(&mockRepository{}, &mockEmail{}, &mockFiles{}, &mockIntegrations{})

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
			`<textarea required name="instruction-2" rows="3" autofocus class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #2" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')"></textarea>`,
			`<button type="button" class="delete-button w-10 h-10 md:w-7 md:h-7 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/2" hx-include="[name^='instruction']"> - </button>`,
			`<button type="button" class="md:w-7 md:h-7 md:right-auto w-10 h-10 text-center bg-green-300 border border-gray-800 rounded-lg hover:bg-green-600 hover:text-white center dark:bg-green-500" title="Shortcut: CTRL + Enter" hx-post="/recipes/add/manual/instruction" hx-target="#instructions-list" hx-swap="beforeend" hx-include="[name^='instruction']"> + </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_AddManualInstructionDelete(t *testing.T) {
	srv := server.NewServer(&mockRepository{}, &mockEmail{}, &mockFiles{}, &mockIntegrations{})

	uri := "/recipes/add/manual/instruction/"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("does not yield input when only one input left", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+"1", formHeader, strings.NewReader("instruction-1="))

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
				`<textarea required name="instruction-1" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #1" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">One</textarea>`,
				`<textarea required name="instruction-2" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #2" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">Two</textarea>`,
				`<textarea required name="instruction-3" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #3" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">Three</textarea>`,
			},
		},
		{
			name:  "delete first entry",
			entry: 1,
			want: []string{
				`<textarea required name="instruction-1" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #1" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">Two</textarea>`,
				`<textarea required name="instruction-2" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #2" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">Three</textarea>`,
				`<textarea required name="instruction-3" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #3" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">''</textarea>`,
			},
		},
		{
			name:  "delete middle entry",
			entry: 3,
			want: []string{
				`<textarea required name="instruction-1" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #1" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">One</textarea>`,
				`<textarea required name="instruction-2" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #2" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">Two</textarea>`,
				`<textarea required name="instruction-3" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #3" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">''</textarea>`,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+strconv.Itoa(tc.entry), formHeader, strings.NewReader("instruction-1=One&instruction-2=Two&instruction-3=Three&instruction-4=''"))

			assertStatus(t, rr.Code, http.StatusOK)
			want := append(tc.want, []string{
				`<button type="button" class="delete-button w-10 h-10 md:w-7 md:h-7 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/1" hx-include="[name^='instruction']">-</button>`,
				`<button type="button" class="delete-button w-10 h-10 md:w-7 md:h-7 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/2" hx-include="[name^='instruction']">-</button>`,
				`<button type="button" class="delete-button w-10 h-10 md:w-7 md:h-7 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/3" hx-include="[name^='instruction']">-</button>`,
			}...)
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}

func TestHandlers_Recipes_AddRequestWebsite(t *testing.T) {
	emailMock := &mockEmail{}
	srv := server.NewServer(&mockRepository{}, emailMock, &mockFiles{}, &mockIntegrations{})

	uri := "/recipes/add/request-website"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("request website successful", func(t *testing.T) {
		originalEmailHitCount := emailMock.hitCount

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader(`website=https://www.eatingbirdfood.com/cinnamon-rolls`))

		assertStatus(t, rr.Code, http.StatusSeeOther)
		assertHeader(t, rr, "HX-Redirect", "/recipes/add")
		if emailMock.hitCount != originalEmailHitCount+1 {
			t.Fatalf("email must have been sent")
		}
	})
}

func TestHandlers_Recipes_AddWebsite(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/add/website"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uri)
	})

	t.Run("add recipe from wrong URL", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader("I love chicken"))

		assertStatus(t, rr.Code, http.StatusBadRequest)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Invalid URI.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("add recipe from an unsupported website", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader("https://www.example.com"))

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "Content-Type", "text/html")
		want := []string{
			`<title hx-swap-oob="true">Unsupported Website | Recipya</title>`,
			`<h3 class="mb-2 text-2xl font-semibold tracking-tight"> This website is not supported </h3>`,
			`<p class="mb-3 text-gray-700"> Unfortunately, we could not extract the recipe from this link. You can either request that our team support this website or go back to the previous page. </p>`,
			`<button name="website" value="https://www.example.com" class="w-full col-span-4 ml-2 p-2 font-semibold text-white bg-blue-500 hover:bg-blue-800"> Request </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("add recipe from supported website error", func(t *testing.T) {
		repo := &mockRepository{RecipesRegistered: make(map[int64]models.Recipes)}
		repo.AddRecipeFunc = func(r *models.Recipe, userID int64) (uint64, error) {
			return 0, errors.New("add recipe error")
		}
		srv.Repository = repo

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader("https://www.eatingbirdfood.com/cinnamon-rolls/"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Recipe could not be added.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("add recipe from a supported website", func(t *testing.T) {
		repo := &mockRepository{RecipesRegistered: make(map[int64]models.Recipes)}
		called := 0
		repo.AddRecipeFunc = func(r *models.Recipe, userID int64) (uint64, error) {
			called += 1
			return 1, nil
		}
		srv.Repository = repo

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader("https://www.eatingbirdfood.com/cinnamon-rolls/"))

		assertStatus(t, rr.Code, http.StatusSeeOther)
		if called != 1 {
			t.Fatal("recipe must have been added to the user's database")
		}
		assertHeader(t, rr, "HX-Redirect", "/recipes/1")
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
		rr := sendHxRequestAsLoggedIn(srv, http.MethodDelete, uri+"/5", noHeader, nil)

		assertStatus(t, rr.Code, http.StatusNotFound)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Recipe not found.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("can delete user's recipe", func(t *testing.T) {
		r := &models.Recipe{ID: 1, Name: "Chicken"}
		_, _ = srv.Repository.AddRecipe(r, 1)

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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Failed to retrieve recipe.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("successful request", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodGet, fmt.Sprintf(uri, 1), noHeader, nil)

		want := []string{
			`<title hx-swap-oob="true">Edit Chicken Jersey | Recipya</title>`,
			`<input type="text" name="title" id="title" placeholder="Title of the recipe*" required autocomplete="off" value="Chicken Jersey" class="w-full py-2 bg-gray-50 font-semibold text-center text-gray-600 placeholder-gray-400 rounded-t-lg dark:bg-gray-900 dark:text-gray-200">`,
			`<img src="/data/images/` + baseRecipe.Image.String() + `.jpg" alt="Image preview of the recipe." class="h-full"><span><input type="file" accept="image/*" name="image" value="/data/images/` + baseRecipe.Image.String() + `.jpg" _="on dragover or dragenter halt the event then set the target's style.background to 'lightgray' on dragleave or drop set the target's style.background to '' on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from next <button/>">`,
			`<input type="text" list="categories" name="category" value="american" class="bg-transparent text-center focus:outline-none" placeholder="Breakfast*" required>`,
			`<input id="yield" type="number" min="1" name="yield" value="12" class="w-24 rounded bg-gray-100 p-2 dark:bg-gray-900">`,
			`<input id="source" type="text" placeholder="Source*" name="source" value="https://example.com/recipes/yummy" class="w-full rounded bg-gray-100 p-2 dark:bg-gray-900">`,
			`<textarea id="description" name="description" rows="10" placeholder="This Thai curry chicken will make you drool." required class="p-2 border border-gray-300 rounded-t-lg dark:bg-gray-800 dark:border-none" >A delicious recipe!</textarea>`,
			`<tbody class="text-right text-gray-500 bg-white divide-y divide-gray-200 dark:text-gray-200 dark:bg-slate-800 dark:divide-slate-600"><tr><td class="py-2 dark:border-gray-800"><p>Calories:</p></td><td class="text-center dark:border-gray-800"><label for="calories" class="hidden"></label><input type="text" id="calories" name="calories" autocomplete="off" placeholder="368kcal" value="354" class="w-3/4 p-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-800 dark:text-gray-200"></td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Total carbs:</p></td><td class="text-center dark:border-gray-800"><label for="total-carbohydrates" class="hidden"></label><input type="text" id="total-carbohydrates" name="total-carbohydrates" autocomplete="off" placeholder="35g" value="7g" class="w-3/4 p-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-800 dark:text-gray-200"></td></tr><tr><td class="py-2 dark:border-gray-800"><p>Sugars:</p></td><td class="text-center dark:border-gray-800"><label for="sugars" class="hidden"></label><input type="text" id="sugars" name="sugars" autocomplete="off" placeholder="3g" value="6g" class="w-3/4 p-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-800 dark:text-gray-200"></td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Protein:</p></td><td class="text-center dark:border-gray-800"><label for="protein" class="hidden"></label><input type="text" id="protein" name="protein" autocomplete="off" placeholder="21g" value="3g" class="w-3/4 p-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-800 dark:text-gray-200"></td></tr><tr><td class="py-2 dark:border-gray-800"><p>Total fat:</p></td><td class="text-center dark:border-gray-800"><label for="total-fat" class="hidden"></label><input type="text" id="total-fat" name="total-fat" autocomplete="off" placeholder="15g" value="8g" class="w-3/4 p-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-800 dark:text-gray-200"></td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Saturated fat:</p></td><td class="text-center dark:border-gray-800"><label for="saturated-fat" class="hidden"></label><input type="text" id="saturated-fat" name="saturated-fat" autocomplete="off" placeholder="1.8g" value="4g" class="w-3/4 p-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-800 dark:text-gray-200"></td></tr><tr><td class="py-2 dark:border-gray-800"><p>Cholesterol:</p></td><td class="text-center dark:border-gray-800"><label for="cholesterol" class="hidden"></label><input type="text" id="cholesterol" name="cholesterol" autocomplete="off" placeholder="1.1mg" value="1g" class="w-3/4 p-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-800 dark:text-gray-200"></td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Sodium:</p></td><td class="text-center dark:border-gray-800"><label for="sodium" class="hidden"></label><input type="text" id="sodium" name="sodium" autocomplete="off" placeholder="100mg" value="5g" class="w-3/4 p-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-800 dark:text-gray-200"></td></tr><tr><td class="py-2 dark:border-gray-800"><p>Fiber:</p></td><td class="text-center dark:border-gray-800"><label for="fiber" class="hidden"></label><input type="text" id="fiber" name="fiber" autocomplete="off" placeholder="8g" value="2g" class="w-3/4 p-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-800 dark:text-gray-200"></td></tr></tbody>`,
			`<input type="text" id="time-preparation" name="time-preparation" value="00:30:00" class="w-full p-1 placeholder-gray-400 bg-white border border-gray-400 html-duration-picker dark:bg-gray-900 dark:text-gray-200"></td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2">Cooking:</td><td class="w-1/2 text-center"><label for="time-cooking" class="hidden"></label><input type="text" id="time-cooking" name="time-cooking" value="01:00:00" class="w-full p-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 html-duration-picker dark:bg-gray-900 dark:text-gray-200" _="on load get <button/> in #times-table then add .dark:bg-gray-200 to it"></td></tr></tbody></table></div></div></div><div class="col-span-6 py-2 border-b border-x border-black md:col-span-2 md:border-r-0"><fieldset class="text-center border-none">`,
			`<li class="pb-2"><label><input type="text" name="ingredient-1" placeholder="Ingredient #1" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="ing1" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')"></label><button type="button" class="w-10 h-10 text-center bg-green-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-green-600 hover:text-white center dark:bg-green-500" title="Shortcut: Enter" hx-post="/recipes/add/manual/ingredient" hx-target="#ingredients-list" hx-swap="beforeend" hx-include="[name^='ingredient']"> + </button><button type="button" class="delete-button w-10 h-10 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/1" hx-include="[name^='ingredient']"> - </button><div class="inline-block h-4 cursor-move handle"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div></li><li class="pb-2"><label><input type="text" name="ingredient-2" placeholder="Ingredient #2" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="ing2" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')"></label><button type="button" class="w-10 h-10 text-center bg-green-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-green-600 hover:text-white center dark:bg-green-500" title="Shortcut: Enter" hx-post="/recipes/add/manual/ingredient" hx-target="#ingredients-list" hx-swap="beforeend" hx-include="[name^='ingredient']"> + </button><button type="button" class="delete-button w-10 h-10 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/2" hx-include="[name^='ingredient']"> - </button><div class="inline-block h-4 cursor-move handle"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div></li><li class="pb-2"><label><input type="text" name="ingredient-3" placeholder="Ingredient #3" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" value="ing3" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')"></label><button type="button" class="w-10 h-10 text-center bg-green-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-green-600 hover:text-white center dark:bg-green-500" title="Shortcut: Enter" hx-post="/recipes/add/manual/ingredient" hx-target="#ingredients-list" hx-swap="beforeend" hx-include="[name^='ingredient']"> + </button><button type="button" class="delete-button w-10 h-10 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/3" hx-include="[name^='ingredient']"> - </button><div class="inline-block h-4 cursor-move handle"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div></li>`,
			`<li class="pt-2 md:pl-0"><div class="flex"><label class="w-[95%]"><textarea required name="instruction-1" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #1" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">ins1</textarea></label><div class="grid"><div class="h-4 cursor-move handle grid place-content-center"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div><button type="button" class="md:w-7 md:h-7 md:right-auto w-10 h-10 text-center bg-green-300 border border-gray-800 rounded-lg hover:bg-green-600 hover:text-white center dark:bg-green-500" title="Shortcut: CTRL + Enter" hx-post="/recipes/add/manual/instruction" hx-target="#instructions-list" hx-swap="beforeend" hx-include="[name^='instruction']"> + </button><button type="button" class="delete-button w-10 h-10 md:w-7 md:h-7 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/1" hx-include="[name^='instruction']"> - </button></div></div></li><li class="pt-2 md:pl-0"><div class="flex"><label class="w-[95%]"><textarea required name="instruction-2" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #2" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">ins2</textarea></label><div class="grid"><div class="h-4 cursor-move handle grid place-content-center"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div><button type="button" class="md:w-7 md:h-7 md:right-auto w-10 h-10 text-center bg-green-300 border border-gray-800 rounded-lg hover:bg-green-600 hover:text-white center dark:bg-green-500" title="Shortcut: CTRL + Enter" hx-post="/recipes/add/manual/instruction" hx-target="#instructions-list" hx-swap="beforeend" hx-include="[name^='instruction']"> + </button><button type="button" class="delete-button w-10 h-10 md:w-7 md:h-7 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/2" hx-include="[name^='instruction']"> - </button></div></div></li><li class="pt-2 md:pl-0"><div class="flex"><label class="w-[95%]"><textarea required name="instruction-3" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #3" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">ins3</textarea></label><div class="grid"><div class="h-4 cursor-move handle grid place-content-center"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div><button type="button" class="md:w-7 md:h-7 md:right-auto w-10 h-10 text-center bg-green-300 border border-gray-800 rounded-lg hover:bg-green-600 hover:text-white center dark:bg-green-500" title="Shortcut: CTRL + Enter" hx-post="/recipes/add/manual/instruction" hx-target="#instructions-list" hx-swap="beforeend" hx-include="[name^='instruction']"> + </button><button type="button" class="delete-button w-10 h-10 md:w-7 md:h-7 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/3" hx-include="[name^='instruction']"> - </button></div></div></li>`,
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
			Keywords:     nil,
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
			assertHeader(t, rr, "HX-Trigger", fmt.Sprintf(`{"showToast":"{\"message\":\"%s\",\"backgroundColor\":\"bg-red-500\"}"}`, tc.want))

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
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Recipe not found.\",\"backgroundColor\":\"bg-red-500\"}"}`)
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
			`<label for="ingredient-0"></label><input type="checkbox" id="ingredient-0" class="mt-1"><span class="pl-2">4 lb chicken</span>`,
			`<label for="ingredient-1"></label><input type="checkbox" id="ingredient-1" class="mt-1"><span class="pl-2">1 cup bread loaf</span>`,
			`<label for="ingredient-2"></label><input type="checkbox" id="ingredient-2" class="mt-1"><span class="pl-2">1 tbsp beef broth</span>`,
			`<label for="ingredient-3"></label><input type="checkbox" id="ingredient-3" class="mt-1"><span class="pl-2">15 cups flour</span>`,
			`<label for="ingredient-4"></label><input type="checkbox" id="ingredient-4" class="mt-1"><span class="pl-2">4 big apples</span>`,
			`<label for="ingredient-5"></label><input type="checkbox" id="ingredient-5" class="mt-1"><span class="pl-2">Lots of big apples</span>`,
			`<label for="ingredient-6"></label><input type="checkbox" id="ingredient-6" class="mt-1"><span class="pl-2">5 slices of bacon</span>`,
			`<label for="ingredient-7"></label><input type="checkbox" id="ingredient-7" class="mt-1"><span class="pl-2">4 2/3 cans of bamboo sticks</span>`,
			`<label for="ingredient-8"></label><input type="checkbox" id="ingredient-8" class="mt-1"><span class="pl-2">3 can of tomato paste</span>`,
			`<label for="ingredient-9"></label><input type="checkbox" id="ingredient-9" class="mt-1"><span class="pl-2">13 1/2 peanut butter jars</span>`,
			`<label for="ingredient-10"></label><input type="checkbox" id="ingredient-10" class="mt-1"><span class="pl-2">15 mL of whiskey</span>`,
			`<label for="ingredient-11"></label><input type="checkbox" id="ingredient-11" class="mt-1"><span class="pl-2">1 1/3 tbsp lemon juice</span>`,
			`<label for="ingredient-12"></label><input type="checkbox" id="ingredient-12" class="mt-1"><span class="pl-2">Ground ginger</span>`,
			`<label for="ingredient-13"></label><input type="checkbox" id="ingredient-13" class="mt-1"><span class="pl-2">6 Large or 8 medium ripe Hass avocados</span>`,
			`<label for="ingredient-14"></label><input type="checkbox" id="ingredient-14" class="mt-1"><span class="pl-2">1/2 tsp salt plus more for seasoning</span>`,
			`<label for="ingredient-15"></label><input type="checkbox" id="ingredient-15" class="mt-1"><span class="pl-2">1 fresh pineapple, cored and cut into 3-inch pieces</span>`,
			`<label for="ingredient-16"></label><input type="checkbox" id="ingredient-16" class="mt-1"><span class="pl-2">Un sac de chips de 1kg</span>`,
			`<label for="ingredient-17"></label><input type="checkbox" id="ingredient-17" class="mt-1"><span class="pl-2">4 15-ounce can Goya beans</span>`,
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
			`<article class="grid gap-4 p-4 text-sm grid-cols-4 md:m-auto md:max-w-7xl md:text-base md:grid-cols-10">`,
			`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the Chinese Firmware recipe" hx-get="/recipes/1" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">Chinese Firmware</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"></span></div><p class="pb-8 text-justify text-sm"></p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/1" hx-target="#content" hx-push-url="true"> View </button></div>`,
			`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the Lovely Canada recipe" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">Lovely Canada</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"></span></div><p class="pb-8 text-justify text-sm"></p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"> View </button></div>`,
			`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the Lovely Ukraine recipe" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">Lovely Ukraine</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"></span></div><p class="pb-8 text-justify text-sm"></p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"> View </button></div>`,
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
				`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the Lovely Canada recipe" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">Lovely Canada</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"></span></div><p class="pb-8 text-justify text-sm"></p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"> View </button></div>`,
				`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the Lovely Ukraine recipe" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">Lovely Ukraine</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"></span></div><p class="pb-8 text-justify text-sm"></p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"> View </button></div>`,
			},
		},
		{
			query: "chi",
			want: []string{
				`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the Chinese Firmware recipe" hx-get="/recipes/1" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">Chinese Firmware</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"></span></div><p class="pb-8 text-justify text-sm"></p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/1" hx-target="#content" hx-push-url="true"> View </button></div>`,
			},
		},
		{
			query: "lovely",
			want: []string{
				`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the Lovely Canada recipe" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">Lovely Canada</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"></span></div><p class="pb-8 text-justify text-sm"></p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/2" hx-target="#content" hx-push-url="true"> View </button></div>`,
				`<div class="relative col-span-2 bg-white rounded-lg shadow-lg dark:bg-neutral-700"><img class="rounded-t-lg cursor-pointer w-full hover:opacity-80 border-b dark:border-b-gray-800 h-48 px-4 pt-4" src="" alt="Image for the Lovely Ukraine recipe" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"><div class="py-2 px-3"><div class="grid grid-flow-col gap-1 pb-2"><p class="font-semibold">Lovely Ukraine</p><span class="grid place-content-center text-xs text-center font-medium select-none p-[0.25rem] bg-indigo-700 text-white rounded-lg h-fit"></span></div><p class="pb-8 text-justify text-sm"></p><button class="absolute bottom-2 border-2 border-gray-800 rounded-lg center w-[90%] hover:bg-gray-800 hover:text-white dark:border-gray-800 hover:dark:bg-neutral-600" hx-get="/recipes/3" hx-target="#content" hx-push-url="true"> View </button></div>`,
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

	app.Config.URL = "https://www.recipya.com"
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
	_, _ = srv.Repository.AddRecipe(recipe, 1)
	link, _ := srv.Repository.AddShareLink(models.Share{RecipeID: 1, CookbookID: -1, UserID: 1})

	t.Run("create valid share link", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri(1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<input type="url" value="example.com/r/33320755-82f9-47e5-bb0a-d1b55cbd3f7b" class="w-full rounded-lg bg-gray-100 px-4 py-2" readonly="readonly">`,
			`<button class="w-24 font-semibold p-2 bg-gray-300 rounded-lg hover:bg-gray-400" title="Copy to clipboard" _="on click if window.navigator.clipboard then call navigator.clipboard.writeText('example.com/r/33320755-82f9-47e5-bb0a-d1b55cbd3f7b') then put 'Copied!' into me then add @title='Copied to clipboard!' then toggle @disabled on me then toggle .cursor-not-allowed .bg-green-600 .text-white .hover:bg-gray-400 on me else call alert('Your browser does not support the clipboard feature. Please copy the link manually.') end"> Copy </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("create invalid share link", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri(10), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Failed to create share link.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("access share link anonymous", func(t *testing.T) {
		rr := sendRequest(srv, http.MethodGet, link, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
			`<button class="ml-2" title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release()`,
			`<div class="flex"><a href="/auth/login" class="flex-auto mr-2 rounded-lg px-2 py-1 hover:text-white hover:bg-green-600"> Log In </a><a href="/auth/register" class="flex-auto mr-4 rounded-lg px-2 py-1 bg-amber-300 dark:bg-orange-600 hover:text-white hover:bg-red-600"> Sign Up </a></div>`,
			`<h1 class="grid col-span-4 py-2 font-semibold place-content-center">Chicken Jersey</h1>`,
			`<div class="grid justify-end col-span-1 place-content-center print:invisible"><button class="mr-2" onclick="print()" title="Print recipe"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"/></svg></button></div>`,
			`<img id="output" class="w-full text-center h-96" alt="Image of the recipe" style="object-fit: scale-down" src="/data/images/` + recipe.Image.String() + `.jpg">`,
			`<span class="text-sm font-normal leading-none">American</span>`,
			`<p class="text-sm text-center">2 servings</p>`,
			`<a class="p-1 border rounded-lg center hover:bg-gray-800 hover:text-white dark:border-gray-800" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank"> Source </a>`,
			`<p class="p-2 text-sm whitespace-pre-line">This is the most delicious recipe!</p>`,
			`<p class="text-xs">Per 100g: 500 calories; total carbohydrates 7g; sugars 6g; protein 3g; total fat 8g; saturated fat 4g; cholesterol 1g; sodium 5g; fiber 2g.</p>`,
			`<table class="min-w-full text-xs divide-y divide-gray-200 table-auto dark:text-gray-200 dark:divide-slate-600"><thead class="h-12 font-medium tracking-wider text-white bg-gray-800"><tr><th scope="col" class="text-right"><p class="text-center">Nutrition<br>(per 100g)</p></th><th scope="col" class="text-center"><p>Amount<br>(optional)</p></th></tr></thead><tbody class="text-right text-gray-500 bg-white divide-y divide-gray-200 dark:text-gray-200 dark:bg-slate-800 dark:divide-slate-600"><tr><td class="py-2 dark:border-gray-800"><p>Calories:</p></td><td class="text-center dark:border-gray-800"> 500 </td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Total carbs:</p></td><td class="text-center dark:border-gray-800"> 7g </td></tr><tr><td class="py-2 dark:border-gray-800"><p>Sugars:</p></td><td class="text-center dark:border-gray-800"> 6g </td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Protein:</p></td><td class="text-center dark:border-gray-800"> 3g </td></tr><tr><td class="py-2 dark:border-gray-800"><p>Total fat:</p></td><td class="text-center dark:border-gray-800"> 8g </td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Saturated fat:</p></td><td class="text-center dark:border-gray-800"> 4g </td></tr><tr><td class="py-2 dark:border-gray-800"><p>Cholesterol:</p></td><td class="text-center dark:border-gray-800"> 1g </td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Sodium:</p></td><td class="text-center dark:border-gray-800"> 5g </td></tr><tr><td class="py-2 dark:border-gray-800"><p>Fiber:</p></td><td class="text-center dark:border-gray-800"> 2g </td></tr></tbody></table></div></div></div><div class="col-span-3 border-b border-r border-black md:col-span-1 print:col-span-1"><div class="inline-block min-w-full overflow-x-auto align-middle"><div class="overflow-hidden border-gray-200">`,
			`<table class="min-w-full text-xs divide-y divide-gray-200 dark:text-gray-200 dark:divide-slate-600"><thead class="h-12 font-medium tracking-wider text-white bg-gray-800 border-l dark:border-l-gray-600 print:h-1"><tr><th scope="col" class="text-right">Time</th><th scope="col" class="text-center">h:m:s</th></tr></thead><tbody class="text-right text-gray-500 bg-white divide-y divide-gray-200 dark:text-gray-200 dark:bg-slate-800 dark:divide-slate-600"><tr><td class="py-2">Preparation:</td><td class="text-center print:py-0"><time datetime="PT05M">5m</time></td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2">Cooking:</td><td class="w-1/2 text-center print:py-0"><time datetime="PT1H05M">1h05m</time></td></tr><tr class="font-semibold"><td class="py-2">Total:</td><td class="w-1/2 text-center print:py-0"><time datetime="PT1H10M">1h10m</time></td></tr></tbody></table>`,
			`<div class="col-span-6 py-2 border-b border-x border-black md:col-span-2 md:border-r-0 dark:border-gray-800 print:hidden"><h2 class="pb-2 m-auto text-sm font-semibold text-center text-gray-600 underline dark:text-gray-200"> Ingredients </h2><ul class="grid px-6 columns-2"><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200 dark:hover:bg-gray-800" _="on mousedown toggle .line-through then toggle @checked on first <input/> in me"><label for="ingredient-0"></label><input type="checkbox" id="ingredient-0" class="mt-1"><span class="pl-2">Ing1</span></li><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200 dark:hover:bg-gray-800" _="on mousedown toggle .line-through then toggle @checked on first <input/> in me"><label for="ingredient-1"></label><input type="checkbox" id="ingredient-1" class="mt-1"><span class="pl-2">Ing2</span></li><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200 dark:hover:bg-gray-800" _="on mousedown toggle .line-through then toggle @checked on first <input/> in me"><label for="ingredient-2"></label><input type="checkbox" id="ingredient-2" class="mt-1"><span class="pl-2">Ing3</span></li></ul></div><div class="col-span-6 py-2 pb-8 border-b border-x border-black rounded-bl-lg md:rounded-bl-none md:col-span-4 dark:border-gray-800 print:hidden"><h2 class="pb-2 m-auto text-sm font-semibold text-center text-gray-600 underline dark:text-gray-200"> Instructions</h2><ol class="grid px-6 list-decimal"><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200 dark:hover:bg-gray-800" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins1</span></li><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200 dark:hover:bg-gray-800" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins2</span></li><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200 dark:hover:bg-gray-800" _="on mousedown toggle .line-through"><span class="whitespace-pre-line">Ins3</span></li></ol></div>`,
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
			rr := tc.sendFunc(srv, http.MethodGet, link, noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
				`<button class="ml-2" title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release()`,
				`<button class="mr-2" title="Add recipe to collection" hx-post="/recipes/1/share/add"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg></button>`,
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
		t.Run("access share link logged "+tc.name, func(t *testing.T) {
			rr := tc.sendFunc(srv, http.MethodGet, link, noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
				`<button class="ml-2" title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release()`,
				`<button class="mr-2" onclick="print()" title="Print recipe">`,
				`<button class="mr-2" hx-delete="/recipes/1" hx-swap="none" title="Delete recipe" hx-confirm="Are you sure you wish to delete this recipe?" hx-indicator="#fullscreen-loader">`,
			}
			notWant := []string{
				`<dialog id="share-dialog" class="p-4 border-4 border-black min-w-[15rem]"><div id="share-dialog-result" class="pb-4"></div>`,
				`<button class="mr-2" title="Add recipe to collection" hx-post="/recipes/1/share/add"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg></button>`,
				`<button class="mr-2" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call #share-dialog.showModal()">`,
				`<button class="mr-2" title="Add recipe to collection" hx-post="/recipes/1/share/add"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg></button>`,
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
			`<title hx-swap-oob="true">Not Found | Recipya</title>`,
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
			_, _ = srv.Repository.AddRecipe(r, 1)

			rr := tc.sendFunc(srv, http.MethodGet, uri+"/1", noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + r.Name + " | Recipya</title>",
				`<button class="ml-2" title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d='M 12.276 18.55 v -0.748 a 4.79 4.79 0 0 1 1.463 -3.458 a 5.763 5.763 0 0 0 1.804 -4.21 a 5.821 5.821 0 0 0 -6.475 -5.778 c -2.779 0.307 -4.99 2.65 -5.146 5.448 a 5.82 5.82 0 0 0 1.757 4.503 a 4.906 4.906 0 0 1 1.5 3.495 v 0.747 a 1.44 1.44 0 0 0 1.44 1.439 h 2.218 a 1.44 1.44 0 0 0 1.44 -1.439 z m -1.058 0 c 0 0.209 -0.17 0.38 -0.38 0.38 h -2.22 c -0.21 0 -0.38 -0.171 -0.38 -0.38 v -0.748 c 0 -1.58 -0.664 -3.13 -1.822 -4.254 A 4.762 4.762 0 0 1 4.98 9.863 c 0.127 -2.289 1.935 -4.204 4.205 -4.455 a 4.762 4.762 0 0 1 5.3 4.727 a 4.714 4.714 0 0 1 -1.474 3.443 a 5.853 5.853 0 0 0 -1.791 4.225 v 0.746 z M 11.45 20.51 H 8.006 a 0.397 0.397 0 1 0 0 0.795 h 3.444 a 0.397 0.397 0 1 0 0 -0.794 z M 11.847 22.162 a 0.397 0.397 0 0 0 -0.397 -0.397 H 8.006 a 0.397 0.397 0 1 0 0 0.794 h 3.444 c 0.22 0 0.397 -0.178 0.397 -0.397 z z z z z z z z M 10.986 23.416 H 8.867 a 0.397 0.397 0 1 0 0 0.794 h 1.722 c 0.22 0 0.397 -0.178 0.397 -0.397 z' to #icon-bulb else call initWakeLock() then add @d='M12.276 18.55v-.748a4.79 4.79 0 0 1 1.463-3.458 5.763 5.763 0 0 0 1.804-4.21 5.821 5.821 0 0 0-6.475-5.778c-2.779.307-4.99 2.65-5.146 5.448a5.82 5.82 0 0 0 1.757 4.503 4.906 4.906 0 0 1 1.5 3.495v.747a1.44 1.44 0 0 0 1.44 1.439h2.218a1.44 1.44 0 0 0 1.44-1.439zm-1.058 0c0 .209-.17.38-.38.38h-2.22c-.21 0-.38-.171-.38-.38v-.748c0-1.58-.664-3.13-1.822-4.254A4.762 4.762 0 0 1 4.98 9.863c.127-2.289 1.935-4.204 4.205-4.455a4.762 4.762 0 0 1 5.3 4.727 4.714 4.714 0 0 1-1.474 3.443 5.853 5.853 0 0 0-1.791 4.225v.746zM11.45 20.51H8.006a.397.397 0 1 0 0 .795h3.444a.397.397 0 1 0 0-.794zM11.847 22.162a.397.397 0 0 0-.397-.397H8.006a.397.397 0 1 0 0 .794h3.444c.22 0 .397-.178.397-.397zM.397 10.125h2.287a.397.397 0 1 0 0-.794H.397a.397.397 0 1 0 0 .794zM19.456 9.728a.397.397 0 0 0-.397-.397h-2.287a.397.397 0 1 0 0 .794h2.287c.22 0 .397-.178.397-.397zM9.331.397v2.287a.397.397 0 1 0 .794 0V.397a.397.397 0 1 0-.794 0zM16.045 2.85 14.43 4.465a.397.397 0 1 0 .561.561l1.617-1.617a.397.397 0 1 0-.562-.56zM5.027 14.429a.397.397 0 0 0-.56 0l-1.618 1.616a.397.397 0 1 0 .562.562l1.617-1.617a.397.397 0 0 0 0-.561zM4.466 5.027a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.56L3.41 2.848a.397.397 0 1 0-.561.561zM16.045 16.607a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.562L14.99 14.43a.397.397 0 1 0-.561.56zM10.986 23.416a.397.397 0 0 0-.397-.397H8.867a.397.397 0 1 0 0 .794h1.722c.22 0 .397-.178.397-.397z' to #icon-bulb end"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor"><path id="icon-bulb" d="M12.276 18.55v-.748a4.79 4.79 0 0 1 1.463-3.458 5.763 5.763 0 0 0 1.804-4.21 5.821 5.821 0 0 0-6.475-5.778c-2.779.307-4.99 2.65-5.146 5.448a5.82 5.82 0 0 0 1.757 4.503 4.906 4.906 0 0 1 1.5 3.495v.747a1.44 1.44 0 0 0 1.44 1.439h2.218a1.44 1.44 0 0 0 1.44-1.439zm-1.058 0c0 .209-.17.38-.38.38h-2.22c-.21 0-.38-.171-.38-.38v-.748c0-1.58-.664-3.13-1.822-4.254A4.762 4.762 0 0 1 4.98 9.863c.127-2.289 1.935-4.204 4.205-4.455a4.762 4.762 0 0 1 5.3 4.727 4.714 4.714 0 0 1-1.474 3.443 5.853 5.853 0 0 0-1.791 4.225v.746zM11.45 20.51H8.006a.397.397 0 1 0 0 .795h3.444a.397.397 0 1 0 0-.794zM11.847 22.162a.397.397 0 0 0-.397-.397H8.006a.397.397 0 1 0 0 .794h3.444c.22 0 .397-.178.397-.397zM.397 10.125h2.287a.397.397 0 1 0 0-.794H.397a.397.397 0 1 0 0 .794zM19.456 9.728a.397.397 0 0 0-.397-.397h-2.287a.397.397 0 1 0 0 .794h2.287c.22 0 .397-.178.397-.397zM9.331.397v2.287a.397.397 0 1 0 .794 0V.397a.397.397 0 1 0-.794 0zM16.045 2.85 14.43 4.465a.397.397 0 1 0 .561.561l1.617-1.617a.397.397 0 1 0-.562-.56zM5.027 14.429a.397.397 0 0 0-.56 0l-1.618 1.616a.397.397 0 1 0 .562.562l1.617-1.617a.397.397 0 0 0 0-.561zM4.466 5.027a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.56L3.41 2.848a.397.397 0 1 0-.561.561zM16.045 16.607a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.562L14.99 14.43a.397.397 0 1 0-.561.56zM10.986 23.416a.397.397 0 0 0-.397-.397H8.867a.397.397 0 1 0 0 .794h1.722c.22 0 .397-.178.397-.397z"/></svg></button>`,
				`<button class="ml-2" title="Edit recipe" hx-get="/recipes/1/edit" hx-push-url="true" hx-target="#content">`,
				`<h1 class="grid col-span-4 py-2 font-semibold place-content-center">Chicken Jersey</h1>`,
				`<button class="mr-2" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call #share-dialog.showModal()">`,
				`<button class="mr-2" onclick="print()" title="Print recipe">`,
				`<button class="mr-2" hx-delete="/recipes/1" hx-swap="none" title="Delete recipe" hx-confirm="Are you sure you wish to delete this recipe?" hx-indicator="#fullscreen-loader">`,
				`<dialog id="share-dialog" class="p-4 border-4 border-black min-w-[15rem]"><div id="share-dialog-result" class="pb-4"></div>`,
				`<img id="output" class="w-full text-center h-96" alt="Image of the recipe" style="object-fit: scale-down" src="/data/images/` + r.Image.String() + `.jpg">`,
				`<span class="text-sm font-normal leading-none">American</span>`,
				`<div class="grid col-span-2 py-2 border-black place-items-center text-sm md:col-span-3 md:border-t dark:border-gray-800"><form autocomplete="off" _="on submit halt the event"><input id="yield" type="number" min="1" name="yield" value="2" class="w-16 text-center rounded bg-gray-100 p-2 grid self-end dark:bg-gray-900" hx-get="/recipes/1/scale" hx-trigger="input" hx-target="#ingredients-instructions-container"><label for="yield" class="grid self-start">servings</label></form></div>`,
				`<a class="p-1 border rounded-lg center hover:bg-gray-800 hover:text-white dark:border-gray-800" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank"> Source </a>`,
				`<p class="p-2 text-sm whitespace-pre-line">This is the most delicious recipe!</p>`,
				`<p class="text-xs">Per 100g: 500 calories; total carbohydrates 7g; sugars 6g; protein 3g; total fat 8g; saturated fat 4g; cholesterol 1g; sodium 5g; fiber 2g.</p>`,
				`<table class="min-w-full text-xs divide-y divide-gray-200 table-auto dark:text-gray-200 dark:divide-slate-600"><thead class="h-12 font-medium tracking-wider text-white bg-gray-800"><tr><th scope="col" class="text-right"><p class="text-center">Nutrition<br>(per 100g)</p></th><th scope="col" class="text-center"><p>Amount<br>(optional)</p></th></tr></thead><tbody class="text-right text-gray-500 bg-white divide-y divide-gray-200 dark:text-gray-200 dark:bg-slate-800 dark:divide-slate-600"><tr><td class="py-2 dark:border-gray-800"><p>Calories:</p></td><td class="text-center dark:border-gray-800"> 500 </td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Total carbs:</p></td><td class="text-center dark:border-gray-800"> 7g </td></tr><tr><td class="py-2 dark:border-gray-800"><p>Sugars:</p></td><td class="text-center dark:border-gray-800"> 6g </td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Protein:</p></td><td class="text-center dark:border-gray-800"> 3g </td></tr><tr><td class="py-2 dark:border-gray-800"><p>Total fat:</p></td><td class="text-center dark:border-gray-800"> 8g </td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Saturated fat:</p></td><td class="text-center dark:border-gray-800"> 4g </td></tr><tr><td class="py-2 dark:border-gray-800"><p>Cholesterol:</p></td><td class="text-center dark:border-gray-800"> 1g </td></tr><tr class="bg-gray-100 dark:bg-slate-700"><td class="py-2 dark:border-gray-800"><p>Sodium:</p></td><td class="text-center dark:border-gray-800"> 5g </td></tr><tr><td class="py-2 dark:border-gray-800"><p>Fiber:</p></td><td class="text-center dark:border-gray-800"> 2g </td></tr></tbody></table>`,
				`<td class="py-2">Preparation:</td><td class="text-center print:py-0"><time datetime="PT05M">5m</time></td></tr><tr class="bg-gray-100 dark:bg-slate-700">`,
				`<td class="py-2">Cooking:</td><td class="w-1/2 text-center print:py-0"><time datetime="PT1H05M">1h05m</time></td></tr><tr class="font-semibold">`,
				`<td class="py-2">Total:</td><td class="w-1/2 text-center print:py-0"><time datetime="PT1H10M">1h10m</time></td>`,
				`<div class="hidden print:grid col-span-6 ml-2 mt-2 mb-4"><h1 class="text-sm"><b>Ingredients</b></h1><ol class="col-span-6 w-full" style="column-count: 1"><li class="text-sm"><label><input type="checkbox"></label><span class="pl-2">Ing1</span></li><li class="text-sm"><label><input type="checkbox"></label><span class="pl-2">Ing2</span></li><li class="text-sm"><label><input type="checkbox"></label><span class="pl-2">Ing3</span></li></ol></div>`,
				`<div class="hidden print:grid col-span-5"><h1 class="text-sm"><b>Instructions</b></h1><ol class="col-span-6 list-decimal w-full ml-6" style="column-count: 2; column-gap: 2.5rem"><li><span class="text-sm whitespace-pre-line">Ins1</span></li><li><span class="text-sm whitespace-pre-line">Ins2</span></li><li><span class="text-sm whitespace-pre-line">Ins3</span></li></ol></div>`,
				`<script> var wakeLock = null; initWakeLock(); function initWakeLock() { navigator.wakeLock?.request("screen") .then((lock) => { wakeLock = lock; wakeLock.onrelease = () => { wakeLock = null; console.info("Screen lock deactivated."); } console.info("Screen lock activated."); }) .catch((err) => { ` + "console.log(`Screen lock error: ${err.name}, ${err.message}`)" + ` }); } </script>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}
