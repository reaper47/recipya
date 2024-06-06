package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeUsapears(root *goquery.Document) (models.RecipeSchema, error) {
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")

	prep := root.Find(".recipe-legend").First().Prev().Text()
	split := strings.Split(prep, " ")
	isMin := strings.Contains(prep, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			prep = "PT" + split[i] + "M"
		}
	}

	nodes := root.Find("li[itemprop='ingredients']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.Join(strings.Fields(sel.Text()), " ")
		ingredients = append(ingredients, s)
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] ol li")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		s = strings.Join(strings.Fields(s), " ")
		instructions = append(instructions, models.NewHowToStep(s))
	})

	return models.RecipeSchema{
		DatePublished: datePublished,
		Description:   &models.Description{Value: description},
		Image:         &models.Image{Value: image},
		Ingredients:   &models.Ingredients{Values: ingredients},
		Instructions:  &models.Instructions{Values: instructions},
		Name:          name,
		PrepTime:      prep,
	}, nil
}
