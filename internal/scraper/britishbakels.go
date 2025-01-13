package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeBritishBakels(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = strings.TrimSpace(root.Find(".content-header ").Text())
	rs.Category.Value = root.Find("h4:contains('Category')").Next().Text()

	var (
		isHeader   bool
		unit       string
		ingredient string
		numRight   int
	)

	root.Find("div.card:contains('Ingredients')").Find("[id^=tab-ingredients] div").Each(func(_ int, sel *goquery.Selection) {
		class := sel.AttrOr("class", "")
		if !strings.Contains(class, "text-xs-") {
			return
		}

		s := strings.TrimSpace(sel.Text())
		if strings.HasPrefix(s, "Group") || s == "" {
			rs.Ingredients.Values = append(rs.Ingredients.Values, s)
			return
		}

		if strings.HasPrefix(s, "Ingredient") {
			isHeader = true
			return
		}

		isLeft := strings.Contains(class, "text-xs-left")
		if isLeft {
			numRight = 0
		} else {
			numRight++
		}

		if numRight == 2 {
			return
		}

		if isHeader {
			if unit == "" && !isLeft {
				unit = s
				return
			} else if isLeft {
				isHeader = false
			} else {
				return
			}
		}

		if isLeft {
			ingredient = s
		}

		if numRight == 1 {
			rs.Ingredients.Values = append(rs.Ingredients.Values, s+" "+unit+" "+ingredient)
		}
	})

	content := root.Find("div.card:contains('Method')")
	nodes := content.Find("ol")
	if nodes.Length() == 0 {
		text := strings.TrimSpace(content.Find(".card-body").Text())
		for _, s := range strings.Split(text, "\n") {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	} else {
		nodes.Each(func(_ int, sel *goquery.Selection) {
			name := strings.TrimSpace(sel.Prev().Text())
			var h models.HowToItem
			if name != "" {
				h.Name = name
			}

			sel.Children().Each(func(_ int, li *goquery.Selection) {
				h.Text = strings.TrimSpace(li.Text())
				rs.Instructions.Values = append(rs.Instructions.Values, h)
			})
		})
	}

	return rs, nil
}
