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
	Url             string          `json:"url"`
}

// ToRecipe transforms the RecipeSchema to a Recipe.
func (m RecipeSchema) ToRecipe() (Recipe, error) {
	if m.AtType.Value != "Recipe" {
		return Recipe{}, fmt.Errorf("RecipeSchema %#v is not based on the Schema or the field is missing", m)
	}

	var category string
	if m.Category.Value == "" {
		category = "uncategorized"
	} else {
		category = strings.TrimSpace(strings.ToLower(m.Category.Value))
	}

	times, err := NewTimes(m.PrepTime, m.CookTime)
	if err != nil {
		return Recipe{}, err
	}

	nutrition, err := m.NutritionSchema.toNutrition()
	if err != nil {
		return Recipe{}, err
	}

	created := m.DateCreated
	if m.DatePublished != "" {
		created = m.DatePublished
	}

	var createdAt time.Time
	if created != "" {
		createdAt, err = time.Parse("2006-01-02", strings.Split(created, "T")[0])
		if err != nil {
			return Recipe{}, fmt.Errorf("could not parse createdAt date %s: '%s'", created, err)
		}
	}

	image, err := uuid.Parse(m.Image.Value)
	if err != nil {
		image = uuid.Nil
	}

	updatedAt := createdAt
	if m.DateModified != "" {
		updatedAt, err = time.Parse("2006-01-02", strings.Split(m.DateModified, "T")[0])
		if err != nil {
			return Recipe{}, fmt.Errorf("could not parse modifiedAt date %s: '%s'", m.DateModified, err)
		}
	}

	r := Recipe{
		ID:           0,
		Name:         m.Name,
		Description:  m.Description.Value,
		Image:        image,
		Url:          m.Url,
		Yield:        m.Yield.Value,
		Category:     category,
		Times:        times,
		Ingredients:  m.Ingredients,
		Nutrition:    nutrition,
		Instructions: m.Instructions.Values,
		Tools:        m.Tools.Values,
		Keywords:     strings.Split(m.Keywords.Values, ","),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	r.Normalize()
	return r, nil
}

// SchemaType holds the type of the schema. It should be "Recipe".
type SchemaType struct {
	Value string
}

// MarshalJSON encodes the schema's type.
func (m SchemaType) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

// UnmarshalJSON decodes the type of the schema.
// The type "Recipe" will be searched for if the data is an array.
func (m *SchemaType) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		m.Value = x
	case []interface{}:
		for _, t := range x {
			if t.(string) == "Recipe" {
				m.Value = "Recipe"
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
func (m Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

// UnmarshalJSON decodes the category according to the schema (https://schema.org/recipeCategory).
// The schema specifies that the expected value is of type Text. However, some
// websites send the category in an array, which explains the need for this function.
func (m *Category) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		m.Value = x
	case []interface{}:
		if len(x) > 0 {
			m.Value = x[0].(string)
		}
	}

	if m.Value != "" {
		m.Value = strings.ReplaceAll(m.Value, "&amp;", "&")
	}
	return nil
}

// CookingMethod holds a recipe's category.
type CookingMethod struct {
	Value string
}

// MarshalJSON encodes the cusine.
func (m CookingMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

// UnmarshalJSON decodes the cooking method according to the schema (https://schema.org/cookingMethod).
// The schema specifies that the expected value is of type Text. However, some
// websites send the cuisine in an array, which explains the need for this function.
func (m *CookingMethod) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		m.Value = x
	case []interface{}:
		if len(x) > 0 {
			m.Value = x[0].(string)
		}
	}
	return nil
}

// Cuisine holds a recipe's category.
type Cuisine struct {
	Value string
}

// MarshalJSON encodes the cusine.
func (m Cuisine) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

// UnmarshalJSON decodes the cuisine according to the schema (https://schema.org/recipeCuisine).
// The schema specifies that the expected value is of type Text. However, some
// websites send the cuisine in an array, which explains the need for this function.
func (m *Cuisine) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		m.Value = x
	case []interface{}:
		if len(x) > 0 {
			m.Value = x[0].(string)
		}
	}
	return nil
}

// Description holds a description.
type Description struct {
	Value string
}

// MarshalJSON encodes the description.
func (m Description) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

// UnmarshalJSON decodes the description field of a recipe.
func (m *Description) UnmarshalJSON(data []byte) error {
	var d string
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}

	d = strings.TrimSpace(d)
	d = strings.ReplaceAll(d, "\u00a0", "")
	d = strings.ReplaceAll(d, "&quot;", "")
	m.Value = d

	return nil
}

// Keywords holds keywords as a single string.
type Keywords struct {
	Values string
}

// MarshalJSON encodes the keywords.
func (m Keywords) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Values)
}

// UnmarshalJSON decodes the recipe's keywords according to the schema (https://schema.org/keywords).
// Some websites store the keywords in an array, which explains the need
// for a custom decoder.
func (m *Keywords) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		m.Values = strings.TrimSpace(x)
	case []interface{}:
		var xs []string
		for _, v := range x {
			xs = append(xs, v.(string))
		}
		m.Values = strings.TrimSpace(strings.Join(xs, ","))
	}
	return nil
}

// Image holds a recipe's image. The JSON fields correspond
type Image struct {
	Value string
}

// MarshalJSON encodes the image.
func (m Image) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

// UnmarshalJSON decodes the image according to the schema (https://schema.org/image).
func (m *Image) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		m.Value = x
	case []interface{}:
		if len(x) > 0 {
			switch y := x[0].(type) {
			case string:
				m.Value = y
			case map[string]interface{}:
				m.Value = y["url"].(string)
			}

		}
	case interface{}:
		url, ok := v.(map[string]interface{})["url"]
		if ok {
			m.Value = url.(string)
		}
	}
	return nil
}

// Ingredients holds a recipe's list of ingredients.
type Ingredients struct {
	Values []string
}

// MarshalJSON encodes the ingredients.
func (m Ingredients) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Values)
}

// UnmarshalJSON unmarshals the ingredients according to the schema (https://schema.org/recipeInstructions).
func (m *Ingredients) UnmarshalJSON(data []byte) error {
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
		m.Values = append(m.Values, str)
	}
	return nil
}

// Instructions holds a recipe's list of instructions.
type Instructions struct {
	Values []string
}

// MarshalJSON encodes the list of instructions.
func (m Instructions) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Values)
}

type section struct {
	AtType string                   `json:"@type"`
	Name   string                   `json:"name"`
	Items  []map[string]interface{} `json:"itemListElement"`
	Text   string                   `json:"text"`
}

// UnmarshalJSON unmarshals the instructions according to the schema (https://schema.org/recipeInstructions).
func (m *Instructions) UnmarshalJSON(data []byte) error {
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
				m.Values = append(m.Values, strings.TrimSpace(s))
			}
		}
	case []interface{}:
		for _, part := range x {
			switch y := part.(type) {
			case string:
				y = strings.ReplaceAll(y, "&quot;", "")
				m.Values = append(m.Values, strings.TrimSpace(y))
			case map[string]interface{}:
				text, ok := y["text"]
				if ok {
					str := strings.TrimSuffix(text.(string), "\n")
					str = strings.ReplaceAll(str, "\u00a0", "")
					str = strings.ReplaceAll(str, "\u2009", "")
					m.Values = append(m.Values, strings.TrimSpace(str))
					continue
				}

				parseSections(part, m)
			case []interface{}:
				for _, sect := range y {
					parseSections(sect, m)
				}
			default:
				parseSections(y, m)
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
func (m Yield) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

// UnmarshalJSON unmarshals the yield according to the schema (https://schema.org/recipeYield).
func (m *Yield) UnmarshalJSON(data []byte) error {
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
				m.Value = int16(i)
				break
			}
		}
	case float64:
		m.Value = int16(x)
	case []interface{}:
		if len(x) > 0 {
			v := strings.Split(x[0].(string), " ")[0]
			i, err := strconv.ParseInt(v, 10, 16)
			if err == nil {
				m.Value = int16(i)
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

func (m NutritionSchema) toNutrition() (Nutrition, error) {
	return Nutrition{
		Calories:           m.Calories,
		TotalCarbohydrates: m.Carbohydrates,
		Sugars:             m.Sugar,
		Protein:            m.Protein,
		TotalFat:           m.Fat,
		SaturatedFat:       m.SaturatedFat,
		Cholesterol:        m.Cholesterol,
		Sodium:             m.Sodium,
		Fiber:              m.Fiber,
	}, nil
}

// UnmarshalJSON unmarshals the nutrition according to the schema
func (m *NutritionSchema) UnmarshalJSON(data []byte) error {
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
			m.Calories = val
		}

		if val, ok := x["carbohydrateContent"].(string); ok {
			m.Carbohydrates = val
		}

		if val, ok := x["cholesterolContent"].(string); ok {
			m.Cholesterol = val
		}

		if val, ok := x["fatContent"].(string); ok {
			m.Fat = val
		}

		if val, ok := x["fiberContent"].(string); ok {
			m.Fiber = val
		}

		if val, ok := x["proteinContent"].(string); ok {
			m.Protein = val
		}

		if val, ok := x["saturatedFatContent"].(string); ok {
			m.SaturatedFat = val
		}

		if val, ok := x["servingSize"].(string); ok {
			m.Servings = val
		}

		if val, ok := x["sodiumContent"].(string); ok {
			m.Sodium = val
		}

		if val, ok := x["sugarContent"].(string); ok {
			m.Sugar = val
		}

		if val, ok := x["transFatContent"].(string); ok {
			m.TransFat = val
		}

		if val, ok := x["unsaturatedFatContent"].(string); ok {
			m.UnsaturatedFat = val
		}
	}
	return nil
}

// Tools holds the list of tools used for a recipe.
type Tools struct {
	Values []string
}

// MarshalJSON encodes the list of tools.
func (m Tools) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Values)
}

// UnmarshalJSON unmarshals the tools according to the schema (https://schema.org/tool).
func (m *Tools) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		m.Values = append(m.Values, x)
	case []map[string]interface{}:
		for _, v := range x {
			m.Values = append(m.Values, v["name"].(string))
		}
	}
	return nil
}
