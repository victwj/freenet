// Functions related to generating/manipulating cryptographic keys
package main

import (
	"crypto/ed25519"
	"crypto/sha1"
	"fmt"
	"hash/fnv"
	"math/rand"
)

// Generates 40 hex characters
// I.e. 40 * 4 = 160 bits / 8 = 20 bytes
func hashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func hashBytes(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Use FNV to generate int from string
func hashFnv(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}

// Using edDSA since RSA Go cannot be made deterministic
func genEdKey(s string) (ed25519.PublicKey, ed25519.PrivateKey) {
	seed := hashFnv(s)
	r := rand.New(rand.NewSource(seed))
	pub, priv, err := ed25519.GenerateKey(r)
	if err != nil {
		panic(err)
	}
	return pub, priv
}

// Generate freenet keyword signed key
func genKeywordSignedKey(descr string) (ed25519.PublicKey, ed25519.PrivateKey, string) {
	pub, priv := genEdKey(descr)
	ksk := hashBytes(pub)
	return pub, priv, ksk[:10]
}
