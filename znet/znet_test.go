package znet

import (
	"testing"

	"zgo.at/zstd/ztest"
)

func TestSafeDialer(t *testing.T) {
	tests := []struct {
		net, addr string
		want      string
	}{
		{"tcp4", "8.8.8.8:80", ""},
		{"tcp6", "[abab::1]:80", ""},

		{"tcp4", "127.0.0.1:80", "public"},
		{"udp4", "8.8.8.8:80", "allowed"},
		{"tcp4", "8.8.8.8", "host/port pair"},
		{"tcp4", "8.8.8.8:81", "allowed"},
		{"tcp4", "8.8.8.300:80", "invalid"},
	}

	sc := socketControl(nil, nil)
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := sc(tt.net, tt.addr, nil)
			if !ztest.ErrorContains(err, tt.want) {
				t.Errorf("\ngot:  %q\nwant: %q", err, tt.want)
			}
		})
	}
}

func TestRemovePort(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"example.com", "example.com"},
		{"example.com:80", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out := RemovePort(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %q\nwant: %q\n", out, tt.want)
			}
		})
	}
}

func TestPrivateIP(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"", true},
		{"localhost", true},
		{"example.com", true},
		{"$!#!@#\"`'.", true},

		{"127.0.0.1", true},
		{"127.99.0.1", true},
		{"::1", true},
		{"0::1", true},
		{"0000:0000::1", true},
		{"0000:0000:0000:0000:0000:0000:0000:0001", true},
		{"fe80:0000:0000:0000:0000:0000:0000:0001", true},

		{"0000:0000:0000:0000:0000:0000:0000:0002", false},
		{"f081::1", false},
		{"::2", false},
		{"8.8.8.8", false},

		// IPv4 mapped address.
		{"::ffff:169.254.169.254", true},
		{"::ffff:8.8.8.8", false},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out := PrivateIPString(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %t\nwant: %t\n", out, tt.want)
			}
		})
	}
}
