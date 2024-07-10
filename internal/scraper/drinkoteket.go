package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeDrinkoteket(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = root.Find("meta[property='og:title']").Last().AttrOr("content", "")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Category.Value = getPropertyContent(root, "article:section")
	rs.Image.Value = getPropertyContent(root, "og:image")
	getIngredients(&rs, root.Find("ul.ingredients li"))
	getInstructions(&rs, root.Find("div[itemprop=recipeInstructions] li"))

	nodes := root.Find("#recipe-utrustning .rbs-img-content")
	rs.Tools.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(strings.TrimSpace(sel.Text())))
	})

	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.PrepTime = getItempropContent(root, "prepTime")
	rs.CookTime = getItempropContent(root, "cookTime")

	return rs, nil
}
