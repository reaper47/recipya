package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeKwestiasmaku(root *goquery.Document) (models.RecipeSchema, error) {
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	description := root.Find("span[itemprop='description']").Text()

	nodes := root.Find(".field-name-field-skladniki li")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\n", "")
		v = strings.ReplaceAll(v, "\t", "")
		ingredients[i] = v
	})

	nodes = root.Find(".field-name-field-przygotowanie li")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\n", "")
		v = strings.ReplaceAll(v, "\t", "")
		instructions = append(instructions, models.NewHowToStep(v))
	})

	yield := findYield(root.Find(".field-name-field-ilosc-porcji").Text())

	return models.RecipeSchema{
		Image:         &models.Image{Value: image},
		DatePublished: datePublished,
		DateModified:  dateModified,
		Description:   &models.Description{Value: description},
		Name:          name,
		Yield:         &models.Yield{Value: yield},
		Ingredients:   &models.Ingredients{Values: ingredients},
		Instructions:  &models.Instructions{Values: instructions},
	}, nil
}
