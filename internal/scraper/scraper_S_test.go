package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_S(t *testing.T) {
	testcases := []testcase{
		{
			name: "saboresajinomoto.com.br",
			in:   "https://www.saboresajinomoto.com.br/receita/pizza-de-pao-amanhecido",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT15M",
				DatePublished: "2023/11/17",
				Image: models.Image{
					Value: "https://www.saboresajinomoto.com.br/uploads/images/recipes/pizza-de-pao-amanhecido.webp",
				},
				Ingredients: models.Ingredients{Values: []string{
					"3 pães tipo francês amanhecidos e picados",
					"meia xícara (chá) de água (100 ml)",
					"1 sachê de Tempero SAZÓN® Laranja",
					"1 ovo",
					"3 colheres (sopa) de queijo tipo parmesão ralado",
					"meia xícara (chá) de polpa de tomate (100 ml)",
					"7 fatias de muçarela",
					"1 tomate pequeno cortado em rodelas finas",
					"6 azeitonas verdes sem caroço",
				}},
				Instructions: models.Instructions{Values: []string{
					"Em uma tigela grande, coloque o pão, a água e o Tempero SAZÓN®, e, com o auxílio de um garfo, amasse até o pão desmanchar completamente e formar uma massa homogênea. Acrescente o ovo e o queijo ralado, e misture.",
					"Espalhe a massa em uma assadeira redonda (30 cm de diâmetro), untada e enfarinhada, e leve ao forno médio (180 graus), preaquecido, por 15 minutos, ou até estar assada e desgrudar do fundo da assadeira.",
					"Retire do forno, espalhe a polpa de tomate pela superfície, por cima a muçarela, as rodelas de tomate e as azeitonas, e volte ao forno médio (180 graus), por mais 10 minutos, ou até o queijo derreter completamente. Retire do forno e sirva em seguida.",
				}},
				Name:     "PIZZA DE PÃO AMANHECIDO",
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://www.saboresajinomoto.com.br/receita/pizza-de-pao-amanhecido",
			},
		},
		{
			name: "sallysbakingaddiction.com",
			in:   "https://sallysbakingaddiction.com/breakfast-pastries/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Breakfast"},
				CookTime:      "PT20M",
				CookingMethod: models.CookingMethod{Value: "Baking"},
				Cuisine:       models.Cuisine{Value: "Danish"},
				DatePublished: "2020-08-01",
				Description: models.Description{
					Value: "These homemade breakfast pastries use a variation of classic Danish pastry dough. We're working the " +
						"butter directly into the dough, which is a different method from laminating it with separate layers " +
						"of butter. Make sure the butter is very cold before beginning. This recipe yields 2 pounds of dough.",
				},
				Keywords: models.Keywords{Values: "breakfast pastries, danishes, pastry"},
				Image: models.Image{
					Value: "https://sallysbakingaddiction.com/wp-content/uploads/2020/06/breakfast-pastries-2-225x225.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1/4 cup (60ml) warm water (between 100-110°F, 38-43°C)",
						"2 and 1/4 teaspoons Platinum Yeast from Red Star (1 standard packet)*",
						"1/4 cup (50g) granulated sugar",
						"1/2 cup (120ml) whole milk, at room temperature (between 68–72°F, 20-22°C)",
						"1 large egg, at room temperature",
						"1 teaspoon salt",
						"14 Tablespoons (196g) unsalted butter, cold",
						"2 and 1/2 cups (313g) all-purpose flour (spooned &amp; leveled), plus more for generously flouring hands, " +
							"surface, and dough",
						"2/3 cup filling (see recipe notes for options &amp; cheese filling)",
						"1 large egg",
						"2 Tablespoons (30ml) milk",
						"1 cup (120g) confectioners’ sugar",
						"2 Tablespoons (30ml) milk or heavy cream",
						"1 teaspoon pure vanilla extract",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"To help guarantee success, I recommend reading through the recipe, watching the video tutorial, and " +
							"reading the explanations below this recipe. (All answer many FAQs.) Do not use an electric mixer " +
							"for this dough. It&#8217;s best if the dough is folded together with a wooden spoon or rubber spatula " +
							"since it is so sticky. There is very minimal mixing required.",
						"Whisk the warm water, yeast, and 1 Tablespoon (6g) of sugar together in a large bowl. Cover and allow to " +
							"rest until foamy on top, about 5 minutes. If the surface doesn&#8217;t have bubbles on top or look " +
							"foamy after 15 minutes (it should if the yeast isn&#8217;t expired), start over with a fresh packet of" +
							" yeast. Whisk in remaining sugar, the milk, egg, and salt. Once these wet ingredients are mixed together, " +
							"lightly cover and set the bowl aside as you work on the next step.",
						"Cut the cold butter into 1/4 inch slices and add to a food processor or blender. Top with 2 and 1/2 cups " +
							"flour. Pulse the mixture 12-15 times, until butter is crumbled into pea-size bits. See photo below " +
							"for a visual. Using a food processor or blender is best for this dough. Keeping that in mind, if you " +
							"don&#8217;t have one, you can use a pastry cutter to work in the butter.",
						"Pour the flour mixture into the wet yeast mixture. Very gently fold everything together using a rubber " +
							"spatula or wooden spoon. Fold *just until* the dry ingredients are moistened. The butter must remain " +
							"in pieces and crumbles, which creates a flaky pastry. Turn the sticky dough out onto a large piece of " +
							"plastic wrap, parchment paper, aluminum foil, or into any container you can tightly cover.",
						"Wrap the dough/cover up tightly and refrigerate for at least 4 hours and up to 48 hours.",
						"Take the dough out of the refrigerator to begin the “rolling and folding” process. If the dough sat for " +
							"more than 4 hours, it may have slightly puffed up and that&#8217;s ok. (It will deflate as you shape " +
							"it, which is also ok.) Very generously flour a work surface. The dough is very sticky, so make sure you " +
							"have more flour nearby as you roll and fold. Using the palm of your hands, gently flatten the dough into " +
							"a small square. Using a rolling pin, roll out into a 15&#215;8 inch rectangle. When needed, flour " +
							"the work surface and dough as you are rolling. Fold the dough into thirds as if it were a business " +
							"letter. (See photos and video tutorial.) Turn it clockwise and roll it out into a 15 inch long rectangle " +
							"again. Then, fold into thirds again. Turn it clockwise. You’ll repeat rolling and folding 1 more time for " +
							"a total of 3 times.",
						"Wrap up/seal tightly and refrigerate for at least 1 hour and up to 24 hours. You can also freeze the dough " +
							"at this point. See freezing instructions.",
						"Line two large baking sheets with parchment paper or silicone baking mats. Rimmed baking sheets are best " +
							"because butter may leak from the dough as it bakes. If you don&#8217;t have rimmed baking sheets, when " +
							"it&#8217;s time to preheat the oven, place another baking sheet on the oven rack below to catch any butter " +
							"that may drip.",
						"Take the dough out of the refrigerator and cut it in half. Wrap 1 half up and keep refrigerated as you " +
							"work with the first half. (You can freeze half of the dough at this point, use the freezing instructions " +
							"below.)",
						"Cut the first half of dough into 8 even pieces. This will be about 1/4 cup of dough per pastry. Roll each " +
							"into balls. Flatten each into a 2.5 inch circle. Use your fingers to create a lip around the edges. See " +
							"photos and video tutorial if needed. Press the center down to flatten the center as much as you can so you " +
							"can fit the filling inside. (Center puffs up as it bakes.) Arrange pastries 3 inches apart on a lined " +
							"baking sheet. Repeat with second half of dough.",
						"Spoon 2 teaspoons of fruity filling or 1 Tablespoon of cheese filling inside each.",
						"Whisk the egg wash ingredients together. Brush on the edges of each shaped pastry.",
						"This step is optional, but I very strongly recommend it. Chill the shaped pastries in the refrigerator, " +
							"covered or uncovered, for at least 15 minutes and up to 1 hour. See recipe note. You can preheat the " +
							"oven as they finish up chilling.",
						"Preheat oven to 400°F (204°C).",
						"Bake for 19-22 minutes or until golden brown around the edges. Some butter may leak from the dough, " +
							"that&#8217;s completely normal and expected. Feel free to remove the baking sheets from the oven halfway " +
							"through baking and brush the dough with any of the leaking butter, then place back in the oven to finish " +
							"baking. (That&#8217;s what I do!)",
						"Remove baked pastries from the oven. Cool for at least 5 minutes before icing/serving.",
						"Whisk the icing ingredients together. If you want a thicker icing, whisk in more confectioners’ sugar. " +
							"If you want a thinner icing, whisk in more milk or cream. Drizzle over warm pastries and serve.",
						"Cover leftover iced or un-iced pastries and store at room temperature for 1 day or in the refrigerator " +
							"for up to 5 days. Or you can freeze them for up to 3 months. Thaw before serving. Before enjoying, feel " +
							"free to reheat leftover iced or un-iced pastries in the microwave for a few seconds until warmed.",
					},
				},
				Name:     "Breakfast Pastries with Shortcut Homemade Dough",
				PrepTime: "PT6H",
				Yield:    models.Yield{Value: 16},
				URL:      "https://sallysbakingaddiction.com/breakfast-pastries/",
			},
		},
		{
			name: "sallys-blog.de",
			in:   "https://sallys-blog.de/rezepte/zwieback-dessert-etimek-tatlisi-no-bake",
			want: models.RecipeSchema{
				AtContext: "https://schema.org",
				AtType:    models.SchemaType{Value: "Recipe"},
				Description: models.Description{
					Value: "Dieses hübsche Dessert sieht nicht nur gut aus, sondern ist auch ganz schnell zubereitet, dann muss es nur noch gekühlt werden und ist bereit zum S",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"200 g Zucker",
						"200 g Wasser",
						"225 g Zwieback",
						"1 l Milch",
						"80 g Zucker",
						"2 TL Vanilleextrakt",
						"90 g Speisestärke",
						"100 g Kokosraspeln",
						"1 TL Zimt",
						"400 g Sahne",
						"1 TL Vanilleextrakt",
						"4 TL San-apart",
						"30 g Pistazien",
						"100 g Himbeeren (frisch oder TK)",
						"100 g Heidelbeeren (frisch oder TK)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Gib den Zucker in einen Topf, erhitze ihn, bis er leicht karamellisiert und gieße ihn mit dem Wasser auf. Koche den Sirup kurz auf, bis sich das Karamell gelöst hat und er ganz leicht dickflüssig wird und nimm ihn dann vom Herd. Verteile den Zwieback in einer Ofenform oder einem Backrahmen und gieße den warmen Sirup darüber, wende dabei die Zweibackscheiben, so kann der Sirup von allen Seiten einziehen.",
						"Verrühre die Milch mit dem Zucker, Vanilleextrakt und der Stärke in einem Topf und lasse den Pudding aufkochen. Koche ihn für etwa 2 Minuten und gieße ihn dann ebenfalls über den Zwieback. Lasse die Masse abkühlen.",
						"Tipp: Streiche den Pudding durch ein Haarsieb, um mögliche Klümpchen zu entfernen.",
						"Röste die Kokosraspeln in einer Pfanne ohne Fett an, bis sie leicht Farbe annehmen, verrühre sie mit dem Zimt und lasse sie auf einem Teller abkühlen. Schlage die Sahne mit Vanilleextrakt und Sanapart steif und streiche sie auf dem Pudding glatt. Hacke die Pistazien klein und streue sie gemeinsam mit den Kokosraspeln über die Sahne. Stelle das Dessert für mindestens 4 Stunden oder über Nacht in den Kühlschrank. Dekoriere das Dessert zum Servieren mit den frischen Beeren. Viel Spaß beim Nachmachen, eure Sally!",
					},
				},
				Name:     "zwieback dessert / etimek tatlisi / no bake ramadan rezept",
				PrepTime: "30",
				URL:      "https://sallys-blog.de/rezepte/zwieback-dessert-etimek-tatlisi-no-bake",
			},
		},
		{
			name: "saltpepperskillet.com",
			in:   "https://saltpepperskillet.com/recipes/creamy-mashed-potatoes/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Side Dish"},
				CookTime:      "PT40M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2021-11-16T16:12:35+00:00",
				Description: models.Description{
					Value: "The creamiest, most luxurious and delicious mashed potatoes. A beloved side dish that can become the star of the meal.",
				},
				Keywords: models.Keywords{Values: "mashed potatoes"},
				Image: models.Image{
					Value: "https://saltpepperskillet.com/wp-content/uploads/creamy-mashed-potatoes-horizontal.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 lbs Yukon Gold potatoes", "1/4 lb unsalted butter ((room temperature))",
						"1 cup heavy cream ((hot))", "Diamond kosher salt", "Freshly ground pepper",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Place the whole unpeeled potatoes in a 4-quart saucepan and cover with a few inches of cold water.",
						"Bring to a simmer and cook until they are very tender when pierced with a paring knife, about 25 to 35 minutes depending on the size of the potatoes.",
						"Strain the cooked potatoes in a colander and carefully peel with a clean kitchen towel and a paring knife. Work quickly so the potatoes stay warm.",
						"Push the potatoes through a potato ricer or food mill back into the warm pot they were cooked in.",
						"Add the room temperature butter and combine using a stiff spatula. Next start adding 3/4 cup of the hot cream and gently fold in. Add more as needed.",
						"Season well with kosher salt and freshly ground pepper. Taste and add more butter, cream or seasoning as needed.",
						"Cover and keep warm until serving with another pat of butter and fresh herbs if desired.",
					},
				},
				Name:     "Mashed Potatoes",
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://saltpepperskillet.com/recipes/creamy-mashed-potatoes/",
			},
		},
		{
			name: "saveur.com",
			in:   "https://www.saveur.com/recipes/varenyky-pierogi-recipe/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "Larisa Frumkin’s Varenyky",
				CookTime:  "PT0D1H30M",
				Description: models.Description{
					Value: "These sweet dumplings, known as pierogi in Poland and varenyky in Ukraine, are a staple of many Slavic cuisines.",
				},
				DatePublished: "2022-04-05 17:23:56",
				Image: models.Image{
					Value: "https://www.saveur.com/uploads/2022/04/HR-Pierogi-Saveur-08-scaled.jpg?auto=webp",
				},
				Yield: models.Yield{Value: 4},
				URL:   "https://www.saveur.com/recipes/varenyky-pierogi-recipe/",
				Ingredients: models.Ingredients{
					Values: []string{
						"2 cups all-purpose flour, plus more for dusting",
						"1 tsp. kosher salt",
						"2 large eggs, separated",
						"1 tbsp. vegetable oil",
						"2 cups (1 lb.) farmer’s cheese",
						"1 large egg yolk",
						"3 tbsp. sugar (or substitute salt to taste for a savory version)",
						"4 tbsp. softened unsalted butter",
						"Sour cherry preserves, sour cream, or crème fraîche, to serve",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Make the dough: To a food processor, add the flour and the salt. With the motor running, add the 2 egg yolks, one-by-one, then drizzle in the oil through the feeder tube. With the motor still running, drizzle in 8–10 tablespoons of cool water, just until the dough begins to form a ball around the blade. Lightly flour a clean work surface, then transfer the dough out onto it and knead just until smooth, about 2 minutes. Cover with a clean kitchen towel and set aside to rest for 30 minutes.",
						"Meanwhile, make the filling: In a medium bowl, mix together the farmer’s cheese, egg yolk, and sugar. Set aside.",
						"Lightly flour a large rimmed baking sheet and set it by your work surface.",
						"Begin shaping the varenyky. Dust your work surface lightly with flour; divide the dough in half and shape into 2 balls. Keep one ball covered with the kitchen towel and, using a lightly floured rolling pin or a hand-crank pasta roller, roll the other ball into a thin sheet, about 1⁄16 -inch thick. Using a 3-inch round cookie cutter, punch out circles of the dough. Place a heaping teaspoon of filling in the center of each circle. Brush the edges of the circles lightly with egg white, then fold into a half moons, pressing the edges firmly together with either your fingers or with the tines of a fork to seal. Place the varenyky on the baking sheet about ½-inch apart and cover with a damp cloth. Roll out the second ball of dough, and repeat, then combine all of the leftover dough scraps to make a third batch.",
						"Fill a large pot two thirds of the way with water and salt generously. Set over medium high heat and bring to a boil. Carefully lower half the varenyky into the pot. Boil, stirring occasionally to prevent sticking, until the dumplings rise to the surface and the dough is cooked through, 6–7 minutes. Using a slotted spoon, transfer the varenyky to a deep bowl and add the butter, tossing gently with the spoon to melt. Keep warm while you cook the remaining dumplings. Divide the varenyky among 4 deep plates, top with sour cherry preserves and sour cream, and serve warm.",
					},
				},
				PrepTime: "PT0D0H0M",
			},
		},
		{
			name: "seriouseats.com",
			in:   "https://www.seriouseats.com/miyeok-guk-korean-seaweed-and-brisket-soup",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Soups and Stews"},
				CookTime:      "PT180M",
				Cuisine:       models.Cuisine{Value: "Korean"},
				DateModified:  "2023-07-12T11:51:57.502-04:00",
				DatePublished: "2020-03-02T08:00:03.000-05:00",
				Description: models.Description{
					Value: "Tender seaweed and pieces of beef brisket come together in this warming, comforting, and nutritious " +
						"Korean soup.",
				},
				Image: models.Image{
					Value: "https://www.seriouseats.com/thmb/BkhWm33gH4ho1SSMFOLwfO8O6Ww=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__serious_eats__seriouseats.com__2020__02__20200128-miyeok-guk-korean-seaweed-soup-vicky-wasik-7-21447f5620914e4b9e19912a78b7306c.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 ounce (30g) dried miyeok seaweed (also sold under the Japanese name wakame )",
						"3 whole medium cloves garlic plus 3 finely minced medium cloves garlic, divided",
						"One 1-inch piece fresh ginger (about 1/3 ounce; 10g), peeled",
						"1/2 of a medium white onion (about 3 ounces; 85g for the half onion)",
						"12 ounces (350g) beef brisket, washed in cold water",
						"2 tablespoons (30ml) Joseon ganjang (Korean soup soy sauce; see note), divided",
						"Kosher or sea salt",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a medium bowl, cover seaweed with at least 3 inches cold water and let stand at room temperature until " +
							"fully softened and hydrated, about 2 hours.",
						"Meanwhile, in a Dutch oven or pot, combine whole garlic cloves, ginger, onion, and brisket with 1 1/2 " +
							"quarts (1 1/2L) cold water and bring to a boil over high heat. Lower heat to maintain a gentle simmer " +
							"and cook, covered, until brisket is tender and broth is slightly cloudy, about 2 hours. Using a slotted " +
							"spoon, remove and discard garlic cloves, ginger, and onion from broth.",
						"Transfer brisket to a work surface and allow to cool slightly, then slice across the grain into bite-size " +
							"pieces. Transfer brisket to a small bowl and toss well with 1 tablespoon soy sauce and remaining " +
							"3 cloves minced garlic. Set aside.",
						"Drain seaweed and squeeze well to remove excess water. Transfer to work surface and roughly chop into " +
							"bite-size pieces.",
						"Return broth to a simmer and add seaweed and seasoned brisket. If the proportion of liquid to solids is " +
							"too low for your taste, you can top up with water and return to a simmer. Add remaining 1 tablespoon " +
							"soy sauce and simmer until seaweed is tender, about 30 minutes. Season to taste with salt.",
						"Ladle soup into bowls and serve alongside hot rice and any banchan (side dishes) of your choosing.",
					},
				},
				Name: "Miyeok-Guk (Korean Seaweed and Brisket Soup)",
				NutritionSchema: models.NutritionSchema{
					Calories:       "173 kcal",
					Carbohydrates:  "2 g",
					Cholesterol:    "60 mg",
					Fat:            "11 g",
					Fiber:          "0 g",
					Protein:        "17 g",
					SaturatedFat:   "4 g",
					Servings:       "4",
					Sodium:         "421 mg",
					Sugar:          "0 g",
					UnsaturatedFat: "0 g",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://www.seriouseats.com/miyeok-guk-korean-seaweed-and-brisket-soup",
			},
		},
		{
			name: "simple-veganista.com",
			in:   "https://simple-veganista.com/blackberry-cobbler/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dessert"},
				CookTime:      "PT40M",
				CookingMethod: models.CookingMethod{Value: "bake"},
				Cuisine:       models.Cuisine{Value: "Southern"},
				DatePublished: "2021-05-14",
				Description: models.Description{
					Value: "With only 7 ingredients this flavorful vegan blackberry cobbler is a great way to use up the season's abundance of blackberries.",
				},
				Image: models.Image{
					Value: "https://simple-veganista.com/wp-content/uploads/2016/06/best-vegan-blackberry-cobbler-9-225x225.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 cup flour", "1/2 cup sugar", "1 heaped teaspoon baking powder",
						"1/4 teaspoon cinnamon",
						"1 cup unsweetened vanilla plant milk, warmed if using coconut oil",
						"1/3 cup coconut oil or vegan butter/margarine, warmed to liquid state",
						"2 &#8211; 3 cups (about 12oz.) blackberries, fresh or frozen",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 350 degrees F.",
						"Batter: In a medium bowl, whisk together the flour, sugar, baking powder, and cinnamon. Be sure to warm your milk if using coconut oil before using or it will harden the coconut oil when combined. Add milk and mix well. Add oil/butter, mix well again.",
						"Assemble: Pour batter into an 8-inch baking dish. Drop blackberries into the batter, distributing evenly all over. Push down some of the blackberries into the batter, then add more over top. Sprinkle a little pure cane sugar over top if you like. I would refrain from using coconut sugar as it may burn. If you want to add coconut sugar, sprinkle some at least halfway through baking time.",
						"Bake: Place in oven on the middle rack and bake for 40 &#8211; 45 minutes. Let cool a few minutes and serve.",
						"Serve: Pair with a scoop of non-dairy vanilla ice cream for dessert and a scoop of vanilla non-dairy yogurt for breakfast (yes, this would be a fine way to start your day!)",
						"Serves 6",
						"Leftovers can be stored in the refrigerator for up to 5 days. Warm serving portions in the microwave or oven set to 350 for 10 &#8211; 15 minutes.",
					},
				},
				Name: "Southern Blackberry Cobbler (Vegan + Easy)",
				NutritionSchema: models.NutritionSchema{
					Calories:       "245 calories",
					Carbohydrates:  "28.3 g",
					Cholesterol:    "0 mg",
					Fat:            "13.5 g",
					Fiber:          "7 g",
					Protein:        "5.5 g",
					SaturatedFat:   "10.1 g",
					Sodium:         "215.6 mg",
					Sugar:          "5.5 g",
					TransFat:       "0 g",
					UnsaturatedFat: "",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://simple-veganista.com/blackberry-cobbler/",
			},
		},
		{
			name: "simply-cookit.com",
			in:   "https://www.simply-cookit.com/de/rezepte/paprikagulasch",
			want: models.RecipeSchema{
				AtContext: "https://schema.org",
				AtType:    models.SchemaType{Value: "Recipe"},
				Category:  models.Category{Value: "Fleisch & Geflügel"},
				CookTime:  "PT1H59M00S",
				Description: models.Description{
					Value: "Deftiges, würziges Gulasch mit bunter Paprika - so wird die Sauce schön fruchtig.",
				},
				Keywords: models.Keywords{Values: "Fleisch & Geflügel"},
				Image: models.Image{
					Value: "https://www.simply-cookit.com/sites/default/files/styles/square/public/assets/image/2021/03/buntes-paprikagulasch_portrait.jpg?h=2e35897f&itok=MwyrxSek",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"300 g Zwiebeln", "20 ml Rapsöl", "750 g Gulasch, gemischt",
						"20 g Tomatenmark", "1 EL Paprikapulver, edelsüß",
						"1 EL Paprikapulver, rosenscharf", "3 Stängel Thymian, frisch",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Das Universalmesser einsetzen. Die Zwiebeln schälen, halbieren und einwiegen. Den Deckel schließen, den Messbecher einsetzen und (Universalmesser | Stufe 14 | 10 Sek.) zerkleinern. Das Universalmesser entnehmen und das Lebensmittel mit dem Spatel abstreifen. Die Zwiebeln umfüllen.",
						"Den 3D-Rührer einsetzen. Das Rapsöl einwiegen, den Deckel schließen, den Messbecher einsetzen und das Öl (3D-Rührer | 200 \xc2\xb0C | 4 Min.) erhitzen. Das Gulasch einwiegen. Den Deckel schlie\xc3\x9fen, den Messbecher entnehmen und (3D-Rührer | Stufe 2 | 200 °C | 9 Min.) scharf anbraten. Die Zwiebeln zugeben, den Deckel schließen und die Zwiebeln (3D-Rührer | Stufe 2 | 160 °C | 2 Min.) kurz mitbraten.",
					},
				},
				Name: "Buntes Paprikagulasch",
				NutritionSchema: models.NutritionSchema{
					Calories:      "577 kcal",
					Carbohydrates: "16 g",
					Fat:           "36 g",
					Protein:       "38 g",
				},
				Yield: models.Yield{Value: 4},
				URL:   "https://www.simply-cookit.com/de/rezepte/paprikagulasch",
			},
		},
		{
			name: "simplyquinoa.com",
			in:   "https://www.simplyquinoa.com/spicy-kimchi-quinoa-bowls/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Course"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "korean"},
				DatePublished: "2021-01-15T07:00:21+00:00",
				Description: models.Description{
					Value: "These spicy kimchi quinoa bowls are the perfect weeknight dinner. They&#039;re quick, easy, and " +
						"super healthy, packed with protein, fermented veggies, and greens!",
				},
				Keywords: models.Keywords{Values: "egg, kimchi, quinoa bowl"},
				Image: models.Image{
					Value: "https://www.simplyquinoa.com/wp-content/uploads/2015/06/spicy-kimchi-quinoa-bowls-3.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 teaspoons toasted sesame oil",
						"1/2 teaspoon freshly grated ginger",
						"1 teaspoon minced garlic",
						"2 cups cooked quinoa (cooled)",
						"1 cup kimchi (chopped)",
						"2 teaspoons kimchi \"juice\" (the liquid from the jar)",
						"2 teaspoons gluten-free tamari",
						"1 teaspoon hot sauce (optional)",
						"2 cups kale (finely chopped)",
						"2 eggs",
						"1/4 cup sliced green onions for garnish (optional)",
						"Fresh ground pepper for garnish (optional)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Heat the oil in a large skillet over medium heat. Add ginger and garlic and saute for 30 - 60 seconds " +
							"until fragrant. Add the quinoa and kimchi and cook until hot, about 2 - 3 minutes. Stir in kimchi " +
							"juice, tamari and hot sauce if using. Turn to low and stir occasionally while you prepare the other " +
							"ingredients.",
						"In a separate skillet, cook the eggs on low until the whites have cooked through but the yolks are " +
							"still runny, about 3 - 5 minutes.",
						"Steam the kale in a separate pot for 30 - 60 seconds until soft.",
						"Assemble the bowls, dividing the kimchi-quinoa mixture and kale evenly between two dishes. Top with green " +
							"onions and fresh pepper if using.",
					},
				},
				Name: "Spicy Kimchi Quinoa Bowls",
				NutritionSchema: models.NutritionSchema{
					Calories:      "359 kcal",
					Carbohydrates: "46 g",
					Cholesterol:   "163 mg",
					Fat:           "12 g",
					Fiber:         "5 g",
					Protein:       "17 g",
					SaturatedFat:  "2 g",
					Servings:      "1",
					Sodium:        "489 mg",
					Sugar:         "1 g",
				},
				PrepTime: "PT3M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://www.simplyquinoa.com/spicy-kimchi-quinoa-bowls/",
			},
		},
		{
			name: "simplyrecipes.com",
			in:   "https://www.simplyrecipes.com/recipes/chicken_tikka_masala/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				CookTime:      "PT20M",
				Cuisine:       models.Cuisine{Value: "British"},
				DateModified:  "2023-09-29T18:15:43.573-04:00",
				DatePublished: "2017-02-27T04:30:56.000-05:00",
				Description: models.Description{
					Value: "This easy stovetop Chicken Tikka Masala tastes just like your favorite Indian take-out and is ready " +
						"in under an hour. Leftovers are even better the next day!",
				},
				Keywords: models.Keywords{
					Values: "Comfort Food, Quick and Easy, Restaurant Favorite, British, Indian, Gluten-Free, Dinner",
				},
				Image: models.Image{
					Value: "https://www.simplyrecipes.com/thmb/pYiHJojfyPYHFzhTQS8OU0GXUlE=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/__opt__aboutcom__coeus__resources__content_migration__simply_recipes__uploads__2017__02__2017-02-27-ChickenTikkaMasala-18-2b30d704a54e4620a0f17fd085afeef5.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"For the chicken:",
						"1 1/4 pounds boneless skinless chicken breasts, thighs, or a mix",
						"6 tablespoons plain whole milk yogurt",
						"1/2 tablespoon grated ginger",
						"3 cloves of garlic, minced",
						"1 teaspoon cumin",
						"1 teaspoon paprika",
						"1 1/4 teaspoons salt",
						"For the tikka masala sauce:",
						"2 tablespoons canola oil, divided",
						"1 small onion, thinly sliced (about 5 ounces, or 1 1/2 cups sliced)",
						"2 teaspoons grated ginger",
						"4 cloves garlic, minced",
						"1 tablespoon ground coriander",
						"2 teaspoons paprika",
						"1 teaspoon garam masala",
						"1/2 teaspoon turmeric",
						"1/2 teaspoon freshly ground black pepper",
						"1 (14-ounce) can crushed fire-roasted tomatoes (regular crushed tomatoes work, too)",
						"6 tablespoons plain whole milk yogurt",
						"1/4 to 1/2 teaspoon cayenne pepper",
						"1/2 teaspoon salt",
						"Cooked rice , to serve",
						"Cilantro, to garnish",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Prepare the chicken: Trim chicken thighs of any extra fat. Chop into bite-sized pieces about 1/2 to 1-inch " +
							"wide. Place the chicken thigh pieces to a medium bowl. Add the yogurt, ginger, garlic, cumin, paprika " +
							"and salt. Using your hands, combine the chicken with the spices until the chicken is evenly coated.",
						"Marinate the chicken: Cover the bowl with plastic wrap and let the chicken marinate in the fridge for at " +
							"least 45 minutes or as long as overnight. (Marinating for 4 to 6 hours is perfect.)",
						"Cook the chicken: In a large skillet, heat 1 tablespoon of canola oil over medium-high heat. Add the chicken" +
							" thigh pieces and cook for about 6 to 7 minutes, until they’re cooked through. Transfer to a plate " +
							"and set aside.",
						"Toast the spices: Wipe down the pan you used to cook the chicken. Heat remaining canola oil over medium " +
							"heat. Add the onions and cook for 5 minutes, until softened, stirring often. Add the grated ginger," +
							" minced garlic, coriander, paprika, garam masala, turmeric, black pepper, salt, and cayenne. Let the " +
							"spices cook until fragrant, about 30 seconds to a minute.",
						"Make the sauce: Add the crushed tomatoes to the pan with the spices and let everything cook for 4 minutes," +
							" stirring often. Add the yogurt and stir to combine.",
						"Simmer the sauce: Reduce the heat to medium-low and let the sauce simmer for another 4 minutes. Add the c" +
							"hicken pieces to the pan and coat with sauce.",
						"Serve: Serve over cooked basmati rice and garnish with cilantro.",
					},
				},
				Name: "Chicken Tikka Masala",
				NutritionSchema: models.NutritionSchema{
					Calories:       "324 kcal",
					Carbohydrates:  "25 g",
					Cholesterol:    "84 mg",
					Fat:            "10 g",
					Fiber:          "3 g",
					Protein:        "34 g",
					SaturatedFat:   "2 g",
					Servings:       "4",
					Sodium:         "828 mg",
					Sugar:          "5 g",
					UnsaturatedFat: "0 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://www.simplyrecipes.com/recipes/chicken_tikka_masala/",
			},
		},
		{
			name: "simplywhisked.com",
			in:   "https://www.simplywhisked.com/dill-pickle-pasta-salad/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Salads"},
				CookTime:      "PT10M",
				CookingMethod: models.CookingMethod{Value: "Stovetop"},
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-03-02",
				Description: models.Description{
					Value: "Looking for something new to bring to your next potluck? This super flavorful dill pickle pasta salad " +
						"is a crowd pleaser, and it's so easy to make. It's loaded with crunchy dill pickles, savory bacon, " +
						"toasted cashews and topped with a creamy dill dressing.",
				},
				Keywords: models.Keywords{
					Values: "dill pickle pasta salad, pasta salad with pickles, dill pasta salad, pasta salad recipe, dairy free " +
						"dill pickle pasta salad, dairy free pasta salad, dairy free macaroni salad",
				},
				Image: models.Image{
					Value: "https://www.simplywhisked.com/wp-content/uploads/2022/01/Dairy-Free-Dill-Pickle-Pasta-Salad-3-225x225.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 pound pasta, cooked and cooled",
						"4 slices bacon, thinly sliced",
						"3/4 cup cashews, chopped",
						"1 cup chopped dill pickles",
						"2 stalks celery, thinly sliced",
						"3 green onions, thinly sliced",
						"2 tablespoons fresh dill, chopped",
						"2 cups mayonnaise",
						"1 tablespoon Dijon mustard",
						"4 tablespoons pickle juice",
						"2 tablespoons water",
						"Coarse salt &amp; black pepper, to taste",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"If needed, cook pasta according to package directions for al dente in a large pot of salted, boiling " +
							"water. Drain pasta and rinse with cold water.",
						"In a medium mixing bowl, whisk the mayonnaise, Dijon mustard, pickle juice, and water until smooth. " +
							"Season salt &amp; pepper, to taste.",
						"Combine salad ingredients in a large mixing bowl, reserving about 1 teaspoon fresh dill. Add dressing " +
							"and stir until evenly coated.",
						"Before serving, adjust seasoning with salt &amp; pepper (to taste) and garnish with remaining dill.",
					},
				},
				Name: "Dill Pickle Pasta Salad",
				NutritionSchema: models.NutritionSchema{
					Calories:      "386 calories",
					Carbohydrates: "25.7 g",
					Cholesterol:   "16.2 mg",
					Fat:           "28.5 g",
					Fiber:         "1.6 g",
					Protein:       "7 g",
					SaturatedFat:  "5.1 g",
					Sodium:        "356.8 mg",
					Sugar:         "1.9 g",
					TransFat:      "0.1 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 16},
				URL:      "https://www.simplywhisked.com/dill-pickle-pasta-salad/",
			},
		},
		{
			name: "skinnytaste.com",
			in:   "https://www.skinnytaste.com/air-fryer-steak/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-03-29T09:04:25+00:00",
				Description: models.Description{
					Value: "Make perfect Air Fryer Steak that is seared on the outside and juicy on the inside. Air frying " +
						"steak is quick and easy with no splatter or mess in the kitchen!",
				},
				Keywords: models.Keywords{
					Values: "Air Fryer Recipes, air fryer steak, sirloin",
				},
				Image: models.Image{
					Value: "https://www.skinnytaste.com/wp-content/uploads/2022/03/Air-Fryer-Steak-6.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 teaspoon garlic powder",
						"1/2 teaspoon sweet paprika",
						"1 teaspoon kosher salt",
						"1/4 teaspoon black pepper",
						"4 sirloin steaks (1 inch thick (1 1/2 lbs total))",
						"olive oil spray",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Combine the spices in a small bowl.",
						"Spray the steak with olive oil and coat both sides with the spices.",
						"Preheat the air fryer so the basket gets hot. For a 1-inch steak, air fry 400F 10 minutes " +
							"turning halfway, for medium rare, for medium, cook 12 minutes, flipping halfway. " +
							"See temp chart below, time may vary slightly with different air fryer models, " +
							"and the thickness of the steaks.",
						"Finish with a pinch of more salt and black pepper.",
						"Let it rest, tented with foil 5 minutes before slicing.",
					},
				},
				Name: "Air Fryer Steak",
				NutritionSchema: models.NutritionSchema{
					Calories:      "221 kcal",
					Carbohydrates: "0.5 g",
					Cholesterol:   "117.5 mg",
					Fat:           "7 g",
					Fiber:         "0.5 g",
					Protein:       "39.5 g",
					SaturatedFat:  "2 g",
					Servings:      "1",
					Sodium:        "391 mg",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.skinnytaste.com/air-fryer-steak/",
			},
		},
		{
			name: "sobors.hu",
			in:   "https://sobors.hu/receptek/karamelles-sajttorta-poharkrem-recept/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT40M",
				DateModified:  "2023-11-24",
				DatePublished: "2023-11-23",
				Description: models.Description{
					Value: "Készíts végtelenül krémes desszertet, amelyet akár több nappal a vendégvárás előtt is összedobhatsz: hamis sajtkrém és mennyei karamellszósz alkotja!",
				},
				Keywords: models.Keywords{
					Values: "sós karamell,vegán,pohárkrém,tejmentes,narancsos mézeskalácsos pohárdesszert,karamellszósz,vegán desszert,sajttorta pohárdesszert,hamis karamell,karamelles sajttorta pohárdesszert",
				},
				Image: models.Image{
					Value: "https://kep.cdn.index.hu/1/0/5187/51872/518723/51872395_3953787_76294cf7da02a7d922b6c376100fbca3_wm.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"150 ml vegán karamell a leírás szerint",
						"250 ml habosítható növényi tejszín", "250 g vegán krémsajt",
						"1 tk vaníliakivonat",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Elősz&ouml;r k&eacute;sz&iacute;ts&uuml;k el a veg&aacute;n s&oacute;s karamellt az ebben a receptben le&iacute;rtak szerint. Haszn&aacute;ljunk fel belőle 150 ml-t, a marad&eacute;kb&oacute;l pedig k&eacute;sz&iacute;ts&uuml;nk karamelles latt&eacute;t, csurgassuk alm&aacute;s pite tetej&eacute;re, vagy fogyasszuk tetsz&eacute;s szerint. Ha elk&eacute;sz&uuml;lt, hagyjuk szobah\xc5\x91m&eacute;rs&eacute;kletűre hűlni.",
						"Egy t&aacute;lban habos&iacute;tsuk fel a n&ouml;v&eacute;nyi tejsz&iacute;nt, majd tegy&uuml;k f&eacute;lre. Tegy&uuml;k egy m&aacute;sik t&aacute;lba a veg&aacute;n kr&eacute;msajtot, &eacute;s a k&eacute;zi mixerrel ezt is habos&iacute;tsuk fel, majd adjuk hozz&aacute; a van&iacute;liakivonatot &eacute;s 100 ml-t a karamellből, &eacute;s dolgozzuk &ouml;ssze.",
						"Egy spatul&aacute;val &oacute;vatosan forgassuk bele a felhabos&iacute;tott tejsz&iacute;nt, majd osszuk a kr&eacute;met 4 egyforma poh&aacute;rba a k&ouml;vetkező m&oacute;don: kanalazzunk a poharak alj&aacute;ra 1-1 kan&aacute;llal a megmaradt karamellből, majd mehet r&aacute; a kr&eacute;m is.",
						"Fogyaszt&aacute;s előtt legal&aacute;bb 1 &oacute;r&aacute;ra tegy&uuml;k hűtőbe, hogy kiss&eacute; megdermedjen, &eacute;s az &iacute;zek &ouml;ssze&aacute;lljanak.",
					},
				},
				Name:     "Karamelles sajttortapohárkrém",
				PrepTime: "PT5M",
				URL:      "https://sobors.hu/receptek/karamelles-sajttorta-poharkrem-recept/",
				Yield:    models.Yield{Value: 1},
			},
		},
		{
			name: "southerncastiron.com",
			in:   "https://southerncastiron.com/creamy-turkey-and-wild-rice-soup/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Dishes"},
				DateModified:  "2023-11-22T20:45:27+00:00",
				DatePublished: "2023-11-22T20:45:27+00:00",
				Description: models.Description{
					Value: "The star of the Thanksgiving meal takes on a new life in this cozy and comforting Creamy Turkey and Wild Rice Soup.",
				},
				Image: models.Image{
					Value: "https://southerncastiron.com/wp-content/uploads/2017/06/sthCastIron-logo-544x180.png",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"4 tablespoons unsalted butter", "½ cup sliced celery",
						"2 large carrots, halved and sliced", "1 medium yellow onion, chopped",
						"2 cloves garlic, chopped", "2 teaspoons fresh thyme leaves",
						"¼ cup all-purpose flour", "1 cup wild rice blend",
						"2 (32-ounce) packages low-sodium chicken broth", "1 teaspoon kosher salt",
						"½ teaspoon ground black pepper", "¼ teaspoon poultry seasoning",
						"3 cups shredded roasted turkey",
						"Garnish: ground black pepper, fresh thyme leaves",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a medium cast-iron Dutch oven, melt butter over medium heat. Add celery, carrot, onion, garlic, and thyme; cook over medium heat, stirring occasionally, until vegetables are crisp-tender, about 10 minutes. Sprinkle flour over vegetables and cook, stirring, until evenly coated and lightly browned, about 3 minutes.",
						"Add rice and gradually stir in broth, salt, pepper, and poultry seasoning. Bring to a boil; reduce heat to low and simmer, stirring occasionally, until vegetables are tender, about 30 minutes.",
						"Add turkey; cook, stirring occasionally, until rice is tender, 10 to 15 minutes. Serve warm. Garnish with pepper and thyme, if desired.",
					},
				},
				Name:  "Creamy Turkey and Wild Rice Soup",
				Yield: models.Yield{Value: 6},
				URL:   "https://southerncastiron.com/creamy-turkey-and-wild-rice-soup/",
			},
		},
		{
			name: "southernliving.com",
			in:   "https://www.southernliving.com/recipes/oven-roasted-corn-on-cob",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Name:          "Oven-Roasted Corn On The Cob",
				Category:      models.Category{Value: "Side Dish"},
				Cuisine:       models.Cuisine{Value: "American"},
				DateModified:  "2023-11-05T21:00:17.338-05:00",
				DatePublished: "2019-05-14T09:02:49.000-04:00",
				Description: models.Description{
					Value: "Great corn doesn&#39;t get much easier than our Oven-Roasted Corn on the Cob recipe. The trick? Flavored butter and foil. See how to bake corn on the cob in the oven.",
				},
				Yield: models.Yield{Value: 4},
				Image: models.Image{
					Value: "https://www.southernliving.com/thmb/-bpB7uavaEqLXMhmTD0mz3Fj9c0=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/20220408_SL_OvenRoastedCornontheCobb_Beauty_1904-ed8011d403984f0aba111ec358359e02.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1/4 cup unsalted butter, softened",
						"1 tablespoon chopped fresh flat-leaf parsley",
						"2 medium garlic cloves, minced (2 tsp.)",
						"1 teaspoon chopped fresh rosemary",
						"1 teaspoon chopped fresh thyme",
						"3/4 teaspoon kosher salt",
						"1/2 teaspoon black pepper",
						"4 ears fresh corn, husks removed",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Make butter mixture: Preheat oven to 425°F. Stir together butter, parsley, garlic, rosemary, thyme, salt, and pepper in a bowl until evenly combined.",
						"Spread butter on corn: Spread 1 tablespoon herb butter on each corn cob.",
						"Wrap corn in foil: Wrap each corn on the cob individually in aluminum foil.",
						"Roast corn in oven: Place foil-wrapped corn on a baking sheet. Bake in preheated oven until corn is soft, 20 to 25 minutes, turning once halfway through cook time. Remove corn from foil, and serve",
					},
				},
				URL: "https://www.southernliving.com/recipes/oven-roasted-corn-on-cob",
			},
		},
		{
			name: "spendwithpennies.com",
			in:   "https://www.spendwithpennies.com/split-pea-soup/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				CookTime:      "PT130M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-10-15T20:32:25+00:00",
				Description: models.Description{
					Value: "Split pea soup is the perfect way to use up leftover ham. Split peas and ham are simmered in a delicious broth to create a thick and hearty soup!",
				},
				Keywords: models.Keywords{
					Values: "best recipe, ham and pea soup, how to make, leftover ham, split pea soup",
				},
				Image: models.Image{
					Value: "https://www.spendwithpennies.com/wp-content/uploads/2023/10/1200-Split-Pea-Soup-SpendWithPennies.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 cups dried split peas (green or yellow (14 oz))",
						"1 meaty ham bone (or 2 cups diced leftover ham)",
						"4 cups chicken broth",
						"4 cups water (or additional broth if desired)",
						"2 teaspoons dried parsley",
						"1 bay leaf",
						"3 ribs celery (diced)",
						"2  carrots (diced)",
						"1 large onion (diced)",
						"½ teaspoon black pepper",
						"½ teaspoon dried thyme",
						"salt to taste",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Sort through the peas to ensure there is no debris. Rinse and drain well.",
						"In a large pot, combine peas, ham, water, broth, parsley, and bay leaf. Bring to a boil, reduce heat " +
							"to low, and simmer covered for 1 hour.",
						"Add in celery, carrots, onion, pepper, thyme, and salt. Cover and simmer 45 minutes more.",
						"Remove ham bone and chop the meat. Return the meat to the soup and cook uncovered until thickened and the peas have broken down and the soup has thickened, about 20 minutes more.",
						"Discard the bay leaf and season with salt and additional pepper to taste.",
					},
				},
				Name: "Split Pea Soup",
				NutritionSchema: models.NutritionSchema{
					Calories:      "365 kcal",
					Carbohydrates: "45 g",
					Cholesterol:   "29 mg",
					Fat:           "9 g",
					Fiber:         "18 g",
					Protein:       "27 g",
					SaturatedFat:  "3 g",
					Sodium:        "900 mg",
					Sugar:         "8 g",
				},
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://www.spendwithpennies.com/split-pea-soup/",
			},
		},
		{
			name: "staysnatched.com",
			in:   "https://www.staysnatched.com/seafood-dressing",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "dinner"},
				CookTime:      "PT55M",
				Cuisine:       models.Cuisine{Value: "Louisiana"},
				DatePublished: "2022-11-14T12:29:00+00:00",
				Description: models.Description{
					Value: "This Seafood Dressing is a classic Southern side dish recipe made with cornbread and the Louisiana Holy Trinity vegetables including celery, green peppers, and onion. Load this with shrimp, crab, and even lobster! This is perfect for soul food dinners, Thanksgiving, or any gathering.",
				},
				Keywords: models.Keywords{
					Values: "how to make seafood dressing, Louisiana seafood dressing, seafood cornbread dressing, seafood dressing, shrimp cornbread dressing",
				},
				Image: models.Image{
					Value: "https://www.staysnatched.com/wp-content/uploads/2022/05/seafood-dressing-with-crab-and-shrimp-6-1.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"Cooked/Day Old Cornbread", "1 tablespoon olive oil",
						"1 cup finely chopped white onion", "1 cup finely chopped celery",
						"1/4 cup finely chopped green peppers", "4 garlic cloves minced",
						"2 teaspoons ground sage (Adjust to taste. 5 leaves if using fresh)",
						"1 teaspoon Old Bay Seasoning (or any seafood seasoning)", "1 teaspoon salt",
						"1 teaspoon ground black pepper", "14 oz cream of chicken soup",
						"2 - 2 1/2 cups broth (Any broth is fine including vegetable broth)",
						"2 eggs (Beaten)", "8 oz lump crab meat",
						"1 pound raw shrimp (Peeled and deveined)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"You can get the recipe for homemade cornbread here.",
						"Preheat oven to 350 degrees.",
						"Heat a skillet on medium-high heat. When hot, add the olive oil, celery, onions, garlic, and green peppers.",
						"Saute until the onions are translucent, fragrant, and the vegetables are soft. Remove the vegetables from the pan. Set aside to cool while prepping the remaining ingredients.",
						"Place day old cornbread in a large mixing bowl. Use a large spoon and/or your hands to break down the cornbread. You want to fully crumble it.",
						"Add in the sauteed veggies, ground sage, Old Bay seasoning, salt, and pepper. Mix well. Stop and taste the mixture here. This is where you want to adjust the seasoning and spice if necessary to suit your taste. Before adding the eggs, taste the dressing.",
						"Add in the eggs and cream of chicken soup. Stir.",
						"Slowly pour in the chicken broth. Start with a little chicken broth and then stir until the dressing is thick. Add more when needed. You should get the consistency of thick oatmeal.",
						"Add in 1/2 of the dressing mixture and spread it throughout the bottom of the pan.",
						"Add the shrimp and crab on top.",
						"Top it with a layer of the remaining dressing.",
						"Bake for 40-45 minutes uncovered on 350 degrees until a toothpick returns clean.",
					},
				},
				Name: "Seafood Dressing Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories:      "289 kcal",
					Carbohydrates: "17 g",
					Fat:           "9 g",
					Protein:       "9 g",
					Servings:      "1",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 10},
				URL:      "https://www.staysnatched.com/seafood-dressing",
			},
		},
		{
			name: "steamykitchen.com",
			in:   "https://steamykitchen.com/4474-korean-style-tacos-with-kogi-bbq-sauce.html",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Course"},
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2020-07-25T05:19:27+00:00",
				Description: models.Description{
					Value: "This is a great way to use your leftover pulled pork or roasted chicken. The BBQ Sauce from Kogi BBQ " +
						"was created by Chef Roy to be strong flavored enough to match the smokiness of BBQ’d pork or roasted " +
						"chicken. You can add use kimchi (spicy pickled Korean cabbage) to top the tacos, or make a quick " +
						"cucumber pickle like I have. The recipe for the quick cucumber pickle is below.",
				},
				Keywords: models.Keywords{Values: "korean bbq, taco"},
				Image: models.Image{
					Value: "https://steamykitchen.com/wp-content/uploads/2009/07/kogi-bbq-taco-151.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 pound cooked pulled pork or cooked shredded chicken",
						"12 corn or flour tortillas",
						"1/4 cup prepared store-bought Korean Kimchi ((optional))",
						"1 large English cucumber (or 2 Japanese cucumbers, sliced very thinly)",
						"2 tablespoons rice vinegar",
						"1/2 teaspoon sugar",
						"1/2 teaspoon finely minced fresh chili pepper (or more depending on your tastes)",
						"1 generous pinch of salt",
						"2 tablespoons Korean fermented hot pepper paste (gochujang)",
						"3 tablespoons sugar",
						"2 tablespoons soy sauce",
						"1 teaspoon rice wine vinegar",
						"2 teaspoons sesame oil",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Make the Quick Cucumber Pickle: Mix together all the Quick Pickle ingredients. You can make this a " +
							"few hours in advance and store in refrigerator. The longer it sits, the less “crunch” you’ll " +
							"have. I like making this cucumber pickle 1 hour prior, storing in refrigerator and serving it " +
							"cold on the tacos for texture and temperature contrast.",
						"Make the Koji BBQ Sauce: Whisk all BBQ sauce ingredients together until sugar has dissolved and mixture " +
							"is smooth. You can make this a few days in advance and store tightly covered in the refrigerator.",
						"Toss the Koji BBQ Sauce with your cooked pulled pork or shredded chicken. Warm the tortillas and serve" +
							" tacos with the Quick Cucumber Pickle.",
					},
				},
				Name: "Korean Style Tacos with Kogi BBQ Sauce Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories:      "503 kcal",
					Carbohydrates: "48 g",
					Cholesterol:   "102 mg",
					Fat:           "20 g",
					Fiber:         "5 g",
					Protein:       "34 g",
					SaturatedFat:  "6 g",
					Servings:      "1",
					Sodium:        "722 mg",
					Sugar:         "11 g",
				},
				PrepTime: "PT60M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://steamykitchen.com/4474-korean-style-tacos-with-kogi-bbq-sauce.html",
			},
		},
		{
			name: "streetkitchen.co",
			in:   "https://streetkitchen.co/recipe/thai-red-duck-curry/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Description: models.Description{
					Value: "This exquisite Thai Red Duck Curry is made with pineapple, cherry tomatoes and authentic red curry spices and coconut.",
				},
				Name:  "Thai Red Duck Curry",
				Yield: models.Yield{Value: 4},
				Image: models.Image{
					Value: "https://streetkitchen.co/wp-content/uploads/2022/10/Thai-Red-Duck-Curry-feature.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 tablespoons vegetable oil",
						"380g (approx. 2) fresh duck breast fillets, pat dry",
						"Salt to taste",
						"1 red onion, cut into thin wedges",
						"1 x 285g packet Street Kitchen Red Thai Curry Kit",
						"150g tinned pineapple slices in juice, drained, quartered",
						"1/2 x 250g punnet cherry tomatoes",
						"Thai basil, mint leaves and sliced red chilli for garnish",
						"Steamed jasmine rice and roti to serve",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Brush the oil over duck breasts and season with salt. Heat a non-stick frying pan over medium heat. When hot, add duck breasts skin-side down. Cook for 4-5 minutes or until the skin is golden and crisp. Turn and cook for 3-4 minutes or until just cooked through. Transfer to a plate. Cut into 1cm slices",
						"Discard excess oil, retaining 1 tbsp. Return pan to medium heat. Add spice pack and cook for 5 seconds. Add onion and cook for 3 minutes or until softened.\u00a0",
						"Stir in curry paste, coconut milk sachet, pineapple, tomatoes and 1/2 cup water. Stir until combined. Bring to the boil and simmer for 2 minutes. Nestle duck into sauce and simmer for 5 minutes or until duck is cooked through and sauce thickens. Spoon rice into serving bowls. Top with curry and garnish with basil and mint leaves and sliced red chilli. Serve with extra grilled roti.",
					},
				},
				URL: "https://streetkitchen.co/recipe/thai-red-duck-curry/",
			},
		},
		{
			name: "sunbasket.com",
			in:   "https://sunbasket.com/recipe/chicken-and-dumplings",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Category:  models.Category{Value: ""},
				CookTime:  "PT35M",
				Description: models.Description{
					Value: "This is Sunbasket’s easy (and gluten-free!) spin on an American classic.",
				},
				Keywords: models.Keywords{Values: ""},
				Image: models.Image{
					Value: "https://cdn.sunbasket.com/c46a59d6-5745-4b86-9574-0e3e4ab4318b.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 cup milk",
						"1 teaspoon apple cider vinegar",
						"2 ounces organic cremini or other button mushrooms",
						"1 or 2 cloves organic peeled fresh garlic",
						"Sunbasket gluten-free dumpling mix (Cup4Cup gluten-free flour - sugar - kosher salt - baking powder - " +
							"baking soda)",
						"Chicken options:",
						"2 to 4 boneless skinless chicken thighs (about 10 ounces total)",
						"2 boneless skinless chicken breasts (about 6 ounces each)",
						"1 cup organic mirepoix (onions - carrots - celery)",
						"1 tablespoon tomato paste",
						"4 or 5 sprigs organic fresh flat-leaf parsley",
						"2 tablespoons Cup4Cup gluten-free flour",
						"1 cup chicken broth",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Prep the milk-vinegar mixture",
						"Prep the vegetables; make the dumpling dough",
						"Prep and brown the chicken",
						"Cook the vegetables",
						"Finish the chicken; cook the dumplings",
						"Serve",
					},
				},
				Name: "Chicken and dumplings",
				NutritionSchema: models.NutritionSchema{
					Calories:     "520",
					Cholesterol:  "135mg",
					Fat:          "21g",
					Fiber:        "2g",
					Protein:      "34g",
					SaturatedFat: "3.5g",
					Sodium:       "830mg",
					Sugar:        "7g",
				},
				Yield: models.Yield{Value: 2},
				URL:   "https://sunbasket.com/recipe/chicken-and-dumplings",
			},
		},
		{
			name: "sundpaabudget.dk",
			in:   "https://sundpaabudget.dk/shawarma-bowl",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Kylling"},
				CookTime:      "PT15M",
				Cuisine:       models.Cuisine{Value: "Middelhavsmad"},
				DatePublished: "2023-01-13T17:24:13+00:00",
				Image: models.Image{
					Value: "https://sundpaabudget.dk/wp-content/uploads/2023/01/20230113125804_IMG_1660-1.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 dl ris", "300 g kebab", "1 ds kikærter", "2 peberfrugter", "1 rødløg",
						"1 spsk olie til stegning", "1 tsk spidskommen", "1 tsk paprika",
						"salt og peber", "4 tomater", "1 agurk", "3 spsk citronsaft", "salt og peber",
						"evt persille",
						"2,5 dl fraiche 5%",
						"1 stort fed hvidløg",
						"0,5 tsk sennep",
						"salt og peber",
						"syltet rødkålssalat ((se opskrift))",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Kog risen som anvist.",
						"Agurk-tomatsalat: vask og hak agurk og tomat i små fine tern. Vend dem sammen citronsaft og salt og peber.",
						"Rør ingredienserne sammen til dressingen og smag til.",
						"Peberfrugt skæres i grove strimler, rødløget skæres i tynde halve skiver, hæld vandet fra kikærterne.",
						"Opvarm lidt olie på en pande og rist først kikærterne sammen spidskommen og paprika samt salt og peber. Tag dem til side og rist peberfrugterne på begge sider så de tager farve.",
						"Steg kebaben som anvist på pakken.",
						"Anret shawarma bowlen med alle ingredienserne. Skal det være ekstra lækkert, kan jeg anbefale at lave syltet rødkålssalat til - det smager skønt.",
					},
				},
				Name: "Shawarma bowl",
				NutritionSchema: models.NutritionSchema{
					Calories: "300 kcal",
					Servings: "1",
				},
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://sundpaabudget.dk/shawarma-bowl",
			},
		},
		{
			name: "sunset.com",
			in:   "https://www.sunset.com/recipe/veggie-chili",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT0H0M",
				DatePublished: "2019-02-25",
				Description:   models.Description{Value: "Veggie Chili"},
				Image: models.Image{
					Value: "https://img.sunset02.com/sunsetm/wp-content-uploads/2019-03-29UTC04/veggie-chili-sun-61056-0219.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 cup each dried pinto beans and dried kidney beans, sorted of debris and rinsed, and soaked separately for two nights in 3 cups water each",
						"1 oz each dry guajillo chile and dry ancho chile, stemmed and seeded, and soaked separately overnight in 1 1/2 cups hot water each",
						"1 tbsp 90/10 canola and olive oil blend",
						"3 &#189; cups cups diced yellow onion", "1 &#189; cups cups diced carrot",
						"2 cups diced bell pepper", "3 oz jalapeño chile, seeded and diced",
						"4 garlic cloves, minced", "2 &#189; tsp salt", "1 tsp sweet paprika",
						"1 tsp smoked paprika", "1 tbsp chili powder", "&#189; tsp cayenne",
						"&#189; tbsp cumin", "1 tsp black pepper", "2 tbsp onion powder",
						"2 tbsp tomato paste",
						"&#189; can crushed tomatoes",
						"4 &#189; tbsp apple cider vinegar",
						"2 oz canned chipotle chiles in adobo",
						"1 &#189; tbsp tbsp. brown sugar",
						"1 cup fresh corn kernels",
						"3 tbsp finely chopped dark chocolate (80% cacao)",
						"Sour cream, sliced green onion, pepitas, parmesan crisps, for garnish (optional)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Drain and rinse beans; set aside. Drain chiles and reserve 1/2 cup of each soaking liquid; set aside.",
						"Heat a large pot over medium-low heat. Add oil, then onion, carrot, bell pepper, jalapeño, garlic, and 1/2 tsp. salt. Stir frequently on medium-low heat until liquid is released and vegetables start to stew down.",
						"Add the paprikas, chili powder, cayenne, cumin, black pepper, onion powder, tomato paste, beans, and tomatoes and stir to incorporate.",
						"In a medium bowl, use an immersion blender to purée the chiles, the reserved chile soaking liquids, 1 1/2 tbsp. vinegar, the chipotle in adobo, brown sugar, 1 1/2 tsp. salt, and 3 tbsp. water. Add mixture to pot, bring to a simmer, and simmer lightly for 1 1/2 hours.",
						"Add corn and 1/3 cup water to pot, bring back to a simmer, and let simmer lightly for 1 hour.",
						"Remove pot from heat. Stir in the chocolate, 1/2 tsp. salt, and the remaining 3 tbsp. vinegar. Serve warm, topped with garnishes, if using.",
					},
				},
				Name:     "Veggie Chili",
				PrepTime: "PT0H0M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://www.sunset.com/recipe/veggie-chili",
			},
		},
		{
			name: "sweetcsdesigns.com",
			in:   "https://sweetcsdesigns.com/roasted-tomato-marinara-sauce/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Pasta"},
				CookTime:      "PT45M",
				Cuisine:       models.Cuisine{Value: "italian"},
				DatePublished: "2022-03-09",
				Description: models.Description{
					Value: "Tired of using those bland jars of spaghetti sauce? This homemade roasted tomato marinara " +
						"sauce is packed with more flavor than those store-bought jars.",
				},
				Keywords: models.Keywords{
					Values: "pasta, sauce, tomato, italian, tomato sauce",
				},
				Image: models.Image{
					Value: "https://sweetcsdesigns.com/wp-content/uploads/2022/03/roasted-tomato-marinara-sauce-recipe-picture-720x720.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 pounds roma tomatoes",
						"1 onion, peeled and quartered",
						"2 tablespoons olive oil",
						"4 tablespoons balsamic vinegar, divided",
						"1 whole head garlic",
						"¼ cup fresh basil leaves, lightly packed",
						"Salt and Pepper",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat your oven to 425 and line a large baking sheet with foil (it will help keep the tomatoes from " +
							"sticking to the pan).",
						"Add the tomatoes and onion to the pan and drizzle over 1.5 tablespoons olive oil along with the balsamic " +
							"vinegar. Sprinkle generously with salt and pepper and toss to combine.",
						"Slice the top off the head of garlic, making sure to leave it intact.",
						"Drizzle the remaining half tablespoon of olive oil over the garlic and wrap it tightly in foil.",
						"Place the garlic cut side up on the baking sheet.",
						"Roast for 35 minutes, or until the tomatoes are caramelized and the onions are golden brown.",
						"Remove the garlic from the foil and squeeze out the cloves-- be careful, it will be hot!",
						"Add the garlic cloves to a food processor, along with the roasted tomatoes, roasted onion, remaining two " +
							"tablespoons balsamic vinegar, and fresh basil.",
						"Puree to your desired consistency and season with salt and pepper to taste.",
						"Store the marinara sauce in an airtight container in the fridge for up to five days or in an airtight " +
							"container in the freezer for up to three months.",
					},
				},
				Name: "Roasted Tomato Marinara Sauce",
				NutritionSchema: models.NutritionSchema{
					Calories:       "68 calories",
					Carbohydrates:  "8 grams carbohydrates",
					Cholesterol:    "0 milligrams cholesterol",
					Fat:            "4 grams fat",
					Fiber:          "2 grams fiber",
					Protein:        "1 grams protein",
					SaturatedFat:   "1 grams saturated fat",
					Sodium:         "46 milligrams sodium",
					Sugar:          "5 grams sugar",
					TransFat:       "0 grams trans fat",
					UnsaturatedFat: "3 grams unsaturated fat",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://sweetcsdesigns.com/roasted-tomato-marinara-sauce/",
			},
		},
		{
			name: "sweetpeasandsaffron.com",
			in:   "https://sweetpeasandsaffron.com/slow-cooker-cilantro-lime-chicken-tacos-freezer-slow-cooker/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Course"},
				CookTime:      "PT240M",
				Cuisine:       models.Cuisine{Value: "Mexican"},
				DatePublished: "2022-03-24T01:10:00+00:00",
				Description: models.Description{
					Value: "Cilantro lime chicken tacos are full of simple, bright flavors: cilantro, lime juice, garlic and a " +
						"touch of honey! No sauteeing required (just 15 min prep!), and easy to meal prep.",
				},
				Keywords: models.Keywords{
					Values: "cilantro lime chicken crockpot tacos, cilantro lime chicken tacos, crockpot tacos, meal prep tacos",
				},
				Image: models.Image{
					Value: "https://sweetpeasandsaffron.com/wp-content/uploads/2017/08/Slow-Cooker-Cilantro-Lime-Chicken-Tacos.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"4 chicken breasts (roughly 2 lbs; boneless skinless chicken thighs may also be used * see note 1)",
						"11.5 oz can of corn kernels (drained; 341 mL)",
						"15 oz can of black beans (drained &amp; rinsed; optional)",
						"1 red onion (sliced into strips)",
						"1/2 cup chicken stock",
						"2 cloves garlic",
						"1/4 teaspoon salt",
						"1/2 teaspoon cumin",
						"1/4 teaspoon ground coriander",
						"1 lime (zested)",
						"2 tablespoons honey (note 2)",
						"1/4 cup packed cilantro leaves",
						"Tortillas (2 small tortillas per person)",
						"shredded cabbage",
						"radish slices",
						"greek yogurt",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Combine - In the base of a 6-quart slow cooker, place the chicken breasts, corn and onion slices.",
						"Blend sauce - Using a stand or immersion blender, blend the sauce ingredients and pour over slow cooker " +
							"contents.",
						"Slow cook - Cover and cook on low for 4-5 hours, until chicken is cooked through.",
						"Serve - Shred chicken, then serve in tortillas topped with yogurt, shredded cabbage, and radish slices. " +
							"* see note 3",
					},
				},
				Name: "Cilantro Lime Chicken Crockpot Tacos",
				NutritionSchema: models.NutritionSchema{
					Calories:      "267 kcal",
					Carbohydrates: "28 g",
					Cholesterol:   "73 mg",
					Fat:           "4 g",
					Fiber:         "6 g",
					Protein:       "31 g",
					SaturatedFat:  "1 g",
					Servings:      "0.5 cup",
					Sodium:        "227 mg",
					Sugar:         "7 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://sweetpeasandsaffron.com/slow-cooker-cilantro-lime-chicken-tacos-freezer-slow-cooker/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
