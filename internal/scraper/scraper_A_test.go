package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_A(t *testing.T) {
	testcases := []testcase{
		{
			name: "abril.com",
			in:   "https://claudia.abril.com.br/receitas/estrogonofe-de-carne/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Carne"},
				CookTime:      "PT30M",
				CookingMethod: models.CookingMethod{Value: "Refogado"},
				Cuisine:       models.Cuisine{Value: "Brasileira"},
				DateModified:  "2020-02-05T07:51:35-0300",
				DatePublished: "2008-10-24",
				Description: models.Description{
					Value: "Derreta a manteiga e refogue a cebola até ficar transparente. Junte a carne e tempere " +
						"com o sal. Mexa até a carne dourar de todos os lados. Acrescente a mostarda, o catchup, " +
						"a pimenta-do-reino e o tomate picado. Cozinhe até formar um molho espesso. Se necessário, " +
						"adicione água quente aos poucos. Quando o molho estiver [&hellip;]",
				},
				Keywords: models.Keywords{
					Values: "Estrogonofe de carne, Refogado, Dia a Dia, Carne, Brasileira, creme de leite, ketchup (ou catchup), pimenta-do-reino",
				},
				Image: models.Image{
					Value: "https://claudia.abril.com.br/wp-content/uploads/2020/02/receita-estrogonofe-de-carne.jpg?" +
						"quality=85&strip=info&w=620&h=372&crop=1",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"500 gramas de alcatra cortada em tirinhas",
						"1/4 xícara (chá) de manteiga",
						"1 unidade de cebola picada",
						"1 colher (sobremesa) de mostarda",
						"1 colher (sopa) de ketchup (ou catchup)",
						"1 pitada de pimenta-do-reino",
						"1 unidade de tomate sem pele picado",
						"1 xícara (chá) de cogumelo variado | variados escorridos",
						"1 lata de creme de leite",
						"sal a gosto",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Derreta a manteiga e refogue a cebola até ficar transparente.Junte a carne e tempere com o sal.Mexa até a carne dourar de todos os lados.Acrescente a mostarda, o catchup, a pimenta-do-reino e o tomate picado.Cozinhe até formar um molho espesso.Se necessário, adicione água quente aos poucos.Quando o molho estiver encorpado e a carne macia, adicione os cogumelos e o creme de leite.Mexa por 1 minuto e retire do fogo.Sirva imediatamente, acompanhado de arroz e batata palha.Dica: Se juntar água ao refogar a carne, frite-a até todo o líquido evaporar.",
					},
				},
				Name:     "Estrogonofe de carne",
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://claudia.abril.com.br/receitas/estrogonofe-de-carne/",
			},
		},
		{
			name: "abuelascounter.com",
			in:   "https://abuelascounter.com/roasted-carrot-soup/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Soups"},
				CookTime:      "PT35M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-10-24T19:45:56+00:00",
				Keywords: models.Keywords{
					Values: "abuelau0026#039;s,cuban,easy recipes,healthy recipes,hosting,roasted carrot soup,soups,thanksgiving recipes,traditional",
				},
				Image: models.Image{
					Value: "https://abuelascounter.com/wp-content/uploads/2023/10/Roasted-Carrot-Soup-Recipe.jpeg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 cups of carrots, that have been peeled and diced",
						"1 ½ cups of butternut squash, that has been peeled and diced",
						"1 apple, peeled and diced (we like to use gala or granny smith apples)",
						"3 shallots, cut in quarters", "6 sprigs of thyme",
						"4 tablespoons of olive oil or avocado oil", "Freshly grated nutmeg",
						"3 ½ to 4 ½ cups of chicken or vegetable stock",
						"Salt and freshly cracked pepper",
						"Garnish: chives, sour cream, Calabrian chili peppers",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 425 degrees.",
						"Line a sheet pan with parchment paper or a silpat cover. Add all the vegetables, along with the apple, thyme and shallot. Toss in oil and sprinkle with 1 teaspoon of salt and freshly cracked pepper.",
						"Roast everything for 15 minutes. Then add 2 cups of stock. Roast for another 15 minutes or until all the vegetables have completely cooked through and are tender.",
						"Add everything to the blender (including the liquid) but make sure you leave out the thyme. Before adding the thyme to the blender, strip the thyme sprigs of their leaves. Discard the stems and only add the leaves to the blender.",
						"Add another 2 cups of stock. Add a small pinch of freshly grated nutmeg.",
						"Blend until completely smooth. Use a rubber spatula to move any chunks or pieces from the sides of the blender. If you want it to be a little thinner add another ½-1 cup of liquid or as much as you need to get it to your preferred consistency.",
						"Add to a pot keep warm on low heat. Serve and garnish.",
					},
				},
				Name:     "Roasted Carrot Soup",
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://abuelascounter.com/roasted-carrot-soup/",
			},
		},
		{
			name: "acouplecooks.com",
			in:   "https://www.acouplecooks.com/shaved-brussels-sprouts/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "Easy Shaved Brussels Sprouts",
				Description: models.Description{
					Value: "This shaved Brussels sprouts recipe make a tasty side dish that's full of texture and flavor! " +
						"Shredded Brussels are quick and crowd-pleasing.",
				},
				Keywords: models.Keywords{
					Values: "Shaved Brussels sprouts, Shaved Brussels sprouts recipe, shredded Brussel sprouts, shredded " +
						"Brussels sprouts",
				},
				Image: models.Image{
					Value: "https://www.acouplecooks.com/wp-content/uploads/2022/03/Shredded-Brussels-Sprouts-001-225x225.jpg",
				},
				URL: "https://www.acouplecooks.com/shaved-brussels-sprouts/",
				Ingredients: models.Ingredients{
					Values: []string{
						"1 pound Brussels sprouts (off the stalk)",
						"2 cloves garlic, minced",
						"1 small shallot, minced",
						"1/4 cup shredded Parmesan cheese (omit for vegan)",
						"½ teaspoon kosher salt, plus more to taste",
						"2 tablespoons olive oil",
						"1/4 cup Italian panko (optional, omit for gluten-free or use GF panko)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Shave the Brussels sprouts:\n\nWith a knife: Remove any tough outer layers with your fingers. " +
							"With a large Chef’s knife, cut the Brussels sprout in half lengthwise. Place the cut " +
							"side down and thinly slice cross-wise to create shreds. Separate the shreds with your " +
							"fingers. Discard the root end.",
						"With a food processor (fastest!): Use a food processor with the shredding disc attachment blade. " +
							"(Here&#8217;s a video.)",
						"With a mandolin: Slice the whole Brussels sprouts with a mandolin, taking proper safety precautions " +
							"to keep your fingers away from the blade. (Here&#8217;s a video.)",
						"In a medium bowl, stir together the minced garlic, shallot, Parmesan cheese, and kosher salt.",
						"In a large skillet, heat the olive oil over medium high heat. Add the Brussels sprouts and cook " +
							"for 4 minutes, stirring only occasionally, until tender and browned. Stir in the Parmesan " +
							"mixture and cook additional 3 to 4 minutes until lightly browned and fragrant. Remove the " +
							"heat and if desired, stir in the panko. Taste and add additional salt as necessary.",
					},
				},
				NutritionSchema: models.NutritionSchema{
					Calories:      "149 calories",
					Carbohydrates: "14.6 g",
					Cholesterol:   "3.6 mg",
					Fat:           "9.2 g",
					Fiber:         "6.5 g",
					Protein:       "6.5 g",
					SaturatedFat:  "2.1 g",
					Sodium:        "271.1 mg",
					Sugar:         "3 g",
					TransFat:      "0 g",
				},
				PrepTime:      "PT10M",
				CookTime:      "PT7M",
				Yield:         models.Yield{Value: 4},
				Category:      models.Category{Value: "Side dish"},
				CookingMethod: models.CookingMethod{Value: "Shredded"},
				Cuisine:       models.Cuisine{Value: "Vegetables"},
				DatePublished: "2022-03-23",
			},
		},
		{
			name: "addapinch.com",
			in:   "https://addapinch.com/easy-grape-jelly-meatballs-recipe/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizer"},
				CookTime:      "PT240M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-11-17T07:01:00+00:00",
				Description: models.Description{
					Value: "Grape Jelly Meatballs make the best appetizer recipe! It is a crowd favorite. Juicy meatballs are slow-cooked until tender in a sweet and spicy sauce. Made with three ingredients and a crowd favorite!",
				},
				Keywords: models.Keywords{Values: "grape jelly meatballs, grape jelly meatballs recipe"},
				Image: models.Image{
					Value: "https://addapinch.com/wp-content/uploads/2023/10/grape-jelly-meatballs-recipe-Grape-Jelly-Meatballs-Recipe-20230004-1.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 (32-ounce) package frozen meatballs", "1 cup grape jelly",
						"1 1/2 cups BBQ sauce (homemade or store-bought)", "green onions (slices)",
						"chives (chopped)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Add the meatballs, grape jelly, and BBQ sauce to the slow cooker. Stir to combine. Cover with the lid. Set the timer for 4 hours on the low setting, stirring after 2 hours. Leave on the warm setting until ready to serve.",
						"Optional topping to serve: Top with green onions or chives for serving.",
					},
				},
				Name: "Easy Grape Jelly Meatballs Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories:       "140 kcal",
					Carbohydrates:  "34 g",
					Fat:            "0.2 g",
					Fiber:          "1 g",
					Protein:        "0.4 g",
					SaturatedFat:   "0.02 g",
					Servings:       "1",
					Sodium:         "377 mg",
					Sugar:          "26 g",
					UnsaturatedFat: "0.08 g",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 12},
				URL:      "https://addapinch.com/easy-grape-jelly-meatballs-recipe/",
			},
		},
		{
			name: "afghankitchenrecipes.com",
			in:   "http://www.afghankitchenrecipes.com/recipe/norinj-palau-rice-with-orange/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Name:          "Norinj Palau (Rice with orange)",
				DatePublished: "2015-01-01",
				Image: models.Image{
					Value: "http://www.afghankitchenrecipes.com/wp-content/uploads/2015/01/afghan_norinj_pilaw-250x212.jpg",
				},
				Yield:    models.Yield{Value: 4},
				PrepTime: "PT10M",
				CookTime: "PT2H0M",
				Description: models.Description{
					Value: "Norinj Palau is one of traditional Afghan dishes and it has a lovely delicate flavour. " +
						"This pilau is prepared with the peel of the bitter (or Seville) oranges. It is quite a sweet dish.",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"450 g long grain white rice, preferably basmati",
						"75 ml vegetable oil",
						"2 medium onions, chopped",
						"1 medium chicken or 700–900 g lamb on the bone cut in pieces",
						"570 ml water, plus 110 ml water",
						"peel of 1 large orange",
						"50 g sugar",
						"50 g blanched and flaked almonds",
						"50 g blanched and flaked pistachios",
						"½ tsp saffron or egg yellow food colour (optional)",
						"25 ml rosewater (optional)",
						"1 tsp ground green or white cardamom seeds (optional)",
						"salt and pepper",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Measure out the rice and rinse several times until the water remains clear.",
						"Add fresh water and leave the rice to soak for at least half an hour.",
						"Heat the oil and add the chopped onions.",
						"Stir and fry them over a medium to high heat until golden brown and soft.",
						"Add the meat and fry until brown, turning frequently.",
						"Add 570 ml of water, salt and pepper and cook gently until the meat is tender.",
						"While the meat is cooking, wash and cut up the zest of a large orange into matchstick-sized pieces, " +
							"removing as much pith as possible.",
						"To remove any bitter taste, put the orange strips into a strainer and dip first in boiling water " +
							"and then in cold.",
						"Repeat this several times. Set aside.",
						"Make a syrup by bringing to the boil 110 ml of water and the 50 g of sugar. Add the orange peel, " +
							"the flaked almonds and pistachios to the boiling syrup.",
						"Boil for about 5 minutes, skimming off the thick froth when necessary. Strain and set aside the peel and nuts.",
						"Add the saffron and rosewater to the syrup and boil again gently for another 3 minutes.",
						"To cook the rice, strain the chicken stock (setting the meat to one side), and add the syrup.",
						"Make the syrup and stock up to 570 ml by adding extra water if necessary.",
						"The oil will be on the surface of the stock and this should also be included in the cooking of the rice.",
						"Bring the liquid to the boil in a large casserole. Drain the rice and then add it to the boiling liquid.",
						"Add salt, the nuts and the peel, reserving about a third for garnishing.",
						"Bring back to the boil, then cover with a tightly fitting lid, turn down the heat to medium and boil for " +
							"about 10 minutes until the rice is tender and all the liquid is absorbed.",
						"Add the meat, the remaining peel and nuts on top of the rice and cover with a tightly fitting lid. Put into " +
							"a preheated oven – 150°C (300°F, mark 2) – for 20–30 minutes. Or cook over a very low heat for the " +
							"same length of time.",
						"When serving, place the meat in the centre of a large dish, mound the rice over the top and then garnish " +
							"with the reserved orange peel and nuts.",
					},
				},
				Category: models.Category{Value: "Rice Dishes"},
				URL:      "http://www.afghankitchenrecipes.com/recipe/norinj-palau-rice-with-orange/",
			},
		},
		{
			name: "ah.nl",
			in:   "https://www.ah.nl/allerhande/recept/R-R1197438/boeuf-bourguignon-uit-de-oven-met-geroosterde-spruiten",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "hoofdgerecht"},
				DatePublished: "2022-09-30",
				Description: models.Description{
					Value: "Proef de authentieke Franse keuken met boeuf bourguignon oftewel: 'Rund op z'n Bourgondisch'. Rode wijn, laurierblad en mosterd zorgen voor diepe winterse smaken. Serveren met geroosterde spruiten.",
				},
				Keywords: models.Keywords{Values: "oven, hoofdgerecht, herfst, winter"},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 tenen knoflook", "2 sjalotten", "300 g winterpeen",
						"700 g runder sukadelappen", "1 el tarwebloem", "30 g ongezouten roomboter",
						"3 el milde olijfolie", "150 g baconreepjes", "70 g tomatenpuree",
						"10 g verse tijm", "300 ml rode wijn", "750 ml water", "2 laurierblaadjes",
						"1 tl dijonmosterd", "400 g champignons", "700 g geschoonde spruitjes", "10 g verse platte peterselie",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Verwarm de oven voor op 175 °C.",
						"Snijd de knoflook fijn. Snijd de sjalotten in halve ringen. Schil de winterpeen en snijd in blokjes van ca. een ½ cm.",
						"Dep de sukadelappen droog met keukenpapier, snijd in stukken van ca. 2 cm.",
						"Bestrooi het vlees rondom met peper en meng de bloem erdoor, zodat het vlees rondom met een laagje bedekt is.",
						"Verhit de braadpan met de boter en olie en braad het vlees 4 min. op hoog vuur tot het meeste vocht verdampt is en het vlees rondom bruin en krokant is.",
						"Voeg de baconreepjes en tomatenpuree toe en bak 4 min. op hoog vuur mee.",
						"Ris ondertussen de blaadjes van de tijm.",
						"Voeg de helft van de tijm, de wijn, het water, de knoflook, sjalot, peen, laurier en mosterd toe aan het vlees.",
						"Zet de braadpan met de deksel erop onder in de oven en stoof het vlees ca. 2½ uur tot het vlees mals is. Haal na 1½ uur de deksel van de pan.",
						"Snijd de champignons in kwarten en voeg de laatste 45 min. van de stooftijd toe aan het vlees.",
						"Test na 2½ uur of het vlees mals is, laat het zo nodig nog 30 min. langer stoven.",
						"Halveer ondertussen de spruiten en verdeel ze over de met bakpapier beklede traybake.",
						"Besprenkel met de rest van de olie en bestrooi met de rest van de tijm.",
						"Breng op smaak met peper en rooster de laatste 20 min. van de stooftijd boven in de oven met de boeuf bourguignon mee.",
						"Breng de stoof op smaak met peper en eventueel zout. Verwijder de laurierblaadjes.",
						"Snijd de peterselie fijn. Bestrooi de boeuf bourguignon met de peterselie en serveer met de geroosterde spruiten.",
					},
				},
				Name: "Boeuf bourguignon uit de oven met geroosterde spruiten",
				NutritionSchema: models.NutritionSchema{
					Calories:      "640 kcal energie",
					Carbohydrates: "23 g koolhydraten",
					Fat:           "33 g vet",
					Protein:       "52 g eiwit",
					SaturatedFat:  "12 g waarvan verzadigd",
				},
				Yield: models.Yield{Value: 4},
				URL:   "https://www.ah.nl/allerhande/recept/R-R1197438/boeuf-bourguignon-uit-de-oven-met-geroosterde-spruiten",
			},
		},
		{
			name: "akispetretzikis.com",
			in:   "https://akispetretzikis.com/recipe/6867/eukolos-mpaklavas",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT120M",
				Cuisine:       models.Cuisine{Value: "Greek"},
				DatePublished: "2022-09-26T10:02:10.000000Z",
				Description: models.Description{
					Value: "Εύκολος μπακλαβάς από τον Άκη Πετρετζίκη. Συνταγή για σιροπιαστό μπακλαβά με φύλλο κρούστας και γέμιση με καρύδια, κανέλα και γαρίφαλο!",
				},
				Keywords: models.Keywords{Values: "Εύκολος,μπακλαβάς"},
				Image: models.Image{
					Value: "https://akispetretzikis.com/photos/122904/eukolos-mpaklavas-26-7-22-ep1-site.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"450 γρ. φύλλα κρούστας για γλυκά",
						"400 γρ. βούτυρο", "1 κ.σ. κανέλα",
						"1 κ.γ. τριμμένο γαρίφαλο", "400 γρ. καρύδια",
						"50 γρ. φρυγανιά τριμμένη",
						"600 γρ. κρυσταλλική ζάχαρη", "400 ml νερό",
						"1 στικ κανέλα", "50 γρ. μέλι", "1 λεμόνι",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Βάζουμε στον πολυκόφτη τα καρύδια, τη φρυγανιά, την κανέλα, το γαρίφαλο, και χτυπάμε μέχρι να διαλυθούν όλα τα υλικά.",
						"Προθερμαίνουμε τον φούρνο στους 150&deg;C στον αέρα.",
						"Βουτυρώνουμε ένα ορθογώνιο ταψί 25x32 εκ.",
						"Απλώνουμε τα φύλλα κρούστας στον πάγκο εργασίας και τα κόβουμε στο μέγεθος του ταψιού. Το φύλλο που περισσεύει το κρατάμε στην άκρη.",
						"Βάζουμε τα μισά φύλλα στο ταψί, ρίχνουμε όλη τη γέμιση και καλύπτουμε με τα φύλλα που είχαμε αφήσει στην άκρη.",
						"Χαράσσουμε τον μπακλαβά σε διαγώνια κομμάτια.",
						"Ρίχνουμε το λιωμένο βούτυρο με μια κουτάλα σταδιακά και το απλώνουμε ως τις άκρες με ένα πινέλο.",
						"Μεταφέρουμε το ταψί στον φούρνο, ψηλά στη σχάρα και ψήνουμε για 2 ώρες μέχρι να γίνει τραγανός.",
						"Λίγα λεπτά πριν βγάλουμε τον μπακλαβά από τον φούρνο ετοιμάζουμε το σιρόπι.",
						"Τοποθετούμε μια κατσαρόλα σε δυνατή φωτιά να κάψει.",
						"Ρίχνουμε τη ζάχαρη, το νερό, το λεμόνι, την κανέλα, το μέλι, και αφήνουμε να πάρει μια βράση μέχρι να διαλυθεί η ζάχαρη.",
						"Μόλις πάρει μια βράση και έχουμε βγάλει από τον φούρνο τον μπακλαβά, σιροπιάζουμε με το καυτό σιρόπι τον καυτό μπακλαβά.",
						"Προσέχουμε και σιροπιάζουμε σταδιακά με την κουτάλα για να μην διαλυθεί ο μπακλαβάς.",
						"Αφήνουμε τον μπακλαβά στην άκρη μέχρι να &lsquo;&rsquo;πιει&rsquo;&rsquo; όλο το σιρόπι.",
						"Κόβουμε σε κομμάτια και σερβίρουμε.",
					},
				},
				Name: "Εύκολος μπακλαβάς",
				NutritionSchema: models.NutritionSchema{
					Calories: "508 calories",
				},
				PrepTime: "PT35M",
				URL:      "https://akispetretzikis.com/recipe/6867/eukolos-mpaklavas",
				Yield:    models.Yield{Value: 1},
			},
		},
		{
			name: "allrecipes.com",
			in:   "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dessert"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "American"},
				DateModified:  "2023-08-28T17:26:15.610-04:00",
				DatePublished: "1998-04-18T16:10:32.000-04:00",
				Description: models.Description{
					Value: "This chocolate chip cookie recipe makes delicious cookies with crisp edges and chewy middles. Try this wildly-popular cookie recipe for yourself!",
				},
				Image: models.Image{
					Value: "https://www.allrecipes.com/thmb/8xwaWAHtl_QLij6D-G0Z4B1HDVA=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/10813-best-chocolate-chip-cookies-mfs-146-4x3-b108aceffa6043a1ac81c3c5a9b034c8.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 cup butter, softened",
						"1 cup white sugar",
						"1 cup packed brown sugar",
						"2 eggs",
						"2 teaspoons vanilla extract",
						"1 teaspoon baking soda",
						"2 teaspoons hot water",
						"0.5 teaspoon salt",
						"3 cups all-purpose flour",
						"2 cups semisweet chocolate chips",
						"1 cup chopped walnuts",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Gather your ingredients, making sure your butter is softened, and your eggs are room temperature.",
						"Preheat the oven to 350 degrees F (175 degrees C).",
						"Beat butter, white sugar, and brown sugar with an electric mixer in a large bowl until smooth.",
						"Beat in eggs, one at a time, then stir in vanilla.",
						"Dissolve baking soda in hot water. Add to batter along with salt.",
						"Stir in flour, chocolate chips, and walnuts.",
						"Drop spoonfuls of dough 2 inches apart onto ungreased baking sheets.",
						"Bake in the preheated oven until edges are nicely browned, about 10 minutes.",
						"Cool on the baking sheets briefly before removing to a wire rack to cool completely.",
						"Store in an airtight container or serve immediately and enjoy!",
					},
				},
				Name: "Best Chocolate Chip Cookies",
				NutritionSchema: models.NutritionSchema{
					Calories:       "146 kcal",
					Carbohydrates:  "19 g",
					Cholesterol:    "10 mg",
					Fat:            "8 g",
					Fiber:          "1 g",
					Protein:        "2 g",
					SaturatedFat:   "4 g",
					Sodium:         "76 mg",
					UnsaturatedFat: "0 g",
				},
				PrepTime: "PT20M",
				Tools:    models.Tools{},
				Yield:    models.Yield{Value: 48},
				URL:      "https://www.allrecipes.com/recipe/10813/best-chocolate-chip-cookies/",
			},
		},
		{
			name: "altonbrown.com",
			in:   "https://altonbrown.com/recipes/the-apple-pie/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Sweets"},
				CookTime:      "PT70M",
				DatePublished: "2023-11-01T10:47:23+00:00",
				Description: models.Description{
					Value: "I can finally say, after tinkering with my original recipe on and off for a decade, this is the apple pie I want when I want apple pie, a craving that rises constant in me from the first rattle of fall leaves. We&#039;ve upped the grains of paradise (see note below), traded out the apple cider for apple cider vinegar, and enhanced the texture by inviting liquid pectin to the party, easily found in most mega marts. We&#039;ve also changed the apple mixture, replacing the original golden delicious with Pink Lady. I realize some of you may not be able to land these where you live, so feel free to replace whose with Winesaps, or just stick with Honeycrisps and Granny Smiths.Some of you will notice that I&#039;ve eschewed the pie bird relied upon in earlier iterations. As much as I dig their retro-homey vibe, truth is, slitting the top crust does as good a job of venting. Yes, the top crust does sometimes crack without the bird there to support the middle, but so what...I&#039;m going to eat it and I trust you will too.This recipe first appeared on altonbrown.com.Photo by Lynne Calamia",
				},
				Keywords: models.Keywords{Values: "Baking, Desserts, Fall, Fruit, Holidays, Thanksgiving, The Apple Pie"},
				Image:    models.Image{Value: "https://altonbrown.com/wp-content/uploads/2020/08/IMG_1623-scaled.jpeg"},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 3/4 cups all-purpose flour, plus extra for dusting",
						"1 teaspoon kosher salt", "1 tablespoon sugar",
						"12 tablespoons unsalted butter, cut into 1/2-inch pieces, chilled",
						"1/4 cup plus 2 teaspoons vegetable shortening, cut into 1/2-inch pieces, chilled",
						"7 tablespoons Laird&#39;s Applejack or apple brandy, such as calvados, chilled",
						"4 1/2 pounds (8 large) apples, mix of Pink Lady, Honeycrisp, and Granny Smith",
						"1/2 cup plus 1 tablespoon sugar, divided", "1/4 cup tapioca flour",
						"2 tablespoons liquid pectin (we like Certo)",
						"1 tablespoon apple cider vinegar", "1/4 teaspoon kosher salt",
						"1 teaspoon freshly ground grains of paradise",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Combine the flour, salt, and sugar in the bowl of your food processor and pulse 3-4 times. Add the butter and pulse until texture looks mealy, 5-6 pulses. Then, add the shortening and pulse until incorporated, another 3-4 pulses. Remove the lid and drizzle in 5 tablespoons of the Applejack. Replace the lid and pulse 5 times. Add the remaining Applejack and pulse until the mixture begins to hold together and pull away from the sides of the bowl.",
						"Dump the mixture onto a clean surface and squeeze together with your hands to form a smooth ball. Divide the ball in half and press each into a disk about 1-inch thick. Wrap each dough in plastic wrap and refrigerate for at least 1 hour. (You can refrigerate longer, even overnight, but the dough will have to sit at room temperature for 15 minutes to be malleable enough to roll.",
						"Peel, core, and slice the apples into 1/4-inch slices and move to a large mixing bowl. Add 1/4 cup of the sugar and toss with your hands to thoroughly coat. Set aside for 45 minutes, tossing halfway through, then transfer to a colander set over a large bowl and set aside to drain for 45 minutes.",
						"Transfer the accumulated juices (you should have about 1/4 cup) to a small saucepan and reduce over medium heat to 2 tablespoons, about 3 minutes. Set aside to cool.",
						"Stir the remaining sugar into the apple slices along with the tapioca flour, liquid pectin, apple cider vinegar, salt, and grains of paradise. Set aside.",
						"Crank your oven to 400℉ and move a rack to the lowest position.",
						`Place a 12" x 24" piece of wax paper on a clean work surface and lightly dust with flour. Remove the dough disks from the refrigerator and allow them to come to room temperature, 15 minutes. Discard the plastic wrap from one and place the dough on the wax paper. Dust with a bit more flour and roll into a 12" x 12" circle. Carefully peel the wax paper off and place the dough into the tart pan, gently pressing it into the edges. (See the note on dough movement below).`,
						"Arrange the apples in the bottom of the pan in concentric circles starting around the edges, working toward the center, which will result in a slight mound shape. Pour any remaining liquid evenly over the apples.",
						"Roll out the second dough disk in the same manner as the first. Place this dough over the apples and seal the edges of the pie, trimming any excess dough. Make a few slits in the top of the crust with a paring knife to give steam a way out. Park the pie on a foil-lined sheet pan and brush the top of the crust with the reduced juice. Bake for 1 hour, 10 minutes.",
						"Transfer the pie to a cooling rack and rest for at least 4 hours before removing from the tart pan and slicing.",
					},
				},
				Name:     "The Apple Pie",
				PrepTime: "PT60M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://altonbrown.com/recipes/the-apple-pie/",
			},
		},
		{
			name: "amazingribs.com",
			in:   "https://amazingribs.com/tested-recipes-chicken-recipes-crispy-grilled-buffalo-wings-recipe/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizer"},
				CookTime:      "PT30M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2020-01-14T18:29:00+00:00",
				Description: models.Description{
					Value: "True Buffalo wings are deep fried, but I love the flavor and convenience of cooking them on " +
						"the grill, and even smoking them first. And there is much less mess. Click here to tweet this",
				},
				Keywords: models.Keywords{
					Values: "barbecue, buffalo chicken wings, Chicken, chicken wings, grill, grilled buffalo chicken wings, " +
						"grilled chicken wings",
				},
				Image: models.Image{
					Value: "https://amazingribs.com/wp-content/uploads/2020/10/buffalo-wings.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 ounces cream cheese", "3 ounces quality blue cheese, crumbled",
						"1/2 cup half and half", "1/4 cup sour cream",
						"1/2 teaspoon Simon &amp; Garfunkel Seasoning or Poultry Seasoning",
						"1/2 cup melted salted butter", "2 cloves minced or pressed garlic",
						"1/2 cup Frank's Original RedHot Sauce",
						"24  whole chicken wings ((about 4 pounds (1.8 kg) for 24 whole wings))",
						"Morton Coarse Kosher Salt", "ground black pepper", "6 stalks of celery",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Prep. To create the blue cheese dip, take the cream cheese and the blue cheese out of the fridge and " +
							"let them come to room temp. Then smush them together with the spices. Mix in the sour cream and " +
							"half and half. Put them in a serving bowl and refrigerate. You can do this a day ahead. Cut up the " +
							"celery into 4-inch (10 cm) sections and put it back in the chiller.",
						"For the Buffalo hot sauce, melt the butter over a low heat and then adding the garlic. Let it simmer " +
							"for about a minute but don't let the garlic brown. Then add the Frank's RedHot sauce. Let them get " +
							"to know each other for at least 3 to 4 minutes. But remember, if you don't want to use the original " +
							"and you want to get creative, try one or more of the other sauces listed above. I'm partial to DC Mumbo " +
							"Sauce. Like the dip, the sauce can be made a day in advance.",
						"When it comes time to prep the wings, note that there are three distinct pieces of different thickness " +
							"and skin to meat ratio: (1) The tips (2) the flats or wingettes in the center, and (3) the drumettes " +
							"on the end that attaches to the shoulders. The thickness differences means they cook at different speeds " +
							"and finish at different times. The best thing to do is separate them into three parts with kitchen shears, " +
							"a sturdy knife, or a Chinese cleaver (my weapon of choice because the ka-chunk noise of chopping them off " +
							"is so very satisfying).",
						"The tips are almost all skin, really thin, and small enough that they often fall through the grates or " +
							"burn to a crisp. You can cook them if you wish, but I freeze them for use in making soup. Separate the " +
							"V shaped piece remaining at the joint between the flat and drumette. You will cook both these parts.",
						"Some folks like to season them with a spice rub. That works most of the time. I find that most commerci" +
							"al rubs are too salty for such thin cuts, and most have too much sugar that tends to burn during the " +
							"crisping phase. Besides, they just get lost under the sauces and dips. So I just season them with salt and " +
							"pepper. As Rachael Ray says: \"Easy peasy.\"",
						"Fire up. You can start them on a smoker if you wish, but I usually grill them. Set up the grill for 2-zone " +
							"cooking with the indirect side at about 325°F (162.8°C) to help crisp the skin and melt the fat. If you " +
							"wish, add wood to the direct side to create smoke. Use a lot of smoke.",
						"Cook. Add the wings to the indirect heat side of the grill and cook with the lid closed until the skins are " +
							"golden. That will probably take about 7 to 10 minutes per side. By then they are pretty close to done.",
						"To crisp the skin, move the wings to the direct heat side of your grill, high heat, lid open, and stand there, " +
							"turning frequently until the skin is dark golden to brown but not burnt, keeping a close eye on the " +
							"skinnier pieces, moving them to the indirect zone when they are done.",
						"Serve. Put the sauce in a big mixing bowl or pot and put it on the grill and get it warm. Stir or whisk well. " +
							"Keep warm. When the wings are done you can serve them with the sauce on the side for dipping, or just dump " +
							"them in with the sauce and toss or stir until they are coated. Then slide them onto a serving platter. Put " +
							"the celery sticks next to them, and serve with a bowl of Blue Cheese Dip. People can scoop some Blue " +
							"Cheese Sauce on their plates, and dip in the celery and wings.",
					},
				},
				Name:     "Crispy Grilled Buffalo Wings Recipe",
				PrepTime: "PT120M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://amazingribs.com/tested-recipes-chicken-recipes-crispy-grilled-buffalo-wings-recipe/",
			},
		},
		{
			name: "ambitiouskitchen.com",
			in:   "https://www.ambitiouskitchen.com/lemon-garlic-salmon/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				CookTime:      "PT18M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-05-05T11:00:05+00:00",
				Description: models.Description{
					Value: "Wonderful honey lemon garlic salmon covered in an easy lemon garlic butter marinade and baked to flaky " +
						"perfection. This flavorful lemon garlic salmon recipe makes a delicious, protein packed dinner served " +
						"with your favorite salad or side dishes, and the marinade is perfect for your go-to proteins.",
				},
				Image: models.Image{
					Value: "https://www.ambitiouskitchen.com/wp-content/uploads/2021/01/Lemon-Garlic-Salmon-5.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 pound salmon", "2 tablespoons butter, melted",
						"2 tablespoons honey (or sub pure maple syrup)",
						"1 teaspoon dijon mustard, preferably grainy dijon", "½ lemon, juiced",
						"Zest from 1 lemon", "½ teaspoon garlic powder (or 3 cloves garlic, minced)",
						"Freshly ground salt and pepper",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat the oven to 400 degrees F. Line a large baking sheet with parchment paper or foil and grease " +
							"lightly with olive oil or nonstick cooking spray. Place salmon skin side down.",
						"In a medium bowl, whisk together the melted butter, honey, dijon, lemon juice, lemon zest, garlic powder " +
							"and salt and pepper. Generously brush the salmon with the marinade.",
						"Place salmon in the oven and bake for 15-20 minutes or until salmon easily flakes with a fork. Mine is " +
							"always perfect at 16-18 minutes. Enjoy immediately.",
					},
				},
				Keywords: models.Keywords{
					Values: "lemon garlic butter salmon, lemon garlic salmon, lemon honey garlic salmon",
				},
				NutritionSchema: models.NutritionSchema{
					Calories:      "279 kcal",
					Carbohydrates: "9.3 g",
					Fat:           "17.2 g",
					Fiber:         "0.1 g",
					Protein:       "21.3 g",
					SaturatedFat:  "7.6 g",
					Servings:      "1",
					Sugar:         "8.2 g",
				},
				Name:     "Honey Lemon Garlic Salmon",
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.ambitiouskitchen.com/lemon-garlic-salmon/",
			},
		},
		{
			name: "archanaskitchen.com",
			in:   "https://www.archanaskitchen.com/karnataka-style-orange-peels-curry-recipe",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Name:          "Karnataka Style Orange Peels Curry Recipe",
				Cuisine:       models.Cuisine{Value: "Karnataka"},
				DateModified:  "2023-11-13T05:30:02+0000",
				DatePublished: "2017-10-05T00:23:00+0000",
				Description: models.Description{
					Value: "Did you know that we can make a yummy curry out of Orange Peels? It is tangy, sweetish, spicy, slightly bitter and bursting with flavors. It is an unique recipe. So next time you have some guests at home, make this recipe and impress your friends and family. It is filled with flavours and tastes delicious with almost everything. Next time you eat an orange, don't throw the peels, make a curry out of it.\nServe Karnataka Style Orange Peels Curry along with Cabbage Thoran and Whole Wheat Lachha Paratha for your weekday meal. It even tastes great with Steamed Rice.\nIf you like this recipe, you can also try other Karnataka recipes such as\n\nMavina Hannina Gojju Recipe\nMavina Hannina Gojju Recipe\nKarnataka Style Bassaru Palya Recipe",
				},
				Image: models.Image{
					Value: "https://www.archanaskitchen.com/images/archanaskitchen/1-Author/Smitha-Kalluraya/" +
						"Karnataka_style_Orange_Peels_Curry_.jpg",
				},
				PrepTime: "PT15M",
				CookTime: "PT15M",
				Yield:    models.Yield{Value: 4},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 Orange",
						"Tamarind, big lemon size",
						"3 tablespoon Jaggery, adjustable",
						"1 teaspoon Rasam Powder",
						"1 teaspoon Red Chilli powder",
						"1 tablespoon Rice flour",
						"Salt, to taste",
						"1 tablespoon Oil",
						"1 teaspoon Mustard seeds (Rai/ Kadugu)",
						"1/2 teaspoon Cumin seeds (Jeera)",
						"1/4 teaspoon Methi Seeds (Fenugreek Seeds)",
						"1 teaspoon White Urad Dal (Split)",
						"1 Dry Red Chilli, broken",
						"1 Green Chilli, slit",
						"Asafoetida (hing), a pinch",
						"1/4 teaspoon Turmeric powder (Haldi)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"To begin making the Karnataka Style Orange Peels Curry recipe, Peel out 2 oranges. Gently scrape the white " +
							"region from the orange peel using a spoon. Remove as much as possible as this will reduce the bitterness.",
						"Chop the orange peels finely.",
						"Soak tamarind in water. After it's soaked well, squeeze out tamarind water and keep aside.",
						"Heat oil in a heavy bottomed pan and temper with mustard seeds, urad dal, cumin seeds & fenugreek seeds. When " +
							"it splutters add curry leaves, red chilli powder, green chilli, hing and turmeric powder. Mix and add " +
							"chopped orange peel. Fry for 4-5 minutes.",
						"Later add tamarind water, salt, jaggery, red chilli powder and rasam powder. Add some water and " +
							"close the lid. Keep the flame on low medium and allow the peel to cook well. Mix in between.",
						"Meanwhile  in a mixing bowl, mix rice flour in water such that no lumps are formed.",
						"When orange peels are cooked well, add the rice flour mix to the curry and allow it to boil for 1-2 minutes. " +
							"Adding rice flour gives nice thickness to the gravy. If the desired consistency is attained, switch off.",
						"Serve Karnataka Style Orange Peels Curry along with Cabbage Thoran and Whole Wheat Lachha Paratha for your " +
							"weekday meal. It even tastes great with Steamed Rice.",
					},
				},
				Category: models.Category{Value: "Indian Curry Recipes"},
				Keywords: models.Keywords{
					Values: "South Indian Recipes,Indian Lunch Recipes,Orange Recipes,Karnataka Recipes",
				},
				URL: "https://www.archanaskitchen.com/karnataka-style-orange-peels-curry-recipe",
			},
		},
		{
			name: "arla.se",
			in:   "https://www.arla.se/recept/kycklingpasta-med-spenat-och-grillade-gronsaker/",
			want: models.RecipeSchema{
				AtContext: "https://schema.org",
				AtType:    models.SchemaType{Value: "Recipe"},
				Category:  models.Category{Value: "Huvudrätt, Kvällsmat, Lunch, Middag"},
				Description: models.Description{
					Value: "Supersnabb vardagspasta med krämig sås. Vill du göra rätten vegetarisk, tillsätt färdigkokta bönor eller kikärtor istället för kyckling. Klart!",
				},
				Keywords: models.Keywords{Values: "Kyckling,Pasta"},
				Image: models.Image{
					Value: "https://cdn-rdb.arla.com/Files/arla-se/2485582835/17d593bd-6d99-44ab-a5f9-4649eb3b2972.jpg?w=1300&h=525&mode=crop&ak=f525e733&hm=697c0698",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"450 g kycklinglårfilé", "2 msk Arla Köket® Smör- & rapsolja",
						"1 tsk salt", "1 krm svartpeppar",
						"2 dl Arla Köket® Lätt crème fraiche parmesan & vitlök",
						"65 g babyspenat", "300 g pasta",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Koka pastan enligt anvisning på förpackningen.",
						"Strimla kycklingköttet. Fräs det i smör-&rapsolja i en stekpanna, salta och peppra.",
						"Rör ner crème fraiche och låt koka ihop ca 5 min. Vänd ner spenaten och blanda med pastan. Servera direkt.",
					},
				},
				Name: "Kycklingpasta med spenat och parmesan",
				NutritionSchema: models.NutritionSchema{
					Calories:      "575 kcal",
					Carbohydrates: "57 g",
					Fat:           "22 g",
					Protein:       "34 g",
				},
				Yield: models.Yield{Value: 4},
				URL:   "https://www.arla.se/recept/kycklingpasta-med-spenat-och-grillade-gronsaker/",
			},
		},
		{
			name: "atelierdeschefs.fr",
			in:   "https://www.atelierdeschefs.fr/fr/recette/17741-boeuf-bourguignon-traditionnel.php",
			want: models.RecipeSchema{
				AtContext:   atContext,
				AtType:      models.SchemaType{Value: "Recipe"},
				Name:        "Bœuf bourguignon traditionnel",
				CookTime:    "P0Y0M0DT0H0M10800S",
				Description: models.Description{Value: "Une vraie recette de la tradition française: des morceaux de bœuf cuits longuement dans un bouillon au vin rouge."},
				Image:       models.Image{Value: "https://adc-dev-images-recipes.s3.eu-west-1.amazonaws.com/bourguignon_3bd.jpg"},
				Ingredients: models.Ingredients{
					Values: []string{
						"1.5 kg Boeuf à braiser ( jumeau, collier, macreuse )",
						"2 pièce(s) Carotte(s)",
						"1 pièce(s) Oignon(s)",
						"30 g Farine de blé",
						"2 pièce(s) Gousse(s) d'ail",
						"1.5 l Vin de Bourgogne",
						"3 cl Huile de tournesol",
						"6 pincée(s) Sel fin",
						"6 tour(s) Moulin à poivre",
						"40 cl Fond de veau",
						"150 g Lardon(s)",
						"150 g Oignon(s) grelot",
						"150 g Champignon(s) de Paris",
						"10 g Sucre en poudre",
						"50 g Beurre doux",
						"3 cl Huile d'olive",
						"0.25 botte(s) Persil plat",
					},
				},
				Instructions: models.Instructions{Values: []string{
					"Couper et dégraisser légèrement la viande. Éplucher et tailler en gros morceaux les carottes et l'oignon. Éplucher et dégermer les gousses d'ail.\nMettre la viande et la garniture dans le vin rouge, et faire mariner toute une nuit au réfrigérateur.",
					"Égoutter la viande et la garniture en conservant le vin. Séparer la garniture et la viande. Effeuiller le persil, conserver les tiges pour la cuisson et les feuilles pour le dressage.\n\nDans une cocotte chaude, mettre l'huile de tournesol et colorer les morceaux de viande environ 1 minute de chaque côté. Ajouter la garniture aromatique, assaisonner de sel fin, puis cuire doucement pendant 3 minutes. Singer (c'est-à-dire ajouter la farine) et cuire à nouveau 1 minute tout en mélangeant pour bien incorporer la farine. Mouiller avec le vin rouge puis avec le fond de veau. Ajouter les tiges de persil et compléter avec de l'eau si nécessaire. Faire bouillir puis baisser le feu et laisser mijoter pendant 2h30.\n\nLorsque la viande est cuite, la retirer de la cocotte. Passer la sauce au chinois pour la filtrer et vérifier sa texture. Si elle est encore trop liquide, la réduire pendant quelques minutes. La goûter et l'assaisonner de sel et de poivre.",
					"Éplucher les champignons au couteau. \n\nDisposer les oignons grelots dans une poêle. Ajouter de l'eau à mi-hauteur, 20 g de beurre et 1 cuillère à soupe de sucre. Couvrir au contact avec un papier sulfurisé et cuire jusqu'à évaporation complète de l'eau. Lorsque le sucre commence à caraméliser, ajouter 1 cuillère à soupe d'eau et bien enrober les oignons de caramel.\n\nDans une casserole d'eau froide, mettre les lardons et faire bouillir pour les blanchir. Bien les égoutter, puis les colorer dans une poêle antiadhésive bien chaude. Réserver ensuite sur du papier absorbant. Dans la même poêle, mettre un filet d'huile d'olive et faire sauter les champignons pour les colorer. Réserver.",
					"Ciseler les feuilles de persil.\nDans un plat, déposer la viande, verser la sauce dessus et disposer les garnitures.",
				}},
				PrepTime: "P0Y0M0DT0H0M1200S",
				URL:      "https://www.atelierdeschefs.fr/fr/recette/17741-boeuf-bourguignon-traditionnel.php",
				Yield:    models.Yield{Value: 6},
			},
		},
		{
			name: "averiecooks.com",
			in:   "https://www.averiecooks.com/slow-cooker-beef-stroganoff/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Slow Cooker"},
				CookTime:      "PT8H",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-03-15",
				Description: models.Description{
					Value: "A comfort food classic that everyone in the family LOVES! Hearty chunks of beef, rich and flavorful beef gravy, and served over a bed of warm noodles to soak up all that goodness! The EASIEST recipe for beef stroganoff ever because your Crock-Pot truly does all the work! Set it and forget it!",
				},
				Image: models.Image{
					Value: "https://www.averiecooks.com/wp-content/uploads/2022/03/beefstroganoff-13-480x480.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 pounds beef stew meat or beef chuck, diced into large bite-sized pieces or chunks",
						"1/2 cup white onion, diced small", "3 to 5 cloves garlic, finely minced",
						"1 teaspoon salt", "1 teaspoon freshly ground pepper",
						"1 teaspoon beef bouillon",
						"2 to 3 sprig fresh thyme OR 1/2 teaspoon dried thyme",
						"two 10-ounce cans cream of mushroom soup",
						"2 cups low sodium beef broth, plus more if desired",
						"1 tablespoon Dijon mustard, optional", "1 tablespoon Worcestershire sauce",
						"1/2 cup heavy cream, optional for a creamier sauce; at room temperature",
						"12 ounces wide egg noodles, cooked according to package directions (or your favorite pasta or mashed potatoes)",
						"Fresh parsley, finely minced; optional for garnishing",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"To a large 7 to 8-quart slow cooker, add the beef, onion, garlic, salt, pepper, beef bouillon, thyme, " +
							"and stir to combine; set aside.",
						"To a medium bowl, add the mushroom soup, beef broth, optional Dijon, Worcestershire sauce, and whisk " +
							"to combine.",
						"Pour the liquid over the contents in the slow cooker, stir to combine, cover the the lid, and cook on " +
							"high for 4 to 5 hours OR on low for 7 to 8 hours, or until done. Tip - At any time the beef " +
							"is slow cooking and you feel like it needs a bit more beef broth, it's fine to add a bit more, to taste.",
						"In the last 15 minutes of slow cooking, cook the egg noodles in a pot of boiling water according to " +
							"package directions; drain and set aside.* (See Notes about why I don't cook the noodles in " +
							"the slow cooker and cook them separately.)",
						"Optionally, if you want a creamier sauce, after the beef stroganoff has cooked and has cooled a bit " +
							"meaning it's not boiling nor bubbling, slowly you can add 1/2-cup heavy cream at room temperature " +
							"while whisking vigorously as you add it. Tip - Do NOT add cold cream to hot beef liquid because the " +
							"dairy proteins can separate, or break, and you will end up with a horribly ugly looking sauce after hours " +
							"and hours slow cooking. So make sure the cream is at room temp and the beef stroganoff has cooled a bit " +
							"if you are adding cream.",
						"Plate a bed of noodles, top with the beef and gravy mixture, and serve immediately. Extra beef " +
							"stroganoff will keep airtight in the fridge for up to 5 days and in the freezer for up to 4 months. " +
							"Tip - Because the noodles will continue to absorb moisture, including any of the beef gravy or sauce, " +
							"it's very important to store the beef mixture and the noodles separately in the fridge or freezer.",
					},
				},
				Name: "Slow Cooker Beef Stroganoff",
				NutritionSchema: models.NutritionSchema{
					Calories:       "575 calories",
					Carbohydrates:  "30 grams carbohydrates",
					Cholesterol:    "191 milligrams cholesterol",
					Fat:            "23 grams fat",
					Fiber:          "2 grams fiber",
					Protein:        "62 grams protein",
					SaturatedFat:   "11 grams saturated fat",
					Servings:       "1",
					Sodium:         "1288 milligrams sodium",
					Sugar:          "6 grams sugar",
					TransFat:       "1 grams trans fat",
					UnsaturatedFat: "12 grams unsaturated fat",
				},
				Keywords: models.Keywords{
					Values: "slow cooker beef stroganoff, easy beef stroganoff recipe, beef stroganoff with cream of mushroom " +
						"soup, beef stroganoff with stew meat, slow cooker beef stroganoff recipe",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://www.averiecooks.com/slow-cooker-beef-stroganoff/",
			},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
