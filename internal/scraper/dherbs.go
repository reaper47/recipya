package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeDherbs(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getItempropContent(root, "name")
	rs.Description.Value = getNameContent(root, "description")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.PrepTime = getItempropContent(root, "prepTime")
	rs.CookTime = getItempropContent(root, "cookTime")
	rs.Yield.Value = findYield(getItempropContent(root, "recipeYield"))

	thumbs, _ := root.Find("img.attachment-post-thumbnail").Attr("srcset")
	before, _, ok := strings.Cut(thumbs, ",")
	if ok {
		rs.ThumbnailURL.Value, _, _ = strings.Cut(before, " ")
	}

	rs.Category.Value = strings.ToLower(strings.TrimSpace(root.Find("span[itemprop='recipeCategory'] a").First().Text()))

	nodes := root.Find("li[itemprop='recipeIngredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("li[itemprop='recipeInstructions']")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	return rs, nil
}
