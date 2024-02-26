package statements

import (
	"github.com/reaper47/recipya/internal/models"
	"strings"
	"testing"
)

func BenchmarkBuildSearchRecipeQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		got := BuildSearchRecipeQuery([]string{"one", "two", "three", "four", "five"}, models.SearchOptionsRecipes{FullSearch: true})
		_ = got
	}
}

func BenchmarkBuildSelectNutrientFDC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		got := BuildSelectNutrientFDC([]string{"one", "two", "three", "four", "five"})
		_ = got
	}
}

func TestBuildBaseSelectRecipe(t *testing.T) {
	testcases := []struct {
		name string
		in   models.Sort
		want string
	}{
		{
			name: "A-Z",
			in:   models.Sort{IsAToZ: true},
			want: "ROW_NUMBER() OVER (ORDER BY recipes.name ASC) AS row_num",
		},
		{
			name: "Z-A",
			in:   models.Sort{IsZToA: true},
			want: "ROW_NUMBER() OVER (ORDER BY recipes.name DESC) AS row_num",
		},
		{
			name: "no options",
			in:   models.Sort{},
			want: "",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildBaseSelectRecipe(tc.in)

			before, after1, _ := strings.Cut(got, "FROM recipes")
			if !strings.Contains(before, tc.want) {
				t.Fatalf("expected %q in SELECT of query", tc.want)
			}

			_, after2, _ := strings.Cut(baseSelectRecipe, "FROM recipes")
			if after1 != after2 {
				t.Fatal("FROM recipes bit from baseRecipes variable not equal")
			}
		})
	}
}

func TestBuildSelectNutrientFDC(t *testing.T) {
	testcases := []struct {
		name        string
		ingredients []string
		want        string
	}{
		{
			name:        "one ingredient",
			ingredients: []string{"one"},
			want:        "WHERE description LIKE '%one%'",
		},
		{
			name:        "multiple ingredients",
			ingredients: []string{"one", "two"},
			want:        "WHERE description LIKE '%one%' AND description LIKE '%two%'",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildSelectNutrientFDC(tc.ingredients)
			if !strings.Contains(got, tc.want) {
				t.Fatalf("expected %q in query", tc.want)
			}
		})
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
		{
			name:    "full search one query sort A-Z",
			queries: []string{"one"},
			options: models.SearchOptionsRecipes{ByName: true, Sort: models.Sort{IsAToZ: true}},
			want:    "SELECT * FROM (" + BuildBaseSelectRecipe(models.Sort{IsAToZ: true}) + ` WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ?) ORDER BY rank) GROUP BY recipes.id)`,
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
