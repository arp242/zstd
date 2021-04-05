package ztime

import (
	"testing"
	"time"
)

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
			have := DurationAs(tt.d, tt.as)
			if have != tt.want {
				t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
			}
		})
	}
}
