package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func IsTokenExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}
