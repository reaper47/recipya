package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeThatVeganDad(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = root.Find("[itemprop=headline]").Text()
	rs.ThumbnailURL.Value = getItempropContent(root, "thumbnailUrl")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.DateModified = getItempropContent(root, "dateModified")

	isInstructions := false
	root.Find("h4:contains('INGREDIENTS')").Parent().Find("ul").Each(func(_ int, ul *goquery.Selection) {
		if isInstructions || strings.Contains(ul.Prev().Text(), "METHOD") {
			isInstructions = true
			return
		}

		var name string
		prev := ul.Prev()
		if goquery.NodeName(prev) == "p" {
			name = strings.TrimSpace(prev.Text())
			rs.Ingredients.Values = append(rs.Ingredients.Values, name)
		}

		ul.Children().Each(func(_ int, sel *goquery.Selection) {
			rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
		})
	})

	var name string
	root.Find("h4:contains('METHOD')").NextAll().Each(func(_ int, sel *goquery.Selection) {
		switch goquery.NodeName(sel) {
		case "p":
			if sel.Children().Length() == 1 {
				name = strings.TrimSpace(sel.Text())
				return
			}

			s := strings.TrimPrefix(sel.Text(), "-")
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s, &models.HowToItem{
				Name: name,
			}))
		case "ul":
			sel.Children().Each(func(_ int, li *goquery.Selection) {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(li.Text()))
			})
		}
	})

	return rs, nil
}
