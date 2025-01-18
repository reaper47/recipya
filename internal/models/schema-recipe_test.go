package models_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/models"

	"slices"

	"github.com/google/uuid"
)

func TestCategory_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Category: &models.Category{Value: "dinner"},
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
		CookingMethod: &models.CookingMethod{Value: "frying"},
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
		Cuisine: &models.Cuisine{Value: "French"},
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
		Description: &models.Description{Value: "ze best chicken"},
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
		Keywords: &models.Keywords{Values: "big,fat,meat"},
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
		Image: &models.Image{Value: "image1.png"},
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
		Ingredients: &models.Ingredients{Values: []string{"1 chicken", "2 eggs"}},
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
		Instructions: &models.Instructions{
			Values: []models.HowToItem{
				{Type: "HowToStep", Text: "preheat"},
				{Type: "HowToStep", Text: "mix all"},
				{Type: "HowToStep", Text: "eat everything"},
			},
		},
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

func TestThumnailURL_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		ThumbnailURL: &models.ThumbnailURL{Value: "thumbnail.png"},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "text",
			data: `{"thumbnailUrl": "thumbnail.png"}`,
		},
		{
			name: "list of strings",
			data: `{"thumbnailUrl": ["thumbnail.png","preview.png"]}`,
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
		Yield: &models.Yield{Value: 4},
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
		AtType: &models.SchemaType{Value: "Recipe"},
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
		NutritionSchema: &models.NutritionSchema{
			Calories:       "420",
			Carbohydrates:  "4",
			Cholesterol:    "3",
			Fat:            "5",
			Fiber:          "7",
			Protein:        "9",
			SaturatedFat:   "10.5",
			Servings:       "4",
			Sodium:         "350",
			Sugar:          "13",
			TransFat:       "2",
			UnsaturatedFat: "54",
		},
	}

	testcases := []struct {
		name string
		data string
	}{
		{
			name: "nutrional info without space characters",
			data: `{"nutrition":{
			  "calories": "420 kcal",
			  "carbohydrateContent": "4g",
			  "cholesterolContent": "3mg",
			  "fatContent": "5gram",
			  "fiberContent": "7grams",
			  "proteinContent": "9g",
			  "saturatedFatContent": "10.5g, saturated fat",
			  "servingSize": "makes 4 loaves",
			  "sodiumContent": "350milligram",
			  "sugarContent": "13g",
			  "transFatContent": "2g",
			  "unsaturatedFatContent": "54g"
			}}`,
		},
		{
			name: "nutrional info with abbreviated suffixes",
			data: `{"nutrition":{
			  "calories": "420 kcal",
			  "carbohydrateContent": "4 g",
			  "cholesterolContent": "3 mg",
			  "fatContent": "5 g",
			  "fiberContent": "7 g",
			  "proteinContent": "9 g",
			  "saturatedFatContent": "10.5 g",
			  "servingSize": "makes 4 loaves",
			  "sodiumContent": "350 mg",
			  "sugarContent": "13 g",
			  "transFatContent": "2 g",
			  "unsaturatedFatContent": "54 g"
			}}`,
		},
		{
			name: "nutrional info with non-abbreviated singular suffixes",
			data: `{"nutrition":{
			  "calories": "420 calories",
			  "carbohydrateContent": "4 gram",
			  "cholesterolContent": "3 milligram",
			  "fatContent": "5 gram",
			  "fiberContent": "7 gram",
			  "proteinContent": "9 gram",
			  "saturatedFatContent": "10.5 gram, saturated fat",
			  "servingSize": "makes 4 loaf",
			  "sodiumContent": "350 milligram",
			  "sugarContent": "13 gram",
			  "transFatContent": "2 gram",
			  "unsaturatedFatContent": " 54 gram"
			}}`,
		},
		{
			name: "nutrional info with non-abbreviated plural suffixes",
			data: `{"nutrition":{
			  "calories": "420 calories",
			  "carbohydrateContent": "4 grams",
			  "cholesterolContent": "3 milligrams",
			  "fatContent": "5 grams",
			  "fiberContent": "7 grams",
			  "proteinContent": "9 grams",
			  "saturatedFatContent": "10,5 grams, saturated fat",
			  "servingSize": "makes 4 loaves",
			  "sodiumContent": "350 milligrams",
			  "sugarContent": "13 grams",
			  "transFatContent": "2 grams",
			  "unsaturatedFatContent": " 54 grams"
			}}`,
		},
		{
			name: "nutrional info without suffixes as strings",
			data: `{"nutrition":{
			  "calories": "420",
			  "carbohydrateContent": "4",
			  "cholesterolContent": "3",
			  "fatContent": "5",
			  "fiberContent": "7",
			  "proteinContent": "9",
			  "saturatedFatContent": "10,5",
			  "servingSize": "4",
			  "sodiumContent": "350",
			  "sugarContent": "13",
			  "transFatContent": "2",
			  "unsaturatedFatContent": "54"
			}}`,
		},
		{
			name: "nutrional info without suffixes as integers",
			data: `{"nutrition":{
			  "calories": 420,
			  "carbohydrateContent": 4,
			  "cholesterolContent": 3,
			  "fatContent": 5,
			  "fiberContent": 7,
			  "proteinContent": 9,
			  "saturatedFatContent": 10.5,
			  "servingSize": 4,
			  "sodiumContent": 350,
			  "sugarContent": 13,
			  "transFatContent": 2,
			  "unsaturatedFatContent": 54
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
	tool := models.NewHowToTool("wok", &models.HowToItem{Quantity: 3})

	got := tool.StringQuantity()

	want := "3 wok"
	if got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestTools_UnmarshalJSON(t *testing.T) {
	want := models.RecipeSchema{
		Tools: &models.Tools{
			Values: []models.HowToItem{
				{
					Type:     "HowToTool",
					Text:     "Saw",
					Quantity: 1,
				},
				{
					Type:     "HowToTool",
					Text:     "Blender",
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
				{"@type":"HowToTool", "text":"Saw","requiredQuantity":1},
				{"@type":"HowToTool", "text":"Blender", "requiredQuantity":2}
			]}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, want)
		})
	}
}

func TestVideoObject_UnmarshalJSON(t *testing.T) {
	testcases := []struct {
		name string
		data string
		want models.RecipeSchema
	}{
		{
			name: "json",
			data: `{"video": {"@type":"VideoObject","name":"one"}}`,
			want: models.RecipeSchema{
				Video: &models.Videos{
					Values: []models.VideoObject{
						{AtType: "VideoObject", Name: "one"},
					},
				},
			},
		},
		{
			name: "json",
			data: `{"video": {"@type": "VideoObject","name": "Boeuf bourguignon met geroosterde spruiten","thumbnailUrl": ["https://allerhande.bbvms.com/mediaclip/4943112/pthumbnail/120/67.jpg","https://allerhande.bbvms.com/mediaclip/4943112/pthumbnail/900/500.jpg"],"contentUrl": "https://d1p9dpblu12ati.cloudfront.net/allerhande/media/2022/09/29/asset-4943112-1664455067445069.mp4","duration": "PT3M49S","uploadDate": "2022-10-06T22:00:00.000","interactionStatistic": {"@type": "InteractionCounter","interactionType": {"@type": "http://schema.org/WatchAction","userInteractionCount": 0}}}}`,
			want: models.RecipeSchema{
				Video: &models.Videos{
					Values: []models.VideoObject{
						{
							AtType:       "VideoObject",
							ContentURL:   "https://d1p9dpblu12ati.cloudfront.net/allerhande/media/2022/09/29/asset-4943112-1664455067445069.mp4",
							Duration:     "PT3M49S",
							Name:         "Boeuf bourguignon met geroosterde spruiten",
							ThumbnailURL: &models.ThumbnailURL{Value: "https://allerhande.bbvms.com/mediaclip/4943112/pthumbnail/120/67.jpg"},
							UploadDate:   time.Date(2022, 10, 6, 22, 0, 0, 0, time.UTC),
						},
					},
				},
			},
		},
		{
			name: "list of json",
			data: `{"video": [{"@type":"VideoObject","name":"one"},{"@type":"VideoObject","name":"two"}]}`,
			want: models.RecipeSchema{
				Video: &models.Videos{
					Values: []models.VideoObject{
						{AtType: "VideoObject", Name: "one"},
						{AtType: "VideoObject", Name: "two"},
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assertRecipeSchema(t, tc.data, tc.want)
		})
	}
}

func TestRecipeSchema_Recipe(t *testing.T) {
	imageID := uuid.New()
	tools := []models.HowToItem{
		{Name: "Saw", Quantity: 1},
		{Name: "Blender", Quantity: 2},
	}
	rs := models.RecipeSchema{
		AtContext:     "@Schema",
		AtType:        &models.SchemaType{Value: "Recipe"},
		Category:      &models.Category{Value: "lunch"},
		CookTime:      "PT3H",
		CookingMethod: nil,
		Cuisine:       &models.Cuisine{Value: "american"},
		DateCreated:   "2022-03-16",
		DateModified:  "2022-03-20",
		DatePublished: "2022-03-16",
		Description:   &models.Description{Value: "description"},
		Keywords:      &models.Keywords{Values: "kw1,kw2,kw3"},
		Image:         &models.Image{Value: imageID.String()},
		Ingredients:   &models.Ingredients{Values: []string{"ing1", "ing2", "ing3"}},
		Instructions: &models.Instructions{
			Values: []models.HowToItem{
				{Type: "HowToStep", Text: "ins1"},
				{Type: "HowToStep", Text: "ins2"},
				{Type: "HowToStep", Text: "ins3"},
			},
		},
		Name: "name",
		NutritionSchema: &models.NutritionSchema{
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
		Tools:    &models.Tools{Values: tools},
		Yield:    &models.Yield{Value: 4},
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

func TestRecipeSchema_Marshal(t *testing.T) {
	imageID := uuid.New()
	thumbnailID := uuid.New()
	tools := []models.HowToItem{
		models.NewHowToTool("Saw"),
		models.NewHowToTool("Blender"),
	}
	rs := models.RecipeSchema{
		AtContext:     "@Schema",
		AtType:        &models.SchemaType{Value: "Recipe"},
		Category:      &models.Category{Value: "lunch"},
		CookTime:      "PT3H",
		CookingMethod: nil,
		Cuisine:       &models.Cuisine{Value: "american"},
		DateCreated:   "2022-03-16",
		DateModified:  "2022-03-20",
		DatePublished: "2022-03-16",
		Description:   &models.Description{Value: "description"},
		Keywords:      &models.Keywords{Values: "kw1,kw2,kw3"},
		Image:         &models.Image{Value: imageID.String()},
		Ingredients:   &models.Ingredients{Values: []string{"ing1", "ing2", "ing3"}},
		Instructions: &models.Instructions{
			Values: []models.HowToItem{
				{Type: "HowToStep", Text: "ins1"},
				{Type: "HowToStep", Text: "ins2"},
				{Type: "HowToStep", Text: "ins3"},
			},
		},
		Name: "name",
		NutritionSchema: &models.NutritionSchema{
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
		PrepTime:     "PT1H",
		ThumbnailURL: &models.ThumbnailURL{Value: thumbnailID.String()},
		TotalTime:    "PT4H",
		Tools:        &models.Tools{Values: tools},
		Yield:        &models.Yield{Value: 4},
		URL:          "https://recipes.musicavis.ca",
	}

	xb, _ := json.Marshal(rs)
	var got map[string]any
	_ = json.Unmarshal(xb, &got)

	for k, v := range got {
		switch k {
		case "@context":
			if v != rs.AtContext {
				t.Errorf("got @context %q; want %q", v, rs.AtContext)
			}
		case "@type":
			if v.(string) != rs.AtType.Value {
				t.Errorf("got @type %q; want %q", v, rs.AtType)
			}
		case "recipeCategory":
			if v.(string) != rs.Category.Value {
				t.Errorf("got recipeCategory %q; want %q", v, rs.Category.Value)
			}
		case "cookTime":
			if v.(string) != rs.CookTime {
				t.Errorf("got recipeCategory %q; want %q", v, rs.CookTime)
			}
		case "cookingMethod":
			if v.(string) != rs.CookingMethod.Value {
				t.Errorf("got cookingMethod %q; want %q", v, rs.CookingMethod.Value)
			}
		case "recipeCuisine":
			if v.(string) != rs.Cuisine.Value {
				t.Errorf("got recipeCuisine %q; want %q", v, rs.Cuisine.Value)
			}
		case "dateCreated":
			if v.(string) != rs.DateCreated {
				t.Errorf("got dateCreated %q; want %q", v, rs.DateCreated)
			}
		case "dateModified":
			if v.(string) != rs.DateModified {
				t.Errorf("got dateModified %q; want %q", v, rs.DateModified)
			}
		case "datePublished":
			if v.(string) != rs.DatePublished {
				t.Errorf("got datePublished %q; want %q", v, rs.DatePublished)
			}
		case "description":
			if v.(string) != rs.Description.Value {
				t.Errorf("got description %q; want %q", v, rs.Description.Value)
			}
		case "keywords":
			if v.(string) != rs.Keywords.Values {
				t.Errorf("got keywords %q; want %q", v, rs.Keywords.Values)
			}
		case "image":
			if v.(string) != rs.Image.Value {
				t.Errorf("got image %q; want %q", v, rs.Image.Value)
			}
		case "recipeIngredient":
			xa := v.([]any)
			xs := make([]string, 0, len(xa))
			for _, a := range xa {
				xs = append(xs, a.(string))
			}
			if !slices.Equal(xs, rs.Ingredients.Values) {
				t.Errorf("got recipeIngredient %q; want %q", v, rs.Ingredients.Values)
			}
		case "recipeInstructions":
			xa := v.([]any)
			xv := make([]models.HowToItem, 0, len(xa))
			for _, a := range xa {
				m := a.(map[string]any)
				s, ok := m["text"]
				if ok {
					xv = append(xv, models.HowToItem{Type: "HowToStep", Text: s.(string)})
				}
			}

			if !slices.Equal(xv, rs.Instructions.Values) {
				t.Errorf("got recipeInstructions %q; want %q", xv, rs.Instructions.Values)
			}
		case "name":
			if v.(string) != rs.Name {
				t.Errorf("got name %q; want %q", v, rs.Name)
			}
		case "nutrition":
			var want map[string]any
			b, _ := json.Marshal(&rs.NutritionSchema)
			_ = json.Unmarshal(b, &want)
			if !cmp.Equal(v.(map[string]any), want) {
				t.Errorf("got nutrition %q; want %q", v, want)
			}
		case "prepTime":
			if v.(string) != rs.PrepTime {
				t.Errorf("got prepTime %q; want %q", v, rs.PrepTime)
			}
		case "thumbnailUrl":
			want := "/data/images/thumbnails/" + rs.ThumbnailURL.Value
			s := v.(string)
			_, after, ok := strings.Cut(v.(string), "/")
			if ok {
				s = "/" + after
			}
			if s != want {
				t.Errorf("got thumbnailURL %q; want %q", v, rs.ThumbnailURL.Value)
			}
		case "tool":
			xa := v.([]any)
			xv := make([]models.HowToItem, 0, len(xa))
			for _, a := range xa {
				m := a.(map[string]any)
				xv = append(xv, models.NewHowToTool(m["text"].(string), &models.HowToItem{
					Quantity: int(m["requiredQuantity"].(float64)),
				}))
			}

			if !slices.Equal(xv, rs.Tools.Values) {
				t.Errorf("got tool %+v; want %+v", xv, rs.Tools.Values)
			}
		case "totalTime":
			if v.(string) != rs.TotalTime {
				t.Errorf("got totalTime %q; want %q", v, rs.TotalTime)
			}
		case "recipeYield":
			if int16(v.(float64)) != rs.Yield.Value {
				t.Errorf("got recipeYield %q; want %q", v, rs.Yield.Value)
			}
		case "url":
			if v.(string) != rs.URL {
				t.Errorf("got url %q; want %q", v, rs.URL)
			}
		default:
			t.Errorf("invalid key %q", k)
		}
	}
}
