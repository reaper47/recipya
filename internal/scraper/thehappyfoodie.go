package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeTheHappyFoodie(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
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

	nodes = root.Find(".hf-ingredients__single-group tr")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.Join(strings.Fields(s.Text()), " "))
	})

	nodes = root.Find(".hf-method__text p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s.Text()))
	})

	return rs, nil
}
