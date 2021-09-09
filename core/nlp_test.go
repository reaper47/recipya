package core

import (
	"testing"

	"github.com/reaper47/recipya/data"
	"github.com/reaper47/recipya/mock"
)

func init() {
	data.Data = &mock.MockDataReader{}
}

func TestNlp(t *testing.T) {
	t.Run(
		"NlpProcessIngredients processes ingredients correctly",
		test_NlpExtractIngredients_ProcessCorrectly,
	)
}

func test_NlpExtractIngredients_ProcessCorrectly(t *testing.T) {
	entries := map[string]string{
		"2 pounds (1 kg) g gal ml cc dl carrots washed and peeled (or unpeeled)": "carrots",
		"1/3 cup butter":                                      "butter",
		"3 tablespoons honey":                                 "honey",
		"4 garlic cloves minced":                              "garlic cloves",
		"cracked black pepper":                                "black pepper",
		"2 tablespoons fresh chopped parsley":                 "parsley",
		"1 (8 ounce) package sliced mushrooms":                "mushrooms",
		"1 large onion, diced":                                "large onion",
		"4 cloves garlic, diced":                              "cloves garlic",
		"salt and ground black pepper to taste":               "salt ground black pepper",
		"1½ pounds beef sirloin steak, cut into 1 inch cubes": "beef sirloin steak",
		"1 cup white wine":                                    "white wine",
		"2 cups beef bouillon":                                "beef bouillon",
		"1 teaspoon crumbled dried thyme":                     "dried thyme",
		"1 teaspoon dried basil":                              "basil",
		"½ teaspoon dried oregano":                            "oregano",
		"2 bay leaves":                                        "bay leaves",
		"½ teaspoon ground black pepper":                      "ground black pepper",
		"¼ cup all-purpose flour":                             "all-purpose flour",
		"2 cups half-and-half cream ":                         "half-and-half cream",
		"1 teaspoon salt":                                     "salt",
		"1 teaspoon baking soda":                              "baking soda",
		"1 teaspoon baking powder":                            "baking powder",
		"3 eggs":                                              "eggs",
		"2¼ cups white sugar":                                 "white sugar",
		"½ cup vegetable oil":                                 "vegetable oil",
		"½ cup canned pumpkin":                                "pumpkin",
		"1 tablespoon vanilla extract":                        "vanilla extract",
		"2½ cups grated zucchini":                             "zucchini",
		"1 cup chopped walnuts (Optional)":                    "walnuts",
		"2 tablespoons vegetable oil":                         "vegetable oil",
		"1 cup sliced fresh mushrooms":                        "mushrooms",
		"1 cup snow peas":                                     "snow peas",
		"¾ cup shredded carrot":                               "carrot",
		"4 medium green onions, sliced":                       "medium green onions",
		"1 clove garlic, minced":                              "clove garlic",
		"¼ cup reduced-sodium soy sauce":                      "reduced-sodium soy sauce",
		"1 teaspoon white sugar":                              "white sugar",
		"¼ teaspoon cayenne pepper":                           "cayenne pepper",
		"8 ounces cooked spaghetti":                           "spaghetti",
		"1 tablespoon toasted sesame seeds":                   "sesame seeds",
		"3 cups all-purpose flour":                            "all-purpose flour",
		"1 tablespoon ground cinnamon":                        "ground cinnamon",
		"3 teaspoons vanilla extract":                         "vanilla extract",
		"1 cup butter, softened":                              "butter",
		"1 cup white sugar":                                   "white sugar",
		"1 cup packed brown sugar":                            "brown sugar",
		"2 cups semisweet chocolate chips":                    "semisweet chocolate chips",
	}
	var ingredients, expected []string
	for key, value := range entries {
		ingredients = append(ingredients, key)
		expected = append(expected, value)
	}

	actual := NlpExtractIngredients(ingredients)

	fail := false
	if len(actual) != len(expected) {
		t.Fatalf("slice lengths unequal: got %v want %v", len(actual), len(expected))
	}
	for i := range actual {
		if actual[i] != expected[i] {
			fail = true
			t.Logf("ingredient is unequal: got '%v' want '%v'", actual[i], expected[i])
		}
	}
	if fail {
		t.Fail()
	}
}
