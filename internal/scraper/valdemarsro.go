package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeValdemarsro(root *html.Node) (models.RecipeSchema, error) {
	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "itemprop", "image")
		s = getAttr(node, "src")
	}()

	chDescription := make(chan string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chDescription <- strings.Join(vals, "\n\n")
		}()

		node := getElement(root, "itemprop", "description")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Data != "p" {
				continue
			}

			var xp []string
			for c := c.FirstChild; c != nil; c = c.NextSibling {
				switch c.Type {
				case html.ElementNode:
					if c.Data == "img" {
						continue
					}

					var xs []string
					for c := c.FirstChild; c != nil; c = c.NextSibling {
						switch c.Type {
						case html.ElementNode:
							xs = append(xs, c.FirstChild.Data)
						case html.TextNode:
							xs = append(xs, c.Data)
						}
					}
					xp = append(xp, strings.Join(xs, ""))
				case html.TextNode:
					xp = append(xp, c.Data)
				}
			}
			vals = append(vals, strings.Join(xp, ""))
		}
	}()

	chPrepTime := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chPrepTime <- s
		}()

		node := traverseAll(root, func(node *html.Node) bool {
			return node.Data == "Arbejdstid"
		})[0]
		parts := strings.Split(node.Parent.NextSibling.FirstChild.Data, " ")
		for _, v := range parts {
			i, err := strconv.Atoi(v)
			if err == nil {
				s = "PT" + strconv.Itoa(i) + "M"
			}
		}
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := traverseAll(root, func(node *html.Node) bool {
			return node.Data == "Antal"
		})[0]
		yield = findYield(strings.Split(node.Parent.NextSibling.FirstChild.Data, " "))
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "itemprop") == "recipeIngredient"
		})
		for _, n := range xn {
			if n.Type != html.ElementNode {
				continue
			}
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

		node := getElement(root, "itemprop", "recipeInstructions")
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
		Name:         <-getItemPropData(root, "headline"),
		Description:  models.Description{Value: <-chDescription},
		Image:        models.Image{Value: <-chImage},
		PrepTime:     <-chPrepTime,
		Yield:        models.Yield{Value: <-chYield},
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
	}, nil
}
