package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeHomeChef(root *html.Node) (models.RecipeSchema, error) {
	content := getElement(root, "id", "mainContent")

	chName := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chName <- s
		}()

		node := getElement(content, "itemprop", "name")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				s = c.FirstChild.Data
				break
			}
		}
	}()

	chDescription := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDescription <- s
		}()

		node := getElement(content, "itemprop", "description")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				s = strings.TrimSpace(c.FirstChild.Data)
				break
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

		servings := <-getItemPropAttr(content, "recipeYield", "content")
		yield = findYield(strings.Split(servings, " "))
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
				if c.Type == html.TextNode {
					v := strings.ReplaceAll(c.Data, "\n", " ")
					v = strings.ReplaceAll(v, "  ", " ")
					xs = append(xs, strings.TrimSpace(v))
				}
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

		node := getElement(content, "itemprop", "recipeInstructions")
		xn := traverseAll(node, func(node *html.Node) bool {
			return getAttr(node, "itemprop") == "itemListElement"
		})
		for _, n := range xn {
			span := getElement(n, "itemprop", "description")

			var xs []string
			for c := span.FirstChild; c != nil; c = c.NextSibling {
				if c.FirstChild == nil {
					continue
				}

				for c := c.FirstChild; c != nil; c = c.NextSibling {
					switch c.Type {
					case html.ElementNode:
						xs = append(xs, strings.TrimSpace(c.FirstChild.Data))
					case html.TextNode:
						xs = append(xs, strings.TrimSpace(c.Data))
					}
				}
			}
			vals = append(vals, strings.Join(xs, " "))
		}
	}()

	return models.RecipeSchema{
		AtContext:   "https://schema.org",
		AtType:      models.SchemaType{Value: "Recipe"},
		Image:       models.Image{Value: <-getItemPropAttr(content, "image", "content")},
		Name:        <-chName,
		Description: models.Description{Value: <-chDescription},
		Yield:       models.Yield{Value: <-chYield},
		NutritionSchema: models.NutritionSchema{
			Calories:       <-getItemPropData(content, "calories"),
			Carbohydrates:  <-getItemPropData(content, "carbohydrateContent"),
			Sugar:          <-getItemPropData(content, "sugarContent"),
			Protein:        <-getItemPropData(content, "proteinContent"),
			Fat:            <-getItemPropData(content, "fatContent"),
			SaturatedFat:   <-getItemPropData(content, "saturatedFatContent"),
			Cholesterol:    <-getItemPropData(content, "cholesterolContent"),
			Sodium:         <-getItemPropData(content, "sodiumContent"),
			Fiber:          <-getItemPropData(content, "fiberContent"),
			TransFat:       <-getItemPropData(content, "transFatContent"),
			UnsaturatedFat: <-getItemPropData(content, "unsaturatedFatContent"),
		},
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
	}, nil
}
