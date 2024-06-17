package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeNigella(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, " — ")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	rs.Yield.Value = findYield(root.Find("p[class='serves']").Text())

	untrimmedDescription, _ := root.Find("meta[property='og:description']").Attr("content")
	description, _ := strings.CutSuffix(untrimmedDescription, "\n\nFor US cup measures, use the toggle at the top of the ingredients list.")
	rs.Description.Value = description
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	keywords, _ := root.Find("meta[itemprop='keywords]").Attr("content")
	rs.Keywords.Values = keywords

	nodes := root.Find("li[itemprop=recipeIngredient]")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("div[itemprop=recipeInstructions]").First().Find("ol li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(sel.Text())))
	})

	return rs, nil
}
