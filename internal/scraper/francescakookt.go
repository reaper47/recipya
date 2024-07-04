package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeFrancescakookt(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getNameContent(root, "description")
	rs.Name = getPropertyContent(root, "og:title")
	rs.Image.Value = getPropertyContent(root, "og:image")

	nodes := root.Find("meta[property='article:tag']")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		v, _ := sel.Attr("content")
		keywords = append(keywords, v)
	})
	rs.Keywords.Values = strings.ToLower(strings.Join(keywords, ","))

	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.DatePublished = getPropertyContent(root, "article:published_time")

	category := root.Find("h2:contains('Categorieën')").Next().Find("a").First().Text()
	before, _, ok := strings.Cut(category, "recepten")
	if ok {
		category = strings.TrimSpace(before)
	}
	rs.Category.Value = category

	var yield int16
	root.Find("h3:contains('personen')").Each(func(_ int, sel *goquery.Selection) {
		if yield != 0 {
			return
		}
		yield = findYield(sel.Text())
	})
	rs.Yield.Value = yield

	getIngredients(&rs, root.Find("div.dynamic-entry-content ul").Last().Find("li"))
	getInstructions(&rs, root.Find("div.dynamic-entry-content ol").Last().Find("li"))

	s := root.Find("p:contains('Bereidingstijd')").Text()
	if s != "" {
		parts := strings.Split(s, " ")
		if len(parts) == 3 && strings.Contains(s, "minuten") {
			rs.PrepTime = "PT" + regex.Digit.FindString(s) + "M"
		}
	}

	return rs, nil
}
