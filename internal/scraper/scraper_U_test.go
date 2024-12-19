package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_U(t *testing.T) {
	testcases := []testcase{
		{
			name: "uitpaulineskeuken.nl",
			in:   "https://uitpaulineskeuken.nl/recept/breekbrood-met-kaasfondue",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookTime:      "PT345M",
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				DateModified:  "2023-10-31T11:09:22+00:00",
				Description:   &models.Description{Value: "Dit breekbrood met kaasfondue heeft 2 vullingen: zongedroogde tomaat en gedroogde vijgen. Een feestelijk borrelrecept voor je volgende party!"},
				Keywords: &models.Keywords{
					Values: "familie recepten,brood en deeg,noten recepten,tomaat,franse recepten,italiaanse recepten,kerst,oud en nieuw,pasen,sinterklaas,verjaardag,oven,vegetarische recepten,herfst recepten,lente recepten,winter recepten,zomer recepten,borrelhapjes,hoofdgerecht",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"50 + 450 gr bloem",
						"7 gr gist",
						"100 gr roomboter",
						"200 ml volle melk",
						"4 + 1 eieren",
						"12 gr zout",
						"2 el olijfolie",
						"1 pak Emmi Fondü",
						"2 takjes rozemarijn",
						"4 takjes tijm",
						"3 zongedroogde tomaatjes",
						"7 groene olijven (zonder pit)",
						"2 gedroogde vijgen",
						"25 gr walnoten",
						"3 el maanzaad",
						"3 el pompoenpitten",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Laat de boter samen met de melk in een steelpannetje smelten. Voeg wanneer de boter gesmolten is 50 gram bloem toe en roer met een garde door de melk tot er een aardappelpuree structuur ontstaat. Schep het over in een bakje en laat volledig afkoelen, bijvoorbeeld in de koelkast."},
						{Type: "HowToStep", Text: "Kneed het deeg met behulp van een keukenmachine. Doe het afgekoelde boter- en bloemmengsel in de kom van de keukenmachine, voeg de eieren, rest van de bloem (450 gram), gist, olijfolie en zout toe en kneed 20 minuten op een lage stand. Zet de machine iets hoger en kneed nog eens 10 minuten.\nTest of het brood genoeg gekneed is door een balletje deeg langzaam uit elkaar te trekken, als het deeg niet scheurt, of bijna niet, is het klaar om te laten rijzen."},
						{Type: "HowToStep", Text: "Laat het deeg afgedekt 30 minuten op kamertemperatuur rijzen. Zet het deeg daarna afgedekt, bijvoorbeeld in een ruime luchtdichte doos 4 uur (of een nacht wat je wilt) in de koelkast om het daar verder te laten rijzen.\nKies een luchtdichte doos of kom waar het deeg 1,5 tot 2 keer zo groot in kan worden."},
						{Type: "HowToStep", Text: "Bereid de vulling voor. Deel eerst het deeg in tweeën. Hak voor de helft van de broodjes de zongedroogde tomaat, olijf en rozemarijn fijn. Kneed door de helft van het deeg. Hak voor de andere helft van de broodjes de vijg, tijm en walnoten fijn. Kneed door de andere helft van het deeg. Vorm in totaal 20 tot 22 bolletjes, elk van 30 gram."},
						{Type: "HowToStep", Text: "Bekleed een bakplaat met bakpapier. Zet een (ovenbestendige) kom waar later de Emmi Fondü in komt op het bakpapier, leg hier de bolletjes omheen. Leg de bolletjes met 1 centimeter uit elkaar, dit geeft ze ruimte om nog iets te rijzen. Laat nog 1,5 uur rijzen onder plastic folie (op kamertemperatuur)."},
						{Type: "HowToStep", Text: "Verwarm de oven voor op 180 graden (elektrisch), 160 graden (hetelucht).\nKlop het ei los en bestrijk de gerezen broodjes met het ei. Bestrooi de zongedroogde tomatenvariant met pompoenpitten en de vijgenvariant met maanzaad."},
						{Type: "HowToStep", Text: "Bak het breekbrood in 25 tot 30 minuten af in de voorverwarmde oven. De (ovenbestendige) kom laat je dus in de oven tussen de broodjes staan! Tijdens bakken gaan de broodjes namelijk verder rijzen, de kom zorgt ervoor dat de krans zijn vorm houdt, en er later genoeg ruimte overblijft om de fondue in te zetten."},
						{Type: "HowToStep", Text: "Maak terwijl het brood bakt de fondue. Doe de Emmi Fondü in een pannetje en laat de kaas langzaam smelten. Leg de broodkrans op een plank of grote schaal, zet de kom in het midden en vul met de fondue. Serveer het breekbrood met kaasfondue meteen."},
					},
				},
				Name:            "Breekbrood met kaasfondue",
				NutritionSchema: &models.NutritionSchema{},
				PrepTime:        "PT45M",
				ThumbnailURL:    &models.ThumbnailURL{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				Yield:           &models.Yield{Value: 1},
				URL:             "https://uitpaulineskeuken.nl/recept/breekbrood-met-kaasfondue",
				Video:           &models.Videos{},
			},
		},
		{
			name: "unsophisticook.com",
			in:   "https://unsophisticook.com/oven-roasted-baby-potatoes/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "Side Dish"},
				CookTime:      "PT25M",
				Cuisine:       &models.Cuisine{Value: "Italian"},
				DatePublished: "2021-07-26T12:46:27+00:00",
				Description: &models.Description{
					Value: "These flavorful Italian oven roasted baby potatoes are my go-to side dish for just about everything... Easy, fast, and made with pantry staples, they&#x27;re a total winner!",
				},
				Keywords: &models.Keywords{
					Values: "garlic baby potatoes, oven roasted baby potatoes, roasted baby red potatoes, roasted baby yukon gold potatoes, roasted whole baby potatoes",
				},
				Image: &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"1 3-lb. bag baby gold potatoes",
						"2 tablespoon avocado oil (or any neutral cooking oil)",
						"1 tablespoon Italian seasoning*", "1 teaspoon kosher salt",
						"1/2 teaspoon garlic powder",
						"freshly ground black pepper (to taste)",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Text: "Preheat oven to 425°F degrees. Slice potatoes lengthwise.", Type: "HowToStep"},
						{Text: "Add potatoes to a large mixing bowl. Toss with avocado oil, Italian seasoning, garlic powder, salt, and freshly ground black pepper.", Type: "HowToStep"},
						{Text: "Position potatoes cut side down on a large baking sheet. Bake at 425°F degrees for 23-25 minutes or until the potatoes are golden brown on the cut side and can easily be pierced with a fork.", Type: "HowToStep"},
						{Text: "Season with additional salt and pepper to taste. Serve warm.", Type: "HowToStep"},
					},
				},
				Name: "Italian Oven Roasted Baby Potatoes",
				NutritionSchema: &models.NutritionSchema{
					Calories:      "148",
					Carbohydrates: "30",
					Fat:           "3.5",
					Fiber:         "3.5",
					Protein:       "4.6",
					SaturatedFat:  "0.5",
					Servings:      "1 g",
					Sodium:        "240",
					Sugar:         "4.6",
				},
				PrepTime:  "PT5M",
				TotalTime: "PT30M",
				Yield:     &models.Yield{Value: 8},
				URL:       "https://unsophisticook.com/oven-roasted-baby-potatoes/",
			},
		},
		{
			name: "usapears.org",
			in:   "https://usapears.org/recipe/pear-vinegar/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        &models.SchemaType{Value: "Recipe"},
				Category:      &models.Category{Value: "uncategorized"},
				CookingMethod: &models.CookingMethod{},
				Cuisine:       &models.Cuisine{},
				DatePublished: "2022-11-15T17:41:17+00:00",
				Description:   &models.Description{Value: "…"},
				Keywords:      &models.Keywords{},
				Image:         &models.Image{Value: anUploadedImage.String()},
				Ingredients: &models.Ingredients{
					Values: []string{
						"pear scraps, peels, cores, and bruised bits",
						"a large jar that fits all the scraps", "water",
						"mother of vinegar (you can purchase one, or use live vinegar – just look for “with the mother” on the bottle)",
						"piece of cloth to cover the jar and rubber band to hold it in place",
					},
				},
				Instructions: &models.Instructions{
					Values: []models.HowToItem{
						{Type: "HowToStep", Text: "Put all of the scraps and peels into your well-cleaned jar, and fill the jar with water"},
						{Type: "HowToStep", Text: "Cover the jar with a cloth rubber band it into place. This cloth also allows it to breathe as it ferments"},
						{Type: "HowToStep", Text: "Set it aside on your counter or in a cabinet away from direct sunlight. After a day or two, little bubbles will appear as the pear and water begin fermenting"},
						{Type: "HowToStep", Text: "Strain out the pears and solids and return the liquid to a clean, wide-mouthed jar"},
						{Type: "HowToStep", Text: "Add the mother of vinegar. Or if you are using a live vinegar, shake up the bottle and pour in about 2 tbsp"},
						{Type: "HowToStep", Text: "Cover the jar with cloth and a rubber band again"},
					},
				},
				Name:            "Pear Vinegar",
				NutritionSchema: &models.NutritionSchema{},
				ThumbnailURL:    &models.ThumbnailURL{},
				Tools:           &models.Tools{Values: []models.HowToItem{}},
				PrepTime:        "PT20M",
				Yield:           &models.Yield{Value: 1},
				URL:             "https://usapears.org/recipe/pear-vinegar/",
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
