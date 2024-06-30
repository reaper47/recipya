package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrape15gram(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.CookTime = getItempropContent(root, "cookTime")
	rs.PrepTime = getItempropContent(root, "prepTime")

	getIngredients(&rs, root.Find("li[itemprop='recipeIngredient']"))
	getInstructions(&rs, root.Find("li[itemprop='recipeInstructions']"))

	nodes := root.Find("span[itemprop='keywords']")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(sel.Text()))
	})

	rs.Description = &models.Description{Value: strings.TrimSpace(root.Find("p[itemprop='description']").Text())}
	rs.Keywords = &models.Keywords{Values: strings.Join(keywords, ",")}
	rs.Name = root.Find("h1[itemprop='name']").Text()
	rs.Yield = &models.Yield{Value: findYield(root.Find("span[itemprop='recipeYield']").Text())}
	return rs, nil
}
