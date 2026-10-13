package models

import (
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

// KYCProfilePayload represents the structured KYC/KYB JSONB payload stored on a user profile.
type KYCProfilePayload struct {
	// Personal KYC fields
	DateOfBirth      string `json:"date_of_birth,omitempty"`
	Nationality      string `json:"nationality,omitempty"`
	IDDocumentType   string `json:"id_document_type,omitempty"`
	IDDocumentNumber string `json:"id_document_number,omitempty"`

	// Business KYB fields
	CompanyLegalName   string `json:"company_legal_name,omitempty"`
	RegistrationNumber string `json:"registration_number,omitempty"`
	TaxID              string `json:"tax_id,omitempty"`
	Industry           string `json:"industry,omitempty"`
	CompanySize        string `json:"company_size,omitempty"`
	Website            string `json:"website,omitempty"`
	AuthorizedRepRole  string `json:"authorized_rep_role,omitempty"`
}

// NotificationPreferencesPayload represents the structured notification settings stored on a user profile.
type NotificationPreferencesPayload struct {
	EmailDeployments bool `json:"email_deployments"`
	EmailBilling     bool `json:"email_billing"`
	EmailSecurity    bool `json:"email_security"`
	MarketingUpdates bool `json:"marketing_updates"`
}

// User represents a user in the application's database.
// It is the canonical representation of a user across all services.
type User struct {
	ID                          uuid.UUID           `json:"id"`
	KeycloakID                  string              `json:"keycloak_id"`
	Username                    string              `json:"username"`
	Email                       string              `json:"email"`
	CreatedAt                   time.Time           `json:"created_at"`
	UpdatedAt                   time.Time           `json:"updated_at"`
	FirstName                   *string             `json:"first_name,omitempty"`
	LastName                    *string             `json:"last_name,omitempty"`
	PhoneNumber                 *string             `json:"phone_number,omitempty"`
	AccountType                 *string             `json:"account_type,omitempty"` // e.g., "personal" or "business"
	PreferredCurrency           *string             `json:"preferred_currency,omitempty" db:"preferred_currency"`
	StreetAddress               *string             `json:"street_address,omitempty"`
	City                        *string             `json:"city,omitempty"`
	State                       *string             `json:"state,omitempty"`
	PostalCode                  *string             `json:"postal_code,omitempty"`
	Country                     *string             `json:"country,omitempty"`
	KycStatus                   *string             `json:"kyc_status,omitempty"` // e.g., "unverified", "pending", "verified", "rejected"
	KycProfile                  json.RawMessage     `json:"kyc_profile,omitempty" db:"kyc_profile"`
	NotificationPreferences     json.RawMessage     `json:"notification_preferences,omitempty" db:"notification_preferences"`
	KycSubmittedAt              *time.Time          `json:"kyc_submitted_at,omitempty" db:"kyc_submitted_at"`
	ProfileCompletionPercentage int                 `json:"profile_completion_percentage" db:"-"`
	MissingProfileFields        []string            `json:"missing_profile_fields,omitempty" db:"-"`
	GiteaOrgName                *string             `json:"gitea_org_name,omitempty"`
	Roles                       []string            `json:"roles,omitempty"`
	ProjectRoles                map[string][]string `json:"project_roles,omitempty"`
	ProjectPermissions          json.RawMessage     `json:"project_permissions,omitempty" db:"project_permissions"`
	SyncStatus                  string              `json:"sync_status,omitempty" db:"sync_status"`
	SyncLockedAt                *time.Time          `json:"sync_locked_at,omitempty" db:"sync_locked_at"`
}

func hasNonEmptyPtr(s *string) bool {
	return s != nil && strings.TrimSpace(*s) != ""
}

func hasNonEmptyStr(s string) bool {
	return strings.TrimSpace(s) != ""
}

// CalculateProfileCompletion computes a deterministic 0-100% completion score and ordered list of missing fields.
// Base Mandatory Fields account for 50% of the score; Extended Address & KYC/KYB Fields account for 50%.
func CalculateProfileCompletion(u *User) (int, []string) {
	if u == nil {
		return 0, nil
	}

	var kyc KYCProfilePayload
	if len(u.KycProfile) > 0 {
		_ = json.Unmarshal(u.KycProfile, &kyc)
	}

	acctType := "personal"
	if hasNonEmptyPtr(u.AccountType) {
		acctType = strings.ToLower(strings.TrimSpace(*u.AccountType))
	}

	type fieldCheck struct {
		name    string
		present bool
	}

	baseChecks := []fieldCheck{
		{"account_type", hasNonEmptyPtr(u.AccountType)},
		{"first_name", hasNonEmptyPtr(u.FirstName)},
		{"last_name", hasNonEmptyPtr(u.LastName)},
		{"phone_number", hasNonEmptyPtr(u.PhoneNumber)},
		{"country", hasNonEmptyPtr(u.Country)},
		{"preferred_currency", hasNonEmptyPtr(u.PreferredCurrency)},
	}
	if acctType == "business" {
		baseChecks = append(baseChecks, fieldCheck{"company_legal_name", hasNonEmptyStr(kyc.CompanyLegalName)})
	}

	extendedChecks := []fieldCheck{
		{"street_address", hasNonEmptyPtr(u.StreetAddress)},
		{"city", hasNonEmptyPtr(u.City)},
		{"state", hasNonEmptyPtr(u.State)},
		{"postal_code", hasNonEmptyPtr(u.PostalCode)},
	}
	if acctType == "business" {
		extendedChecks = append(extendedChecks,
			fieldCheck{"registration_number", hasNonEmptyStr(kyc.RegistrationNumber)},
			fieldCheck{"tax_id", hasNonEmptyStr(kyc.TaxID)},
			fieldCheck{"industry", hasNonEmptyStr(kyc.Industry)},
			fieldCheck{"company_size", hasNonEmptyStr(kyc.CompanySize)},
			fieldCheck{"website", hasNonEmptyStr(kyc.Website)},
			fieldCheck{"authorized_rep_role", hasNonEmptyStr(kyc.AuthorizedRepRole)},
		)
	} else {
		extendedChecks = append(extendedChecks,
			fieldCheck{"date_of_birth", hasNonEmptyStr(kyc.DateOfBirth)},
			fieldCheck{"nationality", hasNonEmptyStr(kyc.Nationality)},
			fieldCheck{"id_document_type", hasNonEmptyStr(kyc.IDDocumentType)},
			fieldCheck{"id_document_number", hasNonEmptyStr(kyc.IDDocumentNumber)},
		)
	}

	var missing []string
	completedBase := 0
	for _, c := range baseChecks {
		if c.present {
			completedBase++
		} else {
			missing = append(missing, c.name)
		}
	}

	completedExtended := 0
	for _, c := range extendedChecks {
		if c.present {
			completedExtended++
		} else {
			missing = append(missing, c.name)
		}
	}

	baseScore := 50.0 * (float64(completedBase) / float64(len(baseChecks)))
	extScore := 50.0 * (float64(completedExtended) / float64(len(extendedChecks)))
	pct := int(math.Round(baseScore + extScore))

	u.ProfileCompletionPercentage = pct
	u.MissingProfileFields = missing
	return pct, missing
}

