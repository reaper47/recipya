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

// SelectCategoryID is the query to get the ID of the specified category.
const SelectCategoryID = `
	SELECT id 
	FROM categories 
	WHERE name = ?`

// SelectCuisineID is the query to get the ID of the specified cuisine.
const SelectCuisineID = `
	SELECT id 
	FROM cuisines 
	WHERE name = ?`

// SelectCountWebsites fetches the number of supported websites.
const SelectCountWebsites = `
	SELECT COUNT(id)
	FROM websites`

// SelectRecipe fetches a user's recipe.
const SelectRecipe = `
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
						   WHERE ingredient_recipe.recipe_id = (SELECT recipe_id
																FROM (SELECT recipe_id,
																			 ROW_NUMBER() OVER (ORDER BY id) AS row_num
																	  FROM user_recipe
																	  WHERE user_id = ?) AS t
																WHERE row_num = ?))), '')  AS ingredients,
		   COALESCE((SELECT GROUP_CONCAT(instruction_name, '<!---->')
					 FROM (SELECT DISTINCT instructions.name AS instruction_name
						   FROM instruction_recipe
									JOIN instructions ON instructions.id = instruction_recipe.instruction_id
						   WHERE instruction_recipe.recipe_id = (SELECT recipe_id
																 FROM (SELECT recipe_id,
																			  ROW_NUMBER() OVER (ORDER BY id) AS row_num
																	   FROM user_recipe
																	   WHERE user_id = ?) AS t
																 WHERE row_num = ?))), '') AS instructions,
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
			 LEFT JOIN times ON time_recipe.time_id = times.id
	WHERE recipes.id = (SELECT recipe_id
						FROM (SELECT recipe_id,
									 ROW_NUMBER() OVER (ORDER BY id) AS row_num
							  FROM user_recipe
							  WHERE user_id = ?) AS t
						WHERE row_num = ?)`

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

// SelectUserPassword fetches the user's password for verification purposes.
const SelectUserPassword = `
	SELECT id, hashed_password
	FROM users
	WHERE email = ?`

// SelectUsers fetches all users from the database.
const SelectUsers = `
	SELECT id, email 
	FROM users`

// SelectWebsites fetches all websites from the database.
const SelectWebsites = `
	SELECT id, host, url
	FROM websites`

// SelectWebsitesSearch fetches all websites that match the query.
const SelectWebsitesSearch = `
	SELECT id, host, url
	FROM websites
	WHERE url
	LIKE ?`
