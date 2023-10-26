package crypto

import (
	"crypto/sha512"
	"fmt"
)

// HashSHA512 Returns the SHA512 hash of a string.
func HashSHA512(s string) string {
	hash := sha512.Sum512([]byte(s))
	return fmt.Sprintf("%x", hash)
}
