//go:build testhttp

package zhttputil

import (
	"fmt"
	"strings"
	"testing"

	"zgo.at/zstd/ztest"
)

// TODO: better to not depend on interwebz; use httptest.Server
func TestFetch(t *testing.T) {
	cases := []struct {
		in, want, wantErr string
	}{
		{"http://example.com", "<html>", ""},
		{"http://fairly-certain-this-doesnt-exist-asdasd12g1ghdfddd.com", "", "cannot download"},
		{"http://httpbin.org/status/400", "", "400"},
		{"http://httpbin.org/status/500", "", "500"},
		// Make sure we return the body as well.
		{"http://httpbin.org/status/418", "teapot", "418"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := Fetch(tc.in)
			if !ztest.ErrorContains(err, tc.wantErr) {
				t.Errorf("wrong error\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !strings.Contains(string(out), tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", string(out), tc.want)
			}
		})
	}
}
