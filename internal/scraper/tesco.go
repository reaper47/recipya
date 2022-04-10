package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeTesco(root *html.Node) (models.RecipeSchema, error) {
	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "class", "recipe-detail__img")
		s = "https://realfood.tesco.com" + getAttr(node, "src")
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(root, "class", "recipe-detail__meta-item recipe-detail__meta-item_servings")
		yieldStr := strings.Fields(node.FirstChild.Data)[1]
		i, err := strconv.Atoi(yieldStr)
		if err == nil {
			yield = int16(i)
		}
	}()

	chNutrition := make(chan models.NutritionSchema)
	go func() {
		var calories string
		defer func() {
			_ = recover()
			chNutrition <- models.NutritionSchema{
				Calories: calories,
				Servings: "1",
			}
		}()

		node := getElement(root, "class", "recipe-detail__meta-item recipe-detail__meta-item_calories")
		parts := strings.Split(node.FirstChild.Data, "/")
		calories = strings.Split(parts[0], " ")[0]
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "class") == "recipe-detail__heading"
		})
		var node *html.Node
		for _, n := range xn {
			if n.Type == html.ElementNode && n.FirstChild.Data == "Ingredients" {
				node = n
				break
			}
		}
		ul := node.NextSibling.NextSibling
		for c := ul.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			v := c.FirstChild.Data
			v = strings.ReplaceAll(v, "\n", "")
			v = strings.Join(strings.Fields(v), " ")
			vals = append(vals, v)
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
			return getAttr(node, "class") == "recipe-detail__heading"
		})
		var node *html.Node
		for _, n := range xn {
			if n.Type == html.ElementNode && n.FirstChild.Data == "Method" {
				node = n
				break
			}
		}
		ol := node.NextSibling.NextSibling.FirstChild.NextSibling
		for c := ol.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			v := c.FirstChild.Data
			v = strings.ReplaceAll(v, "\n", "")
			v = strings.Join(strings.Fields(v), " ")
			vals = append(vals, v)
		}
	}()

	description := <-getElementData(root, "class", "recipe-detail__intro")
	description = strings.ReplaceAll(description, "\n", "")
	description = strings.Join(strings.Fields(description), " ")

	return models.RecipeSchema{
		AtContext:       "https://schema.org",
		AtType:          models.SchemaType{Value: "Recipe"},
		Name:            <-getElementData(root, "class", "recipe-detail__headline"),
		Description:     models.Description{Value: description},
		Yield:           models.Yield{Value: <-chYield},
		Image:           models.Image{Value: <-chImage},
		NutritionSchema: <-chNutrition,
		Ingredients:     models.Ingredients{Values: <-chIngredients},
		Instructions:    models.Instructions{Values: <-chInstructions},
	}, nil
}
