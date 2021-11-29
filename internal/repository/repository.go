package repository

import "github.com/reaper47/recipya/internal/models"

// Repository is the database repository
type Repository interface {
	GetRecipe(id int64) (models.Recipe, error)
	GetAllRecipes() ([]models.Recipe, error)
	InsertNewRecipe(recipe models.Recipe) (int64, error)
	DeleteRecipe(id int64) error

	Close()
}
