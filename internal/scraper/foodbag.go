package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/reaper47/recipya/internal/models"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type foodbag struct {
	DishRecipe struct {
		Name                string  `json:"name"`
		Subtitle            string  `json:"subtitle"`
		Description         string  `json:"description"`
		ImageURL1           string  `json:"imageUrl1"`
		SupplementaryAmount float64 `json:"supplementaryAmount"`
		TotalTime           string  `json:"totalTime"`
		OptionalTag         string  `json:"optionalTag"`
		DishVariants        []struct {
			NumberOfPersons int `json:"numberOfPersons"`
			PreparationTime int `json:"preparationTime"`
			CookingTime     int `json:"cookingTime"`
			TotalTime       int `json:"totalTime"`
			Ingredients     []struct {
				Name                     string      `json:"name"`
				Size                     string      `json:"size"`
				IsBasic                  bool        `json:"isBasic"`
				Organic                  bool        `json:"organic"`
				Allergies                interface{} `json:"allergies"`
				IngredientAllergensIcons interface{} `json:"ingredientAllergensIcons"`
				Quantity                 interface{} `json:"quantity"`
			} `json:"ingredients"`
			Allergies        string `json:"allergies"`
			Directions       string `json:"directions"`
			Tip              string `json:"tip"`
			AllergiesSummary string `json:"allergiesSummary"`
		} `json:"dishVariants"`
		NutriValue100G string `json:"nutriValue100g"`
		NutriValue2P   string `json:"nutriValue2p"`
	} `json:"dishRecipe"`
	PrimaryTags []struct {
		Name      string `json:"name"`
		StyleName string `json:"styleName"`
	} `json:"primaryTags"`
}

type foodbagNutrition struct {
	Items []struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	} `json:"Items"`
}

func (s *Scraper) scrapeFoodbag(rawURL string) (models.RecipeSchema, error) {
	parse, err := url.Parse(rawURL)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	queries := strings.Split(parse.RawQuery, "&")
	if len(queries) == 0 {
		return models.RecipeSchema{}, errors.New("query parameter 'dishId' not found")
	}

	var dishID string
	for _, query := range queries {
		before, after, ok := strings.Cut(query, "=")
		if ok && before == "dishId" {
			dishID = after
		}
	}

	if dishID == "" {
		return models.RecipeSchema{}, errors.New("query parameter 'dishId' not found")
	}

	apiURL := "https://admin.foodbag.be/api/dishrecipe?dishId=" + dishID + "&language=nl"
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	res, err := s.Client.Do(req)
	if err != nil {
		return models.RecipeSchema{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return models.RecipeSchema{}, fmt.Errorf("got status code %d for %q", res.StatusCode, apiURL)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	var fb foodbag
	err = json.Unmarshal(data, &fb)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs := models.NewRecipeSchema()
	rs.URL = rawURL

	keywords := make([]string, 0, len(fb.PrimaryTags))
	for _, tag := range fb.PrimaryTags {
		keywords = append(keywords, tag.Name)
	}
	rs.Keywords.Values = strings.Join(keywords, ",")

	var nutrition models.NutritionSchema
	if fb.DishRecipe.NutriValue100G != "" {
		var n foodbagNutrition
		err = json.Unmarshal([]byte(fb.DishRecipe.NutriValue100G), &n)
		if err != nil {
			return models.RecipeSchema{}, err
		}

		for _, item := range n.Items {
			switch item.Name {
			case "energie":
				nutrition.Calories = item.Value
			case "vetten":
				nutrition.Fat = item.Value
			case "koolhydraten":
				nutrition.Carbohydrates = item.Value
			case "voedingsvezels":
				nutrition.Fiber = item.Value
			case "eiwitten":
				nutrition.Protein = item.Value
			}
		}
	}
	rs.NutritionSchema = &nutrition

	if len(fb.DishRecipe.DishVariants) > 0 {
		variant := fb.DishRecipe.DishVariants[0]

		rs.Yield.Value = int16(variant.NumberOfPersons)
		rs.PrepTime = "PT" + strconv.Itoa(variant.PreparationTime) + "M"
		rs.CookTime = "PT" + strconv.Itoa(variant.CookingTime) + "M"

		parts := strings.Split(variant.Directions, "|")
		rs.Instructions.Values = make([]models.HowToStep, 0, len(parts))
		for _, v := range parts {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(v))
		}

		rs.Ingredients.Values = make([]string, 0, len(variant.Ingredients))
		for _, ing := range variant.Ingredients {
			var quantity string
			switch x := ing.Quantity.(type) {
			case string:
				quantity = x + " " + ing.Size + " "
			}
			rs.Ingredients.Values = append(rs.Ingredients.Values, quantity+ing.Name)
		}
	}

	rs.Description.Value = fb.DishRecipe.Description
	rs.Image.Value = fb.DishRecipe.ImageURL1
	rs.Name = fb.DishRecipe.Name + " " + fb.DishRecipe.Subtitle

	return rs, nil
}
