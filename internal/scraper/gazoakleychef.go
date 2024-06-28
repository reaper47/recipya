package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
	"strings"
)

func scrapeGazoakleychef(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = root.Find(".entry-title").Text()

	root.Find(".entry-quick-info div.row div").Each(func(_ int, sel *goquery.Selection) {
		c := strings.TrimSpace(sel.Text())
		if strings.HasPrefix(strings.ToLower(c), "serve") {
			atoi, err := strconv.ParseInt(regex.Digit.FindString(c), 10, 16)
			if err == nil {
				rs.Yield.Value = int16(atoi)
			}
		} else if strings.HasPrefix(strings.ToLower(c), "cooks in:") {
			parts := strings.Split(strings.TrimPrefix(c, "cooks in:"), ",")
			for _, part := range parts {
				num := regex.Digit.FindString(part)
				if num == "" {
					continue
				}

				minutes := "PT" + num + "M"
				hours := "PT" + num + "H"

				if strings.Contains(strings.ToLower(part), "prep") {
					if strings.Contains(part, "min") {
						rs.PrepTime = minutes
					} else if strings.Contains(part, "hour") {
						rs.PrepTime = hours
					}
				} else if strings.Contains(strings.ToLower(part), "cooking") {
					if strings.Contains(part, "min") {
						rs.CookTime = minutes
					} else if strings.Contains(part, "hour") {
						rs.CookTime = hours
					}
				} else {
					if strings.Contains(part, "min") {
						rs.PrepTime = minutes
					} else if strings.Contains(part, "hour") {
						rs.PrepTime = hours
					}
				}
			}
		}
	})

	rs.Description.Value = strings.TrimSpace(root.Find(".recipe-introduction").Text())
	rs.Category.Value = root.Find(".entry-recipe-categories a").First().Text()

	nodes := root.Find(".recipe-ingredients p")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find(".recipe-method p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	})

	nodes = root.Find(".entry-share-ingredients a")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "" {
			keywords = append(keywords, s)
		}
	})
	rs.Keywords.Values = strings.Join(keywords, ",")

	return rs, nil
}
