package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeStreetKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	article := root.Find("article").First()
	name := article.Find("h1").First().Text()

	description, _ := root.Find("meta[name='description']").Attr("content")
	datePublished, _ := root.Find("meta[property='article:published_time").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time").Attr("content")

	yield := findYield(article.Find(".c-svgicon--servings").Next().Text())

	nodes := article.Find(".ingredients label")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = strings.TrimSpace(s.Text())
	})

	nodes = article.Find(".method-step")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.TrimSuffix(s.Text(), "\n")
		instructions[i] = v
	})

	image, _ := article.Find("img").First().Attr("src")

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Name:          name,
		Yield:         models.Yield{Value: yield},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
	}, nil
}
