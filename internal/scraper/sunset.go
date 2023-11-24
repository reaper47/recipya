package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
	"strings"
)

func scrapeSunset(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	for i, s := range rs.Instructions.Values {
		doc, err := html.Parse(strings.NewReader(s))
		if err == nil {
			node := goquery.NewDocumentFromNode(doc)
			rs.Instructions.Values[i] = node.Find("p").First().Text()
		}
	}

	return rs, nil
}
