package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_F(t *testing.T) {
	testcases := []testcase{
		{
			name: "farmhousedelivery.com",
			in:   "https://recipes.farmhousedelivery.com/green-shakshuka/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Breakfast"},
				DateModified:  "2023-01-20T18:47:32+00:00",
				DatePublished: "2023-01-20T08:13:29+00:00",
				Description: models.Description{
					Value: "Who knew a recipe could get us all excited with just a word? Shakshuka is awfully fun to say, but even more fun when you realize how delicious it is. An easy way to make dinner out of a few Read on!",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 yellow onion, slivered",
						"3 cloves garlic, minced",
						"2 Tbsp. olive oil",
						"1 jalapeno, seeded and minced",
						"4 big handfuls greens (mix kale & spinach), chopped and washed",
						"1/2 cup cream",
						"4 eggs",
						"Salt & pepper",
						"Plain yogurt, for serving",
						"Hot sauce, for serving",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Saute onions & garlic in olive oil in a large skillet until they take on a little color. Add jalapeno and continue cooking for 1-2 minutes. Add chopped greens, season with salt and pepper to taste and cover until greens are just wilted. Add cream and bring to a simmer. Crack eggs on top of greens, cover and cook until eggs are cooked to your preference. Serve in wide bowls with a dollop of yogurt and a drizzle of hot sauce and thick slices of warm bread on the side.",
					},
				},
				Name: "Green Shakshuka",
				URL:  "https://recipes.farmhousedelivery.com/green-shakshuka/",
			},
		},
		{
			name: "farmhouseonboone.com",
			in:   "https://www.farmhouseonboone.com/sourdough-pretzel-buns",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Sourdough"},
				CookTime:      "PT25M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-11-15",
				Description: models.Description{
					Value: "The most delicious sourdough pretzel buns have a soft and fluffy interior with a deep brown exterior that tastes just like the soft pretzels you love.",
				},
				Keywords: models.Keywords{Values: "baking, bread, sourdough pretzel buns"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"½ cup sourdough starter, active and bubbly (113 grams)",
						"1 cup water (236 grams)", "2 teaspoons sugar (8 grams)",
						"1/4 cup butter, softened (57 grams)", "2 teaspoons salt (10 grams)",
						"3 cups unbleached all purpose flour (420 grams)", "8 cups water", "1/2 cup baking soda",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Four to 12 hours before starting this recipe, feed the sourdough starter. It needs to be really active, bubbly, and able to pass the float test.",
						"In a bowl of a stand mixer with a dough hook attachment, add the dough ingredients.",
						"Mix on low speed (speed 1-2 on a Kitchen Aid mixer) until it forms a shiny dough ball that will pass the windowpane test, about 7-10 minutes. Take a small amount of dough and stretch it into a square. If you can stretch it and see through the dough without the dough breaking, it is ready to move on to the next stage",
						"Place the dough into a lightly oiled bowl and cover with a lid, plastic wrap, or beeswax wraps.",
						"Allow it to bulk ferment for 8-12 hours or until doubled. If your house is on the colder side, it may take longer. It should just double in size.",
						"Divide into 10 equal pieces.",
						"Shape the dough into round dough balls. I like to create tension in the dough by pulling the dough into the center and kind of folding it in.",
						"Place on a parchment lined baking sheet and cover with a tea towel.",
						"Allow to rise in a warm spot until doubled in size.",
						"Preheat oven to 425 degrees.",
						"Add 8 cups of water and baking soda to a large pot and bring to a boil.",
						"Gently place each roll into the boiling water and boil for about 2 minutes. Flip and boil on the other side for another 2 minutes.",
						"Transfer the boiled pretzel buns to the parchment lined baking sheet.",
						"Score an X on top of each bun with a lame (a razor blade used in baking) or sharp knife and sprinkle with coarse sea salt.",
						"Bake 20-25 minutes until achieving a deep golden brown color.",
						"Allow to cool, then serve.",
					},
				},
				Name: "Sourdough Pretzel Buns",
				NutritionSchema: models.NutritionSchema{
					Calories:       "202 calories",
					Carbohydrates:  "34 grams carbohydrates",
					Cholesterol:    "12 milligrams cholesterol",
					Fat:            "5 grams fat",
					Fiber:          "1 grams fiber",
					Protein:        "5 grams protein",
					SaturatedFat:   "3 grams saturated fat",
					Servings:       "1",
					Sodium:         "462 milligrams sodium",
					Sugar:          "1 grams sugar",
					TransFat:       "0 grams trans fat",
					UnsaturatedFat: "2 grams unsaturated fat",
				},
				PrepTime: "PT15M",
				Yield:    models.Yield{Value: 10},
				URL:      "https://www.farmhouseonboone.com/sourdough-pretzel-buns",
			},
		},
		{
			name: "fattoincasadabenedetta.it",
			in:   "https://www.fattoincasadabenedetta.it/ricetta/lasagne-al-pistacchio/",
			want: models.RecipeSchema{
				AtContext: "https://schema.org",
				AtType:    models.SchemaType{Value: "Recipe"},
				Category:  models.Category{Value: "Pasta"},
				CookTime:  "P0DT0H40M0S",
				Description: models.Description{
					Value: "La lasagna al pistacchio è una lasagna bianca gustosa, caratterizzata da un contrasto di sapori fantastico. Un piatto della domenica speciale!",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"pistacchi 150 g al naturale", "olio di semi di girasole 80 g",
						"besciamella 500 ml", "sfoglie sottili pronte per lasagne 200 g",
						"mozzarella 250 g a dadini", "formaggio grattugiato 60 g",
						"speck 200 g a fette", "granella di pistacchio",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In un tritatutto mettiamo i pistacchi e l’olio di semi di girasole e frulliamo bene per ottenere un pesto cremoso.",
						"In una ciotola mettiamo la besciamella e uniamo i pistacchi frullati, mescolando con un cucchiaio.",
						"Prepariamo tutti gli ingredienti e assembliamo la nostra lasagna: prendiamo una pirofila da forno 30&#215;20 cm e spalmiamo sul fondo la besciamella al pistacchio.",
						"Sistemiamo sopra la prima sfoglia e mettiamo altra besciamella, poi qualche fetta di speck, la mozzarella a dadini e il formaggio grattugiato.",
						"Proseguiamo con gli altri strati fino a terminare gli ingredienti. In superficie mettiamo uno strato di besciamella al pistacchio, una spolverata di formaggio grattugiato e infine della granella di pistacchio.",
						"Inforniamo e cuociamo in forno ventilato a 170 °C per 35-40 minuti, oppure in forno statico a 180 °C per lo stesso tempo.",
						"Ecco pronta la nostra gustosissima lasagna al pistacchio! Buon appetito!",
					},
				},
				Name:     "Lasagna al pistacchio",
				PrepTime: "P0DT0H10M0S",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.fattoincasadabenedetta.it/ricetta/lasagne-al-pistacchio/",
			},
		},
		{
			name: "fifteenspatulas.com",
			in:   "https://www.fifteenspatulas.com/guacamole/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizer"},
				Cuisine:       models.Cuisine{Value: "Mexican"},
				DatePublished: "2023-06-05T09:44:00+00:00",
				Description: models.Description{
					Value: "This Homemade Guacamole has the perfect texture and combination of flavors, with chunky " +
						"mashed avocados mixed with fresh lime juice, jalapeno, white onion, tomatoes, and cilantro.",
				},
				Keywords: models.Keywords{Values: "guacamole"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 ripe avocados",
						"1/4 cup seeded and diced tomato",
						"2 tbsp finely chopped white onion",
						"2 tbsp minced jalapeno",
						"1 tbsp freshly squeezed lime juice",
						"1/2 tsp salt (or to taste)",
						"2 tbsp chopped cilantro",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Cut the avocados in half and remove the pits.",
						"Scoop the avocado flesh out into a bowl, and mash the avocado with a fork, leaving plenty " +
							"of chunky, unmashed bits of avocado.",
						"Add the tomato, onion, jalapeno, lime juice, and salt, then gently stir to combine.",
						"Gently fold in the cilantro.",
						"Taste the guacamole and adjust to your tastes (you may desire more salt, or more acidity), " +
							"then serve. Enjoy!",
					},
				},
				Name: "Guacamole",
				NutritionSchema: models.NutritionSchema{
					Calories:      "168 kcal",
					Carbohydrates: "10 g",
					Fat:           "15 g",
					Fiber:         "7 g",
					Protein:       "2 g",
					SaturatedFat:  "2 g",
					Servings:      "1",
					Sodium:        "307 mg",
					Sugar:         "2 g",
				},
				PrepTime: "PT10M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.fifteenspatulas.com/guacamole/",
			},
		},
		{
			name: "finedininglovers.com",
			in:   "https://www.finedininglovers.com/recipes/main-course/szechuan-chicken",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Category: models.Category{
					Value: "Main Course",
				},
				Description: models.Description{
					Value: "Szechuan Chicken is a spicy, crispy chicken recipe from Sichuan Region in China: discover the original recipe and try to make it at home.",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"Chicken",
						"White onion",
						"Red pepper",
						"Yellow bell peppers",
						"Chilli",
						"Ginger",
						"Cane sugar",
						"Cornstarch",
						"Garlic powder",
						"Dark soy sauce",
						"Sesame Oil",
						"Chicken stock",
						"Szechuan pepper",
						"Sunflower oil",
						"Kosher salt",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Start with the chicken. Cut it into pieces of equal size that are not too small so that the meat " +
							"can have the same cooking time. Put the cornstarch in a bowl (set aside a spoon) and toss " +
							"the chicken in the flour.\nIn a non-stick pan",
						"pour a drizzle of seed oil and when hot",
						"add the floured chicken. Cook the chicken nuggets for 3/4 minutes so that they become crunchy " +
							"and amber on the outside. In the meantime",
						"wash the peppers",
						"remove the stalk",
						"seeds and internal filaments",
						"then cut them into small pieces. Peel the onion and cut this into small pieces of the same size " +
							"as the peppers.\nWhen the chicken is well browned",
						"remove it from the pan and keep it warm in a covered bowl. In the same pan",
						"add the peppers",
						"onion",
						"whole chillies",
						"garlic powder and freshly grated ginger. Season with salt and bring to the fire to dry the " +
							"vegetables. It will take about 5 minutes on high heat. If you have decided to add it",
						"now is the time to add the Szechuan pepper.\nMeanwhile",
						"prepare the sauce by combining the soy sauce",
						"sesame oil",
						"chicken broth",
						"sugar and a tablespoon of cornstarch in a bowl. Mix very well.\nAdd the chicken to the vegetables " +
							"and cook for another 5 minutes",
						"then add the sauce and mix well. After 30 seconds your super creamy and spicy Szechuan chicken " +
							"is ready. Serve it alone or accompanied with white rice. If you like",
						"you can add chopped fresh parsley or coriander. Tip & Tricks The use of corn starch instead of normal flour",
						"both in the flour and in the preparation of the sauce",
						"gives the dish a creamy texture. Origins This recipe has its origins in the Szechuan region",
						"a place where every dish is prepared with spicy hints. Also typical of this area is Szechuan pepper",
						"a fragrant berry that is added to dishes dried or in powder. Variants For a slightly less " +
							"spicy Szechuan chicken",
						"you can remove the seeds from the peppers or add less during cooking. To give a note of freshness " +
							"you can add the grated lime or lemon zest just before serving.",
					},
				},
				Keywords: models.Keywords{Values: "Tried and Tasted,Channel: Food & Drinks"},
				Name:     "Szechuan Chicken",
				NutritionSchema: models.NutritionSchema{
					Servings: "4",
				},
				Yield: models.Yield{Value: 4},
				URL:   "https://www.finedininglovers.com/recipes/main-course/szechuan-chicken",
			},
		},
		{
			name: "fitmencook.com",
			in:   "https://fitmencook.com/rosemary-blue-cheese-turkey-sliders/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT20M",
				DateModified:  "2023-05-16T19:59:35+00:00",
				DatePublished: "2021-09-05T18:05:31+00:00",
				Keywords:      models.Keywords{Values: "meal prep,meat,turkey"},
				Image:         models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"Meat",
						"1 1/2 lb 93% lean ground turkey",
						"1 1/2 tablespoons The Fit Cook Everyday Blend",
						"2 tablespoons fresh rosemary, finely chopped",
						"1/3 cup blue cheese crumble",
						"pinch of sea salt & pepper",
						"\n",
						"Sliders",
						"8 (wheat) slider buns",
						"2 roma tomatoes, sliced",
						"1 medium cucumber, diagonally sliced",
						"8 tablespoons Dijon mustard",
						"8 Romaine lettuce leaves",
						"\n",
						"Quick Caramelized Onions (OPTIONAL)",
						"3 tablespoons olive oil",
						"1 large white onion, sliced",
						"pinch of sea salt & peppe",
						"\n",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a bowl, mix together the ingredients for the turkey sliders.  Use an ice cream scoop or 1/4 " +
							"cup measuring cup to scoop to measure out each slider to keep them uniform. Alternatively " +
							"you can measure out full sized burgers.",
						"Set a nonstick skillet on medium heat and once hot, spray with avocado or olive oil, then add " +
							"the slider patties. Cook for 4 – 6 minutes on each side, or until the top/bottom are " +
							"browned and the slider is cooked through.",
						"In a carbon steel skillet on medium heat and once hot, add olive oil and onion. As the onions " +
							"saute and caramelize, add a pinch of sea salt to draw out sweetness. If desired, reduce " +
							"the onions in 1/2 cup of (chicken/veggie/beef) broth, white wine, or water. Continue " +
							"cooking until the onions are “wilted,” soft and golden brown, about 15 minutes.",
						"Build the sliders! Toast the buns then add Dijon, lettuce, tomato, cucumber, caramelized onions " +
							"and the burger!",
					},
				},
				Name: "Rosemary Blue Cheese Turkey Sliders",
				NutritionSchema: models.NutritionSchema{
					Calories:      "330cal",
					Carbohydrates: "24g",
					Fat:           "15g",
					Fiber:         "2g",
					Protein:       "21g",
					Sodium:        "670mg",
					Sugar:         "5g",
				},
				PrepTime: "PT5M",
				URL:      "https://fitmencook.com/rosemary-blue-cheese-turkey-sliders/",
			},
		},
		{
			name: "food.com",
			in:   "https://www.food.com/recipe/jim-lahey-s-no-knead-pizza-margherita-382696",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Lunch/Snacks"},
				CookTime:      "PT6M",
				DatePublished: "2009-07-24T13:27Z",
				Description: models.Description{
					Value: "This is a great recipe for a simple, thin crust pizza.  It's from Jim Lahey (of no-knead " +
						"bread fame) who now runs a popular NYC pizzeria called Co.  The recipe was printed in New " +
						"York Magazine (Jul 12, 2009).  If you don't have a pizza stone, this works well in a cast " +
						"iron skillet.  The recipe requires very little time and effort but the dough must be started the day before.",
				},
				Keywords: models.Keywords{
					Values: "Cheese,Vegetable,European,Low Cholesterol,Healthy,Kid Friendly,Kosher,Broil/Grill,< 60 Mins," +
						"Oven,Beginner Cook,Easy,Inexpensive",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"3  cups  all-purpose flour or 3  cups  bread flour, more for dusting",
						"1/4 teaspoon  instant yeast",
						"1 1/2 teaspoons  salt",
						"1 1/4 cups  water",
						"1   vine-ripened tomatoes (about 5 oz.) or 1    heirloom tomato (about 5 oz.)",
						"1  pinch  salt",
						"1/4 teaspoon  extra virgin olive oil",
						"5  tablespoons  tomato sauce",
						"2  ounces  buffalo mozzarella (about 1/4 ball)",
						"basil leaves",
						"1  tablespoon  extra virgin olive oil",
						"salt",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"To make dough: In a large bowl, mix the flour, yeast, and salt. Add water and stir until blended " +
							"(the dough will be very sticky). Cover with plastic wrap and let rest for 12 to 24 hours in a warm " +
							"spot, about 70 degrees.",
						"Place the dough on a lightly floured work surface and sprinkle the top with flour. Fold the dough " +
							"over on itself once or twice, cover loosely with plastic wrap, and let rest for 15 minutes.",
						"Shape the dough into 3 or 4 balls, depending on how thick you want the crust. Generously sprinkle " +
							"a clean cotton towel with flour and cover the dough with it. Let the dough rise for 2 hours.",
						"To make sauce: Blanch tomato for 5 seconds in boiling water and quickly remove. Allow to cool to " +
							"the touch. Peel the skin with your hands and quarter the tomato. Remove the jelly and seeds, and reserve " +
							"in a strainer or fine sieve. Strain the jelly to remove seeds, and combine resulting liquid in a " +
							"bowl with the flesh of the tomatoes. Proceed to crush the tomatoes with your hands. Add salt " +
							"and olive oil and stir.",
						"To make pizza: Place pizza stone on the middle rack of the oven and preheat on high broil. Stretch " +
							"or toss the dough into a disk approximately 10 inches in diameter. Pull rack out of oven and place the " +
							"dough on top of the preheated pizza stone. Drizzle 5 generous tablespoons of sauce over the dough, and " +
							"spread evenly. Try to keep the sauce about &frac12; inch away from the perimeter of the dough. Break " +
							"apart or slice the buffalo mozzarella and arrange over the dough. Return rack and pizza stone to the " +
							"middle of the oven and broil for approximately 6 minutes. Remove and top with basil, olive oil, and salt.",
					},
				},
				Name: "Jim Lahey’s No-Knead Pizza Margherita",
				NutritionSchema: models.NutritionSchema{
					Calories:      "569.4",
					Carbohydrates: "98.9",
					Cholesterol:   "15",
					Fat:           "10.5",
					Fiber:         "4.3",
					Protein:       "17.9",
					SaturatedFat:  "3.4",
					Sodium:        "1472",
					Sugar:         "2.7",
				},
				PrepTime: "PT30M",
				URL:      "https://www.food.com/recipe/jim-lahey-s-no-knead-pizza-margherita-382696",
				Yield:    models.Yield{Value: 1},
			},
		},
		{
			name: "food52.com",
			in:   "https://food52.com/recipes/7930-orecchiette-with-roasted-butternut-squash-kale-and-caramelized-red-onion",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizer"},
				CookTime:      "PT1H0M",
				Cuisine:       models.Cuisine{Value: "Italian"},
				DateModified:  "2021-04-17 19:58:03 -0400",
				DatePublished: "2010-11-21 19:54:38 -0500",
				Description: models.Description{
					Value: "This recipe is for the butternut squash lover. This orecchiette recipe is yet another reason " +
						"to love squash (and kale, and adorably shaped pasta).",
				},
				Keywords: models.Keywords{
					Values: "Onion, Kale, Goat Cheese, Butternut Squash, Sage, Milk/Cream, Nutmeg, Pasta, Fall",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 large butternut squash, cut into small cubes, divided",
						"4 tablespoons extra-virgin olive oil, divided",
						"Kosher salt and pepper",
						"1 pinch cayenne pepper",
						"1/4 teaspoon ground nutmeg",
						"1 red onion, sliced thinly",
						"1/2 pound orecchiette",
						"1 or 2 garlic cloves, minced",
						"2 cups chicken broth, divided",
						"1 bunch kale",
						"1/2 cup white wine",
						"1/2 cup heavy cream",
						"1 ounce goat cheese, optional",
						"1 tablespoon chopped fresh sage",
						"Parmesan cheese, to serve",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 425° F. Toss all but 1 cup of the butternut squash with 1 tablespoon olive " +
							"oil, salt, pepper, a pinch of cayenne pepper, and the nutmeg. Roast until butternut " +
							"squash pieces are tender and caramelized, about 30 minutes. Set aside.",
						"Heat 1 tablespoon olive oil in a medium saucepan over low heat. Cook sliced red onions " +
							"until caramelized, about 30 minutes. Set aside.",
						"Heat a pot of water over high heat until boiling. Salt water generously. Cook orecchiette " +
							"according to package instructions until al dente.",
						"Meanwhile, heat another tablespoon of olive oil in a heavy pan over medium-high heat. Cook " +
							"the remaining cup of butternut squash for approximately 3 minutes. Add garlic and " +
							"cook for another minute. Add 1/2 cup of the chicken broth and cook until broth is almost completely absorbed.",
						"Remove the middle stems from the kale and roughly chop the leaves. Add kale to butternut " +
							"squash and stir until kale has softened. Add caramelized red onions.",
						"Add white wine and cook for 2 minutes. Add remaining chicken broth and reduce, about 10 minutes.",
						"Turn heat to low and add the heavy cream. When the pasta is al dente, add it to the pan with the " +
							"sauce. Add the roasted butternut squash.",
						"Loosen sauce with pasta water if needed. Sprinkle with goat cheese (optional), sage, and Parmesan cheese.",
					},
				},
				Name:     "Orecchiette With Roasted Butternut Squash, Kale, &amp;amp; Caramelized Red Onion",
				PrepTime: "PT0H15M",
				Yield:    models.Yield{Value: 4},
				URL: "https://food52.com/recipes/7930-orecchiette-with-roasted-butternut-squash-kale-" +
					"and-caramelized-red-onion",
			},
		},
		{
			name: "foodandwine.com",
			in:   "https://www.foodandwine.com/recipes/garlic-salmon-with-sheet-pan-potatoes",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Dinner"},
				DateModified:  "2023-08-02T10:26:55.248-04:00",
				DatePublished: "2022-03-25T10:20:00.000-04:00",
				Description: models.Description{
					Value: "For this sheet pan dinner, baby potatoes, red onion, and spring onions get a head start in a hot oven, before they are joined by a side of salmon, slathered with mustard and drizzled with toasted garlic oil, which cooks alongside the vegetables for a seamless final presentation. Sommelier Erin Miller, of Charlie Palmer&#39;s Dry Creek Kitchen in Healdsburg, California, who provided the inspiration for this dish, notes that it tastes even better when served with a great wine. She recommends a glass of Hirsch Vineyards Raschen Ridge Sonoma Coast Pinot Noir, noting, &#34;The bright acidity of the Hirsch Pinot Noir is a perfect foil for the fresh, fatty fish and flavors of garlic and lemon.&#34;",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1.5 pounds baby yellow potatoes, halved lengthwise",
						"1 large red onion, cut into 1/2-inch wedges",
						"5 spring onions (about 6 ounces), trimmed and halved lengthwise",
						"0.333 cup plus 2 tablespoons extra-virgin olive oil, divided",
						"2 teaspoons kosher salt, divided",
						"0.5 teaspoon black pepper, divided",
						"0.25 cup Dijon mustard, divided",
						"1 (2-pound) skin-on or skinless salmon side",
						"3 garlic cloves, thinly sliced (about 1 1/2 tablespoons)",
						"2 tablespoons chopped fresh tarragon",
						"2 tablespoons chopped fresh chives",
						"Lemon wedges, for serving",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 425°F. Toss together potatoes, red onion, spring onions, and 2 tablespoons oil on a " +
							"large rimmed baking sheet; spread in an even layer. Sprinkle evenly with 1 teaspoon salt " +
							"and 1/4 teaspoon pepper. Roast in preheated oven until vegetables begin to brown, about 20 " +
							"minutes. Remove from oven, and reduce oven temperature to 325°F.",
						"Add 1 tablespoon mustard to vegetable mixture on baking sheet; toss to coat. Push vegetables to long " +
							"edges of baking sheet. Place salmon, skin side down, lengthwise in middle of baking sheet. " +
							"Spread salmon with remaining 3 tablespoons mustard; sprinkle with remaining 1 teaspoon salt " +
							"and remaining 1/4 teaspoon pepper.",
						"Heat remaining 1/3 cup oil in a large skillet over medium-high. Add garlic; cook, stirring often, " +
							"until garlic is fragrant and light golden brown, about 2 minutes. Pour hot oil mixture over " +
							"salmon on baking sheet. Roast at 325°F until salmon flakes easily with a fork and vegetables " +
							"are tender, 12 to 15 minutes. Remove from oven. Transfer salmon to a platter; sprinkle with " +
							"tarragon and chives. Transfer vegetables to a bowl. Serve salmon and vegetables with lemon wedges.",
					},
				},
				Name:  "Sizzling Garlic Salmon with Sheet Pan Potatoes",
				Yield: models.Yield{Value: 4},
				URL:   "https://www.foodandwine.com/recipes/garlic-salmon-with-sheet-pan-potatoes",
			},
		},
		{
			name: "foodnetwork.co.uk",
			in:   "https://foodnetwork.co.uk/recipes/waldorf-chicken-boats",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				CookTime:      "PT10M",
				DateCreated:   "2015-11-30T11:55:59+00:00",
				DateModified:  "2020-09-29T11:10:46+00:00",
				DatePublished: "2015-11-30T11:55:59+00:00",
				Description: models.Description{
					Value: "Introduce the kids to salad with this one that's fun to assemble and eat.",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"For the dressing:", "1/2 cup low-fat Greek yogurt",
						"1 tablespoon lemon juice", "1 tablespoon olive oil",
						"2 teaspoons apple cider vinegar", "1/2 teaspoon chopped fresh thyme leaves",
						"Salt and freshly ground black pepper", "For the salad:", "1/4 cup walnuts",
						"1 1/2 cups 1/2-inch pieces rotisserie chicken, skin and bones discarded",
						"1 small crisp apple, cored and cut into 1/2-inch pieces",
						"1/2 cup red seedless grapes, halved",
						"3 spring onions, whites and greens divided and chopped",
						"8 hearts of romaine leaves",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat the oven to 175°C.",
						"For the dressing: Whisk together the yogurt, lemon juice, oil, vinegar, thyme, 1/2 teaspoon salt and 1/4 teaspoon pepper in a small bowl until combined; set aside.",
						"For the salad: Spread the walnuts on a baking sheet, and toast until lightly browned and fragrant, 5 to 7 minutes. Let cool for 5 minutes, then break up into 1/4-inch pieces.",
						"To assemble: Toss the walnuts, chicken, apples, grapes and scallion whites together with the dressing in a large bowl until everything is well coated; season with salt and pepper. Scoop a scant 1/2 cup of the salad into each romaine leaf, and arrange 2 on each serving plate. Garnish with the scallion greens, and serve.",
						"Copyright 2015 Television Food Network, G.P. All rights reserved.",
						"Note: Never leave a child unattended in the kitchen. Limit the child to tasks that are safe and age-appropriate.",
						"From Food Network Kitchen",
					},
				},
				Name:     "Waldorf Chicken Boats",
				PrepTime: "PT20M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://foodnetwork.co.uk/recipes/waldorf-chicken-boats",
			},
		},
		{
			name: "foodrepublic.com",
			in:   "https://www.foodrepublic.com/recipes/hand-cut-burger/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "How To Cook A Hand-Cut Burger",
				Description: models.Description{
					Value: "When you don't have a meat grinder, but still want a nice, juicy burger, this recipe for a hand-cut burger has a trick you'll use over and over again.",
				},
				Image:         models.Image{Value: anUploadedImage.String()},
				DateModified:  "2018-06-07T23:13:57+00:00",
				DatePublished: "2018-06-08T15:00:40+00:00",
				Ingredients: models.Ingredients{
					Values: []string{
						"1 (1 1/2-pound) boneless rib-eye steak (preferably dry-aged)",
						"2 cups unsalted butter or rendered beef tallow, plus 2 tablespoons unsalted butter",
						"1 white onion",
						"1 (5-ounce) piece horseradish (2 to 3 inches)",
						"2 tablespoons buttermilk",
						"kosher salt",
						"4 Pain de Mie Buns or 8 slices soft slab bread",
						"16 to 24 dill pickle slices",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Chill the steak in the freezer until firm to the touch but not frozen, 15 to 20 minutes. Cut the steak into 1⁄4-inch-thick slices, then slice into 1⁄4-inch-thick strips, and then into 1⁄4-inch cubes. Remove the sinew and connective tissue but keep the fat.",
						"Divide the beef into four equal balls. Put a sheet of plastic wrap over a 4-inch ring mold on a cutting board or other hard surface. Put a ball in the middle of the mold and gently press down with the palm of your hand, forming a patty that is 4 inches wide. Pop it out with the plastic wrap. Put the patties on a large dish or small baking sheet and refrigerate until ready to cook.",
						"Melt 2 cups of the butter in a pot over medium heat. (Why yes, that is a lot of butter, but it’s used to fully submerge the onion while it cooks; you will not eat 2 cups of butter in this burger.) Add the onion, turn the heat to low, and gently cook at a bare simmer until the onion is tender, about 20 minutes. The onion should be cooked but still al dente, so there’s some texture and a slight hit of sharpness yet not enough that you’ll taste onion the rest of the day. Remove the onion from the butter and drain on a paper towel.",
						"While the onion cooks, make a horseradish sauce. In a bowl, mix the grated horseradish with the buttermilk and a pinch of salt. Stir to combine and refrigerate until ready to use.",
						"Before you begin cooking the burgers, get the buns toasting. Heat a cast-iron skillet or similar surface over medium-low heat. Slice the buns in half horizontally. Smear the remaining 2 tablespoons of butter on the buns and place, butter side down, on the hot surface, working in batches if necessary. Toast until golden brown, 6 to 8 minutes, adjusting the heat if necessary. You want to do your best to time their completion to the burger cooking.",
						"While the buns toast, cook the patties. Heat a cast-iron skillet or grill over high heat. Use a spatula to handle the patty—it will be loose, so be careful. Salt both sides of each patty and put them on the hot skillet. Cook on one side, about 1 minute, then flip the patties and cook until rare, another minute.",
						"Place a patty on a bottom bun and top with some pickles and onions. Slather 1 1⁄2 teaspoons horseradish sauce on the top bun and cap it off. Repeat.",
					},
				},
				URL: "https://www.foodrepublic.com/recipes/hand-cut-burger/",
			},
		},
		{
			name: "forksoverknives.com",
			in:   "https://www.forksoverknives.com/recipes/vegan-snacks-appetizers/crispy-buffalo-cauliflower-bites/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Category: models.Category{
					Value: "vegan snacks",
				},
				CookTime:      "PT0D0H35M",
				DatePublished: "2017-01-27 15:29:56",
				Description: models.Description{
					Value: "It took a lot of trial and error to find the right coating that would not draw out the moisture " +
						"and would make the florets crisp, so I am pleased that it has turned out to be a very " +
						"simple recipe. You will not need to add salt as the sauces have enough salt to season " +
						"them. Either a smoky barbecue sauce or Frank’s hot sauce would work well, but if you " +
						"are like me and prefer sweet and spicy, then try a little bit of both. Serve with " +
						"ranch or Caesar dressing on the side if you wish, or whip up a batch of Spinach Ranch " +
						"Dip. Note: The buffalo cauliflower bites will get softer once they are coated with " +
						"the sauce, so hold off tossing until the very last minute",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"⅔ cup brown rice flour",
						"2 tablespoons almond flour",
						"1 tablespoon tomato paste",
						"2 teaspoons garlic powder",
						"2 teaspoons onion powder",
						"2 teaspoons smoked paprika",
						"1 teaspoon dried parsley",
						"1 head cauliflower, cut into 2-inch florets",
						"⅓ cup Frank’s hot sauce or barbecue sauce",
						"<a href=\"https://www.forksoverknives.com/recipes/spinach-ranch-dip/#gs.MasyuTc\" target=\"_blank\" " +
							"rel=\"noopener\">Spinach Ranch Dip</a>",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat oven to 450°F. Line 2 baking sheets with parchment paper.",
						"Combine the brown rice flour, almond flour, tomato paste, garlic powder, onion powder, paprika, " +
							"parsley, and ⅔ cup of water in a blender. Puree until the batter is smooth and thick. Transfer " +
							"to a bowl and add the cauliflower florets; toss until the florets are well coated with the batter.",
						"Arrange the cauliflower in a single layer on the prepared baking sheets, making sure that the " +
							"florets do not touch one another. Bake for 20 to 25 minutes, until crisp on the edges. They " +
							"will not get crispy all over while still in the oven.",
						"Remove from the heat and let stand for 3 minutes to crisp up a bit more. Transfer to a bowl " +
							"and drizzle with the sauce. Serve immediately.",
					},
				},
				Name:     "Crispy Buffalo Cauliflower Bites",
				PrepTime: "PT0D0H0M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://www.forksoverknives.com/recipes/vegan-snacks-appetizers/crispy-buffalo-cauliflower-bites/",
			},
		},
		{
			name: "forktospoon.com",
			in:   "https://forktospoon.com/air-fryer-blooming-onion-bites",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetizer"},
				CookTime:      "PT15M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-11-20T11:59:57+00:00",
				Description: models.Description{
					Value: "Air Fryer Blooming Onion Bites -- Step into a world of crispy, flavorful delight with our Air Fryer Blooming Onion Bites recipe!",
				},
				Keywords: models.Keywords{Values: "Air Fryer Blooming Onion Bites"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"12 ounces pearl onions (frozen)", "2 tablespoon olive oil",
						"½ teaspoon paprika", "½ teaspoon oregano", "½ teaspoon ground thyme",
						"½ teaspoon cumin", "½ teaspoon kosher salt", "½ teaspoon black pepper",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Mix the Seasonings: In a large bowl, combine olive oil, paprika, oregano, ground thyme, cumin, salt, and pepper.",
						"Add Onions: Take the frozen pearl onions (do not defrost) and add them to the bowl. Toss them gently to ensure they are evenly coated with the seasoning mix.",
						"Air Fry: Preheat your air fryer. Place the seasoned onions in the air fryer basket and cook at a high temperature (around 400°F or 200°C) for about 10-15 minutes, or until they are golden brown and crispy. Shake the basket halfway through for even cooking.",
						"Serve: Once cooked, let them cool slightly before serving. These onion bites are perfect as a snack or appetizer, served with your favorite dipping sauce.",
					},
				},
				Name: "Air Fryer Blooming Onion Bites",
				NutritionSchema: models.NutritionSchema{
					Calories:       "66 kcal",
					Carbohydrates:  "6 g",
					Fat:            "5 g",
					Fiber:          "1 g",
					Protein:        "1 g",
					SaturatedFat:   "1 g",
					Servings:       "1",
					Sodium:         "197 mg",
					Sugar:          "2 g",
					UnsaturatedFat: "4 g",
				},
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 6},
				URL:      "https://forktospoon.com/air-fryer-blooming-onion-bites",
			},
		},
		{
			name: "franzoesischkochen.de",
			in:   "https://www.franzoesischkochen.de/navettes-aus-marseille/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Alte französische Rezepte"},
				CookTime:      "PT30M",
				Cuisine:       models.Cuisine{Value: "französische Rezepte"},
				DateModified:  "2022-03-09T15:14:09+02:00 ",
				DatePublished: "2022-03-09T15:14:09+02:00 ",
				Description:   models.Description{Value: "Ein einfaches Rezept mit Schritt-für-Schritt-Fotos und vielen Tipps über das Thema: Navettes aus Marseille"},
				Keywords:      models.Keywords{Values: "französisch,frankreich,Alte französische Rezepte,Einfachste Rezepte,In der Boulangerie,Kekse &amp; Plätzchen,Provence,Traditionelle Rezepte,Typisch französische Kuchen"},
				Image:         models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 Ei",
						"160g Mehl T55",
						"20 g Olivenöl",
						"60 g Honig",
						"1 TL Orangenblütenwasser",
						"1 kleine Prise Salz. Milch zum Pinseln.",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Das Ei mit dem Honig, Orangenblütenwasser und Olivenöl rühren....",
						"Kleine Teigkugels (25 bis 30g ) nehmen, wie Knete in die Länge ausrollen. Sie sollten eine " +
							"ovale Form oder besser eine Schiff-Form bekommen. Die geformten Navette auf " +
							"einem mit Backpapier belegten Backblech legen. Die Navettes in der Mitte entlang anschneiden.",
						"Backen: 180°C Umluft für 12 bis 15 Minuten. (es kommt darauf an, ob Ihr eure Navettes brauner " +
							"mögt wie ich oder lieber hell!)",
					},
				},
				Name:     "Navettes aus Marseille",
				PrepTime: "PT60M",
				Yield:    models.Yield{Value: 3},
				URL:      "https://www.franzoesischkochen.de/navettes-aus-marseille/",
			},
		},
		{
			name: "fredriksfika.allas.se",
			in:   "https://fredriksfika.allas.se/fredriks-pajer/knackig-appelpaj-med-rakram/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Fika, äppelpaj"},
				DateCreated:   "2020-08-27T17:18:00.000000Z",
				DateModified:  "2023-01-07T09:01:57.000000Z",
				DatePublished: "2020-08-27T17:18:00.000000Z",
				Description:   models.Description{Value: "Godaste äppelpajen"},
				Image:         models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"5-6 äpplen som har fast fruktkött", "2 tsk kanel", "2 msk strösocker",
						"150 g smör", "1,5 dl strösocker", "0,5 dl ljus sirap", "0,5 dl vispgrädde",
						"3 dl havregryn", "1,5 dl vetemjöl", "0,5 tsk bakpulver", "0,5 tsk salt",
						"2 äggulor", "2 msk strösocker", "1 msk vaniljsocker", "2 dl vispgrädde",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Sätt ugnen på 175 grader.",
						"Skala, kärna ur och klyfta äpplena, blanda med socker och kanel och lägg i en pajform på cirka 20 x 30 cm.",
						"Smält smöret i en stekpanna eller kastrull tillsammans med strösocker, sirap och grädde.",
						"Blanda ner havregryn och vetemjöl blandat med bakpulver och salt och rör till en jämn smet. Ta av från värmen.",
						"Fördela smeten över äpplena i formen.",
						"Grädda i nedre delen av ugnen i ca 35-40 minuter, tills pajen är gyllenbrun. Låt svalna något innan servering.",
						"Råkräm – Vispa äggulor, strösocker och vaniljsocker fluffigt.",
						"Vispa grädden rätt fast.", "Vänd ner äggvispet i grädden.",
					},
				},
				Name:  "Fyllning",
				Yield: models.Yield{Value: 6},
				URL:   "https://fredriksfika.allas.se/fredriks-pajer/knackig-appelpaj-med-rakram/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
