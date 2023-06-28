package templates

// Page is a string alias for the name of an HTML page.
type Page string

// Name of the page. The value is the name of the associated template without the extension.
const (
	AddRecipePage       Page = "add-recipe"
	AddRecipeManualPage Page = "add-recipe-manual"
	HomePage            Page = "home"
	LandingPage         Page = "landing"
	LoginPage           Page = "login"
	RegisterPage        Page = "register"
	Simple              Page = "simple"
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
	case HomePage, LandingPage:
		return "Home"
	case LoginPage:
		return "Login"
	case RegisterPage:
		return "Register"
	case Simple:
		return "<title>"
	default:
		return ""
	}
}
