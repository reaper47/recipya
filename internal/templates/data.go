package templates

import (
	"github.com/reaper47/recipya/internal/models"
)

// Data holds data to pass on to the templates.
type Data struct {
	IsAuthenticated bool // IsAuthenticated says whether the user is authenticated.

	Title string // Title is the text inserted <title> tag's text.

	Content      string // Content is text to insert into the template.
	ContentTitle string // ContentTitle is the header of the Content.

	Scraper ScraperData
	View    *ViewRecipeData
}

// RegisterData is the data to pass on to the user registration template.
type RegisterData struct {
	Email           string
	PasswordConfirm string
}

// ScraperData holds template data related to the recipe scraper.
type ScraperData struct {
	UnsupportedWebsite string
}

// NewViewRecipeData creates and populates a new ViewRecipeData.
func NewViewRecipeData(id int64, recipe *models.Recipe, isFromHost, isShared bool) *ViewRecipeData {
	return &ViewRecipeData{
		FormattedTimes: formattedTimes{
			Cook:          formatDuration(recipe.Times.Cook, false),
			CookDateTime:  formatDuration(recipe.Times.Cook, true),
			Prep:          formatDuration(recipe.Times.Prep, false),
			PrepDateTime:  formatDuration(recipe.Times.Prep, true),
			Total:         formatDuration(recipe.Times.Total, false),
			TotalDateTime: formatDuration(recipe.Times.Total, true),
		},
		ID:          id,
		IsURL:       isURL(recipe.URL),
		IsUUIDValid: isUUIDValid(recipe.Image),
		Recipe:      recipe,
		Share: shareData{
			IsFromHost: isFromHost,
			IsShared:   isShared,
		},
	}
}

// ViewRecipeData holds template data related to viewing a recipe.
type ViewRecipeData struct {
	FormattedTimes formattedTimes
	ID             int64
	IsURL          bool
	IsUUIDValid    bool
	Recipe         *models.Recipe
	Share          shareData
}

type formattedTimes struct {
	Cook          string
	CookDateTime  string
	Prep          string
	PrepDateTime  string
	Total         string
	TotalDateTime string
}

type shareData struct {
	IsFromHost bool
	IsShared   bool
}
