package models

import (
	"encoding/json"
	"fmt"
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
	Yield           Yield           `json:"recipeYield"`
	URL             string          `json:"url"`
}

// ToRecipe transforms the RecipeSchema to a Recipe.
func (r RecipeSchema) ToRecipe() (Recipe, error) {
	if r.AtType.Value != "Recipe" {
		return Recipe{}, fmt.Errorf("RecipeSchema %#v is not based on the Schema or the field is missing", r)
	}

	var category string
	if r.Category.Value == "" {
		category = "uncategorized"
	} else {
		category = strings.TrimSpace(strings.ToLower(r.Category.Value))
	}

	times, err := NewTimes(r.PrepTime, r.CookTime)
	if err != nil {
		return Recipe{}, err
	}

	nutrition, err := r.NutritionSchema.toNutrition()
	if err != nil {
		return Recipe{}, err
	}

	created := r.DateCreated
	if r.DatePublished != "" {
		created = r.DatePublished
	}

	var createdAt time.Time
	if created != "" {
		createdAt, err = time.Parse("2006-01-02", strings.Split(created, "T")[0])
		if err != nil {
			return Recipe{}, fmt.Errorf("could not parse createdAt date %s: '%s'", created, err)
		}
	}

	image, err := uuid.Parse(r.Image.Value)
	if err != nil {
		image = uuid.Nil
	}

	updatedAt := createdAt
	if r.DateModified != "" {
		updatedAt, err = time.Parse("2006-01-02", strings.Split(r.DateModified, "T")[0])
		if err != nil {
			return Recipe{}, fmt.Errorf("could not parse modifiedAt date %s: '%s'", r.DateModified, err)
		}
	}

	recipe := Recipe{
		ID:           0,
		Name:         r.Name,
		Description:  r.Description.Value,
		Image:        image,
		URL:          r.URL,
		Yield:        r.Yield.Value,
		Category:     category,
		Times:        times,
		Ingredients:  r.Ingredients,
		Nutrition:    nutrition,
		Instructions: r.Instructions.Values,
		Tools:        r.Tools.Values,
		Keywords:     strings.Split(r.Keywords.Values, ","),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	recipe.Normalize()
	return recipe, nil
}

// SchemaType holds the type of the schema. It should be "Recipe".
type SchemaType struct {
	Value string
}

// MarshalJSON encodes the schema's type.
func (s SchemaType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}

// UnmarshalJSON decodes the type of the schema.
// The type "Recipe" will be searched for if the data is an array.
func (s *SchemaType) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		s.Value = x
	case []interface{}:
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
func (c Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

// UnmarshalJSON decodes the category according to the schema (https://schema.org/recipeCategory).
// The schema specifies that the expected value is of type Text. However, some
// websites send the category in an array, which explains the need for this function.
func (c *Category) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		c.Value = x
	case []interface{}:
		if len(x) > 0 {
			c.Value = x[0].(string)
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

// MarshalJSON encodes the cusine.
func (c CookingMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

// UnmarshalJSON decodes the cooking method according to the schema (https://schema.org/cookingMethod).
// The schema specifies that the expected value is of type Text. However, some
// websites send the cuisine in an array, which explains the need for this function.
func (c *CookingMethod) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		c.Value = x
	case []interface{}:
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

// MarshalJSON encodes the cusine.
func (c Cuisine) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

// UnmarshalJSON decodes the cuisine according to the schema (https://schema.org/recipeCuisine).
// The schema specifies that the expected value is of type Text. However, some
// websites send the cuisine in an array, which explains the need for this function.
func (c *Cuisine) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		c.Value = x
	case []interface{}:
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
func (d Description) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}

// UnmarshalJSON decodes the description field of a recipe.
func (d *Description) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\u00a0", "")
	s = strings.ReplaceAll(s, "&quot;", "")
	d.Value = s

	return nil
}

// Keywords holds keywords as a single string.
type Keywords struct {
	Values string
}

// MarshalJSON encodes the keywords.
func (k Keywords) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.Values)
}

// UnmarshalJSON decodes the recipe's keywords according to the schema (https://schema.org/keywords).
// Some websites store the keywords in an array, which explains the need
// for a custom decoder.
func (k *Keywords) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		k.Values = strings.TrimSpace(x)
	case []interface{}:
		var xs []string
		for _, v := range x {
			xs = append(xs, v.(string))
		}
		k.Values = strings.TrimSpace(strings.Join(xs, ","))
	}
	return nil
}

// Image holds a recipe's image. The JSON fields correspond
type Image struct {
	Value string
}

// MarshalJSON encodes the image.
func (i Image) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Value)
}

// UnmarshalJSON decodes the image according to the schema (https://schema.org/image).
func (i *Image) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		i.Value = x
	case []interface{}:
		if len(x) > 0 {
			switch y := x[0].(type) {
			case string:
				i.Value = y
			case map[string]interface{}:
				i.Value = y["url"].(string)
			}

		}
	case interface{}:
		url, ok := v.(map[string]interface{})["url"]
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
func (i Ingredients) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Values)
}

// UnmarshalJSON unmarshals the ingredients according to the schema (https://schema.org/recipeInstructions).
func (i *Ingredients) UnmarshalJSON(data []byte) error {
	var xv []interface{}
	err := json.Unmarshal(data, &xv)
	if err != nil {
		return err
	}

	cases := []struct {
		old string
		new string
	}{
		{old: "  ", new: " "},
		{old: "\u00a0", new: " "},
		{old: "&frac12;", new: "½"},
		{old: "&frac34;", new: "¾"},
	}

	for _, v := range xv {
		str := strings.TrimSpace(v.(string))
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
func (i Instructions) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Values)
}

type section struct {
	AtType string                   `json:"@type"`
	Name   string                   `json:"name"`
	Items  []map[string]interface{} `json:"itemListElement"`
	Text   string                   `json:"text"`
}

// UnmarshalJSON unmarshals the instructions according to the schema (https://schema.org/recipeInstructions).
func (i *Instructions) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		x = strings.ReplaceAll(x, "\u00a0", " ")
		parts := strings.Split(strings.TrimSpace(x), "\n")
		for _, s := range parts {
			if s != "" {
				i.Values = append(i.Values, strings.TrimSpace(s))
			}
		}
	case []interface{}:
		for _, part := range x {
			switch y := part.(type) {
			case string:
				y = strings.ReplaceAll(y, "&quot;", "")
				i.Values = append(i.Values, strings.TrimSpace(y))
			case map[string]interface{}:
				text, ok := y["text"]
				if ok {
					str := strings.TrimSuffix(text.(string), "\n")
					str = strings.ReplaceAll(str, "\u00a0", "")
					str = strings.ReplaceAll(str, "\u2009", "")
					i.Values = append(i.Values, strings.TrimSpace(str))
					continue
				}

				parseSections(part, i)
			case []interface{}:
				for _, sect := range y {
					parseSections(sect, i)
				}
			default:
				parseSections(y, i)
			}
		}
	}
	return nil
}

func parseSections(part interface{}, instructions *Instructions) {
	var sect section
	b, err := json.Marshal(part)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(b, &sect)
	if err != nil {
		fmt.Println(err)
		return
	}

	if sect.Text != "" {
		instructions.Values = append(instructions.Values, sect.Text)
	}

	for _, item := range sect.Items {
		text, ok := item["text"]
		if ok {
			str := strings.TrimSuffix(text.(string), "\n")
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
func (y Yield) MarshalJSON() ([]byte, error) {
	return json.Marshal(y.Value)
}

// UnmarshalJSON unmarshals the yield according to the schema (https://schema.org/recipeYield).
func (y *Yield) UnmarshalJSON(data []byte) error {
	var v interface{}
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
	case []interface{}:
		if len(x) > 0 {
			v := strings.Split(x[0].(string), " ")[0]
			i, err := strconv.ParseInt(v, 10, 16)
			if err == nil {
				y.Value = int16(i)
			}
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

func (n NutritionSchema) toNutrition() (Nutrition, error) {
	return Nutrition{
		Calories:           n.Calories,
		TotalCarbohydrates: n.Carbohydrates,
		Sugars:             n.Sugar,
		Protein:            n.Protein,
		TotalFat:           n.Fat,
		SaturatedFat:       n.SaturatedFat,
		Cholesterol:        n.Cholesterol,
		Sodium:             n.Sodium,
		Fiber:              n.Fiber,
	}, nil
}

// UnmarshalJSON unmarshals the nutrition according to the schema
func (n *NutritionSchema) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case []interface{}:
		break
	case map[string]interface{}:
		if val, ok := x["calories"].(string); ok {
			n.Calories = val
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
			n.Servings = val
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
func (t Tools) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Values)
}

// UnmarshalJSON unmarshals the tools according to the schema (https://schema.org/tool).
func (t *Tools) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		t.Values = append(t.Values, x)
	case []map[string]interface{}:
		for _, v := range x {
			t.Values = append(t.Values, v["name"].(string))
		}
	}
	return nil
}
