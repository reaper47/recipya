package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapePanelinha(root *html.Node) (models.RecipeSchema, error) {
	rs, err := scrapeLdJSON(root)

	chIngredients := make(chan []string)
	if len(rs.Ingredients.Values) == 0 {
		go func() {
			var vals []string
			defer func() {
				_ = recover()
				chIngredients <- vals
			}()

			ul := traverseAll(root, func(node *html.Node) bool {
				return strings.Contains(node.Data, "Ingredientes")
			})[0].Parent.NextSibling.FirstChild

			for li := ul.FirstChild; li != nil; li = li.NextSibling {
				if li.Type != html.ElementNode {
					continue
				}
				vals = append(vals, li.FirstChild.Data)
			}
		}()
	} else {
		chIngredients <- rs.Ingredients.Values
	}

	chInstructions := make(chan []string)
	if len(rs.Ingredients.Values) == 0 {
		go func() {
			var vals []string
			defer func() {
				_ = recover()
				chInstructions <- vals
			}()

			ul := traverseAll(root, func(node *html.Node) bool {
				return strings.Contains(node.Data, "Modo de preparo")
			})[0].Parent.NextSibling.FirstChild

			for li := ul.FirstChild; li != nil; li = li.NextSibling {
				if li.Type != html.ElementNode {
					continue
				}
				vals = append(vals, strings.TrimSuffix(li.FirstChild.Data, "\u00a0"))
			}
		}()
	} else {
		chInstructions <- rs.Instructions.Values
	}

	rs.Ingredients = models.Ingredients{Values: <-chIngredients}
	rs.Instructions = models.Instructions{Values: <-chInstructions}
	return rs, err
}
