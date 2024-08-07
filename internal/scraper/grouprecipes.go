package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeGrouprecipes(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getNameContent(root, "description")
	rs.Image.Value = root.Find(".photos img").First().AttrOr("src", "")
	rs.CookTime = root.Find(".cooktime .value-title").AttrOr("title", "")
	rs.Name = root.Find("title").Text()
	rs.Yield.Value = findYield(root.Find(".servings").Text())

	var keywords strings.Builder
	root.Find(".tags_text li").Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		keywords.WriteString(",")
		keywords.WriteString(s)
	})
	rs.Keywords.Values = keywords.String()

	getIngredients(&rs, root.Find(".ingredients li"), []models.Replace{{"shopping list", ""}}...)

	nodes := root.Find(".instructions li")
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
