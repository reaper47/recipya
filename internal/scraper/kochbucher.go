package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeKochbucher(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = root.Find("h1[itemprop=headline]").Text()

	rs.Ingredients.Values = strings.Split(root.Find("p:contains('Zutaten')").First().Next().Text(), "\n")
	for i, ingredient := range rs.Ingredients.Values {
		rs.Ingredients.Values[i] = strings.TrimSpace(ingredient)
	}

	node := root.Find("p:contains('Zubereitung')")
	for {
		node = node.Next()
		if goquery.NodeName(node) != "p" {
			break
		}
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(node.Text()))
	}

	return rs, nil
}
