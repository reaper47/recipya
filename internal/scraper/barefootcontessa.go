package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeBarefootcontessa(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseGraph(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	split := strings.Split(rs.Name, " |")
	if len(split) > 0 {
		rs.Name = strings.TrimSpace(split[0])
	}

	for i, value := range rs.Instructions.Values {
		value = strings.TrimPrefix(value, "<p>")
		value = strings.TrimSuffix(value, "</p>")
		rs.Instructions.Values[i] = value
	}

	if len(rs.Ingredients.Values) == 1 {
		rs.Ingredients.Values = strings.Split(rs.Ingredients.Values[0], "\n")
	}

	before, _, ok := strings.Cut(rs.Description.Value, "<p")
	if ok {
		rs.Description.Value = strings.TrimSpace(before)
	}

	return rs, nil
}
