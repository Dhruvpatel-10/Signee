// internal/domain/auth/types.go
package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignupRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	FirstName       string `json:"fname" binding:"required"`
	LastName        string `json:"lname" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"` // Can be email or username
	Password string `json:"password" binding:"required"`
}

// SecureClaims for JWT with encrypted fields
type SecureClaims struct {
	UserIDEnc       string `json:"uid_enc"` // Encrypted UserID
	RoleEnc         string `json:"rol_enc"` // Encrypted Role
	TokenType       string `json:"typ"`     // "access" or "refresh"
	Fingerprint     string `json:"fp"`      // Browser fingerprint hash
	IPHash          string `json:"iph"`     // IP address hash
	DecryptedUserID string // Not serialized, used after decryption
	DecryptedRole   string // Not serialized, used after decryption
	jwt.RegisteredClaims
}

// SecureCookieConfig for secure cookies
type SecureCookieConfig struct {
	TokenName string
	Domain    string
	Path      string
	MaxAge    uint32
	Secure    bool
	HttpOnly  bool
	SameSite  http.SameSite
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type AuthResponse struct {
	// User         *UserInfo  `json:"user"`
	Tokens       *TokenPair `json:"tokens"`
	RequiresMFA  bool       `json:"requires_mfa,omitempty"`
	MFAChallenge string     `json:"mfa_challenge,omitempty"`
}

type PassConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}
