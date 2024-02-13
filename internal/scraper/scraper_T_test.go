package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_T(t *testing.T) {
	testcases := []testcase{
		{
			name: "tasteofhome.com",
			in:   "https://www.tasteofhome.com/recipes/cast-iron-skillet-steak/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				CookTime:      "PT5M",
				DateModified:  "2024-01-02",
				DatePublished: "2019-02-13",
				Description: models.Description{
					Value: `If you’ve never cooked steak at home before, it can be a little intimidating. That’s why I came up with this simple steak recipe that’s so easy, you could make it any day of the week. —<a href="https://www.tasteofhome.com/author/jschend/">James Schend</a>, <a href="https://www.dairyfreed.com/" target="_blank" rel="noopener">Dairy Freed</a>`,
				},
				Keywords: models.Keywords{Values: ""},
				Image: models.Image{
					Value: "https://www.tasteofhome.com/wp-content/uploads/2019/02/Cast-Iron-Skillet-Steak_EXPS_CIMZ19_235746_B01_15_10b-6.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 beef New York strip or ribeye steak (1 pound), 1 inch thick",
						"3 teaspoons kosher salt, divided",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"To season, start by removing the steak from the refrigerator and generously sprinkle two teaspoons of kosher salt on all sides of the fillet. Let it stand for 45 to 60 minutes. This resting period gives the meat time to come up to room temperature, which helps the steak cook more evenly. It also gives the meat time to absorb some of the salt. Editor's Tip: Use this time to whip up a homemade steak seasoning blend to take it up a notch.",
						"The other key to a delicious steak is heat. And since that signature sear comes from a sizzling hot pan, a cast-iron skillet is the way to go. This hearty pan gets extremely hot and also retains heat for a long time, making it the perfect vessel for steak. You'll want to preheat your pan over high heat for four to five minutes, or until very hot. Then, pat your steak dry with paper towels and sprinkle the remaining teaspoon of salt in the bottom of the skillet. Now you're ready to sear!",
						"Place the steak into the skillet and cook until it's easily moved. This takes between one and two minutes. Carefully flip the steak, placing it in a different section of the skillet. Cook for 30 seconds, and then begin moving the steak around, occasionally pressing slightly to ensure even contact with the skillet. A grill press is great for this. Moving the steak around the pan helps it cook faster and more evenly. Editor's Tip: This step produces a lot of smoke, so make sure you're cooking in a well-ventilated space. It's also a good idea to turn your kitchen vent or fan on.",
						"Continue turning and flipping the steak until it's cooked to your desired degree of doneness. A thermometer inserted in the thickest part of the meat should read: Medium-rare: 135°F Medium: 140°F Medium-well: 145°F Keep in mind that the steak will continue to cook a little bit after it's been removed from the pan, so aim for a few degrees shy of your desired temperature. Let it rest for 10 minutes before slicing.",
					},
				},
				Name: "Cast-Iron Skillet Steak",
				NutritionSchema: models.NutritionSchema{
					Calories:      "494 calories",
					Carbohydrates: "0 carbohydrate (0 sugars",
					Cholesterol:   "134mg cholesterol",
					Fat:           "36g fat (15g saturated fat)",
					Fiber:         "0 fiber)",
					Protein:       "40g protein. ",
					Sodium:        "2983mg sodium",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://www.tasteofhome.com/recipes/cast-iron-skillet-steak/",
			},
		},
		{
			name: "tastesbetterfromscratch.com",
			in:   "https://tastesbetterfromscratch.com/apple-crisp",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dessert"},
				CookTime:      "PT35M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-11-21T06:00:00+00:00",
				Description: models.Description{
					Value: "Our favorite Apple Crisp recipe is made with Granny Smith apples and a delicious oatmeal crumb topping.It&#39;s irresistible!",
				},
				Keywords: models.Keywords{
					Values: "apple crisp, Apple crisp ingredients, apple crisp recipe, Apple crisp recipe 9x9 pan, apple crisp recipe with oats, apple crisp topping, apple crisp with oats, Best apple crisp, easy apple crisp, Homemade apple crisp, how to make apple crisp, Old fashioned apple crisp recipe",
				},
				Image: models.Image{
					Value: "https://tastesbetterfromscratch.com/wp-content/uploads/2017/11/Apple-Crisp-Web-19.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2/3 cup old-fashioned rolled oats", "1/2 cup all-purpose flour",
						"1/2 cup light brown sugar", "1/2 teaspoon ground cinnamon",
						"1/2 teaspoon baking powder",
						"1/2 cup salted butter (, cut into small pieces)",
						"3-4 large Granny Smith apples* (, peeled and thinly sliced)",
						"3 Tablespoons salted butter (, melted)", "2 Tablespoons all-purpose flour",
						"1 Tablespoon lemon juice", "3 Tablespoons milk",
						"1/2 teaspoon vanilla extract", "1/4 cup light brown sugar",
						"1/2 teaspoon ground cinnamon", "1/4 teaspoon ground nutmeg",
						"vanilla Ice Cream", "Homemade Caramel Sauce",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 375 degrees F.",
						"Crumble Toppings: In a medium size bowl combine oats, flour, brown sugar, cinnamon, and baking powder. Add butter and cut in with a pastry blender or fork until well combined. Refrigerate while you prepare the apple filling.",
						"Apple Filling: (I use a Johnny Apple Peeler to peel, core and slice the apples all at once.) In a small bowl stir together melted butter and flour until smooth. Add lemon juice, milk and vanilla and stir. Stir in brown sugar, cinnamon, and nutmeg. Pour butter mixture over apples and toss to coat.",
						"Bake: Pour apple mixture into an 8x8-inch baking dish and spread into an even layer. Sprinkle crumble topping evenly over the apples. Bake for about 35 minutes or until golden brown and top is set. Remove from oven cool at least 15 minutes before serving.",
						"Serve with vanilla ice cream and homemade caramel sauce, if desired.",
					},
				},
				Name: "The BEST Apple Crisp",
				NutritionSchema: models.NutritionSchema{
					Calories:      "382 kcal",
					Carbohydrates: "57 g",
					Cholesterol:   "42 mg",
					Fat:           "16 g",
					Fiber:         "3 g",
					Protein:       "2 g",
					SaturatedFat:  "10 g",
					Servings:      "1",
					Sodium:        "64 mg",
					Sugar:         "38 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://tastesbetterfromscratch.com/apple-crisp",
			},
		},
		{
			name: "tastesoflizzyt.com",
			in:   "https://www.tastesoflizzyt.com/easter-ham-pie/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Brunch"},
				CookTime:      "PT75M",
				Cuisine:       models.Cuisine{Value: "Italian"},
				DatePublished: "2022-04-04T04:28:00+00:00",
				Description: models.Description{
					Value: "A midwestern take on Italian Easter Pie, this Easter Ham Pie is perfect for Sunday brunch. " +
						"It&#039;s a great leftover ham recipe.",
				},
				Keywords: models.Keywords{
					Values: "breakfast and brunch, easter breakfast, leftover ham recipe",
				},
				Image: models.Image{
					Value: "https://www.tastesoflizzyt.com/wp-content/uploads/2022/03/Easter-Ham-Pie-15.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"double pie dough for 9-inch pie pan",
						"6 large eggs (beaten)",
						"2 cups diced ham (about ½ pound)",
						"15 ounces ricotta cheese",
						"2 cups mozzarella cheese",
						"½ cup Parmesan cheese",
						"1 teaspoon oregano",
						"½ teaspoon garlic powder",
						"Salt and pepper (for topping)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Place one of the pie crusts in the bottom of a greased 9” pie plate. Flute the edges.",
						"Chill the pie crust for 30 minutes.",
						"Preheat the oven to 375°F.",
						"Place parchment paper in the bottom of the pie crust, then fill with pie weights or dry beans.",
						"Blind bake the pie crust for 15 minutes or until the edges are starting to brown. Remove the crust " +
							"from the oven and allow it to cool while you prepare the filling.",
						"In a large bowl, whisk the eggs.",
						"Then add the ham, ricotta, mozzarella, parmesan, oregano and garlic powder. Mix well.",
						"Pour the mixture into the baked bottom crust.",
						"Top with the second crust and flute the edges to seal.",
						"Use a sharp knife to make 3 slits across the top of the pie crust.",
						"Sprinkle with sea salt and freshly ground pepper.",
						"Bake at 350°F for one hour.",
						"Store any leftoveres in the refrigerator in an airtight container.",
					},
				},
				Name: "Easter Ham Pie",
				NutritionSchema: models.NutritionSchema{
					Calories:       "474 kcal",
					Carbohydrates:  "24 g",
					Cholesterol:    "192 mg",
					Fat:            "30 g",
					Fiber:          "1 g",
					Protein:        "26 g",
					SaturatedFat:   "14 g",
					Servings:       "1",
					Sodium:         "913 mg",
					Sugar:          "1 g",
					TransFat:       "1 g",
					UnsaturatedFat: "13 g",
				},
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://www.tastesoflizzyt.com/easter-ham-pie/",
			},
		},
		{
			name: "tasty.co",
			in:   "https://tasty.co/recipe/honey-soy-glazed-salmon",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Lunch"},
				Cuisine:       models.Cuisine{Value: "North American"},
				DateModified:  "2022-11-28T23:00:00",
				DatePublished: "2017-05-11T21:21:36",
				Description: models.Description{
					Value: "Two words: honey salmon! Sure, it takes a tiny bit of prep work, but once you marinate your " +
						"salmon, you won’t be able to go back. A simple mix of honey, soy sauce, garlic, and ginger coats " +
						"and flavors your fish for 30 minutes before you throw it on the pan until the outside is perfectly " +
						"crispy. Once that’s done, you heat up and reduce some extra marinade to make a thick, to-die-for " +
						"glaze to pour over your filet. Serve with your favorite veggies or rice and enjoy!",
				},
				Image: models.Image{
					Value: "https://img.buzzfeed.com/video-api-prod/assets/04ff8cfcc4b5428a8bcc6b03099d4492/Thumb_A_FB.jpg?" +
						"resize=1200:*",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"12 oz skinless salmon",
						"1 tablespoon olive oil",
						"4 cloves garlic, minced",
						"2 teaspoons ginger, minced",
						"½ teaspoon red pepper",
						"1 tablespoon olive oil",
						"⅓ cup less sodium soy sauce",
						"⅓ cup honey",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Place salmon in a sealable bag or medium bowl.",
						"In a small bowl or measuring cup, mix marinade ingredients.",
						"Pour half of the marinade on the salmon. Save the other half for later.",
						"Let the salmon marinate in the refrigerator for at least 30 minutes.",
						"In a medium pan, heat oil. Add salmon to the pan, but discard the used marinade. Cook salmon on one " +
							"side for about 2-3 minutes, then flip over and cook for an additional 1-2 minutes.",
						"Remove salmon from pan. Pour in remaining marinade and reduce.",
						"Serve the salmon with sauce and a side of veggies. We used broccoli.",
						"Enjoy!",
					},
				},
				Name: "Honey Soy-Glazed Salmon Recipe by Tasty",
				NutritionSchema: models.NutritionSchema{
					Calories:      "705 calories",
					Carbohydrates: "60 grams",
					Fat:           "35 grams",
					Fiber:         "0 grams",
					Protein:       "37 grams",
					Sugar:         "57 grams",
				},
				Yield: models.Yield{Value: 2},
				URL:   "https://tasty.co/recipe/honey-soy-glazed-salmon",
			},
		},
		{
			name: "tastykitchen.com",
			in:   "https://tastykitchen.com/recipes/main-courses/garlic-shrimp-scampi-with-angel-hair-pasta/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "Garlic Shrimp Scampi with Angel Hair Pasta",
				Category:  models.Category{Value: "Main Courses"},
				Description: models.Description{
					Value: "This shrimp scampi with angel hair pasta is a delicious dinner with plenty of cheese that you can " +
						"make in just 15 minutes! ",
				},
				Image: models.Image{
					Value: "https://tastykitchen.com/recipes/wp-content/uploads/sites/2/2020/09/SHRIMP-SCAMPI-WITH-ANGEL-HAIR-" +
						"PASTA-15-410x308.jpg",
				},
				PrepTime: "PT5M",
				CookTime: "PT10M",
				Ingredients: models.Ingredients{
					Values: []string{
						"3 Tablespoons Butter",
						"4 cloves Garlic, Minced",
						"1 pound Large Shrimp, Peeled & Deveined",
						"½ teaspoons Red Pepper Flakes (or To Taste)",
						"1 pinch Sea Salt",
						"1 piece Lemon, Zested",
						"⅓ cups Freshly Squeezed Lemon Juice",
						"6 cups Vegetable Or Fish Stock",
						"12 ounces, weight Angel Hair Pasta",
						"⅔ cups Freshly Grated Parmesan Cheese",
						"⅓ cups Finely Chopped Fresh Parsley",
						"Salt And Pepper, to taste",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Melt the butter in a large skillet over medium heat. Add the garlic and cook until fragrant, about " +
							"1 minute, stirring constantly. Add in the shrimp and sprinkle with chili flakes and sea salt.",
						"Cook the shrimp for about 1 to 2 minutes per side, or until pink and opaque. Set aside on a plate.",
						"To the same skillet, add the lemon zest and juice. Pour in 5 cups of broth. Keep the remaining 1 cup " +
							"on hand, in case you need it for later. Stir to combine then bring the liquid to a gentle simmer.",
						"Add the pasta, and using a pair of tongs, stir occasionally for a few minutes, or until the pasta starts " +
							"to soften and bend. Fully immerse the pasta into the liquid, and keep stirring frequently to avoid " +
							"sticking to the bottom of the pan.",
						"Cook the pasta according to the timing on the package directions. Angel hair pasta needs no more than 3 " +
							"minutes to cook to al dente.",
						"Once the pasta is cooked to al dente, return the shrimp to the pan and turn the heat off. Stir in the " +
							"Parmesan cheese and chopped parsley, and mix until the cheese has melted into the sauce. The " +
							"pasta will continue to absorb the sauce if not served immediately. So quickly season with salt " +
							"and pepper to your taste, and dig in!",
					},
				},
				Yield: models.Yield{Value: 4},
				URL:   "https://tastykitchen.com/recipes/main-courses/garlic-shrimp-scampi-with-angel-hair-pasta/",
			},
		},
		{
			name: "tesco.com",
			in:   "https://realfood.tesco.com/recipes/salted-honey-and-rosemary-lamb-with-roasties-and-rainbow-carrots.html",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Name:          "Salted honey and rosemary lamb with roasties and rainbow carrots",
				Category:      models.Category{Value: "Dinner"},
				CookTime:      "PT2H0M",
				Cuisine:       models.Cuisine{Value: "British"},
				DateModified:  "2022-04-06T10:26:45Z",
				DatePublished: "2022-03-23T14:55:51Z",
				Description: models.Description{
					Value: "This one-tray wonder is a great option for Sunday roast, Mother's Day or even Easter",
				},
				Keywords: models.Keywords{
					Values: "honey, lamb, easter, mother's Day, sunday roast, sunday lunch, roast dinner, veg, meat and two veg, potatoes, roast potatoes, roasted potatoes",
				},
				Image: models.Image{
					Value: "https://realfood.tesco.com/media/images/1400x919-OneTrayEasterRoast-7f95aa82-f903-4339-ab87-" +
						"6cffda6d26d8-0-1400x919.jpg",
				},
				Yield: models.Yield{Value: 6},
				NutritionSchema: models.NutritionSchema{
					Calories:      "855 calories",
					Carbohydrates: "64.2 grams carbohydrate",
					Cholesterol:   "",
					Fat:           "44 grams fat",
					Fiber:         "10.4 grams fibre",
					Protein:       "51.2 grams protein",
					SaturatedFat:  "16 grams saturated fat",
					Sugar:         "21 grams sugar",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2kg lamb leg joint", "¼ tsp flaky sea salt, plus extra to season",
						"20g pack fresh rosemary, most leaves finely chopped, some whole sprigs",
						"12 garlic cloves, crushed", "7 tbsp olive oil",
						"1.5kg King Edward potatoes, peeled and cut into 4-5cm chunks",
						"2 x 450g packs Tesco Finest rainbow or Imperator carrots, scrubbed, halved lengthways if large",
						"50g clear honey, to taste",
						"20g fresh mint, leaves picked",
						"1 tbsp granulated sugar",
						"3 tbsp white wine vinegar",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat the oven to gas 3, 160 ̊C, fan 140 ̊C. Use a small sharp knife to pierce holes about 4cm deep all " +
							"over the lamb, then season well with sea salt. In a small bowl, combine the chopped rosemary and " +
							"garlic with 6 tbsp olive oil.",
						"Rinse the chopped potatoes and tip into one side of a large, rimmed baking tray. Put the carrots in the " +
							"other half of the tray. Toss the carrots with 1 tbsp oil, then toss the potatoes in half the rosemary- " +
							"garlic oil and spread out in a single layer over half of the tray; add the whole rosemary sprigs. Put " +
							"the lamb in the centre of the tray and rub with the remaining rosemary-garlic oil.",
						"Roast for 1 hr 45 mins. Mix the honey with &frac14; tsp flaky sea salt and dab all over the lamb with a pastry " +
							"brush. Return to the oven, increase the temperature to gas 7, 220 ̊C, fan 200 ̊C, then roast for a " +
							"final 15 mins. Transfer the lamb to a warmed plate, then set aside to rest, covered loosely with " +
							"foil, for at least 20 mins. Keep the veg warm while the lamb rests.",
						"Meanwhile, make the mint sauce. Finely chop the mint leaves with the sugar (this helps to stop oxidisation), " +
							"then mix well in a bowl with 3 tbsp boiling water and the vinegar. Serve alongside the roast lamb, " +
							"potatoes and carrots.",
					},
				},
				PrepTime: "PT1H0M",
				URL:      "https://realfood.tesco.com/recipes/salted-honey-and-rosemary-lamb-with-roasties-and-rainbow-carrots.html",
			},
		},
		{
			name: "theclevercarrot.com",
			in:   "https://www.theclevercarrot.com/2021/10/homemade-sourdough-breadcrumbs/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Sourdough Bread"},
				CookTime:      "PT25M",
				CookingMethod: models.CookingMethod{Value: "Oven-Baked"},
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2021-10-03",
				Description: models.Description{
					Value: "Transform leftover bread into delicious homemade breadcrumbs, with just 10 minutes hands on time and " +
						"minimal effort. I like using sourdough bread, but any bread will do. Breadcrumbs can be stored " +
						"in the freezer for up to 3-6 months.",
				},
				Keywords: models.Keywords{
					Values: "homemade, sourdough bread, breadcrumbs, Italian, seasoned bread crumbs",
				},
				Image: models.Image{
					Value: "https://www.theclevercarrot.com/wp-content/uploads/2021/09/SOURDOUGH-BREADCRUMBS-2-225x225.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 tsp dried garlic powder",
						"1 tsp dried onion powder",
						"1 tsp fine sea salt",
						"2 tsp dried oregano",
						"1 tbsp dried parsley",
						"1/4 cup (30 g) ground Pecorino Romano or Parmesan cheese",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat your oven to 300˚ F (150˚ C). Grab (2x) rimmed baking sheets.",
						"Cut the bread into small cubes, including the crust, about 1-inch in size.",
						"Add the cubes to a food processor or high-powdered blender (you will need to work in batches). Process " +
							"until fine crumbs form. Note: if you cannot get the crumbs small enough at this stage, you " +
							"can process them again after baking.",
						"Divide the crumbs over the (2x) sheet pans in one even layer.",
						"Bake for 15-30 minutes, stirring once, or until the crumbs are crispy. Bake time will vary depending on " +
							"bread type and freshness.",
						"Remove the breadcrumbs from the oven. Allow to cool.",
						"At this point, you can either keep the breadcrumbs plain or season, Italian-style (my preference). Add the " +
							"dried garlic, onion, salt, oregano, parsley and cheese; mix well.",
						"Portion the cooled breadcrumbs into containers or a zip-top bag and freeze until ready to use.",
					},
				},
				Name:     "Homemade Sourdough Breadcrumbs",
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.theclevercarrot.com/2021/10/homemade-sourdough-breadcrumbs/",
			},
		},
		{
			name: "theexpertguides.com",
			in:   "https://theexpertguides.com/recipes/is-guacamole-vegan/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Side Dish"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "Mexican"},
				DatePublished: "2023-09-11T19:31:32+00:00",
				Description: models.Description{
					Value: "This vegan guacamole recipe yields a creamy and tangy dip that you can eat with tortilla chips, bread, or veggies. It's quick and requires a few basic ingredients.",
				},
				Keywords: models.Keywords{Values: "avocado guacamole recipe, guacamole, vegan guacamole, vegan salsa"},
				Image: models.Image{
					Value: "https://theexpertguides.com/wp-content/uploads/2022/09/homemade-vegan-guacamole.webp",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 Avocados", "2 tbsp fresh lime juice", "2 garlic cloves (minced)",
						"10 cherry tomatoes (sliced into small pieces)",
						"1/4 cup red onion (finely chopped)", "1/4 cup fresh cilantro (chopped)",
						"Salt and black pepper to taste",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Slice the avocados into small cubes in a mixing bowl, and mash them with a fork.",
						"Add fresh lime juice and mix in the avocados.",
						"Add sliced tomatoes, garlic, onion, and cilantro. Continue mashing with a fork.",
						"Add salt and black pepper to taste.",
						"Serve with sliced veggies, tortilla chips, or fajitas.",
					},
				},
				Name:     "Homemade Vegan Guacamole Recipe",
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://theexpertguides.com/recipes/is-guacamole-vegan/",
			},
		},
		{
			name: "thehappyfoodie.co.uk",
			in:   "https://thehappyfoodie.co.uk/recipes/leek-and-lentil-gratin/",
			want: models.RecipeSchema{
				AtContext:    atContext,
				AtType:       models.SchemaType{Value: "Recipe"},
				DateModified: "2022-02-07T16:00:36+00:00",
				Name:         "Leek and Puy Lentil Gratin with a Crunchy Feta Topping",
				Yield:        models.Yield{Value: 4},
				Description: models.Description{
					Value: "A creamy, cheesy, and deeply comforting dish, Rukmini Iyer's one-tin vegetarian gratin also happens " +
						"to be packed with nutritious Puy lentils.",
				},
				Image: models.Image{
					Value: "https://thehappyfoodie.co.uk/wp-content/uploads/2022/02/Rukmini-Iyer-Leek-and-Lentil-Gratin-scaled.jpg",
				},
				PrepTime: "PT10M",
				CookTime: "PT40M",
				Keywords: models.Keywords{
					Values: "Vegetarian, Feta, Leek, Lentil, One Pot, One-pot, Dinner, Main Course, Easy, Quick",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"30g butter", "3 cloves of garlic, crushed", "500g leeks, thinly sliced",
						"2 tsp sea salt", "freshly ground black pepper",
						"500g vac-packed cooked Puy lentils", "450ml crème fraîche",
						"125g feta cheese, crumbled", "50g panko or fresh white breadcrumbs",
						"1 tbsp olive oil",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat the oven to 180°C fan/200°C/gas 6. Put the butter and garlic into a roasting tin and pop into " +
							"the oven to melt while you get on with slicing the leeks.",
						"Mix the sliced leeks with the melted garlic butter, season well with the sea salt and black pepper, then " +
							"return to the oven to roast for 20 minutes.",
						"After 20 minutes, stir through the Puy lentils, crème fraîche and another good scatter of sea salt, then " +
							"top with the feta cheese and breadcrumbs. Drizzle with the olive oil, then return to the oven " +
							"for a further 20–25 minutes, until golden brown on top.",
						"Serve the gratin hot, with a mustard or balsamic dressed green salad alongside.",
					},
				},
				URL: "https://thehappyfoodie.co.uk/recipes/leek-and-lentil-gratin/",
			},
		},
		{
			name: "thekitchencommunity.org",
			in:   "https://thekitchencommunity.org/turkey-salad-recipe/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				DateModified:  "2023-12-08T11:47:42-05:00",
				DatePublished: "2023-11-23T12:23:49-05:00",
				Description: models.Description{
					Value: "Are you looking for a great way to use that leftover turkey from your holiday feast? Look no further – we have the perfect turkey salad recipe for you! This easy turkey salad is not only a delicious option, but it's also a fantastic way to ensure that none of your delicious turkey goes to",
				},
				Image: models.Image{
					Value: "https://thekitchencommunity.org/wp-content/uploads/2023/11/Turkey-Salad-1200x800.jpg",
				},
				Name:  "Turkey Salad Recipe",
				URL:   "https://thekitchencommunity.org/turkey-salad-recipe/",
				Yield: models.Yield{Value: 1},
			},
		},
		{
			name: "thekitchenmagpie.com",
			in:   "https://www.thekitchenmagpie.com/blt-pasta-salad/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Salad"},
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2021-06-16T04:00:00+00:00",
				Description: models.Description{
					Value: "Fantastic BLT pasta salad with bacon, lettuce, tomatoes, two types of cheese and tangy Ranch " +
						"dressing! The perfect side dish for BBQ&#039;s, or eat it as a meal!",
				},
				Keywords: models.Keywords{Values: "BLT pasta salad"},
				Image: models.Image{
					Value: "https://www.thekitchenmagpie.com/wp-content/uploads/images/2021/05/BLTpastasalad.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"one 454 gram box Farfalle pasta (bow-tie pasta cooked until al dente)",
						"1 pound bacon diced and cooked",
						"3 cups chopped romaine lettuce",
						"2 cups red cherry tomatoes (halved)",
						"1 cup orange or yellow cherry tomatoes (halved)",
						"1 cup cheddar cheese ( cut into small cubes then measure)",
						"1 cup pepper Jack cheese ( cut into small cubes then measure)",
						"1/4 cup red onion (sliced into very thin almost transparent rings)",
						"2 avocados  (pitted, peeled then diced)",
						"1 1/2 - 2 cups ranch dressing",
						"2 tablespoons fresh parsley chopped",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a large serving bowl place the cooked pasta, then pour one cup of the dressing on top. Stir " +
							"to coat the pasta.",
						"Add in the cooked bacon, romaine lettuce, two types of tomatoes, the two types of cheese, the red " +
							"onion, and the avocados.",
						"Pour the remaining half a cup of dressing over, then gently mix until the other ingredients are " +
							"slightly coated. Taste test, and add more dressing if desired.",
						"Garnish with the fresh parsley and serve.",
					},
				},
				Name: "BLT Pasta Salad",
				NutritionSchema: models.NutritionSchema{
					Calories:      "653 kcal",
					Carbohydrates: "10 g",
					Cholesterol:   "78 mg",
					Fat:           "60 g",
					Fiber:         "3 g",
					Protein:       "20 g",
					SaturatedFat:  "15 g",
					Servings:      "1",
					Sodium:        "1533 mg",
					Sugar:         "3 g",
					TransFat:      "1 g",
				},
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 12},
				URL:      "https://www.thekitchenmagpie.com/blt-pasta-salad/",
			},
		},
		{
			name: "thekitchn.com",
			in:   "https://www.thekitchn.com/how-to-reheat-turkey-and-keep-it-moist-251033",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main dish"},
				DateCreated:   "2017-11-24T12:30:00-05:00",
				DateModified:  "2023-10-24T19:30:00-04:00",
				DatePublished: "2017-11-24T12:30:00-05:00",
				Description: models.Description{
					Value: "Learn the best way to reheat leftover turkey this weekend and a few tips for better microwaved leftovers.",
				},
				Keywords: models.Keywords{
					Values: "anheuser busch thanksgiving,butter,christmas,Cooking Lessons from The Kitchn,dinner,easy,Gluten-Free,How To,Ingredient,Lunch,Main Dish,Meat Tutorials,poultry,recipe,thanksgiving,Tips & Techniques,Turkey,Wellness,Dish Types,team:food,purpose:search,platform:email,platform:homepage,platform:newsfeeds,platform:search,platform:social,course:main dish,diet:under 500 calories,diet:gluten-free,diet:under 200 calories,diet:low-calorie,diet:under 400 calories,diet:low-fat,diet:under 300 calories,mainingredients:poultry,mainingredients:turkey,dishtype:poultry dish,meal:dinner,meal:lunch,updated:2023-10-24",
				},
				Image: models.Image{
					Value: "https://cdn.apartmenttherapy.info/image/upload/f_jpg,q_auto:eco,c_fill,g_auto,w_1500,ar_16:9/k%2Farchive%2F6cc517c3ad0e6cce62ca2954c6f6b400fca13bc1",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 pound cooked turkey (breast meat, thigh meat, or whole drumsticks)",
						"1 cup low-sodium turkey, vegetable, or chicken broth",
						"1 tablespoon unsalted butter, cut into 4 pieces",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Heat the oven to 350°F. Arrange a rack in the middle of the oven and heat to 350°F.",
						"Slice the turkey and spread into a baking dish. Uniform slices or bite-sized pieces reheat quickly and evenly, so slice the turkey as needed. Place the turkey in a a single layer in an 8x8-inch baking dish.",
						"Cover the turkey with broth and dot with the butter. Add the broth and dot the turkey with the pieces of butter. Cover the baking dish with aluminum foil.",
						"Reheat in the oven for 30 to 35 minutes. Bake until the turkey is heated through, 30 to 35 minutes. Remember that all leftovers should be reheated to a minimum of 165°F.",
						"Serving. Uncover and serve the turkey immediately.",
					},
				},
				Name: "The Best Method for Reheating Turkey So It Never Dries Out",
				NutritionSchema: models.NutritionSchema{
					Calories:       "194 cal",
					Carbohydrates:  "2.3 g",
					Cholesterol:    "0 mg",
					Fat:            "9.1 g",
					Fiber:          "0 g",
					Protein:        "26.1 g",
					SaturatedFat:   "4.3 g",
					Servings:       "Serves 4",
					Sodium:         "175.8 mg",
					Sugar:          "1.0 g",
					UnsaturatedFat: "0.0 g",
				},
				Yield: models.Yield{Value: 4},
				URL:   "https://www.thekitchn.com/how-to-reheat-turkey-and-keep-it-moist-251033",
			},
		},
		{
			name: "themagicalslowcooker.com",
			in:   "https://www.themagicalslowcooker.com/slow-cooker-apple-cider-pot-roast/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Course"},
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-11-09T07:25:41+00:00",
				Description: models.Description{
					Value: "Give your pot roast a tasty upgrade by using apple cider - a taste the whole family is sure to enjoy.",
				},
				Keywords: models.Keywords{Values: "apple cider pot roast"},
				Image: models.Image{
					Value: "https://www.themagicalslowcooker.com/wp-content/uploads/2023/10/apple-cider-pot-roast-57.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 lbs. beef chuck roast", "1 tsp. salt", "¼ tsp. pepper",
						"2 Tbsp. vegetable oil", "1 cup apple cider", "1 cup beef broth",
						"2 Tbsp. Worcestershire sauce", "4 cloves garlic (Minced)",
						"1 Tbsp. tomato paste", "1 Tbsp. dijon mustard",
						"2 tsp. chopped fresh rosemary (or 1 teaspoon dried)",
						"2 tsp. chopped fresh thyme leaves (or 1 teaspoon dried)",
						"1 lbs. baby carrots",
						"1 lbs. baby potatoes (halved or quartered depending on size)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Season the chuck roast on both sides with salt and pepper. You can use more or less depending on your preferences.",
						"Heat the vegetable oil in a cast iron skillet or large pan over medium-high heat. Sear the chuck roast until browned on all sides, about 3-4 minutes per side. Remove the roast from the skillet and set aside.",
						"In a large mixing bowl, whisk together the apple cider, beef broth, Worcestershire sauce, minced garlic, tomato paste, dijon mustard, chopped rosemary, and chopped thyme. Whisk until thoroughly combined and set aside.",
						"Place the baby carrots and baby potatoes in the bottom of a slow cooker.",
						"Carefully transfer the seared chuck roast onto the bed of vegetables in the slow cooker.",
						"Pour the apple cider mixture over the roast and vegetables in the slow cooker.",
						"Cover the slow cooker and cook on LOW heat for 8-10 hours or on HIGH heat for 5-6 hours, until the meat is tender and easily shreds with a fork.",
						"Once cooked, remove the chuck roast from the slow cooker and let it rest for a few minutes before slicing or shredding it.",
						"Serve the pot roast with the cooked baby carrots, potatoes and a ladle of the flavorful cooking liquid, Garnish with fresh thyme, parsley and rosemary sprigs. Enjoy!",
					},
				},
				Name: "Slow Cooker Apple Cider Pot Roast",
				NutritionSchema: models.NutritionSchema{
					Calories:       "565 kcal",
					Carbohydrates:  "26 g",
					Cholesterol:    "156 mg",
					Fat:            "31 g",
					Fiber:          "4 g",
					Protein:        "46 g",
					SaturatedFat:   "12 g",
					Servings:       "1",
					Sodium:         "908 mg",
					Sugar:          "9 g",
					TransFat:       "2 g",
					UnsaturatedFat: "19 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://www.themagicalslowcooker.com/slow-cooker-apple-cider-pot-roast/",
			},
		},
		{
			name: "thenutritiouskitchen.co",
			in:   "http://thenutritiouskitchen.co/fluffy-paleo-blueberry-pancakes/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "breakfast"},
				CookTime:      "PT15M",
				CookingMethod: models.CookingMethod{Value: "stove top"},
				Cuisine:       models.Cuisine{Value: "pancakes"},
				DatePublished: "2022-04-08",
				Description: models.Description{
					Value: "Delicious paleo blueberry pancakes made extra fluffy with simple, healthy ingredients all made in " +
						"one bowl! The perfect weekend breakfast or brunch recipe packed with refreshing blueberry flavors, " +
						"no added sugars and 100% dairy-free + gluten-free!",
				},
				Keywords: models.Keywords{
					Values: "paleo, pancakes, gluten-free, dairy-free",
				},
				Image: models.Image{
					Value: "https://thenutritiouskitchen.co/wp-content/uploads/2022/04/paleoblueberrypancakes-225x225.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3/4 cup unsweetened apple sauce",
						"1/4 cup cashew butter",
						"1 large egg",
						"1 cup Bobs Red Mill super fine almond flour",
						"1 cup tapioca flour, or arrowroot flour",
						"2 teaspoons baking powder",
						"Sea salt + cinnamon to taste",
						"1/3 cup fresh blueberries",
						"Pure maple syrup",
						"Coconut cream (optional)",
						"Vegan butter or ghee",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a large bowl whisk the apple sauce, cashew butter and egg until fully combined.",
						"Then add in the almond flour, tapioca flour, baking powder, sea salt, and cinnamon. Combine with a mixing " +
							"spoon until a thick batter forms. It will be thick but still moist!",
						"Heat ghee or butter in a pan over medium heat. Once the pan is hot, scoop the batter using a cookie scoop " +
							"or about 1/4 cup full of batter. I like to cook 3 at a time.",
						"Place some blueberries onto the pancakes while on pan and cook for about 3 minutes. Flip gently, and " +
							"continue cooking for about 2 minutes over medium-low heat.",
						"Repeat with remaining batter then serve with maple syrup, extra butter and toppings of choice. I love " +
							"coconut cream or cashew butter!",
					},
				},
				Name:     "Fluffy Paleo Blueberry Pancakes",
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://thenutritiouskitchen.co/fluffy-paleo-blueberry-pancakes/",
			},
		},
		{
			name: "themodernproper.com",
			in:   "https://themodernproper.com/turkey-pozole-rojo",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "dinner, soup, diary-free, gluten-free, one-pot-meals"},
				CookTime:      "PT1H",
				Cuisine:       models.Cuisine{Value: "Mexican"},
				DateModified:  "2023-11-23",
				DatePublished: "2023-11-23",
				Description: models.Description{
					Value: "This spicy, warming turkey pozole rojo soup is a flavorful way to use all of your turkey leftovers. Easy, and warming, it’s perfect palate-awakener for the day after Thanksgiving.",
				},
				Keywords: models.Keywords{
					Values: "turkey pozole rojo, fall recipe, thanksgiving recipe, winter recipe, diary-free recipe, gluten-free recipe, one-pot-meals",
				},
				Image: models.Image{
					Value: "https://images.themodernproper.com/billowy-turkey/production/posts/TurkeyPozoleRojo_12.jpg?w=960&h=540&q=82&fm=jpg&fit=crop&dm=1700147273&s=7920925d561e1b3ec624e434447b8d39",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 tablespoon extra-virgin olive oil", "1 large yellow onion, diced",
						"7 garlic cloves, smashed", "4 large dried ancho chiles",
						"4 large dried guajillo chiles", "8 cups turkey or chicken stock",
						"1 ½ teaspoons dried oregano, preferably Mexican oregano",
						"2 teaspoons cumin", "2 bay leaves", "1½ teaspoons kosher salt",
						"2 (28-ounce) can Mexican hominy, drained and rinsed",
						"4-6 cups shredded cooked turkey (or chicken or pork)",
						"Seasoned cabbage (see recipe below), for serving",
						"Avocado, diced, for serving", "Cilantro leaves, for serving",
						"Radishes, thinly sliced, for serving",
						"Fresh lime juice, for serving",
						"Tostadas, for serving",
						"3 cups finely shredded cabbage",
						"1 tablespoon vegetable or olive oil",
						"3 tablespoons white vinegar",
						"1 teaspoon kosher salt",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Cut open the dried chiles with kitchen scissors or a knife and discard the seeds",
						"Heat a large dry skillet over medium-high heat",
						"Toast the chiles, moving them occasionally, for 4 minutes",
						"(See Note)Cover the chiles with 2 inches of water",
						"Bring the water to a boil over high heat, then turn the heat off and let the chilis soak until softened, about 15 minutes.Add chilies and garlic along with 1 ½ cups of the cooking liquid to a high speed blender and blend until smooth",
						"Make the seasoned cabbage",
						"In a medium bowl, combine cabbage, olive oil, vinegar, salt, and pepper and toss until combined.Heat the oil in a medium pot over medium heat",
						"Once the oil is glistening, add the onions and cook until translucent, 4-5 minutes",
						"To the same pot over medium heat, add the blended chilis, chicken stock, cumin, oregano, bay leaves and salt",
						"Increase the heat to high and bring to a simmer", "Cook for 30 minutes",
						"Stir in the hominy and the turkey and continue to cook until the turkey is warmed through and the hominy has absorbed the broth, 20 minutes more",
						"Taste and adjust season with salt.Serve the pozole in bowls topped with seasoned cabbage, avocado, cilantro, radishes, and lime juice.",
					},
				},
				Name: "Turkey Pozole Rojo",
				NutritionSchema: models.NutritionSchema{
					Calories:      "406 calories",
					Carbohydrates: "47 grams carbohydrates",
					Cholesterol:   "0 milligrams cholesterol",
					Fat:           "6 grams fat",
					Fiber:         "13 grams fiber",
					Protein:       "8 grams protein",
					SaturatedFat:  "1 grams saturated fat",
					Servings:      "8",
					Sodium:        "1800 milligrams sodium",
					Sugar:         "23 grams sugar",
				},
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://themodernproper.com/turkey-pozole-rojo",
			},
		},
		{
			name: "thepioneerwoman.com",
			in:   "https://www.thepioneerwoman.com/food-cooking/recipes/a8865/eggs-benedict/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "brunch"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "American"},
				DateModified:  "2023-03-22T14:48:00Z EST",
				DatePublished: "2007-10-12T09:33:50Z EST",
				Description: models.Description{
					Value: "Ree Drummond shares her secrets to making perfect eggs Benedict. From flawless poached eggs to velvety Hollandaise sauce, it's the best brunch recipe.",
				},
				Keywords: models.Keywords{Values: "Recipes, Cooking, Food"},
				Image: models.Image{
					Value: "https://hips.hearstapps.com/thepioneerwoman/wp-content/uploads/2007/10/1546875357_506daa8f1c.jpg?crop=0.664xw:1xh;center,top&resize=1200:*",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 whole English muffins",
						"6 slices Canadian bacon",
						"2 sticks butter, plus more for the muffins",
						"6 whole eggs (plus 3 egg yolks)",
						"1 whole lemon, juiced",
						"Cayenne pepper, to taste",
						"Paprika, to garnish",
						"Chopped chives, to garnish",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Bring a pot of water to a boil. While the water's boiling, place the English muffin halves and an equal number of Canadian bacon slices on a cookie sheet. Lightly butter the English muffins and place the sheet under the broiler for just a few minutes, or until the English muffins are very lightly golden. Be careful not to dry out the Canadian bacon.",
						"Now if you do not have an egg poacher you can poach your eggs by doing the following: With a spoon, begin stirring the boiling water in a large, circular motion. When the tornado's really twisting, crack in an egg. The reason for the swirling is so the egg will wrap around itself as it cooks, keeping it together. Cook for about 2 1/2 to 3 minutes. Repeat with the remaining eggs.",
						"Melt 2 sticks of butter in a small saucepan until sizzling, but don't let it burn! Separate three eggs and place the yolks into a blender. Turn the blender on low to allow the yolks to combine, then begin pouring the very hot butter in a thin stream into the blender. The blender should remain on the whole time, and you should be careful to pour in the butter very slowly. Keep pouring butter until it’s all gone, then immediately begin squeezing lemon juice into the blender. If you are going to add cayenne pepper, this is the point at which you would do that.",
						"Place the English muffins on the plate, face up. Next, place a slice of Canadian bacon on each half. Place an egg on top of the bacon and then top with a generous helping of Hollandaise sauce. <em>Vegetarian variation:</em> you can omit the Canadian bacon altogether, or you can wilt fresh spinach and place it on the muffins for Eggs Florentine, which is divine in its own right. Top with more cayenne, or a sprinkle of paprika, and chopped chives if you like.",
					},
				},
				Name:     "Eggs Benedict",
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 3},
				URL:      "https://www.thepioneerwoman.com/food-cooking/recipes/a8865/eggs-benedict/",
			},
		},
		{
			name: "therecipecritic.com",
			in:   "https://therecipecritic.com/avocado-egg-rolls/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizer"},
				CookTime:      "PT5M",
				Cuisine:       models.Cuisine{Value: "Asian American"},
				DatePublished: "2019-05-06T06:00:02+00:00",
				Description: models.Description{
					Value: "Avocado Egg Rolls are crispy on the outside with an avocado mixture&nbsp;inside that is bursting with flavor!",
				},
				Image: models.Image{
					Value: "https://therecipecritic.com/wp-content/uploads/2019/05/avocado_egg_roll7.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 cup vegetable oil (for frying)", "3 avocados (diced)",
						"1/4 cup red onion (diced)", "1 Roma tomato (diced)",
						"3 tbsp chopped fresh cilantro leaves", "1 teaspoon garlic powder",
						"Juice of 1 lime", "salt and pepper to taste", "8 egg roll wrappers",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a large skillet heat the olive oil to medium high heat.",
						"In a medium bowl, add the avocado and mash to desired consistency. Add the onion, tomato, cilantro, garlic powder, lime juice and salt and pepper to taste.",
						"To make the egg rolls: Place the avocado mixture in the center of each wrapper. Using your finger, rub the edges with water. Bring the bottom edge of the wrapper and roll it tightly over the filling. Fold in the sides and continue to roll up the wrapper and press to seal. Repeat until you have used all of the wrappers.",
						"Add the. egg rolls to the hot oil and fry until they are golden brown on all sides for about 2-3 minutes. Remove with a metal tong onto a paper towel lined plate. Serve immediately with favorite dipping sauce.",
					},
				},
				Name: "The Best Avocado Egg Rolls",
				NutritionSchema: models.NutritionSchema{
					Calories:      "187 kcal",
					Carbohydrates: "15 g",
					Cholesterol:   "1 mg",
					Fat:           "14 g",
					Fiber:         "5 g",
					Protein:       "3 g",
					SaturatedFat:  "4 g",
					Servings:      "1",
					Sodium:        "80 mg",
					Sugar:         "1 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://therecipecritic.com/avocado-egg-rolls/",
			},
		},
		{
			name: "thespruceeats.com",
			in:   "https://www.thespruceeats.com/pasta-with-anchovies-and-breadcrumbs-recipe-5215384",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				CookTime:      "PT28M",
				Cuisine:       models.Cuisine{Value: "Italian"},
				DateModified:  "2023-02-20T21:36:14.876-05:00",
				DatePublished: "2022-01-12T16:19:41.204-05:00",
				Description: models.Description{
					Value: "Anchovies and pepper provide complex flavor in this simple Sicilian-style pasta dish. Crisp garlicky " +
						"breadcrumbs are the perfect finishing touch.",
				},
				Keywords: models.Keywords{Values: "anchovie pasta"},
				Image: models.Image{
					Value: "https://www.thespruceeats.com/thmb/_al2WJ0fr7Kt_LFIVpq-jQtimk0=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/pasta-with-anchovies-and-breadcrumbs-recipe-5215384-Hero_01-a52a47010bd04ead814b972d518738ef.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"For the Breadcrumbs:",
						"4 slices day-old bread, about 7 ounces",
						"3 tablespoons olive oil",
						"5 cloves garlic, minced",
						"1/2 teaspoon kosher salt, or to taste",
						"1 teaspoon lemon zest, optional",
						"For the Pasta:",
						"12 ounces dry spaghetti or linguine",
						"1 1/2 tablespoons kosher salt for the cooking water, plus more, to taste",
						"3 tablespoons olive oil",
						"2 cloves garlic, cut in half lengthwise",
						"3 oil-packed anchovy fillets",
						"1/2 teaspoon crushed red pepper flakes, or more, to taste",
						"1 tablespoon lemon juice, optional",
						"1/4 teaspoon freshly ground black pepper, or to taste",
						"For Serving:",
						"3 tablespoons chopped Italian flat-leaf parsley",
						"1/4 cup freshly grated Parmesan cheese, optional",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Gather the ingredients.",
						"Bring a large pot of well salted water to a boil. While the water is heating, remove the crusts from the " +
							"bread, if you like, and tear or cut it into cubes.",
						"Put bread cubes, garlic, and 1/4 teaspoon of kosher salt in a food processor and pulse to make coarse " +
							"breadcrumbs. You should have about 2 cups of crumbs, a bit more if the crusts are not removed.",
						"Heat the 3 tablespoons of the oil in a large Dutch oven or other heavy-duty pot over medium heat. When " +
							"the oil shimmers, add the breadcrumbs. Cook, stirring frequently, until the crumbs are lightly brown " +
							"and crisp, about 5 to 12 minutes, depending on the moisture in the bread.",
						"Transfer the breadcrumbs to a bowl, toss with lemon zest, if using, and set aside.",
						"Meanwhile, cook the pasta according to al dente according to package instructions, reserving 1/2 cup of " +
							"the pasta water. Drain the pasta and set aside.",
						"Wipe out the pot used for the breadcrumbs. Add 3 tablespoons of the oil over medium heat until it shimmers. " +
							"Add the halved garlic cloves and cook until lightly brown, about 2 minutes.",
						"Remove and discard the garlic pieces. Add the anchovies and crushed red pepper to the garlic-flavored oil " +
							"and cook for 1 minute longer. Add the lemon juice, if using.",
						"Add the pasta to the pot and toss, cooking, until warmed through. Add some cooking water to loosen the " +
							"mixture as needed. Taste and adjust the seasonings with salt and more crushed red pepper flakes, if desired.",
						"Plate the pasta in wide, shallow pasta bowls and top with a generous amount of garlic breadcrumbs. Garnish " +
							"with chopped parsley and Parmesan cheese, if using. Serve with lemon wedges, if desired.",
					},
				},
				Name: "Pasta With Anchovies and Breadcrumbs Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories:       "642 kcal",
					Carbohydrates:  "90 g",
					Cholesterol:    "3 mg",
					Fat:            "24 g",
					Fiber:          "4 g",
					Protein:        "17 g",
					SaturatedFat:   "3 g",
					Sodium:         "518 mg",
					Sugar:          "5 g",
					UnsaturatedFat: "0 g",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.thespruceeats.com/pasta-with-anchovies-and-breadcrumbs-recipe-5215384",
			},
		},
		{
			name: "thevintagemixer.com",
			in:   "https://www.thevintagemixer.com/roasted-asparagus-grilled-cheese/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Course"},
				CookTime:      "PT5M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2017-04-03T04:00:45+00:00",
				Description: models.Description{
					Value: "A seasonally fresh take on grilled cheese with asparagus.",
				},
				Image: models.Image{
					Value: "http://d6h7vs5ykbiug.cloudfront.net/wp-content/uploads/2017/04/asparagus-grilled-cheese-9.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"6 spears of asparagus",
						"1 teaspoon olive oil",
						"sea salt and freshly ground pepper",
						"2 ounces of white cheddar cheese*",
						"2 ounces of gruyere cheese*",
						"4 slices of sourdough bread",
						"grass fed butter",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 400 degrees and line a baking sheet with foil.",
						"Remove the woody ends of the asparagus and toss the spears with olive oil, then sprinkle with the sea salt " +
							"and pepper. Spread the asparagus spears out onto the prepared baking sheet and roast for 8-10 " +
							"minutes or until slightly brown.",
						"Meanwhile, slice the cheese and butter the bread. Remove asparagus from oven.",
						"Heat up a large skillet over medium heat. Add two slices of the bread, butter side down, to the pan. " +
							"Add cheese then add 3 spears of asparagus to each sandwich. Top with the other bread slices, butter side up.",
						"Toast for 3-4 minutes then flip and toast for 2 minutes. Remove from heat and serve hot with tomato soup!",
					},
				},
				Name:     "Asparagus Grilled Cheese",
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://www.thevintagemixer.com/roasted-asparagus-grilled-cheese/",
			},
		},
		{
			name: "thewoksoflife.com",
			in:   "https://thewoksoflife.com/fried-wontons/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizers and Snacks"},
				CookTime:      "PT20M",
				Cuisine:       models.Cuisine{Value: "Chinese"},
				DatePublished: "2015-09-05T13:38:44+00:00",
				Description: models.Description{
					Value: "Fried wontons are a easy-to-make crispy, crunchy, delicious appetizer. Your guests will be talking " +
						"about these fried wontons long after the party's over!",
				},
				Keywords: models.Keywords{Values: "fried wontons"},
				Image: models.Image{
					Value: "https://thewoksoflife.com/wp-content/uploads/2015/09/fried-wontons-6-1.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"12 oz. ground pork ((340g))",
						"2 tablespoons finely chopped scallions",
						"1 teaspoon sesame oil",
						"1 tablespoon soy sauce",
						"1 tablespoon shaoxing wine ((or dry sherry))",
						"1/2 teaspoon sugar",
						"1 tablespoon peanut or canola oil",
						"2 tablespoons water ((plus more for sealing the wontons))",
						"1/8 teaspoon white pepper",
						"40-50 Wonton skins ((1 pack, medium thickness))",
						"2 tablespoons apricot jam",
						"1/2 teaspoon soy sauce",
						"1/2 teaspoon rice wine vinegar",
						"2 tablespoons honey",
						"2 tablespoons Sriracha",
						"1 ½ tablespoons light soy sauce",
						"1 tablespoon sugar ((dissolved in 1 tablespoon hot water))",
						"1 teaspoon Worcestershire sauce",
						"1/2 teaspoon toasted sesame seeds ((optional))",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Start by making the filling. Simply combine the ground pork, chopped scallions, sesame oil, soy sauce, " +
							"wine (or sherry), sugar, oil, water, and white pepper in a bowl. Whip everything together by hand " +
							"for 5 minutes or in a food processor for 1 minute. You want the pork to look a bit like a paste.",
						"To make the wontons, take a wrapper, and add about a teaspoon of filling. Overstuffed wontons will pop " +
							"open during the cooking process and make a mess. Use your finger to coat the edges with a little " +
							"water (this helps the two sides seal together).",
						"For shape #1:",
						"Fold the wrapper in half into a rectangle, and press the two sides together so you get a firm seal. Hold " +
							"the bottom two corners of the little rectangle you just made, and bring the two corners together, " +
							"pressing firmly to seal. (Use a little water to make sure it sticks.)",
						"Shape #2:",
						"Fold the wonton in half so you have a triangle shape. Bring together the two outer corners, and press to " +
							"seal (you can use a little water to make sure it sticks).",
						"Keep assembling until all the filling is gone (this recipe should make between 40 and 50 wontons). Place " +
							"the wontons on a baking sheet or plate lined with parchment paper to prevent sticking.",
						"At this point, you can cover the wontons with plastic wrap, put the baking sheet/plate into the freezer, " +
							"and transfer them to Ziploc bags once they’re frozen. They’ll keep for a couple months in the " +
							"freezer and be ready for the fryer whenever you’re ready.",
						"To conserve oil, use a small pot to fry the wontons. Fill it with 2 to 3 inches of oil, making sure the " +
							"pot is deep enough so the oil does not overflow when adding the wontons. Heat the oil to 350 degrees, " +
							"and fry in small batches, turning the wontons occasionally until they are golden brown.",
						"If you have a small spider strainer or slotted spoon, you can use it to keep the wontons submerged when " +
							"frying. This method will give you the most uniform golden brown look without the fuss of turning them. " +
							"Remove the fried wontons to a sheet pan lined with paper towels or a metal cooling rack to drain.",
						"To make one or all of the sauces, simply mix the respective ingredients in a small bowl, and you’re ready " +
							"to eat!",
					},
				},
				Name: "Fried Wontons",
				NutritionSchema: models.NutritionSchema{
					Calories:      "164 kcal",
					Carbohydrates: "15 g",
					Cholesterol:   "23 mg",
					Fat:           "8 g",
					Fiber:         "1 g",
					Protein:       "7 g",
					SaturatedFat:  "2 g",
					Servings:      "1",
					Sodium:        "243 mg",
					Sugar:         "1 g",
				},
				PrepTime: "PT90M",
				Yield:    models.Yield{Value: 12},
				URL:      "https://thewoksoflife.com/fried-wontons/",
			},
		},
		{
			name: "thinlicious.com",
			in:   "https://thinlicious.com/low-carb-greek-yogurt/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Breakfast"},
				CookTime:      "PT720M",
				Cuisine:       models.Cuisine{Value: "Egg free"},
				DatePublished: "2023-07-19T07:00:00+00:00",
				Description: models.Description{
					Value: "This low-carb yogurt recipe is deliciously thick and creamy. Perfect for breakfast, snacks, or even cooking and baking.",
				},
				Keywords: models.Keywords{Values: "low-carb yogurt"},
				Image: models.Image{
					Value: "https://thinlicious.com/wp-content/uploads/2023/04/Low-Carb-Yogurt-Featured-Image-Template-1200x1200-1.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"½ gallon ultra-pasteurized whole milk",
						"2 tbsp. plain Greek yogurt (with live stains)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In the Instant Pot, whisk together a cup of milk and 2 tablespoons of Greek yogurt until well combined. Then whisk in the rest of the milk. This makes it easier to thoroughly combine the yogurt and milk.",
						"Close the Instant Pot lid and press the \"Yogurt\" button, then press \"+\" to adjust the time. The incubation period for Greek yogurt is between 10-12 hours. Choose 10 hours for a mildly tangy yogurt. Increase the time for a tangier flavored yogurt.The Instant Pot will beep and then the timer will start counting up. Do not remove the lid or stir during the incubation period.",
						"After incubation, remove the lid check to make sure your yogurt has set, and taste your yogurt. If it is not tangy enough incubate it for another hour. Your yogurt will be thick, but not yet Greek yogurt consistency.",
						"Now, it is time to cool and strain the yogurt. To strain place a colander into a mixing bowl and line it with a cheesecloth or coffee filters. Pour the yogurt into the colander and place the whole thing into the refrigerator. Let the yogurt cool and strain for 2-4 hours. The longer it strains the thicker the yogurt will be.",
						"When the yogurt reaches your desired thickness store the yogurt in an airtight container in the refrigerator for up to 2 weeks. Reserve a portion of the yogurt to use as a starter for your next batch of yogurt.",
					},
				},
				Name: "Low-Carb Yogurt",
				NutritionSchema: models.NutritionSchema{
					Calories:      "70 kcal",
					Carbohydrates: "4 g",
					Fat:           "4 g",
					Protein:       "5 g",
					Servings:      "12",
					Sodium:        "55 mg",
					Sugar:         "4 g",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 12},
				URL:      "https://thinlicious.com/low-carb-greek-yogurt/",
			},
		},
		{
			name: "tidymom.net",
			in:   "https://tidymom.net/make-ahead-mashed-potatoes/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Side Dish Recipes"},
				CookTime:      "PT45M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2021-09-28",
				Description: models.Description{
					Value: "MAKE-AHEAD MASHED POTATOES are perfect for busy holidays. They will keep up to two days in the refrigerator and come out creamy and buttery, with a slightly crunchy golden crust on top every time!",
				},
				Keywords: models.Keywords{Values: "potatoes, side dish, mashed"},
				Image: models.Image{
					Value: "https://tidymom.net/blog/wp-content/uploads/2021/11/make-ahead-mashed-potatoes-plated-image-480x480.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2½ pounds potatoes, peeled and cubed (Yukon Gold, Russet or combo of the two)",
						"3 ounces cream cheese, room temperature",
						"1¼ cup milk (or half heavy cream and half milk), warmed", "1 teaspoon salt",
						"¼ teaspoon pepper", "1/2 cup unsalted butter (divided)",
						"freshly chopped chives or thyme for garnish (optional)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Spray a large casserole dish with cooking spray (or grease with butter) and set aside.",
						"Boil potatoes until tender but still firm (about 15 minutes). Drain and drain the potatoes and put them back in the pot without a lid. Set them back on a hot/warm burner for a few minutes, shaking the pot, until all the excess moisture on the potatoes has evaporated.",
						"In a large mixing bowl, mix potatoes, cream cheese, milk, 1/4 cup (4 tablespoons) softened butter, and salt. Mix until blended.",
						"Add mashed potatoes to the prepared casserole dish and cover. Store in the refrigerator overnight or for up to two days.",
						"When ready to bake, let the potatoes sit at room temperature for about 1/2 hour and heat the oven to 350° F. Dot the top of the potatoes with the remaining 1/4 cup (4 tablespoons) of butter. Bake covered in the preheated oven for 30 minutes, uncover and bake until hot about another 15-20 minutes. This will give the mashed potatoes a bit of a crispy top.",
						"Garnish with freshly chopped chives or thyme, if desired.",
					},
				},
				Name: "Make-Ahead Mashed Potatoes",
				NutritionSchema: models.NutritionSchema{
					Calories:       "193 calories",
					Carbohydrates:  "22 grams carbohydrates",
					Cholesterol:    "30 milligrams cholesterol",
					Fat:            "11 grams fat",
					Fiber:          "2 grams fiber",
					Protein:        "4 grams protein",
					SaturatedFat:   "7 grams saturated fat",
					Servings:       "1",
					Sodium:         "221 milligrams sodium",
					Sugar:          "3 grams sugar",
					TransFat:       "0 grams trans fat",
					UnsaturatedFat: "3 grams unsaturated fat",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 10},
				URL:      "https://tidymom.net/make-ahead-mashed-potatoes/",
			},
		},
		{
			name: "timesofindia.com",
			in:   "https://recipes.timesofindia.com/recipes/beetroot-cold-soup/rs90713582.cms",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizers"},
				Cuisine:       models.Cuisine{Value: "Vegetarian"},
				DatePublished: "2022-04-07T22:43:02+05:30",
				Description: models.Description{
					Value: "Yearning for a satiating and delicious meal, then try this easy and tasty cold soup made with " +
						"beetroot, curd, hard boiled eggs, coriander leaves and spices. To make this simple soup, just " +
						"follow us through some simple steps and make a sumptuous and enjoy it cold.",
				},
				Keywords: models.Keywords{
					Values: "Beetroot Cold Soup recipe, Vegetarian, cook Beetroot Cold Soup",
				},
				Image: models.Image{
					Value: "https://static.toiimg.com/thumb/90713582.cms?width=1200&height=900",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 Numbers egg",
						"1 cup yoghurt (curd)",
						"1 handfuls coriander leaves",
						"0 As required salt",
						"0 As required black pepper",
						"1/2 teaspoon cumin powder",
						"1 teaspoon oregano",
						"0 As required water",
						"1 Numbers beetroot",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"To make this easy recipe, take a pan and add water, boil the egg with a dash of salt. Once done, remove " +
							"the shell and save for garnishing. In the meantime, wash and chop the beetroot and make a smooth blend.",
						"Next, whisk the curd and add beetroot blend along with spices, chopped coriander leaves and mix it well. " +
							"Garnish with oregano and egg and enjoy.",
					},
				},
				Name: "Beetroot Cold Soup Recipe",
				NutritionSchema: models.NutritionSchema{
					Calories: "189 cal",
					Servings: "1",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 2},
				URL:      "https://recipes.timesofindia.com/recipes/beetroot-cold-soup/rs90713582.cms",
			},
		},
		{
			name: "tine.no",
			in:   "https://www.tine.no/oppskrifter/middag-og-hovedretter/kylling-og-fjarkre/rask-kylling-tikka-masala",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "middag"},
				Cuisine:       models.Cuisine{Value: "indisk"},
				DateModified:  "2022-12-19T09:14:12.161Z",
				DatePublished: "2018-09-05T18:43:37.088Z",
				Description: models.Description{
					Value: "En god og rask oppskrift på en kylling tikka masala. Dette er en rett med små smakseksplosjoner som " +
						"sender tankene til India.",
				},
				Keywords: models.Keywords{Values: "kylling"},
				Image: models.Image{
					Value: "https://www.tine.no/_/recipeimage/w_1080,h_1080,c_fill,x_2880,y_1920,g_xy_center/recipeimage/w1r3ydbmyeqcngqpxatv.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"4 dl basmatiris",
						"400 g kyllingfileter",
						"1 ss TINE® Meierismør til steking",
						"1 stk paprika",
						"0.5 dl chili",
						"3 stk vårløk",
						"1 ts hvitløksfedd",
						"1 ss hakket, frisk ingefær",
						"0.5 dl hakket, frisk koriander",
						"2 ts garam masala",
						"3 dl TINE® Lett Crème Fraîche 18 %",
						"3 ss tomatpuré",
						"0.5 ts salt",
						"0.25 ts pepper",
						"0.5 dl slangeagurk",
						"3 dl TINE® Yoghurt Naturell Gresk Type",
						"0.5 dl frisk mynte",
						"1 ts hvitløksfedd",
						"0.5 ts salt",
						"0.25 ts pepper",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Kok ris etter anvisningen på pakken.",
						"Del kylling i biter. Brun kyllingen i smør i en stekepanne på middels varme.",
						"Rens og hakk paprika, chili, vårløk og hvitløk og ha det i stekepannen sammen med kyllingen. Rens og " +
							"finhakk ingefær og frisk koriander. Krydre med garam masala, koriander og ingefær.",
						"Hell i crème fraîche og tomatpuré, og la småkoke i 5 minutter. Smak til med salt og pepper.",
						"Riv agurk og bland den med yoghurt. Hakk mynte og hvitløk og bland det i. Smak til med salt og pepper.",
					},
				},
				Name:  "Rask kylling tikka masala",
				Yield: models.Yield{Value: 4},
				URL:   "https://www.tine.no/oppskrifter/middag-og-hovedretter/kylling-og-fjarkre/rask-kylling-tikka-masala",
			},
		},
		{
			name: "tudogostoso.com.br",
			in:   "https://www.tudogostoso.com.br/receita/585-rocambole-de-carne-moida.html",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Carnes"},
				DatePublished: "2005-08-06",
				Description: models.Description{
					Value: "Descubra a receita de Rocambole de carne moída para fazer em 30 minutos. Em uma vasilha junte à carne, o pacote de creme de cebola e amasse bem até ficar homogêneo. Em seguida sobre um plástico abra a massa (carne moida com o creme de cebola) no tamanho do tabuleiro que irá assar. Sobre a massa coloque a cebola e o alho picados bem pequenos, a muss…",
				},
				Keywords: models.Keywords{Values: "Receita de Rocambole de carne moída"},
				Image: models.Image{
					Value: "https://static.itdg.com.br/images/1200-675/929a16753e9c9d22f04269c75ec17da7/353783-original.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 kg de carne moída", "1 pacote de creme de cebola", "1 cebola",
						"150 g de mussarela", "150 g de presunto", "3 fatias de bacon", "Alho a gosto",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Em uma vasilha junte à carne, o pacote de creme de cebola e amasse bem até ficar homogêneo.",
						"Em seguida sobre um plástico abra a massa (carne moida com o creme de cebola) no tamanho do tabuleiro que irá assar.",
						"Sobre a massa coloque a cebola e o alho picados bem pequenos, a mussarela em seguida o presunto.",
						"Para enrolar a massa basta começar a levantar a ponta do plástico.",
						"Sobre o rocambole já enrolado coloque as fatias de bacon para decorar.",
					},
				},
				Name:     "Rocambole de carne moída",
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 8},
				URL:      "https://www.tudogostoso.com.br/receita/585-rocambole-de-carne-moida.html",
			},
		},
		{
			name: "twopeasandtheirpod.com",
			in:   "https://www.twopeasandtheirpod.com/easy-chickpea-salad/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Salad"},
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2022-03-28T06:03:00+00:00",
				Description: models.Description{
					Value: "This chickpea salad is a vegetarian version of a classic chicken salad, some refer to it as chickpea " +
						"chicken salad. It&#039;s made with basic ingredients, loaded with flavor, and perfect for picnics, " +
						"work lunches, or a simple, healthy lunch at home.",
				},
				Image: models.Image{
					Value: "https://www.twopeasandtheirpod.com/wp-content/uploads/2022/02/Chickpea-Salad-4.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 (15 oz) cans chickpeas, (drained and rinsed)",
						"3/4 cup diced celery",
						"1/2 cup diced dill pickles",
						"1/2 cup sliced green onion",
						"1/2 cup plain Greek yogurt",
						"1 tablespoon lemon juice",
						"1 to 2 tablespoons Dijon mustard",
						"2 teaspoons red wine vinegar",
						"2 tablespoons freshly chopped dill",
						"2 tablespoons freshly chopped parsley",
						"1/4 teaspoon garlic powder",
						"Kosher salt and black pepper, (to taste)",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Place the chickpeas on a clean dish towel or paper towel. Put another towel on top. Use your hands to " +
							"roll and rub the chickpeas for about 20-30 seconds. This will help the skins come off easier. " +
							"Remove the skins and discard. I try to remove most of the skins, but if you don’t get them all, " +
							"that is ok. And if you don&#39;t have time to remove the skins, just rinse and drain. The salad " +
							"will still be good. Removing the skins just makes the salad a little creamier and smooth.",
						"Place the chickpeas in a large bowl and mash with a fork or potato masher until most of the chickpeas " +
							"are smashed. Stir in the celery, onion, and pickles.",
						"In a small bowl, whisk together the Greek yogurt, lemon juice, mustard, red wine vinegar, dill, parsley, " +
							"garlic powder, salt, and pepper.",
						"Add the sauce to the chickpea mixture and stir until well combined. Taste and adjust ingredients, if " +
							"necessary.",
						"Serve in between two slices of bread to make a sandwich, in pita bread, in a lettuce wrap, in a tortilla, " +
							"with crackers or chips, or on top of a rice cake, or add to a bed of greens to make a salad! The " +
							"options are endless.",
					},
				},
				Keywords: models.Keywords{Values: "vegetarian"},
				Name:     "Chickpea Salad",
				NutritionSchema: models.NutritionSchema{
					Calories:       "215 kcal",
					Carbohydrates:  "32 g",
					Cholesterol:    "1 mg",
					Fat:            "5 g",
					Fiber:          "10 g",
					Protein:        "14 g",
					SaturatedFat:   "1 g",
					Servings:       "1",
					Sodium:         "803 mg",
					Sugar:          "2 g",
					TransFat:       "1 g",
					UnsaturatedFat: "3 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.twopeasandtheirpod.com/easy-chickpea-salad/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
