package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestClaims(t *testing.T) {
	testcases := []struct {
		in   customClaims
		want bool
	}{
		{
			in: customClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(-1024 * time.Hour)},
				},
				SID: "12345",
			},
			want: false,
		},
		{
			in: customClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(14 * 24 * time.Hour)},
				},
			},
			want: false,
		},
		{
			in: customClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(14 * 24 * time.Hour)},
				},
				SID: "12345",
			},
			want: true,
		},
	}
	for i, tc := range testcases {
		if tc.in.IsValid() != tc.want {
			t.Errorf("IsValid for test #%d: %#v, want %v", i, tc.in, tc.want)
		}
	}
}
