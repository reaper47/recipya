package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapePotatoRolls(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = root.Find("[itemprop=name]").Text()
	rs.Image.Value = root.Find("[itemprop=image]").First().AttrOr("src", "")

	getIngredients(&rs, root.Find(".ingredient"))
	getInstructions(&rs, root.Find(".direction p"))

	prep := strings.TrimSpace(root.Find(".icon-clock").First().Parent().Text())
	if prep != "" {
		rs.PrepTime = "PT" + regex.Digit.FindString(prep)
		if strings.Contains(strings.ToLower(prep), "min") {
			rs.PrepTime += "M"
		} else {
			rs.PrepTime += "H"
		}
	}
	return rs, nil
}
