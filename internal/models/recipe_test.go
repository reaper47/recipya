package models_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
	"time"

	"golang.org/x/exp/slices"

	"github.com/google/uuid"
)

func TestRecipe_Schema(t *testing.T) {
	imageUUID := uuid.New()
	r := models.Recipe{
		ID:          1,
		Name:        "name",
		Description: "description",
		Image:       imageUUID,
		URL:         "https://www.google.com",
		Yield:       4,
		Category:    "lunch",
		Times: models.Times{
			Prep:  1 * time.Hour,
			Cook:  2 * time.Hour,
			Total: 3 * time.Hour,
		},
		Ingredients: models.Ingredients{Values: []string{"ing1", "ing2", "ing3"}},
		Nutrition: models.Nutrition{
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

	schema := r.Schema()

	if schema.AtContext != "https://schema.org" {
		t.Errorf("context must be https://schema.org")
	}
	if schema.AtType.Value != "Recipe" {
		t.Errorf("type must be Recipe")
	}
	if schema.Category.Value != "lunch" {
		t.Errorf("wanted category 'lunch' but got '%q'", schema.Category)
	}
	if schema.CookTime != "PT2H0M0S" {
		t.Errorf("wanted cook time PT1H but got %q", schema.CookTime)
	}
	if schema.CookingMethod.Value != "" {
		t.Errorf("wanted an empty cooking method but got %q", schema.CookingMethod)
	}
	if schema.Cuisine.Value != "" {
		t.Errorf("wanted an empty cusine but got %q", schema.Cuisine)
	}
	v := r.CreatedAt.Format(time.DateOnly)
	if schema.DateCreated != v {
		t.Errorf("wanted date created %q but got %q", v, schema.DateCreated)
	}
	v = r.UpdatedAt.Format(time.DateOnly)
	if schema.DateModified != v {
		t.Errorf("wanted date modified %q but got %q", v, schema.DateModified)
	}
	v = r.CreatedAt.Format(time.DateOnly)
	if schema.DatePublished != v {
		t.Errorf("wanted date published %q but got %q", v, schema.DatePublished)
	}
	if schema.Description.Value != "description" {
		t.Errorf("wanted description 'description' but got %q", schema.Description)
	}
	if schema.Keywords.Values != "kw1,kw2,kw3" {
		t.Errorf("wanted keywords 'kw1,kw2,kw3' but got %q", schema.Keywords)
	}
	v = imageUUID.String()
	if schema.Image.Value != v {
		t.Errorf("wanted uuid %q but got %q", v, schema.Image)
	}

	if !slices.Equal(schema.Ingredients.Values, []string{"ing1", "ing2", "ing3"}) {
		t.Errorf("wanted ingredients []string{ing1, ing2, ing3} but got %v", schema.Ingredients)
	}
	if !slices.Equal(schema.Instructions.Values, []string{"ins1", "ins2", "ins3"}) {
		t.Errorf(
			"wanted instructions []string{ins1, ins2, ins3} but got %v",
			schema.Instructions.Values,
		)
	}
	if schema.Name != "name" {
		t.Errorf("wanted name 'name' but got %q", schema.Name)
	}
	if schema.NutritionSchema != r.Nutrition.Schema("4") {
		t.Errorf("wanted nutrition but got %v", schema.NutritionSchema)
	}
	if schema.PrepTime != "PT1H0M0S" {
		t.Errorf("wanted prepTime PT1H0M0S but got %q", schema.PrepTime)
	}
	if !slices.Equal(schema.Tools.Values, []string{"t1", "t2", "t3"}) {
		t.Errorf("wanted tools []string{t1, t2, t3} but got %v", schema.Tools.Values)
	}
	if schema.Yield.Value != 4 {
		t.Errorf("wanted yield 4 but got %d", schema.Yield.Value)
	}
	if schema.URL != "https://www.google.com" {
		t.Errorf("wanted url https://www.google.com but got %q", schema.URL)
	}
}

func TestRecipe_Normalize(t *testing.T) {
	r := models.Recipe{
		Description: "Place the chicken pieces on a baking sheet and bake 1l 1 l 1ml 1 ml until they 425°f (220°c) and golden.",
		Ingredients: models.Ingredients{
			Values: []string{"ing1 1L", "1 L ing2", "ing3 of 1mL stuff", "ing4 of stuff 1 mL"},
		},
		Instructions: []string{
			"ins1 1l",
			"1 l ins2",
			"ins3 of 1ml stuff",
			"ins4 of stuff 1 ml",
		},
	}

	r.Normalize()

	expectedDescription := "Place the chicken pieces on a baking sheet and bake 1L 1 L 1mL 1 mL until they 425°F (220°C) and golden."
	if r.Description != expectedDescription {
		t.Errorf("expected the description to be normalized")
	}

	expected := []string{
		"ing1 1L",
		"1 L ing2",
		"ing3 of 1mL stuff",
		"ing4 of stuff 1 mL",
	}
	for i, v := range r.Ingredients.Values {
		if v != expected[i] {
			t.Errorf("expected '%v' but got '%v'", expected[i], v)
		}
	}

	expected = []string{"ins1 1L", "1 L ins2", "ins3 of 1mL stuff", "ins4 of stuff 1 mL"}
	for i, v := range r.Instructions {
		if v != expected[i] {
			t.Errorf("expected '%v' but got '%v'", expected[i], v)
		}
	}
}

func TestNewTimes(t *testing.T) {
	actual, err := models.NewTimes("PT1H0M0S", "PT2H0M0S")
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
}
