package models

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/donna-legal/word2number"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/units"
	"github.com/reaper47/recipya/internal/utils/duration"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"github.com/reaper47/recipya/internal/utils/regex"
)

var (
	wordConverter *word2number.Converter
)

// Errors to signal specific applications.
var (
	ErrIsAccuChef         = errors.New("accuchef") // TODO: Place errors somewhere else.
	ErrIsEasyRecipeDeluxe = errors.New("easy recipe deluxe")
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

// NewBaseRecipe creates a new, empty Recipe.
func NewBaseRecipe() Recipe {
	return Recipe{
		Category:     "uncategorized",
		Images:       make([]uuid.UUID, 0),
		Ingredients:  make([]string, 0),
		Instructions: make([]string, 0),
		Keywords:     make([]string, 0),
		Tools:        make([]HowToItem, 0),
		Videos:       make([]VideoObject, 0),
		Yield:        1,
	}
}

// Recipe is the struct that holds a recipe's information.
type Recipe struct {
	Category     string
	CreatedAt    time.Time
	Cuisine      string
	Description  string
	ID           int64
	Images       []uuid.UUID
	Ingredients  []string
	Instructions []string
	Keywords     []string
	Name         string
	Nutrition    Nutrition
	Times        Times
	Tools        []HowToItem
	UpdatedAt    time.Time
	URL          string
	Videos       []VideoObject
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
		return r, errors.New("could not determine measurement system")
	} else if currentSystem == to {
		return r, errors.New("system already " + to.String())
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
		instructions[i] = units.ConvertParagraph(s, currentSystem, to)
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

	tools := make([]HowToItem, len(r.Tools))
	copy(tools, r.Tools)

	videos := make([]VideoObject, len(r.Videos))
	copy(videos, r.Videos)

	return Recipe{
		Category:     r.Category,
		CreatedAt:    r.CreatedAt,
		Cuisine:      r.Cuisine,
		Description:  r.Description,
		ID:           r.ID,
		Images:       r.Images,
		Ingredients:  ingredients,
		Instructions: instructions,
		Keywords:     keywords,
		Name:         r.Name,
		Nutrition: Nutrition{
			Calories:           r.Nutrition.Calories,
			Cholesterol:        r.Nutrition.Cholesterol,
			Fiber:              r.Nutrition.Fiber,
			Protein:            r.Nutrition.Protein,
			TotalFat:           r.Nutrition.TotalFat,
			SaturatedFat:       r.Nutrition.SaturatedFat,
			UnsaturatedFat:     r.Nutrition.UnsaturatedFat,
			TransFat:           r.Nutrition.TransFat,
			Sodium:             r.Nutrition.Sodium,
			Sugars:             r.Nutrition.Sugars,
			TotalCarbohydrates: r.Nutrition.TotalCarbohydrates,
		},
		Times: Times{
			Prep:  r.Times.Prep,
			Cook:  r.Times.Cook,
			Total: r.Times.Total,
		},
		Tools:     tools,
		UpdatedAt: r.UpdatedAt,
		URL:       r.URL,
		Videos:    videos,
		Yield:     r.Yield,
	}
}

// IsEmpty verifies whether all the Recipe fields are empty.
func (r *Recipe) IsEmpty() bool {
	return r.Category == "" && r.CreatedAt.Equal(time.Time{}) && r.Cuisine == "" && r.Description == "" &&
		r.ID == 0 && len(r.Images) == 0 && len(r.Ingredients) == 0 && len(r.Instructions) == 0 &&
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

	for i, ingredient := range scaledIngredients {
		scaledIngredients[i] = regex.Digit.ReplaceAllStringFunc(ingredient, func(s string) string {
			f, err := strconv.ParseFloat(s, 64)
			if err == nil {
				return strconv.FormatFloat(f, 'f', -1, 64)
			}
			return s
		})
	}

	r.Ingredients = scaledIngredients
	r.Yield = yield
	r.Normalize()
}

// Schema creates the schema representation of the Recipe.
func (r *Recipe) Schema() RecipeSchema {
	var thumbnail string
	images := make([]string, 0, len(r.Images))
	if len(r.Images) > 0 {
		thumbnail = app.Config.Address() + "/data/images/" + r.Images[0].String()

		for _, img := range r.Images {
			images = append(images, app.Config.Address()+"/data/images/"+img.String()+app.ImageExt)
		}
	}

	instructions := make([]HowToItem, 0, len(r.Instructions))
	for _, ins := range r.Instructions {
		instructions = append(instructions, NewHowToStep(ins))
	}

	video := &Videos{Values: make([]VideoObject, 0, len(r.Videos))}
	for i, v := range r.Videos {
		u := app.Config.Address() + "/data/videos/" + v.ID.String() + app.VideoExt
		if v.ContentURL == "" {
			v.ContentURL = u
		}

		video.Values = append(video.Values, VideoObject{
			AtType:       "VideoObject",
			Name:         "Video #" + strconv.Itoa(i+1),
			Description:  "A video showing how to cook " + r.Name,
			ID:           v.ID,
			ThumbnailURL: v.ThumbnailURL,
			ContentURL:   v.ContentURL,
			EmbedURL:     v.EmbedURL,
			UploadDate:   v.UploadDate,
			Duration:     v.Duration,
			Expires:      time.Now().AddDate(1000, 0, 0),
		})
	}

	schema := RecipeSchema{
		AtContext:       "https://schema.org",
		AtType:          &SchemaType{Value: "Recipe"},
		Category:        &Category{Value: r.Category},
		CookingMethod:   &CookingMethod{},
		CookTime:        formatDuration(r.Times.Cook),
		Cuisine:         &Cuisine{Value: r.Cuisine},
		DateCreated:     r.CreatedAt.Format(time.DateOnly),
		DateModified:    r.UpdatedAt.Format(time.DateOnly),
		DatePublished:   r.CreatedAt.Format(time.DateOnly),
		Description:     &Description{Value: r.Description},
		Keywords:        &Keywords{Values: strings.Join(r.Keywords, ",")},
		Image:           &Image{Value: strings.Join(images, ";")},
		Ingredients:     &Ingredients{Values: r.Ingredients},
		Instructions:    &Instructions{Values: instructions},
		Name:            r.Name,
		NutritionSchema: r.Nutrition.Schema(strconv.Itoa(int(r.Yield))),
		PrepTime:        formatDuration(r.Times.Prep),
		ThumbnailURL:    &ThumbnailURL{Value: thumbnail},
		Tools:           &Tools{Values: r.Tools},
		TotalTime:       formatDuration(r.Times.Total),
		Yield:           &Yield{Value: r.Yield},
		URL:             r.URL,
		Video:           video,
	}

	if schema.CookingMethod.Value == "" {
		schema.CookingMethod = nil
	}

	if schema.Cuisine.Value == "" {
		schema.Cuisine = nil
	}

	if schema.Description.Value == "" {
		schema.Description = nil
	}

	if len(schema.Keywords.Values) == 0 {
		schema.Keywords = nil
	}

	if schema.Image.Value == "" {
		schema.Image = nil
	}

	if len(schema.Tools.Values) == 0 {
		schema.Tools = nil
	}

	if len(schema.Ingredients.Values) == 0 {
		schema.Ingredients = nil
	}

	if len(schema.Instructions.Values) == 0 {
		schema.Instructions = nil
	}

	if schema.NutritionSchema.Equal(NutritionSchema{}) {
		schema.NutritionSchema = nil
	}

	if schema.ThumbnailURL.Value == "" {
		schema.ThumbnailURL = nil
	}

	if schema.Yield.Value == 0 {
		schema.Yield = nil
	}

	return schema
}

// Times holds a variety of intervals.
type Times struct {
	Prep  time.Duration
	Cook  time.Duration
	Total time.Duration
}

// Equal verifies whether the Times struct is equal to the other Times.
func (t Times) Equal(other Times) bool {
	return t.Prep == other.Prep && t.Cook == other.Cook && t.Total == other.Total
}

// NewTimes creates a struct of Times from the Schema Duration fields for prep and cook time.
func NewTimes(prep, cook string) (Times, error) {
	p, err := parseDuration(prep)
	if err != nil {
		slog.Error("Could not parse duration", "prep", prep, "error", err)
		p = 0
	}

	c, err := parseDuration(cook)
	if err != nil {
		slog.Error("Could not parse duration", "cook", cook, "error", err)
		c = 0
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
	if d == 0 {
		return ""
	}

	s := strings.ToUpper(d.Truncate(time.Millisecond).String())

	before, _, ok := strings.Cut(s, "M0S")
	if ok {
		s = before + "M"
	}

	before, _, ok = strings.Cut(s, "H0M")
	if ok {
		s = before + "H"
	}

	return "PT" + s
}

// Nutrition holds nutrition facts.
type Nutrition struct {
	Calories           string
	Cholesterol        string
	Fiber              string
	IsPerServing       bool
	Protein            string
	TotalFat           string
	SaturatedFat       string
	UnsaturatedFat     string
	TransFat           string
	Sodium             string
	Sugars             string
	TotalCarbohydrates string
}

// Clean empties any negligent field of the Nutrition.
func (n *Nutrition) Clean() {
	if strings.HasPrefix(n.Calories, "0") {
		n.Calories = ""
	}

	if strings.HasPrefix(n.Cholesterol, "0") {
		n.Cholesterol = ""
	}

	if strings.HasPrefix(n.Fiber, "0") {
		n.Fiber = ""
	}

	if strings.HasPrefix(n.Protein, "0") {
		n.Protein = ""
	}

	if strings.HasPrefix(n.TotalFat, "0") {
		n.TotalFat = ""
	}

	if strings.HasPrefix(n.SaturatedFat, "0") {
		n.SaturatedFat = ""
	}

	if strings.HasPrefix(n.UnsaturatedFat, "0") {
		n.UnsaturatedFat = ""
	}

	if strings.HasPrefix(n.TransFat, "0") {
		n.TransFat = ""
	}

	if strings.HasPrefix(n.Sodium, "0") {
		n.Sodium = ""
	}

	if strings.HasPrefix(n.Sugars, "0") {
		n.Sugars = ""
	}

	if strings.HasPrefix(n.TotalCarbohydrates, "0") {
		n.TotalCarbohydrates = ""
	}
}

// Equal verifies whether the Nutrition struct is equal to the other.
func (n *Nutrition) Equal(other Nutrition) bool {
	return n.Calories == other.Calories &&
		n.Cholesterol == other.Cholesterol &&
		n.Fiber == other.Fiber &&
		n.Protein == other.Protein &&
		n.TotalFat == other.TotalFat &&
		n.SaturatedFat == other.SaturatedFat &&
		n.UnsaturatedFat == other.UnsaturatedFat &&
		n.TransFat == other.TransFat &&
		n.Sodium == other.Sodium &&
		n.Sugars == other.Sugars &&
		n.TotalCarbohydrates == other.TotalCarbohydrates
}

// Format formats the nutrition.
func (n *Nutrition) Format() string {
	if n.Equal(Nutrition{}) {
		return ""
	}

	var sb strings.Builder
	if n.Calories != "" {
		sb.WriteString("calories ")
		sb.WriteString(EnsureNutritionUnitForString(n.Calories, "Calories"))
		sb.WriteString("; ")
	}

	if n.TotalCarbohydrates != "" {
		sb.WriteString("total carbohydrates ")
		sb.WriteString(EnsureNutritionUnitForString(n.TotalCarbohydrates, "Carbohydrates"))
		sb.WriteString("; ")
	}

	if n.Sugars != "" {
		sb.WriteString("sugar ")
		sb.WriteString(EnsureNutritionUnitForString(n.Sugars, "Sugar"))
		sb.WriteString("; ")
	}
	if n.Protein != "" {
		sb.WriteString("protein ")
		sb.WriteString(EnsureNutritionUnitForString(n.Protein, "Protein"))
		sb.WriteString("; ")
	}

	if n.TotalFat != "" {
		sb.WriteString("total fat ")
		sb.WriteString(EnsureNutritionUnitForString(n.TotalFat, "Fat"))
		sb.WriteString("; ")
	}

	if n.SaturatedFat != "" {
		sb.WriteString("saturated fat ")
		sb.WriteString(EnsureNutritionUnitForString(n.SaturatedFat, "SaturatedFat"))
		sb.WriteString("; ")
	}

	if n.UnsaturatedFat != "" {
		sb.WriteString("unsaturated fat ")
		sb.WriteString(EnsureNutritionUnitForString(n.UnsaturatedFat, "UnsaturatedFat"))
		sb.WriteString("; ")
	}

	if n.TransFat != "" {
		sb.WriteString("trans fat ")
		sb.WriteString(EnsureNutritionUnitForString(n.TransFat, "TransFat"))
		sb.WriteString("; ")
	}

	if n.Cholesterol != "" {
		sb.WriteString("cholesterol ")
		sb.WriteString(EnsureNutritionUnitForString(n.Cholesterol, "Cholesterol"))
		sb.WriteString("; ")
	}

	if n.Sodium != "" {
		sb.WriteString("sodium ")
		sb.WriteString(EnsureNutritionUnitForString(n.Sodium, "Sodium"))
		sb.WriteString("; ")
	}

	if n.Fiber != "" {
		sb.WriteString("fiber ")
		sb.WriteString(EnsureNutritionUnitForString(n.Fiber, "Fiber"))
	}

	return "Per 100g: " + sb.String()
}

// Scale scales the nutrition by the given multiplier.
func (n *Nutrition) Scale(multiplier float64) {
	scale := func(s string) string {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return s
		}
		return extensions.FloatToString(f*multiplier, "%.2f")
	}

	n.Calories = regex.Digit.ReplaceAllStringFunc(n.Calories, scale)
	n.Cholesterol = regex.Digit.ReplaceAllStringFunc(n.Cholesterol, scale)
	n.Fiber = regex.Digit.ReplaceAllStringFunc(n.Fiber, scale)
	n.Protein = regex.Digit.ReplaceAllStringFunc(n.Protein, scale)
	n.TotalFat = regex.Digit.ReplaceAllStringFunc(n.TotalFat, scale)
	n.SaturatedFat = regex.Digit.ReplaceAllStringFunc(n.SaturatedFat, scale)
	n.UnsaturatedFat = regex.Digit.ReplaceAllStringFunc(n.UnsaturatedFat, scale)
	n.TransFat = regex.Digit.ReplaceAllStringFunc(n.TransFat, scale)
	n.Sodium = regex.Digit.ReplaceAllStringFunc(n.Sodium, scale)
	n.Sugars = regex.Digit.ReplaceAllStringFunc(n.Sugars, scale)
	n.TotalCarbohydrates = regex.Digit.ReplaceAllStringFunc(n.TotalCarbohydrates, scale)
}

// Schema creates the schema representation of the Nutrition.
func (n *Nutrition) Schema(servings string) *NutritionSchema {
	return &NutritionSchema{
		Calories:       n.Calories,
		Carbohydrates:  n.TotalCarbohydrates,
		Cholesterol:    n.Cholesterol,
		Fat:            n.TotalFat,
		Fiber:          n.Fiber,
		SaturatedFat:   n.SaturatedFat,
		UnsaturatedFat: n.UnsaturatedFat,
		TransFat:       n.TransFat,
		Protein:        n.Protein,
		Servings:       servings,
		Sodium:         n.Sodium,
		Sugar:          n.Sugars,
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
		totalFat       float64
		saturatedFat   float64
		unsaturatedFat float64
		transFat       float64
		sodium         float64
		sugars         float64
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
			transFat += v
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
		TotalFat:           formatNutrient(totalFat / weight),
		SaturatedFat:       formatNutrient(saturatedFat / weight),
		UnsaturatedFat:     formatNutrient(unsaturatedFat / weight),
		TransFat:           formatNutrient(transFat / weight),
		Sodium:             formatNutrient(sodium / weight),
		Sugars:             formatNutrient(sugars / weight),
		TotalCarbohydrates: formatNutrient(carbs / weight),
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
	Advanced   AdvancedSearch
	CookbookID int64
	Query      string
	Page       uint64
	Sort       Sort
}

// Arg returns combines the field values of the struct, ready for FTS.
func (s *SearchOptionsRecipes) Arg() string {
	var args []string

	fields := map[string]string{
		"category":     s.Advanced.Category,
		"cuisine":      s.Advanced.Cuisine,
		"description":  s.Advanced.Description,
		"ingredients":  s.Advanced.Ingredients,
		"instructions": s.Advanced.Instructions,
		"keywords":     s.Advanced.Keywords,
		"name":         s.Advanced.Name,
		"source":       s.Advanced.Source,
		"tools":        s.Advanced.Tools,
	}
	for col, field := range fields {
		x := toArg(field, col)
		if x != "" {
			args = append(args, x)
		}
	}

	return strings.Join(slices.DeleteFunc(args, func(s string) bool {
		return s == ""
	}), " AND ")
}

func toArg(s, col string) string {
	parts := strings.Split(s, ",")
	if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
		return ""
	}

	sep := " OR "
	if col == "ingredients" || col == "instructions" || col == "keywords" {
		sep = " AND "
	}

	var cat strings.Builder
	cat.WriteString("(")
	if col == "description" || col == "instructions" {
		for i, part := range parts {
			if i > 0 {
				cat.WriteString(sep)
			}

			cat.WriteString(col + ":NEAR(")
			subParts := strings.Split(part, " ")
			for i, subPart := range subParts {
				if i > 0 {
					cat.WriteString(" ")
				}
				cat.WriteString(`"` + subPart + `"`)
			}
			cat.WriteString(")")
		}
	} else {
		for i, part := range parts {
			if i > 0 {
				cat.WriteString(sep)
			}

			subcats := strings.Split(strings.TrimSpace(part), ":")
			if len(subcats) > 0 {
				cat.WriteString(col + `:"`)
				cat.WriteString(strings.Join(subcats, "*+"))
				cat.WriteString(`*"`)
			} else {
				cat.WriteString(col + `:"` + strings.TrimSpace(part) + `*"`)
			}
		}
	}

	cat.WriteString(")")
	return cat.String()
}

// IsBasic verifies whether the search is basic.
func (s *SearchOptionsRecipes) IsBasic() bool {
	return s.Advanced.Category == "" && s.Advanced.Cuisine == "" && s.Advanced.Description == "" &&
		s.Advanced.Ingredients == "" && s.Advanced.Instructions == "" && s.Advanced.Keywords == "" && s.Advanced.Name == "" &&
		s.Advanced.Source == "" && s.Advanced.Tools == ""
}

// AdvancedSearch stores the components of an advanced search query.
type AdvancedSearch struct {
	Category     string
	Cuisine      string
	Description  string
	Ingredients  string
	Instructions string
	Keywords     string
	Name         string
	Source       string
	Text         string
	Tools        string
}

// Sort defines sorting options.
type Sort struct {
	IsAToZ bool
	IsZToA bool

	IsNewestToOldest bool
	IsOldestToNewest bool

	IsDefault bool
	IsRandom  bool
}

// IsSort verifies whether there is a sort option enabled.
func (s *Sort) IsSort() bool {
	return s.IsAToZ || s.IsZToA
}

// String returns a string representation of the sorting order based on the Sort struct.
func (s *Sort) String() string {
	switch {
	case s.IsAToZ:
		return "a-z"
	case s.IsZToA:
		return "z-a"
	case s.IsNewestToOldest:
		return "new-old"
	case s.IsOldestToNewest:
		return "old-new"
	case s.IsRandom:
		return "random"
	default:
		return "default"
	}
}

// NewAdvancedSearch creates an AdvancedSearch object from a search query.
func NewAdvancedSearch(query string) AdvancedSearch {
	var (
		a              AdvancedSearch
		isCat          bool
		isCuisine      bool
		isDescription  bool
		isIngredients  bool
		isInstructions bool
		isKeywords     bool
		isName         bool
		isSource       bool
		isTools        bool
	)

	reset := func() {
		isCat = false
		isCuisine = false
		isDescription = false
		isIngredients = false
		isInstructions = false
		isKeywords = false
		isName = false
		isSource = false
		isTools = false
	}

	xs := strings.Fields(strings.TrimPrefix(query, "q="))
	for _, s := range xs {
		if strings.HasPrefix(s, "cat:") {
			reset()
			isCat = true
			a.Category = strings.TrimPrefix(s, "cat:")
		} else if strings.HasPrefix(s, "cuisine:") {
			reset()
			isCuisine = true
			a.Cuisine = strings.TrimPrefix(s, "cuisine:")
		} else if strings.HasPrefix(s, "desc:") {
			reset()
			isDescription = true
			a.Description = strings.TrimPrefix(s, "desc:")
		} else if strings.HasPrefix(s, "ing:") {
			reset()
			isIngredients = true
			a.Ingredients = strings.TrimPrefix(s, "ing:")
		} else if strings.HasPrefix(s, "ins:") {
			reset()
			isInstructions = true
			a.Instructions = strings.TrimPrefix(s, "ins:")
		} else if strings.HasPrefix(s, "name:") {
			reset()
			isName = true
			a.Name = strings.TrimPrefix(s, "name:")
		} else if strings.HasPrefix(s, "src:") {
			reset()
			a.Source = strings.TrimPrefix(s, "src:")
		} else if strings.HasPrefix(s, "tag:") {
			reset()
			isKeywords = true
			a.Keywords = strings.TrimPrefix(s, "tag:")
		} else if strings.HasPrefix(s, "tool:") {
			reset()
			isTools = true
			a.Tools = strings.TrimPrefix(s, "tool:")
		} else if isCat {
			a.Category += " " + s
		} else if isCuisine {
			a.Cuisine += " " + s
		} else if isDescription {
			a.Description += " " + s
		} else if isIngredients {
			a.Ingredients += " " + s
		} else if isInstructions {
			a.Instructions += " " + s
		} else if isKeywords {
			a.Keywords += " " + s
		} else if isName {
			a.Name += " " + s
		} else if isSource {
			a.Source += " " + s
		} else if isTools {
			a.Tools += " " + s
		} else {
			a.Text += " " + s
		}
	}

	a.Text = normalizeFTSTerm(strings.TrimSpace(a.Text))
	return a
}

func normalizeFTSTerm(s string) string {
	if s == "" {
		return ""
	}
	return `"` + strings.ReplaceAll(s, "'", "''") + `"`
}

// NewSearchOptionsRecipe creates a SearchOptionsRecipe struct configured for the search method.
func NewSearchOptionsRecipe(query url.Values) SearchOptionsRecipes {
	page, err := strconv.ParseUint(query.Get("page"), 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}

	opts := SearchOptionsRecipes{
		Advanced: NewAdvancedSearch(query.Get("q")),
		Page:     page,
	}
	opts.Query = opts.Advanced.Text

	switch query.Get("sort") {
	case "a-z":
		opts.Sort.IsAToZ = true
	case "z-a":
		opts.Sort.IsZToA = true
	case "new-old":
		opts.Sort.IsNewestToOldest = true
	case "old-new":
		opts.Sort.IsOldestToNewest = true
	case "random":
		opts.Sort.IsRandom = true
	default:
		opts.Sort.IsDefault = true
	}

	return opts
}

// NewRecipeFromTextFile extracts the recipe from a text file.
func NewRecipeFromTextFile(r io.Reader) (Recipe, error) {
	recipe := NewBaseRecipe()
	blocks := extractBlocks(r)

	if len(blocks) == 1 {
		if bytes.HasPrefix(blocks[0], []byte("*****AccuChef Import File")) {
			return Recipe{}, ErrIsAccuChef
		}
		return parseSaffron(blocks[0])
	} else if len(blocks) > 1 {
		if bytes.HasSuffix(blocks[0], []byte("___________________________________________________________________________")) {
			return Recipe{}, ErrIsEasyRecipeDeluxe
		}
	}

	var (
		isDescriptionBlock      = true
		isMetaDataBlock         bool
		isIngredientsBlock      bool
		isInstructionsBlock     bool
		isPostInstructionsBlock bool
	)

	lines := bytes.Split(blocks[0], []byte("\n"))
	if len(lines) > 1 {
		recipe.Name = string(lines[0])
		recipe.Description = string(lines[1])
	} else {
		recipe.Name = string(blocks[0])
	}

	for i, block := range blocks[1:] {
		block = bytes.ReplaceAll(block, []byte("•\t"), []byte(""))
		block = bytes.ReplaceAll(block, []byte("* "), []byte(""))

		isKeyValue, isURL := isBlockKeyValues(block, &recipe)
		if isURL {
			continue
		}

		if (recipe.Description != "" && isDescriptionBlock) || (recipe.Description == "" && !isIngredientsBlock && isDescriptionBlockWhenEmpty(block)) {
			isDescriptionBlock = false
			isMetaDataBlock = len(recipe.Instructions) == 0

			if recipe.Description != "" {
				lines := bytes.Split(block, []byte("\n"))

				numColons := 0
				for _, line := range lines {
					numColons += bytes.Count(line, []byte(":"))
				}

				matches := regex.Time.FindStringSubmatch(string(lines[0]))
				if matches != nil {
					matches = slices.DeleteFunc(regex.Time.FindStringSubmatch(string(lines[0])), func(s string) bool { return s == "" })
				}

				lower := bytes.TrimSpace(bytes.ToLower(block))

				switch {
				case bytes.HasPrefix(lower, []byte("tips")):
					isDescriptionBlock = true
					isMetaDataBlock = false
					recipe.Description += "\n\n"
				case matches != nil && matches[0] != "" && len(matches) <= 3:
				case len(lines) > 1 && numColons == len(lines):
				case bytes.Contains(lower, []byte("porsjoner")) || bytes.HasPrefix(lower, []byte("tags")):
				case !isKeyValue && bytes.Contains(block, []byte(":")) || isBlockMostlyIngredients(block):
					isMetaDataBlock = false
					isIngredientsBlock = true
				case len(lines) > 0 && len(lines[0]) > 50 || (len(lines) == 1 && !regex.Digit.Match(lines[0])):
					isDescriptionBlock = true
					isMetaDataBlock = false
					recipe.Description += "\n\n"
				}
			}
		} else if isMetaDataBlock {
			isMetaDataBlock = false
			isIngredientsBlock = true
		} else if isBlockMostlyIngredients(block) {
			isDescriptionBlock = false
			isMetaDataBlock = false
			isIngredientsBlock = true
		} else if isBlockInstructions(block) || (isIngredientsBlock && !isBlockMostlyIngredients(block)) || (isDescriptionBlock && regex.Unit.FindStringIndex(string(block)) != nil && regex.Unit.FindStringIndex(string(block))[0] < 50) {
			isDescriptionBlock = false
			isIngredientsBlock = false
			isInstructionsBlock = true
		} else if isInstructionsBlock {
			isInstructionsBlock = false
			isPostInstructionsBlock = true
		}

		if isDescriptionBlock {
			recipe.Description += string(block)
		} else if isMetaDataBlock {
			parts := bytes.Split(block, []byte("\n"))
			if len(parts) == 1 && len(bytes.Split(parts[0], []byte(";"))) > 1 {
				parts = bytes.Split(parts[0], []byte(";"))
			}

			for _, line := range parts {
				processMetaData(line, &recipe)
			}
		} else if isIngredientsBlock {
			lines := strings.Split(string(block), "\n")
			if len(lines) == 1 && !strings.Contains(strings.ToLower(lines[0]), "ingred") {
				recipe.Description += lines[0]
				continue
			}

			if len(lines) < 3 && regex.Unit.Match(blocks[i+1]) {
				continue
			}

			recipe.Ingredients = append(recipe.Ingredients, lines...)
		} else if isInstructionsBlock {
			dotIndex := bytes.Index(block, []byte("."))
			lines := strings.Split(string(block), "\n")
			numSentences := len(strings.Split(lines[len(lines)-1], "."))
			lower := strings.ToLower(lines[0])

			if !strings.HasPrefix(lower, "instruction") && !strings.HasPrefix(lower, "steps") && !strings.Contains(lower, "directions") && !strings.HasPrefix(lower, "fremgang") && len(lines) < 3 && (dotIndex == -1 || dotIndex > 4) && numSentences < 3 {
				if len(lines) == 1 && bytes.Contains(block, []byte("personer")) {
					processMetaData([]byte(lines[0]), &recipe)
				} else if !bytes.Contains(block, []byte("Slik gjør du")) {
					isIngredientsBlock = true
					recipe.Ingredients = append(recipe.Ingredients, string(bytes.TrimSpace(block)))
				}
				continue
			}

			numWords := len(strings.Split(lines[0], " "))
			if numWords > 1 && numWords < 4 && !strings.Contains(lower, "prep") && !strings.Contains(lower, "step") && !strings.Contains(lower, "slik") && !strings.HasPrefix(lower, "for") {
				isInstructionsBlock = false
				isIngredientsBlock = true
				recipe.Ingredients = append(recipe.Ingredients, lines...)
				continue
			}

			for _, line := range strings.Split(string(block), "\n") {
				before, after, found := strings.Cut(line, ".")
				if found {
					_, err := strconv.ParseInt(before, 10, 64)
					if err == nil {
						recipe.Instructions = append(recipe.Instructions, strings.TrimSpace(after))
					} else {
						recipe.Instructions = append(recipe.Instructions, strings.TrimSpace(line))
					}
				} else {
					recipe.Instructions = append(recipe.Instructions, strings.TrimSpace(line))
				}
			}
		} else if isPostInstructionsBlock {
			parse, err := url.Parse(string(block))
			if err == nil && bytes.HasPrefix(block, []byte("http")) {
				recipe.URL = parse.String()
			} else {
				lines := strings.Split(strings.TrimSpace(string(block)), "\n")
				for i, line := range lines {
					lines[i] = strings.TrimSpace(line)
				}
				recipe.Instructions = append(recipe.Instructions, lines...)
			}
		}
	}

	if recipe.Category == "" {
		recipe.Category = "uncategorized"
	}

	for i, ing := range recipe.Ingredients {
		recipe.Ingredients[i] = strings.TrimSpace(ing)
	}

	recipe.Ingredients = slices.DeleteFunc(recipe.Ingredients, func(s string) bool {
		s = strings.ToLower(s)
		if s == "" {
			return true
		}

		words := []string{"prep", "ingredien", "slik", "tilbe", "instruct", "method"}
		for _, word := range words {
			if strings.HasPrefix(s, word) || len(word) == 0 {
				return true
			}
		}
		return false
	})

	recipe.Instructions = slices.DeleteFunc(recipe.Instructions, func(s string) bool {
		if s == "" {
			return true
		}

		s = strings.ToLower(s)
		words := []string{"instr", "method", "recipe", "direction", "step by step", "slik", "fremg", "framg", "preparation", "steps"}
		for _, word := range words {
			if strings.HasPrefix(s, word) || len(word) == 0 {
				return true
			}
		}
		return false
	})

	if recipe.Yield == 0 {
		recipe.Yield = 1
	}

	return recipe, nil
}

func extractBlocks(r io.Reader) [][]byte {
	var (
		scanner     = bufio.NewScanner(r)
		i           = 0
		blocks      = make([][]byte, 0)
		isFirstLine = true
	)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		isEmpty := bytes.Equal(line, []byte(""))

		if isFirstLine && i == 0 && isEmpty {
			continue
		}
		isFirstLine = false

		if bytes.Equal(line, []byte("")) {
			blocks[i] = bytes.TrimSpace(blocks[i])
			i++
		}

		if len(blocks) <= i {
			blocks = append(blocks, make([]byte, 0))
		}
		blocks[i] = append(blocks[i], append(line, '\n')...)
	}

	blocks[len(blocks)-1] = bytes.TrimSpace(blocks[len(blocks)-1])
	return blocks
}

func isBlockKeyValues(block []byte, recipe *Recipe) (isKeyValue, isURL bool) {
	colonIndex := bytes.Index(block, []byte(":"))
	isKeyValue = colonIndex >= 0

	if colonIndex >= 0 {
		rawURL := bytes.TrimSpace(block)
		before, _, ok := bytes.Cut(rawURL, []byte("\n"))
		if ok {
			rawURL = before
		}
		parse, err := url.Parse(string(rawURL))
		if bytes.HasPrefix(block, []byte("http")) && err == nil {
			recipe.URL = parse.String()
			return false, true
		}

		_, _, found := bytes.Cut(block[colonIndex:], []byte("\n"))
		if found {
			_, after, found := bytes.Cut(block[colonIndex:], []byte(":"))
			if found {
				l := bytes.Split(after, []byte("\n"))[0]
				isKeyValue = len(bytes.TrimSpace(l)) != 0
			}
		}
	}

	return isKeyValue, false
}

func isDescriptionBlockWhenEmpty(block []byte) bool {
	isSmall := len(block) < 50 && regex.Digit.Match(block)

	lowerBlock := bytes.ToLower(block)
	hasMetaKeyword := bytes.Contains(lowerBlock, []byte("yield")) || bytes.HasPrefix(lowerBlock, []byte("serves")) ||
		bytes.Contains(lowerBlock, []byte("porsjoner")) || bytes.HasPrefix(lowerBlock, []byte("prep")) ||
		bytes.HasPrefix(lowerBlock, []byte("course"))

	return isSmall || hasMetaKeyword
}

func isBlockMostlyIngredients(block []byte) bool {
	lines := bytes.Split(block, []byte("\n"))
	numLines := len(lines)
	if numLines < 2 {
		return false
	}
	numLines--

	hits := 0.
	incremental := 0
	for i, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}

		lineStr := units.ReplaceVulgarFractions(string(line))
		idx := regex.Digit.FindStringIndex(lineStr)
		if len(idx) == 2 {
			num, err := strconv.ParseInt(string(lineStr[idx[1]-1]), 10, 64)
			if num == int64(i+1) {
				incremental++
			}

			if err == nil && idx[0] == 0 {
				hits++
				continue
			}
		}
	}

	threshold := 0.6
	if numLines < 4 {
		threshold = 0.3
	}
	return incremental != len(lines[1:]) && hits/float64(numLines) >= threshold
}

func isBlockInstructions(block []byte) bool {
	lines := bytes.Split(block, []byte("\n"))
	hits := 0.
	for _, line := range lines {
		dotIndex := bytes.IndexByte(line, '.')
		before, _, found := bytes.Cut(line, []byte("."))
		if found && dotIndex >= 0 && dotIndex < 4 {
			_, err := strconv.ParseInt(string(before), 10, 64)
			if err == nil {
				hits++
			}
		}
	}
	return hits/float64(len(lines)) >= 0.8
}

func processMetaData(line []byte, recipe *Recipe) {
	line = bytes.ToLower(line)
	if bytes.HasPrefix(line, []byte("course")) || bytes.HasPrefix(line, []byte("type")) {
		_, after, found := bytes.Cut(line, []byte(":"))
		if found {
			before, _, found := bytes.Cut(after, []byte("/"))
			if found {
				recipe.Category = strings.TrimSpace(string(before))
			} else {
				before, _, _ = bytes.Cut(after, []byte(","))
				recipe.Category = strings.TrimSpace(string(before))
			}
		}
	} else if bytes.HasPrefix(line, []byte("cuisine")) || bytes.HasPrefix(line, []byte("opprin")) {
		_, after, found := bytes.Cut(line, []byte(":"))
		if found {
			before, _, _ := bytes.Cut(after, []byte(","))
			recipe.Cuisine = strings.TrimSpace(string(before))
		}
	} else if bytes.Contains(line, []byte("time")) || bytes.HasPrefix(line, []byte("prep")) || bytes.HasPrefix(line, []byte("cook")) ||
		bytes.HasPrefix(line, []byte("stek")) || bytes.Contains(line, []byte("min")) || bytes.HasPrefix(line, []byte("tilber")) {
		before, after, found := bytes.Cut(line, []byte("-"))
		if found && bytes.Contains(after, []byte("min")) && regex.Digit.MatchString(string(before)) {
			dur, err := time.ParseDuration(regex.Digit.FindString(string(before)) + "m")
			if err == nil {
				if recipe.Times.Prep == 0 {
					recipe.Times.Prep += dur
				} else {
					recipe.Times.Cook += dur
				}
			}
			return
		}

		dur := duration.From(string(line))
		if dur > 0 {
			if bytes.HasPrefix(line, []byte("prep")) || bytes.HasPrefix(line, []byte("hands")) {
				recipe.Times.Prep += dur
			} else if bytes.HasPrefix(line, []byte("cook")) || bytes.HasPrefix(line, []byte("oven")) ||
				bytes.HasPrefix(line, []byte("simmer")) || bytes.HasPrefix(line, []byte("stek")) {
				recipe.Times.Cook += dur
			} else if bytes.HasPrefix(line, []byte("total")) || bytes.Contains(line, []byte("timer")) {
				recipe.Times.Cook = dur - recipe.Times.Prep
			} else if recipe.Times.Prep == 0 && recipe.Times.Cook == 0 {
				recipe.Times.Prep += dur
			}
		}
	} else if bytes.Contains(line, []byte("serv")) || bytes.HasPrefix(line, []byte("yield")) || bytes.Contains(line, []byte("portion")) ||
		bytes.Contains(line, []byte("stk")) || bytes.HasPrefix(line, []byte("makes")) || bytes.Contains(line, []byte("pers")) ||
		bytes.Contains(line, []byte("porsjoner")) || bytes.Contains(line, []byte("person")) || bytes.Contains(line, []byte("voksne")) ||
		bytes.Contains(line, []byte("styk")) || bytes.Contains(line, []byte("til")) {
		yield, err := strconv.ParseInt(regex.Digit.FindString(string(line)), 10, 16)
		if err == nil {
			recipe.Yield = int16(yield)
		}
	} else if recipe.Times.Prep == 0 && recipe.Times.Cook == 0 && (bytes.HasPrefix(line, []byte("total")) || bytes.HasPrefix(line, []byte("tid"))) {
		before, after, found := bytes.Cut(line, []byte(","))
		if found {
			dur, err := time.ParseDuration(regex.Digit.FindString(string(before)) + "h")
			if err == nil {
				recipe.Times.Prep += dur
			}

			dur, err = time.ParseDuration(regex.Digit.FindString(string(after)) + "m")
			if err == nil {
				recipe.Times.Prep += dur
			}
		} else if bytes.Contains(line, []byte("min")) {
			dur, err := time.ParseDuration(regex.Digit.FindString(string(line)) + "m")
			if err == nil {
				recipe.Times.Prep += dur
			}
		}
	} else if bytes.HasPrefix(line, []byte("keyword")) || bytes.HasPrefix(line, []byte("hove")) || bytes.HasPrefix(line, []byte("anle")) ||
		bytes.HasPrefix(line, []byte("karak")) || bytes.HasPrefix(line, []byte("tag")) {
		_, after, found := strings.Cut(string(line), ":")
		if found {
			for _, s := range strings.Split(after, ",") {
				recipe.Keywords = append(recipe.Keywords, strings.TrimSpace(s))
			}
		}
	} else if bytes.HasPrefix(line, []byte("tool")) {
		_, after, found := strings.Cut(string(line), ":")
		if found {
			for _, s := range strings.Split(after, ",") {
				recipe.Tools = append(recipe.Tools, NewHowToTool(strings.TrimSpace(s)))
			}
		}
	}
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

// NewRecipesFromRecipeKeeper extracts the recipes from a Recipe Keeper export.
func NewRecipesFromRecipeKeeper(root *goquery.Document) Recipes {
	rs := NewRecipeSchema()

	nodes := root.Find(".recipe-details")
	recipes := make(Recipes, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		nodes := sel.Find("meta[itemprop='recipeCategory']")
		keywords := make([]string, 0, nodes.Length()+1)
		nodes.Each(func(_ int, sub *goquery.Selection) {
			s := sub.AttrOr("content", "")
			if s != "" {
				keywords = append(keywords, s)
			}
		})

		course := sel.Find("span[itemprop='recipeCourse']").Text()
		if course == "" {
			course = "uncategorized"
		}

		var yield int16
		parsed, err := strconv.ParseInt(regex.Digit.FindString(sel.Find("span[itemprop='recipeYield']").Text()), 10, 16)
		if err == nil {
			yield = int16(parsed)
		}

		var prep time.Duration
		rs.PrepTime = sel.Find("meta[itemprop='prepTime']").AttrOr("content", "")
		_, after, ok := strings.Cut(rs.PrepTime, "PT")
		if ok {
			prep, _ = time.ParseDuration(strings.ToLower(after))
		}

		var cook time.Duration
		cookTime := sel.Find("meta[itemprop='cookTime']").AttrOr("content", "")
		_, after, ok = strings.Cut(cookTime, "PT")
		if ok {
			cook, _ = time.ParseDuration(strings.ToLower(after))
		}

		nodes = sel.Find("div[itemprop='recipeIngredients'] p")
		ingredients := make([]string, 0, nodes.Length())
		nodes.Each(func(_ int, sub *goquery.Selection) {
			s := strings.TrimSpace(sub.Text())
			if s != "" {
				ingredients = append(ingredients, s)
			}
		})

		nodes = sel.Find("div[itemprop='recipeDirections'] p")
		instructions := make([]string, 0, nodes.Length())
		nodes.Each(func(_ int, sub *goquery.Selection) {
			s := strings.TrimSpace(sub.Text())
			if s == "" {
				return
			}

			dotIndex := strings.Index(s, ".")
			if dotIndex >= 0 && dotIndex < 5 {
				s = strings.TrimSpace(s[dotIndex+1:])
			}

			instructions = append(instructions, s)
		})

		notes := strings.TrimSpace(sel.Find("div[itemprop='recipeNotes']").Text())
		if notes != "" {
			instructions = append(instructions, "Notes:\n"+notes)
		}

		extractNut := func(selector, unit string) string {
			v := sel.Find("meta[itemprop='"+selector+"']").AttrOr("content", "")
			if v != "" {
				v += unit
			}
			return v
		}

		recipes = append(recipes, Recipe{
			Category:     course,
			CreatedAt:    time.Time{},
			Description:  "Imported from Recipe Keeper.",
			Ingredients:  ingredients,
			Instructions: instructions,
			Keywords:     keywords,
			Name:         sel.Find("h2[itemprop='name']").Text(),
			Nutrition: Nutrition{
				Calories:           extractNut("recipeNutCalories", " kcal"),
				Cholesterol:        extractNut("recipeNutCholesterol", "mg"),
				Fiber:              extractNut("recipeNutDietaryFiber", "g"),
				IsPerServing:       true,
				Protein:            extractNut("recipeNutProtein", "g"),
				SaturatedFat:       extractNut("recipeNutSaturatedFat", "g"),
				Sodium:             extractNut("recipeNutSodium", "mg"),
				Sugars:             extractNut("recipeNutSugars", "g"),
				TotalCarbohydrates: extractNut("recipeNutTotalCarbohydrate", "g"),
				TotalFat:           extractNut("recipeNutTotalFat", "g"),
			},
			Times:     Times{Prep: prep, Cook: cook},
			Tools:     make([]HowToItem, 0),
			UpdatedAt: time.Time{},
			URL:       "Recipe Keeper",
			Yield:     yield,
		})
	})
	return recipes
}

func parseSaffron(block []byte) (Recipe, error) {
	var (
		isIngredients  bool
		isInstructions bool
		isYield        bool

		recipe = NewBaseRecipe()
	)

	for _, b := range bytes.Split(block, []byte("\n")) {
		before, after, ok := bytes.Cut(b, []byte(":"))
		before = bytes.TrimSpace(before)
		after = bytes.TrimSpace(after)

		if !ok {
			v := strings.TrimSpace(string(before))
			if isIngredients {
				recipe.Ingredients = append(recipe.Ingredients, v)
			} else if isInstructions {
				recipe.Instructions = append(recipe.Instructions, v)
			}
			continue
		}

		v := strings.TrimSpace(string(after))
		if strings.TrimSpace(v) == "" {
			if bytes.Equal(before, []byte("Ingredients")) {
				isIngredients = true
			} else if bytes.Equal(before, []byte("Instructions")) {
				isIngredients = false
				isInstructions = true
			}
			continue
		}

		if isYield {
			t := duration.From(v)
			if t > 0 && recipe.Times.Prep == 0 {
				recipe.Times.Prep = t
			} else if t > 0 && recipe.Times.Cook == 0 {
				recipe.Times.Cook = t
			}
		}

		switch string(before) {
		case "Title":
			recipe.Name = v
		case "Description":
			recipe.Description = v
		case "Source", "Original URL":
			recipe.URL = v
		case "Yield":
			isYield = true
			parsed, err := strconv.ParseInt(v, 10, 16)
			if err == nil {
				recipe.Yield = int16(parsed)
			}
		case "Cookbook":
			isYield = false
		}
	}

	return recipe, nil
}

// NewRecipesFromAccuChef extracts all recipes from an exported Accuchef file.
func NewRecipesFromAccuChef(r io.Reader) Recipes {
	blocks := extractBlocks(r)
	if len(blocks) == 0 {
		return nil
	}

	parts := bytes.Split(blocks[0], []byte("*****AccuChef Import File (C)SIVART Software http://www.AccuChef.com"))
	recipes := make(Recipes, 0, len(parts))

	for _, b := range parts[1:] {
		b = bytes.TrimSpace(b)

		var (
			ingredient  string
			instruction string
			isNutrition bool
			nutrition   string
			notes       string
		)

		recipe := NewBaseRecipe()
		recipe.URL = "AccuChef"

		for _, line := range bytes.Split(b, []byte("\n")) {
			if len(line) <= 1 {
				continue
			}

			c := line[0]
			content := string(bytes.TrimSpace(line[1:]))

			switch c {
			case 'A':
				recipe.Name = content
			case 'B':
				recipe.Category = strings.Split(content, ",")[0]
			case 'C':
				recipe.Keywords = append(recipe.Keywords, content)
			case 'D':
				parsed, err := strconv.ParseInt(content, 10, 16)
				if err == nil {
					recipe.Yield = int16(parsed)
				}
			case 'P':
				split := strings.Split(content, ":")
				if len(split) == 2 {
					h := strings.TrimSpace(split[0])
					m := strings.TrimSpace(split[1])
					recipe.Times.Prep, _ = time.ParseDuration(h + "h" + m + "m")
				}
			case 'F':
				notes += content + " "
			case 'H':
				if ingredient != "" {
					recipe.Ingredients = append(recipe.Ingredients, ingredient)
				}
				ingredient = content
			case 'I':
				ingredient += " " + content
			case 'K':
				ingredient += " (" + content + ")"
			case 'J':
				if isNutrition {
					if strings.HasSuffix(content, "~") {
						goto Break
					}
					nutrition += content + " "
					continue
				}

				if ingredient != "" {
					recipe.Ingredients = append(recipe.Ingredients, ingredient)
					ingredient = ""
				}

				if strings.HasPrefix(content, "Per ") {
					instruction = strings.TrimSuffix(strings.TrimSpace(instruction), "~")
					recipe.Instructions = append(recipe.Instructions, strings.Join(strings.Fields(instruction), " "))
					instruction = ""

					_, after, _ := strings.Cut(content, ":")
					isNutrition = true
					nutrition += after + " "
					continue
				}

				if content == "~" {
					instruction = strings.TrimSuffix(strings.TrimSpace(instruction), "~")
					recipe.Instructions = append(recipe.Instructions, strings.Join(strings.Fields(instruction), " "))
					instruction = ""
					continue
				}

				instruction += content + " "
			}
		}

	Break:
		if instruction != "" {
			instruction = strings.TrimSuffix(strings.TrimSpace(instruction), "~")
			recipe.Instructions = append(recipe.Instructions, strings.Join(strings.Fields(instruction), " "))
		}

		split := strings.Split(strings.TrimSpace(nutrition), ";")
		for i, s := range split {
			all := strings.Split(strings.TrimSpace(s), " ")
			if all[0] == "0" || all[0] == "" || len(all) <= 2 {
				continue
			}

			v := all[0] + " " + all[1]

			if i == 0 && len(all) > 3 {
				recipe.Nutrition.Calories = v
				continue
			}

			switch all[2] {
			case "protein":
				recipe.Nutrition.Protein = v
			case "tot":
				recipe.Nutrition.TotalFat = v
			case "sat":
				recipe.Nutrition.SaturatedFat = v
			case "carb":
				recipe.Nutrition.TotalCarbohydrates = v
			case "fiber":
				recipe.Nutrition.Fiber = v
			case "sodium":
				recipe.Nutrition.Sodium = v
			case "cholesterol":
				recipe.Nutrition.Cholesterol = v
			}
		}

		recipe.Instructions = append(recipe.Instructions, strings.TrimSpace(notes))
		recipe.Instructions = slices.DeleteFunc(recipe.Instructions, func(s string) bool { return s == "" })

		recipes = append(recipes, recipe)
	}

	return recipes
}

// NewRecipesFromEasyRecipeDeluxe extracts the recipes from a MasterCook file.
func NewRecipesFromEasyRecipeDeluxe(r io.Reader) Recipes {
	scanner := bufio.NewScanner(r)

	var (
		isIngredients  bool
		isInstructions bool
		isNewRecipe    bool
		isMetadata     = true
		isNotes        bool

		delim        = []byte("___________________________________________________________________________")
		previousLine []byte
		numNewLines  = 0

		recipe  = Recipe{URL: "Easy Recipe Deluxe"}
		recipes Recipes
	)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if numNewLines == 2 {
			isNewRecipe = len(recipe.Instructions) > 0
			numNewLines = 0
		}

		if isMetadata {
			if bytes.Equal(line, delim) && bytes.HasPrefix(previousLine, []byte("Cooking Time:")) {
				isMetadata = false
				continue
			} else if bytes.HasPrefix(line, []byte("Category:")) {
				_, after, _ := bytes.Cut(line, []byte("Category:"))
				recipe.Category = string(bytes.TrimSpace(after))
			} else if bytes.HasPrefix(line, []byte("Main Ingredient:")) {
				_, after, _ := bytes.Cut(line, []byte("Main Ingredient:"))
				recipe.Keywords = []string{string(bytes.TrimSpace(after))}
			} else if bytes.HasPrefix(line, []byte("Cuisine Style:")) {
				_, after, _ := bytes.Cut(line, []byte("Cuisine Style:"))
				recipe.Cuisine = string(bytes.TrimSpace(after))
			} else if bytes.HasPrefix(line, []byte("Yield:")) {
				_, after, _ := bytes.Cut(line, []byte("Yield:"))
				for _, b := range bytes.Split(after, []byte(" ")) {
					parsed, err := strconv.ParseInt(string(b), 10, 16)
					if err == nil {
						recipe.Yield = int16(parsed)
						break
					}
				}
			} else if bytes.HasPrefix(line, []byte("Preparation Time:")) {
				_, after, _ := bytes.Cut(line, []byte("Preparation Time:"))
				after = bytes.TrimSpace(after)
				if !bytes.Equal(after, []byte("")) {
					parts := strings.Split(string(after), ":")
					var t string
					if parts[0] != "" {
						t += parts[0] + "h"
					}
					if parts[1] != "" {
						t += parts[1] + "m"
					}
					recipe.Times.Prep, _ = time.ParseDuration(t)
				}
			} else if bytes.HasPrefix(line, []byte("Cooking Time:")) {
				_, after, _ := bytes.Cut(line, []byte("Cooking Time:"))
				after = bytes.TrimSpace(after)
				if !bytes.Equal(after, []byte("")) {
					parts := strings.Split(string(after), ":")
					var t string
					if parts[0] != "" {
						t += parts[0] + "h"
					}
					if parts[1] != "" {
						t += parts[1] + "m"
					}
					recipe.Times.Cook, _ = time.ParseDuration(t)
				}
			} else if bytes.Equal(line, delim) && !bytes.HasPrefix(previousLine, []byte("Cooking Time:")) {
				recipe.Name = string(previousLine)
			}

			previousLine = line
			continue
		} else if bytes.Equal(line, []byte("[Ingredients]")) {
			isMetadata = false
			isIngredients = true
			numNewLines = 0
			continue
		} else if !isInstructions && bytes.Equal(line, []byte("[Instruction]")) {
			isIngredients = false
			isInstructions = true
			numNewLines = 0
			continue
		} else if isInstructions && bytes.Equal(line, []byte("[Note]")) {
			isInstructions = false
			isNotes = true
			numNewLines = 0
			continue
		} else if bytes.Equal(line, []byte("")) {
			previousLine = line
			numNewLines++
			continue
		}

		if isNewRecipe {
			isNewRecipe = false
			isInstructions = false
			isMetadata = true

			recipes = append(recipes, recipe)
			recipe = Recipe{Name: string(line), URL: "Easy Recipe Deluxe"}
			numNewLines = 0
		} else if isIngredients {
			recipe.Ingredients = append(recipe.Ingredients, string(line))
			numNewLines = 0
			continue
		} else if isInstructions {
			dotIndex := bytes.Index(line, []byte("."))
			if dotIndex >= 0 && dotIndex < 5 {
				line = bytes.TrimSpace(line[dotIndex+1:])
			}
			recipe.Instructions = append(recipe.Instructions, string(line))
			numNewLines = 0
			continue
		} else if isNotes {
			if bytes.Contains(line, []byte("Calories;")) {
				parts := bytes.Split(line, []byte(";"))
				for _, part := range parts {
					sub := bytes.Split(part, []byte(" "))
					v := strings.TrimSpace(string(sub[0])) + " " + strings.TrimSpace(string(sub[1]))

					if bytes.HasSuffix(part, []byte("Calories")) {
						recipe.Nutrition.Calories = strings.TrimSpace(string(sub[len(sub)-1])) + " kcal"
					} else if bytes.HasSuffix(part, []byte("Sat Fat")) {
						recipe.Nutrition.SaturatedFat = v
					} else if bytes.HasSuffix(part, []byte("Fat")) {
						recipe.Nutrition.TotalFat = v
					} else if bytes.HasSuffix(part, []byte("Cholesterol")) {
						recipe.Nutrition.Cholesterol = v
					} else if bytes.HasSuffix(part, []byte("Sodium")) {
						recipe.Nutrition.Sodium = v
					} else if bytes.HasSuffix(part, []byte("Carb")) {
						recipe.Nutrition.TotalCarbohydrates = v
					} else if bytes.HasSuffix(part, []byte("Fiber")) {
						recipe.Nutrition.Fiber = v
					} else if bytes.HasSuffix(part, []byte("Protein.")) {
						recipe.Nutrition.Protein = v
					}
				}
			} else if !bytes.Equal(line, delim) {
				recipe.Instructions = append(recipe.Instructions, string(line))
			}
		}

		if bytes.Equal(line, delim) && !bytes.HasPrefix(previousLine, []byte("Cooking Time:")) {
			isNewRecipe = len(recipe.Ingredients) > 0
			recipe.Name = string(previousLine)
			previousLine = line
			numNewLines = 0
			continue
		}

		previousLine = line
		numNewLines = 0
	}

	return append(recipes, recipe)
}

func init() {
	c, err := word2number.NewConverter("en")
	if err != nil {
		panic(err)
	}
	wordConverter = c
}
