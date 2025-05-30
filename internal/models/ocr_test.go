package models_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAzureVision_Recipe(t *testing.T) {
	testcases := []struct {
		name string
		file string
		want models.Recipe
	}{
		{
			name: "recipe1.jpg",
			file: "recipe1.json",
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "Recipe created using Azure AI Document Intelligence.",
				Images:      []uuid.UUID{},
				Ingredients: []string{
					"3 Tbs. butter",
					"6 slices firm, course-textured bread",
					"6 c. meat or chicken broth",
					"salt",
					"6 eggs",
					"freshly ground pepper",
					"6 Tbs. freshly grated Parmesan cheese",
				},
				Instructions: []string{
					"* melt butter in a saucepan and fry the bread slices on both sides until golden brown. Divide among 6 soup bowls. Place bowls in are over which has been pre heated to 350ºF and turned off.",
					"* Meanwhile, Uring wroth to a boil, adding pal if necessary. Break an egg on to each slice of bread, pour the boiling broth over it an add pepper to taste. Sprinkle with Parme cheese. Harrish w/parsley nut meg.",
				},
				Keywords: []string{},
				Name:     "Zuppa Pavese (Pavia Soup)",
				Tools:    make([]models.HowToItem, 0),
				URL:      "OCR",
				Videos:   []models.VideoObject{},
				Yield:    6,
			},
		},
		{
			name: "recipe2.jpg",
			file: "recipe2.json",
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "Recipe created using Azure AI Document Intelligence.",
				Images:      []uuid.UUID{},
				Ingredients: []string{
					"1 cup of potatoes",
					"1/2 tbsp salt (table)",
					"2 tsp curmin spices",
					"2.5 cups flour",
					"1/2 cup olive oil",
				},
				Instructions: []string{
					"Mix all ingredients together except for olive oil Heat a pan on the oven Pour the olive oil in the pan Add the mix and potatoes to the pan Grill for 20m until the potatoes are crispy and brown Serve hot",
				},
				Keywords: []string{},
				Name:     "Oven-Baked Potatoes",
				Tools:    make([]models.HowToItem, 0),
				URL:      "OCR",
				Videos:   []models.VideoObject{},
				Yield:    1,
			},
		},
		{
			name: "recipe3.jpg",
			file: "recipe3.json",
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "from Shaker Museum -N.Y",
				Images:      []uuid.UUID{},
				Ingredients: []string{
					"1c. butter",
					"1 T. flour",
					"21/2 c. sugar",
					"IT. old bourbon",
					"1 c. heavy cream",
					`I un baked 8"pie shell`,
					"4 egg yolks, well beat",
				},
				Instructions: []string{
					"In small saucepan melt butter, add",
					"sugar and cream. Bring just to the",
					"bail. In mixing bowl beat eggs,",
					"stir in flour and first mixture.",
					"add old bourbon and pour into",
					"pastry lined pietin. Bake at 375'",
					"for 35 minutes or until set,",
				},
				Keywords: []string{},
				Name:     "Kentucky Pudding",
				Tools:    make([]models.HowToItem, 0),
				URL:      "OCR",
				Videos:   []models.VideoObject{},
				Yield:    4,
			},
		},
		{
			name: "recipe4.jpg",
			file: "recipe4.json",
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "Recipe created using Azure AI Document Intelligence.",
				Images:      []uuid.UUID{},
				Ingredients: []string{
					"2 cups flour",
					"1/2 cup sugar",
					"1 cup walnuts and pecans (chopped)",
					"1 tsp vanilla",
					"2 sticks butter",
				},
				Instructions: []string{
					"melt butter and pour over dry ingredients. Spoon drop onto cookie sheet. (They don't Spread)",
					"Bake at 350° for 15-20 min. Cool completely. Coat with powdered sugar.",
					"Keep in airtight container. Do not refrigerate.",
				},
				Keywords: []string{},
				Name:     "Xenia's Polish Cookies",
				Tools:    make([]models.HowToItem, 0),
				URL:      "OCR",
				Videos:   []models.VideoObject{},
				Yield:    1,
			},
		},
		{
			name: "recipe5.jpg",
			file: "recipe5.json",
			want: models.Recipe{
				Category:    "uncategorized",
				Description: "Recipe created using Azure AI Document Intelligence.",
				Images:      []uuid.UUID{},
				Ingredients: []string{
					"2 cups cream corn",
					"1 kg. onion, choppe",
					"2 cups TOMATOES",
					"1 gr. pepper,",
					"8.02 con tomato Sauce",
					"3eggs, besten",
					"1/2cups pitted olives",
					"1 cup MILK.",
					"3/4 cup olive oil",
					"1 cup cornmeal",
					"/ T. sacr",
					"1/2 tspcumin",
				},
				Instructions: []string{
					"Simmer the first 9 ingredients for 20 minutes. Stie often - the cream corn tends to want To make it stick. add the next 4 ingredients. 9x 13 pan 325° for 1 hour.",
				},
				Keywords: []string{},
				Name:     "Tamace Die",
				Tools:    make([]models.HowToItem, 0),
				URL:      "OCR",
				Videos:   []models.VideoObject{},
				Yield:    1,
			},
		},
		{
			name: "recipe6.png",
			file: "recipe6.json",
			want: models.Recipe{
				Category:    "uncategorized",
				Description: ":selected: WHY THIS RECIPE WORKS: For a fast and easy meal with plenty of Mexican-inspired flavor, we turned to quick-cooking boneless chicken breasts. First, we gave the mild chicken a layer of spicy flavor by seasoning it with chili powder as well as salt and pepper. Next, we dredged the breasts in flour, which served two purposes: It created a barrier between the fat in the pan and the moisture in the cutlet so that the fat \"spit\" less, and it helped to produce a consistently brown and crispy crust. We used the same pan to whip up a simple and flavorful side dish from common Mexican ingredients. We toasted corn kernels in a bit of oil, which brought out their sweetness nicely. We then softened some tomatoes (cherry tomatoes were our favorite) and brightened up our salsa with cilantro and fresh lime juice. Garlic and shallot rounded out the flavor of the salsa. The bright salsa perfectly complemented our crispy chicken breasts. Be sure not to stir the corn when cooking in step 4 or it will not brown well. If using fresh corn, you will need three to four ears in order to yield 3 cups of kernels. See the sidebar that follows the recipe.",
				Images:      []uuid.UUID{},
				Ingredients: []string{
					"1/2 cup all-purpose flour",
					"4 (6- to 8-ounce) boneless, skinless chicken breasts, trimmed and pounded to 1/2-inch thickness",
					"1 teaspoon chili powder", "Salt and pepper", "3 tablespoons vegetable oil",
					"3 cups fresh or thawed frozen corn", "1 shallot, minced",
					"2 garlic cloves, minced", "12 ounces cherry tomatoes, halved",
					"1/4 cup minced fresh cilantro", "2 tablespoons lime juice",
				},
				Instructions: []string{
					"Spread flour into shallow dish. Pat chicken dry with paper towels and season with chili powder, salt, and pepper. Working with 1 chicken breast at time, dredge in flour, shaking off excess.",
					"Heat 2 tablespoons oil in 12-inch nonstick skillet over medium-high heat until just smoking. Lay chicken in skillet and cook until well browned on first side, 6 to 8 minutes.",
					"Flip chicken over, reduce heat to medium, and continue to cook until chicken registers 160 degrees, 6 to 8 minutes. Transfer chicken to plate and tent with aluminum foil.",
					"Add remaining 1 tablespoon oil to now-empty skillet and place over medium-high heat until shimmering. Add corn and cook, without stirring, until well browned and roasted, 8 to 10 minutes. Stir in shallot and garlic and cook until fragrant, about 30 seconds. Stir in tomatoes, scraping up any browned bits, and cook until just softened, about 2 minutes.",
					"Off heat, stir in cilantro and lime juice and season with salt and pepper to taste. Transfer vegetables to platter and serve with chicken.",
					"TEST KITCHEN TIP: CUTTING KERNELS OFF THE COB",
					"After removing husk and silk, stand ear upright in large bowl and use paring knife to slice kernels off cob.",
				},
				Keywords: []string{},
				Name:     "Sauteed Chicken with Cherry Tomato and Roasted Corn Salsa",
				Tools:    make([]models.HowToItem, 0),
				URL:      "OCR",
				Videos:   []models.VideoObject{},
				Yield:    4,
			},
		},
		{
			name: "recipe7.png",
			file: "recipe7.json",
			want: models.Recipe{
				Category:    "uncategorized",
				Description: ":selected: WHY THIS RECIPE WORKS: For a fast and easy meal with plenty of Mexican-inspired flavor, we turned to quick-cooking boneless chicken breasts. First, we gave the mild chicken a layer of spicy flavor by seasoning it with chili powder as well as salt and pepper. Next, we dredged the breasts in flour, which served two purposes: It created a barrier between the fat in the pan and the moisture in the cutlet so that the fat \"spit\" less, and it helped to produce a consistently brown and crispy crust. We used the same pan to whip up a simple and flavorful side dish from common Mexican ingredients. We toasted corn kernels in a bit of oil, which brought out their sweetness nicely. We then softened some tomatoes (cherry tomatoes were our favorite) and brightened up our salsa with cilantro and fresh lime juice. Garlic and shallot rounded out the flavor of the salsa. The bright salsa perfectly complemented our crispy chicken breasts. Be sure not to stir the corn when cooking in step 4 or it will not brown well. If using fresh corn, you will need three to four ears in order to yield 3 cups of kernels. See the sidebar that follows the recipe.",
				Images:      []uuid.UUID{},
				Ingredients: []string{
					"1/2 cup all-purpose flour",
					"4 (6- to 8-ounce) boneless, skinless chicken breasts, trimmed and pounded to 1/2-inch thickness",
					"1 teaspoon chili powder", "Salt and pepper", "3 tablespoons vegetable oil",
					"3 cups fresh or thawed frozen corn", "1 shallot, minced",
					"2 garlic cloves, minced", "12 ounces cherry tomatoes, halved",
					"1/4 cup minced fresh cilantro", "2 tablespoons lime juice",
				},
				Instructions: []string{
					"Spread flour into shallow dish. Pat chicken dry with paper towels and season with chili powder, salt, and pepper. Working with 1 chicken breast at time, dredge in flour, shaking off excess.",
					"Heat 2 tablespoons oil in 12-inch nonstick skillet over medium-high heat until just smoking. Lay chicken in skillet and cook until well browned on first side, 6 to 8 minutes.",
					"Flip chicken over, reduce heat to medium, and continue to cook until chicken registers 160 degrees, 6 to 8 minutes. Transfer chicken to plate and tent with aluminum foil.",
					"Add remaining 1 tablespoon oil to now-empty skillet and place over medium-high heat until shimmering. Add corn and cook, without stirring, until well browned and roasted, 8 to 10 minutes. Stir in shallot and garlic and cook until fragrant, about 30 seconds. Stir in tomatoes, scraping up any browned bits, and cook until just softened, about 2 minutes.",
					"Off heat, stir in cilantro and lime juice and season with salt and pepper to taste. Transfer vegetables to platter and serve with chicken.",
				},
				Keywords: []string{},
				Name:     "Sautéed Chicken with Cherry Tomato and Roasted Corn Salsa",
				Tools:    make([]models.HowToItem, 0),
				URL:      "OCR",
				Videos:   []models.VideoObject{},
				Yield:    4,
			},
		},
		{
			name: "recipe8.pdf",
			file: "recipe8.json",
			want: models.Recipe{
				Category:    "uncategorized",
				Cuisine:     "einschub: mitte",
				Description: "Recipe created using Azure AI Document Intelligence.",
				Images:      []uuid.UUID{},
				Ingredients: []string{
					"FÜR DIE SPRINGFORM (Ø 20 CM)",
					"etwas Fett",
					"BISKUITTEIG",
					"4 Eier (Größe M)",
					"60 g Zucker 8 BI. Dr. Oetker Gelatine weiß 1 Pck. Dr. Oetker Vanillin-Zucker 150 g Dr. Oetker Kuvertüre Weiß 400 ml Kokosmilch",
					"80 g Weizenmehl 1 gestr. TL Dr. Oetker Original Backin 20 g Dr. Oetker Gustin Feine Speisestärke",
					"MÜRBETEIG",
					"100 g Weizenmehl",
					"70 g",
					"weiche Butter oder Margarine",
					"40 g",
					"Zucker",
					"1 EL",
					"Wasser",
					"KOKOSCREME",
					"300 g",
					"kalte Schlagsahne",
					"ZUM BESTREICHEN",
					"3 EL helle Konfitüre , z.B. Marille-Maracuja",
					"Dr.Oetker",
					"Himmlische Kokoscreme-Torte ZUBEREITUNG: 60 MIN . ETWAS ÜBUNG ERFORDERLICH . ETWA 12 STÜCK",
					"ZUM TRÄNKEN",
					"etwa 4 EL Kokoslikör",
					"ZUM VERZIEREN",
					"etwa 50 g Kokosraspel",
					"PRODUKTE ANSEHEN",
					"Dr Oetker",
					"Vanillin Zucker",
					"Dr. Oetker Backpulver Original Backin",
					"Gustin",
					"DrOetker",
					"helene Speise- stärke",
					"Blatt Gelatine Qualsiz Gold extre",
					"Kuvertüre",
					"1",
					"Vorbereiten",
				},
				Instructions: []string{
					"Backblech mit Backpapier belegen. Springformboden fetten und mit Backpapier belegen. Backofen vorheizen.",
					"Ober- und Unterhitze: etwa 180 ℃ Heißluft: etwa 160 ℃",
					"Biskuitteig zubereiten",
					"Eier in einer Rührschüssel mit einem Mixer (Rührstäbe) auf höchster Stufe 1 Min. schaumig schlagen. Mit Vanillin-Zucker gemischten Zucker unter Rühren in 1 Min. einstreuen und die Masse weitere 2 Min. schlagen. Mehl mit Backin und Gustin mischen und kurz auf niedrigster Stufe unterrühren. Teig in der Form glatt streichen. Form auf dem Rost in den Backofen schieben.",
					"Einschub: Mitte",
					"Backzeit: etwa 20 Min.",
					"Springformrand lösen und entfernen. Biskuit auf einen mit Backpapier belegten Kuchenrost stürzen und erkalten lassen. Springform säubern und Boden fetten. Backofentemperatur erhöhen.",
					"Ober- und Unterhitze: etwa 200 ℃",
				},
				Keywords: []string{},
				Name:     "Himmlische Kokoscreme-Torte ZUBEREITUNG:60 MIN . ETWAS ÜBUNG ERFORDERLICH · ETWA 12 STÜCK",
				Times:    models.Times{Prep: 1 * time.Hour},
				Tools: []models.HowToItem{
					{Text: "FÜR DAS BACKBLECH", Quantity: 1, Type: "HowToTool"},
					{Text: "Backpapier", Quantity: 1, Type: "HowToTool"},
					{Text: "FÜR DIE SPRINGFORM (Ø 20 CM) Backpapier", Quantity: 1, Type: "HowToTool"},
				},
				URL:    "OCR",
				Videos: []models.VideoObject{},
				Yield:  12,
			},
		},
		{
			name: "recipe9.jpg",
			file: "recipe9.json",
			want: models.Recipe{
				Category:    "poultry",
				Cuisine:     "(java, indonesia)",
				Description: "A gorgeous coconut-milk curry from Java, Indonesia, perfumed with lemongrass, ginger, cinnamon sticks, and coriander. It's one of the benchmark dishes by which Indonesian home cooks are judged. If a young cook's opor ayam is as rich and delicate as it should be, she is well on her way to becoming skilled in the kitchen. The dish is a perfect showcase for a high-quality free-range chicken. A whole bird, cut into small, bone-in serving pieces, will yield the best results, though whole chicken parts can be substituted without really compromising the taste of the dish. If you have the inclina- tion to make them, fried shallots are traditionally strewn on top just before the dish is served.\n\nDaun salam leaves, the seasoning herb prized in Indonesian cooking, helps give this dish its unique aroma. I've often seen bay leaves listed as a substitute for daun salam in recipe books on Indonesian cuisine. While bay leaves have an aggressively mentholated taste, daun salam leaves are subtle, with a faintly foresty flavor. The only thing the two herbs share in common is that they are both green leaves and grow on trees. If you are unable to find daun salam leaves, omit them. The dish will still taste exquisite.",
				Images:      []uuid.UUID{},
				Ingredients: []string{
					"2 tablespoons Crisp-Fried Shallots (page 84),",
					"optional",
					"FOR THE FLAVORING PASTE",
					"1 tablespoon coriander seeds",
					"1 fresh red Holland chile or other fresh long, hot",
					"1 fresh red Holland chile or other fresh long, hot red chile such as Fresno or cayenne, stemmed and coarsely chopped (optional, but provides subtle heat and color) 6 shallots (about 5 ounces/140 grams), coarsely chopped",
					"2 cloves garlic, coarsely chopped",
					"1 piece fresh or thawed, frozen galangal,",
					"11/2 inches (4 centimeters) long, peeled and thinly sliced against the grain (about 11/2 tablespoons; optional)",
					"1 piece fresh ginger, 2 inches (5 centimeters) long, peeled and thinly sliced against the grain (about 2 tablespoons)",
				},
				Instructions: []string{"No instructions found in image."},
				Keywords:     []string{},
				Name:         "JAVANESE CHICKEN CURRY Opor Ayam",
				Nutrition:    models.Nutrition{},
				Times:        models.Times{},
				Tools:        []models.HowToItem{},
				UpdatedAt:    time.Time{},
				URL:          "OCR",
				Videos:       []models.VideoObject{},
				Yield:        4,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			xb, _ := os.ReadFile(filepath.Join("testdata", "ocr", tc.file))
			var av models.AzureDILayout
			_ = json.Unmarshal(xb, &av)

			got := av.Recipe()

			if !cmp.Equal(got, tc.want) {
				t.Log(cmp.Diff(got, tc.want))
				t.Fail()
			}
		})
	}
}
