package auth

import (
	"crypto/sha1"
	"encoding/hex"
)

// Create a sha1 hash from a string
func CreateHash(text string) string {
	h := sha1.New()

	h.Write([]byte(text))

	return hex.EncodeToString(h.Sum(nil))
}
