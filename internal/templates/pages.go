package templates

// Page is a string alias for the name of an HTML page.
type Page string

// Name of the page. The value is the name of the associated template without the extension.
const (
	AddRecipePage           Page = "add-recipe"
	AddRecipeManualPage     Page = "add-recipe-manual"
	CookbookPage            Page = "cookbook"
	CookbooksPage           Page = "cookbooks"
	ForgotPasswordPage      Page = "forgot-password"
	ForgotPasswordResetPage Page = "forgot-password-reset"
	HomePage                Page = "home"
	LandingPage             Page = "landing"
	LoginPage               Page = "login"
	RegisterPage            Page = "register"
	SettingsPage            Page = "settings"
	SimplePage              Page = "simple"
	ViewRecipePage          Page = "view-recipe"
)

// String stringifies the Page.
func (p Page) String() string {
	return string(p)
}

// Title returns the title of Page for the <title> tag.
func (p Page) Title() string {
	switch p {
	case AddRecipePage:
		return "Add Recipe"
	case AddRecipeManualPage:
		return "Manual"
	case CookbookPage:
		return "Cookbook"
	case CookbooksPage:
		return "Cookbooks"
	case ForgotPasswordPage:
		return "Forgot Password"
	case ForgotPasswordResetPage:
		return "Reset Password"
	case HomePage, LandingPage:
		return "Home"
	case LoginPage:
		return "Login"
	case RegisterPage:
		return "Register"
	case SettingsPage:
		return "Settings"
	case SimplePage:
		return "<title>"
	case ViewRecipePage:
		return "View Recipe"
	default:
		return ""
	}
}
