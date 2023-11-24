package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_0to9(t *testing.T) {
	testcases := []testcase{
		{

			name: "101cookbooks.com",
			in:   "https://www.101cookbooks.com/simple-bruschetta/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "Simple Bruschetta",
				Description: models.Description{
					Value: "Good tomatoes are the thing that matters most when it comes to making bruschetta - the classic " +
						"Italian antipasto. It is such a simple preparation that paying attention to the little details matters.",
				},
				Category: models.Category{Value: "Appetizer"},
				Cuisine:  models.Cuisine{Value: "Easy"},
				Image: models.Image{
					Value: "https://images.101cookbooks.com/bruschetta-recipe-h1.jpg?w=1200&auto=format",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"3 fresh tomatoes, ripe",
						"A small handful of basil leaves",
						"1 teaspoon good-tasting white wine vinegar (or red/balsamic), or to taste",
						"1/4 teaspoon sea salt, or to taste",
						"4 tablespoons extra-virgin olive oil, plus more for serving",
						"3 - 4 sourdough or country-style bread slices (at least 1/2-inch thick)",
						"2 cloves garlic, peeled",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Rinse and dry your tomatoes. Halve each of them, use a finger to remove the seeds, and cut out the " +
							"cores. Roughly cut the tomatoes into 1/2-inch pieces and place in a medium bowl. " +
							"Tear the basil into small pieces, and add that as well. Add 2 tablespoons of the olive oil, " +
							"a small splash of vinegar, and a pinch of salt. Gently toss, taste, adjust if necessary, and set aside.",
						"Heat a grill or oven to medium-high. When it’s ready, use the remaining 2 tablespoons of the olive oil to " +
							"brush across the slices of bread. Grill or bake until well-toasted and golden brown with a hint of char. " +
							"Flipping when the first side is done. Remove from grilled and when cool enough to handle, " +
							"rub both sides of each slice of bread with garlic.",
						"Cut each slice of bread in half if you like, and top each segment with the tomato mixture. And a " +
							"finishing drizzle of olive oil is always nice.",
					},
				},
				NutritionSchema: models.NutritionSchema{
					Calories:      "233 kcal",
					Carbohydrates: "19 g",
					Cholesterol:   "",
					Fat:           "15 g",
					Protein:       "6 g",
					Sodium:        "287 mg",
					Servings:      "1",
				},
				Keywords:      models.Keywords{Values: "simple"},
				Yield:         models.Yield{Value: 4},
				PrepTime:      "PT15M",
				CookTime:      "PT5M",
				DatePublished: "2022-06-29T14:40:39+00:00",
				URL:           "https://www.101cookbooks.com/simple-bruschetta/",
			},
		},
		{
			name: "750g.com",
			in:   "https://www.750g.com/carbonara-vegetarienne-r202631.htm",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				Category:      models.Category{Value: "Plats"},
				CookTime:      "PT10M",
				DatePublished: "2017-06-28",
				Description: models.Description{
					Value: "Une variante v&amp;eacute;g&amp;eacute;tarienne de la Carbonara dans laquelle&amp;nbsp;les&amp;nbsp;lardons sont remplac&amp;eacute;s par des l&amp;eacute;gumes.",
				},
				Keywords: models.Keywords{
					Values: "carbonara vegetarienne,carbonara aux légumes,pâtes carbonara sans lardons,spaghetti à la carbonara,pâtes à la carbonara,pâtes aux légumes,spaghetti,pâtes,carottes,courgettes,petits pois,plats végétariens,cuisine estivale,cuisine végétarienne,750green",
				},
				Image: models.Image{
					Value: "https://static.750g.com/images/1200-675/4133b19abe72214a94b66780b61f3973/327789.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"500 g de spaghetti bio Barilla", "4 carottes", "2 courgettes", "2 oeufs",
						"1 oignon", "80 g de parmesan", "1 poignée de petits pois frais",
						"4 c. à s. d'huile d'olive", "4 feuilles de basilic frais", "Poivre",
						"Sel ou sel fin",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Disposez tous les ingr&amp;eacute;dients sur un plateau.",
						"&amp;Eacute;cossez les petits pois.&amp;nbsp;&amp;Eacute;pluchez les carottes et coupez-les en julienne.&amp;nbsp;Coupez les courgettes en julienne.&amp;nbsp;",
						"Coupez l&#039;oignon en lamelles fines.&amp;nbsp;",
						"Chauffez l&#039;huile d&#039;olive dans une po&amp;ecirc;le puis faites revenir l&#039;oignon, les carottes et les courgettes pendant 5 minutes. Ajoutez les petits pois, salez le tout et laissez rissoler pendant encore 5 minutes.&amp;nbsp;",
						"Versez le parmesan et les &amp;oelig;ufs dans un saladier, ajoutez une pinc&amp;eacute;e de sel et du poivre du moulin.",
						"Fouettez et r&amp;eacute;servez.",
						"Portez &amp;agrave; &amp;eacute;bullition un grand volume d&#039;eau dans une casserole &amp;agrave; bords hauts.&amp;nbsp;D&amp;egrave;s l&amp;rsquo;&amp;eacute;bullition, ajoutez 20 g de gros sel, ajoutez les&amp;nbsp;p&amp;acirc;tes et remuez une fois avec une cuill&amp;egrave;re en bois. Faites-les cuire selon le temps de cuisson indiqu&amp;eacute; sur le paquet.&amp;nbsp;&amp;nbsp;",
						`&amp;Eacute;gouttez les spaghetti "al dente" dans une passoire, puis versez-les dans le saladier avec la sauce aux &amp;oelig;ufs et parmesan. Remuez.&amp;nbsp;`,
						"Ajoutez&amp;nbsp;les l&amp;eacute;gumes, remuez bien, ajoutez les feuilles de basilic cisel&amp;eacute;es&amp;nbsp;et servez imm&amp;eacute;diatement, bien chaud.",
						"Bonne d&amp;eacute;gustation !",
					},
				},
				Name:     "Carbonara végétarienne",
				PrepTime: "PT30M",
				Yield:    models.Yield{Value: 4},
				URL:      "https://www.750g.com/carbonara-vegetarienne-r202631.htm",
			},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
