package models

import (
	"encoding/json"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

func TestNormalizePhoneWithCountry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		rawPhone    string
		country     string
		wantPhone   string
		wantCountry string
	}{
		{
			name:        "USA alias with legacy hyphenated number",
			rawPhone:    "123-555-0123",
			country:     "USA",
			wantPhone:   "+1 123 555 0123",
			wantCountry: "United States",
		},
		{
			name:        "Nigeria local number with trunk 0",
			rawPhone:    "08031234567",
			country:     "Nigeria",
			wantPhone:   "+234 803 123 4567",
			wantCountry: "Nigeria",
		},
		{
			name:        "Nigeria ISO code with unprefixed 234 digits",
			rawPhone:    "2348031234567",
			country:     "NG",
			wantPhone:   "+234 803 123 4567",
			wantCountry: "Nigeria",
		},
		{
			name:        "UK alias with local trunk 0 number",
			rawPhone:    "07911 123456",
			country:     "UK",
			wantPhone:   "+44 791 112 3456",
			wantCountry: "United Kingdom",
		},
		{
			name:        "Already prefixed international phone with trunk 0 after dial code",
			rawPhone:    "+23408031234567",
			country:     "Nigeria",
			wantPhone:   "+234 803 123 4567",
			wantCountry: "Nigeria",
		},
		{
			name:        "Already normalized international phone is idempotent",
			rawPhone:    "+234 803 123 4567",
			country:     "Nigeria",
			wantPhone:   "+234 803 123 4567",
			wantCountry: "Nigeria",
		},
		{
			name:        "Empty phone with country alias still normalizes country",
			rawPhone:    "",
			country:     "US",
			wantPhone:   "",
			wantCountry: "United States",
		},
		{
			name:        "Ghana local number",
			rawPhone:    "024 123 4567",
			country:     "Ghana",
			wantPhone:   "+233 241 234 567",
			wantCountry: "Ghana",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotPhone, gotCountry := NormalizePhoneWithCountry(tt.rawPhone, tt.country)
			if gotPhone != tt.wantPhone {
				t.Errorf("NormalizePhoneWithCountry(%q, %q) phone = %q, want %q", tt.rawPhone, tt.country, gotPhone, tt.wantPhone)
			}
			if gotCountry != tt.wantCountry {
				t.Errorf("NormalizePhoneWithCountry(%q, %q) country = %q, want %q", tt.rawPhone, tt.country, gotCountry, tt.wantCountry)
			}
		})
	}
}

func TestNormalizeUserPhoneAndCountry(t *testing.T) {
	t.Parallel()

	u := &User{
		PhoneNumber: strPtr("123-555-0123"),
		Country:     strPtr("USA"),
	}
	if !NormalizeUserPhoneAndCountry(u) {
		t.Fatalf("expected NormalizeUserPhoneAndCountry to return true for legacy USA user")
	}
	if *u.PhoneNumber != "+1 123 555 0123" {
		t.Errorf("got PhoneNumber %q, want %q", *u.PhoneNumber, "+1 123 555 0123")
	}
	if *u.Country != "United States" {
		t.Errorf("got Country %q, want %q", *u.Country, "United States")
	}
	if NormalizeUserPhoneAndCountry(u) {
		t.Errorf("expected second NormalizeUserPhoneAndCountry call to be idempotent and return false")
	}
}

func TestCalculateProfileCompletion(t *testing.T) {
	t.Parallel()

	personalKYC, _ := json.Marshal(KYCProfilePayload{
		DateOfBirth:      "1994-05-12",
		Nationality:      "Nigerian",
		IDDocumentType:   "passport",
		IDDocumentNumber: "A12345678",
	})

	businessBaseKYC, _ := json.Marshal(KYCProfilePayload{
		CompanyLegalName: "Dev Kitchen Cloud Ltd",
	})

	businessFullKYC, _ := json.Marshal(KYCProfilePayload{
		CompanyLegalName:   "Dev Kitchen Cloud Ltd",
		RegistrationNumber: "RC-998877",
		TaxID:              "TIN-112233",
		Industry:           "Cloud Infrastructure",
		CompanySize:        "11-50",
		Website:            "https://devkitchen.cloud",
		AuthorizedRepRole:  "CTO",
	})

	tests := []struct {
		name              string
		user              *User
		wantPct           int
		wantMissingLen    int
		wantContainsField string
	}{
		{
			name:           "nil user returns 0%",
			user:           nil,
			wantPct:        0,
			wantMissingLen: 0,
		},
		{
			name:              "empty user returns 0% and 14 missing fields",
			user:              &User{},
			wantPct:           0,
			wantMissingLen:    14,
			wantContainsField: "account_type",
		},
		{
			name: "base-only personal user returns 50%",
			user: &User{
				AccountType:       strPtr("personal"),
				FirstName:         strPtr("Ada"),
				LastName:          strPtr("Lovelace"),
				PhoneNumber:       strPtr("+2348000000000"),
				Country:           strPtr("Nigeria"),
				PreferredCurrency: strPtr("NGN"),
			},
			wantPct:           50,
			wantMissingLen:    8,
			wantContainsField: "date_of_birth",
		},
		{
			name: "full personal user returns 100%",
			user: &User{
				AccountType:       strPtr("personal"),
				FirstName:         strPtr("Ada"),
				LastName:          strPtr("Lovelace"),
				PhoneNumber:       strPtr("+2348000000000"),
				Country:           strPtr("Nigeria"),
				PreferredCurrency: strPtr("NGN"),
				StreetAddress:     strPtr("42 Marina Road"),
				City:              strPtr("Lagos"),
				State:             strPtr("Lagos"),
				PostalCode:        strPtr("100001"),
				KycProfile:        personalKYC,
			},
			wantPct:        100,
			wantMissingLen: 0,
		},
		{
			name: "base-only business user returns 50%",
			user: &User{
				AccountType:       strPtr("business"),
				FirstName:         strPtr("Grace"),
				LastName:          strPtr("Hopper"),
				PhoneNumber:       strPtr("+2348111111111"),
				Country:           strPtr("Nigeria"),
				PreferredCurrency: strPtr("USD"),
				KycProfile:        businessBaseKYC,
			},
			wantPct:           50,
			wantMissingLen:    10,
			wantContainsField: "registration_number",
		},
		{
			name: "full business user returns 100%",
			user: &User{
				AccountType:       strPtr("business"),
				FirstName:         strPtr("Grace"),
				LastName:          strPtr("Hopper"),
				PhoneNumber:       strPtr("+2348111111111"),
				Country:           strPtr("Nigeria"),
				PreferredCurrency: strPtr("USD"),
				StreetAddress:     strPtr("10 Victoria Island"),
				City:              strPtr("Lagos"),
				State:             strPtr("Lagos"),
				PostalCode:        strPtr("101241"),
				KycProfile:        businessFullKYC,
			},
			wantPct:        100,
			wantMissingLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pct, missing := CalculateProfileCompletion(tt.user)
			if pct != tt.wantPct {
				t.Errorf("CalculateProfileCompletion() pct = %d, want %d (missing=%v)", pct, tt.wantPct, missing)
			}
			if len(missing) != tt.wantMissingLen {
				t.Errorf("CalculateProfileCompletion() len(missing) = %d (%v), want %d", len(missing), missing, tt.wantMissingLen)
			}
			if tt.user != nil && tt.user.ProfileCompletionPercentage != tt.wantPct {
				t.Errorf("user.ProfileCompletionPercentage = %d, want %d", tt.user.ProfileCompletionPercentage, tt.wantPct)
			}
			if tt.wantContainsField != "" {
				found := false
				for _, m := range missing {
					if m == tt.wantContainsField {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected missing fields %v to contain %q", missing, tt.wantContainsField)
				}
			}
		})
	}
}
