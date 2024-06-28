package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
	"strings"
)

func scrapeUitpaulineskeuken(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	s := strings.TrimSpace(root.Find(".fa-stopwatch").First().Parent().Text())
	if s != "" {
		parts := strings.Split(s, "+")
		if len(parts) == 2 {
			v := regex.Digit.FindString(parts[0])
			rs.PrepTime = "PT" + v + "M"

			left, err := strconv.Atoi(v)
			if err == nil {
				right, err := strconv.Atoi(regex.Digit.FindString(parts[1]))
				if err == nil {
					rem := right - left
					if rem > 0 {
						rs.CookTime = "PT" + strconv.Itoa(rem) + "M"
					}
				}
			}
		}
	}

	nodes := root.Find("#ingredienten ul li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("#recept ol li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	nodes = root.Find("#gerelateerd a")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(strings.ToLower(sel.Text())))
	})
	rs.Keywords.Values = strings.Join(keywords, ",")

	return rs, nil
}
