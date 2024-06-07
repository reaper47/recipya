package scraper

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type monsieurCuisineAPI struct {
	Data struct {
		Recipe struct {
			Name          string `json:"name"`
			PublishedDate string `json:"publishedDate"`
			CreatedDate   string `json:"createdDate"`
			LastUpdated   string `json:"lastUpdated"`
			Thumbnail     struct {
				Landscape string `json:"landscape"`
			} `json:"thumbnail"`
			Complexity          string `json:"complexity"`
			PreparationDuration int    `json:"preparationDuration"`
			Duration            int    `json:"duration"`
			Categories          []struct {
				Name string `json:"name"`
			} `json:"categories"`
			Tags         []any `json:"tags"`
			Description  any   `json:"description"`
			ServingSizes []struct {
				ID                  int    `json:"id"`
				Amount              int    `json:"amount"`
				ServingUnit         string `json:"servingUnit"`
				Instruction         string `json:"instruction"`
				PreparationDuration int    `json:"preparationDuration"`
				Duration            int    `json:"duration"`
				Steps               []struct {
					ID          int    `json:"id"`
					Order       int    `json:"order"`
					Name        string `json:"name"`
					Description string `json:"description"`
					Ingredients []any  `json:"ingredients"`
				} `json:"steps"`
				Ingredients []struct {
					Order              int    `json:"order"`
					Amount             string `json:"amount"`
					Unit               string `json:"unit"`
					SystemIngredientID int    `json:"systemIngredientId"`
					IngredientGroupID  int    `json:"ingredientGroupId"`
					Name               string `json:"name"`
					IconURL            string `json:"iconUrl"`
					IngredientCategory struct {
						ID                 int         `json:"id"`
						Name               string      `json:"name"`
						IconURL            string      `json:"iconUrl"`
						IconPressURL       string      `json:"iconPressUrl"`
						SecondIconURL      interface{} `json:"secondIconUrl"`
						SecondIconPressURL interface{} `json:"secondIconPressUrl"`
					} `json:"ingredientCategory"`
				} `json:"ingredients"`
			} `json:"servingSizes"`
			Nutrients []struct {
				Name   string `json:"name"`
				Unit   string `json:"unit"`
				Amount int    `json:"amount"`
			} `json:"nutrients"`
		} `json:"recipe"`
	} `json:"data"`
}

func (s *Scraper) scrapeMonsieurCuisine(root *goquery.Document, rawURL string, files services.FilesService) (models.RecipeSchema, error) {
	js := strings.TrimSpace(root.Find("script:contains('window.siteConfig = JSON.parse(')").Last().Text())
	if js == "" {
		return models.RecipeSchema{}, errors.New("could not find recipe ID")
	}

	_, after, ok := strings.Cut(js, `"recipeId":`)
	if !ok {
		return models.RecipeSchema{}, errors.New("could not fetch recipe ID")
	}

	id, _, ok := strings.Cut(after, ",")
	if !ok {
		return models.RecipeSchema{}, errors.New("could not fetch recipe ID")
	}

	req, err := http.NewRequest(http.MethodGet, "https://mc-api.tecpal.com/api/v2/recipes/"+id, nil)
	if err != nil {
		return models.RecipeSchema{}, err
	}
	req.Header.Set("Accept-Language", "nl-NL")
	req.Header.Set("device-type", "web")

	res, err := s.Client.Do(req)
	if err != nil {
		return models.RecipeSchema{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	var m monsieurCuisineAPI
	err = json.Unmarshal(data, &m)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs := models.NewRecipeSchema()

	rs.DateCreated = m.Data.Recipe.CreatedDate
	rs.DateModified = m.Data.Recipe.LastUpdated
	rs.DatePublished = m.Data.Recipe.PublishedDate
	rs.Name = m.Data.Recipe.Name
	rs.URL = rawURL

	if m.Data.Recipe.PreparationDuration > 0 {
		rs.PrepTime = "PT" + strconv.Itoa(m.Data.Recipe.PreparationDuration) + "M"
	}

	if rs.PrepTime != "" && m.Data.Recipe.Duration > m.Data.Recipe.PreparationDuration {
		rem := m.Data.Recipe.Duration - m.Data.Recipe.PreparationDuration
		rs.CookTime = "PT" + strconv.Itoa(rem) + "M"
	}

	if len(m.Data.Recipe.Categories) > 0 {
		rs.Category.Value = m.Data.Recipe.Categories[0].Name
	}

	if m.Data.Recipe.Description != nil {
		switch x := m.Data.Recipe.Description.(type) {
		case string:
			rs.Description.Value = x
		}
	}

	keywords := make([]string, 0, len(m.Data.Recipe.Tags))
	for _, tag := range m.Data.Recipe.Tags {
		switch x := tag.(type) {
		case string:
			keywords = append(keywords, x)
		}
	}
	rs.Keywords.Values = strings.Join(keywords, ",")

	block := m.Data.Recipe.ServingSizes[0]
	rs.Ingredients.Values = make([]string, 0, len(block.Ingredients))
	for _, ing := range block.Ingredients {
		var sb strings.Builder
		if ing.Amount != "" {
			sb.WriteString(ing.Amount + " ")
			if ing.Unit != "" {
				sb.WriteString(ing.Unit + " ")
			}
		}

		sb.WriteString(ing.Name)
		rs.Ingredients.Values = append(rs.Ingredients.Values, sb.String())
	}

	rs.Yield.Value = int16(block.Amount)

	var ns models.NutritionSchema
	for _, n := range m.Data.Recipe.Nutrients {
		v := strconv.Itoa(n.Amount) + " " + n.Unit

		switch n.Name {
		case "Calories":
			ns.Calories = v
		case "Carbohydrate":
			ns.Carbohydrates = v
		case "Fat":
			ns.Fat = v
		case "Protein":
			ns.Protein = v
		}
	}

	parts := strings.Split(block.Instruction, "\r\n\r\n")
	rs.Instructions.Values = make([]models.HowToItem, 0, len(parts))
	for _, part := range parts {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(part))
	}

	imageUUID, err := files.ScrapeAndStoreImage(m.Data.Recipe.Thumbnail.Landscape)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs.Image.Value = imageUUID.String()
	return rs, nil
}
