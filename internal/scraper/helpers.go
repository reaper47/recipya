package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
	"strings"
)

func findYield(s string) int16 {
	parts := strings.Split(s, " ")
	for _, part := range parts {
		i, err := strconv.ParseInt(part, 10, 16)
		if err == nil {
			return int16(i)
		}
	}
	return 0
}

func getItempropContent(doc *goquery.Document, name string) string {
	s, _ := doc.Find("meta[itemprop='" + name + "']").Attr("content")
	return strings.TrimSpace(s)
}

func getNameContent(doc *goquery.Document, name string) string {
	s, _ := doc.Find("meta[name='" + name + "']").Attr("content")
	return strings.TrimSpace(s)
}

func getPropertyContent(doc *goquery.Document, name string) string {
	s, _ := doc.Find("meta[property='" + name + "']").Attr("content")
	return strings.TrimSpace(s)
}

func getIngredients(rs *models.RecipeSchema, nodes *goquery.Selection, replaceOpts ...models.Replace) {
	if rs.Ingredients == nil {
		rs.Ingredients = &models.Ingredients{}
	}

	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		for _, opt := range replaceOpts {
			s = strings.ReplaceAll(s, opt.Old, opt.New)
			if opt.Old == "useFields" {
				s = strings.Join(strings.Fields(s), " ")
			}
		}

		s = strings.TrimSpace(s)
		if s != "" {
			rs.Ingredients.Values = append(rs.Ingredients.Values, s)
		}
	})
}

func getInstructions(rs *models.RecipeSchema, nodes *goquery.Selection, replaceOpts ...models.Replace) {
	if rs.Instructions == nil {
		rs.Instructions = &models.Instructions{}
	}

	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		for _, opt := range replaceOpts {
			s = strings.ReplaceAll(s, opt.Old, opt.New)
			if opt.Old == "useFields" {
				s = strings.Join(strings.Fields(s), " ")
			}
		}

		if s != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	})
}

func getTime(rs *models.RecipeSchema, node *goquery.Selection, isPrepTime bool) {
	var s string
	t := strings.TrimSpace(node.Text())
	if t != "" {
		s = "PT" + regex.Digit.FindString(t)
		if strings.Contains(strings.ToLower(t), "min") {
			s += "M"
		} else {
			s += "H"
		}
	}

	if isPrepTime {
		rs.PrepTime = s
	} else {
		rs.CookTime = s
	}
}
