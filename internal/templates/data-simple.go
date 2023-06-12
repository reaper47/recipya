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
