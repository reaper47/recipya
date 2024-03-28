package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_D(t *testing.T) {
	testcases := []testcase{
		{
			name: "davidlebovitz.com",
			in:   "https://www.davidlebovitz.com/marcella-hazans-bolognese-sauce-recipe-italian-beef-tomato/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Course"},
				DatePublished: "2023-09-27T10:25:00+00:00",
				Description: models.Description{
					Value: "There are a few Marcella Hazan Pasta Bolognese recipes out there. I tweaked a few of them to come up with this recipe, which is inspired and adapted by her. Note that this sauce will take a while to prepare. It&#39;s mostly downtime. At first, you&#39;re just sauteeing ingredients, stirring until they&#39;re combined, then adding wine and milk, simmering and stirring until those have been absorbed. Once the tomatoes have been added, that&#39;s when you let the sauce cook at the lowest heat possible, stirring every once in a while, until it&#39;s ready. Within an hour, it comes together into a nice paste, but if you cook it another hour, nursing it with water as you go, you&#39;ll get a sauce with a richer flavor. One of Hazan&#39;s recipes says to cook sauce Bolognese for up to 4 hours! In the untraditional category, some people like to grate from Parmesan cheese over finished bowls of pasta.",
				},
				Keywords: models.Keywords{Values: "mushrooms"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 tablespoons butter (salted or unsalted, plus 1 tablespoon for finishing the pasta)",
						"3 tablespoons olive oil", "1/2 cup (60g) diced onions",
						"1/2 cup (65g) diced celery", "1/2 cup (65g) diced carrots ((peeled))",
						"12 ounces (340g) ground beef ((I recommend using one that's at least 15% fat))",
						"1 teaspoon salt (or more to taste)", "freshly-ground black pepper",
						"3/4 cup (180ml) whole milk", "1/8 teaspoon freshly ground nutmeg",
						"1 cup (250ml) dry white wine",
						"1 1/2 cups (350ml) canned plum tomatoes (crushed, with their juice)",
						"1 1/2 tablespoons tomato paste", "1 pound (450g) pasta",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Melt the butter with the olive oil in a soup pot over medium-high heat. Add the onions and cook, stirring frequently, until the onions are soft and translucent, about 3 minutes. Add the celery and carrots and cook, stirring a few times, until they start to wilt, 3 to 4 minutes.",
						"Add the ground beef, salt, and some freshly ground pepper. Lower the heat to medium and cook, stirring, until the beef is no longer raw on the outside. Pour in the milk, and cook at a steady simmer, stirring occasionally, until the milk is absorbed. Add a dusting of nutmeg and the wine, and continue to cook until the wine is mostly absorbed. (These steps can take a bit longer than you think, maybe 10 to 15 minutes of so, but this isn&#39;t a sauce to be rushed.)",
						"Add the tomatoes in their juice and the tomato paste. Let come close to a boil then lower the heat to as low as possible until the sauce is just barely bubbling. Cook the sauce uncovered for 1 hour, stirring every once in a while, until most of the liquid is absorbed but the mixture is still wet, rich and thick. (There&#39;s a picture of it in the spoon, in the post.) You can use the sauce now, or if you want to give it some extra attention, you can cook it for another hour, adding up to 1/2 cup (125ml) of water, little by little as it continues to cook, and stirring occasionally, to make the sauce even more unctuous.",
						"Taste and add salt, if desired.",
						"To serve, cook the pasta in lightly salted water as directed on the package. (Before draining, reserve a little of the pasta cooking water.) Drain the pasta and toss the hot pasta in the Bolognese sauce with 1 tablespoon of butter. If the sauce needs a bit of thinning out, add a parsimonious splash of the reserved pasta water.",
					},
				},
				Name:  "Pasta Bolognese",
				Yield: models.Yield{Value: 4},
				URL:   "https://www.davidlebovitz.com/marcella-hazans-bolognese-sauce-recipe-italian-beef-tomato/",
			},
		},
		{
			name: "delish.com",
			in:   "https://www.delish.com/cooking/recipe-ideas/a24489879/beef-and-broccoli-recipe/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "weeknight meals"},
				CookTime:      "PT0S",
				Cuisine:       models.Cuisine{Value: "American"},
				DateModified:  "2023-06-26T17:34:00Z EST",
				DatePublished: "2018-11-06T17:45:06.218596Z EST",
				Description: models.Description{
					Value: "A classic Chinese-American dish with thinly sliced, velveted flank steak in a rich brown sauce with tender-crisp broccoli.",
				},
				Keywords: models.Keywords{
					Values: "American, Asian, dinner, weeknight meals, beef and broccoli recipe, sirloin recipes, Chinese take " +
						"out recipes, easy weeknight dinner recipes, stir-fry recipe, beef stir-fry",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"2 tbsp. dry sherry or shaoxing wine",
						"2 tbsp. unseasoned rice vinegar",
						"1/2 tsp. kosher salt",
						"1/2 tsp. freshly ground black pepper",
						"1/3 c. plus 1/4 c. low-sodium soy sauce, divided",
						"2 tbsp. plus 1 1/2 tsp. cornstarch, divided",
						"2 tbsp. light brown sugar, divided",
						"1 1/2 lb. flank or skirt steak, sliced very thin against the grain",
						"4 cloves garlic",
						`1 (1/2") piece ginger, peeled`,
						"2 scallions",
						"2 small heads broccoli",
						"2/3 c. low-sodium beef broth",
						"2 tbsp. oyster sauce",
						"2 tsp. sriracha (optional)",
						"3 tbsp. vegetable oil, divided",
						"Toasted sesame seeds andwhite rice, for serving",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"In a medium bowl, combine sherry, vinegar, salt, pepper, 1/3 cup soy sauce, 1 tablespoon plus 1 1/2 teaspoons cornstarch, and 1 tablespoon brown sugar. Add steak and toss to coat. Let sit 20 minutes.",
						"Meanwhile, finely chop garlic and ginger. Slice scallions and separate green and white parts. In a small bowl, combine garlic, ginger, and white scallion parts; reserve green parts for serving. Chop broccoli into florets and transfer to another small bowl.",
						"In a large measuring cup, whisk broth, oyster sauce, sriracha (if using), and remaining 1/4 cup soy sauce, 1 tablespoon cornstarch, and 1 tablespoon brown sugar. When ready to cook, arrange bowls of beef, garlic mixture, broccoli, and stir-fry sauce next to stove.",
						"In a large skillet over medium-high heat, heat 1 tablespoon oil. Add half of beef and cook, undisturbed, 1 minute, then stir and cook until cooked through and starting to char in some spots, about 1 minute more. Transfer to a plate. Repeat with 1 tablespoon oil and remaining beef. Discard excess marinade.",
						"Return skillet to medium heat and heat remaining 1 tablespoon oil. Add garlic mixture and cook, stirring occasionally, until fragrant, about 2 minutes. Add broccoli and cook, stirring frequently, until slightly softened, about 1 minute, then add stir-fry sauce. Cover and cook 3 minutes. Uncover, return beef to skillet, and toss to coat. Cook, tossing frequently, until warmed through and broccoli is crisp-tender, 2 to 3 minutes more.",
						"Divide rice among plates. Top with stir-fry, sesame seeds, and reserved green scallion parts.",
					},
				},
				Name: "Beef & Broccoli",
				NutritionSchema: models.NutritionSchema{
					Calories:      "609 Calories",
					Carbohydrates: "26 g",
					Cholesterol:   "111 mg",
					Fat:           "35 g",
					Fiber:         "7 g",
					Protein:       "46 g",
					SaturatedFat:  "10 g",
					Sodium:        "1918 mg",
					Sugar:         "9 g",
					TransFat:      "1 g",
				},
				PrepTime:  "PT10M",
				TotalTime: "PT40M",
				Yield:     models.Yield{Value: 4},
				URL:       "https://www.delish.com/cooking/recipe-ideas/a24489879/beef-and-broccoli-recipe/",
			},
		},
		{
			name: "ditchthecarbs.com",
			in:   "https://www.ditchthecarbs.com/how-to-make-keto-samosa-air-fryer-oven/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Appetiser"},
				CookTime:      "PT10M",
				Cuisine:       models.Cuisine{Value: "Egg free"},
				DatePublished: "2022-03-17T09:35:35+00:00",
				Description: models.Description{
					Value: "Keto samosas is an Indian vegetarian dish perfect for appitizers, snacks, or even a meal.",
				},
				Keywords: models.Keywords{Values: "Keto Samosas"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 batch keto roti dough",
						"2 cups cauliflower (steamed and chopped)",
						"2 tbsp extra virgin olive oil ((ghee or coconut oil))",
						"½ tsp ground cumin",
						"½ tsp ground ginger",
						"1 tbsp garam masala",
						"2 cloves garlic (minced)",
						"2 tbsp jalapeno (chopped)",
						"2 stems green onions (chopped)",
						"1 tbsp lemon juice",
						"2 tbsp cilantro (chopped)",
						"+/- salt and pepper",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Start by maxing up a batch of keto roti dough. While the dough is resting you will want to make the filling.",
						"To make the filling heat the oil in a skillet. When the oil is hot add the seasoning(cumin, ginger, " +
							"Garam masala), garlic, and jalapeno to the skillet. Saute in the skillet for about 1 minute then " +
							"add the rest of the filling ingredients to the skillet. Let the filling cook for 3-4 minutes, stirring " +
							"occassionally. Once the filling is done remove it from the heat and let it cool completely.",
						"While the filling is cooling prepare and shape the samosas wrappers. To do this cut your roti dough " +
							"into 6 equal sections. Form each section into a ball. Roll each ball into a thin circle about the " +
							"size of a side plate between two sheets of parchment paper using a rolling pin.Then using a knife cut " +
							"each circle in half. This will give you 12 samosa wrappers. Cover and set aside until filling is " +
							"completely cool.",
						"Start by laying out one samosa wrapper with the flat cut side toward the bottom. Use your finger to " +
							"brush water in a line along half of the flat cut side. This will help the wrapper stick together " +
							"as you fold it.",
						"Next, fold the two corner edges up to make a cone. Overlap the corners with the wet edge on top of " +
							"the other. Firmly press the seam down to ensure it sticks and pinch the bottom of the cone closed.",
						"Hold the cone upright and open in your hand with the seam facing your palm. Your palm will help support " +
							"the seam while you add the filling. Spoon the filling into the cone until it is 2/3 of the way full.",
						"Using your finger brush water around the inside edge of the cone. Then press the edges of the cone " +
							"together and fold the bottom under to seal the cone and form a triangular samosa. Place the samosa " +
							"down so that it is standing up with the sealed bottom facing down. Repeast until all 12 samosas are filled.",
						"Finally, brush the samosas with olive oil and cook one of 3 ways.1) Bake in the oven on a baking tray " +
							"at 190°C/375°F for 15-18 minutes or until the samosas are golden brown.2) Arrange the samosas " +
							"in your air fryer basket and air at 190°C/375°F for 7-8 minutes. Depending on the size of your air fryer " +
							"you may need to cook the samosas in batches.3) Add an inch of olive oil to a skillet and fry the " +
							"samosas over medium-high heat for 3-4 minutes on each side. Only cook 3-4 samosas at a time. Let " +
							"the samosas cool for a few minutes after cooking and enjoy.",
					},
				},
				Name: "Easy Keto Samosas Recipe (Air Fryer Recipe)",
				NutritionSchema: models.NutritionSchema{
					Calories:      "69.6 kcal",
					Carbohydrates: "7.5 g",
					Fat:           "4.5 g",
					Fiber:         "3.8 g",
					Protein:       "1.2 g",
					SaturatedFat:  "",
					Servings:      "1",
					Sodium:        "30.7 mg",
					Sugar:         "0.9 g",
				},
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 12},
				URL:      "https://www.ditchthecarbs.com/how-to-make-keto-samosa-air-fryer-oven/",
			},
		},
		{
			name: "domesticate-me.com",
			in:   "https://domesticate-me.com/10-summer-cocktail-recipes/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Drinks"},
				DatePublished: "2021-05-28T16:11:33+00:00",
				Description: models.Description{
					Value: "Made with muddled strawberries, thyme, lemon, vodka, and St. Germain, this refreshing Strawberry " +
						"Thyme Cooler is perfect for all your summer celebrations!",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1/3 cup hulled and diced strawberries",
						"1 sprig of fresh thyme",
						"½ ounce fresh lemon juice",
						"1 ounce St. Germain",
						"2 ounces vodka",
						"2-3 ounces club soda (depending on personal taste)",
						"Sliced strawberries",
						"Thyme sprig",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Muddle the strawberries, thyme sprig, and lemon juice in a cocktail shaker. Add the St. Germain, " +
							"vodka, and some ice, and shake vigorously to combine.",
						"Strain the cocktail into a glass with ice and top with club soda.",
						"Garnish with sliced strawberries and a sprig of thyme. Bottoms up!",
					},
				},
				Name:     "Strawberry Thyme Cooler and 9 Other Summer Cocktail Recipes",
				PrepTime: "PT5M",
				Yield:    models.Yield{Value: 1},
				URL:      "https://domesticate-me.com/10-summer-cocktail-recipes/",
			},
		},
		{
			name: "downshiftology.com",
			in:   "https://downshiftology.com/recipes/baked-chicken-breasts/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Main Course"},
				CookTime:      "PT20M",
				Cuisine:       models.Cuisine{Value: "American"},
				DatePublished: "2023-10-31T08:00:00+00:00",
				Description: models.Description{
					Value: "You&#039;ll love these perfectly baked chicken breasts that are juicy, tender, easy, and flavorful! All you need is a drizzle of olive oil, a special seasoning mix, and a few insider tips for these super tasty, no-fail chicken breasts. Watch the video below to see how I make them in my kitchen!",
				},
				Keywords: models.Keywords{
					Values: "baked chicken, baked chicken breasts, baked chicken recipe",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"4 boneless skinless chicken breasts",
						"1 tablespoon olive oil (or avocado oil)",
						"1 teaspoon kosher salt",
						"1 teaspoon paprika",
						"½ teaspoon garlic powder",
						"½ teaspoon dried thyme (or oregano, parsley or other herbs)",
						"¼ teaspoon freshly ground black pepper",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Preheat your oven to 425°F (220°C). In a small bowl, mix together the paprika, garlic powder, thyme, salt and pepper.",
						"Lightly coat the chicken breasts in olive oil, and then generously rub the spice mix on both sides of the chicken.",
						"Place the chicken breasts in a baking dish and cook for 20 to 25 minutes, depending on size (see chart " +
							"above). Let the chicken rest for a few minutes to allow the juices to redistribute within " +
							"the meat, then serve.",
					},
				},
				Name: "Best Baked Chicken Breast",
				NutritionSchema: models.NutritionSchema{
					Calories:      "163 kcal",
					Carbohydrates: "1 g",
					Cholesterol:   "72 mg",
					Fat:           "7 g",
					Fiber:         "1 g",
					Protein:       "24 g",
					SaturatedFat:  "1 g",
					Servings:      "1",
					Sodium:        "713 mg",
					Sugar:         "1 g",
				},
				PrepTime:  "PT5M",
				TotalTime: "PT25M",
				Yield:     models.Yield{Value: 4},
				URL:       "https://downshiftology.com/recipes/baked-chicken-breasts/",
			},
		},
		{
			name: "dr.dk",
			in:   "https://www.dr.dk/mad/opskrift/nytarskage-med-champagne-kransekagebund-solbaer-og-chokoladepynt",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        models.SchemaType{Value: "Recipe"},
				Name:          "Nytårskage med champagne, kransekagebund, solbær og chokoladepynt",
				DatePublished: "2022-01-01T21:00:00+00:00",
				Image:         models.Image{Value: anUploadedImage.String()},
				Description: models.Description{
					Value: "Smuk nytårskage med urvisere og masser af smag, der passer perfekt til nytårsaften.",
				},
				URL: "https://www.dr.dk/mad/opskrift/nytarskage-med-champagne-kransekagebund-solbaer-og-chokoladepynt",
				Ingredients: models.Ingredients{
					Values: []string{
						"30 g æggehvider",
						"50 g sukker",
						"Fintrevet citronskal",
						"200 g bagemarcipan",
						"15 g sukker",
						"15 g ristede, hakkede smuttede mandler",
						"1 nip salt",
						"50 g mørk chokolade",
						"1 spsk. frysetørret solbærpulver",
						"1 tsk. pufsukker",
						"140 g hvid chokolade",
						"1/2 blad husblas",
						"2 æggeblommer",
						"15 g sukker",
						"80 ml piskefløde",
						"80 g solbærpuré",
						"2 tsk. citronsaft",
						"1/2 tsk. citronsyre",
						"75 g sukker",
						"40 g glukosesirup",
						"30 ml tør champagne",
						"1/2 tsk. citronsyre",
						"50 g æggehvider",
						"1 spsk. sukker",
						"2 tsk. solbærpulver",
						"2 blade husblas",
						"400 g gold chokolade",
						"160 ml tør champagne",
						"40 ml citronsaft",
						"4 dl piskefløde",
						"1 tsk. solbærpulver",
						"5 blade husblas",
						"150 g hvid chokolade",
						"100 g kondenseret mælk",
						"51 g solbærpuré",
						"54 g vand",
						"150 g sukker",
						"100 g glukosesirup",
						"Minimum 800 g mørk chokolade",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Kransekagebund",
						"Varm marcipan kort i mikroovn, rør alle ingredienser sammen i røremaskinen indtil den samler sig.",
						"Fordel dejen i en kagering på ca. 16ø og bag den I ca. 15 minutter ved 200C, til den er let gylden. Lad den køle ned.",
						"\n",
						"Mandel-solbær-knas",
						"Smelt sukkeret ved lav blus i en gryde, indtil den er let gylden. Tag gryden af varme og rør mandler og salt i.",
						"Hæld karamelliseret mandler i et stk bagepapir og lad den køle helt ned. Hak mandlerne fint. Smelt chokoladen op til 45C, tilsæt solbærpulver, pufsukker og de hakket mandler i.",
						"Fordel knaset oven på kransekagebunden til et jævnt fint lag. Stil kagen i fryseren mens du arbejder videre.",
						"\n",
						"Solbærcremeux",
						"Kom chokoladen i en skål og sæt den til side. Sæt husblas i koldt vand i ca. 5 minutter.",
						"Pisk æggeblommer og sukker let sammen i en anden skål.",
						"Varm fløde, citronsyre og solbærpuré op til kogepunktet og hæld den i æggeblandingen under piskning. Hæld blandingen tilbage i gryden og kog cremen op til 85C under konstant omrøring med en silikoneske.",
						"Tag gryden af varmen og sigt cremen.",
						"Vrid husblas fri fra vand og smelt den i cremen. Hæld cremen over chokoladen og rør midt i skålen indtil chokoladen og cremen er blevet homogen.",
						"Smag til med citronsaft og stablen cremeuxen. Fordel den ovenpå kransekagebund (chokolade knas skal være ned i bunden) i bageringen og sæt den i fryseren.",
						"\n",
						"Champagne-flødebolleskum",
						"Bring 75g sukker, glukosesirup, champagne og citronsyre i kog op til 118C.",
						"Pisk æggehvider og 1 spsk sukker næsten stive i en skål og hæld den varme sukkerlage ned i æggehviderne i en tynd stråle under piskning. Pisk videre, til skummet er sejt og fast. Ved slutning tilsæt solbærpulver og pisk færdig.",
						"Fordel skummet over solbærcremeux, gør overfladen glat med en paletkniv og stil kagen i fryseren igen.",
						"\n",
						"Champagnemousse",
						"1/2 tsk citronsyre Udblød husblassen i koldt vand. Smelt chokoladen op til 45C.",
						"Bring champagne, citronsyre og citronsaft til kogepunktet og tag gryden straks af varmen.",
						"Vrid husblas fri for vand og rør den i den varme champagne. Hæld champagnen over chokoladen, mens du rør " +
							"i midten af skålen. Fortsæt indtil massen samler sig.",
						"Pisk fløde til let skum, ved skummet over chokoladen ad 3 omgange.",
						"Fordel champagnemousse i en silikoneforme og tryk solbærindlæg ned i moussen - bunden skal være oppe. Lad " +
							"kagen fryse indtil den skal glazes.",
						"\n",
						"Solbærglaze",
						"Læg husblad i koldt vandbad i ca. 5 minutter. Kom chokolade i en høj kande.",
						"Bring solbærpuré, vand, sukker og glukosesirup i en gryde og koge den op til kogepunktet. Tag gryden af " +
							"varmen og rør kondenseret mælk i.",
						"Vrid husblassen fri for vand og rør den ud i den varme væske. Hæld væsken over chokoladen og lad den " +
							"træk i et par minutter.",
						"Stavblend glaze indtil det er samlet og ensartet. Tilsæt guldstøv og stavblend igen.",
						"Dæk overflade med film og lad gazen køle ned til 32 grader.",
						"Placér den frosne kage på en rist og glaze kagen.",
						"\n",
						"Forberedelse af chokoladepynt",
						"Klip et langt stykke plast på 63 cm til som et bylandskab og et andet stykke plast på 30 cm",
						"Klip to viser af stift kageplast",
						"Forberedelse af en halvkugle Ø10 cm pudses med vat",
						"Temperere guld kakaofarve sammen med guldstøv og dup med en svamp et fyrværkeri mønster i halvkuglen, så " +
							"det ligner en guldregn",
						"Temperere 800 g mørk chokolade",
						"Fyld hele halvkuglen op med chokolade og lade det sidde i 1-2 min for at sørge for, at skallen bliver tyk",
						"Bank alt det overskydende chokolade ud på et stykke bagepapir og sæt den til at størkne",
						"Fordel et jævnt lag chokolade ud på plasten der forestiller et bylandskab og på 30 cm plast",
						"Beklæd en 20 cm kagering med plast og sæt bylandskabet rundt herom, som et stort bånd",
						"Udstik med en Ø3 cm ring små knapper af den næsten størknet chokolade på 30 cm plast og befri dem, når " +
							"chokoladen er kølet helt af",
						"De runde cirkler pensles med guldstøv og efterfølgende skrives der med romertal fra 1-12",
						"Befri halvkuglen fra formen og sæt til side",
						"Lav visere af chokolade",
						"\n",
						"Samling",
						"Den frosne kage befries fra formen og sættes på en drejeskive",
						"Kagen betrækkes med et lag smørcreme, som glattes tyndt ud",
						"Kagen sættes tilbage på frost i 10 min",
						"Betræk med et nyt lag, som glattes helt jævnt ud",
						"Sæt kagen tilbage på frost 10-15 min",
						"Spray kagen med cremet velvet spray farve",
						"Betræk kagen med chokoladebånd",
						"Pynt toppen med guldtal oven på små chokoladekugler og placere halvkuglen i midten af kagen",
					},
				},
			},
		},
		{
			name: "drinkoteket.se",
			in:   "https://drinkoteket.se/recept/limoncello-spritz/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Champagnedrinkar"},
				CookTime:      "PT2M",
				DatePublished: "2023-12-30",
				Description: models.Description{
					Value: "Limoncello Spritz Drinkrecept på Drinkoteket.se. Här hittar du en mängd recept på enkla och goda drinkar och cocktails online. Välkommen in!",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"5 cl Limoncello", "10 cl Bubbel", "4 cl Sodavatten", "1 cl Citronjuice",
						"Rekommenderat bubbel", "Prosecco",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Fyll ett glas med is.",
						"Tillsätt ett torrt mousserande vin, förslagsvis prosecco.",
						"Tillsätt limoncello och en skvätt citronjuice.", "Rör om lite lätt.",
						"Garnera med en citronskiva.",
					},
				},
				Name:     "Limoncello Spritz",
				PrepTime: "PT1M",
				Tools: models.Tools{
					Values: []string{"Barset med shaker", "Shaker", "Jigger", "Citruspress", "Cocktailsil", "Barsked"},
				},
				Yield: models.Yield{Value: 1},
				URL:   "https://drinkoteket.se/recept/limoncello-spritz/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
