package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_E(t *testing.T) {
	testcases := []testcase{
		{
			name: "eatingbirdfood.com",
			in:   "https://www.eatingbirdfood.com/cinnamon-rolls/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Breakast"},
				CookTime:      "PT25M",
				Cuisine:       &models.Cuisine{Value: "American"},
				DatePublished: "2022-12-15T08:05:22+00:00",
				Description: &models.Description{
					Value: "Make cinnamon rolls from scratch with this easy recipe that&#039;s perfect for beginners! They&#039;re soft, gooey, and made with bread flour, which gives them the perfect fluffy texture. Overnight instructions included.",
				},
				Keywords: &models.Keywords{Values: "cinnamon rolls"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1 0.25 oz. pkg. active dry yeast",
						"1 cup warm milk (around 105°–115°F)",
						"3 cups bread flour (plus more for rolling dough)",
						"1 teaspoon sea salt",
						"2 Tablespoons granulated sugar",
						"3 Tablespoons unsalted butter, melted (plus more for greasing)",
						"1 large egg (at room temperature)",
						"½ cup brown sugar (packed )",
						"1 Tablespoon ground cinnamon",
						"¼ cup unsalted butter (softened)",
						"4 oz cream cheese (full fat, softened to room temperature)",
						"¼ cup Greek yogurt (I used plain full fat)",
						"2 Tablespoons maple syrup",
						"1 teaspoon vanilla extract",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Heat milk in a saucepan or in the microwave (about 40-50 seconds) until warm, but not hot. Around (105°–115°F). Stir yeast into warm almond milk until dissolved. Let stand 10 minutes."},
						{Text: "In a stand mixer with the paddle attachment, add flour, salt and sugar. Mix to combine. With mixer running at low speed, add melted butter, egg and yeast mixture. Increase speed to medium-low, and mix 2 minutes until dough starts to form."},
						{Text: "Switch to the dough hook attachment and knead the dough in the stand mixer at medium-low speed for 5 minutes, or until dough is smooth. Increase speed to medium and mix 2 minutes. Kneading is done when dough makes a slapping sound as it hits the side of the bowl. Dough temperature should be close to 90°F. If dough is too sticky add a little more flour."},
						{Text: "Once combined, place dough in oiled mixing bowl and cover with plastic wrap and a warm towel. Let rise 1-2 hours, or until doubled in volume. Time will depend on how warm your house is."},
						{Text: "While dough is rising, make cinnamon sugar filling by stirring together brown sugar and cinnamon in a small bowl. Grease 13 x 9-inch baking pan or round 9.5-inch pie pan with butter."},
						{Text: "Once dough has doubled in size, sprinkle extra flour on your surface and rolling pin and roll dough into 14 x 12-inch rectangle."},
						{Text: "Spread softened butter onto dough with your fingers or a knife."},
						{Text: "Sprinkle cinnamon sugar mixture over butter and press down slightly with your hands."},
						{Text: "Starting at the top, roll the dough toward you into a large log, moving back and forth down the line of dough (in a “typewriter” motion) and always rolling toward you. Roll it tightly as you go so the rolls will be nice and neat. When it’s all rolled, pinch the seam closed and turn the roll over so that the seam is facing down. Cut roll crosswise into 12 1-inch-thick pieces and place on prepared baking pan."},
						{Text: "Cover, and let rise in warm place 45 minutes, or until doubled in size."},
						{Text: "Preheat oven to 350°F. Bake cinnamon rolls for 20-25 minutes, or until golden brown, cooked through and no longer doughy. I baked mine for 22 minutes."},
						{Text: "While cinnamon rolls are cooling, make cream cheese frosting by adding cream cheese, greek yogurt, maple syrup and vanilla to a large mixing bowl. Using a hand mixer on medium speed, whip all the ingredients together until smooth and fluffy, scraping down the sides as needed. Alternatively, you can use a stand mixer."},
						{Text: "Spread frosting over warm cinnamon rolls and serve."},
					},
				},
				Name: "Fluffy Cinnamon Rolls",
				NutritionSchema: &models.NutritionSchema{
					Calories:       "269 kcal",
					Carbohydrates:  "37 g",
					Cholesterol:    "43 mg",
					Fat:            "11 g",
					Fiber:          "1 g",
					Protein:        "6 g",
					SaturatedFat:   "6 g",
					Servings:       "1",
					Sodium:         "252 mg",
					Sugar:          "14 g",
					UnsaturatedFat: "3 g",
				},
				PrepTime:  "PT30M",
				TotalTime: "PT175M",
				Yield:     &models.Yield{Value: 12},
				URL:       "https://www.eatingbirdfood.com/cinnamon-rolls/",
			},
		},
		{
			name: "eatingwell.com",
			in:   "https://www.eatingwell.com/recipe/7887715/lemon-chicken-piccata/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				Cuisine:       &models.Cuisine{Value: ""},
				DateModified:  "2023-09-19T12:07:24.154-04:00",
				DatePublished: "2021-02-04T15:20:10.000-05:00",
				Description: &models.Description{
					Value: "This chicken piccata recipe has a bright, briny flavor, is made from ingredients you likely have on hand, and goes with everything from chicken to tofu to scallops. Bonus: It&#39;s lower in calories than a lot of other pan sauces.",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1.25 pounds boneless, skinless chicken breasts, trimmed",
						".5 teaspoon salt",
						".25 teaspoon ground pepper",
						"2 tablespoons extra-virgin olive oil",
						"1 medium shallot, minced",
						"3 cloves garlic, minced",
						"2 teaspoons all-purpose flour",
						".5 cup low-sodium chicken broth",
						".5 cup dry white wine",
						"2 tablespoons lemon juice",
						"1 tablespoon butter",
						"1 tablespoon capers, rinsed",
						"2 tablespoons chopped fresh parsley",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Remove tenders from chicken and reserve for another use. Place the chicken breasts between 2 pieces of plastic wrap and gently pound with a meat mallet, rolling pin or small skillet to an even thickness of about ½ inch. Pat the chicken dry and sprinkle with salt and pepper."},
						{Text: "Heat oil in a large skillet over medium-high heat. Add the chicken and cook, flipping once, until well browned on both sides, 6 to 8 minutes. Continue to cook, flipping often, until an instant-read thermometer inserted in the thickest part registers 165°F, about 3 minutes more. Transfer to a clean cutting board and tent with foil to keep warm."},
						{Text: "Reduce heat to medium and add shallot to the pan. Cook, stirring often, until softened, 1 to 2 minutes. Add garlic and cook, stirring, until fragrant, about 1 minute. Sprinkle with flour and cook, stirring, for 1 minute. Stir in broth and wine, scraping up any browned bits. Simmer until reduced by half, 3 to 5 minutes. Remove from heat and stir in lemon juice, butter, capers and parsley. Serve the chicken with the sauce."},
					},
				},
				Name: "Lemon Chicken Piccata",
				NutritionSchema: &models.NutritionSchema{
					Calories:       "264 kcal",
					Carbohydrates:  "7 g",
					Cholesterol:    "70 mg",
					Fat:            "13 g",
					Protein:        "24 g",
					SaturatedFat:   "4 g",
					Sodium:         "382 mg",
					Sugar:          "1 g",
					UnsaturatedFat: "0 g",
				},
				TotalTime: "PT20M",
				Yield:     &models.Yield{Value: 4},
				URL:       "https://www.eatingwell.com/recipe/7887715/lemon-chicken-piccata/",
			},
		},
		{
			name: "eatliverun.com",
			in:   "https://www.eatliverun.com/dinner-cheap-garlickly-greens-beans/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Main Course"},
				CookTime:      "PT4M",
				DatePublished: "2017-05-18T05:31:11+00:00",
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1 cup dried lima beans",
						"1 large bunch dark leafy greens (such as collard greens or kale, roughly chopped)",
						"4 tbsp extra virgin olive oil", "3 cloves garlic (minced)", "kosher salt",
						"freshly ground black pepper",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "To quick soak your beans like I usually do, just place beans in a large pot and cover with water. Bring to a high boil and boil for 3 minutes. Turn off burner, cover pot and let beans rest for an hour. After an hour, drain beans and add fresh water with plenty of salt. Bring to a simmer and cook until beans are soft, about 40 minutes."},
						{Text: "After your beans have cooked, drain them gently (lima beans get really soft so be careful not to crush them). Add three tablespoons of the olive oil in a large saute pan over medium heat."},
						{Text: "Add minced garlic and saute for 30 seconds. Add leafy greens and saute for another minute. Stir in beans. Season to taste with salt and pepper and drizzle additional olive oil over top."},
					},
				},
				Name:      "Garlicky Beans & Greens",
				PrepTime:  "PT90M",
				TotalTime: "PT94M",
				Yield:     &models.Yield{Value: 3},
				URL:       "https://www.eatliverun.com/dinner-cheap-garlickly-greens-beans/",
			},
		},
		{
			name: "eatsmarter.com",
			in:   "https://eatsmarter.com/recipes/vietnamese-chicken-cabbage-salad",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Lunch"},
				Cuisine:       &models.Cuisine{Value: "Asian, Vietnamese"},
				DatePublished: "2016-10-07",
				Description: &models.Description{
					Value: "Light and refreshing Vietnamese Chicken Cabbage Salad with crunchy peanuts",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"24 ozs Chicken broth",
						"22 ozs Chicken breasts",
						"3  carrots (10-12 ounces)",
						"1 bunch Radish",
						"1 Red chili pepper",
						"1  Chinese cabbage",
						"1 bunch mixed Fresh herbs (such as mint, basil)",
						"1  Lime",
						"3 Tbsps chopped Peanuts",
						"2 Tbsps vegetable oil",
						"peppers",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Bring broth to a boil, add chicken and simmer for about 20 minutes on low heat."},
						{Text: "Peel carrots and cut into thin, long strips. Rinse and dry radishes, cut into thin slices. Rinse and chop chile pepper."},
						{Text: "Rinse and dry cabbage, cut into fine strips. Rinse and shake dry herbs, pluck off leaves."},
						{Text: "Remove chicken from broth and cool. Squeeze lime juice."},
						{Text: "Cut meat into thin strips, mix with prepared vegetables and herbs, drizzle with lime juice and oil, season with salt and pepper and sprinkle with nuts. Serve."},
					},
				},
				Name:      "Vietnamese Chicken Cabbage Salad",
				TotalTime: "PT40M",
				Yield:     &models.Yield{Value: 4},
				URL:       "https://eatsmarter.com/recipes/vietnamese-chicken-cabbage-salad",
			},
		},
		{
			name: "eatwell101.com",
			in:   "https://www.eatwell101.com/garlic-parmesan-marinated-mushrooms-recipe",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				DateModified:  "2023-06-18T14:01:38+00:00",
				DatePublished: "2023-06-16T22:24:57+00:00",
				Description: &models.Description{
					Value: "This marinated mushroom recipe is the perfect appetizer or side for any occasion. Enjoy the delicious flavor and texture of these mushrooms marinated in olive oil, parmesan, garlic, and herbs! CLICK HERE to Get the Recipes",
				},
				Keywords: &models.Keywords{
					Values: "Garlic Parmesan Marinated Mushrooms Recipe, Easy marinated Mushrooms, how to marinate mushrooms",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1 lb (450g) mushrooms, sliced", "1/3 cup grated parmesan",
						"4 clove garlic, minced", "1/2 teaspoon black pepper",
						"2 tablespoons chopped parsley", "1/2 cup olive oil",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "1. To make the garlic parmesan marinated mushrooms: In a small bowl, combine minced garlic, grated parmesan, chopped parsley, and olive oil."},
						{Text: "2. Arrange the sliced mushrooms on a large shallow plate, and pour the olive oil-parmesan mixture on top. Toss gently to coat the mushrooms and marinate for at least 30 minutes."},
						{Text: "3. Serve the parmesan mushroom dip with toasted bread, crackers, tortilla chips, or pita chips. Enjoy! ❤️"},
						{Text: "Photo credit: © Eatwell101.com"},
						{Text: "More mushroom recipes"},
					},
				},
				Name: "Olive Garlic Parmesan Marinated Mushrooms",
				URL:  "https://www.eatwell101.com/garlic-parmesan-marinated-mushrooms-recipe",
			},
		},
		{
			name: "eatwhattonight.com",
			in:   "https://eatwhattonight.com/2021/11/diced-chicken-with-spicy-chilies-%e8%be%a3%e5%ad%90%e9%b8%a1%e4%b8%81/#wpzoom-recipe-card",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Sides"},
				CookTime:      "PT20M",
				Cuisine:       &models.Cuisine{Value: "Chinese"},
				DatePublished: "2021-11-30T13:14:26+08:00",
				Description:   &models.Description{Value: ""},
				Keywords:      &models.Keywords{Values: "Chicken"},
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"150g diced boneless chicken leg, de-skinned",
						"12 pcs dried chilies, soaked in water to soften, cut into sections",
						"2 tbsp cooking oil",
						"1 tsp Szechuan peppercorns",
						"2 small bulbs of garlic, sliced",
						"5-6 thin slices of ginger",
						"2-3 tbsp water",
						"1 tbsp cooking wine",
						"Some sesame seeds",
						"1 3/4 tsp light soya sauce",
						"1 tbsp ginger juice",
						"1/2 tsp sugar",
						"Pinch of pepper",
						"4 tsp corn flour",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Marinate chicken, cover and set aside in the fridge for an hour."},
						{Text: "Heat up 1 tbsp cooking oil and panfry marinated chicken till cooked and remove from flame (helps to release excess oils from chicken)."},
						{Text: "Add to air-fryer to grill for further 5 mins at 160 degrees C or till browned. Remove and set aside."},
						{Text: "Heat up balance cooking oil and saute ginger and garlic slices."},
						{Text: "Add Szechuan pepper and dried chilies and stir-fry to bring out the aroma."},
						{Text: "Add water if it starts to get too dry. Add chicken pieces and stir-fry to mix well."},
						{Text: "Add cooking wine and stir-fry till chicken are dry. Off heat and sprinkle sesame seeds all over. Mix well and serve hot immediately to enjoy."},
					},
				},
				Name:      "Diced Chicken with Spicy Chilies 辣子鸡丁",
				PrepTime:  "PT10M",
				TotalTime: "PT30M",
				Yield:     &models.Yield{Value: 1},
				URL:       "https://eatwhattonight.com/2021/11/diced-chicken-with-spicy-chilies-%e8%be%a3%e5%ad%90%e9%b8%a1%e4%b8%81/#wpzoom-recipe-card",
			},
		},
		{
			name: "elavegan.com",
			in:   "https://elavegan.com/vegan-moussaka-lentils-gluten-free/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Main Course"},
				CookTime:      "PT50M",
				Cuisine:       &models.Cuisine{Value: "Greek"},
				DatePublished: "2023-11-19T08:13:22+00:00",
				Description: &models.Description{
					Value: "Vegan moussaka with lentils and eggplant! This popular Greek dish can be easily made without meat and still tastes amazing. This healthy casserole is a wonderful comfort meal which is flavorful, satisfying, and very enjoyable. The recipe is plant-based, gluten-free, and fairly easy to make.",
				},
				Keywords: &models.Keywords{Values: "eggplant casserole, lentil moussaka, vegan moussaka"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"2 pounds potatoes (peeled)", "3 large eggplants", "Olive oil (to brush)",
						"Sea salt &amp; pepper (to sprinkle)",
						"3 cups cooked lentils ((about 1 1/2 cups dried))", "2 cups passata",
						"1 cup chopped tomatoes", "1 Tbsp olive oil", "1 large onion (chopped)",
						"2 cloves of garlic (minced)", "2 bay leaves", "1-2 tsp dried thyme",
						"1 tsp oregano", "1 tsp paprika", "1 tsp coconut sugar (or brown sugar)",
						"1 pinch of cinnamon",
						"Sea salt &amp; pepper (to taste)",
						"2 Tbsp vegan butter",
						"2 cups plant-based milk",
						"3 1/2 Tbsp cornstarch",
						"2 Tbsp nutritional yeast",
						"Sea salt &amp; pepper (to taste)",
						"1 pinch of nutmeg",
						"Vegan cheese (to taste (optional))",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Check the video in the post for visual instructions.First, preheat the oven to 390 °F/200 °C and line two baking sheets with parchment paper."},
						{Text: "Then, slice each eggplant into 4-5 lengthwise slices and the potatoes into ½-inch (1 cm) thick slices."},
						{Text: "Once sliced, arrange the slices in a single layer across the baking sheets, lightly brush them with olive oil, and sprinkle with salt and pepper."},
						{Text: "Bake them for 20 minutes or until lightly browned."},
						{Text: "Meanwhile, peel and chop the onion and mince the garlic while heating a little oil in a large skillet over medium heat."},
						{Text: "Once hot, sauté the onion and garlic for 4-5 minutes, until translucent."},
						{Text: "Then add the tomato puree, chopped tomatoes, and all the spices (including salt and pepper to taste). Stir, add the lentils, and let it simmer on low heat for about 5 minutes."},
						{Text: "Add the plant-based milk to a small saucepan/skillet with the cornstarch, nutritional yeast, salt, and pepper, and whisk well."},
						{Text: "Then, add the dairy-free butter and bring the mixture to a boil over medium-high heat. Immediately lower to a simmer and stir/whisk constantly until it thickens. Remove it from the heat and set aside."},
						{Text: "Grease a 9x13-inch (23x33cm) baking dish with vegan butter or oil and arrange half the potato and eggplant slices across the bottom of the dish."},
						{Text: "Top that with the lentil mixture, followed by the remaining eggplant and potato slices."},
						{Text: "Pour the bechamel sauce over the top and spread it evenly. Then, optionally, sprinkle some vegan cheese over the top."},
						{Text: "Transfer the potato eggplant moussaka to the oven to bake for 30 minutes or until golden brown on top with tender eggplant/potato and a bubbling filling. Finally, optionally garnish with herbs, and enjoy!"},
					},
				},
				Name: "Vegan Moussaka",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "400 kcal",
					Carbohydrates: "70 g",
					Fat:           "8 g",
					Fiber:         "22 g",
					Protein:       "19 g",
					SaturatedFat:  "1 g",
					Servings:      "1",
					Sugar:         "16 g",
				},
				PrepTime:  "PT30M",
				TotalTime: "PT80M",
				Yield:     &models.Yield{Value: 6},
				URL:       "https://elavegan.com/vegan-moussaka-lentils-gluten-free/",
			},
		},
		{
			name: "emmikochteinfach.de",
			in:   "https://emmikochteinfach.de/kartoffelgratin-rezept-klassisch-und-einfach/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookTime:      "PT45M",
				DatePublished: "2022-03-16T18:00:40+00:00",
				Description: &models.Description{
					Value: "Das leckere Gratin aus Kartoffeln geht einfacher als man denkt und schmeckt als leckere Beilage zu Fleisch und Fisch oder als vegetarische Hauptspeise mit einem Salat.",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"800 g mittelgroße, &quot;vorwiegend festkochende&quot; Kartoffeln (2-3 mm dicke Scheiben)",
						"200 ml Kochsahne (15% Fett)",
						"100 g französischer Kräuterfrischkäse (z.B. Bresso oder Miree – französische Kräuter)",
						"100 g geriebener Emmentaler (alternativ Comté, Gruyère oder gar kein Käse)",
						"20 g Butter (in Flöckchen, Stückchen)", "1 TL Butter für die Form",
						"1 TL Salz", "0,5 TL weißer oder schwarzer Pfeffer, gemahlen",
						"1/2 TL Muskat (frisch gerieben)",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Den Backofen für das Kartoffelgratin auf 180 Grad Umluft (200 °C Ober- / Unterhitze) vorheizen."},
						{Text: "Die 800 g Kartoffeln abwaschen, schälen, gegebenenfalls nochmals abwaschen und mit einem scharfen Messer oder Hobel in ca. 2-3 mm dicke Scheiben schneiden."},
						{Text: "Eine runde oder rechteckige Auflaufform mit ca. 1 TL Butter einfetten und die Kartoffelscheiben dachziegelartig, reihenartig arrangieren."},
						{Text: "Bei einer runden Auflaufform, wie meiner, arbeitest Du Dich am besten von außen nach innen. Schichtest die Kartoffelscheiben schneckenförmig im Kreis ein. TIPP: Bei rechteckigen Formen, an der schmalen Seite beginnen und reihenartig einschichten."},
						{Text: "Nun verrührst Du mit einem Schneebesen 200 ml Kochsahne, 100 g Kräuterfrischkäse, 1 TL Salz, 0,5 TL weißer oder schwarzer Pfeffer und 0,5 TL frisch geriebene Muskatnuss."},
						{Text: "Die Sahne-Frischkäse-Masse gießt Du nun gleichmäßig über die Kartoffelscheiben und verteilst die 20 g Butter in Flöckchen darüber. Die 100 g Emmentaler reibst Du frisch mit einer Reibe und verteilst ihn ebenfalls gleichmäßig auf den Kartoffeln. Im vorgeheizten Backofen, auf der mittleren Schiene, ca. 45 Minuten backen."},
						{Text: "Wenn Dir der Käse gegen Ende der Backzeit zu braun wird, das Gratin am besten mit Alufolie abdecken."},
						{Text: "Ich wünsche Dir mit meinem Kartoffelgratin Rezept einen guten Appetit!"},
					},
				},
				Name:      "Kartoffelgratin Rezept klassisch und einfach",
				PrepTime:  "PT15M",
				TotalTime: "PT60M",
				Yield:     &models.Yield{Value: 4},
				URL:       "https://emmikochteinfach.de/kartoffelgratin-rezept-klassisch-und-einfach/",
			},
		},
		{
			name: "epicurious.com",
			in:   "https://www.epicurious.com/recipes/food/views/olive-oil-cake",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				DateModified:  "2017-11-01T10:56:00.000-04:00",
				DatePublished: "2017-11-01T10:56:00.000-04:00",
				Description: &models.Description{
					Value: "Even die-hard butter devotees admit that olive oil makes exceptionally good cakes. EVOO is liquid at room temperature, so it lends superior moisture over time. In fact, olive oil cake only improves the longer it sits—this dairy-free version will keep on your counter for days (not that it’ll last that long).",
				},
				Keywords: &models.Keywords{
					Values: "cake,amaretto,vermouth,grand marnier,italian,cake flour,almond flour,lemon,vanilla,snack,breakfast,nut free,baking,stand mixer,bon appétit,web",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1¼ cups plus 2 tablespoons extra-virgin olive oil; plus more for pan",
						"1 cup plus 2 tablespoons sugar; plus more",
						"2 cups cake flour",
						"⅓ cup almond flour or meal or fine-grind cornmeal",
						"2 teaspoons baking powder",
						"½ teaspoon baking soda",
						"½ teaspoon kosher salt",
						"3 tablespoons amaretto, Grand Marnier, sweet vermouth, or other liqueur",
						"1 tablespoon finely grated lemon zest",
						"3 tablespoon fresh lemon juice",
						"2 teaspoons vanilla extract",
						"3 large eggs",
						"A 9\"-diameter springform pan",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Preheat oven to 400°F. Drizzle bottom and sides of pan with oil and use your fingers to coat. Line bottom with a round of parchment paper and smooth to eliminate air bubbles; coat parchment with more oil. Generously sprinkle pan with sugar and tilt to coat in an even layer; tap out excess. Whisk cake flour, almond flour, baking powder, baking soda, and salt in a medium bowl to combine and eliminate any lumps. Stir together amaretto, lemon juice, and vanilla in a small bowl."},
						{Text: "Using an electric mixer on high speed (use whisk attachment if working with a stand mixer), beat eggs, lemon zest, and 1 cup plus 2 Tbsp. sugar in a large bowl until mixture is very light, thick, pale, and falls off the whisk or beaters in a slowly dissolving ribbon, about 3 minutes if using a stand mixer and about 5 minutes if using a hand mixer. With mixer still on high speed, gradually stream in 1¼ cups oil and beat until incorporated and mixture is even thicker. Reduce mixer speed to low and add dry ingredients in 3 additions, alternating with amaretto mixture in 2 additions, beginning and ending with dry ingredients. Fold batter several times with a large rubber spatula, making sure to scrape the bottom and sides of bowl. Scrape batter into prepared pan, smooth top, and sprinkle with more sugar."},
						{Text: "Place cake in oven and immediately reduce oven temperature to 350°F. Bake until top is golden brown, center is firm to the touch, and a tester inserted into the center comes out clean, 40–50 minutes. Transfer pan to a wire rack and let cake cool in pan 15 minutes."},
						{Text: "Poke holes all over top of cake with a toothpick or skewer and drizzle with remaining 2 Tbsp. oil; let it absorb. Run a thin knife around edges of cake and remove ring from pan. Slide cake onto rack and let cool completely. For the best flavor and texture, wrap cake in plastic and let sit at room temperature at least a day before serving."},
						{Text: "Cake can be baked 4 days ahead. Store tightly wrapped at room temperature."},
					},
				},
				Name:  "Olive Oil Cake",
				Yield: &models.Yield{Value: 8},
				URL:   "https://www.bonappetit.com/recipe/olive-oil-cake",
			},
		},
		{
			name: "errenskitchen.com",
			in:   "https://www.errenskitchen.com/baked-or-barbecued-sticky-glazed-ribs/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Main"},
				CookTime:      "PT120M",
				Cuisine:       &models.Cuisine{Value: "American"},
				DatePublished: "2015-07-17T07:32:43+00:00",
				Description: &models.Description{
					Value: "A delightfully sweet and sticky dish, make sure you have napkins at hand!",
				},
				Keywords: &models.Keywords{Values: "bbq ribs recipe, pork ribs recipe, ribs recipe"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"4 lbs spare ribs (cut into individual ribs)", "⅛ cup smoked paprika",
						"⅛ cup sweet paprika", "1 tablespoon garlic powder", "1 teaspoon cumin",
						"1 teaspoons ground black pepper", "5 teaspoons dark brown sugar",
						"1 tablespoon salt", "⅔ cup ketchup", "3 tablespoons soy sauce",
						"1 tablespoon balsamic vinegar", "½ cup honey", "5 tablespoons bourbon",
						"1 teaspoon molasses", "1 teaspoon corn syrup",
						"1 teaspoon smoked paprika (or liquid smoke to taste)",
						"salt (to taste)",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Preheat oven to 300°F/150°C"},
						{Text: "In a small bowl, mix the spices for the rub."},
						{Text: "Cover the ribs with the rub so they’re coated all over."},
						{Text: "Place the ribs on a baking tray, and cover the pan with foil and bake for 2 to 2½ hours, until tender."},
						{Text: "Meanwhile, to make the glaze, add all the ingredients to a saucepan, stir well and bring the mixture to a simmer."},
						{Text: "Simmer for 5 mins until thickened and sticky, taste for seasoning, and add salt to taste, then remove from the heat and set aside until needed."},
						{Text: "At this point, they can be dipped in the glaze and grilled on the BBQ for 10-15 minutes (brushing with the glaze again as needed)"},
						{Text: "When the ribs are done, remove from the oven and increase the heat to 400°F/200°C."},
						{Text: "Using a pair of tongs, dip each rib in the glaze, then return to the rack."},
						{Text: "Place the ribs back in the oven and cook for 10 mins."},
						{Text: "Remove from oven, dip into the glaze again, then return to the oven for another 10-12 minutes until sticky."},
						{Text: "Serve hot."},
					},
				},
				Name: "Baked or Barbecued Sticky Glazed Ribs",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "1539 kcal",
					Carbohydrates: "58 g",
					Cholesterol:   "362 mg",
					Fat:           "107 g",
					Fiber:         "3 g",
					Protein:       "73 g",
					SaturatedFat:  "34 g",
					Servings:      "1",
					Sodium:        "3239 mg",
					Sugar:         "49 g",
				},
				PrepTime:  "PT5M",
				TotalTime: "PT125M",
				Yield:     &models.Yield{Value: 4},
				URL:       "https://www.errenskitchen.com/baked-or-barbecued-sticky-glazed-ribs/",
			},
		},
		{
			name: "expressen.se",
			in:   "https://alltommat.expressen.se/recept/saftiga-choklad--och-apelsinbullar/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    &models.SchemaType{Value: "Recipe"},
				Category:  &models.Category{Value: "uncategorized"},
				Name:      "Saftiga choklad- och apelsinbullar",
				Description: &models.Description{
					Value: `Goda små "fjärilsbullar" med choklad och krämig apelsinfyllning. För att få bullarnas fina fjärilsliknande form skärs degrullen i skivor som trycks ihop i mitten. Spritsa apelsinfyllningen i mitten av varje bulle. Supergott och lyxigt!`,
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"50 g jäst",
						"5 dl mjölk",
						"200 g smör",
						"2 dl strösocker",
						"1 msk hela kardemummakärnor",
						"1.5 tsk salt",
						"16 dl vetemjöl",
						"200 g smör",
						"2 dl strösocker",
						"3 msk kakao",
						"2 tsk vaniljsocker",
						"0.5 dl strösocker",
						"1.25 dl mjölk",
						"0.5 dl apelsin",
						"1 apelsin",
						"17 g majsstärkelse (Maizena)",
						"40 g äggulor",
						"1 krm salt",
						"5 g smör",
						"1 ägg",
						"1 krm salt",
						"droppar vatten",
						"4 msk pärlsocker",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Smula ner jästen i en bunke, gärna tillhörande en köksassistent. Tillsätt mjölken och blanda sedan i resten av ingredienserna. Arbeta degen i ca 10 min. Lägg en bakduk över bunken och låt degen jäsa i 30min. Gör under tiden fyllningen och apelsinkrämen."},
						{Text: "Chokladfyllning: Rör ihop smör, strösocker, kakao och vaniljsocker. Om den känns för tjock kan du mikra den ett par sekunder."},
						{Text: "Apelsinkräm: Blanda ihop alla ingredienserna utom smöret i en kastrull. Låt det sjuda och vispa under tiden. Dra av kastrullen från värmen när krämen börjar tjockna och vispa i smöret. Passera krämen genom en sil. Fyll en spritspåse med apelsinkrämen. Låt den svalna."},
						{Text: "Stjälp upp degen på en mjölad arbetsbänk. Kavla ut degen till en rektangel, 25x65 cm. Bred ut chokladfyllningen på hela ytan. Rulla ihop degen från långsidan till en rulle. Skär rullen i bitar, ca 3 cm breda. Tryck till bitarna på mitten med ett grillspett eller en rund pinne så att snittytorna viks in mot mitten som en fjäril."},
						{Text: "Låt bullarna jäsa under en bakduk i 1–1 ½ timme. Sätt ugnen på 220grader."},
						{Text: "Gör ett hål med fingret i mitten på varje bulle och spritsa i apelsinfyllningen. Vispa ihop ägg, salt och vatten med en gaffel. Pensla bullarna med äggblandningen och strö över pärlsocker. Grädda bullarna i 7–8min, låt dem svalna på ett galler."},
					},
				},
				Keywords: &models.Keywords{Values: "sections/recept"},
				Yield:    &models.Yield{Value: 22},
				URL:      "https://alltommat.expressen.se/recept/saftiga-choklad--och-apelsinbullar/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
