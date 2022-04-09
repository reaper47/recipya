package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestModelRecipe(t *testing.T) {
	t.Run("Recipe ToArgs is correct", func(t *testing.T) {
		r := Recipe{
			ID:          1,
			Name:        "name",
			Description: "description",
			Image:       uuid.New(),
			Url:         "https://www.google.com",
			Yield:       4,
			Category:    "lunch",
			Times: Times{
				Prep:  1 * time.Hour,
				Cook:  2 * time.Hour,
				Total: 3 * time.Hour,
			},
			Ingredients: []string{"ing1", "ing2", "ing3"},
			Nutrition: Nutrition{
				Calories:           "1kcal",
				TotalCarbohydrates: "1g",
				Sugars:             "2g",
				Protein:            "3g",
				TotalFat:           "4g",
				SaturatedFat:       "5g",
				Cholesterol:        "6g",
				Sodium:             "7g",
				Fiber:              "8g",
			},
			Instructions: []string{"ins1", "ins2", "ins3"},
			Tools:        []string{"t1", "t2", "t3"},
			Keywords:     []string{"kw1", "kw2", "kw3"},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now().Add(2 * time.Hour),
		}

		argsWithID := r.ToArgs(true)
		argsWithoutID := r.ToArgs(false)

		numElements := 30
		if len(argsWithID) != numElements {
			t.Errorf("wanted %d elements but got %d", len(argsWithID), numElements)
		}
		if len(argsWithoutID) != numElements-1 {
			t.Errorf("wanted %d elements but got %d", len(argsWithoutID), numElements-1)
		}
	})

	t.Run("Recipe ToSchema transforms the Recipe to its schema successfully", func(t *testing.T) {
		imageUUID := uuid.New()
		r := Recipe{
			ID:          1,
			Name:        "name",
			Description: "description",
			Image:       imageUUID,
			Url:         "https://www.google.com",
			Yield:       4,
			Category:    "lunch",
			Times: Times{
				Prep:  1 * time.Hour,
				Cook:  2 * time.Hour,
				Total: 3 * time.Hour,
			},
			Ingredients: []string{"ing1", "ing2", "ing3"},
			Nutrition: Nutrition{
				Calories:           "1kcal",
				TotalCarbohydrates: "1g",
				Sugars:             "2g",
				Protein:            "3g",
				TotalFat:           "4g",
				SaturatedFat:       "5g",
				Cholesterol:        "6g",
				Sodium:             "7g",
				Fiber:              "8g",
			},
			Instructions: []string{"ins1", "ins2", "ins3"},
			Tools:        []string{"t1", "t2", "t3"},
			Keywords:     []string{"kw1", "kw2", "kw3"},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now().Add(2 * time.Hour),
		}

		schema := r.ToSchema()

		if schema.AtContext != "http://schema.org" {
			t.Errorf("context must be http://schema.org")
		}
		if schema.AtType != "Recipe" {
			t.Errorf("type must be Recipe")
		}
		if schema.Category != "lunch" {
			t.Errorf("wanted category 'lunch' but got '%q'", schema.Category)
		}
		if schema.CookTime != "PT2H0M0S" {
			t.Errorf("wanted cook time PT1H but got %q", schema.CookTime)
		}
		if schema.CookingMethod != "" {
			t.Errorf("wanted an empty cooking method but got %q", schema.CookingMethod)
		}
		if schema.Cuisine != "" {
			t.Errorf("wanted an empty cusine but got %q", schema.Cuisine)
		}
		v := r.CreatedAt.Format("2006-01-02")
		if schema.DateCreated != v {
			t.Errorf("wanted date created %q but got %q", v, schema.DateCreated)
		}
		v = r.UpdatedAt.Format("2006-01-02")
		if schema.DateModified != v {
			t.Errorf("wanted date modified %q but got %q", v, schema.DateModified)
		}
		v = r.CreatedAt.Format("2006-01-02")
		if schema.DatePublished != v {
			t.Errorf("wanted date published %q but got %q", v, schema.DatePublished)
		}
		if schema.Description != "description" {
			t.Errorf("wanted description 'description' but got %q", schema.Description)
		}
		if schema.Keywords != "kw1,kw2,kw3" {
			t.Errorf("wanted keywords 'kw1,kw2,kw3' but got %q", schema.Keywords)
		}
		v = string(imageUUID.String())
		if schema.Image != v {
			t.Errorf("wanted uuid %q but got %q", v, schema.Image)
		}
		// TODO: Uncomment when Fedora updated its Go package to 1.18
		/*if !slices.Equal(schema.Ingredients, []string{"ing1", "ing2", "ing3"}) {
			t.Errorf("wanted ingredients []string{ing1, ing2, ing3} but got %v", schema.Ingredients)
		}
		if !slices.Equal(schema.Instructions.Values, []string{"ins1", "ins2", "ins3"}) {
			t.Errorf("wanted instructions []string{ins1, ins2, ins3} but got %v", schema.Instructions.Values)
		}*/
		if schema.Name != "name" {
			t.Errorf("wanted name 'name' but got %q", schema.Name)
		}
		if schema.NutritionSchema != r.Nutrition.toSchema("4") {
			t.Errorf("wanted nutrition but got %v", schema.NutritionSchema)
		}
		if schema.PrepTime != "PT1H0M0S" {
			t.Errorf("wanted prepTime PT1H0M0S but got %q", schema.PrepTime)
		}
		// TODO: Uncomment when Fedora updated its Go package to 1.18
		/*if !slices.Equal(schema.Tools.Values, []string{"t1", "t2", "t3"}) {
			t.Errorf("wanted tools []string{t1, t2, t3} but got %v", schema.Tools.Values)
		}*/
		if schema.Yield.Value != 4 {
			t.Errorf("wanted yield 4 but got %d", schema.Yield.Value)
		}
		if schema.Url != "https://www.google.com" {
			t.Errorf("wanted url https://www.google.com but got %q", schema.Url)
		}
	})

	t.Run("Creating a new Times parses correctly", func(t *testing.T) {
		actual, err := NewTimes("PT1H0M0S", "PT2H0M0S")
		if err != nil {
			t.Fatal(err)
		}

		if actual.Cook != 2*time.Hour {
			t.Errorf("wanted cooking time 2H but got %v", actual.Cook.String())
		}
		if actual.Prep != 1*time.Hour {
			t.Errorf("wanted prep time 1H but got %v", actual.Prep.String())
		}
		if actual.Total != 3*time.Hour {
			t.Errorf("wanted total time 3H but got %v", actual.Total.String())
		}
	})
}
