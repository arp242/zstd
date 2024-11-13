package zstrconv

import (
	"testing"

	"zgo.at/zstd/ztest"
)

func TestParseInt(t *testing.T) {
	tests := []struct {
		in      string
		want    int32
		wantErr string
	}{
		{"123", 123, ""},
		{"-0x99", -0x99, ""},
		{"asd", 0, "invalid syntax"},
		{"99999999999999999999", 2147483647, "value out of range"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have, err := ParseInt[int32](tt.in, 0)
			if !ztest.ErrorContains(err, tt.wantErr) {
				t.Fatal(err)
			}
			if have != tt.want {
				t.Errorf("\nhave: %d\nwant: %d", have, tt.want)
			}
		})
	}
}
