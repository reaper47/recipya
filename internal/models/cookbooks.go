package models

import (
	"github.com/google/uuid"
)

// Cookbook is the struct that holds information on a cookbook.
type Cookbook struct {
	ID      int64
	Count   int64
	Image   uuid.UUID
	Recipes Recipes
	Title   string
}

// MakeView creates a templates.CookbookView from the Cookbook.
// The index is the position of the cookbook in the list of cookbooks presented to the user.
func (c Cookbook) MakeView(index int64, page uint64) CookbookView {
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
	Recipes     Recipes
	PageNumber  uint64
	PageItemID  int64
	Title       string
}
