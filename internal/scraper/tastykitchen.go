package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeTastyKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Category.Value = root.Find("a[rel='category tag']").First().Text()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = root.Find("h1[itemprop=name]").Text()
	rs.PrepTime = root.Find("time[itemprop=prepTime]").AttrOr("datetime", "")
	rs.CookTime = root.Find("time[itemprop=cookTime]").AttrOr("datetime", "")
	rs.Yield.Value = findYield(root.Find("input[name=servings]").AttrOr("value", ""))
	getIngredients(&rs, root.Find("span[itemprop=ingredient]"), []models.Replace{{"\u00a0", " "}}...)
	getInstructions(&rs, root.Find("span[itemprop=instructions] p"))
	return rs, nil
}
