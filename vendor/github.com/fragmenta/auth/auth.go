// Package auth provides helpers for encryption, hashing and encoding.
package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// TODO: Add rotating cyphers on login (move to scrypt instead of bcrypt)

// HashCost sets the cost of bcrypt hashes
// - if this changes hashed passwords would need to be recalculated.
const HashCost = 10

// TokenLength sets the length of random tokens used for authenticity tokens.
const TokenLength = 32

// CheckPassword compares a password hashed with bcrypt.
func CheckPassword(pass, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}

// HashPassword hashes a password with a random salt using bcrypt.
func HashPassword(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), HashCost)
	return string(hash), err
}

// AuthenticityToken returns a new token for a request,
// and if necessary sets the cookie with our secret.
func AuthenticityToken(writer http.ResponseWriter, request *http.Request) (string, error) {
	// Fetch the session store
	session, err := Session(writer, request)
	if err != nil {
		return "", err
	}
	// Get the secret from the session, or generate if none found
	secret := session.Get(SessionTokenKey)
	if secret == "" {
		secret = BytesToBase64(RandomToken(TokenLength))
		session.Set(SessionTokenKey, secret)
		session.Save(writer)
	}

	// Now from secret, generate a secure token for this request
	token := AuthenticityTokenWithSecret(Base64ToBytes(secret))
	return BytesToBase64(token), nil
}

// CheckAuthenticityToken checks the token against that
// stored in a session cookie, and returns an error if the check fails.
func CheckAuthenticityToken(token string, request *http.Request) error {

	// Fetch the session store
	session, err := SessionGet(request)
	if err != nil {
		return err
	}

	// Get the secret from the session
	secret := session.Get(SessionTokenKey)
	if secret == "" {
		return fmt.Errorf("auth: error fetching authenticity secret from session")
	}

	return CheckAuthenticityTokenWithSecret(Base64ToBytes(token), Base64ToBytes(secret))
}

// CheckAuthenticityTokenWithSecret checks
// an auth token against a secret.
func CheckAuthenticityTokenWithSecret(token, secret []byte) error {

	// Check token length
	if len(token) != TokenLength*2 {
		return fmt.Errorf("auth: error failed - invalid token length %d", len(token))
	}

	// Grab random byte prefix, xor suffix secret against it to get our secret out,
	// and compare result to secret stored in cookie
	s := safeXORBytes(token[TokenLength:], token[:TokenLength])
	if CheckRandomToken(s, secret) {
		return nil
	}

	// If we reach here, CheckRandomToken failed
	return fmt.Errorf("auth: error failed with token")
}

// AuthenticityTokenWithSecret generates a new authenticity token
// from the secret by xoring a new random token with it
// and prepending the random bytes
// See https://github.com/rails/rails/pull/16570
// or gorilla/csrf for justification.
func AuthenticityTokenWithSecret(secret []byte) []byte {
	random := RandomToken(TokenLength)
	return append(random, safeXORBytes(random, secret)...)
}

// safeXORBytes is from https://golang.org/src/crypto/cipher/xor.go.
func safeXORBytes(a, b []byte) []byte {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst
}
