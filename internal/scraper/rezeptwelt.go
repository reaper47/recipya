package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeRezeptwelt(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getItempropContent(root, "name")
	rs.Category.Value = getItempropContent(root, "recipeCategory")
	rs.Description.Value = getItempropContent(root, "description")
	rs.Image.Value = root.Find("img[itemprop=image]").AttrOr("src", "")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.DateModified = getItempropContent(root, "dateModified")
	getIngredients(&rs, root.Find("li[itemprop=recipeIngredient]"))
	rs.PrepTime = strings.Replace(getItempropContent(root, "performTime"), "min", "M", 1)
	rs.Cuisine.Value = getItempropContent(root, "recipeCuisine")
	rs.Keywords.Values = getItempropContent(root, "keywords")
	getInstructions(&rs, root.Find("ol.steps-list li"))
	rs.Yield.Value = findYield(root.Find("span[itemprop=recipeYield]").Text())

	nodes := root.Find("meta[itemprop=tool]")
	rs.Tools.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		t := s.AttrOr("content", "")
		if t != "" {
			rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(t))
		}
	})

	return rs, nil
}
