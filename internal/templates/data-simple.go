package templates

// ErrorConfirm encapsulates the information displayed to the user when confirming an account fails.
var ErrorConfirm = Data{
	Title:        "Confirm",
	ContentTitle: "Confirm Error",
	Content: `An error occurred when you requested to confirm your account.
				The problem has been forwarded to our team automatically. We will look into it and come
                back to you. We apologise for this inconvenience.`,
}

// ErrorTokenExpired encapsulates the information displayed to the user when a token is expired.
var ErrorTokenExpired = Data{
	Title:        "Token Expired",
	ContentTitle: "Expired",
	Content: `The token associated with the URL expired.
				The problem has been forwarded to our team automatically. We will look into it and come
                back to you. We apologise for this inconvenience.`,
}

// ForgotPasswordSuccess encapsulates the information displayed to the user when the user clicks forgot password.
var ForgotPasswordSuccess = Data{
	Title:        "Forgot Password",
	ContentTitle: "Password Reset Requested",
	Content:      "An email with instructions on how to reset your password has been sent to you. Please check your inbox and follow the provided steps to regain access to your account.",
}

// PageNotFound encapsulates the information displayed when the page is 404.
var PageNotFound = Data{
	Title:        "Not Found",
	ContentTitle: "Page Not Found",
	Content:      "The page you requested to view is not found. Please go back to the main page.",
}

// SuccessConfirm encapsulates the information displayed to the user when confirming an account succeeds.
var SuccessConfirm = Data{
	Title:        "Confirm",
	ContentTitle: "Confirmation Successful",
	Content:      "Your account has been confirmed.",
}

// UserLimitReachedError encapsulates the information displayed to the user when the user limit has been reached.
// The limit depends on the SendGrid API key.
var UserLimitReachedError = Data{
	Title:        "User Limit Reached",
	ContentTitle: "User Limit Reached",
	Content: `You cannot register because the user limit has been reached. This limit depends on the SendGrid API key. 
				You can sponsor the author of this project or buy him a coffee for him to have enough money to purchase
				the paid SendGrid plan to increase the limit. You will find the details here: 
				https://github.com/reaper47/heavy-metal-notifier?tab=readme-ov-file#sponsors.`,
}
