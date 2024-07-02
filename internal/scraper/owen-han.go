package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeOwenhan(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.DatePublished = getPropertyContent(root, "datePublished")
	rs.DateModified = getPropertyContent(root, "dateModified")
	rs.Name = getItempropContent(root, "headline")
	rs.Description.Value = getItempropContent(root, "description")
	rs.Category.Value = root.Find(".blog-item-category").First().Text()

	content := root.Find("h4:contains('INGREDIENTS')").Parent()
	getIngredients(&rs, content.Find("ul p"))
	getInstructions(&rs, content.Find("ol p"))

	return rs, nil
}
