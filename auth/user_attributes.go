package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SetUserAttribute safely updates a user's attributes in Keycloak by performing a read-modify-write.
// It first fetches the full user representation, updates the attributes, and then PUTs the entire object back.
// This is done using manual API calls to bypass bugs in some versions of the gocloak library's UpdateUser function.
func SetUserAttribute(ctx context.Context, adminAPIURL, realm, userID, adminAccessToken string, attributes map[string][]string) error {
	userURL := fmt.Sprintf("%s/admin/realms/%s/users/%s", adminAPIURL, realm, userID)

	// 1. GET the full user representation
	getReq, err := http.NewRequestWithContext(ctx, "GET", userURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create get user request: %w", err)
	}
	getReq.Header.Set("Authorization", "Bearer "+adminAccessToken)

	getResp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		return fmt.Errorf("failed to perform get user request: %w", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		return fmt.Errorf("get user failed with status %d", getResp.StatusCode)
	}

	var userRepresentation map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&userRepresentation); err != nil {
		return fmt.Errorf("failed to decode user representation: %w", err)
	}

	// 2. Modify the attributes
	userRepresentation["attributes"] = attributes

	// 3. PUT the full, modified user representation back
	jsonPayload, err := json.Marshal(userRepresentation)
	if err != nil {
		return fmt.Errorf("failed to marshal updated user representation: %w", err)
	}

	putReq, err := http.NewRequestWithContext(ctx, "PUT", userURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create set user attribute request: %w", err)
	}

	putReq.Header.Set("Content-Type", "application/json")
	putReq.Header.Set("Authorization", "Bearer "+adminAccessToken)

	putResp, err := http.DefaultClient.Do(putReq)
	if err != nil {
		return fmt.Errorf("failed to perform set user attribute request: %w", err)
	}
	defer putResp.Body.Close()

	if putResp.StatusCode != http.StatusNoContent && putResp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		_ = json.NewDecoder(putResp.Body).Decode(&errResp)
		return fmt.Errorf("set user attribute failed with status %d: %v", putResp.StatusCode, errResp)
	}

	return nil
}

// AreAttributesEqual compares two attribute maps to see if they are functionally identical.
func AreAttributesEqual(a, b map[string][]string) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valA := range a {
		valB, ok := b[key]
		if !ok {
			return false
		}
		if len(valA) != len(valB) {
			return false
		}
		// Create maps for quick lookups, as order doesn't matter
		mapA := make(map[string]struct{}, len(valA))
		for _, v := range valA {
			mapA[v] = struct{}{}
		}
		for _, v := range valB {
			if _, found := mapA[v]; !found {
				return false
			}
		}
	}

	return true
}
