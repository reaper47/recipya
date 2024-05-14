package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeFoodRepublic(root *goquery.Document) (models.RecipeSchema, error) {
	content := root.Find(".recipe-card")

	node := content.Find(".recipe-card-prep-time")
	prepTime := node.Find(".recipe-card-amount").Text()
	if node.Find(".recipe-card-unit").Text() == "minutes" {
		prepTime = "PT" + prepTime + "M"
	}

	node = content.Find(".recipe-card-cook-time")
	cookTime := node.Find(".recipe-card-amount").Text()
	if node.Find(".recipe-card-unit").Text() == "minutes" {
		cookTime = "PT" + cookTime + "M"
	}

	image, _ := content.Find(".recipe-card-image img").Attr("data-lazy-src")

	nodes := content.Find(".recipe-ingredients li")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = strings.ReplaceAll(s.Text(), "  ", " ")
	})

	nodes = content.Find(".recipe-directions li")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\u00a0", " ")
		instructions[i] = v
	})

	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")

	yieldStr := content.Find(".recipe-card-servings .recipe-card-amount").Text()
	yield, _ := strconv.ParseInt(yieldStr, 10, 16)

	name := content.Find(".recipe-card-title").Text()
	name = strings.TrimLeft(name, "\n")
	name = strings.TrimSpace(name)

	return models.RecipeSchema{
		CookTime:        cookTime,
		DateModified:    dateModified,
		DatePublished:   datePublished,
		Description:     models.Description{Value: content.Find(".recipe-card-description").Text()},
		Image:           models.Image{Value: image},
		Ingredients:     models.Ingredients{Values: ingredients},
		Instructions:    models.Instructions{Values: instructions},
		Name:            name,
		NutritionSchema: models.NutritionSchema{},
		PrepTime:        prepTime,
		Yield:           models.Yield{Value: int16(yield)},
	}, nil
}
