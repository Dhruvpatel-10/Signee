package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/dhruvpatel-10/signee/ca-api/internal/domain/auth"
	"github.com/gin-gonic/gin"
)

// Login function
func (s *AuthService) Login(c *gin.Context) {
	var req auth.LoginRequest // You'll need to define this struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := s.DB.GetUserByEmail(c, req.Email) // You'll need this DB method
	if err != nil && err != sql.ErrNoRows {
		// Only fatal for unexpected DB errors
		log.Printf("GetUserByEmail failed or User does not exist : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Verify password
	isValid, err := verifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		log.Printf("verifyPassword failed: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !isValid {
		log.Printf("Password does not match for email: %s", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token or session (implement this based on your auth strategy)
	// token, err := generateJWTToken(user.ID) // You'll need to implement this
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message":    "login successful",
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})
}
