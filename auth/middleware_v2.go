package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	common_errors "github.com/hkinc45/dev-kitchen-go-common/errors"
)

// CheckPermissionRequest defines the structure for requests to the auth service's check endpoint.
type CheckPermissionRequest struct {
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	Scope        string `json:"scope"`
	SubjectToken string `json:"subject_token"`
}

// ResourceIDExtractor is a function that extracts a resource's ID from the request context.
type ResourceIDExtractor func(c *gin.Context) (string, error)

// RequirePermissionV2 creates a Gin middleware that checks if a user has a specific permission for a dynamic resource.
// It works by calling the internal `/v2/auth/check` endpoint in the auth-service.
//
// - httpClient: An authenticated HTTP client for service-to-service calls.
// - resourceType: The type of resource being checked (e.g., "project", "recipe").
// - idExtractor: A function that extracts the resource's ID from the Gin context.
// - scope: The scope to check for (e.g., "project:read").
func RequirePermissionV2(httpClient *http.Client, resourceType string, idExtractor ResourceIDExtractor, scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get auth service URL from environment
		authServiceURL := os.Getenv("AUTH_SERVICE_URL")
		if authServiceURL == "" {
			c.Error(common_errors.NewInternalServerError("misconfigured authentication service URL"))
			c.Abort()
			return
		}

		// 2. Get the raw user token from the Authorization header.
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Error(common_errors.NewUnauthorizedError("authorization header missing or improperly formatted"))
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Extract the resource ID using the provided extractor function
		resourceID, err := idExtractor(c)
		if err != nil {
			c.Error(common_errors.NewBadRequestError(fmt.Sprintf("failed to extract resource ID for permission check: %v", err)))
			c.Abort()
			return
		}

		// 4. Construct the request to the auth service
		checkReqPayload := CheckPermissionRequest{
			ResourceType: resourceType,
			ResourceID:   resourceID,
			Scope:        scope,
			SubjectToken: token,
		}

		payloadBytes, err := json.Marshal(checkReqPayload)
		if err != nil {
			c.Error(common_errors.NewInternalServerError("failed to construct permission check request"))
			c.Abort()
			return
		}

		// 5. Call the auth service's check endpoint
		checkURL := fmt.Sprintf("%s/internal/v2/auth/check", authServiceURL)
		req, err := http.NewRequestWithContext(c.Request.Context(), "POST", checkURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			c.Error(common_errors.NewInternalServerError("failed to create permission check request"))
			c.Abort()
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			c.Error(common_errors.NewInternalServerError("failed to communicate with authentication service"))
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// 6. Handle the response
		switch resp.StatusCode {
		case http.StatusOK:
			c.Next() // Permission granted
		case http.StatusForbidden:
			// The error message from the auth service is now more generic, so we create a specific one here.
			c.Error(common_errors.NewForbiddenError(fmt.Sprintf("missing required permission: %s on resource %s:%s", scope, resourceType, resourceID)))
			c.Abort()
		default:
			c.Error(common_errors.NewInternalServerError(fmt.Sprintf("unexpected error from authentication service: status %d", resp.StatusCode)))
			c.Abort()
		}
	}
}
