package statements

import (
	"github.com/reaper47/recipya/internal/models"
	"strings"
	"testing"
)

func BenchmarkBuildPaginatedResultsQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		got := buildSelectPaginatedResultsQuery([]string{"one", "two", "three", "four", "five"}, models.SearchOptionsRecipes{IsFullSearch: true})
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
			name: "new to old",
			in:   models.Sort{IsNewestToOldest: true},
			want: "ROW_NUMBER() OVER (ORDER BY recipes.created_at ASC) AS row_num",
		},
		{
			name: "old to new",
			in:   models.Sort{IsOldestToNewest: true},
			want: "ROW_NUMBER() OVER (ORDER BY recipes.created_at DESC) AS row_num",
		},
		{
			name: "default",
			in:   models.Sort{IsDefault: true},
			want: "ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num",
		},
		{
			name: "random",
			in:   models.Sort{IsRandom: true},
			want: "ROW_NUMBER() OVER (ORDER BY RANDOM()) AS row_num",
		},
		{
			name: "no options",
			in:   models.Sort{},
			want: "ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildBaseSelectRecipe(tc.in)

			before, after1, _ := strings.Cut(got, "FROM recipes")
			if !strings.Contains(before, tc.want) {
				t.Fatalf("expected %q in SELECT of query", tc.want)
			}

			_, after2, _ := strings.Cut(baseSelectSearchRecipe, "FROM recipes")
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
			options: models.SearchOptionsRecipes{IsByName: true},
			want:    "SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? ORDER BY rank) GROUP BY recipes.id)",
		},
		{
			name:    "many queries",
			queries: []string{"one", "two", "tHrEe", "FOUR"},
			options: models.SearchOptionsRecipes{IsByName: true},
			want:    "SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ? AND name MATCH ? AND name MATCH ? AND name MATCH ?) ORDER BY rank) GROUP BY recipes.id)",
		},
		{
			name:    "full search one query",
			queries: []string{"one"},
			options: models.SearchOptionsRecipes{IsFullSearch: true},
			want:    "SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ? OR description MATCH ? OR category MATCH ? OR ingredients MATCH ? OR instructions MATCH ? OR keywords MATCH ? OR source MATCH ?) ORDER BY rank) GROUP BY recipes.id)",
		},
		{
			name:    "full search many queries",
			queries: []string{"one", "two"},
			options: models.SearchOptionsRecipes{IsFullSearch: true},
			want:    "SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND ((name MATCH ? AND name MATCH ?) OR (description MATCH ? AND description MATCH ?) OR (category MATCH ? AND category MATCH ?) OR (ingredients MATCH ? AND ingredients MATCH ?) OR (instructions MATCH ? AND instructions MATCH ?) OR (keywords MATCH ? AND keywords MATCH ?) OR (source MATCH ? AND source MATCH ?)) ORDER BY rank) GROUP BY recipes.id)",
		},
		{
			name:    "full search one query sort A-Z",
			queries: []string{"one"},
			options: models.SearchOptionsRecipes{IsByName: true, Sort: models.Sort{IsAToZ: true}},
			want:    "SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.name ASC) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ?) ORDER BY rank) GROUP BY recipes.id)",
		},
		{
			name:    "cookbook search",
			queries: []string{"choco"},
			options: models.SearchOptionsRecipes{IsByName: true, CookbookID: 1},
			want:    "SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ?) ORDER BY rank) AND recipes.id NOT IN (SELECT recipe_id FROM cookbook_recipes WHERE cookbook_id = ?) GROUP BY recipes.id)",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildSearchRecipeQuery(tc.queries, tc.options)
			compareSQL(t, got, tc.want)
		})
	}
}

func TestBuildSelectPaginatedResults(t *testing.T) {
	testcases := []struct {
		name    string
		queries []string
		page    uint64
		options models.SearchOptionsRecipes
		want    string
	}{
		{
			name:    "empty query",
			queries: make([]string, 0),
			page:    1,
			options: models.SearchOptionsRecipes{},
			want:    "WITH results AS (SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? ORDER BY rank) GROUP BY recipes.id)) SELECT * FROM results WHERE row_num BETWEEN 1 AND 15",
		},
		{
			name:    "many queries",
			queries: []string{"one", "two", "tHrEe", "FOUR"},
			options: models.SearchOptionsRecipes{IsByName: true},
			want:    "WITH results AS (SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ? AND name MATCH ? AND name MATCH ? AND name MATCH ?) ORDER BY rank) GROUP BY recipes.id)) SELECT * FROM results WHERE row_num BETWEEN 1 AND 15",
		},
		{
			name:    "full search one query",
			queries: []string{"one"},
			page:    2,
			options: models.SearchOptionsRecipes{IsFullSearch: true},
			want:    "WITH results AS (SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ? OR description MATCH ? OR category MATCH ? OR ingredients MATCH ? OR instructions MATCH ? OR keywords MATCH ? OR source MATCH ?) ORDER BY rank) GROUP BY recipes.id)) SELECT * FROM results WHERE row_num BETWEEN 16 AND 30",
		},
		{
			name:    "full search many queries",
			queries: []string{"one", "two"},
			options: models.SearchOptionsRecipes{IsFullSearch: true},
			want:    "WITH results AS (SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND ((name MATCH ? AND name MATCH ?) OR (description MATCH ? AND description MATCH ?) OR (category MATCH ? AND category MATCH ?) OR (ingredients MATCH ? AND ingredients MATCH ?) OR (instructions MATCH ? AND instructions MATCH ?) OR (keywords MATCH ? AND keywords MATCH ?) OR (source MATCH ? AND source MATCH ?)) ORDER BY rank) GROUP BY recipes.id)) SELECT * FROM results WHERE row_num BETWEEN 1 AND 15",
		},
		{
			name:    "full search one query sort A-Z",
			queries: []string{"one"},
			options: models.SearchOptionsRecipes{IsByName: true, Sort: models.Sort{IsAToZ: true}},
			want:    "WITH results AS (SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.name ASC) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ?) ORDER BY rank) GROUP BY recipes.id)) SELECT * FROM results WHERE row_num BETWEEN 1 AND 15",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildSelectPaginatedResults(tc.queries, tc.page, tc.options)
			compareSQL(t, got, tc.want)
		})
	}
}

func TestBuildSelectSearchResultsCount(t *testing.T) {
	testcases := []struct {
		name    string
		queries []string
		page    uint64
		options models.SearchOptionsRecipes
		want    string
	}{
		{
			name:    "empty query",
			queries: make([]string, 0),
			page:    1,
			options: models.SearchOptionsRecipes{},
			want:    "WITH results AS (SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? ORDER BY rank) GROUP BY recipes.id))SELECT COUNT(*) FROM results",
		},
		{
			name:    "many queries",
			queries: []string{"one", "two", "tHrEe", "FOUR"},
			page:    2,
			options: models.SearchOptionsRecipes{IsByName: true},
			want:    "WITH results AS (SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ? AND name MATCH ? AND name MATCH ? AND name MATCH ?) ORDER BY rank) GROUP BY recipes.id))SELECT COUNT(*) FROM results",
		},
		{
			name:    "full search one query",
			queries: []string{"one"},
			page:    3,
			options: models.SearchOptionsRecipes{IsFullSearch: true},
			want:    "WITH results AS (SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ? OR description MATCH ? OR category MATCH ? OR ingredients MATCH ? OR instructions MATCH ? OR keywords MATCH ? OR source MATCH ?) ORDER BY rank) GROUP BY recipes.id))SELECT COUNT(*) FROM results",
		},
		{
			name:    "full search many queries",
			queries: []string{"one", "two"},
			options: models.SearchOptionsRecipes{IsFullSearch: true},
			want:    "WITH results AS (SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND ((name MATCH ? AND name MATCH ?) OR (description MATCH ? AND description MATCH ?) OR (category MATCH ? AND category MATCH ?) OR (ingredients MATCH ? AND ingredients MATCH ?) OR (instructions MATCH ? AND instructions MATCH ?) OR (keywords MATCH ? AND keywords MATCH ?) OR (source MATCH ? AND source MATCH ?)) ORDER BY rank) GROUP BY recipes.id))SELECT COUNT(*) FROM results",
		},
		{
			name:    "full search one query sort A-Z",
			queries: []string{"one"},
			options: models.SearchOptionsRecipes{IsByName: true, Sort: models.Sort{IsAToZ: true}},
			want:    "WITH results AS (SELECT * FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND (name MATCH ?) ORDER BY rank) GROUP BY recipes.id))SELECT COUNT(*) FROM results",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildSelectSearchResultsCount(tc.queries, tc.options)
			compareSQL(t, got, tc.want)
		})
	}
}

func compareSQL(tb testing.TB, got, want string) {
	tb.Helper()
	got = strings.Join(strings.Fields(strings.TrimSpace(got)), " ")
	want = strings.Join(strings.Fields(strings.TrimSpace(want)), " ")
	if got != want {
		tb.Fatalf("got:\n%q\nbut want:\n%q", got, want)
	}
}
