package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeRecettesDuQuebec(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DatePublished, _ = root.Find("meta[name='cXenseParse:recs:publishtime']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Name, _ = root.Find("meta[name='cXenseParse:title']").Attr("content")

	category := root.Find(".categories h6").Siblings().First().Text()
	rs.Category.Value = strings.TrimSpace(category)

	keywords := root.Find(".tags h6").Siblings().First().Text()
	rs.Keywords.Values = strings.TrimSpace(keywords)

	root.Find("span.cat").Each(func(_ int, sel *goquery.Selection) {
		switch sel.Text() {
		case "Préparation":
			rs.PrepTime = "PT"
			parts := strings.Split(sel.Siblings().First().Text(), "&")
			if len(parts) == 2 {
				rs.PrepTime += strings.Split(parts[0], " ")[0] + "H"
				rs.PrepTime += strings.Split(parts[1], " ")[1] + "M"
			} else {
				xs := strings.Split(parts[0], " ")
				rs.PrepTime += xs[0]
				if strings.Contains(xs[1], "min") {
					rs.PrepTime += "M"
				} else {
					rs.PrepTime += "H"
				}
			}
		case "Cuisson":
			rs.CookTime = "PT"
			parts := strings.Split(sel.Siblings().First().Text(), "&")
			if len(parts) == 2 {
				rs.CookTime += strings.Split(parts[0], " ")[0] + "H"
				rs.CookTime += strings.Split(parts[1], " ")[1] + "M"
			} else {
				xs := strings.Split(parts[0], " ")
				rs.CookTime += xs[0]
				if strings.Contains(xs[1], "min") {
					rs.CookTime += "M"
				} else {
					rs.CookTime += "H"
				}
			}
		case "Portion(s)":
			yieldStr := sel.Siblings().First().Text()
			for _, s := range strings.Split(yieldStr, " ") {
				yield64, err := strconv.ParseInt(s, 10, 16)
				if err == nil {
					rs.Yield.Value = int16(yield64)
					break
				}
			}
		}
	})

	nodes := root.Find(".ingredients ul label")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		s = strings.ReplaceAll(s, "\n", "")
		s = strings.Join(strings.Fields(s), " ")
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find(".method p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	image, _ := root.Find("picture img").Attr("srcset")
	split := strings.Split(image, "?")
	if len(split) > 0 {
		rs.Image.Value = split[0]
	}

	return rs, nil
}
