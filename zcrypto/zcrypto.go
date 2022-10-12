// Package zcrypto implements cryptographic helpers.
package zcrypto

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// Hash returns the base-16 encoded string of value hashed with h.
func Hash(h hash.Hash, value string) string {
	if value == "" {
		return ""
	}
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

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
	m, _, ok := strings.Cut(hash, "-")
	if !ok {
		return false, errors.New("zcrypto.VerifyHash: no hash algorithm in hash string")
	}
	if m != "sha256" {
		return false, fmt.Errorf("zcrypto.VerifyHash: unknown hash algorithm: %q", m)
	}

	ver, err := HashFile(filename)
	if err != nil {
		return false, fmt.Errorf("zcrypto.VerifyHash: %w", err)
	}
	return subtle.ConstantTimeCompare([]byte(ver), []byte(hash)) == 1, nil
}

// Secret number of 256 bits formatted in base36 (~48 bytes).
func Secret256() string { return secret(4) }

// Secret number of 192 bits formatted in base36 (~38 bytes).
func Secret192() string { return secret(3) }

// Secret number of 128 bits formatted in base36 (~25 bytes).
func Secret128() string { return secret(2) }

// Secret number of 64 bits formatted in base36 (~12 bytes).
func Secret64() string { return secret(1) }

const alphabet = "23456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

// Secret string of size  characters selected from sel.
//
// If sel is "" it will use [0-9a-zA-Z], excluding some easily confused
// characters: 0, 1, i, l, o, I, and O.
func SecretString(size int, sel string) string {
	if sel == "" {
		sel = alphabet
	}

	var (
		b = make([]byte, size)
		m = big.NewInt(int64(len(sel)))
	)
	for i := range b {
		n, err := rand.Int(rand.Reader, m)
		if err != nil {
			panic(fmt.Errorf("zcrypto.Secret: %w", err))
		}
		b[i] = sel[n.Int64()]
	}
	return string(b)
}

var max = big.NewInt(0).SetUint64(1e19)

func secret(n int) string {
	var key strings.Builder
	for i := 0; i < n; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(fmt.Errorf("zcrypto.Secret: %w", err))
		}
		_, _ = key.WriteString(strconv.FormatUint(n.Uint64(), 36))
	}
	return key.String()
}
