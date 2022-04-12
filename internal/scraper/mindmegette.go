package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeMindMegette(root *html.Node) (models.RecipeSchema, error) {
	rs, err := scrapeLdJSON(root)

	chYield := make(chan models.Yield)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- models.Yield{Value: yield}
		}()

		node := getElement(root, "class", "spritePortion")
		yieldStr := node.NextSibling.FirstChild.FirstChild.Data
		yield = findYield(strings.Split(yieldStr, " "))
	}()

	rs.Yield = <-chYield
	return rs, err
}
