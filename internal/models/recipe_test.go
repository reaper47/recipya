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
			want:   models.SearchOptionsRecipes{ByName: true},
		},
		{
			name:   "name",
			method: "name",
			want:   models.SearchOptionsRecipes{ByName: true},
		},
		{
			name:   "empty defaults to name",
			method: "full",
			want:   models.SearchOptionsRecipes{FullSearch: true},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			compare(t, models.NewSearchOptionsRecipe(tc.method), tc.want)
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

func TestNewRecipeFromTextFile(t *testing.T) {
	testcases := []struct {
		name string
		buf  *bytes.Buffer
		want models.Recipe
	}{
		{
			name: "recipe 1",
			buf:  bytes.NewBuffer([]byte("Bailey's Irish Cream Chocolate Cake with Kahlua Buttercream\n\nMakes: 12 slices\n\nIngredients\nCake\n1 cup unsalted butter, softened to room temperature\n2 1/2 cups granulated sugar\n4 large eggs, room temperature\n2 teaspoons vanilla extract\n3 cups all-purpose flour\n1 1/2 teaspoons baking soda\n1 1/2 teaspoon baking powder\n1/2 teaspoon salt\n1 cup unsweetened cocoa powder\n2 1/4 cups Baileys Irish Cream\nKahlua Buttercream Frosting\n1 cup unsalted butter, softened to room temperature\n4 cups powdered sugar\n1 teaspoon vanilla\n3 tablespoons kahlua\n\nInstructions\nPreheat oven to 350 degrees.\nGrease and flour the cake pans and then line the bottom of the pans with round parchment paper for easier removal.\nCream the butter and sugars in a large mixing bowl. Add in the eggs one at a time and blend with a handheld mixer.\nAdd in the vanilla and mix to combine.\nCombine flour, baking soda, baking powder, salt and cocoa powder in a large bowl.\nAlternating between the two, slowly add the dry ingredient mixture and Baileys to the bowl with the butter and sugar. Mix on low speed to combine ingredients.\nDivide the cake batter evenly between the three cake pans.\nBake 22-28 minutes or until a toothpick inserted in the center of the cake comes out clean.\nCool and then remove the top of each cake with a cake leveler.\nMix the frosting by creaming the butter and powdered sugar. Add in the vanilla and kahlua, mixing with a hand-held mixer on low speed.\nENJOY!\n\nhttps://myincrediblerecipes.com/baileys-irish-cream-chocolate-cake-with-kahlua-buttercream-frosting/")),
			want: models.Recipe{
				Category: "uncategorized",
				Ingredients: []string{
					"Cake",
					"1 cup unsalted butter, softened to room temperature",
					"2 1/2 cups granulated sugar",
					"4 large eggs, room temperature",
					"2 teaspoons vanilla extract",
					"3 cups all-purpose flour",
					"1 1/2 teaspoons baking soda",
					"1 1/2 teaspoon baking powder",
					"1/2 teaspoon salt",
					"1 cup unsweetened cocoa powder",
					"2 1/4 cups Baileys Irish Cream",
					"Kahlua Buttercream Frosting",
					"1 cup unsalted butter, softened to room temperature",
					"4 cups powdered sugar",
					"1 teaspoon vanilla",
					"3 tablespoons kahlua",
				},
				Instructions: []string{
					"Preheat oven to 350 degrees.",
					"Grease and flour the cake pans and then line the bottom of the pans with round parchment paper for easier removal.",
					"Cream the butter and sugars in a large mixing bowl. Add in the eggs one at a time and blend with a handheld mixer.",
					"Add in the vanilla and mix to combine.",
					"Combine flour, baking soda, baking powder, salt and cocoa powder in a large bowl.",
					"Alternating between the two, slowly add the dry ingredient mixture and Baileys to the bowl with the butter and sugar. Mix on low speed to combine ingredients.",
					"Divide the cake batter evenly between the three cake pans.",
					"Bake 22-28 minutes or until a toothpick inserted in the center of the cake comes out clean.",
					"Cool and then remove the top of each cake with a cake leveler.",
					"Mix the frosting by creaming the butter and powdered sugar. Add in the vanilla and kahlua, mixing with a hand-held mixer on low speed.",
					"ENJOY!",
				},
				Keywords: make([]string, 0),
				Name:     "Bailey's Irish Cream Chocolate Cake with Kahlua Buttercream",
				Tools:    make([]string, 0),
				URL:      "https://myincrediblerecipes.com/baileys-irish-cream-chocolate-cake-with-kahlua-buttercream-frosting/",
				Yield:    12,
			},
		},
		{
			name: "recipe 2",
			buf:  bytes.NewBuffer([]byte("Beef and Guinness Stew\n\nThis hearty beef stew is made with lean boneless chuck that's cooked with carrots, parsnips and turnips and flavored with dark beer. Simmering it in a Dutch oven for about 2 hours makes the meat and vegetables fork tender and delicious.\n\nYield: 8 servings (serving size: about 1 cup)\nTotal: 3 Hours, 18 Minutes\n\nIngredients\n3 tablespoons canola oil, divided\n1/4 cup all-purpose flour $\n2 pounds boneless chuck roast, trimmed and cut into 1-inch cubes\n1 teaspoon salt, divided\n5 cups chopped onion (about 3 onions)\n1 tablespoon tomato paste $\n4 cups fat-free, lower-sodium beef broth $\n1 (11.2-ounce) bottle Guinness Stout\n1 tablespoon raisins\n1 teaspoon caraway seeds\n1/2 teaspoon black pepper\n1 1/2 cups (1/2-inch-thick) diagonal slices carrot (about 8 ounces) $\n1 1/2 cups (1/2-inch-thick) diagonal slices parsnip (about 8 ounces)\n1 cup (1/2-inch) cubed peeled turnip (about 8 ounces)\n2 tablespoons finely chopped fresh flat-leaf parsley\nPreparation\n\n1. Heat 1 1/2 tablespoons oil in a Dutch oven over medium-high heat. Place flour in a shallow dish. Sprinkle beef with 1/2 teaspoon salt; dredge beef in flour. Add half of beef to pan; cook 5 minutes, turning to brown on all sides. Remove beef from pan with a slotted spoon. Repeat procedure with remaining 1 1/2 tablespoons oil and beef.\n2. Add onion to pan; cook 5 minutes or until tender, stirring occasionally. Stir in tomato paste; cook 1 minute, stirring frequently. Stir in broth and beer, scraping pan to loosen browned bits. Return meat to pan. Stir in remaining 1/2 teaspoon salt, raisins, caraway seeds, and pepper; bring to a boil. Cover, reduce heat, and simmer 1 hour, stirring occasionally. Uncover and bring to a boil. Cook 50 minutes, stirring occasionally. Add carrot, parsnip, and turnip. Cover, reduce heat to low, and simmer 30 minutes, stirring occasionally. Uncover and bring to a boil; cook 10 minutes or until vegetables are tender. Sprinkle with parsley.\n\nhttp://www.myrecipes.com/recipe/beef-guinness-stew-10000001963989/")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "This hearty beef stew is made with lean boneless chuck that's cooked with carrots, parsnips and turnips and flavored with dark beer. Simmering it in a Dutch oven for about 2 hours makes the meat and vegetables fork tender and delicious.",
				Ingredients: []string{
					"3 tablespoons canola oil, divided",
					"1/4 cup all-purpose flour $",
					"2 pounds boneless chuck roast, trimmed and cut into 1-inch cubes",
					"1 teaspoon salt, divided",
					"5 cups chopped onion (about 3 onions)",
					"1 tablespoon tomato paste $",
					"4 cups fat-free, lower-sodium beef broth $",
					"1 (11.2-ounce) bottle Guinness Stout",
					"1 tablespoon raisins",
					"1 teaspoon caraway seeds",
					"1/2 teaspoon black pepper",
					"1 1/2 cups (1/2-inch-thick) diagonal slices carrot (about 8 ounces) $",
					"1 1/2 cups (1/2-inch-thick) diagonal slices parsnip (about 8 ounces)",
					"1 cup (1/2-inch) cubed peeled turnip (about 8 ounces)",
					"2 tablespoons finely chopped fresh flat-leaf parsley",
				},
				Instructions: []string{
					"Heat 1 1/2 tablespoons oil in a Dutch oven over medium-high heat. Place flour in a shallow dish. Sprinkle beef with 1/2 teaspoon salt; dredge beef in flour. Add half of beef to pan; cook 5 minutes, turning to brown on all sides. Remove beef from pan with a slotted spoon. Repeat procedure with remaining 1 1/2 tablespoons oil and beef.",
					"Add onion to pan; cook 5 minutes or until tender, stirring occasionally. Stir in tomato paste; cook 1 minute, stirring frequently. Stir in broth and beer, scraping pan to loosen browned bits. Return meat to pan. Stir in remaining 1/2 teaspoon salt, raisins, caraway seeds, and pepper; bring to a boil. Cover, reduce heat, and simmer 1 hour, stirring occasionally. Uncover and bring to a boil. Cook 50 minutes, stirring occasionally. Add carrot, parsnip, and turnip. Cover, reduce heat to low, and simmer 30 minutes, stirring occasionally. Uncover and bring to a boil; cook 10 minutes or until vegetables are tender. Sprinkle with parsley.",
				},
				Keywords: make([]string, 0),
				Name:     "Beef and Guinness Stew",
				Times: models.Times{
					Cook: 3*time.Hour + 18*time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "http://www.myrecipes.com/recipe/beef-guinness-stew-10000001963989/",
				Yield: 8,
			},
		},
		{
			name: "recipe 3",
			buf:  bytes.NewBuffer([]byte("Cauliflower Jalapeno Popper Soup\n\nCauliflower Jalapeno Popper Soup is a rich and creamy low carb soup full of cauliflower, bacon, and jalapenos. This soup really packs the heat.\n\nCourse: Soup\nCuisine: American\nKeyword: cauliflower, low carb\nPrep Time: 15 minutes\nCook Time: 18 minutes\nTotal Time: 33 minutes\nServings: 6\n\nIngredients\n6 slices bacon\n1 cup (150 g) diced onion\n4 to 5 jalapenos, finely diced\n2 garlic cloves, minced\n2 tablespoons all-purpose flour\n3 ½ cups (8.3 dl) chicken broth\n½ teaspoon salt\n½ teaspoon black pepper\n½ teaspoon paprika\n¼ teaspoon garlic powder\n¼ teaspoon onion powder\n1 small head cauliflower, cut into small florets\n1 cup (2.4 dl) heavy cream\n8 ounces (225 g) cream cheese, cut into cubes and softened\n2 cups (225 g) shredded sharp cheddar cheese\n2 green onions, sliced\n\nInstructions\nCook bacon in a large Dutch oven until crispy. Remove bacon and set aside.\nLeave 2 tablespoons bacon grease in pot.\nAdd onion and jalapenos to pot. Cook over medium heat for 4 to 5 minutes to soften.\nAdd garlic and cook 30 seconds. Stir in flour. Cook and stir for 1 minute.\nGradually whisk in chicken broth. Add salt, pepper, paprika, garlic powder, and onion powder. Bring mixture to a simmer.\nAdd cauliflower. Cook for 10 minutes or until cauliflower is soft.\nAdd heavy cream and cream cheese. Turn heat to low and whisk until cream cheese is completely mixed in.\nWhisk in cheese 1/2 cup at a time.\nLadle into bowls. Crumble the bacon and sprinkle on top. Top with green onion.\n\nhttps://skinnysouthernrecipes.com/cauliflower-jalapeno-popper-soup/")),
			want: models.Recipe{
				Category:    "soup",
				Cuisine:     "american",
				Description: "Cauliflower Jalapeno Popper Soup is a rich and creamy low carb soup full of cauliflower, bacon, and jalapenos. This soup really packs the heat.",
				Ingredients: []string{
					"6 slices bacon",
					"1 cup (150 g) diced onion",
					"4 to 5 jalapenos, finely diced",
					"2 garlic cloves, minced",
					"2 tablespoons all-purpose flour",
					"3 ½ cups (8.3 dl) chicken broth",
					"½ teaspoon salt",
					"½ teaspoon black pepper",
					"½ teaspoon paprika",
					"¼ teaspoon garlic powder",
					"¼ teaspoon onion powder",
					"1 small head cauliflower, cut into small florets",
					"1 cup (2.4 dl) heavy cream",
					"8 ounces (225 g) cream cheese, cut into cubes and softened",
					"2 cups (225 g) shredded sharp cheddar cheese",
					"2 green onions, sliced",
				},
				Instructions: []string{
					"Cook bacon in a large Dutch oven until crispy. Remove bacon and set aside.",
					"Leave 2 tablespoons bacon grease in pot.",
					"Add onion and jalapenos to pot. Cook over medium heat for 4 to 5 minutes to soften.",
					"Add garlic and cook 30 seconds. Stir in flour. Cook and stir for 1 minute.",
					"Gradually whisk in chicken broth. Add salt, pepper, paprika, garlic powder, and onion powder. Bring mixture to a simmer.",
					"Add cauliflower. Cook for 10 minutes or until cauliflower is soft.",
					"Add heavy cream and cream cheese. Turn heat to low and whisk until cream cheese is completely mixed in.",
					"Whisk in cheese 1/2 cup at a time.",
					"Ladle into bowls. Crumble the bacon and sprinkle on top. Top with green onion.",
				},
				Keywords: []string{"cauliflower", "low carb"},
				Name:     "Cauliflower Jalapeno Popper Soup",
				Times: models.Times{
					Prep: 15 * time.Minute,
					Cook: 18 * time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "https://skinnysouthernrecipes.com/cauliflower-jalapeno-popper-soup/",
				Yield: 6,
			},
		},
		{
			name: "recipe 4",
			buf:  bytes.NewBuffer([]byte("Classic roast chicken & gravy\n\nA classic roast chicken recipe should be in everyone's repertoire. Make this oven baked chicken with gravy for a satisfying Sunday lunch.\n\nPrep:10 mins\nCook:1 hr and 30 mins\nPlus resting\nEasy\nServes 4\n\nIngredients\n1 onion, roughly chopped\n2 carrots, roughly chopped\n1 free range chicken, about 1.5kg/3lb 5oz\n1 lemon, halved\nsmall bunch thyme (optional)\n25g butter, softened\n\nFor the gravy\n1 tbsp plain flour\n250ml chicken stock (a cube is fine)\n\nMethod\nHeat oven to 190C/fan 170C/gas 5. Have a shelf ready in the middle of the oven without any shelves above it.\nScatter 1 roughly chopped onion and 2 roughly chopped carrots over the base of a roasting tin that fits the whole 1 ½ kg chicken, but doesn’t swamp it.\nSeason the cavity of the chicken liberally with salt and pepper, then stuff with 2 lemon halves and a small bunch of thyme, if using.\nSit the chicken on the vegetables, smother the breast and legs all over with 25g softened butter, then season the outside with salt and pepper.\nPlace in the oven and leave, undisturbed, for 1 hr 20 mins – this will give you a perfectly roasted chicken. To check, pierce the thigh with a skewer and the juices should run clear.\nCarefully remove the tin from the oven and, using a pair of tongs, lift the chicken to a dish or board to rest for 15-20 mins. As you lift the dish, let any juices from the chicken pour out of the cavity into the roasting tin.\nWhile the chicken is resting, make the gravy. Place the roasting tin over a low flame, then stir in 1 tbsp flour and sizzle until you have a light brown, sandy paste.\nGradually pour in 250ml chicken stock, stirring all the time, until you have a thickened sauce.\nSimmer for 2 mins, using a wooden spoon to stir, scraping any sticky bits from the tin.\nStrain the gravy into a small saucepan, then simmer and season to taste. When you carve the bird, add any extra juices to the gravy.\n\nHow long to cook roast chicken\nA 1.5kg chicken will be perfectly roasted after 1 hr 20 mins at 190C/fan 170C/gas 5.\nIt doesn’t matter what you stuff into it, rub or sprinkle over it or put around it, this timing never changes.\nRemember this and you will always be able to roast a chicken.\n\nTips for making the perfect roast chicken:\nAlways leave your chicken to rest for at least 15 mins before carving. This will give you a juicy chicken that is a lot easier to carve.\nFor a more succulent chicken, take it out of the fridge one hour before cooking to bring it up to room temperature. This rule applies to any meat you are roasting.\nDon't worry about turning or basting your chicken as it roasts. Yes, these can give good results, but are fiddly for the beginner, plus every time you open the oven you lose heat.\n\nhttps://www.bbcgoodfood.com/recipes/classic-roast-chicken-gravy")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "A classic roast chicken recipe should be in everyone's repertoire. Make this oven baked chicken with gravy for a satisfying Sunday lunch.",
				Ingredients: []string{
					"1 onion, roughly chopped",
					"2 carrots, roughly chopped",
					"1 free range chicken, about 1.5kg/3lb 5oz",
					"1 lemon, halved",
					"small bunch thyme (optional)",
					"25g butter, softened",
					"For the gravy",
					"1 tbsp plain flour",
					"250ml chicken stock (a cube is fine)",
				},
				Instructions: []string{
					"Heat oven to 190C/fan 170C/gas 5. Have a shelf ready in the middle of the oven without any shelves above it.",
					"Scatter 1 roughly chopped onion and 2 roughly chopped carrots over the base of a roasting tin that fits the whole 1 ½ kg chicken, but doesn’t swamp it.",
					"Season the cavity of the chicken liberally with salt and pepper, then stuff with 2 lemon halves and a small bunch of thyme, if using.",
					"Sit the chicken on the vegetables, smother the breast and legs all over with 25g softened butter, then season the outside with salt and pepper.",
					"Place in the oven and leave, undisturbed, for 1 hr 20 mins – this will give you a perfectly roasted chicken. To check, pierce the thigh with a skewer and the juices should run clear.",
					"Carefully remove the tin from the oven and, using a pair of tongs, lift the chicken to a dish or board to rest for 15-20 mins. As you lift the dish, let any juices from the chicken pour out of the cavity into the roasting tin.",
					"While the chicken is resting, make the gravy. Place the roasting tin over a low flame, then stir in 1 tbsp flour and sizzle until you have a light brown, sandy paste.",
					"Gradually pour in 250ml chicken stock, stirring all the time, until you have a thickened sauce.",
					"Simmer for 2 mins, using a wooden spoon to stir, scraping any sticky bits from the tin.",
					"Strain the gravy into a small saucepan, then simmer and season to taste. When you carve the bird, add any extra juices to the gravy.",
					"How long to cook roast chicken",
					"A 1.5kg chicken will be perfectly roasted after 1 hr 20 mins at 190C/fan 170C/gas 5.",
					"It doesn’t matter what you stuff into it, rub or sprinkle over it or put around it, this timing never changes.",
					"Remember this and you will always be able to roast a chicken.",
					"Tips for making the perfect roast chicken:",
					"Always leave your chicken to rest for at least 15 mins before carving. This will give you a juicy chicken that is a lot easier to carve.",
					"For a more succulent chicken, take it out of the fridge one hour before cooking to bring it up to room temperature. This rule applies to any meat you are roasting.",
					"Don't worry about turning or basting your chicken as it roasts. Yes, these can give good results, but are fiddly for the beginner, plus every time you open the oven you lose heat.",
				},
				Keywords: make([]string, 0),
				Name:     "Classic roast chicken & gravy",
				Times: models.Times{
					Prep: 10 * time.Minute,
					Cook: 1*time.Hour + 30*time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "https://www.bbcgoodfood.com/recipes/classic-roast-chicken-gravy",
				Yield: 4,
			},
		},
		{
			name: "recipe 5",
			buf:  bytes.NewBuffer([]byte("Crispy Thai Pork with Cucumber Salad\n\nThe crispy bits are the key to this pork’s deliciousness. Use a wide, flat spatula for the best effect.\n\n4 Servings\n\nIngredients\n1/2 English hothouse cucumber, halved lengthwise, thinly sliced crosswise\n1 shallot, thinly sliced\n2 red or green Thai chiles, with seeds, thinly sliced, divided\n1/2 cup fresh lime juice, divided\nKosher salt\n3 tablespoons vegetable oil\n6 garlic cloves, thinly sliced\n1 pound ground pork\nFreshly ground black pepper\n1/4 cup low-sodium chicken broth\n2 tablespoons reduced-sodium soy sauce\n1 tablespoon fish sauce (such as nam pla and nuoc nam)\n2 teaspoons light brown sugar\n1 cup fresh basil leaves\n1/2 cup fresh cilantro leaves\n1/2 cup fresh mint leaves\n2 heads Bibb lettuce, leaves separated\nSteamed white rice and lime wedges (for serving)\n\nRecipe Preparation\nToss cucumber, shallot, 1 chile, and 1/4 cup lime juice in a medium bowl; season with salt and set cucumber salad aside.\nHeat oil in a large skillet over high heat. Add garlic and remaining chile and cook, stirring, until garlic is just beginning to turn golden, about 30 seconds. Add pork, season with salt and pepper, and cook, breaking up with a spoon and pressing down firmly to help brown, until cooked through, browned, and crisp in spots, 6–8 minutes.\nAdd broth, soy sauce, fish sauce, and brown sugar and cook, scraping up any brown bits from the bottom of skillet, until liquid is almost completely evaporated, about 2 minutes. Mix in remaining 1/4 cup lime juice.\nToss basil, cilantro, and mint in a medium bowl. Serve pork with herbs, lettuce, rice, lime wedges, and reserved cucumber salad alongside.\n\nhttps://www.bonappetit.com/recipe/crispy-thai-pork-cucumber-salad")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "The crispy bits are the key to this pork’s deliciousness. Use a wide, flat spatula for the best effect.",
				Ingredients: []string{
					"1/2 English hothouse cucumber, halved lengthwise, thinly sliced crosswise",
					"1 shallot, thinly sliced",
					"2 red or green Thai chiles, with seeds, thinly sliced, divided",
					"1/2 cup fresh lime juice, divided",
					"Kosher salt",
					"3 tablespoons vegetable oil",
					"6 garlic cloves, thinly sliced",
					"1 pound ground pork",
					"Freshly ground black pepper",
					"1/4 cup low-sodium chicken broth",
					"2 tablespoons reduced-sodium soy sauce",
					"1 tablespoon fish sauce (such as nam pla and nuoc nam)",
					"2 teaspoons light brown sugar",
					"1 cup fresh basil leaves",
					"1/2 cup fresh cilantro leaves",
					"1/2 cup fresh mint leaves",
					"2 heads Bibb lettuce, leaves separated",
					"Steamed white rice and lime wedges (for serving)",
				},
				Instructions: []string{
					"Toss cucumber, shallot, 1 chile, and 1/4 cup lime juice in a medium bowl; season with salt and set cucumber salad aside.",
					"Heat oil in a large skillet over high heat. Add garlic and remaining chile and cook, stirring, until garlic is just beginning to turn golden, about 30 seconds. Add pork, season with salt and pepper, and cook, breaking up with a spoon and pressing down firmly to help brown, until cooked through, browned, and crisp in spots, 6–8 minutes.",
					"Add broth, soy sauce, fish sauce, and brown sugar and cook, scraping up any brown bits from the bottom of skillet, until liquid is almost completely evaporated, about 2 minutes. Mix in remaining 1/4 cup lime juice.",
					"Toss basil, cilantro, and mint in a medium bowl. Serve pork with herbs, lettuce, rice, lime wedges, and reserved cucumber salad alongside.",
				},
				Keywords: make([]string, 0),
				Name:     "Crispy Thai Pork with Cucumber Salad",
				Tools:    make([]string, 0),
				URL:      "https://www.bonappetit.com/recipe/crispy-thai-pork-cucumber-salad",
				Yield:    4,
			},
		},
		{
			name: "recipe 6",
			buf:  bytes.NewBuffer([]byte("Domada (Gambian Peanut Stew)\n\nThe national dish of Gambia. A thick, saucy stew served over rice.\n\nPrep time 10 mins\nCook time 1 hour\nTotal time 1 hour 10 mins\nCuisine: African\nServes: 4\n\nIngredients\n1 lb beef steak or 1 lb chicken breast, cut into ½ inch chunks (or use bone-in chicken pieces and simmer them in the sauce; once cooked leave the pieces whole or remove the meat from the bones and add it back to the stew.)\n1 large onion, diced\n2 tablespoons olive oil\n3 cloves garlic, minced\n3 Roma tomatoes, diced\n½ can (3 oz) tomato paste\n¾ cup natural, unsweetened peanut butter\n4 Maggi or Knorr tomato bouillon cubes\n3 cups water\nScotch bonnet chilies, diced, according to heat preference\n4 cups pumpkin or sweet potato, diced\nSalt and pepper to taste\n\nInstructions\nHeat the oil in large Dutch oven. Saute the onions until golden. Add the beef and garlic and continue to sauté until the beef is no longer pink. Add the tomatoes and cook for 3 minutes. Add the tomato paste, chilies, peanut butter and stir to combine. Add the water and bouillon cubes. Bring to a boil, reduce heat, cover, and simmer for 15 minutes, stirring occasionally. Add squash, cover, and continue to cook for 35-40 minutes or until the pumpkin is tender, stirring occasionally. Season with salt and pepper.\nServe hot with rice. This stew tastes even better the next day.\n\nhttps://www.daringgourmet.com/domoda-gambian-peanut-stew/")),
			want: models.Recipe{
				Category:    "uncategorized",
				Cuisine:     "african",
				Description: "The national dish of Gambia. A thick, saucy stew served over rice.",
				Ingredients: []string{
					"1 lb beef steak or 1 lb chicken breast, cut into ½ inch chunks (or use bone-in chicken pieces and simmer them in the sauce; once cooked leave the pieces whole or remove the meat from the bones and add it back to the stew.)",
					"1 large onion, diced",
					"2 tablespoons olive oil",
					"3 cloves garlic, minced",
					"3 Roma tomatoes, diced",
					"½ can (3 oz) tomato paste",
					"¾ cup natural, unsweetened peanut butter",
					"4 Maggi or Knorr tomato bouillon cubes",
					"3 cups water",
					"Scotch bonnet chilies, diced, according to heat preference",
					"4 cups pumpkin or sweet potato, diced",
					"Salt and pepper to taste",
				},
				Instructions: []string{
					"Heat the oil in large Dutch oven. Saute the onions until golden. Add the beef and garlic and continue to sauté until the beef is no longer pink. Add the tomatoes and cook for 3 minutes. Add the tomato paste, chilies, peanut butter and stir to combine. Add the water and bouillon cubes. Bring to a boil, reduce heat, cover, and simmer for 15 minutes, stirring occasionally. Add squash, cover, and continue to cook for 35-40 minutes or until the pumpkin is tender, stirring occasionally. Season with salt and pepper.",
					"Serve hot with rice. This stew tastes even better the next day.",
				},
				Keywords:  make([]string, 0),
				Name:      "Domada (Gambian Peanut Stew)",
				Nutrition: models.Nutrition{},
				Times: models.Times{
					Prep: 10 * time.Minute,
					Cook: 1 * time.Hour,
				},
				Tools: make([]string, 0),
				URL:   "https://www.daringgourmet.com/domoda-gambian-peanut-stew/",
				Yield: 4,
			},
		},
		{
			name: "recipe 7",
			buf:  bytes.NewBuffer([]byte("Duck Fat-Roasted Brussels Sprouts\n\nA few tablespoons of duck fat and a very hot oven are all you need to turn some sleepy Brussels sprouts into something much more special.\n\nPrep Time: 20 mins\nCook Time: 20 mins\nTotal Time: 40 mins\nServings: 4\n\nIngredients\n2 tablespoons duck fat, or more as needed\n2 pounds Brussels sprouts, trimmed and halved lengthwise\nsalt and freshly ground black pepper to taste\n1 pinch cayenne pepper, or more to taste\n1 lemon, juiced\n\nDirections\nPreheat the oven to 450 degrees F (230 degrees C). Line a baking sheet with parchment paper or a silicone baking mat.\nHeat duck fat in a small saucepan over low heat until melted. Set aside.\nCombine Brussels sprouts, salt, black pepper, and cayenne pepper in a large bowl. Pour melted duck fat over Brussels sprouts and stir to coat evenly. Spread evenly onto the prepared baking sheet.\nBake in the preheated oven until Brussels sprouts are browned and tender, but still slightly firm, 15 to 20 minutes, flipping sprouts halfway through. Top with freshly squeezed lemon juice.\n\nhttps://www.allrecipes.com/recipe/231351/duck-fat-roasted-brussels-sprouts/")),
			want: models.Recipe{
				Category:    "uncategorized",
				CreatedAt:   time.Time{},
				Description: "A few tablespoons of duck fat and a very hot oven are all you need to turn some sleepy Brussels sprouts into something much more special.",
				Ingredients: []string{
					"2 tablespoons duck fat, or more as needed",
					"2 pounds Brussels sprouts, trimmed and halved lengthwise",
					"salt and freshly ground black pepper to taste",
					"1 pinch cayenne pepper, or more to taste",
					"1 lemon, juiced",
				},
				Instructions: []string{
					"Preheat the oven to 450 degrees F (230 degrees C). Line a baking sheet with parchment paper or a silicone baking mat.",
					"Heat duck fat in a small saucepan over low heat until melted. Set aside.",
					"Combine Brussels sprouts, salt, black pepper, and cayenne pepper in a large bowl. Pour melted duck fat over Brussels sprouts and stir to coat evenly. Spread evenly onto the prepared baking sheet.",
					"Bake in the preheated oven until Brussels sprouts are browned and tender, but still slightly firm, 15 to 20 minutes, flipping sprouts halfway through. Top with freshly squeezed lemon juice.",
				},
				Keywords: make([]string, 0),
				Name:     "Duck Fat-Roasted Brussels Sprouts",
				Times: models.Times{
					Prep: 20 * time.Minute,
					Cook: 20 * time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "https://www.allrecipes.com/recipe/231351/duck-fat-roasted-brussels-sprouts/",
				Yield: 4,
			},
		},
		{
			name: "recipe 8",
			buf:  bytes.NewBuffer([]byte("Gin and tonic lemon & lime cheesecake\n\nThis easy no-bake cheesecake is subtly infused with the flavours of gin, lemon and lime for a deliciously zesty and zingy dessert. \nThe beauty of this pudding (apart from being so easy to make - you don’t even need to turn the oven on!) is that you can make it the day before you want to eat it, and it keeps in the fridge for second servings the next day - if it lasts that long, of course!\nThis fabulous gin and tonic cheesecake is best served with - what else? - a refreshing G&T, in the garden, in the sunshine. Bliss!\n\nFor the base:\n200g digestive biscuits\n100g butter\n\nFor the filling:\n500g cream cheese\n100g icing sugar\n250ml double cream\n50ml gin\nZest and juice of one lime or lemon, or half of each\n\nFor the topping:\n30ml gin\n100ml tonic\n2 tbsp sugar\n1tbsp freshly squeezed lemon juice\nFresh lemon and/or lime slices, to decorate\n\nTo make the base, blitz the digestive biscuits in a food processor or crush them in a sandwich bag with a rolling pin until you have fine crumbs. Melt the butter and combine thoroughly with the biscuit. You want this mixture to be squidgy enough that it will form a base, but not overly wet - if you need to add more butter or biscuit crumbs to get the right texture, play around with the quantities.\nButter and line a 23cm loose-bottomed tin then press the buttery crumb mixture firmly down into the base to create an even layer. Chill in the fridge for an hour or until it is firmly set.\nMeanwhile, mix the cream cheese, icing sugar, gin, lemon and lime zest and juice and cream together in a bowl until completely combined. Spoon the cream mixture onto the biscuit base and smooth the surface with a flat knife or spatula. Leave to set in the fridge for at least six hours or overnight.\nTo make the topping, gently heat the gin and tonic with the lemon juice and sugar and stir gently until all the sugar has dissolved and you’re left with a syrup. Leave to cool.\nOnce your cheesecake is set, carefully slide the cake onto a serving plate. Decorate with slices of lemon and lime (or you could use curls of zest, if you prefer) then drizzle the gin and tonic syrup over the top. Allow the cheesecake to come to room temperature before serving.\n\nhttps://www.craftginclub.co.uk/ginnedmagazine/easy-gin-and-tonic-lemon-lime-cheesecake-recipe")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "This easy no-bake cheesecake is subtly infused with the flavours of gin, lemon and lime for a deliciously zesty and zingy dessert. \nThe beauty of this pudding (apart from being so easy to make - you don’t even need to turn the oven on!) is that you can make it the day before you want to eat it, and it keeps in the fridge for second servings the next day - if it lasts that long, of course!\nThis fabulous gin and tonic cheesecake is best served with - what else? - a refreshing G&T, in the garden, in the sunshine. Bliss!",
				Ingredients: []string{
					"For the base:",
					"200g digestive biscuits",
					"100g butter",
					"For the filling:",
					"500g cream cheese",
					"100g icing sugar",
					"250ml double cream",
					"50ml gin",
					"Zest and juice of one lime or lemon, or half of each",
					"For the topping:",
					"30ml gin",
					"100ml tonic",
					"2 tbsp sugar",
					"1tbsp freshly squeezed lemon juice",
					"Fresh lemon and/or lime slices, to decorate",
				},
				Instructions: []string{
					"To make the base, blitz the digestive biscuits in a food processor or crush them in a sandwich bag with a rolling pin until you have fine crumbs. Melt the butter and combine thoroughly with the biscuit. You want this mixture to be squidgy enough that it will form a base, but not overly wet - if you need to add more butter or biscuit crumbs to get the right texture, play around with the quantities.",
					"Butter and line a 23cm loose-bottomed tin then press the buttery crumb mixture firmly down into the base to create an even layer. Chill in the fridge for an hour or until it is firmly set.",
					"Meanwhile, mix the cream cheese, icing sugar, gin, lemon and lime zest and juice and cream together in a bowl until completely combined. Spoon the cream mixture onto the biscuit base and smooth the surface with a flat knife or spatula. Leave to set in the fridge for at least six hours or overnight.",
					"To make the topping, gently heat the gin and tonic with the lemon juice and sugar and stir gently until all the sugar has dissolved and you’re left with a syrup. Leave to cool.",
					"Once your cheesecake is set, carefully slide the cake onto a serving plate. Decorate with slices of lemon and lime (or you could use curls of zest, if you prefer) then drizzle the gin and tonic syrup over the top. Allow the cheesecake to come to room temperature before serving.",
				},
				Keywords: make([]string, 0),
				Name:     "Gin and tonic lemon & lime cheesecake",
				Tools:    make([]string, 0),
				URL:      "https://www.craftginclub.co.uk/ginnedmagazine/easy-gin-and-tonic-lemon-lime-cheesecake-recipe",
				Yield:    1,
			},
		},
		{
			name: "recipe 9",
			buf:  bytes.NewBuffer([]byte("Heston’s roast potatoes\n\n“In the end isn’t a roast all about the roast potatoes? While we all love a roast, I reckon that deep down we all have this secret shared belief that the best bit of any roast is the roast potato. Of course, there’s a majesty to a tender, juicy, aromatic bird with a lovely browned skin. Of course, a big-flavoured gravy makes a difference. But no matter how delicious everything else is, if the roasties aren’t quite there, then we’re not quite satisfied. That’s probably part of the reason I got a little bit obsessed about getting roast potatoes with exactly the texture I wanted: crispy and crunchy on the outside; soft and fluffy on the inside. So, if that’s how you like your potatoes, you’re in for a treat.”\n\nServes 4-6\nHands-on time 20 min\nOven time 1 hour\nSimmering time 25 min\n\nIngredients\n2kg maris piper potatoes\nVegetable oil or melted goose fat, duck fat or lard to roast\n3 thyme sprigs\n6 garlic cloves\n\nMethod\nPeel the potatoes and cut them into even-size large chunks [1]. Immerse the potatoes in a bowl of cold water as you prepare them to prevent browning. Rinse the potatoes in a colander under cold running water until the water runs clear, to remove excess starch.\nFill a large pan with lightly salted water and add the potatoes. Bring to the boil, then reduce the heat to a simmer and cook for 25 minutes until tender. Drain them very well in a colander, then spread out on a large wire rack set over a tray. Allow to cool.\nIn the meantime, heat the oven to 180°C fan/gas 6. Select a large roasting tray, big enough to take all the potatoes and spread out in a single layer. Add enough oil (or melted fat) to the roasting tray to create a shallow layer, about 5mm deep.\nLay the potatoes out carefully in the roasting tray, including all the smaller, broken-up pieces (they make delicious ultra-crispy bits), and roast in the oven for 20 minutes.\nTurn the potatoes a little and roast them for an additional 20 minutes or until firmed up and lightly golden on all sides. Turn the potatoes once more, especially the sides that look like they may need more time.\nScatter over the thyme and use the back of a knife to smash the garlic cloves. Add these to the tray and return to the oven for another 20 minutes. The potatoes will be golden and crispy. Season with salt and serve immediately, to retain their crispiness.\n\nDelicious Tips\nFloury potatoes have the edge over waxy ones for roasties, as the post-simmer fluffy texture catches fat and creates crunchiness. Ideally, you’re looking for a variety that won’t fall to pieces too easily. Maris piper is a good all-rounder for this.\nYou could cut the potatoes into thirds or quarters, depending on their size. Those pointy corners are good, anyway, as they can catch fat and create crispness.\nAfter simmering you’ll notice that the potatoes have not only softened but also look a little translucent. Some may show cracks or may have broken up a little. This is a key part of the process: those cracks are where fat can collect, which is what creates that crunchy exterior.\nOvercrowding the roasting tray can hinder browning and crisping, so aim for a bit of space between each spud.\nAfter their second spell in the oven, take a good look at your potatoes. Are they firm with a harder skin and golden all over? If not, you’ll need to return them to the oven until they are, before the next stage.\n\nhttps://www.deliciousmagazine.co.uk/recipes/hestons-roast-potatoes/\n\n\n220 grader\nLegg i potetene, bruk en skje og hell det varme fettet over alle potetene.\nSett i ovnen, 220 grader i 20 minutter.\nReduser til 180 grader, snu potetbitene og sett tilbake i ovnen i ca 20 minutter.")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "“In the end isn’t a roast all about the roast potatoes? While we all love a roast, I reckon that deep down we all have this secret shared belief that the best bit of any roast is the roast potato. Of course, there’s a majesty to a tender, juicy, aromatic bird with a lovely browned skin. Of course, a big-flavoured gravy makes a difference. But no matter how delicious everything else is, if the roasties aren’t quite there, then we’re not quite satisfied. That’s probably part of the reason I got a little bit obsessed about getting roast potatoes with exactly the texture I wanted: crispy and crunchy on the outside; soft and fluffy on the inside. So, if that’s how you like your potatoes, you’re in for a treat.”",
				Ingredients: []string{
					"2kg maris piper potatoes",
					"Vegetable oil or melted goose fat, duck fat or lard to roast",
					"3 thyme sprigs",
					"6 garlic cloves",
				},
				Instructions: []string{
					"Peel the potatoes and cut them into even-size large chunks [1]. Immerse the potatoes in a bowl of cold water as you prepare them to prevent browning. Rinse the potatoes in a colander under cold running water until the water runs clear, to remove excess starch.",
					"Fill a large pan with lightly salted water and add the potatoes. Bring to the boil, then reduce the heat to a simmer and cook for 25 minutes until tender. Drain them very well in a colander, then spread out on a large wire rack set over a tray. Allow to cool.",
					"In the meantime, heat the oven to 180°C fan/gas 6. Select a large roasting tray, big enough to take all the potatoes and spread out in a single layer. Add enough oil (or melted fat) to the roasting tray to create a shallow layer, about 5mm deep.",
					"Lay the potatoes out carefully in the roasting tray, including all the smaller, broken-up pieces (they make delicious ultra-crispy bits), and roast in the oven for 20 minutes.",
					"Turn the potatoes a little and roast them for an additional 20 minutes or until firmed up and lightly golden on all sides. Turn the potatoes once more, especially the sides that look like they may need more time.",
					"Scatter over the thyme and use the back of a knife to smash the garlic cloves. Add these to the tray and return to the oven for another 20 minutes. The potatoes will be golden and crispy. Season with salt and serve immediately, to retain their crispiness.",
					"Delicious Tips",
					"Floury potatoes have the edge over waxy ones for roasties, as the post-simmer fluffy texture catches fat and creates crunchiness. Ideally, you’re looking for a variety that won’t fall to pieces too easily. Maris piper is a good all-rounder for this.",
					"You could cut the potatoes into thirds or quarters, depending on their size. Those pointy corners are good, anyway, as they can catch fat and create crispness.",
					"After simmering you’ll notice that the potatoes have not only softened but also look a little translucent. Some may show cracks or may have broken up a little. This is a key part of the process: those cracks are where fat can collect, which is what creates that crunchy exterior.",
					"Overcrowding the roasting tray can hinder browning and crisping, so aim for a bit of space between each spud.",
					"After their second spell in the oven, take a good look at your potatoes. Are they firm with a harder skin and golden all over? If not, you’ll need to return them to the oven until they are, before the next stage.",
					"220 grader",
					"Legg i potetene, bruk en skje og hell det varme fettet over alle potetene.",
					"Sett i ovnen, 220 grader i 20 minutter.",
					"Reduser til 180 grader, snu potetbitene og sett tilbake i ovnen i ca 20 minutter.",
				},
				Keywords: make([]string, 0),
				Name:     "Heston’s roast potatoes",
				Times: models.Times{
					Prep: 20 * time.Minute,
					Cook: 1*time.Hour + 25*time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "https://www.deliciousmagazine.co.uk/recipes/hestons-roast-potatoes/",
				Yield: 4,
			},
		},
		{
			name: "recipe 10",
			buf:  bytes.NewBuffer([]byte("Mango Ice Cream in a Blender! ไอศรีมมะม่วง ง่ายสุดๆ\n\nA ridiculously simple ice cream recipe that is also ridiculously delicious! This mango ice cream is done in a blitz (literally!) using only a blender and 3-4 simple ingredients. Trust me, you’ll be amazed, as I was, at how good this ice cream is given how easy and quick it is to make! Now every time good sweet mango goes on sale I buy a whole bunch and freeze them just for this occasion! Just one note: mango makes up a very large proportion of this recipe so it’s really important that you use good, ripe mango for this! Enjoy!\n\nServes 3-4\n\nINGREDIENTS\n300 g sweet, ripe mango, cut into small cubes and freeze, you will need about 1½ – 2 mangoes (see note)\n¼ cup (60 mL) honey (may need more if your mango isn’t very sweet)\n½ cup (120 mL) Greek yogurt, plain, full-fat (about 10% fat) (if you can’t find Greek yogurt, see note below on how to make your own)\n1-3 tsp lime juice (amount depends on how tart your mango is, and you may not need any at all if the mango is sour)\nNote: When freezing mango, freeze on a plate or tray lined with parchment paper or plastic wrap (so the mango won’t stick to the tray) and spread the mango cubes out so they are not touching each other (so they won’t stick to each other). If your tray isn’t big enough, you can stack the mango in layers, with plastic wrap or parchment paper in between each layer.\n\nTo make your own Greek yogurt from regular, plain, full-fat yogurt: Line a mesh sieve with two layers of cheesecloth or coffee filter, and place it over a bowl. Place at least 1 cup of yogurt in the sieve and let it sit in the fridge for 2 hour, and you will notice a lot of watery liquid (the whey) will be strained out and your yogurt will not be very thick and creamy.\n\nINSTRUCTIONS\nPlace a metal cake pan that you will use for storing the ice cream into the freezer before you start prepping. I suggest using a metal cake pan because it conducts heat faster and will help the ice cream firm up faster.  Note: Faster freezing means less icy and more creamy ice cream.\nPut the greek yogurt into the blender first, followed by the lime juice (if using) and honey. Tip: Try to place the honey into the middle of the blender jug and avoid getting it onto the sides because the part that sticks to the blender may not blend in.  \nAdd the mango and blend on highest speed (I use the pulse function). It will likely be stuck in the beginning, so you have to keep pushing the mixture down to the blade, using either the blender tamper or a wooden spoon, being careful not to let the spoon touch the blade! After a little pushing in the beginning, the mixture will eventually blend smoothly without getting stuck, and as soon as the mixture looks smooth, it is done. Do not over blend, you don’t want to melt the mixture any more than necessary.\nYou can either serve it right out of the blender as a soft-serve style ice cream, or place it into the metal cake pan and freeze for at least 1 hour. If you’ve let it freeze for several hours, you may need to let it sit at room temp for 5 minutes to soften slightly to make scooping easier. For long term storage, cover the container with plastic wrap and then another layer of aluminum foil.\n\nhttps://hot-thai-kitchen.com/mango-ice-cream/")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "A ridiculously simple ice cream recipe that is also ridiculously delicious! This mango ice cream is done in a blitz (literally!) using only a blender and 3-4 simple ingredients. Trust me, you’ll be amazed, as I was, at how good this ice cream is given how easy and quick it is to make! Now every time good sweet mango goes on sale I buy a whole bunch and freeze them just for this occasion! Just one note: mango makes up a very large proportion of this recipe so it’s really important that you use good, ripe mango for this! Enjoy!",
				Ingredients: []string{
					"300 g sweet, ripe mango, cut into small cubes and freeze, you will need about 1½ – 2 mangoes (see note)",
					"¼ cup (60 mL) honey (may need more if your mango isn’t very sweet)",
					"½ cup (120 mL) Greek yogurt, plain, full-fat (about 10% fat) (if you can’t find Greek yogurt, see note below on how to make your own)",
					"1-3 tsp lime juice (amount depends on how tart your mango is, and you may not need any at all if the mango is sour)",
					"Note: When freezing mango, freeze on a plate or tray lined with parchment paper or plastic wrap (so the mango won’t stick to the tray) and spread the mango cubes out so they are not touching each other (so they won’t stick to each other). If your tray isn’t big enough, you can stack the mango in layers, with plastic wrap or parchment paper in between each layer.",
				},
				Instructions: []string{
					"To make your own Greek yogurt from regular, plain, full-fat yogurt: Line a mesh sieve with two layers of cheesecloth or coffee filter, and place it over a bowl. Place at least 1 cup of yogurt in the sieve and let it sit in the fridge for 2 hour, and you will notice a lot of watery liquid (the whey) will be strained out and your yogurt will not be very thick and creamy.",
					"Place a metal cake pan that you will use for storing the ice cream into the freezer before you start prepping. I suggest using a metal cake pan because it conducts heat faster and will help the ice cream firm up faster.  Note: Faster freezing means less icy and more creamy ice cream.",
					"Put the greek yogurt into the blender first, followed by the lime juice (if using) and honey. Tip: Try to place the honey into the middle of the blender jug and avoid getting it onto the sides because the part that sticks to the blender may not blend in.",
					"Add the mango and blend on highest speed (I use the pulse function). It will likely be stuck in the beginning, so you have to keep pushing the mixture down to the blade, using either the blender tamper or a wooden spoon, being careful not to let the spoon touch the blade! After a little pushing in the beginning, the mixture will eventually blend smoothly without getting stuck, and as soon as the mixture looks smooth, it is done. Do not over blend, you don’t want to melt the mixture any more than necessary.",
					"You can either serve it right out of the blender as a soft-serve style ice cream, or place it into the metal cake pan and freeze for at least 1 hour. If you’ve let it freeze for several hours, you may need to let it sit at room temp for 5 minutes to soften slightly to make scooping easier. For long term storage, cover the container with plastic wrap and then another layer of aluminum foil.",
				},
				Keywords: make([]string, 0),
				Name:     "Mango Ice Cream in a Blender! ไอศรีมมะม่วง ง่ายสุดๆ",
				Tools:    make([]string, 0),
				URL:      "https://hot-thai-kitchen.com/mango-ice-cream/",
				Yield:    3,
			},
		},
		{
			name: "recipe 11",
			buf:  bytes.NewBuffer([]byte("Sai Grog - Thai Sausage\n\nMost tourists visiting Thailand remain entirely unaware of just how good Thai Sausages can be. The reason for this is that they seldom find their way onto a restaurant menu. Instead, they are usually purchased from a food market, or from a street vendor, to be eaten as a snack. This is truly a shame, as a good quality Thai Sausage can be very tasty.\nUsually the sausage will be served cut into thin slices, and eaten with a small amount of greens such as cabbage leaves, whole garlic cloves, whole chilli and cucumber. Sometimes the Thai Sausage will also be dipped into a sweet chilli sauce.\nMost Thai Sausages are made from pork, but they can also be made from beef, chicken, or even fish. Depending upon the region in which they were made, they tend to exhibit a slightly different taste. In Bangkok and Central Thailand they will be quite meaty tasting, with a lot of garlic present, in the North East (Isan), more chilli is added, and the quality of the meat, and quantity added, tends to be less, making for a more fatty sausage. In the South of Thailand, they tend to be much more sweat in taste, with sugar added to the recipe.\nWhat really separates a Thai Sausage from its western counterparts is the sheer variety of herbs and spices added to the sausage mix. Typically a Thai Sausage will contain galangal, lemon grass, garlic, coriander, chilli, kaffir lime leaves, white pepper and fish sauce. As we can see, this is quite a mixture of flavours, making the Thai Sausage something quite special.\nThe major difference between a Thai Sausage and Western Sausage is in the form of the filler content, which is mixed with the meat. In the West we tend to add things like bran, corn, or similar crops, whereas a Thai Sausage uses sticky rice instead. This makes for a very heavy, filling sausage, one of the reasons it is so popular as a snack food, as a couple of them can fill you up nicely.\nThai Sausages can be fried, but they tend to taste so much better when they have been barbequed over charcoal. They need to be cooked very slowly, to allow the flavours of the herbs and spices to mix with the meat fat as it cooks, spreading flavours into the whole sausage.\n\nPreparation Time: 10 minutes \nCooking Time: 45 minutes \nReady In: 55 minutes\nPortions: 4 people\n\nIngredients\n150g Minced Pork\n½ Cup Sticky Rice\n50 gms Galangal\n50 gms Lemongrass\n5 cloves Thai Garlic\n3 Thai Coriander Roots\n5 Thai Green Chillies\n4 Kaffir Lime Leaves\n½ tablespoon Thai White Pepper\n2 tablespoons Fish Sauce\n1 teaspoon MSG (optional)\n1 teaspoon Salt\n\nInstructions\nCook the sticky rice by soaking it in hot water for about an hour, and then steaming it in a Thai rice steamer for 10-15 minutes. Then let it cool down.\nPulverise the garlic using a pestle and mortar.\nMix all the ingredients including the minced pork thoroughly. Use a food processer if required. Leave the mixture at room temperature overnight.\nIf you are making traditional sausages, pipe it into the sausage skin. You can also make small meat patties/hamburgers with the mixture, if you don’t want to bother with making sausages.\nCan be either cooked in the oven or barbequed for about 20-25 minutes.\n\nhttps://www.thai-food-online.co.uk/pages/thai-sausage-recipe")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "Most tourists visiting Thailand remain entirely unaware of just how good Thai Sausages can be. The reason for this is that they seldom find their way onto a restaurant menu. Instead, they are usually purchased from a food market, or from a street vendor, to be eaten as a snack. This is truly a shame, as a good quality Thai Sausage can be very tasty.\nUsually the sausage will be served cut into thin slices, and eaten with a small amount of greens such as cabbage leaves, whole garlic cloves, whole chilli and cucumber. Sometimes the Thai Sausage will also be dipped into a sweet chilli sauce.\nMost Thai Sausages are made from pork, but they can also be made from beef, chicken, or even fish. Depending upon the region in which they were made, they tend to exhibit a slightly different taste. In Bangkok and Central Thailand they will be quite meaty tasting, with a lot of garlic present, in the North East (Isan), more chilli is added, and the quality of the meat, and quantity added, tends to be less, making for a more fatty sausage. In the South of Thailand, they tend to be much more sweat in taste, with sugar added to the recipe.\nWhat really separates a Thai Sausage from its western counterparts is the sheer variety of herbs and spices added to the sausage mix. Typically a Thai Sausage will contain galangal, lemon grass, garlic, coriander, chilli, kaffir lime leaves, white pepper and fish sauce. As we can see, this is quite a mixture of flavours, making the Thai Sausage something quite special.\nThe major difference between a Thai Sausage and Western Sausage is in the form of the filler content, which is mixed with the meat. In the West we tend to add things like bran, corn, or similar crops, whereas a Thai Sausage uses sticky rice instead. This makes for a very heavy, filling sausage, one of the reasons it is so popular as a snack food, as a couple of them can fill you up nicely.\nThai Sausages can be fried, but they tend to taste so much better when they have been barbequed over charcoal. They need to be cooked very slowly, to allow the flavours of the herbs and spices to mix with the meat fat as it cooks, spreading flavours into the whole sausage.",
				Ingredients: []string{
					"150g Minced Pork",
					"½ Cup Sticky Rice",
					"50 gms Galangal",
					"50 gms Lemongrass",
					"5 cloves Thai Garlic",
					"3 Thai Coriander Roots",
					"5 Thai Green Chillies",
					"4 Kaffir Lime Leaves",
					"½ tablespoon Thai White Pepper",
					"2 tablespoons Fish Sauce",
					"1 teaspoon MSG (optional)",
					"1 teaspoon Salt",
				},
				Instructions: []string{
					"Cook the sticky rice by soaking it in hot water for about an hour, and then steaming it in a Thai rice steamer for 10-15 minutes. Then let it cool down.",
					"Pulverise the garlic using a pestle and mortar.",
					"Mix all the ingredients including the minced pork thoroughly. Use a food processer if required. Leave the mixture at room temperature overnight.",
					"If you are making traditional sausages, pipe it into the sausage skin. You can also make small meat patties/hamburgers with the mixture, if you don’t want to bother with making sausages.",
					"Can be either cooked in the oven or barbequed for about 20-25 minutes.",
				},
				Keywords:  make([]string, 0),
				Name:      "Sai Grog - Thai Sausage",
				Nutrition: models.Nutrition{},
				Times: models.Times{
					Prep: 10 * time.Minute,
					Cook: 45 * time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "https://www.thai-food-online.co.uk/pages/thai-sausage-recipe",
				Yield: 4,
			},
		},
		{
			name: "recipe 12",
			buf:  bytes.NewBuffer([]byte("Thai Beef Stew เนื้อตุ๋น\n\nA flavourful beef stew that won’t leave you feeling heavy at the end of the meal! The super tender beef is stewed in a broth infused with lots of herbs and spices. You’ll learn about my favourite cut of beef for stew that I promise will become your favourite too. If you’ve got a slow cooker or a crockpot, this is the perfect dish for it!\n\nServes 4\n\nINGREDIENTS\n700 g (1½ lb) beef “digital muscle” (see note) or other stew-friendly beef such as shank, round or chuck, cut into 1-inch thick pieces.\n2 generous pinches of salt\n5 cups unsalted, plain beef stock (see note)\n2 Tbsp soy sauce\n2 Tbsp Golden Mountain sauce\n2 Tbsp oyster sauce\n1 Tbsp black soy sauce (confused by all these sauces? See this video!)\n3 Tbsp dark brown sugar, palm sugar is okay too\n2 pieces star anise\n2 sticks cinnamon\n1 tsp coriander seeds, toasted\n2 bay leaves\n10 slices galangal\n10 slices ginger\n½ tsp black peppercorns, cracked\n½ tsp white peppercorns, cracked\nOptional herbs and spices you can also add: lemongrass tops, white cardamom pods, dried goji berries, sichuan peppercorns\n1 onion, large dice\n8 cloves garlic, smashed\n2 dried shiitake mushrooms, optional\n4-inch piece daikon, cut into 1-inch slices (I didn’t use this in the video, but it’s a nice addition if you have it)\nChopped cilantro and green onions for serving\nOptional condiments: Chili vinegar (any kinds of chilies soaked in 5% white vinegar for at least 15 minutes before using). Fried garlic is great on this too!\n\nNotes: This cut of beef is officially called “super digital flexor muscle” which is a part of the cow’s hind legs. It is unique in the clear marbling of tendons all throughout the meat.\n\nTo make plain beef stock for this recipe, simply simmer 1½ lb of beef bones in about 8 cups of water for 3 hours, skimming off all the scum that floats to the top, adding more water as needed to keep the bones submerged. If you end up with less than 5 cups of stock, simply add more liquid to make up the shortfall. I highly recommend homemade beef stock so that you don’t get extra flavourings and salt often added to commercial beef stock. But if you’re gonna buy it, make sure it is UNSALTED! \n\nINSTRUCTIONS\nSprinkle salt over the meat on both sides. In a pot, add a little bit of oil to coat the bottom and heat over medium high heat. Add the meat and sear without stirring until well-browned. Do not crowd the pan; you will have to do this in batches. Flip and sear the other side.\nOnce you’re done, pour off excess oil but leave the browned bits in the pot; however, if you’ve burned these bits and they’ve turned black, scrub them off before proceeding. Browned bits give a nice flavour, but burned bits do not!\nReturn all beef to the pot, cover with the beef stock and add soy sauce, Golden Mountain sauce, oyster sauce, black soy sauce, and brown sugar. Stir to mix and bring to a simmer.\nWhile the beef is coming to a simmer, make your spice bag: In a muslin soup bag or on a square piece of cheesecloth, place all the spices and herbs except onions, garlic, shiitake mushrooms, and daikon. Close the bag, or if using cheesecloth tie the corners together to secure the spices inside.\nIf you notice any foam on top of the stew at this point, skim it off before adding the spice bag. Add the spice bag and push it in until it is submerged. Add onion, garlic, shiitake mushrooms and daikon (if using), then cover the pot and simmer over low heat for 3 hours or until the meat is fork tender.\nAfter 3 hours, remove the spice bag and taste and adjust seasoning with more salt or sugar as needed. Serve with rice or pour over noodles to make a delicious noodle soup! Top with chopped cilantro and onions, and drizzle with a little chili vinegar if desired. Enjoy!\n\nhttp://hot-thai-kitchen.com/thai-beef-stew/")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "A flavourful beef stew that won’t leave you feeling heavy at the end of the meal! The super tender beef is stewed in a broth infused with lots of herbs and spices. You’ll learn about my favourite cut of beef for stew that I promise will become your favourite too. If you’ve got a slow cooker or a crockpot, this is the perfect dish for it!",
				Ingredients: []string{
					"700 g (1½ lb) beef “digital muscle” (see note) or other stew-friendly beef such as shank, round or chuck, cut into 1-inch thick pieces.",
					"2 generous pinches of salt",
					"5 cups unsalted, plain beef stock (see note)",
					"2 Tbsp soy sauce",
					"2 Tbsp Golden Mountain sauce",
					"2 Tbsp oyster sauce",
					"1 Tbsp black soy sauce (confused by all these sauces? See this video!)",
					"3 Tbsp dark brown sugar, palm sugar is okay too",
					"2 pieces star anise",
					"2 sticks cinnamon",
					"1 tsp coriander seeds, toasted",
					"2 bay leaves",
					"10 slices galangal",
					"10 slices ginger",
					"½ tsp black peppercorns, cracked",
					"½ tsp white peppercorns, cracked",
					"Optional herbs and spices you can also add: lemongrass tops, white cardamom pods, dried goji berries, sichuan peppercorns",
					"1 onion, large dice",
					"8 cloves garlic, smashed",
					"2 dried shiitake mushrooms, optional",
					"4-inch piece daikon, cut into 1-inch slices (I didn’t use this in the video, but it’s a nice addition if you have it)",
					"Chopped cilantro and green onions for serving",
					"Optional condiments: Chili vinegar (any kinds of chilies soaked in 5% white vinegar for at least 15 minutes before using). Fried garlic is great on this too!",
				},
				Instructions: []string{
					"Notes: This cut of beef is officially called “super digital flexor muscle” which is a part of the cow’s hind legs. It is unique in the clear marbling of tendons all throughout the meat.",
					"To make plain beef stock for this recipe, simply simmer 1½ lb of beef bones in about 8 cups of water for 3 hours, skimming off all the scum that floats to the top, adding more water as needed to keep the bones submerged. If you end up with less than 5 cups of stock, simply add more liquid to make up the shortfall. I highly recommend homemade beef stock so that you don’t get extra flavourings and salt often added to commercial beef stock. But if you’re gonna buy it, make sure it is UNSALTED!",
					"Sprinkle salt over the meat on both sides. In a pot, add a little bit of oil to coat the bottom and heat over medium high heat. Add the meat and sear without stirring until well-browned. Do not crowd the pan; you will have to do this in batches. Flip and sear the other side.",
					"Once you’re done, pour off excess oil but leave the browned bits in the pot; however, if you’ve burned these bits and they’ve turned black, scrub them off before proceeding. Browned bits give a nice flavour, but burned bits do not!",
					"Return all beef to the pot, cover with the beef stock and add soy sauce, Golden Mountain sauce, oyster sauce, black soy sauce, and brown sugar. Stir to mix and bring to a simmer.",
					"While the beef is coming to a simmer, make your spice bag: In a muslin soup bag or on a square piece of cheesecloth, place all the spices and herbs except onions, garlic, shiitake mushrooms, and daikon. Close the bag, or if using cheesecloth tie the corners together to secure the spices inside.",
					"If you notice any foam on top of the stew at this point, skim it off before adding the spice bag. Add the spice bag and push it in until it is submerged. Add onion, garlic, shiitake mushrooms and daikon (if using), then cover the pot and simmer over low heat for 3 hours or until the meat is fork tender.",
					"After 3 hours, remove the spice bag and taste and adjust seasoning with more salt or sugar as needed. Serve with rice or pour over noodles to make a delicious noodle soup! Top with chopped cilantro and onions, and drizzle with a little chili vinegar if desired. Enjoy!",
				},
				Keywords: make([]string, 0),
				Name:     "Thai Beef Stew เนื้อตุ๋น",
				Tools:    make([]string, 0),
				URL:      "http://hot-thai-kitchen.com/thai-beef-stew/",
				Yield:    4,
			},
		},
		{
			name: "recipe 13",
			buf:  bytes.NewBuffer([]byte("Thai Tea Ice Cream (No Machine Method) ไอติมชาเย็น\n\nMakes 1 quart\n\nINGREDIENTS\n1 1/2 cup whipping cream\n1/4 cup Thai tea leaves\nA pinch of salt\n1/2 can sweetened condensed milk (150 ml)\n1.5 Tbsp Irish cream liqueur such as Bailey’s (optional)\nOptional Topping\n3 Tbsp sweetened condensed milk\n3 Tbsp evaporated milk\n\nINSTRUCTIONS\nHeat whipping cream in a small pot over medium heat, stirring occasionally, until the cream is steaming. Add tea leaves and a small pinch of salt, and stir just until the cream boils. Remove from heat and steep for 5 minutes (don’t leave it sitting until it cools down or it’ll be really hard to strain!). Strain the cream using a fine mesh strainer into a 2-cup measuring cup*, pressing out as much liquid as you can—you should have about 1 1/4 cup of cream, if you have too little, add more fresh cream until you have 1 1/4 cup; if you have a little more, don’t worry about it. Once the cream is cool enough to go into the fridge, refrigerate for at least 4 hours or until completely cold; I like to do this step 1 day in advance.\n*Note: If you don’t have a 2-cup measure, you can strain it into a 1-cup measure and another ¼ cup measure. Just make sure you don’t let the first cup overflow!\nWhile the cream is cooling, make the topping by stirring the evaporated milk and condensed milk together. Refrigerate until ready to use.\nOnce the cream is chilled completely, whip the cream to soft-peak stage using a stand mixer with a whisk attachment or a hand mixer. (When the beaters start leaving a trail that doesn’t immediately disappear, you’re at soft peaks.) Add condensed milk and the Irish cream liqueur and continue whipping until stiff peaks (when you lift your whisk, the peak that forms maintains its shape).\nTransfer the ice cream into a metal container and freeze for 2-3 hours or until solid. Alternatively, put the ice cream in between your favourite cookies to make an ice cream sandwich!\n\nhttp://hot-thai-kitchen.com/thai-tea-ice-cream/")),
			want: models.Recipe{
				Category: "uncategorized",
				Ingredients: []string{
					"1 1/2 cup whipping cream",
					"1/4 cup Thai tea leaves",
					"A pinch of salt",
					"1/2 can sweetened condensed milk (150 ml)",
					"1.5 Tbsp Irish cream liqueur such as Bailey’s (optional)",
					"Optional Topping",
					"3 Tbsp sweetened condensed milk",
					"3 Tbsp evaporated milk",
				},
				Instructions: []string{
					"Heat whipping cream in a small pot over medium heat, stirring occasionally, until the cream is steaming. Add tea leaves and a small pinch of salt, and stir just until the cream boils. Remove from heat and steep for 5 minutes (don’t leave it sitting until it cools down or it’ll be really hard to strain!). Strain the cream using a fine mesh strainer into a 2-cup measuring cup*, pressing out as much liquid as you can—you should have about 1 1/4 cup of cream, if you have too little, add more fresh cream until you have 1 1/4 cup; if you have a little more, don’t worry about it. Once the cream is cool enough to go into the fridge, refrigerate for at least 4 hours or until completely cold; I like to do this step 1 day in advance.",
					"*Note: If you don’t have a 2-cup measure, you can strain it into a 1-cup measure and another ¼ cup measure. Just make sure you don’t let the first cup overflow!",
					"While the cream is cooling, make the topping by stirring the evaporated milk and condensed milk together. Refrigerate until ready to use.",
					"Once the cream is chilled completely, whip the cream to soft-peak stage using a stand mixer with a whisk attachment or a hand mixer. (When the beaters start leaving a trail that doesn’t immediately disappear, you’re at soft peaks.) Add condensed milk and the Irish cream liqueur and continue whipping until stiff peaks (when you lift your whisk, the peak that forms maintains its shape).",
					"Transfer the ice cream into a metal container and freeze for 2-3 hours or until solid. Alternatively, put the ice cream in between your favourite cookies to make an ice cream sandwich!",
				},
				Keywords: make([]string, 0),
				Name:     "Thai Tea Ice Cream (No Machine Method) ไอติมชาเย็น",
				Tools:    make([]string, 0),
				URL:      "http://hot-thai-kitchen.com/thai-tea-ice-cream/",
				Yield:    1,
			},
		},
		{
			name: "recipe 14",
			buf:  bytes.NewBuffer([]byte("Ultimate roast chicken\n\nChef Margot Henderson shares her recipe for roast chicken. She says: 'A simple roast chicken is one of my all-time favourite dishes. If it comes wrapped in plastic, it’s best to take it out as soon as possible and cover it in paper instead. I like the bird to be dry as possible, for the crispiest skin'\n\nServes: 4-5\nPrep Time: 25 mins\nTotal Time: 1 hr 30 mins, plus resting\n\nIngredients\n1 Taste the Difference free-range whole chicken, about 1.7kg\n1 garlic bulb, unpeeled, halved horizontally\n1 lemon, halved, plus a squeeze for the gravy\na handful of thyme sprigs\n5 fresh bay leaves\n125g soft unsalted butter\nabout 2 tbsp extra-virgin olive oil\n100ml white wine\n3 tbsp plain flour\nabout 400ml fresh chicken stock\n\nStep by step\nPreheat the oven to 200°C, fan 180°C, gas 6. Fill the chicken cavity with half a bulb of garlic, half a lemon and the herbs. Season all over with flaky salt and black pepper.\nPack the butter under the skin on the breast and legs. Place the chicken in a roasting tin, drizzle the oil over the skin, squeeze over the other lemon half and maybe season a bit more. Tuck the other half of the garlic bulb in, too.\nRoast for 15 minutes, then turn the oven down to 180°C, fan 160°C, gas 4 and roast for a further 45 minutes. Check the bird is cooked by spearing between the leg and the breast with a skewer and if any pinkness appears, carry on cooking. Once cooked, leave to rest for 15 minutes; prop the bird up to let the juices run out.\nPour the juices from the roasting tin into a jug and skim off as much of the fat as is practical. Pour the white wine into the roasting tin and stir to release all the lovely sticky bits from the bottom. Scrape into a saucepan and add the roasting juices. Whisk in the plain flour and cook for a minute, then gradually whisk in the chicken stock to your preferred thickness. Simmer for 5-10 minutes. Add a squeeze of lemon juice and season to taste.\nTo serve, joint the bird into about 10 pieces and place on a platter, then pour the gravy over the meat.\n\nhttps://www.sainsburysmagazine.co.uk/recipes/mains/ultimate-roast-chicken")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "Chef Margot Henderson shares her recipe for roast chicken. She says: 'A simple roast chicken is one of my all-time favourite dishes. If it comes wrapped in plastic, it’s best to take it out as soon as possible and cover it in paper instead. I like the bird to be dry as possible, for the crispiest skin'",
				Ingredients: []string{
					"1 Taste the Difference free-range whole chicken, about 1.7kg",
					"1 garlic bulb, unpeeled, halved horizontally",
					"1 lemon, halved, plus a squeeze for the gravy",
					"a handful of thyme sprigs",
					"5 fresh bay leaves",
					"125g soft unsalted butter",
					"about 2 tbsp extra-virgin olive oil",
					"100ml white wine",
					"3 tbsp plain flour",
					"about 400ml fresh chicken stock",
				},
				Instructions: []string{
					"Preheat the oven to 200°C, fan 180°C, gas 6. Fill the chicken cavity with half a bulb of garlic, half a lemon and the herbs. Season all over with flaky salt and black pepper.",
					"Pack the butter under the skin on the breast and legs. Place the chicken in a roasting tin, drizzle the oil over the skin, squeeze over the other lemon half and maybe season a bit more. Tuck the other half of the garlic bulb in, too.",
					"Roast for 15 minutes, then turn the oven down to 180°C, fan 160°C, gas 4 and roast for a further 45 minutes. Check the bird is cooked by spearing between the leg and the breast with a skewer and if any pinkness appears, carry on cooking. Once cooked, leave to rest for 15 minutes; prop the bird up to let the juices run out.",
					"Pour the juices from the roasting tin into a jug and skim off as much of the fat as is practical. Pour the white wine into the roasting tin and stir to release all the lovely sticky bits from the bottom. Scrape into a saucepan and add the roasting juices. Whisk in the plain flour and cook for a minute, then gradually whisk in the chicken stock to your preferred thickness. Simmer for 5-10 minutes. Add a squeeze of lemon juice and season to taste.",
					"To serve, joint the bird into about 10 pieces and place on a platter, then pour the gravy over the meat.",
				},
				Keywords: make([]string, 0),
				Name:     "Ultimate roast chicken",
				Times: models.Times{
					Prep: 25 * time.Minute,
					Cook: 1*time.Hour + 5*time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "https://www.sainsburysmagazine.co.uk/recipes/mains/ultimate-roast-chicken",
				Yield: 4,
			},
		},
		{
			name: "recipe 15",
			buf:  bytes.NewBuffer([]byte("Vin Santo Ice Cream With Cantuccini By Nigella\n\nOne of the loveliest puddings, if it quite counts as that, to order in restaurants in Tuscany, is a glass of vin santo, that resinous, intense, amber-coloured holy wine, with a few almond-studded biscuits to dunk in. The idea for this comes purely from that: the ice cream is further deepened by the addition of treacly muscovado sugar, and the wine in it keeps it voluptuously velvety, even after being frozen. You don't absolutely need to serve the cantuccini biscuits with it, but the combination is pretty well unbeatable.\n\nServes: 8-10\n\nINGREDIENTS\n584 millilitres double cream\n200 millilitres vin santo\n8 egg yolks\n6 tablespoons light brown muscovado sugar\n1 packet cantuccini biscuits\n\nMETHOD\nHeat the double cream in one pan, the vin santo in another. Whisk the yolks and muscovado sugar together and, still whisking, pour first the hot vin santo and then the warmed cream into them. Pour this mixture into a good-sized pan and cook till a velvety custard. I don't bother with a double boiler, and actually don't even keep the heat very low, but you will need to stir constantly, and if you think there's any trouble ahead, plunge the into a sink half filled with water and whisk like mad. It shouldn't take more than 10 minutes, this way, for the custard to cook.\nWhen it has thickened, take it off the heat and allow to cool. Then chill and freeze in an ice-cream maker or put the cooled custard into a covered container, stick it in the freezer and whip it out every hour for 3 hours as it freezes and give it a good beating, either with an electric whisk, by hand or in the processor. That gets rid of any ice crystals that form and that make the ice cream crunchy rather than smooth.\nServe by scooping out into small wine glasses, giving everyone some cantuccini to dip into the deep, almost incense-intense ice cream.\n\nhttps://www.nigella.com/recipes/vin-santo-ice-cream-with-cantuccini")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "One of the loveliest puddings, if it quite counts as that, to order in restaurants in Tuscany, is a glass of vin santo, that resinous, intense, amber-coloured holy wine, with a few almond-studded biscuits to dunk in. The idea for this comes purely from that: the ice cream is further deepened by the addition of treacly muscovado sugar, and the wine in it keeps it voluptuously velvety, even after being frozen. You don't absolutely need to serve the cantuccini biscuits with it, but the combination is pretty well unbeatable.",
				Ingredients: []string{
					"584 millilitres double cream",
					"200 millilitres vin santo",
					"8 egg yolks",
					"6 tablespoons light brown muscovado sugar",
					"1 packet cantuccini biscuits",
				},
				Instructions: []string{
					"Heat the double cream in one pan, the vin santo in another. Whisk the yolks and muscovado sugar together and, still whisking, pour first the hot vin santo and then the warmed cream into them. Pour this mixture into a good-sized pan and cook till a velvety custard. I don't bother with a double boiler, and actually don't even keep the heat very low, but you will need to stir constantly, and if you think there's any trouble ahead, plunge the into a sink half filled with water and whisk like mad. It shouldn't take more than 10 minutes, this way, for the custard to cook.",
					"When it has thickened, take it off the heat and allow to cool. Then chill and freeze in an ice-cream maker or put the cooled custard into a covered container, stick it in the freezer and whip it out every hour for 3 hours as it freezes and give it a good beating, either with an electric whisk, by hand or in the processor. That gets rid of any ice crystals that form and that make the ice cream crunchy rather than smooth.",
					"Serve by scooping out into small wine glasses, giving everyone some cantuccini to dip into the deep, almost incense-intense ice cream.",
				},
				Keywords: make([]string, 0),
				Name:     "Vin Santo Ice Cream With Cantuccini By Nigella",
				Tools:    make([]string, 0),
				URL:      "https://www.nigella.com/recipes/vin-santo-ice-cream-with-cantuccini",
				Yield:    8,
			},
		},
		{
			name: "recipe 16",
			buf:  bytes.NewBuffer([]byte("Bestefars bankekjøtt\n\nNår kjøttet skal koke lenge blir det en typisk søndagsmiddag. Lise Finckenhagen serverer bestefars bankekjøtt med gulrøtter, purre og potetstappe.\n\nEnkel\nType: Middag\nHovedingrediens: Storfe\nKarakteristika: Gryte\n4 porsjoner\n\nIngredienser\n800 - 1000 g bankekjøtt, flatbiff eller rundstek\n4 løk\n6 ss mel\n4-6 dl oksekraft (ev. buljong)\n2 laurbærblader\nsalt og kvernet sort pepper\nsmør til steking\nhakket persille til pynt\n\nSlik gjør du\nSkjær kjøttet i 1-2 cm tykke skiver på tvers av fibrene i kjøttet. Bruk en kjøttbanker eller lignende og bank kjøttet lett, det skal bare bankes og ikke moses. Krydre kjøttet med salt og kvernet sort pepper.\nSkrell og del løken i båter.\nVend kjøttskivene i mel og brun dem på begge sider i godt med smør i en varm stekepanne. Det er lov å være litt raus med smøret. Bruk helst en jernpanne hvis du har det. Stek i flere omganger og legg kjøttet over i en passende gryte etter hvert som det er ferdig brunet. Kok gjerne ut pannen med en skvett vann (eller litt av kraften) etter hver bruning og hell det over kjøttet i gryta.\nFortsett med løken når alt kjøttet er ferdig brunet. Stek løken i godt med smør til løken blir gyllen og begynner å mykne. Hell all løken over kjøttet.\nTilsett laurbærblad og hell i kraft til det akkurat dekker kjøttet. La det hele surre/småputre under lokk i ca. 2 ½ time, eller til kjøttet er helt mørt. Juster smaken med salt og kvernet pepper.\nServeres for eksempel med glaserte gulrøtter og purre, og mandelpotetstappe smakt til med revet pepperrot.\n\nTips\nErstatt 2 dl av kraften med mørkt øl, eller smak sausen til med en klunk Madeira.\n\nhttps://www.nrk.no/mat/bestefars-bankekjott-1.12015168")),
			want: models.Recipe{
				Category:    "middag",
				Description: "Når kjøttet skal koke lenge blir det en typisk søndagsmiddag. Lise Finckenhagen serverer bestefars bankekjøtt med gulrøtter, purre og potetstappe.",
				Ingredients: []string{
					"800 - 1000 g bankekjøtt, flatbiff eller rundstek",
					"4 løk",
					"6 ss mel",
					"4-6 dl oksekraft (ev. buljong)",
					"2 laurbærblader",
					"salt og kvernet sort pepper",
					"smør til steking",
					"hakket persille til pynt",
				},
				Instructions: []string{
					"Skjær kjøttet i 1-2 cm tykke skiver på tvers av fibrene i kjøttet. Bruk en kjøttbanker eller lignende og bank kjøttet lett, det skal bare bankes og ikke moses. Krydre kjøttet med salt og kvernet sort pepper.",
					"Skrell og del løken i båter.",
					"Vend kjøttskivene i mel og brun dem på begge sider i godt med smør i en varm stekepanne. Det er lov å være litt raus med smøret. Bruk helst en jernpanne hvis du har det. Stek i flere omganger og legg kjøttet over i en passende gryte etter hvert som det er ferdig brunet. Kok gjerne ut pannen med en skvett vann (eller litt av kraften) etter hver bruning og hell det over kjøttet i gryta.",
					"Fortsett med løken når alt kjøttet er ferdig brunet. Stek løken i godt med smør til løken blir gyllen og begynner å mykne. Hell all løken over kjøttet.",
					"Tilsett laurbærblad og hell i kraft til det akkurat dekker kjøttet. La det hele surre/småputre under lokk i ca. 2 ½ time, eller til kjøttet er helt mørt. Juster smaken med salt og kvernet pepper.",
					"Serveres for eksempel med glaserte gulrøtter og purre, og mandelpotetstappe smakt til med revet pepperrot.",
					"Tips",
					"Erstatt 2 dl av kraften med mørkt øl, eller smak sausen til med en klunk Madeira.",
				},
				Keywords: []string{"storfe", "gryte"},
				Name:     "Bestefars bankekjøtt",
				Tools:    make([]string, 0),
				URL:      "https://www.nrk.no/mat/bestefars-bankekjott-1.12015168",
				Yield:    4,
			},
		},
		{
			name: "recipe 17",
			buf:  bytes.NewBuffer([]byte("Festkake med lemoncurd og knust ananas\n\nEn frisk og god kake og en fryd for øyet!\n\nPersoner: 10\nTid: 40-60 min\nVanskelighetsgrad: middels\n\nIngredienser\nSukkerbrød\n7.5 stk egg\n250 g sukker\n250 g hvetemel\n1.25 ts bakepulver\nLemoncurd\n12.5 stk eggeplommer\n300 g sukker\n10 stk sitroner (kun saften av)\n2.5 smør\nKaken\n2.5 boks ananas\n\nFremgangsmåte\nSukkerbrød: Pisk egg og sukker til eggedosis.\nSikt mel og bakepulver over eggedosisen og vend forsiktig inn.\nHa deigen i en kakeform med bakepapirkledd bunn.\nStek kaken på rist, midt i ovnen ved 175°C i ca 30-40 min. Avkjøl kaken helt i formen.\nTa kaken ut av formen og fjern bakepapiret. Hvis nødvendig bruk kniv for å løsne kaken fra formen.\nLemoncurd: Pisk sammen sukker og eggeplommer i en kasserolle.\nTilsett sitronsaften.\nVarm opp på svak varme til det tykner. Rør hele tiden.\nSett kasserollen til side og rør inn det romtempererte smøret.\nVips, den deiligste curd!\nKaken: Del i to eller tre bunner. Ønsker du en høy kake, dobler du oppskriften og lager to stykker.\nMonter kaken ved å smøre lemoncurd og ananas på den første bunnen, topp den med krem og repeter med neste lag så høyt du vil.\nPynt gjerne kremkaken med stjernefrukt, sitronskiver lett kandisert i sukker og mynte-blader. Spiselige blomster setter en ekstra piff på det hele. Ønsker du bær og andre deiligheter er det bare å fylle på med dette.\n\nhttps://spar.no/Oppskrifter/Festkake-med-lemoncurd-og-knust-ananas/")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "En frisk og god kake og en fryd for øyet!",
				Ingredients: []string{
					"Sukkerbrød",
					"7.5 stk egg",
					"250 g sukker",
					"250 g hvetemel",
					"1.25 ts bakepulver",
					"Lemoncurd",
					"12.5 stk eggeplommer",
					"300 g sukker",
					"10 stk sitroner (kun saften av)",
					"2.5 smør",
					"Kaken",
					"2.5 boks ananas",
				},
				Instructions: []string{
					"Sukkerbrød: Pisk egg og sukker til eggedosis.",
					"Sikt mel og bakepulver over eggedosisen og vend forsiktig inn.",
					"Ha deigen i en kakeform med bakepapirkledd bunn.",
					"Stek kaken på rist, midt i ovnen ved 175°C i ca 30-40 min. Avkjøl kaken helt i formen.",
					"Ta kaken ut av formen og fjern bakepapiret. Hvis nødvendig bruk kniv for å løsne kaken fra formen.",
					"Lemoncurd: Pisk sammen sukker og eggeplommer i en kasserolle.",
					"Tilsett sitronsaften.",
					"Varm opp på svak varme til det tykner. Rør hele tiden.",
					"Sett kasserollen til side og rør inn det romtempererte smøret.",
					"Vips, den deiligste curd!",
					"Kaken: Del i to eller tre bunner. Ønsker du en høy kake, dobler du oppskriften og lager to stykker.",
					"Monter kaken ved å smøre lemoncurd og ananas på den første bunnen, topp den med krem og repeter med neste lag så høyt du vil.",
					"Pynt gjerne kremkaken med stjernefrukt, sitronskiver lett kandisert i sukker og mynte-blader. Spiselige blomster setter en ekstra piff på det hele. Ønsker du bær og andre deiligheter er det bare å fylle på med dette.",
				},
				Keywords: make([]string, 0),
				Name:     "Festkake med lemoncurd og knust ananas",
				Times: models.Times{
					Prep: 40 * time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "https://spar.no/Oppskrifter/Festkake-med-lemoncurd-og-knust-ananas/",
				Yield: 10,
			},
		},
		{
			name: "recipe 18",
			buf:  bytes.NewBuffer([]byte("Frisk påskekake med sjokolade og pasjonsfrukt\n\nHva er vel ikke bedre enn å kunne servere ei frisk og god kake med nøttebunn til påskekosen? Den hvite sjokolademoussen blir balansert med den friske pasjonsfrukten.\n\nEnkel\nType: Søtt\nHovedingrediens: Nøtter/mandler\nAnledning: Påske\n\nIngredienser\nMandelbunn\n4 eggehviter, mellomstore/store\n200 g sukker, helst ekstra finkornet\n250 g malte mandler\nHvit sjokolademousse\n300 g hvit sjokolade\n2 pasjonsfrukt\n3 plater gelatin\n1/2 vaniljestang\n2 dl melk\n2 eggeplommer\n1 ss sukker\n4 dl fløte\nPasjonsfruktsaus\n6 store pasjonsfrukt\n80 g sukker\n1,5 dl ferskpresset appelsinjuice\n\nSlik gjør du\nMandelbunn\nVarm ovnen til 180 grader. Kle en springform, som er 20–24 centimeter i diameter, med bakepapir. Visp eggehvitene ved middels hastighet til mykt skum. Tilsett sukkeret litt etter litt. Øk hastigheten når alt sukkeret er tilsatt, og pisk videre til en tykk og blank marengs. Vend inn mandlene, hell røren over i formen og bre den jevnt utover. Stek mandelbunnen midt i ovnen i cirka 30–35 minutter. Avkjøl og ta av kakeringen når bunnen er avkjølt.\nVask kakeringen og sett den tilbake rundt bunnen. Kle innsiden av ringen med bakepapir eller plastremse før du strammer den rundt bunnen.\n\nHvit sjokolademousse\nFinhakk sjokoladen og ha den over i en bolle.\nDel pasjonsfrukten og skrap ut fruktkjøttet. Kjør det raskt med en stavmikser, uten å knuse de sorte frøene, og gni pasjonsfrukten gjennom en sil for å fjerne frøene. Vær nøye med å få med all saften.\nBløtlegg gelatinplatene i rikelig med kaldt vann. Skrap ut vaniljefrøene, ha vaniljestang og frø sammen med melken i en liten kasserolle. La det få et sakte oppkok.\nVisp eggeplommene sammen med 1 spiseskje sukker til sukkeret er oppløst. Klem vannet ut av gelatinplatene og løs de opp i den varme kokende melken. Hell den varme melken over eggeplommene mens du visper energisk. Blandingen skal tykne noe.\nSil blandingen over sjokoladen, litt etter litt, mens du hele tiden rører energisk fra midten med en slikkepott. Kjør til slutt sjokoladeblandingen raskt med en stavmikser, slik at blandingen blir helt glatt og jevn. Rør inn pasjonsfrukten og la sjokoladeblandingen stå på kjøkkenbenken til den har romtemperatur, rør om av og til. Pisk fløten til myk krem og vend kremen inn i sjokoladeblandingen.\nHell sjokolademoussen over mandelbunnen i formen.\nLa sjokolademoussen stå i kjøleskapet i minst 6 timer før servering. Kaken kan med fordel fryses.\n\nPasjonsfruktsaus\nDel pasjonsfruktene til sausen og skrap ut fruktkjøttet. Kjør fruktkjøttet raskt med en stavmikser for å frigjøre saften i frøsekkene – uten å knuse frøene.\nHa det over i en kasserolle sammen med sukker og appelsinjuice.\nKok opp og la sausen koke i et par minutter. Jevn den eventuelt med en liten spiseskje maisenna rørt ut i litt kaldt vann.\nAvkjøl sausen før servering.\n\nAnretning\nHvis du har fryst kaken må du huske å la den tine før servering. Fjern form og bakepapir mens den ennå er frossen. Sett kaken over på ønsket serveringsfat, og tin den i kjøleskapet.\nHell litt av sausen over kaken ved servering, og server resten ved siden av. Som ekstra pynt er det fint med for eksempel spiselige blomster og høvlet hvit sjokolade.\n\nhttps://www.nrk.no/mat/frisk-paskekake-med-sjokolade-og-pasjonsfrukt-1.14513946")),
			want: models.Recipe{
				Category:    "søtt",
				Description: "Hva er vel ikke bedre enn å kunne servere ei frisk og god kake med nøttebunn til påskekosen? Den hvite sjokolademoussen blir balansert med den friske pasjonsfrukten.",
				Ingredients: []string{
					"Mandelbunn",
					"4 eggehviter, mellomstore/store",
					"200 g sukker, helst ekstra finkornet",
					"250 g malte mandler",
					"Hvit sjokolademousse",
					"300 g hvit sjokolade",
					"2 pasjonsfrukt",
					"3 plater gelatin",
					"1/2 vaniljestang",
					"2 dl melk",
					"2 eggeplommer",
					"1 ss sukker",
					"4 dl fløte",
					"Pasjonsfruktsaus",
					"6 store pasjonsfrukt",
					"80 g sukker",
					"1,5 dl ferskpresset appelsinjuice",
				},
				Instructions: []string{
					"Mandelbunn",
					"Varm ovnen til 180 grader. Kle en springform, som er 20–24 centimeter i diameter, med bakepapir. Visp eggehvitene ved middels hastighet til mykt skum. Tilsett sukkeret litt etter litt. Øk hastigheten når alt sukkeret er tilsatt, og pisk videre til en tykk og blank marengs. Vend inn mandlene, hell røren over i formen og bre den jevnt utover. Stek mandelbunnen midt i ovnen i cirka 30–35 minutter. Avkjøl og ta av kakeringen når bunnen er avkjølt.",
					"Vask kakeringen og sett den tilbake rundt bunnen. Kle innsiden av ringen med bakepapir eller plastremse før du strammer den rundt bunnen.",
					"Hvit sjokolademousse",
					"Finhakk sjokoladen og ha den over i en bolle.",
					"Del pasjonsfrukten og skrap ut fruktkjøttet. Kjør det raskt med en stavmikser, uten å knuse de sorte frøene, og gni pasjonsfrukten gjennom en sil for å fjerne frøene. Vær nøye med å få med all saften.",
					"Bløtlegg gelatinplatene i rikelig med kaldt vann. Skrap ut vaniljefrøene, ha vaniljestang og frø sammen med melken i en liten kasserolle. La det få et sakte oppkok.",
					"Visp eggeplommene sammen med 1 spiseskje sukker til sukkeret er oppløst. Klem vannet ut av gelatinplatene og løs de opp i den varme kokende melken. Hell den varme melken over eggeplommene mens du visper energisk. Blandingen skal tykne noe.",
					"Sil blandingen over sjokoladen, litt etter litt, mens du hele tiden rører energisk fra midten med en slikkepott. Kjør til slutt sjokoladeblandingen raskt med en stavmikser, slik at blandingen blir helt glatt og jevn. Rør inn pasjonsfrukten og la sjokoladeblandingen stå på kjøkkenbenken til den har romtemperatur, rør om av og til. Pisk fløten til myk krem og vend kremen inn i sjokoladeblandingen.",
					"Hell sjokolademoussen over mandelbunnen i formen.",
					"La sjokolademoussen stå i kjøleskapet i minst 6 timer før servering. Kaken kan med fordel fryses.",
					"Pasjonsfruktsaus",
					"Del pasjonsfruktene til sausen og skrap ut fruktkjøttet. Kjør fruktkjøttet raskt med en stavmikser for å frigjøre saften i frøsekkene – uten å knuse frøene.",
					"Ha det over i en kasserolle sammen med sukker og appelsinjuice.",
					"Kok opp og la sausen koke i et par minutter. Jevn den eventuelt med en liten spiseskje maisenna rørt ut i litt kaldt vann.",
					"Avkjøl sausen før servering.",
					"Anretning",
					"Hvis du har fryst kaken må du huske å la den tine før servering. Fjern form og bakepapir mens den ennå er frossen. Sett kaken over på ønsket serveringsfat, og tin den i kjøleskapet.",
					"Hell litt av sausen over kaken ved servering, og server resten ved siden av. Som ekstra pynt er det fint med for eksempel spiselige blomster og høvlet hvit sjokolade.",
				},
				Keywords: []string{"nøtter/mandler", "påske"},
				Name:     "Frisk påskekake med sjokolade og pasjonsfrukt",
				Tools:    make([]string, 0),
				URL:      "https://www.nrk.no/mat/frisk-paskekake-med-sjokolade-og-pasjonsfrukt-1.14513946",
				Yield:    1,
			},
		},
		{
			name: "recipe 19",
			buf:  bytes.NewBuffer([]byte("Grove polarbrød\n\n18 stk.\n\nIngredienser:\n350 g sammalt hvete, fin\n50 g sammalt hvete, grov\n300 g hvetemel\n2 ts sukker\n1 ts salt\n25 g gjær\n2 ss nøytral olje\n3 dl melk\n2 dl vann\n\nFramgangsmåte:\n1. Bland sammen alt det tørre. Varm melk og vann til det når cirka 37 grader, og rør ut gjæren i væsken.\n2. Elt sammen væske, gjær, olje og melblandingen. Kjør det gjerne i kjøkkenmaskin en stund. La deigen heve i ca 45 minutter.\n3. Del deigen i 18 emner, og form disse til boller før de trykkes eller kjevles flate. Legg på bakepapir, og lag små hull med enden av en skje eller en annen passende gjenstand.\n4. La polarbrødene etterheve i omlag 15 minutter. Stek på 225 grader i ca 9 minutter, og legg deretter til avkjøling på en bakerist.\n\nhttps://coop.no/extra/mat--trender/hjemmelagde-grove-polarbrod/")),
			want: models.Recipe{
				Category: "uncategorized",
				Ingredients: []string{
					"350 g sammalt hvete, fin",
					"50 g sammalt hvete, grov",
					"300 g hvetemel",
					"2 ts sukker",
					"1 ts salt",
					"25 g gjær",
					"2 ss nøytral olje",
					"3 dl melk",
					"2 dl vann",
				},
				Instructions: []string{
					"Bland sammen alt det tørre. Varm melk og vann til det når cirka 37 grader, og rør ut gjæren i væsken.",
					"Elt sammen væske, gjær, olje og melblandingen. Kjør det gjerne i kjøkkenmaskin en stund. La deigen heve i ca 45 minutter.",
					"Del deigen i 18 emner, og form disse til boller før de trykkes eller kjevles flate. Legg på bakepapir, og lag små hull med enden av en skje eller en annen passende gjenstand.",
					"La polarbrødene etterheve i omlag 15 minutter. Stek på 225 grader i ca 9 minutter, og legg deretter til avkjøling på en bakerist.",
				},
				Keywords: make([]string, 0),
				Name:     "Grove polarbrød",
				Tools:    make([]string, 0),
				URL:      "https://coop.no/extra/mat--trender/hjemmelagde-grove-polarbrod/",
				Yield:    18,
			},
		},
		{
			name: "recipe 20",
			buf:  bytes.NewBuffer([]byte("Hellstrøms asiatisk ribbe\n\n1 - 2 timer\n4 Voksne\n\nIngredienser\nEtt ribbestykke ca.3-4 kilo, med ben og svor\n1/2 liter appelsinjuice, fersk presset\n1/4 liter mangojuice eller 2 mango\n1/4 liter ananasjuice eller 1 ananas\n1/4 liter soyasaus\nEn flaske tørr hvitvin\n1/4 liter hønsebuljong/vann\n50 gr fersk ingefær\n2 hvitløk delt i to\n1 purreløk\n1 liten selleristang\n1 ss korianderfrø\n1 ss kummin\n1 ss sechuanpepper\n2 små kanelstenger\n6 stjerneanis\n2 ss honning\n2 røde chilipepper\n2 ss asiatisk østerssaus/oystersauce fra flaske\nSalt og pepper\n\nSalte og pepre ribben gjerne et par døgn i forkant, og lag fine ruter i svoren med en skarp kniv eller et barberblad. Rens frukten og kjør i en blender til juice/pure (hvis du ikke bruker ferdig juice).\nKutt purreløk og selleristang i biter, del en hel hvitløk i to og riv ingefæren på et lite rivjern. Fres de oppkuttede grønnsakene i en stor jerngryte i litt olivenolje. Tilsett resten av ingrediensene og kok opp.\nLegg ribbestykket med svorsiden ned i jerngryta, og sett jerngryta i ovnen på 180grader. Vi steker i cirka 1 time, snur kjøttet, fortsetter nok en time, siler sjyen og øser den over kjøttet den siste halvtimen, slik at vi får et glansfullt, lakkert aspekt, som en pekingand. \nSjyen koker vi litt inn, og øser den over kjøttet, som vi har delt i passende store serveringstykker. Server gjerne med spinat, ananas og bok shoy.\nGod fornøyelse!\n\nhttp://mat.tv3.no/oppskrifter/jul-med-hellstrom/hellstroms-asiatisk-ribbe")),
			want: models.Recipe{
				Category: "uncategorized",
				Ingredients: []string{
					"Ett ribbestykke ca.3-4 kilo, med ben og svor",
					"1/2 liter appelsinjuice, fersk presset",
					"1/4 liter mangojuice eller 2 mango",
					"1/4 liter ananasjuice eller 1 ananas",
					"1/4 liter soyasaus",
					"En flaske tørr hvitvin",
					"1/4 liter hønsebuljong/vann",
					"50 gr fersk ingefær",
					"2 hvitløk delt i to",
					"1 purreløk",
					"1 liten selleristang",
					"1 ss korianderfrø",
					"1 ss kummin",
					"1 ss sechuanpepper",
					"2 små kanelstenger",
					"6 stjerneanis",
					"2 ss honning",
					"2 røde chilipepper",
					"2 ss asiatisk østerssaus/oystersauce fra flaske",
					"Salt og pepper",
				},
				Instructions: []string{
					"Salte og pepre ribben gjerne et par døgn i forkant, og lag fine ruter i svoren med en skarp kniv eller et barberblad. Rens frukten og kjør i en blender til juice/pure (hvis du ikke bruker ferdig juice).",
					"Kutt purreløk og selleristang i biter, del en hel hvitløk i to og riv ingefæren på et lite rivjern. Fres de oppkuttede grønnsakene i en stor jerngryte i litt olivenolje. Tilsett resten av ingrediensene og kok opp.",
					"Legg ribbestykket med svorsiden ned i jerngryta, og sett jerngryta i ovnen på 180grader. Vi steker i cirka 1 time, snur kjøttet, fortsetter nok en time, siler sjyen og øser den over kjøttet den siste halvtimen, slik at vi får et glansfullt, lakkert aspekt, som en pekingand.",
					"Sjyen koker vi litt inn, og øser den over kjøttet, som vi har delt i passende store serveringstykker. Server gjerne med spinat, ananas og bok shoy.",
					"God fornøyelse!",
				},
				Keywords: make([]string, 0),
				Name:     "Hellstrøms asiatisk ribbe",
				Times: models.Times{
					Cook: 1*time.Hour + 2*time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "http://mat.tv3.no/oppskrifter/jul-med-hellstrom/hellstroms-asiatisk-ribbe",
				Yield: 4,
			},
		},
		{
			name: "recipe 21",
			buf:  bytes.NewBuffer([]byte("Olympiatoppens superbrød\n\nSupergodt brød etter oppskriften til Olympiatoppens egen kokk Harald Haugen. Brødet blir ekstra godt og saftig med Biola®. Ingenting slår fersk brød med brunost!\n\nTIPS!: Denne oppskriften gir 2 store eller 3 mindre brød. Lager du 3 små brød stekes de i ca. 40 minutter. 2 større brød trenger ca. 50 minutter i ovnen.\n\nVanskelig: Enkel\nTid: 150 min\nSteketid: 40-50 min\nAntall stykker 2\n\nIngredienser\n4 dl vann\n1 dl solsikkekjerner\n4 dl Biola syrnet lettmelk Naturell\n7.5 dl hvetemel\n7.5 dl sammalt hvete, grov\n1 pakke tørrgjær\n100 g hakkede hasselnøtter\n4 ss kruskakli (kan sløyfes)\n2 ts salt\n\nTil servering\n1  TINE Gudbrandsdalsost G35\n\nSlik gjør du\nKok opp vann og hell det i en bakebolle sammen med solsikkekjernene. La stå i 10 minutter\nNår vannet med solsikkekjerner har stått i 10 minutter blander du inn Biola, mel, gjær, nøtter og kruskakli. Elt på lavest hastighet i 5-6 minutter.\nHa i salt og elt videre i 4 minutter.\nHa gjerne i tørkede tranebær, rosiner eller valnøtter også.\nLa heve i ca. 30 minutter.\nForm fine brød, og hev til deigen kjennes luftig, ca. 30-40 minutter.\nSett stekeovnen på 180 °C.\nStek brødene på nederste rille i ovnen i 40-50 minutter. Server nybakt brød med brunost.\n\nhttp://www.tine.no/oppskrifter/bakst/brod-og-rundstykker/olympiatoppens-superbr%C3%B8d")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "Supergodt brød etter oppskriften til Olympiatoppens egen kokk Harald Haugen. Brødet blir ekstra godt og saftig med Biola®. Ingenting slår fersk brød med brunost!\n\nTIPS!: Denne oppskriften gir 2 store eller 3 mindre brød. Lager du 3 små brød stekes de i ca. 40 minutter. 2 større brød trenger ca. 50 minutter i ovnen.",
				Ingredients: []string{
					"4 dl vann",
					"1 dl solsikkekjerner",
					"4 dl Biola syrnet lettmelk Naturell",
					"7.5 dl hvetemel",
					"7.5 dl sammalt hvete, grov",
					"1 pakke tørrgjær",
					"100 g hakkede hasselnøtter",
					"4 ss kruskakli (kan sløyfes)",
					"2 ts salt",
					"Til servering",
					"1  TINE Gudbrandsdalsost G35",
				},
				Instructions: []string{
					"Kok opp vann og hell det i en bakebolle sammen med solsikkekjernene. La stå i 10 minutter",
					"Når vannet med solsikkekjerner har stått i 10 minutter blander du inn Biola, mel, gjær, nøtter og kruskakli. Elt på lavest hastighet i 5-6 minutter.",
					"Ha i salt og elt videre i 4 minutter.",
					"Ha gjerne i tørkede tranebær, rosiner eller valnøtter også.",
					"La heve i ca. 30 minutter.",
					"Form fine brød, og hev til deigen kjennes luftig, ca. 30-40 minutter.",
					"Sett stekeovnen på 180 °C.",
					"Stek brødene på nederste rille i ovnen i 40-50 minutter. Server nybakt brød med brunost.",
				},
				Keywords: make([]string, 0),
				Name:     "Olympiatoppens superbrød",
				Times: models.Times{
					Prep: 150 * time.Minute,
					Cook: 40 * time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "http://www.tine.no/oppskrifter/bakst/brod-og-rundstykker/olympiatoppens-superbr%C3%B8d",
				Yield: 2,
			},
		},
		{
			name: "recipe 22",
			buf:  bytes.NewBuffer([]byte("Ostekake med pepperkakebunn og klementingelé\n\nMarit Hegles ostekake har smaker som minner om jul. Kaken er et godt alternativ til dessert.\n\nMiddels\nType: Søtt\nKarakteristika: Kake\nLagre som favoritt\n\nIngredienser\n\nPepperkakebunn\n50 g søte havrekjeks (Digestive)\n200 g pepperkaker\n1 ts kanel\n80 g smør\n\nOstefyll\n4 plater gelatin\n250 g kremost naturell\n3 dl seterrømme\n150 g sukker\n2 ts vaniljesukker\n0,5 dl friskpresset klementinjuice\n3 dl kremfløte\n\nKlementingelé\n2,5 dl friskpresset klementinjuice (cirka 8 klementiner)\n3 plater gelatin\n4 ss sukker\n\nSlik gjør du\n\nPepperkakebunn\nKnus kjeksene så godt du kan, gjerne i en hurtigmikser eller ved å plassere kjeksene i en stor plastpose og rull et kjevle over kjeksene til de er helt knuste. Bland inn kanelen. Smelt smøret og bland det med kjeksene.\nPress blandingen ned i en rund kakeform på 22-24 centimeter. Det er en fordel å bruke en form med avtagbar bunn, eventuelt en løs kakering.\n\nOstefyll\nLegg gelatinplatene i kaldt vann. Bland seterrømme, kremost, sukker og vaniljesukker til en jevn masse. Pisk fløten til en løs krem i en egen bolle.\nVarm opp klementinjuicen i en kjele til den så vidt begynner å ryke. Fjern kjelen fra varmen, klem vannet av gelatinplatene og løs dem opp i den varme klementinjuicen. Hell juicen over i en kopp og vent til den er romtemperert. Da heller du den sakte over i osteblandingen mens du rører hurtig.\nBland så halvparten av den piskede kremen inn i ostemassen. Når den er jevnt blandet rører du inn resten av kremen.\nFordel ostekremen over kjeksbunnen. Avkjøl ostekaken i kjøleskap i minst 2 timer før du har på gelé.\n\nKlementingelé\nLøs opp gelatinplatene i kaldt vann i cirka 5 minutter. Varm opp cirka 0,5 dl av klementinjuicen sammen med sukkeret i en kjele til den så vidt begynner å ryke. Fjern kjelen fra varmen, klem vannet av gelatinplatene og løs dem opp i den varme juicen. Når gelatinen er helt oppløst rører du inn resten av klementinjuicen. Rør godt slik at gelatinen blir jevnt fordelt.\nLa geleen stå på benken til du ser at den nesten begynner å stivne. Da heller du den forsiktig over ostekaken. Sett kaken tilbake i kjøleskapet og la den stå til geleen har stivnet helt, det tar et par timer.\nPynt etter eget ønske, pepperkaker og julegodt passer fint.\n\nhttps://www.nrk.no/mat/ostekake-med-pepperkakebunn-og-klementingele-1.14771818")),
			want: models.Recipe{
				Category:    "søtt",
				Description: "Marit Hegles ostekake har smaker som minner om jul. Kaken er et godt alternativ til dessert.",
				Ingredients: []string{
					"Pepperkakebunn",
					"50 g søte havrekjeks (Digestive)",
					"200 g pepperkaker",
					"1 ts kanel",
					"80 g smør",
					"Ostefyll",
					"4 plater gelatin",
					"250 g kremost naturell",
					"3 dl seterrømme",
					"150 g sukker",
					"2 ts vaniljesukker",
					"0,5 dl friskpresset klementinjuice",
					"3 dl kremfløte",
					"Klementingelé",
					"2,5 dl friskpresset klementinjuice (cirka 8 klementiner)",
					"3 plater gelatin",
					"4 ss sukker",
				},
				Instructions: []string{
					"Pepperkakebunn",
					"Knus kjeksene så godt du kan, gjerne i en hurtigmikser eller ved å plassere kjeksene i en stor plastpose og rull et kjevle over kjeksene til de er helt knuste. Bland inn kanelen. Smelt smøret og bland det med kjeksene.",
					"Press blandingen ned i en rund kakeform på 22-24 centimeter. Det er en fordel å bruke en form med avtagbar bunn, eventuelt en løs kakering.",
					"Ostefyll",
					"Legg gelatinplatene i kaldt vann. Bland seterrømme, kremost, sukker og vaniljesukker til en jevn masse. Pisk fløten til en løs krem i en egen bolle.",
					"Varm opp klementinjuicen i en kjele til den så vidt begynner å ryke. Fjern kjelen fra varmen, klem vannet av gelatinplatene og løs dem opp i den varme klementinjuicen. Hell juicen over i en kopp og vent til den er romtemperert. Da heller du den sakte over i osteblandingen mens du rører hurtig.",
					"Bland så halvparten av den piskede kremen inn i ostemassen. Når den er jevnt blandet rører du inn resten av kremen.",
					"Fordel ostekremen over kjeksbunnen. Avkjøl ostekaken i kjøleskap i minst 2 timer før du har på gelé.",
					"Klementingelé",
					"Løs opp gelatinplatene i kaldt vann i cirka 5 minutter. Varm opp cirka 0,5 dl av klementinjuicen sammen med sukkeret i en kjele til den så vidt begynner å ryke. Fjern kjelen fra varmen, klem vannet av gelatinplatene og løs dem opp i den varme juicen. Når gelatinen er helt oppløst rører du inn resten av klementinjuicen. Rør godt slik at gelatinen blir jevnt fordelt.",
					"La geleen stå på benken til du ser at den nesten begynner å stivne. Da heller du den forsiktig over ostekaken. Sett kaken tilbake i kjøleskapet og la den stå til geleen har stivnet helt, det tar et par timer.",
					"Pynt etter eget ønske, pepperkaker og julegodt passer fint.",
				},
				Keywords: []string{"kake"},
				Name:     "Ostekake med pepperkakebunn og klementingelé",
				Tools:    make([]string, 0),
				URL:      "https://www.nrk.no/mat/ostekake-med-pepperkakebunn-og-klementingele-1.14771818",
				Yield:    1,
			},
		},
		{
			name: "recipe 23",
			buf:  bytes.NewBuffer([]byte("Ostekake med røde bær og roseblader\n\nNå kan du lage en kake som har fått VM-gull. Ostekaken til Sverre Sætre var med da han sammen med kokkelandslaget vant VM-gull for noen år tilbake. Oppskriften er en forenklet utgave av Sommerkaken med friske bær og krystalliserte roseblader.\n\nMiddels\nType: Søtt\nKarakteristika: Kake\nAnledning: Fest\n\nIngredienser\n\nKjeksbunn\n100 g havrekjeks\n15 g (2 ss) hakkede pistasjnøtter\n1 ts brunt sukker\n25 g smeltet smør\n1 ts nøtteolje av hasselnøtt eller valnøtt (du kan også bruke rapsolje)\n\nOstekrem\n3 gelatinplater\n1 dl sitronsaft (her kan du også bruke limesaft, pasjonsfruktjuice eller andre syrlige juicer)\n250 g kremost naturell\n150 g sukker\nfrøene fra en vaniljestang\n3 dl kremfløte\n\nTopping\n300 g friske bringebær\nkandiserte roseblader\nurter\n\nSlik gjør du\nBruk en kakering på 22 centimeter i diameter og fire centimeter høy.\n\nKjeksbunn\nKnus kjeksene og bland med pistasjnøttene, sukker og olje. Varm smøret slik at det blir nøttebrunt på farge og bland det med kjeksblandingen til en jevn masse. Sett kakeringen på en tallerken med bakepapir eller bruk en springform med bunn. Trykk ut kjeksmassen i bunnen av kakeformen.\nTips: Kle innsiden av ringen med bakepapir. Da blir det enklere å få ut bunnen.\n\nOstekrem\nBløtlegg gelatinen i kaldt vann i 5 minutter. Kjør kremost, vaniljefrø, sukker og halvparten av juicen til en glatt masse i en matprosessor. Varm resten av juicen til kokepunktet og ta den av platen. Kryst vannet ut av den oppbløtte gelatinen og la den smelte i den varme juicen.\nTilsett den varme juicen i ostemassen, og rør den godt inn. Dette kan gjøres i matprosessoren.\nPisk fløten til krem, og vend kremen inn i ostemassen med en slikkepott. Fyll ostekrem til toppen av ringen, og stryk av med en palett slik at kaken blir helt jevn. Sett kaken i kjøleskapet til den stivner.\nFør servering: Ta kaken ut av kjøleskapet. Dekk toppen av kaken med friske bær. Pynt med sukrede roseblader og urter.\n\nTips\nOstekremen kan også fylles i små glass og serveres med bringebærsaus. Gjør man dette, bør kremen stå 2 timer i kjøleskapet slik at den stivner.\n\nKandiserte roseblader\nKandiserte blomster og blader er nydelige og godt som pynt til kaker og desserter.\nPensle rosebladene med eggehvite.\nDryss sukker på bladene. Jeg pleier å knuse sukkeret i en morter, eller kjøre det i en matprosessor slik at det blir enda finere.\nLegg til tørking over natten.\nDenne teknikken kan brukes på alt av spiselige blomster og urter, som fioler, stemorsblomster, karse, rødkløver, hvitkløver, roseblader, nellik, markjordbær- og hagejordbærblomster, ringblomst, agurkurt, svarthyll, kornblomst, løvetann, mynte,\nsitronmelisse m.m.\nNB! Blomster er stort sett ikke regnet som matvarer, derfor er det ikke tatt hensyn når det gjelder sprøyting. Hvis man kjøper blomster til dette formål, må man altså passe på at de ikke er sprøytet.\n\nhttps://www.nrk.no/mat/ostekake-med-rode-baer-og-roseblader-1.8229671")),
			want: models.Recipe{
				Category:    "søtt",
				Description: "Nå kan du lage en kake som har fått VM-gull. Ostekaken til Sverre Sætre var med da han sammen med kokkelandslaget vant VM-gull for noen år tilbake. Oppskriften er en forenklet utgave av Sommerkaken med friske bær og krystalliserte roseblader.",
				Ingredients: []string{
					"Kjeksbunn",
					"100 g havrekjeks",
					"15 g (2 ss) hakkede pistasjnøtter",
					"1 ts brunt sukker",
					"25 g smeltet smør",
					"1 ts nøtteolje av hasselnøtt eller valnøtt (du kan også bruke rapsolje)",
					"Ostekrem",
					"3 gelatinplater",
					"1 dl sitronsaft (her kan du også bruke limesaft, pasjonsfruktjuice eller andre syrlige juicer)",
					"250 g kremost naturell",
					"150 g sukker",
					"frøene fra en vaniljestang",
					"3 dl kremfløte",
					"Topping",
					"300 g friske bringebær",
					"kandiserte roseblader",
					"urter",
				},
				Instructions: []string{
					"Kjeksbunn",
					"Knus kjeksene og bland med pistasjnøttene, sukker og olje. Varm smøret slik at det blir nøttebrunt på farge og bland det med kjeksblandingen til en jevn masse. Sett kakeringen på en tallerken med bakepapir eller bruk en springform med bunn. Trykk ut kjeksmassen i bunnen av kakeformen.",
					"Tips: Kle innsiden av ringen med bakepapir. Da blir det enklere å få ut bunnen.",
					"Ostekrem",
					"Bløtlegg gelatinen i kaldt vann i 5 minutter. Kjør kremost, vaniljefrø, sukker og halvparten av juicen til en glatt masse i en matprosessor. Varm resten av juicen til kokepunktet og ta den av platen. Kryst vannet ut av den oppbløtte gelatinen og la den smelte i den varme juicen.",
					"Tilsett den varme juicen i ostemassen, og rør den godt inn. Dette kan gjøres i matprosessoren.",
					"Pisk fløten til krem, og vend kremen inn i ostemassen med en slikkepott. Fyll ostekrem til toppen av ringen, og stryk av med en palett slik at kaken blir helt jevn. Sett kaken i kjøleskapet til den stivner.",
					"Før servering: Ta kaken ut av kjøleskapet. Dekk toppen av kaken med friske bær. Pynt med sukrede roseblader og urter.",
					"Tips",
					"Ostekremen kan også fylles i små glass og serveres med bringebærsaus. Gjør man dette, bør kremen stå 2 timer i kjøleskapet slik at den stivner.",
					"Kandiserte roseblader",
					"Kandiserte blomster og blader er nydelige og godt som pynt til kaker og desserter.",
					"Pensle rosebladene med eggehvite.",
					"Dryss sukker på bladene. Jeg pleier å knuse sukkeret i en morter, eller kjøre det i en matprosessor slik at det blir enda finere.",
					"Legg til tørking over natten.",
					"Denne teknikken kan brukes på alt av spiselige blomster og urter, som fioler, stemorsblomster, karse, rødkløver, hvitkløver, roseblader, nellik, markjordbær- og hagejordbærblomster, ringblomst, agurkurt, svarthyll, kornblomst, løvetann, mynte,",
					"sitronmelisse m.m.",
					"NB! Blomster er stort sett ikke regnet som matvarer, derfor er det ikke tatt hensyn når det gjelder sprøyting. Hvis man kjøper blomster til dette formål, må man altså passe på at de ikke er sprøytet.",
				},
				Keywords: []string{"kake", "fest"},
				Name:     "Ostekake med røde bær og roseblader",
				Tools:    make([]string, 0),
				URL:      "https://www.nrk.no/mat/ostekake-med-rode-baer-og-roseblader-1.8229671",
				Yield:    1,
			},
		},
		{
			name: "recipe 24",
			buf:  bytes.NewBuffer([]byte("Ostekake med rømme\n\nBunn:\n•\t200gr bixit, digestive eller annen kjeks (eller pepperkaker)\n•\t100gr smør\n\n1.\tKnus kjeksen. Smelt smøret. Bland sammen, og fordel oppå bakepapir i en form på ca 25cm.\n2.\tSett i kjøleskap mens du lager ostekrem.\n\nOstekrem:\n•\t1 liten kremfløte\n•\t1liten (200g) Philadelphia-ost eller annen naturell kremost,\n•\t1 beger lettrømme\n•\t1ts vaniljesukker\n•\t120gr melis\n•\tSaften av en sitron\n•\t1 pk sitron eller appelsin-gelé til ostekremen\n•\t1 pk gele til lokket\n\n1.\tLag 1 pakke gelé som på pakken, men bruk halv mengde vann. Kjøl ned, men ikke la den bli stiv.\n2.\tBland rømme, ost, sitronsaft, melis og vaniljesukker.\n3.\tPisk kremen.\n4.\tBland kremen i ostebandingen. Forsiktig så du ikke rører ut luften. \n5.\tHell i geleen og bland.\n6.\tHell blandingen over bunnen.\n7.\tLa stivne i kjøleskap.\n8.\tKok rød (eller annen gelé) som på pakken, men bruk igjen halv mengde vann. La avkjøles men ikke stivne.\n9.\tSleng på det du ønsker av frukt og bær (men ikke kiwi eller ananas, for da stivner ikke geleen). Hell forsiktig kald gelé over. \n10.\tLa stivne i kjøleskap.\n\nVoila! Da er det bare å forsyne seg \uF04A Er lurt å starte kvelden før, for hver gelé trenger 4-6 timer på å stivne.")),
			want: models.Recipe{
				Category: "uncategorized",
				Ingredients: []string{
					"Bunn:",
					"200gr bixit, digestive eller annen kjeks (eller pepperkaker)",
					"100gr smør",
					"Ostekrem:",
					"1 liten kremfløte",
					"1liten (200g) Philadelphia-ost eller annen naturell kremost,",
					"1 beger lettrømme",
					"1ts vaniljesukker",
					"120gr melis",
					"Saften av en sitron",
					"1 pk sitron eller appelsin-gelé til ostekremen",
					"1 pk gele til lokket",
				},
				Instructions: []string{
					"Knus kjeksen. Smelt smøret. Bland sammen, og fordel oppå bakepapir i en form på ca 25cm.",
					"Sett i kjøleskap mens du lager ostekrem.",
					"Lag 1 pakke gelé som på pakken, men bruk halv mengde vann. Kjøl ned, men ikke la den bli stiv.",
					"Bland rømme, ost, sitronsaft, melis og vaniljesukker.",
					"Pisk kremen.",
					"Bland kremen i ostebandingen. Forsiktig så du ikke rører ut luften.",
					"Hell i geleen og bland.",
					"Hell blandingen over bunnen.",
					"La stivne i kjøleskap.",
					"Kok rød (eller annen gelé) som på pakken, men bruk igjen halv mengde vann. La avkjøles men ikke stivne.",
					"Sleng på det du ønsker av frukt og bær (men ikke kiwi eller ananas, for da stivner ikke geleen). Hell forsiktig kald gelé over.",
					"La stivne i kjøleskap.",
					"Voila! Da er det bare å forsyne seg \uF04A Er lurt å starte kvelden før, for hver gelé trenger 4-6 timer på å stivne.",
				},
				Keywords: make([]string, 0),
				Name:     "Ostekake med rømme",
				Tools:    make([]string, 0),
				Yield:    1,
			},
		},
		{
			name: "recipe 25",
			buf:  bytes.NewBuffer([]byte("Påleggssalat med karri og kylling\n\nHar du prøvd hjemmelaget påleggssalat? Prøv denne enkle og gode påleggssalaten med kyllingpålegg og karri. Masse smak til brødskiven! Serveres på godt brød eller grove rundstykker. Holder seg fint 2-3 dager i kjøleskapet.\n\n15 min\nLett\nca 4 pers.\n\nIngredienser\n220 g kyllingfilet\n2 dl Lettrømme\n4 ss Majones\n1 ss Karripulver\n0,25 ts Chilipulver\n0,25 ts Sitronsaft\n3 ss Rødløk (finhakket)\n3 ss Rød paprika (finhakket)\n3 ss Mais\n2 ss Eple (finhakket)\nSalt & Pepper\n\nSlik gjør du\nSkjær kyllingpålegget i små biter.\nVisp sammen lettrømme/kesam, majones, karri og chili. Smak dressingen til med sitronsaft, salt og pepper.\nRør inn kyllingpålegg, rødløk, paprika, mais og eple.\nServeres med det samme, men blir bedre når den får stått litt. Oppbevar salaten i en tett beholder i kjøleskapet - så holder den seg fint 2-3 dager.\nServeres på godt brød eller grove rundstykker, gjerne med litt ruccolasalat som ekstra pynt.\n\nhttps://www.prior.no/oppskrifter/paaleggssalat-kylling-karri\n")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "Har du prøvd hjemmelaget påleggssalat? Prøv denne enkle og gode påleggssalaten med kyllingpålegg og karri. Masse smak til brødskiven! Serveres på godt brød eller grove rundstykker. Holder seg fint 2-3 dager i kjøleskapet.",
				Ingredients: []string{
					"220 g kyllingfilet",
					"2 dl Lettrømme",
					"4 ss Majones",
					"1 ss Karripulver",
					"0,25 ts Chilipulver",
					"0,25 ts Sitronsaft",
					"3 ss Rødløk (finhakket)",
					"3 ss Rød paprika (finhakket)",
					"3 ss Mais",
					"2 ss Eple (finhakket)",
					"Salt & Pepper",
				},
				Instructions: []string{
					"Skjær kyllingpålegget i små biter.",
					"Visp sammen lettrømme/kesam, majones, karri og chili. Smak dressingen til med sitronsaft, salt og pepper.",
					"Rør inn kyllingpålegg, rødløk, paprika, mais og eple.",
					"Serveres med det samme, men blir bedre når den får stått litt. Oppbevar salaten i en tett beholder i kjøleskapet - så holder den seg fint 2-3 dager.",
					"Serveres på godt brød eller grove rundstykker, gjerne med litt ruccolasalat som ekstra pynt.",
				},
				Keywords: make([]string, 0),
				Name:     "Påleggssalat med karri og kylling",
				Times: models.Times{
					Prep: 15 * time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "https://www.prior.no/oppskrifter/paaleggssalat-kylling-karri",
				Yield: 4,
			},
		},
		{
			name: "recipe 26",
			buf:  bytes.NewBuffer([]byte("Panang curry med svin\n\nPanang er en storfavoritt i Thailand. I en kald årstid er det lite som gir mer sol og varme enn denne bollen med curry fra Østen. Oppskriften er med svin, men den er like god med kylling.\n\n4 porsjoner\n30 min\nEnkel\nType: Middag\nHovedingrediens: Storfe\nKarakteristika: Panne/wok\nOpprinnelsesland: Thailand\n\nIngredienser\nCurry\n500 g nakkekoteletter uten bein\n25 g palmesukker/smak til med vanlig sukker\n4 dl kokosmelk\n1 dl vann\n4 ss panang curry paste\n1–2 ss fiskesaus\n2 kaffir limeblader\n1 paprika eller spisspaprika i skiver\n1 løk\n1 bunt thai basilikum (kan sløyfes)\n\nTilbehør\nris\n\nSlik gjør du\nStart med å frese tre spiseskjeer curry i cirka ett minutt.\nHa i halvparten av kokosmelken sammen med vann. Bland godt til sausen begynner å putre.\nLegg i kjøttet sammen med kaffir limeblader, fiskesaus og palmesukker (bruker du kompakt, kan du banke opp først. Eller du kan bare bruke vanlig sukker).\nLa det koke i tre-fire minutter til kjøttet er gjennomkokt.\nHa i grønnsaker og den resterende kokosmelken. La det koke i fire-fem minutter til.\nSmak gjerne til med litt mer fiskesaus om du ønsker.\nVend helt til slutt inn thai basilikum. Denne gir en søt og god aroma til curryen.\nServeres varm med ris.\n\nTips\nPanang curry paste kan kjøpes fra asiatiske matbutikker.\nThai basilikum og kaffir lime blader kan kjøpes fra asiatiske matbutikker eller diverse dagligvarebutikker.\nThai basilikum har en søt og god aroma, men finner du ikke dette kan du smake til med vanlig basilikum.\n\nhttps://www.nrk.no/mat/thailandsk-panang-curry-med-nakkekoteletter-1.16274001")),
			want: models.Recipe{
				Category:    "middag",
				Cuisine:     "thailand",
				Description: "Panang er en storfavoritt i Thailand. I en kald årstid er det lite som gir mer sol og varme enn denne bollen med curry fra Østen. Oppskriften er med svin, men den er like god med kylling.",
				Ingredients: []string{
					"Curry",
					"500 g nakkekoteletter uten bein",
					"25 g palmesukker/smak til med vanlig sukker",
					"4 dl kokosmelk",
					"1 dl vann",
					"4 ss panang curry paste",
					"1–2 ss fiskesaus",
					"2 kaffir limeblader",
					"1 paprika eller spisspaprika i skiver",
					"1 løk",
					"1 bunt thai basilikum (kan sløyfes)",
				},
				Instructions: []string{
					"Start med å frese tre spiseskjeer curry i cirka ett minutt.",
					"Ha i halvparten av kokosmelken sammen med vann. Bland godt til sausen begynner å putre.",
					"Legg i kjøttet sammen med kaffir limeblader, fiskesaus og palmesukker (bruker du kompakt, kan du banke opp først. Eller du kan bare bruke vanlig sukker).",
					"La det koke i tre-fire minutter til kjøttet er gjennomkokt.",
					"Ha i grønnsaker og den resterende kokosmelken. La det koke i fire-fem minutter til.",
					"Smak gjerne til med litt mer fiskesaus om du ønsker.",
					"Vend helt til slutt inn thai basilikum. Denne gir en søt og god aroma til curryen.",
					"Serveres varm med ris.",
					"Tips",
					"Panang curry paste kan kjøpes fra asiatiske matbutikker.",
					"Thai basilikum og kaffir lime blader kan kjøpes fra asiatiske matbutikker eller diverse dagligvarebutikker.",
					"Thai basilikum har en søt og god aroma, men finner du ikke dette kan du smake til med vanlig basilikum.",
				},
				Keywords:  []string{"storfe", "panne/wok"},
				Name:      "Panang curry med svin",
				Nutrition: models.Nutrition{},
				Times: models.Times{
					Prep: 30 * time.Minute,
				},
				Tools:     make([]string, 0),
				UpdatedAt: time.Time{},
				URL:       "https://www.nrk.no/mat/thailandsk-panang-curry-med-nakkekoteletter-1.16274001",
				Yield:     4,
			},
		},
		{
			name: "recipe 27",
			buf:  bytes.NewBuffer([]byte("Potetmos med mandelpotet\nEn enkel potetmos på mandelpoteter.\n\n4 personer\n\nIngredienser\n700 g mandelpotet\n1 dl melk\n50 g meierismør eller mer\nsalt etter smak\n\nSlik gjør du\n1. Skrell potetene, ha dem i en gryte ok kok opp. Juster end varmen og la potetene trekke like under kokepunktet til de akkurat er møre, omtrent 15 minutter. Hell av vannet og la potetene dampe seg tørre i gryta.\n2. Mos potetene sammen med smøret. Spe med melk til mosen er passe tjukk. Smak til med salt. Server som den er, eller dryss over urter eller gressløk.")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "En enkel potetmos på mandelpoteter.",
				Ingredients: []string{
					"700 g mandelpotet",
					"1 dl melk",
					"50 g meierismør eller mer",
					"salt etter smak",
				},
				Instructions: []string{
					"Skrell potetene, ha dem i en gryte ok kok opp. Juster end varmen og la potetene trekke like under kokepunktet til de akkurat er møre, omtrent 15 minutter. Hell av vannet og la potetene dampe seg tørre i gryta.",
					"Mos potetene sammen med smøret. Spe med melk til mosen er passe tjukk. Smak til med salt. Server som den er, eller dryss over urter eller gressløk.",
				},
				Keywords: make([]string, 0),
				Name:     "Potetmos med mandelpotet",
				Tools:    make([]string, 0),
				Yield:    4,
			},
		},
		{
			name: "recipe 28",
			buf:  bytes.NewBuffer([]byte("Saftig og holdbart matpakkebrød\n\nHemmeligheten bak dette brødet er at kornet bløtlegges i forkant, da blir brødet veldig saftig og får lang holdbarhet. Det er enkelt å lage, men det krever at du bruker et par minutter kvelden før du skal bake.\n\nEnkel\nOppskriften er stor nok til 3 brød i brødform.\nType: Frokost/Lunsj\nHovedingrediens: Korn\nKarakteristika: Brødmat\n\nIngredienser\nSteg 1\n70 g havregryn, lettkokte\n80 g varierte frø (f.eks. rugflak, sesamfrø, linfrø, gresskarfrø)\n12 dl kaldt vann\n\nSteg 2\n500 sammalt hvete fin\n300 g sammalt grovt rugmel (ev. kan du benytte sammalt hvete grov)\n\nSteg 3\n1 dl romtemperert vann\n1 ss honning\n25 g fersk gjær (ev. ½ pk. tørrgjær)\n2 ts salt\n2 ss raps- eller solsikkeolje\n650 g hvetemel, cirka\n\nSlik gjør du\nDu kan bruke de frøene du har i kjøkkenskapet ditt. Marit Hegle brukte rugflak, gresskarfrø og sesamfrø til dette brødet. Du kan også benytte linfrø eller andre frø.\nDet samme gjelder grovt mel. Du bør bruke 500 gram sammalt hvete fin, men de resterende 300 grammene grovt mel kan bestå av det du finner hjemme, inkludert sammalt hvete fin.\n\nKvelden før\nI en bolle blander du havregryn med frø og kaldt vann. Dekk bollen med plastfolie, og la den stå på kjøkkenbenken over natten.\n\n30 minutter før baking\nBland sammalt hvete fin og sammalt grovt rugmel inn i frøblandingen du bløtla kvelden før. La dette stå i 30 minutter.\nBland honning inn i 1 desiliter romtemperert vann. Løs opp gjæren i dette vannet.\nRør inn salt og olje i frø- og melblandingen. Rør så inn vannet som inneholder gjæren. Elt så inn litt og litt hvetemel til du har en jevn deig, som så vidt er litt klissete. Elt videre i 15 minutter. Dekk bakebollen med plastfolie og la deigen heve i cirka 60 minutter, til den har dobbel størrelse.\nNår deigen er ferdig hevet fordeler du den i 3 smurte brødformer. Dekk formene med plastfolie og la de etterheve i cirka 30 minutter.\nBrødene stekes på nederste rille i cirka 40 minutter, på 200 grader over/undervarme.\nNår brødene er ferdig stekt løsnes de forsiktig fra formene og avkjøles på rist.\n\nhttps://www.nrk.no/mat/saftig-og-holdbart-matpakkebrod-1.14993142")),
			want: models.Recipe{
				Category:    "frokost",
				Description: "Hemmeligheten bak dette brødet er at kornet bløtlegges i forkant, da blir brødet veldig saftig og får lang holdbarhet. Det er enkelt å lage, men det krever at du bruker et par minutter kvelden før du skal bake.",
				Ingredients: []string{
					"Steg 1",
					"70 g havregryn, lettkokte",
					"80 g varierte frø (f.eks. rugflak, sesamfrø, linfrø, gresskarfrø)",
					"12 dl kaldt vann",
					"Steg 2",
					"500 sammalt hvete fin",
					"300 g sammalt grovt rugmel (ev. kan du benytte sammalt hvete grov)",
					"Steg 3",
					"1 dl romtemperert vann",
					"1 ss honning",
					"25 g fersk gjær (ev. ½ pk. tørrgjær)",
					"2 ts salt",
					"2 ss raps- eller solsikkeolje",
					"650 g hvetemel, cirka",
				},
				Instructions: []string{
					"Du kan bruke de frøene du har i kjøkkenskapet ditt. Marit Hegle brukte rugflak, gresskarfrø og sesamfrø til dette brødet. Du kan også benytte linfrø eller andre frø.",
					"Det samme gjelder grovt mel. Du bør bruke 500 gram sammalt hvete fin, men de resterende 300 grammene grovt mel kan bestå av det du finner hjemme, inkludert sammalt hvete fin.",
					"Kvelden før",
					"I en bolle blander du havregryn med frø og kaldt vann. Dekk bollen med plastfolie, og la den stå på kjøkkenbenken over natten.",
					"30 minutter før baking",
					"Bland sammalt hvete fin og sammalt grovt rugmel inn i frøblandingen du bløtla kvelden før. La dette stå i 30 minutter.",
					"Bland honning inn i 1 desiliter romtemperert vann. Løs opp gjæren i dette vannet.",
					"Rør inn salt og olje i frø- og melblandingen. Rør så inn vannet som inneholder gjæren. Elt så inn litt og litt hvetemel til du har en jevn deig, som så vidt er litt klissete. Elt videre i 15 minutter. Dekk bakebollen med plastfolie og la deigen heve i cirka 60 minutter, til den har dobbel størrelse.",
					"Når deigen er ferdig hevet fordeler du den i 3 smurte brødformer. Dekk formene med plastfolie og la de etterheve i cirka 30 minutter.",
					"Brødene stekes på nederste rille i cirka 40 minutter, på 200 grader over/undervarme.",
					"Når brødene er ferdig stekt løsnes de forsiktig fra formene og avkjøles på rist.",
				},
				Keywords: []string{"korn", "brødmat"},
				Name:     "Saftig og holdbart matpakkebrød",
				Tools:    make([]string, 0),
				URL:      "https://www.nrk.no/mat/saftig-og-holdbart-matpakkebrod-1.14993142",
				Yield:    3,
			},
		},
		{
			name: "recipe 29",
			buf:  bytes.NewBuffer([]byte("Salsiccia, hjemmelaget italiensk pølse\n\nSalsiccia er en italiensk pølse med svinekjøtt og urter. Hvis du ikke vil lage pølser kan du like gjerne bruke det malte kjøttet til biffer eller farse. Paul Svensson anbefaler at du nyter denne delikatessen med en herlig grønn salat og sylteagurk.\n\nINGREDIENSER\n1 kg svinekam\n3 sjalottløk\n3 fedd hvitløk\n5 ss finhakkede urter (basilikum, timian, rosmarin)\n1 ss ristede fenikkelfrø\n1/2 fennikel\n4 ss maismel\n4 ss melk\nsalt og nykvernet pepper\nDessuten trenger du\nkjøttkvern\npølsestapper, pølsehorn\n10 meter svinetarmer (tykke pølser)\n10 meter lammetarmer (tynne pølser)\nhyssing\n\nSLIK GJØR DU\n\nGrovmal svinekjøttet to ganger i kjøttkvern. La løk, hvitløk, fenikkelfrø og finhakket fenikkel surre i stekepannen til det er mykt og mal det sammen med kjøtt og urter.\n\nBland maismel og melk og tilsett i farsen. Smak til med salt og pepper.\n\nTre tarmen over pølsehornet. Skru hornet på kvernen. Fyll i farse og lag pølser. Bind for med bomullsgarn og pass på at pølsene ikke blir stappet for fulle, de sprekker lett når de blir stekt. Knytt de ytterste pølsene til slutt, da kan du lettere tilpasse hvor mye du stapper i hver pølse.\n\nPølsene grilles hele, deles etterpå.\n\nhttp://www.nrk.no/mat/salsiccia_-hjemmelaget-italiensk-polse-1.10863582")),
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "Salsiccia er en italiensk pølse med svinekjøtt og urter. Hvis du ikke vil lage pølser kan du like gjerne bruke det malte kjøttet til biffer eller farse. Paul Svensson anbefaler at du nyter denne delikatessen med en herlig grønn salat og sylteagurk.",
				Ingredients: []string{
					"1 kg svinekam",
					"3 sjalottløk",
					"3 fedd hvitløk",
					"5 ss finhakkede urter (basilikum, timian, rosmarin)",
					"1 ss ristede fenikkelfrø",
					"1/2 fennikel",
					"4 ss maismel",
					"4 ss melk",
					"salt og nykvernet pepper",
					"Dessuten trenger du",
					"kjøttkvern",
					"pølsestapper, pølsehorn",
					"10 meter svinetarmer (tykke pølser)",
					"10 meter lammetarmer (tynne pølser)",
					"hyssing",
				},
				Instructions: []string{
					"Grovmal svinekjøttet to ganger i kjøttkvern. La løk, hvitløk, fenikkelfrø og finhakket fenikkel surre i stekepannen til det er mykt og mal det sammen med kjøtt og urter.",
					"Bland maismel og melk og tilsett i farsen. Smak til med salt og pepper.",
					"Tre tarmen over pølsehornet. Skru hornet på kvernen. Fyll i farse og lag pølser. Bind for med bomullsgarn og pass på at pølsene ikke blir stappet for fulle, de sprekker lett når de blir stekt. Knytt de ytterste pølsene til slutt, da kan du lettere tilpasse hvor mye du stapper i hver pølse.",
					"Pølsene grilles hele, deles etterpå.",
				},
				Keywords: make([]string, 0),
				Name:     "Salsiccia, hjemmelaget italiensk pølse",
				Tools:    make([]string, 0),
				URL:      "http://www.nrk.no/mat/salsiccia_-hjemmelaget-italiensk-polse-1.10863582",
				Yield:    1,
			},
		},
		{
			name: "recipe 30",
			buf:  bytes.NewBuffer([]byte("Sitronpotetmos\n\nDette trenger du til 4 personer:\n\n1 kg poteter\n100 g smør\n3 dl melk\nsalt og pepper\n1 knivspiss revet muskat\n1 ss finhakket timian\n2 ts finrevet sitronskall\n\nSlik gjør du:\nSkrell potetene og del dem i store terninger. Kok dem i usaltet vann.\nSmelt smør i en gryte og hell på melk. Krydre med salt, pepper og muskat. Når potetene er ferdigkokte, helles vannet av og de moses sammen med smørmelken.\nTilsett hakket timian og sitronskall rett før servering.\n\nhttps://coop.no/mega/hjemmerestauranten/fransk-gryterett-med-sitronpotetmos/")),
			want: models.Recipe{
				Category: "uncategorized",
				Ingredients: []string{
					"1 kg poteter",
					"100 g smør",
					"3 dl melk",
					"salt og pepper",
					"1 knivspiss revet muskat",
					"1 ss finhakket timian",
					"2 ts finrevet sitronskall",
				},
				Instructions: []string{
					"Skrell potetene og del dem i store terninger. Kok dem i usaltet vann.",
					"Smelt smør i en gryte og hell på melk. Krydre med salt, pepper og muskat. Når potetene er ferdigkokte, helles vannet av og de moses sammen med smørmelken.",
					"Tilsett hakket timian og sitronskall rett før servering.",
				},
				Keywords: make([]string, 0),
				Name:     "Sitronpotetmos",
				Tools:    make([]string, 0),
				URL:      "https://coop.no/mega/hjemmerestauranten/fransk-gryterett-med-sitronpotetmos/",
				Yield:    4,
			},
		},
		{
			name: "recipe 31",
			buf:  bytes.NewBuffer([]byte("Thai biffsalat med lime- og chilidressing\n\nBiffsalat med asiatiske thaismaker som får fram solskinnet. Lise Finckenhagen tilsetter lime, ingefær, chili, fiskesaus, koriander, mynte og peanøtter til sprø grønnsaker.\n\nOm oppskriften\n\nEnkel\nType: Middag\nHovedingrediens: Storfe\nKarakteristika: Salat\n2 porsjoner\n\nIngredienser\n2 møre biffer à ca. 150 g\n1 papaya eller mango\n2 gulrøtter\n1 mellomstor slangeagurk\n1 liten rødløk\n1 bunt frisk koriander\n1 håndfull mynteblader\n\nMarinade\n1 ss soyasaus\n1 ss fiskesaus\n1 ss lysebrunt sukker (ev. brunt eller hvitt)\n1 fedd hvitløk, revet\n\nDressing\n2 ss fiskesaus\nsaften av 2 lime\n2 ss lysebrunt sukker (ev. brunt eller hvitt)\n1 ts revet, frisk ingefær\n1 rød chili\n2 ts sesamolje\n2 ss hakkede peanøtter\n\nSLIK GJØR DU\nBland sammen ingrediensene til marinaden og hell den over i en liten brødpose. Legg i biffene, press luften ut av posen og knyt igjen. Mariner kjøttet i minst en time.\nVisp sammen alle ingrediensene til dressingen.\nTa kjøttet ut av marinaden. «Tørk» kjøttet lett med kjøkkenpapir. Stek biffene i litt olje i en varm stekepanne, ca. 2-3 minutter på hver side. La kjøttet hvile på en tallerken mens du lager salaten.\nRens og skjær papaya (mango) og gulrøtter i tynne strimler. Finsnitt rødløk. Lag lange bånd av slangeagurken (bruk en ostehøvel eller potetskreller), unngå å bruke den bløte kjernen. Bland alle grønnsakene i en bolle, tilsett dressing og friske urter. Bland godt og fordel salaten på tallerkener, eller ha alt på et fat.\nSkjær kjøttet i tynne skiver og legg skivene på eller ved siden av salaten.\n\nhttp://www.nrk.no/mat/biffsalat-med-lime--og-chilidressing-1.12213466")),
			want: models.Recipe{
				Category:    "middag",
				Description: "Biffsalat med asiatiske thaismaker som får fram solskinnet. Lise Finckenhagen tilsetter lime, ingefær, chili, fiskesaus, koriander, mynte og peanøtter til sprø grønnsaker.\n\nOm oppskriften",
				Ingredients: []string{
					"2 møre biffer à ca. 150 g",
					"1 papaya eller mango",
					"2 gulrøtter",
					"1 mellomstor slangeagurk",
					"1 liten rødløk",
					"1 bunt frisk koriander",
					"1 håndfull mynteblader",
					"Marinade",
					"1 ss soyasaus",
					"1 ss fiskesaus",
					"1 ss lysebrunt sukker (ev. brunt eller hvitt)",
					"1 fedd hvitløk, revet",
					"Dressing",
					"2 ss fiskesaus",
					"saften av 2 lime",
					"2 ss lysebrunt sukker (ev. brunt eller hvitt)",
					"1 ts revet, frisk ingefær",
					"1 rød chili",
					"2 ts sesamolje",
					"2 ss hakkede peanøtter",
				},
				Instructions: []string{
					"Bland sammen ingrediensene til marinaden og hell den over i en liten brødpose. Legg i biffene, press luften ut av posen og knyt igjen. Mariner kjøttet i minst en time.",
					"Visp sammen alle ingrediensene til dressingen.",
					"Ta kjøttet ut av marinaden. «Tørk» kjøttet lett med kjøkkenpapir. Stek biffene i litt olje i en varm stekepanne, ca. 2-3 minutter på hver side. La kjøttet hvile på en tallerken mens du lager salaten.",
					"Rens og skjær papaya (mango) og gulrøtter i tynne strimler. Finsnitt rødløk. Lag lange bånd av slangeagurken (bruk en ostehøvel eller potetskreller), unngå å bruke den bløte kjernen. Bland alle grønnsakene i en bolle, tilsett dressing og friske urter. Bland godt og fordel salaten på tallerkener, eller ha alt på et fat.",
					"Skjær kjøttet i tynne skiver og legg skivene på eller ved siden av salaten.",
				},
				Keywords: []string{"storfe", "salat"},
				Name:     "Thai biffsalat med lime- og chilidressing",
				Tools:    make([]string, 0),
				URL:      "http://www.nrk.no/mat/biffsalat-med-lime--og-chilidressing-1.12213466",
				Yield:    2,
			},
		},
		{
			name: "recipe 32",
			buf:  bytes.NewBuffer([]byte("Vaffelpoteter\n\nTilberedningstid 20 minutter\nAntall porsjoner 4\n\nIngredienser\n4 store poteter\n2 ss finhakket frisk kruspersille\n2 ss finskåret frisk gressløk\n30 g smør\nsmak til med salt og pepper\n\nFremgangsmåte\nSkrell potetene før de finrives på råkostjern. Ha raspet potet over i en bolle. Tilsett persille og gressløk. Bland alt godt sammen. Smak til med salt og hvit pepper. Smelt smøret og bland det i potetblandingen. Stek potetene i et vaffeljern.\n\nTips:\nGrønnsaker som gulrot, persillerot, knollselleri og kålrot, m.m. kan blandes med potetene. Vaffelpoteter serveres til stekt kjøtt eller fisk.\n\nhttp://www.aperitif.no/oppskrifter/oppskrift/vaffelpoteter/1991837")),
			want: models.Recipe{
				Category: "uncategorized",
				Ingredients: []string{
					"4 store poteter",
					"2 ss finhakket frisk kruspersille",
					"2 ss finskåret frisk gressløk",
					"30 g smør",
					"smak til med salt og pepper",
				},
				Instructions: []string{
					"Skrell potetene før de finrives på råkostjern. Ha raspet potet over i en bolle. Tilsett persille og gressløk. Bland alt godt sammen. Smak til med salt og hvit pepper. Smelt smøret og bland det i potetblandingen. Stek potetene i et vaffeljern.",
					"Tips:",
					"Grønnsaker som gulrot, persillerot, knollselleri og kålrot, m.m. kan blandes med potetene. Vaffelpoteter serveres til stekt kjøtt eller fisk.",
				},
				Keywords:  make([]string, 0),
				Name:      "Vaffelpoteter",
				Nutrition: models.Nutrition{},
				Times: models.Times{
					Prep: 20 * time.Minute,
				},
				Tools: make([]string, 0),
				URL:   "http://www.aperitif.no/oppskrifter/oppskrift/vaffelpoteter/1991837",
				Yield: 4,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := models.NewRecipeFromTextFile(tc.buf)
			if !cmp.Equal(got, tc.want) {
				t.Log(cmp.Diff(got, tc.want))
				t.Fail()
			}
		})
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
