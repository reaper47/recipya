package integrations_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/integrations"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTandoorImport(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api-token-auth/":
			data := `{"id": 1995, "token": "tda_7648db14_0145_4007_a0e7_edb80b43e7e7", "scope": "read write app", "expires": "2029-03-21T07:19:38.245403Z", "user_id": 2, "test": 2}`
			w.WriteHeader(200)
			_, _ = w.Write([]byte(data))
		case "/api/recipe/":
			data := `{"count":10,"next":null,"previous":null,"results":[{"id":36588,"name":"Asiatisch gebratene Nudeln süß-scharf, vegetarisch","description":"Asiatisch gebratene Nudeln süß-scharf, vegetarisch. Über 285 Bewertungen und für super befunden. Mit ► Portionsrechner ► Kochbuch ► Video-Tipps! Jetzt entdecken und ausprobieren!","image":"https://tandoor-media.s3.amazonaws.com/recipes/39780d95-9fcf-4a93-9049-55a3e011574c_36588.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAVZWN7ZULN5FV7LEM%2F20240415%2Feu-central-1%2Fs3%2Faws4_request&X-Amz-Date=20240415T181400Z&X-Amz-Expires=3600&X-Amz-SignedHeaders=host&X-Amz-Signature=3924df81b351a05b661a17a04bf46431c95f33b45e307db07f6464b4455372ec","keywords":[],"working_time":20,"waiting_time":5,"created_by":2,"created_at":"2023-06-08T23:48:16.232326+02:00","updated_at":"2023-06-08T23:48:18.257060+02:00","internal":true,"servings":4,"servings_text":"Portion(en)","rating":null,"last_cooked":null,"new":false},{"id":40182,"name":"fg","description":null,"image":null,"keywords":[],"working_time":0,"waiting_time":0,"created_by":2,"created_at":"2023-08-17T20:56:20.312998+02:00","updated_at":"2023-08-17T20:56:20.333831+02:00","internal":true,"servings":1,"servings_text":"","rating":null,"last_cooked":null,"new":false},{"id":36566,"name":"Hearty Steak and Potatoes with Balsamic-Cranberry Pan Sauce","description":null,"image":"https://tandoor-media.s3.amazonaws.com/recipes/de4ce656-1289-44e9-a65a-e28d8997c4ab_36566.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAVZWN7ZULN5FV7LEM%2F20240415%2Feu-central-1%2Fs3%2Faws4_request&X-Amz-Date=20240415T181400Z&X-Amz-Expires=3600&X-Amz-SignedHeaders=host&X-Amz-Signature=a253c42c50c8c82ab4d601ec64c5ddab7c106ab919f0feeadbc154013ed47e0c","keywords":[],"working_time":40,"waiting_time":0,"created_by":2,"created_at":"2023-06-07T12:48:26.297623+02:00","updated_at":"2023-06-07T12:48:30.037266+02:00","internal":true,"servings":2,"servings_text":"2","rating":3.8,"last_cooked":"2023-11-23T16:50:00+01:00","new":false}]}`
			w.WriteHeader(200)
			_, _ = w.Write([]byte(data))
		default:
			if strings.HasPrefix(r.URL.Path, "/api/recipe/") {
				data := `{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil","description":null,"image":"","keywords":[{"id":3098,"name":"vegan","label":"vegan","description":"","parent":null,"numchild":0,"created_at":"2022-01-21T17:50:01.273297+01:00","updated_at":"2024-04-07T12:00:58.050370+02:00","full_name":"vegan"},{"id":3777,"name":"lunch","label":"lunch","description":"","parent":null,"numchild":0,"created_at":"2022-02-09T10:02:46.922926+01:00","updated_at":"2024-03-29T17:56:44.467565+01:00","full_name":"lunch"},{"id":4195,"name":" dinner","label":" dinner","description":"","parent":null,"numchild":0,"created_at":"2022-03-10T17:36:18.305741+01:00","updated_at":"2023-12-15T16:53:41.528797+01:00","full_name":" dinner"},{"id":5149,"name":" vegetarian","label":" vegetarian","description":"","parent":null,"numchild":0,"created_at":"2022-05-21T09:08:25.154172+02:00","updated_at":"2023-12-11T18:52:57.232193+01:00","full_name":" vegetarian"},{"id":5171,"name":" asian","label":" asian","description":"","parent":null,"numchild":0,"created_at":"2022-05-21T09:08:43.293171+02:00","updated_at":"2023-11-11T22:08:23.930709+01:00","full_name":" asian"},{"id":5324,"name":" gluten free","label":" gluten free","description":"","parent":null,"numchild":0,"created_at":"2022-05-29T07:46:30.812517+02:00","updated_at":"2023-12-11T10:22:41.345330+01:00","full_name":" gluten free"}],"steps":[{"id":21081,"name":"","instruction":"*Wonderful vegan stir fry noodles made in just 30 minutes for the perfect weeknight dinner! These easy sesame noodles have a homemade stir fry sauce and a boost of plant-based protein from chickpeas. Top with fresh herbs and roasted cashews for a beautiful meal you'll make again and again.*  \\n\\nFirst make your stir fry sauce: in a medium bowl, whisk together the soy sauce, water, garlic, coconut sugar, sesame oil, rice vinegar, fresh ginger, sesame seeds, red pepper flakes and arrowroot starch (or cornstarch). Set aside.  \\nAdd 1 tablespoon sesame oil to a large pot then add in chopped onion and sliced carrots and cook for 2-4 minutes until onions begin to soften. Next add in broccoli and bell pepper and cook, stirring frequently, for an additional 6-8 minutes or until broccoli is slightly tender but still have a bite.  \\nWhile the veggies are cooking, make your stir fry rice noodles according to the directions on the package. Then drain and set aside.  \\nAdd the drained chickpeas to the pot with the cooked veggies. Immediately turn the heat to low and add in the sauce. Cook for an additional 2 minutes over low heat until the sauce begins to thicken a bit. It should be nice and saucy. Stir in rice noodles, fresh basil ribbons and cashews; toss again to combine. Garnish with scallions and cilantro. Serves 4.","ingredients":[{"id":99600,"food":{"id":34718,"name":"For the Sauce:","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"For the Sauce:","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":null,"amount":0,"conversions":[],"note":"","order":0,"is_header":false,"no_amount":false,"original_text":"For the sauce:","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99601,"food":{"id":41165,"name":"low sodium soy sauce or coconut aminos","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"low sodium soy sauce or coconut aminos","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3910,"name":"cup","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":0.33,"conversions":[{"food":"low sodium soy sauce or coconut aminos","unit":"cup","amount":0.3333333333333333}],"note":"","order":0,"is_header":false,"no_amount":false,"original_text":"1/3 cup low sodium soy sauce or coconut aminos","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99602,"food":{"id":18847,"name":"water","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"water","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3910,"name":"cup","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":0.33,"conversions":[{"food":"water","unit":"cup","amount":0.3333333333333333}],"note":"","order":0,"is_header":false,"no_amount":false,"original_text":"⅓ cup water","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99603,"food":{"id":19243,"name":"garlic","plural_name":null,"description":"","recipe":null,"url":"","properties":[{"id":27360,"property_amount":0,"property_type":{"id":9719,"name":"Sugar","unit":"g","description":"Zucker","order":0,"open_data_slug":null,"fdc_id":1011}},{"id":27373,"property_amount":0.38,"property_type":{"id":4249,"name":"Fat","unit":"g","description":"","order":0,"open_data_slug":null,"fdc_id":1004}},{"id":27371,"property_amount":28.2,"property_type":{"id":4250,"name":"Carbohydrates","unit":"g","description":"","order":0,"open_data_slug":null,"fdc_id":1005}},{"id":27374,"property_amount":6.62,"property_type":{"id":4251,"name":"Proteins","unit":"g","description":"","order":0,"open_data_slug":null,"fdc_id":1003}},{"id":27372,"property_amount":143,"property_type":{"id":4252,"name":"Calories","unit":"kcal","description":"","order":0,"open_data_slug":null,"fdc_id":1008}}],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":1104647,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"garlic","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":4086,"name":"cloves","plural_name":null,"description":null,"base_unit":null,"open_data_slug":null},"amount":3,"conversions":[{"food":"garlic","unit":"cloves","amount":3}],"note":"minced","order":0,"is_header":false,"no_amount":false,"original_text":"3 cloves garlic, minced","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99604,"food":{"id":41166,"name":"coconut sugar or brown sugar","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"coconut sugar or brown sugar","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3937,"name":"tablespoons","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":2,"conversions":[{"food":"coconut sugar or brown sugar","unit":"tablespoons","amount":2}],"note":"or sub 1 tablespoon pure maple syrup","order":0,"is_header":false,"no_amount":false,"original_text":"2 tablespoons coconut sugar or brown sugar (or sub 1 tablespoon pure maple syrup)","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99605,"food":{"id":19244,"name":"sesame oil","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"sesame oil","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3908,"name":"tablespoon","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":1,"conversions":[{"food":"sesame oil","unit":"tablespoon","amount":1}],"note":"preferably toasted sesame oil","order":0,"is_header":false,"no_amount":false,"original_text":"1 tablespoon sesame oil (preferably toasted sesame oil)","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99606,"food":{"id":32672,"name":"rice vinegar","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"rice vinegar","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3908,"name":"tablespoon","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":1,"conversions":[{"food":"rice vinegar","unit":"tablespoon","amount":1}],"note":"","order":0,"is_header":false,"no_amount":false,"original_text":"1 tablespoon rice vinegar","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99607,"food":{"id":37716,"name":"fresh grated ginger","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"fresh grated ginger","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3908,"name":"tablespoon","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":1,"conversions":[{"food":"fresh grated ginger","unit":"tablespoon","amount":1}],"note":"","order":0,"is_header":false,"no_amount":false,"original_text":"1 tablespoon fresh grated ginger","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99608,"food":{"id":19240,"name":"sesame seeds","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"sesame seeds","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3908,"name":"tablespoon","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":1,"conversions":[{"food":"sesame seeds","unit":"tablespoon","amount":1}],"note":"or sub 1 tablespoon tahini","order":0,"is_header":false,"no_amount":false,"original_text":"1 tablespoon sesame seeds (or sub 1 tablespoon tahini)","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99609,"food":{"id":22730,"name":"red pepper flakes","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"red pepper flakes","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3909,"name":"teaspoon","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":0.5,"conversions":[{"food":"red pepper flakes","unit":"teaspoon","amount":0.5}],"note":"","order":0,"is_header":false,"no_amount":false,"original_text":"½ teaspoon red pepper flakes","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99610,"food":{"id":41167,"name":"arrowroot starch","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"arrowroot starch","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3908,"name":"tablespoon","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":0.5,"conversions":[{"food":"arrowroot starch","unit":"tablespoon","amount":0.5}],"note":"or sub cornstarch","order":0,"is_header":false,"no_amount":false,"original_text":"½ tablespoon arrowroot starch (or sub cornstarch)","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99611,"food":{"id":41168,"name":"For the veggies & chickpeas:","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"For the veggies & chickpeas:","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":null,"amount":0,"conversions":[],"note":"","order":0,"is_header":false,"no_amount":false,"original_text":"For the veggies & chickpeas:","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99612,"food":{"id":20168,"name":"toasted sesame oil","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"toasted sesame oil","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3908,"name":"tablespoon","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":1,"conversions":[{"food":"toasted sesame oil","unit":"tablespoon","amount":1}],"note":"preferably toasted sesame oil","order":0,"is_header":false,"no_amount":false,"original_text":"1 tablespoon toasted sesame oil (preferably toasted sesame oil)","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99613,"food":{"id":22245,"name":"Onion","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"Onion","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":7154,"name":"white","plural_name":null,"description":null,"base_unit":null,"open_data_slug":null},"amount":0.5,"conversions":[{"food":"Onion","unit":"white","amount":0.5}],"note":"cut into large chunks","order":0,"is_header":false,"no_amount":false,"original_text":"½ white onion, cut into large chunks","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99614,"food":{"id":21826,"name":"carrots","plural_name":null,"description":"","recipe":null,"url":"","properties":[{"id":27569,"property_amount":0.35,"property_type":{"id":4249,"name":"Fat","unit":"g","description":"","order":0,"open_data_slug":null,"fdc_id":1004}},{"id":27568,"property_amount":10.27,"property_type":{"id":4250,"name":"Carbohydrates","unit":"g","description":"","order":0,"open_data_slug":null,"fdc_id":1005}},{"id":27570,"property_amount":0.94,"property_type":{"id":4251,"name":"Proteins","unit":"g","description":"","order":0,"open_data_slug":null,"fdc_id":1003}},{"id":27563,"property_amount":0,"property_type":{"id":9719,"name":"Sugar","unit":"g","description":"Zucker","order":0,"open_data_slug":null,"fdc_id":1011}}],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":2258586,"food_onhand":true,"supermarket_category":{"id":179,"name":"General","description":""},"parent":null,"numchild":0,"inherit_fields":[],"full_name":"carrots","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":4020,"name":"large","plural_name":null,"description":null,"base_unit":null,"open_data_slug":null},"amount":2,"conversions":[{"food":"carrots","unit":"large","amount":2}],"note":"thinly sliced","order":0,"is_header":false,"no_amount":false,"original_text":"2 large carrots, thinly sliced","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99615,"food":{"id":21153,"name":"bell pepper","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"bell pepper","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":5337,"name":"red","plural_name":null,"description":null,"base_unit":null,"open_data_slug":null},"amount":1,"conversions":[{"food":"bell pepper","unit":"red","amount":1}],"note":"chopped","order":0,"is_header":false,"no_amount":false,"original_text":"1 red bell pepper, chopped","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99616,"food":{"id":41169,"name":"head of broccoli","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"head of broccoli","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":4020,"name":"large","plural_name":null,"description":null,"base_unit":null,"open_data_slug":null},"amount":1,"conversions":[{"food":"head of broccoli","unit":"large","amount":1}],"note":"chopped into florets","order":0,"is_header":false,"no_amount":false,"original_text":"1 large head of broccoli, chopped into florets","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99617,"food":{"id":34295,"name":"chickpeas, rinsed and drained","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"chickpeas, rinsed and drained","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":4076,"name":"can","plural_name":null,"description":null,"base_unit":null,"open_data_slug":null},"amount":1,"conversions":[{"food":"chickpeas, rinsed and drained","unit":"can","amount":1}],"note":"15 ounce","order":0,"is_header":false,"no_amount":false,"original_text":"1 (15 ounce) can chickpeas, rinsed and drained","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99618,"food":{"id":41170,"name":"For the noodles:","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"For the noodles:","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":null,"amount":0,"conversions":[],"note":"","order":0,"is_header":false,"no_amount":false,"original_text":"For the noodles:","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99619,"food":{"id":41171,"name":"stir fry rice noodles","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"stir fry rice noodles","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3999,"name":"ounces","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":10,"conversions":[{"food":"stir fry rice noodles","unit":"ounces","amount":10}],"note":"or feel free to sub ramen noodles or soba noodles","order":0,"is_header":false,"no_amount":false,"original_text":"10 ounces stir fry rice noodles (or feel free to sub ramen noodles or soba noodles)","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99620,"food":{"id":20840,"name":"For serving:","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"For serving:","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":null,"amount":0,"conversions":[],"note":"","order":0,"is_header":false,"no_amount":false,"original_text":"For serving:","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99621,"food":{"id":21159,"name":"Basil leaves","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"Basil leaves","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3910,"name":"cup","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":0.5,"conversions":[{"food":"Basil leaves","unit":"cup","amount":0.5}],"note":"ribboned/julienned","order":0,"is_header":false,"no_amount":false,"original_text":"½ cup basil leaves, ribboned/julienned","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99622,"food":{"id":36081,"name":"roasted cashews","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"roasted cashews","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":{"id":3910,"name":"cup","plural_name":null,"description":"","base_unit":null,"open_data_slug":null},"amount":0.5,"conversions":[{"food":"roasted cashews","unit":"cup","amount":0.5}],"note":"chopped","order":0,"is_header":false,"no_amount":false,"original_text":"1/2 cup roasted cashews, chopped","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99623,"food":{"id":18367,"name":"scallions","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":true,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"scallions","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":null,"amount":0,"conversions":[],"note":"green part of the onion only","order":0,"is_header":false,"no_amount":false,"original_text":"Scallions (green part of the onion only)","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false},{"id":99624,"food":{"id":41172,"name":"Extra sesame seeds","plural_name":null,"description":"","recipe":null,"url":"","properties":[],"properties_food_amount":100,"properties_food_unit":null,"fdc_id":null,"food_onhand":false,"supermarket_category":null,"parent":null,"numchild":0,"inherit_fields":[],"full_name":"Extra sesame seeds","ignore_shopping":false,"substitute":[],"substitute_siblings":false,"substitute_children":false,"substitute_onhand":false,"child_inherit_fields":[],"open_data_slug":null},"unit":null,"amount":0,"conversions":[],"note":"","order":0,"is_header":false,"no_amount":false,"original_text":"Extra sesame seeds","used_in_recipes":[{"id":11321,"name":"30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil"}],"always_use_plural_unit":false,"always_use_plural_food":false}],"instructions_markdown":"<p><em>Wonderful vegan stir fry noodles made in just 30 minutes for the perfect weeknight dinner! These easy sesame noodles have a homemade stir fry sauce and a boost of plant-based protein from chickpeas. Top with fresh herbs and roasted cashews for a beautiful meal you'll make again and again.</em>  </p>\\n<p>First make your stir fry sauce: in a medium bowl, whisk together the soy sauce, water, garlic, coconut sugar, sesame oil, rice vinegar, fresh ginger, sesame seeds, red pepper flakes and arrowroot starch (or cornstarch). Set aside.<br>\\nAdd 1 tablespoon sesame oil to a large pot then add in chopped onion and sliced carrots and cook for 2-4 minutes until onions begin to soften. Next add in broccoli and bell pepper and cook, stirring frequently, for an additional 6-8 minutes or until broccoli is slightly tender but still have a bite.<br>\\nWhile the veggies are cooking, make your stir fry rice noodles according to the directions on the package. Then drain and set aside.<br>\\nAdd the drained chickpeas to the pot with the cooked veggies. Immediately turn the heat to low and add in the sauce. Cook for an additional 2 minutes over low heat until the sauce begins to thicken a bit. It should be nice and saucy. Stir in rice noodles, fresh basil ribbons and cashews; toss again to combine. Garnish with scallions and cilantro. Serves 4.</p>","time":0,"order":0,"show_as_header":true,"file":null,"step_recipe":null,"step_recipe_data":null,"show_ingredients_table":true}],"working_time":15,"waiting_time":15,"created_by":2,"created_at":"2022-08-19T04:29:55.207363+02:00","updated_at":"2024-03-13T18:41:58.272557+01:00","source_url":"https://www.ambitiouskitchen.com/vegan-stir-fry-sesame-noodles/","internal":true,"show_ingredient_overview":true,"nutrition":null,"properties":[],"servings":4,"file_path":"","servings_text":"['4', '4 servings']","rating":null,"last_cooked":null,"private":false,"shared":[]}`
				w.WriteHeader(200)
				_, _ = w.Write([]byte(data))
			} else if r.URL.Path == "/image" {
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
		got, err = integrations.TandoorImport(srv.URL, "master", "yoda", srv.Client(), files.UploadImage, c)
		if err != nil {
			return
		}
	}()
	for range c {
	}

	r := models.Recipe{
		Category:  "vegan",
		CreatedAt: time.Date(2022, 8, 19, 2, 29, 55, 207363000, time.UTC),
		Ingredients: []string{
			"For the Sauce:",
			"0.33 cup low sodium soy sauce or coconut aminos",
			"0.33 cup water",
			"3 cloves garlic",
			"2 tablespoons coconut sugar or brown sugar",
			"1 tablespoon sesame oil",
			"1 tablespoon rice vinegar",
			"1 tablespoon fresh grated ginger",
			"1 tablespoon sesame seeds",
			"0.5 teaspoon red pepper flakes",
			"0.5 tablespoon arrowroot starch",
			"For the veggies & chickpeas:",
			"1 tablespoon toasted sesame oil",
			"0.5 white Onion",
			"2 large carrots",
			"1 red bell pepper",
			"1 large head of broccoli",
			"1 can chickpeas, rinsed and drained",
			"For the noodles:",
			"10 ounces stir fry rice noodles",
			"For serving:",
			"0.5 cup Basil leaves",
			"0.5 cup roasted cashews",
			"scallions",
			"Extra sesame seeds",
		},
		Instructions: []string{
			"*Wonderful vegan stir fry noodles made in just 30 minutes for the perfect weeknight dinner! These easy sesame noodles have a homemade stir fry sauce and a boost of plant-based protein from chickpeas. Top with fresh herbs and roasted cashews for a beautiful meal you'll make again and again.*",
			"First make your stir fry sauce: in a medium bowl, whisk together the soy sauce, water, garlic, coconut sugar, sesame oil, rice vinegar, fresh ginger, sesame seeds, red pepper flakes and arrowroot starch (or cornstarch). Set aside.  \\nAdd 1 tablespoon sesame oil to a large pot then add in chopped onion and sliced carrots and cook for 2-4 minutes until onions begin to soften. Next add in broccoli and bell pepper and cook, stirring frequently, for an additional 6-8 minutes or until broccoli is slightly tender but still have a bite.  \\nWhile the veggies are cooking, make your stir fry rice noodles according to the directions on the package. Then drain and set aside.  \\nAdd the drained chickpeas to the pot with the cooked veggies. Immediately turn the heat to low and add in the sauce. Cook for an additional 2 minutes over low heat until the sauce begins to thicken a bit. It should be nice and saucy. Stir in rice noodles, fresh basil ribbons and cashews; toss again to combine. Garnish with scallions and cilantro. Serves 4.",
		},
		Keywords:  []string{"vegan", "lunch", " dinner", " vegetarian", " asian", " gluten free"},
		Name:      "30 Minute Vegan Stir Fry Sesame Noodles with Chickpeas & Basil",
		Times:     models.Times{Prep: 15 * time.Minute, Cook: 15 * time.Minute},
		UpdatedAt: time.Date(2024, 03, 13, 17, 41, 58, 272557000, time.UTC),
		URL:       "https://www.ambitiouskitchen.com/vegan-stir-fry-sesame-noodles/",
		Yield:     4,
	}
	want := models.Recipes{r, r, r}
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}
