package server_test

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"github.com/reaper47/recipya/internal/services"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/internal/units"
	"io"
	"mime/multipart"
	"net/http"
	"slices"
	"strings"
)

func newServerTest() *server.Server {
	repo := &mockRepository{
		AuthTokens:             make([]models.AuthToken, 0),
		RecipesRegistered:      make(map[int64]models.Recipes),
		ShareLinks:             make(map[string]models.ShareRecipe),
		UserSettingsRegistered: make(map[int64]*models.UserSettings),
		UsersRegistered:        make([]models.User, 0),
		UsersUpdated:           make([]int64, 0),
	}
	return server.NewServer(repo, &mockEmail{}, &mockFiles{}, &mockIntegrations{})
}

type mockRepository struct {
	AuthTokens                  []models.AuthToken
	AddRecipeFunc               func(recipe *models.Recipe, userID int64) (int64, error)
	CookbooksFunc               func(userID int64) ([]models.Cookbook, error)
	CookbooksRegistered         map[int64][]models.Cookbook
	MeasurementSystemsFunc      func(userID int64) ([]units.System, models.UserSettings, error)
	RecipeFunc                  func(id, userID int64) (*models.Recipe, error)
	RecipesRegistered           map[int64]models.Recipes
	ShareLinks                  map[string]models.ShareRecipe
	SwitchMeasurementSystemFunc func(system units.System, userID int64) error
	UserSettingsRegistered      map[int64]*models.UserSettings
	UsersRegistered             []models.User
	UsersUpdated                []int64
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

	m.CookbooksRegistered[userID] = append(m.CookbooksRegistered[userID], cookbook)
	return 1, nil
}

func (m *mockRepository) AddRecipe(r *models.Recipe, userID int64) (int64, error) {
	if m.AddRecipeFunc != nil {
		return m.AddRecipeFunc(r, userID)
	}

	if m.RecipesRegistered[userID] == nil {
		m.RecipesRegistered[userID] = make(models.Recipes, 0)
	}

	if !slices.ContainsFunc(m.RecipesRegistered[userID], func(recipe models.Recipe) bool {
		return recipe.ID == r.ID
	}) {
		m.RecipesRegistered[userID] = append(m.RecipesRegistered[userID], *r)
	}
	return int64(len(m.RecipesRegistered)), nil
}

func (m *mockRepository) AddShareLink(share models.ShareRecipe) (string, error) {
	for _, recipes := range m.RecipesRegistered {
		if slices.ContainsFunc(recipes, func(r models.Recipe) bool { return r.ID == share.RecipeID }) {
			for link, s := range m.ShareLinks {
				if s.RecipeID == share.RecipeID && s.UserID == share.UserID {
					return link, nil
				}
			}

			link := "/r/33320755-82f9-47e5-bb0a-d1b55cbd3f7b"
			m.ShareLinks[link] = share
			return link, nil
		}
	}
	return "", errors.New("recipe not found")
}

func (m *mockRepository) Categories(_ int64) ([]string, error) {
	return []string{"breakfast", "lunch", "dinner"}, nil
}

func (m *mockRepository) Confirm(userID int64) error {
	if !slices.ContainsFunc(m.UsersRegistered, func(user models.User) bool {
		return user.ID == userID
	}) {
		return sql.ErrNoRows
	}
	return nil
}

func (m *mockRepository) Cookbooks(userID int64) ([]models.Cookbook, error) {
	if m.CookbooksFunc != nil {
		return m.CookbooksFunc(userID)
	}

	cookbooks, ok := m.CookbooksRegistered[userID]
	if !ok {
		return nil, errors.New("user not registered")
	}
	return cookbooks, nil
}

func (m *mockRepository) DeleteAuthToken(userID int64) error {
	index := slices.IndexFunc(m.AuthTokens, func(token models.AuthToken) bool { return token.UserID == userID })
	if index != -1 {
		m.AuthTokens = slices.Delete(m.AuthTokens, index, index+1)
	}
	return nil
}

func (m *mockRepository) DeleteRecipe(id, userID int64) (int64, error) {
	recipes, ok := m.RecipesRegistered[userID]
	if !ok {
		return -1, errors.New("user not found")
	}

	var rowsAffected int64
	i := slices.IndexFunc(recipes, func(r models.Recipe) bool {
		if r.ID == id {
			rowsAffected++
		}
		return r.ID == id
	})
	if i == -1 {
		return 0, nil
	}

	slices.Delete(recipes, i, i+1)
	recipes = recipes[:]
	return rowsAffected, nil
}

func (m *mockRepository) GetAuthToken(_, _ string) (models.AuthToken, error) {
	return models.AuthToken{UserID: 1}, nil
}

func (m *mockRepository) IsUserExist(email string) bool {
	return slices.ContainsFunc(m.UsersRegistered, func(user models.User) bool {
		return user.Email == email
	})
}

func (m *mockRepository) IsUserPassword(id int64, _ string) bool {
	return slices.IndexFunc(m.UsersRegistered, func(user models.User) bool { return user.ID == id }) != -1
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

func (m *mockRepository) Recipes(userID int64) models.Recipes {
	if recipes, ok := m.RecipesRegistered[userID]; ok {
		return recipes
	}
	return models.Recipes{}
}

func (m *mockRepository) RecipeShared(link string) (*models.ShareRecipe, error) {
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

func (m *mockRepository) SwitchMeasurementSystem(system units.System, userID int64) error {
	if m.SwitchMeasurementSystemFunc != nil {
		return m.SwitchMeasurementSystemFunc(system, userID)
	}

	for i, r := range m.RecipesRegistered[userID] {
		converted, err := r.ConvertMeasurementSystem(system)
		if err != nil {
			return err
		}
		m.RecipesRegistered[userID][i] = *converted
	}
	return nil
}

func (m *mockRepository) UpdateConvertMeasurementSystem(userID int64, isEnabled bool) error {
	if _, ok := m.UserSettingsRegistered[userID]; !ok {
		return errors.New("user not found")
	}
	m.UserSettingsRegistered[userID].ConvertAutomatically = isEnabled
	return nil
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
		newRecipe.Category = updatedRecipe.Category
	}

	if oldRecipe.Cuisine != updatedRecipe.Cuisine {
		newRecipe.Cuisine = updatedRecipe.Cuisine
	}

	if oldRecipe.Description != updatedRecipe.Description {
		newRecipe.Description = updatedRecipe.Description
	}

	if updatedRecipe.Image != uuid.Nil && oldRecipe.Image != updatedRecipe.Image {
		newRecipe.Image = updatedRecipe.Image
	}

	if len(oldRecipe.Ingredients) == len(updatedRecipe.Ingredients) {
		for i, ingredient := range updatedRecipe.Ingredients {
			if oldRecipe.Ingredients[i] != updatedRecipe.Ingredients[i] {
				newRecipe.Ingredients[i] = ingredient
			}
		}
	} else {
		newRecipe.Ingredients = slices.Clone(updatedRecipe.Ingredients)
	}

	if len(oldRecipe.Instructions) == len(updatedRecipe.Instructions) {
		for i, ingredient := range updatedRecipe.Instructions {
			if oldRecipe.Instructions[i] != updatedRecipe.Instructions[i] {
				newRecipe.Instructions[i] = ingredient
			}
		}
	} else {
		newRecipe.Instructions = slices.Clone(updatedRecipe.Instructions)
	}

	if len(oldRecipe.Keywords) == len(updatedRecipe.Keywords) {
		for i, ingredient := range updatedRecipe.Keywords {
			if oldRecipe.Keywords[i] != updatedRecipe.Keywords[i] {
				newRecipe.Keywords[i] = ingredient
			}
		}
	} else {
		newRecipe.Keywords = slices.Clone(updatedRecipe.Keywords)
	}

	if oldRecipe.Name != updatedRecipe.Name {
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
		newRecipe.URL = updatedRecipe.URL
	}

	if oldRecipe.Yield != updatedRecipe.Yield {
		newRecipe.Yield = updatedRecipe.Yield
	}

	m.RecipesRegistered[userID][oldRecipe.ID-1] = newRecipe
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
	if _, ok := m.UserSettingsRegistered[userID]; !ok {
		return models.UserSettings{}, errors.New("user not found")
	}
	return *m.UserSettingsRegistered[userID], nil
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

func (m *mockEmail) Send(_ string, _ templates.EmailTemplate, _ any) {
	m.hitCount += 1
}

type mockFiles struct {
	exportHitCount      int
	extractRecipesFunc  func(fileHeaders []*multipart.FileHeader) models.Recipes
	ReadTempFileFunc    func(name string) ([]byte, error)
	uploadImageHitCount int
}

func (m *mockFiles) ExportRecipes(recipes models.Recipes, _ models.FileType) (string, error) {
	var s string
	for _, recipe := range recipes {
		s += recipe.Name + "-"
	}
	m.exportHitCount++
	return s, nil
}

func (m *mockFiles) ExtractRecipes(fileHeaders []*multipart.FileHeader) models.Recipes {
	if m.extractRecipesFunc != nil {
		return m.extractRecipesFunc(fileHeaders)
	}
	return models.Recipes{}
}

func (m *mockFiles) ReadTempFile(name string) ([]byte, error) {
	if m.ReadTempFileFunc != nil {
		return m.ReadTempFileFunc(name)
	}
	return []byte(name), nil
}

func (m *mockFiles) UploadImage(_ io.ReadCloser) (uuid.UUID, error) {
	m.uploadImageHitCount++
	return uuid.New(), nil
}

type mockIntegrations struct {
	NextcloudImportFunc func(client *http.Client, baseURL, username, password string, files services.FilesService) (*models.Recipes, error)
}

func (m *mockIntegrations) NextcloudImport(client *http.Client, baseURL, username, password string, files services.FilesService) (*models.Recipes, error) {
	if m.NextcloudImportFunc != nil {
		return m.NextcloudImportFunc(client, baseURL, username, password, files)
	}
	return &models.Recipes{
		{ID: 1, Name: "One"},
		{ID: 2, Name: "Two"},
	}, nil
}
