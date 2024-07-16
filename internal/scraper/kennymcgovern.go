package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeKennyMcGovern(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}
	rs.Image.Value = root.Find(".gridfeel-post-thumbnail-single-img").AttrOr("src", "")
	return rs, nil
}
