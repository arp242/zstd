package zhttputil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSafeClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	c := SafeClient()
	resp, err := c.Get(srv.URL)
	if err == nil {
		t.Fatal("err is nil")
	}
	if resp != nil {
		t.Fatal("resp not nil", resp)
	}
}

// Based on httputil.DumpRequest
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func TestDumpBody(t *testing.T) {
	chunk := func(s string) string { return fmt.Sprintf("%x\r\n%s\r\n", len(s), s) }

	tests := []struct {
		Req  http.Request
		Body interface{} // optional []byte or func() io.ReadCloser to populate Req.Body

		WantDump string
		ReadN    int64
		NoBody   bool // if true, set DumpRequest{,Out} body to false
	}{

		// HTTP/1.1 => chunked coding; body; empty trailer
		{
			Req: http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
					Path:   "/search",
				},
				ProtoMajor:       1,
				ProtoMinor:       1,
				TransferEncoding: []string{"chunked"},
			},
			Body:     []byte("abcdef"),
			WantDump: chunk("abcdef") + chunk(""),
			ReadN:    -1,
		},
		{
			Req: http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
					Path:   "/search",
				},
				ProtoMajor:       1,
				ProtoMinor:       1,
				TransferEncoding: []string{"chunked"},
			},
			Body:     []byte("abcdef"),
			WantDump: chunk("a") + chunk(""),
			ReadN:    1,
		},
		{
			Req: http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
					Path:   "/search",
				},
				ProtoMajor:       1,
				ProtoMinor:       1,
				TransferEncoding: []string{"chunked"},
			},
			Body:     []byte("abcdef"),
			WantDump: chunk("ab") + chunk(""),
			ReadN:    2,
		},

		// Request with Body > 8196 (default buffer size)
		{
			Req: http.Request{
				Method: "POST",
				URL: &url.URL{
					Scheme: "http",
					Host:   "post.tld",
					Path:   "/",
				},
				Header: http.Header{
					"Content-Length": []string{"8193"},
				},
				ContentLength: 8193,
				ProtoMajor:    1,
				ProtoMinor:    1,
			},
			Body:     bytes.Repeat([]byte("a"), 8193),
			WantDump: strings.Repeat("a", 8193),
			ReadN:    -1,
		},
		{
			Req: http.Request{
				Method: "POST",
				URL: &url.URL{
					Scheme: "http",
					Host:   "post.tld",
					Path:   "/",
				},
				Header: http.Header{
					"Content-Length": []string{"8193"},
				},
				ContentLength: 8193,
				ProtoMajor:    1,
				ProtoMinor:    1,
			},
			Body:     bytes.Repeat([]byte("a"), 8293),
			WantDump: strings.Repeat("a", 8193),
			ReadN:    8193,
		},
		{
			Req: http.Request{
				Method: "POST",
				URL: &url.URL{
					Scheme: "http",
					Host:   "post.tld",
					Path:   "/",
				},
				Header: http.Header{
					"Content-Length": []string{"10"},
				},
				ContentLength: 10,
				ProtoMajor:    1,
				ProtoMinor:    1,
			},
			Body:     bytes.Repeat([]byte("a"), 10),
			WantDump: strings.Repeat("a", 10),
			ReadN:    15,
		},
		{
			Req: http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme: "http",
					Host:   "www.google.com",
					Path:   "/search",
				},
				ProtoMajor:       1,
				ProtoMinor:       1,
				TransferEncoding: []string{"chunked"},
			},
			Body:     []byte("abcdef"),
			WantDump: chunk("abcdef") + chunk(""),
			ReadN:    100,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			setBody := func() {
				if tt.Body == nil {
					return
				}
				switch b := tt.Body.(type) {
				case []byte:
					tt.Req.Body = io.NopCloser(bytes.NewReader(b))
				case func() io.ReadCloser:
					tt.Req.Body = b()
				default:
					t.Fatalf("unsupported Body of %T", tt.Body)
				}
			}

			setBody()
			if tt.Req.Header == nil {
				tt.Req.Header = make(http.Header)
			}

			setBody()
			dump, err := DumpBody(&tt.Req, tt.ReadN)
			if err != nil {
				t.Fatal(err)
			}
			if string(dump) != tt.WantDump {
				t.Errorf("\nwant: %#v\ngot:  %#v\n", tt.WantDump, string(dump))
			}
		})
	}
}
