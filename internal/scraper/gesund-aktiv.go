package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeGesundAktiv(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Name = strings.TrimSpace(root.Find("h1[itemprop=headline]").Text())
	getIngredients(&rs, root.Find(".news-recipes-indgredients").Last().Find("ul li"))
	getInstructions(&rs, root.Find(".news-recipes-cookingsteps").Last().Find("ol li"))
	return rs, nil
}
