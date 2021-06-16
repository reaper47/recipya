package core

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/reaper47/recipe-hunter/api"
	"github.com/reaper47/recipe-hunter/model"
)

func initRecipesRoutes(r *mux.Router, env *Env) {
	r.HandleFunc(api.RecipeSearch, env.getSearch).Methods(http.MethodGet, http.MethodOptions)
}

func (env *Env) getSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&model.Recipes{
		Objects: recipes,
	})
}

func writeErrorJson(code int, message string, w http.ResponseWriter) {
	payload := api.ErrorJson{
		Objects: api.Error{
			Code:    code,
			Message: message,
			Status:  http.StatusText(code),
		}}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(payload)
}
