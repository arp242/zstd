package zcompress

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"zgo.at/zstd/ztest"
)

type testClose struct {
	io.Reader
	didClose bool
}

func (tc *testClose) Close() error { tc.didClose = true; return nil }

func TestAutoGzip(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		r, err := AutoGzip(strings.NewReader("Hello, world!\n"))
		if err != nil {
			t.Fatal(err)
		}
		have, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(have, []byte("Hello, world!\n")) {
			t.Error(have)
		}
	})

	t.Run("gzip", func(t *testing.T) {
		tc := &testClose{Reader: bytes.NewReader([]byte{0x1f, 0x8b, 0x08, 0x08,
			0xb4, 0xa7, 0xf0, 0x66, 0x00, 0x03, 0x68, 0x00, 0xf3, 0x48, 0xcd,
			0xc9, 0xc9, 0xd7, 0x51, 0x28, 0xcf, 0x2f, 0xca, 0x49, 0x51, 0xe4,
			0x02, 0x00, 0x18, 0xa7, 0x55, 0x7b, 0x0e, 0x00, 0x00, 0x00})}
		r, err := AutoGzip(tc)
		if err != nil {
			t.Fatal(err)
		}
		have, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(have, []byte("Hello, world!\n")) {
			t.Error(have)
		}
		err = r.Close()
		if err != nil {
			t.Fatal(err)
		}
		if !tc.didClose {
			t.Fatal("!didClose")
		}
	})

	t.Run("invalid gzip", func(t *testing.T) {
		r, err := AutoGzip(bytes.NewReader([]byte{0x1f, 0x8b, 0x08, 0x08, 0xb4,
			0xa7, 0xf0, 0x66, 0x00, 0x03, 0x68, 0x00, 0xf3, 0x48, 0xcd, 0xc9,
			0xc9, 0xd7, 0x51, 0x28, 0xce, 0x2f, 0xca, 0x49, 0x51, 0xe4, 0x02,
			0x00, 0x18, 0xa7, 0x55, 0x7b, 0x0e, 0x00, 0x00, 0x00}))
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.ReadAll(r)
		if !ztest.ErrorContains(err, "gzip: invalid checksum") {
			t.Fatal(err)
		}
	})
}
