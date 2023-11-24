package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeSallysblog(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[name='description']").Attr("content")

	prep := root.Find("p:contains('Zubereitungszeit')").Next().Text()
	split := strings.Split(prep, " ")
	isMin := strings.Contains(strings.ToLower(prep), "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			prep = split[i]
		}
	}

	nodes := root.Find(".recipe-description").Next().Find(".hidden").First().Prev().Find("div.text-lg")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		s = strings.TrimSpace(s)
		if s != "" {
			ingredients = append(ingredients, s)
		}
	})

	nodes = root.Find(".recipe-description div p")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		AtContext:    atContext,
		AtType:       models.SchemaType{Value: "Recipe"},
		Description:  models.Description{Value: description},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         strings.ToLower(root.Find("h1").First().Text()),
		PrepTime:     prep,
	}, nil
}
