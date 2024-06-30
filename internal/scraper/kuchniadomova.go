package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeKuchniadomova(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	name := root.Find("h2[itemprop='name']").Text()
	name = strings.ReplaceAll(name, "\n", "")
	rs.Name = strings.ReplaceAll(name, "\t", "")

	yieldStr := root.Find("p[itemprop='recipeYield']").Text()
	yieldStr = strings.ReplaceAll(yieldStr, "-", " ")
	rs.Yield.Value = findYield(yieldStr)

	rs.Category.Value = getItempropContent(root, "recipeCategory")
	rs.Keywords.Values = getNameContent(root, "keywords")
	rs.Cuisine.Value = root.Find("p[itemprop='recipeCuisine']").Text()

	image, _ := root.Find("#article-img-1").Attr("data-src")
	rs.Image.Value = "https://kuchnia-domowa.pl" + image

	description := root.Find("#recipe-description").Text()
	description = strings.TrimPrefix(description, "\n")
	rs.Description.Value = strings.TrimSuffix(description, "\n")

	getIngredients(&rs, root.Find("li[itemprop='recipeIngredient']"))
	getInstructions(&rs, root.Find("#recipe-instructions li"))

	return rs, nil
}
