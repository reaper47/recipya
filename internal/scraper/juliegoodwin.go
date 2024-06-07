package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeJuliegoodwin(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(rs.Name, " - ")
	if ok {
		rs.Name = strings.TrimSpace(before)
	}

	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Description.Value = root.Find(".divider-wrap").First().Next().Text()
	rs.Yield.Value = findYield(root.Find("i.fa-cutlery").Parent().Text())
	rs.Category.Value = root.Find(".meta-category a").First().Text()

	times := root.Find("strong:contains('Prep time')").Parent().Text()
	if times != "" {
		split := strings.Split(times, "\n")
		if len(split) == 2 {
			s := split[0]
			split = strings.Split(s, " ")
			isMin := strings.Contains(s, "min")
			for i, s := range split {
				_, err := strconv.ParseInt(s, 10, 64)
				if err == nil && isMin {
					rs.PrepTime = "PT" + split[i] + "M"
				}
			}

			if len(split) > 1 {
				split = strings.Split(split[1], " ")
				if len(split) > 0 {
					isMin = strings.Contains(split[0], "min")
					for i, s := range split {
						_, err := strconv.ParseInt(s, 10, 64)
						if err == nil && isMin {
							rs.CookTime = "PT" + split[i] + "M"
						}
					}
				}
			}
		}
	}

	rs.Ingredients.Values = strings.Split(root.Find("h4:contains('Ingredients')").Next().Text(), "\n")
	for i, ingredient := range rs.Ingredients.Values {
		_, after, found := strings.Cut(ingredient, "•")
		if found {
			rs.Ingredients.Values[i] = strings.TrimSpace(after)
		}
	}

	nodes := root.Find("h4:contains('Method')").Parent().Find("p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		_, err := strconv.ParseInt(s[:1], 10, 64)
		if err == nil {
			_, after, _ := strings.Cut(s, ".")
			after = strings.Join(strings.Fields(after), " ")
			s = strings.TrimSpace(after)
		}
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	return rs, nil
}
