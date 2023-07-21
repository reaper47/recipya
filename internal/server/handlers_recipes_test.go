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

func TestHandlers_Recipes_New(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/add"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
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
				`<button class="flex justify-center w-full duration-300 border-2 border-gray-800 rounded-lg hover:bg-gray-800 hover:text-white center" hx-target="#content" hx-push-url="/recipes/add/unsupported-website" hx-prompt="Enter the recipe's URL" hx-post="/recipes/add/website" hx-indicator="#fullscreen-loader"> Fetch </button>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}

func TestHandlers_Recipes_AddManual(t *testing.T) {
	repo := &mockRepository{}
	srv := server.NewServer(repo, &mockEmail{}, &mockFiles{})

	uri := "/recipes/add/manual"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
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
				`<form hx-post="/recipes/add/manual" enctype="multipart/form-data" class="grid max-w-6xl grid-cols-6 px-4 pb-4 m-auto mt-4">`,
				`<input type="text" name="title" id="title" placeholder="Title of the recipe*" required class="w-full py-2 font-bold text-center text-gray-600 placeholder-gray-400 rounded-t-lg">`,
				`<label class="grid col-span-6 w-full h-96 border-b border-l border-r border-black place-content-center md:border-r-0 md:col-span-4 text-sm"><img src="" alt="Image preview of the recipe." class="h-full"><span><input type="file" accept="image/*" name="image" required _="on dragover or dragenter halt the event then set the target's style.background to 'lightgray' on dragleave or drop set the target's style.background to '' on drop or change make an FileReader called reader then if event.dataTransfer get event.dataTransfer.files[0] else get event.target.files[0] end then set {src: window.URL.createObjectURL(it)} on previous <img/> then remove .hidden from next <button/>"><button type="button" class="px-2 bg-red-300 border border-gray-800 rounded-lg hover:bg-red-600 hover:text-white hidden" _="on click set {value: ''} on previous <input/> then set {src: ''} on previous <img/> then add .hidden"> Delete </button></span></label>`,
				`<div class="grid col-span-2 py-2 border-black place-items-center text-sm md:col-span-3 md:border-t"><div><label for="yield">Yields</label><input id="yield" type="number" min="1" name="yield" value="4" class="w-24 rounded bg-gray-100 p-2"><span>servings</span></div></div>`,
				`<label class="grid place-content-center p-2 font-medium text-sm text-blue-700 bg-blue-100 border border-blue-300 rounded-full w-fit"><input type="text" list="categories" name="category" class="bg-transparent text-center" placeholder="Breakfast*" required><datalist id="categories"><option>breakfast</option><option>lunch</option><option>dinner</option></datalist></label>`,
				`<textarea id="description" name="description" rows="10" class="p-2 border border-gray-300 rounded-t-lg" placeholder="This Thai curry chicken will make you drool..." required></textarea>`,
				`<th scope="col" class="text-right md:text-center"><p>Nutrition<br>(per 100g)</p></th><th scope="col" class="text-center"><p>Amount<br>(optional)</p></th>`,
				`<th scope="col" class="text-right">Time</th><th scope="col" class="text-center">h:m:s</th>`,
				`<ol id="ingredients-list" class="pl-6 list-decimal">`,
				`<ol id="instructions-list" class="pl-4 list-decimal">`,
				`<input type="text" name="ingredient-1" placeholder="Ingredient #1" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" onkeydown="handleKeyDownIngredient(event)">`,
				`<button type="submit" class="col-span-6 p-2 font-semibold text-white duration-300 bg-blue-500 hover:bg-blue-800"> Submit </button>`,
				`<script src="https://cdn.jsdelivr.net/npm/sortablejs@latest/Sortable.min.js"></script>`,
				`<script src="https://cdn.jsdelivr.net/npm/html-duration-picker@latest/dist/html-duration-picker.min.js"></script>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}

	t.Run("submit recipe", func(t *testing.T) {
		repo = &mockRepository{
			Recipes: make(map[int64]models.Recipes, 0),
		}
		srv.Repository = repo
		originalNumRecipes := len(repo.Recipes)

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
		id := int64(len(repo.Recipes))
		if len(repo.Recipes) != originalNumRecipes+1 {
			t.Fatal("expected one more recipe to be added to the database")
		}
		gotRecipe := repo.Recipes[1][id-1]
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
	srv := server.NewServer(&mockRepository{}, &mockEmail{}, &mockFiles{})

	uri := "/recipes/add/manual/ingredient"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
	})

	t.Run("does not yield input when previous input empty", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("ingredient-1="))

		assertStatus(t, rr.Code, http.StatusUnprocessableEntity)
	})

	t.Run("yields new ingredient input", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("ingredient-1=ingredient1"))

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<input autofocus type="text" name="ingredient-2" placeholder="Ingredient #2" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" onkeydown="handleKeyDownIngredient(event)">`,
			`&nbsp;<button type="button" class="delete-button w-10 h-10 duration-300 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/2" hx-include="[name^='ingredient']">-</button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_AddManualIngredientDelete(t *testing.T) {
	srv := server.NewServer(&mockRepository{}, &mockEmail{}, &mockFiles{})

	uri := "/recipes/add/manual/ingredient/"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
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
				`<input type="text" name="ingredient-1" placeholder="Ingredient #1" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" value="one" onkeydown="handleKeyDownIngredient(event)">`,
				`<input type="text" name="ingredient-2" placeholder="Ingredient #2" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" value="two" onkeydown="handleKeyDownIngredient(event)">`,
				`<input type="text" name="ingredient-3" placeholder="Ingredient #3" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" value="three" onkeydown="handleKeyDownIngredient(event)">`,
			},
		},
		{
			name:  "delete first entry",
			entry: 1,
			want: []string{
				`<input type="text" name="ingredient-1" placeholder="Ingredient #1" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" value="two" onkeydown="handleKeyDownIngredient(event)">`,
				`<input type="text" name="ingredient-2" placeholder="Ingredient #2" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" value="three" onkeydown="handleKeyDownIngredient(event)">`,
				`<input type="text" name="ingredient-3" placeholder="Ingredient #3" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" value="''" onkeydown="handleKeyDownIngredient(event)">`,
			},
		},
		{
			name:  "delete middle entry",
			entry: 3,
			want: []string{
				`<input type="text" name="ingredient-1" placeholder="Ingredient #1" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" value="one" onkeydown="handleKeyDownIngredient(event)"><`,
				`<input type="text" name="ingredient-2" placeholder="Ingredient #2" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" value="two" onkeydown="handleKeyDownIngredient(event)">`,
				`<input type="text" name="ingredient-3" placeholder="Ingredient #3" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" value="''" onkeydown="handleKeyDownIngredient(event)">`,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+strconv.Itoa(tc.entry), formHeader, strings.NewReader("ingredient-1=one&ingredient-2=two&ingredient-3=three&ingredient-4=''"))

			assertStatus(t, rr.Code, http.StatusOK)
			want := append(tc.want, []string{
				`&nbsp;<button type="button" class="delete-button w-10 h-10 duration-300 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/1" hx-include="[name^='ingredient']">-</button>`,
				`&nbsp;<button type="button" class="delete-button w-10 h-10 duration-300 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/2" hx-include="[name^='ingredient']">-</button>`,
				`&nbsp;<button type="button" class="delete-button w-10 h-10 duration-300 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/3" hx-include="[name^='ingredient']">-</button>`,
			}...)
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}

func TestHandlers_Recipes_AddManualInstruction(t *testing.T) {
	srv := server.NewServer(&mockRepository{}, &mockEmail{}, &mockFiles{})

	uri := "/recipes/add/manual/instruction"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
	})

	t.Run("does not yield input when previous input empty", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("instruction-1="))

		assertStatus(t, rr.Code, http.StatusUnprocessableEntity)
	})

	t.Run("yields new instruction input", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("instruction-1=one"))

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<textarea autofocus required name="instruction-2" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #2" onkeydown="handleKeyDownInstruction(event)"></textarea>`,
			`<button type="button" class="delete-button mt-4 md:flex-initial w-10 h-10 right-0.5 md:w-7 md:h-7 md:right-auto duration-300 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/2" hx-include="[name^='instruction']">-</button>`,
			`<button type="button" class="md:flex-initial bottom-0 right-0.5 md:w-7 md:h-7 md:right-auto w-10 h-10 text-center duration-300 bg-green-300 border border-gray-800 rounded-lg hover:bg-green-600 hover:text-white center" title="Shortcut: CTRL + Enter" hx-post="/recipes/add/manual/instruction" hx-target="#instructions-list" hx-swap="beforeend" hx-include="[name^='instruction']">+</button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})
}

func TestHandlers_Recipes_AddManualInstructionDelete(t *testing.T) {
	srv := server.NewServer(&mockRepository{}, &mockEmail{}, &mockFiles{})

	uri := "/recipes/add/manual/instruction/"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
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
				`<textarea required name="instruction-1" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #1" onkeydown="handleKeyDownInstruction(event)">One</textarea>`,
				`<textarea required name="instruction-2" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #2" onkeydown="handleKeyDownInstruction(event)">Two</textarea>`,
				`<textarea required name="instruction-3" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #3" onkeydown="handleKeyDownInstruction(event)">Three</textarea>`,
			},
		},
		{
			name:  "delete first entry",
			entry: 1,
			want: []string{
				`<textarea required name="instruction-1" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #1" onkeydown="handleKeyDownInstruction(event)">Two</textarea>`,
				`<textarea required name="instruction-2" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #2" onkeydown="handleKeyDownInstruction(event)">Three</textarea>`,
				`<textarea required name="instruction-3" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #3" onkeydown="handleKeyDownInstruction(event)">''</textarea>`,
			},
		},
		{
			name:  "delete middle entry",
			entry: 3,
			want: []string{
				`<textarea required name="instruction-1" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #1" onkeydown="handleKeyDownInstruction(event)">One</textarea>`,
				`<textarea required name="instruction-2" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #2" onkeydown="handleKeyDownInstruction(event)">Two</textarea>`,
				`<textarea required name="instruction-3" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #3" onkeydown="handleKeyDownInstruction(event)">''</textarea>`,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri+strconv.Itoa(tc.entry), formHeader, strings.NewReader("instruction-1=One&instruction-2=Two&instruction-3=Three&instruction-4=''"))

			assertStatus(t, rr.Code, http.StatusOK)
			want := append(tc.want, []string{
				`<button type="button" class="mt-4 md:flex-initial w-10 h-10 right-0.5 md:w-7 md:h-7 md:right-auto duration-300 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/1" hx-include="[name^='instruction']">-</button>`,
				`<button type="button" class="mt-4 md:flex-initial w-10 h-10 right-0.5 md:w-7 md:h-7 md:right-auto duration-300 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/2" hx-include="[name^='instruction']">-</button>`,
				`<button type="button" class="mt-4 md:flex-initial w-10 h-10 right-0.5 md:w-7 md:h-7 md:right-auto duration-300 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/3" hx-include="[name^='instruction']">-</button>`,
			}...)
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}

func TestHandlers_Recipes_AddRequestWebsite(t *testing.T) {
	emailMock := &mockEmail{}
	srv := server.NewServer(&mockRepository{}, emailMock, &mockFiles{})

	uri := "/recipes/add/request-website"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
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
		assertMustBeLoggedIn(t, srv, uri)
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
			`<h3 class="mb-2 text-2xl font-bold tracking-tight"> This website is not supported </h3>`,
			`<p class="mb-3 text-gray-700"> You can either request our team to support this website or go back to the previous page. </p>`,
			`<button name="website" value="https://www.example.com" class="w-full col-span-4 ml-2 p-2 font-semibold text-white bg-blue-500 hover:bg-blue-800"> Request </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("add recipe from supported website error", func(t *testing.T) {
		repo := &mockRepository{Recipes: make(map[int64]models.Recipes, 0)}
		repo.AddRecipeFunc = func(r *models.Recipe, userID int64) (int64, error) {
			return -1, errors.New("add recipe error")
		}
		srv.Repository = repo

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader("https://www.eatingbirdfood.com/cinnamon-rolls/"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Recipe could not be added.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("add recipe from a supported website", func(t *testing.T) {
		repo := &mockRepository{Recipes: make(map[int64]models.Recipes, 0)}
		called := 0
		repo.AddRecipeFunc = func(r *models.Recipe, userID int64) (int64, error) {
			called += 1
			return 1, nil
		}
		srv.Repository = repo

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader("https://www.eatingbirdfood.com/cinnamon-rolls/"))

		assertStatus(t, rr.Code, http.StatusSeeOther)
		if called != 1 {
			t.Fatal("recipe must have been added to the user's database")
		}
	})
}

func TestHandlers_recipes_Delete(t *testing.T) {
	repo := &mockRepository{
		Recipes:         map[int64]models.Recipes{1: make(models.Recipes, 0)},
		UsersRegistered: []models.User{{ID: 1, Email: "test@example.com"}},
	}
	srv := newServerTest()
	srv.Repository = repo

	uri := "/recipes"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
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

func TestHandlers_recipes_Share(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes/%d/share"

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
	_ = srv.Repository.AddShareLink("https://www.recipya.com/recipes/1/share", 1)

	t.Run("create valid share link", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, fmt.Sprintf(uri, 1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<input type="url" value="https://www.recipya.com/recipes/1/share" class="w-full rounded-lg bg-gray-100 px-4 py-2" readonly="readonly">`,
			`<button class="w-24 font-semibold p-2 bg-gray-300 rounded-lg hover:bg-gray-400" title="Copy to clipboard" _="on click js if ('clipboard' in window.navigator) { navigator.clipboard.writeText('https://www.recipya.com/recipes/1/share') } end then put 'Copied!' into me then add @title='Copied to clipboard!' then toggle @disabled on me then toggle .cursor-not-allowed .bg-green-600 .text-white .hover:bg-gray-400 on me"> Copy </button>`,
		}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	t.Run("create invalid share link", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, fmt.Sprintf(uri, 10), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Failed to create share link.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("access share link anonymous", func(t *testing.T) {
		rr := sendRequest(srv, http.MethodGet, fmt.Sprintf(uri, 1), noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		want := []string{
			`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
			`<button class="ml-2" title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release()`,
			`<div class="flex"><div class="flex-auto"><a href="/auth/login" class="mr-4 rounded-lg px-2 py-1 text-white hover:bg-green-600">Log In</a></div><div class="flex-auto"><a href="/auth/register" class="mr-4 rounded-lg bg-blue-100 px-2 py-1 hover:bg-red-600 hover:text-white"> Sign Up </a></div></div>`,
			`<h1 class="grid col-span-4 py-2 font-bold place-content-center">` + recipe.Name + `</h1><div class="grid justify-end col-span-1 place-content-center print:invisible"><button class="mr-2" onclick="print()" title="Print recipe"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"/></svg></button>`,
			`<img id="output" class="w-full text-center h-96" alt="Image of the recipe" style="object-fit: scale-down" src="/data/images/` + recipe.Image.String() + `.jpg">`,
			`<span class="text-sm font-normal leading-none">American</span>`,
			`<p class="text-sm text-center">2 servings</p>`,
			`<a class="p-1 duration-300 border rounded-lg hover:bg-gray-800 hover:text-white center" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank"> Source </a>`,
			`<p class="p-2 text-sm whitespace-pre-line">This is the most delicious recipe!</p>`,
			`<p class="text-xs">Per 100g: 500 calories; total carbohydrates 7g; sugars 6g; protein 3g; total fat 8g; saturated fat 4g; cholesterol 1g; sodium 5g; fiber 2g.</p>`,
			`<table class="min-w-full text-xs divide-y divide-gray-200 table-auto"><thead class="h-12 font-medium tracking-wider text-white bg-gray-800"><tr><th scope="col" class="text-right"><p class="text-center">Nutrition<br>(per 100g)</p></th><th scope="col" class="text-center"><p>Amount<br>(optional)</p></th></tr></thead><tbody class="text-right text-gray-500 bg-white divide-y divide-gray-200"><tr><td><p>Calories:</p></td><td class="py-3 text-center"> 500 </td></tr><tr class="bg-gray-100"><td><p>Total carbs:</p></td><td class="py-3 text-center"> 7g </td></tr><tr><td><p>Sugars:</p></td><td class="py-3 text-center"> 6g </td></tr><tr class="bg-gray-100"><td><p>Protein:</p></td><td class="py-3 text-center"> 3g </td></tr><tr><td><p>Total fat:</p></td><td class="py-3 text-center"> 8g </td></tr><tr class="bg-gray-100"><td><p>Saturated fat:</p></td><td class="py-3 text-center"> 4g </td></tr><tr><td><p>Cholesterol:</p></td><td class="py-3 text-center"> 1g </td></tr><tr class="bg-gray-100"><td><p>Sodium:</p></td><td class="py-3 text-center"> 5g </td></tr><tr><td><p>Fiber:</p></td><td class="py-3 text-center"> 2g </td></tr></tbody></table>`,
			`<td>Preparation:</td><td class="py-3 text-center print:py-0"><time datetime="PT05M">0h05</time></td>`,
			`<td>Cooking:</td><td class="w-1/2 py-3 text-center print:py-0"><time datetime="PT1H05M">1h05</time></td>`,
			`<td>Total:</td><td class="w-1/2 py-3 text-center print:py-0"><time datetime="PT1H10M">1h10</time></td>`,
			`<div class="col-span-6 py-2 border-b border-l border-r border-black md:col-span-2 md:border-r-0 print:hidden"><h2 class="pb-2 m-auto text-sm font-bold text-center text-gray-600 underline ">Ingredients</h2><ul class="grid px-6 columns-2"><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me then toggle @checked on first <input/> in me"><label for="ingredient-0"></label><input type="checkbox" id="ingredient-0" class="mt-1"><span class="pl-2">Ing1</span></li><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me then toggle @checked on first <input/> in me"><label for="ingredient-1"></label><input type="checkbox" id="ingredient-1" class="mt-1"><span class="pl-2">Ing2</span></li><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me then toggle @checked on first <input/> in me"><label for="ingredient-2"></label><input type="checkbox" id="ingredient-2" class="mt-1"><span class="pl-2">Ing3</span></li></ul></div>`,
			`<div class="col-span-6 py-2 pb-8 border-b border-l border-r border-black rounded-bl-lg md:rounded-bl-none md:col-span-4 print:hidden"><h2 class="pb-2 m-auto text-sm font-bold text-center text-gray-600 underline">Instructions</h2><ol class="grid px-6 list-decimal"><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me"><span class="whitespace-pre-line">Ins1</span></li><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me"><span class="whitespace-pre-line">Ins2</span></li><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me"><span class="whitespace-pre-line">Ins3</span></li></ol></div>`,
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
			rr := tc.sendFunc(srv, http.MethodGet, fmt.Sprintf(uri, 1), noHeader, nil)

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
			rr := tc.sendFunc(srv, http.MethodGet, fmt.Sprintf(uri, 1), noHeader, nil)

			assertStatus(t, rr.Code, http.StatusOK)
			want := []string{
				`<title hx-swap-oob="true">` + recipe.Name + " | Recipya</title>",
				`<button class="ml-2" title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release()`,
				`<button class="mr-2" onclick="print()" title="Print recipe">`,
				`<button class="mr-2" hx-delete="/recipes/1" hx-swap="none" title="Delete recipe" hx-confirm="Are you sure you wish to delete this recipe?">`,
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

	website1 := `<tr class="border px-8 py-2"><td class="border px-8 py-2">1</td><td class="border px-8 py-2"><a class="underline" href="https://101cookbooks.com" target="_blank">101cookbooks.com</a></td></tr>`
	website2 := `<tr class="border px-8 py-2"><td class="border px-8 py-2">2</td><td class="border px-8 py-2"><a class="underline" href="http://www.afghankitchenrecipes.com" target="_blank">afghankitchenrecipes.com</a></td></tr>`

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
	})

	t.Run("returns list of websites to logged in user", func(t *testing.T) {
		rr := sendRequestAsLoggedIn(srv, http.MethodGet, uri, noHeader, nil)

		assertStatus(t, rr.Code, http.StatusOK)
		assertHeader(t, rr, "Content-Type", "text/html")
		want := []string{website1, website2}
		assertStringsInHTML(t, getBodyHTML(rr), want)
	})

	testcases := []struct {
		name   string
		search string
		want   []string
	}{
		{
			name:   "search is empty",
			search: "",
			want:   []string{website1, website2},
		},
		{
			name:   "search for specific website",
			search: "101cookbooks",
			want:   []string{website1},
		},
		{
			name:   "search all dot coms",
			search: ".com",
			want:   []string{website1, website2},
		},
		{
			name:   "search query not present in any website",
			search: "z",
			want:   []string{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("search="+tc.search))

			assertStatus(t, rr.Code, http.StatusOK)
			assertHeader(t, rr, "Content-Type", "text/html")
			assertStringsInHTML(t, getBodyHTML(rr), tc.want)
		})
	}
}

func TestHandlers_Recipes_View(t *testing.T) {
	srv := newServerTest()

	uri := "/recipes"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, uri)
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
				`<button class="ml-2" title="Toggle screen lock" _="on load if not navigator.wakeLock hide me end on click if wakeLock wakeLock.release() then add @d='M 12.276 18.55 v -0.748 a 4.79 4.79 0 0 1 1.463 -3.458 a 5.763 5.763 0 0 0 1.804 -4.21 a 5.821 5.821 0 0 0 -6.475 -5.778 c -2.779 0.307 -4.99 2.65 -5.146 5.448 a 5.82 5.82 0 0 0 1.757 4.503 a 4.906 4.906 0 0 1 1.5 3.495 v 0.747 a 1.44 1.44 0 0 0 1.44 1.439 h 2.218 a 1.44 1.44 0 0 0 1.44 -1.439 z m -1.058 0 c 0 0.209 -0.17 0.38 -0.38 0.38 h -2.22 c -0.21 0 -0.38 -0.171 -0.38 -0.38 v -0.748 c 0 -1.58 -0.664 -3.13 -1.822 -4.254 A 4.762 4.762 0 0 1 4.98 9.863 c 0.127 -2.289 1.935 -4.204 4.205 -4.455 a 4.762 4.762 0 0 1 5.3 4.727 a 4.714 4.714 0 0 1 -1.474 3.443 a 5.853 5.853 0 0 0 -1.791 4.225 v 0.746 z M 11.45 20.51 H 8.006 a 0.397 0.397 0 1 0 0 0.795 h 3.444 a 0.397 0.397 0 1 0 0 -0.794 z M 11.847 22.162 a 0.397 0.397 0 0 0 -0.397 -0.397 H 8.006 a 0.397 0.397 0 1 0 0 0.794 h 3.444 c 0.22 0 0.397 -0.178 0.397 -0.397 z z z z z z z z M 10.986 23.416 H 8.867 a 0.397 0.397 0 1 0 0 0.794 h 1.722 c 0.22 0 0.397 -0.178 0.397 -0.397 z' to first <path/> else call initWakeLock() then add @d='M12.276 18.55v-.748a4.79 4.79 0 0 1 1.463-3.458 5.763 5.763 0 0 0 1.804-4.21 5.821 5.821 0 0 0-6.475-5.778c-2.779.307-4.99 2.65-5.146 5.448a5.82 5.82 0 0 0 1.757 4.503 4.906 4.906 0 0 1 1.5 3.495v.747a1.44 1.44 0 0 0 1.44 1.439h2.218a1.44 1.44 0 0 0 1.44-1.439zm-1.058 0c0 .209-.17.38-.38.38h-2.22c-.21 0-.38-.171-.38-.38v-.748c0-1.58-.664-3.13-1.822-4.254A4.762 4.762 0 0 1 4.98 9.863c.127-2.289 1.935-4.204 4.205-4.455a4.762 4.762 0 0 1 5.3 4.727 4.714 4.714 0 0 1-1.474 3.443 5.853 5.853 0 0 0-1.791 4.225v.746zM11.45 20.51H8.006a.397.397 0 1 0 0 .795h3.444a.397.397 0 1 0 0-.794zM11.847 22.162a.397.397 0 0 0-.397-.397H8.006a.397.397 0 1 0 0 .794h3.444c.22 0 .397-.178.397-.397zM.397 10.125h2.287a.397.397 0 1 0 0-.794H.397a.397.397 0 1 0 0 .794zM19.456 9.728a.397.397 0 0 0-.397-.397h-2.287a.397.397 0 1 0 0 .794h2.287c.22 0 .397-.178.397-.397zM9.331.397v2.287a.397.397 0 1 0 .794 0V.397a.397.397 0 1 0-.794 0zM16.045 2.85 14.43 4.465a.397.397 0 1 0 .561.561l1.617-1.617a.397.397 0 1 0-.562-.56zM5.027 14.429a.397.397 0 0 0-.56 0l-1.618 1.616a.397.397 0 1 0 .562.562l1.617-1.617a.397.397 0 0 0 0-.561zM4.466 5.027a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.56L3.41 2.848a.397.397 0 1 0-.561.561zM16.045 16.607a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.562L14.99 14.43a.397.397 0 1 0-.561.56zM10.986 23.416a.397.397 0 0 0-.397-.397H8.867a.397.397 0 1 0 0 .794h1.722c.22 0 .397-.178.397-.397z' to first <path/> end"><svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 hover:text-red-600" fill="none" viewBox="0 0 24 24" width="24px" height="24px" stroke="currentColor"><path d="M12.276 18.55v-.748a4.79 4.79 0 0 1 1.463-3.458 5.763 5.763 0 0 0 1.804-4.21 5.821 5.821 0 0 0-6.475-5.778c-2.779.307-4.99 2.65-5.146 5.448a5.82 5.82 0 0 0 1.757 4.503 4.906 4.906 0 0 1 1.5 3.495v.747a1.44 1.44 0 0 0 1.44 1.439h2.218a1.44 1.44 0 0 0 1.44-1.439zm-1.058 0c0 .209-.17.38-.38.38h-2.22c-.21 0-.38-.171-.38-.38v-.748c0-1.58-.664-3.13-1.822-4.254A4.762 4.762 0 0 1 4.98 9.863c.127-2.289 1.935-4.204 4.205-4.455a4.762 4.762 0 0 1 5.3 4.727 4.714 4.714 0 0 1-1.474 3.443 5.853 5.853 0 0 0-1.791 4.225v.746zM11.45 20.51H8.006a.397.397 0 1 0 0 .795h3.444a.397.397 0 1 0 0-.794zM11.847 22.162a.397.397 0 0 0-.397-.397H8.006a.397.397 0 1 0 0 .794h3.444c.22 0 .397-.178.397-.397zM.397 10.125h2.287a.397.397 0 1 0 0-.794H.397a.397.397 0 1 0 0 .794zM19.456 9.728a.397.397 0 0 0-.397-.397h-2.287a.397.397 0 1 0 0 .794h2.287c.22 0 .397-.178.397-.397zM9.331.397v2.287a.397.397 0 1 0 .794 0V.397a.397.397 0 1 0-.794 0zM16.045 2.85 14.43 4.465a.397.397 0 1 0 .561.561l1.617-1.617a.397.397 0 1 0-.562-.56zM5.027 14.429a.397.397 0 0 0-.56 0l-1.618 1.616a.397.397 0 1 0 .562.562l1.617-1.617a.397.397 0 0 0 0-.561zM4.466 5.027a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.56L3.41 2.848a.397.397 0 1 0-.561.561zM16.045 16.607a.396.396 0 0 0 .562 0 .397.397 0 0 0 0-.562L14.99 14.43a.397.397 0 1 0-.561.56zM10.986 23.416a.397.397 0 0 0-.397-.397H8.867a.397.397 0 1 0 0 .794h1.722c.22 0 .397-.178.397-.397z"/></svg></button>`,
				`<button hx-get="/recipes/1/edit" class="ml-2" title="Edit recipe">`,
				`<h1 class="grid col-span-4 py-2 font-bold place-content-center">Chicken Jersey</h1>`,
				`<button class="mr-2" title="Share recipe" hx-post="/recipes/1/share" hx-target="#share-dialog-result" _="on htmx:afterRequest from me if event.detail.successful call #share-dialog.showModal()">`,
				`<button class="mr-2" onclick="print()" title="Print recipe">`,
				`<button class="mr-2" hx-delete="/recipes/1" hx-swap="none" title="Delete recipe" hx-confirm="Are you sure you wish to delete this recipe?">`,
				`<dialog id="share-dialog" class="p-4 border-4 border-black min-w-[15rem]"><div id="share-dialog-result" class="pb-4"></div>`,
				`<img id="output" class="w-full text-center h-96" alt="Image of the recipe" style="object-fit: scale-down" src="/data/images/` + r.Image.String() + `.jpg">`,
				`<span class="text-sm font-normal leading-none">American</span>`,
				`<p class="text-sm text-center">2 servings</p>`,
				`<a class="p-1 duration-300 border rounded-lg hover:bg-gray-800 hover:text-white center" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank"> Source </a>`,
				`<p class="p-2 text-sm whitespace-pre-line">This is the most delicious recipe!</p>`,
				`<p class="text-xs">Per 100g: 500 calories; total carbohydrates 7g; sugars 6g; protein 3g; total fat 8g; saturated fat 4g; cholesterol 1g; sodium 5g; fiber 2g.</p>`,
				`<table class="min-w-full text-xs divide-y divide-gray-200 table-auto"><thead class="h-12 font-medium tracking-wider text-white bg-gray-800"><tr><th scope="col" class="text-right"><p class="text-center">Nutrition<br>(per 100g)</p></th><th scope="col" class="text-center"><p>Amount<br>(optional)</p></th></tr></thead><tbody class="text-right text-gray-500 bg-white divide-y divide-gray-200"><tr><td><p>Calories:</p></td><td class="py-3 text-center"> 500 </td></tr><tr class="bg-gray-100"><td><p>Total carbs:</p></td><td class="py-3 text-center"> 7g </td></tr><tr><td><p>Sugars:</p></td><td class="py-3 text-center"> 6g </td></tr><tr class="bg-gray-100"><td><p>Protein:</p></td><td class="py-3 text-center"> 3g </td></tr><tr><td><p>Total fat:</p></td><td class="py-3 text-center"> 8g </td></tr><tr class="bg-gray-100"><td><p>Saturated fat:</p></td><td class="py-3 text-center"> 4g </td></tr><tr><td><p>Cholesterol:</p></td><td class="py-3 text-center"> 1g </td></tr><tr class="bg-gray-100"><td><p>Sodium:</p></td><td class="py-3 text-center"> 5g </td></tr><tr><td><p>Fiber:</p></td><td class="py-3 text-center"> 2g </td></tr></tbody></table>`,
				`<td>Preparation:</td><td class="py-3 text-center print:py-0"><time datetime="PT05M">0h05</time></td>`,
				`<td>Cooking:</td><td class="w-1/2 py-3 text-center print:py-0"><time datetime="PT1H05M">1h05</time></td>`,
				`<td>Total:</td><td class="w-1/2 py-3 text-center print:py-0"><time datetime="PT1H10M">1h10</time></td>`,
				`<div class="col-span-6 py-2 border-b border-l border-r border-black md:col-span-2 md:border-r-0 print:hidden"><h2 class="pb-2 m-auto text-sm font-bold text-center text-gray-600 underline ">Ingredients</h2><ul class="grid px-6 columns-2"><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me then toggle @checked on first <input/> in me"><label for="ingredient-0"></label><input type="checkbox" id="ingredient-0" class="mt-1"><span class="pl-2">Ing1</span></li><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me then toggle @checked on first <input/> in me"><label for="ingredient-1"></label><input type="checkbox" id="ingredient-1" class="mt-1"><span class="pl-2">Ing2</span></li><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me then toggle @checked on first <input/> in me"><label for="ingredient-2"></label><input type="checkbox" id="ingredient-2" class="mt-1"><span class="pl-2">Ing3</span></li></ul></div>`,
				`<div class="col-span-6 py-2 pb-8 border-b border-l border-r border-black rounded-bl-lg md:rounded-bl-none md:col-span-4 print:hidden"><h2 class="pb-2 m-auto text-sm font-bold text-center text-gray-600 underline">Instructions</h2><ol class="grid px-6 list-decimal"><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me"><span class="whitespace-pre-line">Ins1</span></li><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me"><span class="whitespace-pre-line">Ins2</span></li><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me"><span class="whitespace-pre-line">Ins3</span></li></ol></div>`,
				`<script> var wakeLock = null; initWakeLock(); function initWakeLock() { navigator.wakeLock?.request("screen") .then((lock) => { wakeLock = lock; wakeLock.onrelease = () => { wakeLock = null; console.info("Screen lock deactivated."); } console.info("Screen lock activated."); }) .catch((err) => { ` + "console.log(`Screen lock error: ${err.name}, ${err.message}`)" + ` }); } </script>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}
