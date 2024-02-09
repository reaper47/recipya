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

	counts := make([]count, 0, len(hit))
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
