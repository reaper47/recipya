package core

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/otiai10/gosseract"
	"github.com/reaper47/recipya/api"
	"github.com/reaper47/recipya/data"
	"github.com/reaper47/recipya/model"
)

func initRecipesRoutes(r *mux.Router, env *Env) {
	GET := []string{http.MethodGet, http.MethodOptions}
	POST := []string{http.MethodPost, http.MethodOptions}
	jsonHeader := []string{"Content-Type", "application/json"}

	// GET /api/v1/recipes
	r.HandleFunc(api.Recipes, env.getRecipes).Methods(GET...)

	// POST /api/v1/recipes
	r.HandleFunc(api.Recipes, env.postRecipe).Headers(jsonHeader...).Methods(POST...)

	// GET /api/v1/recipes/info
	r.HandleFunc(api.RecipesInfo, env.getRecipesInfo).Methods(GET...)

	// GET /api/v1/categories
	r.HandleFunc(api.RecipeCategories, env.getCategories).Methods(GET...)

	// GET /api/v1/search
	r.HandleFunc(api.RecipeSearch, env.getSearch).Methods(GET...)

	// POST /api/v1/import/ocr
	r.HandleFunc(api.RecipeOcr, env.postImportOcr).Methods(POST...)

	// POST /api/v1/import/url
	r.HandleFunc(api.RecipeImport, env.postImportRecipe).Headers(jsonHeader...).Methods(POST...)

	// GET /api/v1/import/websites
	r.HandleFunc(api.ImportWebsites, env.getImportWebsites).Methods(GET...)
}

func (env *Env) getRecipes(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var category string
	if c, ok := q["c"]; ok {
		category = c[0]
	}

	page := -1
	if p, ok := q["page"]; ok {
		page, _ = strconv.Atoi(p[0])
		if page <= 0 {
			message := "Page must be >= 1."
			writeErrorJson(http.StatusBadRequest, message, w)
			return
		}
	}

	limit := -1
	if l, ok := q["limit"]; ok {
		limit, _ = strconv.Atoi(l[0])
		if limit <= 0 {
			message := "Limit must be >= 1."
			writeErrorJson(http.StatusBadRequest, message, w)
			return
		}
	}

	if limit > 0 && page == -1 {
		message := "The `limit` parameter must be used when `page` is specified."
		writeErrorJson(http.StatusBadRequest, message, w)
		return
	}

	if limit == -1 {
		limit = 12
	}

	recipes, err := env.recipes.GetRecipes(category, page, limit)
	if err != nil {
		message := "Get recipes: " + err.Error()
		writeErrorJson(http.StatusBadRequest, message, w)
		return
	}
	writeSuccessJson(&model.Recipes{Objects: recipes}, w)
}

func (env *Env) postRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe model.Recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	switch {
	case err == io.EOF:
		message := "No recipe in body of the request."
		writeErrorJson(http.StatusBadRequest, message, w)
		return
	case err != nil:
		message := "Error decoding the JSON body of the request: " + err.Error()
		writeErrorJson(http.StatusInternalServerError, message, w)
		return
	}

	ingredients := env.NlpExtractIngredients(recipe.RecipeIngredient)
	recipe.Nutrition = FetchNutrientsInfo(ingredients)

	recipeID, err := env.recipes.InsertRecipe(&recipe)
	if err != nil {
		message := "Adding the recipe has failed. Err: " + err.Error()
		writeErrorJson(http.StatusBadRequest, message, w)
		return
	}
	writeCreatedJson(&model.ID{Id: recipeID}, w)
}

func (env *Env) getRecipesInfo(w http.ResponseWriter, r *http.Request) {
	info, err := env.recipes.GetRecipesInfo()
	if err != nil {
		message := "Error retrieving recipes info: " + err.Error()
		writeErrorJson(http.StatusInternalServerError, message, w)
		return
	}
	writeSuccessJson(model.RecipesInfoWrapper{Info: info}, w)
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

func (env *Env) getImportWebsites(w http.ResponseWriter, r *http.Request) {
	writeSuccessJson(&model.Websites{Objects: data.Websites}, w)
}

func (env *Env) postImportOcr(w http.ResponseWriter, r *http.Request) {
	buf, err := api.UploadFileBuffer(r, 5, "image")
	if err != nil {
		message := "Error uploading the image: " + err.Error()
		writeErrorJson(http.StatusInternalServerError, message, w)
		return
	}

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(buf.Bytes())
	text, _ := client.Text()
	recipe, err := parseRecipeOcr(text)
	if err != nil {
		message := "Error parsing the recipe: " + err.Error()
		writeErrorJson(http.StatusInternalServerError, message, w)
		return
	}
	writeSuccessJson(&recipe, w)
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
