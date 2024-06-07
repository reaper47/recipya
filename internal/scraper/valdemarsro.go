package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeValdemarsro(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Name = root.Find("h1[itemprop='headline']").Text()
	rs.Yield.Value = findYield(root.Find(".fa-sort").Parent().Text())

	start := root.Find("div[itemprop='description']").Children().First()
	rs.Description.Value = start.NextUntil(".post-recipe").FilterFunction(func(_ int, selection *goquery.Selection) bool {
		return selection.Nodes[0].Data == "p"
	}).Text()

	prepTimeStr := root.Find("span:contains('Tid i alt')").Next().Text()
	parts := strings.Split(prepTimeStr, " ")
	switch len(parts) {
	case 2:
		rs.PrepTime = "PT" + parts[0] + "M"
	case 4:
		rs.PrepTime = "PT" + parts[0] + "H" + parts[2] + "M"
	}

	cookTimeStr := root.Find("span:contains('Arbejdstid')").Next().Text()
	parts = strings.Split(cookTimeStr, " ")
	switch len(parts) {
	case 2:
		rs.CookTime = "PT" + parts[0] + "M"
	case 4:
		rs.CookTime = "PT" + parts[0] + "H" + parts[2] + "M"
	}

	nodes := root.Find("li[itemprop='recipeIngredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, s.Text())
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s.Text()))
	})

	return rs, nil
}
