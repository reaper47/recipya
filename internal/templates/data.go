package templates

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/units"
	"html/template"
	"strings"
)

// Data holds data to pass on to the templates.
type Data struct {
	IsAuthenticated  bool // IsAuthenticated indicates whether the user is authenticated.
	IsDemo           bool // IsDemo indicates whether running instance is the demo version.
	IsHxRequest      bool // IsHxRequest indicates whether the request is an HX one. It is used for oop swaps.
	IsToastWSVisible bool // IsToastWSVisible indicates whether to display the notification for websocket tasks.

	Title        string        // Title is the text inserted <title> tag's text.
	Content      string        // Content is text to insert into the template.
	ContentHTML  template.HTML // ContentHTML is the non-escaped HTML to insert into the template.
	ContentTitle string        // ContentTitle is the header of the Content.

	Functions FunctionsData

	About           AboutData
	CookbookFeature CookbookFeature
	Pagination      Pagination
	Recipes         models.Recipes
	Settings        SettingsData
	Scraper         ScraperData
	View            *ViewRecipeData
}

// AboutData holds general application data.
type AboutData struct {
	Version string
}

// CookbookFeature is the data to pass related to the cookbook feature.
type CookbookFeature struct {
	Cookbooks    []models.Cookbook
	Cookbook     models.CookbookView
	MakeCookbook func(index int64, cookbook models.Cookbook, page uint64) models.CookbookView
	ShareData    ShareData
	ViewMode     models.ViewMode
}

// NewFunctionsData initializes a new FunctionsData.
func NewFunctionsData() FunctionsData {
	return FunctionsData{
		CutString: func(s string, numCharacters int) string {
			if len(s) < numCharacters {
				return s
			}
			return s[:numCharacters] + "…"
		},
		Dec: func(v int64) int64 {
			return v - 1
		},
		Inc: func(v int64) int64 {
			return v + 1
		},
		IsUUIDValid: func(u uuid.UUID) bool {
			return u != uuid.Nil
		},
		MulAll: func(vals ...int64) int64 {
			res := int64(1)
			for _, v := range vals {
				res *= v
			}
			return res
		},
	}
}

// FunctionsData provides functions for use in the templates.
type FunctionsData struct {
	CutString   func(s string, numCharacters int) string
	Dec         func(v int64) int64
	Inc         func(v int64) int64
	IsUUIDValid func(u uuid.UUID) bool
	MulAll      func(vals ...int64) int64
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
		Share: ShareData{
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
	Share          ShareData
}

func newFormattedTimes(times models.Times) formattedTimes {
	cook := formatDuration(times.Cook, false)
	prep := formatDuration(times.Prep, false)

	var prepEdit string
	parts := strings.Split(prep, "h")
	if len(parts) > 0 {
		minutes := strings.Split(parts[1], "m")
		if len(minutes) > 0 {
			prepEdit = fmt.Sprintf("%02s:%02s:00", parts[0], minutes[0])
		}
	}

	var cookEdit string
	parts = strings.Split(cook, "h")
	if len(parts) > 0 {
		minutes := strings.Split(parts[1], "m")
		if len(minutes) > 0 {
			cookEdit = fmt.Sprintf("%02s:%02s:00", parts[0], minutes[0])
		}
	}

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

// ShareData holds information on the entity being shared.
type ShareData struct {
	IsFromHost bool
	IsShared   bool
}
