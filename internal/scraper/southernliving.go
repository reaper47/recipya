package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/regex"
	"golang.org/x/net/html"
)

func scrapeSouthernLiving(root *html.Node) (models.RecipeSchema, error) {
	chDescription := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDescription <- s
		}()

		node := getElement(root, "property", "og:description")
		s = getAttr(node, "content")
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return strings.HasPrefix(node.Data, "Serves")
		})
		yield = findYield(strings.Split(xn[0].Data, " "))
	}()

	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "class", "recipe-tout-image recipe-info-items-2")
		xn := traverseAll(node, func(node *html.Node) bool {
			return node.Data == "noscript"
		})
		s = regex.ImageSrc.FindString(xn[0].FirstChild.Data)
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "class") == "ingredients-item-name elementFont__body"
		})
		for _, n := range xn {
			q := strings.TrimSpace(n.FirstChild.FirstChild.Data)
			name := strings.TrimSpace(n.FirstChild.NextSibling.Data)
			vals = append(vals, q+" "+name)
		}
	}()

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- vals
		}()

		ul := getElement(root, "class", "instructions-section")
		for c := ul.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			p := getElement(c, "class", "paragraph")
			vals = append(vals, strings.TrimSpace(p.FirstChild.FirstChild.Data))
		}
	}()

	return models.RecipeSchema{
		AtContext: "https://schema.org",
		AtType:    models.SchemaType{Value: "Recipe"},
		Name: <-getElementData(
			root,
			"class",
			"headline heading-content elementFont__display",
		),
		Description:  models.Description{Value: <-chDescription},
		Image:        models.Image{Value: <-chImage},
		Yield:        models.Yield{Value: <-chYield},
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
	}, nil
}
