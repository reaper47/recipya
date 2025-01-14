package integrations_test

import (
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/integrations"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMealieImport(t *testing.T) {
	testcases := []struct {
		name       string
		mealieJSON string
		wantRecipe models.Recipe
	}{
		{
			name:       "old API",
			mealieJSON: `{"id":"843b4a6d-6855-48c3-8186-22f096310243","userId":"e72ff251-4693-4e44-ad1d-9d9c2b033541","groupId":"083bba0c-e400-4b84-8055-b01a888b27fd","name":"Roasted Vegetable Bowls with Green Tahini","slug":"roasted-vegetable-bowls-with-green-tahini","image":"Z4Ox","recipeYield":"6 servings","totalTime":"45 minutes","prepTime":"15 minutes","cookTime":null,"performTime":"30 minutes","description":"Roasted Vegetable Bowls! Crispy tender roasted veggies, buttery avocado, all together in a bowl with a drizzle of green tahini sauce.","recipeCategory":[],"tags":[{"id":"70cc7ab9-cc6f-41d0-b8b9-8d16384f857e","name":"Vegetable Bowl Recipe","slug":"vegetable-bowl-recipe"},{"id":"da629ccc-56cb-4400-bce1-55ca0f14905b","name":"Roasted Vegetable Bowls","slug":"roasted-vegetable-bowls"},{"id":"e3184b1f-2bd0-48b8-b766-fa3eaf8285a5","name":"Green Tahini","slug":"green-tahini"}],"tools":[],"rating":4,"orgURL":"https://pinchofyum.com/30-minute-meal-prep-roasted-vegetable-bowls-with-green-tahini","dateAdded":"2024-04-12","dateUpdated":"2024-04-12T18:14:29.168064","createdAt":"2024-04-12T18:06:06.692275","updateAt":"2024-04-12T18:07:56.850947","lastMade":null,"recipeIngredient":[{"quantity":8.0,"unit":null,"food":{"id":"31d502b5-dea9-4580-b8b2-86bfde80f456","name":"carrot","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.825124","updateAt":"2024-03-05T06:54:56.825126"},"note":"large, peeled and chopped","isFood":true,"disableAmount":false,"display":"8 carrot large, peeled and chopped","title":null,"originalText":"8 large carrots, peeled and chopped","referenceId":"0b7f622f-a6f1-4e51-b381-d665ee54da47"},{"quantity":3.0,"unit":null,"food":null,"note":"chopped","isFood":true,"disableAmount":false,"display":"3 chopped","title":null,"originalText":"3 golden potatoes, chopped","referenceId":"d0944e99-f189-48aa-b951-1d2c1aaf7655"},{"quantity":1.0,"unit":null,"food":{"id":"05bb0bf7-dc26-4961-9bce-bd5563f7a6c7","name":"broccoli","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.552679","updateAt":"2024-03-05T06:54:56.552681"},"note":"cut into florets","isFood":true,"disableAmount":false,"display":"1 broccoli cut into florets","title":null,"originalText":"1 head of broccoli, cut into florets","referenceId":"3b1e5c67-5f02-49b8-88d7-6db8a681de05"},{"quantity":1.0,"unit":null,"food":{"id":"1a0beaa6-b6f2-4143-81e9-6709fe00d33a","name":"cauliflower","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.214977","updateAt":"2024-03-05T06:54:56.214981"},"note":"cut into florets","isFood":true,"disableAmount":false,"display":"1 cauliflower cut into florets","title":null,"originalText":"1 head of cauliflower, cut into florets","referenceId":"c0836160-4afd-4157-a780-7ac0a7f41a6f"},{"quantity":0.0,"unit":null,"food":null,"note":"","isFood":true,"disableAmount":false,"display":"","title":null,"originalText":"olive oil and salt","referenceId":"ac21685f-510d-4534-b4cd-e0ac57e04a04"},{"quantity":0.5,"unit":{"id":"56939576-ff3a-4760-98c4-cd7e7ea8418b","name":"cup","pluralName":null,"description":"","extras":{},"fraction":true,"abbreviation":"","pluralAbbreviation":"","useAbbreviation":false,"aliases":[],"createdAt":"2024-03-05T06:59:19.679450","updateAt":"2024-03-05T06:59:19.679454"},"food":{"id":"02bc4201-08ca-45b2-b032-7babfa4346f4","name":"olive oil","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.572704","updateAt":"2024-03-05T06:54:56.572706"},"note":"mild tasting)","isFood":true,"disableAmount":false,"display":"¹/₂ cup olive oil mild tasting)","title":null,"originalText":"1/2 cup olive oil (mild tasting)","referenceId":"af426c53-d5e4-4ab0-bc17-6b80ab2d389d"},{"quantity":0.5,"unit":{"id":"56939576-ff3a-4760-98c4-cd7e7ea8418b","name":"cup","pluralName":null,"description":"","extras":{},"fraction":true,"abbreviation":"","pluralAbbreviation":"","useAbbreviation":false,"aliases":[],"createdAt":"2024-03-05T06:59:19.679450","updateAt":"2024-03-05T06:59:19.679454"},"food":{"id":"c30d5cf5-d7e6-4f8b-8338-a010cae94441","name":"water","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:59:33.037532","updateAt":"2024-03-05T06:59:33.037536"},"note":"","isFood":true,"disableAmount":false,"display":"¹/₂ cup water","title":null,"originalText":"1/2 cup water","referenceId":"5cc3edfb-11c3-403f-afc0-f10d7863791c"},{"quantity":0.25,"unit":{"id":"56939576-ff3a-4760-98c4-cd7e7ea8418b","name":"cup","pluralName":null,"description":"","extras":{},"fraction":true,"abbreviation":"","pluralAbbreviation":"","useAbbreviation":false,"aliases":[],"createdAt":"2024-03-05T06:59:19.679450","updateAt":"2024-03-05T06:59:19.679454"},"food":{"id":"337615d4-f263-4289-80e2-7fe79983c29e","name":"tahini","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.942480","updateAt":"2024-03-05T06:54:56.942482"},"note":"","isFood":true,"disableAmount":false,"display":"¹/₄ cup tahini","title":null,"originalText":"1/4 cup tahini","referenceId":"e77a5636-f06c-479d-932e-86920ff04ae3"},{"quantity":0.0,"unit":null,"food":{"id":"3fce1ca1-3fbe-4a29-b0c8-6411a581be4d","name":"cilantro","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.640255","updateAt":"2024-03-05T06:54:56.640257"},"note":"and/or parsley","isFood":true,"disableAmount":false,"display":"cilantro and/or parsley","title":null,"originalText":"a big bunch of cilantro and/or parsley","referenceId":"eee40c2c-ce2d-42db-abd7-2743db1a017e"},{"quantity":1.0,"unit":{"id":"b9ca3f9e-f7d5-4bce-8181-29ae56939ce6","name":"clove","pluralName":null,"description":"","extras":{},"fraction":true,"abbreviation":"","pluralAbbreviation":"","useAbbreviation":false,"aliases":[],"createdAt":"2024-03-05T08:19:26.421964","updateAt":"2024-03-05T08:19:26.421968"},"food":{"id":"7f2f8ad6-035d-46c7-8503-9e6bf7281165","name":"garlic","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.355030","updateAt":"2024-03-05T06:54:56.355032"},"note":"","isFood":true,"disableAmount":false,"display":"1 clove garlic","title":null,"originalText":"1 clove garlic","referenceId":"84153160-0675-49f6-9363-0f255b0fdf1c"},{"quantity":0.0,"unit":null,"food":null,"note":"(about 2 tablespoons)","isFood":true,"disableAmount":false,"display":"(about 2 tablespoons)","title":null,"originalText":"squeeze of half a lemon (about 2 tablespoons)","referenceId":"777d50a9-ad50-4b03-a5e8-331ae6cc94f1"},{"quantity":0.5,"unit":{"id":"cda9b5eb-21c5-4acf-b65f-7b397e560eb3","name":"teaspoon","pluralName":null,"description":"","extras":{},"fraction":true,"abbreviation":"","pluralAbbreviation":"","useAbbreviation":false,"aliases":[],"createdAt":"2024-03-05T08:19:34.732459","updateAt":"2024-03-05T08:19:34.732462"},"food":{"id":"90a64fd5-ce9d-4774-a1c6-68c65ed5afea","name":"salt","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.265153","updateAt":"2024-03-05T06:54:56.265155"},"note":"","isFood":true,"disableAmount":false,"display":"¹/₂ teaspoon salt","title":null,"originalText":"1/2 teaspoon salt (more to taste)","referenceId":"cd6dc996-55bf-40a7-a280-c6b754157f2d"},{"quantity":6.0,"unit":null,"food":{"id":"b72ab124-26bc-40f1-a59c-3d469fe890b1","name":"eggs","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.442542","updateAt":"2024-03-05T06:54:56.442544"},"note":"hard boiled (or other protein)","isFood":true,"disableAmount":false,"display":"6 eggs hard boiled (or other protein)","title":null,"originalText":"6 hard boiled eggs (or other protein)","referenceId":"b56905f5-1d19-4d20-8498-3fd9285bf8c3"},{"quantity":3.0,"unit":null,"food":{"id":"d20ca2a2-f869-4680-937a-e00349b34893","name":"avocado","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.787485","updateAt":"2024-03-05T06:54:56.787488"},"note":"","isFood":true,"disableAmount":false,"display":"3 avocado","title":null,"originalText":"3 avocados","referenceId":"4e2cbbd8-43b0-4f37-8247-00b547ae99e3"}],"recipeInstructions":[{"id":"4f6932fd-1da1-4264-a7b3-fc8c7f8bcec1","title":"","text":"Prep","ingredientReferences":[]},{"id":"668767d9-1f73-4283-b450-fed81bb9521a","title":"","text":"Preheat the oven to 425 degrees.","ingredientReferences":[]},{"id":"ce35bb0f-aac7-4396-b4c6-73e28a8bf1bc","title":"","text":"Roasted Vegetables","ingredientReferences":[]},{"id":"cc60b45f-014c-4323-9767-4e3e957dead0","title":"","text":"Arrange your vegetables onto a few baking sheets lined with parchment (I keep each vegetable in its own little section). Toss with olive oil and salt. Roast for 25-30 minutes.","ingredientReferences":[]},{"id":"517eb65d-dc49-4d63-9a2d-82dbdee8f728","title":"","text":"Sauce","ingredientReferences":[]},{"id":"039c6817-4e7f-4690-9c84-20574bec62a4","title":"","text":"While the veggies are roasting, blitz up your sauce in the food processor or blender.","ingredientReferences":[]},{"id":"a78f9737-4b5b-4798-aa48-9f087d63c535","title":"","text":"Finish","ingredientReferences":[]},{"id":"ade3a45d-1546-47ca-ab0b-441707c58429","title":"","text":"Voila! Portion and save for the week! Serve with avocado or hard boiled eggs or… anything else that would make your lunch life amazing.","ingredientReferences":[]}],"nutrition":{"calories":"322","fatContent":"24.6","proteinContent":"6.1","carbohydrateContent":"24.7","fiberContent":"6.7","sodiumContent":"302.3","sugarContent":"6.9"},"settings":{"public":true,"showNutrition":false,"showAssets":false,"landscapeView":false,"disableComments":false,"disableAmount":false,"locked":false},"assets":[],"notes":[],"extras":{},"comments":[]}`,
			wantRecipe: models.Recipe{
				Category:    "uncategorized",
				CreatedAt:   time.Date(2024, 4, 12, 0, 0, 0, 0, time.UTC),
				Description: "Roasted Vegetable Bowls! Crispy tender roasted veggies, buttery avocado, all together in a bowl with a drizzle of green tahini sauce.",
				Images:      []uuid.UUID{uuid.MustParse("f2f4b3aa-1e04-42a2-a581-607bf84f6800")},
				Ingredients: []string{
					"8 large carrots, peeled and chopped", "3 golden potatoes, chopped",
					"1 head of broccoli, cut into florets",
					"1 head of cauliflower, cut into florets", "olive oil and salt",
					"1/2 cup olive oil (mild tasting)", "1/2 cup water", "1/4 cup tahini",
					"a big bunch of cilantro and/or parsley",
					"1 clove garlic",
					"squeeze of half a lemon (about 2 tablespoons)",
					"1/2 teaspoon salt (more to taste)",
					"6 hard boiled eggs (or other protein)",
					"3 avocados",
				},
				Instructions: []string{
					"Prep", "Preheat the oven to 425 degrees.", "Roasted Vegetables",
					"Arrange your vegetables onto a few baking sheets lined with parchment (I keep each vegetable in its own little section). Toss with olive oil and salt. Roast for 25-30 minutes.", "Sauce",
					"While the veggies are roasting, blitz up your sauce in the food processor or blender.",
					"Finish",
					"Voila! Portion and save for the week! Serve with avocado or hard boiled eggs or… anything else that would make your lunch life amazing.",
				},
				Keywords: []string{"Vegetable Bowl Recipe", "Roasted Vegetable Bowls", "Green Tahini"},
				Name:     "Roasted Vegetable Bowls with Green Tahini",
				Nutrition: models.Nutrition{
					Calories:           "322",
					Fiber:              "6.7",
					Protein:            "6.1",
					Sodium:             "302.3",
					Sugars:             "6.9",
					TotalCarbohydrates: "24.7",
					TotalFat:           "24.6",
				},
				Times:     models.Times{Prep: 15 * time.Minute, Cook: 30 * time.Minute, Total: 45 * time.Minute},
				Tools:     []models.HowToItem{},
				UpdatedAt: time.Date(2024, 04, 12, 0, 0, 0, 0, time.UTC),
				URL:       "https://pinchofyum.com/30-minute-meal-prep-roasted-vegetable-bowls-with-green-tahini",
				Yield:     6,
			},
		},
		{
			name:       "latest API",
			mealieJSON: `{"id":"134372c8-1b49-41e5-a2d2-3249cc7bdb12","user_id":"499aa092-e3fe-46e7-ac5e-9ddf360c9eee","household_id":"67858612-ba5b-470c-bd0d-20b7fdf9a88a","group_id":"0d7fc2c1-58f4-4c88-8177-083263e467c4","name":"Mini cannelés chorizo comté au thermomix","slug":"mini-canneles-chorizo-comte-au-thermomix","image":"wTJi","recipe_servings":24.0,"recipe_yield_quantity":0.0,"recipe_yield":"","total_time":null,"prep_time":"15 minutes","cook_time":null,"perform_time":"54 minutes","description":"Qui ne connait pas le cannelé, ce petit gâteau Bordelais irresistible à pâte molle et tendre. Le voici revisité en version salée pour un apéro légèrement relevé, avec du chorizo et du comté.\r\n\nL'essayer c'est l'adopter :)","recipe_category":[],"tags":[],"tools":[],"rating":4.0,"org_url":"https://www.cookomix.com/recettes/mini-canneles-chorizo-comte-thermomix/","date_added":"2024-12-18","date_updated":"2024-12-31T10:20:45.610286Z","created_at":"2024-12-18T08:18:02.011499Z","updated_at":"2024-12-31T10:20:45.622388Z","last_made":null,"recipe_ingredient":[{"quantity":1.0,"unit":null,"food":null,"note":"Beurre - 30 grammes","is_food":false,"disable_amount":true,"display":"Beurre - 30 grammes","title":null,"original_text":null,"reference_id":"1c728d6b-fd56-4516-a83a-1e742163d879"},{"quantity":1.0,"unit":null,"food":null,"note":"Lait demi-écrémé - 250 ml","is_food":false,"disable_amount":true,"display":"Lait demi-écrémé - 250 ml","title":null,"original_text":null,"reference_id":"e5ea6b2d-9fb3-4c21-97c8-76c0f46ce407"},{"quantity":1.0,"unit":null,"food":null,"note":"Chorizo - 50 grammes","is_food":false,"disable_amount":true,"display":"Chorizo - 50 grammes","title":null,"original_text":null,"reference_id":"04ae6888-705b-4cb7-be99-bf747248422c"},{"quantity":1.0,"unit":null,"food":null,"note":"Comté - 60 grammes","is_food":false,"disable_amount":true,"display":"Comté - 60 grammes","title":null,"original_text":null,"reference_id":"2ebad54d-ab1a-41e6-87af-84196fc7cbec"},{"quantity":1.0,"unit":null,"food":null,"note":"Oeuf - 2","is_food":false,"disable_amount":true,"display":"Oeuf - 2","title":null,"original_text":null,"reference_id":"a621f5f7-5b78-4a63-9564-cc4df260d0f2"},{"quantity":1.0,"unit":null,"food":null,"note":"Farine - 60 grammes","is_food":false,"disable_amount":true,"display":"Farine - 60 grammes","title":null,"original_text":null,"reference_id":"cd53e2eb-8163-4c84-a593-8fcc025a1431"},{"quantity":1.0,"unit":null,"food":null,"note":"Sel - 1 pincée","is_food":false,"disable_amount":true,"display":"Sel - 1 pincée","title":null,"original_text":null,"reference_id":"2928a65b-12be-4ee1-9c6a-32de40c2917e"},{"quantity":1.0,"unit":null,"food":null,"note":"Poivre - 1 pincée","is_food":false,"disable_amount":true,"display":"Poivre - 1 pincée","title":null,"original_text":null,"reference_id":"e42c649c-12ad-44e2-a4b8-87387442f2ff"}],"recipe_instructions":[{"id":"02e8a3ea-1696-4adc-8282-a96be0df0aec","title":"","summary":"","text":"Préchauffer le four à 180°C.","ingredient_references":[]},{"id":"6a841504-b4ad-408a-9281-0acb48c1573e","title":"","summary":"","text":"Dans une casserolle mettre 30 grammes de beurre coupés en morceaux","ingredient_references":[]},{"id":"a67f84b0-3960-4776-b01f-fba47aa20738","title":"","summary":"","text":"Ajouter le lait","ingredient_references":[]},{"id":"a4f49ac8-1280-4c6b-8146-2b62d5b3a706","title":"","summary":"","text":"Faire chauffer à feu moyen en tournant régulièrement","ingredient_references":[]},{"id":"92adc312-0070-4f74-8e1f-6684c352fbfe","title":"","summary":"","text":"Dans un saladier, ajouter 1 oeuf et le jaune d'oeuf","ingredient_references":[]},{"id":"2aa1d7c3-31aa-41f2-a628-c518b72f604e","title":"","summary":"","text":"Ajouter 60 grammes de farine","ingredient_references":[]},{"id":"f77dba58-d336-470e-8232-cbe5cca1a6bc","title":"","summary":"","text":"Ajouter 1 pincée de sel (à ajuster en fonction des goûts)","ingredient_references":[]},{"id":"b4fff326-6e32-4edf-8515-34af3f6b17c6","title":"","summary":"","text":"Ajouter 1 pincée de poivre (à ajuster en fonction des goûts) ","ingredient_references":[]},{"id":"cb4ac495-af3d-4bca-8452-af6631ac31e9","title":"","summary":"","text":"Mélanger pour obtenir une pâte lisse","ingredient_references":[]},{"id":"88e3cc68-74f5-4800-9b57-e6530cfda543","title":"","summary":"","text":"Ajouter progressivement le mélange lait/beurre","ingredient_references":[]},{"id":"9a6c5dae-d4cf-4edd-aadf-c776b5db415a","title":"","summary":"","text":"Ajouter le comté et le chorizo coupés en dés. Bien mélanger","ingredient_references":[]},{"id":"e63371c4-cbfb-4108-a4c3-bd7ec07eebd4","title":"","summary":"","text":"Transvaser dans des moules à mini-cannelés en remplissant aux 3/4.\nLa quantité doit permettre de faire 2 fournées de 12.","ingredient_references":[]},{"id":"d97c56d4-a7a6-4237-9c15-ccf64f715903","title":"","summary":"","text":"Mettre dans le four pendant 25 min à 180°C. Mode chaleur tournante. Adaptez la cuisson en fonction de vos moules et/ou de votre four.","ingredient_references":[]}],"nutrition":{"calories":"37","carbohydrate_content":null,"cholesterol_content":null,"fat_content":null,"fiber_content":null,"protein_content":null,"saturated_fat_content":null,"sodium_content":null,"sugar_content":null,"trans_fat_content":null,"unsaturated_fat_content":null},"settings":{"public":true,"show_nutrition":false,"show_assets":false,"landscape_view":false,"disable_comments":false,"disable_amount":true,"locked":false},"assets":[],"notes":[],"extras":{},"comments":[]}`,
			wantRecipe: models.Recipe{
				Category:    "uncategorized",
				CreatedAt:   time.Date(2024, 12, 18, 0, 0, 0, 0, time.UTC),
				Description: "Qui ne connait pas le cannelé, ce petit gâteau Bordelais irresistible à pâte molle et tendre. Le voici revisité en version salée pour un apéro légèrement relevé, avec du chorizo et du comté.\r\n\nL'essayer c'est l'adopter :)",
				Images:      []uuid.UUID{uuid.MustParse("f2f4b3aa-1e04-42a2-a581-607bf84f6800")},
				Ingredients: []string{
					"Beurre - 30 grammes", "Lait demi-écrémé - 250 ml", "Chorizo - 50 grammes",
					"Comté - 60 grammes", "Oeuf - 2", "Farine - 60 grammes", "Sel - 1 pincée",
					"Poivre - 1 pincée",
				},
				Instructions: []string{
					"Préchauffer le four à 180°C.",
					"Dans une casserolle mettre 30 grammes de beurre coupés en morceaux",
					"Ajouter le lait", "Faire chauffer à feu moyen en tournant régulièrement",
					"Dans un saladier, ajouter 1 oeuf et le jaune d'oeuf",
					"Ajouter 60 grammes de farine",
					"Ajouter 1 pincée de sel (à ajuster en fonction des goûts)",
					"Ajouter 1 pincée de poivre (à ajuster en fonction des goûts) ",
					"Mélanger pour obtenir une pâte lisse",
					"Ajouter progressivement le mélange lait/beurre",
					"Ajouter le comté et le chorizo coupés en dés. Bien mélanger",
					"Transvaser dans des moules à mini-cannelés en remplissant aux 3/4.\nLa quantité doit permettre de faire 2 fournées de 12.",
					"Mettre dans le four pendant 25 min à 180°C. Mode chaleur tournante. Adaptez la cuisson en fonction de vos moules et/ou de votre four.",
				},
				Keywords: []string{},
				Name:     "Mini cannelés chorizo comté au thermomix",
				Nutrition: models.Nutrition{
					Calories: "37",
				},
				Times: models.Times{
					Prep: 15 * time.Minute,
					Cook: 39 * time.Minute,
				},
				Tools:     []models.HowToItem{},
				UpdatedAt: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
				URL:       "https://www.cookomix.com/recettes/mini-canneles-chorizo-comte-thermomix/",
				Yield:     24,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			testMealie(t, tc.mealieJSON, tc.wantRecipe)
		})
	}
}

func testMealie(t testing.TB, mealieRecipe string, wantRecipe models.Recipe) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/auth/token":
			data := `{"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", "token_type": "bearer"}`
			w.WriteHeader(200)
			_, _ = w.Write([]byte(data))
		case "/api/recipes":
			data := `{"page":2,"per_page":100,"total":450,"total_pages":5,"items":[{"id":"a1860172-5b66-4b5e-aca7-0626a55338fe","userId":"e72ff251-4693-4e44-ad1d-9d9c2b033541","groupId":"083bba0c-e400-4b84-8055-b01a888b27fd","name":"Sourdough Crackers","slug":"sourdough-crackers","image":"I1Wd","recipeYield":"20 servings","totalTime":"1 hour 25 minutes","prepTime":"25 minutes","cookTime":null,"performTime":"25 minutes","description":"Here's the perfect solution to your discarded sourdough dilemma. The rosemary, while optional, complements the tang of the sourdough perfectly. We're obsessed with these crackers, especially when dipped into some healthy hummus.","recipeCategory":[],"tags":[{"id":"c66f3a06-40fc-4cde-836e-829d4e79378a","name":"Crackers","slug":"crackers"},{"id":"61cdb5d7-3475-4616-b169-69fca70e9bc9","name":"Whole Grain","slug":"whole-grain"},{"id":"0bf785b0-9867-4501-96c1-fe79ef3b6083","name":"Party Snacks","slug":"party-snacks"}],"tools":[],"rating":null,"orgURL":"https://www.kingarthurbaking.com/recipes/sourdough-crackers-recipe","dateAdded":"2024-04-06","dateUpdated":"2024-04-07T10:38:36.127382","createdAt":"2024-04-06T19:18:04.359452","updateAt":"2024-04-07T10:38:36.130814","lastMade":"2024-04-07T03:59:59"},{"id":"f3f1d73c-df94-485c-ae49-683dd06154fc","userId":"e72ff251-4693-4e44-ad1d-9d9c2b033541","groupId":"083bba0c-e400-4b84-8055-b01a888b27fd","name":"test meal","slug":"test-meal","image":"228","recipeYield":null,"totalTime":null,"prepTime":null,"cookTime":null,"performTime":null,"description":"","recipeCategory":[],"tags":[],"tools":[],"rating":null,"orgURL":null,"dateAdded":"2024-04-06","dateUpdated":"2024-04-07T02:09:39.803633","createdAt":"2024-04-06T15:35:15.739100","updateAt":"2024-04-06T15:35:15.739102","lastMade":null},{"id":"a98bea7c-45bd-4e64-b8f8-ba695bb64393","userId":"e72ff251-4693-4e44-ad1d-9d9c2b033541","groupId":"083bba0c-e400-4b84-8055-b01a888b27fd","name":"testtest","slug":"testtest","image":null,"recipeYield":null,"totalTime":null,"prepTime":null,"cookTime":null,"performTime":null,"description":"","recipeCategory":[],"tags":[],"tools":[],"rating":null,"orgURL":null,"dateAdded":"2024-04-06","dateUpdated":"2024-04-06T15:02:41.149774","createdAt":"2024-04-06T15:02:41.151518","updateAt":"2024-04-06T15:02:41.151520","lastMade":null}],"next":null,"previous":"/recipes?page=1&perPage=100&orderDirection=desc&requireAllCategories=false&requireAllFoods=false&requireAllTags=false&requireAllTools=false"}`
			w.WriteHeader(200)
			_, _ = w.Write([]byte(data))
		default:
			if strings.HasPrefix(r.URL.Path, "/api/recipes/") {
				w.WriteHeader(200)
				_, _ = w.Write([]byte(mealieRecipe))
			} else if strings.HasPrefix(r.URL.Path, "/api/media/recipes/") {
				w.WriteHeader(200)
				_, _ = w.Write([]byte("OK"))
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		}
	}))
	defer srv.Close()
	files := &mockFiles{}

	c := make(chan models.Progress)
	var (
		got models.Recipes
		err error
	)
	go func() {
		defer close(c)
		got, err = integrations.MealieImport(srv.URL, "master", "yoda", srv.Client(), files.UploadImage, c)
		if err != nil {
			return
		}
	}()
	for range c {
	}

	want := models.Recipes{wantRecipe, wantRecipe, wantRecipe}
	assertRecipes(t, got, want, files)
}
