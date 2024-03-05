package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_I(t *testing.T) {
	testcases := []testcase{
		{
			name: "ica.se",
			in:   "https://www.ica.se/recept/chicken-a-la-king-729980/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Huvudrätt,Middag"},
				CookingMethod: models.CookingMethod{Value: "I gryta"},
				DateModified:  "2023-09-08",
				DatePublished: "2023-09-08",
				Description: models.Description{
					Value: "En krämig och god kycklinggryta är chicken a la king. Grytan görs av smakrika kycklinglårfiléer, champinjoner och paprika. Samt naturligtvis vispgrädde för den perfekta krämigheten. Servera gärna kycklinggrytan med ris och en grön sallad.",
				},
				Image: models.Image{
					Value: "https://assets.icanet.se/t_ICAseAbsoluteUrl/imagevaultfiles/id_250717/cf_259/chicken_à_la_king.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"4 port ris", "1 gul lök", "250 g färska champinjoner",
						"2 röda eller gula paprikor", "500 g kycklinglårfilé", "3 msk smör",
						"3 dl vispgrädde", "2 1/2 dl vatten", "1 1/2 hönsbuljongtärning",
						"4 tsk majsstärkelse", "salt", "svartpeppar", "1/2 citron",
						"1 kruka persilja", "gärna sallad",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Koka riset enligt anvisning p&aring; f&ouml;rpackningen.",
						"Skala och finhacka l&ouml;ken. Ansa och skiva champinjonerna. Dela, k&auml;rna ur och sk&auml;r paprikan i ca 2 cm stora t&auml;rningar.",
						"Dela kycklingl&aring;ren i 4 bitar.",
						"Stek svamp och l&ouml;k i sm&ouml;ret i en stor gryta. N&auml;r svampen b&ouml;rjar f&aring; f&auml;rg, tills&auml;tt kycklingen. R&ouml;r runt kycklingen.",
						"Tills&auml;tt gr&auml;dde, vatten, smulad buljongt&auml;rning och paprika. Koka i ca 10 minuter.",
						"R&ouml;r ut majsst&auml;rkelsen med lite kallt vatten och tills&auml;tt i kycklinggrytan under omr&ouml;ring.",
						"Smaka av grytan med salt, peppar och n&aring;gra droppar citronsaft.",
						"Hacka persiljan.",
						"Till servering: Servera grytan, toppad med persilja och med ris och g&auml;rna sallad.",
					},
				},
				Name: "Chicken à la king",
				NutritionSchema: models.NutritionSchema{
					Calories:      "775 calories",
					Carbohydrates: "58 g",
					Fat:           "45 g",
					Protein:       "34 g",
					Servings:      "4",
				},
				Yield: models.Yield{Value: 4},
				URL:   "https://www.ica.se/recept/chicken-a-la-king-729980/",
			},
		},
		{
			name: "indianhealthyrecipes.com",
			in:   "https://www.indianhealthyrecipes.com/mango-rice-mamidikaya-pulihora/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main"},
				CookTime:      "PT25M",
				Cuisine:       models.Cuisine{Value: "South Indian"},
				DatePublished: "2016-04-01T07:34:51+00:00",
				Description: models.Description{
					Value: "This Mango rice is a traditional South Indian dish made with precooked rice, raw green unripe " +
						"mangoes tempering spices and curry leaves. It tastes slightly tangy, hot and flavorful.",
				},
				Keywords: models.Keywords{Values: "mango rice, mango rice recipe"},
				Image: models.Image{
					Value: "https://www.indianhealthyrecipes.com/wp-content/uploads/2022/04/mango-rice-recipe.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 cups rice",
						"1 raw unripe green mango ((Sour, medium sized))",
						"1 to 2 sprigs curry leaves",
						"3 to 4 green chilies ( slit or chopped)",
						"1 to 2 dried red chilies (broken)",
						"⅖ to 1 teaspoon salt ((adjust to taste))",
						"¼ teaspoon turmeric ((prefer organic))",
						"¼ cup peanuts (or cashewnuts)",
						"1 tablespoon chana dal ((bengal gram))",
						"1 tablespoon urad dal ((skinned split black gram))",
						"1 teaspoon mustard seeds",
						"1 inch ginger (chopped, sliced, grated )",
						"3 tablespoons oil",
						"1 pinch hing ((⅛ teaspoon asafetida))",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Cook rice to grainy texture: Add rice to a bowl &amp; rinse it a few times. Pour 4 cups water &amp; " +
							"place the bowl in a pressure cooker. Cover the bowl &amp; pressure cook for 3 whistles.",
						"When the pressure drops, remove the rice &amp; cool completely.",
						"Wash, peel and grate or chop the mango. You can also grate in a processor. Update: Lately I have started " +
							"to make this with cooked mango. I just peel, cut and add them to a bowl. Cover and place it in " +
							"the pressure cooker over the rice bowl (PIP). Once cooked I mash it and use as mentioned below.",
						"Heat 1 tablespoon oil in a pan and fry the peanuts on a medium heat until aromatic and golden. Remove them " +
							"to a plate for later.",
						"Pour 2 tablespoons more oil and heat it. Add chana dal, urad dal, mustard seeds and dried red chilli.",
						"When the lentils turn light golden, add ginger, green chilies &amp; curry leaves. Fry till the curry leaves " +
							"become crisp, then add hing.",
						"Add mango, salt and turmeric.Saute for 2 to 3 minutes. Cook covered until the mango turns mushy, completely " +
							"soft &amp; pulpy. (skip this with cooked mango.)",
						"Add this to the cooked rice little by little and begin to mix. Taste test and add more mango mixture as " +
							"required. Adjust salt and oil at this stage.",
						"Transfer mango rice to serving plates and garnish with roasted peanuts.",
					},
				},
				Name: "Mango Rice Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories:      "636 kcal",
					Carbohydrates: "83 g",
					Fat:           "28 g",
					Fiber:         "7 g",
					Protein:       "11 g",
					SaturatedFat:  "13 g",
					Servings:      "1",
					Sodium:        "28 mg",
					Sugar:         "1 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.indianhealthyrecipes.com/mango-rice-mamidikaya-pulihora/",
			},
		},
		{
			name: "im-worthy.com",
			in:   "https://im-worthy.com/cranberry-walnut-oatmeal-energy-balls/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizer"},
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-10-25T10:23:17+00:00",
				Description: models.Description{
					Value: "Indulge in these scrumptious no-bake Cranberry-Walnut Energy Balls, a delightful blend of sweet cranberries and brain-boosting walnuts. Crafted with love, these energy-packed wonders offer a quick and easy way to nourish your body, whether you&#39;re starting your day or need an on-the-go breakfast or snack. Satisfy your cravings with wholesome, bite-sized goodness – no baking required!",
				},
				Keywords: models.Keywords{
					Values: "energy balls, energy bites, mediterranean diet desserts, oatmeal energy balls",
				},
				Image: models.Image{
					Value: "https://im-worthy.com/wp-content/uploads/2023/10/Oatmeal-Energy-Balls_blog-17.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3/4 cup raw walnuts", "1/2 cup sweetened dried cranberries",
						"1/4 cup pitted dates", "3/4 cup old-fashioned rolled oats", "2 tbsp tahini",
						"2 tbsp fresh lemon juice", "1 tbsp pure maple syrup",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Begin by placing your raw walnuts, dried cranberries, and pitted dates in a food processor. Give them a whirl until they become finely chopped but still maintain some texture. It&#39;s like creating the canvas for your masterpiece.",
						"Add the old-fashioned rolled oats, tahini, fresh lemon juice, and pure maple syrup to the food processor. These ingredients join the party, creating a vibrant dance of flavors. Pulse everything until the mixture becomes sticky and easily forms into balls.",
						"With clean hands, shape the mixture into bite-sized balls, like crafting your edible work of art. Place them on a baking sheet lined with parchment paper.",
						"Let the oatmeal balls chill in the refrigerator for about 30 minutes. This step is your chance to practice a bit of mindfulness. While they&#39;re cooling, take a moment to savor the delightful anticipation of your delicious creation.",
					},
				},
				Name: "Cranberry-Walnut Oatmeal Energy Balls (No-Bake)",
				NutritionSchema: models.NutritionSchema{
					Calories: "170 kcal",
					Servings: "1",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://im-worthy.com/cranberry-walnut-oatmeal-energy-balls/",
			},
		},
		{
			name: "innit.com",
			in:   "https://www.innit.com/meal/504/8008/Salad%3A%20Coconut-Pineapple-Salad",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Salads and Sides"},
				CookTime:      "PT7M",
				DatePublished: "2022-02-12",
				Image: models.Image{
					Value: "https://www.innit.com/meal-service/en-US/images/Meal-Salad%3A%20Coconut_Pineapple_Salad_" +
						"1529953193419_480x480.png",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 Fresh Mexican Limes",
						"2 cups Pineapple",
						"1/2 cup Mint",
						"2 cups Jasmine Rice",
						"2 cups Canned Coconut Milk",
						"2 tsp Kosher Salt",
						"1/4 tsp Korean Chili Flakes",
						"2 Tbsp Extra Virgin Olive Oil",
						"1 cup Macadamia Nuts",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Combine ingredients in pot & bring to a boil.",
						"Cover & simmer for 15 minutes.",
						"Remove from heat & steam with lid for 5 minutes.",
						"Dice pineapple.",
						"Combine ingredients in a large bowl. Mix well.",
						"Plate on platter or bowl & garnish with macadamia nuts, mint & chili flakes.",
					},
				},
				Name: "Coconut Pineapple Rice",
				NutritionSchema: models.NutritionSchema{
					Calories:       "880 kcal",
					Carbohydrates:  "88 g",
					Cholesterol:    "0 mg",
					Fat:            "56 g",
					Fiber:          "7 g",
					Protein:        "12 g",
					SaturatedFat:   "28 g",
					Sodium:         "1190 mg",
					Sugar:          "10 g",
					UnsaturatedFat: "28 g",
				},
				PrepTime: "PT28M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.innit.com/meal/504/8008/Salad%3A%20Coconut-Pineapple-Salad",
			},
		},
		{
			name: "insanelygoodrecipes.com",
			in:   "https://insanelygoodrecipes.com/chicken-cordon-bleu-casserole/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Chicken"},
				CookTime:      "PT30M",
				DatePublished: "2023-11-14T09:02:16-05:00",
				Description: models.Description{
					Value: "This chicken cordon bleu casserole has everything you love about the classic dish! It's creamy, savory, and full of delicious flavor.",
				},
				Image: models.Image{
					Value: "https://insanelygoodrecipes.com/wp-content/uploads/2023/11/Homemade-Chicken-Cordon-Bleu-with-Egg-Noodles-Ham-and-Cheese.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 (8-ounce) package wide egg noodles", "2 cups chopped cooked chicken breast",
						"8 ounces cooked ham, cubed", "8 ounces Swiss cheese, cubed",
						"1 can cream of chicken soup", "1/2 cup 2% milk", "1/2 cup sour cream",
						"2 tablespoons butter", "1/3 cup seasoned bread crumbs",
						"1/4 cup grated Parmesan cheese",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Gather all your ingredients. Preheat your oven to 350 degrees Fahrenheit (175°C) and lightly grease a 9x13-inch baking dish.",
						"In a large pot, boil lightly salted water. Cook the egg noodles in boiling water, stirring occasionally. Let them cook until they're tender yet firm to bite for 7 to 9 minutes. Drain and set them aside.",
						"Combine the cooked egg noodles, chopped chicken breast, cubed ham, and Swiss cheese in a large bowl. Set it aside.",
						"Mix the cream of chicken soup, milk, and sour cream in a separate bowl until well combined. Pour this mixture over the chicken, ham, and cheese. Stir until everything is well-coated.",
						"Transfer the noodle mixture to your prepared baking dish, spreading it out evenly.",
						"In a small pan, melt the butter over medium heat. Add the seasoned bread crumbs and stir until they're coated in the butter. Cook them for 2 to 3 minutes.",
						"Sprinkle the breadcrumb mixture evenly over the noodle mixture in the baking dish.",
						"Place the casserole in the preheated oven. Bake it for 30 minutes until it's hot and bubbly and the breadcrumb topping is golden brown.",
						"Allow the casserole to cool for 10 minutes, serve, and enjoy!",
					},
				},
				Name: "Easy Chicken Cordon Bleu Casserole Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories: "455 cal",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://insanelygoodrecipes.com/chicken-cordon-bleu-casserole/",
			},
		},
		{
			name: "inspiralized.com",
			in:   "https://inspiralized.com/vegetarian-zucchini-noodle-pad-thai/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT15M",
				DatePublished: "2014-05-05T12:00:03+00:00",
				Description: models.Description{
					Value: "Make quick and healthy zucchini noodle pad thai with eggs, hoisin sauce, peanuts and spiralized " +
						"zucchini for dinner tonight.",
				},
				Image: models.Image{
					Value: "https://inspiralized.com/wp-content/uploads/2014/05/IMG_9863-copy-1.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 whole eggs",
						"1/4 cup roasted salted peanuts",
						"1/2 tbsp peanut oil (or oil of choice)",
						"1 garlic clove (minced)",
						"1 shallot (minced)",
						"1 tbsp coconut flour",
						"1 tbsp roughly chopped cilantro + whole cilantro leaves to garnish",
						"2 medium zucchinis (Blade C)",
						"For the sauce:",
						"2 tbsp freshly squeezed lime juice",
						"1 tbsp fish sauce (or hoisin sauce, if you're strict vegetarian)",
						"1/2 tbsp soy sauce",
						"1 tbsp chili sauce (I used Thai chili garlic sauce)",
						"1 tsp honey",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Scramble the eggs and set aside.",
						"Place all of the ingredients for the sauce into a bowl, whisk together and set aside.",
						"Place the peanuts into a food processor and pulse until lightly ground (no big peanuts should " +
							"remain, but it shouldn't be powdery). Set aside.",
						"Place a large skillet over medium heat. Add in oil, garlic and shallots. Cook for about 1-2 minutes, " +
							"stirring frequently, until the shallots begin to soften. Add in the sauce and whisk quickly " +
							"so that the flour dissolves and the sauce thickens. Cook for 2-3 minutes or until sauce is " +
							"reduced and thick.",
						"Once the sauce is thick, add in the zucchini noodles and cilantro and stir to combine thoroughly.",
						"Cook for about 2 minutes or until noodles soften and then add in the scrambled eggs and ground peanuts. " +
							"Cook for about 30 seconds, tossing to fully combine.",
						"Plate onto dishes and garnish with cilantro leaves. Serve with lime wedges.",
					},
				},
				Name:     "Vegetarian Zucchini Noodle Pad Thai",
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://inspiralized.com/vegetarian-zucchini-noodle-pad-thai/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
