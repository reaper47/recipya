package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeUsapears(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = getPropertyContent(root, "og:title")

	prep := root.Find(".recipe-legend").First().Prev().Text()
	split := strings.Split(prep, " ")
	isMin := strings.Contains(prep, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			prep = "PT" + split[i] + "M"
		}
	}
	rs.PrepTime = prep

	getIngredients(&rs, root.Find("li[itemprop=ingredients]"), []models.Replace{{"useFields", ""}}...)
	getInstructions(&rs, root.Find("div[itemprop=recipeInstructions] ol li"), []models.Replace{{"useFields", ""}}...)

	return rs, nil
}
