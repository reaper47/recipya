package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeAldi(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Category.Value = strings.Replace(root.Find(".tab-nav--item.dropdown--list--item.m-active").Text(), "Recipes", "", 1)
	rs.Category.Value = strings.TrimSpace(strings.Replace(rs.Category.Value, "selected", "", 1))

	content := root.Find("#main-content")
	sidebar := root.Find("#sidebar")

	rs.Description.Value = getNameContent(root, "description")
	rs.Image.Value = content.Find("img").First().AttrOr("src", "")
	rs.Name = strings.TrimSpace(content.Find("h1").Text())
	rs.Yield.Value = findYield(sidebar.Find("p:contains('Serves')").Text())
	getIngredients(&rs, sidebar.Find("ul li"))
	getInstructions(&rs, content.Find("ol li"))
	getTime(&rs, sidebar.Find("p:contains('Prep Time')"), true)

	cook := strings.TrimSpace(sidebar.Find("p:contains('Cook Time')").Text())
	if cook != "" {
		rs.CookTime = "PT" + regex.Digit.FindString(cook)
		if strings.Contains(cook, "min") {
			rs.CookTime += "M"
		} else {
			rs.CookTime += "H"
		}
	}

	return rs, nil
}
