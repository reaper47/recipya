package templates

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/units"
	"strings"
)

// Data holds data to pass on to the templates.
type Data struct {
	IsAuthenticated bool // IsAuthenticated says whether the user is authenticated.

	Title string // Title is the text inserted <title> tag's text.

	Content      string // Content is text to insert into the template.
	ContentTitle string // ContentTitle is the header of the Content.

	Functions FunctionsData

	Recipes  models.Recipes
	Settings SettingsData
	Scraper  ScraperData
	View     *ViewRecipeData
}

// RegisterData is the data to pass on to the user registration template.
type RegisterData struct {
	Email           string
	PasswordConfirm string
}

// SettingsData holds template data related to the user settings.
type SettingsData struct {
	MeasurementSystems []units.System
	UserSettings       models.UserSettings
}

// ScraperData holds template data related to the recipe scraper.
type ScraperData struct {
	UnsupportedWebsite string
}

// NewViewRecipeData creates and populates a new ViewRecipeData.
func NewViewRecipeData(id int64, recipe *models.Recipe, isFromHost, isShared bool) *ViewRecipeData {
	return &ViewRecipeData{
		FormattedTimes: newFormattedTimes(recipe.Times),
		ID:             id,
		Inc: func(n int) int {
			return n + 1
		},
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
	Categories     []string
	FormattedTimes formattedTimes
	ID             int64
	Inc            func(n int) int
	IsURL          bool
	IsUUIDValid    bool
	Recipe         *models.Recipe
	Share          shareData
}

func newFormattedTimes(times models.Times) formattedTimes {
	cook := formatDuration(times.Cook, false)
	prep := formatDuration(times.Prep, false)

	parts := strings.Split(prep, "h")
	prepEdit := fmt.Sprintf("%02s:%02s:00", parts[0], parts[1])

	parts = strings.Split(cook, "h")
	cookEdit := fmt.Sprintf("%02s:%02s:00", parts[0], parts[1])

	return formattedTimes{
		Cook:          cook,
		CookDateTime:  formatDuration(times.Cook, true),
		CookEdit:      cookEdit,
		Prep:          prep,
		PrepDateTime:  formatDuration(times.Prep, true),
		PrepEdit:      prepEdit,
		Total:         formatDuration(times.Total, false),
		TotalDateTime: formatDuration(times.Total, true),
	}
}

type formattedTimes struct {
	Cook          string
	CookDateTime  string
	CookEdit      string
	Prep          string
	PrepDateTime  string
	PrepEdit      string
	Total         string
	TotalDateTime string
}

type shareData struct {
	IsFromHost bool
	IsShared   bool
}

// FunctionsData provides functions for use in the templates.
type FunctionsData struct {
	CutString   func(s string, numCharacters int) string
	IsUUIDValid func(u uuid.UUID) bool
}
