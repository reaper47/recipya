package server_test

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/server"
	"github.com/reaper47/recipya/internal/templates"
	"golang.org/x/exp/slices"
	"io"
	"mime/multipart"
	"strings"
)

func newServerTest() *server.Server {
	repo := &mockRepository{
		AuthTokens:        make([]models.AuthToken, 0),
		RecipesRegistered: make(map[int64]models.Recipes, 0),
		ShareLinks:        make(map[int64]string, 0),
		UsersRegistered:   make([]models.User, 0),
		UsersUpdated:      make([]int64, 0),
	}
	return server.NewServer(repo, &mockEmail{}, &mockFiles{})
}

type mockRepository struct {
	AuthTokens        []models.AuthToken
	AddRecipeFunc     func(recipe *models.Recipe, userID int64) (int64, error)
	RecipeFunc        func(id, userID int64) (*models.Recipe, error)
	RecipesRegistered map[int64]models.Recipes
	ShareLinks        map[int64]string
	UsersRegistered   []models.User
	UsersUpdated      []int64
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

func (m *mockRepository) AddShareLink(link string, recipeID int64) error {
	for _, recipes := range m.RecipesRegistered {
		if slices.ContainsFunc(recipes, func(r models.Recipe) bool { return r.ID == recipeID }) {
			if _, ok := m.ShareLinks[recipeID]; ok {
				return nil
			}
			m.ShareLinks[recipeID] = link
			return nil
		}
	}
	return errors.New("recipe not found")
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

func (m *mockRepository) IsRecipeShared(id int64) bool {
	_, ok := m.ShareLinks[id]
	return ok
}

func (m *mockRepository) IsUserExist(email string) bool {
	return slices.ContainsFunc(m.UsersRegistered, func(user models.User) bool {
		return user.Email == email
	})
}

func (m *mockRepository) IsUserPassword(id int64, password string) bool {
	return slices.IndexFunc(m.UsersRegistered, func(user models.User) bool { return user.ID == id }) != -1
}

func (m *mockRepository) Recipe(id, userID int64) (*models.Recipe, error) {
	if m.RecipeFunc != nil {
		return m.RecipeFunc(id, userID)
	}

	if recipes, ok := m.RecipesRegistered[userID]; ok {
		if int64(len(recipes)) > id {
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

func (m *mockRepository) RecipeUser(recipeID int64) int64 {
	for userID, recipes := range m.RecipesRegistered {
		if i := slices.IndexFunc(recipes, func(r models.Recipe) bool { return r.ID == recipeID }); i != -1 {
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

func (m *mockRepository) UserID(email string) int64 {
	index := slices.IndexFunc(m.UsersRegistered, func(user models.User) bool {
		return user.Email == email
	})
	if index == -1 {
		return -1
	}
	return m.UsersRegistered[index].ID
}

func (m *mockRepository) UpdatePassword(userID int64, _ auth.HashedPassword) error {
	m.UsersUpdated = append(m.UsersUpdated, userID)
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

type mockEmail struct {
	hitCount int64
}

func (m *mockEmail) Send(_ string, _ templates.EmailTemplate, _ any) {
	m.hitCount += 1
}

type mockFiles struct {
	extractRecipesFunc func(fileHeaders []*multipart.FileHeader) models.Recipes
	ReadTempFileFunc   func(name string) ([]byte, error)
}

func (m *mockFiles) ExportRecipes(recipes models.Recipes) (string, error) {
	var s string
	for _, recipe := range recipes {
		s += recipe.Name + "-"
	}
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
	return uuid.New(), nil
}
