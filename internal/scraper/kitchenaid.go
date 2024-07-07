package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"slices"
	"strings"
)

func scrapeKitchenaid(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	summary := root.Find(".blogRecipe-summary")
	rs.Yield.Value = findYield(summary.Find("li:contains('Makes')").Text())

	prep := strings.TrimSpace(summary.Find("li:contains('Time')").Text())
	if prep != "" {
		rs.PrepTime = "PT"
		matches := regex.Time.FindAllString(prep, 2)
		if matches != nil {
			matches = slices.DeleteFunc(matches, func(s string) bool { return s == "" })
		}

		if len(matches) == 2 {
			rs.PrepTime += regex.Digit.FindString(matches[0]) + "H" + regex.Digit.FindString(matches[1]) + "M"
		}
	}

	article := root.Find("article.blogRecipe-article")
	article.Find(".leftPanel ul").Each(func(_ int, ul *goquery.Selection) {
		section := strings.TrimSpace(ul.Prev().Text())
		if section == "" {
			return
		}

		if section != "Ingredients" {
			rs.Ingredients.Values = append(rs.Ingredients.Values, section)
		}

		ul.Children().Each(func(_ int, li *goquery.Selection) {
			if goquery.NodeName(li) == "ul" {
				li.Children().Each(func(_ int, subli *goquery.Selection) {
					rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(subli.Text()))
				})
			} else {
				rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(li.Text()))
			}
		})
	})
	getInstructions(&rs, article.Find(".rightPanel ol li"), []models.Replace{
		{"useFields", ""},
	}...)

	return rs, nil
}
