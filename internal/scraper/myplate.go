package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeMyPlate(root *html.Node) (models.RecipeSchema, error) {
	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "class", "image-style-recipe-525-x-350-")
		s = getAttr(node, "src")
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(root, "class", "mp-recipe-full__detail mp-recipe-full__detail--yield grid-col-6 grid-row flex-column flex-align-center ")
		node = getElement(node, "class", "mp-recipe-full__detail--data")
		s := strings.TrimSpace(node.FirstChild.Data)
		s = strings.ReplaceAll(s, "\n", "")
		yield = findYield(strings.Split(s, " "))
	}()

	chCookTime := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chCookTime <- s
		}()

		node := getElement(root, "class", "mp-recipe-full__detail mp-recipe-full__detail--cook-time grid-col-6 grid-row flex-column flex-align-center")
		node = getElement(node, "class", "mp-recipe-full__detail--data")

		time := strings.TrimSpace(node.FirstChild.Data)
		time = strings.ReplaceAll(time, "\n", "")

		parts := strings.Split(time, " ")
		for _, part := range parts {
			i, err := strconv.Atoi(part)
			if err == nil {
				s = "PT" + strconv.Itoa(i) + "M"
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

		node := getElement(root, "class", "mp-recipe-full__description")
		s = node.FirstChild.NextSibling.FirstChild.Data
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		ul := getElement(root, "class", "field__items ingredients yyyyy")
		for li := ul.FirstChild; li != nil; li = li.NextSibling {
			if li.Type != html.ElementNode {
				continue
			}

			var xs []string
			for c := li.FirstChild; c != nil; c = c.NextSibling {
				switch c.Type {
				case html.ElementNode:
					v := strings.ReplaceAll(c.FirstChild.Data, "\n", "")

					xs = append(xs, strings.TrimSpace(v))
				case html.TextNode:
					v := strings.ReplaceAll(c.Data, "\n", "")
					xs = append(xs, strings.Join(strings.Fields(v), " "))
				}
			}
			vals = append(vals, strings.TrimSpace(strings.Join(xs, " ")))
		}
	}()

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- vals
		}()

		node := getElement(root, "class", "clearfix text-formatted field field--name-field-instructions field--type-text-long field--label-above")
		ol := getElement(node, "class", "field__item").FirstChild
		for li := ol.FirstChild; li != nil; li = li.NextSibling {
			if li.Type != html.ElementNode {
				continue
			}

			vals = append(vals, strings.TrimSpace(li.FirstChild.Data))
		}
	}()

	chNutrition := make(chan models.NutritionSchema)
	go func() {
		var nutrition models.NutritionSchema
		defer func() {
			_ = recover()
			chNutrition <- nutrition
		}()

		content := getElement(root, "class", "mp-recipe-full__nutrition-form")

		cases := []string{
			"Total Calories",
			"Total Fat",
			"Saturated Fat",
			"Cholesterol",
			"Sodium",
			"Carbohydrates",
			"Dietary Fiber",
			"Total Sugars",
			"Protein",
		}
		for _, c := range cases {
			tr := traverseAll(content, func(node *html.Node) bool {
				return node.Data == c
			})

			if len(tr) == 0 {
				continue
			}

			v := tr[len(tr)-1].Parent.NextSibling.NextSibling.FirstChild.Data
			switch c {
			case "Total Calories":
				nutrition.Calories = v
			case "Total Fat":
				nutrition.Fat = v
			case "Saturated Fat":
				nutrition.SaturatedFat = v
			case "Cholesterol":
				nutrition.Cholesterol = v
			case "Sodium":
				nutrition.Sodium = v
			case "Carbohydrates":
				nutrition.Carbohydrates = v
			case "Dietary Fiber":
				nutrition.Fiber = v
			case "Total Sugars":
				nutrition.Sugar = v
			case "Protein":
				nutrition.Protein = v
			}
		}

		servings := getElement(root, "class", "field field--name-field-recipe-serving-size field--type-string field--label-inline")
		servings = getElement(servings, "class", "field__item")
		v := strings.TrimSpace(servings.FirstChild.Data)
		v = strings.ReplaceAll(v, "\n", "")
		nutrition.Servings = v
	}()

	return models.RecipeSchema{
		AtContext:       "https://schema.org",
		AtType:          models.SchemaType{Value: "Recipe"},
		Name:            <-getElementData(root, "class", "field field--name-title field--type-string field--label-hidden"),
		Image:           models.Image{Value: <-chImage},
		Yield:           models.Yield{Value: <-chYield},
		CookTime:        <-chCookTime,
		Description:     models.Description{Value: <-chDescription},
		Ingredients:     models.Ingredients{Values: <-chIngredients},
		Instructions:    models.Instructions{Values: <-chInstructions},
		NutritionSchema: <-chNutrition,
	}, nil
}
