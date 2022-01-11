package aternos_api

import (
	"encoding/base64"
	"math/rand"
)

// randomString generates a random lowercase string.
// E.g. mdlc2c9chx9, mbywjir33mm
func randomString(length int) string {
	charset := []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	s := make([]rune, length)

	for i := range s {
		s[i] = charset[rand.Intn(len(charset))]
	}

	return string(s)
}

func atob(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
