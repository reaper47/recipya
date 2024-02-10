package scraper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeRecipeCommunity(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	description, _ := root.Find("meta[itemprop='description']").Attr("content")
	datePublished := root.Find(".recipe-summary .creation-date").Text()
	if strings.Contains(datePublished, ":") {
		datePublished = strings.Trim(strings.Split(datePublished, ":")[1], " ")
	}
	dateModified := root.Find(".recipe-summary .changed-date").Text()
	if strings.Contains(dateModified, ":") {
		dateModified = strings.Trim(strings.Split(dateModified, ":")[1], " ")
	}
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	yield := findYield(root.Find("span[itemprop='recipeYield']").Parent().Text())

	prepTime, _ := root.Find("#preparation-time-final meta[itemprop='performTime']").Attr("content")
	cookTime, _ := root.Find("#preparation-time-final meta[itemprop='totalTime']").Attr("content")

	nodes := root.Find(".catText")
	allKeywords := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		allKeywords[i] = s.Text()
	})
	keywords := strings.Join(allKeywords, ", ")

	nodes = root.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = strings.Join(strings.Fields(s.Text()), " ")
	})

	nodes = root.Find("ol.steps-list li")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = s.Text()
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		Name:          name,
		DatePublished: datePublished,
		DateModified:  dateModified,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Yield:         models.Yield{Value: yield},
		PrepTime:      prepTime,
		CookTime:      cookTime,
		Keywords:      models.Keywords{Values: keywords},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
	}, nil
}
