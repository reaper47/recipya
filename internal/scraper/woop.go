package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeWoop(root *html.Node) (rs models.RecipeSchema, err error) {
	chDescription := make(chan string)
	go func() {
		var v string
		go func() {
			_ = recover()
			chDescription <- v
		}()

		node := getElement(root, "class", "product attribute recipe-origins")
		v = getElement(node, "class", "value").FirstChild.Data
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		go func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(root, "class", "product attribute serving-amount")
		value := getElement(node, "class", "value")

		parts := strings.Split(value.FirstChild.Data, " ")
		for _, part := range parts {
			i, err := strconv.Atoi(part)
			if err == nil {
				yield = int16(i)
				break
			}
		}
	}()

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		go func() {
			_ = recover()
			chInstructions <- vals
		}()

		node := getElement(root, "class", "product attribute cooking-instructions")
		value := getElement(node, "class", "value")
		for c := value.FirstChild.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode || c.FirstChild == nil {
				continue
			}

			h5 := strings.TrimSpace(c.FirstChild.NextSibling.FirstChild.Data)
			if h5 == "" {
				continue
			}

			p := strings.TrimSpace(c.LastChild.PrevSibling.FirstChild.Data)
			vals = append(vals, strings.TrimSpace(h5+"\n"+p))
		}
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		go func() {
			_ = recover()
			chIngredients <- vals
		}()

		node := getElement(root, "class", "product attribute ingredients")
		value := getElement(node, "class", "value")
		for c := value.FirstChild.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode || c.FirstChild == nil {
				continue
			}
			vals = append(vals, strings.TrimSpace(c.FirstChild.Data))
		}
	}()

	chNutrition := make(chan models.NutritionSchema)
	go func() {
		var m models.NutritionSchema
		go func() {
			_ = recover()
			chNutrition <- m
		}()

		node := getElement(root, "class", "product attribute nutritional-info")
		value := getElement(node, "class", "value")
		for c := value.FirstChild.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			parts := strings.Split(c.FirstChild.Data, ":")
			val := strings.TrimSpace(strings.Join(parts[1:], " "))
			switch parts[0] {
			case "Energy":
				m.Calories = val
			case "Protein":
				m.Protein = val
			case "Carbohydrate":
				m.Carbohydrates = val
			case "Fat":
				m.Fat = val
			}
		}
	}()

	return models.RecipeSchema{
		AtContext:       "https://schema.org",
		AtType:          models.SchemaType{Value: "Recipe"},
		Name:            <-getElementData(root, "data-ui-id", "page-title-wrapper"),
		Description:     models.Description{Value: <-chDescription},
		Ingredients:     models.Ingredients{Values: <-chIngredients},
		Instructions:    models.Instructions{Values: <-chInstructions},
		NutritionSchema: <-chNutrition,
		Yield:           models.Yield{Value: <-chYield},
	}, nil
}
