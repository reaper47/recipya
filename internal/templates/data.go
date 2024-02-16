package templates

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/units"
	"html/template"
	"strings"
)

// Data holds data to pass on to the templates.
type Data struct {
	IsAdmin          bool // IsAdmin indicates whether the user is an administration.
	IsAuthenticated  bool // IsAuthenticated indicates whether the user is authenticated.
	IsAutologin      bool // IsAutologin indicates whether the user enabled autologin.
	IsDemo           bool // IsDemo indicates whether running instance is the demo version.
	IsNoSignups      bool // IsNoSignups indicates whether account registrations are disabled.
	IsHxRequest      bool // IsHxRequest indicates whether the request is an HX one. It is used for oop swaps.
	IsToastWSVisible bool // IsToastWSVisible indicates whether to display the notification for websocket tasks.

	Title        string        // Title is the text inserted <title> tag's text.
	Content      string        // Content is text to insert into the template.
	ContentHTML  template.HTML // ContentHTML is the non-escaped HTML to insert into the template.
	ContentTitle string        // ContentTitle is the header of the Content.

	Functions FunctionsData[int64]

	About           AboutData
	Admin           AdminData
	CookbookFeature CookbookFeature
	Pagination      Pagination
	Recipes         models.Recipes
	Reports         ReportsData
	Settings        SettingsData
	Scraper         ScraperData
	View            *ViewRecipeData
}

// NewAboutData creates a new instance of AboutData.
func NewAboutData() AboutData {
	return AboutData{
		Version: app.Version,
	}
}

// AboutData holds general application data.
type AboutData struct {
	Version string
}

// AdminData holds data for the admin page.
type AdminData struct {
	Users []models.User
}

// CookbookFeature is the data to pass related to the cookbook feature.
type CookbookFeature struct {
	Cookbooks    []models.Cookbook
	Cookbook     CookbookView
	MakeCookbook func(index int64, cookbook models.Cookbook, page uint64) CookbookView
	ShareData    ShareData
	ViewMode     models.ViewMode
}

// MakeCookbookView creates a templates.CookbookView from the Cookbook.
// The index is the position of the cookbook in the list of cookbooks presented to the user.
func MakeCookbookView(c models.Cookbook, index int64, page uint64) CookbookView {
	return CookbookView{
		ID:          c.ID,
		Image:       c.Image,
		IsUUIDValid: c.Image != uuid.Nil,
		NumRecipes:  c.Count,
		PageNumber:  page,
		PageItemID:  index + 1,
		Recipes:     c.Recipes,
		Title:       c.Title,
	}
}

// CookbookView holds data related to viewing a cookbook.
type CookbookView struct {
	ID          int64
	Image       uuid.UUID
	IsUUIDValid bool
	NumRecipes  int64
	Recipes     models.Recipes
	PageNumber  uint64
	PageItemID  int64
	Title       string
}

// NewFunctionsData initializes a new FunctionsData.
func NewFunctionsData[T int64 | uint64]() FunctionsData[T] {
	return FunctionsData[T]{
		CutString: func(s string, numCharacters int) string {
			if len(s) < numCharacters {
				return s
			}
			return s[:numCharacters] + "…"
		},
		IsUUIDValid: func(u uuid.UUID) bool {
			return u != uuid.Nil
		},
		MulAll: func(vals ...T) T {
			res := T(1)
			for _, v := range vals {
				res *= v
			}
			return res
		},
	}
}

// FunctionsData provides functions for use in the templates.
type FunctionsData[T int64 | uint64] struct {
	CutString   func(s string, numCharacters int) string
	IsUUIDValid func(u uuid.UUID) bool
	MulAll      func(vals ...T) T
}

// RegisterData is the data to pass on to the user registration template.
type RegisterData struct {
	Email           string
	PasswordConfirm string
}

// SettingsData holds template data related to the user settings.
type SettingsData struct {
	Backups            []Backup
	MeasurementSystems []units.System
	UserSettings       models.UserSettings
}

// Backup holds data related to backups.
type Backup struct {
	Display string
	Value   string
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
	if len(parts) > 1 {
		minutes := strings.Split(parts[1], "m")
		if len(minutes) > 0 {
			prepEdit = fmt.Sprintf("%02s:%02s:00", parts[0], minutes[0])
		}
	}

	var cookEdit string
	parts = strings.Split(cook, "h")
	if len(parts) > 1 {
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

// ReportsData holds data related to reports.
type ReportsData struct {
	CurrentReport []models.ReportLog
	Imports       []models.Report
	Sort          string
}

// ShareData holds information on the entity being shared.
type ShareData struct {
	IsFromHost bool
	IsShared   bool
}

// ErrorTokenExpired encapsulates the information displayed to the user when a token is expired.
var ErrorTokenExpired = Data{
	Title: "Token Expired",
	Content: `The token associated with the URL expired.
				The problem has been forwarded to our team automatically. We will look into it and come
                back to you. We apologise for this inconvenience.`,
}
