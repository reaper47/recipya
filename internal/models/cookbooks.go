package models

import "github.com/google/uuid"

// Cookbook is the struct that holds information on a cookbook.
type Cookbook struct {
	ID      int64
	Count   int
	Image   uuid.UUID
	Recipes Recipes
	Title   string
}
