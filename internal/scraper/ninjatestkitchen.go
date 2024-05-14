package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeNinjatestkitchen(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[itemprop='name']").Attr("content")
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	image, _ := root.Find("meta[itemprop='image']").Attr("content")
	datePublished, _ := root.Find("meta[itemprop='datePublished']").Attr("content")
	keywords, _ := root.Find("meta[itemprop='keywords']").Attr("content")
	prepTime, _ := root.Find("meta[itemprop='prepTime']").Attr("content")
	recipeCategory, _ := root.Find("meta[itemprop='recipeCategory']").Attr("content")

	recipeIngredient, _ := root.Find("meta[itemprop='recipeIngredient']").Attr("content")
	ingredients := strings.Split(recipeIngredient, ",")
	for i, s := range ingredients {
		ingredients[i] = strings.TrimSpace(s)
	}

	nodes := root.Find(".single-method__method li p")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, sel.Text())
	})

	recipeYield, _ := root.Find("meta[itemprop='recipeYield']").Attr("content")

	return models.RecipeSchema{
		Category:      models.Category{Value: recipeCategory},
		DatePublished: datePublished,
		Description:   models.Description{Value: strings.TrimSpace(description)},
		Keywords:      models.Keywords{Values: keywords},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		PrepTime:      prepTime,
		Yield:         models.Yield{Value: findYield(recipeYield)},
	}, nil
}
