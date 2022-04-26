package repository

import "github.com/reaper47/recipya/internal/models"

// Repository is the interface for the database repository.
type Repository interface {
	// CreateUsers stores a new user in the database.
	CreateUser(username, email, hashedPassword string) (models.User, error)

	// User gets a user from the database based on the username or email.
	User(id string) models.User

	// UserByID gets a user by its ID.
	UserByID(id int64) models.User

	// Recipe fetches a recipe from the database.
	//
	// The returned recipe will be empty if the query returns no row.
	Recipe(id int64) models.Recipe

	// RecipeForUser fetches a recipe for the given user. An error will
	// be returned if the recipe is not tied to the user.
	RecipeForUser(id, userID int64) (models.Recipe, error)

	// RecipeBelongsToUser verifies that the recipe belongs to the user.
	RecipeBelongsToUser(id, userID int64) bool

	// Recipes fetches all of the recipes from the database.
	Recipes(userID int64, page int) ([]models.Recipe, error)

	// RecipesCount returns the number of recipes in the database.
	RecipesCount() (int, error)

	// Categories gets all categories for the given user from the database.
	Categories(userID int64) []string

	// EditCategory changes the name of category 'oldCategory'
	// to 'newCategory' where applicable.
	EditCategory(oldCategory, newCategory string, userID int64) error

	// InsertCategory inserts a category for a user in the database.
	InsertCategory(name string, userID int64) error

	// DeleteCategory deletes the category for the user.
	DeleteCategory(name string, userID int64) error

	// InsertNewRecipe inserts a new recipe into the database.
	//
	// The CreatedAt and UpdatedAt timestamps are not inserted
	// because the database will take care it.
	InsertNewRecipe(r models.Recipe, userID int64) (int64, error)

	// UpdateRecipes update the recipe in the database.
	UpdateRecipe(r models.Recipe) error

	// DeleteRecipe deletes the recipe with the passed id from the database.
	DeleteRecipe(id int64) error

	// Images fetches all distinct image UUIDs for recipes.
	// An empty slice will be returned when an error occurred.
	Images() []string

	// Fetches all of the scraper's supported websites.
	Websites() []models.Website

	// Close closes the database's connection.
	Close()
}
