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
	rs.Category.Value = strings.ToLower(strings.TrimSpace(root.Find("span[itemprop='recipeCategory'] a").First().Text()))
	getIngredients(&rs, root.Find("li[itemprop='recipeIngredient']"))
	getInstructions(&rs, root.Find("li[itemprop='recipeInstructions']"))

	thumbs, _ := root.Find("img.attachment-post-thumbnail").Attr("srcset")
	before, _, ok := strings.Cut(thumbs, ",")
	if ok {
		rs.ThumbnailURL.Value, _, _ = strings.Cut(before, " ")
	}

	return rs, nil
}
