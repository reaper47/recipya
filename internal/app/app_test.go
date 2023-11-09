package app_test

import (
	"github.com/reaper47/recipya/internal/app"
	"testing"
)

func TestConfigFile_Address(t *testing.T) {
	testcases := []struct {
		name string
		in   app.ConfigFile
		want string
	}{
		{
			name: "without port",
			in:   app.ConfigFile{URL: "https://127.0.0.1"},
			want: "https://127.0.0.1",
		},
		{
			name: "with port",
			in:   app.ConfigFile{URL: "https://127.0.0.1", Port: 8078},
			want: "https://127.0.0.1:8078",
		},
		{
			name: "address 0.0.0.0 on linux not affected",
			in:   app.ConfigFile{URL: "https://0.0.0.0", Port: 8078},
			want: "https://0.0.0.0:8078",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.Address()
			if got != tc.want {
				t.Fatalf("got %q but want %q", got, tc.want)
			}
		})
	}
}
