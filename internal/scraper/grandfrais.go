package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeGrandfrais(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[property='og:description']").Attr("content")

	prep := root.Find(".pre-requie-item p").First().Text()
	split := strings.Split(prep, " ")
	isMin := strings.Contains(prep, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			prep = "PT" + split[i] + "M"
		}
	}

	cook := root.Find(".pre-requie-item p").Last().Text()
	split = strings.Split(cook, " ")
	isMin = strings.Contains(cook, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			cook = "PT" + split[i] + "M"
		}
	}

	nodes := root.Find("div[itemprop='ingredients']")
	ingredients := strings.Split(nodes.Text(), "\n")
	for i, ingredient := range ingredients {
		_, after, ok := strings.Cut(ingredient, "- ")
		if ok {
			ingredients[i] = after
		}
	}

	nodes = root.Find("div[itemprop='recipeInstructions'] li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, s)
	})

	image, _ := root.Find("img[itemprop='image']").Attr("src")
	if !strings.HasPrefix(image, "https://") {
		image = "https://www.grandfrais.com" + image
	}

	return models.RecipeSchema{
		AtContext:    atContext,
		AtType:       models.SchemaType{Value: "Recipe"},
		CookTime:     cook,
		Description:  models.Description{Value: description},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         root.Find("h1[itemprop='name']").Text(),
		PrepTime:     prep,
		Yield:        models.Yield{Value: findYield(root.Find("p[itemprop='recipeYield']").Text())},
	}, nil
}
