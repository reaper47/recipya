package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeJuliegoodwin(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, " - ")
	if ok {
		name = strings.TrimSpace(before)
	}

	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	description := root.Find(".divider-wrap").First().Next().Text()

	times := root.Find("strong:contains('Prep time')").Parent().Text()
	var (
		cook string
		prep string
	)
	if times != "" {
		split := strings.Split(times, "\n")
		if len(split) == 2 {
			s := split[0]
			split = strings.Split(s, " ")
			isMin := strings.Contains(s, "min")
			for i, s := range split {
				_, err := strconv.ParseInt(s, 10, 64)
				if err == nil && isMin {
					prep = "PT" + split[i] + "M"
				}
			}

			if len(split) > 1 {
				split = strings.Split(split[1], " ")
				if len(split) > 0 {
					isMin = strings.Contains(split[0], "min")
					for i, s := range split {
						_, err := strconv.ParseInt(s, 10, 64)
						if err == nil && isMin {
							cook = "PT" + split[i] + "M"
						}
					}
				}
			}
		}
	}

	ingredients := strings.Split(root.Find("h4:contains('Ingredients')").Next().Text(), "\n")
	for i, ingredient := range ingredients {
		_, after, found := strings.Cut(ingredient, "•")
		if found {
			ingredients[i] = strings.TrimSpace(after)
		}
	}

	nodes := root.Find("h4:contains('Method')").Parent().Find("p")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		_, err := strconv.ParseInt(s[:1], 10, 64)
		if err == nil {
			_, after, _ := strings.Cut(s, ".")
			after = strings.Join(strings.Fields(after), " ")
			s = strings.TrimSpace(after)
		}
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: root.Find(".meta-category a").First().Text()},
		CookTime:      cook,
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		PrepTime:      prep,
		Yield:         models.Yield{Value: findYield(root.Find("i.fa-cutlery").Parent().Text())},
	}, nil
}
