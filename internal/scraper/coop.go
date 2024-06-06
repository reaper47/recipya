package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCoop(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DateCreated, _ = root.Find("meta[name='creation_date']").Attr("content")

	rs.Name, _ = root.Find("meta[name='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, "|")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	rs.Description.Value, _ = root.Find("meta[name='og:description']").Attr("content")

	nodes := root.Find(".IngredientList-content")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find("ol.List--orderedRecipe")
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	rs.Image.Value, _ = root.Find("picture img").First().Attr("src")
	rs.Image.Value = strings.TrimPrefix(image, "//")
	rs.Yield.Value = findYield(root.Find("span:contains('portioner')").Text())

	return rs, nil
}
