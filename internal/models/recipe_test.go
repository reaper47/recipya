package models_test

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/units"
	"slices"
	"testing"
	"time"
)

func BenchmarkRecipe_ConvertMeasurementSystem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := models.Recipe{
			Description: "Preheat the oven to 177 degrees C (175 degrees C). " +
				"Stir in flour, chocolate chips, and walnuts. " +
				"Drop spoonfuls of dough 30mm apart onto ungreased baking sheets. " +
				"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
			Ingredients: []string{
				"236.58 mL butter, softened",
				"2 eggs",
				"10 mL vanilla extract",
				"5 mL baking soda",
				"709.76 mL all-purpose flour",
				"473.18 mL semisweet chocolate chips",
			},
			Instructions: []string{
				"Preheat the oven to 177 degrees C (175 degrees C).",
				"Stir in flour, chocolate chips, and walnuts.",
				"Drop spoonfuls of dough 30mm apart onto ungreased baking sheets.",
				"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
			},
		}

		converted, err := r.ConvertMeasurementSystem(units.MetricSystem)
		_ = converted
		_ = err
	}
}

func TestRecipe_ConvertMeasurementSystem(t *testing.T) {
	testcases := []struct {
		name string
		in   *models.Recipe
		to   units.System
		want error
	}{
		{
			name: "imperial to imperial",
			in: &models.Recipe{
				Ingredients: []string{
					"2 eggs",
					"1 cup butter, softened",
				},
			},
			to:   units.ImperialSystem,
			want: errors.New("system already imperial"),
		},
		{
			name: "metric to metric",
			in: &models.Recipe{
				Ingredients: []string{
					"2 eggs",
					"1.5 mL butter, softened",
				},
			},
			to:   units.MetricSystem,
			want: errors.New("system already metric"),
		},
	}
	for _, tc := range testcases {
		t.Run("cannot convert "+tc.name, func(t *testing.T) {
			_, err := tc.in.ConvertMeasurementSystem(tc.to)
			if err == nil {
				t.Fatalf("expected error but got %q", err)
			}
		})
	}

	testcases2 := []struct {
		name string
		in   models.Recipe
		to   units.System
		want models.Recipe
	}{
		{
			name: "imperial to metric",
			in: models.Recipe{
				Description: "Preheat the oven to 351 °F (351 °F). " +
					"Stir in flour, chocolate chips, and walnuts. " +
					"Drop spoonfuls of dough 1.18 inches apart onto ungreased baking sheets. " +
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				Ingredients: []string{
					"1 cup butter, softened",
					"2 eggs",
					"2 teaspoons vanilla extract",
					"1 teaspoon baking soda",
					"3 cups all-purpose flour",
					"2 cups semisweet chocolate chips",
					"Salt and pepper",
				},
				Instructions: []string{
					"Preheat the oven to 350 degrees F (175 degrees C).",
					"Stir in flour, chocolate chips, and walnuts.",
					"Drop spoonfuls of dough 2 inches apart onto ungreased baking sheets.",
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				},
			},
			to: units.MetricSystem,
			want: models.Recipe{
				Description: "Preheat the oven to 177 °C (177 °C). " +
					"Stir in flour, chocolate chips, and walnuts. " +
					"Drop spoonfuls of dough 3 cm apart onto ungreased baking sheets. " +
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				Ingredients: []string{
					"2.37 dl butter, softened",
					"2 eggs",
					"10 ml vanilla extract",
					"5 ml baking soda",
					"7.1 dl all-purpose flour",
					"4.73 dl semisweet chocolate chips",
					"Salt and pepper",
				},
				Instructions: []string{
					"Preheat the oven to 177 °C (177 °C).",
					"Stir in flour, chocolate chips, and walnuts.",
					"Drop spoonfuls of dough 5.08 cm apart onto ungreased baking sheets.",
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				},
			},
		},
		{
			name: "metric to imperial",
			in: models.Recipe{
				Description: "Preheat the oven to 177 degrees C (175 degrees C). " +
					"Stir in flour, chocolate chips, and walnuts. " +
					"Drop spoonfuls of dough 30mm apart onto ungreased baking sheets. " +
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				Ingredients: []string{
					"236.58 mL butter, softened",
					"2 eggs",
					"10 mL vanilla extract",
					"5 mL baking soda",
					"709.76 mL all-purpose flour",
					"473.18 mL semisweet chocolate chips",
				},
				Instructions: []string{
					"Preheat the oven to 177 degrees C (175 degrees C).",
					"Stir in flour, chocolate chips, and walnuts.",
					"Drop spoonfuls of dough 30mm apart onto ungreased baking sheets.",
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				},
			},
			to: units.ImperialSystem,
			want: models.Recipe{
				Description: "Preheat the oven to 351 °F (351 °F). " +
					"Stir in flour, chocolate chips, and walnuts. " +
					"Drop spoonfuls of dough 1.18 inches apart onto ungreased baking sheets. " +
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				Ingredients: []string{
					"1 cup butter, softened",
					"2 eggs",
					"2 tsp vanilla extract",
					"1 tsp baking soda",
					"3 cups all-purpose flour",
					"2 cups semisweet chocolate chips",
				},
				Instructions: []string{
					"Preheat the oven to 351 °F (351 °F).",
					"Stir in flour, chocolate chips, and walnuts.",
					"Drop spoonfuls of dough 1.18 inches apart onto ungreased baking sheets.",
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				},
			},
		},
	}
	for _, tc := range testcases2 {
		t.Run("valid "+tc.name, func(t *testing.T) {
			got, _ := tc.in.ConvertMeasurementSystem(tc.to)

			if got.Description != tc.want.Description {
				t.Fatalf("got description:\n%s\nbut want:\n%s", got.Description, tc.want.Description)
			}

			if len(got.Ingredients) != len(tc.want.Ingredients) {
				t.Fatalf("ingredients of different lengths: %#v but want %#v", got.Ingredients, tc.want.Ingredients)
			}
			if len(got.Instructions) != len(tc.want.Instructions) {
				t.Fatalf("instructions of different lengths: %#v but want %#v", got.Ingredients, tc.want.Ingredients)
			}
			for i, s := range got.Ingredients {
				if s != tc.want.Ingredients[i] {
					t.Errorf("got ingredient %q but want %q", s, tc.want.Ingredients[i])
				}
			}
			for i, s := range got.Instructions {
				if s != tc.want.Instructions[i] {
					t.Errorf("got instruction %q but want %q", s, tc.want.Instructions[i])
				}
			}
		})
	}
}

func TestRecipe_Copy(t *testing.T) {
	times, _ := models.NewTimes("PT1H20M", "PT30M")
	original := models.Recipe{
		Category:    "breakfast",
		CreatedAt:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Cuisine:     "American",
		Description: "The best American breakfast you could ever have.",
		ID:          1,
		Image:       uuid.Nil,
		Ingredients: []string{
			"2 eggs",
			"3 cups maple syrup",
			"1/2 tbsp cinnamon",
			"1 toast spread with butter",
		},
		Instructions: []string{
			"Cook the meat on high",
			"Mix the eggs together, whisk and cook on low",
			"Pour maple syrup over meat",
			"Eat",
		},
		Keywords: []string{"eggs", "meat", "juicy"},
		Name:     "Chicken Sauce",
		Nutrition: models.Nutrition{
			Calories:           "1g",
			Cholesterol:        "2g",
			Fiber:              "3g",
			Protein:            "4g",
			SaturatedFat:       "5g",
			Sodium:             "6g",
			Sugars:             "7g",
			TotalCarbohydrates: "8g",
			TotalFat:           "9g",
			UnsaturatedFat:     "10g",
		},
		Times: times,
		Tools: []string{
			"small pan",
			"large pan",
			"spatula",
		},
		UpdatedAt: time.Date(2012, 12, 25, 0, 0, 0, 0, time.UTC),
		URL:       "https://www.example.com",
		Yield:     4,
	}

	copied := original.Copy()

	if !cmp.Equal(original, copied) {
		t.Log(cmp.Diff(original, copied))
		t.Fail()
	}
	copied.Ingredients[0] = "pig"
	if slices.Equal(original.Ingredients, copied.Ingredients) {
		t.Fatal("ingredients slices the same when they should not")
	}
	copied.Instructions[0] = "jesus"
	if slices.Equal(original.Instructions, copied.Instructions) {
		t.Fatal("instructions slices the same when they should not")
	}
	copied.Keywords[0] = "european"
	if slices.Equal(original.Keywords, copied.Keywords) {
		t.Fatal("keywords slices the same when they should not")
	}
	copied.Tools[0] = "saw"
	if slices.Equal(original.Tools, copied.Tools) {
		t.Fatal("tools slices the same when they should not")
	}
}

func TestRecipe_Normalize(t *testing.T) {
	r := models.Recipe{
		Description: "Place the chicken pieces on a baking sheet and bake 1l 1 l 1ml 1 ml until they 425°f (220°c) and golden.",
		Ingredients: []string{"ing1 1L", "1 L ing2", "ing3 of 1mL stuff", "ing4 of stuff 1 mL"},
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
	for i, v := range r.Ingredients {
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

func TestRecipe_Scale(t *testing.T) {
	recipe := models.Recipe{
		ID: 1,
		Ingredients: []string{
			"2 big apples",
			"Lots of big apples",
			"2.5 slices of bacon",
			"2 1/3 cans of bamboo sticks",
			"1½ can of tomato paste",
			"6 ¾ peanut butter jars",
			"7.5mL of whiskey",
			"2 tsp lemon juice",
		},
		Instructions: nil,
		Name:         "Sauce",
		Yield:        4,
	}

	t.Run("double recipe", func(t *testing.T) {
		got, err := recipe.Scale(8)
		if err != nil {
			t.Fatal(err)
		}

		want := recipe.Copy()
		want.Ingredients = []string{
			"4 big apples",
			"Lots of big apples",
			"5 slices of bacon",
			"4 2/3 cans of bamboo sticks",
			"3 can of tomato paste",
			"13 1/2 peanut butter jars",
			"15mL of whiskey",
			"1 1/3 tbsp lemon juice",
		}
		want.Yield = 8
		if !cmp.Equal(got, want) {
			t.Log(cmp.Diff(got, want))
			t.Fail()
		}
	})

	t.Run("decrease recipe by 1.5x", func(t *testing.T) {
		got, err := recipe.Scale(8)
		if err != nil {
			t.Fatal(err)
		}

		want := models.Recipe{
			ID: 1,
			Ingredients: []string{
				"4 big apples",
			},
			Instructions: nil,
			Name:         "Sauce",
			Yield:        4,
		}
		if !cmp.Equal(got, want) {
			t.Log(cmp.Diff(got, want))
			t.Fail()
		}
	})
}

func TestRecipe_Schema(t *testing.T) {
	imageUUID := uuid.New()
	r := models.Recipe{
		Category:     "lunch",
		CreatedAt:    time.Now(),
		Cuisine:      "american",
		Description:  "description",
		ID:           1,
		Image:        imageUUID,
		Ingredients:  []string{"ing1", "ing2", "ing3"},
		Instructions: []string{"ins1", "ins2", "ins3"},
		Keywords:     []string{"kw1", "kw2", "kw3"},
		Name:         "name",
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
		Times: models.Times{
			Prep:  1 * time.Hour,
			Cook:  2 * time.Hour,
			Total: 3 * time.Hour,
		},
		Tools:     []string{"t1", "t2", "t3"},
		UpdatedAt: time.Now().Add(2 * time.Hour),
		URL:       "https://www.google.com",
		Yield:     4,
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
	if schema.Cuisine.Value != "american" {
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
