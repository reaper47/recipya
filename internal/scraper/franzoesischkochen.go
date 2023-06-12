package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeFranzoesischKochen(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil {
		return rs, err
	}

	rs.DateModified = strings.TrimSpace(rs.DateModified)
	rs.DatePublished = strings.TrimSpace(rs.DatePublished)
	return rs, nil
}
