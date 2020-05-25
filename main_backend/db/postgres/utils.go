package postgres

import (
	"crypto/sha256"
	"fmt"
)

func CalculateSha256(values string) string {
	sha256Calculate := sha256.Sum256([]byte(values))
	hash := fmt.Sprintf("%x", sha256Calculate)
	return hash
}
