package integrations

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tandoorLoginResponse struct {
	ID      int       `json:"id"`
	Token   string    `json:"token"`
	Scope   string    `json:"scope"`
	Expires time.Time `json:"expires"`
	UserID  int       `json:"user_id"`
	Test    int       `json:"test"`
}

type tandoorRecipes struct {
	Count   int     `json:"count"`
	Next    *string `json:"next"`
	Results []struct {
		ID int64 `json:"id"`
	} `json:"results"`
}

type tandoorRecipe struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Keywords    []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Label string `json:"label"`
	} `json:"keywords"`
	Steps []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Instruction string `json:"instruction"`
		Ingredients []struct {
			ID   int `json:"id"`
			Food struct {
				Name     string `json:"name"`
				FullName string `json:"full_name"`
			} `json:"food"`
			Unit *struct {
				Name string `json:"name"`
			} `json:"unit"`
			Amount float64 `json:"amount"`
		} `json:"ingredients"`
		InstructionsMarkdown string `json:"instructions_markdown"`
		ShowAsHeader         bool   `json:"show_as_header"`
	} `json:"steps"`
	WorkingTime  int       `json:"working_time"`
	WaitingTime  int       `json:"waiting_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	SourceURL    *string   `json:"source_url"`
	Nutrition    any       `json:"nutrition"`
	Servings     int       `json:"servings"`
	FilePath     string    `json:"file_path"`
	ServingsText string    `json:"servings_text"`
	Rating       float64   `json:"rating"`
	LastCooked   time.Time `json:"last_cooked"`
}

// TandoorImport imports recipes from a Tandoor instance.
func TandoorImport(baseURL, username, password string, client *http.Client, uploadImageFunc func(rc io.ReadCloser) (uuid.UUID, error), progress chan models.Progress) (models.Recipes, error) {
	usernameAttr := slog.String("username", username)

	_, err := url.Parse(baseURL)
	if err != nil {
		slog.Error("Invalid base URL", "url", baseURL, usernameAttr, "error", err)
		return nil, err
	}

	baseURL = strings.TrimSuffix(baseURL, "/")

	// 1. Login
	c := credentials{Username: username, Password: password}
	xb, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api-token-auth/", bytes.NewReader(xb))
	if err != nil {
		slog.Error("Couldn't create new request for auth", "url", baseURL, usernameAttr, "error", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		slog.Error("Couldn't make request for auth", "url", req.Method+" "+req.URL.String(), usernameAttr, "error", err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		_ = res.Body.Close()
		slog.Error("Couldn't send login request", "url", req.Method+" "+req.URL.String(), usernameAttr, "integration", "tandoor", "status", res.StatusCode)
		return nil, err
	}

	var login tandoorLoginResponse
	err = json.NewDecoder(res.Body).Decode(&login)
	if err != nil {
		_ = res.Body.Close()
		slog.Error("Couldn't parse response for auth", "url", req.Method+" "+req.URL.String(), usernameAttr, "error", err)
		return nil, err
	}

	_ = res.Body.Close()
	authHeader := "Bearer " + login.Token

	// 2. Fetch recipe IDs
	var (
		next      = baseURL + "/api/recipe/"
		recipeIDs []int64
	)

	for {
		req, err := http.NewRequest(http.MethodGet, next, nil)
		if err != nil {
			slog.Error("Couldn't create new request for auth", "url", next, usernameAttr, "error", err)
			return nil, err
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", authHeader)

		res, err := client.Do(req)
		if err != nil {
			slog.Error("Couldn't make request", "url", next, "error", err)
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			_ = res.Body.Close()
			slog.Error("Couldn't send request", "url", req.Method+" "+req.URL.String(), usernameAttr, "url", next, "status", res.StatusCode)
			return nil, errors.New("could not send request to /recipe/")
		}

		var t tandoorRecipes
		err = json.NewDecoder(res.Body).Decode(&t)
		if err != nil {
			_ = res.Body.Close()
			slog.Error("Couldn't parse response tandoor recipes", "url", req.Method+" "+req.URL.String(), usernameAttr, "error", err)
			return nil, err
		}
		_ = res.Body.Close()

		if len(recipeIDs) == 0 {
			recipeIDs = make([]int64, 0, t.Count)
		}

		for _, result := range t.Results {
			recipeIDs = append(recipeIDs, result.ID)
		}

		if t.Next == nil {
			break
		}

		next = *t.Next
	}

	// 3. Fetch recipe
	var (
		numRecipes = len(recipeIDs)
		recipes    = make(models.Recipes, 0, len(recipeIDs))
	)

	for i, id := range recipeIDs {
		progress <- models.Progress{
			Value: i,
			Total: numRecipes,
		}

		rawURL := baseURL + "/api/recipe/" + strconv.FormatInt(id, 10) + "/"

		req, err := http.NewRequest(http.MethodGet, rawURL, nil)
		if err != nil {
			slog.Error("Couldn't create new request for auth", "url", rawURL, usernameAttr, "error", err)
			continue
		}
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", authHeader)

		res, err := client.Do(req)
		if err != nil {
			slog.Error("Couldn't make request", "url", rawURL, "error", err)
			continue
		}

		if res.StatusCode != http.StatusOK {
			_ = res.Body.Close()
			slog.Error("Couldn't send request", "url", req.Method+" "+req.URL.String(), usernameAttr, "status", res.StatusCode)
			continue
		}

		var t tandoorRecipe
		err = json.NewDecoder(res.Body).Decode(&t)
		if err != nil {
			_ = res.Body.Close()
			slog.Error("Couldn't parse response tandoor recipe", "url", req.Method+" "+req.URL.String(), usernameAttr, "error", err)
			continue
		}
		_ = res.Body.Close()

		yield := int16(t.Servings)
		if yield == 0 {
			yield = 1
		}

		var img uuid.UUID
		resImage, err := client.Get(t.Image)
		if err == nil {
			img, _ = uploadImageFunc(resImage.Body)
			_ = resImage.Body.Close()
		}

		src := "Tandoor"
		if t.SourceURL != nil {
			src = *t.SourceURL
		}

		keywords := make([]string, 0, len(t.Keywords))
		for _, kw := range t.Keywords {
			keywords = append(keywords, kw.Name)
		}

		category := "uncategorized"
		if len(keywords) > 0 {
			category = keywords[0]
		}

		var (
			ingredients  []string
			instructions = make([]string, 0, len(t.Steps))
		)
		for _, step := range t.Steps {
			for _, s := range strings.Split(step.Instruction, "\\n\\n") {
				instructions = append(instructions, strings.TrimSpace(s))
			}

			for _, ing := range step.Ingredients {
				s := ing.Food.Name
				if ing.Amount > 0 {
					formatted := strconv.FormatFloat(ing.Amount, 'g', -1, 64)
					s = formatted + " " + s
				}

				if ing.Unit != nil && ing.Unit.Name != "" {
					s = s + " " + ing.Unit.Name
				}

				ingredients = append(ingredients, s)
			}
		}

		recipes = append(recipes, models.Recipe{
			Category:     category,
			CreatedAt:    t.CreatedAt.UTC(),
			Description:  t.Description,
			Image:        img,
			Ingredients:  ingredients,
			Instructions: instructions,
			Keywords:     keywords,
			Name:         t.Name,
			Nutrition:    models.Nutrition{},
			Times: models.Times{
				Prep: time.Duration(t.WorkingTime) * time.Minute,
				Cook: time.Duration(t.WaitingTime) * time.Minute,
			},
			UpdatedAt: t.UpdatedAt.UTC(),
			URL:       src,
			Yield:     yield,
		})
	}
	return recipes, nil
}
