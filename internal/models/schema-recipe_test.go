package models_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/constants"
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/exp/slices"
)

func TestModelRecipeSchema(t *testing.T) {
	imageID := uuid.New()
	rs := models.RecipeSchema{
		AtContext:     "@Schema",
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: "lunch"},
		CookTime:      "PT3H",
		CookingMethod: models.CookingMethod{Value: ""},
		Cuisine:       models.Cuisine{Value: ""},
		DateCreated:   "2022-03-16",
		DateModified:  "2022-03-20",
		DatePublished: "2022-03-16",
		Description:   models.Description{Value: "description"},
		Keywords:      models.Keywords{Values: "kw1,kw2,kw3"},
		Image:         models.Image{Value: imageID.String()},
		Ingredients:   models.Ingredients{Values: []string{"ing1", "ing2", "ing3"}},
		Instructions:  models.Instructions{Values: []string{"ins1", "ins2", "ins3"}},
		Name:          "name",
		NutritionSchema: models.NutritionSchema{
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
		Tools:    models.Tools{Values: []string{"t1", "t2", "t3"}},
		Yield:    models.Yield{Value: 4},
		URL:      "https://recipes.musicavis.ca",
	}

	t.Run("ToRecipe transform the schema to a Recipe", func(t *testing.T) {
		created, _ := time.Parse(constants.BasicTimeLayout, "2022-03-16")
		updated, _ := time.Parse(constants.BasicTimeLayout, "2022-03-20")
		expected := models.Recipe{
			Name:        "name",
			Description: "description",
			Image:       imageID,
			URL:         "https://recipes.musicavis.ca",
			Yield:       4,
			Category:    "lunch",
			Times: models.Times{
				Prep:  1 * time.Hour,
				Cook:  3 * time.Hour,
				Total: 4 * time.Hour,
			},
			Ingredients: models.Ingredients{Values: []string{"ing1", "ing2", "ing3"}},
			Nutrition: models.Nutrition{
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
