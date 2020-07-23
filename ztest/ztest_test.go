package ztest

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestErrorContains(t *testing.T) {
	cases := []struct {
		err      error
		str      string
		expected bool
	}{
		{errors.New("Hello"), "Hello", true},
		{errors.New("Hello, world"), "world", true},
		{nil, "", true},

		{errors.New("Hello, world"), "", false},
		{errors.New("Hello, world"), "mars", false},
		{nil, "hello", false},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.err), func(t *testing.T) {
			out := ErrorContains(tc.err, tc.str)
			if out != tc.expected {
				t.Errorf("\nout:      %#v\nexpected: %#v\n", out, tc.expected)
			}
		})
	}
}

func TestTempFile(t *testing.T) {
	f, clean := TempFile(t, "hello\nworld")

	_, err := os.Stat(f)
	if err != nil {
		t.Fatal(err)
	}

	clean()

	_, err = os.Stat(f)
	if err == nil {
		t.Fatal(err)
	}
}

func TestNormalizeIndent(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{
			"\t\twoot\n\t\twoot\n",
			"woot\nwoot",
		},
		{
			"\t\twoot\n\t\t woot",
			"woot\n woot",
		},
		{
			"\t\twoot\n\t\t\twoot",
			"woot\n\twoot",
		},
		{
			"woot\n\twoot",
			"woot\n\twoot",
		},
		{
			"  woot\n\twoot",
			"woot\n\twoot",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := NormalizeIndent(tc.in)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		inOut, inWant string
		want          string
	}{
		{"", "", ""},
		//{nil, nil, ""},

		{"a", "a", ""},
		{"a", "a", ""},
		{"a", "b",
			"\n--- output\n+++ want\n@@ -1 +1 @@\n- a\n+ b\n"},
		{"hello\nworld\nxxx", "hello\nmars\nxxx",
			"\n--- output\n+++ want\n@@ -1,3 +1,3 @@\n  hello\n- world\n+ mars\n  xxx\n"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			out := Diff(tt.inOut, tt.inWant)
			if out != tt.want {
				t.Errorf("\nout:\n%s\nwant:\n%s\nout:  %[1]q\nwant: %[2]q", out, tt.want)
			}
		})
	}
}
