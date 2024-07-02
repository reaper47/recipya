package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
	"time"
)

func scrapeAllClad(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getNameContent(root, "title")
	rs.Description.Value = getNameContent(root, "description")
	rs.Keywords.Values = getNameContent(root, "description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Category.Value, _ = root.Find(".post-categories a").First().Attr("title")
	getIngredients(&rs, root.Find("h2:contains('Ingredients')").Next().Find("ul li"))
	getInstructions(&rs, root.Find("h2:contains('Directions')").Next().Find("ol li"))

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

	return rs, nil
}
