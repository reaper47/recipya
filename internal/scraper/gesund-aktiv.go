package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeGesundAktiv(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[name='description']").Attr("content")

	nodes := root.Find(".news-recipes-indgredients").Last().Find("ul li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		ingredients = append(ingredients, s)
	})

	nodes = root.Find(".news-recipes-cookingsteps").Last().Find("ol li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		Description:  models.Description{Value: description},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         strings.TrimSpace(root.Find("h1[itemprop='headline']").Text()),
	}, nil
}
