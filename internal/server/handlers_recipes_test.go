package server_test

import (
	"errors"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
		want     []string
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
				`<button id="search-button" class="underline" hx-get="/recipes/supported-websites" hx-target="#search-results" > supported </button>`,
				`<button class="flex justify-center w-full text-gray-900 duration-300 border-2 border-gray-800 rounded-lg hover:bg-gray-800 hover:text-white center" hx-target="#content" hx-push-url="/recipes/add/unsupported-website" hx-prompt="Enter the recipe's URL" hx-post="/recipes/add/website" hx-indicator="#fullscreen-loader"> Fetch </button>`,
			}
			assertStringsInHTML(t, getBodyHTML(rr), want)
		})
	}
}

func TestHandlers_Recipes_AddRequestWebsite(t *testing.T) {
	emailMock := &mockEmail{}
	srv := server.NewServer(&mockRepository{}, emailMock)

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
		repo.AddRecipeFunc = func(r *models.Recipe, userID int64) error {
			return errors.New("add recipe error")
		}
		srv.Repository = repo

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uri, promptHeader, strings.NewReader("https://www.eatingbirdfood.com/cinnamon-rolls/"))

		assertStatus(t, rr.Code, http.StatusNoContent)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Recipe could not be added.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("add recipe from a supported website", func(t *testing.T) {
		repo := &mockRepository{Recipes: make(map[int64]models.Recipes, 0)}
		called := 0
		repo.AddRecipeFunc = func(r *models.Recipe, userID int64) error {
			called += 1
			return nil
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
