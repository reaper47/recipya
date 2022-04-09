package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapePaniniHappy(root *html.Node) (models.RecipeSchema, error) {
	content := getElement(root, "class", "hrecipe")

	chPrep := make(chan string)
	chCook := make(chan string)
	go func() {
		var prep string
		var cook string
		defer func() {
			_ = recover()
			chPrep <- prep
			chCook <- cook
		}()

		node := getElement(content, "class", "preptime")
		parts := strings.Split(node.FirstChild.Data, " ")
		for _, part := range parts {
			i, err := strconv.Atoi(part)
			if err == nil {
				prep = "PT" + strconv.Itoa(i) + "M"
				break
			}
		}

		node = getElement(content, "class", "cooktime")
		parts = strings.Split(node.FirstChild.Data, " ")
		for _, part := range parts {
			i, err := strconv.Atoi(part)
			if err == nil {
				cook = "PT" + strconv.Itoa(i) + "M"
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

		node := getElement(content, "class", "yield")
		parts := strings.Split(node.FirstChild.Data, " ")
		for _, part := range parts {
			i, err := strconv.Atoi(part)
			if err == nil {
				yield = int16(i)
				break
			}
		}
	}()

	chName := make(chan string)
	chDescription := make(chan string)
	go func() {
		var name string
		var desc string
		defer func() {
			_ = recover()
			chName <- name
			chDescription <- desc
		}()

		for c := content.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			switch c.Data {
			case "h2":
				if c.FirstChild != nil {
					name = c.FirstChild.Data
				}
			case "p":
				if getAttr(c, "class") == "summary" && c.FirstChild != nil {
					var vals []string
					for c := c.FirstChild.FirstChild; c != nil; c = c.NextSibling {
						switch c.Type {
						case html.ElementNode:
							vals = append(vals, c.FirstChild.Data)
						case html.TextNode:
							vals = append(vals, c.Data)
						}
					}
					desc = strings.Join(vals, "")
				}
			}
		}
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		xn := traverseAll(content, func(node *html.Node) bool {
			return getAttr(node, "class") == "ingredient"
		})
		for _, n := range xn {
			amount := getElement(n, "class", "amount")
			name := getElement(n, "class", "name")
			vals = append(vals, amount.FirstChild.Data+" "+name.FirstChild.Data)
		}
	}()

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- vals
		}()

		xn := traverseAll(content, func(node *html.Node) bool {
			return getAttr(node, "class") == "instruction"
		})
		for _, n := range xn {
			vals = append(vals, n.FirstChild.Data)
		}
	}()

	return models.RecipeSchema{
		AtContext:    "https://schema.org",
		AtType:       models.SchemaType{Value: "Recipe"},
		Name:         <-chName,
		Description:  models.Description{Value: <-chDescription},
		PrepTime:     <-chPrep,
		CookTime:     <-chCook,
		Yield:        models.Yield{Value: <-chYield},
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
	}, nil
}
