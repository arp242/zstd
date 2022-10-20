package ztest

import (
	"errors"
	"os"
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestReplace(t *testing.T) {
	tests := []struct {
		in, want string
		patt     []string
	}{
		{
			"Time: 4.12361 ms", "XXX",
			[]string{`Time: [0-9.]+ ms`},
		},
		{
			"Time: 4.12361 ms", "Time: XXX ms",
			[]string{`Time: ([0-9.]+) ms`},
		},
		{
			"Time: 4.12361 ms", "Time: XXX.XXX ms",
			[]string{`Time: ([0-9]+)\.([0-9]+) ms`},
		},
		{
			"Time: 4.12361 ms", "Time: XXX.XXX XX",
			[]string{`Time: ([0-9]+)\.([0-9]+) ms`, `ms`},
		},
		{
			`
Seq Scan on tbl  (cost=0.00..25.88 rows=6 width=36) (actual time=0.007..0.014 rows=1 loops=1)
  Filter: ((col1)::text = 'hello'::text)
  Rows Removed by Filter: 1
Planning Time: 0.026 ms
Execution Time: 0.055 ms
`,
			`
Seq Scan on tbl  (cost=XXX..XXX rows=6 width=36) (actual time=XXX..XXX rows=1 loops=1)
  Filter: ((col1)::text = 'hello'::text)
  Rows Removed by Filter: 1
Planning Time: XXX ms
Execution Time: XXX ms
`, []string{`([0-9]+.[0-9]+) ms`, `(?:cost|time)=([0-9.]+)\.\.([0-9.]+) `},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Replace(tt.in, tt.patt...)
			if got != tt.want {
				t.Errorf("\ngot:  %s\nwant: %s", got, tt.want)
			}
		})
	}

}

func TestErrorContains(t *testing.T) {
	tests := []struct {
		err  error
		str  string
		want bool
	}{
		{errors.New("Hello"), "Hello", true},
		{errors.New("Hello, world"), "world", true},
		{nil, "", true},

		{errors.New("Hello, world"), "", false},
		{errors.New("Hello, world"), "mars", false},
		{nil, "hello", false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := ErrorContains(tt.err, tt.str)
			if have != tt.want {
				t.Errorf("\nhave: %#v\nwant: %#v\n", have, tt.want)
			}
		})
	}
}

func TestTempFile(t *testing.T) {
	var f string
	t.Run("", func(t *testing.T) {
		f = TempFile(t, "", "hello\nworld")
		_, err := os.Stat(f)
		if err != nil {
			t.Fatalf("stat failed: %s", err)
		}
	})

	_, err := os.Stat(f)
	if err == nil {
		t.Fatalf("stat didn't report any error, but the file should be gone")
	}
}

func TestNormalizeIndent(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := NormalizeIndent(tt.in)
			if have != tt.want {
				t.Errorf("\nhave: %#v\nwant: %#v\n", have, tt.want)
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
			"\n--- have\n+++ want\n@@ -1 +1 @@\n- a\n+ b\n"},
		{"hello\nworld\nxxx", "hello\nmars\nxxx",
			"\n--- have\n+++ want\n@@ -1,3 +1,3 @@\n  hello\n- world\n+ mars\n  xxx\n"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := Diff(tt.inOut, tt.inWant)
			if have != tt.want {
				t.Errorf("\nhave:\n%s\nwant:\n%s\nout:  %[1]q\nwant: %[2]q", have, tt.want)
			}
		})
	}
}

func TestParallel(t *testing.T) {
	tests := []string{"one", "two", "three"}

	var (
		have []string
		lock sync.Mutex
	)
	t.Cleanup(func() {
		sort.Strings(have)
		want := []string{"one", "three", "two"}

		if !reflect.DeepEqual(have, want) {
			t.Errorf("\nhave: %#v\nwant: %#v", have, want)
		}
	})

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			tt := Parallel(t, tt)

			lock.Lock()
			have = append(have, tt)
			lock.Unlock()
		})
	}
}
