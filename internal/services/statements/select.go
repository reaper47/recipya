package statements

import (
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"strings"
)

// RecipesFTSFields lists all columns in the recipes_fts table.
var RecipesFTSFields = []string{"name", "description", "category", "ingredients", "instructions", "keywords", "source"}

// BuildSearchRecipeQuery builds a SQL query for searching recipes.
func BuildSearchRecipeQuery(queries []string, options models.SearchOptionsRecipes) string {
	var sb strings.Builder
	sb.WriteString(baseSelectRecipe)
	sb.WriteString(" WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ?")

	n := len(queries)
	if n > 0 {
		sb.WriteString(" AND (")
		if options.ByName {
			switch n {
			case 1:
				sb.WriteString("name MATCH ?)")
			default:
				sb.WriteString("name MATCH ?")
				for range queries[1:] {
					sb.WriteString(" AND name MATCH ?")
				}
				sb.WriteString(")")
			}
		} else if options.FullSearch {
			nFields := len(RecipesFTSFields)
			for i, field := range RecipesFTSFields {
				if n == 1 {
					sb.WriteString(field + " MATCH ?")
					if i < nFields-1 {
						sb.WriteString(" OR ")
					}
				} else {
					sb.WriteString("(")
					for j := range queries {
						sb.WriteString(field + " MATCH ?")
						if j < n-1 {
							sb.WriteString(" AND ")
						}
					}

					if i < nFields-1 {
						sb.WriteString(") OR ")
					} else {
						sb.WriteString(")")
					}
				}
			}

			sb.WriteString(")")
		}
	}

	sb.WriteString(" ORDER BY rank) GROUP BY recipes.id LIMIT 30")
	return sb.String()
}

// BuildSelectNutrientFDC builds the query to fetch a nutrient from the FDC database.
func BuildSelectNutrientFDC(ingredients []string) string {
	var sb strings.Builder
	for i, ing := range ingredients {
		if i > 0 {
			sb.WriteString(" AND ")
		}
		sb.WriteString("description LIKE '%")
		sb.WriteString(ing)
		sb.WriteString("%'")
	}

	return `
	SELECT food.fdc_id,
		   nutrient.name,
		   food_nutrient.amount,
		   nutrient.unit_name
	FROM food_nutrient
			 INNER JOIN food ON food_nutrient.fdc_id = food.fdc_id
			 INNER JOIN nutrient ON food_nutrient.nutrient_id = nutrient.id
	WHERE food.fdc_id = (SELECT fdc_id
						 FROM food
						 WHERE ` + sb.String() + `
                             AND ((data_type = 'sr_legacy_food' AND (food_category_id = 11 OR food_category_id = 2 OR food_category_id = 18)) OR data_type = 'survey_fndds_food' OR data_type = 'branded_food')
						 ORDER BY 
							 food_category_id DESC, 
							 data_type DESC,
							 description ASC
						 LIMIT 1)
	  AND nutrient.name IN (
							'Energy',
							'Cholesterol',
							'Carbohydrate, by difference',
							'Fiber, total dietary',
							'Protein',
							'Fatty acids, total monounsaturated',
							'Fatty acids, total polyunsaturated',
							'Fatty acids, total trans',
							'Fatty acids, total saturated',
							'Sodium, Na',
							'Sugars, total including NLEA')`
}

// IsRecipeForUserExist checks whether the recipe belongs to the given user.
const IsRecipeForUserExist = `
	SELECT EXISTS(
				   SELECT id
				   FROM user_recipe
				   WHERE user_id = ?
					 AND recipe_id = (SELECT id
									  FROM recipes
									  WHERE name = ?
										AND description = ?
										AND yield = ?
									 	AND url = ?)
			   )`

// SelectAuthToken fetches a non-expired auth token by the selector.
const SelectAuthToken = `
	SELECT id, hash_validator, expires, user_id
	FROM auth_tokens
	WHERE selector = ?
	AND expires > unixepoch('now')`

// SelectCategories is the query to fetch a user's recipe categories.
const SelectCategories = `
	SELECT c.name
	FROM user_category AS uc
	JOIN categories c ON c.id = uc.category_id
	WHERE uc.user_id = ?
	ORDER BY name`

// SelectCookbook is the query to get a user's cookbook.
const SelectCookbook = `
	SELECT c.id, c.title, c.count
	FROM cookbooks AS c
	WHERE id = (SELECT id
				 FROM (SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS row_num
					   FROM cookbooks
					   WHERE user_id = ?)
				 WHERE row_num > (? - 1) *` + templates.ResultsPerPageStr + " + ?)"

// SelectCookbookByID is the query to get a user's cookbook by cookbook ID.
const SelectCookbookByID = `
	SELECT c.id, c.title, c.image, c.count
	FROM cookbooks AS c
	WHERE id = ?
		AND user_id = ?`

// SelectCookbookExists is the query to verify whether the cookbook belongs to the user.
const SelectCookbookExists = `
	SELECT EXISTS (SELECT c.id
				   FROM cookbooks AS c
				   WHERE id = ?
					 AND user_id = ?)`

// SelectCookbookRecipeExists is the query to verify whether the recipe and the cookbook belongs to a user.
const SelectCookbookRecipeExists = `
	SELECT EXISTS (SELECT c.id
				   FROM cookbooks AS c
							JOIN user_recipe AS ur ON c.user_id = ur.user_id
				   WHERE c.id = ?
					 AND c.user_id = ?
					 AND ur.recipe_id = ?);`

// SelectCookbookRecipe is the query to fetch a recipe from a cookbook.
const SelectCookbookRecipe = baseSelectRecipe + `
	JOIN cookbook_recipes AS cr ON recipes.id = cr.recipe_id
	WHERE cr.cookbook_id = ?
		AND cr.recipe_id = ?
	GROUP BY recipes.id`

// SelectCookbookRecipes is the query to fetch the recipes in a cookbook.
const SelectCookbookRecipes = baseSelectRecipe + `
	JOIN cookbook_recipes AS cr ON recipes.id = cr.recipe_id
	WHERE cr.cookbook_id = ?
	GROUP BY recipes.id
	ORDER BY cr.order_index`

// SelectCookbookShared is the query to get a shared cookbook link.
const SelectCookbookShared = `
	SELECT cookbook_id, user_id
	FROM share_cookbooks
	WHERE link = ?`

// SelectCookbookSharedLink is the query to get the link of a shared cookbook.
const SelectCookbookSharedLink = `
	SELECT link 
	FROM share_cookbooks
	WHERE cookbook_id = ?
		AND user_id = ?`

// SelectCookbookUser is the query to get the ID of the user who has the cookbook ID.
const SelectCookbookUser = `
	SELECT user_id
	FROM cookbooks
	WHERE id = ?`

// SelectCookbooks is the query to get a limited number of cookbooks belonging to the user.
var SelectCookbooks = `
	SELECT id, image, title, count
	FROM cookbooks
	WHERE id >= (SELECT id
				 FROM (SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS row_num
					   FROM cookbooks
					   WHERE user_id = ?)
				 WHERE row_num > (? - 1) *` + templates.ResultsPerPageStr + " + " + templates.ResultsPerPageStr + `)
		AND user_id = ?
	LIMIT ` + templates.ResultsPerPageStr

// SelectCounts is the query to get the number of recipes and cookbooks belonging to the user.
const SelectCounts = `
	SELECT cookbooks, recipes
	FROM counts 
	WHERE user_id = ?`

// SelectCountWebsites fetches the number of supported websites.
const SelectCountWebsites = `
	SELECT COUNT(id)
	FROM websites`

// SelectCuisineID is the query to get the ID of the specified cuisine.
const SelectCuisineID = `
	SELECT id 
	FROM cuisines 
	WHERE name = ?`

// SelectMeasurementSystems fetches the units systems along with the user's selected system and settings.
const SelectMeasurementSystems = `
	SELECT ms.name,
		   COALESCE((SELECT GROUP_CONCAT(name)
					 FROM measurement_systems), '') AS systems,
		   us.convert_automatically,
		   us.calculate_nutrition
	FROM measurement_systems AS ms
			 JOIN user_settings AS us ON measurement_system_id = ms.id
	WHERE user_id = ?`

const baseSelectRecipe = `
	SELECT recipes.id                                                                      AS recipe_id,
		   recipes.name                                                                    AS name,
		   recipes.description                                                             AS description,
		   recipes.image                                                                   AS image,
		   recipes.url                                                                     AS url,
		   recipes.yield                                                                   AS yield,
		   recipes.created_at                                                              AS created_at,
		   recipes.updated_at                                                              AS updated_at,
		   categories.name                                                                 AS category,
		   cuisines.name                                                                   AS cuisine,
		   COALESCE((SELECT GROUP_CONCAT(ingredient_name, '<!---->')
					 FROM (SELECT DISTINCT ingredients.name AS ingredient_name
						   FROM ingredient_recipe
									JOIN ingredients ON ingredients.id = ingredient_recipe.ingredient_id
						   WHERE ingredient_recipe.recipe_id = recipes.id
						   ORDER BY ingredient_order)), '')  AS ingredients,
		   COALESCE((SELECT GROUP_CONCAT(instruction_name, '<!---->')
					 FROM (SELECT DISTINCT instructions.name AS instruction_name
						   FROM instruction_recipe
									JOIN instructions ON instructions.id = instruction_recipe.instruction_id
						   WHERE instruction_recipe.recipe_id = recipes.id
						   ORDER BY instruction_order)), '') AS instructions,
		   COALESCE(GROUP_CONCAT(DISTINCT keywords.name), '')                              AS keywords,
		   COALESCE(GROUP_CONCAT(DISTINCT tools.name), '')                                 AS tools,
		   nutrition.calories,
		   nutrition.total_carbohydrates,
		   nutrition.sugars,
		   nutrition.protein,
		   nutrition.total_fat,
		   nutrition.saturated_fat,
		   nutrition.unsaturated_fat,
		   nutrition.cholesterol,
		   nutrition.sodium,
		   nutrition.fiber,
		   times.prep_seconds,
		   times.cook_seconds,
		   times.total_seconds
	FROM recipes
			 LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id
			 LEFT JOIN categories ON category_recipe.category_id = categories.id
			 LEFT JOIN cuisine_recipe ON recipes.id = cuisine_recipe.recipe_id
			 LEFT JOIN cuisines ON cuisine_recipe.cuisine_id = cuisines.id
			 LEFT JOIN ingredient_recipe ON recipes.id = ingredient_recipe.recipe_id
			 LEFT JOIN ingredients ON ingredient_recipe.ingredient_id = ingredients.id
			 LEFT JOIN instruction_recipe ON recipes.id = instruction_recipe.recipe_id
			 LEFT JOIN instructions ON instruction_recipe.instruction_id = instructions.id
			 LEFT JOIN keyword_recipe ON recipes.id = keyword_recipe.recipe_id
			 LEFT JOIN keywords ON keyword_recipe.keyword_id = keywords.id
			 LEFT JOIN tool_recipe ON recipes.id = tool_recipe.recipe_id
			 LEFT JOIN tools ON tool_recipe.tool_id = tools.id
			 LEFT JOIN nutrition ON recipes.id = nutrition.recipe_id
			 LEFT JOIN time_recipe ON recipes.id = time_recipe.recipe_id
			 LEFT JOIN times ON time_recipe.time_id = times.id`

// SelectRecipe fetches a user's recipe.
const SelectRecipe = baseSelectRecipe + `
	WHERE recipes.id = (SELECT recipe_id
						FROM (SELECT recipe_id,
									 ROW_NUMBER() OVER (ORDER BY id) AS row_num
							  FROM user_recipe
							  WHERE user_id = ?) AS t
						WHERE row_num = ?)`

// SelectRecipesAll is the query to fetch all the user's recipes.
const SelectRecipesAll = baseSelectRecipe + `
	WHERE recipes.id IN (SELECT recipe_id FROM user_recipe WHERE user_id = ?)
	GROUP BY recipes.id`

// SelectRecipes is the query to fetch a chunk of the user's recipes.
const SelectRecipes = baseSelectRecipe + `
	WHERE recipes.id >= (SELECT id
			 FROM (SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS row_num
				   FROM user_recipe AS UR
				   WHERE ur.user_id = ?)
			 WHERE row_num > (? - 2) *` + templates.ResultsPerPageStr + " + " + templates.ResultsPerPageStr + `)
	GROUP BY recipes.id
	LIMIT ` + templates.ResultsPerPageStr

// SelectRecipeShared checks whether the recipe is shared.
const SelectRecipeShared = `
	SELECT recipe_id, user_id
	FROM share_recipes
	WHERE link = ?`

// SelectRecipeUser fetches the user whose recipe belongs to.
const SelectRecipeUser = `
	SELECT user_id
	FROM user_recipe
	WHERE recipe_id = ?`

// SelectUserExist checks whether the user is present.
const SelectUserExist = `
	SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE email = ?
	)`

// SelectUserEmail fetches the user's email from their id.
const SelectUserEmail = `
	SELECT email
	FROM users
	WHERE id = ?`

// SelectUserID fetches the user's id from their email.
const SelectUserID = `
	SELECT id
	FROM users
	WHERE email = ?`

// SelectUserSettings is the query to fetch a user's settings.
const SelectUserSettings = `
	SELECT MS.name, convert_automatically, cookbooks_view, calculate_nutrition
	FROM user_settings
	JOIN measurement_systems MS on MS.id = measurement_system_id
	WHERE user_id = ?`

// SelectUserPassword fetches the user's password for verification purposes.
const SelectUserPassword = `
	SELECT id, hashed_password
	FROM users
	WHERE email = ?`

// SelectUserPasswordByID is the query to fetch the user's hashed password by their id.
const SelectUserPasswordByID = `
	SELECT hashed_password
	FROM users
	WHERE id = ?`

// SelectUsers fetches all users from the database.
const SelectUsers = `
	SELECT id, email 
	FROM users`

// SelectWebsites fetches all websites from the database.
const SelectWebsites = `
	SELECT id, host, url
	FROM websites`
