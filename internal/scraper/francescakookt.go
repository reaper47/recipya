package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeFrancescakookt(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	nodes := root.Find("meta[property='article:tag']")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		v, _ := sel.Attr("content")
		keywords = append(keywords, v)
	})
	rs.Keywords.Values = strings.ToLower(strings.Join(keywords, ","))

	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")

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

	nodes = root.Find("div.dynamic-entry-content ul").Last().Find("li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	s := root.Find("p:contains('Bereidingstijd')").Text()
	if s != "" {
		parts := strings.Split(s, " ")
		if len(parts) == 3 && strings.Contains(s, "minuten") {
			rs.PrepTime = "PT" + regex.Digit.FindString(s) + "M"
		}
	}

	nodes = root.Find("div.dynamic-entry-content ol").Last().Find("li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	return rs, nil
}
