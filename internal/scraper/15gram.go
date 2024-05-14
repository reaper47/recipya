package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrape15gram(root *goquery.Document) (models.RecipeSchema, error) {
	image, _ := root.Find("meta[property='og:image']").Attr("content")
	cook, _ := root.Find("meta[itemprop='cookTime']").Attr("content")
	prep, _ := root.Find("meta[itemprop='prepTime']").Attr("content")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, sel.Text())
	})

	nodes = root.Find("li[itemprop='recipeInstructions']")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, sel.Text())
	})

	nodes = root.Find("span[itemprop='keywords']")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(sel.Text()))
	})

	return models.RecipeSchema{
		CookTime:     cook,
		Description:  models.Description{Value: strings.TrimSpace(root.Find("p[itemprop='description']").Text())},
		Keywords:     models.Keywords{Values: strings.Join(keywords, ",")},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         root.Find("h1[itemprop='name']").Text(),
		PrepTime:     prep,
		Yield:        models.Yield{Value: findYield(root.Find("span[itemprop='recipeYield']").Text())},
	}, nil
}
