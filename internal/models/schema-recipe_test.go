package models_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"slices"
)

func TestRecipeSchema_Recipe(t *testing.T) {
	imageID := uuid.New()
	rs := models.RecipeSchema{
		AtContext:     "@Schema",
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: "lunch"},
		CookTime:      "PT3H",
		CookingMethod: models.CookingMethod{Value: ""},
		Cuisine:       models.Cuisine{Value: "american"},
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

	created, _ := time.Parse(time.DateOnly, "2022-03-16")
	updated, _ := time.Parse(time.DateOnly, "2022-03-20")
	expected := models.Recipe{
		Category:     "lunch",
		CreatedAt:    created,
		Cuisine:      "american",
		Description:  "description",
		Image:        imageID,
		Ingredients:  []string{"ing1", "ing2", "ing3"},
		Instructions: []string{"ins1", "ins2", "ins3"},
		Keywords:     []string{"kw1", "kw2", "kw3"},
		Name:         "name",
		Nutrition: models.Nutrition{
			Calories:           "341kcal",
			Cholesterol:        "2g",
			TotalFat:           "27g",
			Fiber:              "4g",
			Protein:            "5g",
			SaturatedFat:       "6g",
			Sodium:             "8g",
			Sugars:             "9g",
			TotalCarbohydrates: "1g",
			UnsaturatedFat:     "11g",
		},
		Times: models.Times{
			Prep:  1 * time.Hour,
			Cook:  3 * time.Hour,
			Total: 4 * time.Hour,
		},
		Tools:     []string{"t1", "t2", "t3"},
		UpdatedAt: updated,
		URL:       "https://recipes.musicavis.ca",
		Yield:     4,
	}
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := rs.Recipe()
	if err != nil {
		t.Fatal(err)
	}
	actualBytes, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(actualBytes, expectedBytes) {
		t.Logf(cmp.Diff(*actual, expected))
		t.Fail()
	}
}
