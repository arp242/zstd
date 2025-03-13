package zio

import (
	"bytes"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type (
	testWriteCloser struct {
		io.Writer
		closed bool
	}
	testReadCloser struct {
		io.Reader
		closed bool
	}
)

func TestDumpReader(t *testing.T) {
	cases := []struct {
		in   io.ReadCloser
		want string
	}{
		{
			io.NopCloser(strings.NewReader("Hello")),
			"Hello",
		},
		{
			io.NopCloser(strings.NewReader("لوحة المفاتيح العربية")),
			"لوحة المفاتيح العربية",
		},
		{
			http.NoBody,
			"",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			outR1, outR2, err := DumpReader(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			out1 := mustRead(t, outR1)
			out2 := mustRead(t, outR2)

			if out1 != tc.want {
				t.Errorf("out1 wrong\nout:  %#v\nwant: %#v\n", out1, tc.want)
			}
			if out2 != tc.want {
				t.Errorf("out2 wrong\nout:  %#v\nwant: %#v\n", out2, tc.want)
			}
		})
	}
}

func mustRead(t *testing.T, r io.Reader) string {
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func TestExists(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{".", true},               // Dir
		{"zio.go", true},          // File
		{"/dev/null", true},       // Device
		{"/proc/1/environ", true}, // Not readable
		{"/etc/localtime", true},  // Symlink

		{"/nonexistent-path", false},
		{"/nonexistent/path", false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Exists(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestTeeReader(t *testing.T) {
	w1, w2 := new(bytes.Buffer), new(bytes.Buffer)
	tee := TeeReader(strings.NewReader("hello"), w1, w2)

	h, _ := io.ReadAll(tee)
	if string(h) != "hello" {
		t.Errorf("read from TeeWriter: %q", string(h))
	}
	if w1.String() != "hello" {
		t.Errorf("read from w1: %q", w1.String())
	}
	if w2.String() != "hello" {
		t.Errorf("read from w2: %q", w2.String())
	}
}

var (
	_ io.Reader = &peekReader{}
	_ io.Closer = &peekReader{}
	_ io.Closer = &testClose{}
)

type testClose struct {
	io.Reader
	didClose bool
}

func (tc *testClose) Close() error { tc.didClose = true; return nil }

func TestPeekReader(t *testing.T) {
	t.Run("read from both", func(t *testing.T) {
		r := PeekReader(strings.NewReader("hello"), []byte("abc"))
		buf := make([]byte, 10)
		n, err := r.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if n != 8 {
			t.Error(n)
		}
		if h := string(buf); h != "abchello\x00\x00" {
			t.Errorf("%q", h)
		}

		buf = make([]byte, 10)
		n, err = r.Read(buf)
		if n != 0 {
			t.Error(n, string(buf))
		}
		if !errors.Is(err, io.EOF) {
			t.Fatal(err)
		}
	})

	t.Run("multiple reads from peeked", func(t *testing.T) {
		r := PeekReader(strings.NewReader("de"), []byte("abc"))
		for i := 0; i < 5; i++ {
			buf := make([]byte, 1)
			n, err := r.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			if n != 1 {
				t.Error(n)
			}
			want := ""
			switch i {
			case 0:
				want = "a"
			case 1:
				want = "b"
			case 2:
				want = "c"
			case 3:
				want = "d"
			case 4:
				want = "e"
			}
			if h := string(buf); h != want {
				t.Error(h)
			}
		}

		buf := make([]byte, 10)
		n, err := r.Read(buf)
		if n != 0 {
			t.Error(n, string(buf))
		}
		if !errors.Is(err, io.EOF) {
			t.Fatal(err)
		}
	})

	t.Run("empty peeked", func(t *testing.T) {
		r := PeekReader(strings.NewReader("hello"), nil)
		buf := make([]byte, 10)
		n, err := r.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if n != 5 {
			t.Error(n)
		}
		if h := string(buf[:n]); h != "hello" {
			t.Error(h)
		}
	})

	t.Run("close", func(t *testing.T) {
		tc := &testClose{Reader: strings.NewReader("hello")}
		r := PeekReader(tc, nil)
		err := r.Close()
		if err != nil {
			t.Fatal(err)
		}
		if !tc.didClose {
			t.Error("!tc.didClose")
		}
	})
}

func TestCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	want := "hello\nworld\nxxx\n"
	tmp := t.TempDir() + "/tmp"
	err := os.WriteFile(tmp, []byte(want), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	{
		fp, err := os.Open(tmp)
		if err != nil {
			t.Fatal(err)
		}
		cnt, err := Count(ctx, fp, []byte{'\n'})
		if err != nil {
			t.Fatal(err)
		}
		if cnt != 3 {
			t.Fatal(cnt)
		}
		have, err := io.ReadAll(fp)
		if err != nil {
			t.Fatal(err)
		}
		fp.Close()
		if string(have) != want {
			t.Fatal(string(have))
		}
	}

	{
		time.Sleep(60 * time.Millisecond)
		fp, err := os.Open(tmp)
		if err != nil {
			t.Fatal(err)
		}
		cnt, err := Count(ctx, fp, []byte{'\n'})
		if err == nil {
			t.Fatal("error is nil")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatal(err)
		}
		if cnt != 0 {
			t.Fatal(cnt)
		}
		have, err := io.ReadAll(fp)
		if err != nil {
			t.Fatal(err)
		}
		fp.Close()
		if string(have) != want {
			t.Fatal(string(have))
		}
	}
}

func (t *testWriteCloser) Close() error { t.closed = true; return nil }
func (t *testReadCloser) Close() error  { t.closed = true; return nil }

func TestHashWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	closer := &testWriteCloser{buf, false}
	w := HashWriter(closer, sha1.New())

	var have []string
	have = append(have, fmt.Sprintf("%x %d", w.Hash(), w.Len()))
	w.Write([]byte("a"))
	have = append(have, fmt.Sprintf("%x %d", w.Hash(), w.Len()))
	w.Write([]byte("b"))
	have = append(have, fmt.Sprintf("%x %d", w.Hash(), w.Len()))
	if w.Close(); !closer.closed {
		t.Fatal("not closed")
	}

	closer = &testWriteCloser{buf, false}
	w = HashWriter(closer, sha256.New())
	have = append(have, fmt.Sprintf("%x %d", w.Hash(), w.Len()))
	w.Write([]byte("a"))
	have = append(have, fmt.Sprintf("%x %d", w.Hash(), w.Len()))
	w.Write([]byte("b"))
	have = append(have, fmt.Sprintf("%x %d", w.Hash(), w.Len()))
	if w.Close(); !closer.closed {
		t.Fatal("not closed")
	}

	want := []string{
		"da39a3ee5e6b4b0d3255bfef95601890afd80709 0",
		"86f7e437faa5a7fce15d1ddcb9eaeaea377667b8 1",
		"da23614e02469a0d7c7bd1bdab5c9c474b1904dc 2",
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 0",
		"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb 1",
		"fb8e20fc2e4c3f248c60c39bd652f3c1347298bb977b8b4d5903b85055620603 2",
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave: %#v\nwant: %#v", have, want)
	}
}

func TestHashReader(t *testing.T) {
	closer := &testReadCloser{strings.NewReader("abab"), false}
	r := HashReader(closer, sha1.New())

	var have []string
	have = append(have, fmt.Sprintf("%x %d", r.Hash(), r.Len()))
	r.Read(make([]byte, 1))
	have = append(have, fmt.Sprintf("%x %d", r.Hash(), r.Len()))
	r.Read(make([]byte, 1))
	have = append(have, fmt.Sprintf("%x %d", r.Hash(), r.Len()))
	if r.Close(); !closer.closed {
		t.Fatal("not closed")
	}

	closer = &testReadCloser{strings.NewReader("abab"), false}
	r = HashReader(closer, sha256.New())
	have = append(have, fmt.Sprintf("%x %d", r.Hash(), r.Len()))
	r.Read(make([]byte, 1))
	have = append(have, fmt.Sprintf("%x %d", r.Hash(), r.Len()))
	r.Read(make([]byte, 1))
	have = append(have, fmt.Sprintf("%x %d", r.Hash(), r.Len()))
	if r.Close(); !closer.closed {
		t.Fatal("not closed")
	}

	want := []string{
		"da39a3ee5e6b4b0d3255bfef95601890afd80709 0",
		"86f7e437faa5a7fce15d1ddcb9eaeaea377667b8 1",
		"da23614e02469a0d7c7bd1bdab5c9c474b1904dc 2",
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 0",
		"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb 1",
		"fb8e20fc2e4c3f248c60c39bd652f3c1347298bb977b8b4d5903b85055620603 2",
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("\nhave: %#v\nwant: %#v", have, want)
	}
}

func TestReadNopCloser(t *testing.T) {
	r := ReadNopCloser(strings.NewReader("abcd"))
	buf := make([]byte, 2)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatal(n)
	}
	err = r.Close()
	if err != nil {
		t.Fatal(err)
	}

	n, err = r.Read(buf)
	if !errors.Is(err, fs.ErrClosed) {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatal(n)
	}
}

func TestWriteNopCloser(t *testing.T) {
	buf := new(bytes.Buffer)
	r := WriteNopCloser(buf)

	n, err := r.Write([]byte("ab"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatal(n)
	}
	err = r.Close()
	if err != nil {
		t.Fatal(err)
	}

	n, err = r.Write([]byte("cd"))
	if !errors.Is(err, fs.ErrClosed) {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatal(n)
	}
}

func TestSlowReader(t *testing.T) {
	t.Run("", func(t *testing.T) {
		start := time.Now()
		r := SlowReader(strings.NewReader("abcd"), 1, 10*time.Millisecond)
		all, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if h := string(all); h != "abcd" {
			t.Error(h)
		}
		if took := time.Since(start); took < time.Millisecond*40 {
			t.Errorf("too fast: %s", took)
		}
	})

	t.Run("", func(t *testing.T) {
		start := time.Now()
		r := SlowReader(strings.NewReader("abcd"), 1, 10*time.Millisecond)

		var all []byte
		for {
			buf := make([]byte, 4)
			n, err := r.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Fatal(err)
			}
			all = append(all, buf[:n]...)
		}
		if h := string(all); h != "abcd" {
			t.Error(h)
		}
		if took := time.Since(start); took < time.Millisecond*40 {
			t.Errorf("too fast: %s", took)
		}
	})

	// t.Run("", func(t *testing.T) {
	// 	r := SlowReader(strings.NewReader("abcd"), 2, 10*time.Millisecond)

	// 	buf := make([]byte, 1)
	// 	n, err := r.Read(buf)
	// 	fmt.Println(n, err, string(buf))
	// })
}
