package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// loadOrGenerateECDSAKey loads ECDSA key or generates and saves a new one
func loadOrGenerateECDSAKey(filename string) (*ecdsa.PrivateKey, error) {
	if _, err := os.Stat(filename); err == nil {
		return loadECDSAKey(filename)
	}

	// Generate new ECDSA key (P-256 curve for ES256)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECDSA key: %v", err)
	}

	// Save the key
	if err := saveECDSAKey(filename, privateKey); err != nil {
		return nil, fmt.Errorf("failed to save ECDSA key: %v", err)
	}

	return privateKey, nil
}

// saveECDSAKey saves ECDSA private key to file with 0600 permissions
func saveECDSAKey(filename string, key *ecdsa.PrivateKey) error {
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshal ECDSA key: %v", err)
	}
	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	}
	return os.WriteFile(filename, pem.EncodeToMemory(pemBlock), 0600)
}

// loadECDSAKey loads ECDSA private key from file
func loadECDSAKey(filename string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read ECDSA key file: %v", err)
	}
	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil || pemBlock.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("invalid PEM block in %s", filename)
	}
	key, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ECDSA key: %v", err)
	}
	return key, nil
}
