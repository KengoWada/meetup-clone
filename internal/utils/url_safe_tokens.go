package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"time"
)

var ErrExpiredToken = errors.New("token has expired")

type TimedTokenData struct {
	CreatedAt string `json:"createdAt"`
	Body      string `json:"body"`
}

func GenerateToken(data string, key []byte) (string, error) {
	gcm, err := getBlockCipher(key)
	if err != nil {
		return "", err
	}

	timedTokenData := TimedTokenData{
		CreatedAt: time.Now().Format(time.RFC3339),
		Body:      data,
	}
	timedTokenDataStr, err := json.Marshal(timedTokenData)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	cipheredText := gcm.Seal(nonce, nonce, timedTokenDataStr, nil)
	return base64.RawURLEncoding.EncodeToString(cipheredText), nil
}

func ValidateToken(token string, key []byte, expiresIn time.Duration) (*TimedTokenData, error) {
	gcm, err := getBlockCipher(key)
	if err != nil {
		return nil, err
	}

	ciphered, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	originalText, err := gcm.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return nil, err
	}

	timedTokenData := TimedTokenData{}
	err = json.Unmarshal([]byte(originalText), &timedTokenData)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now()
	createdAt, err := time.Parse(time.RFC3339, timedTokenData.CreatedAt)
	if err != nil {
		return nil, err
	}

	if timeNow.Sub(createdAt) > expiresIn {
		return nil, ErrExpiredToken
	}

	return &timedTokenData, nil
}

func getBlockCipher(key []byte) (cipher.AEAD, error) {
	sha256Hash := sha256.Sum256(key)
	keyHash := hex.EncodeToString(sha256Hash[:])

	aesBlock, err := aes.NewCipher([]byte(keyHash))
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(aesBlock)
}
