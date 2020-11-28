package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"
)

// RandomToken generates a random token 32 bytes long,
// or at a specified length if arguments are provided.
func RandomToken(args ...int) []byte {
	length := 32
	if len(args) > 0 && args[0] != 0 {
		length = args[0]
	}
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error reading random token:", err)
		return nil
	}
	return b
}

// CheckRandomToken performs a comparison of two tokens
// resistant to timing attacks.
func CheckRandomToken(a, b []byte) bool {
	return (subtle.ConstantTimeCompare(a, b) == 1)
}
