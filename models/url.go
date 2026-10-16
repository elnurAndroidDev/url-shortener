package models

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateShortURL() string {
	randomBytes := make([]byte, 6)
	_, err := rand.Read(randomBytes)

	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(randomBytes)[:6]
}