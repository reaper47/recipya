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

func TestConfigServer_IsCookieSecure(t *testing.T) {
	testcases := []struct {
		name     string
		in       app.ConfigFile
		isSecure bool
	}{
		{
			name:     "is prod localhost",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: false, IsProduction: true, Port: 8078, URL: "http://localhost"}},
			isSecure: true,
		},
		{
			name:     "is demo and prod",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: true, IsProduction: true, Port: 8078, URL: "http://localhost"}},
			isSecure: true,
		},
		{
			name:     "is prod not localhost",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: false, IsProduction: true, Port: 8078, URL: "http://192.168.0.1"}},
			isSecure: false,
		},
		{
			name:     "is not prod",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: false, IsProduction: false, Port: 8078, URL: "http://192.168.0.1"}},
			isSecure: false,
		},
		{
			name:     "is demo not prod",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: true, IsProduction: false, Port: 8078, URL: "http://localhost"}},
			isSecure: false,
		},
		{
			name:     "is hosted website",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: false, IsProduction: true, Port: 8078, URL: "https://www.recipya.com"}},
			isSecure: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.in.IsCookieSecure() != tc.isSecure {
				if tc.isSecure {
					t.Fatal("should have been secure")
				} else {
					t.Fatal("should not have been secure")
				}
			}
		})
	}
}
