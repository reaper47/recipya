package repository

import (
	"github.com/reaper47/recipya/model"
)

type MockRecipeModel struct{}

func (repo *MockRecipeModel) GetCategories() ([]string, error) {
	return []string{"appetizer", "side dish", "dessert"}, nil
}

func (repo *MockRecipeModel) GetRecipe(name string) (*model.Recipe, error) {
	return nil, nil
}

func (repo *MockRecipeModel) GetRecipes(
	category string,
	page int,
	limit int,
) ([]*model.Recipe, error) {
	hasPage := limit == -1
	if !hasPage {
		if category == "" {
			return []*model.Recipe{aRecipe()}, nil
		}
		return []*model.Recipe{otherRecipe()}, nil
	}
	return []*model.Recipe{aRecipe(), otherRecipe()}, nil
}

func (repo *MockRecipeModel) GetRecipesByCategory(c string) ([]*model.Recipe, error) {
	return []*model.Recipe{otherRecipe()}, nil
}

func (repo *MockRecipeModel) GetRecipesInfo() (*model.RecipesInfo, error) {
	return &model.RecipesInfo{Total: 11, TotalPerCategory: map[string]int64{
		"breakfast": 3,
		"lunch":     2,
		"dinner":    4,
		"dessert":   2,
	}}, nil
}

func (repo *MockRecipeModel) InsertRecipe(r *model.Recipe) (int64, error) {
	return 10, nil
}

func (repo *MockRecipeModel) ImportRecipe(url string) (*model.Recipe, error) {
	return aRecipe(), nil
}

func (repo *MockRecipeModel) SearchMaximizeFridge(
	ingredients []string,
	n int,
) ([]*model.Recipe, error) {
	return nil, nil
}

func (repo *MockRecipeModel) SearchMinimizeMissing(
	ingredients []string,
	n int,
) ([]*model.Recipe, error) {
	return nil, nil
}

func (repo *MockRecipeModel) UpdateRecipe(r *model.Recipe, id int64) error {
	return nil
}

func aRecipe() *model.Recipe {
	return &model.Recipe{
		ID:                 1,
		Name:               "carrots",
		Description:        "some delicious carrots",
		Url:                "https://www.example.com",
		PrepTime:           "PT3H30M",
		CookTime:           "PT0H30M",
		TotalTime:          "PT4H0M",
		RecipeCategory:     "side dish",
		Keywords:           "carrots,butter",
		RecipeYield:        4,
		RecipeIngredient:   []string{"1 avocado", "2 carrots"},
		RecipeInstructions: []string{"cut", "cook", "eat"},
		DateModified:       "20210820",
		DateCreated:        "20210820",
	}
}

func otherRecipe() *model.Recipe {
	return &model.Recipe{
		ID:                 2,
		Name:               "super carrots",
		Description:        "some super delicious carrots",
		Url:                "https://www.example.com",
		PrepTime:           "PT3H0M",
		CookTime:           "PT0H30M",
		TotalTime:          "PT3H30M",
		RecipeCategory:     "appetizer",
		Keywords:           "super carrots,butter",
		RecipeYield:        8,
		RecipeIngredient:   []string{"2 avocado", "10 super carrots"},
		RecipeInstructions: []string{"cut", "cook well", "eat"},
		DateModified:       "20210822",
		DateCreated:        "20210821",
	}
}
