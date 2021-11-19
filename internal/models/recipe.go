package models

import (
	"time"

	"github.com/google/uuid"
)

// Recipe holds information on a recipe.
type Recipe struct {
	ID           int64
	Name         string
	Description  string
	Url          uuid.UUID
	Yield        int16
	Category     string
	Times        Times
	Ingredients  []string
	Nutrition    Nutrition
	Instructions []string
	Tools        []string
	Keywords     []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Times holds a variety of intervals.
type Times struct {
	Prep  time.Duration
	Cook  time.Duration
	Total time.Duration
}

// Nutrition holds nutrition facts.
type Nutrition struct {
	ID                 int64
	Calories           string
	TotalCarbohydrates string
	Sugars             string
	Protein            string
	TotalFat           string
	SaturatedFat       string
	Cholesterol        string
	Sodium             string
	Fiber              string
}
