package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
	"time"
)

func scrapePureWow(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Name = root.Find(".hed").Text()
	rs.Image.Value, _ = root.Find(".recipe-img-container img").First().Attr("src")

	root.Find(".recipe-text p").Each(func(_ int, sel *goquery.Selection) {
		rs.Description.Value += sel.Text() + "\n\n"
	})
	rs.Description.Value = strings.TrimSpace(rs.Description.Value)

	pub := strings.TrimSpace(strings.TrimPrefix(root.Find(".pub-date").Text(), "Published"))
	parsed, err := time.Parse("Jan 02, 2006", pub)
	if err == nil {
		rs.DatePublished = parsed.Format(time.DateTime)
	}

	getTime(&rs, root.Find(".prep"), true)
	getTime(&rs, root.Find(".cook"), false)
	getIngredients(&rs, root.Find("[itemprop=recipeIngredient]"))
	getInstructions(&rs, root.Find(".recipe-direction"))

	for i, ins := range rs.Instructions.Values {
		if strings.Index(ins.Text, ".") < 4 {
			_, after, _ := strings.Cut(ins.Text, ".")
			rs.Instructions.Values[i].Text = strings.TrimSpace(after)
		}
	}

	root.Find(".nutrient").Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		n := regex.Digit.FindString(s)
		if strings.HasSuffix(s, "calories") {
			rs.NutritionSchema.Calories = n
		} else if strings.HasSuffix(s, "fat") {
			rs.NutritionSchema.Fat = n
		} else if strings.HasSuffix(s, "carbs") {
			rs.NutritionSchema.Carbohydrates = n
		} else if strings.HasSuffix(s, "protein") {
			rs.NutritionSchema.Protein = n
		} else if strings.HasSuffix(s, "sugars") {
			rs.NutritionSchema.Sugar = n
		}
	})

	return rs, nil
}
