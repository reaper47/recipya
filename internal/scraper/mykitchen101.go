package scraper

import (
	"strings"
	"time"

	"github.com/reaper47/recipya/internal/constants"
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeMyKitchen101(root *html.Node) (models.RecipeSchema, error) {
	chDescription := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDescription <- s
		}()

		node := getElement(root, "class", "post-content entry-content")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			if c.FirstChild.Type == html.TextNode {
				s = c.FirstChild.Data
				break
			}
		}
	}()

	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "class", "post-content entry-content")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			if c.FirstChild.Data == "img" {
				s = getAttr(c.FirstChild, "src")
				break
			}
		}
	}()

	chDatePublished := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDatePublished <- s
		}()

		node := getElement(root, "class", "updated")
		t, _ := time.Parse("Jan 02, 2006", node.FirstChild.Data)
		s = t.Format(constants.BasicTimeLayout)
	}()

	chIngredients := make(chan models.Ingredients)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- models.Ingredients{Values: vals}
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return strings.Contains(node.Data, "材料") || strings.Contains(node.Data, "Ingredients")
		})
		for _, n := range xn {
			for {
				if n.Data != "p" {
					n = n.Parent
					continue
				}
				break
			}
			ul := n.NextSibling.NextSibling
			for li := ul.FirstChild; li != nil; li = li.NextSibling {
				if li.Type != html.ElementNode {
					continue
				}
				vals = append(vals, li.FirstChild.Data)
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

		xn := traverseAll(root, func(node *html.Node) bool {
			return strings.Contains(node.Data, "做法") || strings.Contains(node.Data, "Directions")
		})
		for _, n := range xn {
			for {
				if n.Data != "p" {
					n = n.Parent
					continue
				}
				break
			}

			for c := n.NextSibling.NextSibling; c != nil; c = c.NextSibling {
				if c.Type != html.ElementNode {
					continue
				}

				if c.Data == "script" || c.Data == "ins" || c.Data == "div" ||
					c.FirstChild.Data == "img" {
					continue
				}

				var xs []string
				for p := c.FirstChild; p != nil; p = p.NextSibling {
					switch p.Type {
					case html.ElementNode:
						xs = append(xs, p.FirstChild.Data)
					case html.TextNode:
						xs = append(xs, p.Data)
					}
				}

				vals = append(vals, strings.Join(xs, ""))
			}
		}
	}()

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		DatePublished: <-chDatePublished,
		Name:          <-getElementData(root, "class", "entry-title"),
		Description:   models.Description{Value: <-chDescription},
		Image:         models.Image{Value: <-chImage},
		Ingredients:   <-chIngredients,
		Instructions:  <-chInstructions,
	}, nil
}
