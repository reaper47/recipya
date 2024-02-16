package models_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

var (
	title   = "Chicken"
	message = "Jersey is great."
	action  = "Dismiss"
)

func TestNewErrorToast(t *testing.T) {
	got := models.NewErrorToast(title, message, action)

	want := models.NewErrorToast(title, message, action)
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}

func TestNewInfoToast(t *testing.T) {
	got := models.NewInfoToast(title, message, action)

	want := models.NewInfoToast(title, message, action)
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}

func TestNewWarningToast(t *testing.T) {
	got := models.NewWarningToast(title, message, action)

	want := models.NewWarningToast(title, message, action)
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}

func TestToast_Render(t *testing.T) {
	got := models.NewWarningToast(title, message, action).Render()

	want := `{"showToast":"{\"action\":\"Dismiss\",\"background\":\"alert-warning\",\"message\":\"Jersey is great.\",\"title\":\"Chicken\"}"}`
	if got != want {
		t.Fatalf("got %s but want %s", got, want)
	}
}
