package scraper

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeLidlKochen(root *goquery.Document) (models.RecipeSchema, error) {
	jsonStr := root.Find("script[type='application/ld+json']").Text()
	jsonStr = strings.ReplaceAll(jsonStr, "\n", "")

	var xrs []models.RecipeSchema
	err := json.Unmarshal([]byte(jsonStr), &xrs)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	for _, xr := range xrs {
		if xr.AtType.Value == "Recipe" {
			return xr, nil
		}
	}

	return models.RecipeSchema{}, errors.New("no recipe found")
}
