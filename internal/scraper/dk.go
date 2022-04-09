package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeDk(root *html.Node) (rs models.RecipeSchema, err error) {
	content := getElement(root, "itemtype", "http://schema.org/Recipe")
	description := strings.TrimSpace(<-getItemPropData(content, "description"))

	chYield := make(chan int16)
	go func() {
		var i int
		defer func() {
			_ = recover()
			chYield <- int16(i)
		}()

		yieldStr := getItemPropAttr(content, "recipeYield", "content")
		yield, err := strconv.Atoi(<-yieldStr)
		if err == nil {
			i = int(yield)
		}
	}()

	chIngredients := make(chan models.Ingredients)
	go func() {
		var v models.Ingredients
		defer func() {
			_ = recover()
			chIngredients <- v
		}()

		node := getElement(content, "itemprop", "recipeIngredient").Parent.Parent
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			switch c.Data {
			case "h3":
				v.Values = append(v.Values, "\n", c.FirstChild.Data)
			case "ul":
				for l := c.FirstChild; l != nil; l = l.NextSibling {
					amt := <-getElementData(l, "class", "recipe-ingredients__unit-amount")
					name := <-getElementData(l, "class", "recipe-ingredients__ingredient-instruction")
					v.Values = append(v.Values, strings.TrimSpace(amt)+" "+strings.TrimSpace(name))
				}
			}
		}
		v.Values = v.Values[1:]
	}()

	chInstructions := make(chan models.Instructions)
	go func() {
		var v models.Instructions
		defer func() {
			_ = recover()
			chInstructions <- v
		}()

		node := getElement(content, "itemprop", "recipeInstructions")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			switch c.Data {
			case "h3":
				s := c.FirstChild.Data
				s = strings.ReplaceAll(s, "\n", "")
				v.Values = append(v.Values, "\n", strings.TrimSpace(s))
			case "ol":
				for l := c.FirstChild; l != nil; l = l.NextSibling {
					s := l.FirstChild.Data
					s = strings.ReplaceAll(s, "\n", "")
					v.Values = append(v.Values, strings.TrimSpace(s))
				}
			}
		}
		v.Values = v.Values[1:]
	}()

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		Name:          <-getItemPropData(content, "name"),
		Image:         models.Image{Value: <-getItemPropAttr(content, "url", "content")},
		DatePublished: <-getItemPropAttr(content, "datePublished", "content"),
		Description:   models.Description{Value: description},
		Yield:         models.Yield{Value: <-chYield},
		Ingredients:   <-chIngredients,
		Instructions:  <-chInstructions,
	}, nil
}
