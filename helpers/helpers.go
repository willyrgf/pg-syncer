package helpers

import (
	"crypto/sha256"
	"fmt"
)

// ToSha256 struct to sha256 hash
func ToSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
