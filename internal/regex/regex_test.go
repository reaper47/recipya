package regex

import "testing"

func TestRegex(t *testing.T) {
	t.Run("email regex", func(t *testing.T) {
		testcases := []struct {
			in   string
			want bool
		}{
			{
				in:   "xyz@gmail.com",
				want: true,
			},
			{
				in:   "xyzgmail.com",
				want: false,
			},
			{
				in:   "@gmail.com",
				want: false,
			},
		}
		for _, tc := range testcases {
			actual := Email.MatchString(tc.in)
			if actual != tc.want {
				t.Fatalf("wanted %v for %s but got %v", tc.want, tc.in, actual)
			}
		}
	})
}
