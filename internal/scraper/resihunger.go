package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeReisHunger(root *html.Node) (models.RecipeSchema, error) {
	rs, err := findRecipeLdJSON(root)

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- vals
		}()

		node := getElement(root, "class", "nsection recipe-preparation")
		ol := traverseAll(node, func(node *html.Node) bool {
			return node.Data == "ol"
		})[0]
		for li := ol.FirstChild; li != nil; li = li.NextSibling {
			if li.Type != html.ElementNode {
				continue
			}

			var xs []string
			for c := li.FirstChild.NextSibling.FirstChild; c != nil; c = c.NextSibling {
				switch c.Type {
				case html.ElementNode:
					if c.Data == "br" {
						xs = append(xs, "\n")
					}
				case html.TextNode:
					v := strings.TrimSpace(c.Data)
					v = strings.ReplaceAll(v, "\n", "")
					xs = append(xs, v)
				}
			}
			vals = append(vals, strings.TrimSuffix(strings.Join(xs, ""), "\n"))
		}
	}()

	rs.Instructions = models.Instructions{Values: <-chInstructions}
	return rs, err
}
