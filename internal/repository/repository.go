package repository

import "github.com/reaper47/recipya/internal/models"

// Repository is the database repository
type Repository interface {
	GetRecipe(id int64) (models.Recipe, error)
	GetAllRecipes(page int) ([]models.Recipe, error)
	GetRecipesCount() (int, error)
	InsertNewRecipe(r models.Recipe) (int64, error)
	UpdateRecipe(r models.Recipe) error
	DeleteRecipe(id int64) error

	Close()
}
