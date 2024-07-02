package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeWikiBooks(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	var name string
	split := strings.Split(root.Find("#firstHeading").Text(), ":")
	if len(split) > 1 {
		name = split[1]
	}
	rs.Name = name

	start := root.Find(".mw-parser-output").Children().First()
	if start.Nodes[0].Data == "section" {
		start = root.Find("#mf-section-0").Children().First()
	}
	nodes := start.NextUntil("h2")
	nodes = nodes.FilterFunction(func(_ int, s *goquery.Selection) bool {
		return s.Nodes[0].Data == "p"
	})
	description := nodes.Slice(1, nodes.Length()).Text()
	rs.Description.Value = strings.TrimSuffix(description, "\n")

	rs.Category.Value = root.Find("th:contains('Category')").Next().Text()

	image, _ := root.Find(".infobox-image img").First().Attr("src")
	if image != "" {
		image = "https:" + image
	}
	rs.Image.Value = image

	rs.Yield.Value = findYield(root.Find("th:contains('Servings')").Next().Text())

	start = root.Find("#Ingredients").Parent()
	end := root.Find("#Procedure").Parent()
	getIngredients(&rs, start.NextUntilSelection(end).Find("li"))
	getInstructions(&rs, root.Find("#Procedure").Parent().NextUntil("h2").Find("li"))

	return rs, nil
}
