package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeLekkerenSimpel(root *goquery.Document) (models.RecipeSchema, error) {
	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")
	keywords, _ := root.Find("meta[name='shareaholic='keywords']").Attr("content")
	name := root.Find(".hero__title").Text()

	return models.RecipeSchema{
		DatePublished: datePublished,
		DateModified:  dateModified,
		Name:          name,
		Image:         models.Image{Value: image},
		Keywords:      models.Keywords{Values: keywords},
	}, nil
}
