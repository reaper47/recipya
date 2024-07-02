package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeRezeptwelt(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getItempropContent(root, "name")
	rs.Category.Value = getItempropContent(root, "recipeCategory")
	rs.Description.Value = getItempropContent(root, "description")
	rs.Image.Value, _ = root.Find("img[itemprop=image]").Attr("src")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.DateModified = getItempropContent(root, "dateModified")
	getIngredients(&rs, root.Find("li[itemprop=recipeIngredient]"))
	rs.PrepTime = getItempropContent(root, "performTime")
	rs.Cuisine.Value = getItempropContent(root, "recipeCuisine")
	rs.Keywords.Values = getItempropContent(root, "keywords")
	getInstructions(&rs, root.Find("ol[itemprop=recipeInstructions] li"))
	rs.Yield.Value = findYield(root.Find("span[itemprop=recipeYield]").Text())

	nodes := root.Find("meta[itemprop=tool]")
	rs.Tools.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		t, _ := s.Attr("content")
		rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(t))
	})

	return rs, nil
}
