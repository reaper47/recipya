package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeTheHappyFoodie(root *html.Node) (models.RecipeSchema, error) {
	chName := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chName <- s
		}()

		node := getElement(root, "class", "hf-title__inner")
		s = node.FirstChild.NextSibling.FirstChild.Data
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, "\n", "")
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(root, "class", "hf-metadata__portions d-flex align-items-center")
		yieldStr := node.FirstChild.NextSibling.FirstChild.Data
		var parts []string
		if strings.Contains(yieldStr, "-") {
			parts = strings.Split(yieldStr, "-")
		} else {
			parts = strings.Split(yieldStr, " ")
		}
		yield = findYield(parts)
	}()

	chDescription := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDescription <- s
		}()

		node := getElement(root, "class", "hf-summary__inner")
		s = node.FirstChild.NextSibling.FirstChild.Data
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, "\n", "")
	}()

	chPrepTime := make(chan string)
	go func() {
		var prep string
		defer func() {
			_ = recover()
			chPrepTime <- prep
		}()

		node := getElement(root, "class", "hf-metadata__time-prep")
		s := node.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.Data
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, "\n", "")
		if strings.HasSuffix(s, "min") {
			s = strings.TrimSuffix(s, "min")
			prep = "PT" + s + "M"
		}
	}()

	chCookTime := make(chan string)
	go func() {
		var cook string
		defer func() {
			_ = recover()
			chCookTime <- cook
		}()

		node := getElement(root, "class", "hf-metadata__time-cook")
		s := node.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.Data
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, "\n", "")
		if strings.HasSuffix(s, "min") {
			s = strings.TrimSuffix(s, "min")
			cook = "PT" + s + "M"
		}
	}()

	chKeywords := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chKeywords <- s
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "class") == "hf-tags__single"
		})
		var vals []string
		for _, n := range xn {
			vals = append(vals, n.FirstChild.Data)
		}
		s = strings.Join(vals, ", ")
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		node := getElement(root, "class", "hf-ingredients__single-group")
		tbody := traverseAll(node, func(node *html.Node) bool {
			return node.Data == "tbody"
		})[0]
		for tr := tbody.FirstChild; tr != nil; tr = tr.NextSibling {
			if tr.Type != html.ElementNode {
				continue
			}

			amount := tr.FirstChild.NextSibling.FirstChild.Data
			amount = strings.ReplaceAll(amount, "\n", "")
			amount = strings.ReplaceAll(amount, "\t", "")

			name := tr.LastChild.PrevSibling.FirstChild.Data
			name = strings.ReplaceAll(name, "\n", "")
			name = strings.ReplaceAll(name, "\t", "")

			v := strings.TrimSpace(amount + " " + name)
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

		node := getElement(root, "class", "hf-method__text")
		for p := node.FirstChild; p != nil; p = p.NextSibling {
			if p.Type != html.ElementNode {
				continue
			}
			vals = append(vals, p.FirstChild.Data)
		}
	}()

	return models.RecipeSchema{
		AtContext:    "https://schema.org",
		AtType:       models.SchemaType{Value: "Recipe"},
		Name:         <-chName,
		Description:  models.Description{Value: <-chDescription},
		Yield:        models.Yield{Value: <-chYield},
		PrepTime:     <-chPrepTime,
		CookTime:     <-chCookTime,
		Keywords:     models.Keywords{Values: <-chKeywords},
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
	}, nil
}
