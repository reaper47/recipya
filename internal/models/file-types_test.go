package models_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
)

func TestFileType_Ext(t *testing.T) {
	testcases := []struct {
		name string
		in   models.FileType
		want string
	}{
		{name: "json", in: models.JSON, want: ".json"},
		{name: "pdf", in: models.PDF, want: ".pdf"},
		{name: "mxp", in: models.MXP, want: ".mxp"},
		{name: "txt", in: models.TXT, want: ".txt"},
		{name: "invalid", in: models.InvalidFileType, want: ""},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.Ext()
			if got != tc.want {
				t.Fatalf("got %q but want %q", got, tc.want)
			}
		})
	}
}

func TestNewFileType(t *testing.T) {
	testcases := []struct {
		name string
		in   string
		want models.FileType
	}{
		{name: "json", in: "json", want: models.JSON},
		{name: "pdf", in: "pdf", want: models.PDF},
		{name: "mxp", in: "mxp", want: models.MXP},
		{name: "txt", in: "txt", want: models.TXT},
		{name: "invalid", in: "", want: models.InvalidFileType},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := models.NewFileType(tc.in)
			if got != tc.want {
				t.Fatalf("got %q but want %q", got, tc.want)
			}
		})
	}
}
