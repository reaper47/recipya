package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapePlentyVegan(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Name = root.Find("[itemprop=name].ERSName").Text()
	node := root.Find("[class*=wp-image]").First()
	rs.Image.Value, _ = node.Attr("src")
	rs.Description.Value = strings.TrimSpace(node.Parent().Text())
	rs.PrepTime, _ = root.Find("time[itemprop=prepTime]").Attr("datetime")
	rs.CookTime, _ = root.Find("time[itemprop=cookTime]").Attr("datetime")
	rs.Category.Value = strings.ToLower(root.Find("[itemprop=recipeCuisine]").Text())
	rs.Yield.Value = findYield(root.Find("[itemprop=recipeYield]").Text())
	getIngredients(&rs, root.Find("li[itemprop=ingredients]"))
	getInstructions(&rs, root.Find("li[itemprop=recipeInstructions]"))
	return rs, nil
}
