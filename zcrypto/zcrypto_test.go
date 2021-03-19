package zcrypto

import (
	"testing"
)

func TestHashFile(t *testing.T) {
	f := "../zio/testdata/file1"

	hash, err := HashFile(f)
	if err != nil {
		t.Fatal(err)
	}

	want := "sha256-66a045b452102c59d840ec097d59d9467e13a3f34f6494e539ffd32c1bb35f18"
	if hash != want {
		t.Fatalf("wrong hash: %s", hash)
	}

	ok, err := VerifyHash(f, hash)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("VerifyHash not ok")
	}

	ok, err = VerifyHash(f, hash+"1")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("VerifyHash ok with wrong hash")
	}
}
