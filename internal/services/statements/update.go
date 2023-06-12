package statements

// UpdateIsConfirmed sets the user's account confirmed to true.
const UpdateIsConfirmed = `
	UPDATE users
	SET is_confirmed = 1
	WHERE id = ?`
