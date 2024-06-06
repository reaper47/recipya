package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapePurelyPope(root *goquery.Document) (models.RecipeSchema, error) {
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.Image.Value, _ = root.Find("img[itemprop='image']").Attr("src")
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

	nodes := root.Find("span[itemprop='recipeIngredient']").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.TrimSpace(s.Text()) != ""
	})
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(s.Text()))
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] li")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		instructions = append(instructions, models.NewHowToStep(s.Text()))
	})

	return models.RecipeSchema{
		DateModified:  dateModified,
		DatePublished: datePublished,
		Image:         &models.Image{Value: image},
		Name:          name,
		PrepTime:      prepTime,
		CookTime:      cookTime,
		Yield:         &models.Yield{Value: yield},
		Ingredients:   &models.Ingredients{Values: ingredients},
		Instructions:  &models.Instructions{Values: instructions},
	}, nil
}
