package auth

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// These DEPRECATED functions should not be used
// and will be removed in 2.0

// CheckCSRFToken DEPRECATED
// this function will be removed in 2.0
func CheckCSRFToken(token, b64 string) error {
	// First base64 decode the value
	encrypted := make([]byte, 256)
	_, err := base64.URLEncoding.Decode(encrypted, []byte(b64))
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(encrypted, []byte(token))
}

// CSRFToken DEPRECATED
// this function will be removed in 2.0
func CSRFToken(token string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(token), HashCost)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// EncryptPassword renamed and DEPRECATED
// this function will be removed in 2.0
func EncryptPassword(pass string) (string, error) {
	fmt.Printf("Please use HashPassword instead, auth.EncryptPassword is deprecated")
	return HashPassword(pass)
}
