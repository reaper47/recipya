package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"strings"
)

func scrapeChefnini(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	nodes := root.Find("meta[property='article:tag']")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s, _ := sel.Attr("content")
		xk = append(xk, s)
	})
	rs.Keywords.Values = strings.Join(extensions.Unique(xk), ",")

	categories, _ := root.Find("meta[property='article:section']").Attr("content")
	if categories != "" {
		split := strings.Split(categories, ",")
		if len(split) > 0 {
			rs.Category.Value = strings.TrimSpace(split[0])
		}
	}

	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, " - ")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	description := root.Find("p[itemprop='description']").Text()
	rs.Description.Value = strings.TrimSpace(description)

	nodes = root.Find("li[itemprop='ingredients']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		if s == "" {
			return
		}
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	rs.Yield.Value = findYield(root.Find("h3[itemprop='recipeYield']").Text())

	return rs, nil
}
