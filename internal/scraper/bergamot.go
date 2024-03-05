package scraper

import (
	"encoding/json"
	"errors"
	"github.com/reaper47/recipya/internal/models"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type bergamot struct {
	SourceURL   string `json:"sourceUrl"`
	Lang        string `json:"lang"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserNote    any    `json:"userNote"`
	Ingredients []struct {
		Data []string `json:"data"`
	} `json:"ingredients"`
	Instructions []struct {
		Data []string `json:"data"`
	} `json:"instructions"`
	Time struct {
		PrepTime  int `json:"prepTime"`
		CookTime  any `json:"cookTime"`
		TotalTime int `json:"totalTime"`
	} `json:"time"`
	Nutrition struct {
	} `json:"nutrition"`
	Servings  int       `json:"servings"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Photos    []struct {
		ID                int       `json:"id"`
		RecipeID          int       `json:"recipeId"`
		Reference         string    `json:"reference"`
		Order             int       `json:"order"`
		Status            string    `json:"status"`
		IsUserUploaded    int       `json:"isUserUploaded"`
		SourceURL         string    `json:"sourceUrl"`
		FilenameExtension string    `json:"filenameExtension"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
		DeletedAt         any       `json:"deletedAt"`
		PhotoURL          string    `json:"photoUrl"`
		PhotoThumbURL     string    `json:"photoThumbUrl"`
	} `json:"photos"`
	SourceDomain string `json:"sourceDomain"`
}

func (s *Scraper) scrapeBergamot(url string) (models.RecipeSchema, error) {
	_, after, ok := strings.Cut(url, "https://dashboard.bergamot.app/shared/")
	if !ok {
		return models.RecipeSchema{}, errors.New("could not find bergamot recipe ID")
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.bergamot.app/recipes/shared?r="+after, nil)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	res, err := s.Client.Do(req)
	if err != nil {
		return models.RecipeSchema{}, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	var b bergamot
	err = json.Unmarshal(data, &b)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	var ingredients []string
	for _, d := range b.Ingredients {
		for _, ing := range d.Data {
			ingredients = append(ingredients, ing)
		}
	}

	var instructions []string
	for _, d := range b.Instructions {
		for _, ing := range d.Data {
			instructions = append(instructions, ing)
		}
	}

	var prep string
	if b.Time.PrepTime > 0 {
		prep = "PT" + strconv.Itoa(b.Time.PrepTime) + "M"
	}

	var cook string
	if prep != "" {
		cook = "PT" + strconv.Itoa(b.Time.TotalTime-b.Time.PrepTime) + "M"
	}

	var image string
	if len(b.Photos) > 0 {
		image = b.Photos[0].PhotoThumbURL
	}

	return models.RecipeSchema{
		AtContext:       atContext,
		AtType:          models.SchemaType{Value: "Recipe"},
		CookTime:        cook,
		DateCreated:     b.CreatedAt.String(),
		DateModified:    b.UpdatedAt.String(),
		DatePublished:   b.CreatedAt.String(),
		Description:     models.Description{Value: b.Description},
		Image:           models.Image{Value: image},
		Ingredients:     models.Ingredients{Values: ingredients},
		Instructions:    models.Instructions{Values: instructions},
		Name:            b.Title,
		NutritionSchema: models.NutritionSchema{},
		PrepTime:        prep,
		Yield:           models.Yield{Value: int16(b.Servings)},
		URL:             url,
	}, nil
}
