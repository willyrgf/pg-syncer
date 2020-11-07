package helpers

import (
	"crypto/sha256"
	"fmt"
	"sort"
)

// ToSha256 struct to sha256 hash
func ToSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// ArraysIsEqual return true if both arrays have same content
func ArraysIsEqual(t1 []string, t2 []string) bool {
	a1 := make([]string, len(t1))
	a2 := make([]string, len(t2))
	copy(a1, t1)
	copy(a2, t2)
	sort.Strings(a1)
	sort.Strings(a2)

	if len(a1) != len(a2) {
		return false
	}
	for i, v := range a1 {
		if v != a2[i] {
			return false
		}
	}
	return true
}
