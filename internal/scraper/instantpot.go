package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeInstantPot(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getPropertyContent(root, "og:title")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Image.Value = getPropertyContent(root, "og:image:secure_url")

	overview := root.Find(".article__overview").Children()
	overview.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if strings.HasPrefix(s, "Course") {
			rs.Category.Value = strings.TrimSpace(strings.ToLower(strings.TrimPrefix(s, "Course")))
		} else if strings.HasPrefix(s, "Keywords") {
			rs.Keywords.Values = strings.Join(strings.Fields(strings.TrimPrefix(s, "Keywords")), "")
		} else if strings.HasPrefix(s, "Prep Time") {
			prep := "PT" + regex.Digit.FindString(s)
			if strings.Contains(s, "min") {
				prep += "M"
			} else {
				prep += "H"
			}
			rs.PrepTime = prep
		} else if strings.HasPrefix(s, "Cook Time") {
			cook := "PT" + regex.Digit.FindString(s)
			if strings.Contains(s, "min") {
				cook += "M"
			} else {
				cook += "H"
			}
			rs.CookTime = cook
		} else if strings.HasPrefix(s, "Servings") {
			rs.Yield.Value = findYield(s)
		}
	})

	if rs.PrepTime == "" {
		s := overview.Find("div:contains('Duration')").Next().Text()
		prep := "PT" + regex.Digit.FindString(s)
		if strings.Contains(s, "min") {
			prep += "M"
		} else {
			prep += "H"
		}
		rs.PrepTime = prep
	}

	getIngredients(&rs, root.Find(".article__ingredients").First().Find("li"))
	getInstructions(&rs, root.Find(".article__instructions").First().Find("li"))

	return rs, nil
}
