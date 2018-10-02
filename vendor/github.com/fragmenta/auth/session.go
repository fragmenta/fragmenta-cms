package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// This secure cookie code is based on Gorilla secure cookie
// but with mandatory AES-GCM encryption.

// MaxAge is the age in seconds of a cookie before it expires, default 60 days.
var MaxAge = 86400 * 60

// MaxCookieSize is the maximum length of a cookie in bytes, defaults to 4096.
var MaxCookieSize = 4096

// HMACKey is a 32 byte key for generating HMAC distinct from SecretKey.
var HMACKey []byte

// SecretKey is a 32 byte key for encrypting content with AES-GCM.
var SecretKey []byte

// SessionName is the name of the ssions.
var SessionName = "fragmenta_session"

// SessionUserKey is the session user key.
var SessionUserKey = "user_id"

// SessionTokenKey is the session token key.
var SessionTokenKey = "authenticity_token"

// SecureCookies is true if we use secure https cookies.
var SecureCookies = false

// SessionStore is the interface for a session store.
type SessionStore interface {
	Get(string) string
	Set(string, string)
	Load(request *http.Request) error
	Save(http.ResponseWriter) error
	Clear(http.ResponseWriter)
}

// CookieSessionStore is a concrete version of SessionStore,
// which stores the information encrypted in cookies.
type CookieSessionStore struct {
	values map[string]string
}

// Session loads the current sesions or returns a new blank session.
func Session(writer http.ResponseWriter, request *http.Request) (SessionStore, error) {

	s, err := SessionGet(request)
	if err != nil {
		return s, nil
	}

	return s, nil
}

// SessionGet loads the current session (if any)
func SessionGet(request *http.Request) (SessionStore, error) {

	// Return the current session store from cookie or a new one if none found
	s := &CookieSessionStore{
		values: make(map[string]string, 0),
	}

	if len(HMACKey) == 0 || len(SecretKey) == 0 || len(SessionTokenKey) == 0 {
		return s, errors.New("auth: secrets not initialised")
	}

	// Check if the session exists and load it
	err := s.Load(request)
	if err != nil {
		return s, fmt.Errorf("auth: error loading session: %s", err) // return blank session if none found
	}

	return s, nil
}

// ClearSession clears the current session cookie
func ClearSession(w http.ResponseWriter) {
	// First delete all Set-Cookie headers so we only have one
	w.Header().Del("Set-Cookie")

	cookie := &http.Cookie{
		Name:   SessionName,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}

	http.SetCookie(w, cookie)
}

// Get a value from the session.
func (s *CookieSessionStore) Get(key string) string {
	return s.values[key]
}

// Set a value in the session, this does not save to the cookie.
func (s *CookieSessionStore) Set(key string, value string) {
	s.values[key] = value
}

// Load the session from cookie.
func (s *CookieSessionStore) Load(request *http.Request) error {

	// Return if session name not defined
	if SessionName == "" {
		return fmt.Errorf("auth: error session_name not set")
	}

	cookie, err := request.Cookie(SessionName)
	if err != nil {
		return fmt.Errorf("auth: error getting cookie: %s", err)
	}

	// Read the encrypted values back out into our values in the session.
	err = s.Decode(SessionName, HMACKey, SecretKey, cookie.Value, &s.values)
	if err != nil {
		return fmt.Errorf("auth: error decoding session: %s", err)
	}

	return nil
}

// Save the session to a cookie.
func (s *CookieSessionStore) Save(writer http.ResponseWriter) error {

	// Return error if session name not defined
	if SessionName == "" {
		return fmt.Errorf("auth: error session_name not set")
	}

	encrypted, err := s.Encode(SessionName, s.values, HMACKey, SecretKey)
	if err != nil {
		return fmt.Errorf("auth: error encoding session: %s", err)
	}

	cookie := &http.Cookie{
		Name:     SessionName,
		Value:    encrypted,
		HttpOnly: true,
		Secure:   SecureCookies,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 7), // Expires in seven days
	}

	http.SetCookie(writer, cookie)

	return nil
}

// Clear the session values from the cookie.
func (s *CookieSessionStore) Clear(writer http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   SessionName,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}

	http.SetCookie(writer, cookie)
}

// Encode a given value in the session cookie.
func (s *CookieSessionStore) Encode(name string, value interface{}, hashKey []byte, secretKey []byte) (string, error) {

	if name == "" || hashKey == nil || secretKey == nil {
		return "", errors.New("auth: encode keys not set")
	}

	// Serialize
	b, err := serialize(value)
	if err != nil {
		return "", fmt.Errorf("auth: error serializing value: %s", err)
	}

	// Encrypt with AES/GCM
	b, err = Encrypt(b, secretKey)
	if err != nil {
		return "", fmt.Errorf("auth: error encrypting value: %s", err)
	}

	// Encode to base64
	b = encodeBase64(b)

	// Note Encrypt above also verifies now with GCM.
	// Create MAC for "name|date|value". Extra pipe unused.
	now := time.Now().UTC().Unix()
	b = []byte(fmt.Sprintf("%s|%d|%s|", name, now, b))
	mac := CreateMAC(hmac.New(sha256.New, hashKey), b[:len(b)-1])

	// Append mac, remove name
	b = append(b, mac...)[len(name)+1:]

	// Encode to base64 again
	b = encodeBase64(b)

	// Check length when encoded
	if MaxCookieSize != 0 && len(b) > MaxCookieSize {
		return "", fmt.Errorf("auth: error len over max cookie size: %d", MaxCookieSize)
	}

	// Done, convert to string and return
	return string(b), nil
}

// Decode the value in the session cookie.
func (s *CookieSessionStore) Decode(name string, hashKey []byte, secretKey []byte, value string, dst interface{}) error {

	if name == "" || hashKey == nil || secretKey == nil {
		return errors.New("auth: decode keys not set")
	}

	if MaxCookieSize != 0 && len(value) > MaxCookieSize {
		return errors.New("auth: cookie value is too long")
	}

	// Decode from base64
	b, err := decodeBase64([]byte(value))
	if err != nil {
		return fmt.Errorf("auth: error decoding base 64 value: %s", err)
	}

	// Verify MAC - value is "date|value|mac"
	parts := bytes.SplitN(b, []byte("|"), 3)
	if len(parts) != 3 {
		return errors.New("auth: MAC invalid")
	}
	h := hmac.New(sha256.New, hashKey)
	b = append([]byte(name+"|"), b[:len(b)-len(parts[2])-1]...)
	err = VerifyMAC(h, b, parts[2])
	if err != nil {
		return err
	}

	// Verify date ranges
	timestamp, err := strconv.ParseInt(string(parts[0]), 10, 64)
	if err != nil {
		return errors.New("auth: timestamp invalid")
	}
	now := time.Now().UTC().Unix()
	if MaxAge != 0 && timestamp < now-int64(MaxAge) {
		return errors.New("auth: timestamp expired")
	}

	// Decode from base64
	b, err = decodeBase64(parts[1])
	if err != nil {
		return fmt.Errorf("auth: error decoding value: %s", err)
	}

	// Derypt with AES
	b, err = Decrypt(b, secretKey)
	if err != nil {
		return fmt.Errorf("auth: error decrypting value: %s", err)
	}

	// Deserialize
	err = deserialize(b, dst)
	if err != nil {
		return fmt.Errorf("auth: error deserializing value: %s", err)
	}

	// Done.
	return nil
}
