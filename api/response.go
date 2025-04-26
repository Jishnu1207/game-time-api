package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// ServerError sends a server error response
func ServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// NotFoundError sends a not found error response
func NotFoundError(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// ConflictError sends a conflict error response
func ConflictError(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message)
}

// UnauthorizedError sends an unauthorized error response
func UnauthorizedError(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// ForbiddenError sends a forbidden error response
func ForbiddenError(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}
