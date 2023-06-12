package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeFitMenCook(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil && !strings.HasPrefix(err.Error(), "@type must be Recipe") {
		return rs, err
	}
	rs.AtType = models.SchemaType{Value: "Recipe"}

	image, _ := root.Find("picture > img").Attr("data-lazy-src")
	rs.Image = models.Image{Value: image}

	nodes := root.Find("a[rel='tag']")
	keywords := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		keywords[i] = strings.ToLower(s.Text())
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
	root.Find(".fmc_ingredients ul li").Each(func(i int, s *goquery.Selection) {
		ing := s.Text()
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
	rs.Ingredients = models.Ingredients{Values: ingredients}

	instructions := make([]string, 0)
	root.Find(".fmc_step_content p").Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		if children.Nodes != nil {
			v := strings.ReplaceAll(children.Nodes[0].PrevSibling.Data, "\u00a0", "")
			instructions = append(instructions, strings.TrimSpace(v))
			return
		}

		v := strings.ReplaceAll(s.Text(), "\u00a0", "")
		instructions = append(instructions, strings.TrimSpace(v))
	})
	rs.Instructions = models.Instructions{Values: instructions}

	var nutrition models.NutritionSchema
	nutritionNodes := root.Find(".fmc_macro_cals")
	nutrition.Calories = nutritionNodes.First().Find("span").Text()

	nutritionNodes.NextAll().Each(func(i int, s *goquery.Selection) {
		switch strings.ToLower(s.Nodes[0].FirstChild.Data) {
		case "protein":
			nutrition.Protein = s.Find("span").Text()
		case "fats":
			nutrition.Fat = s.Find("span").Text()
		case "carbs":
			nutrition.Carbohydrates = s.Find("span").Text()
		case "sodium":
			nutrition.Sodium = s.Find("span").Text()
		case "fiber":
			nutrition.Fiber = s.Find("span").Text()
		case "sugar":
			nutrition.Sugar = s.Find("span").Text()
		}
	})
	rs.NutritionSchema = nutrition

	name := root.Find(".fmc_title_1").Text()
	name = strings.ReplaceAll(name, "\n", "")
	name = strings.ReplaceAll(name, "\t", "")
	rs.Name = strings.TrimSpace(name)

	rs.Keywords = models.Keywords{Values: strings.TrimSpace(strings.Join(keywords, ","))}
	return rs, nil
}
