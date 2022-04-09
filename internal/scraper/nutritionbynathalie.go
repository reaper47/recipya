package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeNutritionbynathalie(root *html.Node) (models.RecipeSchema, error) {
	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "id", "viewer-q926v").FirstChild.FirstChild.FirstChild
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			if c.Data == "img" {
				s = getAttr(c, "src")
			}
		}
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return strings.Contains(node.Data, "Ingredients")
		})
		for _, n := range xn {
			for c := n.Parent.Parent.NextSibling; c != nil; c = c.NextSibling {
				if c.Type != html.ElementNode || c.Data != "p" {
					continue
				}

				v := strings.TrimSpace(c.FirstChild.FirstChild.Data)
				if strings.Contains(v, "Directions") || v == "" {
					break
				}
				vals = append(vals, strings.TrimPrefix(v, "• "))
			}
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
			return strings.Contains(node.Data, "Directions")
		})
		if len(xn) == 0 {
			return
		}

		ol := xn[0].Parent.Parent.NextSibling.NextSibling
		for li := ol.FirstChild; li != nil; li = li.NextSibling {
			if li.Type != html.ElementNode {
				continue
			}

			v := strings.TrimSpace(li.FirstChild.FirstChild.Data)
			if strings.Contains(v, "Directions") || v == "" {
				break
			}
			vals = append(vals, strings.TrimPrefix(v, "• "))
		}
	}()

	return models.RecipeSchema{
		AtContext:    "https://schema.org",
		AtType:       models.SchemaType{Value: "Recipe"},
		Name:         <-getElementData(root, "class", "blog-post-title-font blog-post-title-color"),
		Image:        models.Image{Value: <-chImage},
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
	}, nil
}
