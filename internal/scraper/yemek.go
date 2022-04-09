package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeYemek(root *html.Node) (rs models.RecipeSchema, err error) {
	root = getElement(root, "itemtype", "http://schema.org/Recipe")

	chName := make(chan string)
	go func() {
		v := "<Recipe Name>"
		defer func() {
			_ = recover()
			chName <- v
		}()

		node := getElement(root, "id", "malzemeler")
		if node.FirstChild.Data == "h2" {
			v = node.FirstChild.FirstChild.NextSibling.FirstChild.Data
		}
	}()

	chCategory := make(chan string)
	go func() {
		var v string
		defer func() {
			_ = recover()
			chCategory <- v
		}()

		node := getElement(root, "class", "main-category-image")
		v = getAttr(node, "title")
	}()

	chImage := make(chan string)
	go func() {
		node := getElement(root, "itemprop", "image")
		chImage <- getAttr(node, "src")
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		go func() {
			_ = recover()
			chIngredients <- vals
		}()

		node := getElement(root, "itemprop", "recipeIngredient")
		for c := node; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				ing := strings.TrimSpace(c.FirstChild.Data)
				if ing == "" {
					continue
				}
				ing = strings.ReplaceAll(ing, "  ", " ")
				vals = append(vals, ing)
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

		node := getElement(root, "id", "hazir1anis")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "p" {
				ins := strings.TrimSpace(c.LastChild.Data)
				if ins == "" || ins == "fb:like" {
					continue
				}
				ins = strings.ReplaceAll(ins, "  ", " ")

				vals = append(vals, ins)
			}
		}
	}()

	return models.RecipeSchema{
		AtContext:    "https://schema.org",
		AtType:       models.SchemaType{Value: "Recipe"},
		Image:        models.Image{Value: <-chImage},
		Ingredients:  models.Ingredients{Values: <-chIngredients},
		Instructions: models.Instructions{Values: <-chInstructions},
		Name:         <-chName,
		Category:     models.Category{Value: <-chCategory},
	}, nil
}
