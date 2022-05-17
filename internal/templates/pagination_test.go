package templates_test

import (
	"testing"

	"github.com/reaper47/recipya/internal/templates"
	"golang.org/x/exp/slices"
)

func TestNewPagination(t *testing.T) {
	testcases1 := []struct {
		name string
		in   templates.FormErrorData
		want bool
	}{
		{
			name: "all fields empty",
			in:   templates.FormErrorData{},
			want: true,
		},
		{
			name: "one field not empty",
			in:   templates.FormErrorData{Username: "mac"},
			want: false,
		},
		{
			name: "no fields empty",
			in:   templates.FormErrorData{Username: "a", Email: "b", Password: "c"},
			want: false,
		},
	}
	for _, tc := range testcases1 {
		t.Run("test FormErrorData.IsEmpty "+tc.name, func(t *testing.T) {
			actual := tc.in.IsEmpty()
			if actual != tc.want {
				t.Fatalf("IsEmpty: got %v but want %v", actual, tc.want)
			}
		})
	}

	testcases2 := []struct {
		name string
		in   int
		want templates.Pagination
	}{
		{
			name: "negative page number",
			in:   -1,
			want: templates.Pagination{},
		},
		{
			name: "paginate no results in db",
			in:   2,
			want: templates.Pagination{
				Left:       []int{},
				Middle:     []int{},
				Right:      []int{},
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
				Left:       []int{},
				Middle:     []int{},
				Right:      []int{},
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
				Left:       []int{1, 2, 3, 4, 5, 6, 7},
				Middle:     []int{},
				Right:      []int{22},
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
				Left:       []int{1},
				Middle:     []int{9, 10, 11, 12, 13},
				Right:      []int{22},
				Prev:       10,
				Selected:   11,
				Next:       12,
				NumPages:   22,
				NumResults: 258,
			},
		},
		{
			name: "paginate hundreds of results in db select right page",
			in:   20,
			want: templates.Pagination{
				Left:       []int{1},
				Middle:     []int{},
				Right:      []int{17, 18, 19, 20, 21, 22},
				Prev:       19,
				Selected:   20,
				Next:       21,
				NumPages:   22,
				NumResults: 258,
			},
		},
	}
	for _, tc := range testcases2 {
		t.Run("test pagination init "+tc.name, func(t *testing.T) {
			actual := templates.NewPagination(tc.want.NumPages, tc.want.NumResults, tc.in)

			if !slices.Equal(actual.Left, tc.want.Left) {
				t.Errorf("Left: wanted %#v but got %#v", tc.want.Left, actual.Left)
			}
			if !slices.Equal(actual.Middle, tc.want.Middle) {
				t.Errorf("Middle: wanted %#v but got %#v", tc.want.Middle, actual.Middle)
			}
			if !slices.Equal(actual.Right, tc.want.Right) {
				t.Errorf("Right: wanted %#v but got %#v", tc.want.Right, actual.Right)
			}
			if actual.Prev != tc.want.Prev {
				t.Errorf("Prev: wanted %#v but got %#v", tc.want.Prev, actual.Prev)
			}
			if actual.Selected != tc.want.Selected {
				t.Errorf("Selected: wanted %#v but got %#v", tc.want.Selected, actual.Selected)
			}
			if actual.Next != tc.want.Next {
				t.Errorf("Next: wanted %#v but got %#v", tc.want.Next, actual.Next)
			}
			if actual.NumPages != tc.want.NumPages {
				t.Errorf("NumPages: wanted %#v but got %#v", tc.want.NumPages, actual.NumPages)
			}
			if actual.NumResults != tc.want.NumResults {
				t.Errorf(
					"NumResults: wanted %#v but got %#v",
					tc.want.NumResults,
					actual.NumResults,
				)
			}
		})
	}
}
