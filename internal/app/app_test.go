package app_test

import (
	"github.com/reaper47/recipya/internal/app"
	"net"
	"testing"
)

func TestConfigFile_Address(t *testing.T) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()

	testcases := []struct {
		name string
		in   app.ConfigFile
		want string
	}{
		{
			name: "without port",
			in:   app.ConfigFile{Server: app.ConfigServer{URL: "https://127.0.0.1"}},
			want: "https://" + ip,
		},
		{
			name: "with port",
			in:   app.ConfigFile{Server: app.ConfigServer{URL: "https://127.0.0.1", Port: 8078}},
			want: "https://" + net.JoinHostPort(ip, "8078"),
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
