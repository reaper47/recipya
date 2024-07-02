package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_Q(t *testing.T) {
	testcases := []testcase{
		{
			name: "quitoque.fr",
			in:   "https://www.quitoque.fr/recette/12976/saumon-teriyaki-et-riz-rouge",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				CookTime:      "PT20M",
				Description: &models.Description{
					Value: "Le teriyaki, au Japon, désigne un plat poêlé ou grillé et caramélisé dans son assaisonnement. Ici, on l'utilise en sauce douce et enveloppante !",
				},
				Keywords: &models.Keywords{},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"240 g pavé de saumon 2x120g sous vide", "1 poireau pce rdc",
						"150 g INACTIF - riz basmati blanc BIO 150g lighties",
						"20 mL sauce teriyaki 20ml",
						"huile de cuisson (mélange 4 huiles par exemple)", "sel", "poivre",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "La marinade\nRincez et essuyez les pavés de saumon. Poivrez-les.\r\nDéposez les pavés de saumon et enrobez-les avec la sauce teriyaki.\r\nRéservez-les au réfrigérateur 10 min environ.\r\nPendant ce temps, préparez le poireau."},
						{Type: "HowToStep", Text: "Le poireau\nCoupez et retirez la base abîmée du poireau. Incisez-le en deux dans la longueur et rincez-le soigneusement. Emincez-le.\r\nDans une sauteuse, faites chauffer un filet d'huile de cuisson à feu moyen à vif. \r\nFaites revenir le poireau 15 min environ. Salez, poivrez.\r\nA mi-cuisson, ajoutez un fond d'eau et couvrez pour accélérer la cuisson.\r\nEn parallèle, lancez la cuisson du riz."},
						{Type: "HowToStep", Text: "Le riz\nPortez à ébullition une casserole d'eau salée. \r\nFaites cuire le riz selon les indications du paquet.\r\nPendant la cuisson du riz, faites cuire le saumon."},
						{Type: "HowToStep", Text: "Le saumon\nDans une poêle, faites chauffer un filet d’huile de cuisson à feu moyen à vif.\r\nFaites dorer les pavés de saumon avec la sauce (côté peau en premier) 2 à 3 min de chaque côté jusqu'à ce qu'ils soient cuits à coeur. Pour savoir quand les retourner, regardez sur le côté : lorsqu’ils sont cuits à mi-hauteur, retournez-les !\r\nPenchez la sauteuse et à l'aide d'une cuillère, récupérez la sauce pour la verser sur le saumon et bien l'enrober."},
						{Type: "HowToStep", Text: "Dégustez sans attendre votre saumon nappé de sauce teriyaki accompagné de la tombée de poireau et du riz !"},
					},
				},
				Name: "Saumon teriyaki et riz au poireau",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "617.85 kcal",
					Carbohydrates: "70.12 g",
					Fat:           "20.49 g",
					Fiber:         "4.76 g",
					Protein:       "34.94 g",
					SaturatedFat:  "2.81 g",
					Servings:      "2",
					Sodium:        "0.56 g",
					Sugar:         "10.45 g",
				},
				PrepTime:     "PT20M",
				ThumbnailURL: &models.ThumbnailURL{},
				Tools:        &models.Tools{Values: []models.HowToItem{}},
				Yield:        &models.Yield{Value: 2},
				URL:          "https://www.quitoque.fr/recette/12976/saumon-teriyaki-et-riz-rouge",
				Video:        &models.Videos{},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
