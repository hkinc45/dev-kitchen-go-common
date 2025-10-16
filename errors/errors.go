package errors

import (
	"fmt"
	"net/http"
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

// Pre-defined error types
var (
	ErrConflict = NewAPIError(http.StatusConflict, "resource already exists")
)
