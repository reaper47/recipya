package statements

import (
	"github.com/reaper47/recipya/internal/models"
	"strings"
	"testing"
)

func BenchmarkBuildPaginatedResultsQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		got := buildSelectPaginatedResultsQuery(models.SearchOptionsRecipes{Query: "one two three four"})
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
		options models.SearchOptionsRecipes
		want    string
	}{
		{
			name: "no queries",
			want: "SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? ORDER BY rank) GROUP BY recipes.id)",
		},
		{
			name: "advanced category only",
			options: models.SearchOptionsRecipes{
				Advanced: models.AdvancedSearch{Category: "breakfast"},
			},
			want: "SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND recipes_fts MATCH ? ORDER BY rank) GROUP BY recipes.id)",
		},
		{
			name: "advanced multiple categories",
			options: models.SearchOptionsRecipes{
				Advanced: models.AdvancedSearch{Category: "breakfast,dinner"},
			},
			want: "SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND recipes_fts MATCH ? ORDER BY rank) GROUP BY recipes.id)",
		},
		{
			name:    "one query",
			options: models.SearchOptionsRecipes{Query: "one two three four"},
			want:    "SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND recipes_fts MATCH ? ORDER BY rank) GROUP BY recipes.id)",
		},
		{
			name: "one query with advanced search",
			options: models.SearchOptionsRecipes{
				Advanced: models.AdvancedSearch{Category: "breakfast"},
				Query:    "one two three four",
			},
			want: "SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND recipes_fts MATCH ? ORDER BY rank) GROUP BY recipes.id)",
		},
		{
			name:    "cookbook search",
			options: models.SearchOptionsRecipes{Query: "choco", CookbookID: 1},
			want:    "SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND recipes_fts MATCH ? ORDER BY rank) AND recipes.id NOT IN (SELECT recipe_id FROM cookbook_recipes WHERE cookbook_id = ?) GROUP BY recipes.id)",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildSearchRecipeQuery(tc.options)
			compareSQL(t, got, tc.want)
		})
	}
}

func TestBuildSelectPaginatedResults(t *testing.T) {
	testcases := []struct {
		name    string
		options models.SearchOptionsRecipes
		want    string
	}{
		{
			name:    "empty query",
			options: models.SearchOptionsRecipes{Page: 1},
			want:    "WITH results AS (SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? ORDER BY rank) GROUP BY recipes.id)) SELECT * FROM results WHERE row_num BETWEEN 1 AND 15",
		},
		{
			name:    "full search one query",
			options: models.SearchOptionsRecipes{Query: "one two three four", Page: 2},
			want:    "WITH results AS (SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND recipes_fts MATCH ? ORDER BY rank) GROUP BY recipes.id)) SELECT * FROM results WHERE row_num BETWEEN 16 AND 30",
		},
		{
			name:    "with advanced",
			options: models.SearchOptionsRecipes{Query: "one two", Page: 1, Advanced: models.AdvancedSearch{Category: "breakfast"}},
			want:    "WITH results AS (SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND recipes_fts MATCH ? ORDER BY rank) GROUP BY recipes.id)) SELECT * FROM results WHERE row_num BETWEEN 1 AND 15",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildSelectPaginatedResults(tc.options)
			compareSQL(t, got, tc.want)
		})
	}
}

func TestBuildSelectSearchResultsCount(t *testing.T) {
	testcases := []struct {
		name    string
		queries []string
		options models.SearchOptionsRecipes
		want    string
	}{
		{
			name:    "empty query",
			options: models.SearchOptionsRecipes{Page: 1},
			want:    "WITH results AS (SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? ORDER BY rank) GROUP BY recipes.id))SELECT COUNT(*) FROM results",
		},
		{
			name:    "full search one query",
			options: models.SearchOptionsRecipes{Query: "one two three four", Page: 3},
			want:    "WITH results AS (SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND recipes_fts MATCH ? ORDER BY rank) GROUP BY recipes.id))SELECT COUNT(*) FROM results",
		},
		{
			name:    "with advanced",
			options: models.SearchOptionsRecipes{Query: "one two three four", Page: 3, Advanced: models.AdvancedSearch{Category: "breakfast", Text: "one two three four"}},
			want:    "WITH results AS (SELECT recipe_id, name, description, image, created_at, category, row_num FROM ( SELECT recipes.id AS recipe_id, recipes.name AS name, recipes.description AS description, recipes.image AS image, recipes.created_at AS created_at, categories.name AS category, user_id, ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num FROM recipes LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id LEFT JOIN categories ON category_recipe.category_id = categories.id LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? AND recipes_fts MATCH ? ORDER BY rank) GROUP BY recipes.id))SELECT COUNT(*) FROM results",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildSelectSearchResultsCount(tc.options)
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
