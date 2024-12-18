package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeKuchniadomova(root *goquery.Document) (models.RecipeSchema, error) {
	rs, _ := parseWebsite(root)

	yieldStr := root.Find("p[itemprop='recipeYield']").Text()
	yieldStr = strings.ReplaceAll(yieldStr, "-", " ")
	rs.Yield.Value = findYield(yieldStr)

	rs.Category.Value = getItempropContent(root, "recipeCategory")
	if rs.Category.Value == "" {
		rs.Category.Value = strings.TrimSpace(root.Find("ol.breadcrumb li").Last().Text())
	}

	rs.Cuisine.Value = root.Find("p[itemprop=recipeCuisine]").Text()
	getInstructions(&rs, root.Find("#recipe-instructions li"))

	return rs, nil
}
