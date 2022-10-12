package zcrypto

import (
	"crypto/sha1"
	"testing"
)

func TestHash(t *testing.T) {
	if h := Hash(sha1.New(), "ASD"); h != "5271593ca406362d7a2701e331408ab77d5b5b88" {
		t.Error(h)
	}

	if h := Hash(sha1.New(), "XXX"); h != "a9674b19f8c56f785c91a555d0a144522bb318e6" {
		t.Error(h)
	}
}

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
