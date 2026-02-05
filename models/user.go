package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the application's database.
// It is the canonical representation of a user across all services.
type User struct {
	ID            uuid.UUID           `json:"id"`
	KeycloakID    string              `json:"keycloak_id"`
	Username      string              `json:"username"`
	Email         string              `json:"email"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	FirstName     *string             `json:"first_name,omitempty"`
	LastName      *string             `json:"last_name,omitempty"`
	PhoneNumber   *string             `json:"phone_number,omitempty"`
	AccountType   *string             `json:"account_type,omitempty"` // e.g., "personal" or "business"
	StreetAddress *string             `json:"street_address,omitempty"`
	City          *string             `json:"city,omitempty"`
	State         *string             `json:"state,omitempty"`
	PostalCode    *string             `json:"postal_code,omitempty"`
	Country       *string             `json:"country,omitempty"`
	KycStatus     *string             `json:"kyc_status,omitempty"` // e.g., "unverified", "pending", "verified"
	GiteaOrgName  *string             `json:"gitea_org_name,omitempty"`
	Roles         []string            `json:"roles,omitempty"`
	ProjectRoles  map[string][]string `json:"project_roles,omitempty"`
}
