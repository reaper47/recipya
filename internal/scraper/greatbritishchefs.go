package scraper

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

type viteSsrPlugin struct {
	PageContext struct {
		PageProps struct {
			RecipeInfoObject vitePluginRecipeInfo `json:"recipeInfoObject"`
		} `json:"pageProps"`
	} `json:"pageContext"`
}

type vitePluginRecipeInfo struct {
	YieldTextOverride string `json:"yieldTextOverride"`
	IsO8601TimeToCook string `json:"isO8601_TimeToCook"`
	Equipments        []struct {
		Title  string `json:"title"`
		QUnits int    `json:"qUnits"`
	} `json:"equipments"`
	LinkedData struct {
		Context            string   `json:"@context"`
		Name               string   `json:"name"`
		Description        string   `json:"description"`
		RecipeCategory     string   `json:"recipeCategory"`
		TotalTime          string   `json:"totalTime"`
		RecipeIngredient   []string `json:"recipeIngredient"`
		RecipeInstructions []struct {
			Type string `json:"@type"`
			Text string `json:"text"`
		} `json:"recipeInstructions"`
		Keywords string   `json:"keywords"`
		Type     string   `json:"@type"`
		Image    []string `json:"image"`
		Author   struct {
			Type string `json:"@type"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"author"`
		DatePublished string `json:"datePublished"`
		DateModified  string `json:"dateModified"`
	} `json:"linkedData"`
	Author struct {
		FullName     string `json:"fullName"`
		Introduction string `json:"introduction"`
	} `json:"author"`
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
	PageTitle       string `json:"pageTitle"`
	PageDescription string `json:"pageDescription"`
	PublishedDate   string `json:"publishedDate"`
	ModifiedDate    string `json:"modifiedDate"`
}

func scrapeGreatBritishChefs(root *goquery.Document) (models.RecipeSchema, error) {
	var vite viteSsrPlugin
	if err := json.Unmarshal([]byte(root.Find("#vite-plugin-ssr_pageContext").Text()), &vite); err != nil {
		return models.RecipeSchema{}, err
	}
	info := vite.PageContext.PageProps.RecipeInfoObject
	linkedData := info.LinkedData

	yield, _ := strconv.Atoi(strings.TrimSpace(info.YieldTextOverride))

	instructions := make([]string, len(linkedData.RecipeInstructions))
	for i, ins := range linkedData.RecipeInstructions {
		instructions[i] = ins.Text
	}

	tools := make([]string, len(info.Equipments))
	for i, equipment := range info.Equipments {
		tools[i] = strings.TrimSpace(equipment.Title)
	}

	return models.RecipeSchema{
		AtContext:     linkedData.Context,
		AtType:        models.SchemaType{Value: linkedData.Type},
		Category:      models.Category{Value: linkedData.RecipeCategory},
		CookTime:      info.IsO8601TimeToCook,
		DateModified:  linkedData.DateModified,
		DatePublished: linkedData.DatePublished,
		Description:   models.Description{Value: linkedData.Description},
		Keywords:      models.Keywords{Values: linkedData.Keywords},
		Image:         models.Image{Value: linkedData.Image[0]},
		Ingredients:   models.Ingredients{Values: linkedData.RecipeIngredient},
		Instructions:  models.Instructions{Values: instructions},
		Name:          linkedData.Name,
		Tools:         models.Tools{Values: tools},
		Yield:         models.Yield{Value: int16(yield)},
	}, nil
}
