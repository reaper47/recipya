package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"slices"
	"strconv"
	"strings"
)

func scrapeGoodEatings(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	u, ok := root.Find("link[rel='canonical']").Attr("href")
	if ok {
		_, after, _ := strings.Cut(u, "https://goodeatings.com/recipes/")
		before, _, _ := strings.Cut(after, "/")
		rs.Category.Value = strings.ReplaceAll(before, "-", " ")
	}

	rs.Description.Value = getPropertyContent(root, "og:description")
	if rs.Description.Value == "" {
		rs.Description.Value = root.Find("div.post-content").First().Find("p").First().Text()
	}

	rs.DatePublished = getPropertyContent(root, "article:datePublished")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = root.Find("h2").First().Text()
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.DatePublished = getPropertyContent(root, "article:published_time")

	nodes := root.Find("p:contains('INGREDIENTS:')").NextAll()
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("p:contains('METHOD:')").NextAll()
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	nodes = root.Find(".tag-cloud-link")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(sel.Text()))
	})
	rs.Keywords.Values = strings.Join(keywords, ",")

	s := strings.TrimSuffix(root.Find("p:contains('Serves')").Text(), ".")
	parsed, err := strconv.ParseInt(regex.Digit.FindString(s), 10, 16)
	if err == nil {
		rs.Yield.Value = int16(parsed)
	} else {
		node := root.Find("p:contains('PORTIONS:')")

		yieldStr := node.Text()
		before, after, ok := strings.Cut(yieldStr, "TIME")
		if ok {
			yieldStr = before
			matches := regex.Time.FindStringSubmatch(string(after))
			if matches != nil {
				matches = slices.DeleteFunc(regex.Time.FindStringSubmatch(after), func(s string) bool { return s == "" })
			}

			switch len(matches) {
			case 2:
				rs.PrepTime = "PT" + regex.Digit.FindString(matches[1])
				if strings.Contains(matches[1], "h") {
					rs.PrepTime += "H"
				} else if strings.Contains(matches[1], "min") {
					rs.PrepTime += "M"
				}
			case 3:
				rs.PrepTime = "H" + regex.Digit.FindString(matches[2]) + "M"
			}
		}

		parsed, _ = strconv.ParseInt(regex.Digit.FindString(yieldStr), 10, 16)
		rs.Yield.Value = int16(parsed)

		if len(rs.Ingredients.Values) == 0 {
			nodes = node.NextAll()
			rs.Ingredients.Values = make([]string, 0, nodes.Length())
			nodes.Each(func(_ int, sel *goquery.Selection) {
				rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
			})
		}

		if len(rs.Instructions.Values) == 0 {
			nodes = node.ParentsUntil(".wpb_row.row-inner").Parent().Find("div").Last().Find("p")
			rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
			nodes.Each(func(_ int, sel *goquery.Selection) {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
			})
		}
	}

	return rs, nil
}
