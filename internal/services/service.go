package services

import (
	"bytes"
	"github.com/blang/semver"
	"github.com/google/go-github/v59/github"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/internal/units"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// RepositoryService is the interface that describes the methods required for managing the main data store.
type RepositoryService interface {
	// AddAuthToken adds an authentication token to the database.
	AddAuthToken(selector, validator string, userID int64) error

	// AddCookbook adds a cookbook to the database.
	AddCookbook(title string, userID int64) (int64, error)

	// AddCookbookRecipe adds a recipe to the cookbook.
	AddCookbookRecipe(cookbookID, recipeID, userID int64) error

	// AddRecipeCategory adds a custom recipe category for the user.
	AddRecipeCategory(name string, userID int64) error

	// AddRecipes adds recipes to the user's collection.
	AddRecipes(recipes models.Recipes, userID int64, progress chan models.Progress) ([]int64, []models.ReportLog, error)

	// AddReport adds a report to the database.
	AddReport(report models.Report, userID int64)

	// AddShareLink adds a share link for the recipe.
	AddShareLink(share models.Share) (string, error)

	// AddShareRecipe adds a shared recipe to the user's collection.
	AddShareRecipe(recipeID, userID int64) (int64, error)

	// Categories gets all categories in the database.
	Categories(userID int64) ([]string, error)

	// CheckUpdate checks whether there is a new release for Recipya.
	// It returns the latest information on the application.
	CheckUpdate(files FilesService) (models.AppInfo, error)

	// Confirm confirms the user's account.
	Confirm(userID int64) error

	// Cookbook gets a cookbook by its ID.
	Cookbook(id, userID int64) (models.Cookbook, error)

	// CookbookRecipe gets a recipe from a cookbook.
	CookbookRecipe(id, cookbookID int64) (recipe *models.Recipe, userID int64, err error)

	// CookbookShared checks whether the cookbook is shared.
	// It returns a models.Share. Otherwise, an error.
	CookbookShared(id string) (*models.Share, error)

	// Cookbooks gets a limited number of cookbooks belonging to the user.
	Cookbooks(userID int64, page uint64) ([]models.Cookbook, error)

	// CookbooksShared gets the user's shared cookbooks.
	CookbooksShared(userID int64) ([]models.Share, error)

	// CookbooksUser gets all the user's cookbooks.
	CookbooksUser(userID int64) ([]models.Cookbook, error)

	// Counts gets the models.Counts for the user.
	Counts(userID int64) (models.Counts, error)

	// DeleteAuthToken removes an authentication token from the database.
	DeleteAuthToken(userID int64) error

	// DeleteCookbook deletes a user's cookbook.
	DeleteCookbook(id, userID int64) error

	// DeleteRecipe deletes a user's recipe.
	DeleteRecipe(id, userID int64) error

	// DeleteRecipeCategory deletes a user's recipe category.
	DeleteRecipeCategory(name string, userID int64) error

	// DeleteRecipeFromCookbook deletes a recipe from a cookbook. It returns the number of recipes in the cookbook.
	DeleteRecipeFromCookbook(recipeID, cookbookID int64, userID int64) (int64, error)

	// DeleteUser deletes a user and his or her data.
	DeleteUser(id int64) error

	// GetAuthToken gets a non-expired auth token by the selector.
	GetAuthToken(selector, validator string) (models.AuthToken, error)

	// Images fetches all distinct image UUIDs for recipes.
	// An empty slice is returned when an error occurred.
	Images() []string

	// InitAutologin creates a default user for the autologin feature if no users are present.
	InitAutologin() error

	// IsUserExist checks whether the user is present in the database.
	IsUserExist(email string) bool

	// IsUserPassword checks whether the password is the user's password.
	IsUserPassword(id int64, password string) bool

	// MeasurementSystems gets the units systems, along with the one the user selected, in the database.
	MeasurementSystems(userID int64) ([]units.System, models.UserSettings, error)

	// Nutrients gets the nutrients for the ingredients from the FDC database, along with the total weight.
	Nutrients(ingredients []string) (models.NutrientsFDC, float64, error)

	// Recipe gets the user's recipe of the given id.
	Recipe(id, userID int64) (*models.Recipe, error)

	// Recipes gets the user's recipes.
	Recipes(userID int64, opts models.SearchOptionsRecipes) models.Recipes

	// RecipesAll gets all the user's recipes.
	RecipesAll(userID int64) models.Recipes

	// RecipeShared checks whether the recipe is shared.
	// It returns a models.Share. Otherwise, an error.
	RecipeShared(id string) (*models.Share, error)

	// RecipesShared gets all the user's shared recipes.
	RecipesShared(userID int64) ([]models.Share, error)

	// RecipeUser gets the user for which the recipe belongs to.
	RecipeUser(recipeID int64) int64

	// Register adds a new user to the store.
	Register(email string, hashPassword auth.HashedPassword) (int64, error)

	// ReorderCookbookRecipes reorders the recipe indices of a cookbook.
	ReorderCookbookRecipes(cookbookID int64, recipeIDs []uint64, userID int64) error

	// Report gets a report of any type belonging to the user.
	Report(id, userID int64) ([]models.ReportLog, error)

	// ReportsImport gets all import reports.
	ReportsImport(userID int64) ([]models.Report, error)

	// RestoreUserBackup restores the user's data.
	RestoreUserBackup(backup *models.UserBackup) error

	// SearchRecipes searches for recipes based on the configuration.
	// It returns the paginated search recipes, the total number of search results and an error.
	SearchRecipes(opts models.SearchOptionsRecipes, userID int64) (models.Recipes, uint64, error)

	// SwitchMeasurementSystem sets the user's units system to the desired one.
	SwitchMeasurementSystem(system units.System, userID int64) error

	// UpdateCalculateNutrition updates the user's calculate nutrition facts automatically setting.
	UpdateCalculateNutrition(userID int64, isEnabled bool) error

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
	// Queue adds an unsent email to the queue.
	Queue(to string, template templates.EmailTemplate, data any)

	// RateLimits gets the SendGrid API's remaining and reset rate limits.
	RateLimits() (remaining int, resetUnix int64, err error)

	// Send sends an email using the SendGrid API.
	Send(to string, template templates.EmailTemplate, data any) error

	// SendQueue sends emails in the queue until the rate limit has been reached.
	SendQueue() (sent, remaining int, err error)
}

// FilesService is the interface that describes the methods required for manipulating files.
type FilesService interface {
	// BackupGlobal backs up the whole database to the backup directory.
	BackupGlobal() error

	// Backups gets the list of backup dates sorted in descending order for the given user.
	Backups(userID int64) []time.Time

	// BackupUserData backs up a specific user's data to the backup directory.
	BackupUserData(repo RepositoryService, userID int64) error

	// BackupUsersData backs up each user's data to the backup directory.
	BackupUsersData(repo RepositoryService) error

	// ExportCookbook exports the cookbook in the desired file type.
	// It returns the name of file in the temporary directory.
	ExportCookbook(cookbook models.Cookbook, fileType models.FileType) (string, error)

	// ExportRecipes creates a zip containing the recipes to export in the desired file type.
	ExportRecipes(recipes models.Recipes, fileType models.FileType, progress chan int) (*bytes.Buffer, error)

	// ExtractRecipes extracts the recipes from the HTTP files.
	ExtractRecipes(fileHeaders []*multipart.FileHeader) models.Recipes

	// ExtractUserBackup extracts data from the user backup for restoration.
	ExtractUserBackup(date string, userID int64) (*models.UserBackup, error)

	// IsAppLatest checks whether there is a software update.
	IsAppLatest(current semver.Version) (bool, *github.RepositoryRelease, error)

	// MergeImagesToPDF merges images to a PDF file.
	MergeImagesToPDF(images []io.Reader) io.ReadWriter

	// ReadTempFile gets the content of a file in the temporary directory.
	ReadTempFile(name string) ([]byte, error)

	// ScrapeAndStoreImage takes a URL as input and will download and store the image, and return a UUID referencing the image's internal ID
	ScrapeAndStoreImage(rawURL string) (uuid.UUID, error)

	// UpdateApp updates the application to the latest version.
	UpdateApp(current semver.Version) error

	// UploadImage uploads an image to the server.
	UploadImage(rc io.ReadCloser) (uuid.UUID, error)
}

// HttpService is the interface that describes the methods required for preparing and utilizing https requests and responses.
type HttpService interface {
	// PrepareRequestForURL Prepares an HTTP GET request for a given URL.
	// It will apply additional HTTP headers if the host requires it.
	PrepareRequestForURL(url string) (*http.Request, error)

	// GetHost gets the host from the raw URL.
	GetHost(rawURL string) string
}

// IntegrationsService is the interface that describes the methods required for various software integrations.
type IntegrationsService interface {
	// MealieImport imports the recipes from a Mealie instance.
	MealieImport(baseURL, username, password string, files FilesService, progress chan models.Progress) (models.Recipes, error)

	// NextcloudImport imports the recipes from a Nextcloud instance.
	NextcloudImport(baseURL, username, password string, files FilesService, progress chan models.Progress) (models.Recipes, error)

	// ProcessImageOCR processes an image using an OCR service to extract the recipe.
	ProcessImageOCR(files []io.Reader) (models.Recipes, error)

	// TandoorImport imports the recipes from a Tandoor Recipes instance.
	TandoorImport(baseURL, username, password string, files FilesService, progress chan models.Progress) (models.Recipes, error)

	// TestConnection tests the connection of an integration. No error is returned on success.
	TestConnection(api string) error
}
