package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeLekkerenSimpel(root *goquery.Document) (models.RecipeSchema, error) {
	/*var (
		category string
		yield    int16
		prepTime string
		cookTime string
	)
	root.Find(".recipe__meta-title").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			category = strings.TrimSpace(s.Text())
		case 1:
			yield = findYield(strings.Split(s.Text(), " "))
		case 2:
			parts := strings.Split(strings.TrimSpace(s.Text()), " ")
			for i, part := range parts {
				if part == "bereidingstijd" {
					abbrev := "M"
					if !strings.HasPrefix(parts[i-1], "min") {

					}

					prepTime = fmt.Sprintf("PT%d")
				}
			}
		}
	})*/

	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")
	keywords, _ := root.Find("meta[name='shareaholic='keywords']").Attr("content")
	name := root.Find(".hero__title").Text()

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		DatePublished: datePublished,
		DateModified:  dateModified,
		Name:          name,
		/*Category:      models.Category{Value: category},
		Yield:         models.Yield{Value: <-chYield},
		PrepTime:      prepTime,
		CookTime:      cookTime,
		Description:   models.Description{Value: <-chDescription},*/
		Image: models.Image{Value: image},
		/*Ingredients:   models.Ingredients{Values: <-chIngredients},
		Instructions:  models.Instructions{Values: <-chInstructions},*/
		Keywords: models.Keywords{Values: keywords},
	}, nil
}
