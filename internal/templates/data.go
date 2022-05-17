package templates

import (
	"github.com/reaper47/recipya/internal/models"
)

// Data holds general template data.
type Data struct {
	HideSidebar  bool
	HideGap      bool
	HeaderData   HeaderData
	IsViewRecipe bool

	RecipesData RecipesData
	RecipeData  RecipeData
	Categories  []string
	Scraper     Scraper

	FormErrorData FormErrorData
}

// HeaderData holds data for the header.
type HeaderData struct {
	Hide              bool
	IsUnauthenticated bool
	AvatarInitials    string
}

// FormErrorData holds errors related to forms.
type FormErrorData struct {
	Username, Email, Password string
}

// IsEmpty checks whether all of the form's error data fields are empty.
func (f FormErrorData) IsEmpty() bool {
	return f.Username == "" && f.Email == "" && f.Password == ""
}

// RecipesData holds data to pass on to the index template.
type RecipesData struct {
	Recipes    []models.Recipe
	Pagination Pagination
}

// RecipeData holds data to pass to the recipe templates.
type RecipeData struct {
	Recipe           models.Recipe
	HideEditControls bool
}

// Scraper holds template data related to the recipe scraper.
type Scraper struct {
	IsEmailSetUp bool
	Websites     []models.Website
}
