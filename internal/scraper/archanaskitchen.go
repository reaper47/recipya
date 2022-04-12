package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeArchanasKitchen(root *html.Node) (rs models.RecipeSchema, err error) {
	chCategory := make(chan string)
	go func() {
		var v string
		defer func() {
			_ = recover()
			chCategory <- v
		}()

		node := getElement(root, "class", "col-12 recipeCategory")
		v = node.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.FirstChild.Data
	}()

	chCuisine := make(chan string)
	go func() {
		var v string
		defer func() {
			_ = recover()
			chCuisine <- v
		}()

		node := getElement(root, "itemprop", "recipeCuisine")
		v = node.FirstChild.FirstChild.Data
	}()

	chDescription := make(chan string)
	go func() {
		var v string
		defer func() {
			_ = recover()
			chDescription <- v
		}()

		node := getElement(root, "property", "og:description")
		v = getAttr(node, "content")
	}()

	chImage := make(chan string)
	go func() {
		var v string
		defer func() {
			_ = recover()
			chImage <- v
		}()

		node := getElement(root, "itemprop", "image")
		v = "https://www.archanaskitchen.com" + getAttr(node, "src")
	}()

	chKeywords := make(chan string)
	go func() {
		var v string
		defer func() {
			_ = recover()
			chKeywords <- strings.TrimSuffix(v, ",")
		}()

		node := getElement(root, "class", "itemTags")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}
			v += strings.TrimSpace(c.FirstChild.NextSibling.FirstChild.Data) + ","
		}
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		node := getElement(root, "class", "ingredientstitle")
		list := node.NextSibling.NextSibling
		for c := list.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && getAttr(c, "itemprop") == "ingredients" {
				var v string
				for c := c.FirstChild; c != nil; c = c.NextSibling {
					var d string
					if c.Type == html.ElementNode {
						d = c.FirstChild.Data
					} else {
						d = c.Data
					}
					v += " " + strings.TrimSpace(d)
				}
				v = strings.TrimSpace(strings.ReplaceAll(v, " , ", ", "))
				vals = append(vals, v)
			}
		}
	}()

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- vals
		}()

		node := getElement(root, "class", "recipeinstructionstitle")
		list := node.NextSibling.NextSibling
		for c := list.FirstChild; c != nil; c = c.NextSibling {
			var v string
			for c := c.FirstChild.FirstChild; c != nil; c = c.NextSibling {
				var d string
				if c.Type == html.ElementNode {
					d = c.FirstChild.Data
				} else {
					d = c.Data
				}
				v += " " + strings.TrimSpace(d)
			}

			v = strings.ReplaceAll(v, "\u00a0", " ")
			v = strings.TrimSpace(strings.ReplaceAll(v, " .", "."))
			vals = append(vals, v)
		}
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(root, "itemprop", "recipeYield")
		servings := node.FirstChild.NextSibling.FirstChild.Data
		yield = findYield(strings.Split(servings, " "))
	}()

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		Name:          <-getElementData(root, "class", "recipe-title"),
		Description:   models.Description{Value: <-chDescription},
		Image:         models.Image{Value: <-chImage},
		Category:      models.Category{Value: <-chCategory},
		Cuisine:       models.Cuisine{Value: <-chCuisine},
		PrepTime:      <-getItemPropAttr(root, "prepTime", "content"),
		CookTime:      <-getItemPropAttr(root, "cookTime", "content"),
		DatePublished: <-getItemPropAttr(root, "datePublished", "content"),
		DateModified:  <-getItemPropAttr(root, "dateModified", "content"),
		Keywords:      models.Keywords{Values: <-chKeywords},
		Ingredients:   models.Ingredients{Values: <-chIngredients},
		Instructions:  models.Instructions{Values: <-chInstructions},
		Yield:         models.Yield{Value: <-chYield},
	}, nil
}
