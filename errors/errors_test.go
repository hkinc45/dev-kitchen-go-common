package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAPIError(t *testing.T) {
	innerErr := errors.New("database connection failed")
	apiErr := NewAPIErrorWrap(http.StatusInternalServerError, "internal server error", innerErr)

	assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
	assert.Equal(t, "internal server error", apiErr.Message)
	assert.Equal(t, innerErr, apiErr.Unwrap())
	assert.Contains(t, apiErr.Error(), "database connection failed")
}

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("APIError Handling", func(t *testing.T) {
		r := gin.New()
		r.Use(Middleware())
		r.GET("/error", func(c *gin.Context) {
			c.Error(NewNotFoundError("resource not found"))
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "resource not found")
	})

	t.Run("Unexpected Error Handling", func(t *testing.T) {
		r := gin.New()
		r.Use(Middleware())
		r.GET("/panic", func(c *gin.Context) {
			c.Error(errors.New("random error"))
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/panic", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "An unexpected internal error occurred")
	})
}
