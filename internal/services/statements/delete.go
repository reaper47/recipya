package statements

// DeleteAuthToken removes the authentication token associated with the user id from the database.
const DeleteAuthToken = `
	DELETE
	FROM auth_tokens
	WHERE user_id = ?`

// DeleteCookbook is the query to delete a user's cookbook.
const DeleteCookbook = `
	DELETE 
	FROM cookbooks
	WHERE id = (SELECT id
				FROM (SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS row_num
					  FROM cookbooks
					  WHERE user_id = ?) AS t
				WHERE row_num = ?)`

// DeleteRecipe is the query to delete a user's recipe and the recipe itself.
const DeleteRecipe = `
	DELETE
	FROM recipes
	WHERE recipes.id = (SELECT recipe_id
						FROM (SELECT recipe_id,
									 ROW_NUMBER() OVER (ORDER BY id) AS row_num
							  FROM user_recipe
							  WHERE user_id = ?) AS t
						WHERE row_num = ?)`

// DeleteRecipeIngredients is the query for deleting all ingredients from a recipe.
const DeleteRecipeIngredients = `
	DELETE
	FROM ingredient_recipe
	WHERE recipe_id = ?`

// DeleteRecipeInstructions is the query for deleting all instructions from a recipe.
const DeleteRecipeInstructions = `
	DELETE
	FROM instruction_recipe
	WHERE recipe_id = ?`

// DeleteRecipeKeywords is the query for deleting all keywords from a recipe.
const DeleteRecipeKeywords = `
	DELETE
	FROM keyword_recipe
	WHERE recipe_id = ?`

// DeleteRecipeTools is the query for deleting all tools from a recipe.
const DeleteRecipeTools = `
	DELETE
	FROM tool_recipe
	WHERE recipe_id = ?`
