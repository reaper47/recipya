package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
	"time"
)

func scrapeAllClad(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[name='title']").Attr("content")
	description, _ := root.Find("meta[name='description']").Attr("content")
	keywords, _ := root.Find("meta[name='description']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")
	category, _ := root.Find(".post-categories a").First().Attr("title")

	var datePub string
	parse, err := time.Parse("January 02, 2006", root.Find(".post-date span").Last().Text())
	if err == nil {
		datePub = parse.Format(time.DateOnly)
	}

	var (
		cook  string
		prep  string
		yield int16
	)

	meta := strings.Split(root.Find("div:contains('SERVES')").Last().Text(), "\n")
	if len(meta) > 1 {
		for i := 0; i < len(meta)-1; i += 2 {
			v := strings.TrimSpace(meta[i+1])
			switch strings.TrimSpace(meta[i]) {
			case "SERVES":
				yield = findYield(v)
			case "PREP TIME":
				before, _, ok := strings.Cut(v, "MIN")
				if ok {
					v = strings.TrimSpace(before)
				}
				prep = "PT" + v + "M"
			case "COOK TIME":
				before, _, ok := strings.Cut(v, "MIN")
				if ok {
					v = strings.TrimSpace(before)
				}
				cook = "PT" + v + "M"
			}
		}
	}

	nodes := root.Find("h2:contains('Ingredients')").Next().Find("ul li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, sel.Text())
	})

	nodes = root.Find("h2:contains('Directions')").Next().Find("ol li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, sel.Text())
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: category},
		CookTime:      cook,
		DatePublished: datePub,
		Description:   models.Description{Value: description},
		Keywords:      models.Keywords{Values: keywords},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		PrepTime:      prep,
		Yield:         models.Yield{Value: yield},
	}, nil
}
