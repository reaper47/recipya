package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeTheHappyFoodie(root *goquery.Document) (models.RecipeSchema, error) {
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	yield := findYield(root.Find(".hf-metadata__portions p").Text())

	prepTimeStr := root.Find(".hf-metadata__time-prep span").Text()
	var prepTime string
	if prepTimeStr != "" {
		parts := strings.Split(prepTimeStr, " ")
		switch len(parts) {
		case 1:
			minutes := strings.TrimSuffix(parts[0], "min")
			prepTime = "PT" + minutes + "M"
		case 2:
			hour := strings.TrimSuffix(parts[0], "hr")
			minutes := strings.TrimSuffix(parts[1], "min")
			prepTime = "PT" + hour + "H" + minutes + "M"
		}
	}

	cookTimeStr := root.Find(".hf-metadata__time-cook span").Text()
	var cookTime string
	if prepTimeStr != "" {
		parts := strings.Split(cookTimeStr, " ")
		switch len(parts) {
		case 1:
			minutes := strings.TrimSuffix(parts[0], "min")
			cookTime = "PT" + minutes + "M"
		case 2:
			hour := strings.TrimSuffix(parts[0], "hr")
			minutes := strings.TrimSuffix(parts[1], "min")
			cookTime = "PT" + hour + "H" + minutes + "M"
		}
	}

	nodes := root.Find(".hf-tags__single")
	allKeywords := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		allKeywords[i] = s.Text()
	})
	keywords := strings.Join(allKeywords, ", ")

	nodes = root.Find(".hf-ingredients__single-group tr")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = strings.Join(strings.Fields(s.Text()), " ")
	})

	nodes = root.Find(".hf-method__text p")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		instructions = append(instructions, models.NewHowToStep(s.Text()))
	})

	return models.RecipeSchema{
		Name:          name,
		DatePublished: datePublished,
		DateModified:  dateModified,
		Description:   &models.Description{Value: description},
		Image:         &models.Image{Value: image},
		Yield:         &models.Yield{Value: yield},
		PrepTime:      prepTime,
		CookTime:      cookTime,
		Keywords:      &models.Keywords{Values: keywords},
		Ingredients:   &models.Ingredients{Values: ingredients},
		Instructions:  &models.Instructions{Values: instructions},
	}, nil
}
