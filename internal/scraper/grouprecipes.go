package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeGrouprecipes(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	rs.Image.Value, _ = root.Find(".photos img").First().Attr("src")
	rs.CookTime, _ = root.Find(".cooktime .value-title").Attr("title")
	rs.Name = root.Find("title").Text()
	rs.Yield.Value = findYield(root.Find(".servings").Text())

	var keywords strings.Builder
	root.Find(".tags_text li").Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		keywords.WriteString(",")
		keywords.WriteString(s)
	})
	rs.Keywords.Values = keywords.String()

	nodes := root.Find(".ingredients li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		before, _, ok := strings.Cut(s, "\t")
		if ok {
			s = strings.TrimSpace(before)
		}
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find(".instructions li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		before, after, ok := strings.Cut(s, ".")
		_, err := strconv.ParseInt(before[:1], 0, 64)
		if ok && err == nil {
			s = strings.TrimSpace(after)
		}
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	return rs, nil
}
