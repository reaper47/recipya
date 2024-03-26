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
				AtContext: "https://schema.org",
				AtType:    models.SchemaType{Value: "Recipe"},
				Category: models.Category{
					Value: "Familie recepten",
				},
				CookTime:      "PT0H390M",
				Cuisine:       models.Cuisine{Value: "Italiaanse recepten"},
				DatePublished: "23 november 2023",
				Description: models.Description{
					Value: "Dit breekbrood met kaasfondue heeft 2 vullingen: zongedroogde tomaat en gedroogde vijgen. Een feestelijk borrelrecept voor je volgende party!",
				},
				Keywords: models.Keywords{
					Values: "Familie recepten, Brood en deeg, Noten recepten, Tomaat, Franse recepten, Italiaanse recepten, Kerst, Oud en nieuw, Pasen, Sinterklaas, Verjaardag, Oven, Vegetarische recepten, Herfst recepten, Lente recepten, Winter recepten, Zomer recepten, Borrelhapjes, Hoofdgerecht",
				},
				Image: models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"bloem 50 + 450gr", "gist 7gr", "roomboter 100gr", "volle melk 200ml",
						"eieren 4 + 1", "zout 12gr", "olijfolie 2el", "Emmi Fondü 1pak",
						"rozemarijn 2takjes", "tijm 4takjes", "zongedroogde tomaatjes 3",
						"groene olijven 7", "gedroogde vijgen 2", "walnoten 25gr", "maanzaad 3el",
						"pompoenpitten 3el",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Zo maak je het:",
						"1. Laat de boter samen met de melk in een steelpannetje smelten. Voeg wanneer de boter gesmolten is 50 gram bloem toe en roer met een garde door de melk tot er een aardappelpuree structuur ontstaat. Schep het over in een bakje en laat volledig afkoelen, bijvoorbeeld in de koelkast.",
						"2. Kneed het deeg met behulp van een keukenmachine. Doe het afgekoelde boter- en bloemmengsel in de kom van de keukenmachine, voeg de eieren, rest van de bloem (450 gram), gist, olijfolie en zout toe en kneed 20 minuten op een lage stand. Zet de machine iets hoger en kneed nog eens 10 minuten.\nTest of het brood genoeg gekneed is door een balletje deeg langzaam uit elkaar te trekken, als het deeg niet scheurt, of bijna niet, is het klaar om te laten rijzen.",
						"3. Laat het deeg afgedekt 30 minuten op kamertemperatuur rijzen. Zet het deeg daarna afgedekt, bijvoorbeeld in een ruime luchtdichte doos 4 uur (of een nacht wat je wilt) in de koelkast om het daar verder te laten rijzen.\nKies een luchtdichte doos of kom waar het deeg 1,5 tot 2 keer zo groot in kan worden.",
						"4. Bereid de vulling voor. Deel eerst het deeg in tweeën. Hak voor de helft van de broodjes de zongedroogde tomaat, olijf en rozemarijn fijn. Kneed door de helft van het deeg. Hak voor de andere helft van de broodjes de vijg, tijm en walnoten fijn. Kneed door de andere helft van het deeg. Vorm in totaal 20 tot 22 bolletjes, elk van 30 gram.",
						"5. Bekleed een bakplaat met bakpapier. Zet een (ovenbestendige) kom waar later de Emmi Fondü in komt op het bakpapier, leg hier de bolletjes omheen. Leg de bolletjes met 1 centimeter uit elkaar, dit geeft ze ruimte om nog iets te rijzen. Laat nog 1,5 uur rijzen onder plastic folie (op kamertemperatuur).",
						"6. Verwarm de oven voor op 180 graden (elektrisch), 160 graden (hetelucht).\nKlop het ei los en bestrijk de gerezen broodjes met het ei. Bestrooi de zongedroogde tomatenvariant met pompoenpitten en de vijgenvariant met maanzaad.",
						"7. Bak het breekbrood in 25 tot 30 minuten af in de voorverwarmde oven. De (ovenbestendige) kom laat je dus in de oven tussen de broodjes staan! Tijdens bakken gaan de broodjes namelijk verder rijzen, de kom zorgt ervoor dat de krans zijn vorm houdt, en er later genoeg ruimte overblijft om de fondue in te zetten.",
						"8. Maak terwijl het brood bakt de fondue. Doe de Emmi Fondü in een pannetje en laat de kaas langzaam smelten. Leg de broodkrans op een plank of grote schaal, zet de kom in het midden en vul met de fondue. Serveer het breekbrood met kaasfondue meteen.",
					},
				},
				Name:      "Breekbrood met kaasfondue",
				PrepTime:  "PT0H45M",
				TotalTime: "PT0H435M",
				URL:       "https://uitpaulineskeuken.nl/recept/breekbrood-met-kaasfondue",
				Yield:     models.Yield{Value: 1},
			},
		},
		{
			name: "usapears.org",
			in:   "https://usapears.org/recipe/pear-vinegar/",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				DatePublished: "2022-11-15T17:41:17+00:00",
				Description:   models.Description{Value: "…"},
				Image:         models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"pear scraps, peels, cores, and bruised bits",
						"a large jar that fits all the scraps", "water",
						"mother of vinegar (you can purchase one, or use live vinegar – just look for “with the mother” on the bottle)",
						"piece of cloth to cover the jar and rubber band to hold it in place",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Put all of the scraps and peels into your well-cleaned jar, and fill the jar with water",
						"Cover the jar with a cloth rubber band it into place. This cloth also allows it to breathe as it ferments",
						"Set it aside on your counter or in a cabinet away from direct sunlight. After a day or two, little bubbles will appear as the pear and water begin fermenting",
						"Strain out the pears and solids and return the liquid to a clean, wide-mouthed jar",
						"Add the mother of vinegar. Or if you are using a live vinegar, shake up the bottle and pour in about 2 tbsp",
						"Cover the jar with cloth and a rubber band again",
					},
				},
				Name:     "Pear Vinegar",
				PrepTime: "PT20M",
				URL:      "https://usapears.org/recipe/pear-vinegar/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
