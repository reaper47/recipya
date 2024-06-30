package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeTheHappyFoodie(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getPropertyContent(root, "og:title")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Yield.Value = findYield(root.Find(".hf-metadata__portions p").Text())

	rs.PrepTime = root.Find(".hf-metadata__time-prep span").Text()
	if rs.PrepTime != "" {
		parts := strings.Split(rs.PrepTime, " ")
		switch len(parts) {
		case 1:
			minutes := strings.TrimSuffix(parts[0], "min")
			rs.PrepTime = "PT" + minutes + "M"
		case 2:
			hour := strings.TrimSuffix(parts[0], "hr")
			minutes := strings.TrimSuffix(parts[1], "min")
			rs.PrepTime = "PT" + hour + "H" + minutes + "M"
		}
	}

	rs.CookTime = root.Find(".hf-metadata__time-cook span").Text()
	if rs.PrepTime != "" {
		parts := strings.Split(rs.CookTime, " ")
		switch len(parts) {
		case 1:
			minutes := strings.TrimSuffix(parts[0], "min")
			rs.CookTime = "PT" + minutes + "M"
		case 2:
			hour := strings.TrimSuffix(parts[0], "hr")
			minutes := strings.TrimSuffix(parts[1], "min")
			rs.CookTime = "PT" + hour + "H" + minutes + "M"
		}
	}

	nodes := root.Find(".hf-tags__single")
	allKeywords := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		allKeywords[i] = s.Text()
	})
	rs.Keywords.Values = strings.Join(allKeywords, ", ")

	getIngredients(&rs, root.Find(".hf-ingredients__single-group tr"), []models.Replace{{"useFields", ""}}...)
	getInstructions(&rs, root.Find(".hf-method__text p"))

	return rs, nil
}
