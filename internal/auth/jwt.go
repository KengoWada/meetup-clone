package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthenticator handles JWT-based authentication.
// It manages the generation and validation of JWT tokens using a secret, audience, and issuer.
type JWTAuthenticator struct {
	secret string // secret key used to sign the JWT
	aud    string // audience claim for the token
	iss    string // issuer claim for the token
}

// NewJWTAuthenticator creates a new JWTAuthenticator instance.
// It takes a secret key, audience, and issuer as parameters and returns a pointer to JWTAuthenticator.
//
// Parameters:
//   - secret: the secret key used to sign JWT tokens
//   - aud: the audience claim for the tokens
//   - iss: the issuer claim for the tokens
//
// Returns:
//   - *JWTAuthenticator: a pointer to the initialized JWTAuthenticator instance
func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{secret, iss, aud}
}

// GenerateToken creates a signed JWT token with the given claims using the HS256 signing method.
// It uses the configured secret, audience, and issuer to build the token.
//
// The token is signed with HMAC using SHA-256 (`jwt.SigningMethodHS256`).
//
// Parameters:
//   - claims: the claims to embed in the token (e.g., user ID, expiration time)
//
// Returns:
//   - string: the signed JWT token as a string
//   - error: an error if token generation fails
func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken verifies the given JWT token.
// It checks the signature (using HS256), audience, and issuer claims to ensure the token is valid.
//
// Parameters:
//   - token: the JWT token string to validate
//
// Returns:
//   - *jwt.Token: the parsed and validated JWT token
//   - error: an error if the token is invalid or verification fails
func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
