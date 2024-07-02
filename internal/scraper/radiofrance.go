package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeRadioFrance(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Name = strings.TrimSpace(strings.TrimPrefix(root.Find(".commonArticle-title").Text(), "RECETTE -"))

	meta := root.Find(".PopSlotParagraph").Text()
	for _, s := range strings.Split(meta, "|") {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "Pour") {
			rs.NutritionSchema.Servings = strings.TrimSpace(strings.TrimPrefix(s, "Pour"))
		} else if strings.HasPrefix(s, "Préparation") {
			rs.PrepTime = "PT" + regex.Digit.FindString(s)
			if strings.Contains(s, "m") {
				rs.PrepTime += "M"
			} else {
				rs.PrepTime += "H"
			}
		} else if strings.HasPrefix(s, "Cuisson") {
			rs.CookTime = "PT" + regex.Digit.FindString(s)
			if strings.Contains(s, "m") {
				rs.CookTime += "M"
			} else {
				rs.CookTime += "H"
			}
		}
	}

	node := root.Find("ul").FilterFunction(func(_ int, sel *goquery.Selection) bool {
		_, ok := sel.Attr("class")
		return !ok
	})
	getIngredients(&rs, node.Find("li"))
	getInstructions(&rs, node.NextAll())

	return rs, nil
}
