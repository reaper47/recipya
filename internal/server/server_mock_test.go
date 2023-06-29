package server_test

import (
	"database/sql"
	"errors"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"github.com/reaper47/recipya/internal/templates"
	"golang.org/x/exp/slices"
	"mime/multipart"
	"strings"
)

func newServerTest() *server.Server {
	repo := &mockRepository{
		AuthTokens:      make([]models.AuthToken, 0),
		Recipes:         make(map[int64]models.Recipes, 0),
		UsersRegistered: make([]models.User, 0),
	}
	return server.NewServer(repo, &mockEmail{}, &mockFiles{})
}

type mockRepository struct {
	AuthTokens      []models.AuthToken
	AddRecipeFunc   func(recipe *models.Recipe, userID int64) (int64, error)
	Recipes         map[int64]models.Recipes
	UsersRegistered []models.User
}

func (m *mockRepository) AddAuthToken(selector, validator string, userID int64) error {
	token := models.NewAuthToken(int64(len(m.AuthTokens)+1), selector, validator, 10000, userID)
	m.AuthTokens = append(m.AuthTokens, *token)
	return nil
}

func (m *mockRepository) AddRecipe(r *models.Recipe, userID int64) (int64, error) {
	if m.AddRecipeFunc != nil {
		return m.AddRecipeFunc(r, userID)
	}

	if m.Recipes[userID] == nil {
		m.Recipes[userID] = make(models.Recipes, 0)
	}
	m.Recipes[userID] = append(m.Recipes[userID], *r)
	return int64(len(m.Recipes)), nil
}

func (m *mockRepository) Confirm(userID int64) error {
	if !slices.ContainsFunc(m.UsersRegistered, func(user models.User) bool {
		return user.ID == userID
	}) {
		return sql.ErrNoRows
	}
	return nil
}

func (m *mockRepository) DeleteAuthToken(userID int64) error {
	index := slices.IndexFunc(m.AuthTokens, func(token models.AuthToken) bool { return token.UserID == userID })
	if index != -1 {
		m.AuthTokens = slices.Delete(m.AuthTokens, index, index+1)
	}
	return nil
}

func (m *mockRepository) GetAuthToken(_, _ string) (models.AuthToken, error) {
	return models.AuthToken{UserID: 1}, nil
}

func (m *mockRepository) IsUserExist(email string) bool {
	return slices.ContainsFunc(m.UsersRegistered, func(user models.User) bool {
		return user.Email == email
	})
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

func (m *mockRepository) UserInitials(userID int64) string {
	index := slices.IndexFunc(m.UsersRegistered, func(user models.User) bool {
		return user.ID == userID
	})
	if index == -1 {
		return ""
	}
	return string(strings.ToUpper(m.UsersRegistered[index].Email)[0])
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
		{ID: 2, Host: "afghankitchenrecipes.com", URL: "http://www.afghankitchenrecipes.com"},
	}
}

func (m *mockRepository) WebsitesSearch(query string) models.Websites {
	websites := models.Websites{
		{ID: 1, Host: "101cookbooks.com", URL: "https://101cookbooks.com"},
		{ID: 2, Host: "afghankitchenrecipes.com", URL: "http://www.afghankitchenrecipes.com"},
	}

	results := make(models.Websites, 0)
	for _, w := range websites {
		if strings.Contains(w.URL, query) {
			results = append(results, w)
		}
	}
	return results
}

type mockEmail struct {
	hitCount int64
}

func (m *mockEmail) Send(_ string, _ templates.EmailTemplate, _ any) {
	m.hitCount += 1
}

type mockFiles struct {
	extractRecipesFunc func(fileHeaders []*multipart.FileHeader) models.Recipes
}

func (m *mockFiles) ExtractRecipes(fileHeaders []*multipart.FileHeader) models.Recipes {
	if m.extractRecipesFunc != nil {
		return m.extractRecipesFunc(fileHeaders)
	}
	return models.Recipes{}
}
