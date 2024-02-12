package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_B(t *testing.T) {
	testcases := []testcase{
		{
			name: "bakingmischief.com",
			in:   "https://bakingmischief.com/italian-roasted-potatoes/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Side"},
				CookTime:      "PT40M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-02-23T01:00:11+00:00",
				Description: models.Description{
					Value: "These 3-ingredient Italian roasted potatoes are quick and simple to prep. With crispy edges " +
						"and creamy centers, they make an easy side dish that everyone will love.",
				},
				Keywords: models.Keywords{Values: "Italian roasted potatoes"},
				Image: models.Image{
					Value: "https://bakingmischief.com/wp-content/uploads/2021/11/italian-roasted-potatoes-image-square-3.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 pounds red potatoes (*)", "3 tablespoons olive oil",
						"2 teaspoons Italian seasoning (*)", "½ teaspoon salt",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 425°F. Scrub potatoes, remove any blemishes, and cut them into 2-inch pieces.",
						"Pile potatoes on a baking sheet and drizzle with olive oil and sprinkle Italian seasoning and salt " +
							"over the top. Toss until well-coated. Arrange potatoes so they are evenly spaced over the " +
							"baking sheet with a cut side down.",
						"Cover the tray tightly with foil and bake on the center rack of your oven for 15 minutes.",
						"Remove and discard the foil. Raise the oven temperature to 475°F.",
						"Continue to bake the potatoes uncovered for 25 to 30 minutes, rotating the pan once halfway through, " +
							"until the potatoes are fork-tender.", "Remove from the oven, add additional salt if needed, and enjoy!",
					},
				},
				Name: "Italian Roasted Potatoes",
				NutritionSchema: models.NutritionSchema{
					Calories: "194 kcal",
					Servings: "1",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://bakingmischief.com/italian-roasted-potatoes/",
			},
		},
		{
			name: "baking-sense.com",
			in:   "https://www.baking-sense.com/2022/02/23/irish-potato-farls/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Breakfast"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "Irish"},
				DatePublished: "2022-02-23T10:00:00+00:00",
				Description: models.Description{
					Value: "Have you heard of Irish Potato Farls? No? Well, if you love potato pancakes, you&#39;ll love " +
						"potato farls. They&#39;re easy to make with fresh or left over potatoes.",
				},
				Image: models.Image{
					Value: "https://www.baking-sense.com/wp-content/uploads/2022/02/potato-farls-featured.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"24 oz russet potatoes (peeled and cut into 1&quot; cubes.)",
						"1 1/2 teaspoons table salt (divided)",
						"2 oz Irish Butter (room temperature, divided)",
						"3.75 oz all purpose flour (3/4 cup, see note)", "1/2 teaspoon baking powder",
						"1/4 teaspoon ground black pepper", "2 scallions (chopped fine)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Place the potatoes in a pot of water with 1/2 teaspoon of the salt. Boil until the potatoes are tender. Drain.",
						"Pass the potatoes through a ricer back into the pot, or use a potato masher. Add 1 oz (2 tablespoons) " +
							"of the butter and the remaining salt mix until the butter is melted. Add the flour, baking powder " +
							"and pepper and stir until most of the flour is mixed in.",
						"Turn the dough out onto a lightly floured surface and knead in the chopped scallions then form the dough " +
							"into a ball. Divide the dough in half.",
						"Preheat a large skillet over medium heat. While the pan is heating, pat each half of the dough to " +
							"an 8\" round, 1/4”thick. Use flour as needed to prevent sticking. Cut the rounds into quarters. " +
							"You&#39;ll have a total of 8 farls. (See Note)",
						"Melt 1 tablespoon of the remaining butter in the pan. Fry half the farls in the butter until golden brown, " +
							"then flip and fry the other side. Cook until both sides are golden brown and the farl springs " +
							"back when pressed in the center. About 4 minutes per side.",
						"Repeat with the remaining butter and farls. Serve immediately.",
					},
				},
				Keywords: models.Keywords{Values: "pancake, potato"},
				Name:     "Irish Potato Farls Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories:       "165 kcal",
					Carbohydrates:  "25 g",
					Cholesterol:    "15 mg",
					Fat:            "6 g",
					Fiber:          "2 g",
					Protein:        "3 g",
					SaturatedFat:   "4 g",
					Servings:       "1",
					Sodium:         "519 mg",
					Sugar:          "1 g",
					TransFat:       "1 g",
					UnsaturatedFat: "3 g",
				},
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://www.baking-sense.com/2022/02/23/irish-potato-farls/",
			},
		},
		{
			name: "barefootcontessa.com",
			in:   "https://barefootcontessa.com/recipes/brussels-sprouts-lardons",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Sides"},
				DateModified:  "2023-11-14T11:00:31-05:00",
				DatePublished: "2016-07-20T14:17:00-04:00",
				Description: models.Description{
					Value: "Brussels Sprouts Lardons from Barefoot Contessa.",
				},
				Image: models.Image{
					Value: "https://d14iv1hjmfkv57.cloudfront.net/assets/recipes/brussels-sprouts-lardons/_1200x630_crop_center-center_82_none/155-web-horizon.jpg?v=1704301339",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 tablespoons good olive oil",
						"6 ounces Italian pancetta or bacon, 1/4-inch diced",
						"1-1/2 pounds Brussels sprouts (2 containers), trimmed and cut in half",
						"3/4 teaspoon kosher salt",
						"3/4 teaspoon freshly ground black pepper",
						"3/4 cup golden raisins",
						"1-3/4 cups homemade chicken stock or canned broth",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Heat the olive oil in a large (12-inch) sauté pan and add the pancetta. Cook over medium heat, stirring often, until the fat is rendered and the pancetta is golden brown and crisp, 5 to 10 minutes. Remove the pancetta to a plate lined with a paper towel.",
						"Add the Brussels sprouts, salt, and pepper to the fat in the pan and sauté over medium heat for about 5 minutes, until lightly browned. Add the raisins and chicken stock. Lower the heat and cook uncovered, stirring occasionally, until the sprouts are tender when pierced with a knife, about 15 minutes. If the skillet becomes too dry, add a little chicken stock or water. Return the pancetta to the pan, heat through, season to taste, and serve.",
					},
				},
				Name:  "Brussels Sprouts Lardons",
				Yield: models.Yield{Value: 6},
				URL:   "https://barefootcontessa.com/recipes/brussels-sprouts-lardons",
			},
		},
		{
			name: "bbc.co.uk",
			in:   "https://www.bbc.co.uk/food/recipes/healthy_sausage_16132",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Category:  models.Category{Value: "Main course"},
				CookTime:  "PT2H",
				Description: models.Description{
					Value: "Never thought you’d hear the words ‘healthy sausage casserole’? Well, here it is. This " +
						"all-in-one-dish dinner is packed with veggies for a perfect midweek meal.\r\n\r\nEach " +
						"serving provides 348 kcal, 34g protein, 33.5g carbohydrates (of which 18g sugars), 7g " +
						"fat (of which 2g saturates), 9.5g fibre and 1.8g salt.",
				},
				Keywords: models.Keywords{
					Values: "absolute bangers, 400-calorie dinners, cheap stews , comfort food on a budget, easy healthy dinner ideas, easy sausage suppers , healthy and filling, healthy british classics, healthy comfort food, healthy dinner, healthy family meals, healthy meals on a budget, healthy winter food, low-calorie comfort food, low-calorie, making meat go further, sausage suppers, summery sausages, the best sausage, winter stew, autumn, bonfire night, easy family dinners, winter, sausage casserole, sausage, healthy",
				},
				Image: models.Image{
					Value: "https://food-images.files.bbci.co.uk/food/recipes/healthy_sausage_16132_16x9.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 red peppers, seeds removed, cut into chunks",
						"2 carrots, cut into thick slices",
						"2 red onions, cut into wedges",
						"5 garlic cloves, finely sliced",
						"8 lean sausages",
						"400g tin peeled cherry tomatoes",
						"400g tin chickpeas, drained",
						"200ml/7fl oz vegetable stock",
						"1 green chilli, seeds removed, chopped",
						"1 tsp paprika",
						"2 tsp French mustard",
						"100g/3½oz frozen mixed vegetables",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat the oven to 220C/200C Fan/Gas 7.",
						"Put the peppers, carrots, onions and garlic into a deep baking dish and roast for 20 minutes. Add the " +
							"sausages and roast for a further 10 minutes.",
						"Turn the oven down to 200C/180C Fan/Gas 6. Pour the tomatoes and chickpeas into the baking dish, then " +
							"stir in the stock, chilli and paprika. Bake for another 35 minutes.",
						"Stir in the mustard and the frozen mixed veg and return to the oven for 5 minutes. Leave to rest for " +
							"10 minutes before serving.",
					},
				},
				Name: "Healthy sausage casserole",
				NutritionSchema: models.NutritionSchema{
					Calories:      "348kcal",
					Carbohydrates: "33.5g",
					Fat:           "7g",
					Fiber:         "9.5g",
					Protein:       "34g",
					SaturatedFat:  "2g",
					Sugar:         "18g",
				},
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.bbc.co.uk/food/recipes/healthy_sausage_16132",
			},
		},
		{
			name: "bbcgoodfood.com",
			in:   "https://www.bbcgoodfood.com/recipes/three-cheese-risotto",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner, Main course, Side dish, Supper"},
				CookTime:      "PT35M",
				Cuisine:       models.Cuisine{Value: "Italian"},
				DateModified:  "2023-12-15T11:03:00+00:00",
				DatePublished: "2014-12-04T12:17:29+00:00",
				Description: models.Description{
					Value: "Tom Kerridge's indulgently rich and cheesy risotto makes an extra-special side dish for a celebration dinner party",
				},
				Keywords: models.Keywords{
					Values: "cheesy risotto, Indulgent, rice side dish, risotto side dish, Tom Kerridge, Winter",
				},
				Image: models.Image{
					Value: "https://images.immediate.co.uk/production/volatile/sites/30/2020/08/roast-poussin-with-wild-mushroom-sauce_0-8051af1.jpg?resize=768,574",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"25g butter",
						"1 tbsp olive oil",
						"1 onion , finely chopped",
						"2 garlic cloves , finely grated",
						"200g risotto rice",
						"200ml white wine",
						"800ml warm chicken stock",
						"50g fresh parmesan (or vegetarian alternative)",
						"½ ball of mozzarella , diced",
						"pinch of cayenne pepper , to taste",
						"2 tbsp mascarpone",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Melt the butter and olive oil together in a large, shallow saucepan, add the onion and garlic, and cook for 5-10 mins until soft. Add the risotto rice and cook for 2-3 mins, getting a good covering in the fats and giving the rice a slightly toasted flavour.",
						"Add the white wine and cook until it has reduced away. Add the warm chicken stock, a ladleful at a time, and stir into the rice – when it has been absorbed, add more. You may not need to add all the stock, but keep adding until the rice is cooked al dente. It will take around 15 mins to get the risotto to the right consistency.",
						"Take the rice pan off the heat and stir in the cheeses, season and leave to rest for 3-4 mins. Serve with the roasted poussins, morel sauce and some wilted Baby Gem lettuce leaves.",
					},
				},
				Name: "Three-cheese risotto",
				NutritionSchema: models.NutritionSchema{
					Calories:      "451 calories",
					Carbohydrates: "42 grams",
					Fat:           "20 grams",
					Fiber:         "2 grams",
					Protein:       "17 grams",
					SaturatedFat:  "12 grams",
					Servings:      "1",
					Sodium:        "0.9 milligram",
					Sugar:         "4 grams",
				},
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 1},
				URL:      "https://www.bbcgoodfood.com/recipes/three-cheese-risotto",
			},
		},
		{
			name: "bbcgoodfood.com",
			in:   "https://www.bbcgoodfood.com/recipes/three-cheese-risotto",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner, Main course, Side dish, Supper"},
				CookTime:      "PT35M",
				Cuisine:       models.Cuisine{Value: "Italian"},
				DateModified:  "2023-12-15T11:03:00+00:00",
				DatePublished: "2014-12-04T12:17:29+00:00",
				Description: models.Description{
					Value: "Tom Kerridge's indulgently rich and cheesy risotto makes an extra-special side dish for a celebration dinner party",
				},
				Keywords: models.Keywords{
					Values: "cheesy risotto, Indulgent, rice side dish, risotto side dish, Tom Kerridge, Winter",
				},
				Image: models.Image{
					Value: "https://images.immediate.co.uk/production/volatile/sites/30/2020/08/roast-poussin-with-wild-mushroom-sauce_0-8051af1.jpg?resize=768,574",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"25g butter",
						"1 tbsp olive oil",
						"1 onion , finely chopped",
						"2 garlic cloves , finely grated",
						"200g risotto rice",
						"200ml white wine",
						"800ml warm chicken stock",
						"50g fresh parmesan (or vegetarian alternative)",
						"½ ball of mozzarella , diced",
						"pinch of cayenne pepper , to taste",
						"2 tbsp mascarpone",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Melt the butter and olive oil together in a large, shallow saucepan, add the onion and garlic, and cook for 5-10 mins until soft. Add the risotto rice and cook for 2-3 mins, getting a good covering in the fats and giving the rice a slightly toasted flavour.",
						"Add the white wine and cook until it has reduced away. Add the warm chicken stock, a ladleful at a time, and stir into the rice – when it has been absorbed, add more. You may not need to add all the stock, but keep adding until the rice is cooked al dente. It will take around 15 mins to get the risotto to the right consistency.",
						"Take the rice pan off the heat and stir in the cheeses, season and leave to rest for 3-4 mins. Serve with the roasted poussins, morel sauce and some wilted Baby Gem lettuce leaves.",
					},
				},
				Name: "Three-cheese risotto",
				NutritionSchema: models.NutritionSchema{
					Calories:      "451 calories",
					Carbohydrates: "42 grams",
					Fat:           "20 grams",
					Fiber:         "2 grams",
					Protein:       "17 grams",
					SaturatedFat:  "12 grams",
					Servings:      "1",
					Sodium:        "0.9 milligram",
					Sugar:         "4 grams",
				},
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 1},
				URL:      "https://www.bbcgoodfood.com/recipes/three-cheese-risotto",
			},
		},
		{
			name: "bbcgoodfood2.com",
			in:   "https://www.bbcgoodfood.com/recipes/pan-fried-salmon",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner, Main course, Supper"},
				CookTime:      "PT5M",
				DateModified:  "2023-10-03T17:51:48+01:00",
				DatePublished: "2019-02-17T17:02:00+00:00",
				Description: models.Description{
					Value: "Have a perfectly cooked salmon fillet, complete with crisp skin, ready in under 10 minutes. Serve with a side of buttery, seasonal green veg for a filling supper",
				},
				Keywords: models.Keywords{
					Values: "Dinner, Easy, easy salmon, Esther Clark, Fish, Gluten free, how to cook salmon, Omega 3, pan fried salmon recipe, Pan-fried salmon, Salmon, Salmon recipe",
				},
				Image: models.Image{
					Value: "https://images.immediate.co.uk/production/volatile/sites/30/2020/08/salmon-1547b3f.jpg?resize=768,574",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 x 150g salmon fillets (about 4cm thick), skin on",
						"½ tbsp olive oil",
						"20g unsalted butter",
						"½ lemon, juiced",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Generously season the salmon fillets with salt and pepper. Put the oil and butter in a non-stick frying pan over a medium heat, swirling around the pan until melted and foaming, then turn up the heat. Once the butter starts bubbling, add the salmon fillets to the pan, skin-side-down, and fry for 3 mins until crisp. Flip the fillets over, lower the heat and cook for 2 mins more, then drizzle with the lemon juice. Transfer the salmon to a plate and baste with any of the buttery juices left in the pan.",
					},
				},
				Name: "Pan-fried salmon",
				NutritionSchema: models.NutritionSchema{
					Calories:      "524 calories",
					Carbohydrates: "0.3 grams",
					Fat:           "44 grams",
					Fiber:         "0.3 grams",
					Protein:       "31 grams",
					SaturatedFat:  "15 grams",
					Servings:      "2",
					Sodium:        "0.17 milligram",
					Sugar:         "0.3 grams",
				},
				PrepTime: "PT1M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://www.bbcgoodfood.com/recipes/pan-fried-salmon",
			},
		},
		{
			name: "bettycrocker.com",
			in:   "https://www.bettycrocker.com/recipes/spinach-mushroom-quiche/ed3014db-7810-41d6-8e1c-cd4eed7b1db3",
			want: models.RecipeSchema{
				AtContext:    atContext,
				AtType:       models.SchemaType{Value: "Recipe"},
				Category:     models.Category{Value: "Breakfast"},
				Cuisine:      models.Cuisine{Value: "French"},
				DateCreated:  "2011-10-05",
				DateModified: "2013-04-03",
				Description: models.Description{
					Value: "Bisquick® Gluten Free mix crust topped with spinach and mushroom mixture for a tasty breakfast – " +
						"perfect if you love French cuisine.",
				},
				Keywords: models.Keywords{Values: "spinach mushroom quiche"},
				Image: models.Image{
					Value: "https://mojo.generalmills.com/api/public/content/MAHdJv1NBUeLl4-jtMq24g_gmi_hi_res_jpeg.jpeg?v=2e1b9203&t=b5673970ed9e41549a020b29d456506d",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 cup Bisquick™ Gluten Free mix",
						"1/3 cup plus 1 tablespoon shortening",
						"3 to 4 tablespoons cold water",
						"1 tablespoon butter",
						"1 small onion, chopped (1/3 cup)",
						"1 1/2 cups sliced fresh mushrooms (about 4 oz)",
						"4 eggs",
						"1 cup milk",
						"1/8 teaspoon ground red pepper (cayenne)",
						"3/4 cup coarsely chopped fresh spinach",
						"1/4 cup chopped red bell pepper",
						"1 cup gluten-free shredded Italian cheese blend (4 oz)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Heat oven to 425°F. In medium bowl, place Bisquick mix. Cut in shortening, using pastry blender " +
							"(or pulling 2 table knives through ingredients in opposite directions), until particles are size " +
							"of small peas. Sprinkle with cold water, 1 tablespoon at a time, tossing with fork until all flour " +
							"is moistened and pastry almost leaves side of bowl (1 to 2 teaspoons more water can be added if necessary).",
						"Press pastry in bottom and up side of ungreased 9-inch quiche dish or glass pie plate. Bake 12 to " +
							"14 minutes or until crust just begins to brown and is set. Reduce oven temperature to 325°F.",
						"Meanwhile, in 10-inch skillet, melt butter over medium heat. Cook onion and mushrooms in butter about " +
							"5 minutes, stirring occasionally, until tender. In medium bowl, beat eggs, milk and red pepper until " +
							"well blended. Stir in spinach, bell pepper, mushroom mixture and cheese. Pour into partially baked crust.",
						"Bake 40 to 45 minutes or until knife inserted in center comes out clean. Let stand 10 minutes before cutting.",
					},
				},
				Name: "Spinach Mushroom Quiche",
				NutritionSchema: models.NutritionSchema{
					Calories:       "260",
					Carbohydrates:  "16 g",
					Cholesterol:    "120 mg",
					Fat:            "2 ",
					Fiber:          "0 g",
					Protein:        "9 g",
					SaturatedFat:   "6 g",
					Servings:       "1",
					Sodium:         "340 mg",
					Sugar:          "4 g",
					TransFat:       "2 g",
					UnsaturatedFat: "",
				},
				PrepTime: "PT0H30M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://www.bettycrocker.com/recipes/spinach-mushroom-quiche/ed3014db-7810-41d6-8e1c-cd4eed7b1db3",
			},
		},
		{
			name: "biancazapatka.com",
			in:   "https://biancazapatka.com/en/vegan-stuffed-peppers/#recipe",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Lunch & Dinner"},
				CookTime:      "PT45M",
				Cuisine:       models.Cuisine{Value: "Mexican"},
				DateCreated:   "",
				DateModified:  "",
				DatePublished: "2022-06-09T13:47:24+00:00",
				Description: models.Description{
					Value: "These oven-baked Mexican stuffed peppers with rice, beans, vegetables and vegan cheese is a popular meatless recipe that’s easy to make, healthy, gluten-free and so delicious! Serve with vegan aioli, you’ll have the perfect plant-based lunch or dinner for the whole family!",
				},
				Keywords: models.Keywords{
					Values: "Mexican Rice, Mexican Stuffed Peppers, Stuffed Peppers, Vegan Stuffed Peppers",
				},
				Image: models.Image{
					Value: "https://biancazapatka.com/wp-content/uploads/2022/03/stuffed-peppers.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"¾ cup brown rice (uncooked, or other rice or quinoa)",
						"1 ¼ cup vegetable broth (or salted water for cooking the rice)",
						"1 tbsp oil (for frying)", "1 onion (chopped)", "2 garlic cloves (chopped)",
						"1 tsp paprika powder", "1 tsp chili powder (or more if you like it spicier)",
						"1 tsp cumin",
						"1x14 oz can black beans (8.5oz (240g) rinsed &amp; drained, or sub other beans)",
						"1 cup corn (rinsed &amp; drained)", "14 oz chopped tomatoes (or passata)",
						"salt and pepper (to taste)", "4 bell peppers (halved &amp; seeded)",
						"1 cup vegan grated cheese", "1 recipe Vegan Aioli",
						"green hot peppers (or jalapeños, sliced)",
						"scallions (sliced)",
						"parsley or cilantro",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Cook the rice in vegetable broth (or salted water) according to package directions. Then allow to cool, gently fluffing occasionally (can be cooked ahead).",
						"Heat the oil in a large skillet and and sauté the onion for 1-2 minutes until translucent. Then add the garlic and sauté for 30 seconds. Then add the paprika, chili and cumin and sauté briefly.",
						"Turn off the heat and stir in the rice, beans, corn and tomatoes. Season with salt and pepper to taste and fold in half of the vegan grated cheese.",
						"Preheat the oven to 428 °F (220 °C) and cover the bottom of an ovenproof pan or casserole dish with a little water (about ¼ cup).",
						"Place the halved peppers in the baking dish and fill with the Mexican rice. Top with the remaining cheese and cover the pan or dish with a lid (or aluminum foil).",
						"Bake the stuffed peppers for 30 minutes. Then remove the lid, reduce the heat to 392 °F (200 °C) and bake for another 10-15 minutes or until the peppers are tender and lightly golden brown.",
						"Meanwhile, prepare the aioli according to this recipe.",
						"Garnish the stuffed peppers with hot peppers or jalapeños, spring onions and parsley or cilantro as desired and serve with vegan aioli.",
						"Enjoy!",
					},
				},
				Name: "Vegan Mexican Stuffed Peppers",
				NutritionSchema: models.NutritionSchema{
					Calories:      "409.7 kcal",
					Carbohydrates: "69.8 g",
					Fat:           "10.3 g",
					Fiber:         "11.4 g",
					Protein:       "12.4 g",
					SaturatedFat:  "2.4 g",
					Servings:      "1",
					Sodium:        "703.9 mg",
					Sugar:         "8.7 g",
				},
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://biancazapatka.com/en/vegan-stuffed-peppers/#recipe",
			},
		},
		{
			name: "bigoven.com",
			in:   "https://www.bigoven.com/recipe/vegetable-tempura-japanese/19344",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Dish"},
				CookTime:      "PT30M",
				Cuisine:       models.Cuisine{Value: "Japanese"},
				DatePublished: "2004/01/01",
				Description:   models.Description{Value: "not set"},
				Keywords: models.Keywords{
					Values: "nrm, side dish, snacks, vegetables, fry, fall, spring, summer, winter, meatless, vegetarian, " +
						"japanese,  qeethnic, contains white meat, nut free, contains gluten, red meat free, shellfish " +
						"free, contains eggs, dairy free",
				},
				Image: models.Image{
					Value: "https://bigoven-res.cloudinary.com/image/upload/h_320,w_320,c_fill/vegetable-tempura-japanese-e79b5b.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 cup All purpose flour sifted",
						"1 teaspoon Salt",
						"1/8 teaspoon Baking soda",
						"1 large Egg yolk",
						"2 cup Ice water",
						"Vegetable oil for frying",
						"2 medium Zucchini sliced thin",
						"1 medium Green pepper cut into strips",
						"1 large Onion sliced",
						"1/2 pound Button mushrooms",
						"1 cup Broccoli",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Separate the onion into rings. Steam the broccoli 5 minutes (or microwave a few minutes).  In an electric blender combine the flour, salt, baking soda, egg yolk, and water. Blend to mix. Let stand 15 minutes.  Heat 3 - 4 inches of oil in a deep heavy kettle, deep-fat fryer, or electric wok until it registers 375F / 190C on a deep-fat thermometer. Test batter consistency by dipping one piece of vegetable and letting excess drip off. There should be a light coating left on.  Dip and fry, a few at a time, in the hot oil until golden. Drain on paper towels and keep warm in the oven heated to 250F / 130C / Gas Mark  until all are cooked.",
					},
				},
				Name: "Vegetable Tempura - Japanese",
				NutritionSchema: models.NutritionSchema{
					Calories:      "300 calories",
					Carbohydrates: "60.0488662958771 g",
					Cholesterol:   "52.445 mg",
					Fat:           "2.32476678639706 g",
					Fiber:         "5.44774018914959 g",
					Protein:       "11.4823529627363 g",
					SaturatedFat:  "0.619135359073004 g",
					Servings:      "1",
					Sodium:        "50.8475852783613 mg",
					Sugar:         "54.6011261067275 g",
					TransFat:      "0.476087607733719 g",
				},
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.bigoven.com/recipe/vegetable-tempura-japanese/19344",
			},
		},
		{
			name: "blueapron.com",
			in:   "https://www.blueapron.com/recipes/sweet-spicy-pork-belly-fried-rice-with-kimchi-fried-eggs-71501763-72c6-4cc2-86d5-de2bc14d6b7b",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				Cuisine:       models.Cuisine{Value: "Asian"},
				DatePublished: "2023-07-26",
				Description:   models.Description{Value: "Pork belly is an incredibly rich, delicious cut of meat (used to make bacon), whose high fat content allows the meat to turn perfectly tender and flavorful as it cooks and the fat renders. Here, we're coating it with a soy and gochujang glaze, then serving it over fried rice laden with bites of vibrant kimchi and tender bok choy."},
				Image: models.Image{
					Value: "https://media.blueapron.com/recipes/42230/c_main_dish_images/1698080306-30401-0011-1025/0529_FP11_Sweet-Spicy-Glazed-Crispy-Pork-Belly_462_Web_high_feature.jpg",
				},
				Ingredients: models.Ingredients{Values: []string{
					"8 oz No Added Hormones Cooked Pork Belly",
					"2 Pasture-Raised Eggs",
					"½ cup Sushi Rice",
					"¼ cup Cornstarch",
					"2 Scallions",
					"1 Tbsp Soy Sauce",
					"2 tsps Gochujang",
					"3 Tbsps East Asian-Style Sautéed Aromatics",
					"10 oz Baby Bok Choy",
					"1 Tbsp Sesame Oil",
					"⅓ cup Kimchi",
					"1 tsp Black & White Sesame Seeds",
					"2 Tbsps Soy Glaze",
					"1 Tbsp Rice Vinegar",
				}},
				Instructions: models.Instructions{Values: []string{
					"In a small pot, combine the rice, a big pinch of salt, and 3/4 cup of water. Heat to boiling on high. Once boiling, reduce the heat to low. Cover and cook, without stirring, 15 to 17 minutes, or until the water has been absorbed and the rice is tender. Turn off the heat and fluff with a fork.",
					"Meanwhile, wash and dry the fresh produce. Cut off and discard the root ends of the bok choy; thinly slice crosswise. Thinly slice the scallions, separating the white bottoms and hollow green tops. Roughly chop the kimchi. Pat the pork belly dry with paper towels. Place on a cutting board with the fat cap facing up; cut crosswise into 1/2-inch-thick pieces. Place in a bowl. Add the soy sauce; stir to coat. Set aside to marinate, stirring occasionally, at least 10 minutes. To make the glaze, in a separate bowl, combine the vinegar, soy glaze, 1 tablespoon of water, and as much of the gochujang as you'd like, depending on how spicy you'd like the dish to be.",
					"In a medium pan (nonstick, if you have one), heat the sautéed aromatics on medium-high until hot. Add the sliced bok choy, sliced white bottoms of the scallions, and chopped kimchi; season with salt and pepper. Cook, stirring occasionally, 3 to 4 minutes, or until softened and combined. Transfer to a large bowl. Taste, then season with salt and pepper if desired. Wipe out the pan.",
					"In the same pan, heat the sesame oil on medium-high until hot. Add the cooked rice in an even layer. Cook, without stirring, 4 to 5 minutes, or until slightly crispy. Transfer to the bowl of cooked vegetables; stir to combine. Taste, then season with salt and pepper if desired. Cover with foil to keep warm. Wipe out the pan.",
					"In the same pan, heat a drizzle of olive oil on medium-high until hot. Crack the eggs into the pan, keeping them separate; season with salt and pepper. Cook 4 to 5 minutes, or until the whites are set and the yolks are cooked to your desired degree of doneness. Transfer to a plate.",
					"To the bowl of marinated pork belly, add half the cornstarch (you will have extra). Stir to evenly coat. In the same pan, heat a drizzle of olive oil on medium-high until hot. Add the coated pork belly. Cook 2 to 3 minutes per side, or until browned and heated through.* Add the glaze (carefully, as the liquid may splatter). Cook, stirring frequently, 30 seconds to 1 minute, or until the pork is coated. Turn off the heat. Serve the fried rice topped with the glazed pork belly and fried eggs. Garnish with the sesame seeds and sliced green tops of the scallions. Enjoy!",
				}},
				Keywords:        models.Keywords{Values: "Asian"},
				Name:            "Sweet & Spicy Pork Belly Fried Rice",
				NutritionSchema: models.NutritionSchema{Calories: "1140 Cals"},
				URL:             "https://www.blueapron.com/recipes/sweet-spicy-pork-belly-fried-rice-with-kimchi-fried-eggs-71501763-72c6-4cc2-86d5-de2bc14d6b7b",
				Yield:           models.Yield{Value: 2},
			},
		},
		{
			name: "bluejeanchef.com",
			in:   "https://bluejeanchef.com/recipes/apple-dutch-baby-pancake/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Entrées"},
				CookTime:      "PT20M",
				Cuisine:       models.Cuisine{Value: "American"},
				DateCreated:   "",
				DateModified:  "",
				DatePublished: "2023-08-24T11:42:46+00:00",
				Description: models.Description{
					Value: "This recipe for Apple Dutch Baby Pancake is easy to make in a cast iron pan, but be ready to eat it as soon as it comes out of the oven. It waits for no one and will fall quickly. That's ok - it tastes just as delicious.",
				},
				Keywords: models.Keywords{Values: "Breakfast/Brunch"},
				Image:    models.Image{Value: "https://bluejeanchef.com/uploads/2023/07/Apple-Dutch-Baby-1280-9609.jpg"},
				Ingredients: models.Ingredients{
					Values: []string{
						"4 eggs", "⅔ cup half-and-half", "1 teaspoon pure vanilla extract",
						"4 tablespoons melted butter (divided)", "⅔ cup all-purpose flour",
						"¼ teaspoon salt", "1 large honey crisp apple (peeled and thinly sliced)",
						"¼ cup brown sugar", "½ teaspoon ground cinnamon", "Pinch of ground nutmeg",
						"Maple syrup and powdered sugar (for serving)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Place a 10-inch cast iron skillet in the oven and pre-heat the oven to 425°F.",
						"Whisk the eggs, half-and-half, vanilla extract, and 2 tablespoons of the butter together until combined. Add the flour and salt and whisk thoroughly until the mixture is smooth with no lumps. Let the batter rest for at least 10 minutes.",
						"Combine the sliced apples, brown sugar, cinnamon, and a pinch of nutmeg in a bowl and toss to coat the apples. Set the apples aside.",
						"Carefully remove the hot cast iron skillet from the oven and place on the stovetop. (Remember this skillet is hot and you will need to use oven mitts.) Pour the remaining 2 tablespoons of butter into the pre-heated cast iron skillet and carefully scatter the apples over the butter. Pour the batter into the pan evenly over the apples. Use oven mitts to return the pan to the oven.",
						"Bake for 20 to 25 minutes until the pancake is brown and has puffed up. Make sure you don’t open the oven for the first 15 minutes of the cooking process.",
						"Remove the pan from the oven and dust a little powdered sugar on top. The Dutch baby will deflate quickly so have your phone ready so you can catch your Instagram-worthy picture right away! Cut the Dutch baby into wedges and serve immediately with maple syrup.",
					},
				},
				Name: "Apple Dutch Baby Pancake",
				NutritionSchema: models.NutritionSchema{
					Calories:       "248 kcal",
					Carbohydrates:  "25 g",
					Cholesterol:    "139 mg",
					Fat:            "14 g",
					Fiber:          "1 g",
					Protein:        "6 g",
					SaturatedFat:   "8 g",
					Servings:       "1",
					Sodium:         "218 mg",
					Sugar:          "13 g",
					TransFat:       "0.3 g",
					UnsaturatedFat: "5 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://bluejeanchef.com/recipes/apple-dutch-baby-pancake/",
			},
		},
		{
			name: "briceletbaklava.ch",
			in:   "https://briceletbaklava.ch/2023/10/taille-au-sel-de-granges-marnand.html",
			want: models.RecipeSchema{
				AtContext: "https://schema.org",
				AtType:    models.SchemaType{Value: "Recipe"},
				Description: models.Description{
					Value: "Il y a quelques temps je vous avais donné la recette du taillé au sel, celle de ma voisine Claudine, une spécialité vaudoise. Cette recette, je l’avais obtenue, un soir après un bon repas entre amis, chez Claudine. Elle m’avait gribouillé la liste des ingrédients sur un morceau de papier avec quelques explications orales. Quelques semaines plus tard je m’étais lancé et avais alors refait cette recette chez moi mais quelque peu adaptée à ma façon de faire. Après plusieurs essais et content du résultat, je l’avais alors publiée (recette ICI) et il a fait des émules ce taillé, presque le tour du monde, même jusqu’aux oreilles des producteurs de la TV suisse (Couleurs Locales). Ils m’ont contacté car ils désiraient en faire un petit reportage que voici: \n\ninfo : pour visualiser uniquement la petite séquence, allez à la minute 5:15\n\n\n\nAlors, j’ai aussitôt convoqué Claudine et devant les caméras, elle nous a réalisé son taillé, SA recette avec tous les détails qui me manquaient. Mais claudine a bien précisé que ce taillé est une spécialité du village où nous habitons, cette recette qui nous a été transmise de familles en familles et que nous avons envie de faire perdurer. \n\nSa confection est des plus simple et rapide à réaliser. Ici pas de robot, pas de levain et le pétrissage se fait à l’ancienne, c’est-à-dire à la main. Voilà ce que ça donne :",
				},
				Keywords: models.Keywords{Values: "Apéritifs,Boulangerie,Petit-déjeuner,Recettes suisses salées"},
				Image: models.Image{
					Value: "https://image.over-blog.com/xNKjelRLoiHpZIVbWXiqrgKklls=/filters:no_upscale()/image%2F3215825%2F20231014%2Fob_e40698_taille-au-sel-claudine-7874.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"Recette pour 1 taillé d’un diamètre de 26 cm, env. 800 g, 10 à 15 personnes en apéro",
						"Ingrédients",
						"500 g de farine fleur ou 550 ou mélange pour tresse\u00a0(Tabelle des sortes de farine ICI)",
						"16 g de sel",
						"34 cl de lait entier",
						"70 g de beurre", "1 cc de sucre",
						"20 g de levure fraiche de boulanger",
						"50 g de beurre pour le dessus",
						"2 à 3 cuillères à soupe de crème à 35%",
						"Sel pour saupoudrer la surface",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Dans un grand bol, déposer la farine préalablement tamisée.",
						"Ajouter les 16 grammes de sel. Oui je sais, c’est beaucoup, presque trop, mais c’est la recette,",
						"et bien l’incorporer à la farine.",
						"Peser et débiter le beurre en cubes et le mettre dans une petite casserole,",
						"puis verser le lait par dessus et ajouter encore le sucre.",
						"A tout petit feu, tout en brassant, faire fondre le beurre puis ajouter la levure et la faire fondre également.",
						"Mais attention, à ce stade la température du lait ne devra en aucun cas dépasser les 30°C, au risque de tuer les ferments de la levure. Et ce sera le ratage assuré.",
						"Ajouter ensuite ce lait sur la farine",
						"et à la main, commencer le pétrissage pour en faire une pâte assez mole. Corriger en ajoutant un peu de farine si nécessaire.",
						"Après 5 à 10 minutes de pétrissage, la pâte va se former, devenir lisse et se détacher de la paroi du bol.",
						"En faire une boule et la laisser reposer dans le bol, à couvert, pendant 10 à 15 minutes, le temps de beurrer et fariner un moule de cuisson d’un diamètre d’environ 26 cm et d’une hauteur de 6 cm. Ici j’ai utilisé un moule à charnière.",
						"Maintenant reprendre la pâte qui se sera un peu détendue puis la déposer dans le fond du moule, presser afin de bien égaliser la surface.",
						"Couvrir d'un linge ou d’un plastique et laisser lever 1 h à 1½ h à température ambiante, idéalement 22 à 24 degrés.",
						"La pâte devra doubler de volume.",
						"A l’aide d’un petit couteau tranchant, cisailler le dessus du taillé en damier",
						"puis avec vos doigts, enfoncer des petites noix de beurre dans les interstices.",
						"Avec un pinceau, badigeonner légèrement la surface avec la\u00a0crème et la saupoudrer généreusement de sel fin.",
						"Enfourner de suite dans un four préchauffé à 200°C et cuire pendant 30 minutes ou jusqu’à ce que le taillé soit bien doré.",
						"Le démouler aussitôt et le faire refroidir sur une grille.",
						"Attention, ce taillé est complètement addictif, surtout en le servant à l’occasion d’un apéro. Un chasselas sera recommandé.",
						"J’ai aussi entendu que certains le dégustaient au petit déjeuner accompagné d’une tasse de cacao mais ici, je n’ai pas eu le courage d'essayer.",
					},
				},
				Name: "Taillé au sel de Granges-Marnand",
				URL:  "https://briceletbaklava.ch/2023/10/taille-au-sel-de-granges-marnand.html",
			},
		},
		{
			name: "bodybuilding.com",
			in:   "https://www.bodybuilding.com/recipes/beef-teriyaki-rice-and-stir-fry",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT2M",
				DateModified:  "2023-11-07T15:51:57Z",
				DatePublished: "2023-10-02T00:00:00Z",
				Description: models.Description{
					Value: "Our highest-quality Hawaiian teriyaki beef is doused in a special house teriyaki sauce and sliced into scrumptious bite-sized pieces. The beef is then seared on a hot flat griddle to lock in the flavors and create a caramelized coating.",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 serving Hawaiian teriyaki beef", "½ cup brown rice",
						"2 tbsp light soy sauce", "1 tbsp scallions, sliced", "¼ cup green beans",
						"¼ cup broccoli", "¼ cup, chopped red bell pepper, sliced",
						"¼ cup mushrooms, sliced",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Prepare rice to your specifications.",
						"Chop broccoli, red pepper, and green onions.",
						"Add soy sauce and green onions to rice.",
						"Cook broccoli to your specifications.",
						"Remove plastic from thawed Hawaiian Teriyaki Beef package.",
						"In a skillet, add one serving of protein, a portion of the teriyaki marinade, and heat for 2 minutes, stirring every 30 seconds.",
						"Plate up and serve!",
					},
				},
				Keywords: models.Keywords{Values: "Dinner,Lunch"},
				Name:     "Beef Teriyaki Rice and Stir Fry",
				NutritionSchema: models.NutritionSchema{
					Calories:      "352 kcal",
					Carbohydrates: "38 g",
					Fat:           "5 g",
					Protein:       "40 g",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 1},
				URL:      "https://www.bodybuilding.com/recipes/beef-teriyaki-rice-and-stir-fry",
			},
		},
		{
			name: "bonappetit.com",
			in:   "https://www.bonappetit.com/recipe/crispy-chicken-with-zaatar-olive-rice",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				DateModified:  "2022-03-15T05:00:00.000-04:00",
				DatePublished: "2022-03-15T05:00:00.000-04:00",
				Description: models.Description{
					Value: "Give ground chicken the respect it deserves.",
				},
				Keywords: models.Keywords{
					Values: "main,dinner,quick,one-pot meals,easy,weeknight meals,healthyish,gluten-free,nut-free,ground " +
						"chicken,ground turkey,feta,olive,rice,za'atar,spinach,kale,sauté,swiss chard,castelvetrano olive,web",
				},
				Image: models.Image{
					Value: "https://assets.bonappetit.com/photos/6228bc8071b26c82f857f620/16:9/w_6208,h_3492,c_limit/Crispy-Chicken-With-Za%E2%80%99atar-Olive-Rice.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 lb. ground chicken or turkey",
						"½ tsp. smoked paprika",
						"1 tsp. Diamond Crystal or ½ tsp. Morton kosher salt, plus more",
						"Freshly ground black pepper",
						"3 Tbsp. extra-virgin olive oil",
						"1 cup Castelvetrano olives, smashed, pits removed",
						"3 cups cooked rice",
						"1 Tbsp. za’atar, plus more for serving",
						"2 cups coarsely chopped greens (such as spinach, kale, or chard)",
						"Zest and juice of 1 small lemon",
						"2 oz. feta, thinly sliced into planks",
						"Coarsely chopped dill (for serving)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Place chicken in a medium bowl. Sprinkle paprika and 1 tsp. Diamond Crystal or ½ tsp. Morton kosher " +
							"salt over chicken; season with pepper. Gently mix with your hands to combine.",
						"Heat oil in a large nonstick skillet over medium-high. Arrange chicken in pan in a thin, even layer " +
							"and cook, undisturbed, until golden brown and crisp underneath, about 5 minutes. " +
							"Continue to cook, stirring and breaking up into bite-size pieces with a wooden " +
							"spoon, until cooked through, about 1 minute. Using a slotted spoon, transfer chicken to a " +
							"plate, leaving oil and fat behind.",
						"Add olives to same pan and cook, undisturbed, until heated through and blistered, 1–2 minutes. Add " +
							"rice and 1 Tbsp. za’atar and cook, stirring often, until slightly crisp, about 3 " +
							"minutes. Add greens and lemon juice and cook, stirring occasionally, until greens " +
							"are wilted, about 2 minutes. Remove pan from heat; stir in lemon zest, feta, and chicken. " +
							"Taste and season with more salt and pepper if needed.",
						"Transfer chicken and rice to a large shallow bowl; sprinkle with more za’atar and top with dill.",
					},
				},
				Name:  "Crispy Chicken With Za’atar-Olive Rice",
				Yield: models.Yield{Value: 4},
				URL:   "https://www.bonappetit.com/recipe/crispy-chicken-with-zaatar-olive-rice"},
		},
		{
			name: "bongeats.com",
			in:   "https://www.bongeats.com/recipe/chicken-lollipop",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "appetizer"},
				CookTime:      "PT30M",
				Cuisine:       models.Cuisine{Value: "Indian"},
				DatePublished: "Feb 23, 2019",
				Description: models.Description{
					Value: "Learn how to easily make lollipops from chicken wings, then turn them into the hot-sour-crunchy appetiser, drums of heaven",
				},
				Keywords: models.Keywords{
					Values: "chinese chicken recipe, indo-chinese recipe, drums of heaven, chicken drumsticks, chicken lollipop",
				},
				Image: models.Image{
					Value: "https://assets-global.website-files.com/60d34b8627f6e735cf28df18/62a95ea9075f335474fa40bc_Chicken%20Lollipop%20Hero%201.1.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"500 g chicken wings (6 whole wings) ", "\u200d6 g (1 tsp) soy sauce",
						"6 g (1 tsp) green chilli sauce", "15 g (1½ \u00a0tbsp) red chilli sauce",
						"4 g salt", "¼ tsp MSG", "½ tsp ground black pepper", "2 g ginger",
						"2 g garlic", "1 tbsp egg whites", "4 tbsp plain flour (maida)",
						"2 tbsp cornstarch", "vegetable oil for deep-frying", "10 g vegetable oil",
						"50 g onions", "5 g ginger",
						"12 g garlic",
						"3 tbsp chicken stock",
						"30 g (3 tbsp) red chilli sauce",
						"8 g (1 tsp) ketchup",
						"8 g (1 tsp) dark soy sauce\u200d",
						"2 pc spring onions",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Separate the drumette and the flat. Pop the joint between drumette and the flat. Then cut between the joint to separate.",
						"With the flat, hold the tip in your right hand, and flat in the left hand (reverse if you're left-handed). Use force to push out the two small bones. Discard the smaller bone.",
						"You can now discard the wing tip. We needed it for gripping earlier. ",
						"Gather the meat to one side, turning the meat inside-out. This is the first lollipop—the smaller (but tastier) one.",
						"Now, the drumette. We need to expose the joint at the end where we separated the drumette from the flat. ",
						"Cut the connective tissue around the base of the bone, then angle the blade and scrape towards you to form a clearing where you can hold the bone.",
						"Hold the exposed joint with a cloth and scrape away to detach the meat from the bone and gather it at the other end. This will take some force.",
						"Now, hold the loose meat with the cloth and pull to turn the meat inside out to form the second (bigger) lollipop.",
						"Repeat these steps for all the wings. You will end up with double the number of lollipops as the wings you start with.",
						"Mix the soy sauce, green chilli sauce, red chilli sauce, salt, MSG, pepper, ginger, garlic, and egg whites in a mixing bowl.",
						"Add the chicken lollipops, massage them well, and leave them to marinate for as long as it takes to prep the ingredients for the sauce.",
						"The cooking will take place quickly on high heat. So, all your ingredients should be prepped and at hand before the cooking starts.",
						"To make the sauce mix, combine red chilli sauce, ketchup, and dark soy sauce in a bowl.",
						"Heat the kadai until very hot. Add vegetable oil and let it start smoking.",
						"Add the diced onions and cook for one minute.",
						"Add the ginger and fry 30 seconds.",
						"Add garlic and fry 30 seconds.",
						"Add sugar and let it melt and caramelise—about 30 seconds.",
						"Next add the salt, ground pepper, the sauce mix and the chicken stock.",
						"Reduce for a few minutes until the sauce gets syrupy and sticky.\u00a0Turn off the heat.",
						"Heat oil in a kadai or wok until it is very hot (200ºC).",
						"Just before frying, add the plain flour and cornstarch (cornflour) to the marinated chicken and mix to form a dry batter. Adding too much liquid in the batter or the lollipops won't stay crisp.",
						"Add the lollipops to the hot oil, taking care not to overcrowd the pan. Fry in two batches if required. The oil temperature will drop to about 150ºC. The rest of the frying should take place at this temperature to make sure the chicken cooks all the way without getting too brown outside.",
						"Drain on a basket lined with paper towels.",
						"Re-heat the sauce if too cold, and turn off the heat. Now add the chicken lollipops to the pan and toss to coat it evenly. ",
						"Garnish with spring onions and serve hot.",
					},
				},
				Name:     "Chicken Lollipop / Drums of Heaven",
				PrepTime: "PT1H30M",
				Yield:    models.Yield{Value: 12},
				URL:      "https://www.bongeats.com/recipe/chicken-lollipop",
			},
		},
		{
			name: "bowlofdelicious.com",
			in:   "https://www.bowlofdelicious.com/mini-meatloaves/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				CookTime:      "PT25M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2021-01-19T11:41:00+00:00",
				Description: models.Description{
					Value: "Mini Meatloaves cook up in HALF the time of a whole meatloaf - an easy, fast gluten-free recipe " +
						"made with oats instead of breadcrumbs.",
				},
				Keywords: models.Keywords{Values: "Mini meatloaves"},
				Image: models.Image{
					Value: "https://www.bowlofdelicious.com/wp-content/uploads/2014/09/Mini-Meatloaves-square.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 onion (grated)",
						"1 cup quick oats",
						"1/2 cup milk",
						"2 eggs",
						"1 teaspoon kosher salt",
						"1/2 teaspoon black pepper",
						"2 teaspoons soy sauce ((gluten-free if necessary))",
						"2 lbs. ground beef (preferably 80/20)",
						"ketchup or barbecue sauce (for topping, optional)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 400 degrees F. Butter or grease a rimmed baking sheet or casserole dish.",
						"In a large bowl, add the grated onion, quick oats (1 cup), milk (1/2 cup), 2 eggs, kosher salt " +
							"(1 teaspoon), black pepper (1/2 teaspoon), and soy sauce (2 teaspoons) and mix until " +
							"well combined. Allow to sit for 3-5 minutes while the oats absorb some of the moisture.",
						"Add the ground beef (2 lbs.) to the bowl. Mix well- for best results, use hands.",
						"Divide the mixture into eight parts and form them into loaf shapes. (I like to flatten the mixture " +
							"in the bowl and use a knife to portion it out into 8 \"wedges,\" kind of like " +
							"slicing a cake, to get even amounts).",
						"At this point, you can wrap the mini meatloaves in plastic wrap or flash freeze to store in the freezer, or" +
							" refrigerate until you&#039;re ready to bake (see notes for more info on this). Otherwise, proceed to cooking.",
						"Place the mini meatloaves on the prepared baking sheet or casserole dish (as many as you want to cook - " +
							"if frozen, defrost completely before cooking). Bake at 400 degrees F for 25 minutes (or until " +
							"the internal temperature is 160 degrees).",
					},
				},
				Name: "Mini Meatloaves",
				NutritionSchema: models.NutritionSchema{
					Calories:      "356 kcal",
					Carbohydrates: "10 g",
					Cholesterol:   "121 mg",
					Fat:           "25 g",
					Fiber:         "1 g",
					Protein:       "23 g",
					SaturatedFat:  "9 g",
					Servings:      "1",
					Sodium:        "474 mg",
					Sugar:         "2 g",
					TransFat:      "1 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://www.bowlofdelicious.com/mini-meatloaves/",
			},
		},
		{
			name: "budgetbytes.com",
			in:   "https://www.budgetbytes.com/easy-vegetable-stir-fry/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "Asian"},
				DatePublished: "2022-03-15T07:48:39+00:00",
				Description: models.Description{
					Value: "Vegetable stir fry is a quick and easy option for dinner, plus it&#039;s super flexible and a " +
						"great way to use up leftovers from your fridge!",
				},
				Keywords: models.Keywords{Values: "Stir Fry Recipe, vegetable stir fry"},
				Image: models.Image{
					Value: "https://www.budgetbytes.com/wp-content/uploads/2022/03/Easy-Vegetable-Stir-Fry-close.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1/4 cup soy sauce ($0.24)",
						"1/4 cup water ($0.00)",
						"2 Tbsp brown sugar ($0.08)",
						"1 tsp toasted sesame oil ($0.10)",
						"2 cloves garlic, minced ($0.16)",
						"1 tsp grated fresh ginger ($0.10)",
						"1 Tbsp cornstarch ($0.03)",
						"3/4 lb. broccoli ($1.34)",
						"2 carrots ($0.33)",
						"8 oz. mushrooms ($1.69)",
						"8 oz. sugar snap peas ($2.99)",
						"1 small onion ($0.28)",
						"1 red bell pepper ($1.50)", "2 Tbsp cooking oil ($0.16)",
						"1 tsp sesame seeds ($0.06)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Make the stir fry sauce first. Combine the soy sauce, water, brown sugar, sesame oil, garlic, ginger, " +
							"and cornstarch in a small bowl. Set the sauce aside.",
						"Chop the vegetables into similar-sized pieces. It&#39;s up to you whether you slice, dice, or cut into " +
							"any other shape you prefer.",
						"Add the cooking oil to a very large skillet or wok. Heat over medium-high. When the pan and oil are very " +
							"hot (but not smoking), add the hardest vegetables first: carrots and broccoli. Cook and " +
							"stir for about a minute, or just until the broccoli begins to turn bright green.",
						"Next, add the mushrooms and sugar snap peas. Continue to cook and stir for a minute or two more, or just " +
							"until the mushrooms begin to soften.",
						"Finally, add the softest vegetables, bell pepper and onion. Continue to cook and stir just until the onion " +
							"begins to soften.",
						"Give the stir fry sauce another brief stir, then pour it over the vegetables. Continue to cook and stir " +
							"until the sauce begins to simmer, at which point it will thicken and turn glossy. Remove " +
							"the vegetables from the heat, or continue to cook until they are to your desired doneness.",
						"Top the stir fry with sesame seeds and serve!",
					},
				},
				Name: "Easy Vegetable Stir Fry",
				NutritionSchema: models.NutritionSchema{
					Calories:       "209 kcal",
					Carbohydrates:  "27 g",
					Cholesterol:    "",
					Fat:            "9 g",
					Fiber:          "6 g",
					Protein:        "8 g",
					SaturatedFat:   "",
					Sodium:         "869 mg",
					Sugar:          "",
					TransFat:       "",
					UnsaturatedFat: ""},
				PrepTime: "PT15M",
				Tools:    models.Tools{Values: []string(nil)},
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.budgetbytes.com/easy-vegetable-stir-fry/"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
