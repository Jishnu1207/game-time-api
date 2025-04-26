package handlers

import (
	"game-time-api/api"
	"game-time-api/models"
	"game-time-api/repositories"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=5,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// validateEmail checks if the email is valid
func validateEmail(email string) bool {
	// Basic email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validatePassword checks if the password meets the requirements
func validatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}

	// Check for at least one number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return false, "Password must contain at least one number"
	}

	// Check for at least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasSpecial {
		return false, "Password must contain at least one special character"
	}

	return true, ""
}

// Register handles user registration
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	// Trim whitespace from username and email
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	// Validate email format
	if !validateEmail(req.Email) {
		api.ValidationError(c, "Invalid email format")
		return
	}

	// Validate password requirements
	if valid, message := validatePassword(req.Password); !valid {
		api.ValidationError(c, message)
		return
	}

	// Create user repository
	userRepo := repositories.NewUserRepository()

	// Check if username already exists
	usernameExists, err := userRepo.UsernameExists(req.Username)
	if err != nil {
		api.ServerError(c, "Error checking username availability")
		return
	}
	if usernameExists {
		api.ConflictError(c, "Username is already taken")
		return
	}

	// Check if email already exists
	emailExists, err := userRepo.EmailExists(req.Email)
	if err != nil {
		api.ServerError(c, "Error checking email availability")
		return
	}
	if emailExists {
		api.ConflictError(c, "Email is already registered")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		api.ServerError(c, "Error hashing password")
		return
	}

	// Create user
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := userRepo.Create(&user); err != nil {
		api.ServerError(c, "Error creating user")
		return
	}

	// Return user data without password
	response := RegisterResponse{
		Username: user.Username,
		Email:    user.Email,
	}

	api.Success(c, http.StatusCreated, "User registered successfully", response)
}
