package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapePopsugar(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Keywords.Values = getNameContent(root, "keywords")
	rs.Description.Value = getNameContent(root, "description")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = strings.TrimSpace(root.Find("[class^=RecipeDetailsstyles__RecipeTitle]").Text())
	rs.Yield.Value = findYield(root.Find("span:contains('Yield')").Next().Text())
	getTime(&rs, root.Find("span:contains('Cook Time')").Next(), false)
	rs.NutritionSchema.Calories = strings.TrimSpace(root.Find("span:contains('Calories')").Next().Text())
	getIngredients(&rs, root.Find("[class^=RecipeDetailsstyles__RecipeIngredientText] li"))
	getInstructions(&rs, root.Find("[class^=RecipeDetailsstyles__RecipeInstructions] li"))
	return rs, nil
}
