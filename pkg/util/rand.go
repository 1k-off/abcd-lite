package util

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random token")
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length]
}
