package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_R(t *testing.T) {
	testcases := []testcase{
		{
			name: "rachlmansfield.com",
			in:   "https://rachlmansfield.com/delicious-crispy-rice-salad-gluten-free/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT15M",
				DatePublished: "2022-04-03",
				Description: models.Description{
					Value: "This Crispy Rice Salad is such an easy and flavorful recipe to make for lunch or dinner. This salad " +
						"" +
						"is vegan, gluten-free and craveable.",
				},
				Image: models.Image{
					Value: "https://rachlmansfield.com/wp-content/uploads/2022/03/IMG_8796-scaled.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1/2 cup basmati rice, uncooked",
						"1/2 cup broth (code RACHL)",
						"1/4 cup filtered water",
						"1 large head of tuscan kale",
						"Olive oil to massage kale",
						"2 garlic cloves, chopped",
						"1/4 cup chopped sweet onion",
						"1 edamame beans (not in the shells)",
						"1/4 cup scallions, chopped",
						"1/3 cup cherry peppers, sliced",
						"1/2 cup kimchi, chopped",
						"1/3 cup roasted unsalted peanuts",
						"1.5 tablespoons coconut aminos",
						"1 tablespoon sesame oil",
						"Salt and pepper to taste",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Prepare the rice per package but do a mix of the broth and water (optimal flavor!)",
						"While the rice cooks, remove stems from the kale and chop the kale",
						"Add kale to a large mixing bowl and massage with your hands with olive oil to cut the bitterness",
						"Warm a large skillet with oil, garlic and onion and cook for 3-5 minutes or until fragrant",
						"Add in the rice and press down to form a large &#8220;rice pancake&#8221; of sorts",
						"Cook on medium heat for about 8 minutes then start to stir it to crisp the other side of the rice (do " +
							"not cover the rice or it won&#8217;t crisp!)",
						"Remove rice from pan once crisped and add to mixing bowl with the kale and add in the edamame, scallions, " +
							"peppers, kimchi, peanuts and mix",
						"Dress with coconut aminos and sesame oil and salt and pepper and enjoy!",
						"You can also add some cooked salmon, chicken or any additional protein if you&#8217;d like",
					},
				},
				Name:     "Delicious Crispy Rice Salad (gluten-free)",
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://rachlmansfield.com/delicious-crispy-rice-salad-gluten-free/",
			},
		},
		{
			name: "rainbowplantlife.com",
			in:   "https://rainbowplantlife.com/livornese-stewed-beans/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				CookTime:      "PT60M",
				Cuisine:       models.Cuisine{Value: "Italian"},
				DatePublished: "2022-01-04T04:56:15+00:00",
				Description: models.Description{
					Value: "These Tuscan Stewed Beans are the ultimate rustic Italian comfort food! Made with simple pantry-friendly ingredients like onions, garlic, tomato paste and white beans, but big on gourmet Italian flavor. It&#039;s cozy and indulgent yet wholesome, vegan, and gluten-free.",
				},
				Keywords: models.Keywords{Values: "italian beans, italian white bean stew, stewed beans, tuscan beans"},
				Image: models.Image{
					Value: "https://rainbowplantlife.com/wp-content/uploads/2022/01/Livornese-stewed-beans-5-of-5.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1/4 cup (56 mL) extra virgin olive oil",
						"1 medium yellow onion, (chopped)",
						"2 medium or large carrots, (peeled and finely chopped)",
						"2 celery ribs, (diced)",
						"4 garlic cloves, (finely chopped)",
						"½ tsp red pepper flakes",
						"1/4 cup (4g) flat-leaf parsley leaves and tender stems, (minced)",
						"1 tablespoon minced fresh sage",
						"4 1/2 tablespoons (67g) tomato paste ((in a tube, not a can)*)",
						"¾ cup (180 mL) dry white wine**",
						"1 28-ounce (800g) can whole peeled tomatoes, ( crushed by hand)",
						"1 teaspoon kosher salt, (plus more to taste)",
						"Freshly cracked black pepper",
						"1 bay leaf",
						"1 1/2 cups (360 mL) vegetable broth, plus more as desired",
						"2 (15-ounce/425g) cans cannellini beans, drained and rinsed",
						"½ cup (8g) fresh basil, (slivered***)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Heat the olive oil in a Dutch oven over medium heat. Once the oil is hot, add the onion, and season with " +
							"a pinch or two of salt and pepper. Cook for 7 to 8 minutes, until golden, stirring occasionally. " +
							"Add in the carrot, celery, and garlic, with another pinch of salt and cook for 3 to 4 minutes. " +
							"Add the red pepper flakes, parsley, and sage and cook until fragrant, about 1 minute.",
						"Add the tomato paste and cook, stirring almost continuously, for 1 to 2 minutes, until it&#39;s a bit " +
							"darker in color.",
						"Pour the white wine in and deglaze the pan, scraping up any browned bits stuck to the bottom of the pot. " +
							"Allow wine to simmer rapidly for 3 minutes, or until mostly evaporated and it no longer smells like " +
							"wine, stirring often.",
						"Add tomatoes along with their juices, bay leaf, 1 teaspoon kosher salt, and several cracks of black pepper. " +
							"Cook at a rapid simmer, stirring fairly often, until the tomatoes are fully broken down and most of " +
							"the liquid has evaporated, 12 to 13 minutes.",
						"Add the veggie broth and 2 cans of beans. Reduce the heat to low, cover the pan, and maintain a decent " +
							"simmer for 30 minutes, stirring once in a while. If you want the stew to be thicker, towards the end " +
							"of cooking, use the back of a wooden spoon or a spatula to gently smash a small portion of the beans.",
						"Taste, adding a pinch of sugar if needed (if your tomatoes are good-quality, it should not be necessary). " +
							"Remove the bay leaf. Finish with chopped basil. Season to taste, adding salt and pepper as needed.",
					},
				},
				Name: "Tuscan Stewed Beans",
				NutritionSchema: models.NutritionSchema{
					Calories:       "472 kcal",
					Carbohydrates:  "59 g",
					Fat:            "16 g",
					Fiber:          "14 g",
					Protein:        "18 g",
					SaturatedFat:   "2 g",
					Servings:       "1",
					Sodium:         "1117 mg",
					Sugar:          "7 g",
					UnsaturatedFat: "13 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://rainbowplantlife.com/livornese-stewed-beans/",
			},
		},
		{
			name: "realsimple.com",
			in:   "https://www.realsimple.com/food-recipes/browse-all-recipes/sheet-pan-chicken-and-sweet-potatoes",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				DateModified:  "2018-07-05T14:07:52.000-04:00",
				DatePublished: "2016-12-07T11:48:40.000-05:00",
				Description: models.Description{
					Value: "Get the recipe for Sheet Pan Chicken and Sweet Potatoes.",
				},
				Image: models.Image{
					Value: "https://www.realsimple.com/thmb/8gMeQAdUxCc8bTx33CFY4cdH7PU=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/sheet-pan-chicken-sweet-potatoes_0-d610f954ea1e46179f961d536abc8f32.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"4 bone-in, skin-on chicken leg quarters (about 2 lb.)",
						"2 medium sweet potatoes, peeled and cut into 1-in. wedges",
						"1 teaspoon chopped fresh sage",
						"0.75 teaspoon kosher salt, plus more to taste",
						"0.5 teaspoon black pepper, plus more to taste",
						"3 tablespoons olive oil, divided",
						"3 slices bacon",
						"3 cups baby watercress",
						"1 tablespoon fresh lemon juice",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 450°F. Arrange the chicken and sweet potatoes side by side in a single layer on a " +
							"large rimmed baking sheet. Season with the sage, salt, and pepper and drizzle with 2 " +
							"tablespoons of the oil, tossing to coat. Lay the bacon on top of the sweet potatoes.",
						"Roast until a meat thermometer inserted into the thickest portion of a thigh registers 165°F, 20 to " +
							"25 minutes.",
						"Meanwhile, toss together the watercress, lemon juice, and the remaining 1 tablespoon of olive oil and " +
							"season to taste with salt and pepper.",
						"Serve the chicken with the sweet potatoes and salad, with the bacon crumbled over the top.",
					},
				},
				Name: "Sheet Pan Chicken and Sweet Potatoes",
				NutritionSchema: models.NutritionSchema{
					Calories:       "533 kcal",
					Carbohydrates:  "22 g",
					Cholesterol:    "181 mg",
					Fat:            "34 g",
					Protein:        "34 g",
					SaturatedFat:   "9 g",
					Sodium:         "670 mg",
					Sugar:          "5 g",
					UnsaturatedFat: "0 g",
				},
				URL:   "https://www.realsimple.com/food-recipes/browse-all-recipes/sheet-pan-chicken-and-sweet-potatoes",
				Yield: models.Yield{Value: 1},
			},
		},
		{
			name: "receitasnestle.com.br",
			in:   "https://www.receitasnestle.com.br/receitas/pave-de-pessego",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Sobremesa"},
				CookTime:      "PT0M",
				Cuisine:       models.Cuisine{Value: "Local"},
				DatePublished: "2022-10-07T02:44:16-0300",
				Description: models.Description{
					Value: "Sobremesa com creme de Leite MOÇA, biscoito champanhe, pêssegos em calda e Chocolate GALAK",
				},
				Keywords: models.Keywords{
					Values: "Family Meals,Café da Tarde,Sobremesa,Local,New Years,Natal,Sem peixe,Sem crustáceos,Sem carne de porco,Pescador,receita com frutos do mar,receita com peixe,receita sem crustaceos,receita sem carne de porco,receita sem peixe,MyMenuPlan,leite moça,pavê,De outros,Fruta,Frio / montagem",
				},
				Image: models.Image{Value: "/sites/default/files/srh_recipes/38bca0566eef6fb1d272973f2ad9593a.jpg"},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 Lata Leite MOÇA® (lata ou caixinha) 395g",
						"2 xícaras (chá) de Leite Líquido NINHO® Forti+ Integral", "2 gemas",
						"2 colheres (sopa) de amido de milho", "1 pacote de biscoito champanhe",
						"meia xícara (chá) da calda do pêssego",
						"1 lata de pêssego em calda picado", "100 g de raspas de Chocolate GALAK®",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Em uma panela, misture o Leite MOÇA, o Leite NINHO, as gemas e o amido de milho e leve ao fogo, mexendo sempre até engrossar. Deixe esfriar e reserve.",
						"Passe rapidamente parte dos biscoitos na calda do pêssego, acomodando-os no fundo de um refratário retangular (20 x 30 cm).",
						"Distribua metade do Creme e metade dos pêssegos picados.",
						"Repita a camada de biscoitos e coloque o restante do Creme e dos pêssegos.",
						"Decore com as raspas de Chocolate GALAK e com fatias de pêssegos.",
						"Leve à geladeira até o momento de servir.",
					},
				},
				Name: "Pavê de Pêssego",
				NutritionSchema: models.NutritionSchema{
					Calories:     "226",
					Fat:          "10",
					Protein:      "5",
					SaturatedFat: "1",
					Sodium:       "73",
					Sugar:        "11",
				},
				PrepTime: "PT0M",
				Yield:    models.Yield{Value: 12},
				URL:      "https://www.receitasnestle.com.br/receitas/pave-de-pessego",
			},
		},
		{
			name: "recettes.qc.ca",
			in:   "https://www.recettes.qc.ca/recettes/recette/yakisoba-nouille-sautees-a-la-japonaise",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Pâtes alimentaires"},
				CookTime:      "PT20M",
				DatePublished: "2015-07-28T21:44:00-04:00",
				Description: models.Description{
					Value: "Recette de Yakisoba (nouilles sautées à la japonaise)",
				},
				Keywords: models.Keywords{Values: "pates alimentaires"},
				Image: models.Image{
					Value: "https://m1.quebecormedia.com/emp/rqc_prod/recettes_du_quebec-_-45fe466bb6b64f799cc2ce9ab8db72f66d46ef08-_-yakisoba.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"250 g nouilles soba ou à ramen",
						"300 g porc haché",
						"1 cuillère à table huile de sésame",
						"2 cuillères à table huile de pépins de raisin",
						"1 moyen oignon coupé en 8",
						"1 gousse d'ail",
						"500 g chou coupé en fines lanières",
						"1 poivron coupé en fines tranches",
						"2 cuillères à table gingembre rouge mariné (beni-shoga)",
						"2 cuillères à table algues ao-nori séchées en poudre",
						"1 cuillère à table sucre",
						"60 mL mirin",
						"2 cuillères à table saké",
						"60 mL sauce soja japonaise",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Faites cuire les nouilles dans de l’eau bouillante en les gardant fermes. Égouttez-les.",
						"Préparez la sauce : mettez les ingrédients de la sauce dans une petite casserole et chauffez jusqu’à ce que le sucre soit dissous.",
						"Faites chauffer l’huile de sésame et 1 cuillerée d’huile de pépins de raisin dans un wok. Faites-y revenir le porc jusqu’à ce qu’il soit légèrement doré. Réservez.",
						"Ajoutez le reste de l’huile dans le wok et faites sauter l’oignon et l’ail jusqu’à ce que l’oignon blondisse. Ajoutez le chou et le poivron. Cuisez jusqu’à ce qu’ils soient tendres. Ajoutez les nouilles, le porc, le gingembre mariné et la sauce. Mélangez et réchauffez. Servez et parsemez avec les ao-nori.",
					},
				},
				Name:     "Yakisoba (nouilles sautées à la japonaise)",
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.recettes.qc.ca/recettes/recette/yakisoba-nouille-sautees-a-la-japonaise",
			},
		},
		{
			name: "reciperunner.com",
			in:   "https://reciperunner.com/cranberry-apple-sauce/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Side Dishes"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-11-13",
				Description: models.Description{
					Value: "Cranberry apple sauce is sweet and tart and a must make side dish for your holiday table!",
				},
				Keywords: models.Keywords{
					Values: "cranberry apple sauce, cranberry sauce, cranberry applesauce, apple cranberry sauce",
				},
				Image: models.Image{
					Value: "https://reciperunner.com/wp-content/uploads/2023/11/cranberry-apple-sauce-3-scaled-720x720.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"12 ounces fresh or frozen cranberries",
						"1 cup peeled and diced Honeycrisp apples", "1/2 cup maple syrup",
						"1/2 cup apple cider", "1/2 teaspoon ground cinnamon",
						"1/4 teaspoon ground allspice", "1/8 teaspoon salt",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Combine all of the ingredients in a saucepan and bring to a boil over medium-high heat. Once boiling, reduce the heat to a simmer.",
						"Stir occasionally and cook until the cranberries have broken down and the apples have softened, but still have a little bite to them, about 10-15 minutes. The sauce should be compote-like in consistency.",
						"Remove the cranberry apple sauce from the heat and let it cool completely before transferring it to a bowl to chill in the refrigerator until ready to eat.",
					},
				},
				Name: "Cranberry Apple Sauce",
				NutritionSchema: models.NutritionSchema{
					Calories:       "117 calories",
					Carbohydrates:  "30 grams carbohydrates",
					Cholesterol:    "0 milligrams cholesterol",
					Fat:            "0 grams fat",
					Fiber:          "3 grams fiber",
					Protein:        "0 grams protein",
					SaturatedFat:   "0 grams saturated fat",
					Servings:       "1",
					Sodium:         "48 milligrams sodium",
					Sugar:          "23 grams sugar",
					TransFat:       "0 grams trans fat",
					UnsaturatedFat: "0 grams unsaturated fat",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://reciperunner.com/cranberry-apple-sauce/",
			},
		},
		{
			name: "recipetineats.com",
			in:   "https://www.recipetineats.com/chicken-sharwama-middle-eastern/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Chicken"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "Arabic"},
				DatePublished: "2022-02-06T06:47:00+00:00",
				Description: models.Description{
					Value: "Recipe video above. The smell when this is cooking is outrageous! The marinade is very quick to " +
						"prepare and the chicken can be frozen in the marinade, then defrosted prior to cooking. " +
						"Best cooked on the outdoor grill / BBQ, but I usually make it on the stove. Serve with " +
						"Yogurt Sauce (provided) or the Tahini sauce in this recipe. Add a simple salad and " +
						"flatbread laid out on a large platter, then let everyone make their own wraps!",
				},
				Keywords: models.Keywords{Values: "Chicken Shawarma, shawarma"},
				Image: models.Image{
					Value: "https://www.recipetineats.com/wp-content/uploads/2022/02/Chicken-Shawarma-Wrap_2-SQ.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 kg / 2 lb chicken thigh fillets (, skinless and boneless (Note 3))",
						"1 large garlic clove (, minced (or 2 small cloves))",
						"1 tbsp ground coriander",
						"1 tbsp ground cumin",
						"1 tbsp ground cardamon",
						"1 tsp ground cayenne pepper ((reduce to 1/2 tsp to make it not spicy))",
						"2 tsp smoked paprika",
						"2 tsp salt",
						"Black pepper",
						"2 tbsp lemon juice",
						"3 tbsp olive oil",
						"1 cup Greek yoghurt",
						"1 clove garlic (, crushed)",
						"1 tsp cumin",
						"Squeeze of lemon juice",
						"Salt and pepper",
						"4 - 5 flatbreads ((Lebanese or pita bread orhomemade soft flatbreads))",
						"Sliced lettuce ((cos or iceberg))",
						"Tomato slices",
						"Red onion (, finely sliced)",
						"Cheese (, shredded (optional))",
						"Hot sauce of choice ((optional))",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Marinade chicken - Combine the marinade ingredients in a large ziplock bag. Add the chicken, seal, the " +
							"massage from the outside with your hands to make sure each piece is coated. Marinate 24 " +
							"hours (minimum 3 hours).",
						"Yogurt Sauce - Combine the Yogurt Sauce ingredients in a bowl and mix. Cover and put in the fridge " +
							"until required (it will last for 3 days in the fridge).",
						"Preheat stove or BBQ - Heat a large non-stick skillet with 1 tablespoon over medium high heat, or " +
							"lightly brush a BBQ hotplate/grills with oil and heat to medium high. (See notes for baking)",
						"Cook chicken - Place chicken in the skillet or on the grill and cook the first side for 4 to 5 minutes " +
							"until nicely charred. Turn and cook the other side for 3 to 4 minutes (the 2nd side takes less time).",
						"Rest - Remove chicken from the grill and cover loosely with foil. Set aside to rest for 5 minutes.",
						"Slice chicken and pile onto platter alongside flatbreads, Salad and the Yoghurt Sauce (or dairy free " +
							"Tahini sauce from this recipe).",
						"To make a wrap, get a piece of flatbread and smear with Yoghurt Sauce. Top with a bit of lettuce and " +
							"tomato and Chicken Shawarma. Roll up and enjoy!",
					},
				},
				Name: "Chicken Shawarma (Middle Eastern)",
				NutritionSchema: models.NutritionSchema{
					Calories:       "275 kcal",
					Carbohydrates:  "1.1 g",
					Cholesterol:    "140 mg",
					Fat:            "16.2 g",
					Protein:        "32.9 g",
					SaturatedFat:   "3.2 g",
					Servings:       "183 g",
					Sodium:         "918 mg",
					UnsaturatedFat: "13 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.recipetineats.com/chicken-sharwama-middle-eastern/",
			},
		},
		{
			name: "redhousespice.com",
			in:   "https://redhousespice.com/pork-fried-rice/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Course"},
				CookTime:      "PT8M",
				Cuisine:       models.Cuisine{Value: "Chinese"},
				DatePublished: "2022-03-26T17:48:51+00:00",
				Description: models.Description{
					Value: "Delicious pork fried rice made in less than 20 minutes. Enjoy the mix of fluffy rice, tender pork " +
						"and crunchy veggies coated with umami-filled seasoning.",
				},
				Keywords: models.Keywords{Values: "Pork, Rice, Stir-fry"},
				Image: models.Image{
					Value: "https://redhousespice.com/wp-content/uploads/2022/03/chinese-pork-fried-rice-1-scaled.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 tbsp oyster sauce (see note 1 for substitutes)",
						"1 tbsp light soy sauce",
						"½ tbsp dark soy sauce",
						"⅛ tsp ground white pepper",
						"2 tbsp neutral cooking oil (divided)",
						"2 eggs, lightly beaten",
						"1 cup minced pork (about 225g/8oz)",
						"1 small onion, diced",
						"1 tbsp minced garlic",
						"1 tsp minced ginger",
						"½ cup peas (about 50g/1.8oz)",
						"½ cup carrot, diced (about 50g/1.8oz)",
						"3 cups cold cooked white rice (about 400g/14oz (see note 2))",
						"Scallions, finely chopped (for garnishing)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a small bowl, mix oyster sauce, light soy sauce, dark soy sauce and white pepper. Set aside.",
						"Heat an empty, well-seasoned wok over high heat until smoking hot. Add 1 tablespoon of oil (see note 3 " +
							"if using other cookware). Swirl to coat a bigger perimeter.",
						"Pour in the beaten egg. Once it begins to set at the bottom, stir to help the running part flow. Use a " +
							"spatula to scramble so that it turns into small pieces. Transfer out and set aside.",
						"Pour the remaining 1 tablespoon of oil into the wok. Add minced pork. Spread and flatten it to ensure " +
							"maximum contact with the wok. Wait for the bottom part to get lightly browned. Then flip and stir " +
							"to fry it thoroughly.",
						"Once the pork loses its pink colour, add onion, garlic and ginger. Fry until the onion becomes a little " +
							"transparent.",
						"Stir in peas and carrots. Retain the high heat to fry for 30 seconds or so. Add rice and return the egg " +
							"to the wok. Cook for a further 30-40 seconds.",
						"Pour the sauce mixture over. Toss and stir constantly to ensure an even coating. Once all the ingredients " +
							"are piping hot, turn off the heat. Sprinkle scallions over and give everything a final toss.",
					},
				},
				Name: "Pork Fried Rice (猪肉炒饭)",
				NutritionSchema: models.NutritionSchema{
					Calories: "455 kcal",
					Servings: "1",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 3},
				URL:      "https://redhousespice.com/pork-fried-rice/",
			},
		},
		{
			name: "reishunger.de",
			in:   "https://www.reishunger.de/rezepte/rezept/440/chicken-tikka-masala",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Grillen"},
				Cuisine:       models.Cuisine{Value: "Indisch"},
				DatePublished: "2015-07-25T21:55:34+00:00",
				Description: models.Description{
					Value: "Hier kommt der Briten liebstes Gericht: Chicken Tikka Masala zusammen mit einem indischen Biryani " +
						"Reis und Minzjoghurt! Geflügel ist gerade im Sommer eine tolle Sache und dieses Gericht lässt " +
						"sich problemlos für mehrere Personen kochen!",
				},
				Image: models.Image{
					Value: "https://cdn.reishunger.com/chicken-tikka-masala.jpg?quality=85",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"200g Basmati Reis Pusa",
						"2 TL Indian Tikka Marinade",
						"2 EL Bio Cashewkerne",
						"2 EL Rosinen",
						"4-5 Minzblätter",
						"4-5 Minzblätter",
						"2 EL Tomatenmark",
						"300g Hähnchenbrustfilet",
						"4 Hähnchenkeulen",
						"100g Joghurt",
						"40g Butter",
						"4 Knoblauchzehen",
						"3 Schalotten",
						"2 Chillischoten",
						"1 Becher Creme Fraiche",
						"1 kleines Stück Ingwer",
						"1 Limette",
						"0.25 Gurke",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Den Knoblauch und den Ingwer mit einer feinen Küchenreibe in eine große Schüssel reiben. \nDie zwei Schalotten in grobe Scheiben schneiden und die Chilischoten in der Mitte teilen, nun kann man die Kerne entfernen.",
						"Die gehackten Zutaten kommen zu den Übrigen in die Schüssel. Anschließend den Saft einer Limette, 250 Gramm Joghurt, einen guten Schuss Olivenöl sowie die Gewürzmischung samt Salz und Pfeffer dazu. \nAlles zusammen mit dem Hühnchen vermengen und dieses mindestens eine halbe Stunde in den Kühlschrank stellen.",
						"Als Nächstes die marinierten Stücke der Brust aus der Schüssel nehmen und diese auf Spieße stecken. Den Grill auf direkte, starke Hitze vorheizen und die Spieße direkt über die Glut legen, um Sie von allen Seiten anzugrillen. \nSind die Spieße scharf angegrillt (Ihr könnt die Stücke auch in der Pfanne braten), diese vom Grill nehmen und sie kurz bei Seite legen.",
						"In der Zwischenzeit kann man eine Pfanne auf den Grill stellen und dort die Butter, eine grob gehackte Schalotte, das Tomatenmark sowie den Rest der Marinade hineingeben. \nAlles zusammen wird unter direkter starker Hitze angebraten und anschließend mit 500 ml Wasser aufgegossen.",
						"Ist alles etwas eingekocht gebt Ihr die Hähnchenstücke und Euren Becher Creme Fraiche dazu. \nNun wird die Pfanne in den indirekten Bereich gestellt und alles wird etwa 20-25 Minuten bei geschlossenem Deckel geschmort.",
						"Der Reis wird zusammen mit der Gewürzpaste in einen Topf gegeben, anschließend kommt die 1,5 Fache Menge an Wasser dazu (im Zweifel mit Tassen abmessen). \nNun lasst Ihr den Reis zusammen mit dem Wasser und der Gewürzmischung etwa 10 Minuten quellen ehe Ihr Ihn auf den Herd oder die Seitenkochplatte stellt.",
						"Den Reis kurz aufkochen, die Flamme zurückstellen und den Reis unter gelegentlichem Rühren garen lassen. Wenn das Wasser verschwunden ist, ist er fertig.\nJetzt kann man die Rosinen und die Cashewkerne dazu geben und den Reis eventuell mit noch etwas Gewürzpaste und Salz abschmecken.",
						"Für den Minzjoghurt die übrigen 250 Gramm Joghurt nehmen und in eine Schüssel geben. \nDie Gurke in der Mitte durchschneiden und das Kerngehäuse entfernen. Anschließend in feine Würfel schneiden. Die Minze in feine Streifen schneiden und zusammen mit den Gurkenwürfeln zum Joghurt geben. Mit Salz und Pfeffer abschmecken.",
					},
				},
				Name:  "Chicken Tikka Masala",
				Yield: models.Yield{Value: 3},
				URL:   "https://www.reishunger.de/rezepte/rezept/440/chicken-tikka-masala",
			},
		},
		{
			name: "rezeptwelt.de",
			in:   "https://www.rezeptwelt.de/vorspeisensalate-rezepte/haehnchen-nuggets/y3duba6e-e2d56-608317-cfcd2-vjez4wd6",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Cuisine:       models.Cuisine{Value: "Europäisch, Spanisch"},
				Name:          "Hähnchen-Nuggets",
				DateModified:  "2014-05-28",
				DatePublished: "2014-05-26",
				Description: models.Description{
					Value: "Hähnchen-Nuggets, ein Rezept der Kategorie Vorspeisen/Salate. Mehr Thermomix ® Rezepte auf www.rezeptwelt.de",
				},
				Image: models.Image{
					Value: "https://de.rc-cdn.community.thermomix.com/recipeimage/y3duba6e-e2d56-608317-cfcd2-vjez4wd6/57e6b699-53cb-4229-9398-1eb1ec70245e/main/haehnchen-nuggets.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"200 g Hähnchenbrust, ohne haut, in Stücken",
						"1/2 TL Salz",
						"1/2 TL Knoblauch, granuliert, oder eine Knoblauchzehe",
						"2 Scheiben Toastbrot (ohne Kruste)",
						"60 g Frischkäse",
						"60 g Milch",
						"1 Ei, mit 50 g Wasser verquirlt",
						"100 g Paniermehl",
						"Öl zum Frittieren",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Hähnchenfleisch, Salz, Knoblauch in den Mixtopf geben und 5 Sek./Stufe 7 zerkleinern.",
						"Brot in Stücken mit Frischkäse und Milch zugeben und 10 Sek./Stufe 7 vermischen.",
						"Fleischmischung aus dem Mixtopf nehmen, walnussgroße Bällchen formen und leicht mit dem Boden des Messbechers flach drücken. Jedes Nugget zuerst in Ei und dann in Paniermehl wenden. Öl in einer tiefen Pfanne erhitzen. Nuggets darin goldbraun frittieren und auf Küchenkrepp abtropfen lassen.",
					},
				},
				Keywords: models.Keywords{Values: "einfach,europaisch,spanisch,vorspeise,braten,snack,"},
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 20},
				URL:      "https://www.rezeptwelt.de/vorspeisensalate-rezepte/haehnchen-nuggets/y3duba6e-e2d56-608317-cfcd2-vjez4wd6",
			},
		},
		{
			name: "ricetta.it",
			in:   "https://ricetta.it/pan-d-arancio",
			want: models.RecipeSchema{
				AtContext: "https://schema.org",
				AtType:    models.SchemaType{Value: "Recipe"},
				Category:  models.Category{Value: "Torte"},
				Cuisine:   models.Cuisine{Value: "Italiana"},
				Description: models.Description{
					Value: "Il Pan d'arancio è un dolce della tradizione siciliana caratterizzato da un intenso sapore agrumato, dato dall'utilizzo delle arance intere, buccia compresa.",
				},
				Image: models.Image{Value: "https://ricetta.it/Uploads/Imgs/pan-d-arancio_medium.jpg.webp"},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 uova", "100 g di olio", "250 g di zucchero", "300 g di farina 00",
						"100 ml di latte", "1 bustina di lievito", "1 arancia bio (senza semi)",
						"1 cucchiaino di essenza di vaniglia",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Per prima cosa laviamo accuratamente l'arancia e tagliamola a pezzetti, mantenendo la buccia (eliminando con cura gli eventuali semi che risulterebbero amari) [1]. Mettiamola in un mixer e aggiungiamo il latte [2] l'olio [3] e frulliamo fino ad ottenere un composto omogeneo e cremoso [4]. In una boule setacciamo la farina e uniamo lo zucchero [5]; aggiungiamo le uova [6] e il cucchiaino di vaniglia [7]. Aiutandoci con le fruste elettriche amalgamiamo gli ingredienti per qualche secondo. A questo punto uniamo il composto di arancia precedentemente frullata e azioniamo di nuovo le fruste [8]. Aggiungiamo il lievito [9] e mescoliamo per bene [10]. Imburriamo e foderiamo una teglia da 24 cm e trasferiamovi il composto [11]. Cuociamo a 180° per 60 minuti con forno statico (160° per 50 minuti se ventilato). A cottura ultimata, sformiamo la torta, lasciamola raffreddare e cospargiamo la superficie con marmellata di arance [12].",
					},
				},
				Name:  "Pan d'arancio",
				Yield: models.Yield{Value: 6},
				URL:   "https://ricetta.it/pan-d-arancio",
			},
		},
		{
			name: "rosannapansino.com",
			in:   "https://rosannapansino.com/blogs/recipes/rainbow-treats",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				DatePublished: "2023-06-14T19:08:32Z",
				Description:   models.Description{Value: "Check out these colorful and flavorful rainbow treats!"},
				Image: models.Image{
					Value: "https://cdn.shopify.com/s/files/1/0163/5948/9636/files/rainbow_treats_thumbnail_3_2048x2048.jpg?v=1686769264",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 stick (4 ounces) unsalted butter, room temperature", "1 cup sugar",
						"3 egg whites", "2 teaspoons vanilla extract", "1/4 teaspoons almond extract",
						"1 1/2 cups all-purpose flour", "1 1/2 teaspoons baking powder",
						"1/2 teaspoon salt", "1/2 cup whole milk, room temperature",
						"1 cup buttercream", "2 cups white Candy Melts",
						"Food coloring powders (Suncore): red, orange, yellow, green, blue, purple",
						"Rainbow sprinkles (ColorKitchen)", `Two 6" inch cake pans`, "Parchment paper",
						"Offset spatula",
						"Whisk",
						"Rubber spatula",
						"Hand mixer",
						"Melting Pot (Nerdy Nummies)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat the oven to 325F°. Grease and line the cake pans.",
						"In a small bowl, whisk together the flour, baking powder and salt.",
						"In a medium bowl using hand mixer, cream together the butter, sugar and extracts until light and fluffy, 6 to 8 minutes.",
						"Add the egg whites to the butter mixture, mix until combined (do not over mix).",
						"Alternate adding the flour mixture and milk in three addition, starting and ending with the flour. Do not over mix.",
						"Pour batter evenly into prepared pans. Bake for 30 minutes or until toothpick comes out clean.",
						"In a large bowl, crumble the cakes and add the buttercream.Mix with the hand mixer until a smooth dough forms. It should feel like cookie dough.",
						"Separate the cake dough evenly between 6 bowls. Add and mix the colors red, orange, yellow, green, blue and purple into each bowl until desired colors - about 1/2 teaspoons to 1 teaspoon per bowl.",
						"Split all the colored dough into 1/2 teaspoon pieces. Stack 6 pieces of each color starting with the red and ending with the purple and roll into a smooth ball. Repeat with the rest of the dough (about 40 balls). Chill in refrigerator for 1 hour.",
						"Melt the white Candy Melts in the melting pot. Dip each cake ball in the candy and sprinkle the rainbow sprinkles on top.",
					},
				},
				Name:  "Rainbow Treats",
				Yield: models.Yield{Value: 40},
				URL:   "https://rosannapansino.com/blogs/recipes/rainbow-treats",
			},
		},
		{
			name: "rutgerbakt.nl",
			in:   "https://rutgerbakt.nl/basisrecepten/oreo-topping-van-roomkaas/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				DateModified:  "2023/09/19",
				DatePublished: "2022/01/28",
				Description: models.Description{
					Value: "Ben je aan het bakken geslagen met Oreo’s en zoek je een bijpassende vulling of topping? Deze crème is perfect als Oreo cupcake topping of als toppin",
				},
				Image: models.Image{
					Value: "https://rutgerbakt.nl/wp-content/uploads/2022/01/oreo_topping-scaled-1920x1080-c-default.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"10 oreo koekjes", "150 gr boter, op kamertemperatuur", "150 gr poedersuiker",
						"1 tl vanille-extract", "200 gr roomkaas, op kamertemperatuur",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Ben je aan het bakken geslagen met Oreo’s en zoek je een bijpassende vulling of topping? Deze crème is perfect als Oreo cupcake topping of als toppin",
					},
				},
				Name:     "Oreo topping van roomkaas",
				PrepTime: "PT12M",
				URL:      "https://rutgerbakt.nl/basisrecepten/oreo-topping-van-roomkaas/",
				Yield:    models.Yield{Value: 1},
			},
		},
		{
			name: "recipecommunity.com.au",
			in:   "https://www.recipecommunity.com.au/baking-sweet-recipes/flourless-refined-sugar-free-chocolate-cake/1te0mta9-5d0d3-705689-cfcd2-7zd1b4nd",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Name:          "Flourless refined sugar free chocolate cake",
				DatePublished: "2016-06-08",
				DateModified:  "2016-06-11",
				Description: models.Description{
					Value: "Recipe Flourless refined sugar free chocolate cake by Mixing Adventures, learn to make this recipe easily in your kitchen machine and discover other Thermomix recipes in Baking - sweet.",
				},
				Keywords: models.Keywords{Values: "Baking - sweet, recipes, Dessert, Gluten free, Lactose free, Non-dairy"},
				Image: models.Image{
					Value: "https://d2mkh7ukbp9xav.cloudfront.net/recipeimage/1te0mta9-5d0d3-705689-cfcd2-7zd1b4nd/fd1b729e-862e-4c96-b038-18cf8656a293/large/flourless-refined-sugar-free-chocolate-cake.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"270 g Dates (pitted)",
						"125 g Boiling water",
						"1/2 teaspoon Bi-carbonate of Soda",
						"250 g almonds",
						"3 eggs",
						"60 g coconut oil",
						"2 teaspoons vanilla essence",
						"50 g cocoa",
						"1 teaspoon baking powder",
						"pinch salt",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Pre-heat oven to 180 degrees C.  Grease and line a 20 cm round cake pan (I used a silicone ring tin)Mix dates and bi-carb soda in boiling water and leave to soak for a few minutes while you do the next step.",
						"Grind almonds for 10 seconds, speed 9.Set aside.",
						"Put dates, water and bi-carb into mixing bowl, and blend for 30 seconds, speed 6.  You may need to stop and scrape the sides a couple of times, and maybe reduce the speed a little after the first few seconds if necessary.   ",
						"Add eggs, oil, vanilla essence and blend for 20 seconds, speed 6.  Scrape down bowl.",
						"Add almonds, baking powder, cocoa and salt and mix for 20 seconds, speed 6.  Scrape down bowl and repeat. ",
						"Pour into prepared cake tin and bake at 180 degrees for 30 minutes, or until a skewer comes out clean.  Leave until cool before turning onto a serving plate. I covered mine with chocolate ganache and served with whipped cream and raspberry coulis.  ",
					},
				},
				PrepTime: "PT10M",
				CookTime: "PT40M",
				URL:      "https://www.recipecommunity.com.au/baking-sweet-recipes/flourless-refined-sugar-free-chocolate-cake/1te0mta9-5d0d3-705689-cfcd2-7zd1b4nd",
				Yield:    models.Yield{Value: 8},
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
