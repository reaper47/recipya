package models

import (
	"cmp"
	"github.com/google/uuid"
	"slices"
)

type count struct {
	name string
	hits int
}

// Cookbook is the struct that holds information on a cookbook.
type Cookbook struct {
	ID      int64
	Count   int64
	Image   uuid.UUID
	Recipes Recipes
	Title   string
}

// DominantCategories returns the `n` most common categories of recipes in the cookbook.
// If there are fewer than `n` categories, all categories are returned.
func (c Cookbook) DominantCategories(n int) []string {
	hit := make(map[string]int)
	for _, r := range c.Recipes {
		hit[r.Category]++
	}

	var counts []count
	for k, v := range hit {
		counts = append(counts, count{k, v})
	}

	slices.SortFunc(counts, func(a, b count) int {
		return cmp.Compare(b.hits, a.hits)
	})

	if len(counts) < n {
		n = len(counts)
	}

	categories := make([]string, n)
	for i := 0; i < n; i++ {
		categories[i] = counts[i].name
	}
	return categories
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
