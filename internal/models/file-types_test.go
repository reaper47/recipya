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
		{name: "cml", in: models.CML, want: ".cml"},
		{name: "crouton", in: models.Crumb, want: ".crumb"},
		{name: "json", in: models.JSON, want: ".json"},
		{name: "mxp", in: models.MXP, want: ".mxp"},
		{name: "paprika", in: models.Paprika, want: ".paprikarecipes"},
		{name: "pdf", in: models.PDF, want: ".pdf"},
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
		{name: "cml", in: "cml", want: models.CML},
		{name: "crumb", in: "crumb", want: models.Crumb},
		{name: "json", in: "json", want: models.JSON},
		{name: "mxp", in: "mxp", want: models.MXP},
		{name: "paprikarecipes", in: "paprikarecipes", want: models.Paprika},
		{name: "pdf", in: "pdf", want: models.PDF},
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
