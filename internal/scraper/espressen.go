package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeEspressen(root *html.Node) (rs models.RecipeSchema, err error) {
	rs, err = scrapeLdJSONs(root)
	rs.URL = "https://www.expressen.se" + rs.URL
	return rs, err
}
