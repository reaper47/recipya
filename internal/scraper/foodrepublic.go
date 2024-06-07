package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeFoodRepublic(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	content := root.Find(".recipe-card")

	node := content.Find(".recipe-card-prep-time")
	prepTime := node.Find(".recipe-card-amount").Text()
	if node.Find(".recipe-card-unit").Text() == "minutes" {
		prepTime = "PT" + prepTime + "M"
	}
	rs.PrepTime = prepTime

	node = content.Find(".recipe-card-cook-time")
	cookTime := node.Find(".recipe-card-amount").Text()
	if node.Find(".recipe-card-unit").Text() == "minutes" {
		cookTime = "PT" + cookTime + "M"
	}
	rs.CookTime = cookTime

	rs.Image.Value, _ = content.Find(".recipe-card-image img").Attr("data-lazy-src")
	rs.Description.Value = content.Find(".recipe-card-description").Text()

	nodes := content.Find(".recipe-ingredients li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.ReplaceAll(s.Text(), "  ", " "))
	})

	nodes = content.Find(".recipe-directions li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\u00a0", " ")
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(v))
	})

	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")

	yieldStr := content.Find(".recipe-card-servings .recipe-card-amount").Text()
	yield, _ := strconv.ParseInt(yieldStr, 10, 16)
	rs.Yield.Value = int16(yield)

	name := content.Find(".recipe-card-title").Text()
	name = strings.TrimLeft(name, "\n")
	rs.Name = strings.TrimSpace(name)

	return rs, nil
}
