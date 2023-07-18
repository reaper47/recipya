package services

import (
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"mime/multipart"
)

// RepositoryService is the interface that describes the methods required for managing the main data store.
type RepositoryService interface {
	// AddAuthToken adds an authentication token to the database.
	AddAuthToken(selector, validator string, userID int64) error

	// AddRecipe adds a recipe to the user's collection.
	AddRecipe(r *models.Recipe, userID int64) (int64, error)

	// AddShareLink adds a share link for the recipe.
	AddShareLink(link string, recipeID int64) error

	// Confirm confirms the user's account.
	Confirm(userID int64) error

	// DeleteAuthToken removes an authentication token from the database.
	DeleteAuthToken(userID int64) error

	// GetAuthToken gets a non-expired auth token by the selector.
	GetAuthToken(selector, validator string) (models.AuthToken, error)

	// IsRecipeShared checks whether the recipe is shared.
	IsRecipeShared(id int64) bool

	// IsUserExist checks whether the user is present in the database.
	IsUserExist(email string) bool

	// Recipe gets the user's recipe of the given id.
	Recipe(id, userID int64) (*models.Recipe, error)

	// RecipeUser gets the user for which the recipe belongs to.
	RecipeUser(recipeID int64) int64

	// Register adds a new user to the store.
	Register(email string, hashPassword auth.HashedPassword) (int64, error)

	// UpdatePassword updates the user's password.
	UpdatePassword(userID int64, hashedPassword auth.HashedPassword) error

	// UserID gets the user's id from the email. It returns -1 if user not found.
	UserID(email string) int64

	// UserInitials gets the user's initials of maximum two characters.
	UserInitials(userID int64) string

	// Users gets all users in the database.
	Users() []models.User

	// VerifyLogin checks whether the user provided correct login credentials.
	// If yes, their user ID will be returned. Otherwise, -1 is returned.
	VerifyLogin(email, password string) int64

	// Websites gets the list of supported websites from which to extract the recipe.
	Websites() models.Websites

	// WebsitesSearch gets the list of supported websites that match the query.
	WebsitesSearch(query string) models.Websites
}

// EmailService is the interface that describes the methods required for the email client.
type EmailService interface {
	// Send sends an email using the SendGrid API.
	Send(to string, template templates.EmailTemplate, data any)
}

// FilesService is the interface that describes the methods required for manipulating files.
type FilesService interface {
	// ExtractRecipes extracts the recipes from the HTTP files.
	ExtractRecipes(fileHeaders []*multipart.FileHeader) models.Recipes
}
