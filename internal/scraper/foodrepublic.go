package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeFoodRepublic(root *html.Node) (rs models.RecipeSchema, err error) {
	content := getElement(root, "id", "content")

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		yieldStr := <-getItemPropData(content, "recipeYield")
		i, err := strconv.ParseInt(strings.TrimSpace(yieldStr), 10, 16)
		if err == nil {
			yield = int16(i)
		}
	}()

	chIngredients := make(chan models.Ingredients)
	go func() {
		var ingredients models.Ingredients
		defer func() {
			_ = recover()
			chIngredients <- ingredients
		}()

		for li := getElement(content, "itemprop", "recipeIngredient"); li != nil; li = li.NextSibling {
			if li.Type != html.ElementNode {
				continue
			}

			var vals []string
			node := getElement(li, "class", "ingredient-quantity")
			if node != nil {
				vals = append(vals, strings.TrimSpace(node.FirstChild.Data))
			}

			node = getElement(li, "class", "ingredient-quantity-unit")
			if node != nil {
				vals = append(vals, strings.TrimSpace(node.FirstChild.Data))
			}
			vals = append(vals, <-getElementData(li, "class", "ingredient-label"))

			v := strings.Join(vals, " ")
			v = strings.ReplaceAll(v, "\n", "")
			v = strings.ReplaceAll(v, "\t", "")

			ingredients.Values = append(ingredients.Values, strings.TrimSpace(v))
		}
	}()

	chInstructions := make(chan models.Instructions)
	go func() {
		var instructions models.Instructions
		defer func() {
			_ = recover()
			chInstructions <- instructions
		}()

		ol := getElement(content, "itemprop", "recipeInstructions")
		for li := ol.FirstChild; li != nil; li = li.NextSibling {
			for c := li.FirstChild; c != nil; c = c.NextSibling {
				if c.Type != html.ElementNode {
					continue
				}
				instructions.Values = append(instructions.Values, c.FirstChild.Data)
			}
		}
	}()

	chKeywords := make(chan models.Keywords)
	go func() {
		var keywords models.Keywords
		defer func() {
			_ = recover()
			chKeywords <- keywords
		}()

		var vals []string
		div := getElement(content, "class", "tags")
		var ul *html.Node
		for c := div.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				ul = c
				break
			}
		}

		var li *html.Node
		for c := ul.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				li = c
				break
			}
		}

		for c := li; c != nil; c = c.NextSibling {
			a := getElement(c, "rel", "tag")
			if a == nil {
				continue
			}
			vals = append(vals, a.FirstChild.Data)
		}
		keywords.Values = strings.Join(vals, ",")
	}()

	return models.RecipeSchema{
		AtContext: "https://schema.org",
		AtType:    models.SchemaType{Value: "Recipe"},
		Name:      <-getItemPropAttr(content, "name", "content"),
		Description: models.Description{
			Value: <-getItemPropAttr(content, "description", "content"),
		},
		Image:         models.Image{Value: <-getItemPropAttr(content, "image", "content")},
		Yield:         models.Yield{Value: <-chYield},
		DatePublished: <-getItemPropAttr(content, "datePublished", "content"),
		Ingredients:   <-chIngredients,
		Instructions:  <-chInstructions,
		Keywords:      <-chKeywords,
	}, err
}
