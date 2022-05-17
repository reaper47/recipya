package scraper

import (
	"strings"
	"time"

	"github.com/reaper47/recipya/internal/constants"
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeSaveur(root *html.Node) (models.RecipeSchema, error) {
	chDescription := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDescription <- s
		}()

		s = <-getElementData(root, "class", "Article-excerpt")
		s = strings.ReplaceAll(s, "\n", "")
		s = strings.TrimSpace(s)
	}()

	chDatePublished := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDatePublished <- s
		}()

		s = <-getElementData(root, "class", "Article-dateTime")
		s = strings.TrimSpace(s)
		s = strings.TrimPrefix(s, "Published ")

		xs := strings.Split(s, " ")
		if len(xs[1]) == 2 {
			xs[1] = "0" + xs[1]
		}
		s = strings.Join(xs[:3], " ")

		t, err := time.Parse("Jan 02, 2006", s)
		if err == nil {
			s = t.Format(constants.BasicTimeLayout)
		}
	}()

	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		div := getElement(root, "class", "SingleImage-wrapper")
		s = getAttr(div.FirstChild, "src")
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		yieldStr := <-getElementData(root, "property", "recipeYield")
		yield = findYield(strings.Split(yieldStr, " "))
	}()

	chCookTime := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chCookTime <- s
		}()

		s = getAttr(getElement(root, "property", "cookTime"), "content")
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		xn := traverseAll(root, func(node *html.Node) bool {
			return getAttr(node, "property") == "recipeIngredient"
		})
		for _, n := range xn {
			v := strings.ReplaceAll(n.FirstChild.Data, "\n", "")
			v = strings.TrimSpace(v)
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
			return getAttr(node, "property") == "recipeInstructions"
		})
		for _, n := range xn {
			v := strings.ReplaceAll(n.FirstChild.Data, "\n", "")
			v = strings.TrimSpace(v)
			vals = append(vals, v)
		}
	}()

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		Name:          <-getElementData(root, "class", "u-entryTitle"),
		Description:   models.Description{Value: <-chDescription},
		DatePublished: <-chDatePublished,
		Image:         models.Image{Value: <-chImage},
		Yield:         models.Yield{Value: <-chYield},
		CookTime:      <-chCookTime,
		Ingredients:   models.Ingredients{Values: <-chIngredients},
		Instructions:  models.Instructions{Values: <-chInstructions},
	}, nil
}
