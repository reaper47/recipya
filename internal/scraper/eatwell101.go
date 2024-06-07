package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeEatwell101(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	rs.Keywords.Values, _ = root.Find("meta[name='keywords']").Attr("content")
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	nodes := root.Find("h2:contains('Ingredients')").Next().Find("li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find("h2:contains('Directions')")
	for {
		nodes = nodes.Next()
		s := nodes.Text()
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		if goquery.NodeName(nodes) != "p" {
			break
		}
	}

	// The below does not work.
	div := root.Find("#recipecardo p.brandon")

	prep := root.Find("i:contains('Prep Time')").Next().Text()
	split := strings.Split(prep, " ")
	isMin := strings.Contains(prep, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			prep = "PT" + split[i] + "M"
			break
		}
	}

	cook := root.Find("i:contains('Prep Time')").Next().Text()
	split = strings.Split(cook, " ")
	isMin = strings.Contains(cook, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			cook = "PT" + split[i] + "M"
			break
		}
	}

	rs.CookTime = cook
	rs.PrepTime = prep
	rs.Yield.Value = findYield(div.Find("span:contains('servings')").Text())

	return rs, nil
}
