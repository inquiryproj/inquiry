package crypto

import (
	"crypto/rand"
	"fmt"
)

// NewAPIKey generates a new api key using crypto rand package.
func NewAPIKey() (string, error) {
	b := make([]byte, 15)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("inq_%x", b), nil
}
