package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// TokenExchangeResponse represents the successful response from a token exchange request.
type TokenExchangeResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	Scope            string `json:"scope"`
}

// PerformTokenExchange performs a standard RFC 8693 token exchange.
// It uses raw HTTP requests to ensure compatibility with modern Keycloak versions,
// bypassing potential issues with the gocloak library's token exchange implementation.
func PerformTokenExchange(ctx context.Context, tokenURL, clientID, clientSecret, subjectToken, audience string) (*TokenExchangeResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("subject_token", subjectToken)
	data.Set("subject_token_type", "urn:ietf:params:oauth:token-type:access_token")
	data.Set("audience", audience)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token exchange request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform token exchange request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Try to read the error response body for better diagnostics
		var errResp map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("token exchange failed with status %d: %s - response: %v", resp.StatusCode, resp.Status, errResp)
	}

	var tokenResp TokenExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode successful token exchange response: %w", err)
	}

	return &tokenResp, nil
}
