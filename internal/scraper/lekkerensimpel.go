package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
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

	times, _ := root.Find(".entry__container .recipe__meta span:contains('bereidingstijd')").Html()
	parts := strings.Split(times, "<br/>")
	if len(parts) == 1 {
		timeParts := strings.Split(parts[0], " ")
		var time string
		if len(timeParts) == 3 {
			time = "PT" + timeParts[0] + "M"
		} else {
			time = "PT" + timeParts[0] + "H" + timeParts[2] + "M"
		}

		if strings.Contains(parts[0], "oven") {
			rs.CookTime = time
		} else {
			rs.PrepTime = time
		}
	} else if len(parts) == 2 {
		prepTime := strings.Join(strings.Fields(parts[0]), " ")
		prepParts := strings.Split(prepTime, " ")
		if len(prepParts) == 3 {
			rs.PrepTime = "PT" + prepParts[0] + "M"
		} else {
			rs.PrepTime = "PT" + prepParts[0] + "H" + prepParts[2] + "M"
		}

		cookTime := strings.Join(strings.Fields(parts[1]), " ")
		cookParts := strings.Split(cookTime, " ")
		if len(cookParts) == 3 {
			rs.CookTime = "PT" + cookParts[0] + "M"
		} else {
			rs.CookTime = "PT" + cookParts[0] + "H" + cookParts[2] + "M"
		}
	}

	return rs, nil
}
