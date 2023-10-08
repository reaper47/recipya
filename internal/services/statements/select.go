package statements

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

// SelectCookbooks is the query to get a limited number of cookbooks belonging to the user
const SelectCookbooks = `
	SELECT id, image, title, count
	FROM cookbooks
	WHERE user_id = ?`

// SelectCuisineID is the query to get the ID of the specified cuisine.
const SelectCuisineID = `
	SELECT id 
	FROM cuisines 
	WHERE name = ?`

// SelectRecipeCount is the query to get the number of recipes belonging to the user.
const SelectRecipeCount = `
	SELECT recipes
	FROM counts 
	WHERE user_id = ?`

// SelectCountWebsites fetches the number of supported websites.
const SelectCountWebsites = `
	SELECT COUNT(id)
	FROM websites`

// SelectMeasurementSystems fetches the units systems along with the user's selected system and settings.
const SelectMeasurementSystems = `
	SELECT ms.name,
		   COALESCE((SELECT GROUP_CONCAT(name)
					 FROM measurement_systems), '') AS systems,
		   us.convert_automatically
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

// SelectRecipes is the query to fetch all the user's recipes.
const SelectRecipes = baseSelectRecipe + `
	WHERE recipes.id IN (SELECT recipe_id FROM user_recipe WHERE user_id = ?)
	GROUP BY recipes.id`

// SelectRecipeShared checks whether the recipe is shared.
const SelectRecipeShared = `
	SELECT recipe_id, user_id
	FROM share
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
	SELECT MS.name, convert_automatically, cookbooks_view
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
