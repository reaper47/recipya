package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeHeatherChristo(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.PrepTime = root.Find("time[itemprop=prepTime]").AttrOr("datetime", "")
	rs.CookTime = root.Find("time[itemprop=cookTime]").AttrOr("datetime", "")
	rs.Name = strings.TrimSpace(root.Find("div[itemprop=name]").First().Text())
	rs.Yield.Value = findYield(root.Find("span[itemprop=recipeYield]").First().Text())
	getIngredients(&rs, root.Find(".ERSIngredients ul").First().Find("li"))
	getInstructions(&rs, root.Find(".ERSInstructions ol").First().Find("li"))
	return rs, nil
}
