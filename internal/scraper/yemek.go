package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeYemek(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Name = getPropertyContent(root, "og:title")
	rs.Image.Value = getPropertyContent(root, "og:image")
	getIngredients(&rs, root.Find("li[itemprop='recipeIngredient']"))
	getInstructions(&rs, root.Find("p[itemprop='recipeInstructions']"))
	return rs, nil
}
