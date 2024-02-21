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
	compare(t, got, want)
}

func TestNewErrorAuthToast(t *testing.T) {
	got := models.NewErrorAuthToast(message)
	want := models.NewErrorToast("Auth Error", message, "")
	compare(t, got, want)
}

func TestNewErrorDBToast(t *testing.T) {
	got := models.NewErrorDBToast(message)
	want := models.NewErrorToast("Database Error", message, "")
	compare(t, got, want)
}

func TestNewErrorFilesToast(t *testing.T) {
	got := models.NewErrorFilesToast(message)
	want := models.NewErrorToast("Files Error", message, "")
	compare(t, got, want)
}

func TestNewErrorFormToast(t *testing.T) {
	got := models.NewErrorFormToast(message)
	want := models.NewErrorToast("Form Error", message, "")
	compare(t, got, want)
}

func TestNewErrorGeneralToast(t *testing.T) {
	got := models.NewErrorGeneralToast(message)
	want := models.NewErrorToast("General Error", message, "")
	compare(t, got, want)
}

func TestNewErrorReqToast(t *testing.T) {
	got := models.NewErrorReqToast(message)
	want := models.NewErrorToast("Request Error", message, "")
	compare(t, got, want)
}

func TestNewInfoToast(t *testing.T) {
	got := models.NewInfoToast(title, message, action)
	want := models.NewInfoToast(title, message, action)
	compare(t, got, want)
}

func TestNewWarningToast(t *testing.T) {
	got := models.NewWarningToast(title, message, action)
	want := models.NewWarningToast(title, message, action)
	compare(t, got, want)
}

func TestNewWarningWSToast(t *testing.T) {
	got := models.NewWarningWSToast(message)
	want := models.NewWarningToast("Websocket", message, "")
	compare(t, got, want)
}

func TestToast_Render(t *testing.T) {
	got := models.NewWarningToast(title, message, action).Render()

	want := `{"showToast":"{\"action\":\"Dismiss\",\"background\":\"alert-warning\",\"message\":\"Jersey is great.\",\"title\":\"Chicken\"}"}`
	if got != want {
		t.Fatalf("got %s but want %s", got, want)
	}
}

func compare(t *testing.T, got, want any) {
	if !cmp.Equal(got, want) {
		t.Log(cmp.Diff(got, want))
		t.Fail()
	}
}
