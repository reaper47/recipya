package templates_test

import (
	"github.com/reaper47/recipya/internal/templates"
	"testing"
)

var pages = []templates.Page{
	templates.AddRecipePage,
	templates.AddRecipeManualPage,
	templates.CookbooksPage,
	templates.ForgotPasswordPage,
	templates.ForgotPasswordResetPage,
	templates.HomePage,
	templates.LandingPage,
	templates.LoginPage,
	templates.RegisterPage,
	templates.SettingsPage,
	templates.Simple,
	templates.ViewRecipePage,
}

func TestPage_String(t *testing.T) {
	expected := []string{
		"add-recipe",
		"add-recipe-manual",
		"cookbooks",
		"forgot-password",
		"forgot-password-reset",
		"home",
		"landing",
		"login",
		"register",
		"settings",
		"simple",
		"view-recipe",
	}
	for i, p := range pages {
		actual := p.String()
		want := expected[i]
		if actual != want {
			t.Errorf("expected %q but got %q", actual, want)
		}
	}
}

func TestPage_Title(t *testing.T) {
	expected := []string{
		"Add Recipe",
		"Manual",
		"Cookbooks",
		"Forgot Password",
		"Reset Password",
		"Home",
		"Home",
		"Login",
		"Register",
		"Settings",
		"<title>",
		"View Recipe",
	}
	for i, p := range pages {
		actual := p.Title()
		want := expected[i]
		if actual != want {
			t.Errorf("expected %q but got %q", actual, want)
		}
	}
}
