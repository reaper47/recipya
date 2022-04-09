package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/regex"
	"golang.org/x/net/html"
)

func scrapeForksOverKnives(root *html.Node) (rs models.RecipeSchema, err error) {
	rs, err = scrapeLdJSON(root)

	if rs.Category.Value != "" {
		xb := regex.Anchor.ReplaceAll([]byte(rs.Category.Value), []byte(""))
		s := strings.ReplaceAll(string(xb), "</a>", "")
		xs := strings.Split(s, " ")
		rs.Category.Value = xs[len(xs)-1]
	}
	return rs, err
}
