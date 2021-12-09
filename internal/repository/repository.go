package repository

import "github.com/reaper47/recipya/internal/models"

// Repository is the database repository
type Repository interface {
	CreateUser(username, email, hashedPassword string) (models.User, error)
	GetUser(id string) models.User

	GetRecipe(id int64) models.Recipe
	GetRecipes(userID int64, page int) ([]models.Recipe, error)
	GetRecipesCount() (int, error)
	GetCategories() []string
	InsertNewRecipe(r models.Recipe, userID int64) (int64, error)
	UpdateRecipe(r models.Recipe) error
	DeleteRecipe(id int64) error

	Close()
}
