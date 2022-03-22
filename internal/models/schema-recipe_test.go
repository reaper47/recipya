package models

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

func TestModelRecipeSchema(t *testing.T) {
	imageID := uuid.New()
	rs := RecipeSchema{
		AtContext:     "@Schema",
		AtType:        "Recipe",
		Category:      "lunch",
		CookTime:      "PT3H",
		CookingMethod: "",
		Cuisine:       "",
		DateCreated:   "2022-03-16",
		DateModified:  "2022-03-20",
		DatePublished: "2022-03-16",
		Description:   "description",
		Keywords:      "kw1,kw2,kw3",
		Image:         imageID.String(),
		Ingredients:   []string{"ing1", "ing2", "ing3"},
		Instructions:  instructions{Values: []string{"ins1", "ins2", "ins3"}},
		Name:          "name",
		NutritionSchema: NutritionSchema{
			Calories:       "341kcal",
			Carbohydrates:  "1g",
			Cholesterol:    "2g",
			Fat:            "27g",
			Fiber:          "4g",
			Protein:        "5g",
			SaturatedFat:   "6g",
			Servings:       "7g",
			Sodium:         "8g",
			Sugar:          "9g",
			TransFat:       "10g",
			UnsaturatedFat: "11g",
		},
		PrepTime: "PT1H",
		Tools:    tools{Values: []string{"t1", "t2", "t3"}},
		Yield:    yield{Value: 4},
		Url:      "https://recipes.musicavis.ca",
	}

	t.Run("ToRecipe transform the schema to a Recipe", func(t *testing.T) {
		created, _ := time.Parse("2006-01-02", "2022-03-16")
		updated, _ := time.Parse("2006-01-02", "2022-03-20")
		expected := Recipe{
			Name:        "name",
			Description: "description",
			Image:       imageID,
			Url:         "https://recipes.musicavis.ca",
			Yield:       4,
			Category:    "lunch",
			Times: Times{
				Prep:  1 * time.Hour,
				Cook:  3 * time.Hour,
				Total: 4 * time.Hour,
			},
			Ingredients: []string{"ing1", "ing2", "ing3"},
			Nutrition: Nutrition{
				Calories:           "341kcal",
				TotalCarbohydrates: "1g",
				Sugars:             "9g",
				Protein:            "5g",
				TotalFat:           "27g",
				SaturatedFat:       "6g",
				Cholesterol:        "2g",
				Sodium:             "8g",
				Fiber:              "4g",
			},
			Instructions: []string{"ins1", "ins2", "ins3"},
			Tools:        []string{"t1", "t2", "t3"},
			Keywords:     []string{"kw1", "kw2", "kw3"},
			CreatedAt:    created,
			UpdatedAt:    updated,
		}
		expectedBytes, err := json.Marshal(expected)
		if err != nil {
			t.Fatal(err)
		}

		actual, err := rs.ToRecipe()
		if err != nil {
			t.Fatal(err)
		}
		actualBytes, err := json.Marshal(actual)
		if err != nil {
			t.Fatal(err)
		}

		if !slices.Equal(actualBytes, expectedBytes) {
			fmt.Println(created, updated)
			t.Fatalf("expected:\n%s\nbut got\n%s", expectedBytes, actualBytes)
		}
	})
}
