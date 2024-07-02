package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeEatwell101(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getNameContent(root, "description")
	rs.Keywords.Values = getNameContent(root, "keywords")
	rs.Name = getPropertyContent(root, "og:title")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	getIngredients(&rs, root.Find("h2:contains('Ingredients')").Next().Find("li"))

	nodes := root.Find("h2:contains('Directions')")
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
