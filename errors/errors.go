package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError represents a structured error response from a service.
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"error"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: status code %d, message: %s", e.StatusCode, e.Message)
}

func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// Middleware is a Gin middleware for centralized error handling.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process request

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if apiErr, ok := err.(*APIError); ok {
				c.JSON(apiErr.StatusCode, apiErr)
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "An unexpected internal error occurred",
			})
		}
	}
}

// Pre-defined error types
var (
	ErrConflict = NewAPIError(http.StatusConflict, "resource already exists")
)
