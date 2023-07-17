package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeNutritionByNathalie(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")

	ingredients := make([]string, 0)
	start := root.Find("span:contains('Ingredients:')").Last().Parent().Parent().First()
	end := root.Find("span:contains('Directions:')").Last().Parent().Parent()
	nodes := start.NextUntilNodes(end.Nodes...)
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.TrimSpace(s.Text())
		if strings.HasPrefix(v, "Ingredients:") || v == "" {
			return
		}
		v = strings.TrimPrefix(v, "•")
		ingredients = append(ingredients, strings.TrimSpace(v))
	})

	nodes = root.Find("span:contains('Directions:')").Last().Parent().Parent().Next().Next().Find("li")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = strings.TrimSpace(s.Find("p").Text())
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		DatePublished: datePublished,
		DateModified:  dateModified,
		Description:   models.Description{Value: description},
		Name:          name,
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
	}, nil
}
