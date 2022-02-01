package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/utils/duration"
)

// Recipe holds information on a recipe.
type Recipe struct {
	ID           int64
	Name         string
	Description  string
	Image        uuid.UUID
	Url          string
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

// ToArgs adds every field related to a Recipe to an interface{} slice.
func (r Recipe) ToArgs(includeID bool) []interface{} {
	args := []interface{}{}
	if includeID {
		args = append(args, r.ID)
	}

	args = append(args, []interface{}{
		r.Name,
		r.Description,
		r.Image,
		r.Url,
		r.Yield,
		r.Category,
		r.Nutrition.Calories,
		r.Nutrition.TotalCarbohydrates,
		r.Nutrition.Sugars,
		r.Nutrition.Protein,
		r.Nutrition.TotalFat,
		r.Nutrition.SaturatedFat,
		r.Nutrition.Cholesterol,
		r.Nutrition.Sodium,
		r.Nutrition.Fiber,
		r.Times.Prep.String(),
		r.Times.Cook.String(),
	}...)

	arrs := [][]string{r.Ingredients, r.Instructions, r.Keywords, r.Tools}
	for _, arr := range arrs {
		for _, v := range arr {
			args = append(args, v)
		}
	}
	return args
}

// Times holds a variety of intervals.
type Times struct {
	Prep  time.Duration
	Cook  time.Duration
	Total time.Duration
}

// NewTimes creates a Times struct from the Schema Duration fields for prep and cook time.
func NewTimes(prep, cook string) (Times, error) {
	p, err := parseDuration(prep)
	if err != nil {
		return Times{}, err
	}

	c, err := parseDuration(cook)
	if err != nil {
		return Times{}, err
	}

	return Times{Prep: p, Cook: c}, nil
}

func parseDuration(d string) (time.Duration, error) {
	parts := strings.SplitN(d, ":", 3)
	if len(parts) == 3 {
		return time.ParseDuration(parts[0] + "h" + parts[1] + "m" + parts[2] + "s")
	}

	p, err := duration.Parse(d)
	return p.ToTimeDuration(), err
}

// Nutrition holds nutrition facts.
type Nutrition struct {
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
