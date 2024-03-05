package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_P(t *testing.T) {
	testcases := []testcase{
		{
			name: "paleorunningmomma.com",
			in:   "https://www.paleorunningmomma.com/grain-free-peanut-butter-granola-bars-vegan-paleo-option/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Breakfast"},
				Cuisine:       models.Cuisine{Value: "Gluten-free"},
				DatePublished: "2021-03-26T18:11:21+00:00",
				Description: models.Description{
					Value: "These no bake, chewy peanut butter granola bars are a breeze to make and so addicting!  They’re " +
						"gluten free and grain free, with both vegan and paleo options.  Perfect for quick snacks, these " +
						"homemade granola bars are perfect right out of the fridge or freezer.",
				},
				Keywords: models.Keywords{
					Values: "bars, egg-free, no-bake, nut butter, paleo, peanut butter, vegan",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 1/2 cups pecan halves ( or walnuts)",
						"1 1/2 cups almonds",
						"1 cup unsweetened coconut flakes",
						"1/2 tsp salt",
						"2 tablespoons organic coconut oil (melted (use refined for neutral flavor))",
						"2/3 cup peanut butter or other nut butter like almond (cashew, walnut, or sunflower butter (for paleo))",
						"1/4 cup + 2 Tbsp pure maple syrup (or raw honey (for paleo))",
						"1 tsp pure vanilla extract",
						"2/3 cup mini chocolate chips (dairy free )",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Place the nuts in a food processor and pulse several times to “chop” them into a crumbly texture - a " +
							"few larger pieces are a good thing - don’t overmix!",
						"Transfer the nuts to a large mixing bowl and stir in coconut flakes and salt to evenly combine.",
						"Place the melted coconut oil in a medium bowl and whisk in the peanut butter and honey or maple syrup. " +
							"Once mixture is smooth and well combined, stir in the vanilla.",
						"Pour the wet mixture into the large bowl with the dry ingredients and stir to fully combine - I used a " +
							"silicone spatula for this step. Thoroughly mix to make sure all the dry mixture is coated. Once " +
							"coated, gently stir in the chocolate chips.",
						"Line an 8 x 8” or 9 x 9” square pan with parchment paper along the bottom and sides, with extra up the " +
							"sides for easy removal. Transfer mixture in and press down, using your hands, or another piece of " +
							"parchment paper to get it packed tightly into the pan.",
						"Cover the top with parchment or plastic wrap, then set in the freezer for at least 1 hour to firm, or " +
							"longer if you have time.",
						"Remove pan from freezer and grab two ends of the parchment paper to remove the bars, set on a cutting board.",
						"Using a long very sharp knife, cut into 15-20 bars. You can wrap them in parchment individually storing " +
							"in the fridge (for up to two weeks) or freezer for longer.",
						"Bars will start to melt around room temp due to the coconut oil, so they’ll need to be kept chilled to " +
							"stay firm. Enjoy!",
					},
				},
				Name: "Grain Free Peanut Butter Granola Bars {Vegan, Paleo Option}",
				NutritionSchema: models.NutritionSchema{
					Calories:      "249 kcal",
					Carbohydrates: "12 g",
					Cholesterol:   "1 mg",
					Fat:           "21 g",
					Fiber:         "3 g",
					Protein:       "6 g",
					SaturatedFat:  "7 g",
					Servings:      "1",
					Sodium:        "106 mg",
					Sugar:         "7 g",
					TransFat:      "1 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 20},
				URL:      "https://www.paleorunningmomma.com/grain-free-peanut-butter-granola-bars-vegan-paleo-option/",
			},
		},
		{
			name: "panelinha.com.br",
			in:   "https://www.panelinha.com.br/receita/Frango-ao-curry",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Aves"},
				Cuisine:       models.Cuisine{Value: "Prática"},
				DatePublished: "2000-05-13",
				Description: models.Description{
					Value: "A lista de vantagens desta receita é longa: fácil, rápida, tem poucos ingredientes e muito sabor. E " +
						"tem mais: você pode preparar bem antes da hora de servir. Graças ao caldinho delicioso de creme de leite, " +
						"maçãs e especiarias, o frango segue macio, macio mesmo depois de requentado. Agora, o inegociável: investir " +
						"em um curry de qualidade, já que é ele que dá todo o sabor ao preparo.",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 filés de peito de frango (cerca de 500g)",
						"2 maçãs fuji",
						"1 cebola",
						"2 dentes de alho",
						"2 colheres (sopa) de curry",
						"1 xícara (chá) de creme de leite fresco",
						"2 colheres (sopa) de azeite",
						"sal e pimenta-do-reino moída na hora a gosto",
						"folhas de coentro a gosto para servir",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Corte os filés de frango em cubos de 3 cm e transfira para uma tigela. Tempere com 1 colher (chá) de sal, " +
							"o curry e pimenta a gosto e mantenha em temperatura ambiente enquanto prepara os outros ingredientes — assim " +
							"ele absorve melhor o sabor do curry e perde o gelo antes de ir para a panela.",
						"Descasque e pique fino a cebola e os dentes de alho. Descasque e corte as maçãs em fatias de 1,5 cm; " +
							"descarte o miolo com as sementes, corte as fatias em tiras e as tiras em cubos.",
						"Leve uma panela média ao fogo médio. Quando aquecer, regue com o azeite, adicione a cebola e tempere com " +
							"uma pitada de sal. Refogue por cerca de 3 minutos, até murchar. Junte o alho e mexa bem por 1 minuto " +
							"para perfumar.",
						"Acrescente os cubos de frango e refogue por cerca de 1 minuto, até perderem a aparência crua — evite" +
							" refogar em excesso para que o frango não resseque; ele vai terminar de cozinhar com o restante dos " +
							"ingredientes.",
						"Junte os cubos de maçã, regue com o creme de leite e misture bem. Assim que começar a ferver, abaixe o" +
							" fogo e deixe cozinhar por cerca de 10 minutos, até que o frango e a maçã estejam cozidos e o molho " +
							"levemente encorpado. Sirva a seguir com folhas de coentro a gosto.",
					},
				},
				Name:  "Frango ao curry com maçã",
				Yield: models.Yield{Value: 4},
				URL:   "https://www.panelinha.com.br/receita/Frango-ao-curry",
			},
		},
		{
			name: "paninihappy.com",
			in:   "https://paninihappy.com/why-you-need-this-pumpkin-muffin-recipe/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "My Favorite Pumpkin Muffins",
				PrepTime:  "PT10M",
				CookTime:  "PT20M",
				Description: models.Description{
					Value: "Because I’ve made them many times over the years and they’re the best pumpkin muffins I’ve tasted — fluffy, flavorful, unfussy, nice doming.\n\nBecause even though baking with pumpkin can be kind of a seasonal fad, it’s a delicious one, so pumpkin on!\n\nBecause now that it’s October you’re undoubtedly going to need to bring a crowd-pleasing, autumn-appropriate dish to school/work/church/soccer, etc. Or you’re simply going to want one at home on a brisk autumn afternoon.\n\nBecause you can easily double the recipe for a big group — in fact, the orignal recipe from Erin Cooks (one of my earliest favorite food blogs) makes a whole lot of muffins — or make them in mini muffin pans or mini loaf pans for lunch boxes or cute gifts.\n\nBecause pumpkin + cake mix does not equal a recipe (yeah, I said it — sorry, Pinterest!).\n\nBecause you probably have all the ingredients on hand already (especially in October, because pumpkin time).\n\nBecause they can pass for breakfast, dessert or even a side dish — versatility awaits!\n\nBecause friends love to receive the occasional pumpkin muffin surprise on their doorstep.\n\nBecause baking these muffins doubles as an awesome home fragrance for your kitchen.\n\nBecause…oh, just turn on the oven and make ’em!",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Yield: models.Yield{Value: 12},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 1/2 cups flour",
						"1 teaspoon baking soda",
						"1 teaspoon baking powder",
						"1 teaspoon cinnamon",
						"1/2 teaspoon salt",
						"1/4 teaspoon ground ginger",
						"1/4 teaspoon ground nutmeg",
						"2 eggs",
						"3/4 cup sugar",
						"1 cup pumpkin purée",
						"3/4 cup vegetable oil",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Heat the oven to 400ºF.",
						"Combine all of the dry ingredients in a medium bowl. Beat the eggs, sugar, pumpkin, and oil until smooth. " +
							"Pour the pumpkin mixture into the dry ingredients and mix just until blended.",
						"Grease a muffin tin or fill your tin with cupcake papers. Fill the wells with the batter until they are " +
							"2/3 of the way full. Bake for 16-20 minutes. Cool 5 minutes and then complete the cooling process" +
							" on a wire rack.",
					},
				},
				URL: "https://paninihappy.com/why-you-need-this-pumpkin-muffin-recipe/",
			},
		},
		{
			name: "persnicketyplates.com",
			in:   "https://www.persnicketyplates.com/easy-slow-cooker-french-toast-casserole/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Breakfast"},
				CookTime:      "PT240M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-11-13T23:57:32+00:00",
				Description: models.Description{
					Value: "Slow Cooker French Toast Casserole has only a few ingredients, preps overnight, and is perfect for Christmas morning, brunch, or just a weekend breakfast. Full of cinnamon flavor, this easy crockpot casserole is perfect finished off with a dusting of powdered sugar and some fresh fruit.",
				},
				Keywords: models.Keywords{
					Values: "crockpot french toast casserole, overnight slow cooker french toast casserole, slow cooker french toast casserole",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 loaf Italian bread (cut into 1&quot; pieces. (This should be 6-7 cups of bread))",
						"6 large eggs", "2 cups milk (I use 2%)", "2 teaspoons cinnamon (divided)",
						"¼ cup salted butter (at room temperature)", "⅓ cup brown sugar",
						"½ teaspoon vanilla extract", "powdered sugar (for topping)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						`Cube the bread into approximately 1" pieces and spread onto a lipped baking sheet. Bake in a 170°F preheated oven for 30-60 minutes, or until dry. Alternatively, you can cube the bread and let it sit out overnight to stale.`,
						"In a large mixing bowl, whisk together the eggs, milk, vanilla, and 1 teaspoon of the cinnamon.",
						"Add the stale cubed bread to the egg mixture and stir so all pieces of the bread are fully coated. Cover the bowl with plastic wrap and move it to the fridge for 4-8 hours.",
						"In a small mixing bowl, use a pastry cutter or a fork to cut together the butter, brown sugar, and remaining cinnamon. Cover and put in the fridge as well.",
						"When ready, spray a 6 quart slow cooker with non-stick spray and pour in the soaked bread pieces.",
						"Crumble the butter mixture evenly over the top.",
						"Cover and cook on HIGH for 2 hours or LOW for 4 hours, or until the center is set.",
						"Turn the slow cooker to WARM and remove the lid. Let stand for 15 minutes to set before serving.",
						"Top with powdered sugar, syrup, and berries.",
					},
				},
				Name: "Slow Cooker French Toast Casserole",
				NutritionSchema: models.NutritionSchema{
					Calories:       "494 kcal",
					Carbohydrates:  "42 g",
					Cholesterol:    "162 mg",
					Fat:            "31 g",
					Fiber:          "2 g",
					Protein:        "11 g",
					SaturatedFat:   "17 g",
					Servings:       "1 g",
					Sodium:         "335 mg",
					Sugar:          "30 g",
					TransFat:       "0.4 g",
					UnsaturatedFat: "13 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://www.persnicketyplates.com/easy-slow-cooker-french-toast-casserole/",
			},
		},
		{
			name: "pickuplimes.com",
			in:   "https://www.pickuplimes.com/recipe/the-best-vegan-chow-mein-800",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main"},
				CookTime:      "PT15M",
				Cuisine:       models.Cuisine{Value: "Asian-inspired"},
				DatePublished: "Sept. 10, 2023, 3:54 p.m.",
				Description: models.Description{
					Value: "To effortlessly serve up this dish that outshines takeout, the secret is having all your vegetables chopped and prepped. Once you begin cooking, this recipe comes together swiftly, and with everything ready, the entire meal flows seamlessly.",
				},
				Keywords: models.Keywords{Values: "peanut free, tree nut free, vegan, vegetarian, plant-based"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 Tbsp dark soy sauce", "2 Tbsp water", "2 Tbsp toasted sesame oil",
						"&#189; vegetable bouillon cube", "cornstarch", "agave syrup",
						"sriracha hot sauce", "rice vinegar", "&#189; Tbsp vegetable oil",
						"1 cup mock chicken pieces", "1 Tbsp vegetable oil", "1 medium onion",
						"2 clove garlic", "5.291 oz vegan chow mein noodles", "&#188; napa cabbage",
						"2 carrot", "1 head of bok choy",
					},
				},
				Name:     "The Best Vegan Chow Mein",
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 3},
				URL:      "https://www.pickuplimes.com/recipe/the-best-vegan-chow-mein-800",
			},
		},
		{
			name: "pingodoce.pt",
			in:   "https://www.pingodoce.pt/receitas/tarte-de-alho-frances-caramelizado/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				DatePublished: "2021-09-01",
				Description: models.Description{
					Value: "O alho-francês é um ingrediente muito versátil, que fica bem em várias receitas. Aqui, preparámos uma tarte de alho-francês caramelizado, com massa folhada e queijo. Experimente e delicie-se.",
				},
				Keywords: models.Keywords{Values: "alho-francês, tartes"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"800 g alho-francês", "1 c. de sopa azeite Pingo Doce",
						"1 c. de chá flor de sal Pingo Doce", "1 q.b. pimenta",
						"1 q.b. tomilho fresco", "150 ml vinagre de vinho branco Pingo Doce",
						"3 c. de sopa manteiga", "1 c. de sopa açúcar mascavado Pingo Doce",
						"2 c. de sopa mostarda", "3 c. de sopa queijo parmesão ralado Pingo Doce",
						"230 g massa folhada Pingo Doce",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Pré-aqueça o forno a 200 °C.",
						"Corte dois alhos a meio e dois em quatro, eliminando as folhas.",
						"Coloque-os num tabuleiro e regue com o azeite. Tempere com flor de sal, pimenta e tomilho. Leve ao forno por 20 minutos ou até ficarem dourados.",
						"Numa frigideira de ferro ou que possa ir ao forno, adicione o vinagre e deixe cozinhar cerca de 1 minuto ou até reduzir.",
						"Junte a manteiga, o açúcar e mais um pouco de tomilho. Cozinhe durante cerca de 2 minutos.",
						"Retire os alhos do forno e coloque-os na frigideira alinhados sobre o molho de manteiga. Tente ao máximo que fiquem apertadinhos.",
						"Pincele-os com mostarda, polvilhe com o parmesão e cubra com a massa folhada. Empurre as bordas da massa para dentro.",
						"Leve ao forno por 35 minutos ou até a massa estar dourada.",
						"Retire do forno, deixe que arrefeça ligeiramente e vire para um prato.",
						"Sirva ainda quente.",
					},
				},
				Name: "Tarte de alho-francês caramelizado",
				NutritionSchema: models.NutritionSchema{
					Calories: "411 calories",
				},
				PrepTime: "PT45M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.pingodoce.pt/receitas/tarte-de-alho-frances-caramelizado/",
			},
		},
		{
			name: "pinkowlkitchen.com",
			in:   "https://pinkowlkitchen.com/cajun-dirty-rice-with-smoked-sausage/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Course"},
				CookTime:      "PT20M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-04-08T15:41:52+00:00",
				Description: models.Description{
					Value: "This Cajun dirty rice dish is made with chopped veggies, cajun spices, and smoked sausage and is loaded with flavor! Serve this rice with chicken, fish, or shrimp to create the ultimate Cajun meal.",
				},
				Keywords: models.Keywords{Values: "cajun recipes, cajun rice, dirty rice, rice side dishes"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 cup long-grain rice", "2 1/4 cups chicken stock (divided)",
						"1 tablespoon butter (unsalted)", "1 tablespoon olive oil (extra virgin)",
						"12 ounces smoked sausage (about 4 medium links)",
						"1 small yellow onion (chopped)",
						"1 small red bell pepper (seeded and chopped)",
						"1 large stalk celery (chopped into half moons)", "3 cloves garlic (minced)",
						"1 teaspoon onion powder", "1 teaspoon garlic powder",
						"1 teaspoon smoked paprika", "1/2 teaspoon dried thyme",
						"1/2 teaspoon dried oregano", "1/2 teaspoon cayenne pepper",
						"salt and cracked black pepper (to taste)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Cook the rice in 2 cups of chicken stock until tender according to package instructions. Fluff the rice with a fork and set aside.",
						"In a medium heavy-bottomed skillet over medium-high heat, melt butter with the olive oil. Add the smoked sausage to the skillet and cook until the sausage is browned, about 3-5 minutes.",
						"Add the onion, bell pepper, and celery and cook until the veggies begin to soften, about 2-3 minutes. Add the garlic and cook until fragrant, about another minute.",
						"Season the sausage and vegetables with garlic powder, onion powder, smoked paprika, dried thyme, dried oregano, salt, black pepper, and cayenne pepper and stir with a wooden spoon.",
						"Add a little bit of chicken stock to the skillet to deglaze the pan, scraping up any brown bits with the wooden spoon. Add the cooked rice to the skillet and stir to combine the rice with the veggies and sausage. Continue to cook while stirring for another 2 to 3 minutes. Add a little more chicken stock if needed so that the rice isn't dry.",
						"Remove the pan from the heat and transfer the rice to a serving dish. Garnish with fresh chopped parsley if desired. Enjoy!",
					},
				},
				Name: "Cajun Dirty Rice with Smoked Sausage",
				NutritionSchema: models.NutritionSchema{
					Calories:      "389 kcal",
					Carbohydrates: "45.4 g",
					Cholesterol:   "63 mg",
					Fat:           "15.5 g",
					Fiber:         "1.9 g",
					Protein:       "17 g",
					SaturatedFat:  "4.8 g",
					Servings:      "1",
					Sodium:        "942 mg",
					Sugar:         "4.5 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://pinkowlkitchen.com/cajun-dirty-rice-with-smoked-sausage/",
			},
		},
		{
			name: "platingpixels.com",
			in:   "https://www.platingpixels.com/mushroom-tart-recipe/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizer"},
				CookTime:      "PT15M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-06-17T13:15:30+00:00",
				Description: models.Description{
					Value: "With only 15 minutes of prep, this decadent, yet easy mushroom tart is ready in less than 30 minutes.",
				},
				Keywords: models.Keywords{Values: "mushroom flatbread, Mushroom Tart, Mushroom Tart recipe"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 tablespoon olive oil", "1 medium yellow onion (diced)",
						"4-6 cups sliced Baby Bella mushrooms (about 16-ounces)", "½ teaspoon salt",
						"¼ teaspoon ground black pepper", "2 cloves minced garlic",
						"1 tablespoon fresh thyme leaves", "1 sheet Jus-Rol™ Puff Pastry Dough",
						"8- ounces shredded Italian cheese blend (or Mozzarella)",
						"1 medium egg (lightly beaten)",
					},
				}, Instructions: models.Instructions{
					Values: []string{
						"Preheat the oven to 425° F. Remove the Puff Pastry Dough from the fridge at let sit at room temp.",
						"Meanwhile, heat olive oil in a pan to medium heat. Add onions and saute until fragrant and slightly translucent, about 5 minutes. Add the mushrooms, salt, and pepper. Cook stirring occasionally until the mushrooms are browned and have released liquids, then the liquid evaporates, about 8 minutes.",
						"Stir in garlic and thyme and cook 2 minutes more. Remove the pan from the heat.",
						"Unroll the puff pastry sheet on a baking sheet, keeping the parchment on the bottom to prevent sticking. Score the edges using a sharp knife, leaving a 1-inch border. Then use a fork to prick the center puff pastry sheet every 2 to 3 inches.",
						"Top the puff pastry with grated cheese, leaving the 1-inch border of dough, followed by the sauteed mushrooms and onion mixture.",
						"Using a pastry brush, brush the dough border with the beaten egg to help with browning.",
						"Bake for 14-16 minutes, or until the puff pastry edges are golden brown.",
						"Remove from oven and let cool 2-3 minutes. Cut into squares and serve.",
					},
				},
				Name: "Mushroom Tart Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories:       "197 kcal",
					Carbohydrates:  "12 g",
					Cholesterol:    "29 mg",
					Fat:            "14 g",
					Fiber:          "1 g",
					Protein:        "7 g",
					SaturatedFat:   "5 g",
					Servings:       "1",
					Sodium:         "274 mg",
					Sugar:          "1 g",
					TransFat:       "1 g",
					UnsaturatedFat: "8 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 12},
				URL:      "https://www.platingpixels.com/mushroom-tart-recipe/",
			},
		},
		{
			name: "plowingthroughlife.com",
			in:   "https://plowingthroughlife.com/the-best-rich-and-moist-chocolate-cake/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dessert"},
				CookTime:      "PT35M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-01-21T10:07:00+00:00",
				Description: models.Description{
					Value: "Are you craving a delicious chocolate cake? This is the BEST rich and moist chocolate cake from a box mix that is unbelievably EASY to make! This is a bakery quality cake that anyone can make at home with common ingredients!",
				},
				Keywords: models.Keywords{Values: "The Best Rich and Moist Chocolate Cake"},
				Image: models.Image{
					Value: "https://plowingthroughlife.com/wp-content/uploads/2019/01/chocolate-cake-FI.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 box Devil's food cake mix", "1 small package instant chocolate pudding",
						"4 large eggs", "1 cup sour cream", "1/2 cup warm water",
						"1/2 cup vegetable oil", "1 1/2 cups semi-sweet chocolate chips",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 350°F. Grease a 9 x 13 inch pan.",
						"Mix all ingredients except chocolate chips at medium speed for 2 minutes.",
						"Gently fold in chips.",
						"I always use a 9 x 13 pan, by 2 - 9 inch round pans can also be used. Bake for 30 - 40 minutes until a toothpick inserted in the center comes out clean.",
					},
				},
				Name: "The Best Rich and Moist Chocolate Cake",
				NutritionSchema: models.NutritionSchema{
					Calories:      "361 kcal",
					Carbohydrates: "37 g",
					Cholesterol:   "53 mg",
					Fat:           "23 g",
					Fiber:         "2 g",
					Protein:       "5 g",
					SaturatedFat:  "13 g",
					Servings:      "1",
					Sodium:        "360 mg",
					Sugar:         "22 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 15},
				URL:      "https://plowingthroughlife.com/the-best-rich-and-moist-chocolate-cake/",
			},
		},
		{
			name: "popsugar.co.uk",
			in:   "https://www.popsugar.co.uk/food/cinnamon-butter-baked-carrot-recipe-46882533",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Side Dishes"},
				Cuisine:       models.Cuisine{Value: "North American"},
				DateModified:  "2023-11-04T11:05:08+00:00",
				DatePublished: "2023-11-02T13:33:04+00:00",
				Description: models.Description{
					Value: "This easy Christmas carrots recipe features fresh carrots smothered in cinnamon butter, baked until they're perfectly glazed.",
				},
				Keywords: models.Keywords{
					Values: "Christmas, Fall Food, Fall Recipes, Original Recipes, Side Dishes, Carrots, Autumn, Exclusive",
				},
				Image: models.Image{
					Value: "https://media1.popsugar-assets.com/files/thumbor/w2SbXKQCo_24S1wE0XWtTAJniOg/fit-in/2048xorig/filters:format_auto-!!-:strip_icc-!!-/2016/10/19/737/n/35573265/b5db29624d9c2514_cinnamon-butter-baked-carrots-5.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"15 carrots, peeled and stems removed", "75g of butter, softened",
						"100g of sugar", "1 teaspoon kosher salt", "1/2 teaspoon ground cinnamon",
						"80ml boiling water", "2 1/2 teaspoons orange juice",
						"Fresh parsley, chopped, optional",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat the oven to 180°C (or gas mark 4).",
						"Clean carrots and arrange in a 22 x 33 centimetre baking dish.",
						"Using an electric mixer, cream the butter, sugar, salt, and cinnamon. While the mixer is running, slowly add the boiling water and orange juice.",
						"Pour mixture over the carrots and cover the dish with aluminium foil. Bake for 90 minutes.",
						"Remove the foil and transfer carrots to a serving platter. Drizzle carrots with the melted cinnamon butter from the baking dish. Garnish with chopped parsley, if desired.",
					},
				},
				Name: "These Cinnamon-Butter Carrots Are Almost Too Easy to Make",
				NutritionSchema: models.NutritionSchema{
					Calories: "219 per serving",
				},
				Yield: models.Yield{Value: 6},
				URL:   "https://www.popsugar.co.uk/food/cinnamon-butter-baked-carrot-recipe-46882533",
			},
		},
		{
			name: "practicalselfreliance.com",
			in:   "https://practicalselfreliance.com/zucchini-relish/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT10M",
				DatePublished: "2022-08-08",
				Description: models.Description{
					Value: "Zucchini relish is a flavorful topping for summer grilling, and the perfect way to use up extra " +
						"zucchini from the garden.",
				},
				Image: models.Image{
					Value: "https://creativecanning.com/wp-content/uploads/2021/02/Zucchini-Relish-61-720x720.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 cups zucchini, diced (about 3 medium)",
						"1 cup onion, diced (about 1 medium)",
						"1 cup red bell pepper, diced (about 2 small or 1 large)",
						"2 Tablespoons Salt (pickling and canning salt, or kosher salt)",
						"1 3/4 cups sugar",
						"2 teaspoons celery seed (whole)",
						"1 teaspoon mustard seed (whole)",
						"1 cup cider vinegar (5% acidity)",
						"Pickle Crisp Granules (optional, helps veggies stay firm after canning)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Wash vegetables.",
						"Remove stem and blossom ends from zucchini and dice into 1/4 to 1/2 inch pieces. Measure 2 cups.",
						"Peel and dice onion. Measure 1 cup.",
						"Stem and seed peppers, then dice. Measure 1 cup.",
						"Important, don't skip this step! Combine diced vegetables in a large bowl and sprinkle salt over the top. " +
							"Stir gently to distribute the salt, then add water until vegetables are completely submerged. " +
							"Allow the vegetables to soak in the saltwater for 2 hours, then drain completely.",
						"Prepare a water bath canner (optional, only if canning).",
						"In a separate saucepan or stockpot, bring vinegar, sugar, and spices to a gentle simmer (180 degrees F). " +
							"Do not add salt, the salt is only used to soak veggies before draining.",
						"Add drained vegetables to the simmering vinegar/spices and gently simmer for 10 minutes.",
						"Pack hot relish into prepared half-pint or pint jars, leaving 1/2 inch headspace.",
						"If not canning, just seal jars and allow them to cool on the counter before storing in the refrigerator.",
						"If canning, de-bubble jars, wipe rims, and adjust headspace to ensure 1/2 inch. Seal with 2 part canning lids.",
						"Process in a water bath canner for 10 minutes, then turn off the heat. Allow the jars to sit in the canner " +
							"for another 5 minutes to cool slightly, then remove the jars to cool on a towel on the counter.",
						"Leave the jars undisturbed for 24 hours, then check seals. Store any unsealed jars in the refrigerator for " +
							"immediate use. Properly canned and sealed jars should maintain peak quality on the pantry shelf " +
							"for 12-18 months.",
					},
				},
				Name:     "Zucchini Relish Recipe for Canning",
				PrepTime: "PT2H10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://creativecanning.com/zucchini-relish/",
			},
		},
		{
			name: "pressureluckcooking.com",
			in:   "https://pressureluckcooking.com/spanish-omelette-scramble/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Breakfast"},
				CookTime:      "PT20M",
				Cuisine:       models.Cuisine{Value: "Spanish"},
				DatePublished: "2023-10-26",
				Description: models.Description{
					Value: "When I was in Valencia, Spain, I was taught how to make a Spanish omelette. It&#x27;s essentially pan-fried onions and potatoes coated with eggs, seasoned with the basics that are salt and pepper and then cooked in a skillet on both sides for a few minutes where it looks like a super thick, egg pancake. Truth be told, it is one of the most delicious egg dishes you&#x27;ll ever have. \n\nHOWEVER, as basic as it sounds, it can be a messy and somewhat cumbersomething to make with the flipping, removing a half-cooked, runny egg mound toa plate and then returning it to the pan to cook on the other side. At the class I took, three people who volunteered with flipping the half-cooked omelette had itspill onto the counter. So I decided to take the basics of a Spanish omelette, add some optional Spanish-favored cheese and meat and scramble itup to give you that flavor experience, but with a much simpler and fool-proof approach to making it.",
				},
				Keywords: models.Keywords{Values: "Spanish Omelette, Scrambled Eggs"},
				Image: models.Image{
					Value: "https://pressureluckcooking.com/wp-content/uploads/2023/10/Spanish-Omelette-IG-2-scaled-720x720.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1/2 cup extra virgin olive oil",
						"2 pounds Idaho or russet potatoes (about 2-3 of them), peeled, sliced into 1/4-inch pieces and then quartered",
						"1 large Spanish or yellow onion, sliced into thin strips",
						"2 teaspoons seasoned salt or regular salt (you can also use Tony Chachere’s Creole seasoning for a little spice)",
						"1 teaspoon black pepper (optional)", "12 large eggs, lightly beaten",
						"1 cup shredded Manchego cheese (optional, see Jeff's Tips)",
						"4 ounces jamón ibérico, Spanish chorizo (cured/charcuterie-style, NOT raw and the raw and crumbled kind in a tube), prosciutto, or deli ham, diced or sliced into small pieces (optional, see Jeff's Tips)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Add the olive oil to a 12-inch skillet/frying pan or large sauté pan (nonstick is definitely best here) and bring to medium-high heat.",
						"While the oil’s heating, gently crack the eggs into a bowl and whisk until smooth. Set aside.",
						"Once shimmering and heated, about 3 minutes, add the potato, onion, salt and pepper to the skillet and sauté for 15-20 minutes, stirring every 2-3 minutes until the onion and potato have softened and become golden-brown or even slightly charred in color (the potatoes should look like the color of breakfast hash browns/home fries).",
						"Reduce the heat to medium-low. Pour the whisked eggs into the pan with the veggies and oil, top with the cheese and meat (if using either or both) and let it all rest, undisturbed without stirring, for 30-60 seconds. Almost immediately, the egg around the perimeter of the skillet will begin to solidify.",
						"Gently slide a silicone spatula under the eggs. A firm, yet fluffy bottom layer will have formed, and now it's time to scramble the eggs! Using the spatula (and NEVER a fork), gently swirl the eggs, veggies, cheese and meat around in the pan, folding it all together (aka mixing them with a swirly flip) so the egg clings to the veggies and the cheese begins to melt. This should take about 1-2 minutes total, at most. (NOTE: A Spanish omelette is typically served a little runny. Since this is an homage to that, I like to cook my eggs until they're JUST underdone and slightly runny. However, you can absolutely continue to cook your eggs until they're more well-done, if you desire. The oil in the pan will prevent them from drying out.)",
						"Turn off the stove and remove the pan from the heat and serve immediately.",
					},
				},
				Name:     "Spanish Omelette Scramble",
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://pressureluckcooking.com/spanish-omelette-scramble/",
			},
		},
		{
			name: "primaledgehealth.com",
			in:   "https://www.primaledgehealth.com/slow-cooker-crack-chicken/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizers"},
				CookTime:      "PT360M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-03-09T09:12:00+00:00",
				Description: models.Description{
					Value: "Cheesy Ranch chicken, also known as Crack Chicken, is a low-carb slow cooker recipe with 5 " +
						"ingredients! Easily made into a dip or dinner, depending on how you serve it, this " +
						"ultra-creamy combo of chicken, cheese, and bacon is sure to please!",
				},
				Keywords: models.Keywords{
					Values: "Cheesy Ranch Chicken, Keto Crack Chicken, Slow Cooker Crack Chicken",
				},
				Image: models.Image{
					Value: "https://www.primaledgehealth.com/wp-content/uploads/2022/03/crack-chicken.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 pounds chicken thigh (or breast)",
						"12 ounces cream cheese",
						"1 cup cheddar cheese (shredded)",
						"8 ounces bacon (cooked and crumbled)",
						"2 1-ounce Ranch seasoning packets (or follow the DIY option below)",
						"2 teaspoons parsley (dried)",
						"2 teaspoons dill (dried)",
						"2 teaspoons chives (dried)",
						"2 teaspoons onion powder",
						"2 teaspoons garlic powder",
						"1 teaspoon salt",
						"½ teaspoon ground black pepper",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Add chicken to slow cooker and cover with seasoning. Top with cream cheese. Use the back of a spoon " +
							"to smear the cream cheese over the meat.",
						"Cover the slow cooker with a lid. Cook on LOW for 6-8 hours or HIGH 3-4 hours.",
						"Once cooking time finishes, shred the chicken by using two forks and pulling the meat apart. Remove " +
							"bones if needed. Stir well and thoroughly coat the chicken with sauce.",
						"Fry bacon in a pan over medium heat on the stove. Remove and chop or crumble into small pieces. Cover " +
							"the chicken with bacon and shredded cheese. Put the lid back on and continue cooking on HIGH for " +
							"15 minutes until cheese melts.",
						"Remove from heat. Mix bacon and cheese into the chicken or serve as is with them resting on top. Garnish " +
							"with a tablespoon or two of fresh parsley if desired!",
					},
				},
				Name: "Slow Cooker Crack Chicken (Cheesy Ranch Chicken)",
				NutritionSchema: models.NutritionSchema{
					Calories:       "603 kcal",
					Carbohydrates:  "3 g",
					Cholesterol:    "247 mg",
					Fat:            "39 g",
					Fiber:          "1 g",
					Protein:        "37 g",
					SaturatedFat:   "23 g",
					Servings:       "1",
					Sodium:         "689 mg",
					Sugar:          "1 g",
					TransFat:       "1 g",
					UnsaturatedFat: "30 g",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://www.primaledgehealth.com/slow-cooker-crack-chicken/",
			},
		},
		{
			name: "projectgezond.nl",
			in:   "https://www.projectgezond.nl/italiaanse-kiprollade-met-gremolata/",
			want: models.RecipeSchema{
				AtContext:    "https://schema.org",
				AtType:       models.SchemaType{Value: "Recipe"},
				Category:     models.Category{Value: "Diner"},
				DateModified: "2023-12-04T13:31:53+00:00",
				Description: models.Description{
					Value: "Ben je gek op kip pesto en valt een rollade ook altijd goed in de smaak? Maar vind je het op z’n tijd ook leuk en lekker om eens iets anders te proberen? Dan is dit recept echt iets voor jou! Een heerlijke zelfgemaakte kiprollade op Italiaanse wijze. De pesto vervang je door gremolata. Een …",
				},
				Image: models.Image{Value: "https://www.projectgezond.nl/wp-content/uploads/2023/09/Italiaanse-Kiprollade-Gremolata1600-810x1080.jpg"},
				Ingredients: models.Ingredients{
					Values: []string{
						"400 gr kipfilet", "40 gr prosciutto", "40 gr zuivelspread",
						"40 gr platte peterselie", "1 teentje knoflook", "1 citroen",
						"10 ml olijfolie",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Snijd de kipfilets doormidden en leg de kipfilets op een stuk huishoudfolie.",
						"Hoe platter de kipfilets zijn hoe beter. Daarom gaan we ze platslaan. Bedek de kipfilets met nog een stuk huishoudfolie en sla ze voorzichtig plat met bijvoorbeeld een koekenpan.",
						"Verwijder de bovenste laag huishoudfolie en leg de kipfilets tegen elkaar aan.",
						"Bedek de kip met de prosciutto en bedek weer met een stuk huishoudfolie. Draai het geheel nu om zodat de prosciutto onderop komt te liggen.",
						"Smeer de kipfilet in met de zuivelspread.",
						"Maak de gremolata door de peterselie, knoflook, sap van de citroen en de olijfolie met elkaar te blenden in een keukenmachine of hakmolen.",
						"Verdeel de gremolata nu ook over de kipfilet.",
						"Rol de kipfilet met de prosciutto aan de buitenkant nu zo strak mogelijk in de folie en leg voor 1 uur in de koelkast.",
						"Verwarm de oven voor op 150 graden.",
						"Verwijder het folie van de kiprollade en bindt hem eventueel op met rolladetouw om te voorkomen dat hij uit elkaar valt.",
						"Leg de kiprollade in een ovenschaal en bak in ongeveer 45 à 60 minuten gaar.",
						"Snijd de kiprollade in mooie plakken en serveer ze direct.",
					},
				},
				Name: "Italiaanse kiprollade met gremolata",
				NutritionSchema: models.NutritionSchema{
					Calories:      "125 kcal",
					Carbohydrates: "1 gram",
					Fat:           "5 gram",
					Fiber:         "1 gram",
					Protein:       "18 gram",
				},
				URL: "https://www.projectgezond.nl/italiaanse-kiprollade-met-gremolata/",
			},
		},
		{
			name: "przepisy.pl",
			in:   "https://www.przepisy.pl/przepis/placki-ziemniaczane",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Placki ziemniaczane"},
				DatePublished: "2010-11-12T14:29:28.000Z",
				Description: models.Description{
					Value: "Jedni zajadają się plackami ziemniaczanymi posypanymi cukrem, inni z ochotą dodają do nich gęstą, kwaśną śmietanę, a jeszcze inni najbardziej cenią przepisy na placki ziemniaczane polane mięsnym gulaszem. Niewątpliwie najlepsze placki ziemniaczane to te usmażone na złoto, których chrupiąca, lekko przypieczona skórka skrywa miękki i delikatny środek. Niezależnie od tego, jakie dodatki należą do Twoich kulinarnych faworytów, poznaj naszą wersję przyrządzania tych przysmaków. Rumiane placki ziemniaczane to coś zdecydowanie więcej niż tylko potrawa charakterystyczna dla barów mlecznych.\r\nDla kogo? \r\nPrzepis na placki ziemniaczane przypadnie do gustu smakoszom ziemniaczanych potraw. Do stołu chętnie zasiądą również zwolennicy tradycyjnych dań. Kuchnia jak u mamy? Z plackami ziemniaczanymi według naszego przepisu ten efekt uzyskasz bez problemów. Masz ochotę na odrobinę nowości? Do chrupiących placuszków dodaj łososia albo krewetki – nie pożałujesz!\r\nNa jaką okazję? \r\nPlacki ziemniaczane to idealna propozycja na poskromienie większego głodu. Świetnie sprawdzają się w jesiennych i zimowych miesiącach jako rozgrzewający obiad. Jak zrobić placki ziemniaczane w takiej odsłonie? Sos grzybowy, żurawina i oscypek to sezonowe, obowiązkowe dodatki.\r\nCzy wiesz, że? \r\nZiemniaków nie musisz trzeć na drobnej tarce. Wystarczy, że zetrzesz je na grubych oczkach. Jak zrobić placki ziemniaczane tak, aby nie rozpadły się podczas smażenia, do masy ziemniaczanej dodaj jajko i odrobinę mąki. Ich postrzępione brzegi nie tylko pięknie się prezentują, ale też genialnie smakują i chrupią.\r\nDla urozmaicenia: \r\nPamiętaj, że przepis na placki ziemniaczane możesz dowolnie urozmaicić, serwując je z ulubionymi składnikami. Jak zrobić placki ziemniaczane, które zaskoczą wszystkich? Jeśli lubisz eksperymentować, do masy (oprócz czosnku) dorzuć ulubione zioła lub drobno pokrojone warzywa. Dzięki szpinakowi, papryce chili lub cebuli będą aromatyczne i kolorowe.",
				},
				Keywords: models.Keywords{
					Values: "Na co dzień, Ziemniaki, Warzywa, Jajka, Łagodne, Bez mięsa",
				},
				Image: models.Image{
					Value: "https://s3.przepisy.pl/przepisy3ii/img/variants/800x0/placki-ziemniaczane.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 kilogram ziemniaki",
						"1 sztuka cebula",
						"2 sztuki jajka",
						"1 sztuka Przyprawa w Mini kostkach Czosnek Knorr",
						"1 szczypta gałka muszkatołowa",
						"1 szczypta sól",
						"3 łyżki mąka",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Obierz ziemniaki, zetrzyj na tarce. Odsącz masę przez sito. Zetrzyj cebulę na tarce.",
						"Dodaj do ziemniaków cebulę, jajka, gałkę muszkatołową oraz mini kostkę Knorr.",
						"Wymieszaj wszystko dobrze, dodaj mąkę, aby nadać masie odpowiednią konsystencję.",
						"Rozgrzej na patelni olej, nakładaj masę łyżką. Smaż placki z obu stron na złoty brąz i od razu podawaj.",
					},
				},
				Name:     "Placki ziemniaczane",
				PrepTime: "PT40M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://www.przepisy.pl/przepis/placki-ziemniaczane",
			},
		},
		{
			name: "purelypope.com",
			in:   "https://purelypope.com/sweet-chili-brussel-sprouts/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				DateModified:  "2020-09-17T22:05:26+00:00",
				DatePublished: "2020-05-21T00:35:12+00:00",
				Name:          "Sweet Chili Brussel Sprouts",
				Yield:         models.Yield{Value: 4},
				PrepTime:      "PT10M",
				CookTime:      "PT32M",
				Image: models.Image{
					Value: "https://i0.wp.com/purelypope.com/wp-content/uploads/2020/05/IMG_5412-1-scaled.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 cups brussel sprouts, stems removed & cut in half",
						"2 tbsp coconut aminos",
						"1 tbsp sriracha",
						"1/2 tbsp maple syrup",
						"1 tsp sesame oil",
						"Everything bagel seasoning or sesame seeds, to top",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 350 degrees.",
						"Whisk the sauce (coconut aminos, sriracha, maple syrup & sesame oil) together in a large bowl.",
						"Toss in brussel sprouts and coat mixture evenly over the brussels.",
						"Roast for 30 minutes.",
						"Turn oven to broil for 2-3 minutes to crisp (watch carefully to not burn.)",
						"Top with everything or sesame seeds.",
					},
				},
				URL: "https://purelypope.com/sweet-chili-brussel-sprouts/",
			},
		},
		{
			name: "purplecarrot.com",
			in: "https://www.purplecarrot.com/recipe/gnocchi-al-pesto-with-charred-green-beans-lemon-zucchini-bc225f0b-" +
				"1985-4d94-b05b-a78de295b2da?plan=chefs_choice",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				DatePublished: "2020-05-13T06:59:57.162-04:00",
				Image: models.Image{
					Value: "https://images.purplecarrot.com/uploads/product/image/2017/_1400_700_GnocchiAlPestowithCharredGreenBeans_LemonZucchini_WEBHERO-5d97b356980badd112987864879c4f71.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 zucchini, trimmed and peeled lengthwise into ribbons",
						"1 lemon, half juiced, half cut into wedges (divided)",
						"1 tsp Aleppo pepper flakes",
						"6 oz green beans, cut in half",
						"10 oz fresh gnocchi",
						"¼ cup vegan basil pesto",
						"1 tbsp + 2 tsp olive oil*",
						"Salt and pepper*",
						"*Not included",
						"Ingredients are listed for 2 servings. If you're making 4 servings, please double your ingredients by using both meal kit bags provided.",
						"For full ingredient list, see Nutrition.",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"1 - Prepare the zucchini.",
						"2 - Char the green beans.",
						"3 - Cook the gnocchi.",
						"4 - Finish the gnocchi.",
						"5 - Serve.",
					},
				},
				Name: "Gnocchi Al Pesto with Charred Green Beans & Lemon Zucchini",
				NutritionSchema: models.NutritionSchema{
					Calories: "540 cal",
					Fat:      "22.0 g",
				},
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://www.purplecarrot.com/recipe/gnocchi-al-pesto-with-charred-green-beans-lemon-zucchini-bc225f0b-1985-4d94-b05b-a78de295b2da?plan=chefs_choice",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
