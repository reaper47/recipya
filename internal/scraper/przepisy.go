package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapePrzepisy(root *html.Node) (models.RecipeSchema, error) {
	rs, err := findRecipeLdJSON(root)

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(root, "class", "person-count")
		yield = findYield(strings.Split(node.FirstChild.Data, " "))
	}()

	chingredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chingredients <- vals
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "class") == "ingredients-list-content-item"
		})
		for _, n := range xn {
			node := getElement(n, "class", "ingredient-name")
			node = traverseAll(node, func(node *html.Node) bool {
				return strings.Contains(getAttr(node, "class"), "text-bg-white")
			})[0]
			name := strings.TrimSpace(node.FirstChild.Data)

			node = getElement(n, "class", "quantity")
			quantity := strings.TrimSpace(node.FirstChild.FirstChild.Data)

			vals = append(vals, quantity+" "+name)
		}
	}()

	rs.Yield = models.Yield{Value: <-chYield}
	rs.Ingredients = models.Ingredients{Values: <-chingredients}
	return rs, err
}
