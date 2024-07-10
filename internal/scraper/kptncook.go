package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeKptncook(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Image.Value = root.Find("image[itemprop=image]").AttrOr("src", "")
	rs.Name = root.Find("title").Text()
	rs.Yield.Value = findYield(root.Find(".kptn-person-count").Text())

	clock := root.Find("img[alt='clockIcon']")
	if clock != nil && len(clock.Nodes) > 0 && clock.Nodes[0].NextSibling != nil {
		data := clock.Nodes[0].NextSibling.Data
		isMin := strings.Contains(data, "min")
		split := strings.Split(data, " ")
		for i, s := range split {
			_, err := strconv.ParseInt(s, 10, 64)
			if err == nil && isMin {
				rs.PrepTime = "PT" + split[i] + "M"
				break
			}
		}
	}

	// It is possible to extract the category, but I am not sure how.
	node := root.Find("img[src='/assets/images/img_wave_grey4.png']").Parent().Parent()
	for {
		node = node.Next()
		if node.Nodes == nil {
			break
		}
		s := strings.Join(strings.Fields(node.Text()), " ")
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	}

	nodes := root.Find(".kptn-step-title")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		_, err := strconv.ParseInt(s[:1], 10, 64)
		if err == nil {
			_, after, _ := strings.Cut(s, ".")
			after = strings.Join(strings.Fields(after), " ")
			s = strings.TrimSpace(after)
		}
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	return rs, nil
}
