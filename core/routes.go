package core

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/reaper47/recipya/api"
	"github.com/reaper47/recipya/model"
)

func initRecipesRoutes(r *mux.Router, env *Env) {
	r.HandleFunc(api.Recipes, env.getRecipes).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(api.RecipeCategories, env.getCategories).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(api.RecipeSearch, env.getSearch).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc(api.RecipeImport, env.postImportRecipe).Headers("Content-Type", "application/json").Methods(http.MethodPost, http.MethodOptions)
}

func (env *Env) getCategories(w http.ResponseWriter, r *http.Request) {
	c, err := env.recipes.GetCategories()
	if err != nil {
		message := "Get categories: " + err.Error()
		writeErrorJson(http.StatusBadRequest, message, w)
		return
	}
	writeSuccessJson(&model.Categories{Objects: c}, w)
}

func (env *Env) getRecipes(w http.ResponseWriter, r *http.Request) {
	var category string
	if c, ok := r.URL.Query()["c"]; ok {
		category = c[0]
	}

	recipes, err := env.recipes.GetRecipes(category)
	if err != nil {
		message := "Get recipes: " + err.Error()
		writeErrorJson(http.StatusBadRequest, message, w)
		return
	}
	writeSuccessJson(&model.Recipes{Objects: recipes}, w)
}

func (env *Env) getSearch(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	var err error

	ingredientsVar := vars["ingredients"]
	if len(ingredientsVar) == 0 {
		message := "Query parameter: 'ingredients' must be specified."
		writeErrorJson(http.StatusBadRequest, message, w)
		return
	}

	ingredients := strings.Split(ingredientsVar[0], ",")
	if len(ingredients) == 0 || len(ingredients) == 1 && strings.TrimSpace(ingredients[0]) == "" {
		message := "Query parameter: 'ingredients' must have one or more ingredients separated by a comma, e.g. avocado,garlic,chicken."
		writeErrorJson(http.StatusBadRequest, message, w)
		return
	}

	numVar := vars["num"]
	limit := 10
	if len(numVar) == 1 {
		limit, err = strconv.Atoi(numVar[0])
		if err != nil || limit < 1 || limit > 30 {
			message := "Query parameter: 'num' must be an integer between 1-30."
			writeErrorJson(http.StatusBadRequest, message, w)
			return
		}
	}

	modeVar := vars["mode"]
	mode := 2
	if len(modeVar) == 1 {
		mode, err = strconv.Atoi(modeVar[0])
		if err != nil || mode < 1 || mode > 2 {
			message := "Query parameter: 'mode' must be an integer either 1 or 2."
			writeErrorJson(http.StatusBadRequest, message, w)
			return
		}
	}

	recipes, err := Search(ingredients, mode, limit)
	if err != nil {
		message := "Error while searching: " + err.Error()
		writeErrorJson(http.StatusInternalServerError, message, w)
		return
	}
	if recipes == nil {
		recipes = make([]*model.Recipe, 0)
	}

	writeSuccessJson(&model.Recipes{Objects: recipes}, w)
}

func (env *Env) postImportRecipe(w http.ResponseWriter, r *http.Request) {
	var website model.Website
	err := json.NewDecoder(r.Body).Decode(&website)
	switch {
	case err == io.EOF:
		message := "No URL in body of the request."
		writeErrorJson(http.StatusBadRequest, message, w)
		return
	case err != nil:
		message := "Error decoding the JSON body of the request."
		writeErrorJson(http.StatusInternalServerError, message, w)
		return
	}

	recipe, err := env.recipes.ImportRecipe(website.Url)
	if err != nil {
		message := "Import recipe from website failed. The website might not be supported. Err: " + err.Error()
		writeErrorJson(http.StatusBadRequest, message, w)
		return
	}
	writeCreatedJson(&recipe, w)
}
