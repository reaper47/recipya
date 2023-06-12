package statements

// InsertAuthToken is the query to add an authentication token to the database.
const InsertAuthToken = `
	INSERT INTO auth_tokens (selector, hash_validator, user_id)
	VALUES (?, ?, ?)`

// InsertUser is the query to add a user to the database.
const InsertUser = `
	INSERT INTO users (email, hashed_password)
	VALUES (?, ?)`
