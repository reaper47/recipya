package integrations

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/duration"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type mealieLoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type mealieRecipesResponse struct {
	Items []struct {
		ID   string `json:"id"`
		Slug string `json:"slug"`
	} `json:"items"`
	Next *string `json:"next"`
}

type mealieRecipeResponse struct {
	ID             string  `json:"id"`
	UserID         string  `json:"userId"`
	GroupID        string  `json:"groupId"`
	Name           string  `json:"name"`
	Slug           string  `json:"slug"`
	Image          string  `json:"image"`
	RecipeYield    string  `json:"recipeYield"`
	TotalTime      string  `json:"totalTime"`
	PrepTime       string  `json:"prepTime"`
	CookTime       *string `json:"cookTime"`
	PerformTime    string  `json:"performTime"`
	Description    string  `json:"description"`
	RecipeCategory []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"recipeCategory"`
	Tags []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"tags"`
	Tools            []string `json:"tools"`
	Rating           int      `json:"rating"`
	OrgURL           string   `json:"orgURL"`
	DateAdded        string   `json:"dateAdded"`
	DateUpdated      string   `json:"dateUpdated"`
	CreatedAt        string   `json:"createdAt"`
	UpdateAt         string   `json:"updateAt"`
	RecipeIngredient []struct {
		Quantity float64     `json:"quantity"`
		Unit     interface{} `json:"unit"`
		Food     struct {
			ID          string      `json:"id"`
			Name        string      `json:"name"`
			PluralName  interface{} `json:"pluralName"`
			Description string      `json:"description"`
			Extras      struct {
			} `json:"extras"`
			LabelID   interface{}   `json:"labelId"`
			Aliases   []interface{} `json:"aliases"`
			Label     interface{}   `json:"label"`
			CreatedAt string        `json:"createdAt"`
			UpdateAt  string        `json:"updateAt"`
		} `json:"food"`
		Note          string      `json:"note"`
		IsFood        bool        `json:"isFood"`
		DisableAmount bool        `json:"disableAmount"`
		Display       string      `json:"display"`
		Title         interface{} `json:"title"`
		OriginalText  string      `json:"originalText"`
		ReferenceID   string      `json:"referenceId"`
	} `json:"recipeIngredient"`
	RecipeInstructions []struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"recipeInstructions"`
	Nutrition struct {
		Calories            string `json:"calories"`
		FatContent          string `json:"fatContent"`
		ProteinContent      string `json:"proteinContent"`
		CarbohydrateContent string `json:"carbohydrateContent"`
		FiberContent        string `json:"fiberContent"`
		SodiumContent       string `json:"sodiumContent"`
		SugarContent        string `json:"sugarContent"`
	} `json:"nutrition"`
}

// MealieImport imports recipes from a Mealie instance.
func MealieImport(baseURL, username, password string, client *http.Client, uploadImageFunc func(rc io.ReadCloser) (uuid.UUID, error), progress chan models.Progress) (models.Recipes, error) {
	usernameAttr := slog.String("username", username)

	_, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Invalid base URL", "url", baseURL, usernameAttr, "error", err)
		return nil, err
	}

	baseURL = strings.TrimSuffix(baseURL, "/")
	if !strings.HasSuffix(baseURL, "/api") {
		baseURL = baseURL + "/api"
	}

	// 1. Login
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest(http.MethodPost, baseURL+"/auth/token", strings.NewReader(data.Encode()))
	if err != nil {
		slog.Error("Couldn't create new request for auth", "url", baseURL+"/auth/token", usernameAttr, "error", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		slog.Error("Couldn't send login request", "url", req.URL.String(), usernameAttr, "error", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("Couldn't send login request", "url", req.URL.String(), usernameAttr, "integration", "mealie", "status", res.StatusCode)
		return nil, err
	}

	var login mealieLoginResponse
	err = json.NewDecoder(res.Body).Decode(&login)
	if err != nil {
		slog.Error("Couldn't decode Mealie login response", "url", req.URL.String(), usernameAttr, "error", err)
		return nil, err
	}
	authHeader := "Bearer " + login.AccessToken

	// 2. Fetch slugs
	params := url.Values{}
	params.Add("page", "1")
	params.Add("perPage", "100")
	params.Add("orderDirection", "desc")
	params.Add("requireAllCategories", "false")
	params.Add("requireAllTags", "false")
	params.Add("requireAllTools", "false")
	params.Add("requireAllFoods", "false")
	next := "/recipes?" + params.Encode()

	var slugs []string
	for {
		rawURL := baseURL + next
		rawURLAttr := slog.String("url", rawURL)

		req, err = http.NewRequest(http.MethodGet, rawURL, nil)
		if err != nil {
			slog.Error("Couldn't create new request for fetch recipes", usernameAttr, rawURLAttr, "error", err)
			return nil, err
		}
		req.Header.Add("Authorization", authHeader)
		req.Header.Add("Accept", "application/json")

		res, err = client.Do(req)
		if err != nil {
			slog.Error("Couldn't send request", "url", req.URL.String(), usernameAttr, rawURLAttr, "error", err)
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			slog.Error("Couldn't send request", "url", req.URL.String(), usernameAttr, rawURLAttr, "status", res.StatusCode)
			return nil, err
		}

		var m mealieRecipesResponse
		err = json.NewDecoder(res.Body).Decode(&m)
		if err != nil {
			_ = res.Body.Close()
			slog.Error("Couldn't decode request", "url", req.URL.String(), usernameAttr, rawURLAttr, "error", err)
			return nil, err
		}

		_ = res.Body.Close()

		for _, item := range m.Items {
			slugs = append(slugs, item.Slug)
		}

		if m.Next == nil {
			break
		}
		next = *m.Next
	}

	// 3. Fetch recipes
	var (
		mu       sync.Mutex
		numSlugs = len(slugs)
		recipes  = make(models.Recipes, 0, numSlugs)
		wg       sync.WaitGroup
	)

	wg.Add(numSlugs)

	for i, slug := range slugs {
		func() {
			defer func() {
				progress <- models.Progress{
					Value: i,
					Total: numSlugs,
				}
				wg.Done()
			}()

			rawURL := baseURL + "/recipes/" + slug
			rawURLAttr := slog.String("url", rawURL)

			req, err := http.NewRequest(http.MethodGet, rawURL, nil)
			if err != nil {
				slog.Error("Couldn't create new request", "url", rawURL, usernameAttr, rawURLAttr, "error", err)
				return
			}

			mu.Lock()
			req.Header.Add("Authorization", authHeader)
			req.Header.Add("Accept", "application/json")
			mu.Unlock()

			res, err := client.Do(req)
			if err != nil {
				slog.Error("Couldn't send request", "url", req.URL.String(), usernameAttr, rawURLAttr, "error", err)
				return
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				slog.Error("Couldn't send request", "url", req.URL.String(), usernameAttr, rawURLAttr, "status", res.StatusCode)
				return
			}

			var m mealieRecipeResponse
			err = json.NewDecoder(res.Body).Decode(&m)
			if err != nil {
				_ = res.Body.Close()
				slog.Error("Couldn't decode request", "url", req.URL.String(), usernameAttr, rawURLAttr, "error", err)
				return
			}

			category := "uncategorized"
			if len(m.RecipeCategory) > 0 {
				category = m.RecipeCategory[0].Name
			}

			source := m.OrgURL
			if source == "" {
				source = "Mealie"
			}

			var dateCreated time.Time
			before, _, ok := strings.Cut(m.CreatedAt, "T")
			if ok {
				dateCreated, _ = time.Parse(time.DateOnly, before)
			}

			var dateModified time.Time
			before, _, ok = strings.Cut(m.DateUpdated, "T")
			if ok {
				dateModified, _ = time.Parse(time.DateOnly, before)
			}

			var yield int16
			for _, s := range strings.Split(m.RecipeYield, " ") {
				parsed, err := strconv.ParseInt(s, 10, 16)
				if err == nil {
					yield = int16(parsed)
				}
			}

			ingredients := make([]string, 0, len(m.RecipeIngredient))
			for _, s := range m.RecipeIngredient {
				var v string
				if s.OriginalText != "" {
					v = s.OriginalText
				} else if s.Display != "" {
					v = s.Display
				}

				if v != "" {
					ingredients = append(ingredients, v)
				}
			}

			instructions := make([]string, 0, len(m.RecipeInstructions))
			for _, s := range m.RecipeInstructions {
				instructions = append(instructions, s.Text)
			}

			keywords := make([]string, 0, len(m.Tags))
			for _, tag := range m.Tags {
				keywords = append(keywords, tag.Name)
			}

			// 3.1. Fetch image
			var img uuid.UUID
			rawURL = baseURL + "/media/recipes/" + m.ID + "/images/min-original.webp"
			rawURLAttr = slog.String("url", rawURL)

			req, err = http.NewRequest(http.MethodGet, rawURL, nil)
			if err != nil {
				slog.Error("Couldn't create new request", "url", req.URL.String(), usernameAttr, rawURLAttr, "error", err)
				return
			}

			mu.Lock()
			req.Header.Add("Authorization", authHeader)
			req.Header.Add("Accept", "application/json")
			mu.Unlock()

			resImage, err := client.Do(req)
			if err != nil {
				slog.Error("Couldn't send request", "url", req.URL.String(), usernameAttr, rawURLAttr, "error", err)
				return
			}
			defer resImage.Body.Close()

			img, err = uploadImageFunc(resImage.Body)
			if err != nil {
				slog.Error("Couldn't upload image", "url", req.URL.String(), usernameAttr, rawURLAttr, "error", err)
			}

			var images []uuid.UUID
			if img != uuid.Nil {
				images = append(images, img)
			}

			recipe := models.Recipe{
				Category:     category,
				CreatedAt:    dateCreated,
				Description:  m.Description,
				Images:       images,
				Ingredients:  ingredients,
				Instructions: instructions,
				Keywords:     keywords,
				Name:         m.Name,
				Nutrition: models.Nutrition{
					Calories:           extractNut(m.Nutrition.Calories, " kcal"),
					Fiber:              extractNut(m.Nutrition.FiberContent, "g"),
					Protein:            extractNut(m.Nutrition.ProteinContent, "g"),
					Sodium:             extractNut(m.Nutrition.SodiumContent, "g"),
					Sugars:             extractNut(m.Nutrition.SugarContent, "g"),
					TotalCarbohydrates: extractNut(m.Nutrition.CarbohydrateContent, "g"),
					TotalFat:           extractNut(m.Nutrition.FatContent, "g"),
				},
				Times: models.Times{
					Prep: duration.From(m.PrepTime),
					Cook: duration.From(m.PerformTime),
				},
				Tools:     m.Tools,
				UpdatedAt: dateModified,
				URL:       source,
				Yield:     yield,
			}

			mu.Lock()
			recipes = append(recipes, recipe)
			mu.Unlock()
		}()
	}

	wg.Wait()
	return recipes, nil
}

func extractNut(v, unit string) string {
	if v == "" {
		return ""
	}
	return v + unit
}
