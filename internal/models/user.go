package models

import "github.com/reaper47/recipya/internal/units"

// User holds data related to a user.
type User struct {
	ID    int64
	Email string
}

// UserSettings holds the user's settings.
type UserSettings struct {
	CalculateNutritionFact bool
	ConvertAutomatically   bool
	CookbooksViewMode      ViewMode
	MeasurementSystem      units.System
}

// IsCalculateNutrition verifies whether the nutrition facts should be calculated for the recipe.
func (u *UserSettings) IsCalculateNutrition(recipe *Recipe) bool {
	return u.CalculateNutritionFact && recipe.Nutrition.Equal(Nutrition{})
}
