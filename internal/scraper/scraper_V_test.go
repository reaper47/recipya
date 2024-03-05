package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_V(t *testing.T) {
	testcases := []testcase{
		{
			name: "valdemarsro.dk",
			in:   "https://www.valdemarsro.dk/butter_chicken/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Name:          "Butter chicken",
				CookTime:      "PT30M",
				DateModified:  "2023-03-30T10:56:42+00:00",
				DatePublished: "2014-06-25T05:51:40+00:00",
				Description: models.Description{
					Value: "Krydret og lækker butter chicken med ris, indisk salat og godt brød til – det er virkelig en skøn ret.Selv om butter chicken måske umiddelbart tager lidt tid, så er det mest af alt tiden hvor kyllingen skal marinere, for selve arbejdstiden er hverdagsvenlig. Sæt evt kyllingen i marinade allerede aftenen før eller fra morgenstunden – det bliver den blot bedre af.Det er også en herlig weekendret og god til gæster, hvor man kan servere flere skål med lækkert indisk mad at samles om.Prøv også: min bedste opskrift på lækkert naan brød >>",
				},
				Image:    models.Image{Value: anUploadedImage.String()},
				PrepTime: "PT1H30M",
				Yield:    models.Yield{Value: 4},
				Ingredients: models.Ingredients{
					Values: []string{
						"100 g græsk yoghurt 10 %",
						"2 tsk chiliflager",
						"0,50 tsk stødt nellike",
						"2 tsk stødt spidskommen",
						"1 tsk stødt kardemomme",
						"2 tsk garam masala",
						"2 fed hvidløg, finthakkede",
						"1 tsk stødt gurkemeje",
						"1 spsk ingefær, friskrevet",
						"1 dåse hakkede tomater",
						"500 g kyllingebryst, skåret i tern",
						"1 løg, skåret i ringe",
						"0,50 dl piskefløde",
						"50 g smør",
						"2 spsk olivenolie",
						"flagesalt",
						"sort peber, friskkværnet",
						"3 dl basmati ris, kogt efter anvisning på emballagen",
						"2 håndfulde frisk koriander",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Blend yoghurt, chili, nellike, spidskommen, kardemomme, garam masala, gurkemeje, ingefær, hvidløg og " +
							"hakkede tomater sammen til en lind sauce.",
						"Hæld den over kyllingestykkerne, dæk dem til og lad dem trække i min. 30 minutter i køleskabet og gerne " +
							"natten over eller fra morgenstund til aftensmadstid.",
						"Smelt smør og olie i en sauterpande. Sautér løgene, til de bliver bløde og blanke.",
						"Tilsæt kylling, sauce og fløde.",
						"Lad det simre ved lav varme i 30-35 minutter, eller til kyllingen er mør. Server med ris og et drys frisk " +
							"koriander",
					},
				},
				URL: "https://www.valdemarsro.dk/butter_chicken/",
			},
		},
		{
			name: "vanillaandbean.com",
			in:   "https://vanillaandbean.com/carrot-cake-bread/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Brunch"},
				CookTime:      "PT60M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-04-08T17:54:24+00:00",
				Description: models.Description{
					Value: "Carrot Cake Bread  is cake in carrot cake loaf form! It&#039;s just as easy to make as my moist Carrot " +
						"Cake Recipe with Pineapple, but because carrot cake quick bread is made in a loaf pan, it&#039;s a bit " +
						"more casual. Share with an orange spiked cream cheese frosting, or enjoy without. Either way, it&#039;s " +
						"perfect for all your occasions and general snacking!*Time above does not include two hours to cool the " +
						"cake completely before icing. See blog post for more tips.",
				},
				Keywords: models.Keywords{
					Values: "Carrot Cake Bread, Carrot Cake Loaf Recipe",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"3/4 cup (80 grams) Pecans (or Walnuts)",
						"1 1/4 cup (175 grams) Unbleached All Purpose Flour (I use Bob&#039;s Red Mill)",
						"1 tablespoon Corn Starch",
						"1/2 teaspoon Ground Ginger",
						"1 teaspoon Ground Cinnamon",
						"1/2 teaspoon Fine Sea Salt",
						"1 1/4 teaspoon Baking Powder",
						"1/2 teaspoon Baking Soda",
						"1/3 cup (70 grams) Dark Brown Sugar (or light, packed)",
						"1/3 cup (75 grams) Cane Sugar",
						"3/4 cup (150 grams) Coconut Oil (melted and warm to touch, or Canola Oil)",
						"1 Large Orange (zested and 2 tablespoons of orange juice)",
						"2 Whole Eggs (room temperature* see note)",
						"2 teaspoons Vanilla Extract",
						"1/3 cup (75 grams) Unsweetened Whole Milk Yogurt or Greek Yogurt",
						"1 1/4 cups (160 grams) Finely Shredded Carrots (packed, about 3 medium carrots)",
						"5 tablespoons Unsalted Butter (room temperature)",
						"7 ounces (200 grams) Cream Cheese (room temperature)",
						"1/3 cup (40g) Powdered Sugar (sifted)",
						"1 Large Orange (zested)",
						"1/2 Vanilla Bean (scraped, or use 1/2 tsp vanilla bean paste or 1 teaspoon vanilla extract)",
						"Orange Juice (fresh squeezed, a few teaspoons for loosening the icing if needed.)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Arrange an oven rack in the center of the oven. Preheat the oven to 350F (180C). Place the nuts on a small " +
							"sheet pan and put them in the oven while the oven is preheating. Toast the nuts for about 15 minutes, " +
							"or until the nuts are fragrant and starting to turn slightly darker. Chop fine when cool. Set aside.",
						"Grease a 9 inch by 5 inch (23 centimeters by 13 centimeters) loaf pan thoroughly on the all sides and b" +
							"ottom. Be sure to get the corners good! Line the pan with parchment paper, just one strip across the " +
							"bottom and sides will act as a handle to remove the bread once baked. Clip the sides of the parchment " +
							"to the pan if desired.",
						"In a medium mixing bowl, whisk the flour, corn starch, ginger, cinnamon, salt, baking powder, and baking " +
							"soda. Set aside.",
						"In a large mixing bowl, whisk the brown sugar, cane sugar and oil until ingredients are combined, about 30 " +
							"seconds. Add the eggs, orange zest and juice and vanilla. Whisk until the ingredients are emulsified, " +
							"about 20 seconds.",
						"Add the dry ingredients to the egg/sugar mixture. Using a silicone spatula, gently mix/fold until no flour " +
							"streaks remain. If using coconut oil, the batter will start to stiffen at this point because the coconut " +
							"oil is cooling (solidifying). Fold in the yogurt until no white streaks remain.Working quickly, fold in " +
							"the shredded carrots and nuts. If using canola oil, the batter should be thick but loose and almost pourable " +
							"(pictured in the photos). For coconut oil, the batter should be stiff and thick.",
						"Transfer the cake batter to the loaf pan and using an offset spatula, smooth the batter into an even layer. " +
							"If using coconut oil, spread the batter while pressing it into the pan. Tap the pan on the counter to " +
							"disperse any air pockets. Remove the parchment clips if using.",
						"Bake the cake for 60-70 minutes. The loaf is ready when it&#039;s golden, a toothpick poked in the center " +
							"comes out clean, the cake slightly springs back under gentle pressure at center, and the edges of the " +
							"loaf are just starting to pull away from the sides of the pan. Allow cake to cool in the pan for 20 " +
							"minutes on a cooling rack. Lift cake out and transfer to a cooking rack. Cool completely at room " +
							"temperature before icing or covering (about 2 hours).",
						"Start with ingredients at room temperature. Place butter in mixer and beat on high with the paddle " +
							"attachment until light and fluffy, about three minutes. Scrape down your bowl several times to make sure " +
							"the butter is getting whipped. Add the cream cheese and beat another minute. Scrape down the bowl.",
						"Sift in the powder sugar, add the vanilla, and orange zest then mix on low until incorporated, about 30 " +
							"seconds. Taste for sweet adjustment and add another tablespoon of powered sugar if desired. To loosen the " +
							"icing, you can add a few teaspoons of orange juice.",
						"Store the icing in a lidded container, in the refrigerator, if not icing the bread after it cools. Bring " +
							"to room temperature before slathering!",
						"Because cream cheese needs to be refrigerated, consider how/when you&#039;ll share the cake before icing it. " +
							"An iced cake needs to be consumed, refrigerated or frozen within three hours. An uniced cake can set " +
							"covered at room temperature for up to three days.Ice the Cake:Once the loaf cake is cool and just before " +
							"serving, spread the icing evenly over the top of the cake. Slice and share within three hours.Individual " +
							"Pieces:Once the loaf is cool, slice the loaf and share with a dollop or slather of cream cheese icing. " +
							"*Uniced carrot cake bread is lovely toasted or warmed, then slathered with cream cheese icing!",
						"Room Temperature: Store an iced loaf cake at room temperature (70F or less) for up to three hours. " +
							"Afterwards, store the cake in the fridge, covered for up to three days. Before serving, pull the cake " +
							"from the fridge and rest at room temperature for about 30 minutes to allow the fats to soften. Food safety " +
							"says cheese should set out for no more than three hours. An uniced loaf cake can be stored covered at room " +
							"temperature for up to three days.",
						"To Freeze: This cake freezes beautifully, iced or uniced. Simply allow the cake to cool completely, ice the " +
							"cake or not, then freeze individual pieces on a sheet pan. Once frozen wrap snugly in plastic wrap or " +
							"store in a lidded container. Store in freezer then thaw at room temperature for about two hours before " +
							"enjoying. " +
							"If iced, take the plastic wrap off before thawing so the icing doesn&#039;t stick.I&#039;ve only tested " +
							"freezing " +
							"the iced cake for up to two days. Uniced will freeze well for up to two weeks. Unwrap and thaw at room " +
							"temperature, " +
							"covered with a cake dome.",
					},
				},
				Name: "Carrot Cake Bread Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories:       "476 kcal",
					Carbohydrates:  "38 g",
					Cholesterol:    "68 mg",
					Fat:            "35 g",
					Fiber:          "2 g",
					Protein:        "6 g",
					SaturatedFat:   "22 g",
					Servings:       "1",
					Sodium:         "263 mg",
					Sugar:          "22 g",
					TransFat:       "1 g",
					UnsaturatedFat: "11 g",
				},
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 10},
				URL:      "https://vanillaandbean.com/carrot-cake-bread/",
			},
		},
		{
			name: "vegetarbloggen.no",
			in:   "https://www.vegetarbloggen.no/2023/07/15/peanottkake/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "bakst"},
				CookTime:      "PT30M",
				Cuisine:       models.Cuisine{Value: "amerikansk"},
				DatePublished: "2023-07-15T11:33:37+00:00",
				Keywords:      models.Keywords{Values: "peanøttkake"},
				Image:         models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"100 g vegansk margarin", "250 g peanøttsmør", "200 g melkefri sjokolade",
						"200 g brunt sukker", "50 g hvitt sukker",
						"1,5 dl vegansk yoghurt (naturell eller med vaniljesmak)",
						"1 ts ren vanilje/vaniljepulver", "1 ts maisstivelse", "1 ts bakepulver",
						"250 g hvetemel",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Forvarm stekeovnen til 175 grader, og kle en form på ca. 20*20 cm med bakepapir.",
						"Smelt margarin og peanøttsmør i en liten gryte. Grovhakk sjokoladen.",
						"Ha margarin og peanøttsmør over i en bolle. Rør inn brunt og hvitt sukker, yoghurt, vanilje, maisstivelse og bakepulver. Jobb det godt sammen til en helt jevn røre.",
						"Ha i hvetemel, og rør det raskt inn. Det er viktig å ikke blande for mye akkurat her. Vend deretter inn den hakka sjokoladen.",
						"Ha deigen over i formen, og press den ut slik at den får en jevn overflate. Sett i ovnen og stek i 25-30 minutter.",
						"Avkjøl litt før du skjærer opp i serveringsstykker (disse er gode både lune og kalde).",
					},
				},
				Name:     "Peanøttkake",
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 16},
				URL:      "https://www.vegetarbloggen.no/2023/07/15/peanottkake/",
			},
		},
		{
			name: "vegolosi.it",
			in:   "https://www.vegolosi.it/ricette-vegane/pancake-vegani-senza-glutine-alla-quinoa-e-cocco/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT15M",
				DatePublished: "2018-08-02",
				Description: models.Description{
					Value: "I pancake vegani senza glutine alla quinoa e cocco sono una deliziosa colazione che permette anche " +
						"a chi deve evitare il glutine di gustarsi dei morbidi e golosissimi pancake, completati in questo " +
						"caso dall'immancabile sciroppo d'acero, fragole fresche e cocco in scaglie.",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"120 g di farina di quinoa",
						"60 g di farina di cocco",
						"20 g di amido di mais",
						"30 g di zucchero di canna grezzo",
						"2 cucchiaini di lievito (cremor tartaro)",
						"2 cucchiai di farina di semi di lino",
						"½ cucchiaino di bicarbonato",
						"400 g di latte di soia",
						"30 g di olio di semi di girasole",
						"1 cucchiaino di aceto di mele",
						"Fragole",
						"Sciroppo d&#8217;acero",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Mescolate la farina di semi di lino con 6 cucchiai di acqua e lasciate riposare 10 minuti fino a che si " +
							"sarà formato un composto gelatinoso. In una ciotolina versate il latte di soia e l&#8217;aceto d" +
							"i mele e lasciate cagliare per 5 minuti. Riunite in una ciotola la farina di quinoa, la farina di " +
							"cocco, l&#8217;amido di mais, lo zucchero di canna, il lievito e il bicarbonato e mescolate. " +
							"Aggiungete agli ingredienti secchi il latte di soia, il composto di semi di lino e l&#8217;olio e " +
							"amalgamate bene il tutto fino ad ottenere un composto omogeneo e abbastanza denso.",
						"Scaldate una padella antiaderente e ungetela leggermente con un pezzo di carta assorbente imbevuto di " +
							"olio di semi. Versate 1-2 cucchiaiate di impasto per ciascun pancake e lasciate cuocere a fiamma " +
							"medio-bassa per 3-4 minuti per lato. Man mano che i vostri pancake saranno pronti disponeteli su un " +
							"piatto, e completate poi ciascuna porzione con sciroppo d&#8217;acero a piacere, fragole fresche e cocco " +
							"in scaglie.",
					},
				},
				Name:     "Pancake vegani senza glutine alla quinoa e cocco",
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.vegolosi.it/ricette-vegane/pancake-vegani-senza-glutine-alla-quinoa-e-cocco/",
			},
		},
		{
			name: "vegrecipesofindia.com",
			in:   "https://www.vegrecipesofindia.com/paneer-butter-masala/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Course"},
				CookTime:      "PT30M",
				Cuisine:       models.Cuisine{Value: "North Indian"},
				DatePublished: "2022-10-31T23:35:31+00:00",
				Description: models.Description{
					Value: "Paneer Butter Masala Recipe is one of India’s most popular paneer preparation. This restaurant style recipe with soft paneer cubes dunked in a creamy, lightly spiced tomato sauce or gravy is a best one that I have been making for a long time. This rich dish is best served with roti or chapati, paratha, naan or rumali roti.",
				},
				Keywords: models.Keywords{Values: "Paneer Butter Masala"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"18 to 20 cashews (- whole)",
						"⅓ cup hot water (- for soaking cashews)",
						"2 cups tomatoes (- diced, or 300 grams tomatoes or 4 to 5 medium sized, pureed)",
						"1 inch ginger (- peeled and roughly chopped)",
						"3 to 4 garlic cloves (- small to medium-sized, peeled)",
						"2 tablespoons Butter (or 1 tablespoon oil + 1 or 2 tablespoons butter)",
						"1 tej patta ((Indian bay leaf), optional)",
						"½ to 1 teaspoon kashmiri red chili powder (or deghi mirch or ¼ to ½ teaspoon cayenne pepper or paprika)",
						"1.5 cups water (or add as required)",
						"1 inch ginger (- peeled and julienned, reserve a few for garnish)",
						"1 or 2 green chili (- slit, reserve a few for garnish)",
						"200 to 250 grams Paneer ((Indian cottage cheese) - cubed or diced)",
						"1 teaspoon dry fenugreek leaves ((kasuri methi) - optional)",
						"½ to 1 teaspoon Garam Masala (or tandoori masala)",
						"2 to 3 tablespoons light cream (or half &amp; half or 1 to 2 tablespoons heavy cream - optional)",
						"¼ to 1 teaspoon sugar (- optional, add as required depending on the sourness of the tomatoes)",
						"salt (as required)",
						"1 to 2 tablespoons coriander leaves (- chopped, (cilantro) - optional)",
						"1 inch ginger (- peeled and julienned)",
						"1 tablespoon light cream (or 1 tablespoon heavy cream - optional)",
						"1 to 2 teaspoons Butter (- optional)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Soak cashews in a hot water for 20 to 30 minutes. When the cashews are soaking, you can prep the other " +
							"ingredients like chopping tomatoes, preparing ginger-garlic paste, slicing paneer etc.",
						"Then drain and add the soaked cashews in a blender or mixer-grinder.",
						"Add 2 to 3 tablespoons water and blend to a smooth and fine paste without any tiny bits or chunks of cashews.",
						"In the same blender add the roughly chopped tomatoes. No need to blanch the tomatoes before blending.",
						"Blend to a smooth tomato puree. Set aside. Don’t add any water while blending the tomatoes.",
						"Melt butter in a pan on a low heat. Add tej patta and fry for 2 to 3 seconds or till the oil become fragrant.",
						"Add ginger-garlic paste and sauté for about 10 to 12 seconds till the raw aroma disappears.",
						"Add the tomato puree and stir well. Cook for 5 to 6 minutes stirring a few times.",
						"Next add kashmiri red chili powder and stir again. Continue to sauté till the oil starts to leave the sides of the tomato paste. The tomato paste will thicken considerably and will start coming together as one whole lump.",
						"Then add cashew paste and stir well. Sauté the cashew paste for a few minutes till the oil begins to leave " +
							"the sides of the masala paste.",
						"The cashew paste will begin to cook fast. Approx 3 to 4 minutes on a low heat. So keep stirring non-stop.",
						"Add water and mix very well. Simmer on a low to medium-low heat.",
						"The curry will come to a boil.",
						"After 2 to 3 minutes of boiling, add ginger julienne. Reserve a few for garnishing. The curry will also " +
							"begin to thicken.",
						"Add julienned ginger and green chillies, salt and sugar and simmer till the curry begins to thicken.",
						"After 3 to 4 minutes, add slit green chillies.also add salt as per taste and ½ to 1 teaspoon sugar (optional).",
						"You can vary the sugar quantity from ¼ tsp to 1 teaspoon or more depending on the sourness of the tomatoes. " +
							"Sugar is optional and you can skip it too. If you add cream, then you will need to add less sugar.",
						"Mix very well and simmer for a minute.",
						"After the gravy thickens to your desired consistency, then add the paneer cubes and stir gently.I keep the " +
							"gravy to a medium consistency.",
						"After that add crushed kasuri methi (dry fenugreek leaves), garam masala and cream. Gently mix and then " +
							"switch off the heat.",
						"Garnish the curry with coriander leaves and ginger julienne.",
						"You can even dot the gravy with some butter or drizzle some cream.",
						"Serve Paneer Butter Masala hot with plain naan, garlic naan, roti, paratha or steamed basmati or jeera rice " +
							"or even peas pulao.",
						"Side accompaniments can be an onion-cucumber salad or some pickle. Also serve some lemon wedges by the side.",
					},
				},
				Name: "Paneer Butter Masala Recipe (Restaurant Style)",
				NutritionSchema: models.NutritionSchema{
					Calories:      "307 kcal",
					Carbohydrates: "9 g",
					Cholesterol:   "66 mg",
					Fat:           "27 g",
					Fiber:         "2 g",
					Protein:       "9 g",
					SaturatedFat:  "15 g",
					Servings:      "1",
					Sodium:        "493 mg",
					Sugar:         "4 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.vegrecipesofindia.com/paneer-butter-masala/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
