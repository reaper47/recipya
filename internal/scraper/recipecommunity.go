package scraper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeRecipeCommunity(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getPropertyContent(root, "og:title")
	rs.Description.Value = getItempropContent(root, "description")
	rs.DatePublished = root.Find(".recipe-summary .creation-date").Text()
	if strings.Contains(rs.DatePublished, ":") {
		rs.DatePublished = strings.Trim(strings.Split(rs.DatePublished, ":")[1], " ")
	}
	if strings.Contains(rs.DatePublished, "/") {
		rs.DatePublished = strings.ReplaceAll(rs.DatePublished, "/", "-")
	}
	rs.DateModified = root.Find(".recipe-summary .changed-date").Text()
	if strings.Contains(rs.DateModified, ":") {
		rs.DateModified = strings.Trim(strings.Split(rs.DateModified, ":")[1], " ")
	}
	if strings.Contains(rs.DateModified, "/") {
		rs.DateModified = strings.ReplaceAll(rs.DateModified, "/", "-")
	}
	rs.Image.Value = getPropertyContent(root, "og:image")

	rs.Yield.Value = findYield(root.Find("span[itemprop=recipeYield]").Parent().Text())

	rs.PrepTime = strings.Replace(root.Find("meta[itemprop=performTime]").AttrOr("content", ""), "min", "M", 1)
	rs.CookTime = strings.Replace(root.Find("meta[itemprop=totalTime]").AttrOr("content", ""), "min", "M", 1)

	nodes := root.Find(".catText")
	allKeywords := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		allKeywords[i] = s.Text()
	})
	rs.Keywords.Values = strings.Join(allKeywords, ", ")

	getIngredients(&rs, root.Find("li[itemprop=recipeIngredient]"), []models.Replace{{"useFields", ""}}...)
	getInstructions(&rs, root.Find("ol.steps-list li"))

	return rs, nil
}
