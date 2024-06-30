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

	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.DatePublished = getPropertyContent(root, "article:published_time")

	name := getPropertyContent(root, "og:title")
	before, _, ok := strings.Cut(name, " - ")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	description := root.Find("p[itemprop='description']").Text()
	rs.Description.Value = strings.TrimSpace(description)
	getIngredients(&rs, root.Find("li[itemprop='ingredients']"))
	getInstructions(&rs, root.Find("div[itemprop='recipeInstructions'] p"))
	rs.Yield.Value = findYield(root.Find("h3[itemprop='recipeYield']").Text())

	return rs, nil
}
