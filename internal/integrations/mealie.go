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

// MealieRecipe represents the structure of a Mealie JSON file.
type MealieRecipe struct {
	ID                  string           `json:"id"`
	UserID              string           `json:"user_id"`
	HouseholdID         string           `json:"household_id"`
	GroupID             string           `json:"group_id"`
	Name                string           `json:"name"`
	Slug                string           `json:"slug"`
	Image               string           `json:"image"`
	RecipeServings      float64          `json:"recipe_servings"`
	RecipeYieldQuantity float64          `json:"recipe_yield_quantity"`
	RecipeYield         string           `json:"recipe_yield"`
	TotalTime           string           `json:"total_time"`
	PrepTime            string           `json:"prep_time"`
	CookTime            *string          `json:"cook_time"`
	PerformTime         string           `json:"perform_time"`
	Description         string           `json:"description"`
	RecipeCategory      []recipeCategory `json:"recipe_category"`
	Tags                []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"tags"`
	Tools              []models.HowToItem  `json:"tools"`
	Rating             float64             `json:"rating"`
	OrgURL             string              `json:"org_url"`
	DateAdded          string              `json:"date_added"`
	DateUpdated        string              `json:"date_updated"`
	CreatedAt          string              `json:"created_at"`
	UpdateAt           string              `json:"update_at"`
	LastMade           *string             `json:"last_made"`
	RecipeIngredient   []recipeIngredient  `json:"recipe_ingredient"`
	RecipeInstructions []recipeInstruction `json:"recipe_instructions"`
	Nutrition          nutrition           `json:"nutrition"`
}

func (m *MealieRecipe) Schema() models.RecipeSchema {
	category := "uncategorized"
	if len(m.RecipeCategory) > 0 {
		category = m.RecipeCategory[0].Name
	}

	var yield int16
	if m.RecipeYieldQuantity > 0 {
		v, _ := strconv.ParseInt(m.RecipeYield, 10, 16)
		yield = int16(v)
	} else if m.RecipeYieldQuantity > 0 {
		yield = int16(m.RecipeYieldQuantity)
	} else if m.RecipeServings > 0 {
		yield = int16(m.RecipeServings)
	}

	extractTimeMinutes := func(s string) int64 {
		s = strings.TrimPrefix(s, "PT")

		var t int64

		before, after, found := strings.Cut(s, "H")
		if found {
			v, err := strconv.ParseInt(before, 10, 16)
			if err == nil {
				t += v
			}

			before, after, _ = strings.Cut(after, "M")
		} else {
			before, after, _ = strings.Cut(s, "M")
		}

		v, err := strconv.ParseInt(before, 10, 16)
		if err == nil {
			t += v
		}

		return t
	}

	var cookTime string
	if m.CookTime != nil {
		cookTime = *m.CookTime
	} else if m.PerformTime != "" && m.PrepTime != "" {
		minutesCook := extractTimeMinutes(m.PerformTime) - extractTimeMinutes(m.PrepTime)

		if minutesCook > 60 {
			hours := minutesCook / 60
			minutes := minutesCook % 60
			cookTime = "PT" + strconv.FormatInt(hours, 10) + "H" + strconv.FormatInt(minutes, 10) + "M"
		} else {
			cookTime = "PT" + strconv.FormatInt(minutesCook, 10) + "M"
		}
	}

	dateUpdated := m.DateUpdated
	if m.UpdateAt != "" {
		dateUpdated = m.UpdateAt
	}

	keywords := make([]string, 0, len(m.Tags))
	for _, tag := range m.Tags {
		keywords = append(keywords, tag.Name)
	}

	ingredients := make([]string, 0, len(m.RecipeIngredient))
	for _, ing := range m.RecipeIngredient {
		ingredients = append(ingredients, ing.Display)
	}

	instructions := make([]models.HowToItem, 0, len(m.RecipeInstructions))
	for _, ins := range m.RecipeInstructions {
		instructions = append(instructions, models.NewHowToStep(ins.Text))
	}

	tools := make([]models.HowToItem, 0, len(m.Tools))
	for _, tool := range m.Tools {
		tools = append(tools, models.NewHowToTool(tool.Text))
	}

	return models.RecipeSchema{
		AtContext:       "https://schema.org",
		AtType:          &models.SchemaType{Value: "Recipe"},
		Category:        &models.Category{Value: category},
		CookTime:        cookTime,
		DateCreated:     m.DateAdded,
		DateModified:    dateUpdated,
		DatePublished:   m.DateAdded,
		Description:     &models.Description{Value: m.Description},
		Keywords:        &models.Keywords{Values: strings.Join(keywords, ",")},
		Ingredients:     &models.Ingredients{Values: ingredients},
		Instructions:    &models.Instructions{Values: instructions},
		Name:            m.Name,
		NutritionSchema: m.Nutrition.Schema(),
		PrepTime:        m.PrepTime,
		Tools:           &models.Tools{Values: tools},
		TotalTime:       m.TotalTime,
		Yield:           &models.Yield{Value: yield},
		URL:             m.OrgURL,
	}
}

func (m *MealieRecipe) UnmarshalJSON(data []byte) error {
	var temp struct {
		ID string `json:"id"`

		UserID    string `json:"user_id"`
		UserIDOld string `json:"userId"`

		HouseholdID string `json:"household_id"`

		GroupID    string `json:"group_id"`
		GroupIDOld string `json:"groupId"`

		Name  string `json:"name"`
		Slug  string `json:"slug"`
		Image string `json:"image"`

		RecipeServings      float64 `json:"recipe_servings"`
		RecipeYieldQuantity float64 `json:"recipe_yield_quantity"`

		RecipeYield    string `json:"recipe_yield"`
		RecipeYieldOld string `json:"recipeYield"`

		TotalTime    string `json:"total_time"`
		TotalTimeOld string `json:"totalTime"`

		PrepTime    string `json:"prep_time"`
		PrepTimeOld string `json:"prepTime"`

		CookTime    *string `json:"cook_time"`
		CookTimeOld *string `json:"cookTime"`

		PerformTime    string `json:"perform_time"`
		PerformTimeOld string `json:"performTime"`

		Description string `json:"description"`

		RecipeCategory    []recipeCategory `json:"recipe_category"`
		RecipeCategoryOld []recipeCategory `json:"recipeCategory"`

		Tags []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"tags"`
		Tools  []models.HowToItem `json:"tools"`
		Rating float64            `json:"rating"`

		OrgURL    string `json:"org_url"`
		OrgURLOld string `json:"orgURL"`

		DateAdded    string `json:"date_added"`
		DateAddedOld string `json:"dateAdded"`

		DateUpdated    string `json:"date_updated"`
		DateUpdatedOld string `json:"dateUpdated"`

		CreatedAt    string `json:"created_at"`
		CreatedAtOld string `json:"createdAt"`

		UpdateAt    string `json:"update_at"`
		UpdateAtOld string `json:"updateAt"`

		LastMade *string `json:"last_made"`

		RecipeIngredient    []recipeIngredient `json:"recipe_ingredient"`
		RecipeIngredientOld []recipeIngredient `json:"recipeIngredient"`

		RecipeInstructions    []recipeInstruction `json:"recipe_instructions"`
		RecipeInstructionsOld []recipeInstruction `json:"recipeInstructions"`

		Nutrition nutrition `json:"nutrition"`
	}

	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	m.ID = temp.ID

	if temp.UserID != "" {
		m.UserID = temp.UserID
	} else if temp.UserIDOld != "" {
		m.UserID = temp.UserIDOld
	}

	m.HouseholdID = temp.HouseholdID

	if temp.GroupID != "" {
		m.GroupID = temp.GroupID
	} else if temp.GroupIDOld != "" {
		m.GroupID = temp.GroupIDOld
	}

	m.Name = temp.Name
	m.Slug = temp.Slug
	m.Image = temp.Image

	m.RecipeServings = temp.RecipeServings
	m.RecipeYieldQuantity = temp.RecipeYieldQuantity

	if temp.RecipeYield != "" {
		m.RecipeYield = temp.RecipeYield
	} else if temp.RecipeYieldOld != "" {
		m.RecipeYield = temp.RecipeYieldOld
	} else if int(temp.RecipeServings) > 0 {
		m.RecipeYield = strconv.FormatFloat(temp.RecipeServings, 'f', -1, 64)
	} else if int(temp.RecipeYieldQuantity) > 0 {
		m.RecipeYield = strconv.FormatFloat(temp.RecipeYieldQuantity, 'f', -1, 64)
	}

	normalizeTime := func(s string) string {
		if s == "" || strings.HasPrefix(s, "PT") {
			return s
		}

		t := "PT"
		parts := strings.Split(s, " ")
		if len(parts) == 2 {
			if strings.Contains(s, "min") {
				t += parts[0] + "M"
			} else {
				t += parts[0] + "H"
			}
		} else if len(parts) == 4 {
			t += parts[0] + "H" + parts[1] + "M"
		}

		return t
	}

	if temp.TotalTime != "" {
		m.TotalTime = normalizeTime(temp.TotalTime)
	} else if temp.TotalTimeOld != "" {
		m.TotalTime = normalizeTime(temp.TotalTimeOld)
	}

	if temp.PrepTime != "" {
		m.PrepTime = normalizeTime(temp.PrepTime)
	} else if temp.PrepTimeOld != "" {
		m.PrepTime = normalizeTime(temp.PrepTimeOld)
	}

	if temp.CookTime != nil {
		s := normalizeTime(*temp.CookTime)
		m.CookTime = &s
	} else if temp.CookTimeOld != nil {
		s := normalizeTime(*temp.CookTimeOld)
		m.CookTime = &s
	}

	if temp.PerformTime != "" {
		m.PerformTime = normalizeTime(temp.PerformTime)
	} else if temp.PerformTimeOld != "" {
		m.PerformTime = normalizeTime(temp.PerformTimeOld)
	}

	m.Description = temp.Description

	if len(temp.RecipeCategory) > 0 {
		m.RecipeCategory = temp.RecipeCategory
	} else if len(temp.RecipeCategoryOld) > 0 {
		m.RecipeCategory = temp.RecipeCategoryOld
	}

	m.Tags = temp.Tags
	m.Tools = temp.Tools
	m.Rating = temp.Rating

	if temp.OrgURL != "" {
		m.OrgURL = temp.OrgURL
	} else if temp.OrgURLOld != "" {
		m.OrgURL = temp.OrgURLOld
	}

	if temp.DateAdded != "" {
		m.DateAdded = temp.DateAdded
	} else if temp.DateAddedOld != "" {
		m.DateAdded = temp.DateAddedOld
	}

	if temp.DateUpdated != "" {
		m.DateUpdated = temp.DateUpdated
	} else if temp.DateUpdatedOld != "" {
		m.DateUpdated = temp.DateUpdatedOld
	}

	if temp.CreatedAt != "" {
		m.CreatedAt = temp.CreatedAt
	} else if temp.CreatedAtOld != "" {
		m.CreatedAt = temp.CreatedAtOld
	}

	if temp.UpdateAt != "" {
		m.UpdateAt = temp.UpdateAt
	} else if temp.UpdateAtOld != "" {
		m.UpdateAt = temp.UpdateAtOld
	}

	m.LastMade = temp.LastMade

	if temp.RecipeIngredient != nil {
		m.RecipeIngredient = temp.RecipeIngredient
	} else if temp.RecipeIngredientOld != nil {
		m.RecipeIngredient = temp.RecipeIngredientOld
	}

	if temp.RecipeInstructions != nil {
		m.RecipeInstructions = temp.RecipeInstructions
	} else if temp.RecipeInstructionsOld != nil {
		m.RecipeInstructions = temp.RecipeInstructionsOld
	}

	m.Nutrition = temp.Nutrition

	return nil
}

type recipeCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type recipeIngredient struct {
	Quantity      float64 `json:"quantity"`
	Unit          any     `json:"unit"`
	Food          food    `json:"food"`
	Note          string  `json:"note"`
	IsFood        bool    `json:"isFood"`
	DisableAmount bool    `json:"disable_amount"`
	Display       string  `json:"display"`
	Title         any     `json:"title"`
	OriginalText  string  `json:"original_text"`
	ReferenceID   string  `json:"reference_id"`
}

func (r *recipeIngredient) UnmarshalJSON(data []byte) error {
	var temp struct {
		Quantity float64 `json:"quantity"`
		Unit     any     `json:"unit"`
		Food     food    `json:"food"`
		Note     string  `json:"note"`

		IsFood    bool `json:"isFood"`
		IsFoodOld bool `json:"is_food"`

		DisableAmount    bool `json:"disable_amount"`
		DisableAmountOld bool `json:"disableAmount"`

		Display string `json:"display"`
		Title   any    `json:"title"`

		OriginalText    string `json:"original_text"`
		OriginalTextOld string `json:"originalText"`

		ReferenceID    string `json:"reference_id"`
		ReferenceIDOld string `json:"referenceId"`
	}

	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	r.Quantity = temp.Quantity
	r.Unit = temp.Unit
	r.Food = temp.Food
	r.Note = temp.Note
	r.IsFood = temp.IsFood || temp.IsFoodOld
	r.DisableAmount = temp.DisableAmount || temp.DisableAmountOld
	r.Display = temp.Display
	r.Title = temp.Title

	if temp.OriginalText != "" {
		r.OriginalText = temp.OriginalText
	} else if temp.OriginalTextOld != "" {
		r.OriginalText = temp.OriginalTextOld
	}

	if temp.ReferenceID != "" {
		r.ReferenceID = temp.ReferenceID
	} else if temp.ReferenceIDOld != "" {
		r.ReferenceID = temp.ReferenceIDOld
	}
	return nil
}

type food struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PluralName  any    `json:"plural_name"`
	Description string `json:"description"`
	Extras      struct {
	} `json:"extras"`
	LabelID   any    `json:"label_id"`
	Aliases   []any  `json:"aliases"`
	Label     any    `json:"label"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
}

func (f *food) UnmarshalJSON(data []byte) error {
	var temp struct {
		ID   string `json:"id"`
		Name string `json:"name"`

		PluralName    any `json:"plural_name"`
		PluralNameOld any `json:"pluralName"`

		Description string `json:"description"`
		Extras      struct {
		} `json:"extras"`

		LabelID    any `json:"label_id"`
		LabelIDOld any `json:"labelId"`

		Aliases []any `json:"aliases"`
		Label   any   `json:"label"`

		CreatedAt    string `json:"created_at"`
		CreatedAtOld string `json:"createdAt"`

		UpdateAt    string `json:"update_at"`
		UpdateAtOld string `json:"updateAt"`
	}

	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	f.ID = temp.ID
	f.Name = temp.Name

	if temp.PluralName != nil {
		f.PluralName = temp.PluralName
	} else if temp.PluralNameOld != nil {
		f.PluralName = temp.PluralNameOld
	}

	f.Description = temp.Description
	f.Extras = temp.Extras

	if temp.LabelID != nil {
		f.LabelID = temp.LabelID
	} else if temp.LabelIDOld != nil {
		f.LabelID = temp.LabelIDOld
	}

	f.Aliases = temp.Aliases
	f.Label = temp.Label

	if temp.CreatedAt != "" {
		f.CreatedAt = temp.CreatedAt
	} else if temp.CreatedAtOld != "" {
		f.CreatedAt = temp.CreatedAtOld
	}

	if temp.UpdateAt != "" {
		f.UpdateAt = temp.UpdateAt
	} else if temp.UpdateAtOld != "" {
		f.UpdateAt = temp.UpdateAtOld
	}

	return nil
}

type recipeInstruction struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Summary              string `json:"summary"`
	Text                 string `json:"text"`
	IngredientReferences []any  `json:"ingredient_references"`
}

type nutrition struct {
	Calories              string
	CarbohydrateContent   string
	CholesterolContent    string
	FatContent            string
	FiberContent          string
	ProteinContent        string
	SaturatedFatContent   string
	SodiumContent         string
	SugarContent          string
	TransFatContent       string
	UnsaturatedFatContent string
}

// Schema converts the Mealie nutrition struct to a NutritionSchema one.
func (n *nutrition) Schema() *models.NutritionSchema {
	return &models.NutritionSchema{
		Calories:       n.Calories,
		Carbohydrates:  n.CarbohydrateContent,
		Cholesterol:    n.CholesterolContent,
		Fat:            n.FatContent,
		Fiber:          n.FiberContent,
		Protein:        n.ProteinContent,
		SaturatedFat:   n.SaturatedFatContent,
		Sodium:         n.SodiumContent,
		Sugar:          n.SugarContent,
		TransFat:       n.TransFatContent,
		UnsaturatedFat: n.UnsaturatedFatContent,
	}
}

func (n *nutrition) UnmarshalJSON(data []byte) error {
	var temp struct {
		Calories *string `json:"calories"`

		CarbohydrateContent    *string `json:"carbohydrate_content"`
		CarbohydrateContentOld *string `json:"carbohydrateContent"`

		CholesterolContent    *string `json:"cholesterol_content"`
		CholesterolContentOld *string `json:"cholesterolContent"`

		FatContent    *string `json:"fat_content"`
		FatContentOld *string `json:"fatContent"`

		FiberContent    *string `json:"fiber_content"`
		FiberContentOld *string `json:"fiberContent"`

		ProteinContent    *string `json:"protein_content"`
		ProteinContentOld *string `json:"proteinContent"`

		SaturatedFatContent    *string `json:"saturated_fat_content"`
		SaturatedFatContentOld *string `json:"saturatedFatContent"`

		SodiumContent    *string `json:"sodium_content"`
		SodiumContentOld *string `json:"sodiumContent"`

		SugarContent    *string `json:"sugar_content"`
		SugarContentOld *string `json:"sugarContent"`

		TransFatContent    *string `json:"trans_fat_content"`
		TransFatContentOld *string `json:"transFatContent"`

		UnsaturatedFatContent    *string `json:"unsaturated_fat_content"`
		UnsaturatedFatContentOld *string `json:"unsaturatedFatContent"`
	}

	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	if temp.Calories != nil {
		n.Calories = *temp.Calories
	}

	if temp.CarbohydrateContent != nil {
		n.CarbohydrateContent = *temp.CarbohydrateContent
	} else if temp.CarbohydrateContentOld != nil {
		n.CarbohydrateContent = *temp.CarbohydrateContentOld
	}

	if temp.CholesterolContent != nil {
		n.CholesterolContent = *temp.CholesterolContent
	} else if temp.CholesterolContentOld != nil {
		n.CholesterolContent = *temp.CholesterolContentOld
	}

	if temp.FatContent != nil {
		n.FatContent = *temp.FatContent
	} else if temp.FatContentOld != nil {
		n.FatContent = *temp.FatContentOld
	}

	if temp.FiberContent != nil {
		n.FiberContent = *temp.FiberContent
	} else if temp.FiberContentOld != nil {
		n.FiberContent = *temp.FiberContentOld
	}

	if temp.ProteinContent != nil {
		n.ProteinContent = *temp.ProteinContent
	} else if temp.ProteinContentOld != nil {
		n.ProteinContent = *temp.ProteinContentOld
	}

	if temp.SaturatedFatContent != nil {
		n.SaturatedFatContent = *temp.SaturatedFatContent
	} else if temp.SaturatedFatContentOld != nil {
		n.SaturatedFatContent = *temp.SaturatedFatContentOld
	}

	if temp.SodiumContent != nil {
		n.SodiumContent = *temp.SodiumContent
	} else if temp.SodiumContentOld != nil {
		n.SodiumContent = *temp.SodiumContentOld
	}

	if temp.SugarContent != nil {
		n.SugarContent = *temp.SugarContent
	} else if temp.SugarContentOld != nil {
		n.SugarContent = *temp.SugarContentOld
	}

	if temp.TransFatContent != nil {
		n.TransFatContent = *temp.TransFatContent
	} else if temp.TransFatContentOld != nil {
		n.TransFatContent = *temp.TransFatContentOld
	}

	if temp.UnsaturatedFatContent != nil {
		n.UnsaturatedFatContent = *temp.UnsaturatedFatContent
	} else if temp.UnsaturatedFatContentOld != nil {
		n.UnsaturatedFatContent = *temp.UnsaturatedFatContentOld
	}

	return nil
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

			var m MealieRecipe
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

			times := models.Times{
				Prep: duration.From(m.PrepTime),
				Cook: duration.From(m.PerformTime),
			}

			if m.CookTime == nil && m.PrepTime != "" {
				if m.TotalTime != "" {
					times.Cook = duration.From(m.TotalTime) - duration.From(m.PrepTime)
				} else if m.PerformTime != "" {
					times.Cook = duration.From(m.PerformTime) - duration.From(m.PrepTime)
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
					Calories:           m.Nutrition.Calories,
					Fiber:              m.Nutrition.FiberContent,
					Protein:            m.Nutrition.ProteinContent,
					Sodium:             m.Nutrition.SodiumContent,
					Sugars:             m.Nutrition.SugarContent,
					TotalCarbohydrates: m.Nutrition.CarbohydrateContent,
					TotalFat:           m.Nutrition.FatContent,
				},
				Times:     times,
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
