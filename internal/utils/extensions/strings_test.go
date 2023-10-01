package extensions_test

import (
	"github.com/reaper47/recipya/internal/utils/extensions"
	"testing"
)

func TestFloatToString(t *testing.T) {
	testcases := []struct {
		name   string
		in     float64
		format string
		want   string
	}{
		{
			name:   "plain float",
			in:     3,
			format: "%f",
			want:   "3",
		},
		{
			name:   "decimal with trailing zeroes",
			in:     3.140000,
			format: "%.2f",
			want:   "3.14",
		},
		{
			name:   "no trailing zeroes",
			in:     3.14159,
			format: "%f",
			want:   "3.14159",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := extensions.FloatToString(tc.in, tc.format)
			if got != tc.want {
				t.Fatalf("got %q but want %q", got, tc.want)
			}
		})
	}
}
