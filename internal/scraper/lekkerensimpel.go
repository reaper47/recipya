package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeLekkerenSimpel(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Keywords.Values, _ = root.Find("meta[name='shareaholic='keywords']").Attr("content")
	rs.Name = root.Find(".hero__title").Text()

	return rs, nil
}
