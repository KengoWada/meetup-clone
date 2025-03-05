package testutils

import (
	"time"

	"github.com/KengoWada/meetup-clone/internal/auth"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type ClaimsDetails struct {
	Sub       int64
	Exp       int64
	IssuedAt  int64
	NotBefore int64
}

func GenerateTesAuthToken(authenticator auth.Authenticator, cfg config.AuthConfig, isValid bool, ID int64) (string, error) {
	claims := generateJWTClaims(cfg, isValid, ID)
	token, err := authenticator.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateJWTClaims(cfg config.AuthConfig, isValid bool, ID int64) jwt.MapClaims {
	var (
		exp = time.Hour * time.Duration(cfg.Exp)
		iat = time.Now()
		aud = time.Now()
	)

	if !isValid {
		exp = -exp
		iat = iat.Add(-exp)
		aud = aud.Add(-exp)
	}

	return jwt.MapClaims{
		"sub": ID,
		"exp": time.Now().Add(exp).Unix(),
		"iat": iat.Unix(),
		"nbf": aud.Unix(),
		"aud": cfg.Audience,
		"iss": cfg.Issuer,
	}
}
