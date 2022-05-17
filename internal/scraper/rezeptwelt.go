package scraper

import (
	"strconv"
	"strings"
	"time"

	"github.com/reaper47/recipya/internal/constants"
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeRezeptwelt(root *html.Node) (models.RecipeSchema, error) {
	chName := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chName <- s
		}()

		node := getElement(root, "class", "breadcrumb")
		node = getElement(node.LastChild, "itemprop", "name")
		s = node.FirstChild.Data
	}()

	chDatePublished := make(chan string)
	chDateCreated := make(chan string)
	chDateModified := make(chan string)
	go func() {
		var (
			pub     string
			created string
			mod     string
		)
		defer func() {
			_ = recover()
			chDatePublished <- pub
			chDateCreated <- created
			chDateModified <- mod
		}()

		pub = <-getItemPropAttr(root, "datePublished", "content")

		xn := traverseAll(root, func(node *html.Node) bool {
			return strings.Contains(getAttr(node, "class"), "creation-date")
		})
		if len(xn) > 0 {
			parts := strings.Split(xn[0].FirstChild.Data, ":")
			for _, v := range parts {
				t, err := time.Parse("02.01.2006", strings.TrimSpace(v))
				if err == nil {
					created = t.Format(constants.BasicTimeLayout)
					break
				}
			}
		}

		xn = traverseAll(root, func(node *html.Node) bool {
			return strings.Contains(getAttr(node, "class"), "changed-date")
		})
		if len(xn) > 0 {
			parts := strings.Split(xn[0].FirstChild.Data, ":")
			for _, v := range parts {
				t, err := time.Parse("02.01.2006", strings.TrimSpace(v))
				if err == nil {
					mod = t.Format(constants.BasicTimeLayout)
					break
				}
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

		i, err := strconv.ParseInt(<-getItemPropData(root, "recipeYield"), 10, 16)
		if err == nil {
			yield = int16(i)
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
			return getAttr(node, "itemprop") == "recipeIngredient"
		})
		for _, n := range xn {
			if n.Type != html.ElementNode {
				continue
			}

			var xs []string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type != html.ElementNode {
					continue
				}
				xs = append(xs, c.FirstChild.Data)
			}
			vals = append(vals, strings.Join(xs, ""))
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
		span := getElement(node, "itemprop", "text")

		var s string
		nbr := 0
		for c := span.FirstChild; c != nil; c = c.NextSibling {
			switch c.Type {
			case html.ElementNode:
				if c.Data == "br" {
					if nbr > 1 {
						s = strings.ReplaceAll(strings.TrimSpace(s), "  ", " ")
						if s != "" {
							vals = append(vals, s)
							s = ""
							nbr = 0
							continue
						}
					} else {
						s += " "
					}
					nbr++
				} else {
					s += c.FirstChild.Data
				}
			case html.TextNode:
				s += strings.ReplaceAll(c.Data, "\n", "")
			}
		}
		vals = append(vals, strings.ReplaceAll(strings.TrimSpace(s), "  ", " "))
	}()

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		Name:          <-chName,
		DatePublished: <-chDatePublished,
		DateCreated:   <-chDateCreated,
		DateModified:  <-chDateModified,
		PrepTime:      <-getItemPropAttr(root, "performTime", "content"),
		Yield:         models.Yield{Value: <-chYield},
		Ingredients:   models.Ingredients{Values: <-chIngredients},
		Instructions:  models.Instructions{Values: <-chInstructions},
	}, nil
}
