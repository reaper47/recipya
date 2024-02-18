package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeGrimGrains(root *goquery.Document) (models.RecipeSchema, error) {
	h2 := root.Find("h2").Text()
	var (
		prep  string
		yield int16
	)
	before, after, ok := strings.Cut(h2, " — ")
	if ok {
		split := strings.Split(before, " ")
		for _, s := range split {
			parseInt, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				yield = int16(parseInt)
			}
		}

		isMin := strings.Contains(after, "min")
		split = strings.Split(after, " ")
		for i, s := range split {
			_, err := strconv.ParseInt(s, 10, 64)
			if err == nil && isMin {
				prep = "PT" + split[i] + "M"
			}
		}
	}

	image, _ := root.Find("img").First().Attr("src")
	if strings.HasPrefix(image, "../") {
		image = "https://grimgrains.com/" + strings.TrimPrefix(image, "../")
	}

	col2 := root.Find(".col2").Text()

	nodes := root.Find(".ingredients dt")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	nodes = root.Find(".instructions li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		AtContext:    atContext,
		AtType:       models.SchemaType{Value: "Recipe"},
		Description:  models.Description{Value: col2},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         root.Find("h1").Text(),
		PrepTime:     prep,
		Yield:        models.Yield{Value: yield},
	}, nil
}
