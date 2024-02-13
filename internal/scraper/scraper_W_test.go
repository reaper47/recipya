package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_W(t *testing.T) {
	testcases := []testcase{
		{
			name: "waitrose.com",
			in:   "https://www.waitrose.com/ecom/recipe/the-best-macaroni-cheese",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main meal"},
				CookTime:      "PT15M",
				DatePublished: "2022-10-12",
				Description: models.Description{
					Value: "Martha Collison's version is creamy and strong in cheese flavour, but not too intense. A salad on the side makes a great accompaniment.",
				},
				Keywords: models.Keywords{
					Values: "Pasta, Quick and easy, Serves 6, Cheese, cheddar, c, Main meal, Graham Norton show, Mac And Cheese",
				},
				Image: models.Image{
					Value: "https://waitrose-prod.scene7.com/is/image/waitroseprod/macaroni-cheese?uuid=ceb08341-8d96-45cc-a2f9-7427489eed93&$Waitrose-Default-Image-Preset$",
				},
				Ingredients: models.Ingredients{Values: []string{"350 g De Cecco Chifferi Rigati, or similar pasta"}},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat the oven to 200\n°C, gas mark 6.\nBring a large pan of water to the boil and add\na big pinch of salt. Add the pasta and cook\nfor 1 minute less than pack instructions (it will\ncontinue to cook in the oven later). Drain and\nset aside.",
						"While the pasta cooks, make the cheese\nsauce. Melt the butter in a large saucepan.\nOnce melted, add the flour and stir until a thick\npaste forms. Cook the mixture for 1 minute,\nstirring constantly, to cook the flour out. Add\nthe paprika, mustard powder, garlic granules\nand pepper, then stir again.",
						"Gradually add the milk, little by little, until the\nsauce is completely smooth. Once bubbling,\nstir in the evaporated milk and add ¾ of all the\ngrated cheese. Mix well until the cheese has\nmelted, then drain the pasta and add it to the\nsauce. Tip into a large ovenproof dish and top\nwith the remaining cheese.",
						"Bake for 10 minutes, then switch the oven to\ngrill and cook for 5 minutes more, or until the\ncheese on top is bubbling and golden. Allow\nto stand for a few minutes before serving.",
					},
				},
				Name: "The best macaroni cheese",
				NutritionSchema: models.NutritionSchema{
					Calories:      "724 calories",
					Carbohydrates: "60 g",
					Fat:           "39 g",
					Fiber:         "2.6 g",
					Protein:       "32 g",
					SaturatedFat:  "24 g",
					Sodium:        "1.6 g",
					Sugar:         "8.7 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://www.waitrose.com/ecom/recipe/the-best-macaroni-cheese",
			},
		},
		{
			name: "watchwhatueat.com",
			in:   "https://www.watchwhatueat.com/healthy-fried-brown-rice/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main or Side"},
				CookTime:      "PT30M",
				Cuisine:       models.Cuisine{Value: "Asian Inspired"},
				DatePublished: "2019-04-26T13:31:54+00:00",
				Description: models.Description{
					Value: "Learn how to make healthy fried brown rice with fresh vegetables. Simple, delicious and wholesome " +
						"vegetable fried rice recipe that is better than takeout.",
				},
				Keywords: models.Keywords{Values: "fried rice, healthy fried rice"},
				Image: models.Image{
					Value: "https://www.watchwhatueat.com/wp-content/uploads/2019/04/Healthy-Fried-Brown-Rice-6.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 cup Jasmine brown rice",
						"3 1/2 cup water",
						"1/2 tbsp sesame oil",
						"3-4 medium garlic cloves (finely chopped)",
						"1/2\"  ginger (peeled and finely chopped)",
						"1/2 bell pepper diced (red and green each)",
						"1 large carrot (peeled and diced)",
						"1 cup green peas (fresh or frozen)",
						"3-4 springs of green onion ((scallions))",
						"1 1/2 tbsp soy sauce",
						"1 tbsp rice vinegar",
						"salt to taste",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a medium pot add rice and water and bring mixture to boil. Then simmer on low heat with cover for 25 " +
							"mins. Let it cool down completely. Using a fork fluff the rice to separate the grains before using.",
						"In a large skillet or wok heat oil on medium to high heat.",
						"Add chopped garlic, ginger and the white portion of scallions (reserve green for garnishing). Cook it for " +
							"1-2 mins.",
						"Then add diced peppers, carrot and peas. Stir fry them for few minutes (see notes).",
						"Now add soy sauce, rice vinegar and mix well.",
						"Finally, add cold rice and mix well with vegetables. Season with salt if necessary.",
						"Garnish with green onions and serve warm.",
					},
				},
				Name: "Healthy Fried Brown Rice With Vegetables",
				NutritionSchema: models.NutritionSchema{
					Calories: "336 kcal",
					Servings: "1",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 5},
				URL:      "https://www.watchwhatueat.com/healthy-fried-brown-rice/",
			},
		},
		{
			name: "wearenotmartha.com",
			in:   "https://wearenotmartha.com/western-omelet/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Breakfast"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-11-07T14:57:29+00:00",
				Description: models.Description{
					Value: "Start your day with a hearty and delicious Western Omelet, a breakfast that&#39;s easy enough to make on both weekends and weekday mornings. Packed with savory ham, green bell peppers, onions, and cheddar cheese, this omelet is sure to become a breakfast favorite in your house!",
				},
				Keywords: models.Keywords{Values: "Egg Recipes, Ham Recipe, Omelet Recipes"},
				Image: models.Image{
					Value: "https://wearenotmartha.com/wp-content/uploads/western-omelet-featured-2.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 large eggs", "1/4 tsp salt, (divided)", "1/4 tsp pepper, (divided)",
						"1 Tbsp unsalted butter, (divided)", "1/3 cup diced ham",
						"1/3 cup diced green pepper", "1/4 cup diced onion",
						"1/3 cup shredded cheddar cheese",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a medium bowl, whisk together eggs, 1/8 tsp salt, and 1/8 tsp pepper. Set aside.",
						"Add 1/2 Tbsp butter to a nonstick 7- or 8-inch nonstick skillet over medium heat.",
						"Add ham, green pepper, onion, 1/8 tsp salt, and 1/8 tsp pepper to the skillet and cook until veggies are softened and ham is just starting to brown, about 3-5 minutes. Remove to a bowl and set aside.",
						"With skillet still over medium heat, add remaining 1/2 Tbsp butter and let coat bottom of pan.",
						"Pour the beaten eggs into the skillet and let them cook undisturbed for a couple minutes until they start to set around the edges.",
						"Using a spatula, gently push the cooked edges toward the center, tilting the pan to let the uncooked eggs flow to the edges to cook. Continue until there’s no liquid egg left, but the top is still slightly runny. You can gently flip the omelet at this point, but you don&#39;t have to.",
						"Add cooked veggies and ham to one side of the omelet and sprinkle with shredded cheese.",
						"Carefully fold the omelet in half, covering the filling and let cook for another minute to melt cheese.",
					},
				},
				Name:     "Western Omelet",
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 1},
				URL:      "https://wearenotmartha.com/western-omelet/",
			},
		},
		{
			name: "weightwatchers.com",
			in:   "https://www.weightwatchers.com/us/recipe/pepperoni-flatbread-pizza/646d12d61a843705da13cb7f",
			want: models.RecipeSchema{
				AtContext: "https://schema.org",
				AtType:    models.SchemaType{Value: "Recipe"},
				CookTime:  "PT12M",
				Image: models.Image{
					Value: "https://cmx.weightwatchers.com/assets-proxy/weight-watchers/image/upload/t_WINE_EXTRALARGE/ewh3p3r7s8g3iwix6u9s.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 spray(s) Cooking spray",
						"1 piece(s) Damascus Bakeries All natural plain roll-ups",
						"2 Tbsp Store-bought pizza sauce", "3 Tbsp Part skim mozzarella cheese",
						"6 slice(s) Turkey pepperoni",
						"3 item(s) Sweet mini baby bell peppers seeded and thinly sliced",
						"1.5 tsp Grated Parmesan cheese",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 425°F. Line a baking sheet with parchment paper. Coat the paper with cooking spray.",
						"Place the flatbread on the lined baking sheet. Spread the sauce all over flatbread. Top with the mozzarella, pepperoni, and then peppers. Lightly coat top with cooking spray. Sprinkle with the Parmesan.",
						"Bake 10 to 15 minutes, until edges are golden brown and cheese is melted. Season with crushed red pepper, Italian seasoning, and/or garlic powder, if desired.",
					},
				},
				Name: "Pepperoni flatbread pizza",
				NutritionSchema: models.NutritionSchema{
					Calories: "484 kcal",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 1},
				URL:      "https://www.weightwatchers.com/us/recipe/pepperoni-flatbread-pizza/646d12d61a843705da13cb7f",
			},
		},
		{
			name: "wellplated.com",
			in:   "https://www.wellplated.com/energy-balls/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Snack"},
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2017-09-11T05:05:00+00:00",
				Description: models.Description{
					Value: "The only Energy Ball recipe you'll ever need, plus six no-bake energy ball flavors! Start with this easy base recipe, then add any of your favorite mix-ins.",
				},
				Keywords: models.Keywords{Values: "Easy Snack Recipe, No Bake Oatmeal Energy Balls"},
				Image: models.Image{
					Value: "https://www.wellplated.com/wp-content/uploads/2017/09/How-to-Make-Energy-Balls-no-text.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 1/4 cups old fashioned rolled oats (you can also swap quick oats or a blend of half quick, half old fashioned)",
						`2 tablespoons "power mix-ins"`,
						"1/2 cup nut butter of choice (peanut butter is my go-to)",
						"1/3 cup sticky liquid sweetener of choice (honey or maple syrup)",
						"1 teaspoon pure vanilla extract", "1/4 teaspoon kosher salt",
						"1/2 cup mix-ins (see below for flavor options)",
						"Any nut butter (honey, 1/2 cup chocolate chips)",
						"Peanut butter (honey, 3 tablespoons chocolate chips, 3 tablespoons chopped peanuts, 2 tablespoons raisins)",
						"Almond butter ( or cashew butter, honey, 1/4 cup dried cranberries, 1/4 cup white chocolate chips)",
						"Replace 1/2 cup of the oatmeal with 1/2 cup unsweetened coconut flakes",
						"Any nut butter (any sweetener, 1/2 cup mini chocolate chips, ADD 2 tablespoons cocoa powder)",
						"Almond butter ( or cashew butter, maple syrup, 1/2 cup raisins, ADD 1/4 teaspoon cinnamon)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Place all of the ingredients in a large mixing bowl: oats, power mix-ins, nut butter, sweetener, vanilla extract, salt, mix-ins, and any other spices you'd like to add. Stir to combine. If the mixture seems too wet, add a bit more oats. If it's too dry, add a bit more nut butter. It should resemble a somewhat sticky dough that holds together when lightly squeezed. Place the bowl in the refrigerator for 30 minutes to set (this will make the balls easier to roll later on).",
						"Remove the bowl from the refrigerator and portion the dough into balls of desired size. (I use a cookie scoop to make mine approximately 1 inch in diameter). Enjoy!",
					},
				},
				Name: "Energy Balls",
				NutritionSchema: models.NutritionSchema{
					Calories:      "131 kcal",
					Carbohydrates: "18 g",
					Fat:           "5 g",
					Fiber:         "3 g",
					Protein:       "4 g",
					SaturatedFat:  "1 g",
					Servings:      "1",
					Sugar:         "6 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 18},
				URL:      "https://www.wellplated.com/energy-balls/",
			},
		},
		{
			name: "whatsgabycooking.com",
			in:   "https://whatsgabycooking.com/pea-prosciutto-spring-pizza/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Breakfast"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "Italian"},
				DatePublished: "2022-03-31T07:00:00+00:00",
				Description: models.Description{
					Value: "When in doubt of what to do with a bunch of fab farmers market produce, put it all on a pizza and " +
						"slap an egg on it, duh!",
				},
				Image: models.Image{
					Value: "https://whatsgabycooking.com/wp-content/uploads/2015/05/ALDI-Spring-Pea-Pizza-2-copy-2.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 recipe fresh pizza dough",
						"1/3 cup Basil Vinaigrette",
						"1 cup mozzarella (shredded or fresh mozzarella sliced )",
						"1/2 cup fresh peas (blanched)",
						"1/2 cup sugar snap peas (sliced thin )",
						"1 bunch asparagus (tips only, blanched)",
						"2 eggs",
						"4 ounces prosciutto",
						"Kosher salt and freshly cracked black pepper to taste",
						"Fresh Basil (torn into pieces)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Pre-heat oven to 475 degrees F.",
						"Shape the dough into 2 medium-ish pizzas while on a clean floured surface. Let the dough sit for 5 " +
							"minutes and then re-form to make sure it's as big as you'd like. Place the pizza dough on a " +
							"lightly floured rimless baking sheet, or pizza peel.",
						"Spread the basil vinaigrette over the top of each pizza. Top with the mozzarella, scatter the peas and " +
							"asparagus on top of the cheese. Transfer to an oven and cook for about 5 minutes. Remove the " +
							"pizza and add the egg on top of each pizza and transfer them back into the oven to continue " +
							"to cook until the egg white is set and the yolk still runny .",
						"Remove from the oven, Add the prosciutto on top and garnish with basil.",
					},
				},
				Keywords: models.Keywords{Values: "homemade pizza, spring pizza"},
				Name:     "Pea Prosciutto Spring Pizza",
				NutritionSchema: models.NutritionSchema{
					Calories:       "610 kcal",
					Carbohydrates:  "56 g",
					Cholesterol:    "123 mg",
					Fat:            "33 g",
					Fiber:          "5 g",
					Protein:        "24 g",
					SaturatedFat:   "11 g",
					Servings:       "1",
					Sodium:         "1124 mg",
					Sugar:          "11 g",
					TransFat:       "0.04 g",
					UnsaturatedFat: "19 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://whatsgabycooking.com/pea-prosciutto-spring-pizza/",
			},
		},
		{
			name: "wholefoodsmarket.co.uk",
			in:   "https://www.wholefoodsmarket.co.uk/recipes/pollo-al-ajillo",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				DateModified:  "2023-06-15T17:16:37+0100",
				DatePublished: "2023-06-13T15:12:14+0100",
				Description: models.Description{
					Value: "Seared chicken on the bone, slow cooked onions &amp; whole cloves of garlic finished with white wine, rosemary &amp; tender, roasted new potatoes.",
				},
				Image: models.Image{
					Value: "http://static1.squarespace.com/static/630492b401e7102ee99cc184/6332c38456a18f062e96eabd/648874cd317f411dacbbb426/1686845797937/CaptureGarlicChicken.PNG?format=1500w",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"100ml Olive oil", "1kg chicken drumsticks & thighs, skin on, bone in",
						"500g Spanish onion, peeled & finely sliced",
						"10 sprigs of thyme, leaves picked", "1 large head of garlic, cloves peeled",
						"1 sprig of rosemary, leaves picked & chopped", "2 bay leaves",
						"200ml white wine", "1kg New potatoes", "Salt & freshly ground black pepper",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Pre-heat the oven to 180°c. Place the new potatoes in a pan with cold water and season generously with salt.",
						"Place on a high heat and bring to the boil, once boiling reduce heat to medium and cook for approximately 15-20 minutes until tender, once cooked drain and set aside.",
						"Place a large, ovenproof, heavy bottomed pan or casserole dish on a high heat and season the chicken generously with salt & pepper.",
						"Once pan is hot, add the olive oil, lay the chicken pieces in the pan and cook on all sides until golden brown, turning every couple of minutes, this should take 10-15 minutes. Remove the chicken from the pan and set aside on a plate.",
						"Add the onions, garlic cloves, rosemary, bay leaves & thyme, turn heat to medium and cook for a further 10 minutes until the onions are starting to brown, stir every minute or so. Add the white wine to deglaze the pan and cook for a further 2 minutes to reduce. Add the potatoes and chicken and gently stir to mix the ingredients together.",
						"Place in the middle shelf of the pre-heated oven for approximately 10-15 minutes to allow the potatoes to lightly brown. Check the chicken is fully cooked with no pink meat and it is steaming hot in the centre.",
						"Allow to rest for 5 minutes, serve, and enjoy.", "See our Terms of Service.",
					},
				},
				Name: "Pollo al Ajillo",
				URL:  "https://www.wholefoodsmarket.co.uk/recipes/pollo-al-ajillo",
			},
		},
		{
			name: "wikibooks.org",
			in:   "https://en.wikibooks.org/wiki/Cookbook:Creamed_Spinach",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "Creamed Spinach",
				Category:  models.Category{Value: "Sauce recipes"},
				Description: models.Description{
					Value: "Creamed spinach makes a nutritious sauce that goes well with fish and meat dishes. In Swedish cuisine it has traditionally been used with boiled potatoes and fish or with chipolata, the Swedish \"prince-sausage\". Creamed spinach can be done in two different ways: either with whole spinach or with chopped spinach.",
				},
				Image: models.Image{
					Value: "https://upload.wikimedia.org/wikipedia/commons/thumb/1/11/Chipolata.jpg/300px-Chipolata.jpg",
				},
				Instructions: models.Instructions{
					Values: []string{
						"Boil the fresh spinach leaves and the salt in the water for about 5 minutes. If you use frozen spinach, follow the instructions on the box.",
						"Drain away the water, and let the spinach dry for a minute or so.",
						"Put the spinach back into the pan and add the cream. Simmer for a few more minutes.",
						"Add salt and pepper to taste.",
						"Boil the fresh spinach leaves as in variation I, or gently thaw the frozen chopped spinach in a little bit of water in a pan.",
						"Mix the flour with the chopped spinach in the pan, and add the milk.",
						"Bring the spinach to the boil, and simmer gently for 3–5 minutes.",
						"Add salt and pepper to taste.",
					},
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"500 g fresh spinach (or whole frozen spinach leaves)",
						"2 dl water",
						"½ dl cream or 1 tbsp butter",
						"½ tsp salt",
						"Black pepper",
						"400 g of chopped spinach, fresh or frozen",
						"1 ½ tbsp flour",
						"1 dl milk",
						"½ tsp salt",
						"Black pepper",
					},
				},
				Yield: models.Yield{Value: 4},
				URL:   "https://en.wikibooks.org/wiki/Cookbook:Creamed_Spinach",
			},
		},
		{
			name: "wikibooks.org_mobile",
			in:   "https://en.m.wikibooks.org/wiki/Cookbook:Creamed_Spinach",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "Creamed Spinach",
				Category:  models.Category{Value: "Sauce recipes"},
				Description: models.Description{
					Value: "Creamed spinach makes a nutritious sauce that goes well with fish and meat dishes. In Swedish cuisine it has traditionally been used with boiled potatoes and fish or with chipolata, the Swedish \"prince-sausage\". Creamed spinach can be done in two different ways: either with whole spinach or with chopped spinach.",
				},
				Image: models.Image{
					Value: "https://upload.wikimedia.org/wikipedia/commons/thumb/1/11/Chipolata.jpg/300px-Chipolata.jpg",
				},
				Instructions: models.Instructions{
					Values: []string{
						"Boil the fresh spinach leaves and the salt in the water for about 5 minutes. If you use frozen spinach, follow the instructions on the box.",
						"Drain away the water, and let the spinach dry for a minute or so.",
						"Put the spinach back into the pan and add the cream. Simmer for a few more minutes.",
						"Add salt and pepper to taste.",
						"Boil the fresh spinach leaves as in variation I, or gently thaw the frozen chopped spinach in a little bit of water in a pan.",
						"Mix the flour with the chopped spinach in the pan, and add the milk.",
						"Bring the spinach to the boil, and simmer gently for 3–5 minutes.",
						"Add salt and pepper to taste.",
					},
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"500 g fresh spinach (or whole frozen spinach leaves)",
						"2 dl water",
						"½ dl cream or 1 tbsp butter",
						"½ tsp salt",
						"Black pepper",
						"400 g of chopped spinach, fresh or frozen",
						"1 ½ tbsp flour",
						"1 dl milk",
						"½ tsp salt",
						"Black pepper",
					},
				},
				Yield: models.Yield{Value: 4},
				URL:   "https://en.m.wikibooks.org/wiki/Cookbook:Creamed_Spinach",
			},
		},
		{
			name: "woop.co.nz",
			in:   "https://woop.co.nz/thai-marinated-beef-sirlion-344-2-f.html",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "Thai marinated beef sirlion",
				Description: models.Description{
					Value: "with crispy noodle salad",
				},
				Yield: models.Yield{Value: 2},
				Image: models.Image{
					Value: "https://woop.co.nz/media/catalog/product/f/-/f-marinated-thai-beef-sirloin_mrypusp3a6h8fzas.jpg?quality=80&bg-color=255,255,255&fit=bounds&height=&width=800",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 pack of marinated beef sirloin steak",
						"1 pot of Thai dressing",
						"1 pack of crispy noodles",
						"1 sachet of roasted peanuts",
						"1 bag of baby leaves",
						"Cucumber",
						"1 tomato",
						"1 red onion",
						"1 bag of coriander",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"1. TO PREPARE THE SALAD:   Using half the cucumber cut in half lengthways and scoop out the seeds with a " +
							"teaspoon then slice on the diagonal into ½ cm halfmoons. Dice the tomato into 1 cm cubes then " +
							"thinly slice ½ of the red onion.",
						"2. TO FINISH THE SALAD:   Place the cucumber, tomato and red onion into a serving bowl along with the crispy " +
							"noodles and baby leaves. Roughly chop the coriander leaves and stalk then add to the serving bowl. " +
							"Add half of the Thai dressing and season with salt and pepper and toss before serving.",
						"3. TO COOK THE BEEF SIRLOIN STEAK:   Remove the marinated beef sirloin steaks from their packaging and pat " +
							"dry with a paper towel and season with salt and pepper. Heat a drizzle of oil in a non-stick frying " +
							"pan over a medium-high heat. Once hot cook the beef for 2-3 mins each side for medium-rare – a little " +
							"longer for well done. Remove from the pan and allow to rest for a few mins before slicing thinly. BBQ " +
							"Instructions: Heat your BBQ up to a medium heat. Once hot cook beef steaks for 2-3 mins each side. " +
							"Remove from the BBQ and allow to rest for a few mins before slicing thinly.",
						"TO SERVE:   Divide salad between plates then top with sliced beef. Drizzle over remaining Thai dressing and " +
							"sprinkle with roasted peanuts.",
					},
				},
				Keywords: models.Keywords{Values: "Magento, Varien, E-commerce"},
				NutritionSchema: models.NutritionSchema{
					Calories:      "2536kj (606Kcal)",
					Carbohydrates: "43g",
					Protein:       "44g",
					Fat:           "28g",
				},
				URL: "https://woop.co.nz/thai-marinated-beef-sirlion-344-2-f.html",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
