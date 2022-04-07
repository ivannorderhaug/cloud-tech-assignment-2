package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

var Secret []byte

// Hash */
func Hash(id string) string {

	b := []byte(id)
	hash := hmac.New(sha256.New, Secret)
	hash.Write(b)

	return hex.EncodeToString(hash.Sum(nil))
}
