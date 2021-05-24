package api

import (
	"fmt"
	"log"

	"github.com/reaper47/recipe-hunter/model"
	"github.com/reaper47/recipe-hunter/repository"
)

func Search(ingredients []string, mode int, limit int) {
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
		log.Println(err)
	}
	fmt.Println(recipes)
}
