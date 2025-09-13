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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid input provided.",
				"status":  http.StatusBadRequest,
			},
		})
		return
	}

	// Get user by email
	// Get user by email
	user, err := s.DB.GetUserByEmail(c, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			log.Printf("User not found for email: %s", req.Email)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_CREDENTIALS",
					"message": "Invalid email or password.",
					"status":  http.StatusUnauthorized,
				},
			})
			return
		}
		// Only fatal for unexpected DB errors
		log.Printf("GetUserByEmail failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Something went wrong. Please try again later.",
				"status":  http.StatusInternalServerError,
			},
		})
		return
	}

	// Verify password
	switch isValid, err := verifyPassword(req.Password, user.PasswordHash); {
	case err != nil:
		log.Printf("verifyPassword failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Something went wrong. Please try again later.",
				"status":  http.StatusInternalServerError,
			},
		})
		return
	case !isValid:
		log.Printf("Password does not match for email: %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Invalid email or password.",
				"status":  http.StatusUnauthorized,
			},
		})
		return
	}

	// Generate JWT token or session (implement this based on your auth strategy)
	// token, err := generateJWTToken(user.ID) // You'll need to implement this
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": gin.H{
	// 			"code":    "INTERNAL_SERVER_ERROR",
	// 			"message": "Something went wrong. Please try again later.",
	// 			"status":  http.StatusInternalServerError,
	// 		},
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message":  "Welcome Back",
		"username": user.FirstName + " " + user.LastName,
	})
}
