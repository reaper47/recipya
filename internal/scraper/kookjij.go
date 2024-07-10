package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeKookjij(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Name = strings.TrimSpace(root.Find("h1.recipe-name").Text())
	rs.Image.Value = root.Find("a.photo").First().Find("img").AttrOr("src", "")
	rs.Cuisine.Value = root.Find("div.details strong:contains('Herkomst')").Next().Text()
	rs.DatePublished = root.Find("time").First().AttrOr("datetime", "")

	str := strings.TrimSpace(root.Find("div.details strong:contains('Personen')").Parent().Text())
	rs.Yield.Value = findYield(strings.ReplaceAll(str, "\n", " "))

	cat := root.Find("li[itemprop=recipeCategory]").First().Text()
	rs.Category.Value = strings.TrimSpace(strings.TrimPrefix(cat, ">"))

	getIngredients(&rs, root.Find("[itemprop=ingredients]"))
	getInstructions(&rs, root.Find("[itemprop=recipeInstructions] li"))
	getTime(&rs, root.Find("div.details strong:contains('Bereidingstijd')").Parent(), true)

	nodes := root.Find("div.recipe-accessoires label")
	rs.Tools.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "" {
			rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(s))
		}
	})

	return rs, nil
}
