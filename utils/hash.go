package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func ComparePasswords(hashedPassword, password string) bool {
	return HashPassword(password) == hashedPassword
}
