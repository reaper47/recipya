package statements

import (
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"strconv"
	"strings"
)

// BuildSelectPaginatedResults builds a SQL query for paginated search results.
func BuildSelectPaginatedResults(opts models.SearchOptionsRecipes) string {
	var sb strings.Builder
	sb.WriteString(buildSelectPaginatedResultsQuery(opts))
	sb.WriteString(" SELECT * FROM results WHERE row_num BETWEEN ")
	sb.WriteString(strconv.FormatUint((opts.Page-1)*templates.ResultsPerPage+1, 10))
	sb.WriteString(" AND ")
	sb.WriteString(strconv.FormatUint((opts.Page-1)*templates.ResultsPerPage+15, 10))
	return sb.String()
}

// BuildSelectSearchResultsCount builds a SQL query for fetching the number of paginated results.
func BuildSelectSearchResultsCount(options models.SearchOptionsRecipes) string {
	var sb strings.Builder
	options.Sort = models.Sort{}
	sb.WriteString(buildSelectPaginatedResultsQuery(options))
	sb.WriteString("SELECT COUNT(*) FROM results")
	return sb.String()
}

func buildSelectPaginatedResultsQuery(options models.SearchOptionsRecipes) string {
	var sb strings.Builder
	sb.WriteString("WITH results AS (")
	sb.WriteString(buildSearchRecipeQuery(options))
	sb.WriteString(")")
	return sb.String()
}

func buildSearchRecipeQuery(opts models.SearchOptionsRecipes) string {
	var sb strings.Builder

	sb.WriteString("SELECT recipe_id, name, description, image, created_at, category, keywords, row_num FROM (" + BuildBaseSelectRecipe(opts.Sort))
	sb.WriteString(" WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ?")

	if opts.Query != "" || !opts.IsBasic() {
		sb.WriteString(" AND recipes_fts MATCH ?")
	}

	sb.WriteString(" ORDER BY rank)")
	if opts.CookbookID > 0 {
		sb.WriteString(" AND recipes.id NOT IN (SELECT recipe_id FROM cookbook_recipes WHERE cookbook_id = ?)")
	}
	sb.WriteString(" GROUP BY recipes.id)")
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
					 AND recipe_id IN (SELECT id
									  FROM recipes
									  WHERE name = ?
										AND description = ?
										AND yield = ?
									 	AND url = ?)
			   )`

// SelectAppInfo fetches general information on the application.
const SelectAppInfo = `
	SELECT is_update_available, updated_at, update_last_checked_at
	FROM app
	WHERE id = 1`

// SelectAuthToken fetches a non-expired auth token by the selector.
const SelectAuthToken = `
	SELECT id, hash_validator, expires, user_id
	FROM auth_tokens
	WHERE selector = ?
	AND expires > unixepoch('now')`

// SelectCategories fetches a user's recipe categories.
const SelectCategories = `
	SELECT c.name
	FROM user_category AS uc
	JOIN categories c ON c.id = uc.category_id
	WHERE uc.user_id = ?
	ORDER BY name`

// SelectCookbook gets a user's cookbook by cookbook ID.
const SelectCookbook = `
	SELECT c.id, c.title, c.image, c.count
	FROM cookbooks AS c
	WHERE id = ?
		AND user_id = ?`

// SelectCookbookExists verifies whether the cookbook belongs to the user.
const SelectCookbookExists = `
	SELECT EXISTS (SELECT c.id
				   FROM cookbooks AS c
				   WHERE id = ?
					 AND user_id = ?)`

// SelectCookbookRecipeExists verifies whether the recipe and the cookbook belongs to a user.
const SelectCookbookRecipeExists = `
	SELECT EXISTS (SELECT c.id
				   FROM cookbooks AS c
							JOIN user_recipe AS ur ON c.user_id = ur.user_id
				   WHERE c.id = ?
					 AND c.user_id = ?
					 AND ur.recipe_id = ?);`

// SelectCookbookRecipe fetches a recipe from a cookbook.
const SelectCookbookRecipe = baseSelectRecipe + `
	JOIN cookbook_recipes AS cr ON recipes.id = cr.recipe_id
	WHERE cr.cookbook_id = ?
		AND cr.recipe_id = ?
	GROUP BY recipes.id`

// SelectCookbookRecipes fetches the recipes in a cookbook.
const SelectCookbookRecipes = baseSelectRecipe + `
	JOIN cookbook_recipes AS cr ON recipes.id = cr.recipe_id
	WHERE cr.cookbook_id = ?
	GROUP BY recipes.id
	ORDER BY cr.order_index`

// SelectCookbookShared gets a shared cookbook link.
const SelectCookbookShared = `
	SELECT cookbook_id, user_id
	FROM share_cookbooks
	WHERE link = ?`

// SelectCookbookSharedLink gets the link of a shared cookbook.
const SelectCookbookSharedLink = `
	SELECT link 
	FROM share_cookbooks
	WHERE cookbook_id = ?
		AND user_id = ?`

// SelectCookbooksShared gets the user's shared cookbooks.
const SelectCookbooksShared = `
	SELECT link, cookbook_id
	FROM share_cookbooks
	WHERE user_id = ?`

// SelectCookbookUser gets the ID of the user who has the cookbook ID.
const SelectCookbookUser = `
	SELECT user_id
	FROM cookbooks
	WHERE id = ?`

// SelectCookbooks gets a limited number of cookbooks belonging to the user.
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

// SelectCookbooksUser gets all cookbooks belonging to the user.
const SelectCookbooksUser = `
	SELECT id, title, image, count
	FROM cookbooks
	WHERE user_id = 2`

// SelectCounts gets the number of recipes and cookbooks belonging to the user.
const SelectCounts = `
	SELECT cookbooks, recipes
	FROM counts 
	WHERE user_id = ?`

// SelectCountWebsites fetches the number of supported websites.
const SelectCountWebsites = `
	SELECT COUNT(id)
	FROM websites`

// SelectCuisineID gets the ID of the specified cuisine.
const SelectCuisineID = `
	SELECT id 
	FROM cuisines 
	WHERE name = ?`

// SelectDistinctImages gets all distinct image UUIDs from the recipes table.
const SelectDistinctImages = `
	SELECT DISTINCT image
	FROM recipes
	UNION
	SELECT DISTINCT image
	FROM cookbooks`

// SelectKeywords fetches all keywords.
const SelectKeywords = `
	SELECT name 
	FROM keywords
	ORDER BY name`

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

// BuildBaseSelectRecipe builds from the options.
func BuildBaseSelectRecipe(sorts models.Sort) string {
	var s string
	if sorts.IsAToZ {
		s = "recipes.name ASC"
	} else if sorts.IsZToA {
		s = "recipes.name DESC"
	} else if sorts.IsNewestToOldest {
		s = "recipes.created_at DESC"
	} else if sorts.IsOldestToNewest {
		s = "recipes.created_at ASC"
	} else if sorts.IsRandom {
		s = "RANDOM()"
	} else {
		return baseSelectSearchRecipe
	}

	before, after, _ := strings.Cut(baseSelectSearchRecipe, "ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num")
	var sb strings.Builder
	sb.WriteString(before)
	sb.WriteString(" ROW_NUMBER() OVER (ORDER BY " + s + ") AS row_num")
	sb.WriteString(after)
	return sb.String()
}

const baseSelectRecipe = `
	SELECT recipes.id                                                                   AS recipe_id,
		   recipes.name                                                                 AS name,
		   recipes.description                                                          AS description,
		   recipes.image                                                                AS image,
		   COALESCE((SELECT GROUP_CONCAT(other_image, ';')
					 FROM (SELECT image AS other_image
						   FROM additional_images_recipe
						   WHERE additional_images_recipe.recipe_id = recipes.id)), '') AS other_images,
		   recipes.url                                                                  AS url,
		   recipes.yield                                                                AS yield,
		   recipes.created_at                                                           AS created_at,
		   recipes.updated_at                                                           AS updated_at,
		   categories.name                                                              AS category,
		   cuisines.name                                                                AS cuisine,
		   COALESCE((SELECT GROUP_CONCAT(ingredient_name, '<!---->')
					 FROM (SELECT DISTINCT ingredients.name AS ingredient_name
						   FROM ingredient_recipe
									JOIN ingredients ON ingredients.id = ingredient_recipe.ingredient_id
						   WHERE ingredient_recipe.recipe_id = recipes.id
						   ORDER BY ingredient_order)), '')                             AS ingredients,
		   COALESCE((SELECT GROUP_CONCAT(instruction_name, '<!---->')
					 FROM (SELECT DISTINCT instructions.name AS instruction_name
						   FROM instruction_recipe
									JOIN instructions ON instructions.id = instruction_recipe.instruction_id
						   WHERE instruction_recipe.recipe_id = recipes.id
						   ORDER BY instruction_order)), '')                            AS instructions,
		   GROUP_CONCAT(DISTINCT keywords.name)                                         AS keywords,
		   (SELECT GROUP_CONCAT(name)
			FROM (SELECT tool_recipe.quantity || ' ' || tools.name AS name
				  FROM tool_recipe
						   JOIN tools ON tool_recipe.tool_id = tools.id
				  WHERE tool_recipe.recipe_id = recipes.id
				  ORDER BY tool_recipe.tool_order))                                     AS tools,
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
		   nutrition.is_per_serving,
		   times.prep_seconds,
		   times.cook_seconds,
		   times.total_seconds,
		   ROW_NUMBER() OVER (ORDER BY recipes.id)                                      AS row_num
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

const baseSelectSearchRecipe = `
	SELECT recipes.id                                                                      AS recipe_id,
		   recipes.name                                                                    AS name,
		   recipes.description                                                             AS description,
		   recipes.image                                                                   AS image,
		   recipes.created_at                                                              AS created_at,
		   categories.name                                                                 AS category,
		   GROUP_CONCAT(DISTINCT keywords.name)  AS keywords,
		   user_id,
		   ROW_NUMBER() OVER (ORDER BY recipes.id) AS row_num
	FROM recipes 
			 LEFT JOIN category_recipe ON recipes.id = category_recipe.recipe_id
			 LEFT JOIN categories ON category_recipe.category_id = categories.id
	    	 LEFT JOIN keyword_recipe ON recipes.id = keyword_recipe.recipe_id
			 LEFT JOIN keywords ON keyword_recipe.keyword_id = keywords.id
			 LEFT JOIN user_recipe ON recipes.id = user_recipe.recipe_id`

// SelectRecipe fetches a user's recipe.
const SelectRecipe = baseSelectRecipe + `
	INNER JOIN user_recipe AS ur ON ur.recipe_id = recipes.id
	WHERE recipes.id = ?
		AND ur.user_id = ?
	LIMIT 1`

// SelectRecipesAll fetches all the user's recipes.
const SelectRecipesAll = baseSelectRecipe + `
	WHERE recipes.id IN (SELECT recipe_id FROM user_recipe WHERE user_id = ?)
	GROUP BY recipes.id`

// SelectRecipes fetches a chunk of the user's recipes.
const SelectRecipes = `
	WITh results AS (
		SELECT recipe_id, name, description, image, created_at, category, keywords, row_num FROM (
			` + baseSelectSearchRecipe + `
			WHERE user_recipe.user_id = ?
			GROUP BY recipes.id
		)
	) SELECT * FROM results WHERE row_num BETWEEN (?-1)*` + templates.ResultsPerPageStr + `+1 AND (?-1)*` + templates.ResultsPerPageStr + `+` + templates.ResultsPerPageStr

// SelectRecipeShared checks whether the recipe is shared.
const SelectRecipeShared = `
	SELECT recipe_id, user_id
	FROM share_recipes
	WHERE link = ?`

// SelectRecipeSharedFromRecipeID gets the user the shared recipe belongs to.
const SelectRecipeSharedFromRecipeID = `
	SELECT user_id
	FROM share_recipes
	WHERE recipe_id = ?`

// SelectRecipesShared gets the recipes the user shared.
const SelectRecipesShared = `
	SELECT link, recipe_id
	FROM share_recipes
	WHERE user_id = ?`

// SelectRecipeUser fetches the user whose recipe belongs to.
const SelectRecipeUser = `
	SELECT user_id
	FROM user_recipe
	WHERE recipe_id = ?`

// SelectRecipeUserExist checks whether the shared recipe belongs to the current user.
const SelectRecipeUserExist = `
	SELECT EXISTS(
		SELECT 1
		FROM user_recipe
		WHERE recipe_id = ?
			AND user_id = ?
	)`

// SelectReport fetches the report of the given ID belonging to the user.
const SelectReport = `
	SELECT id, title, success, error_reason
	FROM report_logs
	WHERE report_id = (SELECT id FROM reports WHERE id = ? AND user_id = ?)`

// SelectReports fetches the import reports for the user.
const SelectReports = `
	SELECT 
    r.id,  
    r.created_at,  
    r.exec_time_ns,  
    GROUP_CONCAT(l.id, ';') AS log_ids
FROM  
    reports r
LEFT JOIN  
    report_logs l ON r.id = l.report_id
WHERE r.report_type = ? AND r.user_id = ?
GROUP BY r.id`

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

// SelectUserSettings fetchs a user's settings.
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

// SelectUserPasswordByID fetches the user's hashed password by their id.
const SelectUserPasswordByID = `
	SELECT hashed_password
	FROM users
	WHERE id = ?`

// SelectUserOne fetches the first user.
const SelectUserOne = `
	SELECT id
	FROM users
	WHERE id = 1`

// SelectUsers fetches all users from the database.
const SelectUsers = `
	SELECT id, email 
	FROM users
	ORDER BY id`

// SelectWebsites fetches all websites from the database.
const SelectWebsites = `
	SELECT id, host, url
	FROM websites`
