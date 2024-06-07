package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeWholefoodsmarket(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[itemprop='headline']").Attr("content")
	rs.DateModified, _ = root.Find("meta[itemprop='dateModified']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[itemprop='datePublished']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[itemprop='image']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[itemprop='description']").Attr("content")

	p := root.Find(".image-subtitle p").Last().Text()
	for _, s := range strings.Split(p, "|") {
		if strings.Contains(strings.ToLower(s), "serves") {
			rs.Yield.Value = findYield(s)
		}
	}

	nodes := root.Find("h4:contains('Ingredients')").Parent().Find("p")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.Join(strings.Fields(sel.Text()), " ")
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find("h4:contains('Method')").ParentsUntil(".sqs-col-6").Last().Parent().Find("p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.Join(strings.Fields(sel.Text()), " ")
		if s != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}

	})

	return rs, nil
}
