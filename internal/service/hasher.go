package service

import (
	"crypto/sha256"
	"encoding/hex"
)

// SimpleHasher implements PasswordHasher interface.
type SimpleHasher struct {
	salt string
}

func NewSimpleHasher(salt string) *SimpleHasher {
	return &SimpleHasher{salt: salt}
}

// Hash hashes the password with SHA256 and a random salt.
// Returns hash:salt
func (h *SimpleHasher) Hash(password string) string {
	hash := sha256.Sum256([]byte(password + h.salt))
	return hex.EncodeToString(hash[:])
}
