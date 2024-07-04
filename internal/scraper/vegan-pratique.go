package scraper

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"strings"
)

func scrapeVeganPratique(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Video.Values = slices.DeleteFunc(rs.Video.Values, func(o models.VideoObject) bool {
		return o.AtType != "VideoObject"
	})

	getIngredients(&rs, root.Find(".elementor-widget-ucaddon_recipe_ingredients ul li"))
	getInstructions(&rs, root.Find(".elementor-widget-ucaddon_recipe_instructions p"))

	nodes := root.Find(".elementor-widget-video")
	rs.Video = &models.Videos{Values: make([]models.VideoObject, 0, nodes.Length())}
	nodes.Each(func(_ int, sel *goquery.Selection) {
		j := sel.AttrOr("data-settings", "")
		if j == "" {
			return
		}

		var m map[string]string
		err = json.Unmarshal([]byte(j), &m)
		if err != nil {
			return
		}

		yt, ok := m["youtube_url"]
		if !ok {
			return
		}

		rs.Video.Values = append(rs.Video.Values, models.VideoObject{
			AtType:   "VideoObject",
			EmbedURL: "https://www.youtube.com/embed/" + strings.TrimPrefix(yt, "https://www.youtube.com/watch?v="),
			IsIFrame: true,
		})
	})

	return rs, nil
}
