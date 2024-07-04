package server_test

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/blang/semver"
	"github.com/google/go-github/v59/github"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"github.com/reaper47/recipya/internal/services"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/internal/units"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

var mutex *sync.Mutex

func init() {
	mutex = &sync.Mutex{}
}

func newServerTest() *server.Server {
	srv := server.NewServer(&mockRepository{
		AuthTokens:             make([]models.AuthToken, 0),
		categories:             map[int64][]string{1: {"chicken"}},
		CookbooksRegistered:    map[int64][]models.Cookbook{1: {{ID: 1}}},
		RecipesRegistered:      make(map[int64]models.Recipes),
		Reports:                make(map[int64][]models.Report),
		ShareLinks:             make(map[string]models.Share),
		UserSettingsRegistered: make(map[int64]*models.UserSettings),
		UsersRegistered:        make([]models.User, 0),
		UsersUpdated:           make([]int64, 0),
	})
	srv.Email = &mockEmail{}
	srv.Files = &mockFiles{}
	srv.Integrations = &mockIntegrations{}
	srv.Scraper = &mockScraper{}

	_ = os.Remove("sessions.csv")
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	return srv
}

type mockRepository struct {
	AuthTokens                         []models.AuthToken
	AddRecipeCategoryFunc              func(name string, userID int64) error
	AddRecipesFunc                     func(recipes models.Recipes, userID int64, progress chan models.Progress) ([]int64, []models.ReportLog, error)
	AddShareRecipeFunc                 func(recipeID, userID int64) (int64, error)
	categories                         map[int64][]string
	CookbooksFunc                      func(userID int64) ([]models.Cookbook, error)
	CookbooksRegistered                map[int64][]models.Cookbook
	DeleteCategoryFunc                 func(name string, userID int64) error
	DeleteCookbookFunc                 func(id, userID int64) error
	IsUserPasswordFunc                 func(userID int64, password string) bool
	MeasurementSystemsFunc             func(userID int64) ([]units.System, models.UserSettings, error)
	RecipeFunc                         func(id, userID int64) (*models.Recipe, error)
	RecipesRegistered                  map[int64]models.Recipes
	Reports                            map[int64][]models.Report
	ReportsFunc                        func(userID int64) ([]models.Report, error)
	RestoreUserBackupFunc              func(backup *models.UserBackup) error
	ShareLinks                         map[string]models.Share
	SwitchMeasurementSystemFunc        func(system units.System, userID int64) error
	UpdateCookbookImageFunc            func(id int64, image uuid.UUID, userID int64) error
	UpdateConvertMeasurementSystemFunc func(userID int64, isEnabled bool) error
	UpdateCalculateNutritionFunc       func(userID int64, isEnabled bool) error
	UserSettingsRegistered             map[int64]*models.UserSettings
	UsersRegistered                    []models.User
	UsersUpdated                       []int64
}

func (m *mockRepository) AddRecipes(xr models.Recipes, userID int64, progress chan models.Progress) ([]int64, []models.ReportLog, error) {
	if xr == nil {
		return nil, nil, errors.New("recipe is nil")
	}

	if m.AddRecipesFunc != nil {
		return m.AddRecipesFunc(xr, userID, progress)
	}

	if m.RecipesRegistered[userID] == nil {
		m.RecipesRegistered[userID] = make(models.Recipes, 0)
	}

	recipeIDs := make([]int64, 0, len(xr))
	for i, recipe := range xr {
		if recipe.Category == "" {
			recipe.Category = "uncategorized"
		} else {
			delimiters := []string{",", ";", "|"}
			for _, delim := range delimiters {
				if strings.Contains(recipe.Category, delim) {
					recipe.Category = strings.Split(recipe.Category, delim)[0]
				}
			}
		}

		if recipe.Yield == 0 {
			recipe.Yield = 1
		}

		if recipe.URL == "" {
			recipe.URL = "Unknown"
		}

		if !slices.ContainsFunc(m.RecipesRegistered[userID], func(r models.Recipe) bool { return recipe.Name == r.Name }) {
			mutex.Lock()
			recipe.ID = int64(i + 1)
			m.RecipesRegistered[userID] = append(m.RecipesRegistered[userID], recipe)
			recipeIDs = append(recipeIDs, recipe.ID)
			mutex.Unlock()
		}

		if progress != nil {
			progress <- models.Progress{Value: i, Total: len(xr)}
		}
	}

	return recipeIDs, nil, nil
}

func (m *mockRepository) AddShareLink(share models.Share) (string, error) {
	if share.CookbookID != -1 {
		for _, cookbooks := range m.CookbooksRegistered {
			if slices.ContainsFunc(cookbooks, func(c models.Cookbook) bool { return c.ID == share.CookbookID }) {
				for link, s := range m.ShareLinks {
					if s.CookbookID == share.CookbookID && s.UserID == share.UserID {
						return link, nil
					}
				}

				link := "/c/33320755-82f9-47e5-bb0a-d1b55cbd3f7b"
				m.ShareLinks[link] = share
				return link, nil
			}
		}
	} else if share.RecipeID != -1 {
		for _, recipes := range m.RecipesRegistered {
			if slices.ContainsFunc(recipes, func(r models.Recipe) bool { return r.ID == share.RecipeID }) {
				for link, s := range m.ShareLinks {
					if s.RecipeID == share.RecipeID && s.UserID == share.UserID {
						return link, nil
					}
				}

				if m.ShareLinks == nil {
					m.ShareLinks = make(map[string]models.Share)
				}

				link := "/r/33320755-82f9-47e5-bb0a-d1b55cbd3f7b"
				m.ShareLinks[link] = share
				return link, nil
			}
		}
	}
	return "", errors.New("cookbook or recipe not found")
}

func (m *mockRepository) AddShareRecipe(recipeID, userID int64) (int64, error) {
	if m.AddShareRecipeFunc != nil {
		return m.AddShareRecipeFunc(recipeID, userID)
	}

	recipes, ok := m.RecipesRegistered[userID]
	if !ok {
		return 0, errors.New("user not registered")
	}

	okRecipe := slices.ContainsFunc(recipes, func(r models.Recipe) bool {
		return r.ID == recipeID
	})
	if okRecipe {
		return 0, errors.New("recipe exists")
	}

	return 2, nil
}

func (m *mockRepository) AddRecipeCategory(name string, userID int64) error {
	if m.AddRecipeCategoryFunc != nil {
		return m.AddRecipeCategoryFunc(name, userID)
	}

	categories, err := m.Categories(userID)
	if err != nil {
		return err
	}

	if slices.Contains(categories, name) {
		return errors.New("category in use")
	}

	m.categories[userID] = append(m.categories[userID], name)
	return nil
}

func (m *mockRepository) AddReport(report models.Report, userID int64) {
	_, ok := m.Reports[userID]
	if !ok {
		panic("reports for user not initialized")
	}

	m.Reports[userID] = append(m.Reports[userID], report)
}

func (m *mockRepository) CookbooksShared(_ int64) ([]models.Share, error) {
	return make([]models.Share, 0), nil
}

func (m *mockRepository) CookbooksUser(_ int64) ([]models.Cookbook, error) {
	return make([]models.Cookbook, 0), nil
}

func (m *mockRepository) AddAuthToken(selector, validator string, userID int64) error {
	token := models.NewAuthToken(int64(len(m.AuthTokens)+1), selector, validator, 10000, userID)
	m.AuthTokens = append(m.AuthTokens, *token)
	return nil
}

func (m *mockRepository) AddCookbook(title string, userID int64) (int64, error) {
	cookbook := models.Cookbook{
		Recipes: make(models.Recipes, 0),
		Title:   title,
	}

	cookbooks, ok := m.CookbooksRegistered[userID]
	if !ok {
		m.CookbooksRegistered[userID] = []models.Cookbook{cookbook}
		return 1, nil
	}

	isExists := slices.ContainsFunc(cookbooks, func(cookbook models.Cookbook) bool {
		return cookbook.Title == title
	})
	if isExists {
		return -1, errors.New("cookbook exists")
	}

	cookbook.ID = int64(len(cookbooks) + 1)
	m.CookbooksRegistered[userID] = append(m.CookbooksRegistered[userID], cookbook)
	return 1, nil
}

func (m *mockRepository) AddCookbookRecipe(cookbookID, recipeID, userID int64) error {
	cookbooks, ok := m.CookbooksRegistered[userID]
	if !ok {
		return errors.New("cookbook does not belong to user")
	}

	if cookbooks == nil {
		return errors.New("user cookbooks registered is nil")
	}

	cookbookIndex := slices.IndexFunc(cookbooks, func(c models.Cookbook) bool {
		return c.ID == cookbookID
	})
	if cookbookIndex == -1 {
		return errors.New("cookbook not found")
	}

	recipes := m.RecipesRegistered[userID]
	if recipes == nil {
		return errors.New("user recipes is nil")
	}

	recipeIndex := slices.IndexFunc(recipes, func(r models.Recipe) bool {
		return r.ID == recipeID
	})
	if recipeIndex == -1 {
		return errors.New("recipe not found")
	}

	cookbooks[cookbookIndex].Recipes = append(cookbooks[cookbookIndex].Recipes, recipes[recipeID])
	return nil
}

func (m *mockRepository) Categories(userID int64) ([]string, error) {
	categories, ok := m.categories[userID]
	if !ok {
		return nil, errors.New("categories not found")
	}
	return categories, nil
}

func (m *mockRepository) CheckUpdate(_ services.FilesService) (models.AppInfo, error) {
	lastCheckedAt, _ := time.Parse(time.DateTime, "2021-06-18 20:30:05")
	lastUpdatedAt, _ := time.Parse(time.DateTime, "2021-02-24 15:04:05")

	return models.AppInfo{
		IsUpdateAvailable:   false,
		LastUpdatedAt:       lastUpdatedAt,
		LastCheckedUpdateAt: lastCheckedAt,
	}, nil
}

func (m *mockRepository) Confirm(userID int64) error {
	if !slices.ContainsFunc(m.UsersRegistered, func(user models.User) bool {
		return user.ID == userID
	}) {
		return sql.ErrNoRows
	}
	return nil
}

func (m *mockRepository) Cookbook(id, userID int64) (models.Cookbook, error) {
	cookbooks, ok := m.CookbooksRegistered[userID]
	if !ok {
		return models.Cookbook{}, errors.New("user does not have cookbooks")
	}

	i := slices.IndexFunc(cookbooks, func(c models.Cookbook) bool {
		return c.ID == id
	})
	if i == -1 {
		return models.Cookbook{}, errors.New("cookbook not found")
	}

	return cookbooks[i], nil
}

func (m *mockRepository) CookbookRecipe(id int64, cookbookID int64) (recipe *models.Recipe, userID int64, err error) {
	for userID, cookbooks := range m.CookbooksRegistered {
		i := slices.IndexFunc(cookbooks, func(c models.Cookbook) bool {
			return c.ID == cookbookID
		})
		if i == -1 {
			continue
		}

		recipeI := slices.IndexFunc(cookbooks[i].Recipes, func(r models.Recipe) bool {
			return r.ID == id
		})
		if recipeI == -1 {
			break
		}

		cookbook := cookbooks[i]
		recipes := cookbook.Recipes
		if recipes == nil {
			return nil, 0, errors.New("user recipes in cookbook is nil")
		}
		return &recipes[recipeI], userID, nil
	}
	return nil, -1, errors.New("recipe not found")
}

func (m *mockRepository) CookbookShared(id string) (*models.Share, error) {
	share, ok := m.ShareLinks[id]
	if !ok {
		return nil, errors.New("link not found")
	}
	return &share, nil
}

func (m *mockRepository) Cookbooks(userID int64, _ uint64) ([]models.Cookbook, error) {
	if m.CookbooksFunc != nil {
		return m.CookbooksFunc(userID)
	}

	cookbooks, ok := m.CookbooksRegistered[userID]
	if !ok {
		return nil, errors.New("user not registered")
	}
	return cookbooks, nil
}

func (m *mockRepository) Counts(userID int64) (models.Counts, error) {
	var counts models.Counts
	recipes, ok := m.RecipesRegistered[userID]
	if ok {
		counts.Recipes = uint64(len(recipes))
	}

	cookbooks, ok := m.CookbooksRegistered[userID]
	if ok {
		counts.Cookbooks = uint64(len(cookbooks))
	}

	return counts, nil
}

func (m *mockRepository) DeleteAuthToken(userID int64) error {
	index := slices.IndexFunc(m.AuthTokens, func(token models.AuthToken) bool { return token.UserID == userID })
	if index != -1 {
		m.AuthTokens = slices.Delete(m.AuthTokens, index, index+1)
	}
	return nil
}

func (m *mockRepository) DeleteRecipeCategory(name string, userID int64) error {
	if m.DeleteCategoryFunc != nil {
		return m.DeleteCategoryFunc(name, userID)
	}

	if name == "uncategorized" {
		return errors.New("cannot delete 'uncategorized'")
	}

	recipes, ok := m.RecipesRegistered[userID]
	if !ok {
		return errors.New("user not registered")
	}

	for i, r := range recipes {
		if r.Category == name {
			r.Category = "uncategorized"
		}
		recipes[i] = r
	}

	m.RecipesRegistered[userID] = recipes
	return nil
}

func (m *mockRepository) DeleteCookbook(id, userID int64) error {
	if m.DeleteCookbookFunc != nil {
		return m.DeleteCookbookFunc(id, userID)
	}

	m.CookbooksRegistered[userID] = slices.DeleteFunc(m.CookbooksRegistered[userID], func(c models.Cookbook) bool {
		return c.ID == id
	})
	return nil
}

func (m *mockRepository) DeleteRecipe(id, userID int64) error {
	recipes, ok := m.RecipesRegistered[userID]
	if !ok {
		return errors.New("user not found")
	}

	i := slices.IndexFunc(recipes, func(r models.Recipe) bool {
		return r.ID == id
	})
	if i == -1 {
		return nil
	}

	m.RecipesRegistered[userID] = slices.Delete(recipes, i, i+1)
	return nil
}

func (m *mockRepository) DeleteRecipeFromCookbook(recipeID, cookbookID int64, userID int64) (int64, error) {
	cookbooks, ok := m.CookbooksRegistered[userID]
	if !ok {
		return -1, nil
	}

	if cookbooks == nil {
		return 0, errors.New("user cookbooks registered is nil")
	}

	i := slices.IndexFunc(cookbooks, func(c models.Cookbook) bool {
		return c.ID == cookbookID
	})
	if i == -1 {
		return -1, nil
	}
	cookbook := cookbooks[i]

	cookbook.Recipes = slices.DeleteFunc(cookbook.Recipes, func(r models.Recipe) bool {
		return r.ID == recipeID
	})

	cookbooks[i] = cookbook
	return int64(len(cookbook.Recipes)), nil
}

func (m *mockRepository) DeleteUser(id int64) error {
	m.UsersRegistered = slices.DeleteFunc(m.UsersRegistered, func(user models.User) bool {
		return user.ID == id
	})
	return nil
}

func (m *mockRepository) GetAuthToken(_, _ string) (models.AuthToken, error) {
	return models.AuthToken{UserID: 1, Expires: time.Now().Add(1 * time.Hour)}, nil
}

func (m *mockRepository) Media() (images, videos []string) {
	return make([]string, 0), make([]string, 0)
}

func (m *mockRepository) InitAutologin() error {
	return nil
}

func (m *mockRepository) IsUserExist(email string) bool {
	return slices.ContainsFunc(m.UsersRegistered, func(user models.User) bool {
		return user.Email == email
	})
}

func (m *mockRepository) IsUserPassword(id int64, password string) bool {
	if m.IsUserPasswordFunc != nil {
		return m.IsUserPasswordFunc(id, password)
	}
	return slices.IndexFunc(m.UsersRegistered, func(user models.User) bool { return user.ID == id }) != -1
}

func (m *mockRepository) Keywords() ([]string, error) {
	return []string{"big"}, nil
}

func (m *mockRepository) MeasurementSystems(userID int64) ([]units.System, models.UserSettings, error) {
	if m.MeasurementSystemsFunc != nil {
		return m.MeasurementSystemsFunc(userID)
	}
	return []units.System{units.ImperialSystem, units.MetricSystem}, models.UserSettings{
		ConvertAutomatically: false,
		MeasurementSystem:    units.MetricSystem,
	}, nil
}

func (m *mockRepository) Nutrients(_ []string) (models.NutrientsFDC, float64, error) {
	return models.NutrientsFDC{}, 0, nil
}

func (m *mockRepository) Recipe(id, userID int64) (*models.Recipe, error) {
	if m.RecipeFunc != nil {
		return m.RecipeFunc(id, userID)
	}

	if recipes, ok := m.RecipesRegistered[userID]; ok {
		if int64(len(recipes)) < id {
			return nil, errors.New("recipe not found")
		}
		return &recipes[id-1], nil
	}
	return nil, errors.New("recipe not found")
}

func (m *mockRepository) Recipes(userID int64, opts models.SearchOptionsRecipes) models.Recipes {
	if recipes, ok := m.RecipesRegistered[userID]; ok {
		return recipes
	}
	return models.Recipes{}
}

func (m *mockRepository) RecipesAll(userID int64) models.Recipes {
	if recipes, ok := m.RecipesRegistered[userID]; ok {
		return recipes
	}
	return models.Recipes{}
}

func (m *mockRepository) RecipeShared(link string) (*models.Share, error) {
	share, ok := m.ShareLinks[link]
	if !ok {
		return nil, errors.New("recipe not found")
	}
	return &share, nil
}

func (m *mockRepository) RecipeUser(recipeID int64) int64 {
	for userID, recipes := range m.RecipesRegistered {
		i := slices.IndexFunc(recipes, func(r models.Recipe) bool {
			return r.ID == recipeID
		})
		if i != -1 {
			return userID
		}
	}
	return -1
}

func (m *mockRepository) Register(email string, _ auth.HashedPassword) (int64, error) {
	if slices.ContainsFunc(m.UsersRegistered, func(user models.User) bool {
		return user.Email == email
	}) {
		return -1, errors.New("email taken")
	}

	userID := int64(len(m.UsersRegistered) + 1)
	m.UsersRegistered = append(m.UsersRegistered, models.User{
		ID:    userID,
		Email: email,
	})
	return userID, nil
}

func (m *mockRepository) ReorderCookbookRecipes(_ int64, _ []uint64, _ int64) error {
	return nil
}

func (m *mockRepository) RecipesShared(_ int64) ([]models.Share, error) {
	return make([]models.Share, 0), nil
}

func (m *mockRepository) Report(id, userID int64) ([]models.ReportLog, error) {
	reports, ok := m.Reports[userID]
	if !ok {
		return []models.ReportLog{}, nil
	}

	i := slices.IndexFunc(reports, func(r models.Report) bool { return r.ID == id })
	if i == -1 {
		return []models.ReportLog{}, errors.New("report not found")
	}

	return reports[i].Logs, nil
}

func (m *mockRepository) ReportsImport(userID int64) ([]models.Report, error) {
	if m.ReportsFunc != nil {
		return m.ReportsFunc(userID)
	}

	reports, ok := m.Reports[userID]
	if !ok {
		return []models.Report{}, nil
	}
	return reports, nil
}

func (m *mockRepository) RestoreBackup(_ string) error {
	return nil
}

func (m *mockRepository) RestoreUserBackup(backup *models.UserBackup) error {
	if m.RestoreUserBackupFunc != nil {
		return m.RestoreUserBackupFunc(backup)
	}
	return nil
}

func (m *mockRepository) SearchRecipes(opts models.SearchOptionsRecipes, userID int64) (models.Recipes, uint64, error) {
	recipes, ok := m.RecipesRegistered[userID]
	if !ok {
		return nil, 0, errors.New("user not found")
	}

	q := strings.ReplaceAll(opts.Query, `"`, "")
	var results models.Recipes
	for _, r := range recipes {
		if strings.Contains(strings.ToLower(r.Name), q) || strings.Contains(strings.ToLower(r.Category), q) || strings.Contains(strings.ToLower(r.Description), q) {
			results = append(results, r)
		}
	}
	return results, uint64(len(results)), nil
}

func (m *mockRepository) SwitchMeasurementSystem(system units.System, userID int64) error {
	if m.SwitchMeasurementSystemFunc != nil {
		return m.SwitchMeasurementSystemFunc(system, userID)
	}

	for i, r := range m.RecipesRegistered[userID] {
		converted, err := r.ConvertMeasurementSystem(system)
		if err != nil {
			return err
		}

		recipes := m.RecipesRegistered[userID]
		if recipes == nil || converted == nil {
			return errors.New("user recipes or converted recipe is nil")
		}
		recipes[i] = *converted
	}
	return nil
}

func (m *mockRepository) UpdateCalculateNutrition(userID int64, isEnabled bool) error {
	if m.UpdateCalculateNutritionFunc != nil {
		return m.UpdateCalculateNutritionFunc(userID, isEnabled)
	}

	settings, ok := m.UserSettingsRegistered[userID]
	if !ok {
		return errors.New("user not found")
	}

	if settings == nil {
		return errors.New("settings for user is empty")
	}

	settings.CalculateNutritionFact = isEnabled
	return nil
}

func (m *mockRepository) UpdateConvertMeasurementSystem(userID int64, isEnabled bool) error {
	if m.UpdateConvertMeasurementSystemFunc != nil {
		return m.UpdateConvertMeasurementSystemFunc(userID, isEnabled)
	}

	settings, ok := m.UserSettingsRegistered[userID]
	if !ok {
		return errors.New("user not found")
	}

	if settings == nil {
		return errors.New("settings for user is empty")
	}

	settings.ConvertAutomatically = isEnabled
	return nil
}

func (m *mockRepository) UpdateCookbookImage(id int64, image uuid.UUID, userID int64) error {
	if m.UpdateCookbookImageFunc != nil {
		return m.UpdateCookbookImageFunc(id, image, userID)
	}

	_, ok := m.CookbooksRegistered[userID]
	if !ok {
		return errors.New("cookbook not found")
	}

	cookbooks := m.CookbooksRegistered[userID]
	if cookbooks == nil {
		return errors.New("CookbooksRegistered nil for user")
	}

	for i, cookbook := range cookbooks {
		if cookbook.ID == id {
			c := cookbooks[i]
			userCookbooks := m.CookbooksRegistered[userID]
			if userCookbooks == nil {
				return errors.New("user has no cookbooks registered")
			}

			userCookbooks[i] = models.Cookbook{
				ID:      c.ID,
				Count:   c.Count,
				Image:   image,
				Recipes: c.Recipes,
				Title:   c.Title,
			}
			return nil
		}
	}
	return errors.New("cookbook not found")
}

func (m *mockRepository) UpdatePassword(userID int64, _ auth.HashedPassword) error {
	m.UsersUpdated = append(m.UsersUpdated, userID)
	return nil
}

func (m *mockRepository) UpdateRecipe(updatedRecipe *models.Recipe, userID int64, recipeNum int64) error {
	oldRecipe, err := m.Recipe(recipeNum, userID)
	if err != nil {
		return err
	}

	newRecipe := *oldRecipe

	if oldRecipe.Category != updatedRecipe.Category {
		if updatedRecipe.Category == "" {
			updatedRecipe.Category = "uncategorized"
		} else {
			delimiters := []string{",", ";", "|"}
			for _, delim := range delimiters {
				if strings.Contains(updatedRecipe.Category, delim) {
					updatedRecipe.Category = strings.Split(updatedRecipe.Category, delim)[0]
				}
			}
		}
		newRecipe.Category = updatedRecipe.Category
	}

	if oldRecipe.Cuisine != updatedRecipe.Cuisine {
		newRecipe.Cuisine = updatedRecipe.Cuisine
	}

	if oldRecipe.Description != updatedRecipe.Description {
		newRecipe.Description = updatedRecipe.Description
	}

	if len(updatedRecipe.Images) > 0 && !slices.Equal(oldRecipe.Images, updatedRecipe.Images) {
		newRecipe.Images = updatedRecipe.Images
	}

	if newRecipe.Ingredients != nil && len(oldRecipe.Ingredients) == len(updatedRecipe.Ingredients) {
		for i, ingredient := range updatedRecipe.Ingredients {
			if oldRecipe.Ingredients[i] != updatedRecipe.Ingredients[i] {
				newRecipe.Ingredients[i] = ingredient
			}
		}
	} else {
		if len(updatedRecipe.Ingredients) == 0 {
			updatedRecipe.Ingredients = newRecipe.Ingredients
		}

		cloned := slices.Clone(updatedRecipe.Ingredients)
		if cloned != nil {
			newRecipe.Ingredients = cloned
		}
	}

	if newRecipe.Instructions != nil && len(oldRecipe.Instructions) == len(updatedRecipe.Instructions) {
		for i, ingredient := range updatedRecipe.Instructions {
			if oldRecipe.Instructions[i] != updatedRecipe.Instructions[i] {
				newRecipe.Instructions[i] = ingredient
			}
		}
	} else {
		if len(updatedRecipe.Instructions) == 0 {
			updatedRecipe.Instructions = newRecipe.Instructions
		}

		cloned := slices.Clone(updatedRecipe.Instructions)
		if cloned != nil {
			newRecipe.Instructions = cloned
		}
	}

	if newRecipe.Keywords != nil && len(oldRecipe.Keywords) == len(updatedRecipe.Keywords) {
		for i, ingredient := range updatedRecipe.Keywords {
			if oldRecipe.Keywords[i] != updatedRecipe.Keywords[i] {
				newRecipe.Keywords[i] = ingredient
			}
		}
	} else {
		cloned := slices.Clone(updatedRecipe.Keywords)
		if cloned != nil {
			newRecipe.Keywords = cloned
		}
	}

	if oldRecipe.Name != updatedRecipe.Name {
		if updatedRecipe.Name == "" {
			updatedRecipe.Name = oldRecipe.Name
		}
		newRecipe.Name = updatedRecipe.Name
	}

	// To save some lines...
	newRecipe.Nutrition = updatedRecipe.Nutrition

	if oldRecipe.Times.Prep != updatedRecipe.Times.Prep {
		newRecipe.Times.Prep = updatedRecipe.Times.Prep
	}

	if oldRecipe.Times.Cook != updatedRecipe.Times.Cook {
		newRecipe.Times.Cook = updatedRecipe.Times.Cook
	}

	if oldRecipe.Times.Total != updatedRecipe.Times.Total {
		newRecipe.Times.Total = updatedRecipe.Times.Total
	}

	if oldRecipe.URL != updatedRecipe.URL {
		if updatedRecipe.URL == "" {
			updatedRecipe.URL = "Unknown"
		}
		newRecipe.URL = updatedRecipe.URL
	}

	if oldRecipe.Yield != updatedRecipe.Yield {
		if updatedRecipe.Yield == 0 {
			updatedRecipe.Yield = 1
		}
		newRecipe.Yield = updatedRecipe.Yield
	}

	recipes := m.RecipesRegistered[userID]
	if recipes == nil {
		return errors.New("user has no recipes")
	}

	recipes[oldRecipe.ID-1] = newRecipe
	return nil
}

func (m *mockRepository) UpdateUserSettingsCookbooksViewMode(userID int64, mode models.ViewMode) error {
	settings, ok := m.UserSettingsRegistered[userID]
	if !ok {
		return errors.New("user not found")
	}

	settings.CookbooksViewMode = mode
	return nil
}

func (m *mockRepository) UpdateVideo(_ uuid.UUID, _ int) error {
	return nil
}

func (m *mockRepository) UserInitials(userID int64) string {
	index := slices.IndexFunc(m.UsersRegistered, func(user models.User) bool {
		return user.ID == userID
	})
	if index == -1 {
		return ""
	}
	return string(strings.ToUpper(m.UsersRegistered[index].Email)[0])
}

func (m *mockRepository) UserID(email string) int64 {
	index := slices.IndexFunc(m.UsersRegistered, func(user models.User) bool {
		return user.Email == email
	})
	if index == -1 {
		return -1
	}
	return m.UsersRegistered[index].ID
}

func (m *mockRepository) UserSettings(userID int64) (models.UserSettings, error) {
	settings, ok := m.UserSettingsRegistered[userID]
	if !ok {
		return models.UserSettings{}, errors.New("user not found")
	}

	if settings == nil {
		return models.UserSettings{}, errors.New("settings is nil")
	}

	return *settings, nil
}

func (m *mockRepository) Users() []models.User {
	return m.UsersRegistered
}

func (m *mockRepository) VerifyLogin(email, _ string) int64 {
	index := slices.IndexFunc(m.UsersRegistered, func(user models.User) bool {
		return user.Email == email
	})

	if index == -1 {
		return -1
	}
	return m.UsersRegistered[index].ID
}

func (m *mockRepository) Websites() models.Websites {
	return models.Websites{
		{ID: 1, Host: "101cookbooks.com", URL: "https://101cookbooks.com"},
		{ID: 2, Host: "afghankitchenrecipes.com", URL: "https://www.afghankitchenrecipes.com"},
	}
}

type mockEmail struct {
	hitCount int64
}

func (m *mockEmail) Queue(_ string, _ templates.EmailTemplate, _ any) {}

func (m *mockEmail) RateLimits() (remaining int, resetUnix int64, err error) {
	return remaining, resetUnix, nil
}

func (m *mockEmail) SendQueue() (sent, remaining int, err error) {
	return sent, remaining, nil
}

func (m *mockEmail) Send(_ string, _ templates.EmailTemplate, _ any) error {
	m.hitCount++
	return nil
}

type mockFiles struct {
	backupUserDataFunc    func(repo services.RepositoryService, userID int64) error
	exportHitCount        int
	extractRecipesFunc    func(fileHeaders []*multipart.FileHeader) models.Recipes
	extractUserBackupFunc func(date string, userID int64) (*models.UserBackup, error)
	ReadTempFileFunc      func(name string) ([]byte, error)
	updateAppFunc         func(current semver.Version) error
	uploadImageHitCount   int
	uploadImageFunc       func(rc io.ReadCloser) (uuid.UUID, error)
}

func (m *mockFiles) ExtractUserBackup(date string, userID int64) (*models.UserBackup, error) {
	if m.extractUserBackupFunc != nil {
		return m.extractUserBackupFunc(date, userID)
	}
	return &models.UserBackup{
		DeleteSQL:  "delete SQL",
		ImagesPath: "/path/to/images",
		InsertSQL:  "insert SQL",
		Recipes:    make(models.Recipes, 0),
		UserID:     userID,
	}, nil
}

func (m *mockFiles) BackupGlobal() error {
	return nil
}

func (m *mockFiles) Backups(_ int64) []time.Time {
	return nil
}

func (m *mockFiles) BackupUserData(repo services.RepositoryService, userID int64) error {
	if m.backupUserDataFunc != nil {
		return m.backupUserDataFunc(repo, userID)
	}
	return nil
}

func (m *mockFiles) BackupUsersData(_ services.RepositoryService) error {
	return nil
}

func (m *mockFiles) ExportCookbook(cookbook models.Cookbook, fileType models.FileType) (string, error) {
	m.exportHitCount++
	return cookbook.Title + fileType.Ext(), nil
}

func (m *mockFiles) ExportRecipes(recipes models.Recipes, _ models.FileType, _ chan int) (*bytes.Buffer, error) {
	var b bytes.Buffer
	for _, recipe := range recipes {
		b.WriteString(recipe.Name + "-")
	}
	m.exportHitCount++
	return &b, nil
}

func (m *mockFiles) ExtractRecipes(fileHeaders []*multipart.FileHeader) models.Recipes {
	if m.extractRecipesFunc != nil {
		return m.extractRecipesFunc(fileHeaders)
	}
	return models.Recipes{}
}

func (m *mockFiles) IsAppLatest(_ semver.Version) (bool, *github.RepositoryRelease, error) {
	return true, nil, nil
}

func (m *mockFiles) ReadTempFile(name string) ([]byte, error) {
	if m.ReadTempFileFunc != nil {
		return m.ReadTempFileFunc(name)
	}
	return []byte(name), nil
}

func (m *mockFiles) MergeImagesToPDF(images []io.Reader) io.ReadWriter {
	return nil
}

func (m *mockFiles) ScrapeAndStoreImage(_ string) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (m *mockFiles) UpdateApp(current semver.Version) error {
	if m.updateAppFunc != nil {
		return m.updateAppFunc(current)
	}
	return nil
}

func (m *mockFiles) UploadImage(rc io.ReadCloser) (uuid.UUID, error) {
	if m.uploadImageFunc != nil {
		return m.uploadImageFunc(rc)
	}
	m.uploadImageHitCount++
	return uuid.New(), nil
}

func (m *mockFiles) UploadVideo(_ io.ReadCloser, _ services.RepositoryService) (uuid.UUID, error) {
	m.uploadImageHitCount++
	return uuid.New(), nil
}

type mockIntegrations struct {
	importFunc          func(baseURL, username, password string, files services.FilesService) (models.Recipes, error)
	processImageOCRFunc func(file []io.Reader) (models.Recipes, error)
	testConnectionFunc  func(api string) error
}

func (m *mockIntegrations) MealieImport(baseURL, username, password string, files services.FilesService, progress chan models.Progress) (models.Recipes, error) {
	return m.importIntegration(baseURL, username, password, files)
}

func (m *mockIntegrations) NextcloudImport(baseURL, username, password string, files services.FilesService, _ chan models.Progress) (models.Recipes, error) {
	return m.importIntegration(baseURL, username, password, files)
}

func (m *mockIntegrations) TandoorImport(baseURL, username, password string, files services.FilesService, _ chan models.Progress) (models.Recipes, error) {
	return m.importIntegration(baseURL, username, password, files)
}

func (m *mockIntegrations) importIntegration(baseURL, username, password string, files services.FilesService) (models.Recipes, error) {
	if username == "" || password == "" || baseURL == "" {
		return nil, errors.New("invalid username, password or URL")
	}

	if m.importFunc != nil {
		return m.importFunc(baseURL, username, password, files)
	}

	return models.Recipes{
		{ID: 1, Name: "One"},
		{ID: 2, Name: "Two"},
	}, nil
}

func (m *mockIntegrations) ProcessImageOCR(f []io.Reader) (models.Recipes, error) {
	if m.processImageOCRFunc != nil {
		return m.processImageOCRFunc(f)
	}
	return models.Recipes{{ID: 1}}, nil
}

func (m *mockIntegrations) TestConnection(api string) error {
	if m.testConnectionFunc != nil {
		return m.testConnectionFunc(api)
	}

	switch api {
	case "azure-di", "sg":
		return nil
	default:
		return errors.New("invalid api")
	}
}

type mockScraper struct {
	scraperFunc func(url string, files services.FilesService) (models.RecipeSchema, error)
}

func (m *mockScraper) Scrape(url string, files services.FilesService) (models.RecipeSchema, error) {
	if m.scraperFunc != nil {
		return m.scraperFunc(url, files)
	}

	return models.RecipeSchema{
		AtContext: "https://schema.org",
		AtType:    &models.SchemaType{Value: "Recipe"},
		Name:      url,
		URL:       url,
	}, nil
}
