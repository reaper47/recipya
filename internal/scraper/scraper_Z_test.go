package scraper_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestScraper_Z(t *testing.T) {
	testcases := []testcase{
		{
			name: "zeit.de",
			in:   "https://www.zeit.de/zeit-magazin/wochenmarkt/2021-08/kohlrabi-fenchel-carpaccio-fior-di-latte-rezept",
			want: models.RecipeSchema{
				AtContext:     "https://schema.org",
				AtType:        models.SchemaType{Value: "Recipe"},
				DateModified:  "2021-08-16T12:03:26+02:00",
				DatePublished: "2021-08-16T12:03:26+02:00",
				Description: models.Description{
					Value: "Am besten lässt man den Kohlrabi roh und hobelt ihn in hauchdünne Scheiben. Für ein vegetarisches Carpaccio ganz in Weiß kommen dann noch Fenchel und Fior di Latte hinzu.",
				},
				Keywords: models.Keywords{Values: "Gemüse, Leichte Küche, Sommer, Salate"},
				Image:    models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"1 Kohlrabi (klein)", "1 Fenchelknolle", "1/2 Bund Basilikum",
						"2 EL Pistazienkerne", "4 EL Olivenöl", "2 EL Zitronensaft", "1 TL Honig",
						"Salz", "Pfeffer",
					},
				},
				Instructions: models.Instructions{Values: []string{
					"\"In meinem Hochbeet wächst der Kohlrabi gerade so gut und ich weiß gar nicht, was ich damit kochen soll\", erzählte mir neulich eine Kollegin. Verständlich. Die meisten unter uns kennen Kohlrabi nur verkocht in einer Mehlschwitze schwimmend, als Beilage zu Kartoffel und Fleisch. Dabei schmeckt das knackige Gemüse auch roh und hat einen frischen Kohlgeschmack. Zusammen mit dem fein-süßlichen Geschmack von Fenchel und der säuerlich-cremigen Fior di Latte, der kleinen Schwester der Burrata, entsteht auf dem Teller im Nu ein kleines Fest in Weiß.",
					"Den Kohlrabi schälen, den Fenchel waschen. Beide vom unteren Strunk befreien. Den Kohlrabi einmal von der Spitze zum Boden in der Mitte durchschneiden. Nun beide Gemüse mit einem scharfen und stabilen Messer hauchdünn schneiden, bessere Ergebnisse erzielt man auf einer Mandoline. Idealerweise sind die Scheiben so dünn, dass man fast hindurchschauen kann.",
					"Nun aus dem Olivenöl, dem Zitronensaft und Honig ein Dressing zusammenrühren, mit Salz und Pfeffer abschmecken. Die dünnen Gemüsescheiben fächerartig auf einen großen Teller legen, den Käse zerreißen und darauf verteilen. Das Dressing darübergeben und mit Basilikum und Pistazien servieren.",
				}},
				Name:  "Kohlrabi-Fenchel-Carpaccio: Kohlrabi hat Besseres verdient als Mehlschwitze",
				Yield: models.Yield{Value: 2},
				URL:   "https://www.zeit.de/zeit-magazin/wochenmarkt/2021-08/kohlrabi-fenchel-carpaccio-fior-di-latte-rezept",
			},
		},
		{
			name: "zenbelly.com",
			in:   "https://www.zenbelly.com/short-ribs/",
			want: models.RecipeSchema{
				AtContext: atContext,
				AtType:    models.SchemaType{Value: "Recipe"},
				Name:      "pressure cooker honey balsamic short ribs",
				Description: models.Description{
					Value: "rich and decadent short ribs are cooked in a fraction of the time in the instant pot.",
				},
				Category:      models.Category{Value: "beef"},
				CookingMethod: models.CookingMethod{Value: "pressure cook"},
				Image:         models.Image{Value: anUploadedImage.String()},
				Ingredients: models.Ingredients{
					Values: []string{
						"3.5 pounds bone-in short ribs (try to get ones that are at least 1.5 inches thick)",
						"salt and pepper",
						"1 tablespoon extra virgin olive oil",
						"1 medium onion, sliced",
						"4-5 cloves garlic, smashed and roughly chopped",
						"1/2 cup balsamic vinegar",
						"1/4 cup tamari or coconut aminos",
						"2 tablespoons honey",
						"1 tablespoon dijon mustard",
						"1 cup chicken or beef broth or stock",
						"2 sprigs thyme",
						"1 lemon, zested and juiced",
						"fresh parsley, roughly chopped",
					},
				},
				Instructions: models.Instructions{
					Values: []string{
						"Liberally season the short ribs with salt and pepper &#8211; ideally a day in advance and refrigerated, " +
							"but if not, up to two hours at room temperature.",
						"Turn your Instant Pot on Sauté (high, if it's an option). Once it reads HOT, add the oil. Brown " +
							"the short ribs in batches, until well browned on the meaty side (you don't " +
							"need to brown them on all sides)",
						"Remove the short ribs to a plate and add the onions. Sauté for 5-6 minutes, until browned and softened.",
						"Add the garlic and saut\u00e9 for another minute, until fragrant.",
						"Add the balsamic vinegar, tamari, honey, mustard, broth, and thyme. Once it starts to bubble, add the " +
							"short ribs back in, meaty side down.",
						"Hit Cancel on the Instant Pot. Lock on the lid, making sure the seal is in place. Set to cook for 40 " +
							"minutes at high pressure, making sure the valve is set to sealing.",
						"When the time is up, release the pressure and remove the lid once it unlocks. Remove the ribs to a " +
							"platter and turn the Instant Pot back to Saut\u00e9. Reduce the sauce until it's about 3 cups in volume.",
						"Pour the sauce over the ribs. Pour over the lemon juice and sprinkle with lemon zest and parsley.",
						"Serve over polenta, mashed potatoes, or whatever you'd like."},
				},
				Keywords:      models.Keywords{Values: "short ribs"},
				Yield:         models.Yield{Value: 4},
				PrepTime:      "PT20M",
				CookTime:      "PT1H10M",
				DatePublished: "2020-09-01",
				TotalTime:     "PT-473507H34M18S",
				URL:           "https://www.zenbelly.com/short-ribs/",
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			test(t, tc)
		})
	}
}
