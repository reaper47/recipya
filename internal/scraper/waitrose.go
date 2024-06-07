package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeWaitrose(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	for i, value := range rs.Instructions.Values {
		rs.Instructions.Values[i].Text = strings.ReplaceAll(value.Text, "\r", "")
	}

	return rs, nil
}
