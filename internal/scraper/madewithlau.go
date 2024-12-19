package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/reaper47/recipya/internal/models"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type madeWithLau struct {
	Result struct {
		Data struct {
			JSON struct {
				RecipeKeywords   []string `json:"recipeKeywords"`
				RecipeCategory   string   `json:"recipeCategory"`
				IngredientsArray []struct {
					Type    string  `json:"_type"`
					Section string  `json:"section"`
					Key     string  `json:"_key"`
					Amount  float64 `json:"amount"`
					Item    string  `json:"item"`
					Unit    string  `json:"unit"`
				} `json:"ingredientsArray"`
				TotalTime      int    `json:"totalTime"`
				Servings       int    `json:"servings"`
				SeoDescription string `json:"seoDescription"`
				CreatedAt      string `json:"_createdAt"`
				MainImage      struct {
					Asset struct {
						URL string `json:"url"`
					} `json:"asset"`
				} `json:"mainImage"`
				PrepTime          int    `json:"prepTime"`
				TaglineSummary    string `json:"taglineSummary"`
				InstructionsArray []struct {
					RecipeCardDescription []struct {
						Children []struct {
							Text string `json:"text"`
						} `json:"children"`
					} `json:"recipeCardDescription"`
				} `json:"instructionsArray"`
				Title                   string `json:"title"`
				YoutubeTutorialVideoURL string `json:"youtubeTutorialVideoURL"`
			} `json:"json"`
		} `json:"data"`
	} `json:"result"`
}

func (s *Scraper) scrapeMadeWithLau(rawURL string) (models.RecipeSchema, error) {
	parts := strings.Split(rawURL, "/")
	apiURL := "https://www.madewithlau.com/api/trpc/recipe.bySlug,recipe.latest,course.inlinePreview,recipe.getRecipeYouTubeDetails?batch=1&input=%7B%220%22%3A%7B%22json%22%3A%7B%22slug%22%3A%22" + parts[len(parts)-1] + "%22%7D%7D%2C%221%22%3A%7B%22json%22%3A%7B%22numberOfRecipes%22%3A7%7D%7D%2C%222%22%3A%7B%22json%22%3A%7B%22slug%22%3A%22elements-of-flavor%22%7D%7D%2C%223%22%3A%7B%22json%22%3A%7B%22videoId%22%3A%22lWJpa0MRHAs%22%7D%7D%7D"

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

	var xl []madeWithLau
	err = json.Unmarshal(data, &xl)
	if err != nil && len(xl) == 0 {
		return models.RecipeSchema{}, err
	}

	content := xl[0].Result.Data.JSON
	rs := models.NewRecipeSchema()
	rs.URL = rawURL
	rs.Name = content.Title
	rs.Keywords.Values = strings.Join(content.RecipeKeywords, ",")
	rs.Category.Value = content.RecipeCategory
	rs.Yield.Value = int16(content.Servings)
	rs.Description.Value = content.SeoDescription
	rs.Image.Value = content.MainImage.Asset.URL

	before, _, ok := strings.Cut(content.CreatedAt, "+")
	if ok {
		_, err = time.Parse(time.DateTime, strings.TrimSpace(before))
		if err == nil {
			rs.DateCreated = before
		}
	} else {
		rs.DateCreated = content.CreatedAt
	}

	if content.YoutubeTutorialVideoURL != "" {
		var u string
		if strings.HasPrefix(content.YoutubeTutorialVideoURL, "https://www.youtube.com") {
			before, _, _ = strings.Cut(strings.TrimPrefix(content.YoutubeTutorialVideoURL, "https://www.youtube.com/watch?v="), "&&")
			u = "https://www.youtube.com/embed/" + before
		} else if strings.HasPrefix(content.YoutubeTutorialVideoURL, "https://youtu.be/") {
			u = "https://www.youtube.com/embed/" + strings.TrimPrefix(content.YoutubeTutorialVideoURL, "https://youtu.be/")
		}

		rs.Video.Values = append(rs.Video.Values, models.VideoObject{EmbedURL: u, IsIFrame: true})
	}

	rs.Ingredients.Values = make([]string, 0, len(content.IngredientsArray))
	for _, ing := range content.IngredientsArray {
		amount := strconv.FormatFloat(ing.Amount, 'f', -1, 64)
		if amount == "0" {
			amount = ""
		} else {
			amount += " " + ing.Unit
		}

		v := strings.TrimSpace(amount + " " + ing.Item)
		if v != "" {
			rs.Ingredients.Values = append(rs.Ingredients.Values, v)
		}
	}

	for _, block := range content.InstructionsArray {
		var step strings.Builder
		for _, ins := range block.RecipeCardDescription {
			for _, child := range ins.Children {
				step.WriteString(child.Text)
			}
			step.WriteString("\n\n")
		}
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(step.String()))
	}

	if content.PrepTime > 0 {
		rs.PrepTime = "PT" + strconv.Itoa(content.PrepTime) + "M"

		cook := content.TotalTime - content.PrepTime
		if cook > 0 {
			rs.CookTime = "PT" + strconv.Itoa(cook) + "M"
		}
	}

	return rs, nil
}
