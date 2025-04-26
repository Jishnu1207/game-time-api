package middleware

import (
	"game-time-api/api"
	"game-time-api/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware verifies the JWT token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			api.UnauthorizedError(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			api.UnauthorizedError(c, "Authorization header must be in the format: Bearer <token>")
			c.Abort()
			return
		}

		// Get the token
		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			// Get the secret key from environment
			secretKey := config.GetEnv("JWT_SECRET", "your_jwt_secret_key")
			return []byte(secretKey), nil
		})

		if err != nil {
			api.UnauthorizedError(c, "Invalid token")
			c.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			api.UnauthorizedError(c, "Invalid token")
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			api.UnauthorizedError(c, "Invalid token claims")
			c.Abort()
			return
		}

		// Set user ID in context for later use
		userID, ok := claims["user_id"].(float64)
		if !ok {
			api.UnauthorizedError(c, "Invalid user ID in token")
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("user_id", uint(userID))
		c.Set("email", claims["email"])

		// Continue to the next handler
		c.Next()
	}
}
