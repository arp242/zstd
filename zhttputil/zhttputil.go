// Package zhttputil provides HTTP utility functions.
package zhttputil

// TODO: this should really be named "zhttp", but the reason it's not is because
// I have a zgo.at/zhttp package alread 😅 I need to rename that one.

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"zgo.at/zstd/zio"
	"zgo.at/zstd/znet"
)

// SafeClient returns a HTTP client that is only allowed to connect on the given
// ports to non-private addresses.
//
// This also sets the Timeout to 30 seconds, instead of the no timeout.
//
// See: https://www.agwa.name/blog/post/preventing_server_side_request_forgery_in_golang
//
// Also see SafeTransport() and znet.SafeDialer().
func SafeClient() *http.Client {
	return &http.Client{
		Transport: SafeTransport(nil),

		// Set a timeout too; this should be enough for most purposes.
		Timeout: 30 * time.Second,
	}

}

// SafeTransport returns a HTTP transport that is only allowed to connect on the
// given ports to non-private addresses.
//
// Port 80 and 443 are used if the list is empty.
//
// Also see SafeClient() and znet.SafeDialer().
func SafeTransport(ports []int) *http.Transport {
	if len(ports) == 0 {
		ports = []int{80, 443}
	}
	return &http.Transport{
		DialContext: znet.SafeDialer([]string{"tcp4", "tcp6"}, ports).DialContext,

		// Defaults from net/http.DefaultTransport
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// DumpBody reads the body of a HTTP request without consuming it so it can be
// read again later.
//
// It will read at most maxSize of bytes. Use -1 to read everything.
//
// It's based on httputil.DumpRequest.
//
// Copyright 2009 The Go Authors. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file:
// https://golang.org/LICENSE
func DumpBody(r *http.Request, maxSize int64) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	save, body, err := zio.DumpReader(r.Body)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	var dest io.Writer = &b

	chunked := len(r.TransferEncoding) > 0 && r.TransferEncoding[0] == "chunked"
	if chunked {
		dest = httputil.NewChunkedWriter(dest)
	}

	if maxSize < 0 {
		_, err = io.Copy(dest, body)
	} else {
		_, err = io.CopyN(dest, body, maxSize)
		if err == io.EOF {
			err = nil
		}
	}
	if err != nil {
		return nil, err
	}
	if chunked {
		_ = dest.(io.Closer).Close()
		_, _ = io.WriteString(&b, "\r\n")
	}

	r.Body = save
	return b.Bytes(), nil
}

// ErrNotOK is used when the status code is not 200 OK.
type ErrNotOK struct {
	URL string
	Err string
}

func (e ErrNotOK) Error() string {
	return fmt.Sprintf("code %v while downloading %v", e.Err, e.URL)
}

// Fetch the contents of an HTTP URL.
//
// This is not intended to cover all possible use cases  for fetching files,
// only the most common ones. Use the net/http package for more advanced usage.
func Fetch(url string) ([]byte, error) {
	client := http.Client{Timeout: 60 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot download %q: %w", url, err)
	}
	defer response.Body.Close() // nolint: errcheck

	// TODO: Maybe add sanity check to bail out of the Content-Length is very
	// large?
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body of %q: %w", url, err)
	}

	if response.StatusCode != http.StatusOK {
		return data, ErrNotOK{
			URL: url,
			Err: fmt.Sprintf("%v %v", response.StatusCode, response.Status),
		}
	}

	return data, nil
}

// Save an HTTP URL to the directory dir with the filename.
//
// The filename can be generated from the URL if empty.
//
// It will return the full path to the save file. Note that it may create both a
// file *and* return an error (e.g. in cases of non-200 status codes).
//
// This is not intended to cover all possible use cases  for fetching files,
// only the most common ones. Use the net/http package for more advanced usage.
func Save(url string, dir string, filename string) (string, error) {
	// Use last path of url if filename is empty
	if filename == "" {
		tokens := strings.Split(url, "/")
		filename = tokens[len(tokens)-1]
	}
	path := filepath.FromSlash(dir + "/" + filename)

	client := http.Client{Timeout: 60 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("cannot download %q: %w", url, err)
	}
	defer response.Body.Close() // nolint: errcheck

	output, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("cannot create %q: %w", path, err)
	}
	defer output.Close() // nolint: errcheck

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return path, fmt.Errorf("cannot read body of %q in to %q: %w", url, path, err)
	}

	if response.StatusCode != http.StatusOK {
		return path, ErrNotOK{
			URL: url,
			Err: fmt.Sprintf("%v %v", response.StatusCode, response.Status),
		}
	}

	return path, nil
}

// NopWriter is a http.ResponseWriter that doesn't write any body content.
type NopWriter struct{ http.ResponseWriter }

// Write is a no-op.
func (nop *NopWriter) Write(in []byte) (int, error) { return len(in), nil }
