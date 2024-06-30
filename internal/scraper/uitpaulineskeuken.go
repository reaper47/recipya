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

	rs.Name = getPropertyContent(root, "og:title")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")

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

	getIngredients(&rs, root.Find("#ingredienten ul li"))
	getInstructions(&rs, root.Find("#recept ol li"))

	nodes := root.Find("#gerelateerd a")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(strings.ToLower(sel.Text())))
	})
	rs.Keywords.Values = strings.Join(keywords, ",")

	return rs, nil
}
