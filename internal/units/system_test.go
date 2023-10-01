package units_test

import (
	"github.com/reaper47/recipya/internal/units"
	"testing"
)

func TestSystem_String(t *testing.T) {
	testcases := []struct {
		in   units.System
		want string
	}{
		{units.ImperialSystem, "imperial"},
		{units.InvalidSystem, "invalid"},
		{units.MetricSystem, "metric"},
	}
	for _, tc := range testcases {
		t.Run(tc.want, func(t *testing.T) {
			got := tc.in.String()

			if got != tc.want {
				t.Errorf("got %q but want %q", got, tc.want)
			}
		})
	}
}

func TestNewSystem(t *testing.T) {
	testcases := []struct {
		in   string
		want units.System
	}{
		{"imperial", units.ImperialSystem},
		{"space", units.InvalidSystem},
		{"metric", units.MetricSystem},
	}
	for _, tc := range testcases {
		t.Run(tc.in, func(t *testing.T) {
			got := units.NewSystem(tc.in)

			if got != tc.want {
				t.Errorf("got %q but want %q", got, tc.want)
			}
		})
	}
}
