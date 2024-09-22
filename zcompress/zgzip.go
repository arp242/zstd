package zcompress

import (
	"bytes"
	"compress/gzip"
	"io"

	"zgo.at/zstd/zio"
)

// AutoGzip automatically decompressed the reader if it looks like it's
// compressed with gzip.
//
// Unlike gzip.Reader, Close() will close both the gzip.Reader and the
// underlying reader (if it has a Close method).
func AutoGzip(r io.Reader) (io.ReadCloser, error) {
	head := make([]byte, 3)
	_, err := r.Read(head)
	if err != nil {
		return nil, err
	}

	p := zio.PeekReader(r, head)
	if !bytes.Equal(head, []byte{0x1f, 0x8b, 0x08}) {
		return p, nil
	}

	g, err := gzip.NewReader(p)
	if err != nil {
		return nil, err
	}
	return &gzipCloser{g, r}, nil
}

type gzipCloser struct {
	*gzip.Reader
	r io.Reader
}

// Close the [gzip.Reader] and the underlying [io.Reader] if it implements a
// Close method.
func (g *gzipCloser) Close() error {
	err := g.Reader.Close()
	if rc, ok := g.r.(io.ReadCloser); ok {
		return rc.Close()
	}
	return err
}
