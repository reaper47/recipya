package statements

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

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
