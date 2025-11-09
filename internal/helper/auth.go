package helper

import (
	"crypto/sha256"
	"math/rand"
)

const randomCodeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomCode(length int) string {
	// TODO better
	randomBytes := make([]byte, length)
	for i := range randomBytes {
		randomBytes[i] = randomCodeCharset[rand.Intn(len(randomCodeCharset))]
	}

	return string(randomBytes)
}

func Hash(target string) (string, error) {
	// TODO better

	hashed := sha256.Sum256([]byte(target))

	return string(hashed[:]), nil
}
