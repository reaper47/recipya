package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapePotatoRolls(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = root.Find("[itemprop=name]").Text()
	rs.Image.Value = root.Find("[itemprop=image]").First().AttrOr("src", "")
	getIngredients(&rs, root.Find(".ingredient"))
	getInstructions(&rs, root.Find(".direction p"))
	getTime(&rs, root.Find(".icon-clock").First().Parent(), true)
	return rs, nil
}
