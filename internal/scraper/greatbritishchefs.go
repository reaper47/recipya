package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeGreatBritishChefs(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	yieldStr := root.Find(".gbcicon-serves").Next().Text()
	yield, _ := strconv.ParseInt(strings.TrimSpace(yieldStr), 10, 16)
	rs.Yield = &models.Yield{Value: int16(yield)}

	node := root.Find(".gbcicon-clock").Next().Text()
	split := strings.Split(node, " ")
	if len(split) > 1 && split[1] == "minutes" {
		rs.CookTime = "PT" + split[0] + "M"
	}

	return rs, nil
}
