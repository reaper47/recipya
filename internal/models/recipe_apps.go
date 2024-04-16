package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"strings"
	"time"
)

type crouton struct {
	Name            string        `json:"name"`
	WebLink         string        `json:"webLink"`
	UUID            string        `json:"uuid"`
	CookingDuration int           `json:"cookingDuration"`
	Images          []string      `json:"images"`
	Serves          int           `json:"serves"`
	FolderIDs       []interface{} `json:"folderIDs"`
	DefaultScale    int           `json:"defaultScale"`
	Ingredients     []struct {
		Order      int `json:"order"`
		Ingredient struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"ingredient"`
		UUID string `json:"uuid"`
	} `json:"ingredients"`
	Steps []struct {
		Step      string `json:"step"`
		Order     int    `json:"order"`
		IsSection bool   `json:"isSection"`
		UUID      string `json:"uuid"`
	} `json:"steps"`
	Tags []struct {
		UUID  string `json:"uuid"`
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"tags"`
	IsPublicRecipe bool   `json:"isPublicRecipe"`
	SourceImage    string `json:"sourceImage"`
	Nutrition      string `json:"neutritionalInfo"`
	SourceName     string `json:"sourceName"`
	Duration       int    `json:"duration"`
}

// NewRecipeFromCrouton create a Recipe from the content of a Crouton file.
func NewRecipeFromCrouton(rc io.Reader, uploadImageFunc func(rc io.ReadCloser) (uuid.UUID, error)) Recipe {
	var c crouton
	err := json.NewDecoder(rc).Decode(&c)
	if err != nil {
		slog.Error("Failed to decode crouton recipe", "error", err)
		return Recipe{}
	}

	src := c.WebLink
	if src != "" {
		src = "Crouton"
	}

	var img uuid.UUID
	if len(c.Images) > 0 {
		decode, err := base64.StdEncoding.DecodeString(c.Images[0])
		if err == nil {
			img, err = uploadImageFunc(io.NopCloser(bytes.NewReader(decode)))
			if err != nil {
				slog.Error("Failed to upload Crouton image", "src", c.SourceImage, "file", c.Name, "error", err)
			}
		}
	}

	ingredients := make([]string, 0, len(c.Ingredients))
	for _, ing := range c.Ingredients {
		ingredients = append(ingredients, ing.Ingredient.Name)
	}

	instructions := make([]string, 0, len(c.Steps))
	for _, step := range c.Steps {
		instructions = append(instructions, step.Step)
	}

	keywords := make([]string, 0, len(c.Tags))
	for _, tag := range c.Tags {
		keywords = append(keywords, tag.Name)
	}

	category := "uncategorized"
	if len(keywords) > 0 {
		category = keywords[0]
	}

	var n Nutrition
	for _, s := range strings.Split(c.Nutrition, ",\n") {
		before, after, ok := strings.Cut(s, ":")
		if !ok {
			continue
		}

		after = strings.TrimSpace(after)

		switch strings.ToLower(before) {
		case "carbohydrates":
			n.TotalCarbohydrates = after
		case "calories":
			n.Calories = after
		case "fat":
			n.TotalFat = after
		case "sugar":
			n.Sugars = after
		case "cholesterol":
			n.Cholesterol = after
		case "fiber":
			n.Fiber = after
		case "saturated fat":
			n.SaturatedFat = after
		case "sodium":
			n.Sodium = after
		case "protein":
			n.Protein = after
		}
	}

	return Recipe{
		Category:     category,
		Description:  "Imported from Crouton",
		Image:        img,
		Ingredients:  ingredients,
		Instructions: instructions,
		Keywords:     keywords,
		Name:         c.Name,
		Nutrition:    n,
		Times: Times{
			Prep: time.Duration(c.Duration) * time.Minute,
			Cook: time.Duration(c.CookingDuration) * time.Minute,
		},
		Tools: nil,
		URL:   src,
		Yield: int16(c.Serves),
	}
}
