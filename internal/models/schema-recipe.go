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
	AtType          string          `json:"@type"`
	Category        string          `json:"recipeCategory"`
	CookTime        string          `json:"cookTime"`
	CookingMethod   string          `json:"cookingMethod"`
	Cuisine         string          `json:"recipeCuisine"`
	DateCreated     string          `json:"dateCreated"`
	DateModified    string          `json:"dateModified"`
	DatePublished   string          `json:"datePublished"`
	Description     string          `json:"description"`
	Keywords        string          `json:"keywords"`
	Image           string          `json:"image"`
	Ingredients     []string        `json:"recipeIngredient"`
	Instructions    instructions    `json:"recipeInstructions"`
	Name            string          `json:"name"`
	NutritionSchema NutritionSchema `json:"nutrition"`
	PrepTime        string          `json:"prepTime"`
	Tools           tools           `json:"tool"`
	Yield           yield           `json:"recipeYield"`
	Url             string          `json:"url"`
}

// ToRecipe transforms the RecipeSchema to a Recipe.
func (m RecipeSchema) ToRecipe() (Recipe, error) {
	if m.AtType != "Recipe" {
		return Recipe{}, fmt.Errorf("RecipeSchema %#v is not based on the Schema or the field is missing", m)
	}

	var category string
	if m.Category == "" {
		category = "uncategorized"
	} else {
		category = strings.TrimSpace(strings.ToLower(m.Category))
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

	image, err := uuid.Parse(m.Image)
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

	return Recipe{
		ID:           0,
		Name:         m.Name,
		Description:  m.Description,
		Image:        image,
		Url:          m.Url,
		Yield:        m.Yield.Value,
		Category:     category,
		Times:        times,
		Ingredients:  m.Ingredients,
		Nutrition:    nutrition,
		Instructions: m.Instructions.Values,
		Tools:        m.Tools.Values,
		Keywords:     strings.Split(m.Keywords, ","),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

type instructions struct {
	Values []string
}

// MarshalJSON encodes the list of instructions.
func (m instructions) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Values)
}

// UnmarshalJSON unmarshals the instructions according to the schema (https://schema.org/recipeInstructions).
func (m *instructions) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		parts := strings.Split(v.(string), ".")
		for _, part := range parts {
			val := strings.TrimSpace(part + ".")
			if val != "." {
				m.Values = append(m.Values, val)
			}
		}
	case []any:
		for _, part := range x {
			m.Values = append(m.Values, part.(string))
		}
	}
	return nil
}

type yield struct {
	Value int16
}

// MarshalJSON encodes the value of the yield.
func (m yield) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Value)
}

// UnmarshalJSON unmarshals the yield according to the schema (https://schema.org/recipeYield).
func (m *yield) UnmarshalJSON(data []byte) error {
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
				m.Value = int16(i)
				break
			}
		}
	case float64:
		m.Value = int16(x)
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

type tools struct {
	Values []string
}

// MarshalJSON encodes the list of tools.
func (m tools) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Values)
}

// UnmarshalJSON unmarshals the tools according to the schema (https://schema.org/tool).
func (m *tools) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		m.Values = append(m.Values, x)
	case []map[string]any:
		for _, v := range x {
			m.Values = append(m.Values, v["name"].(string))
		}
	}
	return nil
}
