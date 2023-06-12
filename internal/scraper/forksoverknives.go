package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeForksOverKnives(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	if rs.Category.Value != "" {
		xb := regex.Anchor.ReplaceAll([]byte(rs.Category.Value), []byte(""))
		s := strings.ReplaceAll(string(xb), "</a>", "")
		xs := strings.Split(s, " ")
		rs.Category.Value = xs[len(xs)-1]
	}
	return rs, nil
}
