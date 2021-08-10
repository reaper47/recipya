package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/reaper47/recipya/model"
	"github.com/reaper47/recipya/repository"
)

const baseUrl = "/api/v1"

var env = Env{recipes: &repository.MockRecipeModel{}}

func TestRoutes(t *testing.T) {
	t.Run("Get categories returns all categories", test_GetCategories)
	t.Run("Get all recipes", test_GetRecipes_All)
	t.Run("Get all recipes of a category", test_GetRecipes_Category)
	t.Run("Import a recipe given a valid URL", test_PostImportRecipe_ValidUrl)
}

/*
GET
*/

func test_GetCategories(t *testing.T) {
	rr := sendRequest("GET", "/categories", nil, 200, env.getCategories, t)

	expected := `{"categories":["appetizer","side dish","dessert"]}`
	actual := strings.TrimSuffix(rr.Body.String(), "\n")
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func test_GetRecipes_All(t *testing.T) {
	rr := sendRequest("GET", "/recipes", nil, 200, env.getRecipes, t)

	expected := `{"recipes":[{"id":1,"name":"carrots","description":"some delicious carrots","url":"https://www.example.com","image":"","prepTime":"PT3H30M","cookTime":"PT0H30M","totalTime":"PT4H0M","recipeCategory":"side dish","keywords":"carrots,butter","recipeYield":4,"tool":null,"recipeIngredient":["1 avocado","2 carrots"],"recipeInstructions":["cut","cook","eat"],"nutrition":null,"dateModified":"20210820","dateCreated":"20210820"}]}`
	actual := strings.TrimSuffix(rr.Body.String(), "\n")
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func test_GetRecipes_Category(t *testing.T) {
	rr := sendRequest("GET", "/recipes?c=appetizer", nil, 200, env.getRecipes, t)

	expected := `{"recipes":[{"id":2,"name":"super carrots","description":"some super delicious carrots","url":"https://www.example.com","image":"","prepTime":"PT3H0M","cookTime":"PT0H30M","totalTime":"PT3H30M","recipeCategory":"appetizer","keywords":"super carrots,butter","recipeYield":8,"tool":null,"recipeIngredient":["2 avocado","10 super carrots"],"recipeInstructions":["cut","cook well","eat"],"nutrition":null,"dateModified":"20210822","dateCreated":"20210821"}]}`
	actual := strings.TrimSuffix(rr.Body.String(), "\n")
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

/*
POST
*/

func test_PostImportRecipe_ValidUrl(t *testing.T) {
	postBody, _ := json.Marshal(map[string]string{
		"url": "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/",
	})
	rr := sendRequest("POST", "/import/url", postBody, 201, env.postImportRecipe, t)

	var recipe model.Recipe
	err := json.NewDecoder(rr.Body).Decode(&recipe)
	if err != nil {
		t.Fatal(err)
	}
}

func sendRequest(
	method string,
	endpoint string,
	body []byte,
	code int,
	h http.HandlerFunc,
	t *testing.T,
) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, baseUrl+endpoint, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if status := rr.Code; status != code {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	return rr
}