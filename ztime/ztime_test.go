package ztime

import (
	"testing"
	"time"
)

func TestTakes(t *testing.T) {
	tt := Takes(func() { time.Sleep(50 * time.Millisecond) })
	if tt < 50*time.Millisecond || tt > 52*time.Millisecond {
		t.Error(tt)
	}
}

func TestDurationAs(t *testing.T) {
	tests := []struct {
		d, as time.Duration
		want  string
	}{
		{50 * time.Millisecond, time.Microsecond, "50000"},
		{50 * time.Microsecond, time.Millisecond, "0.05"},
		{1, time.Hour, "0.0000000000002"},
		{1261616533, time.Millisecond, "1261.616533"},
		{time.Duration(1261616533).Round(time.Microsecond), time.Millisecond, "1261.617"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := DurationAs(tt.d, tt.as)
			if got != tt.want {
				t.Errorf("\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}

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
		t.Run("", func(t *testing.T) {
			out := FormatDuration(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %q\nwant: %q", out, tt.want)
			}
		})
	}
}
