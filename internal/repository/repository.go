package repository

import "github.com/reaper47/recipya/internal/models"

// Repository is the database repository
type Repository interface {
	CreateUser(username, email, hashedPassword string) (models.User, error)
	User(id string) models.User

	Recipe(id int64) models.Recipe
	Recipes(userID int64, page int) ([]models.Recipe, error)
	RecipesCount() (int, error)
	Categories(userID int64) []string
	InsertCategory(name string, userID int64) error
	InsertNewRecipe(r models.Recipe, userID int64) (int64, error)
	UpdateRecipe(r models.Recipe) error
	DeleteRecipe(id int64) error

	Close()
}
