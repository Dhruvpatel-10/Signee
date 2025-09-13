package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/chacha20poly1305"
)

// loadOrGenerateEncryptionKey loads or generates ChaCha20-Poly1305 key
func loadOrGenerateEncryptionKey(filename string) ([]byte, error) {
	if _, err := os.Stat(filename); err == nil {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read encryption key file: %v", err)
		}
		if len(data) != chacha20poly1305.KeySize {
			return nil, fmt.Errorf("invalid encryption key length")
		}
		return data, nil
	}

	// Generate new 32-byte key
	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate encryption key: %v", err)
	}

	// Save the key
	if err := os.WriteFile(filename, key, 0600); err != nil {
		return nil, fmt.Errorf("failed to save encryption key: %v", err)
	}

	return key, nil
}

// encryptData encrypts data with ChaCha20-Poly1305
func encryptData(data, key []byte) (string, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AEAD cipher: %v", err)
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}
	ciphertext := aead.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptData decrypts data with ChaCha20-Poly1305
func decryptData(encrypted string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted data: %v", err)
	}
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD cipher: %v", err)
	}
	if len(data) < aead.NonceSize() {
		return nil, fmt.Errorf("invalid ciphertext length")
	}
	nonce, ciphertext := data[:aead.NonceSize()], data[aead.NonceSize():]
	return aead.Open(nil, nonce, ciphertext, nil)
}

// generateFingerprint creates a browser fingerprint hash
func (j *JWTManager) generateFingerprint(r *http.Request) string {
	fingerprint := r.Header.Get("User-Agent") +
		r.Header.Get("Accept-Language") +
		r.Header.Get("Accept-Encoding") // Simplified; use secure derivation in production
	hash := sha256.Sum256([]byte(fingerprint))
	return base64.StdEncoding.EncodeToString(hash[:16])
}

// hashIP hashes the client IP address
func (j *JWTManager) hashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip + "your-secret-salt")) // Use secure salt in production
	return base64.StdEncoding.EncodeToString(hash[:16])
}

// getRealIP extracts the client's real IP address
func getRealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func generateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
