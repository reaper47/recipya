package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeRecettesDuQuebec(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DatePublished = getNameContent(root, "cXenseParse:recs:publishtime")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Name = getNameContent(root, "cXenseParse:title")

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

	getIngredients(&rs, root.Find(".ingredients ul label"), []models.Replace{
		{"\n", ""},
		{"useFields", ""},
	}...)

	getInstructions(&rs, root.Find(".method p"))

	image, _ := root.Find("picture img").Attr("srcset")
	split := strings.Split(image, "?")
	if len(split) > 0 {
		rs.Image.Value = split[0]
	}

	return rs, nil
}
