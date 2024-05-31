package models_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"slices"
)

func TestCategory_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Category: models.Category{Value: "dinner"},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "simple string",
			data: `{"recipeCategory": "dinner"}`,
		},
		{
			name: "list of categories",
			data: `{"recipeCategory": ["dinner","lunch"]}`,
		},
		{
			name: "map of values",
			data: `{"recipeCategory": {"Value":"dinner"}}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestCookingMethod_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		CookingMethod: models.CookingMethod{Value: "frying"},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "simple string",
			data: `{"cookingMethod": "frying"}`,
		},
		{
			name: "list of methods",
			data: `{"cookingMethod": ["frying","steaming"]}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestCuisine_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Cuisine: models.Cuisine{Value: "French"},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "simple string",
			data: `{"recipeCuisine": "French"}`,
		},
		{
			name: "list of methods",
			data: `{"recipeCuisine": ["French","English"]}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestDescription_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Description: models.Description{Value: "ze best chicken"},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "text",
			data: `{"description": "ze best chicken"}`,
		},
		{
			name: "map of values",
			data: `{"description": {"Value":"ze best chicken"}}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestKeywords_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Keywords: models.Keywords{Values: "big,fat,meat"},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "text",
			data: `{"keywords": "big,fat,meat"}`,
		},
		{
			name: "list of values",
			data: `{"keywords": ["big", "fat", "meat"]}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestImage_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Image: models.Image{Value: "image1.png"},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "text",
			data: `{"image": "image1.png"}`,
		},
		{
			name: "list of text values",
			data: `{"image": ["image1.png"]}`,
		},
		{
			name: "list of maps",
			data: `{"image": [{"url": "image1.png"}]}`,
		},
		{
			name: "map of values",
			data: `{"image": {"url": "image1.png"}}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestIngredient_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Ingredients: models.Ingredients{Values: []string{"1 chicken", "2 eggs"}},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "list of text",
			data: `{"recipeIngredient": ["1 chicken", "2 eggs"]}`,
		},
		{
			name: "map of values",
			data: `{"recipeIngredient": {"Values": ["1 chicken", "2 eggs"]}}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestInstruction_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Instructions: models.Instructions{Values: []string{"preheat", "mix all", "eat everything"}},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "text",
			data: `{"recipeInstructions": "preheat\nmix all\neat everything"}`,
		},
		{
			name: "list of text values",
			data: `{"recipeInstructions": ["preheat", "mix all", "eat everything"]}`,
		},
		{
			name: "map of values",
			data: `{"recipeInstructions": {"Values": ["preheat", "mix all", "eat everything"]}}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestYield_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Yield: models.Yield{Value: 4},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "text",
			data: `{"recipeYield": "serves 4 people"}`,
		},
		{
			name: "number",
			data: `{"recipeYield": 4}`,
		},
		{
			name: "list of number values",
			data: `{"recipeYield": [4, 5, 6]}`,
		},
		{
			name: "list of text values",
			data: `{"recipeYield": ["makes 4 loaves"]}`,
		},
		{
			name: "map of values",
			data: `{"recipeYield": {"Value": 4}}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestSchemaType_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		AtType: models.SchemaType{Value: "Recipe"},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "simple string",
			data: `{"@type": "Recipe"}`,
		},
		{
			name: "list of types",
			data: `{"@type": ["Recipe","NewsArticle"]}`,
		},
		{
			name: "map of values",
			data: `{"@type": {"Value":"Recipe"}}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestNutritionSchema_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		NutritionSchema: models.NutritionSchema{
			Calories:       "420 kcal",
			Carbohydrates:  "4g",
			Cholesterol:    "3mg",
			Fat:            "5g",
			Fiber:          "7g",
			Protein:        "9g",
			SaturatedFat:   "10g",
			Servings:       "4",
			Sodium:         "350mg",
			Sugar:          "13g",
			TransFat:       "2g",
			UnsaturatedFat: "54g",
		},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "simple string",
			data: `{"nutrition":{
			  "calories": "420 kcal",
			  "carbohydrateContent": "4g",
			  "cholesterolContent": "3mg",
			  "fatContent": "5g",
			  "fiberContent": "7g",
			  "proteinContent": "9g",
			  "saturatedFatContent": "10g",
			  "servingSize": "makes 4 loaves",
			  "sodiumContent": "350mg",
			  "sugarContent": "13g",
			  "transFatContent": "2g",
			  "unsaturatedFatContent": "54g"
			}}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func assertRecipeSchema(t testing.TB, data string, want models.RecipeSchema) {
	t.Helper()

	var got models.RecipeSchema
	err := json.Unmarshal([]byte(data), &got)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}

func TestTool_String(t *testing.T) {
	tool := models.Tool{Quantity: 3, Name: "wok"}

	got := tool.String()

	want := "3 wok"
	if got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestTools_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Tools: models.Tools{
			Values: []models.Tool{
				{
					AtType:   "HowToTool",
					Name:     "Saw",
					Quantity: 1,
				},
				{
					AtType:   "HowToTool",
					Name:     "Blender",
					Quantity: 2,
				},
			},
		},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "text",
			data: `{"tool": "1 Saw, 2 Blender"}`,
		},
		{
			name: "text",
			data: `{"tool": "1 Saw\n2 Blender"}`,
		},
		{
			name: "modified HowToTool",
			data: `{"tool": [
				{"@type":"HowToTool", "item":"Saw", "requiredQuantity":1},
				{"@type":"HowToTool", "item":"Blender", "requiredQuantity":2}
			]}`,
		},
		{
			name: "schema HowToTool",
			data: `{"tool": [
				{"@type":"HowToTool", "name":"Saw","requiredQuantity":1},
				{"@type":"HowToTool", "name":"Blender", "requiredQuantity":2}
			]}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestRecipeSchema_Recipe(t *testing.T) {
	imageID := uuid.New()
	tools := []models.Tool{
		{Name: "Saw", Quantity: 1},
		{Name: "Blender", Quantity: 2},
	}
	rs := models.RecipeSchema{
		AtContext:     "@Schema",
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: "lunch"},
		CookTime:      "PT3H",
		CookingMethod: models.CookingMethod{Value: ""},
		Cuisine:       models.Cuisine{Value: "american"},
		DateCreated:   "2022-03-16",
		DateModified:  "2022-03-20",
		DatePublished: "2022-03-16",
		Description:   models.Description{Value: "description"},
		Keywords:      models.Keywords{Values: "kw1,kw2,kw3"},
		Image:         models.Image{Value: imageID.String()},
		Ingredients:   models.Ingredients{Values: []string{"ing1", "ing2", "ing3"}},
		Instructions:  models.Instructions{Values: []string{"ins1", "ins2", "ins3"}},
		Name:          "name",
		NutritionSchema: models.NutritionSchema{
			Calories:       "341kcal",
			Carbohydrates:  "1g",
			Cholesterol:    "2g",
			Fat:            "27g",
			Fiber:          "4g",
			Protein:        "5g",
			SaturatedFat:   "6g",
			Servings:       "1",
			Sodium:         "8g",
			Sugar:          "9g",
			TransFat:       "10g",
			UnsaturatedFat: "11g",
		},
		PrepTime: "PT1H",
		Tools:    models.Tools{Values: tools},
		Yield:    models.Yield{Value: 4},
		URL:      "https://recipes.musicavis.ca",
	}

	created, _ := time.Parse(time.DateOnly, "2022-03-16")
	updated, _ := time.Parse(time.DateOnly, "2022-03-20")
	expected := models.Recipe{
		Category:     "lunch",
		CreatedAt:    created,
		Cuisine:      "american",
		Description:  "description",
		Images:       []uuid.UUID{imageID},
		Ingredients:  []string{"ing1", "ing2", "ing3"},
		Instructions: []string{"ins1", "ins2", "ins3"},
		Keywords:     []string{"kw1", "kw2", "kw3"},
		Name:         "name",
		Nutrition: models.Nutrition{
			Calories:           "341kcal",
			Cholesterol:        "2g",
			TotalFat:           "27g",
			Fiber:              "4g",
			IsPerServing:       true,
			Protein:            "5g",
			SaturatedFat:       "6g",
			Sodium:             "8g",
			Sugars:             "9g",
			TotalCarbohydrates: "1g",
			UnsaturatedFat:     "11g",
		},
		Times: models.Times{
			Prep:  1 * time.Hour,
			Cook:  3 * time.Hour,
			Total: 4 * time.Hour,
		},
		Tools:     tools,
		UpdatedAt: updated,
		URL:       "https://recipes.musicavis.ca",
		Yield:     4,
	}
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := rs.Recipe()
	if err != nil {
		t.Fatal(err)
	}
	actualBytes, err := json.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(actualBytes, expectedBytes) {
		t.Logf(cmp.Diff(*actual, expected))
		t.Fail()
	}
}
