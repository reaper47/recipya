package statements

// SelectAuthToken fetches a non-expired auth token by the selector.
const SelectAuthToken = `
	SELECT id, hash_validator, expires, user_id
	FROM auth_tokens
	WHERE selector = ?
	AND expires > unixepoch('now')`

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
