package statements

// DeleteAuthToken removes the authentication token associated with the user id from the database.
const DeleteAuthToken = `
	DELETE
	FROM auth_tokens
	WHERE user_id = ?`

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
