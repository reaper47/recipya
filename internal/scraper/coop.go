package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCoop(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DateCreated = getNameContent(root, "creation_date")
	rs.DatePublished = rs.DateCreated
	rs.Description.Value = getNameContent(root, "og:description")
	getIngredients(&rs, root.Find(".IngredientList-content"))

	name := getNameContent(root, "og:title")
	before, _, ok := strings.Cut(name, "|")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	nodes := root.Find("ol.List--orderedRecipe")
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	image, _ := root.Find("picture img").First().Attr("src")
	rs.Image.Value = strings.TrimPrefix(image, "//")
	rs.Yield.Value = findYield(root.Find("span:contains('portioner')").Text())

	return rs, nil
}
