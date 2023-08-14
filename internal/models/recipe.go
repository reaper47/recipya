package models

import (
	"errors"
	"github.com/reaper47/recipya/internal/units"
	"github.com/reaper47/recipya/internal/utils/duration"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Recipes is the type for a slice of recipes.
type Recipes []Recipe

// Categories returns the category of every recipe. The resulting slice has no duplicates.
func (r Recipes) Categories() []string {
	xs := make([]string, len(r))
	for i, recipe := range r {
		xs[i] = recipe.Category
	}
	return extensions.Unique(xs)
}

// Recipe is the struct that holds a recipe's information.
type Recipe struct {
	Category     string
	CreatedAt    time.Time
	Cuisine      string
	Description  string
	ID           int64
	Image        uuid.UUID
	Ingredients  []string
	Instructions []string
	Keywords     []string
	Name         string
	Nutrition    Nutrition
	Times        Times
	Tools        []string
	UpdatedAt    time.Time
	URL          string
	Yield        int16
}

// ConvertMeasurementSystem converts a recipe to another units.System.
func (r *Recipe) ConvertMeasurementSystem(to units.System) (*Recipe, error) {
	currentSystem := units.InvalidSystem
	for _, s := range r.Ingredients {
		system := units.DetectMeasurementSystemFromSentence(s)
		if system != units.InvalidSystem {
			currentSystem = system
			break
		}
	}

	if currentSystem == units.InvalidSystem {
		return nil, errors.New("could not determine measurement system")
	} else if currentSystem == to {
		return nil, errors.New("system already " + to.String())
	}

	ingredients := make([]string, len(r.Ingredients))
	for i, s := range r.Ingredients {
		v, err := units.ConvertSentence(s, currentSystem, to)
		if err != nil {
			ingredients[i] = s
			continue
		}
		ingredients[i] = v
	}

	instructions := make([]string, len(r.Instructions))
	for i, s := range r.Instructions {
		v, err := units.ConvertSentence(s, currentSystem, to)
		if err != nil {
			instructions[i] = s
			continue
		}
		instructions[i] = v
	}

	return &Recipe{
		Category:     r.Category,
		CreatedAt:    r.CreatedAt,
		Cuisine:      r.Cuisine,
		Description:  units.ConvertParagraph(r.Description, currentSystem, to),
		ID:           r.ID,
		Image:        r.Image,
		Ingredients:  ingredients,
		Instructions: instructions,
		Keywords:     r.Keywords,
		Name:         r.Name,
		Nutrition:    r.Nutrition,
		Times:        r.Times,
		Tools:        r.Tools,
		UpdatedAt:    r.UpdatedAt,
		URL:          r.URL,
		Yield:        r.Yield,
	}, nil
}

// Normalize normalizes texts for readability.
// It normalizes quantities, i.e. 1l -> 1L and 1 ml -> 1 mL.
func (r *Recipe) Normalize() {
	r.Description = regex.Quantity.ReplaceAllStringFunc(r.Description, normalizeQuantity)

	for i, v := range r.Ingredients {
		r.Ingredients[i] = regex.Quantity.ReplaceAllStringFunc(v, normalizeQuantity)
	}

	for i, v := range r.Instructions {
		r.Instructions[i] = regex.Quantity.ReplaceAllStringFunc(v, normalizeQuantity)
	}
}

func normalizeQuantity(s string) string {
	xr := []rune(s)
	for i, v := range xr {
		switch v {
		case 'l':
			xr[i] = 'L'
		case 'f':
			xr[i] = 'F'
		case 'c':
			xr[i] = 'C'
		}
	}
	return string(xr)
}

// Schema creates the schema representation of the Recipe.
func (r *Recipe) Schema() RecipeSchema {
	return RecipeSchema{
		AtContext:       "https://schema.org",
		AtType:          SchemaType{Value: "Recipe"},
		Category:        Category{Value: r.Category},
		CookTime:        formatDuration(r.Times.Cook),
		Cuisine:         Cuisine{Value: r.Cuisine},
		DateCreated:     r.CreatedAt.Format(time.DateOnly),
		DateModified:    r.UpdatedAt.Format(time.DateOnly),
		DatePublished:   r.CreatedAt.Format(time.DateOnly),
		Description:     Description{Value: r.Description},
		Keywords:        Keywords{Values: strings.Join(r.Keywords, ",")},
		Image:           Image{Value: r.Image.String()},
		Ingredients:     Ingredients{Values: r.Ingredients},
		Instructions:    Instructions{Values: r.Instructions},
		Name:            r.Name,
		NutritionSchema: r.Nutrition.Schema(strconv.Itoa(int(r.Yield))),
		PrepTime:        formatDuration(r.Times.Prep),
		Tools:           Tools{Values: r.Tools},
		Yield:           Yield{Value: r.Yield},
		URL:             r.URL,
	}
}

// Times holds a variety of intervals.
type Times struct {
	Prep  time.Duration
	Cook  time.Duration
	Total time.Duration
}

// NewTimes creates a struct of Times from the Schema Duration fields for prep and cook time.
func NewTimes(prep, cook string) (Times, error) {
	p, err := parseDuration(prep)
	if err != nil {
		return Times{}, err
	}

	c, err := parseDuration(cook)
	if err != nil {
		return Times{}, err
	}

	return Times{Prep: p, Cook: c, Total: p + c}, nil
}

func parseDuration(d string) (time.Duration, error) {
	parts := strings.SplitN(d, ":", 3)
	if len(parts) == 3 {
		return time.ParseDuration(parts[0] + "h" + parts[1] + "m" + parts[2] + "s")
	}

	p, err := duration.Parse(d)
	if err != nil {
		return time.Duration(-1), err
	}
	return p.ToTimeDuration(), nil
}

func formatDuration(d time.Duration) string {
	return "PT" + strings.ToUpper(d.Truncate(time.Millisecond).String())
}

// Nutrition holds nutrition facts.
type Nutrition struct {
	Calories           string
	Cholesterol        string
	Fiber              string
	Protein            string
	SaturatedFat       string
	Sodium             string
	Sugars             string
	TotalCarbohydrates string
	TotalFat           string
	UnsaturatedFat     string
}

// Equal verifies whether the Nutrition struct is equal to the other.
func (n *Nutrition) Equal(other Nutrition) bool {
	return n.Calories == other.Calories &&
		n.Cholesterol == other.Cholesterol &&
		n.Fiber == other.Fiber &&
		n.Protein == other.Protein &&
		n.SaturatedFat == other.SaturatedFat &&
		n.Sodium == other.Sodium &&
		n.Sugars == other.Sugars &&
		n.TotalCarbohydrates == other.TotalCarbohydrates &&
		n.TotalFat == other.TotalFat &&
		n.UnsaturatedFat == other.UnsaturatedFat
}

// Schema creates the schema representation of the Nutrition.
func (n Nutrition) Schema(servings string) NutritionSchema {
	return NutritionSchema{
		Calories:       n.Calories,
		Carbohydrates:  n.TotalCarbohydrates,
		Cholesterol:    n.Cholesterol,
		Fat:            n.TotalFat,
		Fiber:          n.Fiber,
		Protein:        n.Protein,
		SaturatedFat:   n.SaturatedFat,
		Servings:       servings,
		Sodium:         n.Sodium,
		Sugar:          n.Sugars,
		UnsaturatedFat: n.UnsaturatedFat,
	}
}
