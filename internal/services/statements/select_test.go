package statements

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func BenchmarkBuildSearchRecipeQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		got := BuildSearchRecipeQuery([]string{"one", "two", "three", "four", "five"}, models.SearchOptionsRecipes{FullSearch: true})
		_ = got
	}
}

func TestSelectSearchRecipe(t *testing.T) {
	testcases := []struct {
		name    string
		queries []string
		options models.SearchOptionsRecipes
		want    string
	}{
		{
			name:    "no queries",
			queries: []string{},
			options: models.SearchOptionsRecipes{ByName: true},
			want:    baseSelectRecipe + ` WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? ORDER BY rank) GROUP BY recipes.id LIMIT 30`,
		},
		{
			name:    "many queries",
			queries: []string{"one", "two", "tHrEe", "FOUR"},
			options: models.SearchOptionsRecipes{ByName: true},
			want:    baseSelectRecipe + ` WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ? AND name MATCH ? AND name MATCH ? AND name MATCH ?) ORDER BY rank) GROUP BY recipes.id LIMIT 30`,
		},
		{
			name:    "full search one query",
			queries: []string{"one"},
			options: models.SearchOptionsRecipes{FullSearch: true},
			want:    baseSelectRecipe + ` WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ? OR description MATCH ? OR category MATCH ? OR ingredients MATCH ? OR instructions MATCH ? OR keywords MATCH ? OR source MATCH ?) ORDER BY rank) GROUP BY recipes.id LIMIT 30`,
		},
		{
			name:    "full search many queries",
			queries: []string{"one", "two"},
			options: models.SearchOptionsRecipes{FullSearch: true},
			want:    baseSelectRecipe + ` WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND ((name MATCH ? AND name MATCH ?) OR (description MATCH ? AND description MATCH ?) OR (category MATCH ? AND category MATCH ?) OR (ingredients MATCH ? AND ingredients MATCH ?) OR (instructions MATCH ? AND instructions MATCH ?) OR (keywords MATCH ? AND keywords MATCH ?) OR (source MATCH ? AND source MATCH ?)) ORDER BY rank) GROUP BY recipes.id LIMIT 30`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildSearchRecipeQuery(tc.queries, tc.options)
			if got != tc.want {
				t.Fatalf("got:\n%q\nbut want:\n%q", got, tc.want)
			}
		})
	}
}
