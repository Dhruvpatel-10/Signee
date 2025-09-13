package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/dhruvpatel-10/signee/ca-api/db"
	"github.com/dhruvpatel-10/signee/ca-api/internal/domain/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Your existing signup function (looks good!)
func (s *AuthService) Signup(c *gin.Context) {
	var req auth.SignupRequest
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

	emailUsr, err := s.DB.GetUserByEmail(c, req.Email)
	if err != nil && err != sql.ErrNoRows {
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

	if err == nil && req.Email == emailUsr.Email {
		// Found existing user
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "USER_ALREADY_EXISTS",
				"message": "A user with this email already exists.",
				"status":  http.StatusConflict,
			},
		})
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Something went wrong. Please try again later.",
				"status":  http.StatusInternalServerError,
			},
		})
		return
	}

	user, err := s.DB.CreateUser(c, db.CreateUserParams{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		MfaSecret:    sql.NullString{},
		MfaEnabled:   sql.NullBool{},
		CreatedBy:    uuid.NullUUID{},
	})

	if err != nil {
		log.Printf("CreateUser failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Something went wrong. Please try again later.",
				"status":  http.StatusInternalServerError,
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Congratulations! You've successfully registered",
		"username": user.FirstName + " " + user.LastName,
	})
}
