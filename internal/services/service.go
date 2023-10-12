package services

import (
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/internal/units"
	"io"
	"mime/multipart"
	"net/http"
)

// RepositoryService is the interface that describes the methods required for managing the main data store.
type RepositoryService interface {
	// AddAuthToken adds an authentication token to the database.
	AddAuthToken(selector, validator string, userID int64) error

	// AddCookbook adds a cookbook to the database.
	AddCookbook(title string, userID int64) (int64, error)

	// AddRecipe adds a recipe to the user's collection.
	AddRecipe(r *models.Recipe, userID int64) (int64, error)

	// AddShareLink adds a share link for the recipe.
	AddShareLink(share models.ShareRecipe) (string, error)

	// Categories gets all categories in the database.
	Categories(userID int64) ([]string, error)

	// Confirm confirms the user's account.
	Confirm(userID int64) error

	// Cookbooks gets a limited number of cookbooks belonging to the user.
	Cookbooks(userID int64) ([]models.Cookbook, error)

	// DeleteAuthToken removes an authentication token from the database.
	DeleteAuthToken(userID int64) error

	// DeleteCookbook deletes a user's cookbook.
	DeleteCookbook(id, userID int64) error

	// DeleteRecipe deletes a user's recipe. It returns the number of rows affected.
	DeleteRecipe(id, userID int64) (int64, error)

	// GetAuthToken gets a non-expired auth token by the selector.
	GetAuthToken(selector, validator string) (models.AuthToken, error)

	// IsUserExist checks whether the user is present in the database.
	IsUserExist(email string) bool

	// IsUserPassword checks whether the password is the user's password.
	IsUserPassword(id int64, password string) bool

	// MeasurementSystems gets the units systems, along with the one the user selected, in the database.
	MeasurementSystems(userID int64) ([]units.System, models.UserSettings, error)

	// Recipe gets the user's recipe of the given id.
	Recipe(id, userID int64) (*models.Recipe, error)

	// Recipes gets the user's recipes.
	Recipes(userID int64) models.Recipes

	// RecipeShared checks whether the recipe is shared.
	// It returns a models.ShareRecipe. Otherwise, an error.
	RecipeShared(id string) (*models.ShareRecipe, error)

	// RecipeUser gets the user for which the recipe belongs to.
	RecipeUser(recipeID int64) int64

	// Register adds a new user to the store.
	Register(email string, hashPassword auth.HashedPassword) (int64, error)

	// SwitchMeasurementSystem sets the user's units system to the desired one.
	SwitchMeasurementSystem(system units.System, userID int64) error

	// UpdateConvertMeasurementSystem updates the user's convert automatically setting.
	UpdateConvertMeasurementSystem(userID int64, isEnabled bool) error

	// UpdateCookbookImage updates the image of a user's cookbook.
	UpdateCookbookImage(id int64, image uuid.UUID, userID int64) error

	// UpdatePassword updates the user's password.
	UpdatePassword(userID int64, hashedPassword auth.HashedPassword) error

	// UpdateRecipe updates the recipe with its new values.
	UpdateRecipe(updatedRecipe *models.Recipe, userID int64, recipeNum int64) error

	// UpdateUserSettingsCookbooksViewMode updates the user's preferred cookbooks viewing mode.
	UpdateUserSettingsCookbooksViewMode(userID int64, mode models.ViewMode) error

	// UserID gets the user's id from the email. It returns -1 if user not found.
	UserID(email string) int64

	// UserSettings gets the user's settings.
	UserSettings(userID int64) (models.UserSettings, error)

	// UserInitials gets the user's initials of maximum two characters.
	UserInitials(userID int64) string

	// Users gets all users in the database.
	Users() []models.User

	// VerifyLogin checks whether the user provided correct login credentials.
	// If yes, their user ID will be returned. Otherwise, -1 is returned.
	VerifyLogin(email, password string) int64

	// Websites gets the list of supported websites from which to extract the recipe.
	Websites() models.Websites
}

// EmailService is the interface that describes the methods required for the email client.
type EmailService interface {
	// Send sends an email using the SendGrid API.
	Send(to string, template templates.EmailTemplate, data any)
}

// FilesService is the interface that describes the methods required for manipulating files.
type FilesService interface {
	// ExportRecipes creates a zip containing the recipes to export in the given file type.
	// It returns the name of file in the temporary directory.
	ExportRecipes(recipes models.Recipes, fileType models.FileType) (string, error)

	// ExtractRecipes extracts the recipes from the HTTP files.
	ExtractRecipes(fileHeaders []*multipart.FileHeader) models.Recipes

	// ReadTempFile gets the content of a file in the temporary directory.
	ReadTempFile(name string) ([]byte, error)

	// UploadImage uploads an image to the server.
	UploadImage(rc io.ReadCloser) (uuid.UUID, error)
}

// IntegrationsService is the interface that describes the methods required for various software integrations.
type IntegrationsService interface {
	NextcloudImport(client *http.Client, baseURL, username, password string, files FilesService) (*models.Recipes, error)
}
