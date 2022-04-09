package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapePurelyPope(root *html.Node) (models.RecipeSchema, error) {
	chPrepTime := make(chan string)
	chCookTime := make(chan string)
	go func() {
		var prep string
		var cook string
		defer func() {
			_ = recover()
			chPrepTime <- prep
			chCookTime <- cook
		}()

		node := getElement(root, "itemprop", "prepTime")
		prep = strings.ReplaceAll(getAttr(node, "datetime"), " ", "")

		node = getElement(root, "itemprop", "cookTime")
		cook = strings.ReplaceAll(getAttr(node, "datetime"), " ", "")
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		i, err := strconv.Atoi(<-getItemPropData(root, "recipeYield"))
		if err == nil {
			yield = int16(i)
		}
	}()

	chingredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chingredients <- vals
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "itemprop") == "recipeIngredient"
		})
		for _, n := range xn {
			v := strings.TrimSpace(n.FirstChild.Data)
			if v != "" {
				vals = append(vals, v)
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

		node := getElement(root, "itemprop", "recipeInstructions")
		for c := node.LastChild.PrevSibling.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}
			vals = append(vals, c.FirstChild.Data)
		}
	}()

	return models.RecipeSchema{
		AtContext:    "https://schema.org",
		AtType:       models.SchemaType{Value: "Recipe"},
		Name:         <-getItemPropData(root, "name"),
		PrepTime:     <-chPrepTime,
		CookTime:     <-chCookTime,
		Yield:        models.Yield{Value: <-chYield},
		Ingredients:  models.Ingredients{Values: <-chingredients},
		Instructions: models.Instructions{Values: <-chInstructions},
	}, nil
}
