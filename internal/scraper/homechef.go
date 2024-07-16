package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeHomechef(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getNameContent(root, "description")
	rs.Yield.Value = findYield(root.Find("meta[itemprop=recipeYield]").AttrOr("content", ""))
	rs.Image.Value = root.Find("div img").First().AttrOr("data-srcset", "")
	rs.Name = root.Find("h1").First().Text()

	getIngredients(&rs, root.Find("li[itemprop=recipeIngredient]"), []models.Replace{
		{"\n", " "},
		{"Info", ""},
		{"useFields", ""},
	}...)

	getInstructions(&rs, root.Find("li[itemprop=itemListElement]"), []models.Replace{
		{"\n", " "},
		{"useFields", ""},
	}...)

	return rs, nil
}
