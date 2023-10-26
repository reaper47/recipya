package models

// Share stores the ID of a recipe and the ID of the user who shared it.
type Share struct {
	CookbookID int64
	RecipeID   int64
	UserID     int64
}
