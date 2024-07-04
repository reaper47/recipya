package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeWholefoodsmarket(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getItempropContent(root, "headline")
	rs.DateModified = getItempropContent(root, "dateModified")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.Image.Value = getItempropContent(root, "image")
	rs.Description.Value = getItempropContent(root, "description")

	p := root.Find(".image-subtitle p").Last().Text()
	for _, s := range strings.Split(p, "|") {
		if strings.Contains(strings.ToLower(s), "serves") {
			rs.Yield.Value = findYield(s)
		}
	}

	getIngredients(&rs, root.Find("h4:contains('Ingredients')").Parent().Find("p"), []models.Replace{{"useFields", ""}}...)
	getInstructions(&rs, root.Find("h4:contains('Method')").ParentsUntil(".sqs-col-6").Last().Parent().Find("p"))

	return rs, nil
}
