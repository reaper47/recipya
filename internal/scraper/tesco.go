package scraper

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeTesco(root *goquery.Document) (models.RecipeSchema, error) {
	j := root.Find("script[type='application/ld+json']").First().Text()
	j = strings.ReplaceAll(j, "\n", "")
	var rs models.RecipeSchema
	if err := json.Unmarshal([]byte(j), &rs); err != nil {
		return models.RecipeSchema{}, err
	}

	rs.AtType.Value = "Recipe"
	return rs, nil
}
