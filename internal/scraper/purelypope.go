package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapePurelyPope(root *goquery.Document) (models.RecipeSchema, error) {
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")
	image, _ := root.Find("img[itemprop='image']").Attr("src")
	split := strings.Split(image, "?")
	if len(split) > 0 {
		image = split[0]
	}

	name := root.Find("h2[itemprop='name']").Text()
	yield := findYield(root.Find("span[itemprop='recipeYield']").Text())

	prepTime, _ := root.Find("time[itemprop='prepTime']").Attr("datetime")
	prepTime = strings.ReplaceAll(prepTime, " ", "")

	cookTime, _ := root.Find("time[itemprop='cookTime']").Attr("datetime")
	cookTime = strings.ReplaceAll(cookTime, " ", "")

	nodes := root.Find("span[itemprop='recipeIngredient']").FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.TrimSpace(s.Text()) != ""
	})
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = strings.TrimSpace(s.Text())
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] li")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = s.Text()
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		DateModified:  dateModified,
		DatePublished: datePublished,
		Image:         models.Image{Value: image},
		Name:          name,
		PrepTime:      prepTime,
		CookTime:      cookTime,
		Yield:         models.Yield{Value: yield},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
	}, nil
}
