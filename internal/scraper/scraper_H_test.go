package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_H(t *testing.T) {
	testcases := []testcase{
		{
			name: "halfbakedharvest.com",
			in:   "https://www.halfbakedharvest.com/louisiana-style-chicken-and-rice/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Main Course"},
				CookTime:      "PT50M",
				Cuisine:       &models.Cuisine{Value: "American"},
				DatePublished: "2022-03-23T02:00:00+00:00",
				Description: &models.Description{
					Value: "One Skillet Louisiana Style Chicken and Rice: has a variety of flavors and textures, yet it&#39;s all made in ONE skillet with pantry staple ingredients!",
				},
				Keywords: &models.Keywords{Values: "one skillet"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"2 tablespoons extra virgin olive oil",
						"1 pound boneless chicken breasts or thighs",
						"2 tablespoons cajun seasoning",
						"kosher salt and black pepper",
						"6 tablespoons salted butter",
						"1 lemon, sliced",
						"1/2 cup dry broken spaghetti or angel hair pasta",
						"1 cup long grain rice",
						"1 medium yellow onion, sliced",
						"2 bell peppers, sliced",
						"3-4 cups low sodium chicken broth",
						"3 cloves garlic, chopped",
						"1/2 cup fresh tenders herbs, cilantro + parsley",
						"chili flakes",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Preheat the oven to 425° F."},
						{Text: "In a large oven-safe skillet, combine the olive oil, chicken, and cajun seasoning, toss to coat. Set the skillet over high heat. Sear on both sides until golden, 3-5 minutes. During the last 2 minutes of cooking, add 1 tablespoon of butter and lemon slices. Remove everything from the skillet."},
						{Text: "Add the rice and pasta. Cook until the rice is toasted, about 1 minute. Add the onion and peppers and continue to cook another 3-4 minutes, then pour in 3 cups broth. Season with salt and pepper. Bring to a boil."},
						{Text: "Slide the chicken, lemon slices, and any juices left on the plate back into the skillet. Bring to a boil. Cover the skillet and turn the heat down to the lowest setting possible. Allow the rice to cook 10 minutes, until most of the liquid has cooked into the rice, but not all of it. If needed add more broth. Bake, uncovered for 10-15 minutes or until the chicken is cooked through."},
						{Text: "Meanwhile, melt together 5 tablespoons butter, the garlic, and a pinch of chili flakes. Cook until the butter is browning. Stir in the mixed herbs."},
						{Text: "Serve the chicken and rice drizzled with garlic butter and topped with fresh herbs."},
					},
				},
				Name: "Skillet Louisiana Style Chicken and Rice",
				NutritionSchema: &models.NutritionSchema{
					Calories: "547 kcal",
					Servings: "1",
				},
				PrepTime:  "PT20M",
				TotalTime: "PT70M",
				Yield:     &models.Yield{Value: 6},
				URL:       "https://www.halfbakedharvest.com/louisiana-style-chicken-and-rice/",
			},
		},
		{
			name: "handletheheat.com",
			in:   "https://handletheheat.com/peanut-butter-pie/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Dessert"},
				CookTime:      "PT12M",
				Cuisine:       &models.Cuisine{Value: "American"},
				DatePublished: "2023-07-10T03:00:31+00:00",
				Description: &models.Description{
					Value: "This homemade Peanut Butter Pie is made from scratch with just a few ingredients and will have everyone coming back for seconds! It's CRAZY good! With step-by-step video.",
				},
				Keywords: &models.Keywords{Values: "peanut butter cups, peanut butter pie, peanut butter pie recipe"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"14 whole chocolate graham crackers (196 grams)",
						"1 tablespoon light brown sugar",
						"7 tablespoons (99 grams) unsalted butter, (melted)",
						"8 ounces (227 grams) cream cheese, (at room temperature)",
						"3/4 cup (94 grams) powdered sugar plus 2 tablespoons, (divided)",
						"1 cup (270 grams) creamy conventional peanut butter",
						"1 cup (240 grams) heavy whipping cream", "1 teaspoon vanilla extract",
						"Melted peanut butter", "Melted chocolate", "Mini Reese’s cups",
						"Peanut butter chips",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Preheat the oven to 325°F."},
						{Text: "In the bowl of a food processor, process the crackers and sugar until finely ground. Add the butter and pulse until moistened. Use the bottom of a measuring cup, glass, or ramekin to press the crust mixture into the bottom and up the sides of a 9-inch pie plate. Bake until fragrant, about 10 to 12 minutes. Cool completely on a wire rack."},
						{Text: "In a large bowl, use an electric mixer to beat the cream cheese, 3/4 cup powdered sugar, and the peanut butter until light and fluffy, about 3 minutes."},
						{Text: "In a separate bowl, use an electric mixer with the whisk attachment to whip the heavy cream until thick and light. Add in the remaining 2 tablespoons of powdered sugar and the vanilla extract and continue to whip until stiff peaks form."},
						{Text: "Gently fold the whipped cream into the peanut butter mixture. Pour into the prepared pie shell and freeze for 3 hours or chill in the fridge at least 6 hours."},
						{Text: "Place your melted peanut butter and melted chocolate in separate resealable bags or piping bags. Snip off a tiny corner of the bag&#39;s tip. Squeeze slightly to drizzle the melted peanut butter and melted chocolate over your pie. Top with mini Reese&#39;s cups and peanut butter chips. Serve frozen or refrigerated."},
						{Text: "Store in the fridge, covered, for up to 3 days, or in the freezer for up to 1 month."},
					},
				},
				Name:      "Peanut Butter Pie",
				PrepTime:  "PT20M",
				TotalTime: "PT272M",
				Yield:     &models.Yield{Value: 8},
				URL:       "https://handletheheat.com/peanut-butter-pie/",
			},
		},
		{
			name: "hassanchef.com",
			in:   "https://www.hassanchef.com/2022/10/dragon-chicken.html",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "appetizer"},
				CookTime:      "PT10M",
				Cuisine:       &models.Cuisine{Value: "Indo Chinese"},
				DatePublished: "2022-10-14",
				Description: &models.Description{
					Value: "Dragon Chicken an appetizer or snacks of Indian Chinese cuisines where deep fried chicken strips are stir fried with a spicy combination of sauces and herbs",
				},
				Keywords: &models.Keywords{Values: "Dragon Chicken"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1 large Chicken breast",
						"1 teaspoon red chilly paste",
						"1 teaspoon ginger garlic paste",
						"1/4 teaspoon green chilly paste",
						"2 tablespoon cornflour",
						"1 tablespoon beaten egg or egg white",
						"1/2 of each red, yellow and green bell pepper cut into julienne",
						"1/2 of a onion cut into thin slices",
						"Oil for deep frying",
						"1/2 teaspoon chopped garlic",
						"1/2 teaspoon chopped ginger",
						"1 whole red chilly cut into pieces",
						"10 - 12 roasted or golden fried cashew nuts",
						"1/2 teaspoon red chilly sauce",
						"1/2 teaspoon red chilly paste",
						"1/3 teaspoon pepper powder",
						"1/3 teaspoon madras curry powder",
						"Salt as taste",
						"Seasoning powder(optional)",
						"Some chopped green spring onions",
						"1 tablespoon slurry (cornflour)",
						"2 tablespoon cooking oil",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Dragon Chicken is made using boneless chicken breast pieces. Here I have taken only one single chicken breast and sliced it into two parts. Then cut them into thin strips around half inch width. Now add half teaspoon red chilly paste, ginger garlic paste, around 1/4 teaspoon green chilli paste, 2 tablespoon corn flour, 2 tablespoon beaten egg and salt as per your taste. Mix them nicely so that all the chicken strips are get evenly coated. Cover and let them marinated for 5 to 6 minutes to absorb the flavours."},
						{Text: "Cut the bell peppers and capsicum into julienne shape and onion into thin slices. Finely chopped the ginger and garlic. You can avoid to use the bell peppers if not available at your pantry, capsicum and onion also yield a good result. You can either roast the cashew nuts or fry them till golden brown. Roughly chopped the green spring onions and cut the whole red chilly into two or three parts."},
						{Text: "Heat a kadai or small deep vessel pan with enough cooking oil under medium heat flame. When the oil become hot lower the gas flame to low and carefully add the chicken strips one by one into the hot oil. Fry them until become crisp golden from all sides. Remove them on a absorbent paper to absorb any excess oil."},
						{Text: "Heat a non stick pan or Chinese wok with one tablespoon oil under medium flame heat gas fire. Add the bell peppers, capsicum, onion and pieces of whole red chilly. Stir fry them for few seconds in high flame heat. Then add chopped onion, chopped ginger, garlic, green chilly paste and cashew nuts. Stir fry them for one minute."},
						{Text: "Now add red chilly paste, red chilly sauce, tomato sauce and give them a good mix."},
						{Text: "Add some water and cook the sauces for few seconds. Then add soya sauce, white pepper powder, salt and madras curry powder and cook further for few seconds. Now add the fried chicken strips and mix well with the sauces so that the chicken pieces are get well coated with them."},
						{Text: "Pour 2 tablespoon corn flour slurry over the chicken mixture and mix well. Finally sprinkle some chopped spring onions and give a final mix. Remove on a serving plate and serve hot."},
					},
				},
				Name: "Dragon Chicken",
				NutritionSchema: &models.NutritionSchema{
					Calories: "640 cal",
					Fat:      "34 g",
					Servings: "1",
				},
				PrepTime:  "PT20M",
				TotalTime: "PT30M",
				Yield:     &models.Yield{Value: 1},
				URL:       "https://www.hassanchef.com/2022/10/dragon-chicken.html",
			},
		},
		{
			name: "headbangerskitchen.com",
			in:   "https://headbangerskitchen.com/recipe/keto-chicken-adobo/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Main Dish"},
				CookTime:      "PT25M",
				Cuisine:       &models.Cuisine{Value: "Filipino"},
				DatePublished: "2021-10-06T17:23:00+00:00",
				Description: &models.Description{
					Value: "A Keto version of a classic Filipino chicken adobo",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"6 Chicken Thighs (Bone in and Skin On)",
						"1/2 tbsp Avocado Oil",
						"8 cloves garlic (I used small cloves)",
						"1/2 tbsp Whole Black Peppercorns",
						"4 small bay leaves",
						"1/2 tsp Black Pepper Powder",
						"60 ml dark soya sauce or coconut aminos",
						"80 ml cane vinegar or apple cider vinegar",
						"2 tsp Keto sweetener (1:1 sugar substitute)",
						"spring onion greens for garnish",
						"salt as needed",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Get a large, deep skillet on the stove and add a good glug of oil to it. I’m using avocado oil."},
						{Text: "Once the oil is hot, place your chicken thighs in the pan skin side down. Let the chicken cook till the skin is crispy, about 3-4 minutes."},
						{Text: "Using the flat of a knife, smash the eight cloves of garlic."},
						{Text: "Flip the chicken over and add in the aromats – the smashed garlic, the whole bay leaves and peppercorns – and let everything fry till fragrant, about a minute or two."},
						{Text: "Deglaze the pan with a splash of water, then add in the black pepper, soy sauce and apple cider vinegar and give it all a good mix. Flip the chicken again and let everything come to a boil."},
						{Text: "Cook the chicken for about 10 minutes, flipping the pieces halfway through."},
						{Text: "Check the seasoning in the sauce, then add about two teaspoons worth of Keto sweetener. You may need to adjust this depending on the sweetener you’re using, so taste, taste, taste before adding any more."},
						{Text: "After your 10 minutes are up, remove the chicken from the sauce, and let the sauce continue reducing. You want it to become a sticky, syrupy consistency."},
						{Text: "Once the sauce is reduced, add the chicken back in and flip it a few times till it’s completely basted and bathed in that sauce. Finish with some spring onion greens and serve over plain cauli rice."},
					},
				},
				Name:      "Keto Chicken Adobo",
				PrepTime:  "PT5M",
				TotalTime: "PT30M",
				Yield:     &models.Yield{Value: 6},
				URL:       "https://headbangerskitchen.com/recipe/keto-chicken-adobo/",
			},
		},
		{
			name: "heatherchristo.com",
			in:   "https://heatherchristo.com/2020/05/03/13229/",
			want: models.RecipeSchema{
				AtContext:   atContext,
				AtType:      &models.SchemaType{Value: "Recipe"},
				Category:    &models.Category{Value: "uncategorized"},
				CookTime:    "PT12M",
				Description: &models.Description{Value: "Good morning! I am up early on a Sunday morning of what has not really been my finest week. This quarantine thing has definitely been fascinating and I am sure each and every household could be really and truly great control groups for different psychological experiments at this point LOL! But after a few gloomy …"},
				Image:       &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"Lemon Sugar:",
						"¾ cup sugar",
						"2 tablespoons lemon zest",
						"Scones:",
						"6 tablespoons non-dairy milk (I used full-fat oat milk)",
						"2 tablespoons fresh lemon juice",
						"2 cups all-purpose gluten-free baking flour (i used bob's red mill in dark red bag)",
						"½ cup lemon sugar",
						"2 teaspoons baking powder",
						"½ teaspoon baking soda",
						"½ teaspoon xanthum gum",
						"½ teaspoon kosher salt",
						"6 tablespoons vegan butter or butter flavored shortening",
						"½ cup blueberries",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Preheat the oven to 400 degrees and prepare a sheet pan with a silpat or with parchment paper."},
						{Text: "For the Lemon Sugar:"},
						{Text: "In a small bowl, combine the sugar and the 2 Tablespoons of finely grated zest of lemon (this takes about two large lemons and I like to use a rasp.) Gently stir together- the sugar will get clumpy, and set aside."},
						{Text: "For the Scones:"},
						{Text: "In a small bowl, combine the non-dairy milk and the lemon juice and set aside."},
						{Text: "In the bowl of a food processor, add all of the dry ingredients, including the ½ cup of the lemon sugar. Pulse everything together until well combined."},
						{Text: "Add in the 6 tablespoons of vegan butter or shortening and pulse until you have a coarse crumb, about 15-20 seconds. Then with the processor running, stream in the milk and lemon juice combination and pulse until the dough is uniform. You will have a wet sticky dough."},
						{Text: "I dropped the dough by big “rustic” spoonfuls onto the prepared baking sheets- I got 6 large scones. Then I stuck (with my finger pushing them in!) about 8 blueberries into each scone, making sure some of them were all the way inside the scone, and some were sticking out of the top or sides."},
						{Text: "Generously divide the remaining ¼ cup of lemon sugar over the tops of the scones and then bake them for 12-13 minutes until golden around the edges. I let them cool on the sheet pan for a few minutes before serving them. THey are wonderful all on their own, or served with berry jam."},
					},
				},
				Name:     "Lemon Blueberry Scones with Lemon Sugar (Vegan and Gluten-free)",
				PrepTime: "PT15M",
				Yield:    &models.Yield{Value: 6},
				URL:      "https://heatherchristo.com/2020/05/03/13229/",
			},
		},
		{
			name: "hellofresh.com",
			in:   "https://www.hellofresh.com/recipes/creamy-shrimp-tagliatelle-5a8f0fcbae08b52f161b5832",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "main course"},
				Cuisine:       &models.Cuisine{Value: "American"},
				DatePublished: "2018-02-22T18:45:31+00:00",
				Description: &models.Description{
					Value: "Pronto! Pronto! You can make this dinner recipe with the lightning speed of an Italian race car. Thanks to fresh tagliatelle, which cooks faster than the dried kind, you arrive at al dente perfection in a matter of minutes. The shrimp and heirloom tomatoes only need a quick toss in the pan, too, becoming tender on the count of uno, due, tre.",
				},
				Keywords: &models.Keywords{Values: "Spicy,Dinner Ideas"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"2 clove Garlic",
						"2 unit Scallions",
						"1 unit Chili Pepper",
						"10 ounce Heirloom Grape Tomatoes",
						"1 unit Lemon",
						"10 ounce Shrimp",
						"9 ounce Tagliatelle Pasta",
						"4 tablespoon Sour Cream",
						"1 teaspoon Olive Oil",
						"2 tablespoon Butter",
						"Salt",
						"Pepper",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Wash and dry all produce. Bring a large pot of salted water to a boil. Mince garlic. Trim, then thinly slice scallions, keeping greens and whites separate. Finely mince chili, removing seeds and ribs for less heat. Halve tomatoes. Cut lemon into wedges. Rinse shrimp and pat dry with a paper towel."},
						{Text: "Heat a drizzle of olive oil in large pan over medium-high heat. Add garlic, scallion whites, and chili (to taste). Cook until fragrant, about 30 seconds. Add shrimp and cook, tossing, until starting to turn pink but not quite cooked through, 1-2 minutes. Season with salt and pepper."},
						{Text: "Once water is boiling, add tagliatelle to pot. (TIP: If any noodles are stuck together, separate them first.) Cook, stirring occasionally, until al dente, 4-5 minutes. Carefully scoop out and reserve ¼ cup pasta cooking water, then drain."},
						{Text: "Meanwhile, add tomatoes to pan with shrimp. Cook, tossing, until wilted and juicy, 2-3 minutes. Season with salt and pepper. Remove from heat and set aside until pasta is ready. TIP: If you like it extra hot, add any remaining chili (to taste) at this point."},
						{Text: "Once tagliatelle is done cooking, return pan with shrimp and tomatoes to medium heat and add tagliatelle and 2 TBSP butter. Toss to combine and melt butter. Season with salt and pepper."},
						{Text: "Remove pan from heat and stir in sour cream, a squeeze of lemon, and as much pasta cooking water as needed to reach a saucy consistency. Season with salt and pepper. Divide between plates or bowls and garnish with scallion greens. Serve with lemon wedges on the side for squeezing over."},
					},
				},
				Name: "Creamy Shrimp Tagliatelle with Heirloom Tomatoes, Garlic, and Chili",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "750 kcal",
					Carbohydrates: "86 g",
					Cholesterol:   "350 mg",
					Fat:           "27 g",
					Fiber:         "5 g",
					Protein:       "50 g",
					SaturatedFat:  "12 g",
					Sodium:        "880 mg",
					Sugar:         "9 g",
				},
				TotalTime: "PT20M",
				Yield:     &models.Yield{Value: 2},
				URL:       "https://www.hellofresh.com/recipes/creamy-shrimp-tagliatelle-5a8f0fcbae08b52f161b5832",
			},
		},
		{
			name: "homechef.com",
			in:   "https://www.homechef.com/meals/farmhouse-fried-chicken",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    &models.SchemaType{Value: "Recipe"},
				Category:  &models.Category{Value: "uncategorized"},
				Image:     &models.Image{Value: anUploadedImage.String()},
				Name:      "Farmhouse Fried Chicken",
				URL:       "https://www.homechef.com/meals/farmhouse-fried-chicken",
				Description: &models.Description{
					Value: "This stick-to-your-ribs satisfying country classic is an indulgence you've earned. The crispy comfort that only fried chicken can supply is accompanied by mashed potatoes and sweet corn. While that “other” chicken has you eating out of a bucket, this homey treat transports you to an idyllic country farmhouse on the prairie. Yee-Haw! Tip: Make potatoes easier to cut by making a large slice that leaves a flat surface. Place the flat surface on the cutting board and get to dicing!",
				},
				Yield: &models.Yield{Value: 2},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: `1 Make the Mashed Potatoes Cut potato into 1/2" pieces.Bring a small pot with potato pieces and enough water to cover to a boil. Reduce to a simmer and cook until fork-tender, 12-15 minutes.Drain potatoes in a colander and return to pot. Add half the butter, 1/4 the cream (reserve remaining of each for gravy), 1/2 tsp. olive oil, and a pinch of salt. Mash until desired consistency is reached. Cover and set aside.While potato cooks, prepare ingredients.`},
						{Text: "2 Prepare the Ingredients Trim and thinly slice green onions on an angle.Heat canola oil in a medium pan over medium heat, 5 minutes.While oil heats, pat chicken breasts dry, and season both sides with a pinch of pepper.Combine mayonnaise and 2 tsp. water in a mixing bowl. Place chicken breading in another mixing bowl.Dip one chicken breast in mayonnaise-water mixture, then coat completely in chicken breading, shaking off any excess. Repeat with second chicken breast."},
						{Text: "3 Fry the Chicken Line a plate with a paper towel. Test oil temperature by adding a pinch of chicken breading to it. It should sizzle gently. If it browns immediately, turn heat down and let oil cool. If it doesn't brown, increase heat.Lay chicken breasts away from you in hot oil and flip every 3-5 minutes until golden brown and chicken reaches a minimum internal temperature of 165 degrees, 10-14 minutes.Transfer chicken to towel-lined plate. Rest at least 5 minutes.While chicken rests, cook corn."},
						{Text: "4 Cook the Corn Place another small pot over medium heat. Add 1 tsp. olive oil and corn to hot pot. Stir occasionally until warmed through, 4-5 minutes.Transfer corn to a plate and season with a pinch of salt and pepper.Wipe pot clean and reserve."},
						{Text: `5 Make Gravy and Finish Dish Return pot used to cook corn to medium heat. Add green onions (reserve a pinch for garnish) and remaining cream and bring to a simmer. Once simmering, stir often until slightly thickened, 3-5 minutes.Remove from burner and swirl in remaining butter. Season with a pinch of pepper.If desired, slice chicken into 1/2" pieces.Plate dish as pictured on front of card, pouring gravy over chicken and garnishing potatoes with reserved green onions. Bon appétit!`},
					},
				},
				Ingredients: &models.Ingredients{
					Values: []string{
						"2 Russet Potatoes",
						"13 oz. Boneless Skinless Chicken Breasts",
						"6 fl. oz. Canola Oil",
						"4 oz. Light Cream",
						"½ cup Chicken Breading",
						"3 oz. Corn Kernels",
						"2 Green Onions",
						"0.84 oz. Mayonnaise",
						"⅗ oz. Butter",
					},
				},
			},
		},
		{
			name: "hostthetoast.com",
			in:   "https://hostthetoast.com/guinness-beef-stew-with-cheddar-herb-dumplings/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookTime:      "PT3H",
				DatePublished: "2014-03-18",
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"¼ pound bacon",
						"2 pounds boneless beef chuck, chopped into bite-sized pieces",
						"Kosher salt and black pepper",
						"4 sticks celery, chopped",
						"3 large carrots, chopped",
						"1 large onion, chopped",
						"4 cloves garlic, minced",
						"2 large potatoes or parsnips, diced",
						"1 turnip, diced",
						"3 ounces tomato paste",
						"1 (12 ounce) bottle Guinness",
						"4 cups low sodium chicken broth",
						"2 tablespoons Worcestershire sauce",
						"1 bay leaf",
						"3 sprigs thyme",
						"1 tablespoon cornstarch, or as needed",
						"½ pound cremini mushrooms, sliced (optional)",
						"Chopped parsley",
						"1 ½ cups self-rising flour",
						"1/2 teaspoon garlic powder",
						"1/3 cup shortening",
						"3/4 cup shredded Irish sharp cheddar",
						"2/3 cup milk",
						"2 tablespoons mixed fresh herbs such as parsley, chives, and thyme, chopped",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToStep{
						{Text: "Cook the bacon in a large, oven-safe, heavy-based pot or high-walled saute pan over medium heat."},
						{Text: "Remove the bacon, crumble, and set aside, but leave the bacon fat in the pot. Season the beef with salt and pepper and fry in the bacon fat until browned on all sides. Remove the beef from the pan and set aside."},
						{Text: "In the same pot, fry the onion, celery, and carrots until soft and fragrant, adding a little oil if necessary."},
						{Text: "Add garlic and fry for another 30 seconds. Stir in the tomato paste."},
						{Text: "Pour in the Guinness and Worcestershire sauce. Allow to come to a simmer and stir with a wooden spoon, scraping up the browned bits from the bottom of the pot."},
						{Text: "Add the beef back to the pot and pour in the chicken broth. Add the bay leaf and thyme."},
						{Text: "Reduce to a simmer and cover. Simmer for 1 1/2 hours. Add the potatoes or parsnips and the turnip. Simmer for another ½ hour, or until the vegetables are tender."},
						{Text: "Remove the bay leaf and thyme branches. If the stew is still thin, mix a tablespoon of cornstarch with a tablespoon of cold water to form a slurry. Mix the slurry into the stew and bring the mixture to a boil. Reduce to a simmer again, stirring occasionally, and add in the mushrooms if desired. Cook for 10 minutes, uncovered, until the stew thickens and the mushrooms are cooked through. Stir the bacon back in. Preheat the oven to 350°F."},
						{Text: "Stir together the self-rising flour and garlic powder in a medium bowl. Cut in the shortening until mixture resembles coarse crumbs. Stir in the cheddar cheese, then add the milk and stir until the dry ingredients are moistened."},
						{Text: "Make small balls with the dough and place them on top of the stew, leaving them room to expand-- they grow a lot as they cook. Place the stew in the oven uncovered and bake until the dumplings are browned and cooked through, about 30 to 40 minutes."},
						{Text: "Garnish the stew with parsley and serve."},
					},
				},
				Name:      "Guinness Beef Stew with Cheddar Herb Dumplings",
				PrepTime:  "PT30M",
				TotalTime: "PT3H30M",
				Yield:     &models.Yield{Value: 6},
				URL:       "https://hostthetoast.com/guinness-beef-stew-with-cheddar-herb-dumplings/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
