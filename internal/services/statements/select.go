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
