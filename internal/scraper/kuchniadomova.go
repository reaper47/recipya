package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeKuchniadomova(root *goquery.Document) (models.RecipeSchema, error) {
	name := root.Find("h2[itemprop='name']").Text()
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\t", "")

	yieldStr := root.Find("p[itemprop='recipeYield']").Text()
	yieldStr = strings.ReplaceAll(yieldStr, "-", " ")
	yield := findYield(yieldStr)

	category, _ := root.Find("meta[itemprop='recipeCategory']").Attr("content")
	keywords, _ := root.Find("meta[name='keywords']").Attr("content")
	image, _ := root.Find("#article-img-1").Attr("data-src")

	description := root.Find("#recipe-description").Text()
	description = strings.TrimPrefix(description, "\n")
	description = strings.TrimSuffix(description, "\n")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = s.Text()
	})

	nodes = root.Find("#recipe-instructions li")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = s.Text()
	})

	return models.RecipeSchema{
		AtContext:    atContext,
		AtType:       models.SchemaType{Value: "Recipe"},
		Name:         name,
		Category:     models.Category{Value: category},
		Cuisine:      models.Cuisine{Value: root.Find("p[itemprop='recipeCuisine']").Text()},
		Keywords:     models.Keywords{Values: keywords},
		Image:        models.Image{Value: "https://kuchnia-domowa.pl" + image},
		Description:  models.Description{Value: description},
		Yield:        models.Yield{Value: yield},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
	}, nil
}
