package scraper_test

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
				AtContext:     atContext,
				AtType:        &models.SchemaType{Value: "Recipe"},
				CookingMethod: &models.CookingMethod{},
				Name:          "Simple Bruschetta",
				Description: &models.Description{
					Value: "Good tomatoes are the thing that matters most when it comes to making bruschetta - the classic Italian antipasto. It is such a simple preparation that paying attention to the little details matters.",
				},
				Category: &models.Category{Value: "Appetizer"},
				Cuisine:  &models.Cuisine{Value: "Easy"},
				Image:    &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
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
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Rinse and dry your tomatoes. Halve each of them, use a finger to remove the seeds, and cut out the cores. Roughly cut the tomatoes into 1/2-inch pieces and place in a medium bowl. Tear the basil into small pieces, and add that as well. Add 2 tablespoons of the olive oil, a small splash of vinegar, and a pinch of salt. Gently toss, taste, adjust if necessary, and set aside."},
						{Type: "HowToStep", Text: "Heat a grill or oven to medium-high. When it’s ready, use the remaining 2 tablespoons of the olive oil to brush across the slices of bread. Grill or bake until well-toasted and golden brown with a hint of char. Flipping when the first side is done. Remove from grilled and when cool enough to handle, rub both sides of each slice of bread with garlic."},
						{Type: "HowToStep", Text: "Cut each slice of bread in half if you like, and top each segment with the tomato mixture. And a finishing drizzle of olive oil is always nice."},
					},
				},
				NutritionSchema: &models.NutritionSchema{
					Calories:      "233 kcal",
					Carbohydrates: "19 g",
					Cholesterol:   "",
					Fat:           "15 g",
					Protein:       "6 g",
					Sodium:        "287 mg",
					Servings:      "1",
				},
				Keywords:      &models.Keywords{Values: "simple"},
				Yield:         &models.Yield{Value: 4},
				PrepTime:      "PT15M",
				CookTime:      "PT5M",
				DatePublished: "2022-06-29T14:30:00+00:00",
				ThumbnailURL:  &models.ThumbnailURL{},
				Tools:         &models.Tools{Values: []models.HowToItem{}},
				TotalTime:     "PT20M",
				URL:           "https://www.101cookbooks.com/simple-bruschetta/",
				Video:         &models.Videos{},
			},
		},
		{
			name: "15gram.be",
			in:   "https://15gram.be/recepten/mac-n-cheese-met-gehakt-en-pompoen",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				CookTime:      "PT30M",
				Description:   &models.Description{Value: "We nemen je mee op skivakantie! Of toch naar de après-ski maaltijd. De pompoenblokjes zijn al voorgesneden en mixen we door de saus. Daardoor kleurt die ook mooi oranje!"},
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"400 gr pompoenblokjes",
						"1 teentje knoflook",
						"300 gr rund-varkensgehakt",
						"200 gr tortiglioni pasta",
						"150 gr zure room",
						"75 gr geraspte cheddar",
						"olijfolie",
						"zout",
						"zwarte peper",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Breng een ruime pot gezouten water aan de kook voor de pasta."},
						{Type: "HowToStep", Text: "Verwarm de oven voor op 220°C en bedek een bakplaat met bakpapier."},
						{Type: "HowToStep", Text: "Schik de pompoenblokjes op de bakplaat en pers de knoflook erbij. Hussel door elkaar met olijfolie, zout en zwarte peper. Rooster in 20 min. gaar in de oven."},
						{Type: "HowToStep", Text: `Verhit een scheutje olijfolie in een pan op hoog vuur. Bak het gehakt in 6-8 min. Prak met een spatel in "chunks", het hoeft niet helemaal fijn te zijn.`},
						{Type: "HowToStep", Text: "Kook de pasta volgens de instructies op de verpakking."},
						{Type: "HowToStep", Text: "Schik de geroosterde pompoen in een blender of maatbeker en mix, samen met de zure room en de helft van de geraspte cheddar, tot een gladde saus. Proef en breng op smaak met extra zout of zwarte peper (zie tip)."},
						{Type: "HowToStep", Text: "Giet de pasta af en schik opnieuw in de pot. Meng met de roomsaus en het gehakt. Stort in een ovenschaal en strooi de rest van de kaas erover. Gratineer nog 5-10 min. in de oven voor een mooi kaaskorstje."},
						{Type: "HowToStep", Text: "TIP: Geef de saus wat extra punch met paprikapoeder, chilivlokken of bouillon."},
					},
				},
				Keywords:        &models.Keywords{Values: "Foodbag,gehakt,pasta,pompoen,skikost"},
				Name:            "Mac 'n cheese met gehakt en pompoen",
				NutritionSchema: &models.NutritionSchema{},
				ThumbnailURL:    &models.ThumbnailURL{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				Yield:           &models.Yield{Value: 2},
				URL:             "https://15gram.be/recepten/mac-n-cheese-met-gehakt-en-pompoen",
				Video:           &models.Videos{},
			},
		},
		{
			name: "750g.com",
			in:   "https://www.750g.com/carbonara-vegetarienne-r202631.htm",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Plats"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				CookTime:      "PT10M",
				DatePublished: "2017-06-28",
				Description: &models.Description{
					Value: "Une variante v&amp;eacute;g&amp;eacute;tarienne de la Carbonara dans laquelle&amp;nbsp;les&amp;nbsp;lardons sont remplac&amp;eacute;s par des l&amp;eacute;gumes.",
				},
				Keywords: &models.Keywords{
					Values: "carbonara vegetarienne,carbonara aux légumes,pâtes carbonara sans lardons,spaghetti à la carbonara,pâtes à la carbonara,pâtes aux légumes,spaghetti,pâtes,carottes,courgettes,petits pois,plats végétariens,cuisine estivale,cuisine végétarienne,750green",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"500 g de spaghetti bio Barilla", "4 carottes", "2 courgettes", "2 oeufs",
						"1 oignon", "80 g de parmesan", "1 poignée de petits pois frais",
						"4 c. à s. d'huile d'olive", "4 feuilles de basilic frais", "Poivre",
						"Sel ou sel fin",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Disposez tous les ingr&amp;eacute;dients sur un plateau."},
						{Type: "HowToStep", Text: "&amp;Eacute;cossez les petits pois.&amp;nbsp;&amp;Eacute;pluchez les carottes et coupez-les en julienne.&amp;nbsp;Coupez les courgettes en julienne.&amp;nbsp;"},
						{Type: "HowToStep", Text: "Coupez l&#039;oignon en lamelles fines.&amp;nbsp;"},
						{Type: "HowToStep", Text: "Chauffez l&#039;huile d&#039;olive dans une po&amp;ecirc;le puis faites revenir l&#039;oignon, les carottes et les courgettes pendant 5 minutes. Ajoutez les petits pois, salez le tout et laissez rissoler pendant encore 5 minutes.&amp;nbsp;"},
						{Type: "HowToStep", Text: "Versez le parmesan et les &amp;oelig;ufs dans un saladier, ajoutez une pinc&amp;eacute;e de sel et du poivre du moulin."},
						{Type: "HowToStep", Text: "Fouettez et r&amp;eacute;servez."},
						{Type: "HowToStep", Text: "Portez &amp;agrave; &amp;eacute;bullition un grand volume d&#039;eau dans une casserole &amp;agrave; bords hauts.&amp;nbsp;D&amp;egrave;s l&amp;rsquo;&amp;eacute;bullition, ajoutez 20 g de gros sel, ajoutez les&amp;nbsp;p&amp;acirc;tes et remuez une fois avec une cuill&amp;egrave;re en bois. Faites-les cuire selon le temps de cuisson indiqu&amp;eacute; sur le paquet.&amp;nbsp;&amp;nbsp;"},
						{Type: "HowToStep", Text: `&amp;Eacute;gouttez les spaghetti "al dente" dans une passoire, puis versez-les dans le saladier avec la sauce aux &amp;oelig;ufs et parmesan. Remuez.&amp;nbsp;`},
						{Type: "HowToStep", Text: "Ajoutez&amp;nbsp;les l&amp;eacute;gumes, remuez bien, ajoutez les feuilles de basilic cisel&amp;eacute;es&amp;nbsp;et servez imm&amp;eacute;diatement, bien chaud."},
						{Type: "HowToStep", Text: "Bonne d&amp;eacute;gustation !"},
					},
				},
				Name:            "Carbonara végétarienne",
				NutritionSchema: &models.NutritionSchema{},
				PrepTime:        "PT30M",
				ThumbnailURL:    &models.ThumbnailURL{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				TotalTime:       "PT40M",
				Yield:           &models.Yield{Value: 4},
				URL:             "https://www.750g.com/carbonara-vegetarienne-r202631.htm",
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
