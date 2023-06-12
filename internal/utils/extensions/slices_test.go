package extensions_test

import (
	"github.com/reaper47/recipya/internal/utils/extensions"
	"golang.org/x/exp/slices"
	"testing"
)

func TestUnique(t *testing.T) {
	testcases := []struct {
		name string
		in   []any
		want []any
	}{
		{
			name: "integers",
			in:   []any{1, 2, 1, 3, 4, 5, 3, 6, 7, 8, 5},
			want: []any{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name: "only duplicates",
			in:   []any{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			want: []any{1},
		},
		{
			name: "strings",
			in:   []any{"test", "one", "test", "test", "testy", "two"},
			want: []any{"test", "one", "testy", "two"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := extensions.Unique(tc.in)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("got %v but want %v", got, tc.want)
			}
		})
	}
}
