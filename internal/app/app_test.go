package app_test

import (
	"github.com/reaper47/recipya/internal/app"
	"strconv"
	"strings"
	"testing"
)

func TestConfigFile_Address(t *testing.T) {
	testcases := []struct {
		name           string
		in             app.ConfigFile
		want           string
		wantRandomPort bool
	}{
		{
			name:           "without port",
			in:             app.ConfigFile{Server: app.ConfigServer{URL: "https://localhost"}},
			want:           "https://localhost",
			wantRandomPort: true,
		},
		{
			name:           "with port",
			in:             app.ConfigFile{Server: app.ConfigServer{URL: "https://127.0.0.1", Port: 8078}},
			want:           "https://127.0.0.1:8078",
			wantRandomPort: false,
		},
		{
			name:           "hosted somewhere",
			in:             app.ConfigFile{Server: app.ConfigServer{URL: "https://recipya.com", Port: 8078, IsProduction: true}},
			want:           "https://recipya.com",
			wantRandomPort: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.Address()
			if got != tc.want {
				if tc.wantRandomPort {
					split := strings.Split(got, ":")
					port, err := strconv.ParseInt(split[len(split)-1], 10, 64)
					if err != nil || port == 0 {
						t.Fatal("port should not be 0")
						return
					}
				} else {
					t.Fatalf("got %q but want %q", got, tc.want)
				}
			}
			if tc.wantRandomPort && tc.in.Server.Port == 0 {
				t.Fatalf("expected random port %t %d", tc.wantRandomPort, app.Config.Server.Port)
			}
		})
	}
}
