package scraper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeRecipeCommunity(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[itemprop='description']").Attr("content")
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
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	rs.Yield.Value = findYield(root.Find("span[itemprop='recipeYield']").Parent().Text())

	rs.PrepTime, _ = root.Find("#preparation-time-final meta[itemprop='performTime']").Attr("content")
	rs.CookTime, _ = root.Find("#preparation-time-final meta[itemprop='totalTime']").Attr("content")

	nodes := root.Find(".catText")
	allKeywords := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		allKeywords[i] = s.Text()
	})
	rs.Keywords.Values = strings.Join(allKeywords, ", ")

	nodes = root.Find("li[itemprop='recipeIngredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.Join(strings.Fields(s.Text()), " "))
	})

	nodes = root.Find("ol.steps-list li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s.Text()))
	})

	return rs, nil
}
