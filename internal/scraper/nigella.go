package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeNigella(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	name := getPropertyContent(root, "og:title")
	before, _, ok := strings.Cut(name, " — ")
	if ok {
		name = strings.TrimSpace(before)
	}
	if name == "" {
		name = root.Find("h1[itemprop=name]").Text()
	}
	rs.Name = name

	rs.Yield.Value = findYield(root.Find("p[class=serves]").Text())

	untrimmedDescription := getPropertyContent(root, "og:description")
	description, _ := strings.CutSuffix(untrimmedDescription, "\n\nFor US cup measures, use the toggle at the top of the ingredients list.")
	if description == "" {
		description = getNameContent(root, "description")
	}
	rs.Description.Value = description

	rs.Image.Value = getPropertyContent(root, "og:image")
	if rs.Image.Value == "" {
		img := root.Find("img[itemprop=image]").AttrOr("src", "")
		if strings.HasPrefix(img, "/assets") {
			img = "https://www.nigella.com" + img
		}
		rs.Image.Value = img
	}

	rs.Keywords.Values = root.Find("meta[itemprop=keywords]").AttrOr("content", "")

	getIngredients(&rs, root.Find("li[itemprop=recipeIngredient]"))
	getInstructions(&rs, root.Find("div[itemprop=recipeInstructions]").First().Find("ol li"))

	return rs, nil
}
