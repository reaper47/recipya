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

	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Image.Value, _ = root.Find("picture img").First().Attr("src")

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

	nodes := root.Find("ul li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("ol li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(sel.Text())))
	})

	return rs, nil
}
