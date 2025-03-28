package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/KengoWada/meetup-clone/internal"
)

var ErrExpiredToken = errors.New("token has expired")

// TimedTokenData represents the data structure used for storing the content
// and timestamp of a token. The struct contains information about when
// the token was created and the body of the data included in the token.
//
// Fields:
//   - CreatedAt: The timestamp representing when the token was created.
//     This is stored as a string, using time.RFC3339 format (e.g., "2006-01-02T15:04:05Z07:00").
//   - Body: The main content or data associated with the token. This can be any string value
//     that needs to be securely transmitted or encrypted within the token.
//
// Example:
//
//	{
//	  "createdAt": "2025-03-01T12:34:56Z",
//	  "body": "some important data"
//	}
type TimedTokenData struct {
	CreatedAt string `json:"createdAt"`
	Body      string `json:"body"`
}

// GenerateToken encrypts the provided data along with the current timestamp using
// the provided key, and then encodes the result to a Base64 URL-safe string.
//
// Parameters:
//   - data: The string data to be encrypted. This is the main content to be included in the token.
//   - key: The secret key used for encryption. This key is combined with the data to create the token.
//
// Returns:
//   - A Base64 URL-safe encoded string that represents the encrypted token,
//     which includes the data and the timestamp.
//   - An error if there is any issue during the encryption or encoding process.
//
// Note:
//   - The current timestamp is included in the token to ensure its validity or time-based expiry.
//   - Base64 URL encoding ensures that the token can be safely transmitted over URLs without
//     any special character issues.
func GenerateToken(data string, key []byte) (string, error) {
	return generateToken(data, key, time.Now().UTC().Format(internal.DateTimeFormat))
}

// GenerateTestToken creates an encrypted test token with the provided data and timestamp.
//
// This function is a wrapper around `generateToken`, useful for generating tokens in test scenarios.
// It encrypts the given data and timestamp using AES-GCM and returns a base64 URL-encoded token.
//
// Parameters:
//   - data (string): The data to include in the token.
//   - key ([]byte): The encryption key (must be 32 bytes long).
//   - createdAt (string): The timestamp to embed in the token (e.g., time.Now().Format(time.RFC3339)).
//
// Returns:
//   - string: The generated token as a base64 URL-encoded string.
//   - error: An error if token generation fails.
//
// Example usage:
//
//	key := []byte("your-32-byte-secret-key-----------------")
//	token, err := GenerateTestToken("testData", key, "2025-03-01T12:00:00Z")
//	if err != nil {
//	    log.Fatalf("Failed to generate test token: %v", err)
//	}
//	fmt.Println("Test Token:", token)
func GenerateTestToken(data string, key []byte, createdAt string) (string, error) {
	return generateToken(data, key, createdAt)
}

// ValidateToken validates a given token by decrypting it using the provided key and checks
// if it is still valid based on the provided expiration duration. It returns the decrypted
// token data and an error if any issues arise during decryption or validation.
//
// Parameters:
//   - token: The token string to be validated, which is expected to be in Base64 URL-safe format.
//   - key: The secret key used to decrypt the token and verify its authenticity.
//   - expiresIn: The duration after which the token is considered expired. This helps in
//     validating whether the token has expired based on the timestamp stored in the token.
//
// Returns:
//   - A pointer to a `TimedTokenData` struct containing the decrypted `createdAt` timestamp
//     and `body` of the token, if the token is valid and has not expired.
//   - An error if the token is invalid, expired, or if any issues occur during the decryption process.
//
// Note:
//   - The function checks if the timestamp stored in the token is within the allowed expiration window.
//   - If the token has expired or cannot be decrypted, the function will return ErrExpiredToken error.
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
	createdAt, err := time.Parse(internal.DateTimeFormat, timedTokenData.CreatedAt)
	if err != nil {
		return nil, err
	}

	if timeNow.Sub(createdAt) > expiresIn {
		return nil, ErrExpiredToken
	}

	return &timedTokenData, nil
}

// getBlockCipher creates a new AES-GCM block cipher using the provided key.
// The key is first hashed using SHA-256 to create a 32-byte key suitable for AES encryption.
// The resulting block cipher is returned as an AEAD (Authenticated Encryption with
// Associated Data) cipher for use in secure communication or data encryption.
//
// Parameters:
//   - key: The key to be used in the AES cipher. It will be hashed using SHA-256 to
//     generate the final 32-byte key for AES encryption.
//
// Returns:
//   - A cipher.AEAD object that can be used for encryption and decryption.
//   - An error if the cipher cannot be created due to issues with the key or
//     any failure during cipher creation.
//
// Note:
//   - SHA-256 is used here to generate a secure 32-byte key from the input. This is a
//     strong cryptographic hash function that provides enhanced security compared to MD5.
//   - The AES-GCM (Galois/Counter Mode) cipher is used for authenticated encryption,
//     providing both confidentiality and data integrity.
func getBlockCipher(key []byte) (cipher.AEAD, error) {
	sha256Hash := sha256.Sum256(key)

	aesBlock, err := aes.NewCipher(sha256Hash[:])
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(aesBlock)
}

// generateToken encrypts the provided data along with a timestamp and returns a base64 URL-encoded string.
//
// It uses AES-GCM for encryption with a randomly generated nonce. The function serializes the data and timestamp
// into a JSON object, encrypts it, and encodes the result in a URL-safe format.
//
// Parameters:
//   - data (string): The data to be included in the token.
//   - key ([]byte): The encryption key. Must be 32 bytes long for AES-256.
//   - createdAt (string): The timestamp (e.g., time.Now().Format(time.RFC3339)).
//
// Returns:
//   - string: The base64 URL-encoded encrypted token.
//   - error: An error if encryption or encoding fails.
//
// Example usage:
//
//	key := []byte("your-32-byte-secret-key-----------------")
//	token, err := generateToken("user123", key, time.Now().Format(time.RFC3339))
//	if err != nil {
//	    log.Fatalf("Failed to generate token: %v", err)
//	}
//	fmt.Println("Generated Token:", token)
func generateToken(data string, key []byte, createdAt string) (string, error) {
	gcm, err := getBlockCipher(key)
	if err != nil {
		return "", err
	}

	timedTokenData := TimedTokenData{
		CreatedAt: createdAt,
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
