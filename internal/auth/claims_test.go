package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestAuthClaims(t *testing.T) {
	testcases := []struct {
		name string
		in   customClaims
		want bool
	}{
		{
			name: "claim is expired",
			in: customClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(-1024 * time.Hour)},
				},
				SID: "12345",
			},
			want: false,
		},
		{
			name: "claim has no SID",
			in: customClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(14 * 24 * time.Hour)},
				},
			},
			want: false,
		},
		{
			name: "claim is valid",
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
		t.Run(tc.name, func(t *testing.T) {
			if tc.in.IsValid() != tc.want {
				t.Errorf("IsValid for test #%d: %#v, want %v", i, tc.in, tc.want)
			}
		})
	}
}
