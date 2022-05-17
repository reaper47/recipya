package regex_test

import (
	"regexp"
	"testing"

	"github.com/reaper47/recipya/internal/regex"
	"golang.org/x/exp/slices"
)

func TestRegex(t *testing.T) {
	testcases := []struct {
		name          string
		regex         *regexp.Regexp
		in            string
		want          bool
		wantedMatches []string
	}{
		{
			name:  "email is valid",
			regex: regex.Email,
			in:    "xyz@gmail.com",
			want:  true,
		},
		{
			name:  "email is invalid 1",
			regex: regex.Email,
			in:    "xyzgmail.com",
			want:  false,
		},
		{
			name:  "email is invalid 2",
			regex: regex.Email,
			in:    "@gmail.com",
			want:  false,
		},
		{
			name:  "anchor tag is valid",
			regex: regex.Anchor,
			in:    `<a slot="guide-links-primary" href="https://www.youtube.com/about/press/" style="display: none;">`,
			want:  true,
		},
		{
			name:          "hour minutes is valid",
			regex:         regex.HourMinutes,
			in:            "3:5",
			wantedMatches: []string{"3:5"},
		},
		{
			name:          "hour minutes is valid",
			regex:         regex.HourMinutes,
			in:            "3:50",
			wantedMatches: []string{"3:50"},
		},
		{
			name:          "hour minutes is valid",
			regex:         regex.HourMinutes,
			in:            "03:50",
			wantedMatches: []string{"03:50"},
		},
		{
			name:  "hour minutes is invalid",
			regex: regex.HourMinutes,
			in:    ":50",
			want:  false,
		},
		{
			name:  "hour minutes is invalid",
			regex: regex.HourMinutes,
			in:    "3h5m",
			want:  false,
		},
		{
			name:  "hour minutes is invalid",
			regex: regex.HourMinutes,
			in:    "3:60",
			want:  false,
		},
		{
			name:  "image source is valid",
			regex: regex.ImageSrc,
			in:    `<img src="https://imagesvc.meredithcorp.io/v3/mm/image?url=https%3A%2F%2Fstatic.onecms.io%2Fwp-content%2Fuploads%2Fsites%2F43%2F2022%2F03%2F29%2FBohemian-Orange-Chicken-2000.jpg&amp;q=60" alt="Eight chicken thighs topped with orange zest and a sweet and sour orange sauce in a cast iron skillet" title="Bohemian Orange Chicken" width="250">`,
			wantedMatches: []string{
				"https://imagesvc.meredithcorp.io/v3/mm/image?url=https%3A%2F%2Fstatic.onecms.io%2Fwp-content%2Fuploads%2Fsites%2F43%2F2022%2F03%2F29%2FBohemian-Orange-Chicken-2000.jpg&amp;q=60",
			},
		},
		{
			name:  "image source is invalid",
			regex: regex.ImageSrc,
			in:    `<img src="" alt="Eight chicken thighs topped with orange zest and a sweet and sour orange sauce in a cast iron skillet" title="Bohemian Orange Chicken" width="250">`,
			want:  false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.wantedMatches) > 0 {
				actual := tc.regex.FindAllString(tc.in, -1)
				for _, v := range tc.wantedMatches {
					if slices.Index(actual, v) == -1 {
						t.FailNow()
					}
				}
			} else {
				actual := tc.regex.MatchString(tc.in)
				if actual != tc.want {
					t.Fatalf("got %v but want %v for %s", actual, tc.want, tc.in)
				}
			}
		})
	}
}
