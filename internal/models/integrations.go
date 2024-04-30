package models

import (
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/utils/duration"
	"strconv"
	"strings"
	"time"
)

// NextcloudRecipes holds a Nextcloud recipe's metadata obtained from the Nextcloud Cookbook's /recipes endpoint.
type NextcloudRecipes struct {
	Category            string `json:"category"`
	DateCreated         string `json:"dateCreated"`
	DateModified        string `json:"dateModified"`
	ID                  int64  `json:"recipe_id"`
	ImagePlaceholderURL string `json:"imagePlaceholderUrl"`
	ImageURL            string `json:"imageUrl"`
	Keywords            string `json:"keywords"`
}

// PaprikaRecipe holds a Paprika recipe's JSON data.
type PaprikaRecipe struct {
	UID             string   `json:"uid"`
	Created         string   `json:"created"`
	Hash            string   `json:"hash"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Ingredients     string   `json:"ingredients"`
	Directions      string   `json:"directions"`
	Notes           string   `json:"notes"`
	NutritionalInfo string   `json:"nutritional_info"`
	PrepTime        string   `json:"prep_time"`
	CookTime        string   `json:"cook_time"`
	TotalTime       string   `json:"total_time"`
	Difficulty      string   `json:"difficulty"`
	Servings        string   `json:"servings"`
	Rating          int      `json:"rating"`
	Source          string   `json:"source"`
	SourceURL       string   `json:"source_url"`
	Photo           string   `json:"photo"`
	PhotoLarge      any      `json:"photo_large"`
	PhotoHash       string   `json:"photo_hash"`
	ImageURL        string   `json:"image_url"`
	Categories      []string `json:"categories"`
	PhotoData       string   `json:"photo_data"`
	Photos          []struct {
		Name     string `json:"name"`
		Filename string `json:"filename"`
		Hash     string `json:"hash"`
		Data     string `json:"data"`
	} `json:"photos"`
}

// Recipe converts a PaprikaRecipe to a Recipe.
func (p PaprikaRecipe) Recipe(image uuid.UUID) Recipe {
	category := "uncategorized"
	if len(p.Categories) > 0 {
		category = p.Categories[0]
	}

	var dateCreated time.Time
	dateCreated, _ = time.Parse(time.DateTime, p.Created)

	var yield int16
	parts := strings.Split(p.Servings, " ")
	for _, part := range parts {
		parsed, err := strconv.ParseInt(part, 10, 16)
		if err == nil {
			yield = int16(parsed)
			break
		}
	}

	source := p.SourceURL
	if source == "" {
		source = p.Source
	}

	description := p.Description
	if description == "" {
		description = "Imported from Paprika"
	}

	var images []uuid.UUID
	if image != uuid.Nil {
		images = append(images, image)
	}

	return Recipe{
		Category:     category,
		CreatedAt:    dateCreated,
		Description:  description,
		Images:       images,
		Ingredients:  strings.Split(p.Ingredients, "\n"),
		Instructions: strings.Split(p.Directions, "\n\n"),
		Keywords:     append(p.Categories, "paprika"),
		Name:         p.Name,
		Nutrition:    Nutrition{},
		Times: Times{
			Prep: duration.From(p.PrepTime),
			Cook: duration.From(p.CookTime),
		},
		Tools:     make([]string, 0),
		UpdatedAt: dateCreated,
		URL:       source,
		Yield:     yield,
	}
}
