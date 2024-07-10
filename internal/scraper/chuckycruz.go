package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
	"time"
)

func scrapeChuckycruz(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getPropertyContent(root, "og:title")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Image.Value = root.Find("picture img").First().AttrOr("src", "")

	parse, err := time.Parse("Jan 01, 2006", root.Find(".pencraft.pc-display-flex.pc-gap-4.pc-reset").Last().Text())
	if err == nil {
		rs.DatePublished = parse.Format(time.DateOnly)
	}

	var yield int16
	root.Find("p").Each(func(_ int, sel *goquery.Selection) {
		s := strings.ToLower(sel.Text())
		if yield == 0 && strings.HasPrefix(s, "makes") {
			yield = findYield(regex.Digit.FindString(s))
		} else if strings.HasPrefix(s, "time") {
			if !strings.Contains(s, "hour") {
				rs.PrepTime = "PT" + regex.Digit.FindString(s) + "M"
			}
		}
	})
	rs.Yield.Value = yield

	getIngredients(&rs, root.Find("ul li"))
	getInstructions(&rs, root.Find("ol li"))

	return rs, nil
}
