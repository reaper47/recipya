package core

import (
	"database/sql"

	"github.com/reaper47/recipe-hunter/model"
	"github.com/reaper47/recipe-hunter/repository"
)

// Env stores environment variables for use throughout the program.
type Env struct {
	recipes interface {
		InsertRecipe(r *model.Recipe) error
		UpdateRecipe(r *model.Recipe, id int64) error
		GetRecipe(name string) (*model.Recipe, error)
		SearchMaximizeFridge(ingredients []string, n int) ([]*model.Recipe, error)
		SearchMinimizeMissing(ingredients []string, n int) ([]*model.Recipe, error)
	}
}

// InitEnv initializes the Environment struct
func InitEnv(db *sql.DB) *Env {
	return &Env{
		recipes: &repository.Repository{DB: db},
	}
}
