package models_test

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
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

func TestNewRecipesFromRecipeKeeper(t *testing.T) {
	data := "<!DOCTYPE html><html><head><style type=\"text/css\">html {}body {font-family: sans-serif; font-size: 14px; line-height:1.4; color: #333; background-color: #fff}h2, h3 {font-weight: 500; line-height: 1.1; margin-top: 20px; margin-bottom: 10px}h2 {font-size: 24px}h3 {font-size: 14px}.recipe-details h2 {color: #d24400}.recipe-details h3 {color: #d24400}.recipe-ingredients p {margin: 0}.recipe-notes p {margin: 0}.recipe-photo {width: 250px; height: 250px; margin-top: 20px; object-fit: cover}.recipe-photos-div {width: 125px; height: 125px; margin-right: 5px; margin-bottom: 5px; display: inline-block}.recipe-photos {width: 100%; height: 100%; object-fit: contain}</style></head><body><div class=\"recipe-details\"><meta content=\"a2489b23-050f-5c8a-8e5d-12fdb9fdeea3\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/a2489b23-050f-5c8a-8e5d-12fdb9fdeea3_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Blueberry Cheesecake</h2><div>Courses: <span itemprop=\"recipeCourse\">Dessert</span></div><div>Categories: <span>Cake</span><meta content=\"Cake\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">Serves 6-8</span></div><div>Preparation time: <span>10 mins</span><meta content=\"PT10M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>1 hour </span><meta content=\"PT1H\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>1 1/4 cups graham cracker crumbs / 8 digestive biscuit crumbs</p><p>1/4 cup/50g butter, melted </p><p>20 oz./600g cream cheese </p><p>2 tablespoons all-purpose flour </p><p>3/4 cup/175g caster sugar </p><p>2 eggs, plus 1 yolk</p><p>small pot soured cream</p><p>vanilla extract  </p><p>1 cup / 300g fresh or frozen blueberries</p><p>icing sugar </p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. Melt the butter and pour onto the cracker crumbs and mix well. Press into the bottom of a 9-inch spring form pan. Bake at 325&#176;F until the crust is set, about 10-12 minutes. Allow to cool.</p><p></p><p>2. In a large bowl beat the cream cheese with the flour, caster sugar, eggs, soured cream and vanilla extract with an electric mixer until light and fluffy.</p><p></p><p>3. Pour the mixture into the pan and bake for 35-40 minutes until set. Remove from the oven and leave to cool.</p><p></p><p>4. Heat half the blueberries in a pan with 2 tablespoons icing sugar and stir gently until juicy. Squash the blueberries with a fork then continue to cook for a few minutes. Add the remaining blueberries, remove from the heat and allow to cool. </p><p></p><p>5. Pour the blueberries over the cheesecake just before serving.</p><p></p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Serving size: 1 slice<meta content=\"1 slice\" itemprop=\"recipeNutServingSize\"></div><div>Calories: 571<meta content=\"571\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 41.5g<meta content=\"41.5\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 24.7g<meta content=\"24.7\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 173mg<meta content=\"173\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 451mg<meta content=\"451\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 42g<meta content=\"42\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 1.5g<meta content=\"1.5\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 27.1g<meta content=\"27.1\" itemprop=\"recipeNutSugars\"></div><div>Protein: 10.1g<meta content=\"10.1\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/a2489b23-050f-5c8a-8e5d-12fdb9fdeea3_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"686a4630-2212-43f3-8cb3-87af6cd1ec2c\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"3\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/686a4630-2212-43f3-8cb3-87af6cd1ec2c_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Calabrian Chili Orzo with Vegetables</h2><div>Courses: <span itemprop=\"recipeCourse\"></span></div><div>Categories: <span>Chicken, Vegan, Vegetarian</span><meta content=\"Chicken\" itemprop=\"recipeCategory\"><meta content=\"Vegan\" itemprop=\"recipeCategory\"><meta content=\"Vegetarian\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\"><a href=\"https://www.allrecipes.com/calabrian-chili-orzo-with-vegetables-recipe-8623914\">https://www.allrecipes.com/calabrian-chili-orzo-with-vegetables-recipe-8623914</a></span></div><div>Serving size: <span itemprop=\"recipeYield\">4</span></div><div>Preparation time: <span>10 mins</span><meta content=\"PT10M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>15 mins</span><meta content=\"PT15M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>6 ounces orzo</p><p>2 tablespoons creme fraiche</p><p>1 teaspoon Calabrian chili paste, or more to taste</p><p>1 teaspoon olive oil</p><p>1 zucchini, or 2 if small, cut lengthwise and sliced into half moons</p><p>1/2 cup sliced sweet peppers</p><p>2 cloves garlic, minced, or more to taste</p><p>1/2 teaspoon Italian seasoning</p><p>salt and freshly ground black pepper to taste</p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>Bring a pot of lightly salted water to a boil, add orzo, and cook, uncovered, until tender, 8 to 10 minutes. Turn off the heat, drain thoroughly, and return orzo to the pot. Add cr&#232;me fra&#238;che and Calabrian chili paste. Set aside and keep warm.</p><p></p><p>Heat a skillet over medium-high heat. Add olive oil. Once oil is shimmering, add vegetables, and saute for 4 to 5 minutes. Stir in garlic, and saut&#233; until fragrant, about 45 seconds. Add Italian seasoning and season with salt and pepper. Remove vegetables from heat.</p><p></p><p>Gently stir vegetables into orzo. Serve immediately.</p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/686a4630-2212-43f3-8cb3-87af6cd1ec2c_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"b0bcddc4-23e8-50cc-a879-22b1b2e63919\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/b0bcddc4-23e8-50cc-a879-22b1b2e63919_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Chocolate Chip Cookies</h2><div>Courses: <span itemprop=\"recipeCourse\">Snack</span></div><div>Categories: <span>Cookie</span><meta content=\"Cookie\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">12</span></div><div>Preparation time: <span>5 mins</span><meta content=\"PT5M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>12 mins</span><meta content=\"PT12M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>2 1/4 cups all-purpose flour</p><p>1 teaspoon baking soda</p><p>1 teaspoon salt</p><p>1 cup butter</p><p>1 cup caster sugar</p><p>1 cup soft brown sugar</p><p>1 teaspoon vanilla extract</p><p>2 eggs</p><p>2 cups dark chocolate, broken into small pieces</p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. In a large bowl combine the flour, baking soda and salt.</p><p></p><p>2. In a separate bowl, mix the butter, caster sugar, brown sugar and vanilla extract until smooth. </p><p></p><p>3. Add the eggs and the flour to the mixture and beat to combine. </p><p></p><p>4. Add the chocolate pieces and stir.</p><p></p><p>5. Drop well rounded spoonfuls of dough onto a greased cookie sheet. </p><p></p><p>6. Bake at 375F for 8-10 minutes.</p><p></p><p>7. Remove from the oven and place cookies on a wire rack to cool.</p><p></p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Serving size: 1 cookie<meta content=\"1 cookie\" itemprop=\"recipeNutServingSize\"></div><div>Calories: 491<meta content=\"491\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 24.6g<meta content=\"24.6\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 15.8g<meta content=\"15.8\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 74mg<meta content=\"74\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 443mg<meta content=\"443\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 63.2g<meta content=\"63.2\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 1.6g<meta content=\"1.6\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 43g<meta content=\"43\" itemprop=\"recipeNutSugars\"></div><div>Protein: 5.7g<meta content=\"5.7\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/b0bcddc4-23e8-50cc-a879-22b1b2e63919_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"3259bd0c-30db-4beb-a8a8-2aeecc392824\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/3259bd0c-30db-4beb-a8a8-2aeecc392824_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Guyanese Gojas</h2><div>Courses: <span itemprop=\"recipeCourse\"></span></div><div>Categories: <span></span></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\"><a href=\"https://www.simplyrecipes.com/guyanese-gojas-recipe-5221034\">https://www.simplyrecipes.com/guyanese-gojas-recipe-5221034</a></span></div><div>Serving size: <span itemprop=\"recipeYield\">6 servings</span></div><div>Preparation time: <span>1 hour 5 mins</span><meta content=\"PT1H5M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>35 mins</span><meta content=\"PT35M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>For the goja dough:</p><p>2 cups (240g) all-purpose flour</p><p>2 teaspoons (8g) granulated sugar</p><p>1/2 teaspoon (2g) instant yeast</p><p>2 tablespoons diced cold butter, or vegetable shortening</p><p>1 cup cold whole milk</p><p>1 tablespoon (15g) all-purpose flour, to sprinkle while kneading</p><p>1/8 teaspoon neutral cooking oil, such as canola or vegetable for rubbing the dough ball</p><p>For the goja filling:</p><p>2 cups sweetened flaked coconut</p><p>1/2 teaspoon ground cinnamon</p><p>1/2 teaspoon freshly grated nutmeg</p><p>1 tablespoon light brown sugar</p><p>1 tablespoon fresh grated ginger</p><p>2 tablespoons water</p><p>2 tablespoons butter, melted</p><p>2 teaspoons vanilla extract</p><p>2 tablespoons neutral cooking oil such as canola or vegetable, for cooking filling</p><p>1/4 cup raisins</p><p>2 tablespoons water, for cooking filling</p><p>For shaping and frying the gojas:</p><p>1/4 cup water, for sealing the pastry</p><p>1/4 cup (34g) all-purpose flour, for crimping the gojas</p><p>3 cups canola or vegetable oil, for frying</p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>Combine the dry ingredients :</p><p></p><p>In a large bowl, combine the flour, sugar, and yeast. Add the butter. Using your hands or fingertips, rub the butter into flour until a coarse meal forms.</p><p></p><p>Alica Ramkirpal-Senhouse / Simply Recipes</p><p></p><p>Alica Ramkirpal-Senhouse / Simply Recipes</p><p></p><p>Pour milk into bowl:</p><p></p><p>Make a well in the center of the bowl, pour in the milk. Using a rubber spatula, stir until the dough forms. At this point, the dough will be a little sticky, sprinkle 1 tablespoon flour on the dough and knead it into dough with your hands in the bowl until the dough is no longer sticky.</p><p></p><p>Set dough aside to rest:</p><p></p><p>Rub the top of the dough with the oil and cover with a damp paper towel. Set aside for 15 to 20 minutes.</p><p></p><p>Make the filling:</p><p></p><p>In the bowl of your food processor, add the coconut, cinnamon, nutmeg, brown sugar, ginger, water, butter, and vanilla. Pulse on high until coconut becomes fine and pasty.</p><p></p><p>Cook the filling:</p><p></p><p>Heat a heavy-bottomed pan over low heat. Add the oil, coconut filling, raisins, and 2 tablespoons of water. Cook, stirring occasionally, until the sugar melts and the coconut looks more toasted and slightly darker in color, about 5 minutes. Remove from heat and let cool for a few minutes before assembling the gojas.</p><p></p><p>Alica Ramkirpal-Senhouse / Simply Recipes</p><p></p><p>Weigh and divide the dough:</p><p></p><p>Weigh the dough, then divide the weight by 12 to get the weight for each piece. Now, cut 12 small pieces of dough and weigh each. Add or remove small pieces until you get the exact weight you’re looking for.</p><p></p><p>If you’re not using a scale, divide the dough into 12 pieces using a knife or pastry cutter. Try to eyeball it so they’re all the same size.</p><p></p><p>Alica Ramkirpal-Senhouse / Simply Recipes</p><p></p><p>Roll the goja dough:</p><p></p><p>Round off each dough ball between your palms to form a ball, gently tucking dough under itself to make the top smooth. Once you’ve done this, cover all the dough balls with a damp paper towel to keep it from drying out and crusting.</p><p></p><p>Sprinkle flour on the surface of the dough ball you are working with. Working with one dough ball at a time, flatten slightly with your hands, then roll into a circle 1/8 inch thick and about 5 inches in diameter.</p><p></p><p>Flour your surface as needed as you go along.</p><p></p><p>Repeat with remaining balls of dough, being sure to keep them covered as you work.</p><p></p><p>Dip your pointer finger in water and run it around the outer edges of the dough. Place 2 tablespoons filing in the bottom half of the dough and bring the top half over to seal. Using a fork crimp the edges closed being sure to dip the fork in flour to keep from sticking while crimping. Place assembled gojas on a baking sheet lined with parchment paper.</p><p></p><p>Repeat this step for the rest of the batch.</p><p></p><p>Set up a plate or deep serving platter with a few paper towels to place gojas on after they’re done frying.</p><p></p><p>Heat a medium sized deep pot over medium-low heat. Add the oil and once it’s anywhere between 350-375&#176;F, fry the gojas for 2 to 3 minutes, you’ll have to cook these in batches, being sure to not overcrowd the pot. Use a slotted spoon or tongs to flip the gojas once halfway through cooking. Remove from oil once it is light golden brown and drain on paper towels.</p><p></p><p>Repeat with remaining gojas until they are all fried.</p><p></p><p>Enjoy warm.</p><p></p><p>Alica Ramkirpal-Senhouse / Simply Recipes</p></div></td></tr></table><h3>Notes</h3><div class=\"recipe-notes\" itemprop=\"recipeNotes\"><p></p><p></p><p>(per serving)</p><p>543 Calories 30g Fat 62g Carbs 8g Protein</p><p>Nutrition Facts</p><p>Servings: 6</p><p>Amount per serving</p><p>Calories 543</p><p>% Daily Value*</p><p>Total Fat 30g 38%</p><p>Saturated Fat 13g 67%</p><p>Cholesterol 17mg 6%</p><p>Sodium 132mg 6%</p><p>Total Carbohydrate 62g 23%</p><p>Dietary Fiber 5g 16%</p><p>Total Sugars 20g</p><p>Protein 8g</p><p>Vitamin C 0mg 1%</p><p>Calcium 66mg 5%</p><p>Iron 3mg 16%</p><p>Potassium 271mg 6%</p><p>*The % Daily Value (DV) tells you how much a nutrient in a food serving contributes to a daily diet. 2,000 calories a day is used for general nutrition advice.</p></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/3259bd0c-30db-4beb-a8a8-2aeecc392824_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"31adb931-743f-5e9d-b2a2-9269c0775a4a\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/31adb931-743f-5e9d-b2a2-9269c0775a4a_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Lemon Tart</h2><div>Courses: <span itemprop=\"recipeCourse\">Dessert</span></div><div>Categories: <span>Tart</span><meta content=\"Tart\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">8</span></div><div>Preparation time: <span>1 hour </span><meta content=\"PT1H\" itemprop=\"prepTime\"></div><div>Cooking time: <span>30 mins</span><meta content=\"PT30M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p><b>For the tart shell</b></p><p>1 1/2 cups all-purpose flour </p><p>1/2 cup icing sugar </p><p>2/3 cup softened butter </p><p>pinch of salt </p><p>1 egg yolk</p><p></p><p><b>For the lemon curd</b></p><p>6 lemons </p><p>6 large eggs </p><p>1 1/2 cups caster sugar </p><p>1 1/2 cups cream </p><p></p><p><b><i>To serve</i></b></p><p>Icing sugar</p><p>Cream</p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. In a large bowl mix the softened butter and icing sugar to a cream with a wooden spoon then beat in the egg yolk. Add the flour and salt and rub the butter mixture and flour together with your fingers until crumbly.</p><p></p><p>2. Add the egg yolk to the mixture and knead briefly until it forms a firm dough. Wrap in plastic wrap and leave to chill in the fridge for 30 minutes.</p><p></p><p>3. Roll out the pastry until very thin and line the quiche tin allowing a small amount of pastry to overlap the edges. Prick the base of the pastry with a fork and bake at 400&#176;F for 20 minutes.</p><p></p><p>4. Grate the zest from the lemons into a bowl then add the juice from the lemons. Break the eggs into a large bowl and add the caster sugar and mix well. Add the lemon juice, zest and the cream and whisk gently.</p><p></p><p>5. Pour the mixture into the pastry case and bake at 350&#176;F for 25 to 30 minutes until set.</p><p></p><p>6. Sprinkle with a little icing sugar and serve with cream.</p><p></p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Serving size: 1 slice<meta content=\"1 slice\" itemprop=\"recipeNutServingSize\"></div><div>Calories: 480<meta content=\"480\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 22.4g<meta content=\"22.4\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 12.7g<meta content=\"12.7\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 215mg<meta content=\"215\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 197mg<meta content=\"197\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 64.7g<meta content=\"64.7\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 0.6g<meta content=\"0.6\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 46.1g<meta content=\"46.1\" itemprop=\"recipeNutSugars\"></div><div>Protein: 8g<meta content=\"8\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/31adb931-743f-5e9d-b2a2-9269c0775a4a_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"b05310a1-15ec-560e-896d-e7a54223d0a0\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/b05310a1-15ec-560e-896d-e7a54223d0a0_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Pancakes</h2><div>Courses: <span itemprop=\"recipeCourse\">Breakfast</span></div><div>Categories: <span></span></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">8-12 pancakes</span></div><div>Preparation time: <span>5 mins</span><meta content=\"PT5M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>15 mins</span><meta content=\"PT15M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>3 cups/375g all-purpose flour</p><p>3 teaspoons baking powder</p><p>1 tablespoon caster sugar</p><p>1 1/2 cups/375ml milk</p><p>3 eggs</p><p>1/2 teaspoon salt</p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. Sift the flour, baking powder, salt and caster sugar into a bowl. </p><p></p><p>2. Break the eggs into a separate bowl and whisk together with the milk.</p><p></p><p>3. Gradually add the milk and egg mixture to the flour mixture and whisk to a smooth batter.</p><p></p><p>4. Heat a frying pan over a medium heat and melt a small knob of butter. Pour the batter into the pan, using approximately 1/4 cup for each pancake.</p><p></p><p>5. When the top of the pancake begins to bubble, turn and cook the other side until golden brown.</p><p></p><p>6. Serve with butter and maple syrup.</p><p></p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Serving size: 1 pancake<meta content=\"1 pancake\" itemprop=\"recipeNutServingSize\"></div><div>Calories: 180<meta content=\"180\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 2.4g<meta content=\"2.4\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 0.9g<meta content=\"0.9\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 52mg<meta content=\"52\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 154mg<meta content=\"154\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 32.4g<meta content=\"32.4\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 1g<meta content=\"1\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 3.1g<meta content=\"3.1\" itemprop=\"recipeNutSugars\"></div><div>Protein: 6.7g<meta content=\"6.7\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/b05310a1-15ec-560e-896d-e7a54223d0a0_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"7f2decec-2e7f-594b-90ac-5fed6ed1b4c7\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/7f2decec-2e7f-594b-90ac-5fed6ed1b4c7_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Pasta Puttanesca</h2><div>Courses: <span itemprop=\"recipeCourse\">Main Dish</span></div><div>Categories: <span>Pasta</span><meta content=\"Pasta\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">5</span></div><div>Preparation time: <span>5 mins</span><meta content=\"PT5M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>10 mins</span><meta content=\"PT10M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>4 tablespoons olive oil</p><p>1/8 cup butter</p><p>2 cloves garlic</p><p>4 anchovy fillets</p><p>1lb plum tomatoes</p><p>1 tablespoon tomato puree</p><p>1/2 cup black olives</p><p>1/2 tablespoon capers</p><p>1lb Spaghetti or any other pasta</p><p>Parmesan cheese</p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. Cook the pasta as directed on the instructions. </p><p></p><p>2. At the same time, heat the oil and butter in a separate pan and gently fry the garlic and anchovies for 3-4 minutes.</p><p></p><p>3. Add the tomatoes, tomato puree, olives, capers and fry for 7-8 minutes, stirring from time to time.</p><p></p><p>4. Drain the pasta, pour over the sauce and serve with freshly grated Parmesan cheese.</p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Serving size: 1 bowl<meta content=\"1 bowl\" itemprop=\"recipeNutServingSize\"></div><div>Calories: 459<meta content=\"459\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 20.9g<meta content=\"20.9\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 5.9g<meta content=\"5.9\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 85mg<meta content=\"85\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 500mg<meta content=\"500\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 54.9g<meta content=\"54.9\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 1.3g<meta content=\"1.3\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 2.3g<meta content=\"2.3\" itemprop=\"recipeNutSugars\"></div><div>Protein: 14g<meta content=\"14\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/7f2decec-2e7f-594b-90ac-5fed6ed1b4c7_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"2f98fbfc-bd5f-4a55-8774-48b29bd2c1ce\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Pesto Pull Apart-Bread</h2><div>Courses: <span itemprop=\"recipeCourse\"></span></div><div>Categories: <span>Cake</span><meta content=\"Cake\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\"><a href=\"https://sallysbakingaddiction.com/pesto-pull-apart-bread/#tasty-recipes-129499\">https://sallysbakingaddiction.com/pesto-pull-apart-bread/#tasty-recipes-129499</a></span></div><div>Serving size: <span itemprop=\"recipeYield\">Yield: 1 loaf</span></div><div>Preparation time: <span>3 hours </span><meta content=\"PT3H\" itemprop=\"prepTime\"></div><div>Cooking time: <span>50 mins</span><meta content=\"PT50M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>Dough</p><p>2 teaspoons instant or active dry yeast</p><p>1 Tablespoon granulated sugar</p><p>3/4 cup (6 1/4 ounces/ml) whole milk, warmed to about 110&#176;F (43&#176;C)</p><p>3 Tablespoons (1 1/2 ounces) unsalted butter, softened to room temperature</p><p>1 large egg, at room temperature</p><p>2 1/3 cups (10 1/4 ounces) all-purpose flour (spooned &amp; leveled), plus more as needed*</p><p>1 teaspoon salt</p><p>1 teaspoon garlic powder</p><p>1/2 teaspoon dried basil</p><p>Filling</p><p>1/2 cup (4 1/2 ounces) basil pesto (I recommend my homemade pesto)</p><p>1 cup (4 1/2 ounces / 4 ounces) shredded mozzarella cheese</p><p>Topping</p><p>2 Tablespoons (0.99 ounce) unsalted butter, melted</p><p>1/4 teaspoon garlic powder</p><p>2 Tablespoons (0.53 ounce) freshly grated or shredded parmesan cheese</p><p>optional for garnish: extra pesto</p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>Make the dough: Place the yeast and sugar in the bowl of a stand mixer fitted with a dough hook or paddle attachment. Or, if you do not own a stand mixer, a regular large mixing bowl. Whisk in the warm milk, then loosely cover with a clean kitchen towel and allow to sit for 5-10 minutes. The mixture will be frothy after 5-10 minutes.</p><p></p><p>If you do not have a mixer, you can mix the dough together with a wooden spoon or silicone spatula in this step. Add the butter, egg, flour, salt, garlic powder, and dried basil. Beat on low speed for 3 minutes. Dough will be soft.</p><p></p><p>Knead the dough: Keep the dough in the mixer (and switch to the dough hook if using the paddle) and beat for an additional 5 full minutes, or knead by hand on a lightly floured surface for 5 full minutes. (If you’re new to bread-baking, my How to Knead Dough video tutorial can help here.) If the dough becomes too sticky during the kneading process, sprinkle 1 teaspoon of flour at a time on the dough or on the work surface/in the bowl to make a soft, slightly tacky dough. Do not add more flour than you need because you do not want a dry dough. After kneading, the dough should still feel a little soft. Poke it with your finger—if it slowly bounces back, your dough is ready to rise. You can also do a “windowpane test” to see if your dough has been kneaded long enough: tear off a small (roughly golfball-size) piece of dough and gently stretch it out until it’s thin enough for light to pass through it. Hold it up to a window or light. Does light pass through the stretched dough without the dough tearing first? If so, your dough has been kneaded long enough and is ready to rise. If not, keep kneading until it passes the windowpane test.</p><p></p><p>1st Rise: Shape the kneaded dough into a ball. Place the dough in a greased bowl (I use nonstick spray to grease) and cover with plastic wrap or aluminum foil. Place in a slightly warm environment to rise until doubled in size, around 60-90 minutes. (If desired, use my warm oven trick for rising. See my answer to Where Should Dough Rise? in my Baking with Yeast Guide.)</p><p></p><p>As the dough rises, grease a 9&#215;5-inch loaf pan and prepare the pesto.</p><p></p><p>Assemble &amp; fill the bread: Punch down the dough to release the air. Place dough on a lightly floured work surface. Divide it into 12 equal pieces, about 1/4 cup of dough or 1 3/4 ounces each (a little larger than a golf ball). Using lightly floured hands, flatten each into a circle that’s about 4 inches in diameter. The circle doesn’t have to be perfectly round. I do not use a rolling pan to flatten, but you certainly can if you want. Spread 1-2 teaspoons of pesto onto each. Sprinkle each with 1 heaping Tablespoon of mozzarella cheese. Fold circles in half and line in prepared baking pan, round side up. See photos above for a visual.</p><p></p><p>2nd Rise: Cover with plastic wrap or aluminum foil and allow to rise once again in a slightly warm environment until puffy, about 45 minutes. Do not extend this 2nd rise, as the bread could puff up too much and spill over the sides while baking.</p><p></p><p>Adjust the oven rack to the lower third position then preheat oven to 350&#176;F (177&#176;C).</p><p></p><p>Bake until golden brown, about 50 minutes. If you find the top of the loaf is browning too quickly, tent with aluminum foil. Remove from the oven and place the pan on a cooling rack.</p><p></p><p>Make the topping: Mix the melted butter and garlic butter together. Brush on the warm bread and sprinkle with parmesan cheese. If desired, drop a couple spoonfuls of fresh pesto on top (or serve with extra pesto.) Cool for 10 minutes in the pan, then remove from the pan and serve warm.</p><p></p><p>Cover and store leftovers at room temperature for up to 2 days or in the refrigerator for up to 1 week. Since the bread is extra crispy on the exterior, it will become a little hard after day 1. Reheat in a 300&#176;F (149&#176;C) oven for 10-15 minutes until interior is soft again or warm in the microwave.</p></div></td></tr></table><h3>Notes</h3><div class=\"recipe-notes\" itemprop=\"recipeNotes\"><p>Make Ahead Instructions: Freeze baked and cooled bread for up to 3 months. Thaw at room temperature or overnight in the refrigerator and warm in the oven to your liking. The dough can be prepared through step 4, then after it has risen, punch it down to release the air, cover it tightly, then place in the refrigerator for up to 2 days. Continue with step 5. To freeze the dough, prepare it through step 4. After it has risen, punch it down to release the air. Wrap in plastic wrap and place in a freezer-friendly container for up to 3 months. When ready to use, thaw the dough overnight in the refrigerator. Then let the dough sit at room temperature for about 30 minutes before continuing with step 5. (You may need to punch it down again if it has some air bubbles.)</p><p></p><p>Special Tools (affiliate links): Electric Stand Mixer or Large Glass Mixing Bowl with Wooden Spoon / Silicone Spatula | 9&#215;5-inch Loaf Pan | Cooling Rack | Pastry Brush</p><p></p><p>Yeast: You can use instant or active dry yeast. The rise times may be slightly longer using active dry yeast. Reference my Baking with Yeast Guide for answers to common yeast FAQs.</p><p></p><p>Flour: Feel free to use the same amount of bread flour instead of all-purpose flour.</p><p></p><p>Can I substitute the pesto? Instead of pesto, you can use your favorite tomato sauce, or try this rosemary garlic pull apart bread.</p></div><h3>Photos</h3><hr /></div><div class=\"recipe-details\"><meta content=\"d017e2c8-24e6-5eee-8b68-931c0196bf5a\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/d017e2c8-24e6-5eee-8b68-931c0196bf5a_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Tuna Nicoise Salad</h2><div>Courses: <span itemprop=\"recipeCourse\">Main Dish</span></div><div>Categories: <span>Salad</span><meta content=\"Salad\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">4</span></div><div>Preparation time: <span>10 mins</span><meta content=\"PT10M\" itemprop=\"prepTime\"></div><div>Cooking time: <span></span><meta content=\"PT0S\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>2 cooked tuna steaks or 2 cans of tuna</p><p>12 small potatoes</p><p>5oz fine French beans</p><p>4 tomatoes</p><p>1 large romaine lettuce</p><p>1 red onion, finely sliced </p><p>4 hard-boiled eggs, peeled and sliced</p><p>20 black olives</p><p>Chopped fresh parsley</p><p>8 tablespoons extra virgin olive oil</p><p>3 tablespoons red wine vinegar</p><p>2 garlic cloves, peeled and finely chopped </p><p>1 teaspoon salt</p><p>1 teaspoon ground black pepper</p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. Cook the potatoes until just tender, cool in ice water and peel and quarter.</p><p></p><p>2. Top and tail the beans and boil in water for 5 minutes, cool in ice water.</p><p></p><p>3. Tear the lettuce into small pieces and arrange on a large plate.</p><p></p><p>4. Chop the tomatoes into quarters and add to the plate.</p><p></p><p>5. Cut the tuna into large chunks and add to the the plate.</p><p></p><p>6. Add the potatoes, beans, slice onion, eggs, olives and scatter over the chopped parsley.</p><p></p><p>7. In a small bowl mix the oil, vinegar, garlic, salt and pepper. Pour over the salad and serve.</p><p></p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Calories: 978<meta content=\"978\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 40.4g<meta content=\"40.4\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 6.9g<meta content=\"6.9\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 193mg<meta content=\"193\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 921mg<meta content=\"921\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 118.6g<meta content=\"118.6\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 25.4g<meta content=\"25.4\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 12.6g<meta content=\"12.6\" itemprop=\"recipeNutSugars\"></div><div>Protein: 41.4g<meta content=\"41.4\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/d017e2c8-24e6-5eee-8b68-931c0196bf5a_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div></body></html>"
	buf := bytes.NewBuffer([]byte(data))
	root, _ := goquery.NewDocumentFromReader(buf)

	got := models.NewRecipesFromRecipeKeeper(root)[:3]

	want := models.Recipes{
		{
			Category:    "Dessert",
			Description: "Imported from Recipe Keeper.",
			Ingredients: []string{
				"1 1/4 cups graham cracker crumbs / 8 digestive biscuit crumbs",
				"1/4 cup/50g butter, melted", "20 oz./600g cream cheese",
				"2 tablespoons all-purpose flour", "3/4 cup/175g caster sugar",
				"2 eggs, plus 1 yolk", "small pot soured cream", "vanilla extract",
				"1 cup / 300g fresh or frozen blueberries",
				"icing sugar",
			},
			Instructions: []string{
				"Melt the butter and pour onto the cracker crumbs and mix well. Press into the bottom of a 9-inch spring form pan. Bake at 325°F until the crust is set, about 10-12 minutes. Allow to cool.",
				"In a large bowl beat the cream cheese with the flour, caster sugar, eggs, soured cream and vanilla extract with an electric mixer until light and fluffy.",
				"Pour the mixture into the pan and bake for 35-40 minutes until set. Remove from the oven and leave to cool.",
				"Heat half the blueberries in a pan with 2 tablespoons icing sugar and stir gently until juicy. Squash the blueberries with a fork then continue to cook for a few minutes. Add the remaining blueberries, remove from the heat and allow to cool.",
				"Pour the blueberries over the cheesecake just before serving.",
			},
			Keywords: []string{"Cake"},
			Name:     "Blueberry Cheesecake",
			Nutrition: models.Nutrition{
				Calories:           "571 kcal",
				Cholesterol:        "173mg",
				Fiber:              "1.5g",
				IsPerServing:       true,
				Protein:            "10.1g",
				SaturatedFat:       "24.7g",
				Sodium:             "451mg",
				Sugars:             "27.1g",
				TotalCarbohydrates: "42g",
				TotalFat:           "41.5g",
			},
			Times: models.Times{Prep: 10 * time.Minute, Cook: 1 * time.Hour},
			Tools: []string{},
			URL:   "Recipe Keeper",
			Yield: 6,
		},
		{
			Category:    "uncategorized",
			Description: "Imported from Recipe Keeper.",
			Ingredients: []string{
				"6 ounces orzo", "2 tablespoons creme fraiche",
				"1 teaspoon Calabrian chili paste, or more to taste", "1 teaspoon olive oil",
				"1 zucchini, or 2 if small, cut lengthwise and sliced into half moons",
				"1/2 cup sliced sweet peppers", "2 cloves garlic, minced, or more to taste",
				"1/2 teaspoon Italian seasoning",
				"salt and freshly ground black pepper to taste",
			},
			Instructions: []string{
				"Bring a pot of lightly salted water to a boil, add orzo, and cook, uncovered, until tender, 8 to 10 minutes. Turn off the heat, drain thoroughly, and return orzo to the pot. Add crème fraîche and Calabrian chili paste. Set aside and keep warm.",
				"Heat a skillet over medium-high heat. Add olive oil. Once oil is shimmering, add vegetables, and saute for 4 to 5 minutes. Stir in garlic, and sauté until fragrant, about 45 seconds. Add Italian seasoning and season with salt and pepper. Remove vegetables from heat.",
				"Gently stir vegetables into orzo. Serve immediately.",
			},
			Keywords: []string{"Chicken", "Vegan", "Vegetarian"},
			Name:     "Calabrian Chili Orzo with Vegetables",
			Nutrition: models.Nutrition{
				IsPerServing: true,
			},
			Times: models.Times{Prep: 10 * time.Minute, Cook: 15 * time.Minute},
			Tools: []string{},
			URL:   "Recipe Keeper",
			Yield: 4,
		},
		{
			Category:    "Snack",
			Description: "Imported from Recipe Keeper.",
			Ingredients: []string{
				"2 1/4 cups all-purpose flour", "1 teaspoon baking soda", "1 teaspoon salt",
				"1 cup butter", "1 cup caster sugar", "1 cup soft brown sugar",
				"1 teaspoon vanilla extract", "2 eggs",
				"2 cups dark chocolate, broken into small pieces",
			},
			Instructions: []string{
				"In a large bowl combine the flour, baking soda and salt.",
				"In a separate bowl, mix the butter, caster sugar, brown sugar and vanilla extract until smooth.",
				"Add the eggs and the flour to the mixture and beat to combine.",
				"Add the chocolate pieces and stir.",
				"Drop well rounded spoonfuls of dough onto a greased cookie sheet.",
				"Bake at 375F for 8-10 minutes.",
				"Remove from the oven and place cookies on a wire rack to cool.",
			},
			Keywords: []string{"Cookie"},
			Name:     "Chocolate Chip Cookies",
			Nutrition: models.Nutrition{
				Calories:           "491 kcal",
				Cholesterol:        "74mg",
				Fiber:              "1.6g",
				IsPerServing:       true,
				Protein:            "5.7g",
				SaturatedFat:       "15.8g",
				Sodium:             "443mg",
				Sugars:             "43g",
				TotalCarbohydrates: "63.2g",
				TotalFat:           "24.6g",
			},
			Times: models.Times{Prep: 5 * time.Minute, Cook: 12 * time.Minute},
			Tools: []string{},
			URL:   "Recipe Keeper",
			Yield: 12,
		},
	}
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
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
