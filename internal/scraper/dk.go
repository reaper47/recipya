package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeDk(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	content := root.Find("section[itemtype='http://schema.org/Recipe']")

	yieldStr := content.Find("section[itemprop=recipeYield]").AttrOr("content", "")
	yield, _ := strconv.ParseInt(yieldStr, 10, 16)
	rs.Yield.Value = int16(yield)

	getIngredients(&rs, content.Find("li[itemprop=recipeIngredient]"), []models.Replace{
		{"\n", ""},
		{"useFields", ""},
	}...)

	rs.Instructions.Values = make([]models.HowToItem, 0)
	content.Find("div[itemprop=recipeInstructions] h3,div[itemprop=recipeInstructions] li").Each(func(i int, s *goquery.Selection) {
		if i > 0 && s.Nodes[0].Data == "h3" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep("\n"))
		}

		v := strings.ReplaceAll(s.Text(), "\n", "")
		v = strings.ReplaceAll(v, "\u00a0", "")
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(v))
	})

	description := content.Find("p[itemprop=description]").Text()
	rs.Description.Value = strings.TrimSpace(strings.ReplaceAll(description, "\n", ""))

	rs.Image.Value = getItempropContent(root, "url")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.Name = content.Find("h1[itemprop=name]").Text()

	return rs, nil
}
