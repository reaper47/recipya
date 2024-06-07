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
				data := `{"id":"843b4a6d-6855-48c3-8186-22f096310243","userId":"e72ff251-4693-4e44-ad1d-9d9c2b033541","groupId":"083bba0c-e400-4b84-8055-b01a888b27fd","name":"Roasted Vegetable Bowls with Green Tahini","slug":"roasted-vegetable-bowls-with-green-tahini","image":"Z4Ox","recipeYield":"6 servings","totalTime":"45 minutes","prepTime":"15 minutes","cookTime":null,"performTime":"30 minutes","description":"Roasted Vegetable Bowls! Crispy tender roasted veggies, buttery avocado, all together in a bowl with a drizzle of green tahini sauce.","recipeCategory":[],"tags":[{"id":"70cc7ab9-cc6f-41d0-b8b9-8d16384f857e","name":"Vegetable Bowl Recipe","slug":"vegetable-bowl-recipe"},{"id":"da629ccc-56cb-4400-bce1-55ca0f14905b","name":"Roasted Vegetable Bowls","slug":"roasted-vegetable-bowls"},{"id":"e3184b1f-2bd0-48b8-b766-fa3eaf8285a5","name":"Green Tahini","slug":"green-tahini"}],"tools":[],"rating":4,"orgURL":"https://pinchofyum.com/30-minute-meal-prep-roasted-vegetable-bowls-with-green-tahini","dateAdded":"2024-04-12","dateUpdated":"2024-04-12T18:14:29.168064","createdAt":"2024-04-12T18:06:06.692275","updateAt":"2024-04-12T18:07:56.850947","lastMade":null,"recipeIngredient":[{"quantity":8.0,"unit":null,"food":{"id":"31d502b5-dea9-4580-b8b2-86bfde80f456","name":"carrot","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.825124","updateAt":"2024-03-05T06:54:56.825126"},"note":"large, peeled and chopped","isFood":true,"disableAmount":false,"display":"8 carrot large, peeled and chopped","title":null,"originalText":"8 large carrots, peeled and chopped","referenceId":"0b7f622f-a6f1-4e51-b381-d665ee54da47"},{"quantity":3.0,"unit":null,"food":null,"note":"chopped","isFood":true,"disableAmount":false,"display":"3 chopped","title":null,"originalText":"3 golden potatoes, chopped","referenceId":"d0944e99-f189-48aa-b951-1d2c1aaf7655"},{"quantity":1.0,"unit":null,"food":{"id":"05bb0bf7-dc26-4961-9bce-bd5563f7a6c7","name":"broccoli","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.552679","updateAt":"2024-03-05T06:54:56.552681"},"note":"cut into florets","isFood":true,"disableAmount":false,"display":"1 broccoli cut into florets","title":null,"originalText":"1 head of broccoli, cut into florets","referenceId":"3b1e5c67-5f02-49b8-88d7-6db8a681de05"},{"quantity":1.0,"unit":null,"food":{"id":"1a0beaa6-b6f2-4143-81e9-6709fe00d33a","name":"cauliflower","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.214977","updateAt":"2024-03-05T06:54:56.214981"},"note":"cut into florets","isFood":true,"disableAmount":false,"display":"1 cauliflower cut into florets","title":null,"originalText":"1 head of cauliflower, cut into florets","referenceId":"c0836160-4afd-4157-a780-7ac0a7f41a6f"},{"quantity":0.0,"unit":null,"food":null,"note":"","isFood":true,"disableAmount":false,"display":"","title":null,"originalText":"olive oil and salt","referenceId":"ac21685f-510d-4534-b4cd-e0ac57e04a04"},{"quantity":0.5,"unit":{"id":"56939576-ff3a-4760-98c4-cd7e7ea8418b","name":"cup","pluralName":null,"description":"","extras":{},"fraction":true,"abbreviation":"","pluralAbbreviation":"","useAbbreviation":false,"aliases":[],"createdAt":"2024-03-05T06:59:19.679450","updateAt":"2024-03-05T06:59:19.679454"},"food":{"id":"02bc4201-08ca-45b2-b032-7babfa4346f4","name":"olive oil","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.572704","updateAt":"2024-03-05T06:54:56.572706"},"note":"mild tasting)","isFood":true,"disableAmount":false,"display":"¹/₂ cup olive oil mild tasting)","title":null,"originalText":"1/2 cup olive oil (mild tasting)","referenceId":"af426c53-d5e4-4ab0-bc17-6b80ab2d389d"},{"quantity":0.5,"unit":{"id":"56939576-ff3a-4760-98c4-cd7e7ea8418b","name":"cup","pluralName":null,"description":"","extras":{},"fraction":true,"abbreviation":"","pluralAbbreviation":"","useAbbreviation":false,"aliases":[],"createdAt":"2024-03-05T06:59:19.679450","updateAt":"2024-03-05T06:59:19.679454"},"food":{"id":"c30d5cf5-d7e6-4f8b-8338-a010cae94441","name":"water","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:59:33.037532","updateAt":"2024-03-05T06:59:33.037536"},"note":"","isFood":true,"disableAmount":false,"display":"¹/₂ cup water","title":null,"originalText":"1/2 cup water","referenceId":"5cc3edfb-11c3-403f-afc0-f10d7863791c"},{"quantity":0.25,"unit":{"id":"56939576-ff3a-4760-98c4-cd7e7ea8418b","name":"cup","pluralName":null,"description":"","extras":{},"fraction":true,"abbreviation":"","pluralAbbreviation":"","useAbbreviation":false,"aliases":[],"createdAt":"2024-03-05T06:59:19.679450","updateAt":"2024-03-05T06:59:19.679454"},"food":{"id":"337615d4-f263-4289-80e2-7fe79983c29e","name":"tahini","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.942480","updateAt":"2024-03-05T06:54:56.942482"},"note":"","isFood":true,"disableAmount":false,"display":"¹/₄ cup tahini","title":null,"originalText":"1/4 cup tahini","referenceId":"e77a5636-f06c-479d-932e-86920ff04ae3"},{"quantity":0.0,"unit":null,"food":{"id":"3fce1ca1-3fbe-4a29-b0c8-6411a581be4d","name":"cilantro","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.640255","updateAt":"2024-03-05T06:54:56.640257"},"note":"and/or parsley","isFood":true,"disableAmount":false,"display":"cilantro and/or parsley","title":null,"originalText":"a big bunch of cilantro and/or parsley","referenceId":"eee40c2c-ce2d-42db-abd7-2743db1a017e"},{"quantity":1.0,"unit":{"id":"b9ca3f9e-f7d5-4bce-8181-29ae56939ce6","name":"clove","pluralName":null,"description":"","extras":{},"fraction":true,"abbreviation":"","pluralAbbreviation":"","useAbbreviation":false,"aliases":[],"createdAt":"2024-03-05T08:19:26.421964","updateAt":"2024-03-05T08:19:26.421968"},"food":{"id":"7f2f8ad6-035d-46c7-8503-9e6bf7281165","name":"garlic","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.355030","updateAt":"2024-03-05T06:54:56.355032"},"note":"","isFood":true,"disableAmount":false,"display":"1 clove garlic","title":null,"originalText":"1 clove garlic","referenceId":"84153160-0675-49f6-9363-0f255b0fdf1c"},{"quantity":0.0,"unit":null,"food":null,"note":"(about 2 tablespoons)","isFood":true,"disableAmount":false,"display":"(about 2 tablespoons)","title":null,"originalText":"squeeze of half a lemon (about 2 tablespoons)","referenceId":"777d50a9-ad50-4b03-a5e8-331ae6cc94f1"},{"quantity":0.5,"unit":{"id":"cda9b5eb-21c5-4acf-b65f-7b397e560eb3","name":"teaspoon","pluralName":null,"description":"","extras":{},"fraction":true,"abbreviation":"","pluralAbbreviation":"","useAbbreviation":false,"aliases":[],"createdAt":"2024-03-05T08:19:34.732459","updateAt":"2024-03-05T08:19:34.732462"},"food":{"id":"90a64fd5-ce9d-4774-a1c6-68c65ed5afea","name":"salt","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.265153","updateAt":"2024-03-05T06:54:56.265155"},"note":"","isFood":true,"disableAmount":false,"display":"¹/₂ teaspoon salt","title":null,"originalText":"1/2 teaspoon salt (more to taste)","referenceId":"cd6dc996-55bf-40a7-a280-c6b754157f2d"},{"quantity":6.0,"unit":null,"food":{"id":"b72ab124-26bc-40f1-a59c-3d469fe890b1","name":"eggs","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.442542","updateAt":"2024-03-05T06:54:56.442544"},"note":"hard boiled (or other protein)","isFood":true,"disableAmount":false,"display":"6 eggs hard boiled (or other protein)","title":null,"originalText":"6 hard boiled eggs (or other protein)","referenceId":"b56905f5-1d19-4d20-8498-3fd9285bf8c3"},{"quantity":3.0,"unit":null,"food":{"id":"d20ca2a2-f869-4680-937a-e00349b34893","name":"avocado","pluralName":null,"description":"","extras":{},"labelId":null,"aliases":[],"label":null,"createdAt":"2024-03-05T06:54:56.787485","updateAt":"2024-03-05T06:54:56.787488"},"note":"","isFood":true,"disableAmount":false,"display":"3 avocado","title":null,"originalText":"3 avocados","referenceId":"4e2cbbd8-43b0-4f37-8247-00b547ae99e3"}],"recipeInstructions":[{"id":"4f6932fd-1da1-4264-a7b3-fc8c7f8bcec1","title":"","text":"Prep","ingredientReferences":[]},{"id":"668767d9-1f73-4283-b450-fed81bb9521a","title":"","text":"Preheat the oven to 425 degrees.","ingredientReferences":[]},{"id":"ce35bb0f-aac7-4396-b4c6-73e28a8bf1bc","title":"","text":"Roasted Vegetables","ingredientReferences":[]},{"id":"cc60b45f-014c-4323-9767-4e3e957dead0","title":"","text":"Arrange your vegetables onto a few baking sheets lined with parchment (I keep each vegetable in its own little section). Toss with olive oil and salt. Roast for 25-30 minutes.","ingredientReferences":[]},{"id":"517eb65d-dc49-4d63-9a2d-82dbdee8f728","title":"","text":"Sauce","ingredientReferences":[]},{"id":"039c6817-4e7f-4690-9c84-20574bec62a4","title":"","text":"While the veggies are roasting, blitz up your sauce in the food processor or blender.","ingredientReferences":[]},{"id":"a78f9737-4b5b-4798-aa48-9f087d63c535","title":"","text":"Finish","ingredientReferences":[]},{"id":"ade3a45d-1546-47ca-ab0b-441707c58429","title":"","text":"Voila! Portion and save for the week! Serve with avocado or hard boiled eggs or… anything else that would make your lunch life amazing.","ingredientReferences":[]}],"nutrition":{"calories":"322","fatContent":"24.6","proteinContent":"6.1","carbohydrateContent":"24.7","fiberContent":"6.7","sodiumContent":"302.3","sugarContent":"6.9"},"settings":{"public":true,"showNutrition":false,"showAssets":false,"landscapeView":false,"disableComments":false,"disableAmount":false,"locked":false},"assets":[],"notes":[],"extras":{},"comments":[]}`
				w.WriteHeader(200)
				_, _ = w.Write([]byte(data))
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

	img, _ := uuid.Parse("f2f4b3aa-1e04-42a2-a581-607bf84f6800")
	r := models.Recipe{
		Category:    "uncategorized",
		CreatedAt:   time.Date(2024, 4, 12, 0, 0, 0, 0, time.UTC),
		Description: "Roasted Vegetable Bowls! Crispy tender roasted veggies, buttery avocado, all together in a bowl with a drizzle of green tahini sauce.",
		Images:      []uuid.UUID{img},
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
			Calories:           "322 kcal",
			Fiber:              "6.7g",
			Protein:            "6.1g",
			Sodium:             "302.3g",
			Sugars:             "6.9g",
			TotalCarbohydrates: "24.7g",
			TotalFat:           "24.6g",
		},
		Times:     models.Times{Prep: 15 * time.Minute, Cook: 30 * time.Minute},
		Tools:     []models.HowToItem{},
		UpdatedAt: time.Date(2024, 04, 12, 0, 0, 0, 0, time.UTC),
		URL:       "https://pinchofyum.com/30-minute-meal-prep-roasted-vegetable-bowls-with-green-tahini",
		Yield:     6,
	}
	want := models.Recipes{r, r, r}
	assertRecipes(t, got, want, files)
}
