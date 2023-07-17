package server_test

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
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

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, formHeader, strings.NewReader("title=Salsa&description=The best&calories=666&total-carbohydrates=31g&sugars=0.1mg&protein=5g&total-fat=0g&saturated-fat=0g&cholesterol=256mg&sodium=777mg&fiber=2g&time-preparation=00%3A15%3A30&time-cooking=00%3A30%3A15&ingredient-1=ing1&ingredient-2=ing2&instruction-1=ins1&instruction-2=ins2"))

		assertStatus(t, rr.Code, http.StatusCreated)
		id := int64(len(repo.Recipes))
		if len(repo.Recipes) != originalNumRecipes+1 {
			t.Fatal("expected one more recipe to be added to the database")
		}
		want := models.Recipe{
			Category:     "",
			CreatedAt:    time.Time{},
			Cuisine:      "",
			Description:  "The best",
			Image:        uuid.UUID{},
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
			URL:       "",
			Yield:     0,
		}
		gotRecipe := repo.Recipes[1][id-1]
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

		assertStatus(t, rr.Code, http.StatusNoContent)
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

		assertStatus(t, rr.Code, http.StatusNoContent)
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
				`<a href="/recipes/1/edit" class="ml-2" title="Edit recipe">`,
				`<h1 class="grid col-span-4 py-2 font-bold place-content-center">Chicken Jersey</h1>`,
				`<button class="mr-2" onclick="print()" title="Print recipe">`,
				`<button class="mr-2" hx-delete="/recipes/1" hx-swap="none" title="Delete recipe" hx-confirm="Are you sure you wish to delete this recipe?">`,
				`<img id="output" class="w-full text-center h-96" alt="Image of the recipe" style="object-fit: scale-down" src="/data/images/` + r.Image.String() + `.jpg">`,
				`<span class="text-sm font-normal leading-none">American</span>`,
				`<p class="text-sm text-center">2 servings</p>`,
				`<a class="p-1 text-sm duration-300 border rounded-lg hover:bg-gray-800 hover:text-white center" href="https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/" target="_blank"> Source </a>`,
				`<p class="p-2 text-sm whitespace-pre-line">This is the most delicious recipe!</p>`,
				`<p class="text-xs">Per 100g: 500 calories; total carbohydrates 7g; sugars 6g; protein 3g; total fat 8g; saturated fat 4g; cholesterol 1g; sodium 5g; fiber 2g.</p>`,
				`<table class="min-w-full text-xs divide-y divide-gray-200 table-auto"><thead class="h-12 font-medium tracking-wider text-white bg-gray-800"><tr><th scope="col" class="text-right"><p class="text-center">Nutrition<br>(per 100g)</p></th><th scope="col" class="text-center"><p>Amount<br>(optional)</p></th></tr></thead><tbody class="text-right text-gray-500 bg-white divide-y divide-gray-200"><tr><td><p>Calories:</p></td><td class="py-3 text-center"> 500 </td></tr><tr class="bg-gray-100"><td><p>Total carbs:</p></td><td class="py-3 text-center"> 7g </td></tr><tr><td><p>Sugars:</p></td><td class="py-3 text-center"> 6g </td></tr><tr class="bg-gray-100"><td><p>Protein:</p></td><td class="py-3 text-center"> 3g </td></tr><tr><td><p>Total fat:</p></td><td class="py-3 text-center"> 8g </td></tr><tr class="bg-gray-100"><td><p>Saturated fat:</p></td><td class="py-3 text-center"> 4g </td></tr><tr><td><p>Cholesterol:</p></td><td class="py-3 text-center"> 1g </td></tr><tr class="bg-gray-100"><td><p>Sodium:</p></td><td class="py-3 text-center"> 5g </td></tr><tr><td><p>Fiber:</p></td><td class="py-3 text-center"> 2g </td></tr></tbody></table>`,
				`<td>Preparation:</td><td class="py-3 text-center print:py-0"><time datetime="PT05M">0h05</time></td>`,
				`<td>Cooking:</td><td class="w-1/2 py-3 text-center print:py-0"><time datetime="PT1H05M">1h05</time></td>`,
				`<td>Total:</td><td class="w-1/2 py-3 text-center print:py-0"><time datetime="PT1H10M">1h10</time></td>`,
				`<div class="col-span-6 py-2 border-b border-l border-r border-black md:col-span-2 md:border-r-0 print:hidden"><h2 class="pb-2 m-auto text-sm font-bold text-center text-gray-600 underline ">Ingredients</h2><ul class="grid px-6 columns-2"><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me then toggle @checked on first <input/> in me"><label for="ingredient-0"></label><input type="checkbox" id="ingredient-0" class="mt-1"><span class="pl-2">Ing1</span></li><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me then toggle @checked on first <input/> in me"><label for="ingredient-1"></label><input type="checkbox" id="ingredient-1" class="mt-1"><span class="pl-2">Ing2</span></li><li class="min-w-full py-2 pl-4 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me then toggle @checked on first <input/> in me"><label for="ingredient-2"></label><input type="checkbox" id="ingredient-2" class="mt-1"><span class="pl-2">Ing3</span></li></ul></div>`,
				`<div class="col-span-6 py-2 pb-8 border-b border-l border-r border-black rounded-bl-lg md:rounded-bl-none md:col-span-4 print:hidden"><h2 class="pb-2 m-auto text-sm font-bold text-center text-gray-600 underline">Instructions</h2><ol class="grid px-6 list-decimal"><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me"><span class="whitespace-pre-line">Ins1</span></li><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me"><span class="whitespace-pre-line">Ins2</span></li><li class="min-w-full py-2 text-sm select-none hover:bg-gray-200" _="on mousedown toggle .line-through on me"><span class="whitespace-pre-line">Ins3</span></li></ol></div>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}
