package templates

import (
	"testing"
	"time"
)

func TestTemplatesFunctions(t *testing.T) {
	testcases1 := []struct {
		name string
		in   bool
		want string
	}{
		{name: "is datetime", in: true, want: "PT2H45M"},
		{name: "is not datetime", in: false, want: "2h45m"},
	}
	for _, tc := range testcases1 {
		t.Run("format duration "+tc.name, func(t *testing.T) {
			d, _ := time.ParseDuration("2h45m")
			actual := formatDuration(d, tc.in)

			expected := tc.want
			if actual != expected {
				t.Fatalf("wanted %s but got %s", expected, actual)
			}
		})
	}

	testcases2 := []struct {
		name string
		in   string
		want bool
	}{
		{name: "is true", in: "https://www.google.com", want: true},
		{name: "is false", in: "google.com", want: false},
	}
	for _, tc := range testcases2 {
		t.Run("is url "+tc.name, func(t *testing.T) {
			actual := isURL(tc.in)

			expected := tc.want
			if actual != expected {
				t.Fatalf("wanted %t but got %t", expected, actual)
			}
		})
	}
}
