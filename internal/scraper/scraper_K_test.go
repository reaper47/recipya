package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
	"time"
)

func TestScraper_K(t *testing.T) {
	testcases := []testcase{
		{
			name: "kennymcgovern.com",
			in:   "https://kennymcgovern.com/chicken-noodle-soup",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Soup"},
				CookTime:      "PT5M",
				Cuisine:       &models.Cuisine{Value: "Chinese"},
				DatePublished: "2022-03-27T18:12:02+00:00",
				Description:   &models.Description{},
				Keywords: &models.Keywords{
					Values: "noodles, Soup, noodle soup, chicken, chicken noodle soup, chicken soup",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"40g thin lucky boat noodles (soaked and drained, drained weight)",
						"275 ml light chicken stock",
						"Dash light soy sauce",
						"1/4 teaspoon sea salt",
						"1/4 teaspoon MSG",
						"Pinch white pepper",
						"50 grams raw chicken breast, thinly sliced (or 1 small handful cooked shredded chicken breast or thigh meat)",
						"1 spring onion (finely sliced)",
						"Dash sesame oil (optional, see notes)",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Put the noodle nest in a bowl and cover with boiling water. Let stand for 2 minutes, mix briefly and drain the water. Place 40g of prepared noodles in a soup bowl (one nest of lucky boat noodles will do about 4-5 takeaway sized portions of soup). Set aside."},
						{Type: "HowToStep", Text: "Put the light chicken stock, light soy sauce, sea salt, MSG and white pepper in a saucepan, bring to the boil then reduce to a simmer. Add the sliced chicken and simmer for about 3 minutes until the chicken is cooked through."},
						{Type: "HowToStep", Text: "Pour the chicken soup over the prepared noodles in the bowl. Garnish with the sliced spring onion, drizzle with sesame oil (if using) and serve."},
					},
				},
				Name:      "Chicken Noodle Soup",
				PrepTime:  "PT5M",
				TotalTime: "PT10M",
				Yield:     &models.Yield{Value: 1},
				URL:       "https://kennymcgovern.com/chicken-noodle-soup",
			},
		},
		{
			name: "kingarthurbaking.com",
			in:   "https://www.kingarthurbaking.com/recipes/sourdough-zucchini-bread-recipe",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "bread:quick bread"},
				CookingMethod: &models.CookingMethod{},
				CookTime:      "PT1H5M",
				Cuisine:       &models.Cuisine{},
				DatePublished: "2021-06-03",
				Description: &models.Description{
					Value: "This delicious whole grain zucchini bread makes wonderful use of excess sourdough starter you might otherwise discard. Paired with summer’s avalanche of zucchini, it’s one loaf that solves two kitchen conundrums!",
				},
				Keywords: &models.Keywords{
					Values: "Quick bread, Lemon, Raisin, Sourdough, Spice, Breakfast & brunch",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"3/4 cup (170g) sourdough starter fed (ripe) or unfed (discard)",
						"1/2 cup (99g) granulated sugar",
						"1/4 cup (85g) honey",
						"6 tablespoons (75g) vegetable oil",
						"2 large eggs",
						"1/4 teaspoon nutmeg",
						"1 1/2 teaspoons lemon zest (grated rind)",
						"1 1/2 teaspoons King Arthur Pure Vanilla Extract",
						"1 cup (113g) King Arthur White Whole Wheat Flour",
						"3/4 cup (90g) King Arthur Unbleached All-Purpose Flour",
						"1/2 teaspoon baking soda",
						"1 teaspoon baking powder",
						"1 teaspoon table salt",
						"2 cups (242g to 300g) grated zucchini somewhere between firmly and lightly packed",
						"3/4 cup (85g) chopped walnuts lightly toasted",
						"3/4 cup (128g) raisins currants or dried cranberries",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Preheat the oven to 350°F. Lightly grease a 9” x 5” quick bread pan or 12” x 4” tea loaf pan."},
						{Type: "HowToStep", Text: "In a large bowl, stir together the starter, sugar, honey, oil, eggs, nutmeg, lemon zest, and vanilla until thoroughly combined."},
						{Type: "HowToStep", Text: "In a separate medium bowl, whisk together the flours, baking soda, baking powder, and salt; stir into the wet ingredients."},
						{Type: "HowToStep", Text: "Stir in the grated zucchini, then the nuts and fruit. Transfer the batter to the prepared pan, smoothing the top."},
						{Type: "HowToStep", Text: "Bake the bread in the 9” x 5” pan for 45 minutes. Tent with foil and bake for an additional 20 minutes, until a thin paring knife inserted in the center comes out clean. For bread in a tea loaf pan, bake for 40 minutes before tenting, then bake for another 20 minutes, or until the loaf tests done."},
						{Type: "HowToStep", Text: "Remove the bread from the oven and cool in the pan on a rack."},
						{Type: "HowToStep", Text: "Store bread, well wrapped, at room temperature for up to three days; freeze for longer storage."},
					},
				},
				Name: "Sourdough Zucchini Bread",
				NutritionSchema: &models.NutritionSchema{
					Calories:       "279 calories",
					Carbohydrates:  "33g",
					Cholesterol:    "23g",
					Fat:            "13g",
					Fiber:          "2g",
					Protein:        "5g",
					SaturatedFat:   "2g",
					Servings:       "",
					Sodium:         "202mg",
					Sugar:          "21g",
					TransFat:       "0g",
					UnsaturatedFat: "",
				},
				PrepTime:  "PT30M",
				Tools:     &models.Tools{},
				TotalTime: "PT1H30M",
				Yield:     &models.Yield{Value: 16},
				URL:       "https://www.kingarthurbaking.com/recipes/sourdough-zucchini-bread-recipe",
			},
		},
		{
			name: "kingarthurbaking2.com",
			in:   "https://www.kingarthurbaking.com/recipes/fudge-brownies-recipe",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "cookies:brownies"},
				CookTime:      "PT32M",
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				DatePublished: "2009-12-31",
				Description:   &models.Description{Value: "Fudgy, cakey, fudgy, cakey... can't make up your mind? If you're looking for a brownie that's right in between those two styles, you've found it. These brownies combine a fudge brownie's ultra-moist texture with a subtle cake-like rise, for the best of both worlds."},
				Keywords:      &models.Keywords{Values: "Brownies & blondies, Chocolate, Quick & easy"},
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"4 large eggs",
						"1 1/4 cups (106g) Dutch-process cocoa",
						"1 teaspoon table salt",
						"1 teaspoon baking powder",
						"1 teaspoon espresso powder optional for enhanced flavor",
						"1 tablespoon King Arthur Pure Vanilla Extract",
						"16 tablespoons (227g) unsalted butter",
						"2 1/4 cups (447g) granulated sugar",
						"1 1/2 cups (180g) King Arthur Unbleached All-Purpose Flour",
						"2 cups (340g) chocolate chips",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Text: "Preheat the oven to 350°F. Lightly grease a 9\" x 13\" pan. If you plan to turn the whole sheet of brownies out of the pan at once, grease the pan, line it with parchment, and grease the parchment.", Type: "HowToStep"},
						{Text: "Crack the 4 eggs into a bowl, and beat them at medium speed with the cocoa, salt, baking powder, espresso powder, and vanilla for about 1 minute, or until smooth. You can do this while you're melting your butter (next step).", Type: "HowToStep"},
						{Text: "In a medium-sized microwave-safe bowl, or in a saucepan set over low heat, melt the butter, then add the sugar and stir to combine. Or simply combine the butter and sugar, and heat, stirring, until the butter is melted. Continue to heat (or microwave) briefly, just until the mixture is hot (about 110°F to 120°F), but not bubbling; it'll become shiny looking as you stir it. Heating the mixture to this point will dissolve more of the sugar, which will help produce a shiny top crust on your brownies.", Type: "HowToStep"},
						{Text: "Add the hot butter/sugar mixture to the egg/cocoa mixture, stirring until smooth.", Type: "HowToStep"},
						{Text: "Add the flour and chips, stirring until smooth. Again, adding the chips helps produce a shiny top crust.", Type: "HowToStep"},
						{Text: "Spoon the batter into a lightly greased 9\" x 13\" pan.", Type: "HowToStep"},
						{Text: "Bake the brownies for 28 to 32 minutes, until\u00a0the edges feel set, and the center should look very moist, but not uncooked. When testing to see if brownies are done, take a toothpick or the tip of a sharp knife and carefully poke it into the center of the pan, digging around just enough to see the interior. You should see moist crumbs, but no uncooked batter. Yes, you'll be left with a small divot in the center of your brownies; just cut around it when you're cutting the brownies into squares.", Type: "HowToStep"},
						{Text: "Remove them from the oven and cool on a rack before cutting and serving.", Type: "HowToStep"},
						{Text: "Store any leftovers, well-wrapped, at room temperature for 5 to 6 days. Freeze for longer storage.", Type: "HowToStep"},
					},
				},
				Name: "Fudge Brownies",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "260 calories",
					Carbohydrates: "36g",
					Cholesterol:   "55mg",
					Fat:           "17g",
					Fiber:         "2g",
					Protein:       "3g",
					SaturatedFat:  "8g",
					Servings:      "1",
					Sodium:        "130mg",
					Sugar:         "27g",
					TransFat:      "0g",
				},
				PrepTime:  "PT12M",
				Tools:     &models.Tools{},
				TotalTime: "PT40M",
				Yield:     &models.Yield{Value: 24},
				URL:       "https://www.kingarthurbaking.com/recipes/fudge-brownies-recipe",
			},
		},
		{
			name: "kitchenstories.com",
			in:   "https://www.kitchenstories.com/de/rezepte/valencianische-paella",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Vegetarische Gerichte"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{Value: "spanisch und portugiesisch"},
				DatePublished: "2015-03-12T03:00:00+0000",
				Description: &models.Description{
					Value: "Wer sagt das man Meeresfrüchte für eine gute Paella braucht. Lerne wie du eine perfekte Paella nur mit Gemüse machst – schnell & einfach zum Nachkochen!",
				},
				Keywords: &models.Keywords{
					Values: "kinderfreundlich,Brand Content,Alltagsgerichte,vegetarisch,vegan,Vorspeise,Beilagen,Hauptgericht,Party Food,street food,pescetarisch,Wohlfühlessen,laktosefrei,Gewürze,Fleischlos,Le Creuset,thermohauser,spanisch und portugiesisch,herzhaft,für vier,Alkohol,Kräuter,Gemüse,anschwitzen",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"100 g Oliven", "350 g Cherrytomaten", "2 Knoblauch", "2 Frühlingszwiebeln",
						"1 Zwiebel (rot)", "2 Paprika (rot)", "1 Zucchini", "1 Aubergine",
						"250 g Reis", "2 Zitronen", "0.5 TL Safran", "100 ml Weißwein",
						"300 ml Gemüsebrühe", "300 g Erbsen", "Salz", "Pfeffer",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Oliven abgießen. Cherrytomaten halbieren, Knoblauch hacken und Frühlingszwiebeln in feine Ringe schneiden. Zwiebel, Paprika, Zucchini und Aubergine in walnussgroße Stücke schneiden."},
						{Type: "HowToStep", Text: "Etwas Olivenöl in die Pfanne geben und Knoblauch, Zwiebel, Paprika, Zucchini und Aubergine bei mittlerer Hitze ca. 4 – 6 Min. anbraten."},
						{Type: "HowToStep", Text: "Reis in die Pfanne geben und gut umrühren. Für weitere 1 – 2 Min. anbraten."},
						{Type: "HowToStep", Text: "Zitronen filitieren."},
						{Type: "HowToStep", Text: "Safran in die Pfanne geben und gut umrühren, um alles zu vermengen. Für weitere 1 – 2 Min. braten."},
						{Type: "HowToStep", Text: "Mit Weißwein ablöschen. Gemüsebrühe hinzugeben, bis alle Zutaten mit Flüssigkeit bedeckt sind. Nach Geschmack mit Salz und Pfeffer würzen. Alles aufkochen lassen, dann auf ein Köcheln reduzieren und ca. 15 – 20 Min. köcheln lassen, bis der Reis bissfest ist. Gelegentlich umrühren."},
						{Type: "HowToStep", Text: "Vorsichtig die Zitronenfilets, Tomaten, Oliven, Frühlingszwiebeln und Erbsen unterheben. Für weitere 5 – 6 Min. braten. Genießen!"},
					},
				},
				Name: "Vegetarische Paella mit Zucchini und Aubergine",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "448 cal",
					Carbohydrates: "69 g",
					Fat:           "7 g",
					Protein:       "15 g",
					Servings:      "1",
				},
				PrepTime: "PT45M",
				ThumbnailURL: &models.ThumbnailURL{
					Value: "https://images.kitchenstories.io/recipeImages/RP02_18_06_valencianPaella_titlePicture.jpg",
				},
				Tools:     &models.Tools{Values: []models.HowToItem{}},
				TotalTime: "PT45M",
				URL:       "https://www.kitchenstories.com/de/rezepte/valencianische-paella",
				Yield:     &models.Yield{Value: 4},
			},
		},
		{
			name: "kochbar.de",
			in:   "https://www.kochbar.de/rezept/465773/Spargelsalat-Fruchtig.html",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Hauptspeise"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{Value: "Internationale Küche"},
				DatePublished: "2013-04-20T18:30:20+02:00",
				Description:   &models.Description{Value: "lauwarmer Spargel-Salat"},
				Keywords: &models.Keywords{
					Values: "Spargelsalat Fruchtig, Spargel grün frisch, Spargel weiss frisch, Mango frisch",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1 kg Spargel grün frisch",
						"1 Kg Spargel weiss frisch",
						"1 Stück Mango frisch",
						"1 Stück Orange frisch",
						"4 El Olivenöl",
						"Gourmet-Pfeffer aus meinem KB",
						"Salz",
						"Zucker",
						"Räucherlachs",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "für 2 Personen als Hauptspeise \r\nfür 4 Personen als Vorspeise"},
						{Type: "HowToStep", Text: "1. Spargel schälen und in 4-5 Stücke schneiden"},
						{Type: "HowToStep", Text: "3. Spargel inm Salzwasser und wenig Zucker bissfest Kochen. Wasser wegschütten"},
						{Type: "HowToStep", Text: "4. leicht abkühlen lassen und in der Zwischenzeit die Mango schälen und in kleine Würfel schneiden"},
						{Type: "HowToStep", Text: "Dressing"},
						{Type: "HowToStep", Text: "Saft von 1 Orange in eine Schüssel geben und das Olivenöl hinzufügen gut verrühren und mit Salz und Pfeffer abschmecken. Die Spargeln darin wenden und ein wenig ziehen lassen."},
						{Type: "HowToStep", Text: "Schön Anrichten und mit Lachs garnieren ANSTELLE Lachs passen auch wunderbar Crevetten dazu."},
					},
				},
				Name: "Spargelsalat Fruchtig",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "97 kcal",
					Carbohydrates: "1,87273 g",
					Fat:           "9,23273 g",
					Protein:       "1,78182 g",
					Servings:      "100 g",
				},
				Tools: &models.Tools{Values: []models.HowToItem{}},
				Yield: &models.Yield{Value: 2},
				URL:   "https://www.kochbar.de/rezept/465773/Spargelsalat-Fruchtig.html",
			},
		},
		{
			name: "kochbucher.com",
			in:   "https://kochbucher.com/eierlikor-pralinen/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				Description: &models.Description{
					Value: "Zutaten : 200 g Schokolade, weiße 25 g Butter 125 g Mandel(n), gemahlene 75 ml Eierlikör Zubereitung : Arbeitszeit: ca. 1 Std. Ruhezeit: ca. 1 Tag Schwierigkeitsgrad: simpel Kalorien p. P.: keine Angabe Die Schokolade mit der Butter langsam schmelzen. Einen Teil der Mandeln mit unterheben, dann den Eierlikör unterrühren. Am besten über Nacht erkalten […]",
				},
				Keywords: &models.Keywords{},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"200 g Schokolade, weiße", "25 g Butter", "125 g Mandel(n), gemahlene",
						"75 ml Eierlikör",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Arbeitszeit: ca. 1 Std. Ruhezeit: ca. 1 Tag\nSchwierigkeitsgrad: simpel\nKalorien p. P.: keine Angabe"},
						{Type: "HowToStep", Text: "Die Schokolade mit der Butter langsam schmelzen. Einen Teil der Mandeln mit unterheben, dann den Eierlikör unterrühren. Am besten über Nacht erkalten lassen."},
						{Type: "HowToStep", Text: "Mit Hilfe von einem Teelöffel kleine Mengen abstechen und zu Kugeln formen. Die Kugeln anschließend in den restlichen gemahlenen Mandeln wälzen und in Papierförmchen setzen. Kühl aufbewahren!"},
						{Type: "HowToStep", Text: "Das Rezept ergibt ca. 24 Pralinen."},
					},
				},
				Name:            "Eierlikör – Pralinen",
				NutritionSchema: &models.NutritionSchema{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				Yield:           &models.Yield{Value: 1},
				URL:             "https://kochbucher.com/eierlikor-pralinen/",
			},
		},
		{
			name: "koket.se",
			in:   "https://www.koket.se/mitt-kok/tommy-myllymaki/myllymakis-toast-skagen",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    &models.SchemaType{Value: "Recipe"},
				Category:  &models.Category{Value: "uncategorized"},
				Description: &models.Description{
					Value: "Toast skagen är en klassisk förrätt på årets festdag - nyårsafton. Tommys variant görs med hemslagen majonnäs, pepparrot och löjrom.",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1 kg räkor med skal (gärna färska av fin kvalitet)",
						"2 äggulor",
						"2 tsk senap",
						"1 msk vitvinsvinäger",
						"6 dl matolja",
						"1 kruka dill",
						"10 cm färsk pepparrot, skalad",
						"4 skivor vitt bröd (ej levain)",
						"smör, till stekning",
						"50 g löjrom",
						"1 citron",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Skala alla räkor och ställ åt sidan."},
						{Type: "HowToStep", Text: "Gör en majonnäs genom att lägga ner äggulor, senapen och vinägern i en bunke. Tillsätt matoljan i en tunn stråle medan du vispar hela tiden. Använd elvisp eller handvisp. När majonnäsen är tjock och du ser dragen/spåren av vispen i majonnäsen är den klar."},
						{Type: "HowToStep", Text: "Lägg alla räkor i en bunke, tillsätt fint plockad dill och blanda ner lite majonnäs i taget."},
						{Type: "HowToStep", Text: "Tillsätt lite riven pepparrot och smaka av. Slå på mer majonnäs för en rinnigare röra eller mer pepparrot för mer sting."},
						{Type: "HowToStep", Text: "Ta fram brödet och skär ut önskad form utan att ta med kanterna, använd en skål eller ett glas som mall om ni vill ha runda bröd. Stek sedan gyllene i smör."},
						{Type: "HowToStep", Text: "Lägg upp bröden på tallrik, toppa med skagenröra och en rejäl klick löjrom. Avsluta med en dillkvist och en citronskiva."},
					},
				},
				Name:      "Myllymäkis toast skagen",
				TotalTime: "PT25M",
				URL:       "https://www.koket.se/mitt-kok/tommy-myllymaki/myllymakis-toast-skagen",
				Video: &models.Videos{
					Values: []models.VideoObject{
						{
							AtType:       "VideoObject",
							ContentURL:   "https://www.tv4play.se/iframe/video/20277435?autoload=true",
							Description:  "Myllymäkis toast skagen",
							EmbedURL:     "https://www.tv4play.se/iframe/video/20277435?autoload=true",
							Name:         "Myllymäkis toast skagen",
							ThumbnailURL: &models.ThumbnailURL{Value: "https://img.koket.se/standard-mega/myllymakis-toast-skagen-2.jpg"},
							UploadDate:   time.Date(2014, 12, 29, 5, 0, 0, 0, time.UTC),
						},
					},
				},
				Yield: &models.Yield{Value: 4},
			},
		},
		{
			name: "kptncook.com",
			in:   "https://mobile.kptncook.com/recipe/pinterest/empanadas-mit-wuerziger-tomaten-salsa/3f1e5736",
			want: models.RecipeSchema{
				AtContext:       "https://schema.org",
				AtType:          &models.SchemaType{Value: "Recipe"},
				Category:        &models.Category{Value: "uncategorized"},
				CookingMethod:   &models.CookingMethod{},
				NutritionSchema: &models.NutritionSchema{},
				Cuisine:         &models.Cuisine{},
				Description:     &models.Description{},
				Keywords:        &models.Keywords{},
				Image:           &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"0 red onion", "1 chili pepper", "1 cup(s) cilantro, fresh", "1 lime",
						"1 cup(s) cheese, shredded", "1 tomato", "egg", "butter", "salt", "pepper",
						"vegetable oil", "wheat flour", "water",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "All set?"},
						{Type: "HowToStep", Text: "Knead butter, flour, water, egg, and salt into a uniform dough."},
						{Type: "HowToStep", Text: "Wrap in cling film and refrigerate."},
					},
				},
				Name:     "Cheese Empanadas with Fresh Tomato Salsa",
				PrepTime: "PT30M",
				Tools:    &models.Tools{Values: []models.HowToItem{}},
				Yield:    &models.Yield{Value: 2},
				URL:      "https://mobile.kptncook.com/recipe/pinterest/empanadas-mit-wuerziger-tomaten-salsa/3f1e5736",
			},
		},
		{
			name: "kuchnia-domowa.pl",
			in:   "https://www.kuchnia-domowa.pl/przepisy/dodatki-do-dan/548-mizeria",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Dodatki do dań"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{Value: "Polska"},
				Description: &models.Description{
					Value: "Lekka surówka do obiadu ze świeżego ogórka, śmietany lub jogurtu oraz koperku. Bardzo prosta, idealnie nadająca się do wielu dań obiadowych. Mizeria najsmaczniejsza jest z ziemniakami najlepiej młodymi i jakimś mięsem np. kotletem mielonym lub schabowym.\nMy najbardziej lubimy kremową mizerię z miękkimi, cienkimi plasterkami ogórka doprawioną nie tylko solą i pieprzem, ale również (aby była słodko- winna) sokiem z cytryny i cukrem. A jak u Ciebie przygotowuje się mizerię?",
				},
				Keywords: &models.Keywords{Values: "przepis, mizeria, surówka z ogórków, mizeria z octem i śmietaną, tradycyjna mizeria, klasyczna mizeria, domowa mizeria"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"600 g świeżych ogórków gruntowych (lub długich, szklarniowych)*",
						"300 g gęstej, kwaśnej śmietany 18% lub jogurtu typu greckiego",
						"1 łyżeczka soli",
						"1 łyżka soku z cytryny (lub niepełna łyżka octu jabłkowego)",
						"1 łyżeczka cukru",
						"czarny pieprz mielony",
						"1 łyżka drobno posiekanego koperku",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Ogórki umyć, osuszyć, obrać i pokroić w jak najcieńsze plasterki."},
						{Type: "HowToStep", Text: "Plasterki umieścić w misce i posypać 1 łyżeczką soli. Wymieszać i pozostawić na ok. 15 minut."},
						{Type: "HowToStep", Text: "W międzyczasie śmietanę przełożyć do miseczki. Przyprawić sokiem z cytryny, cukrem, pieprzem i posiekanym koperkiem. Wymieszać."},
						{Type: "HowToStep", Text: "Po 15 minutach odlać wodę, którą puściły ogórki. (Lekko je odcisnąć, ale nie za mocno, aby mizeria nie wyszła za sucha)."},
						{Type: "HowToStep", Text: "Dodać przygotowaną śmietanę i wymieszać."},
					},
				},
				Name:            "Mizeria",
				NutritionSchema: &models.NutritionSchema{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				Yield:           &models.Yield{Value: 4},
				URL:             "https://www.kuchnia-domowa.pl/przepisy/dodatki-do-dan/548-mizeria",
			},
		},
		{
			name: "kuchynalidla.sk",
			in:   "https://kuchynalidla.sk/recepty/bravcova-rolada-so-syrom-a-sunkou",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				DatePublished: "2018-06-21T00:25:43+02:00",
				Description: &models.Description{
					Value: "Zatočte s mäsom a zasýťte hladné bruchá svojich najmilších. Pripravte im na obed bravčovú roládu so syrom a šunkou s lahodným zemiakovým pyré podľa Marcela. Budú sa zalizovať až za ušami!",
				},
				Keywords: &models.Keywords{},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"10 plátkov bravčového karé", "morská soľ v mlynčeku Kania",
						"čierne korenie v mlynčeku Kania", "200 g prosciutto cotto Dulano Selection",
						"čerstvá šalvia", "150 g ementálu Milbona (plátky)",
						"olivový olej Primadonna",
						"hladká múka Castello (na obalenie a do výpeku)", "1 lyžica masla Pilos",
						"100 ml suchého bieleho vína", "trochu vývaru alebo vody na podliatie",
						"1 kg zemiakov", "300 ml plnotučného mlieka Pilos", "80 g masla Pilos",
						"čerstvá pažítka", "125 g rukoly", "trochu citrónovej šťavy",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Postup"},
						{Type: "HowToStep", Text: "POTREBUJEME"},
						{Type: "HowToStep", Text: "Bravčové rolády"},
						{Type: "HowToStep", Text: "Bravčové karé prikryjeme fóliou a naklepeme na tenké plátky. Nožom narežeme okraje,"},
						{Type: "HowToStep", Text: "aby mäso počas opekania držalo svoj tvar. Ochutíme soľou a korením. Na každý plátok položíme plátok šunky, lístky šalvie a plátok syra. Zrolujeme a uzatvoríme pomocou špáradiel."},
						{Type: "HowToStep", Text: "Pripravené rolády obalíme v hladkej múke a opečieme na rozohriatom olivovom oleji"},
						{Type: "HowToStep", Text: "po celom obvode do zlatohneda. K výpeku pridáme lyžicu masla, lyžicu hladkej múky, podlejeme vínom a\u00a0trochou vody. Varíme, kým sa z vína odparí alkohol. Potom pridáme nadrobno nasekané lístky šalvie a dochutíme soľou a korením. Stiahneme z ohňa, prikryjeme pokrievkou a necháme chvíľu odpočívať."},
						{Type: "HowToStep", Text: "Zemiaková kaša"},
						{Type: "HowToStep", Text: "Umyté a očistené zemiaky nakrájame na menšie kúsky a uvaríme v osolenej vode domäkka. Potom ich scedíme. V kastróle zohrejeme mlieko spolu s maslom, ktoré postupne prilievame k uvareným"},
						{Type: "HowToStep", Text: "a scedeným zemiakom. Mixérom vymiešame hladkú kašu. Na záver dochutíme nadrobno nasekanou pažítkou."},
						{Type: "HowToStep", Text: "Servírovanie"},
						{Type: "HowToStep", Text: "Na tanier naservírujeme zemiakovú kašu a rolády nakrájané na kolieska. Pokvapkáme výpekom. Podávame spolu s rukolou, ktorú ochutíme soľou, korením, olivovým olejom a citrónovou šťavou."},
					},
				},
				Name:            "Bravčová roláda so syrom a šunkou",
				NutritionSchema: &models.NutritionSchema{},
				PrepTime:        "PT1H0M",
				Tools: &models.Tools{
					Values: []models.HowToItem{
						{Type: "HowToTool", Text: "tĺčik na mäso", Quantity: 1},
						{Type: "HowToTool", Text: "potravinovú fóliu", Quantity: 1},
						{Type: "HowToTool", Text: "špáradlá, mixér", Quantity: 1},
						{Type: "HowToTool", Text: "panvicu s\u00a0pokrievkou", Quantity: 1},
					},
				},
				Yield: &models.Yield{Value: 5},
				URL:   "https://kuchynalidla.sk/recepty/bravcova-rolada-so-syrom-a-sunkou",
			},
		},
		{
			name: "kwestiasmaku.com",
			in:   "https://www.kwestiasmaku.com/przepis/muffiny-czekoladowe-z-maslem-orzechowym",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				DateModified:  "2022-11-05T09:43:22+01:00",
				DatePublished: "2022-03-24T19:55:17+01:00",
				Description: &models.Description{
					Value: "Mocno kakaowe muffiny wzmocnione dodatkową dawką czekolady w postaci dropsów czekoladowych (lub posiekanej czekolady). Dla miłośników masła orzechowego dodajemy do nich po łyżeczce masła orzechowego i rozprowadzamy je w czekoladowej masie za pomocą wykałaczki.\nZ przepisu otrzymamy od 14 do 16 muffinków. Nakładamy do foremek tyle ciasta aby nie wypływało na zewnątrz podczas pieczenia i nie robił się \"grzybek\". W związku z tym, że możemy mieć różne wielkości foremek, najlepiej wypełniać foremki surowym ciastem do 2/3 ich objętości. Pozostawiamy w ten sposób miejsce na wyrośnięcie ciasta i otrzymamy kształtne babeczki.\n",
				},
				Keywords: &models.Keywords{},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"150 g masła",
						"150 g dropsów czekoladowych (np. z ciemnej czekolady) lub 150 g czekolady deserowej lub gorzkiej",
						"300 g mąki",
						"2 łyżeczki proszku do pieczenia",
						"1/2 łyżeczki sody oczyszczonej",
						"3 łyżki kakao",
						"1 szklanka (200 g) cukru",
						"1 łyżka cukru wanilinowego",
						"2 duże jajka (L)",
						"200 ml mleka",
						"ok. 5 - 6 łyżek masła orzechowego",
						"15 - 18 papilotek",
						"metalowa forma na muffiny z wgłębieniami",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Piekarnik nagrzać do 180 stopni C.\u00a0Masło roztopić i przestudzić. Czekoladę pokroić na kawałeczki."},
						{Type: "HowToStep", Text: "Mąkę przesiać do miski razem z proszkiem do pieczenia, sodą i kakao, dokładnie wymieszać. Dodać cukier oraz cukier wanilinowy i ponownie wymieszać."},
						{Type: "HowToStep", Text: "W drugiej misce rozmiksować jajka z mlekiem (rózgą lub mikserem)."},
						{Type: "HowToStep", Text: "Do sypkich składników dodać masę jajeczną i krótko zamieszać łyżką. Dodać roztopione masło i wymieszać do połączenia składników, pod koniec dodając 2/3 ilości dropsów czekoladowych."},
						{Type: "HowToStep", Text: "Masę wyłożyć do papilotek umieszczonych w formie na muffiny, na wierzch wyłożć po łyżeczce masła orzechowego na każdą muffinkę."},
						{Type: "HowToStep", Text: "Wykałaczką zrobić \"ósemkę\" w cieście mieszając delikatnie masę czekoladową z masłem orzechowym. Wierzch posypać pozostałą 1/3 dropsów czekoladowych."},
						{Type: "HowToStep", Text: "Wstawić do piekarnika (można piec na raty, w 2 partiach) i piec\u00a0przez około 20 -\u00a023 minuty, do suchego patyczka."},
					},
				},
				Name:            "Muffiny czekoladowe z masłem orzechowym",
				NutritionSchema: &models.NutritionSchema{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				Yield:           &models.Yield{Value: 15},
				URL:             "https://www.kwestiasmaku.com/przepis/muffiny-czekoladowe-z-maslem-orzechowym",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
