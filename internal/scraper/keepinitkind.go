package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeKeepinItKind(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Name = strings.TrimSpace(root.Find("div[itemprop=name]").Text())
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Yield.Value = findYield(root.Find("*[itemprop=recipeYield]").Text())
	getIngredients(&rs, root.Find("li[itemprop=ingredients]"))
	getInstructions(&rs, root.Find("li[itemprop=recipeInstructions]"))
	return rs, nil
}
