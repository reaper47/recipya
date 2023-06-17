package services

import (
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
)

// RepositoryService is the interface that describes the methods required for managing the main data store.
type RepositoryService interface {
	// AddAuthToken adds an authentication token to the database.
	AddAuthToken(selector, validator string, userID int64) error

	// AddRecipe adds a recipe to the user's collection.
	AddRecipe(r *models.Recipe, userID int64) error

	// Confirm confirms the user's account.
	Confirm(userID int64) error

	// DeleteAuthToken removes an authentication token from the database.
	DeleteAuthToken(userID int64) error

	// GetAuthToken gets a non-expired auth token by the selector.
	GetAuthToken(selector, validator string) (models.AuthToken, error)

	// IsUserExist checks whether the user is present in the database.
	IsUserExist(email string) bool

	// Register adds a new user to the store.
	Register(email string, hashPassword auth.HashedPassword) (int64, error)

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
