package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeSoutherncastiron(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Yield.Value = findYield(root.Find("div[itemprop='description']").Text())

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, found := strings.Cut(name, " - ")
	if found {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	rs.Category.Value = strings.TrimSpace(root.Find(".td-crumb-container a").Last().Text())

	prep := root.Find(".recipe-legend").First().Prev().Text()
	split := strings.Split(prep, " ")
	isMin := strings.Contains(prep, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			prep = "PT" + split[i] + "M"
		}
	}
	rs.PrepTime = prep

	nodes := root.Find("li[itemprop='ingredients']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.Join(strings.Fields(sel.Text()), " ")
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find("li[itemprop='recipeInstructions']")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		s = strings.Join(strings.Fields(s), " ")
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	return rs, nil
}
