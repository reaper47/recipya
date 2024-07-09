package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeOkokorecepten(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Keywords.Values = getNameContent(root, "keywords")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = getPropertyContent(root, "og:title")
	rs.Yield.Value = findYield(root.Find("li:contains('personen')").Text())
	rs.Category.Value, _, _ = strings.Cut(strings.TrimPrefix(getPropertyContent(root, "og:url"), "https://www.okokorecepten.nl/recept/"), "/")

	prep := strings.TrimSpace(root.Find("li:contains('min')").Text())
	if prep != "" {
		rs.PrepTime = "PT" + regex.Digit.FindString(prep) + "M"
	}

	before, _, ok := strings.Cut(rs.Name, "- recept - okoko recepten")
	if ok {
		rs.Name = strings.TrimSpace(before)
	}

	getIngredients(&rs, root.Find("ul.list-ingredienten li"))
	getInstructions(&rs, root.Find("h2:contains('Bereiden')").NextAll())

	return rs, nil
}
