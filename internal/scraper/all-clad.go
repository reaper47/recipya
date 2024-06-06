package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
	"time"
)

func scrapeAllClad(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[name='title']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	rs.Keywords.Values, _ = root.Find("meta[name='description']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Category.Value, _ = root.Find(".post-categories a").First().Attr("title")

	parse, err := time.Parse("January 02, 2006", root.Find(".post-date span").Last().Text())
	if err == nil {
		rs.DatePublished = parse.Format(time.DateOnly)
	}

	meta := strings.Split(root.Find("div:contains('SERVES')").Last().Text(), "\n")
	if len(meta) > 1 {
		for i := 0; i < len(meta)-1; i += 2 {
			v := strings.TrimSpace(meta[i+1])
			switch strings.TrimSpace(meta[i]) {
			case "SERVES":
				rs.Yield.Value = findYield(v)
			case "PREP TIME":
				before, _, ok := strings.Cut(v, "MIN")
				if ok {
					v = strings.TrimSpace(before)
				}
				rs.PrepTime = "PT" + v + "M"
			case "COOK TIME":
				before, _, ok := strings.Cut(v, "MIN")
				if ok {
					v = strings.TrimSpace(before)
				}
				rs.CookTime = "PT" + v + "M"
			}
		}
	}

	nodes := root.Find("h2:contains('Ingredients')").Next().Find("ul li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, sel.Text())
	})

	nodes = root.Find("h2:contains('Directions')").Next().Find("ol li")
	rs.Instructions.Values = make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	return rs, nil
}
