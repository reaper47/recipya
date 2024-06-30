package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeAfghanKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	content := root.Find("#content")
	info := content.Find(".recipe-info")
	about := content.Find("div[itemprop='about']")
	rs.Yield.Value = findYield(info.Find(".servings .value").Text())

	time := info.Find(".prep-time .value").Text()
	if strings.Contains(time, "m") {
		rs.PrepTime = "PT" + strings.TrimSuffix(time, "m") + "M"
	} else if strings.Contains(time, "h") {
		time = strings.TrimSuffix(time, " h")
		parts := strings.Split(time, ":")
		if len(parts) > 1 {
			rs.PrepTime = "PT" + parts[0] + "H" + parts[1] + "M"
		}
	}

	time = info.Find(".cook-time .value").Text()
	if strings.Contains(time, "m") {
		rs.CookTime = "PT" + strings.TrimSuffix(time, "m") + "M"
	} else if strings.Contains(time, "h") {
		time = strings.TrimSuffix(time, " h")
		parts := strings.Split(time, ":")
		if len(parts) > 1 {
			rs.CookTime = "PT" + parts[0] + "H" + parts[1] + "M"
		}
	}

	if len(about.Nodes) > 0 && about.Nodes[0].FirstChild != nil && about.Nodes[0].FirstChild.NextSibling != nil && about.Nodes[0].FirstChild.NextSibling.FirstChild != nil {
		s := about.Nodes[0].FirstChild.NextSibling.FirstChild.Data
		s = strings.ReplaceAll(s, "\n", "")
		s = strings.ReplaceAll(s, "\u00a0", " ")
		rs.Description.Value = strings.TrimSpace(s)
	}

	getIngredients(&rs, about.Find("li.ingredient"), []models.Replace{{"  ", " "}}...)
	getInstructions(&rs, about.Find("p.instructions"), []models.Replace{{"  ", " "}}...)

	rs.DatePublished, _ = content.Find("meta[itemprop='datePublished']").Attr("content")
	rs.Image.Value, _ = content.Find("meta[itemprop='image']").Attr("content")
	rs.Name = content.Find("h2[itemprop='name']").Text()
	rs.Category = &models.Category{Value: about.Find(".type a").Text()}
	return rs, nil
}
