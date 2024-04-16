package models_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"testing"
	"time"
)

func TestPaprikaRecipe_Recipe(t *testing.T) {
	data := `{"uid": "CF0EB3F6-A5AE-49CC-BDEF-D5B406F6C34C","created": "2024-04-09 05:42:31","hash": "08821A86BC2BE19727844C3A2A59C002A70B979DF1A3CE4A64DB5AD24233D611","name": "Guyanese Gojas","description": "","ingredients": "For the goja dough:\n2 cups (240g) all-purpose flour\n2 teaspoons (8g) granulated sugar\n1/2 teaspoon (2g) instant yeast\n2 tablespoons diced cold butter, or vegetable shortening\n1 cup cold whole milk\n1 tablespoon (15g) all-purpose flour, to sprinkle while kneading\n1/8 teaspoon neutral cooking oil, such as canola or vegetable for rubbing the dough ball\nFor the goja filling:\n2 cups sweetened flaked coconut\n1/2 teaspoon ground cinnamon\n1/2 teaspoon freshly grated nutmeg\n1 tablespoon light brown sugar\n1 tablespoon fresh grated ginger\n2 tablespoons water\n2 tablespoons butter, melted\n2 teaspoons vanilla extract\n2 tablespoons neutral cooking oil such as canola or vegetable, for cooking filling\n1/4 cup raisins\n2 tablespoons water, for cooking filling\nFor shaping and frying the gojas:\n1/4 cup water, for sealing the pastry\n1/4 cup (34g) all-purpose flour, for crimping the gojas\n3 cups canola or vegetable oil, for frying","directions": "Combine the dry ingredients :\n\nIn a large bowl, combine the flour, sugar, and yeast. Add the butter. Using your hands or fingertips, rub the butter into flour until a coarse meal forms.\n\nAlica Ramkirpal-Senhouse / Simply Recipes\n\nAlica Ramkirpal-Senhouse / Simply Recipes\n\nPour milk into bowl:\n\nMake a well in the center of the bowl, pour in the milk. Using a rubber spatula, stir until the dough forms. At this point, the dough will be a little sticky, sprinkle 1 tablespoon flour on the dough and knead it into dough with your hands in the bowl until the dough is no longer sticky.\n\nSet dough aside to rest:\n\nRub the top of the dough with the oil and cover with a damp paper towel. Set aside for 15 to 20 minutes.\n\nMake the filling:\n\nIn the bowl of your food processor, add the coconut, cinnamon, nutmeg, brown sugar, ginger, water, butter, and vanilla. Pulse on high until coconut becomes fine and pasty.\n\nCook the filling:\n\nHeat a heavy-bottomed pan over low heat. Add the oil, coconut filling, raisins, and 2 tablespoons of water. Cook, stirring occasionally, until the sugar melts and the coconut looks more toasted and slightly darker in color, about 5 minutes. Remove from heat and let cool for a few minutes before assembling the gojas.\n\nAlica Ramkirpal-Senhouse / Simply Recipes\n\nWeigh and divide the dough:\n\nWeigh the dough, then divide the weight by 12 to get the weight for each piece. Now, cut 12 small pieces of dough and weigh each. Add or remove small pieces until you get the exact weight you’re looking for.\n\nIf you’re not using a scale, divide the dough into 12 pieces using a knife or pastry cutter. Try to eyeball it so they’re all the same size.\n\nAlica Ramkirpal-Senhouse / Simply Recipes\n\nRoll the goja dough:\n\nRound off each dough ball between your palms to form a ball, gently tucking dough under itself to make the top smooth. Once you’ve done this, cover all the dough balls with a damp paper towel to keep it from drying out and crusting.\n\nSprinkle flour on the surface of the dough ball you are working with. Working with one dough ball at a time, flatten slightly with your hands, then roll into a circle 1/8 inch thick and about 5 inches in diameter.\n\nFlour your surface as needed as you go along.\n\nRepeat with remaining balls of dough, being sure to keep them covered as you work.\n\nDip your pointer finger in water and run it around the outer edges of the dough. Place 2 tablespoons filing in the bottom half of the dough and bring the top half over to seal. Using a fork crimp the edges closed being sure to dip the fork in flour to keep from sticking while crimping. Place assembled gojas on a baking sheet lined with parchment paper.\n\nRepeat this step for the rest of the batch.\n\nSet up a plate or deep serving platter with a few paper towels to place gojas on after they’re done frying.\n\nHeat a medium sized deep pot over medium-low heat. Add the oil and once it’s anywhere between 350-375°F, fry the gojas for 2 to 3 minutes, you’ll have to cook these in batches, being sure to not overcrowd the pot. Use a slotted spoon or tongs to flip the gojas once halfway through cooking. Remove from oil once it is light golden brown and drain on paper towels.\n\nRepeat with remaining gojas until they are all fried.\n\nEnjoy warm.\n\nAlica Ramkirpal-Senhouse / Simply Recipes","notes": "","nutritional_info": "(per serving)\n543 Calories 30g Fat 62g Carbs 8g Protein\nNutrition Facts\nServings: 6\nAmount per serving\nCalories 543\n% Daily Value*\nTotal Fat 30g 38%\nSaturated Fat 13g 67%\nCholesterol 17mg 6%\nSodium 132mg 6%\nTotal Carbohydrate 62g 23%\nDietary Fiber 5g 16%\nTotal Sugars 20g\nProtein 8g\nVitamin C 0mg 1%\nCalcium 66mg 5%\nIron 3mg 16%\nPotassium 271mg 6%\n*The % Daily Value (DV) tells you how much a nutrient in a food serving contributes to a daily diet. 2,000 calories a day is used for general nutrition advice.","prep_time": "65 mins","cook_time": "35 mins","total_time": "","difficulty": "Medium","servings": "6 servings","rating": 0,"source": "Simplyrecipes.com","source_url": "https://www.simplyrecipes.com/guyanese-gojas-recipe-5221034","photo": "57C135F8-0CF2-40F2-89C0-0F7BCFFCF0F4.jpg","photo_large": null,"photo_hash": "3F1AB2130588CB3AEA727B0AAF0F24358921197A12036CF283BEBA9E0C60102C","image_url": "https://www.simplyrecipes.com/thmb/07pz2tqpdKhQdNId5mCQVE58fVw=/750x0/filters:no_upscale():max_bytes(150000):strip_icc():format(webp)/Simply-Recipes-Guyanese-Gojas-LEAD-09-712f8f32ffff41b4b4cdd2f91c0f8b74.jpg","categories": [],"photos": []}`
	var p models.PaprikaRecipe
	_ = json.Unmarshal([]byte(data), &p)

	got := p.Recipe(uuid.Nil)

	want := models.Recipe{
		Category:    "uncategorized",
		CreatedAt:   got.CreatedAt,
		Description: "Imported from Paprika",
		Image:       uuid.Nil,
		Ingredients: []string{
			"For the goja dough:", "2 cups (240g) all-purpose flour",
			"2 teaspoons (8g) granulated sugar", "1/2 teaspoon (2g) instant yeast",
			"2 tablespoons diced cold butter, or vegetable shortening",
			"1 cup cold whole milk",
			"1 tablespoon (15g) all-purpose flour, to sprinkle while kneading",
			"1/8 teaspoon neutral cooking oil, such as canola or vegetable for rubbing the dough ball",
			"For the goja filling:", "2 cups sweetened flaked coconut",
			"1/2 teaspoon ground cinnamon", "1/2 teaspoon freshly grated nutmeg",
			"1 tablespoon light brown sugar", "1 tablespoon fresh grated ginger",
			"2 tablespoons water", "2 tablespoons butter, melted",
			"2 teaspoons vanilla extract",
			"2 tablespoons neutral cooking oil such as canola or vegetable, for cooking filling",
			"1/4 cup raisins",
			"2 tablespoons water, for cooking filling",
			"For shaping and frying the gojas:",
			"1/4 cup water, for sealing the pastry",
			"1/4 cup (34g) all-purpose flour, for crimping the gojas",
			"3 cups canola or vegetable oil, for frying",
		},
		Instructions: []string{
			"Combine the dry ingredients :",
			"In a large bowl, combine the flour, sugar, and yeast. Add the butter. Using your hands or fingertips, rub the butter into flour until a coarse meal forms.",
			"Alica Ramkirpal-Senhouse / Simply Recipes",
			"Alica Ramkirpal-Senhouse / Simply Recipes", "Pour milk into bowl:",
			"Make a well in the center of the bowl, pour in the milk. Using a rubber spatula, stir until the dough forms. At this point, the dough will be a little sticky, sprinkle 1 tablespoon flour on the dough and knead it into dough with your hands in the bowl until the dough is no longer sticky.",
			"Set dough aside to rest:",
			"Rub the top of the dough with the oil and cover with a damp paper towel. Set aside for 15 to 20 minutes.",
			"Make the filling:",
			"In the bowl of your food processor, add the coconut, cinnamon, nutmeg, brown sugar, ginger, water, butter, and vanilla. Pulse on high until coconut becomes fine and pasty.",
			"Cook the filling:",
			"Heat a heavy-bottomed pan over low heat. Add the oil, coconut filling, raisins, and 2 tablespoons of water. Cook, stirring occasionally, until the sugar melts and the coconut looks more toasted and slightly darker in color, about 5 minutes. Remove from heat and let cool for a few minutes before assembling the gojas.",
			"Alica Ramkirpal-Senhouse / Simply Recipes",
			"Weigh and divide the dough:",
			"Weigh the dough, then divide the weight by 12 to get the weight for each piece. Now, cut 12 small pieces of dough and weigh each. Add or remove small pieces until you get the exact weight you’re looking for.",
			"If you’re not using a scale, divide the dough into 12 pieces using a knife or pastry cutter. Try to eyeball it so they’re all the same size.",
			"Alica Ramkirpal-Senhouse / Simply Recipes",
			"Roll the goja dough:",
			"Round off each dough ball between your palms to form a ball, gently tucking dough under itself to make the top smooth. Once you’ve done this, cover all the dough balls with a damp paper towel to keep it from drying out and crusting.",
			"Sprinkle flour on the surface of the dough ball you are working with. Working with one dough ball at a time, flatten slightly with your hands, then roll into a circle 1/8 inch thick and about 5 inches in diameter.",
			"Flour your surface as needed as you go along.",
			"Repeat with remaining balls of dough, being sure to keep them covered as you work.",
			"Dip your pointer finger in water and run it around the outer edges of the dough. Place 2 tablespoons filing in the bottom half of the dough and bring the top half over to seal. Using a fork crimp the edges closed being sure to dip the fork in flour to keep from sticking while crimping. Place assembled gojas on a baking sheet lined with parchment paper.",
			"Repeat this step for the rest of the batch.",
			"Set up a plate or deep serving platter with a few paper towels to place gojas on after they’re done frying.",
			"Heat a medium sized deep pot over medium-low heat. Add the oil and once it’s anywhere between 350-375°F, fry the gojas for 2 to 3 minutes, you’ll have to cook these in batches, being sure to not overcrowd the pot. Use a slotted spoon or tongs to flip the gojas once halfway through cooking. Remove from oil once it is light golden brown and drain on paper towels.",
			"Repeat with remaining gojas until they are all fried.",
			"Enjoy warm.",
			"Alica Ramkirpal-Senhouse / Simply Recipes",
		},
		Keywords:  []string{"paprika"},
		Name:      "Guyanese Gojas",
		Times:     models.Times{Prep: 1*time.Hour + 5*time.Minute, Cook: 35 * time.Minute},
		Tools:     []string{},
		UpdatedAt: got.UpdatedAt,
		URL:       "https://www.simplyrecipes.com/guyanese-gojas-recipe-5221034",
		Yield:     6,
	}
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}
