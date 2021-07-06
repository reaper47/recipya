package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/reaper47/recipe-hunter/repository"
)

const baseUrl = "/api/v1"

var env = Env{recipes: &repository.MockRecipeModel{}}

func TestRoutes(t *testing.T) {
	t.Run("Get categories returns all categories", test_GetCategories)
	t.Run("Get all recipes", test_GetRecipes_All)
	t.Run("Get all recipes of a category", test_GetRecipes_Category)
}

func test_GetCategories(t *testing.T) {
	rr := sendRequest("GET", "/categories", env.getCategories, t)

	expected := `{"categories":["appetizer","side dish","dessert"]}`
	actual := strings.TrimSuffix(rr.Body.String(), "\n")
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func test_GetRecipes_All(t *testing.T) {
	rr := sendRequest("GET", "/recipes", env.getRecipes, t)

	expected := `{"recipes":[{"id":1,"name":"carrots","description":"some delicious carrots","url":"https://www.example.com","image":"","prepTime":"PT3H30M","cookTime":"PT0H30M","totalTime":"PT4H0M","recipeCategory":"side dish","keywords":"carrots,butter","recipeYield":4,"tool":null,"recipeIngredient":["1 avocado","2 carrots"],"recipeInstructions":["cut","cook","eat"],"nutrition":null,"dateModified":"20210820","dateCreated":"20210820"}]}`
	actual := strings.TrimSuffix(rr.Body.String(), "\n")
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func test_GetRecipes_Category(t *testing.T) {
	rr := sendRequest("GET", "/recipes?c=appetizer", env.getRecipes, t)

	expected := `{"recipes":[{"id":2,"name":"super carrots","description":"some super delicious carrots","url":"https://www.example.com","image":"","prepTime":"PT3H0M","cookTime":"PT0H30M","totalTime":"PT3H30M","recipeCategory":"appetizer","keywords":"super carrots,butter","recipeYield":8,"tool":null,"recipeIngredient":["2 avocado","10 super carrots"],"recipeInstructions":["cut","cook well","eat"],"nutrition":null,"dateModified":"20210822","dateCreated":"20210821"}]}`
	actual := strings.TrimSuffix(rr.Body.String(), "\n")
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func sendRequest(
	method string,
	endpoint string,
	h http.HandlerFunc,
	t *testing.T,
) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, baseUrl+endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	return rr
}
