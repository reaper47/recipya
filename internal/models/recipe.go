package models

import (
	"bytes"
	"errors"
	"github.com/donna-legal/word2number"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/units"
	"github.com/reaper47/recipya/internal/utils/duration"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"github.com/reaper47/recipya/internal/utils/regex"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	wordConverter *word2number.Converter
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
		system := units.DetectMeasurementSystem(s)
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

	recipe := r.Copy()
	recipe.Description = units.ConvertParagraph(r.Description, currentSystem, to)
	recipe.Ingredients = ingredients
	recipe.Instructions = instructions
	return &recipe, nil
}

// Copy deep copies the Recipe.
func (r *Recipe) Copy() Recipe {
	ingredients := make([]string, len(r.Ingredients))
	copy(ingredients, r.Ingredients)

	instructions := make([]string, len(r.Instructions))
	copy(instructions, r.Instructions)

	keywords := make([]string, len(r.Keywords))
	copy(keywords, r.Keywords)

	tools := make([]string, len(r.Tools))
	copy(tools, r.Tools)

	return Recipe{
		Category:     r.Category,
		CreatedAt:    r.CreatedAt,
		Cuisine:      r.Cuisine,
		Description:  r.Description,
		ID:           r.ID,
		Image:        r.Image,
		Ingredients:  ingredients,
		Instructions: instructions,
		Keywords:     keywords,
		Name:         r.Name,
		Nutrition: Nutrition{
			Calories:           r.Nutrition.Calories,
			Cholesterol:        r.Nutrition.Cholesterol,
			Fiber:              r.Nutrition.Fiber,
			Protein:            r.Nutrition.Protein,
			SaturatedFat:       r.Nutrition.SaturatedFat,
			Sodium:             r.Nutrition.Sodium,
			Sugars:             r.Nutrition.Sugars,
			TotalCarbohydrates: r.Nutrition.TotalCarbohydrates,
			TotalFat:           r.Nutrition.TotalFat,
			UnsaturatedFat:     r.Nutrition.UnsaturatedFat,
		},
		Times: Times{
			Prep:  r.Times.Prep,
			Cook:  r.Times.Cook,
			Total: r.Times.Total,
		},
		Tools:     tools,
		UpdatedAt: r.UpdatedAt,
		URL:       r.URL,
		Yield:     r.Yield,
	}
}

// IsEmpty verifies whether all the Recipe fields are empty.
func (r *Recipe) IsEmpty() bool {
	return r.Category == "" && r.CreatedAt.Equal(time.Time{}) && r.Cuisine == "" && r.Description == "" &&
		r.ID == 0 && r.Image == uuid.Nil && len(r.Ingredients) == 0 && len(r.Instructions) == 0 &&
		len(r.Keywords) == 0 && r.Name == "" && r.Nutrition.Equal(Nutrition{}) &&
		r.Times.Equal(Times{}) && len(r.Tools) == 0 && r.UpdatedAt.Equal(time.Time{}) &&
		r.URL == "" && r.Yield == 0
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

// Scale scales the recipe to the given yield.
func (r *Recipe) Scale(yield int16) {
	multiplier := float64(yield) / float64(r.Yield)
	scaledIngredients := make([]string, len(r.Ingredients))

	var wg sync.WaitGroup
	wg.Add(len(r.Ingredients))

	for i, ingredient := range r.Ingredients {
		go func(ing string, i int) {
			defer wg.Done()
			ing = units.ReplaceVulgarFractions(ing)
			system := units.DetectMeasurementSystem(ing)

			switch system {
			case units.MetricSystem, units.ImperialSystem:
				m, err := units.NewMeasurementFromString(ing)
				if err != nil {
					return
				}
				m = m.Scale(multiplier)
				scaled := regex.Unit.ReplaceAllString(ing, m.String())
				scaledIngredients[i] = units.ReplaceDecimalFractions(scaled)
			case units.InvalidSystem:
				if regex.BeginsWithWord.MatchString(ing) {
					ing = regex.BeginsWithWord.ReplaceAllStringFunc(ing, func(s string) string {
						f := wordConverter.Words2Number(s) * multiplier
						if f > 0 {
							return strconv.FormatFloat(f, 'g', 2, 64) + " "
						}
						return s
					})
				} else {
					ing = extensions.ScaleString(ing, multiplier)
					ing = units.ReplaceDecimalFractions(ing)
				}
				scaledIngredients[i] = ing
			}
		}(ingredient, i)
	}
	wg.Wait()

	r.Ingredients = scaledIngredients
	r.Yield = yield
	r.Normalize()
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

// Equal verifies whether the Times is equal to the other Times.
func (t Times) Equal(other Times) bool {
	return t.Prep == other.Prep && t.Cook == other.Cook && t.Total == other.Total
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
func (n *Nutrition) Schema(servings string) NutritionSchema {
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

// NutrientsFDC is a type alias for a slice of NutrientFDC.
type NutrientsFDC []NutrientFDC

// NutritionFact calculates the nutrition facts from the nutrients.
// The values in the nutrition table are per 100 grams. The weight is
// the sum of all ingredient quantities in grams.
func (n NutrientsFDC) NutritionFact(weight float64) Nutrition {
	var (
		calories       float64
		carbs          float64
		cholesterol    float64
		fiber          float64
		protein        float64
		saturatedFat   float64
		sodium         float64
		sugars         float64
		totalFat       float64
		unsaturatedFat float64
	)

	for _, nutrient := range n {
		v := nutrient.Value()

		switch nutrient.Name {
		case "Carbohydrates":
			carbs += v
		case "Cholesterol":
			cholesterol += v
		case "Energy":
			if nutrient.UnitName == "KCAL" {
				calories += v
			}
		case "Fatty acids, total monounsaturated", "Fatty acids, total polyunsaturated":
			unsaturatedFat += v
			totalFat += v
		case "Fatty acids, total trans":
			totalFat += v
		case "Fatty acids, total saturated":
			saturatedFat += v
			totalFat += v
		case "Fiber, total dietary":
			fiber += v
		case "Protein":
			protein += v
		case "Sodium, Na":
			sodium += v
		case "Sugars, total including NLEA":
			sugars += v
		}
	}

	weight *= 1e-2
	return Nutrition{
		Calories:           strconv.FormatFloat(calories/weight, 'f', 0, 64) + " kcal",
		Cholesterol:        formatNutrient(cholesterol / weight),
		Fiber:              formatNutrient(fiber / weight),
		Protein:            formatNutrient(protein / weight),
		SaturatedFat:       formatNutrient(saturatedFat / weight),
		Sodium:             formatNutrient(sodium / weight),
		Sugars:             formatNutrient(sugars / weight),
		TotalCarbohydrates: formatNutrient(carbs / weight),
		TotalFat:           formatNutrient(totalFat / weight),
		UnsaturatedFat:     formatNutrient(unsaturatedFat / weight),
	}
}

func formatNutrient(value float64) string {
	if value == 0 {
		return "-"
	}

	unit := "g"

	if value < 1 {
		value *= 1e3
		unit = "mg"
	}

	if value < 1 {
		value *= 1e3
		unit = "ug"
	}

	if value > 1e3 {
		value *= 1e-3
		unit = "kg"
	}

	return strconv.FormatFloat(value, 'f', 2, 64) + " " + unit
}

// NutrientFDC holds the nutrient data from the FDC database.
type NutrientFDC struct {
	ID        int64
	Name      string
	Amount    float64
	UnitName  string
	Reference units.Measurement
}

// Value calculates the amount of the nutrient scaled to the reference, in grams.
func (n NutrientFDC) Value() float64 {
	perGram := n.Amount * 1e-2
	switch n.UnitName {
	case "UG":
		perGram *= 1e-6
	case "MG":
		perGram *= 1e-3
	case "G":
		break
	case "KG":
		perGram *= 1e3
	}

	m, err := n.Reference.Convert(units.Gram)
	if err != nil {
		m, err = n.Reference.Convert(units.Millilitre)
		if err != nil {
			return 0
		}
	}
	return m.Quantity * perGram
}

// SearchOptionsRecipes defines the options for searching recipes.
type SearchOptionsRecipes struct {
	ByName     bool
	FullSearch bool
}

// NewRecipesFromMasterCook extracts the recipes from a MasterCook file.
func NewRecipesFromMasterCook(r io.Reader) Recipes {
	var recipes Recipes

	all, err := io.ReadAll(r)
	if err != nil {
		return nil
	}

	baseRecipe := Recipe{URL: "Imported from MasterCook"}

	var (
		isCategory     bool
		isIngredients  bool
		isInstructions bool
		isStartRecipe  bool
	)

	recipe := baseRecipe
	for _, line := range bytes.Split(all, []byte("\n")) {
		if !isStartRecipe && bytes.Contains(line, []byte("*  Exported from  MasterCook  *")) {
			if recipe.Yield == 0 {
				recipe.Yield = 4
			}

			if len(recipe.Ingredients) > 0 && recipe.Ingredients[0] != "***Information***" {
				recipes = append(recipes, recipe)
			}

			recipe = baseRecipe
			isStartRecipe = true
			isInstructions = false
			continue
		}

		if isStartRecipe && len(line) > 1 {
			recipe.Name = string(bytes.Join(bytes.Fields(line), []byte(" ")))
			isStartRecipe = false
			continue
		}

		if recipe.Yield == 0 && bytes.HasPrefix(line, []byte("Serving Size")) {
			before, prep, _ := bytes.Cut(line, []byte("   "))

			_, after, _ := bytes.Cut(before, []byte(":"))
			yield, err := strconv.ParseInt(string(bytes.TrimSpace(after)), 10, 16)
			if err == nil {
				recipe.Yield = int16(yield)
			}

			_, prep, ok := bytes.Cut(prep, []byte(":"))
			if ok {
				h, m, _ := bytes.Cut(prep, []byte(":"))
				times, err := NewTimes("PT"+string(h)+"H"+string(m), "PT0H0M")
				if err == nil {
					recipe.Times = times
				}
			}
			continue
		}

		if !isCategory && recipe.Category == "" && bytes.HasPrefix(line, []byte("Categories")) {
			isCategory = true
			_, after, _ := bytes.Cut(line, []byte(":"))
			before, after, _ := bytes.Cut(after, []byte("  "))
			recipe.Category = string(bytes.TrimSpace(before))
			recipe.Keywords = append(recipe.Keywords, recipe.Category, strings.TrimSpace(string(after)))
			continue
		}

		isIngredientsStart := bytes.HasPrefix(line, []byte("  Amount"))
		if isCategory && len(line) > 0 && !isIngredientsStart {
			split := bytes.Split(line, []byte("   "))
			for _, b := range split {
				s := string(bytes.TrimSpace(b))
				if s != "" {
					recipe.Keywords = append(recipe.Keywords, s)
				}
			}
		}

		if !isIngredients && isIngredientsStart {
			isCategory = false
			isIngredients = true
			continue
		}

		if isIngredients && len(line) > 1 && !bytes.HasPrefix(line, []byte("-")) {
			recipe.Ingredients = append(recipe.Ingredients, strings.Join(strings.Fields(string(line)), " "))
			continue
		}

		if isIngredients && len(line) < 2 {
			isIngredients = false
			isInstructions = true
			continue
		}

		isEnd := bytes.Contains(line, []byte("- - - - - - -"))
		if isInstructions && len(line) > 0 && !isEnd {
			s := string(bytes.TrimSpace(line))
			before, after, found := strings.Cut(s, ".")
			if found {
				_, err := strconv.Atoi(before)
				if err == nil {
					s = strings.TrimSpace(after)
				}
			}
			if s != "" {
				recipe.Instructions = append(recipe.Instructions, s)
			}
		}

		if isEnd {
			isInstructions = false
			isStartRecipe = false
		}
	}

	if recipe.Yield == 0 {
		recipe.Yield = 4
	}

	recipes = append(recipes, recipe)
	return recipes
}

func init() {
	c, err := word2number.NewConverter("en")
	if err != nil {
		panic(err)
	}
	wordConverter = c
}
