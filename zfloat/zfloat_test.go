package zfloat

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func TestRound(t *testing.T) {
	tests := []struct {
		in   float64
		want float64
	}{
		{123.4999, 123},
		{123.5, 124},
		{123.999, 124},
		{-123.5, -124},
	}

	for _, tt := range tests {
		got := Round(tt.in)
		if got != tt.want {
			t.Errorf("Round(%f) => %f, want %f", tt.in, got, tt.want)
		}
	}

}

func TestRoundPlus(t *testing.T) {
	tests := []struct {
		in        float64
		precision int
		want      float64
	}{
		{123.554999, 3, 123.555},
		{123.555555, 3, 123.556},
		{123.558, 2, 123.56},
		{-123.555555, 3, -123.556},
	}

	for _, tt := range tests {
		got := RoundPlus(tt.in, tt.precision)
		if got != tt.want {
			t.Errorf("Round(%f) => %f, want %f", tt.in, got, tt.want)
		}
	}

}

func TestIsSignedZero(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"1", false},
		{"0", false},
		{"-1", false},
		{"-0", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			f, err := strconv.ParseFloat(tt.in, 64)
			if err != nil {
				t.Fatal(err)
			}

			out := IsSignedZero(f)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestByte(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{500, "500.0B"},
		{1023, "1023.0B"},
		{1024, "1.0KiB"},
		{1424, "1.4KiB"},
		{152310, "148.7KiB"},
		{1024 * 1190, "1.2MiB"},
		{(math.Pow(1024, 5) * 3) + (math.Pow(1024, 4) * 400), "3.4PiB"},
		{(math.Pow(1024, 6) * 3) + (math.Pow(1024, 5) * 400), "3472.0PiB"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out := Byte(tc.in).String()
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}
