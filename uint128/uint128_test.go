package uint128

import (
	"testing"
)

func TestUint128(t *testing.T) {
	n := []byte{44, 25, 67, 129, 231, 4, 77, 72, 157, 135, 42, 180, 126, 162, 176, 131}

	i, err := New(n)
	if err != nil {
		t.Fatal(err)
	}

	b, err := i.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	if string(n) != string(b) {
		t.Errorf("Bytes()")
	}

	if i.Format(16) != "2c194381e7044d489d872ab47ea2b083" {
		t.Errorf("Format(16)")
	}
	if i.Format(10) != "317764523729225658411351088340517695619" {
		t.Errorf("Format(10)")
	}

	parsed, err := Parse("2c194381e7044d489d872ab47ea2b083", 16)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Format(16) != "2c194381e7044d489d872ab47ea2b083" {
		t.Errorf("Format(16)")
	}

	parsed, err = Parse("317764523729225658411351088340517695619", 10)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Format(16) != "2c194381e7044d489d872ab47ea2b083" {
		t.Errorf("Format(10)")
	}

	//n2, err := New([]byte{})
}
