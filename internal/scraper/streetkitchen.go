package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeStreetKitchen(root *html.Node) (models.RecipeSchema, error) {
	content := getElement(root, "id", "Main")

	chName := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chName <- s
		}()

		h1 := traverseAll(content, func(node *html.Node) bool {
			return node.Data == "h1"
		})[0]
		s = h1.FirstChild.Data
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(root, "class", "c-svgicon c-svgicon--servings ")
		i, err := strconv.ParseInt(node.NextSibling.NextSibling.FirstChild.Data, 10, 16)
		if err == nil {
			yield = int16(i)
		}
	}()

	chPrepTime := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chPrepTime <- s
		}()

		node := getElement(root, "class", "c-svgicon c-svgicon--prep-time ")
		t := strings.ReplaceAll(node.NextSibling.Data, "\n", "")
		t = strings.ReplaceAll(t, "\t", "")
		if t != "" {
			s = "PT" + t + "M"
		}
	}()

	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "class", "c-single-recipe--container")
		s = getAttr(node.FirstChild.NextSibling, "src")
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		node := traverseAll(root, func(node *html.Node) bool {
			return node.Data == "Ingredients"
		})[0]
		for c := node.Parent.NextSibling; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}
			vals = append(vals, c.FirstChild.Data)
		}
	}()

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- vals
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "class") == "c-single-recipe__method_copy"
		})
		for _, n := range xn {
			if n.Type != html.ElementNode {
				continue
			}
			vals = append(vals, n.FirstChild.FirstChild.Data)
		}
	}()

	return models.RecipeSchema{
		AtContext:    "https://schema.org",
		AtType:       models.SchemaType{Value: "Recipe"},
		Name:         <-chName,
		Yield:        models.Yield{Value: <-chYield},
		PrepTime:     <-chPrepTime,
		Image:        models.Image{Value: <-chImage},
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
	}, nil
}
