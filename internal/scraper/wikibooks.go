package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeWikiBooks(root *goquery.Document) (models.RecipeSchema, error) {
	var name string
	split := strings.Split(root.Find("#firstHeading").Text(), ":")
	if len(split) > 1 {
		name = split[1]
	}

	start := root.Find(".mw-parser-output").Children().First()
	if start.Nodes[0].Data == "section" {
		start = root.Find("#mf-section-0").Children().First()
	}
	nodes := start.NextUntil("h2")
	nodes = nodes.FilterFunction(func(_ int, s *goquery.Selection) bool {
		return s.Nodes[0].Data == "p"
	})
	description := nodes.Slice(1, nodes.Length()).Text()
	description = strings.TrimSuffix(description, "\n")

	category := root.Find("th:contains('Category')").Next().Text()

	image, _ := root.Find(".infobox-image img").First().Attr("src")
	if image != "" {
		image = "https:" + image
	}

	yield := findYield(root.Find("th:contains('Servings')").Next().Text())

	start = root.Find("#Ingredients").Parent()
	end := root.Find("#Procedure").Parent()
	nodes = start.NextUntilSelection(end).Find("li")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = s.Text()
	})

	nodes = root.Find("#Procedure").Parent().NextUntil("h2").Find("li")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = s.Text()
	})

	return models.RecipeSchema{
		AtContext:    atContext,
		AtType:       models.SchemaType{Value: "Recipe"},
		Name:         name,
		Description:  models.Description{Value: description},
		Image:        models.Image{Value: image},
		Category:     models.Category{Value: category},
		Yield:        models.Yield{Value: yield},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
	}, nil
}
