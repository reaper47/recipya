package models_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

var (
	title   = "Chicken"
	message = "Jersey is great."
)

func TestNewErrorToast(t *testing.T) {
	got := models.NewErrorToast(title, message)

	want := models.NewErrorToast(title, message)
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}

func TestNewInfoToast(t *testing.T) {
	got := models.NewInfoToast(title, message)

	want := models.NewInfoToast(title, message)
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}

func TestNewWarningToast(t *testing.T) {
	got := models.NewWarningToast(title, message)

	want := models.NewWarningToast(title, message)
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}

func TestToast_Render(t *testing.T) {
	got := models.NewWarningToast(title, message).Render()

	want := `{"showToast":"{\"background\":\"alert-warning\",\"message\":\"Jersey is great.\",\"title\":\"Chicken\"}"}`
	if got != want {
		t.Fatalf("got %s but want %s", got, want)
	}
}
