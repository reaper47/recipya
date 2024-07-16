package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeStreetKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	article := root.Find("article").First()
	rs.Name = article.Find("h1").First().Text()
	rs.Description.Value = getNameContent(root, "description")
	rs.DatePublished = root.Find("meta[property='article:published_time").AttrOr("content", "")
	rs.DateModified = root.Find("meta[property='article:modified_time").AttrOr("content", "")
	rs.Image.Value = article.Find("img").First().AttrOr("src", "")
	rs.Yield.Value = findYield(article.Find(".c-svgicon--servings").Next().Text())
	getIngredients(&rs, article.Find(".ingredients label"))
	getInstructions(&rs, article.Find(".method-step"))

	return rs, nil
}
