package handlers

import (
	"game-time-api/api"
	"game-time-api/models"
	"game-time-api/repositories"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	User      *models.User `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// Login handles user authentication
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	// Create user repository
	userRepo := repositories.NewUserRepository()

	// Find user by email
	user, err := userRepo.FindByEmail(req.Email)
	if err != nil {
		api.UnauthorizedError(c, "Invalid email or password")
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		api.UnauthorizedError(c, "Invalid email or password")
		return
	}

	// Generate token
	token, expiresAt, err := generateToken(*user)
	if err != nil {
		api.ServerError(c, "Error generating token")
		return
	}

	// Return response
	response := LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}

	api.Success(c, http.StatusOK, "Login successful", response)
}

// generateToken creates a new JWT token for the user
func generateToken(user models.User) (string, time.Time, error) {
	// Set expiration time
	expiresAt := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     expiresAt.Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
