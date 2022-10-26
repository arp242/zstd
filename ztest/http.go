package ztest

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Code checks if the error code in the recoder matches the desired one, and
// will stop the test with t.Fatal() if it doesn't.
func Code(t *testing.T, recorder *httptest.ResponseRecorder, want int) {
	t.Helper()
	if recorder.Code != want {
		t.Errorf("wrong response code\nhave: %d %s\nwant: %d %s\nbody: %v",
			recorder.Code, http.StatusText(recorder.Code),
			want, http.StatusText(want),
			elideLeft(recorder.Body.String(), 500))
	}
}

// Default values for NewRequest()
var (
	DefaultHost        = "example.com"
	DefaultContentType = "application/json"
)

// NewRequest creates a new request with some sensible defaults set.
func NewRequest(method, target string, body io.Reader) *http.Request {
	r := httptest.NewRequest(method, target, body)
	if r.Host == "" || r.Host == "example.com" {
		r.Host = DefaultHost
	}
	if r.Header.Get("Content-Type") == "" {
		r.Header.Set("Content-Type", DefaultContentType)
	}
	return r
}

// Body returns the JSON representation as an io.Reader. This is useful for
// creating a request body. For example:
//
//	NewRequest("POST", "/", ztest.Body(someStruct{
//	    Foo: "bar",
//	}))
func Body(a any) *bytes.Reader {
	j, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(j)
}

// HTTP sets up a HTTP test. A GET request will be made if r is nil.
//
// For example:
//
//	rr := ztest.HTTP(t, nil, MyHandler)
//
// Or for a POST request:
//
//	r, err := zhttp.NewRequest("POST", "/v1/email", nil)
//	if err != nil {
//	    t.Fatal(err)
//	}
//	rr := ztest.HTTP(t, r, MyHandler)
func HTTP(t *testing.T, r *http.Request, h http.Handler) *httptest.ResponseRecorder {
	t.Helper()

	rr := httptest.NewRecorder()
	if r == nil {
		var err error
		r, err = http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatalf("cannot make request: %v", err)
		}
	}

	h.ServeHTTP(rr, r)
	return rr
}

// MultipartForm writes the keys and values from params to a multipart form.
//
// The first input parameter is used for "multipart/form-data" key/value
// strings, the optional second parameter is used creating file parts.
//
// Don't forget to set the Content-Type from the return value:
//
//	r.Header.Set("Content-Type", contentType)
func MultipartForm(params ...map[string]string) (b *bytes.Buffer, contentType string, err error) {
	b = &bytes.Buffer{}
	w := multipart.NewWriter(b)

	for k, v := range params[0] {
		field, err := w.CreateFormField(k)
		if err != nil {
			return nil, "", err
		}
		_, err = field.Write([]byte(v))
		if err != nil {
			return nil, "", err
		}
	}

	if len(params) > 1 {
		for k, v := range params[1] {
			field, err := w.CreateFormFile(k, k)
			if err != nil {
				return nil, "", err
			}
			_, err = field.Write([]byte(v))
			if err != nil {
				return nil, "", err
			}
		}
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return b, w.FormDataContentType(), nil
}

func elideLeft(s string, n int) string {
	ss := sub(s, 0, n)
	if len(s) != len(ss) {
		return ss + "…"
	}
	return s
}

func sub(s string, start, end int) string {
	var (
		nchars    int
		startbyte = -1
	)
	for bytei := range s {
		if nchars == start {
			startbyte = bytei
		}
		if nchars == end {
			return s[startbyte:bytei]
		}
		nchars++
	}
	if startbyte == -1 {
		return ""
	}
	return s[startbyte:]
}
