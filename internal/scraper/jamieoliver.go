package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeJamieOliver(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	if len(rs.Instructions.Values) == 1 {
		getInstructions(&rs, root.Find("ol li"))
	}

	return rs, nil
}
