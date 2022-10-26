// Package zio implements some I/O utility functions.
package zio

import (
	"bytes"
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
