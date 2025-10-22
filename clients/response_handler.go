package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hkinc45/dev-kitchen-go-common/errors"
)

// HandleResponse handles decoding HTTP responses from other services.
// It decodes either the success body or an APIError.
func HandleResponse(resp *http.Response, successBody interface{}) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Special handling for 409 Conflict
		if resp.StatusCode == http.StatusConflict {
			return errors.ErrConflict
		}
		if resp.StatusCode == http.StatusNotFound {
			return errors.NewNotFoundError("resource not found")
		}

		var apiErr errors.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			// If we can't decode the error, return a generic one.
			bodyBytes, _ := io.ReadAll(resp.Body)
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
