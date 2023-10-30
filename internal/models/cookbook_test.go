package models_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"testing"
)

func TestCookbook_DominantCategories(t *testing.T) {
	testcases := []struct {
		name     string
		cookbook models.Cookbook
		num      int
		want     []string
	}{
		{
			name:     "no recipes",
			cookbook: models.Cookbook{},
			num:      5,
			want:     []string{},
		},
		{
			name: "all same categories",
			cookbook: models.Cookbook{
				Recipes: models.Recipes{
					{Category: "breakfast"},
					{Category: "breakfast"},
					{Category: "breakfast"},
					{Category: "breakfast"},
					{Category: "breakfast"},
					{Category: "breakfast"},
					{Category: "breakfast"},
					{Category: "breakfast"},
				},
			},
			num:  5,
			want: []string{"breakfast"},
		},
		{
			name: "many recipes with different categories",
			cookbook: models.Cookbook{
				Recipes: models.Recipes{
					{Category: "breakfast"},
					{Category: "breakfast"},
					{Category: "dinner"},
					{Category: "lunch"},
					{Category: "dinner"},
					{Category: "dinner"},
					{Category: "lunch"},
					{Category: "lunch"},
					{Category: "lunch"},
					{Category: "lunch"},
					{Category: "appetizer"},
					{Category: "meat"},
					{Category: "appetizer"},
					{Category: "breakfast"},
					{Category: "breakfast"},
					{Category: "meat"},
				},
			},
			num:  3,
			want: []string{"lunch", "breakfast", "dinner"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.cookbook.DominantCategories(tc.num)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("got %v but want %v", got, tc.want)
			}
		})
	}
}

func TestCookbook_MakeView(t *testing.T) {
	cookbook := models.Cookbook{
		ID:      1,
		Count:   2,
		Image:   uuid.Nil,
		Recipes: models.Recipes{{ID: 1}, {ID: 2}},
		Title:   "Lovely Ukraine",
	}

	got := cookbook.MakeView(1, 1)
	want := models.CookbookView{
		ID:          1,
		Image:       uuid.Nil,
		IsUUIDValid: false,
		NumRecipes:  2,
		Recipes:     models.Recipes{{ID: 1}, {ID: 2}},
		PageNumber:  1,
		PageItemID:  2,
		Title:       "Lovely Ukraine",
	}

	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}
