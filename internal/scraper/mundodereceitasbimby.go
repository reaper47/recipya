package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"slices"
	"strings"
)

func scrapeMundodereceitasbimby(root *goquery.Document) (models.RecipeSchema, error) {
	rs.Description.Value, _ = root.Find("meta[name='og:description']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	var dateCreated string
	_, after, ok := strings.Cut(root.Find("span.creation-date").First().Text(), ":")
	if ok {
		parts := strings.Split(strings.TrimSpace(after), ".")
		if len(parts) == 3 {
			slices.Reverse(parts)
			dateCreated = strings.Join(parts, "-")
		}
	}

	var dateModified string
	_, after, ok = strings.Cut(root.Find("span.changed-date").First().Text(), ":")
	if ok {
		parts := strings.Split(strings.TrimSpace(after), ".")
		if len(parts) == 3 {
			slices.Reverse(parts)
			dateModified = strings.Join(parts, "-")
		}
	}

	nodes := root.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("ol[itemprop='recipeInstructions'] div[itemprop='itemListElement']")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, models.NewHowToStep(strings.TrimSpace(sel.Text())))
	})

	var prep string
	parts := strings.Split(root.Find(".media h5.media-heading").Text(), " ")
	switch len(parts) {
	case 2:
		prep = "PT" + regex.Digit.FindString(parts[0]) + "H" + regex.Digit.FindString(parts[1]) + "M"
	case 1:
		prep = "PT" + regex.Digit.FindString(parts[0]) + "M"
	}

	return models.RecipeSchema{
		Category:      &models.Category{Value: root.Find("span[itemprop='recipeCategory']").First().Text()},
		DateCreated:   dateCreated,
		DateModified:  dateModified,
		DatePublished: dateCreated,
		Description:   &models.Description{Value: description},
		Image:         &models.Image{Value: image},
		Ingredients:   &models.Ingredients{Values: ingredients},
		Instructions:  &models.Instructions{Values: instructions},
		Name:          root.Find("a[itemprop='item'] span[itemprop='name']").Last().Text(),
		PrepTime:      prep,
		Yield:         &models.Yield{Value: findYield(root.Find("span[itemprop='recipeYield']").Text())},
	}, nil
}
