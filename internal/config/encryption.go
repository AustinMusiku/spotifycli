package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// EncryptToken encrypts a token using AES-256-GCM
func EncryptToken(plaintext, key string) (string, error) {
	// Hash the key to get 32 bytes
	hash := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptToken decrypts a token using AES-256-GCM
func DecryptToken(ciphertext, key string) (string, error) {
	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Hash the key to get 32 bytes
	hash := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], string(data[nonceSize:])

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, []byte(ciphertext), nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// generateEncryptionKey generates a random encryption key
func generateEncryptionKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		// Fallback to a deterministic key based on user home directory
		homeDir, _ := os.UserHomeDir()
		hash := sha256.Sum256([]byte(homeDir + "spotify-cli-key"))
		return base64.StdEncoding.EncodeToString(hash[:])
	}
	return base64.StdEncoding.EncodeToString(key)
}
