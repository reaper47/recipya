package models_test

import (
	"bytes"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/units"
	"math"
	"slices"
	"testing"
	"time"
)

func BenchmarkRecipe_ConvertMeasurementSystem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := models.Recipe{
			Description: "Preheat the oven to 177 degrees C (175 degrees C). " +
				"Stir in flour, chocolate chips, and walnuts. " +
				"Drop spoonfuls of dough 30mm apart onto ungreased baking sheets. " +
				"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
			Ingredients: []string{
				"236.58 mL butter, softened",
				"2 eggs",
				"10 mL vanilla extract",
				"5 mL baking soda",
				"709.76 mL all-purpose flour",
				"473.18 mL semisweet chocolate chips",
			},
			Instructions: []string{
				"Preheat the oven to 177 degrees C (175 degrees C).",
				"Stir in flour, chocolate chips, and walnuts.",
				"Drop spoonfuls of dough 30mm apart onto ungreased baking sheets.",
				"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
			},
		}
		converted, err := r.ConvertMeasurementSystem(units.MetricSystem)
		_ = converted
		_ = err
	}
}

func BenchmarkRecipe_Scale(b *testing.B) {
	for i := 0; i < b.N; i++ {
		recipe := models.Recipe{
			ID: 1,
			Ingredients: []string{
				"2 big apples",
				"Lots of big apples",
				"2.5 slices of bacon",
				"2 1/3 cans of bamboo sticks",
				"1½can of tomato paste",
				"6 ¾ peanut butter jars",
				"7.5mL of whiskey",
				"2 tsp lemon juice",
				"Ground ginger",
				"3 Large or 4 medium ripe Hass avocados",
				"1/4-1/2 teaspoon salt plus more for seasoning",
				"1/2 fresh pineapple, cored and cut into 1 1/2-inch pieces",
				"Un sac de chips de 1kg",
				"Two 15-ounce can Goya beans",
				"2 big apples",
				"Lots of big apples",
				"2.5 slices of bacon",
				"2 1/3 cans of bamboo sticks",
				"1½can of tomato paste",
				"6 ¾ peanut butter jars",
				"7.5mL of whiskey",
				"2 tsp lemon juice",
				"Ground ginger",
				"3 Large or 4 medium ripe Hass avocados",
				"1/4-1/2 teaspoon salt plus more for seasoning",
				"1/2 fresh pineapple, cored and cut into 1 1/2-inch pieces",
				"Un sac de chips de 1kg",
				"Two 15-ounce can Goya beans",
			},
			Instructions: nil,
			Name:         "Sauce",
			Yield:        4,
		}
		recipe.Scale(4)
	}
}

func BenchmarkNewRecipeFromTextFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte("Ostekake med røde bær og roseblader\n\nNå kan du lage en kake som har fått VM-gull. Ostekaken til Sverre Sætre var med da han sammen med kokkelandslaget vant VM-gull for noen år tilbake. Oppskriften er en forenklet utgave av Sommerkaken med friske bær og krystalliserte roseblader.\n\nMiddels\nType: Søtt\nKarakteristika: Kake\nAnledning: Fest\n\nIngredienser\n\nKjeksbunn\n100 g havrekjeks\n15 g (2 ss) hakkede pistasjnøtter\n1 ts brunt sukker\n25 g smeltet smør\n1 ts nøtteolje av hasselnøtt eller valnøtt (du kan også bruke rapsolje)\n\nOstekrem\n3 gelatinplater\n1 dl sitronsaft (her kan du også bruke limesaft, pasjonsfruktjuice eller andre syrlige juicer)\n250 g kremost naturell\n150 g sukker\nfrøene fra en vaniljestang\n3 dl kremfløte\n\nTopping\n300 g friske bringebær\nkandiserte roseblader\nurter\n\nSlik gjør du\nBruk en kakering på 22 centimeter i diameter og fire centimeter høy.\n\nKjeksbunn\nKnus kjeksene og bland med pistasjnøttene, sukker og olje. Varm smøret slik at det blir nøttebrunt på farge og bland det med kjeksblandingen til en jevn masse. Sett kakeringen på en tallerken med bakepapir eller bruk en springform med bunn. Trykk ut kjeksmassen i bunnen av kakeformen.\nTips: Kle innsiden av ringen med bakepapir. Da blir det enklere å få ut bunnen.\n\nOstekrem\nBløtlegg gelatinen i kaldt vann i 5 minutter. Kjør kremost, vaniljefrø, sukker og halvparten av juicen til en glatt masse i en matprosessor. Varm resten av juicen til kokepunktet og ta den av platen. Kryst vannet ut av den oppbløtte gelatinen og la den smelte i den varme juicen.\nTilsett den varme juicen i ostemassen, og rør den godt inn. Dette kan gjøres i matprosessoren.\nPisk fløten til krem, og vend kremen inn i ostemassen med en slikkepott. Fyll ostekrem til toppen av ringen, og stryk av med en palett slik at kaken blir helt jevn. Sett kaken i kjøleskapet til den stivner.\nFør servering: Ta kaken ut av kjøleskapet. Dekk toppen av kaken med friske bær. Pynt med sukrede roseblader og urter.\n\nTips\nOstekremen kan også fylles i små glass og serveres med bringebærsaus. Gjør man dette, bør kremen stå 2 timer i kjøleskapet slik at den stivner.\n\nKandiserte roseblader\nKandiserte blomster og blader er nydelige og godt som pynt til kaker og desserter.\nPensle rosebladene med eggehvite.\nDryss sukker på bladene. Jeg pleier å knuse sukkeret i en morter, eller kjøre det i en matprosessor slik at det blir enda finere.\nLegg til tørking over natten.\nDenne teknikken kan brukes på alt av spiselige blomster og urter, som fioler, stemorsblomster, karse, rødkløver, hvitkløver, roseblader, nellik, markjordbær- og hagejordbærblomster, ringblomst, agurkurt, svarthyll, kornblomst, løvetann, mynte,\nsitronmelisse m.m.\nNB! Blomster er stort sett ikke regnet som matvarer, derfor er det ikke tatt hensyn når det gjelder sprøyting. Hvis man kjøper blomster til dette formål, må man altså passe på at de ikke er sprøytet.\n\nhttps://www.nrk.no/mat/ostekake-med-rode-baer-og-roseblader-1.8229671"))
		_, _ = models.NewRecipeFromTextFile(buf)
	}
}

func TestRecipe_ConvertMeasurementSystem(t *testing.T) {
	testcases := []struct {
		name string
		in   *models.Recipe
		to   units.System
		want error
	}{
		{
			name: "imperial to imperial",
			in: &models.Recipe{
				Ingredients: []string{
					"2 eggs",
					"1 cup butter, softened",
				},
			},
			to:   units.ImperialSystem,
			want: errors.New("system already imperial"),
		},
		{
			name: "metric to metric",
			in: &models.Recipe{
				Ingredients: []string{
					"2 eggs",
					"1.5 mL butter, softened",
				},
			},
			to:   units.MetricSystem,
			want: errors.New("system already metric"),
		},
	}
	for _, tc := range testcases {
		t.Run("cannot convert "+tc.name, func(t *testing.T) {
			_, err := tc.in.ConvertMeasurementSystem(tc.to)
			if err == nil {
				t.Fatalf("expected error but got %q", err)
			}
		})
	}

	testcases2 := []struct {
		name string
		in   models.Recipe
		to   units.System
		want models.Recipe
	}{
		{
			name: "imperial to metric",
			in: models.Recipe{
				Description: "Preheat the oven to 351 °F (351 °F). " +
					"Stir in flour, chocolate chips, and walnuts. " +
					"Drop spoonfuls of dough 1.18 inches apart onto ungreased baking sheets. " +
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				Ingredients: []string{
					"1 cup butter, softened",
					"2 eggs",
					"2 teaspoons vanilla extract",
					"1 teaspoon baking soda",
					"3 cups all-purpose flour",
					"2 cups semisweet chocolate chips",
					"Salt and pepper",
				},
				Instructions: []string{
					"Preheat the oven to 350 degrees F (175 degrees C).",
					"Stir in flour, chocolate chips, and walnuts.",
					"Drop spoonfuls of dough 2 inches apart onto ungreased baking sheets.",
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				},
			},
			to: units.MetricSystem,
			want: models.Recipe{
				Description: "Preheat the oven to 177 °C (177 °C). " +
					"Stir in flour, chocolate chips, and walnuts. " +
					"Drop spoonfuls of dough 3 cm apart onto ungreased baking sheets. " +
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				Ingredients: []string{
					"2.37 dl butter, softened",
					"2 eggs",
					"10 ml vanilla extract",
					"5 ml baking soda",
					"7.1 dl all-purpose flour",
					"4.73 dl semisweet chocolate chips",
					"Salt and pepper",
				},
				Instructions: []string{
					"Preheat the oven to 177 °C (177 °C).",
					"Stir in flour, chocolate chips, and walnuts.",
					"Drop spoonfuls of dough 5.08 cm apart onto ungreased baking sheets.",
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				},
			},
		},
		{
			name: "metric to imperial",
			in: models.Recipe{
				Description: "Preheat the oven to 177 degrees C (175 degrees C). " +
					"Stir in flour, chocolate chips, and walnuts. " +
					"Drop spoonfuls of dough 30mm apart onto ungreased baking sheets. " +
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				Ingredients: []string{
					"236.58 mL butter, softened",
					"2 eggs",
					"10 mL vanilla extract",
					"5 mL baking soda",
					"709.76 mL all-purpose flour",
					"473.18 mL semisweet chocolate chips",
				},
				Instructions: []string{
					"Preheat the oven to 177 degrees C (175 degrees C).",
					"Stir in flour, chocolate chips, and walnuts.",
					"Drop spoonfuls of dough 30mm apart onto ungreased baking sheets.",
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				},
			},
			to: units.ImperialSystem,
			want: models.Recipe{
				Description: "Preheat the oven to 351 °F (351 °F). " +
					"Stir in flour, chocolate chips, and walnuts. " +
					"Drop spoonfuls of dough 1.18 inches apart onto ungreased baking sheets. " +
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				Ingredients: []string{
					"1 cup butter, softened",
					"2 eggs",
					"2 tsp vanilla extract",
					"1 tsp baking soda",
					"3 cups all-purpose flour",
					"2 cups semisweet chocolate chips",
				},
				Instructions: []string{
					"Preheat the oven to 351 °F (351 °F).",
					"Stir in flour, chocolate chips, and walnuts.",
					"Drop spoonfuls of dough 1.18 inches apart onto ungreased baking sheets.",
					"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
				},
			},
		},
	}
	for _, tc := range testcases2 {
		t.Run("valid "+tc.name, func(t *testing.T) {
			got, _ := tc.in.ConvertMeasurementSystem(tc.to)

			if got.Description != tc.want.Description {
				t.Fatalf("got description:\n%s\nbut want:\n%s", got.Description, tc.want.Description)
			}

			if len(got.Ingredients) != len(tc.want.Ingredients) {
				t.Fatalf("ingredients of different lengths: %#v but want %#v", got.Ingredients, tc.want.Ingredients)
			}
			if len(got.Instructions) != len(tc.want.Instructions) {
				t.Fatalf("instructions of different lengths: %#v but want %#v", got.Ingredients, tc.want.Ingredients)
			}
			for i, s := range got.Ingredients {
				if s != tc.want.Ingredients[i] {
					t.Errorf("got ingredient %q but want %q", s, tc.want.Ingredients[i])
				}
			}
			for i, s := range got.Instructions {
				if s != tc.want.Instructions[i] {
					t.Errorf("got instruction %q but want %q", s, tc.want.Instructions[i])
				}
			}
		})
	}
}

func TestNewSearchOptionsRecipe(t *testing.T) {
	testcases := []struct {
		name   string
		method string
		want   models.SearchOptionsRecipes
	}{
		{
			name:   "empty defaults to name",
			method: "",
			want:   models.SearchOptionsRecipes{IsByName: true, Page: 1, Sort: models.Sort{IsDefault: true}},
		},
		{
			name:   "name",
			method: "name",
			want:   models.SearchOptionsRecipes{IsByName: true, Page: 1, Sort: models.Sort{IsDefault: true}},
		},
		{
			name:   "empty defaults to name",
			method: "full",
			want:   models.SearchOptionsRecipes{IsFullSearch: true, Page: 1, Sort: models.Sort{IsDefault: true}},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			compare(t, models.NewSearchOptionsRecipe(tc.method, "", 1), tc.want)
		})
	}
}

func TestRecipe_IsEmpty(t *testing.T) {
	t.Run("is empty", func(t *testing.T) {
		var r models.Recipe
		if !r.IsEmpty() {
			t.Fail()
		}
	})

	testcases := []struct {
		name   string
		recipe models.Recipe
	}{
		{name: "has category only", recipe: models.Recipe{Category: "breakfast"}},
		{name: "has created at only", recipe: models.Recipe{CreatedAt: time.Now()}},
		{name: "has cuisine only", recipe: models.Recipe{Cuisine: "american"}},
		{name: "has description only", recipe: models.Recipe{Description: "american"}},
		{name: "has ID only", recipe: models.Recipe{ID: 1}},
		{name: "has image only", recipe: models.Recipe{Image: uuid.New()}},
		{name: "has Ingredients only", recipe: models.Recipe{Ingredients: []string{"one"}}},
		{name: "has Instructions only", recipe: models.Recipe{Instructions: []string{"one"}}},
		{name: "has Keywords only", recipe: models.Recipe{Keywords: []string{"one"}}},
		{name: "has Name only", recipe: models.Recipe{Name: "one"}},
		{name: "has Nutrition only", recipe: models.Recipe{Nutrition: models.Nutrition{Calories: "666 kcal"}}},
		{name: "has times only", recipe: models.Recipe{Times: models.Times{Prep: 5 * time.Hour}}},
		{name: "has Tools only", recipe: models.Recipe{Tools: []string{"hose"}}},
		{name: "has UpdatedAt only", recipe: models.Recipe{UpdatedAt: time.Now()}},
		{name: "has URL only", recipe: models.Recipe{URL: "mama"}},
		{name: "has Yield only", recipe: models.Recipe{Yield: 5}},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.recipe.IsEmpty() {
				t.Fail()
			}
		})
	}
}

func TestNutrientFDC_ValuePerGram(t *testing.T) {
	testcases := []struct {
		name     string
		nutrient models.NutrientFDC
		want     float64
	}{
		{
			name: "ug",
			nutrient: models.NutrientFDC{
				Amount:    777,
				UnitName:  "UG",
				Reference: units.Measurement{Quantity: 3, Unit: units.Teaspoon},
			},
			want: 0.00011655,
		},
		{
			name: "mg",
			nutrient: models.NutrientFDC{
				Amount:    78,
				UnitName:  "MG",
				Reference: units.Measurement{Quantity: 2, Unit: units.Cup},
			},
			want: 0.369,
		},
		{
			name: "g",
			nutrient: models.NutrientFDC{
				Amount:    128,
				UnitName:  "G",
				Reference: units.Measurement{Quantity: 1, Unit: units.Pound},
			},
			want: 580.60,
		},
		{
			name: "kg",
			nutrient: models.NutrientFDC{
				Amount:    23,
				UnitName:  "KG",
				Reference: units.Measurement{Quantity: 27, Unit: units.Tablespoon},
			},
			want: 91825.8417,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.nutrient.Value()
			if math.Abs(got-tc.want) > 1e-2 {
				t.Fatalf("got %g but want %g", got, tc.want)
			}
		})
	}
}

func TestNutrientsFDC_CalculateNutrition(t *testing.T) {
	nutrients := models.NutrientsFDC{
		{Name: "Carbohydrates", Amount: 61, UnitName: "G", Reference: units.Measurement{Quantity: 1, Unit: units.Pound}},
		{Name: "Energy", Amount: 478, UnitName: "KCAL", Reference: units.Measurement{Quantity: 100, Unit: units.Gram}},
		{Name: "Carbohydrates", Amount: 478, UnitName: "MG", Reference: units.Measurement{Quantity: 0.5, Unit: units.Teaspoon}},
		{Name: "Cholesterol", Amount: 1.2, UnitName: "MG", Reference: units.Measurement{Quantity: 0.5, Unit: units.Teaspoon}},
		{Name: "Cholesterol", Amount: 1.2, UnitName: "MG", Reference: units.Measurement{Quantity: 1, Unit: units.Tablespoon}},
		{Name: "Fiber, total dietary", Amount: 3, UnitName: "G", Reference: units.Measurement{Quantity: 1, Unit: units.Cup}},
		{Name: "Protein", Amount: 12, UnitName: "G", Reference: units.Measurement{Quantity: 6, Unit: units.Cup}},
		{Name: "Fatty acids, total monounsaturated", Amount: 0.6, UnitName: "MG", Reference: units.Measurement{Quantity: 2, Unit: units.Pound}},
		{Name: "Fatty acids, total polyunsaturated", Amount: 1.2, UnitName: "MG", Reference: units.Measurement{Quantity: 5, Unit: units.Pound}},
		{Name: "Fatty acids, total trans", Amount: 56, UnitName: "UG", Reference: units.Measurement{Quantity: 5, Unit: units.Tablespoon}},
		{Name: "Fatty acids, total saturated", Amount: 128, UnitName: "UG", Reference: units.Measurement{Quantity: 12, Unit: units.Tablespoon}},
		{Name: "Sodium, Na", Amount: 1286, UnitName: "MG", Reference: units.Measurement{Quantity: 3, Unit: units.Cup}},
		{Name: "Sugars, total including NLEA", Amount: 90, UnitName: "G", Reference: units.Measurement{Quantity: 7, Unit: units.Cup}},
	}

	got := nutrients.NutritionFact(100)
	want := models.Nutrition{
		Calories:           "478 kcal",
		Cholesterol:        "207.44 ug",
		Fiber:              "7.10 g",
		Protein:            "170.34 g",
		SaturatedFat:       "227.12 ug",
		Sodium:             "9.13 g",
		Sugars:             "1.49 kg",
		TotalCarbohydrates: "276.70 g",
		TotalFat:           "32.93 mg",
		UnsaturatedFat:     "32.66 mg",
	}

	if !got.Equal(want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}

func TestRecipe_Copy(t *testing.T) {
	times, _ := models.NewTimes("PT1H20M", "PT30M")
	original := models.Recipe{
		Category:    "breakfast",
		CreatedAt:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Cuisine:     "American",
		Description: "The best American breakfast you could ever have.",
		ID:          1,
		Image:       uuid.Nil,
		Ingredients: []string{
			"2 eggs",
			"3 cups maple syrup",
			"1/2 tbsp cinnamon",
			"1 toast spread with butter",
		},
		Instructions: []string{
			"Cook the meat on high",
			"Mix the eggs together, whisk and cook on low",
			"Pour maple syrup over meat",
			"Eat",
		},
		Keywords: []string{"eggs", "meat", "juicy"},
		Name:     "Chicken Sauce",
		Nutrition: models.Nutrition{
			Calories:           "1g",
			Cholesterol:        "2g",
			Fiber:              "3g",
			Protein:            "4g",
			SaturatedFat:       "5g",
			Sodium:             "6g",
			Sugars:             "7g",
			TotalCarbohydrates: "8g",
			TotalFat:           "9g",
			UnsaturatedFat:     "10g",
		},
		Times: times,
		Tools: []string{
			"small pan",
			"large pan",
			"spatula",
		},
		UpdatedAt: time.Date(2012, 12, 25, 0, 0, 0, 0, time.UTC),
		URL:       "https://www.example.com",
		Yield:     4,
	}

	copied := original.Copy()

	if !cmp.Equal(original, copied) {
		t.Log(cmp.Diff(original, copied))
		t.Fail()
	}
	copied.Ingredients[0] = "pig"
	if slices.Equal(original.Ingredients, copied.Ingredients) {
		t.Fatal("ingredients slices the same when they should not")
	}
	copied.Instructions[0] = "jesus"
	if slices.Equal(original.Instructions, copied.Instructions) {
		t.Fatal("instructions slices the same when they should not")
	}
	copied.Keywords[0] = "european"
	if slices.Equal(original.Keywords, copied.Keywords) {
		t.Fatal("keywords slices the same when they should not")
	}
	copied.Tools[0] = "saw"
	if slices.Equal(original.Tools, copied.Tools) {
		t.Fatal("tools slices the same when they should not")
	}
}

func TestRecipe_Normalize(t *testing.T) {
	r := models.Recipe{
		Description: "Place the chicken pieces on a baking sheet and bake 1l 1 l 1ml 1 ml until they 425°f (220°c) and golden.",
		Ingredients: []string{"ing1 1L", "1 L ing2", "ing3 of 1mL stuff", "ing4 of stuff 1 mL"},
		Instructions: []string{
			"ins1 1l",
			"1 l ins2",
			"ins3 of 1ml stuff",
			"ins4 of stuff 1 ml",
		},
	}

	r.Normalize()

	expectedDescription := "Place the chicken pieces on a baking sheet and bake 1L 1 L 1mL 1 mL until they 425°F (220°C) and golden."
	if r.Description != expectedDescription {
		t.Errorf("expected the description to be normalized")
	}

	expected := []string{
		"ing1 1L",
		"1 L ing2",
		"ing3 of 1mL stuff",
		"ing4 of stuff 1 mL",
	}
	for i, v := range r.Ingredients {
		if v != expected[i] {
			t.Errorf("expected '%v' but got '%v'", expected[i], v)
		}
	}

	expected = []string{"ins1 1L", "1 L ins2", "ins3 of 1mL stuff", "ins4 of stuff 1 mL"}
	for i, v := range r.Instructions {
		if v != expected[i] {
			t.Errorf("expected '%v' but got '%v'", expected[i], v)
		}
	}
}

func TestRecipe_Scale(t *testing.T) {
	recipe := models.Recipe{
		ID: 1,
		Ingredients: []string{
			"2 big apples",
			"Lots of big apples",
			"2.5 slices of bacon",
			"2 1/3 cans of bamboo sticks",
			"1½can of tomato paste",
			"6 ¾ peanut butter jars",
			"7.5mL of whiskey",
			"2 tsp lemon juice",
			"Ground ginger",
			"3 Large or 4 medium ripe Hass avocados",
			"1/4-1/2 teaspoon salt plus more for seasoning",
			"1/2 fresh pineapple, cored and cut into 1 1/2-inch pieces",
			"Un sac de chips de 1kg",
			"Two 15-ounce can Goya beans",
			"1 c. soupe bovril boeuf",
		},
		Instructions: nil,
		Name:         "Sauce",
		Yield:        4,
	}

	t.Run("double recipe", func(t *testing.T) {
		got := recipe.Copy()
		got.Scale(8)

		want := recipe.Copy()
		want.Ingredients = []string{
			"4 big apples",
			"Lots of big apples",
			"5 slices of bacon",
			"4 2/3 cans of bamboo sticks",
			"3 can of tomato paste",
			"13 1/2 peanut butter jars",
			"15 mL of whiskey",
			"1 1/3 tbsp lemon juice",
			"Ground ginger",
			"6 Large or 8 medium ripe Hass avocados",
			"1/2 tsp salt plus more for seasoning",
			"1 fresh pineapple, cored and cut into 3-inch pieces",
			"Un sac de chips de 1kg",
			"4 15-ounce can Goya beans",
			"2 cups. soupe bovril boeuf",
		}
		want.Yield = 8
		assertStructsEqual(t, got, want)
	})

	t.Run("decrease recipe by 1.5x", func(t *testing.T) {
		got := recipe.Copy()
		got.Scale(1)

		want := recipe.Copy()
		want.Ingredients = []string{
			"1/2 big apples",
			"Lots of big apples",
			"5/8 slices of bacon",
			"0.583 cans of bamboo sticks",
			"3/8 can of tomato paste",
			"1.687 peanut butter jars",
			"1.880 mL of whiskey",
			"1/2 tsp lemon juice",
			"Ground ginger",
			"3/4 Large or 1 medium ripe Hass avocados",
			"0.060 tsp salt plus more for seasoning",
			"1/8 fresh pineapple, cored and cut into 3/8-inch pieces",
			"Un sac de chips de 1kg",
			"0.5 15-ounce can Goya beans",
			"4 tbsp. soupe bovril boeuf",
		}
		want.Yield = 1
		assertStructsEqual(t, got, want)
	})
}

func TestRecipe_Schema(t *testing.T) {
	imageUUID := uuid.New()
	r := models.Recipe{
		Category:     "lunch",
		CreatedAt:    time.Now(),
		Cuisine:      "american",
		Description:  "description",
		ID:           1,
		Image:        imageUUID,
		Ingredients:  []string{"ing1", "ing2", "ing3"},
		Instructions: []string{"ins1", "ins2", "ins3"},
		Keywords:     []string{"kw1", "kw2", "kw3"},
		Name:         "name",
		Nutrition: models.Nutrition{
			Calories:           "1kcal",
			TotalCarbohydrates: "1g",
			Sugars:             "2g",
			Protein:            "3g",
			TotalFat:           "4g",
			SaturatedFat:       "5g",
			Cholesterol:        "6g",
			Sodium:             "7g",
			Fiber:              "8g",
		},
		Times: models.Times{
			Prep:  1 * time.Hour,
			Cook:  2 * time.Hour,
			Total: 3 * time.Hour,
		},
		Tools:     []string{"t1", "t2", "t3"},
		UpdatedAt: time.Now().Add(2 * time.Hour),
		URL:       "https://www.google.com",
		Yield:     4,
	}

	schema := r.Schema()

	if schema.AtContext != "https://schema.org" {
		t.Errorf("context must be https://schema.org")
	}
	if schema.AtType.Value != "Recipe" {
		t.Errorf("type must be Recipe")
	}
	if schema.Category.Value != "lunch" {
		t.Errorf("wanted category 'lunch' but got '%q'", schema.Category)
	}
	if schema.CookTime != "PT2H0M0S" {
		t.Errorf("wanted cook time PT1H but got %q", schema.CookTime)
	}
	if schema.CookingMethod.Value != "" {
		t.Errorf("wanted an empty cooking method but got %q", schema.CookingMethod)
	}
	if schema.Cuisine.Value != "american" {
		t.Errorf("wanted an empty cusine but got %q", schema.Cuisine)
	}
	v := r.CreatedAt.Format(time.DateOnly)
	if schema.DateCreated != v {
		t.Errorf("wanted date created %q but got %q", v, schema.DateCreated)
	}
	v = r.UpdatedAt.Format(time.DateOnly)
	if schema.DateModified != v {
		t.Errorf("wanted date modified %q but got %q", v, schema.DateModified)
	}
	v = r.CreatedAt.Format(time.DateOnly)
	if schema.DatePublished != v {
		t.Errorf("wanted date published %q but got %q", v, schema.DatePublished)
	}
	if schema.Description.Value != "description" {
		t.Errorf("wanted description 'description' but got %q", schema.Description)
	}
	if schema.Keywords.Values != "kw1,kw2,kw3" {
		t.Errorf("wanted keywords 'kw1,kw2,kw3' but got %q", schema.Keywords)
	}
	v = imageUUID.String()
	if schema.Image.Value != v {
		t.Errorf("wanted uuid %q but got %q", v, schema.Image)
	}

	if !slices.Equal(schema.Ingredients.Values, []string{"ing1", "ing2", "ing3"}) {
		t.Errorf("wanted ingredients []string{ing1, ing2, ing3} but got %v", schema.Ingredients)
	}
	if !slices.Equal(schema.Instructions.Values, []string{"ins1", "ins2", "ins3"}) {
		t.Errorf(
			"wanted instructions []string{ins1, ins2, ins3} but got %v",
			schema.Instructions.Values,
		)
	}
	if schema.Name != "name" {
		t.Errorf("wanted name 'name' but got %q", schema.Name)
	}
	if schema.NutritionSchema != r.Nutrition.Schema("4") {
		t.Errorf("wanted nutrition but got %v", schema.NutritionSchema)
	}
	if schema.PrepTime != "PT1H0M0S" {
		t.Errorf("wanted prepTime PT1H0M0S but got %q", schema.PrepTime)
	}
	if !slices.Equal(schema.Tools.Values, []string{"t1", "t2", "t3"}) {
		t.Errorf("wanted tools []string{t1, t2, t3} but got %v", schema.Tools.Values)
	}
	if schema.Yield.Value != 4 {
		t.Errorf("wanted yield 4 but got %d", schema.Yield.Value)
	}
	if schema.URL != "https://www.google.com" {
		t.Errorf("wanted url https://www.google.com but got %q", schema.URL)
	}
}

func TestNewTimes(t *testing.T) {
	actual, err := models.NewTimes("PT1H0M0S", "PT2H0M0S")
	assertNoError(t, err)
	if err != nil {
		t.Fatal(err)
	}

	if actual.Cook != 2*time.Hour {
		t.Errorf("wanted cooking time 2H but got %v", actual.Cook.String())
	}
	if actual.Prep != 1*time.Hour {
		t.Errorf("wanted prep time 1H but got %v", actual.Prep.String())
	}
	if actual.Total != 3*time.Hour {
		t.Errorf("wanted total time 3H but got %v", actual.Total.String())
	}
}

func TestNewRecipesFromMasterCook(t *testing.T) {
	testcases := []struct {
		name string
		buf  *bytes.Buffer
		want models.Recipes
	}{
		{
			name: "standard mxp file structure",
			buf:  bytes.NewBuffer([]byte("                     *  Exported from  MasterCook  *\n\n                                !Read Me!\n\nRecipe By     : Bill Wight\nServing Size  : 1    Preparation Time :0:00\nCategories    : *Afghan                          Info\n\n  Amount  Measure       Ingredient -- Preparation Method\n--------  ------------  --------------------------------\n                        ***Information***\n\nThe recipes in this collection were collected, corrected, edited and categorized by Bill Wight.\n\nYou may redistribute this collection in any form as long as this message is included.\n\nFor comments or suggestions, email me at:   Bill_Wight_CA@yahoo.com\n\nSome answers to frequently asked questions about my recipe collection:\n\nI notice that you have edited and corrected many of the recipe archives you have on your page.  What did you edit out and correct?\n\nMany of the recipes in this collection were typed, formatted and posted in the late 1980's and early 1990's to various bulletin boards before the Internet came into being.  These recipes were often formatted for a recipe management program called MealMaster.  During the years, the recipes have been posted and reposted to various mailing lists and newsgroups on the Internet.  They often pick up errors during reformatting and reposting.  I have also converted many of these recipes from Web page format and MealMaster format to MasterCook format.  This conversion process often introduces errors into the recipe.  I have attempted to correct the conversion errors.  I have also removed the names of the people who typed, formatted and posted these recipes.  I have kept the recipe author or creator where known.  I have also removed duplicate recipes and I have attempted to make some sense out of the recipes categories. \n\nI notice that you have removed all the nutrition data from the recipes in your collection.  Why did you do this?\n\nSince many of the recipes in this archive were typed and formatted several years ago and they often contain conflicting nutrition information, I have decided to strip out the nutritional data from all the recipes I edit.  If you have a current version of MasterCook (6.0) it will calculate nutritional data more accurately than the older versions.  So after you import the recipes into your copy of MasterCook, you will have up-to-date nutritional data available.  If nutritional data is important to you, make sure that the serving size looks reasonable.  Many recipes have a serving size set to the MasterCook default setting of '1'.  So nutritional data would assume the entire recipe is one serving and report the data incorrectly.  I have tired, as I edited the recipes, to correct the serving size if it was set to \"1\" to a reasonable serving size based on the recipes ingredient amounts.\n\n                   - - - - - - - - - - - - - - - - - - \n\n\n                     *  Exported from  MasterCook  *\n\n                      Abraysham Kabaub (Silk Kebab)\n\nRecipe By     : \nServing Size  : 30   Preparation Time :0:00\nCategories    : *Afghan                          Desserts\n\n  Amount  Measure       Ingredient -- Preparation Method\n--------  ------------  --------------------------------\n                        ***SYRUP***\n   1 1/2  cups          granulated sugar\n   1      teaspoon      lemon juice\n   1      cup           water\n     1/4  teaspoon      saffron threads -- (optional)\n                        ***OMELET***\n   8                    eggs\n   1      pinch         salt\n                        ***TO FINISH***\n   2      cups          oil\n     1/2  teaspoon      ground cardamom\n     3/4  cup           finely chopped pistachios *\n\n*Note: Instead of pistachio nuts, walnuts may be used if desired. \n\nDissolve sugar in water in heavy pan over medium heat. Bring to the boil, add lemon juice and saffron and boil for 10 minutes. Cool and strain into a 25 cm (10 inch) pie plate. Keep aside. Break eggs into a casserole dish about 20 cm (8 inches) in diameter. \n\nThe size and flat base are important. Add salt and mix eggs with fork until yolks and whites are thoroughly combined - do not beat as eggs must not be foamy. Heat oil in an electric frying pan to 190 C (375 F) or in a 25 cm (10 inch) frying pan placed on a thermostatically controlled hot plate or burner. Have ready nearby a long skewer, the plate of syrup, a baking sheet and the nuts mixed with the cardamom. \n\nA bowl of water and a cloth for drying hands are also necessary. Hold dish with eggs in one hand next to the pan of oil and slightly above it. Put hand into egg, palm down, so that egg covers back of hand. Lift out hand, curling fingers slightly inwards, then open out over hot oil, fingers pointing down. \n\nMove hand across surface of oil so that egg falls in streams from fingertips. Dip hand in egg again and make more strands across those already in pan. Repeat 3 to 4 times until about an eighth of the egg is used. There should be a closely meshed layer of egg strands about 20 cm (8 inches) across. Work quickly so that the last lot of egg is added not long after the first lot. Rinse hands quickly and dry. Take skewer and slide under bubbling omelet, lift up and turn over to lightly brown other side. \n\nThe first side will be bubbly, the underside somewhat smoother. When golden brown lift out with skewer and drain over pan. Place omelet flat in the syrup, spoon syrup over the top and lift out with skewer onto baking sheet. Roll up with bubbly side inwards. \n\nFinish roll should be about 3 cm (1 1/4 inches) in diameter. Put to one side and sprinkle with nuts. Repeat with remaining egg, making 7 to 8 rolls in all. Though depth of egg diminishes, you will become so adept that somehow you will get it in the pan in fine strands. \n\nWhen cool, cut kabaubs into 4-5 cm (1 1/2 to 2 inch pieces and serve. These keep well in a sealed container in a cool place.\n\n\n\n                   - - - - - - - - - - - - - - - - - - \n\n\n                     *  Exported from  MasterCook  *\n\n                              Afghan Chicken\n\nRecipe By     : San Francisco Examiner, 6/2/93.\nServing Size  : 6    Preparation Time :0:00\nCategories    : Middle East                      Chicken\n                *Afghan                          On-The-Grill\n\n  Amount  Measure       Ingredient -- Preparation Method\n--------  ------------  --------------------------------\n   2      large   clov  garlic\n     1/2  teaspoon      salt\n   2      cups          plain whole-milk yogurt\n   4      tablespoons   juice and pulp of 1 large lemon\n     1/2  teaspoon      cracked black pepper\n   2      large   whol  chicken breasts -- about 2 pounds\n\nLong, slow marinating in garlicky yogurt tenderizes, moistens and adds deep flavor, so you end up with skinless grilled chicken that's as delicious as it is nutritionally correct. Serve with soft pita or Arab flatbread and fresh yogurt.\n\nPut the salt in a wide, shallow non-reactive bowl with the garlic and mash them together until you have paste. Add yogurt, lemon and pepper.\n\nSkin the chicken breasts, remove all visible fat and separate the halves. Bend each backward to break the bones so the pieces win lie flat. Add to the yogurt and turn so all surfaces are well-coated.\n\nCover the bowl tightly and refrigerate. Allow to marinate at least overnight, up to a day and a half. Turn when you think of it.\n\nTo cook, remove breasts from marinade and wipe off all but a thin film. Broil or grill about 6 inches from the heat for 6 to 8 minutes a side, or until thoroughly cooked. Meat will brown somewhat but should not char. Serve at once.\n\n\n\n\n\n                   - - - - - - - - - - - - - - - - - - \n\n\n                     *  Exported from  MasterCook  *\n\n                          Afghan Chicken Kebobs\n\nRecipe By     : \nServing Size  : 4    Preparation Time :0:00\nCategories    : Chicken                          *Afghan\n                On-The-Grill\n\n  Amount  Measure       Ingredient -- Preparation Method\n--------  ------------  --------------------------------\n   1      cup           yogurt\n   1 1/2  teaspoons     salt\n     1/2  teaspoon      ground red or black pepper\n   3      centiliters   garlic -- finely minced\n   1 1/2  pounds        chicken breasts -- boneless,\n                        skinless -- cut into kebob\n                         -- ¥\n                        flatbread such as lavash\n                        pita or flour tortillas\n   3                    tomatoes -- sliced\n   2                    onions -- sliced\n                        cilantro to taste\n   2                    lemons or 4 limes -- quartered\n\n1. Mix yogurt, salt, pepper and garlic in a bowl. Mix chicken with yogurt and marinate 1 to 2 hours at room temperature, up to 2 days refrigerated.\n\n2. Thread chicken on skewers and grill over medium hot coals.\n\n3. Place warmed pita bread on plates (if using tortillas, toast briefly over flame), divide meat among them, top with tomato and onion slices and cilantro and fold bread over. Serve with lemon or lime quarters for squeezing.\n\n\n\n                   - - - - - - - - - - - - - - - - - - \n\n\n                     *  Exported from  MasterCook  *\n\n              Afghan Pumpkins Kadu Bouranee (Sweet Pumpkin)\n\nRecipe By     : \nServing Size  : 1    Preparation Time :0:00\nCategories    : *Afghan                          Vegetables\n\n  Amount  Measure       Ingredient -- Preparation Method\n--------  ------------  --------------------------------\n   2      pounds        fresh pumpkin or squash\n     1/4  cup           corn oil\n                        ***SWEET TOMATO SAUCE***\n   1      teaspoon      crushed garlic\n   1      cup           water\n     1/2  teaspoon      salt\n     1/2  cup           sugar\n   4      ounces        tomato sauce\n     1/2  teaspoon      ginger root -- chopped fine\n   1      teaspoon      freshly ground coriander\n                        seeds\n     1/4  teaspoon      black pepper\n                        ***YOGURT SAUCE***\n     1/4  teaspoon      crushed garlic\n     1/4  teaspoon      salt\n     3/4  cup           plain yogurt\n                        ***GARNISH***\n                        dry mint leaves -- crushed\n\nPeel the pumpkin and cut into 2-3\" cubes; set aside. Heat oil in a large frying pan that has a lid. Fry the pumpkins on both sides for a couple of minutes until lightly browned. Mix together ingredients for Sweet Tomato Sauce in a bowl then add to pumpkin mixture in fry pan. \n\nCover and cook 20-25 minutes over low heat until the pumpkin is cooked and most of the liquid has evaporated. (I don't know how it's going to evaporate if the pan is covered....-B.) \n\nMix together the ingredients for the yogurt sauce. To serve: Spread half the yogurt sauce on a plate and lay the pumpkin on top. Top with remaining yogurt and any cooking juices left over. Sprinkle with dry mint. \n\nMay be served with chalow (basmati rice) and naan or pita bread.\n\nFrom Afghani Cooking, the cookbook from Da Afghan Restaurant, Bloomington/Minneapolis, MN.\n\n\n                   - - - - - - - - - - - - - - - - - - \n")),
			want: models.Recipes{
				{
					Category: "*Afghan",
					Ingredients: []string{
						"***SYRUP***",
						"1 1/2 cups granulated sugar",
						"1 teaspoon lemon juice",
						"1 cup water",
						"1/4 teaspoon saffron threads -- (optional)",
						"***OMELET***",
						"8 eggs",
						"1 pinch salt",
						"***TO FINISH***",
						"2 cups oil",
						"1/2 teaspoon ground cardamom",
						"3/4 cup finely chopped pistachios *",
					},
					Instructions: []string{
						"*Note: Instead of pistachio nuts, walnuts may be used if desired.",
						"Dissolve sugar in water in heavy pan over medium heat. Bring to the boil, add lemon juice and saffron and boil for 10 minutes. Cool and strain into a 25 cm (10 inch) pie plate. Keep aside. Break eggs into a casserole dish about 20 cm (8 inches) in diameter.",
						"The size and flat base are important. Add salt and mix eggs with fork until yolks and whites are thoroughly combined - do not beat as eggs must not be foamy. Heat oil in an electric frying pan to 190 C (375 F) or in a 25 cm (10 inch) frying pan placed on a thermostatically controlled hot plate or burner. Have ready nearby a long skewer, the plate of syrup, a baking sheet and the nuts mixed with the cardamom.",
						"A bowl of water and a cloth for drying hands are also necessary. Hold dish with eggs in one hand next to the pan of oil and slightly above it. Put hand into egg, palm down, so that egg covers back of hand. Lift out hand, curling fingers slightly inwards, then open out over hot oil, fingers pointing down.",
						"Move hand across surface of oil so that egg falls in streams from fingertips. Dip hand in egg again and make more strands across those already in pan. Repeat 3 to 4 times until about an eighth of the egg is used. There should be a closely meshed layer of egg strands about 20 cm (8 inches) across. Work quickly so that the last lot of egg is added not long after the first lot. Rinse hands quickly and dry. Take skewer and slide under bubbling omelet, lift up and turn over to lightly brown other side.",
						"The first side will be bubbly, the underside somewhat smoother. When golden brown lift out with skewer and drain over pan. Place omelet flat in the syrup, spoon syrup over the top and lift out with skewer onto baking sheet. Roll up with bubbly side inwards.",
						"Finish roll should be about 3 cm (1 1/4 inches) in diameter. Put to one side and sprinkle with nuts. Repeat with remaining egg, making 7 to 8 rolls in all. Though depth of egg diminishes, you will become so adept that somehow you will get it in the pan in fine strands.",
						"When cool, cut kabaubs into 4-5 cm (1 1/2 to 2 inch pieces and serve. These keep well in a sealed container in a cool place.",
					},
					Keywords: []string{"*Afghan", "Desserts"},
					Name:     "Abraysham Kabaub (Silk Kebab)",
					URL:      "Imported from MasterCook",
					Yield:    30,
				},
				{
					Category: "Middle East",
					Ingredients: []string{
						"2 large clov garlic",
						"1/2 teaspoon salt",
						"2 cups plain whole-milk yogurt",
						"4 tablespoons juice and pulp of 1 large lemon",
						"1/2 teaspoon cracked black pepper",
						"2 large whol chicken breasts -- about 2 pounds",
					},
					Instructions: []string{
						"Long, slow marinating in garlicky yogurt tenderizes, moistens and adds deep flavor, so you end up with skinless grilled chicken that's as delicious as it is nutritionally correct. Serve with soft pita or Arab flatbread and fresh yogurt.",
						"Put the salt in a wide, shallow non-reactive bowl with the garlic and mash them together until you have paste. Add yogurt, lemon and pepper.",
						"Skin the chicken breasts, remove all visible fat and separate the halves. Bend each backward to break the bones so the pieces win lie flat. Add to the yogurt and turn so all surfaces are well-coated.",
						"Cover the bowl tightly and refrigerate. Allow to marinate at least overnight, up to a day and a half. Turn when you think of it.",
						"To cook, remove breasts from marinade and wipe off all but a thin film. Broil or grill about 6 inches from the heat for 6 to 8 minutes a side, or until thoroughly cooked. Meat will brown somewhat but should not char. Serve at once.",
					},
					Keywords: []string{"Middle East", "Chicken", "*Afghan", "On-The-Grill"},
					Name:     "Afghan Chicken",
					URL:      "Imported from MasterCook",
					Yield:    6,
				},
				{
					Category: "Chicken",
					Ingredients: []string{
						"1 cup yogurt",
						"1 1/2 teaspoons salt",
						"1/2 teaspoon ground red or black pepper",
						"3 centiliters garlic -- finely minced",
						"1 1/2 pounds chicken breasts -- boneless,",
						"skinless -- cut into kebob",
						"-- ¥",
						"flatbread such as lavash",
						"pita or flour tortillas",
						"3 tomatoes -- sliced",
						"2 onions -- sliced",
						"cilantro to taste",
						"2 lemons or 4 limes -- quartered",
					},
					Instructions: []string{
						"Mix yogurt, salt, pepper and garlic in a bowl. Mix chicken with yogurt and marinate 1 to 2 hours at room temperature, up to 2 days refrigerated.",
						"Thread chicken on skewers and grill over medium hot coals.",
						"Place warmed pita bread on plates (if using tortillas, toast briefly over flame), divide meat among them, top with tomato and onion slices and cilantro and fold bread over. Serve with lemon or lime quarters for squeezing.",
					},
					Keywords: []string{"Chicken", "*Afghan", "On-The-Grill"},
					Name:     "Afghan Chicken Kebobs",
					URL:      "Imported from MasterCook",
					Yield:    4,
				},
				{
					Category: "*Afghan",
					Ingredients: []string{
						"2 pounds fresh pumpkin or squash",
						"1/4 cup corn oil",
						"***SWEET TOMATO SAUCE***",
						"1 teaspoon crushed garlic",
						"1 cup water",
						"1/2 teaspoon salt",
						"1/2 cup sugar",
						"4 ounces tomato sauce",
						"1/2 teaspoon ginger root -- chopped fine",
						"1 teaspoon freshly ground coriander",
						"seeds",
						"1/4 teaspoon black pepper",
						"***YOGURT SAUCE***",
						"1/4 teaspoon crushed garlic",
						"1/4 teaspoon salt",
						"3/4 cup plain yogurt",
						"***GARNISH***",
						"dry mint leaves -- crushed",
					},
					Instructions: []string{
						`Peel the pumpkin and cut into 2-3" cubes; set aside. Heat oil in a large frying pan that has a lid. Fry the pumpkins on both sides for a couple of minutes until lightly browned. Mix together ingredients for Sweet Tomato Sauce in a bowl then add to pumpkin mixture in fry pan.`,
						"Cover and cook 20-25 minutes over low heat until the pumpkin is cooked and most of the liquid has evaporated. (I don't know how it's going to evaporate if the pan is covered....-B.)",
						"Mix together the ingredients for the yogurt sauce. To serve: Spread half the yogurt sauce on a plate and lay the pumpkin on top. Top with remaining yogurt and any cooking juices left over. Sprinkle with dry mint.",
						"May be served with chalow (basmati rice) and naan or pita bread.",
						"From Afghani Cooking, the cookbook from Da Afghan Restaurant, Bloomington/Minneapolis, MN.",
					},
					Keywords: []string{"*Afghan", "Vegetables"},
					Name:     "Afghan Pumpkins Kadu Bouranee (Sweet Pumpkin)",
					URL:      "Imported from MasterCook",
					Yield:    1,
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := models.NewRecipesFromMasterCook(tc.buf)
			if !cmp.Equal(got, tc.want) {
				t.Log(cmp.Diff(got, tc.want))
				t.Fail()
			}
		})
	}
}

func TestSort_IsSort(t *testing.T) {
	t.Run("no sort", func(t *testing.T) {
		s := models.Sort{}
		if s.IsSort() {
			t.Fail()
		}
	})

	t.Run("is sort", func(t *testing.T) {
		testcases := []struct {
			name string
			in   models.Sort
		}{
			{name: "A to Z enabled", in: models.Sort{IsAToZ: true}},
			{name: "Z to A enabled", in: models.Sort{IsZToA: true}},
		}
		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				if !tc.in.IsSort() {
					t.Fail()
				}
			})
		}
	})
}

func assertNoError(tb testing.TB, got error) {
	tb.Helper()
	if got != nil {
		tb.Fatal("got an error but expected none")
	}
}

func assertStructsEqual[T models.Recipe](tb testing.TB, got, want T) {
	tb.Helper()
	if !cmp.Equal(got, want) {
		tb.Log(cmp.Diff(got, want))
		tb.Fail()
	}
}
