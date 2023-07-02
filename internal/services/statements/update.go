package statements

// UpdateIsConfirmed sets the user's account confirmed to true.
const UpdateIsConfirmed = `
	UPDATE users
	SET is_confirmed = 1
	WHERE id = ?`

// UpdatePassword sets the user's new password.
const UpdatePassword = `
	UPDATE users
	SET hashed_password = ?, updated_at = CURRENT_TIMESTAMP 
	WHERE id = ?`
