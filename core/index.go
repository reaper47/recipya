package core

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

const pattern = "+_+"

func Index() {
	env := InitEnv(repository.Db())

	recipes, err := getRecipes()
	if err != nil {
		log.Println(err)
	}

	var numInserted, numIndexed, numSkipped int64
	for _, r := range recipes {
		recipe, err := env.recipes.GetRecipe(r.Name)
		if err != nil {
			log.Printf("Error getting recipe: '%v'. Err: %v\n", r.Name, err)
			numSkipped++
			continue
		}

		if recipe != nil {
			if r.IsCreatedSameTime(recipe) {
				if r.IsModified(recipe) {
					if err = env.recipes.UpdateRecipe(r, recipe.ID); err != nil {
						log.Printf("Error while indexing recipe: '%v'. Err: %v\n", r.Name, err)
						numSkipped++
						continue
					}
					numIndexed++
					continue
				}
				numSkipped++
				continue
			}
			r.Name = pattern + r.Name
		}

		if err = env.recipes.InsertRecipe(r); err != nil {
			log.Printf("Error while inserting recipe: '%v'. Err: %v\n", r.Name, err)
			numSkipped++
			continue
		}
		numInserted++
	}

	log.Printf(
		"Number of recipes: %v. Inserted %v, indexed %v and skipped %v",
		len(recipes),
		numInserted,
		numIndexed,
		numSkipped,
	)
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
