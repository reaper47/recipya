package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeWikiBooks(root *html.Node) (models.RecipeSchema, error) {
	chName := make(chan string)
	go func() {
		var s string
		go func() {
			_ = recover()
			chName <- s
		}()

		node := getElement(root, "id", "firstHeading")
		parts := strings.Split(node.FirstChild.Data, ":")
		s = parts[len(parts)-1]
	}()

	chDescription := make(chan string)
	go func() {
		var s string
		go func() {
			_ = recover()
			chDescription <- s
		}()

		xp := traverseAll(root, func(node *html.Node) bool {
			return node.Data == "p"
		})
		for _, n := range xp[1:] {
			if n.Type != html.ElementNode {
				continue
			}

			if n.PrevSibling.PrevSibling.Data == "h2" {
				break
			}

			var vals []string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				switch c.Type {
				case html.ElementNode:
					vals = append(vals, c.FirstChild.Data)
				case html.TextNode:
					vals = append(vals, c.Data)
				}
			}
			s = strings.Join(vals, "")
			s = strings.ReplaceAll(s, "\n", "")
			break
		}
	}()

	chCategory := make(chan string)
	go func() {
		var s string
		go func() {
			_ = recover()
			chCategory <- s
		}()

		node := traverseAll(root, func(node *html.Node) bool {
			return node.Data == "Category"
		})[0]
		s = node.Parent.NextSibling.NextSibling.FirstChild.FirstChild.Data
	}()

	chImage := make(chan string)
	go func() {
		var s string
		go func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "class", "image").FirstChild
		s = "https:" + getAttr(node, "src")
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		go func() {
			_ = recover()
			chYield <- yield
		}()

		node := traverseAll(root, func(node *html.Node) bool {
			return node.Data == "Servings"
		})[0]
		i, err := strconv.ParseInt(node.Parent.NextSibling.NextSibling.FirstChild.Data, 10, 16)
		if err == nil {
			yield = int16(i)
		}
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		go func() {
			_ = recover()
			chIngredients <- vals
		}()

		node := getElement(root, "id", "Ingredients").Parent.NextSibling.NextSibling
		if node.Data == "h2" {
			node = node.PrevSibling.FirstChild
		}

		var xn []*html.Node
		for c := node; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			if c.Data == "h2" {
				break
			}

			if c.Data == "ul" {
				xn = append(xn, c)
			}
		}

		for _, ul := range xn {
			for li := ul.FirstChild; li != nil; li = li.NextSibling {
				if li.Type != html.ElementNode {
					continue
				}

				var xs []string
				for c := li.FirstChild; c != nil; c = c.NextSibling {
					switch c.Type {
					case html.ElementNode:
						xs = append(xs, c.FirstChild.Data)
					case html.TextNode:
						xs = append(xs, c.Data)
					}
				}
				v := strings.Join(xs, "")
				v = strings.ReplaceAll(v, "\n", "")
				vals = append(vals, strings.TrimSpace(v))
			}
		}
	}()

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		go func() {
			_ = recover()
			chInstructions <- vals
		}()

		node := getElement(root, "id", "Procedure").Parent.NextSibling.NextSibling
		switch node.Data {
		case "p":
			for c := node; c != nil; c = c.NextSibling {
				if c.Type != html.ElementNode {
					continue
				}

				if c.Data == "h2" {
					break
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
				v := strings.Join(xs, "")
				v = strings.ReplaceAll(v, "\n", "")
				vals = append(vals, strings.TrimSpace(v))
			}
		case "h2":
			node = node.PrevSibling
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				if c.Type != html.ElementNode {
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

				v := strings.Join(xs, "")
				v = strings.ReplaceAll(v, "\n", "")
				vals = append(vals, strings.TrimSpace(v))
			}
		}
	}()

	return models.RecipeSchema{
		AtContext:    "https://schema.org",
		AtType:       models.SchemaType{Value: "Recipe"},
		Name:         <-chName,
		Description:  models.Description{Value: <-chDescription},
		Image:        models.Image{Value: <-chImage},
		Category:     models.Category{Value: <-chCategory},
		Yield:        models.Yield{Value: <-chYield},
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
	}, nil
}
