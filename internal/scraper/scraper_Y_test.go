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
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				Description:   &models.Description{},
				Keywords:      &models.Keywords{},
				Name:          "Walnut Turkish Baklava Recipe",
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
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
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Firstly, prepare your dessert syrup."},
						{Type: "HowToStep", Text: "For the syrup: Put water and sugar in a saucepan, stir until sugar is dissolved. Turn down the bottom after boiling syrup. Boil for 20 minutes. Add lemon juice into the syrup. Boil over low heat for 10 minutes."},
						{Type: "HowToStep", Text: "Remove syrup from heat, leave to cool. Then, melt the butter in a small pan."},
						{Type: "HowToStep", Text: "On the other hand, take a filo pastry. Put on the kitchen counter. Apply melted butter with a brush onto filo pastry. Sprinkle with crushed walnuts on to filo pastry. Slowly wrap in roll form. Then, follow the insructions at photo 9-10-11-12.  Put into a greased baking tray. Repeat the same process."},
						{Type: "HowToStep", Text: "Heat the remaining butter. Slowly pour butter onto baklava."},
						{Type: "HowToStep", Text: "Give a preheated 180 degree oven. Cook until golden brown (about 30 minutes)."},
						{Type: "HowToStep", Text: "Remove cooked baklava from the oven. Leave to cool for 8-10 minutes. Then, cut the baklava with a knife or spatula."},
						{Type: "HowToStep", Text: "Finally, pour the cool syrup onto baklava. Rest for 20-25 minutes. Then, service."},
					},
				},
				NutritionSchema: &models.NutritionSchema{},
				ThumbnailURL:    &models.ThumbnailURL{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				Yield:           &models.Yield{Value: 1},
				URL:             "https://ye-mek.net/recipe/walnut-turkish-baklava-recipe",
				Video:           &models.Videos{},
			},
		},
		{
			name: "yumelise.fr",
			in:   "https://www.yumelise.fr/crepes-jambon/",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Plat principal"},
				CookingMethod: &models.CookingMethod{},
				CookTime:      "PT40M",
				Cuisine:       &models.Cuisine{Value: "Française"},
				Description: &models.Description{
					Value: "Des crêpes moelleuses garnies de béchamel au fromage et de jambon puis cuites à la poêle pour apprécier le croustillant.",
				},
				Keywords: &models.Keywords{Values: "jambon,lait entier,parmesan"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"300 grammes de farine fluide T45", "3 oeufs",
						"2 cuillères à soupe d'huile tournesol ici",
						"60 centilitres de lait entier 600g",
						"1 cuillère à soupe de sucre en poudre blanc", "70 grammes de beurre doux",
						"70 grammes de farine", "50 centilitres de lait entier 500g",
						"100 grammes de parmesan fraîchement râpé ou comté, emmental, etc.", "sel",
						"poivre", "muscade moulue facultatif",
						"8 tranches de jambon cuit, doré, sans couenne, 300g en tout ici",
						"beurre demi-sel, pour la cuisson des crêpes au jambon",
						"huile d'olive pour la cuisson des crêpes au jambon",
						"huile de tournesol ici, pour la cuisson des crêpes natures",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Dans un saladier, battre les oeufs en omelette, avec un fouet. Puis ajouter 1/3 de la farine (bien fouetter de l'intérieur vers l'extérieur pour éviter les grumeaux). Ajouter progressivement un deuxième tiers de farine. Puis verser la moitié du lait pour détendre la pâte, toujours de l'intérieur vers l'extérieur. Ajouter les reste de farine, fouetter, et enfin le reste de lait."},
						{Type: "HowToStep", Text: "Ajouter le sucre et l'huile. Bien fouetter. Le sucre, dans ces crêpes salées, apportera du moelleux qui ira très bien avec le jambon et la béchamel au fromage. Cyril Lignac en met 3 cuillères mais une seule suffit pour moi."},
						{Type: "HowToStep", Text: "Laisser reposer la pâte 10 minutes, le temps de préparer la béchamel au fromage."},
						{Type: "HowToStep", Text: "Dans une grande casserole, sur feu moyen, faire fondre le beurre, ajouter la farine, mélanger au fouet et laisser cuire le roux quelques minutes, 4 minutes environ (le laisser cuire évite d'avoir un goût farineux). Verser le lait (en une seule fois), mélanger sur feu moyen, puis cuire jusqu’à ce que la sauce épaississe. Elle doit être bien épaisse. Retirer du feu, saler légèrement (car il y aura du fromage), poivrer, saupoudrer d'un peu de muscade (facultatif), fouetter. Ajouter le fromage râpé. Bien fouetter."},
						{Type: "HowToStep", Text: "Réserver le temps de cuire les crêpes."},
						{Type: "HowToStep", Text: "Chauffer une crêpière. Huiler avec un papier absorbant."},
						{Type: "HowToStep", Text: "Verser une louche de pâte (il ne faut pas une crêpe trop fine, car elles sont meilleures et bien moelleuses plus épaisses, je prends une grande louche pleine pour une crêpière de 28 cm de diamètre), laisser cuire jusqu’à ce que la crêpe soit légèrement dorée, la retourner. Cuire l'autre face (ne laisser que quelques secondes sur la seconde face : le fait de pas trop cuire la seconde face facilitera la “soudure” des crêpes une fois qu’elles seront fourrées) et réserver dans une assiette.Des crêpes trop cuites seraient sèches et moins moelleuses."},
						{Type: "HowToStep", Text: "Cuire toutes les crêpes de la même façon (j'en fais 8)."},
						{Type: "HowToStep", Text: "Déposer une crêpe sur le plan de travail propre, verser au centre une bonne cuillère à soupe de béchamel (la fouetter un peu avant pour l'assouplir) et répartir avec une petite spatule en cercle en épargnant quelques centimètres du bord."},
						{Type: "HowToStep", Text: "Déposer au dessus de la béchamel une tranche de jambon déchirée en morceaux de manière à couvrir le cercle de béchamel."},
						{Type: "HowToStep", Text: "Commencer à rouler 1/3 de la crêpe, puis plier les bords et terminer de rouler."},
						{Type: "HowToStep", Text: "Faire autant de crêpes jambon béchamel souhaitées."},
						{Type: "HowToStep", Text: "Une fois les crêpes roulées, les déposer dans un plat et couvrir d'un film alimentaire. Si elles sont préparées à l'avance, mais cuites ensuite dans la journée, les conserver à température ambiante. Si elles sont cuites le lendemain, les réserver au réfrigérateur mais les mettre quelques heures à température ambiante avant cuisson. On peut également les emballer d'un film alimentaire et les congeler."},
						{Type: "HowToStep", Text: "Dans une grande poêle, chauffer à feu moyen 25g de beurre demi-sel et une cuillère à soupe d’huile d’olive. Déposer les crêpes côté jointure en premier (je peux en mettre 4 dans la poêle), commencer à colorer (2 min), retourner. Il faut colorer tous les côtés (2 min à chaque face), pour que les crêpes au jambon soit bien colorées et croustillantes. Terminer la cuisson par la jointure. Il faut compter 10 minutes de cuisson à feu moyen afin que les crêpes soient bien chaudes à l'intérieur.(Dans une secondes poêle moyenne j'en cuis deux dans 12g de beurre et 1/2 cuil à soupe d'huile d'olive)"},
						{Type: "HowToStep", Text: "Servir de suite."},
					},
				},
				Name:            "Crêpes au jambon",
				NutritionSchema: &models.NutritionSchema{},
				PrepTime:        "PT20M",
				ThumbnailURL:    &models.ThumbnailURL{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				TotalTime:       "",
				Yield:           &models.Yield{},
				URL:             "https://www.yumelise.fr/crepes-jambon/",
				Video:           &models.Videos{},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
