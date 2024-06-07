package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeGrimGrains(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	h2 := root.Find("h2").Text()

	before, after, ok := strings.Cut(h2, " — ")
	if ok {
		split := strings.Split(before, " ")
		for _, s := range split {
			parseInt, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				rs.Yield.Value = int16(parseInt)
			}
		}

		isMin := strings.Contains(after, "min")
		split = strings.Split(after, " ")
		for i, s := range split {
			_, err := strconv.ParseInt(s, 10, 64)
			if err == nil && isMin {
				rs.PrepTime = "PT" + split[i] + "M"
			}
		}
	}

	rs.Image.Value, _ = root.Find("img").First().Attr("src")
	if strings.HasPrefix(rs.Image.Value, "../") {
		rs.Image.Value = "https://grimgrains.com/" + strings.TrimPrefix(rs.Image.Value, "../")
	}

	rs.Description.Value = root.Find(".col2").Text()
	rs.Name = root.Find("h1").Text()

	nodes := root.Find(".ingredients dt")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find(".instructions li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	return rs, nil
}
