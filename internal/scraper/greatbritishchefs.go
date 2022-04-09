package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeGreatBritishChefs(root *html.Node) (models.RecipeSchema, error) {
	chName := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chName <- s
		}()

		node := getElement(root, "class", "Header__Title")
		node = getElement(node, "itemprop", "name")
		s = node.FirstChild.Data
	}()

	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "id", "head-media")
		s = getAttr(node, "src")
	}()

	chCategory := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chCategory <- s
		}()

		node := getElement(root, "itemprop", "recipeCategory")
		node = getElement(node, "class", "header-attribute-text text-capitalize")
		s = node.FirstChild.Data
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(root, "itemprop", "recipeYield")
		node = getElement(node, "class", "header-attribute-text")
		iStr := strings.ReplaceAll(node.FirstChild.Data, "\n", "")
		i, err := strconv.Atoi(strings.TrimSpace(iStr))
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
			var xs []string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type != html.ElementNode {
					continue
				}
				v := strings.ReplaceAll(c.FirstChild.Data, "\n", "")
				xs = append(xs, strings.TrimSpace(v))
			}
			vals = append(vals, strings.Join(xs, " "))
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
			return getAttr(node, "itemprop") == "recipeInstructions"
		})
		for _, n := range xn {
			v := strings.ReplaceAll(n.FirstChild.Data, "\n", "")
			vals = append(vals, strings.TrimSpace(v))
		}
	}()

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		DatePublished: <-getItemPropAttr(root, "datePublished", "content"),
		DateModified:  <-getItemPropAttr(root, "dateModified", "content"),
		Description:   models.Description{Value: <-getItemPropAttr(root, "description", "content")},
		Name:          <-chName,
		Image:         models.Image{Value: <-chImage},
		Category:      models.Category{Value: <-chCategory},
		Yield:         models.Yield{Value: <-chYield},
		CookTime:      <-getItemPropData(root, "cookTime"),
		PrepTime:      <-getItemPropData(root, "prepTime"),
		Ingredients:   models.Ingredients{Values: <-chIngredients},
		Instructions:  models.Instructions{Values: <-chInstructions},
	}, nil
}
