package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeMeljoulwan(root *goquery.Document) (models.RecipeSchema, error) {
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")

	nodes := root.Find("h5:contains('Ingredients')").Parent().Find("li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, strings.TrimSpace(s))
	})

	nodes = root.Find("h5:contains('Directions')").Parent().Find("p")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, models.NewHowToStep(sel.Text()))
	})

	var category string
	root.Find("div.post-category a").Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "Blog" && category == "" {
			category = s
		}
	})

	var sb strings.Builder
	root.Find("div.post-tage a").Each(func(_ int, sel *goquery.Selection) {
		sb.WriteString(sel.Text())
		sb.WriteString(",")
	})
	keywords := strings.TrimSuffix(sb.String(), ",")

	return models.RecipeSchema{
		Category:      &models.Category{Value: category},
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   &models.Description{Value: description},
		Keywords:      &models.Keywords{Values: keywords},
		Image:         &models.Image{Value: image},
		Ingredients:   &models.Ingredients{Values: ingredients},
		Instructions:  &models.Instructions{Values: instructions},
		Name:          name,
		Yield:         &models.Yield{Value: findYield(root.Find("p:contains('Serves')").Text())},
	}, nil
}
