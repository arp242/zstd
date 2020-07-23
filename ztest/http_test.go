package ztest

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCode(t *testing.T) {
	Code(t, &httptest.ResponseRecorder{Code: 200}, 200)

	// TODO: how to test that t.Fatalf() was called?
	//TestCode(t, &httptest.ResponseRecorder{Code: 200}, 201)
}

func TestMultipartForm(t *testing.T) {
	tests := []struct {
		inParams, inFiles map[string]string
		wantCt, wantErr   string
		wantBody          []string
	}{
		{
			map[string]string{
				"key":  "value",
				"w00t": "woot",
			},
			nil,
			"multipart/form-data",
			"",
			[]string{
				NormalizeIndent(`
					--::BOUNDARY::
					Content-Disposition: form-data; name="w00t"

					woot
				`),
				NormalizeIndent(`
					--::BOUNDARY::
					Content-Disposition: form-data; name="key"

					value
				`),
			},
		},

		{
			nil,
			map[string]string{
				"key":  "value",
				"w00t": "woot",
			},
			"multipart/form-data",
			"",
			[]string{
				NormalizeIndent(`
					--::BOUNDARY::
					Content-Disposition: form-data; name="key"; filename="key"
					Content-Type: application/octet-stream

					value
				`),
				NormalizeIndent(`
					--::BOUNDARY::
					Content-Disposition: form-data; name="w00t"; filename="w00t"
					Content-Type: application/octet-stream

					woot
				`),
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			outB, outCt, err := MultipartForm(tt.inParams, tt.inFiles)
			if !ErrorContains(err, tt.wantErr) {
				t.Fatal(err)
			}
			if err != nil {
				return
			}

			bi := strings.Index(outCt, "boundary=")
			tt.wantCt = tt.wantCt + "; " + outCt[bi:]
			if outCt != tt.wantCt {
				t.Errorf("wrong Content-Type\nout:  %#v\nwant: %#v\n", outCt, tt.wantCt)
			}

			// Can't compare body directly as output order isn't guaranteed.
			for _, wantBody := range tt.wantBody {
				wantBody = strings.Replace(wantBody, "::BOUNDARY::", outCt[bi+9:], -1)
				wantBody = strings.Replace(wantBody, "\n", "\r\n", -1) + "\r\n"

				b := outB.String()
				if !strings.Contains(b, wantBody) {
					t.Errorf("wrong body\nOUT:\n%q\nWANT:\n%q\n", b, tt.wantBody)
				}
			}
		})
	}
}
