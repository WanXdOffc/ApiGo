package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	// KeyPrefix is standard for production API keys
	KeyPrefix = "sk_live_"
	// KeyRandomByteLength specifies 32 bytes (256 bits) of cryptographic entropy
	KeyRandomByteLength = 32
)

// GenerateAPIKey generates a cryptographically secure random API key string with the standard prefix
func GenerateAPIKey() (string, error) {
	randomBytes := make([]byte, KeyRandomByteLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate secure random bytes: %w", err)
	}

	randomHex := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("%s%s", KeyPrefix, randomHex), nil
}

// HashAPIKey computes the SHA-256 hash of a raw API key string
func HashAPIKey(rawKey string) string {
	hasher := sha256.New()
	hasher.Write([]byte(rawKey))
	return hex.EncodeToString(hasher.Sum(nil))
}

// MaskAPIKey returns a masked representation suitable for UI display (e.g., sk_live_****abcd)
func MaskAPIKey(rawKey string) string {
	if len(rawKey) <= len(KeyPrefix)+4 {
		return "sk_live_****"
	}
	last4 := rawKey[len(rawKey)-4:]
	return fmt.Sprintf("%s****%s", KeyPrefix, last4)
}
