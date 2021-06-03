package core

import (
	"log"

	"github.com/reaper47/recipe-hunter/model"
	"github.com/reaper47/recipe-hunter/repository"
)

// Search searches for recipes based on the ingredients.
func Search(ingredients []string, mode int, limit int) ([]*model.Recipe, error) {
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
		log.Printf("Error fetching the recipes: %v\n", err)
		return nil, err
	}

	return recipes, nil
}
