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

	/*t.Run("Recipe ToSchema transforms the Recipe to its schema successfully", func(t *testing.T) {
		t.Fail()
	})

	t.Run("Creating a new Times parses correctly", func(t *testing.T) {
		t.Fail()
	})

	t.Run("formatDuration formats correctly", func(t *testing.T) {
		t.Fail()
	})*/
}
