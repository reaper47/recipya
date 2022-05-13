package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeKwestiasmaku(root *html.Node) (models.RecipeSchema, error) {
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

	chName := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chName <- s
		}()

		node := getElement(root, "itemprop", "name")
		s = node.FirstChild.FirstChild.Data
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(
			root,
			"class",
			"field field-name-field-ilosc-porcji field-type-text field-label-hidden",
		)
		yield = findYield(strings.Split(strings.TrimSpace(node.FirstChild.Data), " "))
	}()

	chIngredients := make(chan models.Ingredients)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- models.Ingredients{Values: vals}
		}()

		node := getElement(
			root,
			"class",
			"field field-name-field-skladniki field-type-text-long field-label-hidden",
		)
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "ul" {
				for li := c.FirstChild; li != nil; li = li.NextSibling {
					if li.Type == html.ElementNode {
						vals = append(vals, strings.TrimSpace(li.FirstChild.Data))
					}
				}
				break
			}
		}
	}()

	chInstructions := make(chan models.Instructions)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- models.Instructions{Values: vals}
		}()

		node := getElement(
			root,
			"class",
			"field field-name-field-przygotowanie field-type-text-long field-label-above",
		)
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "ul" {
				for li := c.FirstChild; li != nil; li = li.NextSibling {
					if li.Type != html.ElementNode {
						continue
					}

					var xs []string
					for c2 := li.FirstChild; c2 != nil; c2 = c2.NextSibling {
						switch c2.Type {
						case html.ElementNode:
							xs = append(xs, c2.FirstChild.Data)
						case html.TextNode:
							xs = append(xs, c2.Data)
						}
					}

					v := strings.Join(xs, " ")
					v = strings.ReplaceAll(v, "\n", "")
					v = strings.ReplaceAll(v, "\t", "")
					vals = append(vals, v)
				}
				break
			}
		}
	}()

	return models.RecipeSchema{
		AtContext:    "https://schema.org",
		AtType:       models.SchemaType{Value: "Recipe"},
		Image:        models.Image{Value: <-chImage},
		Description:  models.Description{Value: <-getItemPropData(root, "description")},
		Name:         <-chName,
		Yield:        models.Yield{Value: <-chYield},
		Ingredients:  <-chIngredients,
		Instructions: <-chInstructions,
	}, nil
}
