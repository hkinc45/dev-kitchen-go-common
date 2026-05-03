package auth

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	common_errors "github.com/hkinc45/dev-kitchen-go-common/errors"
	"github.com/stretchr/testify/assert"
)

// RoundTripFunc allows mocking http.Client
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestRequirePermissionV2(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("AUTH_SERVICE_URL", "http://auth-service")
	defer os.Unsetenv("AUTH_SERVICE_URL")

	resourceID := "proj-123"
	extractor := func(c *gin.Context) (string, error) {
		return resourceID, nil
	}

	t.Run("Success", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: RoundTripFunc(func(req *http.Request) *http.Response {
				assert.Equal(t, "http://auth-service/internal/v2/auth/check", req.URL.String())
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"status":"permitted"}`)),
					Header:     make(http.Header),
				}
			}),
		}

		r := gin.New()
		r.Use(common_errors.Middleware())
		r.Use(RequirePermissionV2(mockClient, "project", extractor, "project:read"))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: RoundTripFunc(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error":"forbidden"}`)),
					Header:     make(http.Header),
				}
			}),
		}

		r := gin.New()
		r.Use(common_errors.Middleware())
		r.Use(RequirePermissionV2(mockClient, "project", extractor, "project:read"))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Missing Header", func(t *testing.T) {
		mockClient := &http.Client{}
		r := gin.New()
		r.Use(common_errors.Middleware())
		r.Use(RequirePermissionV2(mockClient, "project", extractor, "project:read"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
