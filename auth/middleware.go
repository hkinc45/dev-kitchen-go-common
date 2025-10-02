package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/hkinc45/dev-kitchen-go-common/models"
)

// Middleware holds the OIDC token verifier and other configuration for auth checks.
type Middleware struct {
	Verifier       *oidc.IDTokenVerifier
	ClientID       string
	AuthServiceURL string
}

// NewMiddleware creates a new OIDC-based authentication middleware.
func NewMiddleware(ctx context.Context, providerURL, clientID, authServiceURL string) (*Middleware, error) {
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})

	return &Middleware{
		Verifier:       verifier,
		ClientID:       clientID,
		AuthServiceURL: authServiceURL,
	}, nil
}

// UserAuth is a middleware for validating tokens from end-users.
// It now performs a JIT provisioning step by calling the auth-service.
func (m *Middleware) UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		idToken, err := m.Verifier.Verify(c.Request.Context(), tokenString)
		if err != nil {
			log.Printf("ERROR: Token verification failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			log.Printf("ERROR: Failed to extract claims from token: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract claims from token"})
			return
		}

		if !m.isAudienceValid(claims) {
			log.Printf("ERROR: Token audience validation failed. Expected '%s' in audience %v", m.ClientID, claims["aud"])
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Token not valid for this service"})
			return
		}

		// JIT Provisioning: Call the auth-service's /me endpoint to get the full user object.
		// This ensures the user exists in the auth-service DB and we get the canonical Application ID.
		user, err := m.jitProvisionUser(c.Request.Context(), authHeader)
		if err != nil {
			log.Printf("ERROR: JIT provisioning failed: %v", err)
			c.AbortWithStatusJSON(http.StatusFailedDependency, gin.H{"error": "Failed to retrieve user profile from auth service"})
			return
		}

		// Set the full user object in the context.
		c.Set("user", user)

		log.Println("User token validated and user object set successfully.")
		c.Next()
	}
}

// jitProvisionUser calls the auth-service's /me endpoint to get the user object.
func (m *Middleware) jitProvisionUser(ctx context.Context, authHeader string) (*models.User, error) {
	meURL := fmt.Sprintf("%s/api/v1/me", m.AuthServiceURL)
	req, err := http.NewRequestWithContext(ctx, "GET", meURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to auth-service: %w", err)
	}
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to auth-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth-service returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user object from auth-service: %w", err)
	}

	return &user, nil
}


// ServiceAuth is a middleware for validating tokens from other services.
// It checks that the token has the required `internal-comm` role.
func (m *Middleware) ServiceAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		idToken, err := m.Verifier.Verify(c.Request.Context(), tokenString)
		if err != nil {
			log.Printf("ERROR: Token verification failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			log.Printf("ERROR: Failed to extract claims from token: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract claims from token"})
			return
		}

		// For service tokens, we check for the 'internal-comm' role.
		if !m.hasInternalCommRole(claims) {
			log.Printf("ERROR: Service token is missing 'internal-comm' role.")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: internal-comm role required"})
			return
		}

		azp, _ := claims["azp"].(string)
		log.Printf("Service token from '%%s' validated successfully.", azp)
		c.Next()
	}
}

// isAudienceValid checks if the service's ClientID is present in the 'aud' claim.
// It handles both string and []string formats for the audience claim.
func (m *Middleware) isAudienceValid(claims map[string]interface{}) bool {
	aud, ok := claims["aud"]
	if !ok {
		return false
	}

	switch v := aud.(type) {
	case string:
		return v == m.ClientID
	case []interface{}:
		for _, a := range v {
			if s, ok := a.(string); ok && s == m.ClientID {
				return true
			}
		}
	}
	return false
}

// hasInternalCommRole checks if the 'internal-comm' role is present in the token.
func (m *Middleware) hasInternalCommRole(claims map[string]interface{}) bool {
	realmAccess, ok := claims["realm_access"].(map[string]interface{})
	if !ok {
		return false
	}

	roles, ok := realmAccess["roles"].([]interface{})
	if !ok {
		return false
	}

	for _, r := range roles {
		if role, ok := r.(string); ok && role == "internal-comm" {
			return true
			}
	}
	return false
}