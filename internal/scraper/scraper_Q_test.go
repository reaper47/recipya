package scraper_test

import (
	"testing"

	"github.com/reaper47/recipya/internal/models"
)

func TestScraper_Q(t *testing.T) {
	testcases := []testcase{
		{
			name: "quitoque.fr",
			in:   "https://www.quitoque.fr/recettes/betterave-de-couleur-et-carotte-roties-au-miel-feta-et-pesto-de-sauge",
			want: models.RecipeSchema{
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				CookTime:      "PT15M",
				Description: &models.Description{
					Value: "Cette recette est pleine de goûts et de couleurs : betteraves, carotte, féta, miel, sauge... De quoi apporter du soleil dans vos assiettes !",
				},
				Keywords: &models.Keywords{Values: "Découverte,Végétarien"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1x betterave rouge chioggia (250g)",
						"1x betterave crue (280g)",
						"2x carotte (200g)",
						"20 ml miel Bio",
						"2x pomme de terre (440g)",
						"0.5 sauge (botte)",
						"15 g parmigiano reggiano râpé AOP",
						"0.5 gousse d'ail",
						"80 g féta AOP",
						"25 g amandes crues",
						"3 càs huile d'olive",
						"3 càs eau",
						"sel",
						"poivre",
						"beurre",
						"lait",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Étape 1\n1. Les légumes\n\n    Préchauffez votre four à 200°C en chaleur tournante !\n    Pendant ce temps, épluchez les betteraves et les carottes.\n    Coupez les betteraves en demi-rondelles et les carottes en frites larges.\n    Déposez les légumes sur une plaque allant au four. Versez dessus le miel et un bon filet d'huile d'olive. Salez, poivrez. Mélangez bien pour tout enrober.\n    Enfournez 25 à 30 min. Remuez régulièrement.\n    Pendant ce temps, préparez la purée.\n\n\n    \n        Voir toute la recette"},
						{Type: "HowToStep", Text: "Étape 2\n2. La purée\n\n    Épluchez et coupez les pommes de terre en morceaux.\n    Déposez-les dans une casserole et recouvrez-les très largement d'eau. Salez.\n    Faites-les cuire 15 à 20 min.\n    Egouttez et réduisez les pommes de terre en purée. Salez, poivrez.\n    Ajoutez une noix de beurre et un filet de lait pour une purée plus onctueuse.\n    En parallèle, préparez le pesto de sauge."},
						{Type: "HowToStep", Text: "Étape 3\n3. Le pesto de sauge\n\n    Dans le bol d'un mixeur, déposez les ingrédients suivants :\n    Effeuillez la sauge.\n    Pelez l'ail.\n    Ajoutez le parmesan, la moitié des amandes, l'huile et l'eau. Salez légèrement, poivrez.\n    Mixez jusqu'à obtenir une texture homogène.\n    Si nécessaire, ajoutez de l'eau pour un pesto plus liquide.\n    Goûtez et rectifiez l'assaisonnement si nécessaire."},
						{Type: "HowToStep", Text: "Étape 4\n4. A table !\n\n    Hachez grossièrement le reste des amandes.\n    Émiettez la féta.\n    Dans des assiettes, déposez au fond la purée. Répartissez dessus les légumes confits. Nappez le tout de pesto de sauge et parsemez d'amandes et de féta."},
					},
				},
				Name: "Betterave de couleur et carotte rôties au miel, féta et pesto de sauge",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "84",
					Carbohydrates: "10.54",
					Cholesterol:   "",
					Fat:           "2.85",
					Fiber:         "2.14",
					Protein:       "3.27",
					SaturatedFat:  "1.27",
					Sodium:        "0.31",
					Sugar:         "3.76",
				},
				PrepTime:     "PT25M",
				ThumbnailURL: &models.ThumbnailURL{},
				Tools: &models.Tools{Values: []models.HowToItem{
					{Type: "HowToTool", Text: "four", Quantity: 1},
					{Type: "HowToTool", Text: "plaque allant au four", Quantity: 1},
					{Type: "HowToTool", Text: "casserole", Quantity: 1},
					{Type: "HowToTool", Text: "passoire", Quantity: 1},
					{Type: "HowToTool", Text: "mixeur", Quantity: 1},
				}},
				Yield: &models.Yield{Value: 2},
				URL:   "https://www.quitoque.fr/recettes/betterave-de-couleur-et-carotte-roties-au-miel-feta-et-pesto-de-sauge",
				Video: &models.Videos{},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
