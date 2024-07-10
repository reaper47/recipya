package scraper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
)

func scrapeFitMenCook(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil && !strings.HasPrefix(err.Error(), "@type must be Recipe") {
		return rs, err
	}

	rs.Image = &models.Image{Value: root.Find("picture > img").AttrOr("data-lazy-src", "")}

	nodes := root.Find("a[rel='tag']")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		keywords = append(keywords, strings.ToLower(s.Text()))
	})

	prepTimeNode := root.Find(".fmc_prep .fmc_amount")
	if prepTimeNode != nil {
		t := prepTimeNode.Text()
		num := regex.Letters.ReplaceAllString(t, "")
		if strings.Contains(strings.ToLower(t), "min") {
			rs.PrepTime = "PT" + num + "M"
		}
	}

	cookTimeNode := root.Find(".fmc_cook")
	if cookTimeNode != nil {
		t := cookTimeNode.Text()
		num := regex.Letters.ReplaceAllString(t, "")
		if strings.Contains(strings.ToLower(t), "min") {
			rs.CookTime = "PT" + strings.TrimSpace(num) + "M"
		}
	}

	ingredients := make([]string, 0)
	root.Find(".fmc_ingredients").Last().Find("ul li").Each(func(_ int, s *goquery.Selection) {
		ing := strings.ReplaceAll(s.Text(), "\t", "")
		if !strings.Contains(ing, "\n\n") {
			return
		}

		split := strings.Split(ing, "\n\n")
		ingredients = append(ingredients, split[0])
		for _, s2 := range strings.Split(split[1], "\n") {
			s2 = strings.ReplaceAll(s2, "\u00a0", " ")
			ingredients = append(ingredients, strings.TrimSpace(s2))
		}
		ingredients = append(ingredients, "\n")
	})
	rs.Ingredients = &models.Ingredients{Values: ingredients}

	var n models.NutritionSchema
	nNodes := root.Find(".fmc_macro_cals")
	n.Calories = regex.Digit.FindString(nNodes.First().Find("span").Text())

	nNodes.NextAll().Each(func(_ int, s *goquery.Selection) {
		switch strings.ToLower(s.Nodes[0].FirstChild.Data) {
		case "protein":
			n.Protein = regex.Digit.FindString(s.Find("span").Text())
		case "fats":
			n.Fat = regex.Digit.FindString(s.Find("span").Text())
		case "carbs":
			n.Carbohydrates = regex.Digit.FindString(s.Find("span").Text())
		case "sodium":
			n.Sodium = regex.Digit.FindString(s.Find("span").Text())
		case "fiber":
			n.Fiber = regex.Digit.FindString(s.Find("span").Text())
		case "sugar":
			n.Sugar = regex.Digit.FindString(s.Find("span").Text())
		}
	})
	rs.NutritionSchema = &n

	rs.Keywords = &models.Keywords{Values: strings.TrimSpace(strings.Join(keywords, ","))}
	return rs, nil
}
