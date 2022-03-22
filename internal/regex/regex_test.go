package regex

import "testing"

func TestRegex(t *testing.T) {
	t.Run("email regex", func(t *testing.T) {
		testcases := []struct {
			name string
			in   string
			want bool
		}{
			{
				name: "email is valid",
				in:   "xyz@gmail.com",
				want: true,
			},
			{
				name: "email is invalid 1",
				in:   "xyzgmail.com",
				want: false,
			},
			{
				name: "email is invalid 2",
				in:   "@gmail.com",
				want: false,
			},
		}
		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				actual := Email.MatchString(tc.in)
				if actual != tc.want {
					t.Fatalf("wanted %v for %s but got %v", tc.want, tc.in, actual)
				}
			})
		}
	})
}
