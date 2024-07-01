package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeLoveAndLemons(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = strings.TrimSpace(root.Find("h1.entry-title").Text())
	rs.Keywords.Values = getItempropContent(root, "keywords")
	rs.Category.Value = getItempropContent(root, "recipeCategory")
	rs.Yield.Value = findYield(root.Find("span[itemprop='recipeYield']").Text())

	root.Find(".ERSIngredients").First().Children().Next().Each(func(_ int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "ul" {
			sel.Find("li").Each(func(_ int, li *goquery.Selection) {
				rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(li.Text()))
			})
		} else {
			s := strings.TrimSpace(sel.Text())
			if s != "" {
				rs.Ingredients.Values = append(rs.Ingredients.Values, s)
			}
		}
	})

	getInstructions(&rs, root.Find("li[itemprop='recipeInstructions']"))

	notes := strings.TrimSpace(root.Find(".ERSNotesDiv").Text())
	if notes != "" {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(notes, &models.HowToItem{Name: "notes"}))
	}

	return rs, nil
}
