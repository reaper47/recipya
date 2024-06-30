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

	getIngredients(&rs, content.Find(".recipe-ingredients li"), []models.Replace{
		{"  ", " "},
	}...)

	getInstructions(&rs, content.Find(".recipe-directions li"), []models.Replace{
		{"\u00a0", " "},
	}...)

	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.DatePublished = getPropertyContent(root, "article:published_time")

	yieldStr := content.Find(".recipe-card-servings .recipe-card-amount").Text()
	yield, _ := strconv.ParseInt(yieldStr, 10, 16)
	rs.Yield.Value = int16(yield)

	name := content.Find(".recipe-card-title").Text()
	name = strings.TrimLeft(name, "\n")
	rs.Name = strings.TrimSpace(name)

	return rs, nil
}
