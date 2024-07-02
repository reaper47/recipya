package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeKwestiasmaku(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = getPropertyContent(root, "og:title")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Description.Value = root.Find("span[itemprop=description]").Text()
	rs.Yield.Value = findYield(root.Find(".field-name-field-ilosc-porcji").Text())

	getIngredients(&rs, root.Find(".field-name-field-skladniki li"), []models.Replace{
		{"\n", ""},
		{"\t", ""},
	}...)

	getInstructions(&rs, root.Find(".field-name-field-przygotowanie li"), []models.Replace{
		{"\n", ""},
		{"\t", ""},
	}...)

	return rs, nil
}
