package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeTastyKitchen(root *html.Node) (models.RecipeSchema, error) {
	chName := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chName <- s
		}()

		node := getElement(root, "class", "recipe-title")
		s = <-getItemPropData(node, "name")
	}()

	chCategory := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chCategory <- s
		}()

		node := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "rel") == "category tag"
		})[0]
		s = node.FirstChild.Data
	}()

	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "class", "the_recipe_image")
		s = getAttr(node, "src")
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "name") == "servings"
		})[0]
		i, err := strconv.Atoi(getAttr(node, "value"))
		if err == nil {
			yield = int16(i)
		}
	}()

	chDescription := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDescription <- s
		}()

		node := getElement(root, "itemprop", "summary")
		s = node.FirstChild.FirstChild.Data
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "itemprop") == "ingredient"
		})
		for _, n := range xn {
			if n.Type != html.ElementNode {
				continue
			}

			amount := <-getItemPropData(n, "amount")
			name := <-getItemPropData(n, "name")
			vals = append(vals, strings.TrimSpace(amount+" "+name))
		}
	}()

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- vals
		}()

		node := getElement(root, "itemprop", "instructions")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}
			vals = append(vals, c.FirstChild.Data)
		}
	}()

	return models.RecipeSchema{
		AtContext:    "https://schema.org",
		AtType:       models.SchemaType{Value: "Recipe"},
		Name:         <-chName,
		Category:     models.Category{Value: <-chCategory},
		Description:  models.Description{Value: <-chDescription},
		Image:        models.Image{Value: <-chImage},
		PrepTime:     <-getItemPropAttr(root, "prepTime", "datetime"),
		CookTime:     <-getItemPropAttr(root, "cookTime", "datetime"),
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
		Yield:        models.Yield{Value: <-chYield},
	}, nil
}
