package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeKwestiasmaku(root *goquery.Document) (models.RecipeSchema, error) {
	image, _ := root.Find("meta[property='og:image']").Attr("content")
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	description := root.Find("span[itemprop='description']").Text()

	nodes := root.Find(".field-name-field-skladniki li")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\n", "")
		v = strings.ReplaceAll(v, "\t", "")
		ingredients[i] = v
	})

	nodes = root.Find(".field-name-field-przygotowanie li")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\n", "")
		v = strings.ReplaceAll(v, "\t", "")
		instructions[i] = v
	})

	yield := findYield(root.Find(".field-name-field-ilosc-porcji").Text())

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		Image:         models.Image{Value: image},
		DatePublished: datePublished,
		DateModified:  dateModified,
		Description:   models.Description{Value: description},
		Name:          name,
		Yield:         models.Yield{Value: yield},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
	}, nil
}
