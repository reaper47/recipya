package models

// ShareRecipe represents a struct that stores the ID of a recipe and the ID of the user who shared it.
type ShareRecipe struct {
	RecipeID int64
	UserID   int64
}
