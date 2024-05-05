package models_test

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/units"
	"io"
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
		buf := bytes.NewBufferString("Ostekake med røde bær og roseblader\n\nNå kan du lage en kake som har fått VM-gull. Ostekaken til Sverre Sætre var med da han sammen med kokkelandslaget vant VM-gull for noen år tilbake. Oppskriften er en forenklet utgave av Sommerkaken med friske bær og krystalliserte roseblader.\n\nMiddels\nType: Søtt\nKarakteristika: Kake\nAnledning: Fest\n\nIngredienser\n\nKjeksbunn\n100 g havrekjeks\n15 g (2 ss) hakkede pistasjnøtter\n1 ts brunt sukker\n25 g smeltet smør\n1 ts nøtteolje av hasselnøtt eller valnøtt (du kan også bruke rapsolje)\n\nOstekrem\n3 gelatinplater\n1 dl sitronsaft (her kan du også bruke limesaft, pasjonsfruktjuice eller andre syrlige juicer)\n250 g kremost naturell\n150 g sukker\nfrøene fra en vaniljestang\n3 dl kremfløte\n\nTopping\n300 g friske bringebær\nkandiserte roseblader\nurter\n\nSlik gjør du\nBruk en kakering på 22 centimeter i diameter og fire centimeter høy.\n\nKjeksbunn\nKnus kjeksene og bland med pistasjnøttene, sukker og olje. Varm smøret slik at det blir nøttebrunt på farge og bland det med kjeksblandingen til en jevn masse. Sett kakeringen på en tallerken med bakepapir eller bruk en springform med bunn. Trykk ut kjeksmassen i bunnen av kakeformen.\nTips: Kle innsiden av ringen med bakepapir. Da blir det enklere å få ut bunnen.\n\nOstekrem\nBløtlegg gelatinen i kaldt vann i 5 minutter. Kjør kremost, vaniljefrø, sukker og halvparten av juicen til en glatt masse i en matprosessor. Varm resten av juicen til kokepunktet og ta den av platen. Kryst vannet ut av den oppbløtte gelatinen og la den smelte i den varme juicen.\nTilsett den varme juicen i ostemassen, og rør den godt inn. Dette kan gjøres i matprosessoren.\nPisk fløten til krem, og vend kremen inn i ostemassen med en slikkepott. Fyll ostekrem til toppen av ringen, og stryk av med en palett slik at kaken blir helt jevn. Sett kaken i kjøleskapet til den stivner.\nFør servering: Ta kaken ut av kjøleskapet. Dekk toppen av kaken med friske bær. Pynt med sukrede roseblader og urter.\n\nTips\nOstekremen kan også fylles i små glass og serveres med bringebærsaus. Gjør man dette, bør kremen stå 2 timer i kjøleskapet slik at den stivner.\n\nKandiserte roseblader\nKandiserte blomster og blader er nydelige og godt som pynt til kaker og desserter.\nPensle rosebladene med eggehvite.\nDryss sukker på bladene. Jeg pleier å knuse sukkeret i en morter, eller kjøre det i en matprosessor slik at det blir enda finere.\nLegg til tørking over natten.\nDenne teknikken kan brukes på alt av spiselige blomster og urter, som fioler, stemorsblomster, karse, rødkløver, hvitkløver, roseblader, nellik, markjordbær- og hagejordbærblomster, ringblomst, agurkurt, svarthyll, kornblomst, løvetann, mynte,\nsitronmelisse m.m.\nNB! Blomster er stort sett ikke regnet som matvarer, derfor er det ikke tatt hensyn når det gjelder sprøyting. Hvis man kjøper blomster til dette formål, må man altså passe på at de ikke er sprøytet.\n\nhttps://www.nrk.no/mat/ostekake-med-rode-baer-og-roseblader-1.8229671")
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
		{name: "has image only", recipe: models.Recipe{Images: []uuid.UUID{uuid.New()}}},
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
		Images:      make([]uuid.UUID, 0),
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
		Images:       []uuid.UUID{imageUUID},
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
	if schema.Image.Value != v+".jpg" {
		t.Errorf("wanted uuid %q but got %q", v, schema.Image.Value)
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
			buf:  bytes.NewBufferString("                     *  Exported from  MasterCook  *\n\n                                !Read Me!\n\nRecipe By     : Bill Wight\nServing Size  : 1    Preparation Time :0:00\nCategories    : *Afghan                          Info\n\n  Amount  Measure       Ingredient -- Preparation Method\n--------  ------------  --------------------------------\n                        ***Information***\n\nThe recipes in this collection were collected, corrected, edited and categorized by Bill Wight.\n\nYou may redistribute this collection in any form as long as this message is included.\n\nFor comments or suggestions, email me at:   Bill_Wight_CA@yahoo.com\n\nSome answers to frequently asked questions about my recipe collection:\n\nI notice that you have edited and corrected many of the recipe archives you have on your page.  What did you edit out and correct?\n\nMany of the recipes in this collection were typed, formatted and posted in the late 1980's and early 1990's to various bulletin boards before the Internet came into being.  These recipes were often formatted for a recipe management program called MealMaster.  During the years, the recipes have been posted and reposted to various mailing lists and newsgroups on the Internet.  They often pick up errors during reformatting and reposting.  I have also converted many of these recipes from Web page format and MealMaster format to MasterCook format.  This conversion process often introduces errors into the recipe.  I have attempted to correct the conversion errors.  I have also removed the names of the people who typed, formatted and posted these recipes.  I have kept the recipe author or creator where known.  I have also removed duplicate recipes and I have attempted to make some sense out of the recipes categories. \n\nI notice that you have removed all the nutrition data from the recipes in your collection.  Why did you do this?\n\nSince many of the recipes in this archive were typed and formatted several years ago and they often contain conflicting nutrition information, I have decided to strip out the nutritional data from all the recipes I edit.  If you have a current version of MasterCook (6.0) it will calculate nutritional data more accurately than the older versions.  So after you import the recipes into your copy of MasterCook, you will have up-to-date nutritional data available.  If nutritional data is important to you, make sure that the serving size looks reasonable.  Many recipes have a serving size set to the MasterCook default setting of '1'.  So nutritional data would assume the entire recipe is one serving and report the data incorrectly.  I have tired, as I edited the recipes, to correct the serving size if it was set to \"1\" to a reasonable serving size based on the recipes ingredient amounts.\n\n                   - - - - - - - - - - - - - - - - - - \n\n\n                     *  Exported from  MasterCook  *\n\n                      Abraysham Kabaub (Silk Kebab)\n\nRecipe By     : \nServing Size  : 30   Preparation Time :0:00\nCategories    : *Afghan                          Desserts\n\n  Amount  Measure       Ingredient -- Preparation Method\n--------  ------------  --------------------------------\n                        ***SYRUP***\n   1 1/2  cups          granulated sugar\n   1      teaspoon      lemon juice\n   1      cup           water\n     1/4  teaspoon      saffron threads -- (optional)\n                        ***OMELET***\n   8                    eggs\n   1      pinch         salt\n                        ***TO FINISH***\n   2      cups          oil\n     1/2  teaspoon      ground cardamom\n     3/4  cup           finely chopped pistachios *\n\n*Note: Instead of pistachio nuts, walnuts may be used if desired. \n\nDissolve sugar in water in heavy pan over medium heat. Bring to the boil, add lemon juice and saffron and boil for 10 minutes. Cool and strain into a 25 cm (10 inch) pie plate. Keep aside. Break eggs into a casserole dish about 20 cm (8 inches) in diameter. \n\nThe size and flat base are important. Add salt and mix eggs with fork until yolks and whites are thoroughly combined - do not beat as eggs must not be foamy. Heat oil in an electric frying pan to 190 C (375 F) or in a 25 cm (10 inch) frying pan placed on a thermostatically controlled hot plate or burner. Have ready nearby a long skewer, the plate of syrup, a baking sheet and the nuts mixed with the cardamom. \n\nA bowl of water and a cloth for drying hands are also necessary. Hold dish with eggs in one hand next to the pan of oil and slightly above it. Put hand into egg, palm down, so that egg covers back of hand. Lift out hand, curling fingers slightly inwards, then open out over hot oil, fingers pointing down. \n\nMove hand across surface of oil so that egg falls in streams from fingertips. Dip hand in egg again and make more strands across those already in pan. Repeat 3 to 4 times until about an eighth of the egg is used. There should be a closely meshed layer of egg strands about 20 cm (8 inches) across. Work quickly so that the last lot of egg is added not long after the first lot. Rinse hands quickly and dry. Take skewer and slide under bubbling omelet, lift up and turn over to lightly brown other side. \n\nThe first side will be bubbly, the underside somewhat smoother. When golden brown lift out with skewer and drain over pan. Place omelet flat in the syrup, spoon syrup over the top and lift out with skewer onto baking sheet. Roll up with bubbly side inwards. \n\nFinish roll should be about 3 cm (1 1/4 inches) in diameter. Put to one side and sprinkle with nuts. Repeat with remaining egg, making 7 to 8 rolls in all. Though depth of egg diminishes, you will become so adept that somehow you will get it in the pan in fine strands. \n\nWhen cool, cut kabaubs into 4-5 cm (1 1/2 to 2 inch pieces and serve. These keep well in a sealed container in a cool place.\n\n\n\n                   - - - - - - - - - - - - - - - - - - \n\n\n                     *  Exported from  MasterCook  *\n\n                              Afghan Chicken\n\nRecipe By     : San Francisco Examiner, 6/2/93.\nServing Size  : 6    Preparation Time :0:00\nCategories    : Middle East                      Chicken\n                *Afghan                          On-The-Grill\n\n  Amount  Measure       Ingredient -- Preparation Method\n--------  ------------  --------------------------------\n   2      large   clov  garlic\n     1/2  teaspoon      salt\n   2      cups          plain whole-milk yogurt\n   4      tablespoons   juice and pulp of 1 large lemon\n     1/2  teaspoon      cracked black pepper\n   2      large   whol  chicken breasts -- about 2 pounds\n\nLong, slow marinating in garlicky yogurt tenderizes, moistens and adds deep flavor, so you end up with skinless grilled chicken that's as delicious as it is nutritionally correct. Serve with soft pita or Arab flatbread and fresh yogurt.\n\nPut the salt in a wide, shallow non-reactive bowl with the garlic and mash them together until you have paste. Add yogurt, lemon and pepper.\n\nSkin the chicken breasts, remove all visible fat and separate the halves. Bend each backward to break the bones so the pieces win lie flat. Add to the yogurt and turn so all surfaces are well-coated.\n\nCover the bowl tightly and refrigerate. Allow to marinate at least overnight, up to a day and a half. Turn when you think of it.\n\nTo cook, remove breasts from marinade and wipe off all but a thin film. Broil or grill about 6 inches from the heat for 6 to 8 minutes a side, or until thoroughly cooked. Meat will brown somewhat but should not char. Serve at once.\n\n\n\n\n\n                   - - - - - - - - - - - - - - - - - - \n\n\n                     *  Exported from  MasterCook  *\n\n                          Afghan Chicken Kebobs\n\nRecipe By     : \nServing Size  : 4    Preparation Time :0:00\nCategories    : Chicken                          *Afghan\n                On-The-Grill\n\n  Amount  Measure       Ingredient -- Preparation Method\n--------  ------------  --------------------------------\n   1      cup           yogurt\n   1 1/2  teaspoons     salt\n     1/2  teaspoon      ground red or black pepper\n   3      centiliters   garlic -- finely minced\n   1 1/2  pounds        chicken breasts -- boneless,\n                        skinless -- cut into kebob\n                         -- ¥\n                        flatbread such as lavash\n                        pita or flour tortillas\n   3                    tomatoes -- sliced\n   2                    onions -- sliced\n                        cilantro to taste\n   2                    lemons or 4 limes -- quartered\n\n1. Mix yogurt, salt, pepper and garlic in a bowl. Mix chicken with yogurt and marinate 1 to 2 hours at room temperature, up to 2 days refrigerated.\n\n2. Thread chicken on skewers and grill over medium hot coals.\n\n3. Place warmed pita bread on plates (if using tortillas, toast briefly over flame), divide meat among them, top with tomato and onion slices and cilantro and fold bread over. Serve with lemon or lime quarters for squeezing.\n\n\n\n                   - - - - - - - - - - - - - - - - - - \n\n\n                     *  Exported from  MasterCook  *\n\n              Afghan Pumpkins Kadu Bouranee (Sweet Pumpkin)\n\nRecipe By     : \nServing Size  : 1    Preparation Time :0:00\nCategories    : *Afghan                          Vegetables\n\n  Amount  Measure       Ingredient -- Preparation Method\n--------  ------------  --------------------------------\n   2      pounds        fresh pumpkin or squash\n     1/4  cup           corn oil\n                        ***SWEET TOMATO SAUCE***\n   1      teaspoon      crushed garlic\n   1      cup           water\n     1/2  teaspoon      salt\n     1/2  cup           sugar\n   4      ounces        tomato sauce\n     1/2  teaspoon      ginger root -- chopped fine\n   1      teaspoon      freshly ground coriander\n                        seeds\n     1/4  teaspoon      black pepper\n                        ***YOGURT SAUCE***\n     1/4  teaspoon      crushed garlic\n     1/4  teaspoon      salt\n     3/4  cup           plain yogurt\n                        ***GARNISH***\n                        dry mint leaves -- crushed\n\nPeel the pumpkin and cut into 2-3\" cubes; set aside. Heat oil in a large frying pan that has a lid. Fry the pumpkins on both sides for a couple of minutes until lightly browned. Mix together ingredients for Sweet Tomato Sauce in a bowl then add to pumpkin mixture in fry pan. \n\nCover and cook 20-25 minutes over low heat until the pumpkin is cooked and most of the liquid has evaporated. (I don't know how it's going to evaporate if the pan is covered....-B.) \n\nMix together the ingredients for the yogurt sauce. To serve: Spread half the yogurt sauce on a plate and lay the pumpkin on top. Top with remaining yogurt and any cooking juices left over. Sprinkle with dry mint. \n\nMay be served with chalow (basmati rice) and naan or pita bread.\n\nFrom Afghani Cooking, the cookbook from Da Afghan Restaurant, Bloomington/Minneapolis, MN.\n\n\n                   - - - - - - - - - - - - - - - - - - \n"),
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

func TestNewRecipeFromCrouton(t *testing.T) {
	data := `{"tags":[{"name":"Dessert","color":"#3398EE","uuid":"FDB8E4F3-90B1-4C95-9837-05446EEB6EA1"}],"cookingDuration":60,"steps":[{"step":"Sponge","order":0,"uuid":"E4AF0B09-A955-4A33-96BC-B1913F67C8F0","isSection":true},{"isSection":false,"step":"Preheat the oven to 140°C fan forced\/ 150°C convection","uuid":"15A62BD6-BFBE-4EDE-AA5E-4FD1686ABF28","order":1},{"order":2,"isSection":false,"uuid":"F15E0B56-1532-4573-B8D1-43FD12E0D8FA","step":"Line the bottom of an 8-inch cake tin with parchment paper"},{"step":"In a medium-sized bowl whisk the egg yolks, mil, and oil","uuid":"1F08D3C2-CE28-49FA-BDC7-15A3150A07FB","isSection":false,"order":3},{"uuid":"4C84C506-A1D4-41FF-83F3-8140A7DA4190","step":"Sift the flour into the egg yolk mixture and mix until combined","order":4,"isSection":false},{"order":5,"isSection":false,"uuid":"2CD4AA8C-0A48-44F7-BB48-2F900F6CA2E0","step":"In another bowl with an electric whisk, or in the bowl of a stand mixer fitted with a whisk attachment, whip the egg whites with sugar until stiff peaks"},{"isSection":false,"uuid":"ACA0BA95-4B3D-40AF-AB4D-C0222264B813","order":6,"step":"Add 1\/3 of the meringue into the egg yolk mixture and mix until smooth"},{"step":"Transfer the lightened egg yolk mixture to the remaining meringue and fold carefully until just combined","order":7,"isSection":false,"uuid":"E238E7C4-9725-4C10-A472-236EA3CF6CE4"},{"step":"Transfer the batter to the cake tin","uuid":"EC598325-ACED-4569-8633-5373C8556863","order":8,"isSection":false},{"order":9,"isSection":false,"step":"Place the cake tin in a water bath (a tray\/tin of boiling water) and bake for 70 minutes","uuid":"8065C1A7-8A98-476F-9D4A-E2B0467B9539"},{"isSection":false,"order":10,"step":"Remove from the oven and allow it to cool completely","uuid":"3CAB45F7-3464-4432-A963-42DB04F01030"},{"order":11,"step":"Once cooled run a knife around the edge of the cake tin and invert the pan","isSection":false,"uuid":"A6554B1A-655A-47A6-81D2-0153C146E374"},{"order":12,"uuid":"FB052DFF-C561-419B-A2C8-24E7B042E910","step":"Wrap in cling wrap and place in the fridge until assembly","isSection":false},{"uuid":"8582E221-5939-46AD-A181-48DD2A2EF552","isSection":true,"step":"Whipped Cream","order":13},{"order":14,"step":"Whip the cream with an electric whisk and slowly stream the sugar in","uuid":"D18FEFD1-C531-4946-AEAA-A206E16B2FA9","isSection":false},{"step":"Beat until stiff peaks","isSection":false,"order":15,"uuid":"83A2E4D6-FB8A-44A9-A4FA-5CA32D220F48"},{"uuid":"049DCA4E-143F-4329-B45B-B2131C44208F","isSection":true,"step":"Assembly","order":16},{"step":"Combine the sugar and water in a small bowl and microwave for 30 seconds until melted, cool","order":17,"isSection":false,"uuid":"C31838BF-B597-44F6-AA3E-79B9EA7C0D14"},{"isSection":false,"step":"Slice half the punnet of strawberries","order":18,"uuid":"B219234F-29FE-4C6C-AFA2-BA16CCA902F2"},{"isSection":false,"order":19,"uuid":"4BF92E88-5197-40F1-9874-23BD146B51C9","step":"Slice the cooled cake into three layers"},{"order":20,"isSection":false,"uuid":"36D04CD1-F6E0-4091-AABA-5CA9FFD96EEA","step":"Lay one layer of cake down and brush with the sugar syrup"},{"uuid":"707D33E3-5C06-49DB-8CE8-80DF9DF61DB8","order":21,"step":"Spread on a layer of cream, a layer of strawberries and then cover with another layer of cream, repeat","isSection":false},{"isSection":false,"order":22,"uuid":"53BEB495-A4CE-4B84-879B-5B9A7FC0D56E","step":"Place the last layer of sponge on top and give the cake a thin crumb coat before icing the entire cake with cream"},{"uuid":"9EA0D9EB-BFCF-402E-9F64-16A0450198E7","isSection":false,"order":23,"step":"Place star tip into a piping bag and fill it with the remaining cream"},{"order":24,"isSection":false,"uuid":"3E4C558F-1EEE-48C2-A436-4D859CA1F5D7","step":"Pipe a border around the edge of the cake and decorate with the remaining strawberries"}],"ingredients":[{"uuid":"496A3A16-A103-488F-ABBA-52CF39B87333","ingredient":{"name":"large eggs","uuid":"4477FC10-D61D-41A4-8608-3A44A0EE5CA1"},"order":0,"quantity":{"quantityType":"ITEM","amount":4}},{"quantity":{"amount":60,"quantityType":"GRAMS"},"ingredient":{"uuid":"79A89102-7EFA-40F3-AB22-97A1D24BD1A7","name":"Whole milk (1\/4 cup)"},"order":1,"uuid":"9D3F4B5E-B5EB-4D27-9D6A-EF74E4FFE50E"},{"ingredient":{"name":"Vegetable oil (3 tbsp)","uuid":"15F3B99B-A359-4E1A-B830-A418D55D08D4"},"uuid":"A5D35402-02BE-43D7-B803-77358089B3DE","order":2,"quantity":{"quantityType":"MILLS","amount":45}},{"ingredient":{"name":"Cornstarch (1\/3 cup 2 tbsp)","uuid":"6E40F929-ECDF-40BC-9583-2CB9671A14A1"},"quantity":{"amount":55,"quantityType":"GRAMS"},"order":3,"uuid":"E015FBCD-607C-4CB2-B5FD-1A9EF7EC7973"},{"order":4,"ingredient":{"name":"All-purpose flour (1\/3 cup 2 tbsp)","uuid":"7E31F4E7-F455-4AAE-8824-F5B9336541CD"},"uuid":"5959FFDA-652D-4CA6-8B7F-F56719DBEDE3","quantity":{"amount":55,"quantityType":"GRAMS"}},{"uuid":"F0432996-3F35-4603-AF74-E0DD5CB8DAAA","quantity":{"quantityType":"GRAMS","amount":90},"ingredient":{"uuid":"96FB138B-BC4E-4961-ACFA-7EAC99178FFB","name":"Granulated sugar (1\/3 cup 2 tbsp)"},"order":5},{"quantity":{"quantityType":"GRAMS","amount":65},"ingredient":{"uuid":"1C7489E8-A4C9-430C-8990-99291FD2B5A7","name":"Granulated sugar (1\/3 cup)"},"order":6,"uuid":"972D2A4A-12AE-443A-87B7-2FAF057D1778"},{"uuid":"23498090-11C7-4A65-99C9-188008094BC7","quantity":{"quantityType":"MILLS","amount":80},"order":7,"ingredient":{"uuid":"0EADB895-371B-4264-91CF-D03A2E2EC30C","name":"Water (1\/3 cup)"}},{"ingredient":{"name":"Whipped Cream (2 1\/2 cups)","uuid":"74E74F16-AFD4-4810-9CA4-A25654816EB0"},"quantity":{"amount":600,"quantityType":"MILLS"},"uuid":"12DB939D-04B1-40A8-802C-A9798A6697DF","order":8},{"order":9,"uuid":"9D25BA89-98ED-4702-BFD7-1C9CB73C463F","ingredient":{"name":"Granulated sugar (1\/2 cup)","uuid":"8D8B87C2-118C-4936-8D36-7E97C4B7E866"},"quantity":{"quantityType":"GRAMS","amount":100}},{"ingredient":{"name":"Vanilla extract","uuid":"D3E8BA17-9EED-4A78-92C7-F8E4AC25BAED"},"uuid":"292AD6FF-78CC-4581-9216-CF5222F6A17F","quantity":{"amount":1,"quantityType":"TEASPOON"},"order":10},{"quantity":{"amount":370,"quantityType":"GRAMS"},"ingredient":{"uuid":"B2182781-DEAE-4A4F-84DE-BB04B7D1B60A","name":"Strawberries (13oz)"},"order":11,"uuid":"CEE9A508-5838-4ECF-9C16-BEA30CCAFA0A"}],"defaultScale":1,"images":["\/9j\/4AAQSkZJRgABAQAASABIAAD\/4QCMRXhpZgAATU0AKgAAAAgABQESAAMAAAABAAEAAAEaAAUAAAABAAAASgEbAAUAAAABAAAAUgEoAAMAAAABAAIAAIdpAAQAAAABAAAAWgAAAAAAAABIAAAAAQAAAEgAAAABAAOgAQADAAAAAQABAACgAgAEAAAAAQAAAOGgAwAEAAAAAQAAAOEAAAAA\/8AAEQgA4QDhAwEiAAIRAQMRAf\/EAB8AAAEFAQEBAQEBAAAAAAAAAAABAgMEBQYHCAkKC\/\/EALUQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29\/j5+v\/EAB8BAAMBAQEBAQEBAQEAAAAAAAABAgMEBQYHCAkKC\/\/EALURAAIBAgQEAwQHBQQEAAECdwABAgMRBAUhMQYSQVEHYXETIjKBCBRCkaGxwQkjM1LwFWJy0QoWJDThJfEXGBkaJicoKSo1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoKDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uLj5OXm5+jp6vLz9PX29\/j5+v\/bAEMAAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAf\/bAEMBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAf\/dAAQAD\/\/aAAwDAQACEQMRAD8A\/roooor0DnIvtR\/un\/x2qv2p\/wC6Kv1XrnOgo\/6V\/tf+O1NVf7T\/ALf6f\/bKKACiiq9c4Fii17fh\/wCy1Xqxa9vw\/wDZaALFFFFABVej7T9p7Z\/T2xn3z6e\/fLFr2\/D\/ANloAsUttajn\/wCK6dP9k++R8pOeDwCyVoUAfSfhH7N\/Ztv9f8a6S5trXNxx\/wCPjP8A6Cf55785ArivBH\/IOP8AnvXfUAc3c6bbfZbz\/Rz6549+vy5Pr6D1bOU8c022\/wBG1S7\/APrn89x9uq89crkpXtniS5trbQ7z+h\/LK7RjI68tnnDdDXktyLTTvDdxd5xjnP8AkZ\/\/AFcqwCsvOB89HW9S+1eIP9M7f\/qxtb0wefdeDlq+ZdS8beLf+Jp\/xNOP8\/j6\/pnGNq++8f2Hqmsd\/wCn0\/pvH+8u3a3yXqV1\/wASzVPz9vpzyOOv55XINfEYnEr+le1+2t\/y32dmfXYbC7dn+H\/kuvb7O3l7vJHx\/wCLNSuT9r1Q569\/x5xg9v4RkZ4X7y+o6bqWbX\/Iz\/46p6YP3m5GSFyErwHNtbf6X269cdO\/br9T7jmvbfCVra6jbW95eemT0\/QZbOR7gj+6zfK3rZJiLKs2rprr5Str1jf1l1Tepy5l9r5nSfaff9Kq\/wBoj0\/Va3P+Ka\/2q5HUtEubb\/SrT9ABjnr0P06577jy7\/Q\/WcN3\/A876q\/5n\/5KaVzc9cE+y9f64GeR0b1OdyovNXXQfT+op32k\/Zuh9e35fl+FYlzc56flx836t7j7wx1yQQUCCOis+iug5z\/\/0P6JfEnjbxbon2i7\/wCEotPsfX8\/wPTrz+Ho3lzftseAPCJuP+Eh1yz+2f59jjj2Pf7\/AN5u\/wDH\/wAG\/BPi37Zm4vPQ8jn9ffOc\/TOdy\/lT8bP2L\/7Subj\/AIR69vB\/u\/5U8\/72T04AJr8gxHEWePEte2emu1\/w5lp93ls+b9Mw2Aym6\/dP73a93bXk6eXNtbsz78+G37bFz8WvH9n4e8F6H\/aOh\/bP+JlqWM8H69O\/fnqCuCK\/SD7UPs3Qfr+XX8P6dq\/n+\/Y5tvGv7M2pappWq6X9vs\/+gl\/LPIx2HUdOjYr9cfDfx+0TW\/s32zI5x9e+f4QDnt0XoC+cN+g8O51\/sTeIq\/v9LK\/\/ANq+j\/vK+9rnyedYWOGxN6FHmw+l3d6\/PZ\/+S+r3l7\/RXN23iS2uc\/ZD157fgOA3OP4uAf7i\/NR\/af8AtV7x8+dJRWH9p9\/0q19qH91f\/HqANKis+tCgDQoqvRQAfZv9j9f\/ALXVm36\/j\/8AE1WrStvs2D9P6\/5\/WucBfsp\/vH\/x2rNt\/n\/x6i2ucdfy4+X9V9h94565ABL6lB0HufhL\/kGn6tXW1zfhv7N\/Ydt\/nvSa34k0TRNN1TVby9s\/+JfZ+mPb1I6d8\/n0roOczvFv+k+RaWn+HTt\/ET1xggdepwwbzb4gD7PolnpXT2Ht6ngj8N2eCduTuPg58afB3xlttQ1vwzcfbLOzvdQ033\/tDT+uOnt0OPpkmrWt239palcXZHfqB29OCvX0\/wDHhtKt52I\/3Z+r\/NHThsSv6Vr2763\/AD22V0fLvxIuf7F8N2elWeP5\/e+hGfxI+o+Y18ueNra203Q7Oz\/PPP4A5H6hvbFfVfi23tta1u8urz\/jy08fy\/JuD6AccfNkqvyp42ubbUtc6Z+wdO\/+Hc9QT6bDgyL+YZl9Z+s\/WLdtfN9dPT+b15rc0ft8L9Vv89fx\/Q+efFn\/AB7D\/PajRfEuo\/2b9jx39uOOn3l6n2PcHGPm8k8bfEj+2\/EeueEvCX\/Ew+wWvTAHbODz2xj7v8WARks3E\/AHxt4k+JNt4g+x2\/8AyA7waaefvY9ORjuP4cdATk7uPC8VYKhiU8RWfVe6tLbed353++zF9Wl\/0Co+tba5x1\/Lj5f1X2H3jnrkAEv7p4A1L+29OvLO868\/57nOP64x92vij4geJfEngD+y\/wDiV\/2h9vvG\/pjvn8x\/31nav0f8E9Ru9S1s9D7n359ee\/BC+gJwdv2+XZjgszwyr0K3639Ftu+67tvaPBiMNicK\/wDaaF9NNLJeT9\/rvf3fJO65bWt\/adNuri09edvp7D5mz\/4768ZKLmfak\/umuk+JH+j+JPcce3Pf+D6AfL07HivAdb8AfEfxH420vxD4euCdH\/4l\/fP5tzzjuyn2DYYV14jMsPhXbEpP0t+Puv7lFW87XOd4dYx\/7PRfVN8u\/Rdfv29Vpy+k\/ah\/dX\/x6j7UP7q\/+PVP\/wAIl42\/6Bho\/wCES8bf9Aw14H+uWWf8\/l90f\/lhwfVcR\/Mv\/Jj\/0f6tdS03RLa2uP0yOuf++sgevyem1sZbwHxJb6bbfaP9C46f3iffP7vGD0G0\/Q\/xe2XPiTQx\/n+XTPXrj24Oa8B8beJNNH2gWZ7fXOffv+f124r8dxOF37L8f\/JdO32t\/P3vv8Pt8v8A5E8V8Sf2bj8f\/wBXXjGPTv1\/gr5t8S639mufbHHTn6Eg4P6eqrhVr0jxJqVz\/XO7+mF5GPU+wGW3+BeLP7N4\/wDrf05\/Ln0rxoV8LhvLy6rr0aXlqn3R6X1aXmej+Cf2mdS8N3X2TWLzp+GR\/MYz\/eb2Az8v3V4A+Lem+LtNt7v8PXr3A47eh\/LIr8YdS1LTbb6evTH6HPfjCem8Z319Ifs46lqVzqX\/AB+\/6Hke+PfPX3PH8XG7DCvcyTip\/WKOGXSW+9r9N1be+72vzKxwZhkmHX+0razvfb5e6tfm+jb1Z+tFtqX6fhx64yev+83TJIJArStroc\/\/ABPXp\/tD3yfmIxwOQF8lttb\/AE98foASe\/cEds5eukttb755647f+gqBjHAweuBgkM36ufCnpP2n3\/SrH2n3\/SuJttS7fr16emcevXn7uTnAKbf2p\/7ooA6X7UP7q\/8Aj1VdS1K503RdUu7T\/iYf2faX\/wDL1yf5A+uazba56ZJ916f1wccDovqMbWRtL7T7\/pQB+Gnxr\/4Lz\/BP4Oajqngm58MeKj8QNAvP7N1Lw12+p+YD1OcsOwUYLP6j+xf+2N42\/aH8fXHxj+x\/YPh9qB\/5BvQ\/qoIyP9pc9SVwteGftw\/8Eu\/H\/wC0x+0hoeq\/D34d2eoaPfXmj6l4l8S5OP5n+bZ7bSfm+xPEv7E\/x+\/Z4+Etv4T+Dnw7\/tHp7Y9P4XHQdx32jdk7P5y8QeO+KMgxlH6hgMRmGAw2mK+raNW6PR6KV+krd1ex+m+F3AtDjbPf7JxGa4bL9dcViK3sW33et9LJ9PV3939YPDdxa+JCbnSh\/aA+v6fdPqTwT7btoKdJc6Jc\/wBm3H2u0\/0L7H0xgficZ6HBOPcZYsq\/gx8HP24vi1+zN4\/g8PfG\/wAB+K9CtP8AmJf2l\/QYHc9N3\/Ahuwn7zfC79pDwB8Y9Nt9V8JXlpf2eoWfX6en546c9SGxtX9C4e8Tsuzah\/B205t02vVRtbu791a9z2+O\/B3POCsSniPq+Jw+7xWGftaG2vSV9u2q2v9n8xfEn7bHgr4XePtc8FXnxE\/s8aeen5fTn8eexTJavCfBP7fnwT+OvjbxR4J8PfES8v\/sH\/IS03PTnHTnnn5sFs+oymzpf+Cr37CvgDxroeufE7w+P7P8AGen2nb+fbOf95PxxmvwY\/wCCR2neEvtPxs8beNdKze6Fr9\/pvb+oOMgdw+TnOzK7fwjjLjDNMg\/tHMM+xV8vv\/sqw22vTve1\/TSyekj9H\/1M4Pzihwxh+CsLia+YYh\/VszwuJobN77S63Wn5X93+grwz8W5fhX9m0H4RWv8AxK\/7Y868845\/4l+o+3bnA68dTuyAv646J4ktvEfhLSzpV5\/pYs\/fsPrz0GW49wN29fwH8N\/tM\/By2+0WmlXlnp5+uO3QfKQevPB9Ruya8u8b\/toeP\/hLrcHi3wTrn9oYIH9m9Sdvfpx+G7rjtilw54yfUcEsR7V5hhld7a3fy1162XyP0DEfRpzPiDEUsA8pfD+Y4Z3SlviO+t19+ttmj9Kf29Pi1c\/AHwlpdoP+PzX7z+Q684\/Lb+XSvHP2ePjZoni3w4f7Ytsf2h6D+uWOen17ocKV\/I\/9vP8A4KU+Ef2svDfwv8PWdteWHxA0O80c6lpv\/wBbqT+PuMYO7yT4b\/GzW\/Deif6Jecafe\/Tp+Z4ye30znFfj\/HXipmmG4w\/tf22J\/sevhv8AdfT733X2b9LrU\/r7w6+izlWJ8KKOX47L\/YcTWrfWcV37a2Wu+tvu1Z+\/HiTTfhL4JtvFGsf2XZfbNQs\/TH4ffOOMn+L3xyE\/Ia6\/a003wl4k1XwT8MdE\/wBM1+8\/4lvbt2GJME8\/xNjGTv5C+OfGz9rTU\/8AhG7i0u7j\/TL+zxx\/huOM5\/2uw3c5X5v\/AGZtbubb4kW\/xX8V6Xn7BZ\/8S3TPunP5v+fGeyj7tcuYZzi+LcRR+oVP7Oy+hrisVJXv3vslpdq+3S9nzejwp4B8L+HGRYvN83y9Zvm9df7Lhd7ve1tL2Vrrm+zfq4x\/pR+Dfwlth4A0vxt8b\/FH+mahZH\/iW\/8A1sduvfjoVyQvOfGz4tab8E9EvPFnwo8H3njzXPsedN\/s3j9ct1\/3BjnAbI2\/KfhLxbc+PtSPjb4x+MP+ET8L6f8A8g3w52+vUn14wvrzg19RWv7Z3wK0S1uLPSrfQvslh1Ofr7seSf6hVztr6LJfE7LOGMPrmGGy7D7a\/va9e997WXktZd9LH88eJngVxfxQquPr4T6zj+2GoujRwF7dXzL3fVfrH+cjxJ\/wUP8A2vfGv7SGl+CfFWl\/8In\/AMT3\/iZeG\/r6jGAB6fMe2TjL\/wBEv7M3jjxYdEOq+NtUJz6Z4x\/wLrn6Z7s3SL+Yr\/gqJ8Y9O8N\/tDeD\/if8Pfsen654g5\/zgjOBjjbjjqcEv+hv7D2t\/H79qzTdK8JDXP7P0f8A4l\/\/AAkniTTO3Ttzn\/x7p\/Dgu3TieIsz4o\/sriZVniKGJ9th8Lhd7Xfaytba19Ntbo8XwlybI8DkHE3CWbUvYV8uf1jF5piaP8Cgte2l9Oj173P34\/4aH+G\/\/QwD\/P8A2zo\/4aH+G\/8A0MA\/z\/2zrxP\/AIdv\/C\/\/AKHC9\/JqP+Hb\/wAL\/wDocL38mr0Obif\/AJ8YX\/wol\/8AKzzf7M8GP+hti\/8Awl\/+6H\/\/0v6QNb8E+Lfs0\/b07f8AszenfbkfwjBrxPUvBPiS5+1\/0\/P+XB9Otfpn4ludOttE78ev4fUD07j0xw7\/ADL4k1LTba274HHv37YyfoAn47Rv\/FMww2Fw0l+\/tbS2jduuq5bvre3pblfN+lYbEYrZ4fXu1\/X329VLTl+Fdb+G\/i3\/AEi7\/wDHvXJ7DB9u4+hxmvAvEngm5+0\/8fOf8\/jj\/PTg1906l4t\/0a8H6Y+v+03t1\/8AHvvL8l+Ldb\/0m49c57f5PX0X1GcKrfEZn9U\/yPZw31r\/AC+\/T8LHgR8AW1t06e4z+HBz\/wCPNkcfLuKRej\/C7W7bw5qX+R7475+mPbnBD+XeLfEtz9luPr\/9f3x09fz\/AIuA0TxJc6bc2V5efyI\/PlvboeOnO4FcssxCw2MoLDXbsntpf0crJ2v9p33091SuvQ+s4VL7+r1\/PddI9td4\/q\/pviTqf89jk4X6D73p8vLGutttb\/0a4u+cdOvX6+p\/D8R1r8oPjD+2l8NP2d\/CR8Q+O9QzJ\/y5aDZj7RqGoZ4569R2yu3++2cVm\/8ABLX+w\/2pfGHxV1W\/j+MHw7+A8c3\/AAmGm+A\/GHjg6jp8Goaxzd\/Ze\/WyHO4cf8hQHDI361nPGmHyjD0NvrOIdsNhXV9jvbXtvq22uj6csvzdf7TjKOCoUvb17cvT\/t7V66Wey805fZ\/QXQf2zPh7\/wAJne+DPE2j6x4cvN2j\/wBm69DcaVqPh+4\/4Sjp\/pQve\/IOBhvfGF+z9N1zTdSP+jXFnf8AU\/8AEtHTP4L265Hy9ABlhXyH4\/8A2Y\/gPomsWHjTQ9Os4\/7H6f8AcP8ATjnt\/D6HK4zXkvxA\/aHtfgn4b\/tTwn4X0HUfC4vP+JkfDfP9eMcdCemP3eArfBZL4pYnB4Zf6zfV\/YYdpfW8N5Ju+6tfsl03STP0jJvCTP8AiDMvq+X0f9oxL1WJfsU099NL9NNO\/Np736Y\/ak\/umut8N21rretaXo\/\/AEELzuevsByB29fbdjCfnj8BP21fgx8fvK0rwxrh0\/xZJD\/yKuv\/ALjUP++gw+2+n+gIPUAjiL9PP2Z\/s2peJNUvLu0\/0LT7Mf8AEywPx7lh36bv90ZLV+x5fnOCx2CpY3A4rD4mg3o7bq3zXmtI6rVs+K4h4VzvhfM6mT5\/l+JwGYYfbC4i33bO3ra3pY+xfDWiaJ4K0TFmMDJ\/LsSOuf59gu3FWf8AhLdNubm8H\/1zj16joD7Z68gEp5d4k1I\/abk2f\/Hnj6\/rznHf5B+G0M3l2panx3\/Ht\/45z\/5D9D2KfmGccf4jC4i+Huru7139NJaq3bv8Nzvy3hvC\/ev8\/Tpb8tbXD4kW3w38f3Oq2nizwXoPiDQ\/sfXn\/Eev95OnG75yvyFc\/szeCvAFrqviH4DfbfCeuafaf2n\/AGdxn8OBzn0V89sfdr6i+1D7L26\/p+eff09u9Z1t4j03w39s1XVf+PP7J2+UDHtznP19+er\/AJl\/rDisTia2Ix+t\/dV\/RJeW3Z+elkfoOHr4vA4f6vh62L9jo3hNb\/PT1\/mtv1cT8zPEn7SGm\/FHw1qlp4hH9n6x4fs9X03xJ4b1LPPXvxz7nk5A2Jgivy7\/AOCfX7E\/j\/4kftD\/ABk0fSrw6f8ABfX9d\/tLU88fjw3OOR8rDoMMrNlfrTxJ8LtN+LX7Vnifxr4S8Pa9\/wAIvfnPB25HT+7J0A\/un\/gPWv3C\/ZC+Dmh\/C7Q9U\/sbSzp\/f3+h6fmA2OpUcBvm8syahxni6tDN6X1nKMR5Pfbeza\/Dvy7o\/SMTxViuCcEsRlVZ4fMFd7N21Wm+l1p5dYys3H88vjp\/wRi\/Z41LQ5z4T1TXvCnjM2ecjsPb\/DYPx\/i\/HLUv+CS\/7dNzbap\/Y+q+FfEFnYf2h\/Zw1Lr\/ACGO\/OM9Mx8ts\/rj8W6lc6jrdx\/+rGPqOcj\/AHev8W0LLmaZrepWttcWdpz+WT+jjjj+9nqVGAj9+YcO8H4ivd5e8NQb1+r6WTd\/i1e337pO3LL6Th76RvivlGG\/2jNVj+l8Q7a6f4tNO9l1vflP4j\/gz\/wTu\/a30H9pbwXr3jP4P2dz4Xs\/ElvZ+MP9JtPs\/wDY\/wDy9XJzek9PRTjuX+5F+mP7Y3\/BLv4x6J4tt\/G3wH8P2eo6H4g\/sb+0vDf164+6PTk9OoH3kf8Ao\/8AtSf3TXcaJqVr\/odnyB+R\/kR2x97Bznaf9U\/B\/q7wvm2CWUYil+4fTz\/C+nl01tZH20fpX+ImGzPC5xRpYb\/oHxWFt+4xF9ur19bdurR\/Ff8AED\/gnj8UPg59j+IXx3t\/7fvNQ\/5Bvgjw3p5H67gBxjnGD0O379e7fsqfsPfHT9pnXTeWfhe88B+GP+hk8SaOfb3YHOTjlPfoGf8AsOtvDXhvxbqX9sanodnqX1HTHryvQZ7+3Gd9dtpvhrTdEtjaaVb2en85\/D0HzAgcdQV9MHAavsMN4MZFm7wv7x\/2Tt9W64jXpfTztez391ajxH04OL8PldXDPKcN\/b6f+y4r\/mHoX0016evrdu8fx8+Hv\/BIL4M6I4vPij4k1jxrJ\/yx0681D7Pp\/f8A59gt7xjoNQx23AjFeOfFr\/giv8C7nxLqni3wlreveE\/C5s\/+Jb4Z03Wfwx0Xr7H2GOa\/afxb4tttE1vS\/Dt1b\/2h\/b\/O3HHHvlTnjH8Oe3l7QjVvEngj7VolxaWl3j15yeemTtUdjwfxxkF\/bxPh1wh7Ctl+AyrDf7N8V019797ffbTazvzR\/Cl9IXxX+uPH4\/iXE2r7vX2Oj01V0n0tq33Vmz+DP\/gpV\/wSq+KPhu50Oz8JH\/hKrMn\/AJGTUvm\/9l5Jx\/dU9ynGytL9nj9o\/wAN\/sPfC63+Gml3lprnxBz\/AMTL5c\/T+IEEZ\/vH\/gJYlf6Cv2\/fDfx++Nvw41z4Ifs923\/E8PXnHbp6d8fdGP4c7iU\/h78Vfsl\/ttw\/F3XPhxYfDvxJ4o8UaH++1LUdAg\/tDT7e3\/4+f+PrheuB\/p4+oGCV\/I1wrjMRhlk+UVll+X4Z2t2fm7dVfWz8r6H9ReF2bcH8TZTVx\/GmY4bTFJvCt+yr4j\/p7W2Sd2tubvoo2P2w\/wCHwfxR\/un\/AL7\/APtdH\/D4P4o\/3T\/33\/8Aa6\/Lv\/h3P\/wVD\/6IR4po\/wCHc\/8AwVD\/AOiEeKa8H\/iGGdf9BmI\/8Gn6p\/b3g7\/0CZd9\/wD9zP\/T\/p3+JHi3\/iSWf+k\/z9T\/ANMm7ehJ9M8FflzxJ42xbf8AH5146f1BP8gTzjG1t\/beLtN+0+G\/yz\/9brgf99+hxha+VLnRNS4\/f+v5+vXrg9R9MdTX835l9a+sr0\/4f8D9ky7DYW3p3023Vra7b+7ve2lo1tT8bWv+kf6RwOemef68Dp\/jXjmt6l9p+2cE\/p+ecbfXnf8Ad4zyiZupaJbf2lqn\/H7\/AHff6\/eIH1Bf04zurktb1C5ttEuOx\/P3x0JPXj7n\/syeOeocB4kntftP\/H7\/AE\/9qfhnPvtbOV4nxIxttNFzj+n49D\/X1G75krzfxtqd1\/x+Z\/z9dp6\/kuMFXyDXFDxbrHie8sPD2m25vbzUNSsLKH\/sMah\/xK7X3A9OuevGdtZ4D+NL5foBwGpf8K31Hx+fFniHQ7PxB4o\/4l\/\/ABMvEY\/tfR9Hx1\/sHQNe\/qW9tuCX+6f+E\/03RPgl4n+x3l59tN54Y1I\/8I0eM6P+XXB55z+OKzvFX\/BDf9r\/AFjxaYbX4m\/B+yuP+fP+0PGfp6fpx0692r9Q\/gf\/AME7vhFon7O3jD4J\/FLXPCusftHmbxRFo\/xO0H+3jp+nef8A8gH\/AI+RZA\/Yv+X\/ADp8f9pcjbgb118TMuwfFnB+a8M4ethqFfENvrZYijP2io9G79dd103j9H4M5hguEuPeHs+x+ExWJwGHxP1nF\/VtbUP4fttZRb6b+e2vN+I1t+1n428Aalbn+1L0Weof8w7UuD\/I5\/ErjdxnI29tc\/EjUdSuf+Eh0u4+3+F\/EFmTqXhvHft\/Eoxn6fe52AAt+b\/7UXhH4pfs8fFHVfD3xC0PXv7c0\/7f\/wAS3+WT789vyILN5v8ADf4ya5bW32y08Qf2frg\/5huNv4c7sk+mcewwHb+QslyXGcPYPR4j6ve\/1V+79X\/Te\/SWmr6H+73EXBmQcU0KWb5WsNiK9fDN\/WcO\/wDePn7z+W\/V3sox9b\/aA+HmueANS0v4u\/BfXNYtJtLvf9DvNM1D7PqHhe4+mwgeudzepVcmv7Mv+CQX7Yeh\/tWfso6P4jmtxb\/ETw1qX\/CH\/E6zx\/zENP8AsVz9q+t7Y3thfcZ6kZXG1v48NN+JdhrHh7xDYa5eG7\/tD+z\/ALZoPH+gZ\/59eDj64+mcBH9s\/wCCX37QOpfBn9pm4+CFt4w8VaPb\/HjUvD+m2l5oF\/plvp32j\/j0tftWM\/8AP6T9767sAV+y+FHGdfKVVw\/sv3Cd3hUmr9v5dVbX8ErWl\/GP0n\/DvE8W8GUa\/s1iM3yfXCYq3790Nv733u3lbVx\/ty+OXxF8DfDG58P3Pie\/+x2niC9+xfbDYapcf+kpGPTG3nGSRhkr45+Hv7VGg\/E5\/H9zodl5ej+F\/GGoeD9B1L+0P+Ro\/s\/7GBdWvcWWb3P5AKuCa4nxt8LfgDp32yz8Q6peAX97\/aWo\/wBpeKdePHr0UDp2ZvvYCDaQnFeLbrwT\/YkF54S\/sEXmn3n\/ACMnhv8Az3988fQBfrc54jw2Mx1dYerh8N\/1CYarZ3W7V+t\/OL2fY\/zw4e4Gz3Fex+sYXEewxC1UaO1l0vZdrPS3kbXjn9t7Qfh6+of8Jh4OvP7Hs7zH\/Eh1G1n\/APivzx3wAcHd638Lvj98Gfjzpt5L8N\/G+j+KI47P\/TNNyBqGn8dDa88nBH3n\/wB4cCvxG+OXhjTfHl3O9teS6P4gkP2M\/bMXGkH\/AMe4x9F\/3+GRfx88beMPjf8AsyfEiw8beD49Y8FeKPD975uJbf8A4l\/iiw5Oep+22QHGMr9Xwdn4zhuKs7xOd1cPXlhv7PurdE9e130vv+N7y\/vTLvor8L8QcL\/WMBmH9n8TLDJf7R\/AxHy31t336q65f7sNEHhLw2f9Dt9B08fw8\/8A6s5P+0vH8RxmuktviRbab9oH2jIz\/P8AHqfYjPcDC7fzx\/Zg+NPhj9qL4S+F\/jB4Pnmt7bV7P\/TdNPP9n6hp+f8AiV6r9l29TjoM\/QDa3pXjbxv4S+Gum3Wr+LdU\/sCzH\/QS4x+hzk4HVMA9skV9SuLc1wqeHoUWu6tr6abPVdZX3XLqpfxNiODf9rrYDML\/AF\/XD\/VLrr10u36ddtE+Y+odS8f6b9pPHB7ccf8AoJHp90dT9zGxsn\/hN9Or84rn9uL9m+2+2f8AFxNC\/wBr8e\/TPr2OOpyMsul\/w3D8Abb7P\/xXmhdee4HHbvn3Cj12jAFeJiM7z7Zt+u+v3R7bJb9XZ832GG8OsZt\/ZuJ37denddNdfu0R+kP\/AAltt6fpVjW\/H9tbeGtc1b\/nwtL\/APp35zjOcYT8OlfAmiftjfB3W9St9K0rxDZ6hecjHT8Okn\/oRz\/fGd9e+\/FLx\/baJ8Hdc1W8s+DoOsc\/0P3iCAfX2PUFfDxOYcQ\/VqvsO\/N3s+2l9b6ba3s73JrcG\/V8ThHiMHbo7a37Lrp3ulv\/AOAfbP7PHi3Utb+H+l3mq\/8AH5qHQZ4BPU\/dyfyXpzjOF+kPtR\/un\/x2vzn\/AGFfiRc+P\/g54f8A7Us\/+Jxp9mPpx745\/Tr2xmvuG3ufs1z0\/LH5enAA\/vemFzvb+n\/DTOsdheH8Bhsw0f1bro9erelrdOvR3vY\/AuMclWGz3H4d0vYN4n0tpo72dtull0tI8l+IGpXP\/C4fB9n\/AGBd6gfseoD+0jz6+358++5c\/N79Rcf2bc232nt059eTn\/Vj9frxkJVevrsky3E5Ti8wxH1v26zDFfWV19h8vdeu\/R2+48TEYh4rDYSh7L6v7D\/Zl20erenVP+ZJdE7tH5q\/tRfEjWvBVzqvw9+FHh8\/2x4gs\/8AiY6jng+\/TnHTovX+H+L4M\/Z4udE+F2peKP8AhItU\/tDXNQ17Gpabj7v+eOy+5GNrfcX7Wv7TPgDw3qXiDw74es7M+NNPs\/8AiZf2lx369SD9R+uMP+fPw\/8Ajr8HvDdzqniLxF4L17xZ4n1G8\/4mWpfYB+QO1c9ugGM9WwGb8SzniHLFnVbD4fMPbug6qt\/y4w+mnS+tuy87XPsMsw2LWW\/wWliNXu79uuvZW2t9qyZ9+f8AC8\/BHpe\/9\/P\/ALRR\/wALz8Eel7\/38\/8AtFfJn\/DY3wm\/6Jhef+C6j\/hsb4Tf9EwvP\/BdW3+tGE\/6Dv8AylUPD\/sTEf8APrE\/+Af\/AHQ\/\/9T+knxJc234fTPy8+7Zx0wTn64IT5v8SfZrb7T+49s\/1z3z9B0yMYIX3\/xF3rwnxb9m+zXH+mY79f0yIsYPc\/hluWX+e8RXwrjutrbrX8G9dd7dlvzQ\/WsN9b0\/r+uh8qeJLnTbb7R\/oZ9Px\/DB7j+IY6\/NtKV4V4k\/tLW7W8\/svTPt\/I7+h\/3Wznp\/D2wSQFb1Hxt9l1LUrPw9o+ftmoXvH+eBznt0xnngP9ffD7wBpvgrw5pVneWePw\/T\/WMc44+8c8ZxuYLzZdhvr3RbdVbfrvK2q9fW9o9KxP1bd2Xdrpu19u9vR2872Px78SfszfH7xta3n\/CPaGdP+hwPUfwn8uPTd2r+kn\/glf8AB\/RPBn7P3h7\/AISv4Z+FfDHxU0KHWNB8YeI7LT7T+0PFGj\/29fappdzdXQ2nmxvbAE5GOhwSGbpPg54bubnTbTVdV0uzsbPHXp+ny49Rz69cV96eEdYttY02NLPRBo9vZ\/jYZ\/2enzE\/Tp1bpX3+W8OYbLf3\/wDl5\/aW2\/RS\/E+VzLMvrL+qpNe7ZuXT10S37L7rXNr7Of7S\/tT7PZ\/bPsX\/ACEs9\/pt\/HP\/AI7XhVz8CvhLp2o634ttvAeg6\/4o8QXmsal\/xUnXj8Mcjvh\/Qrxub3WsS66D6f1FcGYYbB4nTEw+s37eXdXb+d9bWVrGtCvXwq\/cVMXh9Ou\/pfS2ttF+NnI+Q\/Gf7MHwl+K+j2f\/AAv\/AOCfw\/8AFv8AYV5\/aWm2cuj2uv29hcd8C63de\/XP+znav54ftjf8Egv2CfGfwI8X6l8Lvh3pvwr+IlxDp8PhDx54P1DXj\/Z+oW\/ODaD\/AEL7F\/z\/AP8AoJJ7HScMa\/cK5\/z\/AOO1+Xn7V0mpfDrxDZ38OqeZ4c8QTf6HoI\/6CEH\/AALke+0Z6jG0ivzjj3MKHCWQYrHvC+3f+7b6t\/fq7aaq3Tr7v6bwLxNxRhc0o4fAZ1icB54evUVBW09WvkvK97R\/iV8ef8E\/f29vgwms6xN8H9Z8aeF7D99\/wlXg\/wD4qH\/R+pxaWp+2jjjmwXv1\/j+Qvh7+0Df\/AAl+NPw\/+Jn\/AB7XHw\/8Vafq954POn\/8hAWH+lcg7f8A0LJ\/2c7V\/wBK34FXPhvRPAF74i8WXlnp+hafZaxqPiTUtSxx+PHPBOMnrwOBX89mifAD9nf4x\/taa58ZPCX7Od5qA8efGD+09S1P\/hB\/+JPq\/wDbGvemCTz16e+MZb8szLMMjyvB5felfEZxdfVMMtOl\/Vvq\/e9N2f1DmXjHjPEbJMxyLi3C\/wCz5fbE\/WsPW9l7f2T\/AOXvfz17drS8l+Nn7SesePPDfw\/+Ivwu8YQ3Hgvx5CPGFn\/b2f8AiX\/Z\/wDkKaX\/AKLsH20f6fYYKe42gAP8l3X7XtzofAuL3r\/xLT4b1jt\/48B3\/kGJK7f7Mv2gf2Av2W\/jb8Pb34ey\/Cjwr4feSzI03XvB+jWmn6j4YJx0NtZHj7wJKjaRk7gArfyPfGP\/AIN3\/wBsnwjrGsf8IZ8RPBHjXQ8\/6H9s+16bj6c3w9eMNn\/nocgL8Pifo05plOMrf8KH1\/D\/AL76ri8Np\/5Nqul\/h30dtGfv3gl9K\/w3xPD2EyjibDf2dmGHvhsX9Yo9OzfR97\/ozyU\/tXeHvEK+bqOoXmoR\/wDL510\/UOmPYdv72P8AZO5jXkvj7xbofxJ03VPD3\/CQWeoeF7A\/2lpp1I\/2PrGkH10DXv8AiqFPfj5fq2Nq+yeD\/wDggt+3\/rmpRW0kfw+0iP8A487y7\/4SDVOenb7AQMD0Le+c\/LW+Iv8AwRw\/aB+FHiT+wfE3jHw3d+ZZ\/wDH3Zabqn9nnPT\/AJdTjHPoPXoRV4nw0zXKU8diKz+r31xTs9tNPx2asr7393+in45+CeKxH1fKeIMNa+i\/C3XZvvL8GZn\/AASm\/aEv\/wBm\/wDa80r4R69rn\/Fu\/jRNp\/hvyelgfE+of6L4Xux\/ptqeb28\/s\/8Ag\/5CQHPWv6Uv2tP2MdN\/ai8N2\/h3xD4wvNPsrDjoMEevDN2zgY6djtNfzbfDf\/gkv42ufEdnZ+IfFOv9tSzpun5\/Dt2926Y2rjNf2IWPwW8R2fwr8Npa6hrN34k0PQdPhu7y9\/5fxpx\/6euT+OB69AG+8yV4jNaFbD0aTxH9nu2L8vuUbPS+7evVp838PfSF4q4Xw3F+VcTcN51h3mGIoLD4rF+y5ml8rXb9YWvbXc\/ni8Sf8EK\/BP2a4+x\/EW8\/Efr29u\/OMcZr8u\/jp\/wSY\/aZ8AfaP+EJFl4rs\/qT1+oUcnuQn4bcV\/Xp9m+1D\/TNUx6D19e+OuOoyeoLAFVLn+w9EtTd6reWf\/cS549OuP8Ax1fqcYX4WtxFRs\/YZe290v1+wvld9HqcvDvi1xfg8TRjiMWswW7VvTTay166J\/y6Xl\/LN+xN\/wAEu\/2l9b1LS\/EXjbxheeBLP7Ye39cn1z3z0B28r\/UPq\/wi0e6+Cen+ATrl5f8A9h2dhD9smzknTvXkADj3656coeGviB8N9a12y8PaV4o0Ea5gf8S3TeO397JJ79vqTuDt9IHwVqWpaJeYt\/8A6\/6D65\/RcqK+fedZ9evfCvEPS3m\/uSfTovO+vN5nHniFjOIMbhK+P+rZf9XxX1jC4R7aeWjtt6edmj0j9lTwlbeEvBVn\/o\/9n+mRuxz1\/hHTA6NjqN2SrfTNzpvf9OvT1xj16c\/dwM4ITzfw2bbw34b8P2l5nT+v+cbf02jvnbjYvb21z9mzkYxwf6dduPw\/MZxX9a8KJYXI8JQxyb27266teuv2rdU7Xl\/HfFeJxGOzPFY3\/oJ+1f5PdddVe8E3ta\/vVftNzbfaLTPI9P6fezk98jOcYGfl0rXofp\/U14n8e\/jToHwY8Hz+J9c\/0yTobOyGTqHofTrxglcepJ3vieAPjZ4J+LWiaX4h+Hniiz1uz1C1\/wCQbxxnp8xxx6navTnGQye1iMRhso25muvl+fL92nW9mcuGybNMRg\/r\/wBTxTw\/+7\/WX5\/JO+lt4+VrLm4n4tfAH4FW2t6p428Vf8fmud9S\/UdAexPT0IyTXy9qVt+yp4b8\/wCx232\/6HP9FzzkdffjmvUP22dbzbeC9H9ejDoPbkqOnzdc+pX+L8zdSuf9JuP1\/L+pz2\/766V8JiMuyvC4yr\/Z+Xrffe93ppZW67X9Qw+Ixf1ah7fF9rq93\/evqr7K1\/wveP1r\/wALC\/Z5\/wCidfq1H\/Cwv2ef+idfq1fFtFeV7KH\/AEDr7v8A7oeR7ar\/AM\/X\/wCBx\/8AlZ\/\/1f6LvEtzpo\/5fP8A2X8OPMXqfu8+0jfdX5m8bf2J\/pmbz9cD\/wBFdsH0z8vAw+3601L4FXOo+d\/pn2Dj\/wCv79uOS3ptXGKxPDXwu+BXhvUry88WeMNB1C854\/tAYx+h+meSMHjdsT8AxPCeLzPEJ+wWHs7PTpL7Vvdu\/nG29n9n9nw2Z4bD4X+NrdbddujT+68bau7tynxz8HPBOm\/8Jbd+Lbyz\/wCQdn+zuPpjq3B49WzySo2sW\/V\/wT8LrbxdbW\/iHxDb8djnk546fl3GOgzyV5LwB4A03Ubm48Wi3P8AY\/2w\/wBm\/wAOMe\/Pr746\/NkmvrXRLm2\/0i05\/wDQc5PX7j4\/Pj0b7yfofD2S4bKKDw\/o73vb1+G+1t\/uvaXzOZ5ksUujWnM73+a2t8m\/wsbXgnw34buba8+12ZxYXmRpv\/1ufrgt04Gcg16Rc3Nr9mNpbZ6\/Tr6AKVzj3GccHjc\/N+G7b7Lbap\/1+aPtH1x167SPYN6HbyV4HW\/H\/wDYlz9jtLP0\/wA\/\/rLfqa83iriLD5SmsRWvhuy0aXX8V5Por393my7JMRm+Jbw1+9+tkn3tf7tPPWJ6jWPVK3uPtP2f\/ny+x49Qf5Z9CNo\/Hk1a\/wCXW3\/0nHqNo5\/MHHuATnPG3d8nmfWfrV7XtvZ33v5Wv6pQ80yfqq\/mj\/5MVrn\/AD\/47XgHjz9nHwf8XPEmn+IfHd5rN5Z6PDiz0L72n5H8Xv1z94e5OENe\/XP+jfbRc+ILPk\/jzyT04Hbrz7Y+Y\/tLTbf\/AJePr\/gD8n1+8nvuyErxs5yXLM2wyoY+j9Yw3\/QJidbLpps3d+S19JR68HiMbhr4rDvR9\/lq73X530vbSUeb0f4S\/D2y0f8AsR9D03UdLf8A11nLp+l89OP4e\/p9Dtx83b6J4J8N6J\/pdrb2dh\/Z\/wD04Yz+Hz7CcDn2xtfh082ufi1p1qLi08PaV9v\/AKe2QefxBxjAIzhjTfiR4k+041aztPsY\/wCgb7\/iB+p78fdWvMw2c8IYHE0MPQo\/7Re3+zUE27X\/AMLV1ZW12a0vzS9d5bneJwy+W9VWfXza06W++1j0C51y2+03Cn0\/+v3Xtx1z6ZTO58y58SDnr+P\/AOsenbPTPH3a4rU\/G\/hu51u3s\/tH+evo2Mf5JzsXEufGuh\/ary0tP9P\/ALPu8\/n\/AMCGOfQH1wv3K6f9acJ\/0FYf\/wAl\/wDlZH9i4n\/oFR2tzrdz69e3B\/EnYuOO+B1+7klVxDbW11\/x+W9nff8AYS7dwSRgZ7fcA9SclqxLnx\/on+jG38P3nXp0\/qACSfu456nbnbLpf8JJoem6bb6rqv8AxL7T7Zx\/nb6E4+983POCVxfEODxXtl9db+r6Lr2b2v6acvnfYz\/sWt\/0CotW3w38AXXiS38Q\/wBh2mn3n2P\/AJhpz+HoeuM7Vz3xgBe18S22m6j4b1Tw9\/yD7P7Hnr\/M4JJ\/4Cnp82d9c3\/wkmm\/6N\/pncf\/AF+w59f5Ctv7TbalbXC\/aCc8emefUZ5+m3GOhzhPSw2Jwn736vf\/AGh3s91fv10Xe\/bpaOGI+u3w\/wBava2l97+fT+vQ\/Cv9p34Z\/Gab4qWdz8LvEn9j+HtQ0e3s9e+2cwf2hp\/bsCR9RwoA3Yw\/jnj\/APZd+NnxS0SztvEPxw\/s\/wBtN3H9cDHrwsnoC3Vv1y\/aH8N2+neCfFHjb7P\/AMgO80\/+0vf+f48L6jOdzfBfgD4k6J4+1K40rwneHUefy\/Pbz+PPbAG5v5e47y95Tmlb6vhMTif7Rf1jXtd9beqS0XZO\/u\/1Z4c8c5m8kwf1ek\/+Ef8A2f637Hsui0266vb\/AMD8B\/ZL\/YwtvgD8Wf8AhZ3jX4iXvjz8MfjzuP5Yz1Oduyv3etvH+m634b+2f8g0+v09uh4yOGHvuySvxfpvgnU+t5b8Z5x\/ngHHpk9iuPm9a0Tw1421vUtLJz\/YZP8Ae6fgdw9ejEdxtyyr8VkvGmeZR7XDUcttd9Xda6+XVdvuuc\/FuY\/6242lm+cYv\/aLN2WuvR6crVu6fry2Pmf9r62+Ov8AwkngPxv4S1y8\/wCFf+H7z8\/r90Z69D256gt+r\/hv7TqPhLQLn0tNH\/n9D9D19dpxuTmrbTdNudD\/ALH+x\/6Gff8AU8fXj9Rty3Sab3Nnef6H9j\/5Bo6Y9OuR9QxPbI4Ff0LwFkuHw9fFZh\/aGJxCx\/sE8Ktfq9ej8+\/S2vld835txnxm8\/yPJMn\/ALOwuH\/sf29sVh7r2\/tvlJXX+L1ve8fyy\/4KW\/D3xzeaX4f+Ifg\/7Zd3HhMf8gHr\/aHn9B29xnjGP4sYbc\/4J4fC7\/iSar8V9U8H2fhPxnrlnp+mjTdO4\/XAP4nP+5k1+i3i7RLbW9NuLS444wPw\/EZ5J7DOOoxiuT8I2uneErm50qzt7P8A4mF368fzyPbO77v3Oa+hy1YfLs2xSra0a+usXZ93u7312X3Ht4jjvE47gLCcNYekksPb\/asNyt\/V+t3bTd\/Zlfa0bpy+L\/20Ln7T4t0q0x\/y5\/h169eMAdn78g5Xb8CXNr19B+v\/AI7x6fxdff5vs79rTUrq5+JGp\/l39u3GOP8Aab14yQ\/xfc\/8fN3+H86T\/wB+reh+OV\/93o+r\/Uy6Kk+0+\/6Ufaff9KZkf\/\/WP20P+C3Pjbx\/puqeCfgPb\/2f6aj69x\/Ev6qTzj5MIzdd\/wAEYP2J\/wBpn42ePbj9qT9qXxh421DwXp95\/wAUf4a1LUM9OeT\/AF5\/A8r6z+yt\/wAEgofHnjyz17WtP+z+B9HvBnzTk6jt\/wCArgE44w344bd\/Tt\/YeieCfCVn4T8J2f8AZ+j2Fpo4+gz043ZA9Q3PUqcgL8Zl+HxSX1nELfl36r539btadnb3v0HEV8O17DD+9rZLrbstHrt008rGJqet22mm3s9K\/wCPP0PX6dP5hR\/tc7K9R8AH7TbXecjr\/wDq77e\/8T\/Uchfni2t\/tOt29pyPTuPfnHXGO6Y67Tkon0Polzbabn+Xt+R56Z5X8eRXfhv6\/wDJjxj1vRfs32XVP8j8f8\/SvG\/+EJuftNxd3n\/Hn9s\/tL+zewz1OeT3\/H2xmvZNF+zfZdU\/yPx\/z9KzP+Ek0P7T9m+0Wf8AXHTpn14\/p3r4XjLLcJjMZS+vVuZvs7arZbW667b3bbR7fDeJxuHw9Z0N3q7WVvvcultbaedlI4n7Nrf\/AB9mf\/jw2f2dn8+Ppx1bt97qKJvGd1o9hLf+MLjw3Z2dv\/y+f9A\/26jPt93PfrtrN8beLfDeiabql5eapef8S+z1HUunX+X6O\/X7i5XZ\/IF+2T8b\/wBvf9t\/Wbzwr8IvBHxC8F\/BaeHUIbOf7P8AZ9H8U6f4X4\/4mmujdZcXv2A\/YhqJz0\/4muMN8FmOZZVw\/wD7ti74ivd\/VXW+99NvTT+99r+kPB3wUzPxexFV4jF4bJsoy7\/esViP3P1fpZp813d7+ui0cf6kfCv7T\/wz+OupappXwW8SeFfGl5pc1vLeXln\/AMw\/7Rj6Dr6HkcZGAqepab8NtSuba5u\/EN5eajedfQe46gD8eRjoejfw0fsYfGz46\/Arx\/8AZNH0P4weE\/GoGf8AkR8fptBHfqW64AXAav7jv2b\/AB\/qXj\/4S+FvG3iHS7zwnrniCz\/4mfhzUs4x2x24\/wD11OX4fB8UYhYjEVvrH\/UJr+NmvyV7Xum2j0fHXwlxHhHiKP8AZFX2+U4hp4XF\/wDL\/wCa2\/F\/K15dHbeAPs2m\/ZOfthu9R6\/pxyBn2PHU5xmj\/hCdS+y29p9oswOnfj37fTgc++N7ejfaU\/57in\/aff8ASve\/1Y4c\/nf\/AIAfzb\/bWdHkmpeEtb\/0b7LZ4\/0zI\/4mPXn5efKPGO+xfQhj89c3c6bc6b9o\/wCJHr3+Gex6ce3fr8udq+\/\/AGn3\/StD\/Rq8DEcD4PE\/7vivYV+l2n59bbpLorXvaWh2YbiLEYRr61R37\/f2T9Pv02j843NzbW3\/AB+XH9n\/AF7\/AKn2\/luOMvh3P2fxJ9t\/48\/+Jfede\/PqOR267j6FVyRX1F\/o1cnc+EvDeo21x9s0O01AZ5\/TGfnBPPpjGcEjCivPzLw6xkmlRxf1hp6rdru2lJ3v53v5WTPXw3FFDm\/hYlbed3q\/5lfXf4e9rJo8b+zXP+kf8xD\/AJCH\/IN\/\/WOg9h9R0qsPtP2a4\/0zNn6bs\/n+65\/Jcds4bd6RqXw28NXX2v8A4\/NP\/wBD1f8A5Bus4P06DP8AkfNnc3zh4\/8AElz8P\/Enh+zu7zXjZn\/iZYP\/ABOMe47enbknOF\/i+TzDhvM8FhliFR+Ff8w1bVLouz6\/y+Ur35fUw+d5didfvuk7r8Oy2St\/eTPnf9u3WfHOm\/AuP4dfC6O8\/t\/4kalp+gXmvfb9V\/4lGn6f\/pV1z8mOnBIb2J6N+dH7Cn7Lvi34JXOp+IvG3im9+26h1Ueo69uMnv26YUKS36QfGT4keAfEni3S\/D3h7xR4W8QfYNBGf7O5\/wCJxrH\/ABPv+J9\/q+3fcMHnavBbwHUviBomnXNxZ\/aOBgf3v8OnTp37YFfmfFWZYnA4tYfD07\/V\/wDZtVb\/AOS6apaX8to\/svCvESwvB9bIcPeh\/aK+sYnFv+PpdWemun+H9D7q0TxLpv2T\/S7j1zkf4f0\/pXtnw+8f6JdXNvo+kdevUfh3GOOcBTzxnvX4+XPxs03TftBNxn+noccg9s5PHXt81r9nj9tnTbbxtqmleIdDvNPtOv8AaXTp74Hf24xjByDXyuW8aZnlGNo5h7Fuh\/zFvblv0v37+75K2x5mJ4Vw2b0atCgvrGI6tv8AC+iad+i18rI\/cP4yeP8Awl8JfAGqeNvEP\/HnpvAzjPPtt59Ox47YWvm34KftIeAfjrolzq\/gnXP7Q3ddN4BHoc\/OD+S+gxmvWvG\/hzw58fvg5qvh28\/4mGha7Z+nceh\/76\/ujsAclV\/DT9l34S+Nv2Z\/2zv+EJ\/0weF\/EFpqB78+46AY6dWz1HQV+v8AGfEmaYXO8kxFCj\/wkZg6OFawuvsK23nb1Vv1PU8MvDrhfiDgziVYjGfV+LsnVfE4TCYm\/wBXxFCj36XSfbW3S6P3d8XfFLwT4b023\/tXVLPT7y\/9P8hvY\/pjkJnW3\/Ey1HS9V62n2Pr\/APXGT\/wHLdP4dzBfjD4k3PgvxN4\/uPG2lH\/hLP8AhWd5\/wATLTdNyP5knPPrj0znFfXtv4\/tda\/4Rf8Asq3\/ANDv\/wCxxwvJ9yQ5HA9z\/wACydn3GJx9DG4ej9YrOhXw7Ud2\/v0XTy6Wd7M\/NsRkrwODo+wo\/wC84ar9avbp20fTy162+GPwD+0f\/wAlK12vlS66D6f1FfRH7Q+pfafH+ucHn8c\/Thdvr\/y0xnPz7dr\/ADfddB9P6iuuh8Nb\/uJ+bPgqf8Sj\/gl\/6Ug+zW\/vR9mt\/eoqK2POP\/\/X\/sz0y4021tvs2l2ebPT7P3A7Z6bgO\/PPT7oydvN3P\/EyFx\/XkfTnaG692QdhjBas3w3qX2b8OPqPwHP\/AI77bvvPo6lrdtbabcf4\/wBMAD82z1GMmvnz6g8d\/wBJtvFxuvtmbLt\/Tpk9c8YTGRwcll918Nm2ufr\/AA8\/55yBnj2+XOW+ZtcuPst0frnP1PplsYwcHLf8B+83f\/D+5ubq7szzjv8Aj3\/hBxz3XjjjiuPDf1\/5MbYnC79l+P8A5Lp2+1v5+99aab9m03TdU\/Lj5ePfO\/PHpn0wcoK83ubfTdb1Kz+1n+zrz7Hp57+h46Ed\/f0yQWFekabx\/atp+u4\/meAf0G3oSzFXby651PRLrxbeWd5\/xLzp9no\/r+vGPQ5z9AuCX\/O+PPYfXKPt9tL7bX18u3n26HucJ831bEew5t43v\/LZfjaxZufBNz9mx9mGO+T\/APq6c\/XrzjFeb\/8ACJaZpv2i0Olf2fj9eenBP0zt59vndF\/aH+NvhL4FeAP+Fr+LLzXv+EX0\/b\/yDfbqcHnj1Vv++sbW\/LG6\/wCC5n7HH2a41i88P\/EjBvNH8Nrnn9dzceoJ46ZOFr4rMcvytYmP754atbX6xX6PTfXfXpr1vZn7jwXwbx1xRg\/r+VZTiMfh9NcNRvzLr+P\/AAL3vH9O\/wDhCdFudSt\/9G4sOvT8cZJHp12+oxjdWlb6H\/1Fde+xdv8AOCR+X5cV8qaJ\/wAFBf2FtS\/0sfGzQdDvPEH2DH9p8HtjnoPU\/e9Cq9H+z\/h\/43+CfjT\/AJEnx5oPjs39nYcabf8AB+v3vfqU9QB829vhbM8Tu79Varb71d3a6\/rvLxM4\/trKW3m+XYqh1tiKLd79279dbafO6K39l\/8AX9Vb+zbj\/n5P\/fH\/ANsrc+It94G8APFrHia4+z6fJD\/pl7\/aHkW40\/T8+pHGB1wvPBCgYXiPBPxI+APj+20u78PfETQb8Cz\/AOJb\/Zur54\/PPPuB6DOdzT\/qBmuJte93e9qz6b2vZLV7Wfq7rl+YwudJ4b6x7K2Gt1o2+fw62fT7mr3l1v8AxMrb7YfP\/s4f56dQMcccH3OcVa\/tLxH\/AKR\/xNM\/z+nQc+g4\/HG5+s\/4QnTbm2g\/4qC89vz7\/dB+vPsDku9r\/hEbn\/SPseqfhqR6fjnAIHT5ue27Hy8v+r2e4W3+83vq8PWSa\/S\/VWv3W1o8n9s5P\/L\/AOlf\/Kzz77T4j\/571p23iTXLa2uDn7fjHb39fnIx\/wABzz0zivE\/iv8AHX4dfCZ7iz8T63Z3F5HD\/ptnZwapqGP5du\/zZ6bB0r5w0r\/goj+y29t\/pPjgRk\/8\/uja8foPujODn\/lqoPbbu2L8tXzv+yMQsPiM2eHr4bZ1623X+89\/73peyZ+g5LwLnfEGGpYjA5Licwof9BOGpeum7uvut23P0EtvG1z\/AKH9rt\/6f59Oi+v8RFcn8SPEttqOh2f9l6WP7Y+2D\/647flxjnk4U1wHgD4o\/Df4o6J\/avgnxfoPi3R\/+gjpuoY659mGfqOe5I3BfQa7FxHmn1f6v7Xd\/wA+rt58r117\/J\/a+VxHDuFwOIX1jCOjXw+y2Sta3Tfbtbdt3Pxp\/aH+Ctz8PdX8OeOvPH2jxhDf2V5zznTvsPQDb1+2+h9Pn+5L4T\/Ydzrlz6f\/AF\/yzj\/d44BJxvf69\/4Kha34hsNB+H2jeGND1jWLPztQ17Uryy\/5h40\/\/Rcnnn7Z9t6YX2YZDL+c\/wCzd4t8Sa3c3H9qj\/I\/P\/0Jc7t3JBRfj+Isnft7pf8AMPr06arSV9156a66qX6dk2FWLyL6\/d3+7y31b1\/uruu0Pprw38NvtNtjPXjOP0+8P5n2OMV9DfDf4FaJc3P\/ABNLey1DPQe\/1OP\/AB0N7MuQrZunfZfstv16\/wCffPXOO\/419EfDfxINOudL6D1BB\/Tlvfqc9stn5fkMNwo8Riv9poyu+3ut83n6a7ad4683l1s7xGF\/3d2vfbW+nbXze\/ry2XL9MeEvElt8JfDdtpXiL\/kD+n9fveo6Z7Zwc\/P3+paJ4b1H\/irRbWeoXn2T\/kJZ9cc\/xe\/8X1Bz8vzf8QPEnhvxJ5\/h27zqFn9ex78569xycfxLnNVrb4o\/8IB8P9UtAf7Q0fT7T\/iXe3uuQM9uv5jbhP1rD5j\/AGR50KGFXTt5d1slza9b7ngYfDvE+x+rv\/b8RitE1ft6tWvfaV\/k+T8GfG1x+0N+w78fvHvxCvDr3iD4LeLPiRf6nkfr6nqfTtgZBr9uf2Zv2h\/AH7Q\/hvS\/EPgo8f8AQN7D2wW7+zHp\/F0b2PUvBPhv4tfC7+yvEOlWeueGPFlp9M+mPvev945745C\/Cv7CX7G0P7MHxF+Jusf8JJ9ptPEl5b\/8IfZ\/9A\/T9OPbk5J+2\/3e3AOSzceGy\/ELG0cRRqr6hXX1i2Jd\/YO\/+dtE321t736\/xJxVkPE\/B9V4+j9Q4nyf\/Z8L9WfL9YoPfe9\/lbWy6+4fGy5+0+P\/ABR73moA+\/8ALp7kZ6ArlmrxO66D6f1FekfFrUvtPi3XP+vy\/wCvAxjryG7Y9c5+8uVZfLvtNt\/f\/T\/7ZX1WG\/r\/AMmP5YLlFZf9p\/7VH9p\/7VeoB\/\/Q\/p\/03U9b\/wBDuvtn+hYz0\/T\/ADjHXnBNeo3IudS+x59Pr1\/LJ9Onpxk145omtW1wLPSvTr2PToPv5xgf89O4yc5Xv\/7RudOtTZ\/r9ffvj09uOpr4rC7v1\/Rn3Bzni77Tp1t9kwffvnHp1JOD6rjvuwS3unwU8Ja59l\/tjVR\/Z\/PoP5cHOfwGOM5zW34B+HVrbJ\/b2tj7Zqn\/AD5\/9A\/2HLHP1AA\/2sEP7r\/x9fUfyJ\/HGP8AgQI\/uY+Xqw2H\/wBp+s9LvT+tb766+at7kfNxGZR8vLfvfW9tvlfpy+6itb6JbW3+l\/6Z+v4f3vfndjkj5sb1+X9S8S6jrX9qXnT7Beahpupf6Bj\/AJA+vdD95s9B\/MN99fru5uTbW3THP09z3OMcdmz\/ALGfl\/Kr\/go18W\/HfwH+EX\/C5fgt8O7PxpcWd5\/xWHnZH9n6eD\/vMPp09y2AzfkXiTh8RimsRRrfwH\/tWE2\/rVW6bbK1j9R8LclxGb5nhcooX+sZi\/q2ExeJ0tp5Xvq9by77tNy991ZtNvPCWqeHPEOhw6h4T1iz1CLUtO1L\/SLe\/wBP1D\/RPsvUAfie+NvBZPyP\/au\/4JOfD34o6rb+Lfgb\/Zvw7js7P\/ifeA7wDTtPGM5Fr9mPPfjB\/DAC\/l1on\/BX39pzWLz7TN4P+Hsmj4\/5cm1S2\/TbfgYz6DrkkY+b6Zuf+Cwni251Kf8A4SH4D2moWf8AYOn6bjGMfj82OPfjqWOdi\/hH+ufD+Mbw+PrYb\/Z1\/strv\/5Fu9+19Nb2fN\/oPw14F+NnhxjVmHDT+r\/WW\/rWFw1X7l+G9\/uuafgD\/gk54P1vTLN\/iR8RLOze2P8Aoeg6Dp\/\/ACDx\/wBfV13+th9C3O39nv2Y\/hJ8D\/2eNLuLb4baJ\/p\/2O3h1LXf+W+oevqecev55zX5z\/Af\/gqD+yh4qsLPwx4nkvfhXeXFpcS3n\/CYaf8A8S\/\/AEfof7e0vdknH\/L9\/ZmOm5sFK\/Qb4b\/Ej4b+P9S0Pw98KPHnw38WnxBoP9pf8S3xRkj1\/wBk5A5BHtk5Jr2cmxOMw7orIqX1jES0\/wBmXt9fR8vbv6ctnzfkni3iPETHKtgeLauJoYe\/+6SXsaD8re8+nf77+7mft1fBz\/hrT4S3HhPSvGF54T1zT9e0f\/iZabzjt93v6cY\/4EVO78WPAf8AwSl+IHgXxtB4w8MfFjTbfXLOH\/Q\/GHjC31TxT9gHcf2Dpd54Xz6\/8f4x7EbW\/pLuPBGt2tz\/AGVafYzeCz7j64z1z9CB3+9gK3kfi34heBvhtqWn+HvGHinwroesah9o8nTtS8QaXb6jc\/Z911ql1pdpdXhvR\/oPf5ucD5cM7ems64uWLWPxFJ4evh\/lv1v5WeuvfS9j8z4V4yxmUZJW4Zymr7ehiXWxLwm+iW93G7u3282pWXL7t8OfiBa+GPCmn2nibUPtmqR2ZhvL2y03Vrj+0NR\/B\/Tjp7\/Lg7vWbbx\/omo\/6JZ6pZn3x2\/8dzjd0xnpnbkbvgzw3+0z8CvF323\/AIR74oeFdQ\/s\/wD6mj8ugXt6b89SF\/h9kttS03Uul5Z6jx6Y9+On4\/puyC3qLxBz9e9XwuGXRRfL9+7tfTT8Xqfn2I4DpYW3t6OIwz0Ttq1bfstddvPR2SPzn\/aE+FH7QNzYanDZ\/DvWNYki1LUPJvNBOl6h\/aGn6h6Wlqft3r7deDldn4afELwN8UfCP2iG6+F\/xgstHtJrf\/iZS\/DfxRBp6\/2hn7Lj\/Qx+PDZ6ZTOK\/r+ttR1K14tNUvPTt+fOcnDD5sgei5JNWf8AhLfEn\/Pv4W1\/6d\/fdtPPHTHPXccEN+H5xwJwfxRjHmGOxmZUMQ7a\/wAe97fZv09e2ruz+qPDn6R2fcA4P6jQyPDZhh3r\/H9jpvpqm7X119bbH8WHwu\/aK+M37N82r+JPhdJeW9neXmn\/AGvTfF\/gbVRo+of2f7W19aY\/vf8AISXnOcY2N\/Tb\/wAE5\/22PBP7cPgnXDZ+C\/FXhP4geA\/7I\/4STTunH1AY9\/YnHGCS6faFzremala3Ft4h8CaEfyX8P4+Pl4yE9ecpvseAPDHwZ+HWuax4q8F\/Dz\/hHNf8Ww6fD4wvLP8A0f8AtA6f9t+y9iePt2odCPav1Hg3Lsg4f9hh\/wC2\/rGX\/wDQLiKPl01T27penSXw\/i14pYTxKwVXEYjhp4DP9fqmKwtXe\/bR2t6q263UT5d1vxvdeJNc8UWd3pd5\/wAS\/XtQ03+zfEnT+yNH\/wC\/OM5\/iJK4xtbG9fJdS+Evw31G5uLvSvD58P66Op03pwOMcsMdeh56c87fW\/jXqX2Xx\/qnP9n2f2PT\/wCzf4hz2H0x6EH0IHzebW2t\/p74\/QAk9+4I7Zy9ft2Hy7BY7LKOHxFHDYjD3tddn3v+Oi20tqpfzBhsTiMNifrFBfVuzb5d9Htf81ftocT41Ot6JbaXZ\/2Z\/of\/AEEtN7fjsbPQgcjj+FsDbW\/4STUtS+z8cfj69fvDjnHBHqSuDXsdtrdt059PX\/a46Y784PpgEF68l8beG\/7O+x6tpf8Ax6fyHpnAz24wvp2G75DOeEsThb4jA3dB63227fnv6tX932slzmhi0qGJX796a6PTd9La\/wDAu9C1\/wAJL9lznp7fN9D\/AAe3QA9sHADWPDPiTTfiB4J8a+EjcdP4fUDr9MEZ657YGQW4C3\/0m2+yfh0x079f6HH\/AD1Od7c3omia34A1LW7u0POoXmeRj8uV4G7PzbuwA+8G+OxDxGJ\/4TrNrv1af3durVu7+1+o4bDZVh8HVxG+J\/c\/VdPLvf5acvbsz7r+Det3PgDw3pXhLWLz+0ND+xH+zdS578c8Eev9w+7E5rvvEtzbab4k8P6rZ8+w6Dn8QcE+q7eo35+T4V0T4o3P2a3tPs\/t9MfgQe\/XHpnhq6S5+KNtolzolrqniGz\/ALD\/AOJwFzyNI\/Q+3ofVRkhtcNk2J+r0Hh6X+7ex0v36vTXTXZaa2duU+BzHEvDV6uIxGrxKeu1lvva29t0r+R8h+P8Axta\/8Jb4gwf+XzWO\/wDvf7APXHf2+9yvm1z4\/wA546+\/r\/8Ar\/8AHfaviD41\/tV+G7b4geKLbw9b69rn\/E+1D\/mH\/wAyFGMHjhcHnoCEXwi4\/aP8Waj\/AMgvwhe+\/wDaQ4\/Ac4PP3gW+7zwfk+8y7JMTp\/s2\/wCHotfJq78tdz83r5hhH+t7b+nMltr08tve\/UT\/AIT\/AP2KP+E\/\/wBivyX\/AOF2fEj\/AKBlrR\/wuz4kf9Ay1rv\/ALPxHeX3xPP\/ALRl2\/GR\/9H+pr+xP7O9+n4Mfbrx8p64\/wB2rXwmubnxJ4\/1S7u7c\/2P4E0EajxnOeOvOBn\/AHv+Atu2Imh6bc\/Zv7V\/0zUPp+RJzu6gnDBePfBK+6\/CTw3pttc+KHtgftHiSG3Pr\/yD\/tvGN\/OP+AZ3ZJOF2fOYfDJPf\/gf+k7afZ16ctmpfXYnFb9n+H\/k2vf7O3l7voGm6ln\/AEvpzknIPT\/gK4Oe3zZI7gI7dtptzbXOP59f8M\/mvcZGfm8StvtNt7jOPz59VB\/8d92bGW7fTbnv+X\/6unbu3AbuQTVnknW+P\/tP\/CI6r9k9f\/1\/h\/8AX718uaVY6lrFhcaJc23\/ABK4\/tBvPtnTUO2PfH97K\/8AXMBgF+l9Suf7R03+yri4x+n4dMfrxxkc5ryi28AeJNO+2D7TZ6gPtZGMnHt\/yzGfpgD6bg6fhniBkeeYjPaOYZfQxGIoPDa\/V9Ov5Poml0V+kv1HgzMcHg8rrYevVw3t3iVpZtO2297L\/hla9j4D\/wCHZX7FtzrEuq3P7P8AaH7RN9s\/0PWNe0\/T\/b\/j1vLKy\/X88Fm8t8f\/APBIX9ifxaby7s7P4kfDMevhvxR\/+1kn1C\/TdnC\/rSbbUrb\/AES70u8GP8j1+nXt1bO9vzM+Ovi3W9Nurj4Y\/tM3mvfCK80\/XtQ1L4M\/tLfDca9\/Y+kegPX9DzjHy5Gz5LD5P\/FeYZfh27\/7UsRSeuu979dHa2m2tny\/reG8Y+PMNiaOIwvE2ZW3\/i6rbprflv3T7uWh+Wdn\/wAEnP2fvjB4Y8S6l+zX+0prHjO\/0OY2eo6D4w0\/Oo6f6c2v2HHHf+z291wPn8k\/ZC\/4Ie\/tR+CfiR4D+J2sftCXnwzstPvNP8bf2b4bP6dF29P75xnGV27n+x\/ij4R8Sf8ACW+F\/Ft548svhH8dPsf\/ABRP7XXwT\/4nHwP+N\/T\/AIkPxaGg7v8ALYOOHX6h8DftkzeFf7P8B\/tWaV\/wrv4h2cNvpt54v0C3urjwf8QLjH+k6ppZ\/wCP+z7n\/kHccA5BOmrH1bLMpxCx2AyrC7+9LDV9F5\/Ztv31S6WtL73\/AIj54iZvllTKMdmvt1iVo69GnWtr3tFf57WW8v098Sabc+JPDfijSvEVnn\/Q9Q7jp+S45J\/v5\/iJ3B3\/AIZ\/i7p2vah4\/wDEkPjq816XxZoesajps3\/CSf2oNY\/4l\/f\/AImhN7jODt3enJwTX9melala6t9m1XR\/En2+zkP+hmG\/\/tDT\/wCz\/fnHr8zZxnI3YCtxPjb4S\/CXx\/qU+rfEL4T\/AA38eXhs7D\/kJeF8Z+6evGD2+63qNuQG+A46w2C4\/wANg7Yv+z8RgNfrXw27X+PV7duuuqPR8DPEleEONzD2+Xf2jQzDVX087\/a9Fdff9n+JY\/BvTftN59m0zjP4fTsePfPX3YV7H4b8bftD\/C7Uvtnw9+NnxI0IdP8AkaAOemOjZ6Z3ZPHHzYLL\/XFon7Kn7IWm3Ju\/+FH2fh\/103Tvx4J5\/lxnBzmq3iT9j39jbxzpVzp\/if4T\/wBl+Zeeb\/b2gf8AEu6gd7ZbLrn1GOuT\/D8Jh+BM\/wAM\/wDZ+K8P6Ln0fzv5t6WW2uiP6qX0tuA8dfD5twziq9Db+DSrdOuru1vq38vs\/lD+wr\/wUm03xPrfgv4I\/H7xf42k8e6xrNxaWfjyX+y\/7P1D+0PsP2XS\/wDRmsfsX8R+630OcV\/Q3c+ANT+zXlzacfp+Pcc5Hvg54JzX4xfED\/gjn+zL4kubPVfCXxk+JHgP+zj\/AGlpp5\/T65xn5R3wcYb90\/DXie20Tw34e0G88SWeuXel6DYabeal\/wBg\/QfcEE49cdhzgtX7Fw7k1DF4K3E1bD+3uvquLw9bp+rWu+n3n8deOuc8GYrMsJm3h28Th6GI9v8AW8LiKXsfYebtfX7rbXeqPJf+Ea1v\/SP9HvCef8+p\/Pv3qzc6Hc2\/2P8A4mlnqHrpuf5fKR79R6YPL17Zba3bWv8Ay8Wen2f2z\/PZefouB7Yw2n4k03TdS8N6rqv2P\/TNPs9Q1Lpnr6Dt+O71+TJDenX8PsPisPVWAzG1fm7W1XzabXTvpe11I\/AHxVidPrNFNeeyXm9726uye2t\/d\/Hz46+JPtPxQ1zEw1A6fZ6fpv8Ad24\/Fs5+h\/3VyNvm9trdzj\/vkev142AD7pH3nHzZ+Xbvfm\/i1rlt\/wAJ9qlpa3P9of2eP+Jlz+fpjoec47cdVzNDuNSuf+XPPbjPf8Sc\/wDfX0OQV++yTEv6tRtt9VT0vayvpu7rpt92iOP6t5fie2abc9e3P\/1vbIz6Z6n0CP22pQ203hLxG+qXH+hx6TqOpDj\/AJ8f9K7vx3\/iB9xkBc3wT4S1G6+z4t+o2+mPf26eq56A8itv466lonw3+G+qWl5\/yHPEFp\/Zv9m7t30wdowc9T8vQgZ6t9lhsKsTh+vTTtbfve6v1jfs7M8TE4lf0r2v21v+W+zsz8zdE\/bz\/Yw\/6OT+Ff4aj\/nvx+navSbj9vz9h65+x2d5+0x8LfTr364P3ifTgjp94YBXm\/hv+yX+y7qVtb\/8WT+G\/wDtY8MZ59jhD1Kjpzgk8ALX1pon7F\/7Klt\/x6fBP4b\/AFHhf26dePm+vIHXO6uHD+H2WYn\/AGlYzEu7d2+nZ3720uk\/lqpelX45zXDL\/cnd73j\/AMNbfXey83Y\/LHW\/22P2ZtO8SXH9lfGTwr4g0P8A6hoxjtn+POPbyx2wcAJpX\/7av7Met2ESaZe+NtYvJIp\/9Ds\/A3jI45zxdXNjY2ZOP4QePQ5BX9ltE\/Zv+Beif8gv4T+CdP5O7Phj5hk+mxf5tjqd3DL6jp3w38E6b\/x5+ENBtMH\/AKA+l\/h\/y5A\/+heudpyvp4fw5ylb4vr2v+i13s7+qVkceYeIuLx2H5vquH06X6\/jZW6cvre1j+S+503UvG2t3t34T+B\/xy8Qfb7z\/oR9v\/sp9evO7gbVwdnrXh79nL9o3XpreHRf2VvG0byQ+d\/xV+oaXoHP0uLqyvgffnrnyxgq\/wDVHbaJbW3\/AB52\/uecD6k+XuB5xypOfuqcgVt\/Zvb9a+6w3DeA7tL\/AIHa9+677u2yPz3+0sR5f+Tf\/Kz+Yj\/hif8AbG\/6N08Lf+HQaj\/hif8AbG\/6N08Lf+HQav6c6KX9hYb\/AJ9x\/wDBZh9dxP8APL7j\/9L+oD4b+N\/tFzefbPsX2Tntn+o7juM\/xDO1hXtet65dWt1pV34d4vPD95k6l\/jnH9c\/8BzX55fs0+PPhX8QfB+j\/FHwf4wvLzQ9fhg86HOP7Px\/y6\/x\/TI449SBX2x4b1vw3qNtefZP+XD\/AJBo4x+YYZ47fy5FfO0UsThrXeiV\/wCrS0dv73V3Wx9ficMsLibYnR7X+Lt5dr\/zert73vGh+JvCfxbt73VvCs4\/4TPT8ab4u8Ody3fI+XgYGDvJOdpUbcNf\/wBJtrm2\/wCYfeD7ucHP\/Acxj19j15JBT8e\/jr\/wmvw38bXnxC+HvjC80PxR0\/tLnjPfqwOP8D82QF9s+CH\/AAVH+HvjF7PwV+0poll4K1iQ\/wDI1EH+z\/w5vbznpwGz0Ozo0fWcL5f+SjWW4hpfVrPXp0378t1rv999o\/pn\/af+x+v\/ANrrSGp\/Zv8AP5d29cjO3OMfLkMnL6L\/AGL4k0SDxF8PfGGj+PNDv\/XWOoxnPUnP4jJ4KtkPVz\/iZacf9Lt\/sB+76Z69OSOMfU7uMZ2Jqc51ttc\/8ed378\/\/AFuOPwXtnsoXE8f+EvBPxI8N6p4S8baHZ+INH1H8PbtnPuOP+BZZ6zbb7SLa4HFh3wfX9c4z6DnnnG1tL7Tc+v6UAfg18Uf2FPij+zf4t8UfEL9m+4\/4SD4Z+LL0al43+FfiT\/kT\/F\/X\/kPaF1Pifj\/kbQF9ecMrfMtzpngr4k6bqfhL4e\/bNPs9P\/tHTfEn7OvxIxo+seF9YH\/IePgPX9eCZx9Dn2w1f0yalpup3X2j1+vBx68DZ6f8tOnJOVNfnR+0b+xV4K+J0P8AbD6Hd2HiC5OLLxh8PfstvrGn+3+lWV5ZXns2oWDD5uq\/wfmfEPBmHxT+s4elp1tdt29V1X+H5WvH9AyXiuvf2GIraczTbvrbZd+vWy6dby\/Lvwj4\/wDGvwS1K30v4IePP+Eg0fT7PUNS1L4D\/EjTzpHjH0\/4kPRe\/wDdQj\/aySv2x8G\/2\/8A4afEhItH8QwaP8O9Xj\/s\/wA+z8Yaha\/2f9dL166eyGcH7t9\/Z2fUfIqfnR8ffgV+0j4R\/wCJVpdvoPx4\/wCP86a2p\/8AEo8Yf9wHQNex4079QM\/7Az8nwta6J42udSuP+FyfDfX\/AAHoen2eNR1Hxtp2u+ENY0f669oT579zjjOGwyV+ZZjws8S39Ywn1Zfzdn663tp9n7\/il+oYfEYfFYRYmhX+sW37O2\/RNb72dusU2z+vzRLnXNb02fVbPwX\/AMJBoef+JbqXgnxToPjD2zjK5447464bKhT\/AISTRPtQtLjw\/wCKtC+x7f8AkJeF9f7d9wAPJz356DGWKfzbab8UfBP7PHhy88bf8J54p8J6Hp9nn\/iifHGvaP8A2vjt\/EvPp97uc4KrneAP2x\/G37TPiP7HbfHDxXp\/hf8A4l39nf2l8aPHejf+D7QNBDeNO3Z\/f5uK87C8GYfEWS+tYdt6yxPs+ndXvpvdW2tZ2vLxVmP1eySwuJW9qFo2frZrZdtfI\/pj03xb4S1sD7Hrmg6ef+on8v8AQj9Mnvjmj+2\/CVzc2\/2PXPBPbgeKG+vuPUdePfivyO8JfsT+CfH\/APpeq\/taaDp93qHIOm\/B\/QfGGPb+3viyPFHjT26se20ZUt7Jo\/8AwTY+GlskiRftOa9rAk6jUvg\/8LtXx\/263V5d2XYf8w8ZPTrmvew3hl9ZoL\/hR9vp187dPe1v+d1fRHiV+KksSv3OJw1tb7v0a1t2\/VWufpBcnwlbWw+03nU5HP4Y+9155OAPZtqbfAbn9pOz\/t\/xB4Y+FtxqXiO\/0+G40G8\/svR9U1H9\/wB\/+PYvn5c9lwT6hGr54vP+CfV5L9sXwx8WPCviC3uf3Pk6x8OPGWkf+V7wb400O94+rY6ZOSV5r\/hm\/wDaZ+Dlt\/xT3jzxv4U0M\/2hn\/hCPjB\/xJ\/r\/YHjvwZ7\/L7dduMt4WJ8HM0xOMwmIoZgqDoYla4ah+\/9Fqk\/l6Xd2zfDcaZX9WtiKP1myX+8apJd\/wCWze1v0cq3hr9mb42eLdSuNX8Q6GdPGoXn\/MSwf6Zzx78ccg4X7Y+G\/wCyX\/Yps\/7UOPyP6cZ6dwnrzg1+b\/8Aw1p8f9Ntbz\/hE\/2gPFVhe6HZg6lpvx++D+NH+mvfEfQf+EpOMHOQPTgZG3mvFv7dP7cGtW174f8A+EY8L+PNDtLMf2nqPwAOheMNX1fAxkDQf+EX8Yrx\/dVcnBOTgp+95bk+V4J2xNF\/k7Prpvo768r7S\/l+Er4rGYr\/AHeu0n27rZ6Lr5vza1ufqx8Wvi18K\/gVo8iJcaPrHiPyf9T\/ANA8+2QDz9WwBzyw2\/id8UPjHqPxP8SXGq3lxz264HHp0I7gHGcdTkq3wZ8Uf2h9S1LxLe\/8LNnvPCd4b3UDp3hv4kWGvfD3SM6P6bvM9s89+QuF22fBvxk8AeNhbWdt4w0HSO32uy8UaXwBx\/db1Bzlc+33qrMa+IxKUsPT9hh473ur\/p3aTXfdNFZbh1hrJ1XicQnZ7v56vW\/\/AG7pvsmfqJ8G9Suev2z+z+P6dejD\/gOBn1GAyffnhu5P2Y555\/IDtjLY688+mM87fgz4G3ngDZbv\/wALD8K3iDjyf+Ep0G4\/AHacccDA+gGNlfoL4budE765oWf+wxoPsfbn09ep2YCP7mT\/AMFnzeZf7yvR\/kd\/pn8X+fWtz+zP9muJ\/wCE28E\/2l9m\/wCE98E+v\/Iz6D\/8nY\/HpXa211baj9iOj\/8AEwP2PjUegIPX+FunuUzyA\/8Ac93C7P0\/VnhDK2Kbbabz\/pdxef0x9CcHuTgjGMjFbVemeeVvsqf3jR9lT+8a0aK6AP\/T\/Ef9hH\/goR8VP2RPFUj6beHWfBdzN\/xPvhjef6Rp2ofaOt11zY8\/L\/F6DH31\/rz\/AGTv29\/hp8ddAs\/EXw7uLy5vJMm88Hy\/8hDT+eynYOeh+Zc7ekfSv87XTc3FznPH8P8Ann3z9\/Pbb9yvuH9i\/wCM3jn4SfHj4eeIfAkf9oapN4q0\/Tf7Ns\/+X\/7f\/ootegH+mn2O7vjA2\/i+TZljcninr9Xfl17v4bteke+u0f8ARHjzw8yPizDVMww\/+z5wrPV3+sfJ30Xkt+juub\/Qd8babrfxItv+Pc\/25YXn5Y+hOMcd\/X5VYmvgz4kfC7w3c\/2naeIZ83n+nj+zScddvXoPXnhfXbwa+4f+GeP2vvij\/Zf9q\/GT4b\/AfwX\/AMSf+0fDfgnRv+Fg\/EjWP+w\/r2vf8IweeB0HHHODXrXg3\/gmn8EdPv7fXvH+ueNvjBrsd55upzfEDxDjT7j2\/sHwxYeFyPy78kkHd+gYnJK+Jtfmf92NrO+6+1t5rS2705f4m\/trA4f17adV1tZL8fRbn4IaV48+IvwKvLxPhd8YPFfhu484Tf8ACH+D\/wC1PEP2j7P\/AMfP\/FB6W2s2V7en\/sHanqQ6cYr9KfgD\/wAFTf20D9n0r4nfsz3fizwv\/wBDvqX\/ABb7WPr0x\/6CT0IbI2fqxpXwB8AeDtN\/srwT4H8N+GLPyOToGj6Xo\/Tof9Fs7ID34OM5ILBUbzXxJ8E7W5H\/AB58dFB5x7Y+Xr7t3zlSCGrDZPiMsWtbl66N3S7atJ3\/AKtvEr53gcc7+w5b9lF7b2230\/m+VkzT8E\/8FBvgD4lNv\/wmvh\/xV4DvP+olo\/54\/dqDn\/eGc8qMAP794b+Nn7PHiP8A5F74h6DnuP7Q7Yxgj644PfHIxhvzw8SfAG2xcDyMH8Tke3ycH1GWU\/ga8K8SfAG2xc\/8Svnr\/nkfp+Y4rP6xiMMv4CTvfXX7vh+7t2s3LRZfgMQ\/3FZbbbfj7t3p2620s+b9y9N1Lw3qf2f+yvGdnnjr\/nn1+8uOvYhdDFp\/0H7L82r+bDUvgzqWj\/af7KvdY0uP\/qG32qW\/44tmb5uuVDLxy2Mqi81cyfHXQftH9m\/Fj4m2nqf+E58U\/r\/pR3cf7SfgRlsP7Zl\/0AIP7K\/6j\/8AyT\/7of0ua1pXh7xLFcWfiGy8K+ILP\/njrFvpWoad\/wCA2qC9Jz7n6Fclq8tufgp8AftN5eXnw2+FeoZOf+RW\/wDrd+2AMe2Wr+b\/AFLx\/wDtRW32j7H8ePin3\/5mjH\/sxx+Y9MHO+vJdc+KX7Xvb9o34w\/8Agw\/r8v8ALnknBJ3R\/b2A\/wCgf8P\/ALoZf2Jif+gh\/ef07a38B\/2W9V1T+2dV+EXwTvNZz\/x+Xvwv8G3H5\/6EueRn7vvkZxRpnhv9mbwT9o\/4R74d+CdP\/wCxb+H4\/PHzD+fTJ5OK\/lB1Lx\/+1pqP\/H5+0D8cvw8cZ49OcD9W9yCF3+O+JtY+OWq\/6NqvxQ+J2sdsXnjjXrjsecd8Y\/plcKHX9t4P\/oHf3HP\/AGJif+gh\/ef2Q3\/x4+Gfhv7kFpbSY+b\/AEjS4Dj1z8wHPcD6qgARvlzx3\/wU1+BHhj7ak3jf4V2vl5\/1XiG21D8D9lW+PHPAdffG5HT+S+4+Empa3\/yFBd656fw\/hjD549QmP7wzVm2+CffyPzBz+eP5r\/dJ24+bm\/1j\/wCnETp\/sfBf8\/pf+DI\/\/Kz9wvHn\/BarwpvuLXRLzxVrv\/Tbwf4fOn6fj3Oqbb0Y44XT2yeu3I3\/AJ9\/E\/8A4KT\/ABl+IK3MfhXw7puh2f8Az+a9qWq6xqHQ9eLKzJH\/AF4L9Blg3zdpvwcuf+fPOR+Xbk+vH91eOwxhfQNE+DmpH\/lz9ueMZ4+n59xn5cFWw\/tPGYro5PTfS2tr9W+uy\/8AArcsksNhMM1\/nv8A+Sq343XRbHkmr+Lfid8QfN\/4TPxhrGsWck\/nfYv9Rp\/f\/mA6W1pZD8l49du+X0DwT\/wkmiCz\/srxBeaeeOPf3P8AnrnsWb6G8N\/s8a3c\/wDLp\/X6\/wBwHp6\/XG4lffvDf7PH2X\/j7tvT36ewOTnucj2JI2KYXLcTiut1fdrX5fDu\/X52tLo+sx8jifBPxg+Kmsafb+GNe1Cy8YaR9jMP9geL7c+L\/wDR\/wDsA69Za8Rzxzv79erWbn9hXwT8WtO\/5Er4b+HrwWeocf8ACD7v+Qx\/xPv+JCNA\/wCEXBH4qR23Y3V9oeEvgna232f\/AJhw9Sf8S3\/fJQY\/H5fr3wn4AtrYWf8ApGcHp0Hr938Pb6ndsb2styZ4VLfdpXfVfd03+LvZ7Hn4rO1itPRb+XTRdHv87LQ\/Ji2\/4JU22m6b9s0vw\/eeH\/XUfgF8Udf0fWDjj\/iffDn4sf8ACT+vYc9mXJZfknxl+xz8ePALfZfDHxY0HxJql59o+2+FvipYWvgbUNYP\/UB\/tMDRLzof+P8A1DTMd9wIFf1N+G9Etbbp\/wATHPv6dh1Jzz\/AnvuxXbXOiabqOm3FvqtnZ+ILPr\/ZuonB\/kx75HyORyTgEBurH8O\/XcOnZ379vN+uujsvvvLmy7iNYXE3V3r2vbfTdJbdV8mnY\/in+JHgD4\/eANDuNV8W\/sN2eoaP\/wAxLxH\/AMINoP8AZA9\/7e0JuOfTP0GSq\/LugftIeIfg+LO6+FXizXfBd550GpQ\/8Ir8Sde0D+z9Q0\/\/AKCn2ZbP7bjHYc5x8uFDf3QTfs3\/AAHubLUNNl+HfhuOzu7zzruz00apYQH\/AMBBZDr+n8Jzleb8N\/sc\/sheG\/8AS9K\/Z\/8AhXp956\/8IuR+OfkH4fN1687U8\/LeEvqu7bto+lvlru2+jtbZfb9XE8R4XFWXe9m7X9fW76JdtLJn8+n\/AAT0\/wCCsX7TvjPxdoXgDxh4G8SfHDw5qmvafoOpeMLPw\/8A6R4XOodbq6urWy+xfYjjH+nsuO55Ar+pL7Kf7x\/8dqtomleHvDWm29n4e0vR9DsyebPTLDTLfA92tFIO7qPucjkDLlND7XZ+q\/l\/9sr7TC4b6stK+ytu0nf8NPTV66WPh6+JWJs\/Ya\/4draaPTt\/nfRDvs3t+tH2b2\/Wj7T7\/pR9p9\/0rU5D\/9T8xPil\/wAEgdLsP9J8H65eWdxxzZah9p7f7uPw\/VdpVtL9iH9lrw3+z98XdP8AiR8c5P7YvPAesf2l4P0H7BawaB\/xL\/8ASvtWqfatvc+vOO2cN+olz+0z4J8R23+kn+z\/AFx26fiOp\/ifsBjBd+JuPil4Q8S6J\/xNdUtNQu\/+Jht\/4SXTz7nA+6Sfy9toJWv81Mu8Yszymv8AWcNiv7Qw\/wAX1XE67vTrv07rpzWtH\/fPOfBLDcQZZUwOPy\/+z8S3pisP+6aXVdPXVrtre5\/Tb+z98Xbz4o6Db+IobjR4o\/J\/1Nnf6XqH+HH+ecYr7e0P\/D+lfw2aJ431Lwjqdn4g+FHii88B659s\/wCZI1DH0\/hzu6ng89yuA7frR8Af+CtfxL8AQW+j\/HLw3\/wtTw\/1vPGGgfZdP8Q2P\/AeLG9wDjmRvX+01ya\/pDgz6UHD+bYh4bNcJ\/Z9dd17XDva7Vrrrayfm2\/tfwX4mfQV4w4ew9XG8I4n+2MPp\/s2tGu\/LZ3u7fzeTVz+mS36\/j\/8TQdE0256fl0\/Hqff+HHqRkmvmb4BftY\/BD9oCwjm+F3xE0jWL97PzrzQZh\/Z+v6ePT7LdHJ68cr6EjAZfqO16n6\/0Nf0\/hsywePw3t8PUw1ajbXqrbaqy8uvre1z+Fs6ybMuH8ZVwWcYXEZfXoWthcRu\/la7b0t8NvvRyOpeCba5x6Htx\/8AqOec8\/KfTBrzfUvhfbXRuP8AR+DyG\/yU478FenU5LP7\/AP8AH16Yz9Tz\/wB9rtI\/67Z3dmHz3Kj6tH+V\/wBf9wzl9tX\/AJ4\/fL\/5WfEet\/BS29P84z\/n357V5drfwBtrr7R\/of65x1Pov659Pmxvr9Hao3Om6af\/AK4\/+t6H8M45wGXCvgMMtHa1tH117qy1Xrpf7V7S6f7SxHl\/5N\/8rPyX1v8AZvubr\/lz9P8AZz+vHXuv4LnK+b3P7M32kn\/Q\/ocdM+vOM\/VfqV27G\/ZW58N6bz7f5\/ofy9lrDuPCOm\/5Oc\/oB6dz69Aa8v8AsTCfyf1\/4MOv+26v8\/8A6V\/8rPxG1v8AZlufs1x\/o\/Tp\/D+fpx375BG7BSuA1L9l65\/59\/xH\/wCyBn3yc9uua\/d258E6bi39\/wCnXuSd3PAKemQG2Vnf8K\/037Tcc9s\/h6\/wkZJHOG6ZKHNctfhyL3S0Vvu3ey3em7NcLna1b0e\/b8bO1nbv5WtFn4Z6J+yLc3lh\/pNvyf7Q\/H+z\/Tj1B\/vAdcj7i9tpv7GGM\/6P9T19+PqDnovX+6Aq\/s9pvgnTf3\/+jcDr2H4\/d\/2ehPXoc5YufCWm\/p69f6DqeePTnI29OGyXC4Vde6i+lumzXbo35OzRGIzyve9nb1s1a3Tk81\/N5bpn5Qab+yFpvH+hn+X5ZGR+X4Z4r1HRP2XdOtv+XcYHoevr7\/Tp75x8v6QW3hq2\/wCfP8ePf6dfXHX0yKP+EcHt\/wCPV6n1bDYaS8lbV9tPlv15vne55Ptq\/wDPH75f\/Kz5C8N\/AHRNOP8Ax7\/z\/nz\/AOhqR0JbG6vSLb4OaJyM+3p+mBkdO\/bgrlDX0J\/YVSf2Z\/s16H1aXmcn1mP8z\/r\/ALiHhdt8JdN\/5cx17D5uOvX5Omeh289GGCr1f+EJubb7R6j0\/wAOMemM\/d4+XjZ799m9v1o1K26XOOv+eox65OMY64bop9Wl5iPnf+zdStsY4+vPvzjbnB+nv1YVX\/tLUrbocdOo654HUnHHPbHGd2EavdPs3+x+v\/2uq9zpum9vxH6Y6n36AcnPY1zAeAXPiT7N3I49M+nB6A4+mfTGa4H\/AIWQff8A8dr33UvCVtrf+iWdn7\/2j6d+Wyc5J4yE9Ai5ArzfUvhLpptun3enPX\/x3jA9x\/vHJKgHAW3xR+09\/wCmPwy+MkY+\/wDjx8\/bW3i37Vj9QPz+n057clQQ78TqXwk\/+vnn\/wBmXv8A7J9921iy6b4I1LTeM\/0Gfwz39Q3XgxZO7oA9H\/4STUvWy\/Oj\/hJNS9bL865T+xdRo\/sXUa6DnP\/V\/Ef4bf8AL3+H9a988Of8vn4V4H8Nv+Xv8P61754c\/wCXz8K\/ya4r\/wB4r+jP+sfFf7yvX9Edd4S\/4+7yvW9S\/wCPMf8AXka8k8Jf8fd5Xrepf8eY\/wCvI1+TZp\/yNqX\/AG\/\/AOlHiYn+v\/JT7g\/4I+f8nxeF\/wDsnnij+Vf2dWv\/AB8n6\/0Nfxi\/8EfP+T4vC\/8A2TzxR\/Kv7OrX\/j5P1\/oa\/wBS\/AX\/AJJCX\/YTH8kf4V\/Ts\/5OnS\/7B6H\/AKSzVtv8\/wDj1XKp23+f\/HquV+9n8ThWBcdfx\/8Aiq36wLjr+P8A8VXOAlYVv1vfr\/St2sK363v1\/pQBDa\/8fh+n+NW6qWv\/AB+H6f41brmw+3y\/+ROgbbdbz\/rz07\/0Km3P+f8Ax2nW3W8\/689O\/wDQqbc\/5\/8AHa6TnLH\/AC7f59KkqP8A5dv8+lSV0HOV6r1YqvXXhd36\/owCm3v\/ACDrf\/r81n+dOpt7\/wAg63\/6\/NZ\/nUAczpf\/AB86p\/nuK5zW\/wDj00r\/ALDo\/mK6PS\/+PnVP89xXOa3\/AMemlf8AYdH8xXOB1\/8Ay7f59KwtQ6z\/AFP8q3f+Xb\/PpWFqHWf6n+VdAHD3XS8+n+NYFb910vPp\/jWBXQBXooooA\/\/Z"],"webLink":"https:\/\/zhangcatherine.com\/japanese-strawberry-shortcake\/","sourceName":"zhangcatherine.com","folderIDs":[],"sourceImage":"\/9j\/4AAQSkZJRgABAQAASABIAAD\/4QCMRXhpZgAATU0AKgAAAAgABQESAAMAAAABAAEAAAEaAAUAAAABAAAASgEbAAUAAAABAAAAUgEoAAMAAAABAAIAAIdpAAQAAAABAAAAWgAAAAAAAABIAAAAAQAAAEgAAAABAAOgAQADAAAAAQABAACgAgAEAAAAAQAAACCgAwAEAAAAAQAAACAAAAAA\/8AAEQgAIAAgAwEiAAIRAQMRAf\/EAB8AAAEFAQEBAQEBAAAAAAAAAAABAgMEBQYHCAkKC\/\/EALUQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29\/j5+v\/EAB8BAAMBAQEBAQEBAQEAAAAAAAABAgMEBQYHCAkKC\/\/EALURAAIBAgQEAwQHBQQEAAECdwABAgMRBAUhMQYSQVEHYXETIjKBCBRCkaGxwQkjM1LwFWJy0QoWJDThJfEXGBkaJicoKSo1Njc4OTpDREVGR0hJSlNUVVZXWFlaY2RlZmdoaWpzdHV2d3h5eoKDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uLj5OXm5+jp6vLz9PX29\/j5+v\/bAEMAAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAf\/bAEMBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAf\/dAAQAAv\/aAAwDAQACEQMRAD8A\/v4rhPiF488O\/DDwR4r+Ini+7+w+HPBmh6n4h1e7HLfYNNtvtLQWqlx9pu7o7bSytSc3188caBXZErwr4\/ftQaN8Edd0DwTB4en8SeM\/FOmXGsacNQ1nTPCvhHSdLtrr7GbvX\/EGoG8vcm7O2zsNC0DXr9nG\/UBpOnE6jXw38Y\/HnhP9oH4e6p4Y8dfEfxf4sv8AXfs95\/wiHgj\/AIRjwh8OrG30jXbG6u9MOmaV46vPGPiHWB9h22N74p1rxHp9lfiPU9M8O6VqIKVnWqONJezv00\/Pvt5t\/Kx3UY0qmJpfWH+5f+9ew008207O3W2tulj4Hi\/bf\/be\/aG\/bS+FA+Dep\/8ACMfDvXPG3hdofgpe+KbWw0rUNO8L2us6jr1v4o1\/7JZt9l8R+GbO+vNbtLFtQsft39mDTNN1LUv7NVP6jK\/nb+EvgbwV+zZ4y8N\/Fbw\/4c1+78U+CbO31iafxTYGeDxBpFvoXiS1+I1rpWq\/2xZeG\/DWsax4avb6+0KzsdF1DUP7f8Padpp8RHTNS1LTdS\/ofjcOqOp+Vxkep+XofvdCPVfqcHd5GQRzRYK2bYrDV8a9V7CzodfTTzvfp7vxH1vHtbg+vnlSfBOFzDAZG8PR5v7Qf7\/2z0q99l3slt\/eP\/\/Q\/qc\/b5uviDpvxJ+Hl18JH1IfEQfDrxkbODTNN1bURqGnt4x8BwXRu20zw54h+yNalt1ob660Gw35U6qxDWD+H+D\/APhqK6+Enjm+8Ua38Wo\/je+sfYvhxoVno\/jK38EHT\/8AQvtF14ouv+EE0e9ssD7f9hvbHXH0\/nTP7RZtrtX6u\/E74GeFPipq+ieIdW1jxr4c1\/w9p2raNpmseCPFGpeFtQOk6xc6Xd32nXjWZcXdkb3SLC6AKj50wWZShrlIv2W\/BC20Frf+NfjfqkdvEIlln+NPxFtJzzzk6FrGld+cjr6Dhl8PH5ViMRifbQeH6N\/WF8u6d+mvJt1vaX0uDzWnhMGsPPC3X1hYh4nz\/wCfXxK1PX+T7tD8r\/h3J8X5vBvxL039r2w8dy\/aNN1jR\/hlNp1x4nt9IGvnwt4l\/wCRoubrVrPR72z3XdgNDzY3\/wDxMMksMOr\/ALleG2J8N+Hifnc6JpBJ6ZJs7XJ6N3+nTtkMvz9f\/sgfAzV23a9pnjzX0\/cSeRr3xm+MurWSmH7pFrc+O3s+OpzG2fvYVQd30rb28NpDDbW0aQW1vFDFDDGNsMEMIGAo6Y24UH1GSzEfL3ZfhKuEwroVF69dX923bme+r0TPPzPHPG46piuT6trbqt\/Ozevpbq7XR\/\/Z","name":"Japanese Strawberry Shortcake","duration":45,"isPublicRecipe":false,"neutritionalInfo":"Cholesterol: 86.1 mg,\nCarbohydrates: 39.9 g,\nTrans Fat: 0 g,\nServing Size: 1 Slice,\nFiber: 0.9 g,\nSodium: 33.7 mg,\nSugar: 28.9 g,\nCalories: 258 calories,\nProtein: 4 g,\nSaturated Fat: 6.1 g,\nFat: 9.6 g","serves":1,"uuid":"FF2DF2F5-F36A-4FBA-B210-2E0D1896D953"}`
	buf := bytes.NewBuffer([]byte(data))

	got := models.NewRecipeFromCrouton(buf, func(rc io.ReadCloser) (uuid.UUID, error) {
		return uuid.Nil, nil
	})

	want := models.Recipe{
		Category:    "Dessert",
		Description: "Imported from Crouton",
		Ingredients: []string{
			"large eggs", "Whole milk (1/4 cup)", "Vegetable oil (3 tbsp)",
			"Cornstarch (1/3 cup 2 tbsp)", "All-purpose flour (1/3 cup 2 tbsp)",
			"Granulated sugar (1/3 cup 2 tbsp)", "Granulated sugar (1/3 cup)",
			"Water (1/3 cup)", "Whipped Cream (2 1/2 cups)", "Granulated sugar (1/2 cup)",
			"Vanilla extract", "Strawberries (13oz)",
		},
		Instructions: []string{
			"Sponge", "Preheat the oven to 140°C fan forced/ 150°C convection",
			"Line the bottom of an 8-inch cake tin with parchment paper",
			"In a medium-sized bowl whisk the egg yolks, mil, and oil",
			"Sift the flour into the egg yolk mixture and mix until combined",
			"In another bowl with an electric whisk, or in the bowl of a stand mixer fitted with a whisk attachment, whip the egg whites with sugar until stiff peaks",
			"Add 1/3 of the meringue into the egg yolk mixture and mix until smooth",
			"Transfer the lightened egg yolk mixture to the remaining meringue and fold carefully until just combined",
			"Transfer the batter to the cake tin",
			"Place the cake tin in a water bath (a tray/tin of boiling water) and bake for 70 minutes",
			"Remove from the oven and allow it to cool completely",
			"Once cooled run a knife around the edge of the cake tin and invert the pan",
			"Wrap in cling wrap and place in the fridge until assembly", "Whipped Cream",
			"Whip the cream with an electric whisk and slowly stream the sugar in",
			"Beat until stiff peaks",
			"Assembly",
			"Combine the sugar and water in a small bowl and microwave for 30 seconds until melted, cool",
			"Slice half the punnet of strawberries",
			"Slice the cooled cake into three layers",
			"Lay one layer of cake down and brush with the sugar syrup",
			"Spread on a layer of cream, a layer of strawberries and then cover with another layer of cream, repeat",
			"Place the last layer of sponge on top and give the cake a thin crumb coat before icing the entire cake with cream",
			"Place star tip into a piping bag and fill it with the remaining cream",
			"Pipe a border around the edge of the cake and decorate with the remaining strawberries",
		},
		Keywords: []string{"Dessert"},
		Name:     "Japanese Strawberry Shortcake",
		Nutrition: models.Nutrition{
			Calories:           "258 calories",
			Cholesterol:        "86.1 mg",
			Fiber:              "0.9 g",
			Protein:            "4 g",
			SaturatedFat:       "6.1 g",
			Sodium:             "33.7 mg",
			Sugars:             "28.9 g",
			TotalCarbohydrates: "39.9 g",
			TotalFat:           "9.6 g",
		},
		Times: models.Times{Prep: 45 * time.Minute, Cook: 1 * time.Hour},
		URL:   "Crouton",
		Yield: 1,
	}
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}

func TestNewRecipesFromRecipeKeeper(t *testing.T) {
	data := "<!DOCTYPE html><html><head><style type=\"text/css\">html {}body {font-family: sans-serif; font-size: 14px; line-height:1.4; color: #333; background-color: #fff}h2, h3 {font-weight: 500; line-height: 1.1; margin-top: 20px; margin-bottom: 10px}h2 {font-size: 24px}h3 {font-size: 14px}.recipe-details h2 {color: #d24400}.recipe-details h3 {color: #d24400}.recipe-ingredients p {margin: 0}.recipe-notes p {margin: 0}.recipe-photo {width: 250px; height: 250px; margin-top: 20px; object-fit: cover}.recipe-photos-div {width: 125px; height: 125px; margin-right: 5px; margin-bottom: 5px; display: inline-block}.recipe-photos {width: 100%; height: 100%; object-fit: contain}</style></head><body><div class=\"recipe-details\"><meta content=\"a2489b23-050f-5c8a-8e5d-12fdb9fdeea3\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/a2489b23-050f-5c8a-8e5d-12fdb9fdeea3_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Blueberry Cheesecake</h2><div>Courses: <span itemprop=\"recipeCourse\">Dessert</span></div><div>Categories: <span>Cake</span><meta content=\"Cake\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">Serves 6-8</span></div><div>Preparation time: <span>10 mins</span><meta content=\"PT10M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>1 hour </span><meta content=\"PT1H\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>1 1/4 cups graham cracker crumbs / 8 digestive biscuit crumbs</p><p>1/4 cup/50g butter, melted </p><p>20 oz./600g cream cheese </p><p>2 tablespoons all-purpose flour </p><p>3/4 cup/175g caster sugar </p><p>2 eggs, plus 1 yolk</p><p>small pot soured cream</p><p>vanilla extract  </p><p>1 cup / 300g fresh or frozen blueberries</p><p>icing sugar </p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. Melt the butter and pour onto the cracker crumbs and mix well. Press into the bottom of a 9-inch spring form pan. Bake at 325&#176;F until the crust is set, about 10-12 minutes. Allow to cool.</p><p></p><p>2. In a large bowl beat the cream cheese with the flour, caster sugar, eggs, soured cream and vanilla extract with an electric mixer until light and fluffy.</p><p></p><p>3. Pour the mixture into the pan and bake for 35-40 minutes until set. Remove from the oven and leave to cool.</p><p></p><p>4. Heat half the blueberries in a pan with 2 tablespoons icing sugar and stir gently until juicy. Squash the blueberries with a fork then continue to cook for a few minutes. Add the remaining blueberries, remove from the heat and allow to cool. </p><p></p><p>5. Pour the blueberries over the cheesecake just before serving.</p><p></p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Serving size: 1 slice<meta content=\"1 slice\" itemprop=\"recipeNutServingSize\"></div><div>Calories: 571<meta content=\"571\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 41.5g<meta content=\"41.5\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 24.7g<meta content=\"24.7\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 173mg<meta content=\"173\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 451mg<meta content=\"451\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 42g<meta content=\"42\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 1.5g<meta content=\"1.5\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 27.1g<meta content=\"27.1\" itemprop=\"recipeNutSugars\"></div><div>Protein: 10.1g<meta content=\"10.1\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/a2489b23-050f-5c8a-8e5d-12fdb9fdeea3_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"686a4630-2212-43f3-8cb3-87af6cd1ec2c\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"3\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/686a4630-2212-43f3-8cb3-87af6cd1ec2c_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Calabrian Chili Orzo with Vegetables</h2><div>Courses: <span itemprop=\"recipeCourse\"></span></div><div>Categories: <span>Chicken, Vegan, Vegetarian</span><meta content=\"Chicken\" itemprop=\"recipeCategory\"><meta content=\"Vegan\" itemprop=\"recipeCategory\"><meta content=\"Vegetarian\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\"><a href=\"https://www.allrecipes.com/calabrian-chili-orzo-with-vegetables-recipe-8623914\">https://www.allrecipes.com/calabrian-chili-orzo-with-vegetables-recipe-8623914</a></span></div><div>Serving size: <span itemprop=\"recipeYield\">4</span></div><div>Preparation time: <span>10 mins</span><meta content=\"PT10M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>15 mins</span><meta content=\"PT15M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>6 ounces orzo</p><p>2 tablespoons creme fraiche</p><p>1 teaspoon Calabrian chili paste, or more to taste</p><p>1 teaspoon olive oil</p><p>1 zucchini, or 2 if small, cut lengthwise and sliced into half moons</p><p>1/2 cup sliced sweet peppers</p><p>2 cloves garlic, minced, or more to taste</p><p>1/2 teaspoon Italian seasoning</p><p>salt and freshly ground black pepper to taste</p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>Bring a pot of lightly salted water to a boil, add orzo, and cook, uncovered, until tender, 8 to 10 minutes. Turn off the heat, drain thoroughly, and return orzo to the pot. Add cr&#232;me fra&#238;che and Calabrian chili paste. Set aside and keep warm.</p><p></p><p>Heat a skillet over medium-high heat. Add olive oil. Once oil is shimmering, add vegetables, and saute for 4 to 5 minutes. Stir in garlic, and saut&#233; until fragrant, about 45 seconds. Add Italian seasoning and season with salt and pepper. Remove vegetables from heat.</p><p></p><p>Gently stir vegetables into orzo. Serve immediately.</p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/686a4630-2212-43f3-8cb3-87af6cd1ec2c_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"b0bcddc4-23e8-50cc-a879-22b1b2e63919\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/b0bcddc4-23e8-50cc-a879-22b1b2e63919_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Chocolate Chip Cookies</h2><div>Courses: <span itemprop=\"recipeCourse\">Snack</span></div><div>Categories: <span>Cookie</span><meta content=\"Cookie\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">12</span></div><div>Preparation time: <span>5 mins</span><meta content=\"PT5M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>12 mins</span><meta content=\"PT12M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>2 1/4 cups all-purpose flour</p><p>1 teaspoon baking soda</p><p>1 teaspoon salt</p><p>1 cup butter</p><p>1 cup caster sugar</p><p>1 cup soft brown sugar</p><p>1 teaspoon vanilla extract</p><p>2 eggs</p><p>2 cups dark chocolate, broken into small pieces</p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. In a large bowl combine the flour, baking soda and salt.</p><p></p><p>2. In a separate bowl, mix the butter, caster sugar, brown sugar and vanilla extract until smooth. </p><p></p><p>3. Add the eggs and the flour to the mixture and beat to combine. </p><p></p><p>4. Add the chocolate pieces and stir.</p><p></p><p>5. Drop well rounded spoonfuls of dough onto a greased cookie sheet. </p><p></p><p>6. Bake at 375F for 8-10 minutes.</p><p></p><p>7. Remove from the oven and place cookies on a wire rack to cool.</p><p></p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Serving size: 1 cookie<meta content=\"1 cookie\" itemprop=\"recipeNutServingSize\"></div><div>Calories: 491<meta content=\"491\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 24.6g<meta content=\"24.6\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 15.8g<meta content=\"15.8\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 74mg<meta content=\"74\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 443mg<meta content=\"443\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 63.2g<meta content=\"63.2\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 1.6g<meta content=\"1.6\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 43g<meta content=\"43\" itemprop=\"recipeNutSugars\"></div><div>Protein: 5.7g<meta content=\"5.7\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/b0bcddc4-23e8-50cc-a879-22b1b2e63919_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"3259bd0c-30db-4beb-a8a8-2aeecc392824\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/3259bd0c-30db-4beb-a8a8-2aeecc392824_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Guyanese Gojas</h2><div>Courses: <span itemprop=\"recipeCourse\"></span></div><div>Categories: <span></span></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\"><a href=\"https://www.simplyrecipes.com/guyanese-gojas-recipe-5221034\">https://www.simplyrecipes.com/guyanese-gojas-recipe-5221034</a></span></div><div>Serving size: <span itemprop=\"recipeYield\">6 servings</span></div><div>Preparation time: <span>1 hour 5 mins</span><meta content=\"PT1H5M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>35 mins</span><meta content=\"PT35M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>For the goja dough:</p><p>2 cups (240g) all-purpose flour</p><p>2 teaspoons (8g) granulated sugar</p><p>1/2 teaspoon (2g) instant yeast</p><p>2 tablespoons diced cold butter, or vegetable shortening</p><p>1 cup cold whole milk</p><p>1 tablespoon (15g) all-purpose flour, to sprinkle while kneading</p><p>1/8 teaspoon neutral cooking oil, such as canola or vegetable for rubbing the dough ball</p><p>For the goja filling:</p><p>2 cups sweetened flaked coconut</p><p>1/2 teaspoon ground cinnamon</p><p>1/2 teaspoon freshly grated nutmeg</p><p>1 tablespoon light brown sugar</p><p>1 tablespoon fresh grated ginger</p><p>2 tablespoons water</p><p>2 tablespoons butter, melted</p><p>2 teaspoons vanilla extract</p><p>2 tablespoons neutral cooking oil such as canola or vegetable, for cooking filling</p><p>1/4 cup raisins</p><p>2 tablespoons water, for cooking filling</p><p>For shaping and frying the gojas:</p><p>1/4 cup water, for sealing the pastry</p><p>1/4 cup (34g) all-purpose flour, for crimping the gojas</p><p>3 cups canola or vegetable oil, for frying</p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>Combine the dry ingredients :</p><p></p><p>In a large bowl, combine the flour, sugar, and yeast. Add the butter. Using your hands or fingertips, rub the butter into flour until a coarse meal forms.</p><p></p><p>Alica Ramkirpal-Senhouse / Simply Recipes</p><p></p><p>Alica Ramkirpal-Senhouse / Simply Recipes</p><p></p><p>Pour milk into bowl:</p><p></p><p>Make a well in the center of the bowl, pour in the milk. Using a rubber spatula, stir until the dough forms. At this point, the dough will be a little sticky, sprinkle 1 tablespoon flour on the dough and knead it into dough with your hands in the bowl until the dough is no longer sticky.</p><p></p><p>Set dough aside to rest:</p><p></p><p>Rub the top of the dough with the oil and cover with a damp paper towel. Set aside for 15 to 20 minutes.</p><p></p><p>Make the filling:</p><p></p><p>In the bowl of your food processor, add the coconut, cinnamon, nutmeg, brown sugar, ginger, water, butter, and vanilla. Pulse on high until coconut becomes fine and pasty.</p><p></p><p>Cook the filling:</p><p></p><p>Heat a heavy-bottomed pan over low heat. Add the oil, coconut filling, raisins, and 2 tablespoons of water. Cook, stirring occasionally, until the sugar melts and the coconut looks more toasted and slightly darker in color, about 5 minutes. Remove from heat and let cool for a few minutes before assembling the gojas.</p><p></p><p>Alica Ramkirpal-Senhouse / Simply Recipes</p><p></p><p>Weigh and divide the dough:</p><p></p><p>Weigh the dough, then divide the weight by 12 to get the weight for each piece. Now, cut 12 small pieces of dough and weigh each. Add or remove small pieces until you get the exact weight you’re looking for.</p><p></p><p>If you’re not using a scale, divide the dough into 12 pieces using a knife or pastry cutter. Try to eyeball it so they’re all the same size.</p><p></p><p>Alica Ramkirpal-Senhouse / Simply Recipes</p><p></p><p>Roll the goja dough:</p><p></p><p>Round off each dough ball between your palms to form a ball, gently tucking dough under itself to make the top smooth. Once you’ve done this, cover all the dough balls with a damp paper towel to keep it from drying out and crusting.</p><p></p><p>Sprinkle flour on the surface of the dough ball you are working with. Working with one dough ball at a time, flatten slightly with your hands, then roll into a circle 1/8 inch thick and about 5 inches in diameter.</p><p></p><p>Flour your surface as needed as you go along.</p><p></p><p>Repeat with remaining balls of dough, being sure to keep them covered as you work.</p><p></p><p>Dip your pointer finger in water and run it around the outer edges of the dough. Place 2 tablespoons filing in the bottom half of the dough and bring the top half over to seal. Using a fork crimp the edges closed being sure to dip the fork in flour to keep from sticking while crimping. Place assembled gojas on a baking sheet lined with parchment paper.</p><p></p><p>Repeat this step for the rest of the batch.</p><p></p><p>Set up a plate or deep serving platter with a few paper towels to place gojas on after they’re done frying.</p><p></p><p>Heat a medium sized deep pot over medium-low heat. Add the oil and once it’s anywhere between 350-375&#176;F, fry the gojas for 2 to 3 minutes, you’ll have to cook these in batches, being sure to not overcrowd the pot. Use a slotted spoon or tongs to flip the gojas once halfway through cooking. Remove from oil once it is light golden brown and drain on paper towels.</p><p></p><p>Repeat with remaining gojas until they are all fried.</p><p></p><p>Enjoy warm.</p><p></p><p>Alica Ramkirpal-Senhouse / Simply Recipes</p></div></td></tr></table><h3>Notes</h3><div class=\"recipe-notes\" itemprop=\"recipeNotes\"><p></p><p></p><p>(per serving)</p><p>543 Calories 30g Fat 62g Carbs 8g Protein</p><p>Nutrition Facts</p><p>Servings: 6</p><p>Amount per serving</p><p>Calories 543</p><p>% Daily Value*</p><p>Total Fat 30g 38%</p><p>Saturated Fat 13g 67%</p><p>Cholesterol 17mg 6%</p><p>Sodium 132mg 6%</p><p>Total Carbohydrate 62g 23%</p><p>Dietary Fiber 5g 16%</p><p>Total Sugars 20g</p><p>Protein 8g</p><p>Vitamin C 0mg 1%</p><p>Calcium 66mg 5%</p><p>Iron 3mg 16%</p><p>Potassium 271mg 6%</p><p>*The % Daily Value (DV) tells you how much a nutrient in a food serving contributes to a daily diet. 2,000 calories a day is used for general nutrition advice.</p></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/3259bd0c-30db-4beb-a8a8-2aeecc392824_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"31adb931-743f-5e9d-b2a2-9269c0775a4a\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/31adb931-743f-5e9d-b2a2-9269c0775a4a_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Lemon Tart</h2><div>Courses: <span itemprop=\"recipeCourse\">Dessert</span></div><div>Categories: <span>Tart</span><meta content=\"Tart\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">8</span></div><div>Preparation time: <span>1 hour </span><meta content=\"PT1H\" itemprop=\"prepTime\"></div><div>Cooking time: <span>30 mins</span><meta content=\"PT30M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p><b>For the tart shell</b></p><p>1 1/2 cups all-purpose flour </p><p>1/2 cup icing sugar </p><p>2/3 cup softened butter </p><p>pinch of salt </p><p>1 egg yolk</p><p></p><p><b>For the lemon curd</b></p><p>6 lemons </p><p>6 large eggs </p><p>1 1/2 cups caster sugar </p><p>1 1/2 cups cream </p><p></p><p><b><i>To serve</i></b></p><p>Icing sugar</p><p>Cream</p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. In a large bowl mix the softened butter and icing sugar to a cream with a wooden spoon then beat in the egg yolk. Add the flour and salt and rub the butter mixture and flour together with your fingers until crumbly.</p><p></p><p>2. Add the egg yolk to the mixture and knead briefly until it forms a firm dough. Wrap in plastic wrap and leave to chill in the fridge for 30 minutes.</p><p></p><p>3. Roll out the pastry until very thin and line the quiche tin allowing a small amount of pastry to overlap the edges. Prick the base of the pastry with a fork and bake at 400&#176;F for 20 minutes.</p><p></p><p>4. Grate the zest from the lemons into a bowl then add the juice from the lemons. Break the eggs into a large bowl and add the caster sugar and mix well. Add the lemon juice, zest and the cream and whisk gently.</p><p></p><p>5. Pour the mixture into the pastry case and bake at 350&#176;F for 25 to 30 minutes until set.</p><p></p><p>6. Sprinkle with a little icing sugar and serve with cream.</p><p></p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Serving size: 1 slice<meta content=\"1 slice\" itemprop=\"recipeNutServingSize\"></div><div>Calories: 480<meta content=\"480\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 22.4g<meta content=\"22.4\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 12.7g<meta content=\"12.7\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 215mg<meta content=\"215\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 197mg<meta content=\"197\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 64.7g<meta content=\"64.7\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 0.6g<meta content=\"0.6\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 46.1g<meta content=\"46.1\" itemprop=\"recipeNutSugars\"></div><div>Protein: 8g<meta content=\"8\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/31adb931-743f-5e9d-b2a2-9269c0775a4a_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"b05310a1-15ec-560e-896d-e7a54223d0a0\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/b05310a1-15ec-560e-896d-e7a54223d0a0_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Pancakes</h2><div>Courses: <span itemprop=\"recipeCourse\">Breakfast</span></div><div>Categories: <span></span></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">8-12 pancakes</span></div><div>Preparation time: <span>5 mins</span><meta content=\"PT5M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>15 mins</span><meta content=\"PT15M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>3 cups/375g all-purpose flour</p><p>3 teaspoons baking powder</p><p>1 tablespoon caster sugar</p><p>1 1/2 cups/375ml milk</p><p>3 eggs</p><p>1/2 teaspoon salt</p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. Sift the flour, baking powder, salt and caster sugar into a bowl. </p><p></p><p>2. Break the eggs into a separate bowl and whisk together with the milk.</p><p></p><p>3. Gradually add the milk and egg mixture to the flour mixture and whisk to a smooth batter.</p><p></p><p>4. Heat a frying pan over a medium heat and melt a small knob of butter. Pour the batter into the pan, using approximately 1/4 cup for each pancake.</p><p></p><p>5. When the top of the pancake begins to bubble, turn and cook the other side until golden brown.</p><p></p><p>6. Serve with butter and maple syrup.</p><p></p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Serving size: 1 pancake<meta content=\"1 pancake\" itemprop=\"recipeNutServingSize\"></div><div>Calories: 180<meta content=\"180\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 2.4g<meta content=\"2.4\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 0.9g<meta content=\"0.9\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 52mg<meta content=\"52\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 154mg<meta content=\"154\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 32.4g<meta content=\"32.4\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 1g<meta content=\"1\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 3.1g<meta content=\"3.1\" itemprop=\"recipeNutSugars\"></div><div>Protein: 6.7g<meta content=\"6.7\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/b05310a1-15ec-560e-896d-e7a54223d0a0_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"7f2decec-2e7f-594b-90ac-5fed6ed1b4c7\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/7f2decec-2e7f-594b-90ac-5fed6ed1b4c7_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Pasta Puttanesca</h2><div>Courses: <span itemprop=\"recipeCourse\">Main Dish</span></div><div>Categories: <span>Pasta</span><meta content=\"Pasta\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">5</span></div><div>Preparation time: <span>5 mins</span><meta content=\"PT5M\" itemprop=\"prepTime\"></div><div>Cooking time: <span>10 mins</span><meta content=\"PT10M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>4 tablespoons olive oil</p><p>1/8 cup butter</p><p>2 cloves garlic</p><p>4 anchovy fillets</p><p>1lb plum tomatoes</p><p>1 tablespoon tomato puree</p><p>1/2 cup black olives</p><p>1/2 tablespoon capers</p><p>1lb Spaghetti or any other pasta</p><p>Parmesan cheese</p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. Cook the pasta as directed on the instructions. </p><p></p><p>2. At the same time, heat the oil and butter in a separate pan and gently fry the garlic and anchovies for 3-4 minutes.</p><p></p><p>3. Add the tomatoes, tomato puree, olives, capers and fry for 7-8 minutes, stirring from time to time.</p><p></p><p>4. Drain the pasta, pour over the sauce and serve with freshly grated Parmesan cheese.</p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Serving size: 1 bowl<meta content=\"1 bowl\" itemprop=\"recipeNutServingSize\"></div><div>Calories: 459<meta content=\"459\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 20.9g<meta content=\"20.9\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 5.9g<meta content=\"5.9\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 85mg<meta content=\"85\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 500mg<meta content=\"500\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 54.9g<meta content=\"54.9\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 1.3g<meta content=\"1.3\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 2.3g<meta content=\"2.3\" itemprop=\"recipeNutSugars\"></div><div>Protein: 14g<meta content=\"14\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/7f2decec-2e7f-594b-90ac-5fed6ed1b4c7_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div><div class=\"recipe-details\"><meta content=\"2f98fbfc-bd5f-4a55-8774-48b29bd2c1ce\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Pesto Pull Apart-Bread</h2><div>Courses: <span itemprop=\"recipeCourse\"></span></div><div>Categories: <span>Cake</span><meta content=\"Cake\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\"><a href=\"https://sallysbakingaddiction.com/pesto-pull-apart-bread/#tasty-recipes-129499\">https://sallysbakingaddiction.com/pesto-pull-apart-bread/#tasty-recipes-129499</a></span></div><div>Serving size: <span itemprop=\"recipeYield\">Yield: 1 loaf</span></div><div>Preparation time: <span>3 hours </span><meta content=\"PT3H\" itemprop=\"prepTime\"></div><div>Cooking time: <span>50 mins</span><meta content=\"PT50M\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>Dough</p><p>2 teaspoons instant or active dry yeast</p><p>1 Tablespoon granulated sugar</p><p>3/4 cup (6 1/4 ounces/ml) whole milk, warmed to about 110&#176;F (43&#176;C)</p><p>3 Tablespoons (1 1/2 ounces) unsalted butter, softened to room temperature</p><p>1 large egg, at room temperature</p><p>2 1/3 cups (10 1/4 ounces) all-purpose flour (spooned &amp; leveled), plus more as needed*</p><p>1 teaspoon salt</p><p>1 teaspoon garlic powder</p><p>1/2 teaspoon dried basil</p><p>Filling</p><p>1/2 cup (4 1/2 ounces) basil pesto (I recommend my homemade pesto)</p><p>1 cup (4 1/2 ounces / 4 ounces) shredded mozzarella cheese</p><p>Topping</p><p>2 Tablespoons (0.99 ounce) unsalted butter, melted</p><p>1/4 teaspoon garlic powder</p><p>2 Tablespoons (0.53 ounce) freshly grated or shredded parmesan cheese</p><p>optional for garnish: extra pesto</p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>Make the dough: Place the yeast and sugar in the bowl of a stand mixer fitted with a dough hook or paddle attachment. Or, if you do not own a stand mixer, a regular large mixing bowl. Whisk in the warm milk, then loosely cover with a clean kitchen towel and allow to sit for 5-10 minutes. The mixture will be frothy after 5-10 minutes.</p><p></p><p>If you do not have a mixer, you can mix the dough together with a wooden spoon or silicone spatula in this step. Add the butter, egg, flour, salt, garlic powder, and dried basil. Beat on low speed for 3 minutes. Dough will be soft.</p><p></p><p>Knead the dough: Keep the dough in the mixer (and switch to the dough hook if using the paddle) and beat for an additional 5 full minutes, or knead by hand on a lightly floured surface for 5 full minutes. (If you’re new to bread-baking, my How to Knead Dough video tutorial can help here.) If the dough becomes too sticky during the kneading process, sprinkle 1 teaspoon of flour at a time on the dough or on the work surface/in the bowl to make a soft, slightly tacky dough. Do not add more flour than you need because you do not want a dry dough. After kneading, the dough should still feel a little soft. Poke it with your finger—if it slowly bounces back, your dough is ready to rise. You can also do a “windowpane test” to see if your dough has been kneaded long enough: tear off a small (roughly golfball-size) piece of dough and gently stretch it out until it’s thin enough for light to pass through it. Hold it up to a window or light. Does light pass through the stretched dough without the dough tearing first? If so, your dough has been kneaded long enough and is ready to rise. If not, keep kneading until it passes the windowpane test.</p><p></p><p>1st Rise: Shape the kneaded dough into a ball. Place the dough in a greased bowl (I use nonstick spray to grease) and cover with plastic wrap or aluminum foil. Place in a slightly warm environment to rise until doubled in size, around 60-90 minutes. (If desired, use my warm oven trick for rising. See my answer to Where Should Dough Rise? in my Baking with Yeast Guide.)</p><p></p><p>As the dough rises, grease a 9&#215;5-inch loaf pan and prepare the pesto.</p><p></p><p>Assemble &amp; fill the bread: Punch down the dough to release the air. Place dough on a lightly floured work surface. Divide it into 12 equal pieces, about 1/4 cup of dough or 1 3/4 ounces each (a little larger than a golf ball). Using lightly floured hands, flatten each into a circle that’s about 4 inches in diameter. The circle doesn’t have to be perfectly round. I do not use a rolling pan to flatten, but you certainly can if you want. Spread 1-2 teaspoons of pesto onto each. Sprinkle each with 1 heaping Tablespoon of mozzarella cheese. Fold circles in half and line in prepared baking pan, round side up. See photos above for a visual.</p><p></p><p>2nd Rise: Cover with plastic wrap or aluminum foil and allow to rise once again in a slightly warm environment until puffy, about 45 minutes. Do not extend this 2nd rise, as the bread could puff up too much and spill over the sides while baking.</p><p></p><p>Adjust the oven rack to the lower third position then preheat oven to 350&#176;F (177&#176;C).</p><p></p><p>Bake until golden brown, about 50 minutes. If you find the top of the loaf is browning too quickly, tent with aluminum foil. Remove from the oven and place the pan on a cooling rack.</p><p></p><p>Make the topping: Mix the melted butter and garlic butter together. Brush on the warm bread and sprinkle with parmesan cheese. If desired, drop a couple spoonfuls of fresh pesto on top (or serve with extra pesto.) Cool for 10 minutes in the pan, then remove from the pan and serve warm.</p><p></p><p>Cover and store leftovers at room temperature for up to 2 days or in the refrigerator for up to 1 week. Since the bread is extra crispy on the exterior, it will become a little hard after day 1. Reheat in a 300&#176;F (149&#176;C) oven for 10-15 minutes until interior is soft again or warm in the microwave.</p></div></td></tr></table><h3>Notes</h3><div class=\"recipe-notes\" itemprop=\"recipeNotes\"><p>Make Ahead Instructions: Freeze baked and cooled bread for up to 3 months. Thaw at room temperature or overnight in the refrigerator and warm in the oven to your liking. The dough can be prepared through step 4, then after it has risen, punch it down to release the air, cover it tightly, then place in the refrigerator for up to 2 days. Continue with step 5. To freeze the dough, prepare it through step 4. After it has risen, punch it down to release the air. Wrap in plastic wrap and place in a freezer-friendly container for up to 3 months. When ready to use, thaw the dough overnight in the refrigerator. Then let the dough sit at room temperature for about 30 minutes before continuing with step 5. (You may need to punch it down again if it has some air bubbles.)</p><p></p><p>Special Tools (affiliate links): Electric Stand Mixer or Large Glass Mixing Bowl with Wooden Spoon / Silicone Spatula | 9&#215;5-inch Loaf Pan | Cooling Rack | Pastry Brush</p><p></p><p>Yeast: You can use instant or active dry yeast. The rise times may be slightly longer using active dry yeast. Reference my Baking with Yeast Guide for answers to common yeast FAQs.</p><p></p><p>Flour: Feel free to use the same amount of bread flour instead of all-purpose flour.</p><p></p><p>Can I substitute the pesto? Instead of pesto, you can use your favorite tomato sauce, or try this rosemary garlic pull apart bread.</p></div><h3>Photos</h3><hr /></div><div class=\"recipe-details\"><meta content=\"d017e2c8-24e6-5eee-8b68-931c0196bf5a\" itemprop=\"recipeId\"><meta content=\"\" itemprop=\"recipeShareId\"><meta content=\"False\" itemprop=\"recipeIsFavourite\"><meta content=\"0\" itemprop=\"recipeRating\"><table><tr><td><img src=\"images/d017e2c8-24e6-5eee-8b68-931c0196bf5a_0.jpg\" class=\"recipe-photo\"/></td><td style=\"vertical-align:top\"><h2 itemprop=\"name\">Tuna Nicoise Salad</h2><div>Courses: <span itemprop=\"recipeCourse\">Main Dish</span></div><div>Categories: <span>Salad</span><meta content=\"Salad\" itemprop=\"recipeCategory\"></div><div>Collections: <span></span></div><div>Source: <span itemprop=\"recipeSource\">Recipe Keeper</span></div><div>Serving size: <span itemprop=\"recipeYield\">4</span></div><div>Preparation time: <span>10 mins</span><meta content=\"PT10M\" itemprop=\"prepTime\"></div><div>Cooking time: <span></span><meta content=\"PT0S\" itemprop=\"cookTime\"></div></td></tr><tr><td style=\"vertical-align:top;width:250px\"><h3>Ingredients</h3><div class=\"recipe-ingredients\" itemprop=\"recipeIngredients\"><p>2 cooked tuna steaks or 2 cans of tuna</p><p>12 small potatoes</p><p>5oz fine French beans</p><p>4 tomatoes</p><p>1 large romaine lettuce</p><p>1 red onion, finely sliced </p><p>4 hard-boiled eggs, peeled and sliced</p><p>20 black olives</p><p>Chopped fresh parsley</p><p>8 tablespoons extra virgin olive oil</p><p>3 tablespoons red wine vinegar</p><p>2 garlic cloves, peeled and finely chopped </p><p>1 teaspoon salt</p><p>1 teaspoon ground black pepper</p><p></p></div></td><td style=\"vertical-align:top\"><h3>Directions</h3><div itemprop=\"recipeDirections\"><p>1. Cook the potatoes until just tender, cool in ice water and peel and quarter.</p><p></p><p>2. Top and tail the beans and boil in water for 5 minutes, cool in ice water.</p><p></p><p>3. Tear the lettuce into small pieces and arrange on a large plate.</p><p></p><p>4. Chop the tomatoes into quarters and add to the plate.</p><p></p><p>5. Cut the tuna into large chunks and add to the the plate.</p><p></p><p>6. Add the potatoes, beans, slice onion, eggs, olives and scatter over the chopped parsley.</p><p></p><p>7. In a small bowl mix the oil, vinegar, garlic, salt and pepper. Pour over the salad and serve.</p><p></p></div></td></tr></table><div class=\"recipe-notes\" itemprop=\"recipeNotes\"></div><h3>Nutrition</h3><div>Amount per serving</div><div>Calories: 978<meta content=\"978\" itemprop=\"recipeNutCalories\"></div><div>Total Fat: 40.4g<meta content=\"40.4\" itemprop=\"recipeNutTotalFat\"></div><div>Saturated Fat: 6.9g<meta content=\"6.9\" itemprop=\"recipeNutSaturatedFat\"></div><div>Cholesterol: 193mg<meta content=\"193\" itemprop=\"recipeNutCholesterol\"></div><div>Sodium: 921mg<meta content=\"921\" itemprop=\"recipeNutSodium\"></div><div>Total Carbohydrate: 118.6g<meta content=\"118.6\" itemprop=\"recipeNutTotalCarbohydrate\"></div><div>Dietary Fiber: 25.4g<meta content=\"25.4\" itemprop=\"recipeNutDietaryFiber\"></div><div>Sugars: 12.6g<meta content=\"12.6\" itemprop=\"recipeNutSugars\"></div><div>Protein: 41.4g<meta content=\"41.4\" itemprop=\"recipeNutProtein\"></div><h3>Photos</h3><div class=\"recipe-photos-div\"><img src=\"images/d017e2c8-24e6-5eee-8b68-931c0196bf5a_0.jpg\" class=\"recipe-photos\" itemprop=photo0 /></div><hr /></div></body></html>"
	buf := bytes.NewBufferString(data)
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

func TestNutrition_Clean(t *testing.T) {
	got := models.Nutrition{
		Calories:           "0 calories",
		Cholesterol:        "0mg",
		Fiber:              "0 g fibre",
		Protein:            "0 g protein",
		SaturatedFat:       "0 g saturated fat",
		Sodium:             "0 g salt",
		Sugars:             "0 g sugar",
		TotalCarbohydrates: "0 g carbohydrate",
		TotalFat:           "0 g fat",
		UnsaturatedFat:     "0g",
	}
	got.Clean()

	want := models.Nutrition{}
	if !got.Equal(want) {
		t.Fatalf("got %+v but want %+v", got, want)
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
