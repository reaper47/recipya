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

	rs.Name = root.Find("h1[itemprop='name']").Text()
	rs.PrepTime, _ = root.Find("time[itemprop='prepTime']").Attr("datetime")
	rs.CookTime, _ = root.Find("time[itemprop='cookTime']").Attr("datetime")

	yieldStr, _ := root.Find("input[name='servings']").Attr("value")
	rs.Yield.Value = findYield(yieldStr)

	getIngredients(&rs, root.Find("span[itemprop='ingredient']"), []models.Replace{{"\u00a0", " "}}...)
	getInstructions(&rs, root.Find("span[itemprop='instructions'] p"))

	return rs, nil
}
