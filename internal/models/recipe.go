package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/regex"
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
	Ingredients  Ingredients
	Nutrition    Nutrition
	Instructions []string
	Tools        []string
	Keywords     []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ToArgs adds every field related to a Recipe to an any slice.
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

	arrs := [][]string{r.Ingredients.Values, r.Instructions, r.Keywords, r.Tools}
	for _, arr := range arrs {
		for _, v := range arr {
			args = append(args, v)
		}
	}
	return args
}

// ToSchema creates the schema representation of the Recipe.
func (r Recipe) ToSchema() RecipeSchema {
	return RecipeSchema{
		AtContext:       "http://schema.org",
		AtType:          SchemaType{Value: "Recipe"},
		Category:        Category{Value: r.Category},
		CookTime:        formatDuration(r.Times.Cook),
		Cuisine:         Cuisine{},
		DateCreated:     r.CreatedAt.Format("2006-01-02"),
		DateModified:    r.UpdatedAt.Format("2006-01-02"),
		DatePublished:   r.CreatedAt.Format("2006-01-02"),
		Description:     Description{Value: r.Description},
		Keywords:        Keywords{Values: strings.Join(r.Keywords, ",")},
		Image:           Image{Value: string(r.Image.String())},
		Ingredients:     r.Ingredients,
		Instructions:    Instructions{Values: r.Instructions},
		Name:            r.Name,
		NutritionSchema: r.Nutrition.toSchema(strconv.Itoa(int(r.Yield))),
		PrepTime:        formatDuration(r.Times.Prep),
		Tools:           Tools{Values: r.Tools},
		Yield:           Yield{Value: r.Yield},
		Url:             r.Url,
	}
}

// Normalize normalizes texts for readability.
//
// It normalizes quanities, i.e. 1l -> 1L and 1 ml -> 1 mL.
func (r *Recipe) Normalize() {
	r.Description = regex.Quantity.ReplaceAllStringFunc(r.Description, normalizeQuantity)

	for i, v := range r.Ingredients.Values {
		r.Ingredients.Values[i] = regex.Quantity.ReplaceAllStringFunc(v, normalizeQuantity)
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

	return Times{Prep: p, Cook: c, Total: p + c}, nil
}

func parseDuration(d string) (time.Duration, error) {
	parts := strings.SplitN(d, ":", 3)
	if len(parts) == 3 {
		return time.ParseDuration(parts[0] + "h" + parts[1] + "m" + parts[2] + "s")
	}

	p, err := duration.Parse(d)
	return p.ToTimeDuration(), err
}

func formatDuration(d time.Duration) string {
	return "PT" + strings.ToUpper(d.Truncate(time.Millisecond).String())
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

func (m Nutrition) toSchema(servings string) NutritionSchema {
	return NutritionSchema{
		Calories:      m.Calories,
		Carbohydrates: m.TotalCarbohydrates,
		Cholesterol:   m.Cholesterol,
		Fat:           m.TotalFat,
		Fiber:         m.Fiber,
		Protein:       m.Protein,
		SaturatedFat:  m.SaturatedFat,
		Servings:      servings,
		Sodium:        m.Sodium,
		Sugar:         m.Sugars,
	}
}
