package api

import (
	"encoding/json"
	"log"

	"github.com/reaper47/recipe-hunter/model"
	"github.com/reaper47/recipe-hunter/repository"
)

func Search(ingredients []string, mode int, limit int) []byte {
	env := InitEnv(repository.Db())

	var (
		recipes []*model.Recipe
		err     error
	)

	if mode == 1 {
		recipes, err = env.recipes.SearchMinimizeMissing(ingredients, limit)
	} else {
		recipes, err = env.recipes.SearchMaximizeFridge(ingredients, limit)
	}

	if err != nil {
		log.Fatalf("Error fetching the recipes: %v\n", err)
	}

	recipesJson, err := json.Marshal(model.Recipes{Objects: recipes})
	if err != nil {
		log.Fatalf("Error marhsaling the recipes: %v\n", err)
	}
	return recipesJson
}
