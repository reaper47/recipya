package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeDk(root *goquery.Document) (models.RecipeSchema, error) {
	content := root.Find("section[itemtype='http://schema.org/Recipe']")

	yieldStr, _ := content.Find("section[itemprop='recipeYield']").Attr("content")
	yield, _ := strconv.ParseInt(yieldStr, 10, 16)

	nodes := content.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.TrimSpace(s.Text())
		v = strings.ReplaceAll(v, "\n", "")
		ingredients[i] = strings.Join(strings.Fields(v), " ")
	})

	instructions := make([]string, 0)
	content.Find("div[itemprop='recipeInstructions'] h3,div[itemprop='recipeInstructions'] li").Each(func(i int, s *goquery.Selection) {
		if i > 0 && s.Nodes[0].Data == "h3" {
			instructions = append(instructions, "\n")
		}

		v := strings.ReplaceAll(s.Text(), "\n", "")
		v = strings.ReplaceAll(v, "\u00a0", "")
		instructions = append(instructions, strings.TrimSpace(v))
	})

	/*chInstructions := make(chan models.Instructions)
	go func() {
		var v models.Instructions
		defer func() {
			_ = recover()
			chInstructions <- v
		}()

		node := getElement(content, "itemprop", "recipeInstructions")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			switch c.Data {
			case "h3":
				s := c.FirstChild.Data
				s = strings.ReplaceAll(s, "\n", "")
				v.Values = append(v.Values, "\n", strings.TrimSpace(s))
			case "ol":
				for l := c.FirstChild; l != nil; l = l.NextSibling {
					s := l.FirstChild.Data
					s = strings.ReplaceAll(s, "\n", "")
					v.Values = append(v.Values, strings.TrimSpace(s))
				}
			}
		}
		v.Values = v.Values[1:]
	}()*/

	description := content.Find("p[itemprop='description']").Text()
	description = strings.ReplaceAll(description, "\n", "")

	image, _ := content.Find("meta[itemprop='url']").Attr("content")
	datePublished, _ := content.Find("meta[itemprop='datePublished']").Attr("content")

	return models.RecipeSchema{
		Name:          content.Find("h1[itemprop='name']").Text(),
		Image:         models.Image{Value: image},
		DatePublished: datePublished,
		Description:   models.Description{Value: strings.TrimSpace(description)},
		Yield:         models.Yield{Value: int16(yield)},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
	}, nil
}
