package statements

// DeleteAuthToken removes the authentication token associated with the user id from the database.
const DeleteAuthToken = `
	DELETE
	FROM auth_tokens
	WHERE user_id = ?`

// DeleteCookbook deletes a user's cookbook.
const DeleteCookbook = `
	DELETE
	FROM cookbooks
	WHERE id = ?
		AND user_id = ?`

// DeleteCookbookRecipe deletes a recipe from a user's cookbook.
const DeleteCookbookRecipe = `
	DELETE
	FROM cookbook_recipes
	WHERE cookbook_id = (SELECT id FROM cookbooks WHERE id = ? AND user_id = ?)
		AND recipe_id = ?`

// DeleteCookbooks deletes all the user's cookbooks.
const DeleteCookbooks = `
	DELETE
	FROM cookbooks
	WHERE user_id = ?`

// DeleteRecipe deletes a user's recipe and the recipe itself.
const DeleteRecipe = `
	DELETE
	FROM recipes
	WHERE recipes.id = (SELECT recipe_id
						FROM user_recipe
						WHERE user_id = ?
							AND recipe_id = ?)`

// DeleteRecipeIngredients deletes all ingredients from a recipe.
const DeleteRecipeIngredients = `
	DELETE
	FROM ingredient_recipe
	WHERE recipe_id = ?`

// DeleteRecipeInstructions deletes all instructions from a recipe.
const DeleteRecipeInstructions = `
	DELETE
	FROM instruction_recipe
	WHERE recipe_id = ?`

// DeleteRecipeImages deletes the images of user's recipe.
const DeleteRecipeImages = `
	DELETE
	FROM additional_images_recipe
	WHERE recipe_id = (SELECT recipe_id FROM user_recipe WHERE recipe_id = ? AND user_id = ?)`

// DeleteRecipesUser deletes the user's recipes.
const DeleteRecipesUser = `
	DELETE
	FROM recipes
	WHERE id IN (SELECT recipe_id
				 FROM user_recipe
				 WHERE user_id = ?)`

// DeleteUser deletes a user from the users table.
const DeleteUser = `
	DELETE
	FROM users
	WHERE id = ?`

// DeleteUserCategory deletes a user's recipe category.
const DeleteUserCategory = `
	DELETE
	FROM user_category
	WHERE user_id = ?
	  AND category_id = (SELECT id FROM categories WHERE name = ?)
	RETURNING category_id`
