package scraper

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

type foodNetwork struct {
	Props struct {
		Resource struct {
			Instructions struct {
				MethodLegacy string `json:"method_legacy"`
			} `json:"instructions"`
		} `json:"resource"`
	} `json:"props"`
}

func scrapeFoodNetwork(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	j, ok := root.Find("#app").Attr("data-page")
	if ok {
		var fn foodNetwork
		err = json.Unmarshal([]byte(j), &fn)
		if err != nil {
			return models.RecipeSchema{}, err
		}

		instructions := strings.Split(fn.Props.Resource.Instructions.MethodLegacy, "</p>")
		for _, ins := range instructions {
			if ins == "" {
				continue
			}

			ins = strings.Replace(ins, "<p>", "", 1)
			ins = strings.Replace(ins, "</p>", "", 1)
			ins = strings.ReplaceAll(ins, "&deg;", "°")

			rs.Instructions.Values = append(rs.Instructions.Values, strings.TrimSpace(ins))
		}
	}

	return rs, nil
}
