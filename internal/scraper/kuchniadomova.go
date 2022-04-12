package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeKuchniadomova(root *html.Node) (models.RecipeSchema, error) {
	content := getElement(root, "class", "item-page")

	chImage := make(chan string)
	go func() {
		var v string
		defer func() {
			_ = recover()
			chImage <- v
		}()

		node := getElement(content, "id", "article-img-1")
		v = getAttr(node, "data-src")
		v = "https:" + v
	}()

	chDescription := make(chan string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chDescription <- strings.Join(vals, "\n\n")
		}()

		node := getElement(content, "itemprop", "description")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}
			vals = append(vals, c.FirstChild.Data)
		}
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		servings := getElement(content, "itemprop", "recipeYield")
		yield = findYield(strings.Split(servings.FirstChild.Data, " "))
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		node := getElement(content, "id", "recipe-ingredients")
		xn := traverseAll(node, func(node *html.Node) bool {
			return getAttr(node, "itemprop") == "recipeIngredient"
		})
		for _, n := range xn {
			vals = append(vals, n.FirstChild.Data)
		}
	}()

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- vals
		}()

		node := getElement(content, "itemprop", "recipeInstructions")
		ol := node.FirstChild
		for c := ol.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}
			vals = append(vals, c.FirstChild.Data)
		}
	}()

	return models.RecipeSchema{
		AtContext:    "https://schema.org",
		AtType:       models.SchemaType{Value: "Recipe"},
		Name:         <-getItemPropAttr(content, "name", "content"),
		Category:     models.Category{Value: <-getItemPropAttr(content, "recipeCategory", "content")},
		Image:        models.Image{Value: <-chImage},
		Description:  models.Description{Value: <-chDescription},
		Yield:        models.Yield{Value: <-chYield},
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
	}, nil
}
