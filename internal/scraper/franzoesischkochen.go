package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeFranzoesischKochen(root *html.Node) (rs models.RecipeSchema, err error) {
	rs, err = findRecipeLdJSON(root)
	rs.DateModified = strings.TrimSpace(rs.DateModified)
	rs.DatePublished = strings.TrimSpace(rs.DatePublished)
	return rs, err
}
