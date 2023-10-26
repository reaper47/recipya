package models_test

import (
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"testing"
)

func TestViewModeFromInt(t *testing.T) {
	testcases := []struct {
		in   int64
		want models.ViewMode
	}{
		{in: 0, want: models.GridViewMode},
		{in: 1, want: models.ListViewMode},
		{in: 2, want: models.GridViewMode},
	}
	for _, tc := range testcases {
		t.Run(strconv.Itoa(int(tc.in)), func(t *testing.T) {
			got := models.ViewModeFromInt(tc.in)
			if got != tc.want {
				t.Fatalf("got %d but want %d", got, tc.want)
			}
		})
	}
}

func TestViewModeFromString(t *testing.T) {
	testcases := []struct {
		in   string
		want models.ViewMode
	}{
		{in: "grid", want: models.GridViewMode},
		{in: "list", want: models.ListViewMode},
		{in: "hello", want: models.GridViewMode},
	}
	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			got := models.ViewModeFromString(tc.in)
			if got != tc.want {
				t.Fatalf("got %d but want %d", got, tc.want)
			}
		})
	}
}
