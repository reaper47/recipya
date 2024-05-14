package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeTastyKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	category := root.Find("a[rel='category tag']").First().Text()
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	name := root.Find("h1[itemprop='name']").Text()
	prepTime, _ := root.Find("time[itemprop='prepTime']").Attr("datetime")
	cookTime, _ := root.Find("time[itemprop='cookTime']").Attr("datetime")

	yieldStr, _ := root.Find("input[name='servings']").Attr("value")
	yield := findYield(yieldStr)

	nodes := root.Find("span[itemprop='ingredient']")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\u00a0", " ")
		ingredients[i] = v
	})

	nodes = root.Find("span[itemprop='instructions'] p")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = s.Text()
	})

	return models.RecipeSchema{
		Name:         name,
		Category:     models.Category{Value: category},
		Description:  models.Description{Value: description},
		Image:        models.Image{Value: image},
		PrepTime:     prepTime,
		CookTime:     cookTime,
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Yield:        models.Yield{Value: yield},
	}, nil
}
