package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeSoutherncastiron(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, found := strings.Cut(name, " - ")
	if found {
		name = strings.TrimSpace(before)
	}

	category := strings.TrimSpace(root.Find(".td-crumb-container a").Last().Text())

	prep := root.Find(".recipe-legend").First().Prev().Text()
	split := strings.Split(prep, " ")
	isMin := strings.Contains(prep, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			prep = "PT" + split[i] + "M"
		}
	}

	nodes := root.Find("li[itemprop='ingredients']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := strings.Join(strings.Fields(sel.Text()), " ")
		ingredients = append(ingredients, s)
	})

	nodes = root.Find("li[itemprop='recipeInstructions']")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		s = strings.Join(strings.Fields(s), " ")
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: category},
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		PrepTime:      prep,
		Yield:         models.Yield{Value: findYield(root.Find("div[itemprop='description']").Text())},
	}, nil
}
