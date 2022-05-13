package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeGlobo(root *html.Node) (rs models.RecipeSchema, err error) {
	chDescription := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDescription <- s
		}()

		div := getElement(root, "itemprop", "description")
		p := getElement(div, "class", "content-text__container")
		s = strings.TrimSpace(p.FirstChild.Data)
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
			vals = append(vals, n.FirstChild.Data)
		}
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		yieldStr := <-getItemPropAttr(root, "recipeYield", "content")
		i, err := strconv.ParseInt(yieldStr, 10, 16)
		if err == nil {
			yield = int16(i)
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
			n = getElement(n, "class", "recipeInstruction__text")
			vals = append(vals, strings.TrimSpace(n.FirstChild.Data))
		}
	}()

	return models.RecipeSchema{
		AtContext: "https://schema.org",
		AtType:    models.SchemaType{Value: "Recipe"},
		Name:      <-getItemPropAttr(root, "name", "content"),
		Image:     models.Image{Value: <-getItemPropAttr(root, "image", "content")},
		Yield:     models.Yield{Value: <-chYield},
		Keywords:  models.Keywords{Values: <-getItemPropAttr(root, "keywords", "content")},
		Category:  models.Category{Value: <-getItemPropAttr(root, "recipeCategory", "content")},
		CookingMethod: models.CookingMethod{
			Value: <-getItemPropAttr(root, "cookingMethod", "content"),
		},
		Cuisine:       models.Cuisine{Value: <-getItemPropAttr(root, "recipeCuisine", "content")},
		DatePublished: <-getItemPropAttr(root, "datePublished", "content"),
		DateModified:  getAttr(getElement(root, "itemprop", "dateModified"), "datetime"),
		Ingredients:   models.Ingredients{Values: <-chIngredients},
		Instructions:  models.Instructions{Values: <-chInstructions},
		Description:   models.Description{Value: <-chDescription},
	}, err
}
