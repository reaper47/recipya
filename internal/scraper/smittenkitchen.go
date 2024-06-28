package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
	"strings"
)

func scrapeSmittenKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	parsed, err := strconv.ParseInt(regex.Digit.FindString(root.Find("*[itemprop='recipeYield']").Text()), 10, 16)
	if err == nil {
		rs.Yield.Value = int16(parsed)
	}

	rs.PrepTime, _ = root.Find("time[itemprop='totalTime']").Attr("datetime")
	rs.Description.Value = root.Find(".smittenkitchen-print-hide p").First().Text()

	nodes := root.Find(".jetpack-recipe-ingredients").Last().Find("ul")
	children := nodes.Children()
	rs.Ingredients.Values = make([]string, 0, children.Length())
	children.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find(".jetpack-recipe-directions").Last()
	children = nodes.Children()
	rs.Instructions.Values = make([]models.HowToItem, 0, children.Length())
	children.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	})

	rs.Image.Value, _ = root.Find(".post-thumbnail-container").First().Find("img").Attr("src")
	rs.DatePublished, _ = root.Find(".entry-date.published").First().Attr("datetime")
	rs.DateModified, _ = root.Find(".updated").First().Attr("datetime")

	title := root.Find("h3[itemprop='name']").Text()
	if title == "" {
		title = strings.TrimSpace(root.Find(".entry-title").First().Text())
	}
	rs.Name = title

	if len(rs.Ingredients.Values) == 0 && len(rs.Instructions.Values) == 0 {
		root.Find(".smittenkitchen-print-hide").First().NextAll().Each(func(_ int, sel *goquery.Selection) {
			if goquery.NodeName(sel) == "p" {
				s := strings.TrimSpace(sel.Text())
				if strings.HasPrefix(strings.ToLower(s), "serve") {
					parsed, err = strconv.ParseInt(regex.Digit.FindString(s), 10, 16)
					if err == nil {
						rs.Yield.Value = int16(parsed)
					}
					return
				}

				idx := regex.Unit.FindStringIndex(s)
				isAtStart := idx != nil && idx[0] < 8

				idx = regex.Digit.FindStringIndex(s)
				isStartWithNumber := idx != nil && idx[0] == 0

				if len(rs.Ingredients.Values) == 0 && (isAtStart || isStartWithNumber) && sel.Has("br").Length() > 0 {
					xs := strings.Split(sel.Text(), "\n")
					rs.Ingredients.Values = make([]string, 0, len(xs))
					for _, x := range xs {
						rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(x))
					}
				} else if len(rs.Ingredients.Values) > 0 {
					rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
				}
			}
		})
	}

	return rs, nil
}
