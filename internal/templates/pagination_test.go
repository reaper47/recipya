package templates_test

import (
	"github.com/reaper47/recipya/internal/templates"
	"slices"
	"testing"
)

func TestNewPagination(t *testing.T) {
	testcases := []struct {
		name string
		in   uint64
		want templates.Pagination
	}{
		{
			name: "paginate no results in db",
			in:   2,
			want: templates.Pagination{
				Left:       []uint64{},
				Middle:     []uint64{},
				Right:      []uint64{},
				Prev:       1,
				Selected:   2,
				Next:       3,
				NumPages:   0,
				NumResults: 0,
			},
		},
		{
			name: "paginate couple of results in db",
			in:   2,
			want: templates.Pagination{
				Left:       []uint64{},
				Middle:     []uint64{},
				Right:      []uint64{},
				Prev:       1,
				Selected:   2,
				Next:       3,
				NumPages:   0,
				NumResults: 0,
			},
		},
		{
			name: "paginate hundreds of results in db select left page",
			in:   4,
			want: templates.Pagination{
				Left:       []uint64{1, 2, 3, 4, 5, 6, 7},
				Middle:     []uint64{},
				Right:      []uint64{22},
				Prev:       3,
				Selected:   4,
				Next:       5,
				NumPages:   22,
				NumResults: 258,
			},
		},
		{
			name: "paginate hundreds of results in db select middle page",
			in:   11,
			want: templates.Pagination{
				Left:       []uint64{1},
				Middle:     []uint64{9, 10, 11, 12, 13},
				Right:      []uint64{22},
				Prev:       10,
				Selected:   11,
				Next:       12,
				NumPages:   22,
				NumResults: 258,
			},
		},
		{
			name: "paginate some results",
			in:   1,
			want: templates.Pagination{
				Left:       []uint64{1, 2},
				Middle:     []uint64{},
				Right:      []uint64{},
				Prev:       1,
				Selected:   1,
				Next:       2,
				NumPages:   2,
				NumResults: 20,
			},
		},
		{
			name: "paginate hundreds of results in db select right page",
			in:   20,
			want: templates.Pagination{
				Left:       []uint64{1},
				Middle:     []uint64{},
				Right:      []uint64{17, 18, 19, 20, 21, 22},
				Prev:       19,
				Selected:   20,
				Next:       21,
				NumPages:   22,
				NumResults: 258,
			},
		},
	}
	for _, tc := range testcases {
		t.Run("test pagination init "+tc.name, func(t *testing.T) {
			got := templates.NewPagination(tc.in, tc.want.NumPages, tc.want.NumResults, 12, "/recipes", "", templates.PaginationHtmx{})

			if !slices.Equal(got.Left, tc.want.Left) {
				t.Errorf("Left: wanted %#v but got %#v", tc.want.Left, got.Left)
			}
			if !slices.Equal(got.Middle, tc.want.Middle) {
				t.Errorf("Middle: wanted %v but got %v", tc.want.Middle, got.Middle)
			}
			if !slices.Equal(got.Right, tc.want.Right) {
				t.Errorf("Right: wanted %v but got %v", tc.want.Right, got.Right)
			}
			if got.Prev != tc.want.Prev {
				t.Errorf("Prev: wanted %v but got %v", tc.want.Prev, got.Prev)
			}
			if got.Selected != tc.want.Selected {
				t.Errorf("Selected: wanted %v but got %v", tc.want.Selected, got.Selected)
			}
			if got.Next != tc.want.Next {
				t.Errorf("Next: wanted %v but got %v", tc.want.Next, got.Next)
			}
			if got.NumPages != tc.want.NumPages {
				t.Errorf("NumPages: wanted %v but got %v", tc.want.NumPages, got.NumPages)
			}
			if got.NumResults != tc.want.NumResults {
				t.Errorf("NumResults: wanted %v but got %v", tc.want.NumResults, got.NumResults)
			}
		})
	}
}
