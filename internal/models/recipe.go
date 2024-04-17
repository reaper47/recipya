package models

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/donna-legal/word2number"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/units"
	"github.com/reaper47/recipya/internal/utils/duration"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"github.com/reaper47/recipya/internal/utils/regex"
	"io"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	wordConverter *word2number.Converter
)

// ErrIsAccuChef is the error for when a file is an AccuChef one.
var ErrIsAccuChef = errors.New("accuchef") // TODO: Place errors somewhere else.

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
		Image:        uuid.Nil,
		Ingredients:  make([]string, 0),
		Instructions: make([]string, 0),
		Keywords:     make([]string, 0),
		Tools:        make([]string, 0),
	}
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
	IsPerServing       bool
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
	n.SaturatedFat = regex.Digit.ReplaceAllStringFunc(n.SaturatedFat, scale)
	n.Sodium = regex.Digit.ReplaceAllStringFunc(n.Sodium, scale)
	n.Sugars = regex.Digit.ReplaceAllStringFunc(n.Sugars, scale)
	n.TotalCarbohydrates = regex.Digit.ReplaceAllStringFunc(n.TotalCarbohydrates, scale)
	n.TotalFat = regex.Digit.ReplaceAllStringFunc(n.TotalFat, scale)
	n.UnsaturatedFat = regex.Digit.ReplaceAllStringFunc(n.UnsaturatedFat, scale)
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
	CookbookID   int64
	IsByName     bool
	IsFullSearch bool
	Page         uint64
	Sort         Sort
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

// NewSearchOptionsRecipe creates a SearchOptionsRecipe struct configured for the search method.
func NewSearchOptionsRecipe(mode, sort string, page uint64) SearchOptionsRecipes {
	opts := SearchOptionsRecipes{Page: page}

	switch mode {
	case "name":
		opts.IsByName = true
	case "full":
		opts.IsFullSearch = true
	default:
		opts.IsByName = true
	}

	switch sort {
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
				recipe.Tools = append(recipe.Tools, strings.TrimSpace(s))
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
	nodes := root.Find(".recipe-details")
	recipes := make(Recipes, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		nodes := sel.Find("meta[itemprop='recipeCategory']")
		keywords := make([]string, 0, nodes.Length()+1)
		nodes.Each(func(_ int, sub *goquery.Selection) {
			s, _ := sub.Attr("content")
			keywords = append(keywords, s)
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
		prepTime, _ := sel.Find("meta[itemprop='prepTime']").Attr("content")
		_, after, ok := strings.Cut(prepTime, "PT")
		if ok {
			prep, _ = time.ParseDuration(strings.ToLower(after))
		}

		var cook time.Duration
		cookTime, _ := sel.Find("meta[itemprop='cookTime']").Attr("content")
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
			v, _ := sel.Find("meta[itemprop='" + selector + "']").Attr("content")
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
			Tools:     make([]string, 0),
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

		recipe = Recipe{Category: "uncategorized", Yield: 1}
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

			notes  string
			recipe = Recipe{
				Category: "uncategorized",
				Keywords: make([]string, 0),
				URL:      "From AccuChef",
				Yield:    1,
			}
		)

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

func init() {
	c, err := word2number.NewConverter("en")
	if err != nil {
		panic(err)
	}
	wordConverter = c
}
