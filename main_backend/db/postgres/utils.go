package postgres

import (
	"crypto/sha256"
	"fmt"
)

func CalculateSha256(value string) string {
	sha256Calculate := sha256.Sum256([]byte(value))
	hash := fmt.Sprintf("%x", sha256Calculate)
	return hash
}
