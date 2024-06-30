package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeNinjatestkitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getItempropContent(root, "name")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Image.Value = getItempropContent(root, "image")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.Keywords.Values = getItempropContent(root, "keywords")
	rs.PrepTime = getItempropContent(root, "prepTime")
	rs.Category.Value = getItempropContent(root, "recipeCategory")

	recipeIngredient, _ := root.Find("meta[itemprop='recipeIngredient']").Attr("content")
	ingredients := strings.Split(recipeIngredient, ",")
	rs.Ingredients.Values = make([]string, 0, len(ingredients))
	for _, s := range ingredients {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(s))
	}

	getInstructions(&rs, root.Find(".single-method__method li p"))

	yieldStr, _ := root.Find("meta[itemprop='recipeYield']").Attr("content")
	rs.Yield.Value = findYield(yieldStr)

	return rs, nil
}
