package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeALittleBitYummy(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	getInstructions(&rs, root.Find("#method ol li"))
	rs.TotalTime = ""

	if rs.Tools == nil {
		rs.Tools = &models.Tools{}
	}
	root.Find("ul.recipe-equipment-list").First().Find("ul li").Each(func(_ int, sel *goquery.Selection) {
		rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(sel.Text()))
	})

	root.Find("div.ingredients-tab-content").First().Children().Each(func(_ int, sel *goquery.Selection) {
		name := goquery.NodeName(sel)
		if name == "h4" || name == "div" {
			s := strings.Join(strings.Fields(sel.Text()), " ")
			rs.Ingredients.Values = append(rs.Ingredients.Values, s)
		}
	})

	for i, s := range rs.Instructions.Values {
		rs.Instructions.Values[i].Text = strings.Join(strings.Fields(s.Text), " ")
	}

	return rs, nil
}
