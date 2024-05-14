package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeKochbucher(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	ingredients := strings.Split(root.Find("p:contains('Zutaten')").First().Next().Text(), "\n")
	for i, ingredient := range ingredients {
		ingredients[i] = strings.TrimSpace(ingredient)
	}

	var instructions []string
	node := root.Find("p:contains('Zubereitung')")
	for {
		node = node.Next()
		if goquery.NodeName(node) != "p" {
			break
		}
		instructions = append(instructions, strings.TrimSpace(node.Text()))
	}

	return models.RecipeSchema{
		DateCreated:   "",
		DateModified:  "",
		DatePublished: "",
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          root.Find("h1[itemprop='headline']").Text(),
	}, nil
}
