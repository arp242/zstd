// Package zio implements some I/O utility functions.
package zio

import (
	"bytes"
	"context"
	"errors"
	"hash"
	"io"
	"net/http"
	"os"
)

// DumpReader reads all of b to memory and then returns two equivalent
// ReadClosers which will yield the same bytes.
//
// This is useful if you want to read data from an io.Reader more than once.
//
// It returns an error if the initial reading of all bytes fails. It does not
// attempt to make the returned ReadClosers have identical error-matching
// behavior.
//
// This is based on httputil.DumpRequest, see zio.DumpBody() for an example
// usage.
//
// Copyright 2009 The Go Authors. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file:
// https://golang.org/LICENSE
func DumpReader(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}

	if err = b.Close(); err != nil {
		return nil, b, err
	}

	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// Exists reports if a path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Newer reports if the file's mtime is more recent than base.
func ChangedFrom(file, base string) bool {
	// TODO: change the arguments to:
	//
	//    ChangedFrom(base, file string, files ...string)
	//
	// Makes it easier to check multiple files:
	//
	//  if !zio.ChangedFrom("./handlers/api.go", "./tpl/api.json") &&
	//      !zio.ChangedFrom("./kommentaar.conf", "./tpl/api.json") {
	filest, err := os.Stat(file)
	if err != nil {
		return true
	}

	basest, err := os.Stat(base)
	if err != nil {
		return true
	}

	return filest.ModTime().After(basest.ModTime())
}

type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }

// NopCloser returns a WriteCloser with a no-op Close method.
func NopCloser(r io.Writer) io.WriteCloser { return nopCloser{r} }

// NopWriter is an io.Writer that does nothing.
type NopWriter struct{}

// Write is a stub.
func (w *NopWriter) Write(b []byte) (int, error) { return len(b), nil }

// TeeReader returns a [Reader] that writes to w what it reads from r.
//
// All reads from r performed through it are matched with corresponding writes
// to w. There is no internal buffering - the write must complete before the
// read completes. Any error encountered while writing is reported as a read
// error.
//
// This is simular to [io.TeeReader], except that it supports multiple writers.
func TeeReader(r io.Reader, w ...io.Writer) io.Reader {
	return &teeReader{r, w}
}

type teeReader struct {
	r io.Reader
	w []io.Writer
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		for _, ww := range t.w {
			if n, err := ww.Write(p[:n]); err != nil {
				return n, err
			}
		}
	}
	return
}

// Count the number of occurrences of find in the stream.
//
// It will try to seek back to the original position after counting, so that
// something like this will just work:
//
//	fp, _ := os.Open("..")
//	lines, _ := zio.Count(context.Background(), fp, []byte{'\n'})
//
//	scan := bufio.NewScanner(fp)
//	var i int
//	for scan.Scan() {
//	  fmt.Printf("line %d / %d\n", i+1, lines)
//	  i++
//	}
func Count(ctx context.Context, fp io.ReadSeeker, find []byte) (int, error) {
	pos, err := fp.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	var (
		result = make(chan struct {
			cnt int
			err error
		})
	)
	go func() {
		var (
			buf = make([]byte, 1024*128)
			cnt int
		)
		for {
			n, err := fp.Read(buf)
			cnt += bytes.Count(buf[:n], find)
			if err != nil {
				if errors.Is(err, io.EOF) {
					err = nil
				}
				_, seekErr := fp.Seek(pos, io.SeekStart)
				if err == nil {
					err = seekErr
				}
				result <- struct {
					cnt int
					err error
				}{cnt, err}
				break
			}
		}
	}()

	select {
	case r := <-result:
		return r.cnt, r.err
	case <-ctx.Done():
		fp.Seek(pos, io.SeekStart)
		return 0, ctx.Err()
	}
}

// PeekReader returns a reader that first returns data from peeked, before
// reading from the reader r.
//
// This is useful in cases where you want to "peek" a bit data from a reader
// that doesn't support seeking to determine if the compression or file format.
func PeekReader(r io.Reader, peeked []byte) io.ReadCloser {
	return &peekReader{r, peeked}
}

type peekReader struct {
	r      io.Reader
	peeked []byte
}

func (r *peekReader) Read(d []byte) (int, error) {
	if len(r.peeked) == 0 {
		return r.r.Read(d)
	}

	n := copy(d, r.peeked)
	r.peeked = r.peeked[n:]
	if len(r.peeked) > 0 {
		return n, nil
	}
	r.peeked = nil

	n2, err := r.r.Read(d[n:])
	return n + n2, err
}

// Close the underlying reader if it implements a Close method.
func (r *peekReader) Close() error {
	if rc, ok := r.r.(io.ReadCloser); ok {
		return rc.Close()
	}
	return nil
}

// LimitReader returns a Reader that reads from r but stops with EOF after n
// bytes. The underlying implementation is a *LimitedReader.
//
// This is identical to [io.LimitReader], except that it accepts a ReadCloser
// (and supports Close).
func LimitReader(r io.ReadCloser, n int64) io.ReadCloser { return &LimitedReader{r, n} }

// A LimitedReader reads from R but limits the amount of data returned to just N
// bytes. Each call to Read updates N to reflect the new amount remaining.
//
// Read returns EOF when N <= 0 or when the underlying R returns EOF.
//
// This is identical to [io.LimitedReader], except that it accepts a ReadCloser
// (and supports Close).
type LimitedReader struct {
	R io.ReadCloser // underlying reader
	N int64         // max bytes remaining
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	return
}

func (l *LimitedReader) Close() error {
	return l.R.Close()
}

type hashWriter struct {
	w io.WriteCloser
	h hash.Hash
	l int
}

// HashWriter writes to both the writer and hash.
func HashWriter(w io.WriteCloser, h hash.Hash) *hashWriter {
	return &hashWriter{w, h, 0}
}

// Write to the underlying writer and hash.
func (w *hashWriter) Write(b []byte) (int, error) {
	w.h.Write(b)
	w.l += len(b)
	return w.w.Write(b)
}

// Close the underlying writer.
func (w *hashWriter) Close() error { return w.w.Close() }

// Hash sums the hash.
func (w *hashWriter) Hash() []byte { return w.h.Sum(nil) }

// Len gets the total number of bytes written thus far.
func (w *hashWriter) Len() int { return w.l }

type hashReader struct {
	r io.ReadCloser
	h hash.Hash
	l int
}

// HashReader writes to the hash for all data it reads.
func HashReader(r io.ReadCloser, h hash.Hash) *hashReader {
	return &hashReader{r, h, 0}
}

// Read the underlying reader and write to the hash.
func (r *hashReader) Read(b []byte) (int, error) {
	n, err := r.r.Read(b)
	r.h.Write(b[:n])
	r.l += n
	return n, err
}

// Close the underlying reader.
func (r *hashReader) Close() error { return r.r.Close() }

// Hash sums the hash.
func (r *hashReader) Hash() []byte { return r.h.Sum(nil) }

// Len gets the total number of bytes read thus far.
func (r *hashReader) Len() int { return r.l }
