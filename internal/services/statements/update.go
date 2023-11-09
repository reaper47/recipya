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
