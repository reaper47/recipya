package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
	"strings"
)

func scrapeUitpaulineskeuken(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	dateMod, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	var (
		prep string
		cook string
	)
	s := strings.TrimSpace(root.Find(".fa-stopwatch").First().Parent().Text())
	if s != "" {
		parts := strings.Split(s, "+")
		if len(parts) == 2 {
			v := regex.Digit.FindString(parts[0])
			prep = "PT" + v + "M"

			left, err := strconv.Atoi(v)
			if err == nil {
				right, err := strconv.Atoi(regex.Digit.FindString(parts[1]))
				if err == nil {
					rem := right - left
					if rem > 0 {
						cook = "PT" + strconv.Itoa(rem) + "M"
					}
				}
			}
		}
	}

	nodes := root.Find("#ingredienten ul li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("#recept ol li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("#gerelateerd a")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(strings.ToLower(sel.Text())))
	})

	return models.RecipeSchema{
		AtContext:    atContext,
		AtType:       models.SchemaType{Value: "Recipe"},
		Category:     models.Category{Value: "uncategorized"},
		CookTime:     cook,
		DateModified: dateMod,
		Description:  models.Description{Value: description},
		Keywords:     models.Keywords{Values: strings.Join(keywords, ",")},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         name,
		PrepTime:     prep,
		Yield:        models.Yield{Value: 1},
	}, nil
}
