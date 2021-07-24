package model_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/reaper47/recipya/model"
)

func TestRecipe(t *testing.T) {
	t.Run("Unmarshal JSON", test_UnmarshalJSON)
	t.Run("Is created at same time", test_IsCreatedAtSameTime)
	t.Run("Recipe is modified", test_IsModified)
}

func test_UnmarshalJSON(t *testing.T) {
	aRecipeJSON := `{
		"id": 72639,
		"name": "Honey Garlic Butter Roasted Carrots",
		"description": "Honey Garlic Butter Roasted Carrots are the best side dish to add to your dinner table!",
		"url": "https:\/\/cafedelites.com\/honey-garlic-butter-roasted-carrots\/#wprm-recipe-container-51027",
		"image": "\/recipes_\/images\/Honey-Garlic-Butter-Carrots.jpg",
		"prepTime": "PT0H10M",
		"cookTime": "PT0H20M",
		"totalTime": "PT0H30M",
		"recipeCategory": "Side Dish",
		"keywords": "oven,roasted,carrots",
		"recipeYield": 4,
		"tool": ["Knife", "Peeler"],
		"recipeIngredient": [
			"2 pounds (1 kg) carrots washed and peeled (or unpeeled)", 
			"1\/3 cup butter", "3 tablespoons honey", "4 garlic cloves minced"
		],
		"recipeInstructions": [
			"Preheat oven to 425\u00b0F (220\u00b0C).",
			"Lightly grease a large baking sheet with nonstick cooking oil spray; set aside.\n", 
			"Trim ends of carrots and cut into thirds."
		],
		"nutrition": {
			"calories": "281kcal",
			"carbohydrateContent": "35g",
			"proteinContent": "2g",
			"fatContent": "15g",
			"saturatedFatContent": "9g",
			"cholesterolContent": "40mg",
			"sodiumContent": "306mg",
			"fiberContent": "6g",
			"sugarContent": "9g"
		},
		"@context": "http:\/\/schema.org",
		"@type": "model.Recipe",
		"dateModified": "2021-03-30T00:25:51+0000",
		"dateCreated": "2021-03-29T20:19:47+0000",
		"printImage": false,
		"imageUrl": "\/index.php\/apps\/cookbook\/recipes\/72639\/image?size=full"
	}`

	var recipeActual model.Recipe
	if err := json.Unmarshal([]byte(aRecipeJSON), &recipeActual); err != nil {
		t.Fatal(err)
	}

	recipeExpected := aRecipe()
	if !cmp.Equal(recipeActual, recipeExpected, cmpopts.IgnoreFields(model.Recipe{}, "ID")) {
		t.Fatal(cmp.Diff(recipeActual, recipeExpected))
	}
}

func test_IsCreatedAtSameTime(t *testing.T) {
	recipe1 := aRecipe()
	recipe2 := aRecipe()

	if !recipe1.IsCreatedSameTime(&recipe2) {
		t.Fatalf(
			"recipe1 is not created at the=same time as recipe2: %v != %v\n",
			recipe1.DateCreated,
			recipe2.DateCreated,
		)
	}

	recipe1.DateCreated += "9"
	if recipe1.IsCreatedSameTime(&recipe2) {
		t.Fatalf(
			"recipe1 should not have been created at the same time as recipe2: %v == %v\n",
			recipe1.DateCreated,
			recipe2.DateCreated,
		)
	}
}

func test_IsModified(t *testing.T) {
	recipe1 := aRecipe()
	recipe2 := aRecipe()

	if recipe1.IsModified(&recipe2) {
		t.Fatalf(
			"Recipe should not have been modified: %v == %v\n",
			recipe1.DateModified,
			recipe2.DateModified,
		)
	}

	recipe2.DateModified += "9"
	if !recipe1.IsModified(&recipe2) {
		t.Fatalf(
			"Recipe should have been modified: : %v != %v\n",
			recipe1.DateModified,
			recipe2.DateModified,
		)
	}
}

func aRecipe() model.Recipe {
	return model.Recipe{
		Name:           "Honey Garlic Butter Roasted Carrots",
		Description:    "Honey Garlic Butter Roasted Carrots are the best side dish to add to your dinner table!",
		Url:            "https://cafedelites.com/honey-garlic-butter-roasted-carrots/#wprm-recipe-container-51027",
		Image:          "/recipes_/images/Honey-Garlic-Butter-Carrots.jpg",
		PrepTime:       "PT0H10M",
		CookTime:       "PT0H20M",
		TotalTime:      "PT0H30M",
		RecipeCategory: "Side Dish",
		Keywords:       "oven,roasted,carrots",
		RecipeYield:    4,
		Tool:           []string{"Knife", "Peeler"},
		RecipeIngredient: []string{
			"2 pounds (1 kg) carrots washed and peeled (or unpeeled)",
			"1/3 cup butter", "3 tablespoons honey", "4 garlic cloves minced",
		},
		RecipeInstructions: []string{
			"Preheat oven to 425\u00b0F (220\u00b0C).",
			"Lightly grease a large baking sheet with nonstick cooking oil spray; set aside.\n",
			"Trim ends of carrots and cut into thirds.",
		},
		Nutrition: &model.NutritionSet{
			Calories:     "281kcal",
			Carbohydrate: "35g",
			Protein:      "2g",
			Fat:          "15g",
			SaturatedFat: "9g",
			Cholesterol:  "40mg",
			Sodium:       "306mg",
			Fiber:        "6g",
			Sugar:        "9g",
		},
		DateModified: "2021-03-30T00:25:51+0000",
		DateCreated:  "2021-03-29T20:19:47+0000",
	}
}
