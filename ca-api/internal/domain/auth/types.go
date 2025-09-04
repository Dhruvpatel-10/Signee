// internal/domain/auth/types.go
package auth

import (
	"time"

	"github.com/google/uuid"
)

type SignupRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"` // Can be email or username
	Password string `json:"password" binding:"required"`
}

type Session struct {
	ID           string    `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	Roles        []string  `json:"roles"`
	Permissions  []string  `json:"permissions"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
	ExpiresAt    time.Time `json:"expires_at"`
	MFAVerified  bool      `json:"mfa_verified"`
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
