//go:build windows

package app_test

import (
	"github.com/reaper47/recipya/internal/app"
	"net"
	"testing"
)

func TestConfigFile_Address_Windows(t *testing.T) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	testcases := []struct {
		name string
		in   app.ConfigFile
		want string
	}{
		{
			name: "hosted locally",
			in:   app.ConfigFile{Server: app.ConfigServer{URL: "http://localhost", Port: 8078}},
			want: "http://localhost:8078",
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
