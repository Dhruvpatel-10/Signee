package config

import (
	"net/http"
	"time"

	"github.com/dhruvpatel-10/signee/ca-api/internal/domain/auth"
)

func DefaultCookieConfig() auth.SecureCookieConfig {
	return auth.SecureCookieConfig{
		TokenName: "__Host-refresh-token",
		Domain:    "",
		Path:      "/auth/refresh",
		MaxAge:    15 * 60,
		Secure:    true,
		HttpOnly:  true,
		SameSite:  http.SameSiteStrictMode,
	}
}

// SetTokenCookies sets secure cookies for tokens
func SetTokenCookies(w http.ResponseWriter, refreshToken string, config auth.SecureCookieConfig) {

	refreshCookie := &http.Cookie{
		Name:     config.TokenName,
		Value:    refreshToken,
		Domain:   config.Domain,
		Path:     "/auth/refresh",
		MaxAge:   int(7 * 24 * time.Hour / time.Second),
		Secure:   config.Secure,
		HttpOnly: config.HttpOnly,
		SameSite: config.SameSite,
	}
	http.SetCookie(w, refreshCookie)
}
