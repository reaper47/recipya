package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeCuly(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Name = getPropertyContent(root, "og:title")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Yield.Value = findYield(root.Find("span:contains('personen')").Text())
	getIngredients(&rs, root.Find("div.ingredients li"))
	getInstructions(&rs, root.Find("ol li"))

	nodes := root.Find("meta[name='cXenseParse:mhu-article_tag']")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s, _ := sel.Attr("content")
		xk = append(xk, s)
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	prep := strings.TrimSpace(root.Find("span:contains('Voorbereiding')").Parent().Text())
	if prep != "" {
		rs.PrepTime = "PT" + regex.Digit.FindString(prep)
		if strings.Contains(prep, "min") {
			rs.PrepTime += "M"
		} else {
			rs.PrepTime += "H"
		}
	}

	cook := strings.TrimSpace(root.Find("span:contains('Kooktijd')").Parent().Text())
	if cook != "" {
		rs.CookTime = "PT" + regex.Digit.FindString(cook)
		if strings.Contains(cook, "min") {
			rs.CookTime += "M"
		} else {
			rs.CookTime += "H"
		}
	}

	return rs, nil
}
