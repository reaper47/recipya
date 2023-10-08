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
	IsAuthenticated bool // IsAuthenticated indicates whether the user is authenticated.
	IsHxRequest     bool // IsHxRequest indicates whether the request is an HX one. It is used for oop swaps.

	Title string // Title is the text inserted <title> tag's text.

	Content      string // Content is text to insert into the template.
	ContentTitle string // ContentTitle is the header of the Content.

	Functions FunctionsData

	CookbookFeature CookbookFeature
	Recipes         models.Recipes
	Settings        SettingsData
	Scraper         ScraperData
	View            *ViewRecipeData
}

// CookbookFeature is the data to pass related to the cookbook feature.
type CookbookFeature struct {
	Cookbooks []models.Cookbook
	ViewMode  string
}

// FunctionsData provides functions for use in the templates.
type FunctionsData struct {
	CutString   func(s string, numCharacters int) string
	Inc         func(v int64) int64
	IsUUIDValid func(u uuid.UUID) bool
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
	minutes := strings.Split(parts[1], "m")[0]
	prepEdit := fmt.Sprintf("%02s:%02s:00", parts[0], minutes)

	parts = strings.Split(cook, "h")
	minutes = strings.Split(parts[1], "m")[0]
	cookEdit := fmt.Sprintf("%02s:%02s:00", parts[0], minutes)

	prep = strings.TrimPrefix(prep, "0h")
	prep = strings.TrimPrefix(prep, "0")

	cook = strings.TrimPrefix(cook, "0h")
	cook = strings.TrimPrefix(cook, "0")

	total := formatDuration(times.Total, false)
	total = strings.TrimPrefix(total, "0h")
	total = strings.TrimPrefix(total, "0")

	return formattedTimes{
		Cook:          cook,
		CookDateTime:  formatDuration(times.Cook, true),
		CookEdit:      cookEdit,
		Prep:          prep,
		PrepDateTime:  formatDuration(times.Prep, true),
		PrepEdit:      prepEdit,
		Total:         total,
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
