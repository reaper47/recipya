package core

import (
	"database/sql"

	"github.com/reaper47/recipya/model"
	"github.com/reaper47/recipya/repository"
)

// Env stores environment variables for use throughout the program.
type Env struct {
	recipes interface {
		GetRecipe(name string) (*model.Recipe, error)
		GetRecipes(category string, page int, limit int) ([]*model.Recipe, error)
		GetRecipesInfo() (*model.RecipesInfo, error)
		GetCategories() ([]string, error)
		InsertRecipe(r *model.Recipe) (int64, error)
		ImportRecipe(url string) (*model.Recipe, error)
		SearchMaximizeFridge(ingredients []string, n int) ([]*model.Recipe, error)
		SearchMinimizeMissing(ingredients []string, n int) ([]*model.Recipe, error)
		UpdateRecipe(r *model.Recipe, id int64) error
	}
	data interface {
		GetBlacklistIngredients() (map[string]int8, error)
		GetFruitsVeggies() (map[string]int8, error)
	}
}

// InitEnv initializes the Environment struct.
func InitEnv(db *sql.DB) *Env {
	repo := &repository.Repository{DB: db}
	return &Env{
		recipes: repo,
		data:    repo,
	}
}
