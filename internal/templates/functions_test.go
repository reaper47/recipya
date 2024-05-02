package templates

import (
	"github.com/google/uuid"
	"slices"
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

	testcases3 := []struct {
		name string
		in   uuid.UUID
		want bool
	}{
		{name: "invalid", in: uuid.UUID{}, want: false},
		{name: "valid", in: uuid.New(), want: true},
	}
	for _, tc := range testcases3 {
		t.Run("UUID is "+tc.name, func(t *testing.T) {
			actual := isUUIDsValid([]uuid.UUID{tc.in})

			expected := tc.want
			if !slices.Equal(actual, []bool{expected}) {
				t.Fatalf("wanted %t but got %t", expected, actual)
			}
		})
	}
}
