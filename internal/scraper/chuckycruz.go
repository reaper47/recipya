package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
	"time"
)

func scrapeChuckycruz(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	image, _ := root.Find("picture img").First().Attr("src")

	var datePub string
	parse, err := time.Parse("Jan 01, 2006", root.Find(".pencraft.pc-display-flex.pc-gap-4.pc-reset").Last().Text())
	if err == nil {
		datePub = parse.Format(time.DateOnly)
	}

	var (
		prep  string
		yield int16
	)

	root.Find("p").Each(func(_ int, sel *goquery.Selection) {
		s := strings.ToLower(sel.Text())
		if yield == 0 && strings.HasPrefix(s, "makes") {
			yield = findYield(regex.Digit.FindString(s))
		} else if strings.HasPrefix(s, "time") {
			if !strings.Contains(s, "hour") {
				prep = "PT" + regex.Digit.FindString(s) + "M"
			}
		}
	})

	nodes := root.Find("ul li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("ol li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, strings.TrimSpace(sel.Text()))
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		DatePublished: datePub,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		PrepTime:      prep,
		Yield:         models.Yield{Value: yield},
	}, nil
}
