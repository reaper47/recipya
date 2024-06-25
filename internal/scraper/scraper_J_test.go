package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
	"time"
)

func TestScraper_J(t *testing.T) {
	testcases := []testcase{
		{
			name: "jamieoliver.com",
			in:   "https://www.jamieoliver.com/recipes/chicken-recipes/thai-green-chicken-curry/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Mains"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{Value: "asian"},
				DatePublished: "2015-09-16",
				Description: &models.Description{
					Value: "This deliciously fragrant Thai green curry really packs a flavour punch.",
				},
				Keywords: &models.Keywords{
					Values: "chicken, mushroom, dairy-free, poultry, vegetable, thai green, curry, chicken thighs, paste, chicken curry, thai, thai green curry, vegetables, One-pan recipes, Curry, Chicken, Stewing, Dinner Party",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"750 g skinless free-range chicken thighs",
						"groundnut oil",
						"400 g mixed oriental mushrooms",
						"1 x 400g tin of light coconut milk",
						"1 organic chicken stock cube",
						"6 lime leaves",
						"200 g mangetout",
						"½ a bunch fresh Thai basil",
						"2 limes",
						"4 cloves of garlic",
						"2 shallots",
						"5cm piece of ginger",
						"2 lemongrass stalks",
						"4 green Bird's eye chillies",
						"1 teaspoon ground cumin",
						"½ a bunch of fresh coriander",
						"2 tablespoons fish sauce",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "<ol class=\"recipeSteps\"><li>To make the curry paste, peel, roughly chop and place the garlic, shallots and ginger into a food processor. </li><li>Trim the lemongrass, remove the tough outer leaves, then finely chop and add to the processor. Trim and add the chillies along with the cumin and half the coriander (stalks and all). Blitz until finely chopped, add the fish sauce and blitz again. </li><li>Slice the chicken into 2.5cm strips. Heat 1 tablespoon of oil in a large pan on a medium heat, add the chicken and fry for 5 to 7 minutes, or until just turning golden, then transfer to a plate. </li><li>Tear the mushrooms into even pieces. Return the pan to a medium heat, add the mushrooms and fry for 4 to 5 minutes, or until golden. Transfer to a plate using a slotted spoon. </li><li>Reduce the heat to medium-low and add the Thai green paste for 4 to 5 minutes, stirring occasionally. </li><li>Pour in the coconut milk and 200ml of boiling water, crumble in the stock cube and add the lime leaves. Turn the heat up and bring gently to the boil, then simmer for 10 minutes, or until reduced slightly.</li><li>Stir in the chicken and mushrooms, reduce the heat to low and cook for a further 5 minutes, or until the chicken is cooked through, adding the mangetout for the final 2 minutes. </li><li>Season carefully to taste with sea salt and freshly ground black pepper. Pick, roughly chop and stir through the basil leaves and remaining coriander leaves. Serve with lime wedges and steamed rice.</li></ol>"},
					},
				},
				Name: "Thai green chicken curry",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "285 calories",
					Carbohydrates: "6.1 g carbohydrate",
					Fat:           "16.2 g fat",
					Fiber:         "2.2 g fibre",
					Protein:       "28.9 g protein",
					SaturatedFat:  "6.5 g saturated fat",
					Sodium:        "1.0 g salt",
					Sugar:         "4.2 g sugar",
				},
				Tools:     &models.Tools{Values: []models.HowToItem{}},
				TotalTime: "PT50M",
				URL:       "https://www.jamieoliver.com/recipes/chicken-recipes/thai-green-chicken-curry/",
				Video: &models.Videos{
					Values: []models.VideoObject{
						{
							AtType:      "VideoObject",
							ContentUrl:  "https://www.youtube.com/watch?v=wno_qUB02lM",
							Description: "This deliciously fragrant Thai green curry really packs a flavour punch.",
							EmbedUrl:    "https://www.youtube.com/embded/wno_qUB02lM",
							Name:        "Thai green chicken curry",
							ThumbnailURL: &models.ThumbnailURL{
								Value: "https://cdn.jamieoliver.com/recipe-database/oldImages/xtra_med/1575_2_1437576282.jpg",
							},
						},
					},
				},
				Yield: &models.Yield{Value: 6},
			},
		},
		{
			name: "jaimyskitchen.nl",
			in:   "https://jaimyskitchen.nl/recepten/sajoer-lodeh-indonesisch-groente-recept-met-tofu",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Hoofdgerechten"},
				CookTime:      "",
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{Value: "Aziatisch"},
				DatePublished: "2023-07-30 07:00:02",
				Description:   &models.Description{Value: "Sajour Lodeh is een heerlijk groentegerecht met kokosmelk afkomstig uit Indonesië. Een heerlijk vegetarisch (en veganistisch!) gerecht voor bij een rijsttaf..."},
				Keywords:      &models.Keywords{Values: "Hoofdgerechten, Indonesisch, Aziatisch, Vegetarisch, Vegan, Tofu &amp; tempeh, tofu, spitskool, wortel, sperziebonen, bloemkool, kokosmelk"},
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"200 gr wortel, in plakjes", "200 gr sperziebonen, in stukjes",
						"200 gr spitskool, gesneden", "200 gr bloemkool, in roosjes", "375 gr tofu",
						"2 el limoensap", "400 ml kokosmelk", "1 groentenbouillonblokje",
						"zonnebloemolie", "title Boemboe:", "1 ui", "2 tenen knoflook",
						"2 tl sambal badjak", "1 tl komijnpoeder", "1.5 tl korianderpoeder",
						"1 tl laos poeder",
						"1 tl tamarinde pasta",
						"1 stengel citroengras",
						"1 tl palmsuiker",
						"title Keukenhulpjes:",
						"Hapjespan",
						"Steelpan",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Laat de tofu uitlekken en snijd vervolgens in wat dunnere plakken. Snijd de plakken vervolgens in kleine blokjes of driehoekjes. Verhit zonnebloemolie in een hapjespan en bak de tofu in ongeveer 5 minuten goudbruin. Haal uit de pan en laat de olie in de pan."},
						{Type: "HowToStep", Text: "Kook de bloemkool en sperziebonen gaar in een steelpan in ongeveer 8 minuten. Giet af en bewaar de groenten."},
						{Type: "HowToStep", Text: "Maak in de tussentijd de boemboe door de ingrediënten heel fijn te snijden (ui, knoflook) en samen te voegen. Kneus de stengel citroengras. Fruit de boemboe inclusief de stengel citroengras aan in een steelpan met een beetje zonnebloemolie. Voeg na ongeveer 3 minuten kokosmelk en het bouillonblokje toe en laat op laag vuur nog 5 minuten sudderen."},
						{Type: "HowToStep", Text: "Voeg de wortel in de hapjespan met zonnebloemolie. Roerbak ongeveer 5 minuten en voeg dan ook de spitskool toe. Na ongeveer 5 minuten is de spitskool een beetje geslonken, voeg dan ook de bloemkool, sperziebonen en de boemboe met kokosmelk (verwijder de stengel citroengras) toe. Maak af met limoensap naar smaak."},
					},
				},
				Name:            "Sajoer Lodeh - Indonesisch groente recept met tofu",
				NutritionSchema: &models.NutritionSchema{},
				PrepTime:        "PT40M",
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				TotalTime:       "PT40M",
				Yield:           &models.Yield{Value: 4},
				URL:             "https://jaimyskitchen.nl/recepten/sajoer-lodeh-indonesisch-groente-recept-met-tofu",
			},
		},
		{
			name: "jaroflemons.com",
			in:   "https://www.jaroflemons.com/vegetarian-hot-honey-pizza/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Dinner"},
				CookTime:      "PT45M",
				Cuisine:       &models.Cuisine{Value: "American"},
				DatePublished: "2022-01-31T03:00:00+00:00",
				Description:   &models.Description{Value: "This Vegetarian Hot Honey Pizza is what dreams are made of! An easy-to-make (or store bought) pizza crust topped with roasted mushrooms and butternut squash, parmesan, feta, red onions, and arugula. Drizzled with hot honey and then baked to perfection for a delicious balance of savory, spicy, and a touch of sweet!"},
				Keywords:      &models.Keywords{Values: "Vegetarian Hot Honey Pizza"},
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1 cup butternut squash (diced into small cubes)",
						"1/4 cup olive oil (divided)",
						"8 ounces sliced mushrooms",
						"1 batch pizza dough",
						"2 teaspoons minced garlic",
						"4 cups arugula",
						"1 cup grated parmesan",
						"1 cup red onion (sliced)",
						"salt/pepper (to taste)",
						"3/4 cup feta cheese (crumbled)",
						"2 Tablespoons hot honey",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Preheat the oven to 425 degrees F."},
						{Type: "HowToStep", Text: "Dice the butternut squash and place on a baking sheet or roasting pan with 1/4 cup of water (for a more tender texture)."},
						{Type: "HowToStep", Text: "Drizzle with 2 Tablespoons of oil and add salt and pepper to taste."},
						{Type: "HowToStep", Text: "Roast the squash for 15 minutes."},
						{Type: "HowToStep", Text: "Add the sliced mushrooms to the same pan and roast for another 10 minutes."},
						{Type: "HowToStep", Text: "While the vegetables are roasting, prepare and press down the pizza dough into a pan."},
						{Type: "HowToStep", Text: "Drizzle the remaining oil over the pizza dough."},
						{Type: "HowToStep", Text: "Top with the garlic, arugula, grated parmesan, and sliced red onions."},
						{Type: "HowToStep", Text: "Add the roasted squash, mushrooms, salt/pepper (to taste), and crumbled feta."},
						{Type: "HowToStep", Text: "Drizzle with hot honey."},
						{Type: "HowToStep", Text: "Bake the pizza for about 15-20 minutes."},
						{Type: "HowToStep", Text: "Serve and enjoy!"},
					},
				},
				Name: "Vegetarian Hot Honey Pizza",
				NutritionSchema: &models.NutritionSchema{
					Calories:       "193 kcal",
					Carbohydrates:  "11 g",
					Cholesterol:    "24 mg",
					Fat:            "13 g",
					Fiber:          "1 g",
					Protein:        "8 g",
					SaturatedFat:   "5 g",
					Servings:       "1",
					Sodium:         "355 mg",
					Sugar:          "7 g",
					TransFat:       "",
					UnsaturatedFat: "8 g",
				},
				PrepTime:  "PT10M",
				TotalTime: "PT55M",
				Yield:     &models.Yield{Value: 8},
				URL:       "https://www.jaroflemons.com/vegetarian-hot-honey-pizza/",
			},
		},
		{
			name: "jimcooksfoodgood.com",
			in:   "https://jimcooksfoodgood.com/recipe-weeknight-pad-thai/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Main Dish"},
				CookingMethod: &models.CookingMethod{},
				CookTime:      "PT15M",
				Cuisine:       &models.Cuisine{Value: "American"},
				DatePublished: "2023-05-09T12:58:13+00:00",
				Description:   &models.Description{Value: "Quick easy and delicious"},
				Keywords:      &models.Keywords{Values: "#healthyrecipe"},
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"8 ounces rice noodles",
						"3 tablespoons tamarind ((or 2 tablespoons more of both lime juice and brown sugar))",
						"1/2 cups soy sauce",
						"4 tablespoons brown sugar",
						"2 tablespoons Sriracha",
						"2 limes ((one for juice, one for wedges))",
						"2 green onions",
						"2 shallots",
						"3 eggs",
						"4 garlic cloves",
						"1 cup bean sprouts",
						"2 cups Chopped Broccoli",
						"1/2 c roasted peanuts ((coarsely chopped))",
						"3 tablespoons cooking oil",
						"Optional: 1 pound of cooked protein ((shrimp, tofu, etc))",
						"Optional: Sriracha Mayo",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Bring a pot of water to boil and cook rice noodles according to package directions, shy just one minute. Drain and set aside."},
						{Type: "HowToStep", Text: "In a separate bowl, combine tamarind, soy sauce, brown sugar, sriracha, and the juice of one lime together."},
						{Type: "HowToStep", Text: "In a very large pan over medium heat, add one tablespoon of the oil. Add the eggs and scramble until just set, and set aside."},
						{Type: "HowToStep", Text: "Slice the shallots and green onions thinly, and mince the garlic. In the same pan, add the remainder of the oil, still over medium heat. Add green onions, shallots, garlic and broccoli. Sauté until broccoli is cooked through, 4-5 minutes."},
						{Type: "HowToStep", Text: "Add the noodles to the pan and pour on the sauce. Toss to coat all noodles. Add the eggs, bean sprout, and your cooked protein. Sprinkle peanuts on top, and serve along with a wedge of lime and Sriracha Mayo."},
					},
				},
				Name: "Pad Thai",
				NutritionSchema: &models.NutritionSchema{
					Calories: "389 kcal",
					Servings: "4",
				},
				PrepTime: "PT15M",
				Tools:    &models.Tools{Values: []models.HowToItem{}},
				Yield:    &models.Yield{Value: 4},
				URL:      "https://jimcooksfoodgood.com/recipe-weeknight-pad-thai/",
			},
		},
		{
			name: "joyfoodsunshine.com",
			in:   "https://joyfoodsunshine.com/peanut-butter-frosting/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "condiment"},
				Cuisine:       &models.Cuisine{Value: "American"},
				DatePublished: "2022-02-16T02:36:00+00:00",
				Description: &models.Description{
					Value: "This peanut butter frosting recipe is easy to make in 5 minutes. It&#039;s silky smooth, made with more peanut butter than butter and is flavored with vanilla &amp; sea salt. It pipes well and tastes delicious on top of chocolate cupcakes and brownies and chocolate cake.",
				},
				Keywords: &models.Keywords{
					Values: "how to make peanut butter frosting, peanut butter buttercream, peanut butter frosting, peanut butter frosting recipe",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"¾ cup creamy peanut butter",
						"½ cup salted butter (softened)",
						"½ teaspoon pure vanilla extract",
						"¼ teaspoon fine sea salt",
						"2 cups powdered sugar",
						"1-2 Tablespoons whole milk (room temperature)",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "In the bowl of a standing mixer fitted with the paddle attachment, or in a large bowl with a handheld mixer, beat together peanut butter and butter until smooth."},
						{Type: "HowToStep", Text: "Add vanilla and sea salt and beat until combined."},
						{Type: "HowToStep", Text: "Add powdered sugar, 1 cup at a time, and beat until fully incorporated after each addition."},
						{Type: "HowToStep", Text: "Add 1 tablespoon whole milk and beat. If necessary, add an additional 1 tablespoon milk to achieve your desired consistency."},
						{Type: "HowToStep", Text: "Use to frost a chocolate cake, chocolate cupcakes, brownies, etc."},
					},
				},
				Name: "Peanut Butter Frosting Recipe",
				NutritionSchema: &models.NutritionSchema{
					Calories:       "181 kcal",
					Carbohydrates:  "17 g",
					Cholesterol:    "15 mg",
					Fat:            "12 g",
					Fiber:          "1 g",
					Protein:        "3 g",
					SaturatedFat:   "5 g",
					Servings:       "2 TBS",
					Sodium:         "143 mg",
					Sugar:          "16 g",
					TransFat:       "1 g",
					UnsaturatedFat: "6 g",
				},
				PrepTime:  "PT5M",
				TotalTime: "PT5M",
				URL:       "https://joyfoodsunshine.com/peanut-butter-frosting/",
				Video: &models.Videos{
					Values: []models.VideoObject{
						{
							AtType:      "VideoObject",
							ContentUrl:  "https://mediavine-res.cloudinary.com/video/upload/t_original/v1642947749/tzdoglhnbiiiup85wmps.mp4",
							Description: "This peanut butter frosting recipe (peanut butter buttercream) is easy to make in 5 minutes. It's silky smooth, made with more peanut butter than butter and is flavored with vanilla & sea salt. It pipes well and tastes delicious on top of chocolate cupcakes and brownies and chocolate cake.",
							Duration:    "PT53S",
							EmbedUrl:    "https://video.mediavine.com/videos/tzdoglhnbiiiup85wmps.js",
							Name:        "Peanut Butter Frosting Recipe",
							ThumbnailURL: &models.ThumbnailURL{
								Value: "https://mediavine-res.cloudinary.com/image/upload/s--CxD4o5LH--/c_limit,f_auto,fl_lossy,h_1080,q_auto,w_1920/v1642947695/haqpzgp58bql3ypnhnrv.jpg",
							},
							UploadDate: time.Date(2022, 1, 23, 14, 22, 36, 0, time.UTC),
						},
					},
				},
				Yield: &models.Yield{Value: 16},
			},
		},
		{
			name: "juliegoodwin.com.au",
			in:   "https://juliegoodwin.com.au/white-chocolate-and-raspberry-muffins/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Desserts"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				DateModified:  "2021-09-29T06:10:04+00:00",
				DatePublished: "2022-03-27T19:05:59+00:00",
				Description: &models.Description{
					Value: "White Chocolate and Raspberry Muffins\u00a0| This is a very indulgent recipe, absolutely delicious, dense and moist. What I would call a “sometimes food” for sure – but every now and again, with a very good coffee, lovely.",
				},
				Keywords: &models.Keywords{},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"¾ cup white choc bits", "¾ cup caster sugar", "2 cups self-raising flour",
						"1 ½ cups raspberries (frozen can be used, but fresh in season…wow.)",
						"2 eggs", "½ cup vegetable oil", "½ cup milk",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Preheat oven to 180°C. Line 6 large muffin pans (180ml) with paper cases."},
						{Type: "HowToStep", Text: "Combine dry ingredients in a large bowl. Chop about ¾ cup raspberries and mix whole and chopped raspberries with the dry ingredients."},
						{Type: "HowToStep", Text: "Whisk egg, oil and milk together. Pour into a well in the centre of the dry ingredients."},
						{Type: "HowToStep", Text: "Using a wooden spoon or spatula, gently stir the wet ingredients into the dry ingredients until just combined. Too much mixing at this stage will result in tough, chewy muffins. However you do need to ensure there are no lumps"},
						{Type: "HowToStep", Text: "Spoon the mixture among the muffin pans and bake for 20 minutes or until golden on top and springy to touch. 6. Turn out of the muffin pan and serve warm to grateful recipients."},
					},
				},
				Name:            "White Chocolate and Raspberry Muffins",
				NutritionSchema: &models.NutritionSchema{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				Yield:           &models.Yield{Value: 6},
				URL:             "https://juliegoodwin.com.au/white-chocolate-and-raspberry-muffins/",
			},
		},
		{
			name: "justataste.com",
			in:   "https://www.justataste.com/mini-sour-cream-doughnut-muffins-recipe/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Breakfast"},
				CookTime:      "PT16M",
				Cuisine:       &models.Cuisine{Value: "American"},
				DatePublished: "2022-03-10T09:59:00+00:00",
				Description: &models.Description{
					Value: "Two breakfast favorites join forces in a family-friendly recipe for Mini Sour Cream Doughnut Muffins rolled in cinnamon-sugar.",
				},
				Keywords: &models.Keywords{
					Values: "cinnamon, doughnut, sour cream, sugar, vanilla extract",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"Cooking spray",
						"1 cup all-purpose flour",
						"1 teaspoon baking powder",
						"1/4 teaspoon baking soda",
						"1/4 teaspoon salt",
						"3 Tablespoons unsalted butter, at room temp",
						"3 Tablespoons vegetable oil",
						"1 cup sugar, divided",
						"1 large egg",
						"1/2 cup sour cream",
						"1 teaspoon vanilla extract",
						"1 1/2 teaspoons ground cinnamon",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Preheat the oven to 350°F. Grease a nonstick mini muffin pan with cooking spray."},
						{Type: "HowToStep", Text: "In a medium bowl, whisk together the flour, baking powder, baking soda and salt. Set the mixture aside."},
						{Type: "HowToStep", Text: "In the bowl of a stand mixer fitted with the paddle attachment, beat together the butter, vegetable oil and 1/2 cup sugar until well combined, about 2 minutes. Add the egg and beat until combined."},
						{Type: "HowToStep", Text: "Add the flour mixture, sour cream and vanilla extract and beat just until combined."},
						{Type: "HowToStep", Text: "Using a small ice cream scoop (or two spoons), scoop out heaping 1-tablespoon portions of the batter into the prepared muffin pan."},
						{Type: "HowToStep", Text: "Bake the muffins for 16 to 22 minutes until pale golden. While the muffins bake, in a medium bowl, whisk together the remaining 1/2 cup sugar and cinnamon."},
						{Type: "HowToStep", Text: "Remove the muffins from the oven and let them cool for 2 minutes in the pan before transferring them in batches into the cinnamon-sugar mixture, tossing to coat. Repeat the coating process with the remaining muffins then serve."},
						{Type: "HowToStep", Text: "It’s important to toss the muffins in the cinnamon-sugar mixture while they are hot to ensure the cinnamon-sugar will stick."},
						{Type: "HowToStep", Text: "★Did you make this recipe? Don&#39;t forget to give it a star rating below!"},
					},
				},
				Name: "Mini Sour Cream Doughnut Muffins",
				NutritionSchema: &models.NutritionSchema{
					Calories:       "122 kcal",
					Carbohydrates:  "17 g",
					Cholesterol:    "17 mg",
					Fat:            "6 g",
					Fiber:          "1 g",
					Protein:        "1 g",
					SaturatedFat:   "2 g",
					Servings:       "1",
					Sodium:         "57 mg",
					Sugar:          "11 g",
					TransFat:       "1 g",
					UnsaturatedFat: "3 g",
				},
				PrepTime:  "PT10M",
				TotalTime: "PT26M",
				Yield:     &models.Yield{Value: 18},
				URL:       "https://www.justataste.com/mini-sour-cream-doughnut-muffins-recipe/",
			},
		},
		{
			name: "justbento.com",
			in:   "https://justbento.com/handbook/recipe-collection-mains/sushi-roll-bento-make-sushi-rolls-without-sushi-mat",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Mains"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{Value: "Japanese"},
				DateModified:  "2019-06-11T06:20:29+09:00",
				Description: &models.Description{
					Value: "Here is something that I had in my archives - a sushi roll bento, made with ingredients that you might not have thought belong in a sushi. Plus, how to make a fat sushi roll without a sushi mat!",
				},
				Keywords: &models.Keywords{},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"A sheet of nori seaweed",
						"Properly cooked and prepared sushi rice. It should still be slightly warm, not stone cold, for maximum stick-together-ness.",
						"Fillings of your choice", "A clean, non-fuzzy kitchen towel",
						"A sharp kitchen knife",
						"A bowl filled with water with a little vinegar in it",
						"Impeccably clean hands",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Moisten the kitchen towel and then wring it out tightly. It should just be moist, not dripping."},
						{Type: "HowToStep", Text: "Lay down the kitchen towel flat. Put the nori sheet shiny side down on the towel; the long edge should be on the edge of the towel. (You put the nori sheet shiny side down because that side has a slightly less tendency to split, and also for aesthetic reasons.)"},
						{Type: "HowToStep", Text: "Moisten your fingers with the vinegar-water. Place a fairly thin, even layer of sushi rice on the nori seaweed, up to about an inch (2 cm) away from the far edge of the nori."},
						{Type: "HowToStep", Text: "Place the fillings in the middle of the rice, starting with any flat ingredients like lettuce or shiso leaves, then following up with the other things like julienned vegetables."},
						{Type: "HowToStep", Text: "This roll has lettuce, cucumber, carrots, ham and cheese."},
						{Type: "HowToStep", Text: "Now you're ready to roll! Re-moisten your fingers with the vinegar water. Grab the edge of the towel with the nori and roll it over the filling, holding in the filling with your fingertips. Be brave here - quick and decisive movement will have better results than hesitation."},
						{Type: "HowToStep", Text: "Roll the nori and rice over the filling, as you pull on the edge of the kitchen towel on the other side. If you compare it to the sushi mat method, you'll notice that the method is basically the same."},
						{Type: "HowToStep", Text: "Keep rolling and pulling on the towel, evenly over the length of the roll."},
						{Type: "HowToStep", Text: "Once the roll is completely rolled, apply gentle but firm, even pressure over the whole thing."},
						{Type: "HowToStep", Text: "Here's how the roll looks when it's completed, before the towel is removed."},
						{Type: "HowToStep", Text: "Now you're ready to cut. If there's any rice sticking to your fingers, rinse them off. Moisten the knife and your fingers with the vinegar water."},
						{Type: "HowToStep", Text: "Cut the roll into even pieces. If the knife gets sticky, just re-moisten it with the vinegar water."},
						{Type: "HowToStep", Text: "And that's it! It's not as hard as you might have thought, is it?"},
						{Type: "HowToStep", Text: "Here's the end of the roll that I actually used in the bento above. You can tuck in the raggedy end bits good side up in a bento, or just pop them in your mouth as you make them!"},
						{Type: "HowToStep", Text: "I hope this will inspire you to come up with your own fat sushi roll combinations. Not only are they great for individual bentos, they're a nice change-of-pace carb for a barbeque a picnic too."},
					},
				},
				Name:            "A sushi roll bento, plus how to make sushi rolls without a sushi mat",
				NutritionSchema: &models.NutritionSchema{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				Yield:           &models.Yield{Value: 10},
				URL:             "https://justbento.com/handbook/recipe-collection-mains/sushi-roll-bento-make-sushi-rolls-without-sushi-mat",
			},
		},
		{
			name: "justonecookbook.com",
			in:   "https://www.justonecookbook.com/teriyaki-tofu-bowl/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Main Course"},
				CookTime:      "PT25M",
				Cuisine:       &models.Cuisine{Value: "Japanese"},
				DatePublished: "2022-03-21T05:00:00+00:00",
				Description: &models.Description{
					Value: "Smothered with sweet-savory homemade teriyaki sauce, this crispy Pan-Fried Teriyaki Tofu Bowl is amazingly easy and delicious!  It‘s also a great way to incorporate tofu into your weekly menu rotation.",
				},
				Keywords: &models.Keywords{Values: "teriyaki sauce, tofu"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"14 oz medium-firm tofu (momen dofu) ((1 block))",
						"⅓ cup potato starch or cornstarch",
						"3 Tbsp neutral oil ((divided))",
						"¼ cup sake",
						"¼ cup mirin",
						"¼ cup soy sauce",
						"4 tsp sugar",
						"2 servings cooked Japanese short-grain rice ((typically 1⅔ cups (250 g) per donburi serving))",
						"1 green onion/scallion",
						"½ tsp toasted white sesame seeds",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Before You Start: For the steamed rice, please note that 1½ cups (300 g, 2 rice cooker cups) of uncooked Japanese short-grain rice yield 4⅓ cups (660 g) of cooked rice, enough for 2 donburi servings (3⅓ cups, 500 g). See how to cook short-grain rice with a rice cooker, pot over the stove, Instant Pot, or donabe."},
						{Type: "HowToStep", Text: "Open the package of 14 oz medium-firm tofu (momen dofu) and drain out the water.Next, wrap the tofu block in a paper towel (or tea towel) and place it on a plate or tray. Now, press the tofu: First, put another tray or plate or even a cutting board on top of the tofu block to evenly distribute the weight. Then, place a heavy item* (I used a marble mortar but a can of food works) on top to apply pressure.Let it sit for at least 30 minutes before using. *The weighted item should not be so heavy that it will crumble or crush the tofu block but heavy enough that it will press out the tofu&#39;s liquid."},
						{Type: "HowToStep", Text: "While draining the tofu, you can cook the rice or a side dish. For this recipe, I also prepare this blanched broccoli recipe."},
						{Type: "HowToStep", Text: "Gather all the ingredients."},
						{Type: "HowToStep", Text: "To make the homemade teriyaki sauce, whisk the ¼ cup sake, ¼ cup mirin, ¼ cup soy sauce, and 4 tsp sugar in a (microwave-safe) medium bowl. If the sugar doesn‘t dissolve easily, microwave it for 30 seconds and whisk well. Set aside."},
						{Type: "HowToStep", Text: "Cut 1 green onion/scallion diagonally into thin slices."},
						{Type: "HowToStep", Text: "After 30 minutes of draining the tofu, remove the paper towel and transfer the tofu to the cutting board. First, cut the tofu block in half widthwise."},
						{Type: "HowToStep", Text: "Next, cut the tofu into roughly ¾-inch (2-cm) cubes."},
						{Type: "HowToStep", Text: "Put ⅓ cup potato starch or cornstarch in a shallow tray or bowl and gently coat the tofu cubes with the potato starch. Set aside."},
						{Type: "HowToStep", Text: "Heat a large frying pan on medium-high heat. When it‘s hot, add 1½ Tbsp of the 3 Tbsp neutral oil (keep the rest for the next batch) and distribute it evenly. Add the first batch of tofu cubes to the pan, placing them about 1 inch (2.5 cm) apart from each other so it‘s easy to rotate the tofu cubes without sticking to each other."},
						{Type: "HowToStep", Text: "Fry the cubes on one side until golden brown, then turn them to fry the next side. Repeat until all sides are brown and crispy. Transfer the fried tofu cubes to a plate or tray lined with a paper towel."},
						{Type: "HowToStep", Text: "Add the next batch of uncooked tofu to the pan and fry until crispy and golden brown on all sides. Add more of the remaining oil as needed to help brown the tofu faster."},
						{Type: "HowToStep", Text: "Remove all the fried tofu to the plate/tray."},
						{Type: "HowToStep", Text: "Wipe off any remaining oil in the pan with a paper towel. Then, transfer the tofu back into the pan."},
						{Type: "HowToStep", Text: "Add the teriyaki sauce to the pan; the sauce will start to thicken immediately. Quickly toss the tofu cubes in the sauce to coat, then turn off the heat and remove the pan from the stove. Tip: The sauce will continue to thicken with the residual heat, so if you want to keep some sauce in the pan, be sure to turn off the heat as soon as the tofu is coated."},
						{Type: "HowToStep", Text: "Divide 2 servings cooked Japanese short-grain rice into individual large (donburi) bowls. Serve the tofu and blanched broccoli over the steamed rice. Garnish the tofu with green onions and ½ tsp toasted white sesame seeds."},
						{Type: "HowToStep", Text: "You can keep the leftovers in an airtight container and store in the refrigerator for 3 days. Since the texture of the tofu changes when frozen, I don‘t recommend storing the tofu in the freezer."},
					},
				},
				Name: "Pan-Fried Teriyaki Tofu Bowl",
				NutritionSchema: &models.NutritionSchema{
					Calories:       "443 kcal",
					Carbohydrates:  "27 g",
					Fat:            "23 g",
					Fiber:          "3 g",
					Protein:        "21 g",
					SaturatedFat:   "3 g",
					Servings:       "1",
					Sodium:         "979 mg",
					Sugar:          "10 g",
					TransFat:       "1 g",
					UnsaturatedFat: "20 g",
				},
				PrepTime:  "PT5M",
				TotalTime: "PT60M",
				Yield:     &models.Yield{Value: 2},
				URL:       "https://www.justonecookbook.com/teriyaki-tofu-bowl/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
