package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeNinjatestkitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[itemprop='name']").Attr("content")
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	rs.Description.Value = strings.TrimSpace(description)
	rs.Image.Value, _ = root.Find("meta[itemprop='image']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[itemprop='datePublished']").Attr("content")
	rs.Keywords.Values, _ = root.Find("meta[itemprop='keywords']").Attr("content")
	rs.PrepTime, _ = root.Find("meta[itemprop='prepTime']").Attr("content")
	rs.Category.Value, _ = root.Find("meta[itemprop='recipeCategory']").Attr("content")

	recipeIngredient, _ := root.Find("meta[itemprop='recipeIngredient']").Attr("content")
	ingredients := strings.Split(recipeIngredient, ",")
	rs.Ingredients.Values = make([]string, 0, len(ingredients))
	for _, s := range ingredients {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(s))
	}

	nodes := root.Find(".single-method__method li p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	yieldStr, _ := root.Find("meta[itemprop='recipeYield']").Attr("content")
	rs.Yield.Value = findYield(yieldStr)

	return rs, nil
}
