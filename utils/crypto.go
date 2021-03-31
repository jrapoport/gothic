package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// SecureToken creates a new random token
func SecureToken() string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err.Error()) // rand should never fail
	}
	return removePadding(base64.URLEncoding.EncodeToString(b))
}

// HashPassword generates a hashed password from a plaintext string
func HashPassword(password string) []byte {
	// we can safely ignore any error because we control the cost
	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return pw
}

// CheckPassword checks to see if the password matches the hashed password.
func CheckPassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}

func removePadding(token string) string {
	return strings.TrimRight(token, "=")
}
