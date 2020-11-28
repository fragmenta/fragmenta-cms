# Package Auth
Package auth provides helpers for encryption, hashing and encoding.

### Setup

Setup the package on startup

```Go 
  auth.HMACKey = auth.HexToBytes("myhmac_key_from_config")
  auth.SecretKey = auth.HexToBytes("my_secret_key_from_config")
  auth.SessionName = "my_cookie_name"
  auth.SecureCookies = true
```


### Hashed Passwords

Use auth.HashPassword to encrypt and auth.CheckPassword to check hashed passwords (with bcrypt)

```Go 
  user.HashedPassword, err = auth.HashPassword(params.Get("password")
  if err != nil {
    return err
  }
  err = auth.CheckPassword(params.Get("password"), user.HashedPassword)
```

### Encrypted Sessions

Use auth.Session to set and get values from cookies, encrypted with AES GCM. 

```Go 
  // Build the session from the secure cookie, or create a new one
  session, err := auth.Session(writer, request)
  if err != nil {
    return err
  }
  
  // Store something in the session
  session.Set("my_key","my_value")
  session.Save(writer)
```


### Random Tokens

Generate and compare random tokens in constant time using the crypto/rand and crypto/subtle packages. 

```Go 
// Generate a new token
token := auth.RandomToken(32)

// Check tokens
if auth.CheckRandomToken(tok1,tok2) {
  // Tokens match
}
```

## Authorisation

You can use auth/can (separately) to authorise access to resources. 

To authorise actions:

```Go 
// Add an authorisation for admins to manage the pages resource
can.Authorise(role.Admin, can.ManageResource, "pages")
```

To check authorisation in handlers:

```Go 
// Check whether resource (conforming to can.Resource)
// can be managed by user (conforming to can.User) 
can.Manage(resource,user)
```


```Go 
// Interfaces for Users and Resources

// User defines the interface for users which must have numeric roles
type User interface {
	RoleID() int64 // for role check
	UserID() int64 // for ownership check
}

// Resource defines the interface for resources
type Resource interface {
	OwnedBy(int64) bool // for ownership check, passed a UserID
	ResourceID() string // for check against abilities registered on this resource
}
```
