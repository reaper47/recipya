package scraper_test

import (
	"testing"
	"time"

	"github.com/reaper47/recipya/internal/models"
)

func TestScraper_I(t *testing.T) {
	testcases := []testcase{
		{
			name: "ica.se",
			in:   "https://www.ica.se/recept/chicken-a-la-king-729980/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Huvudrätt,Middag"},
				CookingMethod: &models.CookingMethod{Value: "I gryta"},
				Cuisine:       &models.Cuisine{},
				DateModified:  "2023-09-08",
				DatePublished: "2023-09-08",
				Description: &models.Description{
					Value: "En krämig och god kycklinggryta är chicken a la king. Grytan görs av smakrika kycklinglårfiléer, champinjoner och paprika. Samt naturligtvis vispgrädde för den perfekta krämigheten. Servera gärna kycklinggrytan med ris och en grön sallad.",
				},
				Keywords: &models.Keywords{},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"4 port ris", "1 gul lök", "250 g färska champinjoner",
						"2 röda eller gula paprikor", "500 g kycklinglårfilé", "3 msk smör",
						"3 dl vispgrädde", "2 1/2 dl vatten", "1 1/2 hönsbuljongtärning",
						"4 tsk majsstärkelse", "salt", "svartpeppar", "1/2 citron",
						"1 kruka persilja", "gärna sallad",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Koka riset enligt anvisning på förpackningen."},
						{Type: "HowToStep", Text: "Skala och finhacka löken. Ansa och skiva champinjonerna. Dela, kärna ur och skär paprikan i ca 2 cm stora tärningar."},
						{Type: "HowToStep", Text: "Dela kycklinglåren i 4 bitar."},
						{Type: "HowToStep", Text: "Stek svamp och lök i smöret i en stor gryta. När svampen börjar få färg, tillsätt kycklingen. Rör runt kycklingen."},
						{Type: "HowToStep", Text: "Tillsätt grädde, vatten, smulad buljongtärning och paprika. Koka i ca 10 minuter."},
						{Type: "HowToStep", Text: "Rör ut majsstärkelsen med lite kallt vatten och tillsätt i kycklinggrytan under omröring."},
						{Type: "HowToStep", Text: "Smaka av grytan med salt, peppar och några droppar citronsaft."},
						{Type: "HowToStep", Text: "Hacka persiljan."},
						{Type: "HowToStep", Text: "Till servering: Servera grytan, toppad med persilja och med ris och gärna sallad."},
					},
				},
				Name: "Chicken à la king",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "775",
					Carbohydrates: "58",
					Fat:           "45",
					Protein:       "34",
					Servings:      "4",
				},
				ThumbnailURL: &models.ThumbnailURL{},
				Tools:        &models.Tools{Values: []models.HowToItem{}},
				TotalTime:    "PT45M",
				Yield:        &models.Yield{Value: 4},
				URL:          "https://www.ica.se/recept/chicken-a-la-king-729980/",
				Video:        &models.Videos{},
			},
		},
		{
			name: "im-worthy.com",
			in:   "https://im-worthy.com/cranberry-walnut-oatmeal-energy-balls/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Appetizer"},
				Cuisine:       &models.Cuisine{Value: "American"},
				DatePublished: "2023-10-25T10:23:17+00:00",
				Description: &models.Description{
					Value: "Indulge in these scrumptious no-bake Cranberry-Walnut Energy Balls, a delightful blend of sweet cranberries and brain-boosting walnuts. Crafted with love, these energy-packed wonders offer a quick and easy way to nourish your body, whether you&#39;re starting your day or need an on-the-go breakfast or snack. Satisfy your cravings with wholesome, bite-sized goodness – no baking required!",
				},
				Keywords: &models.Keywords{
					Values: "energy balls, energy bites, mediterranean diet desserts, oatmeal energy balls",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"3/4 cup raw walnuts", "1/2 cup sweetened dried cranberries",
						"1/4 cup pitted dates", "3/4 cup old-fashioned rolled oats", "2 tbsp tahini",
						"2 tbsp fresh lemon juice", "1 tbsp pure maple syrup",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Begin by placing your raw walnuts, dried cranberries, and pitted dates in a food processor. Give them a whirl until they become finely chopped but still maintain some texture. It&#39;s like creating the canvas for your masterpiece."},
						{Type: "HowToStep", Text: "Add the old-fashioned rolled oats, tahini, fresh lemon juice, and pure maple syrup to the food processor. These ingredients join the party, creating a vibrant dance of flavors. Pulse everything until the mixture becomes sticky and easily forms into balls."},
						{Type: "HowToStep", Text: "With clean hands, shape the mixture into bite-sized balls, like crafting your edible work of art. Place them on a baking sheet lined with parchment paper."},
						{Type: "HowToStep", Text: "Let the oatmeal balls chill in the refrigerator for about 30 minutes. This step is your chance to practice a bit of mindfulness. While they&#39;re cooling, take a moment to savor the delightful anticipation of your delicious creation."},
					},
				},
				Name: "Cranberry-Walnut Oatmeal Energy Balls (No-Bake)",
				NutritionSchema: &models.NutritionSchema{
					Calories: "170",
					Servings: "1",
				},
				PrepTime:  "PT10M",
				TotalTime: "PT10M",
				Yield:     &models.Yield{Value: 8},
				URL:       "https://im-worthy.com/cranberry-walnut-oatmeal-energy-balls/",
			},
		},
		{
			name: "inbloombakery.com",
			in:   "https://inbloombakery.com/the-best-cinnamon-rolls-ever",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtGraph:       nil,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Appetizer"},
				CookTime:      "PT30M",
				CookingMethod: nil,
				Cuisine:       &models.Cuisine{Value: "American"},
				DateCreated:   "",
				DateModified:  "",
				DatePublished: "2023-08-07T17:07:48+00:00",
				Description: &models.Description{
					Value: "These are the best cinnamon rolls EVER. They&#39;re soft and fluffy with a gooey brown sugar cinnamon filling, frosted with cream cheese frosting all while drenched in a rich, buttery brown sugar cinnamon caramel sauce on the bottom.",
				},
				Keywords: &models.Keywords{Values: "best cinnamon rolls, cinnamon, cinnamon rolls, cream cheese frosting"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1 1/4 cups (300 ml) whole milk, warmed to about 110 degrees",
						"2 1/4 tsp active dry yeast",
						"1 tsp granulated white sugar (to bloom the yeast)",
						"4 3/4 cups (593 g) all-purpose flour, spooned and leveled", "1 1/2 tsp salt",
						"2 tbsp (25 g) granulated white sugar", "2 eggs, whisked", "1 tbsp vanilla",
						"1/2 cup (112 g) unsalted butter, very softened",
						"1/2 cup (112 g) unsalted butter, very softened",
						"1 cup (220 g) light brown sugar, packed",
						"2 tsp cinnamon",
						"1/3 cup (116 g) honey",
						"1/2 tbsp vanilla",
						"3 tbsp (45 ml) heavy cream",
						"1/4 tsp salt",
						"1/2 cup (112 g) unsalted butter, very softened",
						"1 cup (220 g) light brown sugar, packed",
						"2 tbsp cinnamon",
						"1/4 tsp salt",
						"1/4 cup (60 ml) heavy cream (for pouring in between the rolls)",
						"6 tbsp (84 g) unsalted butter, very softened",
						"6 oz (170 g) cream cheese, cold",
						"3/4 cup (97 g) powdered sugar",
						"1/2 tbsp vanilla",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Text: "To start, bloom your yeast. Mix your active dry yeast, sugar, and warm milk together and let sit for 10 minutes until foamy on top.", Type: "HowToStep"},
						{Text: "Next, in a large bowl of a stand mixer with a whisk attachment, mix flour, salt and sugar together.", Type: "HowToStep"},
						{Text: "To the dry ingredients add the whisked eggs, vanilla and softened butter and mix.", Type: "HowToStep"},
						{Text: "Add in the yeast mixture and mix until everything is blended. Switch to a dough hook and knead on medium speed for 7-10 minutes until the dough pulls away from the edges of the bowl, forms a ball and springs back when pressed.", Type: "HowToStep"},
						{Text: "Form the dough into a ball. Place the dough in a large, greased bowl. Cover with plastic wrap or a kitchen towel and place in a warm place to rise, (outside or a cold oven with the light on works well for this), for 1- 1 ½ hours, or until the dough doubles in size.", Type: "HowToStep"},
						{Text: "While the dough rises, make the sticky caramel sauce for the bottoms of the rolls, by mixing the butter, brown sugar, cinnamon, honey, vanilla, heavy cream and salt together in a medium bowl.", Type: "HowToStep"},
						{Text: "Grease a 9x13 inch casserole dish. Spread the mixture evenly over the bottom of the pan.", Type: "HowToStep"},
						{Text: "While the dough rises, also make the filling by mixing the butter, brown sugar, cinnamon and salt together in a small bowl. (If the filling is too firm and not spreadable, add heavy cream a tsp at a time until it&#39;s easily spreadable.)", Type: "HowToStep"},
						{
							Text: "Once the dough has had a chance to rise, remove from bowl, punch dough to release the air, and roll out on alightly floured surface. Roll the dough into about an 18 x 12 inch rectangle. It should be about 1/4 inch thick.",
							Type: "HowToStep",
						},
						{
							Text: "Sprinkle the cinnamon sugar filling over the dough and spread it evenly with an offset spatula.",
							Type: "HowToStep",
						},
						{
							Text: "Roll up the dough tightly into a long log shape starting from the end closest to you. Cut of a bit off the ends to make the log even. Cut 12 rolls about 1 1/2 inches wide using unflavored floss or a very sharp knife.",
							Type: "HowToStep",
						},
						{
							Text: "Place the rolls in the prepared casserole dish with the cinnamon caramel sauce mixture and pour room temperature heavy cream in between each roll. Cover them with plastic wrap and let them proof in a warm spot for about an hour, or until they have doubled in size.",
							Type: "HowToStep",
						},
						{
							Text: "Preheat the oven to 350 degrees about 15 minutes before the rolls are done proofing.",
							Type: "HowToStep",
						},
						{
							Text: "Once the rolls have completed their second rise, bake them for 29-32 minutes until golden brown. (30 minutes is usually perfect!)(Cover with aluminum foil for the last 5 minutes if they&#39;re getting too brown.)",
							Type: "HowToStep",
						},
						{
							Text: "While they bake, make the cream cheese frosting by adding the softened butter to a medium bowl and mixing it on high speed with an electric mixer until pale and fluffy.",
							Type: "HowToStep",
						},
						{
							Text: "Add in the cream cheese and mix until combined on medium speed.",
							Type: "HowToStep",
						},
						{
							Text: "Sift the powdered sugar into the mixture a little at a time and mix on low speed until all is combined.(You can add more powdered sugar if you want the frosting to be sweeter.)",
							Type: "HowToStep",
						},
						{
							Text: "Then add in the vanilla and combine on medium speed until the frosting is smooth.",
							Type: "HowToStep",
						},
						{
							Text: "Let the rolls cool for 10 minutes, then cover the warm rolls with the cream cheese icing. Then serve!",
							Type: "HowToStep",
						},
					},
				},
				Name:            "The Best Cinnamon Rolls Ever",
				NutritionSchema: nil,
				PrepTime:        "PT45M",
				ThumbnailURL:    nil,
				Tools:           nil,
				TotalTime:       "PT225M",
				Yield:           &models.Yield{Value: 12},
				URL:             "https://inbloombakery.com/the-best-cinnamon-rolls-ever",
			},
		},
		{
			name: "indianhealthyrecipes.com",
			in:   "https://www.indianhealthyrecipes.com/mango-rice-mamidikaya-pulihora/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Main"},
				CookTime:      "PT25M",
				Cuisine:       &models.Cuisine{Value: "South Indian"},
				DatePublished: "2016-04-01T07:34:51+00:00",
				Description: &models.Description{
					Value: "This Mango rice is a traditional South Indian dish made with precooked rice, raw green unripe mangoes tempering spices and curry leaves. It tastes slightly tangy, hot and flavorful.",
				},
				Keywords: &models.Keywords{Values: "mango rice, mango rice recipe"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
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
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Cook rice to grainy texture: Add rice to a bowl &amp; rinse it a few times. Pour 4 cups water &amp; place the bowl in a pressure cooker. Cover the bowl &amp; pressure cook for 3 whistles."},
						{Type: "HowToStep", Text: "When the pressure drops, remove the rice &amp; cool completely."},
						{Type: "HowToStep", Text: "Wash, peel and grate or chop the mango. You can also grate in a processor. Update: Lately I have started to make this with cooked mango. I just peel, cut and add them to a bowl. Cover and place it in the pressure cooker over the rice bowl (PIP). Once cooked I mash it and use as mentioned below."},
						{Type: "HowToStep", Text: "Heat 1 tablespoon oil in a pan and fry the peanuts on a medium heat until aromatic and golden. Remove them to a plate for later."},
						{Type: "HowToStep", Text: "Pour 2 tablespoons more oil and heat it. Add chana dal, urad dal, mustard seeds and dried red chilli."},
						{Type: "HowToStep", Text: "When the lentils turn light golden, add ginger, green chilies &amp; curry leaves. Fry till the curry leaves become crisp, then add hing."},
						{Type: "HowToStep", Text: "Add mango, salt and turmeric.Saute for 2 to 3 minutes. Cook covered until the mango turns mushy, completely soft &amp; pulpy. (skip this with cooked mango.)"},
						{Type: "HowToStep", Text: "Add this to the cooked rice little by little and begin to mix. Taste test and add more mango mixture as required. Adjust salt and oil at this stage."},
						{Type: "HowToStep", Text: "Transfer mango rice to serving plates and garnish with roasted peanuts."},
					},
				},
				Name: "Mango Rice Recipe",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "636",
					Carbohydrates: "83",
					Fat:           "28",
					Fiber:         "7",
					Protein:       "11",
					SaturatedFat:  "13",
					Servings:      "1",
					Sodium:        "28",
					Sugar:         "1",
				},
				PrepTime:  "PT10M",
				TotalTime: "PT35M",
				URL:       "https://www.indianhealthyrecipes.com/mango-rice-mamidikaya-pulihora/",
				Video: &models.Videos{
					Values: []models.VideoObject{
						{
							AtType:       "VideoObject",
							ContentURL:   "https://www.youtube.com/watch?v=E4wtxm42u9Y",
							Description:  "For complete recipe of mango rice, you can check this link https://www.indianhealthyrecipes.com/mango-rice-mamidikaya-pulihora/",
							Duration:     "PT1M38S",
							EmbedURL:     "https://www.youtube.com/embed/E4wtxm42u9Y?feature=oembed",
							Name:         "Mango rice recipe | Mamidikaya pulihora | South Indian mango rice",
							ThumbnailURL: &models.ThumbnailURL{Value: "https://i.ytimg.com/vi/E4wtxm42u9Y/hqdefault.jpg"},
							UploadDate:   time.Date(2017, 4, 29, 9, 25, 16, 0, time.UTC),
						},
					},
				},
				Yield: &models.Yield{Value: 4},
			},
		},
		{
			name: "innit.com",
			in:   "https://www.innit.com/meal/504/8008/Salad%3A%20Coconut-Pineapple-Salad",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Salads and Sides"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				CookTime:      "PT7M",
				DatePublished: "2022-02-12",
				Description:   &models.Description{},
				Keywords:      &models.Keywords{},
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
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
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Combine ingredients in pot & bring to a boil."},
						{Type: "HowToStep", Text: "Cover & simmer for 15 minutes."},
						{Type: "HowToStep", Text: "Remove from heat & steam with lid for 5 minutes."},
						{Type: "HowToStep", Text: "Dice pineapple."},
						{Type: "HowToStep", Text: "Combine ingredients in a large bowl. Mix well."},
						{Type: "HowToStep", Text: "Plate on platter or bowl & garnish with macadamia nuts, mint & chili flakes."},
					},
				},
				Name: "Coconut Pineapple Rice",
				NutritionSchema: &models.NutritionSchema{
					Calories:       "880",
					Carbohydrates:  "88",
					Cholesterol:    "0",
					Fat:            "56",
					Fiber:          "7",
					Protein:        "12",
					SaturatedFat:   "28",
					Sodium:         "1190",
					Sugar:          "10",
					UnsaturatedFat: "28",
				},
				PrepTime:     "PT28M",
				ThumbnailURL: &models.ThumbnailURL{},
				Tools:        &models.Tools{Values: []models.HowToItem{}},
				TotalTime:    "PT35M",
				URL:          "https://www.innit.com/meal/504/8008/Salad%3A%20Coconut-Pineapple-Salad",
				Video: &models.Videos{
					Values: []models.VideoObject{
						{
							ContentURL: "https://www.innit.com/meal-service/en-US/videos/MealTask-Salad%3A%20Coconut_Pineapple_Salad_1529612492439_1920x1080.mp4",
							Name:       "Coconut Pineapple Rice ",
							ThumbnailURL: &models.ThumbnailURL{
								Value: "https://www.innit.com/meal-service/en-US/images/MealTask-Salad%3A%20Coconut_Pineapple_Salad_1529612492439_EXTRACTED_MIDDLE_480x480.png",
							},
						},
					},
				},
				Yield: &models.Yield{Value: 4},
			},
		},
		{
			name: "insanelygoodrecipes.com",
			in:   "https://insanelygoodrecipes.com/chicken-cordon-bleu-casserole/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Chicken"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				CookTime:      "PT30M",
				DatePublished: "2023-11-14T09:02:16-05:00",
				Description: &models.Description{
					Value: "This chicken cordon bleu casserole has everything you love about the classic dish! It's creamy, savory, and full of delicious flavor.",
				},
				Keywords: &models.Keywords{},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1 package (8-ounce) wide egg noodles", "2 cups chopped cooked chicken breast",
						"8 ounces cooked ham, sliced or cubed", "8 ounces Swiss cheese, cubed",
						"1 can cream of chicken soup", "1/2 cup 2% milk", "1/2 cup sour cream",
						"2 tablespoons butter", "1/3 cup seasoned bread crumbs",
						"1/4 cup grated Parmesan cheese",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Preheat the oven to 350 degrees Fahrenheit (175°C) and lightly grease a 9x13-inch baking dish."},
						{Type: "HowToStep", Text: "In a large pot of well-salted water, boil the egg noodles and cook to 2 minutes before al dente per the package instructions (usually about 5 minutes)."},
						{Type: "HowToStep", Text: "Drain the noodles then add them to a large bowl. Mix in the chopped chicken breast, ham, and Swiss cheese. Set aside."},
						{Type: "HowToStep", Text: "In a separate bowl, mix the cream of chicken soup, milk, and sour cream until smooth. Pour this over the noodles and chicken and stir until everything is well-coated."},
						{Type: "HowToStep", Text: "Transfer the noodle mixture to the prepared baking dish, spreading it out evenly."},
						{Type: "HowToStep", Text: "In a small pan, melt the butter over medium heat. Add the seasoned bread crumbs and Parmesan cheese and stir until they're coated in the butter. Cook them for 2 to 3 minutes."},
						{Type: "HowToStep", Text: "Sprinkle the breadcrumb mixture evenly over the noodles and bake for 30 minutes or until it's hot and bubbly and the breadcrumbs are golden brown."},
						{Type: "HowToStep", Text: "Allow the casserole to cool for 10 minutes, serve, and enjoy!"},
					},
				},
				Name: "Easy Chicken Cordon Bleu Casserole Recipe",
				NutritionSchema: &models.NutritionSchema{
					Calories: "455",
				},
				PrepTime:     "PT10M",
				ThumbnailURL: &models.ThumbnailURL{},
				Tools:        &models.Tools{Values: []models.HowToItem{}},
				TotalTime:    "PT40M",
				Yield:        &models.Yield{Value: 8},
				URL:          "https://insanelygoodrecipes.com/chicken-cordon-bleu-casserole/",
				Video:        &models.Videos{},
			},
		},
		{
			name: "inspiralized.com",
			in:   "https://inspiralized.com/vegetarian-zucchini-noodle-pad-thai/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookTime:      "PT15M",
				DatePublished: "2014-05-05T12:00:03+00:00",
				Description: &models.Description{
					Value: "Make quick and healthy zucchini noodle pad thai with eggs, hoisin sauce, peanuts and spiralized zucchini for dinner tonight.",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
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
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Scramble the eggs and set aside."},
						{Type: "HowToStep", Text: "Place all of the ingredients for the sauce into a bowl, whisk together and set aside."},
						{Type: "HowToStep", Text: "Place the peanuts into a food processor and pulse until lightly ground (no big peanuts should remain, but it shouldn't be powdery). Set aside."},
						{Type: "HowToStep", Text: "Place a large skillet over medium heat. Add in oil, garlic and shallots. Cook for about 1-2 minutes, stirring frequently, until the shallots begin to soften. Add in the sauce and whisk quickly so that the flour dissolves and the sauce thickens. Cook for 2-3 minutes or until sauce is reduced and thick."},
						{Type: "HowToStep", Text: "Once the sauce is thick, add in the zucchini noodles and cilantro and stir to combine thoroughly."},
						{Type: "HowToStep", Text: "Cook for about 2 minutes or until noodles soften and then add in the scrambled eggs and ground peanuts. Cook for about 30 seconds, tossing to fully combine."},
						{Type: "HowToStep", Text: "Plate onto dishes and garnish with cilantro leaves. Serve with lime wedges."},
					},
				},
				Name:      "Vegetarian Zucchini Noodle Pad Thai",
				PrepTime:  "PT10M",
				TotalTime: "PT25M",
				Yield:     &models.Yield{Value: 2},
				URL:       "https://inspiralized.com/vegetarian-zucchini-noodle-pad-thai/",
			},
		},
		{
			name: "instantpot.com",
			in:   "https://instantpot.com/blogs/recipes/cilantro-lime-rice",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "side dish"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				Description: &models.Description{
					Value: "Instant® is a well-known brand of kitchen and home appliances. Find pressure cookers, air purifiers, air fryers & cookware. Ships to USA and Canada.",
				},
				Keywords: &models.Keywords{
					Values: "15-30min,cilantro,easy,InstantPotCulinaryTeam,limes,pressurecook,rice,sidedish,vegetarian,water",
				},
				Image:       &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{Values: []string{"3 cups rice", "3 cups Water", "3  limes", "1/4 cup cilantro"}},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Text: "Pour rice and water into your Instant Pot.", Type: "HowToStep"},
						{Text: "Use the Rice Function and Pressure Cook.", Type: "HowToStep"},
						{Text: "Do a Quick Release. Juice the limes into the rice, use as much or as little to your taste. Sprinkle cilantro and stir the rice until thoroughly mixed.", Type: "HowToStep"},
						{Text: "Enjoy!", Type: "HowToStep"},
					},
				},
				Name:            "Cilantro Lime Rice",
				NutritionSchema: &models.NutritionSchema{},
				PrepTime:        "PT15M",
				ThumbnailURL:    &models.ThumbnailURL{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				Yield:           &models.Yield{Value: 1},
				URL:             "https://instantpot.com/blogs/recipes/cilantro-lime-rice",
				Video:           &models.Videos{},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
