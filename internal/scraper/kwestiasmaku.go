package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeKwestiasmaku(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Description.Value = root.Find("span[itemprop='description']").Text()
	rs.Yield.Value = findYield(root.Find(".field-name-field-ilosc-porcji").Text())

	nodes := root.Find(".field-name-field-skladniki li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\n", "")
		v = strings.ReplaceAll(v, "\t", "")
		rs.Ingredients.Values = append(rs.Ingredients.Values, v)
	})

	nodes = root.Find(".field-name-field-przygotowanie li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\n", "")
		v = strings.ReplaceAll(v, "\t", "")
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(v))
	})

	return rs, nil
}
