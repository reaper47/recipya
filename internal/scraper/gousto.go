package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type gousto struct {
	Status string `json:"status"`
	Data   struct {
		Entry struct {
			URL        string `json:"url"`
			Title      string `json:"title"`
			Categories []struct {
				Title string `json:"title"`
				URL   string `json:"url"`
				UID   string `json:"uid"`
			} `json:"categories"`
			Media struct {
				Images []struct {
					Image string `json:"image"`
					Width int    `json:"width"`
				} `json:"images"`
			} `json:"media"`
			Rating struct {
				Average float64 `json:"average"`
				Count   int     `json:"count"`
			} `json:"rating"`
			Description string `json:"description"`
			PrepTimes   struct {
				For2 int `json:"for_2"`
				For4 int `json:"for_4"`
			} `json:"prep_times"`
			Cuisine struct {
				Slug  string `json:"slug"`
				Title string `json:"title"`
			} `json:"cuisine"`
			Ingredients []struct {
				Label string `json:"label"`
				Title string `json:"title"`
				UID   string `json:"uid"`
				Name  string `json:"name"`
				Media struct {
					Images []struct {
						Image string `json:"image"`
						Width int    `json:"width"`
					} `json:"images"`
				} `json:"media"`
				Allergens struct {
					Allergen []struct {
						Slug string `json:"slug"`
					} `json:"allergen"`
				} `json:"allergens"`
			} `json:"ingredients"`
			Basics []struct {
				Title string `json:"title"`
				Slug  string `json:"slug"`
			} `json:"basics"`
			CookingInstructions []struct {
				Instruction string `json:"instruction"`
				Order       int    `json:"order"`
				Media       struct {
					Images []struct {
						Image string `json:"image"`
						Width int    `json:"width"`
					} `json:"images"`
				} `json:"media"`
			} `json:"cooking_instructions"`
			Allergens []struct {
				Title string `json:"title"`
				Slug  string `json:"slug"`
			} `json:"allergens"`
			Seo struct {
				Title          string        `json:"title"`
				Description    string        `json:"description"`
				Robots         []interface{} `json:"robots"`
				Canonical      string        `json:"canonical"`
				OpenGraphImage string        `json:"open_graph_image"`
			} `json:"seo"`
			Tags                   []interface{} `json:"tags"`
			UID                    string        `json:"uid"`
			Version                int           `json:"_version"`
			NutritionalInformation struct {
				PerHundredGrams struct {
					EnergyKcal     int `json:"energy_kcal"`
					EnergyKj       int `json:"energy_kj"`
					FatMg          int `json:"fat_mg"`
					FatSaturatesMg int `json:"fat_saturates_mg"`
					CarbsMg        int `json:"carbs_mg"`
					CarbsSugarsMg  int `json:"carbs_sugars_mg"`
					FibreMg        int `json:"fibre_mg"`
					ProteinMg      int `json:"protein_mg"`
					SaltMg         int `json:"salt_mg"`
					NetWeightMg    int `json:"net_weight_mg"`
				} `json:"per_hundred_grams"`
				PerPortion struct {
					EnergyKcal     int `json:"energy_kcal"`
					EnergyKj       int `json:"energy_kj"`
					FatMg          int `json:"fat_mg"`
					FatSaturatesMg int `json:"fat_saturates_mg"`
					CarbsMg        int `json:"carbs_mg"`
					CarbsSugarsMg  int `json:"carbs_sugars_mg"`
					FibreMg        int `json:"fibre_mg"`
					ProteinMg      int `json:"protein_mg"`
					SaltMg         int `json:"salt_mg"`
					NetWeightMg    int `json:"net_weight_mg"`
				} `json:"per_portion"`
			} `json:"nutritional_information"`
			PortionSizes []struct {
				Portions        int  `json:"portions"`
				IsOffered       bool `json:"is_offered"`
				IngredientsSkus []struct {
					ID         string `json:"id"`
					Code       string `json:"code"`
					Quantities struct {
						InBox int `json:"in_box"`
					} `json:"quantities"`
				} `json:"ingredients_skus"`
			} `json:"portion_sizes"`
		} `json:"entry"`
	} `json:"data"`
}

func (s *Scraper) scrapeGousto(rawURL string) (models.RecipeSchema, error) {
	parts := strings.Split(rawURL, "/")
	apiURL := "https://production-api.gousto.co.uk/cmsreadbroker/v1/recipe/" + parts[len(parts)-1]

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	res, err := s.HTTP.Do(req)
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

	var g gousto
	err = json.Unmarshal(data, &g)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	entry := g.Data.Entry
	rs := models.NewRecipeSchema()
	rs.Name = entry.Title

	if len(entry.Categories) > 0 {
		cat := strings.ToLower(entry.Categories[len(entry.Categories)-1].Title)
		before, _, ok := strings.Cut(cat, "recipe")
		if ok {
			cat = strings.TrimSpace(before)
		}
		rs.Category.Value = cat
	}

	if len(entry.Media.Images) > 0 {
		rs.Image.Value = entry.Media.Images[len(entry.Media.Images)-1].Image
	}

	rs.Description.Value = entry.Description

	if entry.PrepTimes.For2 > 0 {
		rs.PrepTime = "PT" + strconv.Itoa(entry.PrepTimes.For2) + "M"
	}

	rs.Cuisine.Value = strings.ToLower(entry.Cuisine.Title)

	rs.Ingredients.Values = make([]string, 0, len(entry.Ingredients)+len(entry.Basics))
	for _, ing := range entry.Ingredients {
		rs.Ingredients.Values = append(rs.Ingredients.Values, ing.Label)
	}
	for _, basic := range entry.Basics {
		rs.Ingredients.Values = append(rs.Ingredients.Values, basic.Title)
	}

	rs.Instructions.Values = make([]models.HowToItem, 0, len(entry.CookingInstructions))
	for _, ins := range entry.CookingInstructions {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(ins.Instruction))
		if err != nil {
			continue
		}
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(doc.Text()))
	}

	rs.Yield.Value = 2

	n := entry.NutritionalInformation.PerHundredGrams
	rs.NutritionSchema = &models.NutritionSchema{
		Calories:      strconv.Itoa(n.EnergyKcal) + " kcal",
		Carbohydrates: strconv.FormatFloat(float64(n.CarbsMg)/1000., 'f', 1, 64) + "g",
		Fat:           strconv.FormatFloat(float64(n.FatMg)/1000., 'f', 1, 64) + "g",
		Fiber:         strconv.FormatFloat(float64(n.FibreMg)/1000., 'f', 1, 64) + "g",
		Protein:       strconv.FormatFloat(float64(n.ProteinMg)/1000., 'f', 1, 64) + "g",
		SaturatedFat:  strconv.FormatFloat(float64(n.FatSaturatesMg)/1000., 'f', 1, 64) + "g",
		Sugar:         strconv.FormatFloat(float64(n.CarbsSugarsMg)/1000., 'f', 1, 64) + "g",
	}

	rs.URL = rawURL

	return rs, nil
}
