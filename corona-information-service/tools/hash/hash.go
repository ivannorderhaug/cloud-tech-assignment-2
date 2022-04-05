package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

var Secret = []byte{0, 4, 0, 2, 2, 0, 0, 0}

// Hash */
func Hash(id string) string {

	b := []byte(id)
	hash := hmac.New(sha256.New, Secret)
	hash.Write(b)

	return hex.EncodeToString(hash.Sum(nil))
}
