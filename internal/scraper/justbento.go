package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeJustbento(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	dateModified, _ := root.Find("meta[property='og:updated_time']").Attr("content")
	datePublished, _ := root.Find("meta[property='og:published_time']").Attr("content")

	category := root.Find("nav.breadcrumb").Find("a:contains('Recipe collection:')").Text()
	_, after, found := strings.Cut(category, ":")
	if found {
		category = strings.TrimSpace(after)
	} else {
		category = ""
	}

	nodes := root.Find(".field-name-body li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	var instructions []string
	nodes = root.Find(".field-name-body ul").Last()
	for {
		nodes = nodes.Next()
		if nodes.Nodes == nil {
			break
		}

		if goquery.NodeName(nodes) != "p" {
			continue
		}

		s := strings.TrimSpace(nodes.Text())
		if s != "" {
			instructions = append(instructions, s)
		}
	}

	image, _ := root.Find(".field-name-body img").First().Attr("src")

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: category},
		Cuisine:       models.Cuisine{Value: "Japanese"},
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		Yield:         models.Yield{Value: findYield(root.Find("*:contains('portions')").Text())},
	}, nil
}
