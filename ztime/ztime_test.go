package ztime

import (
	"fmt"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		in   time.Duration
		want string
	}{
		{10 * time.Second, "10s"},
		{130 * time.Second, "2m10s"},
		{1606 * time.Second, "27m"},
		{3664 * time.Second, "1h1m"},
		{(86400 + 3664) * time.Second, "1d1h"},
		{(86400*15 + 3664) * time.Second, "15d1h"},
		{(86400*17 + 3664*10) * time.Second, "17d10h"},
		{(86400 * 100) * time.Second, "100d"},
		{(86400*204 + 3664*2) * time.Second, "204d"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.in), func(t *testing.T) {
			out := FormatDuration(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %q\nwant: %q", out, tt.want)
			}
		})
	}
}
