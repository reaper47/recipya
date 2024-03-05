package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_Y(t *testing.T) {
	testcases := []testcase{
		{
			name: "ye-mek.net",
			in:   "https://ye-mek.net/recipe/walnut-turkish-baklava-recipe",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "Walnut Turkish Baklava Recipe",
				Image: models.Image{
					Value: "https://cdn.ye-mek.com/img/f/hazir-yufkadan-buzme-burma-baklava.jpg",
				},
				Ingredients: models.Ingredients{
					Values: []string{
						"11-12 sheets filo pastry",
						"170 g butter",
						"Crushed walnut",
						"For Syrup:",
						"2 cups water",
						"2 cups granulated sugar",
						"Quarter lemon juice",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Firstly, prepare your dessert syrup.",
						"For the syrup: Put water and sugar in a saucepan, stir until sugar is dissolved. Turn down the bottom after boiling syrup. Boil for 20 minutes. Add lemon juice into the syrup. Boil over low heat for 10 minutes.",
						"Remove syrup from heat, leave to cool. Then, melt the butter in a small pan.",
						"On the other hand, take a filo pastry. Put on the kitchen counter. Apply melted butter with a brush onto filo pastry. Sprinkle with crushed walnuts on to filo pastry. Slowly wrap in roll form. Then, follow the insructions at photo 9-10-11-12.  Put into a greased baking tray. Repeat the same process.",
						"Heat the remaining butter. Slowly pour butter onto baklava.",
						"Give a preheated 180 degree oven. Cook until golden brown (about 30 minutes).",
						"Remove cooked baklava from the oven. Leave to cool for 8-10 minutes. Then, cut the baklava with a knife or spatula.",
						"Finally, pour the cool syrup onto baklava. Rest for 20-25 minutes. Then, service.",
					},
				},
				URL: "https://ye-mek.net/recipe/walnut-turkish-baklava-recipe",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
