package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeFineDiningLovers(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	nodes := root.Find(".field--name-field-srh-ingredients p")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find(".paragraph--type-srh-step")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		s = strings.ReplaceAll(s, "  ", "")
		s = strings.ReplaceAll(s, "\u00a0", "\n")
		s = strings.ReplaceAll(s, "\n\n\n", "\n\n")
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	return rs, nil
}
