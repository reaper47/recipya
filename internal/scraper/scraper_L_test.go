package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_L(t *testing.T) {
	testcases := []testcase{
		{
			name: "latelierderoxane.com",
			in:   "https://www.latelierderoxane.com/blog/recette-cake-marbre",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT45M",
				DatePublished: "2023-11-12T07:54:40+01:00",
				Description: models.Description{
					Value: "Une recette facile, rapide et adorée des enfants : le cake marbré moelleux au chocolat façon Savane. Un cake parfumé à la vanille et au chocolat. ",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 œufs", "70 g de sucre", "70 g de beurre  fondu",
						"1 sachet de levure chimique", "250 g de farine", "150 g de lait",
						"150 g de chocolat noir fondu", "1 càc d'arôme ou poudre de vanille",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Préchauffe le four à 165°.",
						"Commence par fouetter les œufs et le sucre, à l’aide de ton robot ou batteur électrique, pendant 10 minutes : ton mélange doit s’éclaircir et doubler de volume !",
						"Ajoute le beurre fondu, la levure, la farine et fouette brièvement.",
						"Verse le lait et fouette jusqu’à l’obtention d’un mélange homogène.",
						"Sépare la préparation obtenue dans deux bols.",
						"Dans un des deux bols, ajoute l’arôme ou la poudre de vanille.",
						"Fais fondre ton chocolat, au bain-marie ou au micro-ondes et incorpore-le dans le second bol à l’aide d’une maryse.",
						"Récupère ton moule à cake et beurre-le généreusement.",
						"Verse, dans le fond de ton moule, la moitié de la pâte à la vanille puis la moitié de celle au chocolat.",
						"Répète l’opération une deuxième fois.", "Enfourne pendant 45 min.",
						"Tu peux vérifier la cuisson à l’aide d’un couteau, plante-le au centre de ton cake : ta lame doit ressortir sèche.",
						"À la sortie du four, laisse tiédir ton cake afin de faciliter son démoulage.",
						"À manger sans modération !",
					},
				},
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://www.latelierderoxane.com/blog/recette-cake-marbre",
			},
		},
		{
			name: "leanandgreenrecipes.net",
			in:   "https://leanandgreenrecipes.net/recipes/italian/spaghetti-squash-lasagna/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: `<a href="/recipes/category/main-course/" hreflang="en">Main course</a>`},
				CookTime:      "P0Y0M0DT0H64M0S",
				Cuisine:       models.Cuisine{Value: `<a href="/recipes/italian/" hreflang="en">Italian</a>`},
				DatePublished: "2021-03-10T13:26:59-0600",
				Description: models.Description{
					Value: "<p>If you're not familiar with spaghetti squash then it's time you get acquainted! This simple to prepare Spaghetti Squash Lasagna recipe will have you wondering why you haven't been eating spaghetti this way your entire life. Light, delicious and 100% on plan!</p>\n<p> </p>\n<p>Tip: If you do not like spicy heat you can reduce or omit the crushed red pepper flakes. This recipe call for several of the spices to be divided.</p>",
				},
				Image: models.Image{
					Value: "https://leanandgreenrecipes.net/sites/default/files/2021-03/Spaghetti-Squash-Lasagna.jpeg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 medium Spaghetti Squash", "4 tsp Olive Oil", "1 tsp Salt",
						"1 tsp Black Pepper", "2 tsp Garlic", "1 lbs Lean Ground Turkey",
						"1(14.5oz.) can Diced Tomatoes", "1/2 tsp Basil", "1 tsp Whole Leaf Oregano",
						"1/2 cup Part-skim Ricotta", "1/2 cup 1% Cottage Cheese",
						"1 tsp Crushed Red Pepper Flakes", "1 cup Low-fat Mozzarela",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 400*F.", "Prepare the spaghetti squash. Cut in half",
						"and remove seeds and pulp strands. Rub 1 teaspoon of olive oil into each half of squash and season each half with 1/4 teaspoon each of salt (optional) and pepper. Place each spaghetti squash half face down in a large baking dish and bake for 40 to 60 minutes",
						"cook until the middle is tender and pulls apart easily.",
						"While the spaghetti squash is cooking in a large saucepan",
						"saute garlic in remaining olive oil over a medium heat until fragrant. Add Turkey. Season with 1/4 teaspoon of each salt (optional) and pepper",
						"and then cook until it has browned.", "Add tomatoes", "onion powder",
						"and 1/2 teaspoon each of basil and oregano. When the sauce starts to bubble",
						"reduce heat to a simmer until it has thickened (about 3 to 4 minutes).",
						"Combine the ricotta and cottage cheese into a medium bow. Season with crushed red pepper flakes (optional) and the remaining basil",
						"oregano", "salt (optional)", "and pepper. Lightly mix until combined.",
						"When spaghetti squash is fully cooked",
						"flip it in the baking dish so that it is now skin-side down. Lightly scrape flesh of squash with a fork to create spaghetti-like strands.",
						"Evenly divide the ricotta and cheese mixture between each squash half. Repeat with the meat sauce. Top each half with 1/2 cup of mozzarella cheese.",
						"Turn oven to broil",
						"and cook for an additional 2 minutes",
						"or until cheese is browned and bubbling. Watch that it does not burn due to different oven temperatures. Serve immediately",
					},
				},
				Name: "Healthy Spaghetti Squash Lasagna Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories:      "337",
					Carbohydrates: "25",
					Fat:           "16",
					Protein:       "26",
				},
				PrepTime: "P0Y0M0DT0H10M0S",
				Yield:    models.Yield{Value: 4},
				URL:      "https://leanandgreenrecipes.net/recipes/italian/spaghetti-squash-lasagna/",
			},
		},
		{
			name: "lecker.de",
			in:   "https://www.lecker.de/gemuesepfanne-mit-haehnchen-zuckerschoten-und-brokkoli-79685.html",
			want: models.RecipeSchema{
				AtContext:    "https://schema.org",
				AtType:       models.SchemaType{Value: "Recipe"},
				Category:     models.Category{Value: "Hauptgericht"},
				CookTime:     "PT25M",
				DateModified: "2023-02-23T13:12:02.767Z",
				Description: models.Description{
					Value: "Unser beliebtes Rezept für Gemüsepfanne mit Hähnchen, Zuckerschoten und Brokkoli und mehr als 45.000 weitere kostenlose Rezepte auf LECKER.de.",
				},
				Keywords: models.Keywords{
					Values: "Hähnchen,Geflügel,Fleisch,Zutaten,Mittagessen,Mahlzeit,Rezepte,Abendbrot,Low Carb,Gesundes Essen,Hauptgerichte,Menüs,Brokkoli,Kohl,Gemüse,Zuckerschoten,Gemüsepfanne,Pfannengerichte",
				},
				Image: models.Image{
					Value: "https://images.lecker.de/gemusepfanne-mit-hahnchen-zuckerschoten-und-brokkoli/1x1,id=7e976162,b=lecker,w=1600,h=,ca=12.42,0,87.58,100,rm=sk.jpeg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 Brokkoli", "150 g Zuckerschoten", "2 Lauchzwiebeln", "Salz",
						"4 kleine Hähnchenbrustfilets (à ca. 140 g)", "2 EL Sonnenblumenöl",
						"4 EL Sojasoße", "50 ml Gemüsebrühe", "6 Stiele Koriander",
						"1 EL Sesamsaat",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Brokkoli putzen und in kleinen Röschen vom Strunk schneiden. Zuckerschoten putzen. Lauchzwiebeln putzen und in ca. 10 cm lange Stücke schneiden.",
						"Brokkoli, Zuckerschoten und Lauchzwiebeln in reichlich kochendem Salzwasser 2–3 Minuten garen. Herausnehmen und mit kaltem Wasser abspülen. Auf ein Sieb gießen und gut abtropfen lassen.",
						"Fleisch trocken tupfen und in Würfel (ca. 1,5 x 1,5 cm) schneiden. Öl in einer weiten Pfanne oder einem Wok erhitzen. Fleisch darin, unter Wenden, ca. 3 Minuten kräftig anbraten. Hitze reduzieren und weitere ca. 2 Minuten braten. Gemüse, Sojasoße und Brühe zufügen und 3–4 Minuten dünsten. Dabei gelegentlich wenden.",
						"Koriander waschen und trocken schütteln und, bis auf einige Blätter zum Garnieren, samt Stiel grob hacken. Sesam und Koriander zur Gemüsepfanne geben und untermischen. Mit Salz und Pfeffer abschmecken und mit Koriander garnieren.",
					},
				},
				Name: "Gemüsepfanne mit Hähnchen, Zuckerschoten und Brokkoli",
				NutritionSchema: models.NutritionSchema{
					Calories:      "260 kcal",
					Carbohydrates: "7 g",
					Fat:           "7 g",
					Protein:       "38 g",
					Servings:      "1",
				},
				PrepTime: "PT0M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.lecker.de/gemuesepfanne-mit-haehnchen-zuckerschoten-und-brokkoli-79685.html",
			},
		},
		/*{
			name: "lekkerensimpel.com",
			in:   "https://www.lekkerensimpel.com/gougeres/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "Gougères",
				Category:  models.Category{Value: "Snacks"},
				Yield:     models.Yield{Value: 4},
				URL:       "https://www.lekkerensimpel.com/gougeres/",
				PrepTime:  "PT20M",
				CookTime:  "PT25M",
				Description: models.Description{
					Value: "Vandaag een receptje uit de Franse keuken, namelijk deze gougères. Gougères zijn een soort " +
						"hartige kaassoesjes, erg lekker! We hadden een tijdje geleden een stuk gruyère kaas " +
						"gekocht bij de kaasboer, meer uit nood want parmezaanse had hij op dat moment even niet. " +
						"Inmiddels lag de kaas al een tijdje in de koelkast en moesten we er toch echt wat mee gaan " +
						"doen. Iemand tipte ons dat we echt eens gougères moesten maken en eerlijk gezegd hadden we " +
						"er nooit eerder van gehoord. Een kleine speurtocht bracht ons uiteindelijk bij een recept " +
						"van ‘The Guardian – How to make the perfect gougères‘. We zijn ermee aan de slag gegaan en " +
						"zie hier het resultaat! \n\nNog meer van dit soort lekkere snacks en borrelhapjes vind je " +
						"in onze categorie tapas recepten en tussen de high-tea recepten.",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"90 gr gruyere kaas",
						"125 ml water",
						"40 gr boter",
						"75 gr bloem",
						"3 eieren",
						"nootmuskaat",
						"zout",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Verwarm de oven voor op 200 graden. Klop vervolgens 2 eieren los in een beker. Doe het water, de boter " +
							"en een snuf zout in een pan en laat de boter al roerend smelten. Zet het ‘vuur’ laag en doe " +
							"de bloem erbij. Zelf doen we de bloem eerst door een zeef zodat er geen kleine klontjes meer " +
							"inzitten. Roer de bloem door het botermengsel totdat er een soort deeg ontstaat. Haal de pan " +
							"van het vuur en mix het deeg, bij voorkeur met een mixer, een minuut of 3-4. Voeg dan de " +
							"helft van het losgeklopte ei toe, even goed mengen en dan kan de andere helft erbij. Mix " +
							"daarna nog de nootmuskaat en geraspte gruyère door het deeg. Bekleed een bakplaat met " +
							"bakpapier. Schep met twee lepels kleine bolletjes deeg op de bakplaat of gebruik hiervoor " +
							"een spuitzak. Smeer de bovenkant in met een beetje losgeklopt ei, bestrooi met nog geraspte " +
							"gruyère kaas en dan kan de bakplaat de oven in voor 20-25 minuten. Eet smakelijk!",
						"Bewaar dit recept op Pinterest !",
					},
				},
				DatePublished: "2021-09-28T04:00:00+00:00",
				DateModified:  "2021-09-21T08:22:19+00:00",
			},
		},*/
		{
			name: "lifestyleofafoodie.com",
			in:   "https://lifestyleofafoodie.com/chick-fil-a-peppermint-milkshake/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Drinks"},
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-11-18T18:26:00+00:00",
				Description: models.Description{
					Value: "These Chick-Fil-A peppermint milkshakes are an easy and delicious way to enjoy your favorite holiday drink year round!",
				},
				Keywords: models.Keywords{Values: "chick fil a peppermint milkshake, peppermint milkshake"},
				Image: models.Image{
					Value: "https://lifestyleofafoodie.com/wp-content/uploads/2023/11/Chick-fil-a-peppermint-milkshake-5.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 cups vanilla ice cream", "1/4 cup milk (whole or 2% for creamier results)",
						"1/4 teaspoon peppermint extract (adjust to taste)",
						"1/3 cup crushed peppermint candies (plus extra for garnish)",
						"1/4 cup chocolate chips",
						"2-3 drops of red food coloring (optional, for a pinkish color)",
						"Whipped cream (for topping)", "Maraschino cherries (for garnish)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a blender, combine the vanilla ice cream, milk, peppermint extract, crushed peppermint candies (roughly crushed if using a high-speed blender), chocolate chips, and red food coloring (if using).",
						"Blend the mixture until it&#39;s smooth and all the ingredients are well combined. You can adjust the milk or ice cream to achieve your desired milkshake consistency.",
						"Taste the milkshake and adjust the amount of peppermint extract and crushed candies to your liking. Some people prefer a stronger peppermint flavor, while others like it milder.",
						"Pour the milkshakes into your serving glasses, top them with whipped cream and a maraschino cherry and enjoy.",
					},
				},
				Name: "Chick Fil A Peppermint Milkshake",
				NutritionSchema: models.NutritionSchema{
					Calories:       "649 kcal",
					Carbohydrates:  "90 g",
					Cholesterol:    "88 mg",
					Fat:            "28 g",
					Fiber:          "1 g",
					Protein:        "8 g",
					SaturatedFat:   "17 g",
					Servings:       "1",
					Sodium:         "171 mg",
					Sugar:          "76 g",
					UnsaturatedFat: "7 g",
				},
				PrepTime: "PT1M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://lifestyleofafoodie.com/chick-fil-a-peppermint-milkshake/",
			},
		},
		{
			name: "littlespicejar.com",
			in:   "https://littlespicejar.com/starbucks-pumpkin-loaf/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Bread & Baking"},
				CookTime:      "PT55M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2021-11-09",
				Description: models.Description{
					Value: "Learn how to make an easy delicious copycat Starbucks Pumpkin Loaf right at home! This pumpkin " +
						"bread is studded with roasted pepitas and loaded with spices and so much pumpkin goodness!",
				},
				Keywords: models.Keywords{Values: ""},
				Image: models.Image{
					Value: "https://littlespicejar.com/wp-content/uploads/2021/11/Copycat-Starbucks-Pumpkin-Loaf-8-720x720.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 ½ cups all-purpose flour",
						"1 tablespoon cinnamon",
						"2 teaspoons EACH: baking soda AND ground ginger",
						"1 teaspoon EACH: baking powder AND ground allspice",
						"½ teaspoon EACH: ground cloves, ground cardamom, AND ground nutmeg",
						"¾ teaspoon kosher salt",
						"1 cup EACH: granulated sugar AND coconut oil (or other; see post)",
						"1 ½ cups light brown sugar",
						"1 (15-ounce) can pumpkin puree",
						"4 large eggs, room temperature*",
						"1 tablespoon vanilla extract",
						"Zest of 1 orange (optional)",
						"4-5 tablespoons roasted pepitas, for topping",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"PREP: Position a rack in the center of the oven and preheat the oven to 350°F. Spray two 8 ½ x 4 ½ " +
							"(or 9x5) bread pans with cooking spray, you can also line with parchment if you’d like; set aside for now.",
						"DRY INGREDIENTS: Add the dry ingredients: flour, baking soda baking powder, all the spices, and salt " +
							"to a medium bowl. Whisk to combine; set aside for now.",
						"WET INGREDIENTS: Add the granulated sugar, brown sugar, and oil to a large bowl. Whisk to combine, then " +
							"add the pumpkin puree, eggs, vanilla, and orange zest and combine to whisk until all the eggs have " +
							"been incorporated into the wet batter. Don't be alarmed if the batter splits or curdles! It's totally fine!",
						"BREAD BATTER: Add the dry ingredients into the wet ingredients in two batches, stirring just long enough " +
							"so each batch of flour is incorporated. Do not over-mix or you’ll end up with dry bread!",
						"BAKE: Divide the batter into the to pans, taking care to only fill each pan about ¾ of the way full. The " +
							"bread will rise significantly! Smooth out the batter then sprinkle with the pepitas. Bake the bread for " +
							"52-62 minutes or until a toothpick inserted in the center of the loaf comes out clean. Cool the pans on a " +
							"wire baking rack for at least 10 minutes before removing from the pan and allowing the bread to cool further.",
					},
				},
				Name:     "The Best Starbucks Pumpkin Loaf Recipe (Copycat)",
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://littlespicejar.com/starbucks-pumpkin-loaf/",
			},
		},
		{
			name: "livelytable.com",
			in:   "https://livelytable.com/bbq-ribs-on-the-charcoal-grill/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "main dish"},
				CookTime:      "PT2H30M",
				CookingMethod: models.CookingMethod{Value: "grilled"},
				Cuisine:       models.Cuisine{Value: "BBQ"},
				DatePublished: "2019-07-25",
				Description: models.Description{
					Value: "Nothing says summer like grilled BBQ ribs! These baby back ribs on the charcoal grill are " +
						"simple, delicious, and sure to please a crowd! (gluten-free, dairy-free, nut-free)",
				},
				Keywords: models.Keywords{Values: "BBQ ribs, ribs on the charcoal grill"},
				Image: models.Image{
					Value: "https://livelytable.com/wp-content/uploads/2019/07/ribs-on-charcoal-grill-2-225x225.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 rack baby back pork ribs",
						"1/3 cup BBQ spice rub",
						"water",
						"BBQ sauce of choice (optional)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Prepare fire in the charcoal grill. Remove the grates, place a pile of charcoal on one side of the grill " +
							"only. On the other side, place a small foil pan filled with water. Start the fire and return " +
							"the grates to the grill. Let the grill get to a low temperature (about 275°F.) You may also add " +
							"pieces of wood to the charcoal for a more smoky flavor.",
						"While the fire is heating, prepare ribs. Turn ribs over so that the bone side is facing up. Remove the " +
							"membrane along the back by sliding a dull knife (such as a butter knife) under the membrane " +
							"along the last bone until you get under the membrane. Hold on tight, and pull it until the whole " +
							"thing is removed from the rack of ribs.",
						"Rub ribs all over with spice rub. Once fire is ready, place the ribs on indirect heat - the side of the " +
							"grill that has the foil pan. Cover and cook about 2 hours, watching to make sure the fire is " +
							"maintained at a steady low temperature, adding charcoal as needed, and rotating the rack of ribs " +
							"roughly every 30 minutes so that different edges of the rack are turned toward the hot side.",
						"After 1 1/2 to 2 hours, remove ribs and wrap in foil. Return to the grill for another 30 minutes or so.",
						"When ribs are done, you can either remove them from the foil and place back on the grill, meat side down, " +
							"for a little char, or place them meat side up and brush with barbecue sauce in layers, waiting " +
							"about 5 minutes between layers. Or simply remove them from the grill to a cutting board, slice, and serve!",
					},
				},
				Name: "BBQ Ribs on the Charcoal Grill",
				NutritionSchema: models.NutritionSchema{
					Calories:      "416 calories",
					Carbohydrates: "8.9 g",
					Cholesterol:   "122.9 mg",
					Fat:           "26.3 g",
					Fiber:         "0.8 g",
					Protein:       "36.1 g",
					SaturatedFat:  "9.1 g",
					Servings:      "2",
					Sodium:        "512.8 mg",
					Sugar:         "5.1 g",
					TransFat:      "0.2 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 1},
				URL:      "https://livelytable.com/bbq-ribs-on-the-charcoal-grill/",
			},
		},
		{
			name: "lovingitvegan.com",
			in:   "https://lovingitvegan.com/vegan-buffalo-chicken-dip/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizer"},
				CookTime:      "PT20M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-01-21T14:31:28+00:00",
				Description: models.Description{
					Value: "This baked vegan buffalo chicken dip is rich, creamy and so cheesy. It&#39;s packed with spicy " +
						"flavor and makes the perfect crowd pleasing party dip.",
				},
				Keywords: models.Keywords{
					Values: "vegan buffalo chicken dip, vegan buffalo dip",
				},
				Image: models.Image{
					Value: "https://lovingitvegan.com/wp-content/uploads/2022/01/Vegan-Buffalo-Chicken-Dip-Square.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 1/2 cups Raw Cashews ((225g) Soaked in hot water for 1 hour)",
						"2 Tablespoons Lemon Juice (Freshly Squeezed)",
						"1/2 cup Canned Coconut Cream ((120ml) Unsweetened)",
						"1 teaspoon Distilled White Vinegar",
						"1 teaspoon Salt",
						"1 teaspoon Onion Powder",
						"1 teaspoon Vegan Chicken Spice (or Vegan Poultry Seasoning)",
						"1/2 cup Vegan Buffalo Sauce ((120ml))",
						"3/4 cup Nutritional Yeast ((45g))",
						"14 ounce Can Artichoke Hearts (in Brine or Water, (1 can) Drained and sliced into quarters)",
						"1/3 cup Spring Onions (Chopped)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Soak the cashews. Place the cashews into a bowl. Pour boiling hot water from the kettle over the top of " +
							"the cashews to submerge them. Leave the cashews to soak for 1 hour and then drain and rinse.",
						"Preheat the oven to 375°F (190°C).",
						"Add the soaked cashews, lemon juice, coconut cream, distilled white vinegar, salt, onion powder, vegan " +
							"chicken spice, vegan buffalo sauce and nutritional yeast to the blender and blend until smooth.",
						"Transfer the blended mix to a mixing bowl.",
						"Add chopped artichoke hearts and chopped spring onions and gently fold them in.",
						"Transfer to an oven safe 9-inch round dish and smooth down.",
						"Bake for 20 minutes until lightly browned on top.",
						"Serve topped with chopped spring onions with tortilla chips, crackers, breads or veggies for dipping.",
					},
				},
				Name: "Vegan Buffalo Chicken Dip",
				NutritionSchema: models.NutritionSchema{
					Calories:       "214 kcal",
					Carbohydrates:  "13 g",
					Fat:            "16 g",
					Fiber:          "3 g",
					Protein:        "8 g",
					SaturatedFat:   "7 g",
					Servings:       "1",
					Sodium:         "938 mg",
					Sugar:          "2 g",
					UnsaturatedFat: "8 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://lovingitvegan.com/vegan-buffalo-chicken-dip/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
