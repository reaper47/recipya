package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"slices"
	"strings"
)

func scrapeMundodereceitasbimby(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[name='og:description']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Category.Value = root.Find("span[itemprop='recipeCategory']").First().Text()
	rs.Name = root.Find("a[itemprop='item'] span[itemprop='name']").Last().Text()
	rs.Yield.Value = findYield(root.Find("span[itemprop='recipeYield']").Text())

	_, after, ok := strings.Cut(root.Find("span.creation-date").First().Text(), ":")
	if ok {
		parts := strings.Split(strings.TrimSpace(after), ".")
		if len(parts) == 3 {
			slices.Reverse(parts)
			rs.DateCreated = strings.Join(parts, "-")
			rs.DatePublished = rs.DateCreated
		}
	}

	_, after, ok = strings.Cut(root.Find("span.changed-date").First().Text(), ":")
	if ok {
		parts := strings.Split(strings.TrimSpace(after), ".")
		if len(parts) == 3 {
			slices.Reverse(parts)
			rs.DateModified = strings.Join(parts, "-")
		}
	}

	nodes := root.Find("li[itemprop='recipeIngredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("ol[itemprop='recipeInstructions'] div[itemprop='itemListElement']")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	parts := strings.Split(root.Find(".media h5.media-heading").Text(), " ")
	switch len(parts) {
	case 2:
		rs.PrepTime = "PT" + regex.Digit.FindString(parts[0]) + "H" + regex.Digit.FindString(parts[1]) + "M"
	case 1:
		rs.PrepTime = "PT" + regex.Digit.FindString(parts[0]) + "M"
	}

	return rs, nil
}
