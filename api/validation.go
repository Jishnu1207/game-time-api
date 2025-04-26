package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// FormatValidationError formats validation errors into a user-friendly message
func FormatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			// Convert field name to lowercase for better readability
			field = strings.ToLower(field)

			// Create user-friendly error messages based on validation tags
			switch tag {
			case "required":
				errorMessages = append(errorMessages, field+" is required")
			case "email":
				errorMessages = append(errorMessages, "invalid email format")
			case "min":
				errorMessages = append(errorMessages, field+" must be at least "+e.Param()+" characters")
			case "max":
				errorMessages = append(errorMessages, field+" must not exceed "+e.Param()+" characters")
			default:
				errorMessages = append(errorMessages, field+" failed on "+tag+" validation")
			}
		}
		return strings.Join(errorMessages, "; ")
	}

	// If it's not a validation error, return the original error message
	return err.Error()
}

// HandleValidationError handles validation errors and sends a standardized response
func HandleValidationError(c *gin.Context, err error) {
	message := FormatValidationError(err)
	ValidationError(c, message)
}
