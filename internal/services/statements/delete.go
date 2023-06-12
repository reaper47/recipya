package statements

// DeleteAuthToken removes the authentication token associated with the user id from the database.
const DeleteAuthToken = `
	DELETE
	FROM auth_tokens
	WHERE user_id = ?`
