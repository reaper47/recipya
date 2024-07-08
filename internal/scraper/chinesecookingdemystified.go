package scraper

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
	"time"
)

func scrapeChinesecookingdemystified(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = root.Find(".post-header h1").First().Text()
	rs.Description.Value = root.Find("h3.subtitle").Text()
	thumbnailSrc, exists := root.Find("iframe .ytp-cued-thumbnail-overlay-image").First().Attr("src")
	if exists {
		rs.Image.Value = regex.URL.FindString(thumbnailSrc)
	}

	parse, err := time.Parse("Jan 02, 2006", root.Find(".pencraft.pc-display-flex.pc-gap-4.pc-reset").Eq(1).Text())
	if err == nil {
		rs.DatePublished = parse.Format(time.DateOnly)
	}

	root.Find("h3:contains('Ingredients')").Next().Children().Each(func(_ int, sel *goquery.Selection) {
		sel.Find("p").Each(func(_ int, subsel *goquery.Selection) {
			rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(subsel.Text()))
		})
	})

	getInstructions(&rs, root.Find("h3:contains('Process')").NextUntil(":not(p)"))

	video := root.Find(".youtube-wrap").AttrOr("data-attrs", "")
	if video != "" {
		var m map[string]string
		err = json.Unmarshal([]byte(video), &m)
		if err == nil {
			v, ok := m["videoId"]
			if ok {
				rs.Video.Values = append(rs.Video.Values, models.VideoObject{
					AtType:   "VideoObject",
					EmbedURL: "https://youtube.com/embed/" + v,
					IsIFrame: true,
				})
			}
		}
	}

	return rs, nil
}
