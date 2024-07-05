package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeLekkerenSimpel(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = root.Find(".hero__title").Text()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Category.Value = strings.ToLower(strings.TrimSpace(root.Find(".fa-utensils").First().Parent().Next().Text()))
	rs.Yield.Value = findYield(strings.TrimSpace(root.Find(".fa-user").Last().Parent().Parent().Text()))

	getIngredients(&rs, root.Find("div.recipe__necessities ul").First().Find("li"), []models.Replace{
		{"useFields", ""},
	}...)

	getInstructions(&rs, root.Find(".entry__content p"))

	times := root.Find("span:contains('bereidingstijd')").First().Text()
	if times == "" {
		times = root.Find("span:contains('oventijd')").First().Text()
	}

	if times != "" {
		parts := strings.Split(strings.TrimSpace(times), "  ")

		unit := "H"
		if strings.Contains(parts[0], "min") {
			unit = "M"
		}

		if len(parts) == 2 {
			str := "PT" + regex.Digit.FindString(parts[0]) + unit
			if strings.Contains(parts[0], "bereidingstijd") {
				rs.PrepTime = str
			} else {
				rs.CookTime = str
			}

			unit = "H"
			if strings.Contains(parts[1], "min") {
				unit = "M"
			}
			str = "PT" + regex.Digit.FindString(parts[1]) + unit
			if strings.Contains(parts[1], "bereidingstijd") {
				rs.PrepTime = str
			} else {
				rs.CookTime = str
			}
		} else {
			str := "PT" + regex.Digit.FindString(parts[0]) + unit
			if strings.Contains(parts[0], "bereidingstijd") {
				rs.PrepTime = str
			} else {
				rs.CookTime = str
			}
		}
	}

	return rs, nil
}
