package api

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

	"github.com/reaper47/recipe-hunter/config"
	"github.com/reaper47/recipe-hunter/model"
	"github.com/reaper47/recipe-hunter/repository"
)

func Index() {
	env := InitEnv(repository.Db())

	recipes, err := getRecipes()
	if err != nil {
		log.Println(err)
	}
	for _, r := range recipes{
		if err := env.recipes.InsertRecipe(r); err != nil {
			log.Println(err)
		}
	}

	log.Printf("Indexed %v recipes", len(recipes)) 
}

func getRecipes() ([]*model.Recipe, error) {
	var recipes []*model.Recipe
	re := regexp.MustCompile(`"nutrition":\[\]`)

	err := filepath.WalkDir(
		config.Config.RecipesDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(path) == ".json" {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				data = re.ReplaceAll(data, []byte(`"nutrition":{}`))

				var recipe *model.Recipe
				if err = json.Unmarshal(data, &recipe); err != nil {
					return err
				}
				recipes = append(recipes, recipe)
			}
			return nil
		},
	)
	return recipes, err
}
