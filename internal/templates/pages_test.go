package templates_test

import (
	"github.com/reaper47/recipya/internal/templates"
	"testing"
)

var pages = []templates.Page{
	templates.AddRecipePage,
	templates.HomePage,
	templates.LandingPage,
	templates.LoginPage,
	templates.RegisterPage,
	templates.Simple,
}

func TestPage_String(t *testing.T) {
	expected := []string{
		"add-recipe",
		"home",
		"landing",
		"login",
		"register",
		"simple",
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
		"Home",
		"Home",
		"Login",
		"Register",
		"<title>",
	}
	for i, p := range pages {
		actual := p.Title()
		want := expected[i]
		if actual != want {
			t.Errorf("expected %q but got %q", actual, want)
		}
	}
}
