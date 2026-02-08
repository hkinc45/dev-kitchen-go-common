package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hkinc45/dev-kitchen-go-common/errors"
)

// HandleResponse handles decoding HTTP responses from other services.
// It decodes either the success body or an APIError.
func HandleResponse(resp *http.Response, successBody interface{}) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read the full body to log it for debugging non-2xx responses.
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading non-2xx response body: %v", err)
			return errors.NewAPIError(resp.StatusCode, "failed to read error response body")
		}
		// Log the detailed error response.
		log.Printf("Downstream service returned non-2xx response. Status: %d, Body: %s", resp.StatusCode, string(bodyBytes))

		// Replace the response body with a new reader so it can be read again by the JSON decoder.
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if resp.StatusCode == http.StatusConflict {
			return errors.ErrConflict
		}
		if resp.StatusCode == http.StatusNotFound {
			return errors.NewNotFoundError("resource not found")
		}

		var apiErr errors.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			// If we can't decode a structured error, use the raw body we already read.
			return errors.NewAPIError(resp.StatusCode, fmt.Sprintf("unknown error: %s", string(bodyBytes)))
		}
		apiErr.StatusCode = resp.StatusCode // Ensure status code is set
		return &apiErr
	}

	if successBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(successBody); err != nil {
			return errors.NewAPIError(http.StatusInternalServerError, fmt.Sprintf("failed to decode success response: %v", err))
		}
	}

	return nil
}
