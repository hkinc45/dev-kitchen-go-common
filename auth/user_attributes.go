package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserAttribute represents a simple structure for updating user attributes.
type UserAttribute struct {
	Attributes map[string][]string `json:"attributes"`
}

// SetUserAttribute updates a user's attributes in Keycloak using a manual API call
// to bypass bugs in the gocloak library's UpdateUser function.
func SetUserAttribute(ctx context.Context, adminAPIURL, realm, userID, adminAccessToken string, attributes map[string][]string) error {
	url := fmt.Sprintf("%s/admin/realms/%s/users/%s", adminAPIURL, realm, userID)

	payload := UserAttribute{
		Attributes: attributes,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal user attribute payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create set user attribute request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminAccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform set user attribute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("set user attribute failed with status %d: %v", resp.StatusCode, errResp)
	}

	return nil
}
