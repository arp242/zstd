// Package zcrypto implements cryptographic helpers.
package zcrypto

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// HashFile creates a SHA-256 hash for a file, as "sha256-hash".
func HashFile(filename string) (string, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("zcrypto.HashFile: %w", err)
	}
	defer fp.Close()

	h := sha256.New()
	_, err = io.Copy(h, fp)
	if err != nil {
		return "", fmt.Errorf("zcrypto.HashFile: %w", err)
	}

	return fmt.Sprintf("sha256-%x", h.Sum(nil)), nil
}

// VerifyHash verifies a file with a hash from HashFile().
func VerifyHash(filename, hash string) (bool, error) {
	ver, err := HashFile(filename)
	if err != nil {
		return false, fmt.Errorf("zcrypto.VerifyHash: %w", err)
	}
	return ver == hash, nil
}
