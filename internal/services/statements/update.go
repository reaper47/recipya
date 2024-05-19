package statements

// UpdateCalculateNutrition is the query to update the user's calculate nutrition setting.
const UpdateCalculateNutrition = `
	UPDATE user_settings
	SET calculate_nutrition = ?
	WHERE user_id = ?`

// UpdateConvertAutomatically is the query to update the user's convert automatically setting.
const UpdateConvertAutomatically = `
	UPDATE user_settings
	SET convert_automatically = ?
	WHERE user_id = ?`

// UpdateCookbookImage is the query to update the image of a user's cookbook.
const UpdateCookbookImage = `
UPDATE cookbooks
	SET image = ?
	WHERE user_id = ?
	 AND id = ?`

// UpdateCookbookRecipesReorder is the query to reorder recipes in a cookbook.
const UpdateCookbookRecipesReorder = `
	UPDATE cookbook_recipes
	SET order_index = ?
	WHERE cookbook_id = ?
		AND recipe_id = ?`

// UpdateIsConfirmed sets the user's account confirmed to true.
const UpdateIsConfirmed = `
	UPDATE users
	SET is_confirmed = 1
	WHERE id = ?`

// UpdateMeasurementSystem is the query to update the user's preferred measurement system.
const UpdateMeasurementSystem = `
	UPDATE user_settings
	SET measurement_system_id = (SELECT id FROM measurement_systems WHERE name = ?)
	WHERE user_id = ?`

// UpdateNutrition updates the recipe's nutrition.
const UpdateNutrition = `
	UPDATE nutrition 
	SET calories = ?,
	    total_carbohydrates = ?,
	    sugars = ?, 
	    protein = ?, 
	    total_fat = ?, 
	    saturated_fat = ?, 
	    unsaturated_fat = ?, 
	    cholesterol = ?, 
	    sodium = ?, 
	    fiber = ?,
	    is_per_serving = ?
	WHERE recipe_id = ?`

// UpdatePassword sets the user's new password.
const UpdatePassword = `
	UPDATE users
	SET hashed_password = ?, updated_at = CURRENT_TIMESTAMP 
	WHERE id = ?`

// UpdateRecipeCategory is the query to update a recipe's category.
const UpdateRecipeCategory = `
	UPDATE category_recipe
	SET category_id = ?
	WHERE id = ?`

// UpdateRecipeCategoryReset is the query to reset the category of the user's affected recipes.
const UpdateRecipeCategoryReset = `
	UPDATE category_recipe
	SET category_id = 1
	WHERE recipe_id IN ((SELECT r.id
						 FROM recipes AS r
								  INNER JOIN category_recipe AS cr ON cr.recipe_id = r.id
								  INNER JOIN user_recipe AS ur ON ur.recipe_id = r.id
						 WHERE cr.category_id = ?
						   AND ur.user_id = ?))`

// UpdateRecipeDescription is the query to update a recipe's description.
const UpdateRecipeDescription = `
	UPDATE recipes
	SET description = ?
	WHERE id = ?`

// UpdateRecipeIngredient is the query to update a recipe's ingredient.
const UpdateRecipeIngredient = `
	UPDATE ingredient_recipe
	SET ingredient_id = ?
	WHERE ingredient_id = (SELECT id FROM ingredients WHERE name = ?)
	  AND recipe_id = ?`

// UpdateRecipeInstruction is the query to update a recipe's instruction.
const UpdateRecipeInstruction = `
	UPDATE instruction_recipe
	SET instruction_id = ?
	WHERE instruction_id = (SELECT id FROM instructions WHERE name = ?)
	  AND recipe_id = ?`

// UpdateRecipeTimes is the query to update a recipe's times.
const UpdateRecipeTimes = `
	UPDATE time_recipe
	SET time_id = ?
	WHERE recipe_id = ?`

// UpdateUserSettingsCookbooksViewMode is the query to update the cookbooks_view column of a user's settings.
const UpdateUserSettingsCookbooksViewMode = `
	UPDATE user_settings
	SET cookbooks_view = ?
	WHERE user_id = ?`

// UpdateIsUpdateAvailable is the query to flag whether a release update is available.
const UpdateIsUpdateAvailable = `
	UPDATE app
	SET is_update_available = ?
	WHERE id = 1
	RETURNING updated_at, update_last_checked_at`
