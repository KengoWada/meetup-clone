// Package auth provides tools for handling authentication within the API.
// It includes implementations for generating and validating JWT tokens,
// as well as the Authenticator interface for flexible authentication strategies.
package auth

import "github.com/golang-jwt/jwt/v5"

// Authenticator defines the interface for handling JWT authentication.
// It provides methods to generate and validate JWT tokens.
type Authenticator interface {
	// GenerateToken generates a signed JWT token with the provided claims.
	// It returns the generated token as a string or an error if token creation fails.
	GenerateToken(claims jwt.Claims) (string, error)

	// ValidateToken verifies the provided JWT token and returns the parsed token.
	// It returns an error if the token is invalid or expired.
	ValidateToken(token string) (*jwt.Token, error)
}
