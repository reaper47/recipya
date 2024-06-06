package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeKuchniadomova(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	name := root.Find("h2[itemprop='name']").Text()
	name = strings.ReplaceAll(name, "\n", "")
	rs.Name = strings.ReplaceAll(name, "\t", "")

	yieldStr := root.Find("p[itemprop='recipeYield']").Text()
	yieldStr = strings.ReplaceAll(yieldStr, "-", " ")
	rs.Yield.Value = findYield(yieldStr)

	rs.Category.Value, _ = root.Find("meta[itemprop='recipeCategory']").Attr("content")
	rs.Keywords.Values, _ = root.Find("meta[name='keywords']").Attr("content")
	rs.Cuisine.Value = root.Find("p[itemprop='recipeCuisine']").Text()

	image, _ := root.Find("#article-img-1").Attr("data-src")
	rs.Image.Value = "https://kuchnia-domowa.pl" + image

	description := root.Find("#recipe-description").Text()
	description = strings.TrimPrefix(description, "\n")
	rs.Description.Value = strings.TrimSuffix(description, "\n")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, s.Text())
	})

	nodes = root.Find("#recipe-instructions li")
	rs.Instructions.Values = make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s.Text()))
	})

	return rs, nil
}
