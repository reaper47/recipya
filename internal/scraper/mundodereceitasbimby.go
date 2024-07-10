package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeMundodereceitasbimby(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getNameContent(root, "og:description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Category.Value = root.Find("span[itemprop=recipeCategory]").First().Text()
	rs.Name = getItempropContent(root, "name")
	rs.Yield.Value = findYield(root.Find("span[itemprop=recipeYield]").Text())
	rs.DateModified = getItempropContent(root, "dateModified")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.DateCreated = rs.DatePublished
	rs.CookTime = strings.ReplaceAll(getItempropContent(root, "cookTime"), "min", "M")
	getIngredients(&rs, root.Find("li[itemprop=recipeIngredient]"))
	getInstructions(&rs, root.Find("ol[itemprop=recipeInstructions] div[itemprop=itemListElement]"))

	nodes := root.Find("meta[itemprop=tool]")
	rs.Tools.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(sel.AttrOr("content", "")))
	})

	return rs, nil
}
