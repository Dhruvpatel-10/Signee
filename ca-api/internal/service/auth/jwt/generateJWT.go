package jwt

import (
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dhruvpatel-10/signee/ca-api/internal/domain/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// JWT Manager with elliptic keys for better security
type JWTManager struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	accessExpiration  time.Duration
	refreshExpiration time.Duration
	issuer            string
	encryptionKey     []byte
}

func NewJWTManager(issuer string) (*JWTManager, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	// Load or generate ECDSA key
	privateKeyPath := os.Getenv("ECDSA_PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "ecdsa_private.pem" // Default path
	}
	privateKey, err := loadOrGenerateECDSAKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ECDSA key: %v", err)
	}

	// Load or generate ChaCha20-Poly1305 key
	encryptionKeyPath := os.Getenv("CHACHA20_KEY_PATH")
	if encryptionKeyPath == "" {
		encryptionKeyPath = "chacha20_key.bin" // Default path
	}
	encryptionKey, err := loadOrGenerateEncryptionKey(encryptionKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize encryption key: %v", err)
	}

	return &JWTManager{
		privateKey:        privateKey,
		publicKey:         &privateKey.PublicKey,
		accessExpiration:  15 * time.Minute,
		refreshExpiration: 7 * 24 * time.Hour,
		issuer:            issuer,
		encryptionKey:     encryptionKey,
	}, nil
}

func (j *JWTManager) GenerateTokens(userID, username, role string, r *http.Request) (string, string, error) {
	now := time.Now()
	fingerprint := j.generateFingerprint(r)
	ipHash := j.hashIP(getRealIP(r))

	// Encrypt UserID and Role
	userIDEnc, err := encryptData([]byte(userID), j.encryptionKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt UserID: %v", err)
	}
	roleEnc, err := encryptData([]byte(role), j.encryptionKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt Role: %v", err)
	}

	// Access Token
	accessClaims := auth.SecureClaims{
		UserIDEnc:   userIDEnc,
		RoleEnc:     roleEnc,
		TokenType:   "access",
		Fingerprint: fingerprint,
		IPHash:      ipHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   userIDEnc,
			ID:        generateJTI(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodES256, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %v", err)
	}

	// Refresh Token
	refreshClaims := auth.SecureClaims{
		UserIDEnc:   userIDEnc,
		TokenType:   "refresh",
		Fingerprint: fingerprint,
		IPHash:      ipHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   userIDEnc,
			ID:        generateJTI(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(j.privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %v", err)
	}

	return accessTokenString, refreshTokenString, nil
}

// ValidateAccessToken validates an access token
func (j *JWTManager) ValidateAccessToken(tokenString string, r *http.Request) (*auth.SecureClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.SecureClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*auth.SecureClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.TokenType != "access" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Decrypt UserID and Role
	userID, err := decryptData(claims.UserIDEnc, j.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt userID: %v", err)
	}
	claims.DecryptedUserID = string(userID)

	role, err := decryptData(claims.RoleEnc, j.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt role: %v", err)
	}
	claims.DecryptedRole = string(role)

	// Verify fingerprint
	if claims.Fingerprint != j.generateFingerprint(r) {
		return nil, fmt.Errorf("token fingerprint mismatch")
	}

	// Log IP mismatch (no logout)
	currentIPHash := j.hashIP(getRealIP(r))
	if claims.IPHash != currentIPHash {
		fmt.Printf("Warning: IP mismatch for user %s\n", claims.DecryptedUserID)
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (j *JWTManager) ValidateRefreshToken(tokenString string, r *http.Request) (*auth.SecureClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.SecureClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(*auth.SecureClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Decrypt UserID
	userID, err := decryptData(claims.UserIDEnc, j.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt userID: %v", err)
	}
	claims.DecryptedUserID = string(userID)

	if claims.Fingerprint != j.generateFingerprint(r) {
		return nil, fmt.Errorf("refresh token fingerprint mismatch")
	}

	// Log IP mismatch (no logout)
	currentIPHash := j.hashIP(getRealIP(r))
	if claims.IPHash != currentIPHash {
		fmt.Printf("Warning: IP mismatch for user %s\n", claims.DecryptedUserID)
	}

	return claims, nil
}
