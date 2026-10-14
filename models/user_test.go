package models

import (
	"encoding/json"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

func TestSplitAndNormalizePhone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		rawPhoneCode    string
		rawPhone        string
		rawCountry      string
		wantDialCode    string
		wantLocalNumber string
		wantCountry     string
		wantE164        string
	}{
		{
			name:            "explicit PhoneCountryCode +234 and local number with trunk 0",
			rawPhoneCode:    "+234",
			rawPhone:        "08031234567",
			rawCountry:      "Nigeria",
			wantDialCode:    "+234",
			wantLocalNumber: "8031234567",
			wantCountry:     "Nigeria",
			wantE164:        "+2348031234567",
		},
		{
			name:            "legacy combined +234 803 123 4567 with NG ISO code",
			rawPhoneCode:    "",
			rawPhone:        "+234 803 123 4567",
			rawCountry:      "NG",
			wantDialCode:    "+234",
			wantLocalNumber: "8031234567",
			wantCountry:     "Nigeria",
			wantE164:        "+2348031234567",
		},
		{
			name:            "cross-border user in United Kingdom with explicit +234 phone country code",
			rawPhoneCode:    "+234",
			rawPhone:        "8031234567",
			rawCountry:      "United Kingdom",
			wantDialCode:    "+234",
			wantLocalNumber: "8031234567",
			wantCountry:     "United Kingdom",
			wantE164:        "+2348031234567",
		},
		{
			name:            "NANP +1 10-digit subscriber number preserves all 10 digits",
			rawPhoneCode:    "+1",
			rawPhone:        "4165550199",
			rawCountry:      "Canada",
			wantDialCode:    "+1",
			wantLocalNumber: "4165550199",
			wantCountry:     "Canada",
			wantE164:        "+14165550199",
		},
		{
			name:            "NANP +1 10-digit legacy USA number starting with 1 preserves all 10 digits",
			rawPhoneCode:    "",
			rawPhone:        "123-555-0123",
			rawCountry:      "USA",
			wantDialCode:    "+1",
			wantLocalNumber: "1235550123",
			wantCountry:     "United States",
			wantE164:        "+11235550123",
		},
		{
			name:            "NANP +1 11-digit number with leading country code 1 strips leading 1",
			rawPhoneCode:    "+1",
			rawPhone:        "14165550199",
			rawCountry:      "Canada",
			wantDialCode:    "+1",
			wantLocalNumber: "4165550199",
			wantCountry:     "Canada",
			wantE164:        "+14165550199",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotDial, gotLocal, gotCountry := SplitAndNormalizePhone(tt.rawPhoneCode, tt.rawPhone, tt.rawCountry)
			if gotDial != tt.wantDialCode {
				t.Errorf("SplitAndNormalizePhone() dialCode = %q, want %q", gotDial, tt.wantDialCode)
			}
			if gotLocal != tt.wantLocalNumber {
				t.Errorf("SplitAndNormalizePhone() localNumber = %q, want %q", gotLocal, tt.wantLocalNumber)
			}
			if gotCountry != tt.wantCountry {
				t.Errorf("SplitAndNormalizePhone() canonicalCountry = %q, want %q", gotCountry, tt.wantCountry)
			}
			gotE164 := FormatE164Phone(gotDial, gotLocal)
			if gotE164 != tt.wantE164 {
				t.Errorf("FormatE164Phone() = %q, want %q", gotE164, tt.wantE164)
			}
		})
	}
}

func TestNormalizeUserPhoneAndCountry(t *testing.T) {
	t.Parallel()

	// 1. Legacy USA user with unprefixed phone
	u := &User{
		PhoneNumber: strPtr("123-555-0123"),
		Country:     strPtr("USA"),
		KycProfile:  json.RawMessage(`{"authorized_rep_phone":"+234 0803 123 4567","country_of_incorporation":"NG"}`),
	}
	if !NormalizeUserPhoneAndCountry(u) {
		t.Fatalf("expected NormalizeUserPhoneAndCountry to return true for legacy USA user")
	}
	if u.PhoneCountryCode == nil || *u.PhoneCountryCode != "+1" {
		t.Errorf("got PhoneCountryCode %v, want +1", u.PhoneCountryCode)
	}
	if u.PhoneNumber == nil || *u.PhoneNumber != "1235550123" {
		t.Errorf("got PhoneNumber %v, want 1235550123", u.PhoneNumber)
	}
	if u.Country == nil || *u.Country != "United States" {
		t.Errorf("got Country %v, want United States", u.Country)
	}

	var kyc KYCProfilePayload
	if err := json.Unmarshal(u.KycProfile, &kyc); err != nil {
		t.Fatalf("failed to unmarshal KycProfile: %v", err)
	}
	if kyc.AuthorizedRepPhoneCountryCode != "+234" {
		t.Errorf("got AuthorizedRepPhoneCountryCode %q, want +234", kyc.AuthorizedRepPhoneCountryCode)
	}
	if kyc.AuthorizedRepPhone != "8031234567" {
		t.Errorf("got AuthorizedRepPhone %q, want 8031234567", kyc.AuthorizedRepPhone)
	}

	if NormalizeUserPhoneAndCountry(u) {
		t.Errorf("expected second NormalizeUserPhoneAndCountry call to be idempotent and return false")
	}

	// 2. Cross-border user in United Kingdom with explicit +234 PhoneCountryCode
	crossBorder := &User{
		PhoneCountryCode: strPtr("+234"),
		PhoneNumber:      strPtr("08031234567"),
		Country:          strPtr("United Kingdom"),
	}
	if !NormalizeUserPhoneAndCountry(crossBorder) {
		t.Fatalf("expected NormalizeUserPhoneAndCountry to strip trunk 0 on cross-border phone")
	}
	if *crossBorder.Country != "United Kingdom" {
		t.Errorf("got Country %q, want United Kingdom", *crossBorder.Country)
	}
	if *crossBorder.PhoneCountryCode != "+234" {
		t.Errorf("got PhoneCountryCode %q, want +234", *crossBorder.PhoneCountryCode)
	}
	if *crossBorder.PhoneNumber != "8031234567" {
		t.Errorf("got PhoneNumber %q, want 8031234567", *crossBorder.PhoneNumber)
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
