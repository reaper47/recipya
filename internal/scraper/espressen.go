package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeEspressen(root *html.Node) (rs models.RecipeSchema, err error) {
	rs, err = scrapeLdJSONs(root)
	rs.Url = "https://www.expressen.se" + rs.Url
	return rs, err
}
