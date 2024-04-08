package models

import (
	"encoding/json"
	"fmt"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RecipeSchema is a representation of the Recipe schema (https://schema.org/Recipe).
type RecipeSchema struct {
	AtContext       string          `json:"@context"`
	AtType          SchemaType      `json:"@type"`
	Category        Category        `json:"recipeCategory"`
	CookTime        string          `json:"cookTime"`
	CookingMethod   CookingMethod   `json:"cookingMethod"`
	Cuisine         Cuisine         `json:"recipeCuisine"`
	DateCreated     string          `json:"dateCreated"`
	DateModified    string          `json:"dateModified"`
	DatePublished   string          `json:"datePublished"`
	Description     Description     `json:"description"`
	Keywords        Keywords        `json:"keywords"`
	Image           Image           `json:"image"`
	Ingredients     Ingredients     `json:"recipeIngredient"`
	Instructions    Instructions    `json:"recipeInstructions"`
	Name            string          `json:"name"`
	NutritionSchema NutritionSchema `json:"nutrition"`
	PrepTime        string          `json:"prepTime"`
	Tools           Tools           `json:"tool"`
	TotalTime       string          `json:"totalTime"`
	Yield           Yield           `json:"recipeYield"`
	URL             string          `json:"url"`
}

// Recipe transforms the RecipeSchema to a Recipe.
func (r *RecipeSchema) Recipe() (*Recipe, error) {
	if r.AtType.Value != "Recipe" {
		return nil, fmt.Errorf("RecipeSchema %#v is not based on the Schema or the field is missing", r)
	}

	var category string
	if r.Category.Value == "" {
		category = "uncategorized"
	} else {
		category = strings.TrimSpace(strings.ToLower(r.Category.Value))
	}

	times, err := NewTimes(r.PrepTime, r.CookTime)
	if err != nil {
		return nil, err
	}

	nutrition, err := r.NutritionSchema.nutrition()
	if err != nil {
		return nil, err
	}

	created := r.DateCreated
	if r.DatePublished != "" {
		created = r.DatePublished
	}

	var createdAt time.Time
	if created != "" {
		before, _, found := strings.Cut(created, "T")
		if !found {
			before, _, _ = strings.Cut(created, " ")
		}

		createdAt, err = time.Parse(time.DateOnly, before)
		if err != nil {
			return nil, fmt.Errorf("could not parse createdAt date %s: %w", created, err)
		}
	}

	image, err := uuid.Parse(r.Image.Value)
	if err != nil {
		image = uuid.Nil
	}

	updatedAt := createdAt
	if r.DateModified != "" {
		before, _, found := strings.Cut(r.DateModified, "T")
		if !found {
			before, _, _ = strings.Cut(r.DateModified, " ")
		}

		updatedAt, err = time.Parse(time.DateOnly, before)
		if err != nil {
			return nil, fmt.Errorf("could not parse modifiedAt date %s: %w", r.DateModified, err)
		}
	}

	recipe := Recipe{
		Category:     category,
		CreatedAt:    createdAt,
		Cuisine:      r.Cuisine.Value,
		Description:  r.Description.Value,
		ID:           0,
		Image:        image,
		Ingredients:  r.Ingredients.Values,
		Instructions: r.Instructions.Values,
		Keywords:     extensions.Unique(strings.Split(r.Keywords.Values, ",")),
		Name:         r.Name,
		Nutrition:    nutrition,
		Times:        times,
		Tools:        r.Tools.Values,
		UpdatedAt:    updatedAt,
		URL:          r.URL,
		Yield:        r.Yield.Value,
	}

	recipe.Normalize()
	return &recipe, nil
}

// SchemaType holds the type of the schema. It should be "Recipe".
type SchemaType struct {
	Value string
}

// MarshalJSON encodes the schema's type.
func (s *SchemaType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}

// UnmarshalJSON decodes the type of the schema.
// The type "Recipe" will be searched for if the data is an array.
func (s *SchemaType) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		m := make(map[string]string)
		err = json.Unmarshal(data, &m)
		if err != nil {
			return err
		}

		val, ok := m["Value"]
		if !ok {
			return fmt.Errorf("could not decode description %q", data)
		}
		v = val
	}

	switch x := v.(type) {
	case string:
		s.Value = x
	case map[string]any:
		v, ok := x["Value"]
		if ok {
			s.Value = v.(string)
		}
	case []any:
		for _, t := range x {
			if t.(string) == "Recipe" {
				s.Value = "Recipe"
			}
		}
	}

	return nil
}

// Category holds a recipe's category.
type Category struct {
	Value string
}

// MarshalJSON encodes the category.
func (c *Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

// UnmarshalJSON decodes the category according to the schema (https://schema.org/recipeCategory).
// The schema specifies that the expected value is of type Text. However, some
// websites send the category in an array, which explains the need for this function.
func (c *Category) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		c.Value = x
	case []any:
		if len(x) > 0 {
			c.Value = x[0].(string)
		}
	case map[string]any:
		v, ok := x["Value"]
		if ok {
			c.Value = v.(string)
		}
	}

	if c.Value != "" {
		c.Value = strings.ReplaceAll(c.Value, "&amp;", "&")
	}
	return nil
}

// CookingMethod holds a recipe's category.
type CookingMethod struct {
	Value string
}

// MarshalJSON encodes the cuisine.
func (c *CookingMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

// UnmarshalJSON decodes the cooking method according to the schema (https://schema.org/cookingMethod).
// The schema specifies that the expected value is of type Text. However, some
// websites send the cuisine in an array, which explains the need for this function.
func (c *CookingMethod) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		c.Value = x
	case []any:
		if len(x) > 0 {
			c.Value = x[0].(string)
		}
	}
	return nil
}

// Cuisine holds a recipe's category.
type Cuisine struct {
	Value string
}

// MarshalJSON encodes the cuisine.
func (c *Cuisine) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

// UnmarshalJSON decodes the cuisine according to the schema (https://schema.org/recipeCuisine).
// The schema specifies that the expected value is of type Text. However, some
// websites send the cuisine in an array, which explains the need for this function.
func (c *Cuisine) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		c.Value = x
	case []any:
		if len(x) > 0 {
			c.Value = x[0].(string)
		}
	}
	return nil
}

// Description holds a description.
type Description struct {
	Value string
}

// MarshalJSON encodes the description.
func (d *Description) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}

// UnmarshalJSON decodes the description field of a recipe.
func (d *Description) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		m := make(map[string]string)
		err = json.Unmarshal(data, &m)
		if err != nil {
			return err
		}

		v, ok := m["Value"]
		if !ok {
			return fmt.Errorf("could not decode description %q", data)
		}
		s = v
	}

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\u202f", "")
	s = strings.ReplaceAll(s, "\u00a0", "")
	s = strings.ReplaceAll(s, "&quot;", "")
	s = strings.ReplaceAll(s, "”", `"`)
	s = strings.ReplaceAll(s, "\u00ad", "")
	d.Value = s

	return nil
}

// Keywords holds keywords as a single string.
type Keywords struct {
	Values string
}

// MarshalJSON encodes the keywords.
func (k *Keywords) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.Values)
}

// UnmarshalJSON decodes the recipe's keywords according to the schema (https://schema.org/keywords).
// Some websites store the keywords in an array, which explains the need
// for a custom decoder.
func (k *Keywords) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		k.Values = strings.TrimSpace(x)
	case []any:
		var xs []string
		for _, v := range x {
			xs = append(xs, v.(string))
		}
		k.Values = strings.TrimSpace(strings.Join(xs, ","))
	}

	k.Values = strings.TrimRight(k.Values, ",")
	return nil
}

// Image holds a recipe's image. The JSON fields correspond.
type Image struct {
	Value string
}

// MarshalJSON encodes the image.
func (i *Image) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Value)
}

// UnmarshalJSON decodes the image according to the schema (https://schema.org/image).
func (i *Image) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		i.Value = x
	case []any:
		if len(x) > 0 {
			switch y := x[0].(type) {
			case string:
				i.Value = y
			case map[string]any:
				u, ok := y["url"]
				if ok {
					i.Value = u.(string)
				}
			}
		}
	case map[string]any:
		url, ok := v.(map[string]any)["url"]
		if ok {
			i.Value = url.(string)
			break
		}

		val, ok := v.(map[string]any)["Value"]
		if ok {
			i.Value = val.(string)
		}
	case any:
		url, ok := v.(map[string]any)["url"]
		if ok {
			i.Value = url.(string)
		}
	}
	return nil
}

// Ingredients holds a recipe's list of ingredients.
type Ingredients struct {
	Values []string
}

// MarshalJSON encodes the ingredients.
func (i *Ingredients) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Values)
}

// UnmarshalJSON decodes the ingredients according to the schema (https://schema.org/recipeInstructions).
func (i *Ingredients) UnmarshalJSON(data []byte) error {
	var xv []any
	err := json.Unmarshal(data, &xv)
	if err != nil {
		m := make(map[string][]any)
		err := json.Unmarshal(data, &m)
		if err != nil {
			return err
		}

		v, ok := m["Values"]
		if !ok {
			return fmt.Errorf("could not decode description %q", data)
		}
		xv = v
	}

	cases := []struct {
		old string
		new string
	}{
		{old: "  ", new: " "},
		{old: "\u00a0", new: " "},
		{old: "&frac12;", new: "½"},
		{old: "&frac34;", new: "¾"},
		{old: "&apos;", new: "'"},
		{old: "&nbsp;", new: ""},
		{old: "&#224;", new: "à"},
		{old: "&#8217;", new: "'"},
		{old: "&#339;", new: "œ"},
		{old: "&#233;", new: "é"},
		{old: "&#239;", new: "ï"},
	}

	for _, v := range xv {
		str := v.(string)
		if str == " " {
			continue
		}

		str = strings.TrimSpace(v.(string))
		for _, c := range cases {
			str = strings.ReplaceAll(str, c.old, c.new)
		}
		i.Values = append(i.Values, str)
	}
	return nil
}

// Instructions holds a recipe's list of instructions.
type Instructions struct {
	Values []string
}

// MarshalJSON encodes the list of instructions.
func (i *Instructions) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Values)
}

// UnmarshalJSON decodes the instructions according to the schema (https://schema.org/recipeInstructions).
func (i *Instructions) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		m := make(map[string][]any)
		err = json.Unmarshal(data, &m)
		if err != nil {
			return err
		}

		xv, ok := m["Values"]
		if !ok {
			return fmt.Errorf("could not decode description %q", data)
		}
		v = xv
	}

	switch x := v.(type) {
	case string:
		parts := strings.Split(strings.TrimSpace(x), "\n")
		for _, s := range parts {
			if s != "" {
				i.Values = append(i.Values, strings.TrimSpace(s))
			}
		}
	case []any:

		for _, part := range x {
			switch y := part.(type) {
			case string:
				i.Values = append(i.Values, strings.TrimSpace(y))
			case map[string]any:
				text, ok := y["text"]
				if ok {
					str := strings.TrimSuffix(text.(string), "\n")
					i.Values = append(i.Values, strings.TrimSpace(str))
					continue
				}

				parseSections(part, i)
			case []any:
				for _, sect := range y {
					parseSections(sect, i)
				}
			default:
				parseSections(y, i)
			}
		}
	case map[string]any:
		xv, ok := x["Values"]
		if ok {
			for _, s := range xv.([]any) {
				i.Values = append(i.Values, s.(string))
			}
		}
	}

	cases := map[string]string{
		"\u00a0":  " ",
		"\u2009":  "",
		"&apos;":  "'",
		"&nbsp;":  " ",
		"<br>":    " ",
		"&quot;":  `"`,
		"º":       "°",
		"&#233;":  "é",
		"&#232;":  "è",
		"&#8217;": "'",
		"&#224;":  "à",
		"&#239;":  "ï",
		"&#244;":  "ô",
		"&#x27;":  "'",
	}
	for i2, value := range i.Values {
		for old, newValue := range cases {
			value = strings.ReplaceAll(value, old, newValue)
		}

		value = strings.ReplaceAll(value, "  ", " ")
		i.Values[i2] = strings.TrimSpace(value)
	}

	return nil
}

type section struct {
	AtType string           `json:"@type"`
	Name   string           `json:"name"`
	Items  []map[string]any `json:"itemListElement"`
	Text   string           `json:"text"`
}

func parseSections(part any, instructions *Instructions) {
	b, err := json.Marshal(part)
	if err != nil {
		slog.Error("Marshal part failed", "error", err)
		return
	}

	var sect section
	err = json.Unmarshal(b, &sect)
	if err != nil {
		slog.Error("Marshal section failed", "error", err)
		return
	}

	if sect.Text != "" {
		instructions.Values = append(instructions.Values, sect.Text)
	}

	for _, item := range sect.Items {
		text, ok := item["text"]
		if ok {
			var str string
			switch x := text.(type) {
			case string:
				str = x
			case []any:
				xs := make([]string, 0, len(x))
				for _, v := range x {
					xs = append(xs, v.(string))
				}
				str = strings.Join(xs, ", ")
			default:
				continue
			}

			str = strings.TrimSuffix(str, "\n")
			str = strings.ReplaceAll(str, "\u00a0", "")
			str = strings.ReplaceAll(str, "\u2009", "")
			instructions.Values = append(instructions.Values, strings.TrimSpace(str))
		}
	}
}

// Yield holds a recipe's yield.
type Yield struct {
	Value int16
}

// MarshalJSON encodes the value of the yield.
func (y *Yield) MarshalJSON() ([]byte, error) {
	return json.Marshal(y.Value)
}

// UnmarshalJSON decodes the yield according to the schema (https://schema.org/recipeYield).
func (y *Yield) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		parts := strings.Split(v.(string), " ")
		for _, part := range parts {
			i, err := strconv.ParseInt(part, 10, 8)
			if err == nil {
				y.Value = int16(i)
				break
			}
		}
	case float64:
		y.Value = int16(x)
	case []any:
		if len(x) == 0 {
			break
		}

		split := strings.Split(x[0].(string), " ")
		if len(split) > 0 {
			v := split[0]
			i, err := strconv.ParseInt(v, 10, 16)
			if err == nil {
				y.Value = int16(i)
			}
		}
	case map[string]any:
		v, ok := x["Value"]
		if ok {
			y.Value = int16(v.(float64))
		}
	}
	return nil
}

// NutritionSchema is a representation of the nutrition schema (https://schema.org/NutritionInformation).
type NutritionSchema struct {
	Calories       string `json:"calories"`
	Carbohydrates  string `json:"carbohydrateContent"`
	Cholesterol    string `json:"cholesterolContent"`
	Fat            string `json:"fatContent"`
	Fiber          string `json:"fiberContent"`
	Protein        string `json:"proteinContent"`
	SaturatedFat   string `json:"saturatedFatContent"`
	Servings       string `json:"servingSize"`
	Sodium         string `json:"sodiumContent"`
	Sugar          string `json:"sugarContent"`
	TransFat       string `json:"transFatContent"`
	UnsaturatedFat string `json:"unsaturatedFatContent"`
}

func (n *NutritionSchema) nutrition() (Nutrition, error) {
	nutrition := Nutrition{
		Calories:           n.Calories,
		Cholesterol:        n.Cholesterol,
		Fiber:              n.Fiber,
		Protein:            n.Protein,
		SaturatedFat:       n.SaturatedFat,
		Sodium:             n.Sodium,
		Sugars:             n.Sugar,
		TotalCarbohydrates: n.Carbohydrates,
		TotalFat:           n.Fat,
		UnsaturatedFat:     n.UnsaturatedFat,
	}

	if n.Servings != "" {
		servings, err := strconv.ParseInt(n.Servings, 10, 16)
		if err != nil {
			return Nutrition{}, nil
		}

		nutrition.IsPerServing = true
		nutrition.Scale(1 / float64(servings))
	}

	return nutrition, nil
}

// UnmarshalJSON decodes the nutrition according to the schema.
func (n *NutritionSchema) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case []any:
		break
	case map[string]any:
		if val, ok := x["calories"].(string); ok {
			n.Calories = strings.TrimSpace(val)
		}

		if val, ok := x["carbohydrateContent"].(string); ok {
			n.Carbohydrates = val
		}

		if val, ok := x["cholesterolContent"].(string); ok {
			n.Cholesterol = val
		}

		if val, ok := x["fatContent"].(string); ok {
			n.Fat = val
		}

		if val, ok := x["fiberContent"].(string); ok {
			n.Fiber = val
		}

		if val, ok := x["proteinContent"].(string); ok {
			n.Protein = val
		}

		if val, ok := x["saturatedFatContent"].(string); ok {
			n.SaturatedFat = val
		}

		if val, ok := x["servingSize"].(string); ok {
			xs := strings.Split(val, " ")
			if len(xs) == 2 && len(xs[1]) < 4 {
				n.Servings = val
			} else {
				for _, s := range xs {
					_, err := strconv.Atoi(s)
					if err == nil {
						n.Servings = s
						break
					}
				}
			}
		}

		if val, ok := x["sodiumContent"].(string); ok {
			n.Sodium = val
		}

		if val, ok := x["sugarContent"].(string); ok {
			n.Sugar = val
		}

		if val, ok := x["transFatContent"].(string); ok {
			n.TransFat = val
		}

		if val, ok := x["unsaturatedFatContent"].(string); ok {
			n.UnsaturatedFat = val
		}
	}
	return nil
}

// Tools holds the list of tools used for a recipe.
type Tools struct {
	Values []string
}

// MarshalJSON encodes the list of tools.
func (t *Tools) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Values)
}

// UnmarshalJSON decodes the tools according to the schema (https://schema.org/tool).
func (t *Tools) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		t.Values = append(t.Values, x)
	case []map[string]any:
		for _, v := range x {
			t.Values = append(t.Values, v["name"].(string))
		}
	}
	return nil
}
