package models

import (
	"encoding/json"
	"math"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// countryDialEntry maps an ISO 3166-1 alpha-2 code to its canonical display name and ITU-T E.164 dial code.
type countryDialEntry struct {
	ISO      string
	Name     string
	DialCode string
}

var allCountryDialEntries = []countryDialEntry{
	{ISO: "AD", Name: "Andorra", DialCode: "+376"},
	{ISO: "AE", Name: "United Arab Emirates", DialCode: "+971"},
	{ISO: "AF", Name: "Afghanistan", DialCode: "+93"},
	{ISO: "AG", Name: "Antigua & Barbuda", DialCode: "+1"},
	{ISO: "AI", Name: "Anguilla", DialCode: "+1"},
	{ISO: "AL", Name: "Albania", DialCode: "+355"},
	{ISO: "AM", Name: "Armenia", DialCode: "+374"},
	{ISO: "AO", Name: "Angola", DialCode: "+244"},
	{ISO: "AR", Name: "Argentina", DialCode: "+54"},
	{ISO: "AS", Name: "American Samoa", DialCode: "+1"},
	{ISO: "AT", Name: "Austria", DialCode: "+43"},
	{ISO: "AU", Name: "Australia", DialCode: "+61"},
	{ISO: "AW", Name: "Aruba", DialCode: "+297"},
	{ISO: "AZ", Name: "Azerbaijan", DialCode: "+994"},
	{ISO: "BA", Name: "Bosnia & Herzegovina", DialCode: "+387"},
	{ISO: "BB", Name: "Barbados", DialCode: "+1"},
	{ISO: "BD", Name: "Bangladesh", DialCode: "+880"},
	{ISO: "BE", Name: "Belgium", DialCode: "+32"},
	{ISO: "BF", Name: "Burkina Faso", DialCode: "+226"},
	{ISO: "BG", Name: "Bulgaria", DialCode: "+359"},
	{ISO: "BH", Name: "Bahrain", DialCode: "+973"},
	{ISO: "BI", Name: "Burundi", DialCode: "+257"},
	{ISO: "BJ", Name: "Benin", DialCode: "+229"},
	{ISO: "BM", Name: "Bermuda", DialCode: "+1"},
	{ISO: "BN", Name: "Brunei", DialCode: "+673"},
	{ISO: "BO", Name: "Bolivia", DialCode: "+591"},
	{ISO: "BR", Name: "Brazil", DialCode: "+55"},
	{ISO: "BS", Name: "Bahamas", DialCode: "+1"},
	{ISO: "BT", Name: "Bhutan", DialCode: "+975"},
	{ISO: "BW", Name: "Botswana", DialCode: "+267"},
	{ISO: "BY", Name: "Belarus", DialCode: "+375"},
	{ISO: "BZ", Name: "Belize", DialCode: "+501"},
	{ISO: "CA", Name: "Canada", DialCode: "+1"},
	{ISO: "CD", Name: "Congo - Kinshasa", DialCode: "+243"},
	{ISO: "CF", Name: "Central African Republic", DialCode: "+236"},
	{ISO: "CG", Name: "Congo - Brazzaville", DialCode: "+242"},
	{ISO: "CH", Name: "Switzerland", DialCode: "+41"},
	{ISO: "CI", Name: "Côte d’Ivoire", DialCode: "+225"},
	{ISO: "CK", Name: "Cook Islands", DialCode: "+682"},
	{ISO: "CL", Name: "Chile", DialCode: "+56"},
	{ISO: "CM", Name: "Cameroon", DialCode: "+237"},
	{ISO: "CN", Name: "China", DialCode: "+86"},
	{ISO: "CO", Name: "Colombia", DialCode: "+57"},
	{ISO: "CR", Name: "Costa Rica", DialCode: "+506"},
	{ISO: "CU", Name: "Cuba", DialCode: "+53"},
	{ISO: "CV", Name: "Cape Verde", DialCode: "+238"},
	{ISO: "CW", Name: "Curaçao", DialCode: "+599"},
	{ISO: "CY", Name: "Cyprus", DialCode: "+357"},
	{ISO: "CZ", Name: "Czechia", DialCode: "+420"},
	{ISO: "DE", Name: "Germany", DialCode: "+49"},
	{ISO: "DJ", Name: "Djibouti", DialCode: "+253"},
	{ISO: "DK", Name: "Denmark", DialCode: "+45"},
	{ISO: "DM", Name: "Dominica", DialCode: "+1"},
	{ISO: "DO", Name: "Dominican Republic", DialCode: "+1"},
	{ISO: "DZ", Name: "Algeria", DialCode: "+213"},
	{ISO: "EC", Name: "Ecuador", DialCode: "+593"},
	{ISO: "EE", Name: "Estonia", DialCode: "+372"},
	{ISO: "EG", Name: "Egypt", DialCode: "+20"},
	{ISO: "ER", Name: "Eritrea", DialCode: "+291"},
	{ISO: "ES", Name: "Spain", DialCode: "+34"},
	{ISO: "ET", Name: "Ethiopia", DialCode: "+251"},
	{ISO: "FI", Name: "Finland", DialCode: "+358"},
	{ISO: "FJ", Name: "Fiji", DialCode: "+679"},
	{ISO: "FK", Name: "Falkland Islands", DialCode: "+500"},
	{ISO: "FM", Name: "Micronesia", DialCode: "+691"},
	{ISO: "FO", Name: "Faroe Islands", DialCode: "+298"},
	{ISO: "FR", Name: "France", DialCode: "+33"},
	{ISO: "GA", Name: "Gabon", DialCode: "+241"},
	{ISO: "GB", Name: "United Kingdom", DialCode: "+44"},
	{ISO: "GD", Name: "Grenada", DialCode: "+1"},
	{ISO: "GE", Name: "Georgia", DialCode: "+995"},
	{ISO: "GF", Name: "French Guiana", DialCode: "+594"},
	{ISO: "GG", Name: "Guernsey", DialCode: "+44"},
	{ISO: "GH", Name: "Ghana", DialCode: "+233"},
	{ISO: "GI", Name: "Gibraltar", DialCode: "+350"},
	{ISO: "GL", Name: "Greenland", DialCode: "+299"},
	{ISO: "GM", Name: "Gambia", DialCode: "+220"},
	{ISO: "GN", Name: "Guinea", DialCode: "+224"},
	{ISO: "GP", Name: "Guadeloupe", DialCode: "+590"},
	{ISO: "GQ", Name: "Equatorial Guinea", DialCode: "+240"},
	{ISO: "GR", Name: "Greece", DialCode: "+30"},
	{ISO: "GT", Name: "Guatemala", DialCode: "+502"},
	{ISO: "GU", Name: "Guam", DialCode: "+1"},
	{ISO: "GW", Name: "Guinea-Bissau", DialCode: "+245"},
	{ISO: "GY", Name: "Guyana", DialCode: "+592"},
	{ISO: "HK", Name: "Hong Kong SAR China", DialCode: "+852"},
	{ISO: "HN", Name: "Honduras", DialCode: "+504"},
	{ISO: "HR", Name: "Croatia", DialCode: "+385"},
	{ISO: "HT", Name: "Haiti", DialCode: "+509"},
	{ISO: "HU", Name: "Hungary", DialCode: "+36"},
	{ISO: "ID", Name: "Indonesia", DialCode: "+62"},
	{ISO: "IE", Name: "Ireland", DialCode: "+353"},
	{ISO: "IL", Name: "Israel", DialCode: "+972"},
	{ISO: "IM", Name: "Isle of Man", DialCode: "+44"},
	{ISO: "IN", Name: "India", DialCode: "+91"},
	{ISO: "IO", Name: "British Indian Ocean Territory", DialCode: "+246"},
	{ISO: "IQ", Name: "Iraq", DialCode: "+964"},
	{ISO: "IR", Name: "Iran", DialCode: "+98"},
	{ISO: "IS", Name: "Iceland", DialCode: "+354"},
	{ISO: "IT", Name: "Italy", DialCode: "+39"},
	{ISO: "JE", Name: "Jersey", DialCode: "+44"},
	{ISO: "JM", Name: "Jamaica", DialCode: "+1"},
	{ISO: "JO", Name: "Jordan", DialCode: "+962"},
	{ISO: "JP", Name: "Japan", DialCode: "+81"},
	{ISO: "KE", Name: "Kenya", DialCode: "+254"},
	{ISO: "KG", Name: "Kyrgyzstan", DialCode: "+996"},
	{ISO: "KH", Name: "Cambodia", DialCode: "+855"},
	{ISO: "KI", Name: "Kiribati", DialCode: "+686"},
	{ISO: "KM", Name: "Comoros", DialCode: "+269"},
	{ISO: "KN", Name: "St. Kitts & Nevis", DialCode: "+1"},
	{ISO: "KP", Name: "North Korea", DialCode: "+850"},
	{ISO: "KR", Name: "South Korea", DialCode: "+82"},
	{ISO: "KW", Name: "Kuwait", DialCode: "+965"},
	{ISO: "KY", Name: "Cayman Islands", DialCode: "+1"},
	{ISO: "KZ", Name: "Kazakhstan", DialCode: "+7"},
	{ISO: "LA", Name: "Laos", DialCode: "+856"},
	{ISO: "LB", Name: "Lebanon", DialCode: "+961"},
	{ISO: "LC", Name: "St. Lucia", DialCode: "+1"},
	{ISO: "LI", Name: "Liechtenstein", DialCode: "+423"},
	{ISO: "LK", Name: "Sri Lanka", DialCode: "+94"},
	{ISO: "LR", Name: "Liberia", DialCode: "+231"},
	{ISO: "LS", Name: "Lesotho", DialCode: "+266"},
	{ISO: "LT", Name: "Lithuania", DialCode: "+370"},
	{ISO: "LU", Name: "Luxembourg", DialCode: "+352"},
	{ISO: "LV", Name: "Latvia", DialCode: "+371"},
	{ISO: "LY", Name: "Libya", DialCode: "+218"},
	{ISO: "MA", Name: "Morocco", DialCode: "+212"},
	{ISO: "MC", Name: "Monaco", DialCode: "+377"},
	{ISO: "MD", Name: "Moldova", DialCode: "+373"},
	{ISO: "ME", Name: "Montenegro", DialCode: "+382"},
	{ISO: "MG", Name: "Madagascar", DialCode: "+261"},
	{ISO: "MH", Name: "Marshall Islands", DialCode: "+692"},
	{ISO: "MK", Name: "North Macedonia", DialCode: "+389"},
	{ISO: "ML", Name: "Mali", DialCode: "+223"},
	{ISO: "MM", Name: "Myanmar (Burma)", DialCode: "+95"},
	{ISO: "MN", Name: "Mongolia", DialCode: "+976"},
	{ISO: "MO", Name: "Macao SAR China", DialCode: "+853"},
	{ISO: "MP", Name: "Northern Mariana Islands", DialCode: "+1"},
	{ISO: "MQ", Name: "Martinique", DialCode: "+596"},
	{ISO: "MR", Name: "Mauritania", DialCode: "+222"},
	{ISO: "MS", Name: "Montserrat", DialCode: "+1"},
	{ISO: "MT", Name: "Malta", DialCode: "+356"},
	{ISO: "MU", Name: "Mauritius", DialCode: "+230"},
	{ISO: "MV", Name: "Maldives", DialCode: "+960"},
	{ISO: "MW", Name: "Malawi", DialCode: "+265"},
	{ISO: "MX", Name: "Mexico", DialCode: "+52"},
	{ISO: "MY", Name: "Malaysia", DialCode: "+60"},
	{ISO: "MZ", Name: "Mozambique", DialCode: "+258"},
	{ISO: "NA", Name: "Namibia", DialCode: "+264"},
	{ISO: "NC", Name: "New Caledonia", DialCode: "+687"},
	{ISO: "NE", Name: "Niger", DialCode: "+227"},
	{ISO: "NG", Name: "Nigeria", DialCode: "+234"},
	{ISO: "NI", Name: "Nicaragua", DialCode: "+505"},
	{ISO: "NL", Name: "Netherlands", DialCode: "+31"},
	{ISO: "NO", Name: "Norway", DialCode: "+47"},
	{ISO: "NP", Name: "Nepal", DialCode: "+977"},
	{ISO: "NR", Name: "Nauru", DialCode: "+674"},
	{ISO: "NU", Name: "Niue", DialCode: "+683"},
	{ISO: "NZ", Name: "New Zealand", DialCode: "+64"},
	{ISO: "OM", Name: "Oman", DialCode: "+968"},
	{ISO: "PA", Name: "Panama", DialCode: "+507"},
	{ISO: "PE", Name: "Peru", DialCode: "+51"},
	{ISO: "PF", Name: "French Polynesia", DialCode: "+689"},
	{ISO: "PG", Name: "Papua New Guinea", DialCode: "+675"},
	{ISO: "PH", Name: "Philippines", DialCode: "+63"},
	{ISO: "PK", Name: "Pakistan", DialCode: "+92"},
	{ISO: "PL", Name: "Poland", DialCode: "+48"},
	{ISO: "PR", Name: "Puerto Rico", DialCode: "+1"},
	{ISO: "PS", Name: "Palestinian Territories", DialCode: "+970"},
	{ISO: "PT", Name: "Portugal", DialCode: "+351"},
	{ISO: "PW", Name: "Palau", DialCode: "+680"},
	{ISO: "PY", Name: "Paraguay", DialCode: "+595"},
	{ISO: "QA", Name: "Qatar", DialCode: "+974"},
	{ISO: "RE", Name: "Réunion", DialCode: "+262"},
	{ISO: "RO", Name: "Romania", DialCode: "+40"},
	{ISO: "RS", Name: "Serbia", DialCode: "+381"},
	{ISO: "RU", Name: "Russia", DialCode: "+7"},
	{ISO: "RW", Name: "Rwanda", DialCode: "+250"},
	{ISO: "SA", Name: "Saudi Arabia", DialCode: "+966"},
	{ISO: "SB", Name: "Solomon Islands", DialCode: "+677"},
	{ISO: "SC", Name: "Seychelles", DialCode: "+248"},
	{ISO: "SD", Name: "Sudan", DialCode: "+249"},
	{ISO: "SE", Name: "Sweden", DialCode: "+46"},
	{ISO: "SG", Name: "Singapore", DialCode: "+65"},
	{ISO: "SH", Name: "St. Helena", DialCode: "+290"},
	{ISO: "SI", Name: "Slovenia", DialCode: "+386"},
	{ISO: "SK", Name: "Slovakia", DialCode: "+421"},
	{ISO: "SL", Name: "Sierra Leone", DialCode: "+232"},
	{ISO: "SM", Name: "San Marino", DialCode: "+378"},
	{ISO: "SN", Name: "Senegal", DialCode: "+221"},
	{ISO: "SO", Name: "Somalia", DialCode: "+252"},
	{ISO: "SR", Name: "Suriname", DialCode: "+597"},
	{ISO: "SS", Name: "South Sudan", DialCode: "+211"},
	{ISO: "ST", Name: "São Tomé & Príncipe", DialCode: "+239"},
	{ISO: "SV", Name: "El Salvador", DialCode: "+503"},
	{ISO: "SX", Name: "Sint Maarten", DialCode: "+1"},
	{ISO: "SY", Name: "Syria", DialCode: "+963"},
	{ISO: "SZ", Name: "Eswatini", DialCode: "+268"},
	{ISO: "TC", Name: "Turks & Caicos Islands", DialCode: "+1"},
	{ISO: "TD", Name: "Chad", DialCode: "+235"},
	{ISO: "TG", Name: "Togo", DialCode: "+228"},
	{ISO: "TH", Name: "Thailand", DialCode: "+66"},
	{ISO: "TJ", Name: "Tajikistan", DialCode: "+992"},
	{ISO: "TL", Name: "Timor-Leste", DialCode: "+670"},
	{ISO: "TM", Name: "Turkmenistan", DialCode: "+993"},
	{ISO: "TN", Name: "Tunisia", DialCode: "+216"},
	{ISO: "TO", Name: "Tonga", DialCode: "+676"},
	{ISO: "TR", Name: "Türkiye", DialCode: "+90"},
	{ISO: "TT", Name: "Trinidad & Tobago", DialCode: "+1"},
	{ISO: "TV", Name: "Tuvalu", DialCode: "+688"},
	{ISO: "TW", Name: "Taiwan", DialCode: "+886"},
	{ISO: "TZ", Name: "Tanzania", DialCode: "+255"},
	{ISO: "UA", Name: "Ukraine", DialCode: "+380"},
	{ISO: "UG", Name: "Uganda", DialCode: "+256"},
	{ISO: "US", Name: "United States", DialCode: "+1"},
	{ISO: "UY", Name: "Uruguay", DialCode: "+598"},
	{ISO: "UZ", Name: "Uzbekistan", DialCode: "+998"},
	{ISO: "VA", Name: "Vatican City", DialCode: "+379"},
	{ISO: "VC", Name: "St. Vincent & Grenadines", DialCode: "+1"},
	{ISO: "VE", Name: "Venezuela", DialCode: "+58"},
	{ISO: "VG", Name: "British Virgin Islands", DialCode: "+1"},
	{ISO: "VI", Name: "U.S. Virgin Islands", DialCode: "+1"},
	{ISO: "VN", Name: "Vietnam", DialCode: "+84"},
	{ISO: "VU", Name: "Vanuatu", DialCode: "+678"},
	{ISO: "WS", Name: "Samoa", DialCode: "+685"},
	{ISO: "YE", Name: "Yemen", DialCode: "+967"},
	{ISO: "YT", Name: "Mayotte", DialCode: "+262"},
	{ISO: "ZA", Name: "South Africa", DialCode: "+27"},
	{ISO: "ZM", Name: "Zambia", DialCode: "+260"},
	{ISO: "ZW", Name: "Zimbabwe", DialCode: "+263"},
}

var (
	isoToCountryDial  map[string]countryDialEntry
	nameToCountryDial map[string]countryDialEntry
	countryAliases    map[string]string
	sortedDialCodes   []string
)

func init() {
	isoToCountryDial = make(map[string]countryDialEntry, len(allCountryDialEntries))
	nameToCountryDial = make(map[string]countryDialEntry, len(allCountryDialEntries))
	dialSet := make(map[string]struct{}, len(allCountryDialEntries))

	for _, entry := range allCountryDialEntries {
		isoToCountryDial[strings.ToUpper(entry.ISO)] = entry
		nameToCountryDial[strings.ToLower(entry.Name)] = entry
		dialSet[entry.DialCode] = struct{}{}
	}

	countryAliases = map[string]string{
		"usa":                      "United States",
		"us":                       "United States",
		"u.s.":                     "United States",
		"u.s.a.":                   "United States",
		"united states of america": "United States",
		"uk":                       "United Kingdom",
		"u.k.":                     "United Kingdom",
		"great britain":            "United Kingdom",
		"england":                  "United Kingdom",
		"uae":                      "United Arab Emirates",
		"ivory coast":              "Côte d’Ivoire",
		"cote d'ivoire":            "Côte d’Ivoire",
		"turkey":                   "Türkiye",
		"czech republic":           "Czechia",
		"drc":                      "Congo - Kinshasa",
	}

	sortedDialCodes = make([]string, 0, len(dialSet))
	for code := range dialSet {
		sortedDialCodes = append(sortedDialCodes, code)
	}
	sort.Slice(sortedDialCodes, func(i, j int) bool {
		if len(sortedDialCodes[i]) != len(sortedDialCodes[j]) {
			return len(sortedDialCodes[i]) > len(sortedDialCodes[j])
		}
		return sortedDialCodes[i] < sortedDialCodes[j]
	})
}

// ResolveCountryAndDialCode resolves a country name, ISO-3166 alpha-2 code, or common alias
// into its canonical country name (matching the frontend COUNTRIES list) and ITU-T E.164 dial code.
func ResolveCountryAndDialCode(country string) (canonicalCountry string, dialCode string) {
	trimmed := strings.TrimSpace(country)
	if trimmed == "" {
		return "", "+1"
	}

	lower := strings.ToLower(trimmed)
	if aliasTarget, ok := countryAliases[lower]; ok {
		if entry, found := nameToCountryDial[strings.ToLower(aliasTarget)]; found {
			return entry.Name, entry.DialCode
		}
	}

	upper := strings.ToUpper(trimmed)
	if len(upper) == 2 {
		if entry, ok := isoToCountryDial[upper]; ok {
			return entry.Name, entry.DialCode
		}
	}

	if entry, ok := nameToCountryDial[lower]; ok {
		return entry.Name, entry.DialCode
	}

	return trimmed, "+1"
}

func extractDigitsOnly(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func formatLocalPhoneDigits(rawDigits string) string {
	digits := extractDigitsOnly(rawDigits)
	if len(digits) > 15 {
		digits = digits[:15]
	}
	if digits == "" {
		return ""
	}
	if len(digits) <= 3 {
		return digits
	}
	if len(digits) <= 6 {
		return digits[:3] + " " + digits[3:]
	}
	if len(digits) <= 10 {
		return digits[:3] + " " + digits[3:6] + " " + digits[6:]
	}
	return digits[:3] + " " + digits[3:7] + " " + digits[7:]
}

// SplitAndNormalizePhone separates and normalizes a phone country dial code ("+234"),
// a clean local subscriber phone number ("8031234567" — digits only, without dial code or trunk 0),
// and a canonical country name ("Nigeria").
func SplitAndNormalizePhone(rawPhoneCode, rawPhone, rawCountry string) (dialCode string, localNumber string, canonicalCountry string) {
	canonicalCountry, countryDialCode := ResolveCountryAndDialCode(rawCountry)

	cleanedCodeDigits := extractDigitsOnly(rawPhoneCode)
	if cleanedCodeDigits != "" {
		dialCode = "+" + cleanedCodeDigits
	}

	trimmedPhone := strings.TrimSpace(rawPhone)
	if trimmedPhone == "" {
		if dialCode == "" && strings.TrimSpace(rawCountry) != "" {
			dialCode = countryDialCode
		}
		return dialCode, "", canonicalCountry
	}

	compact := strings.Map(func(r rune) rune {
		switch r {
		case ' ', '-', '(', ')', '.':
			return -1
		default:
			return r
		}
	}, trimmedPhone)

	if strings.HasPrefix(compact, "+") {
		if dialCode != "" && strings.HasPrefix(compact, dialCode) {
			localNumber = extractDigitsOnly(compact[len(dialCode):])
		} else {
			matchedDial := ""
			for _, code := range sortedDialCodes {
				if strings.HasPrefix(compact, code) {
					matchedDial = code
					break
				}
			}
			if matchedDial != "" {
				if dialCode == "" {
					dialCode = matchedDial
				}
				localNumber = extractDigitsOnly(compact[len(matchedDial):])
			} else {
				rest := compact[1:]
				prefixLen := 0
				for prefixLen < len(rest) && prefixLen < 3 && unicode.IsDigit(rune(rest[prefixLen])) {
					prefixLen++
				}
				if prefixLen > 0 {
					if dialCode == "" {
						dialCode = "+" + rest[:prefixLen]
					}
					localNumber = extractDigitsOnly(rest[prefixLen:])
				} else {
					localNumber = extractDigitsOnly(rest)
				}
			}
		}
		if strings.HasPrefix(localNumber, "0") && (dialCode != "+1" || len(localNumber) > 10) {
			localNumber = strings.TrimLeft(localNumber, "0")
		}
	} else {
		digitsOnly := extractDigitsOnly(compact)
		if dialCode == "" {
			dialCode = countryDialCode
		}
		codeDigits := strings.TrimPrefix(dialCode, "+")
		if codeDigits != "" && len(digitsOnly) >= len(codeDigits)+10 && strings.HasPrefix(digitsOnly, codeDigits) {
			localNumber = digitsOnly[len(codeDigits):]
			if strings.HasPrefix(localNumber, "0") && (dialCode != "+1" || len(localNumber) > 10) {
				localNumber = strings.TrimLeft(localNumber, "0")
			}
		} else if strings.HasPrefix(digitsOnly, "0") && (dialCode != "+1" || len(digitsOnly) > 10) {
			localNumber = strings.TrimLeft(digitsOnly, "0")
		} else {
			localNumber = digitsOnly
		}
	}

	if len(localNumber) > 14 {
		localNumber = localNumber[:14]
	}
	if localNumber == "" {
		return dialCode, "", canonicalCountry
	}

	return dialCode, localNumber, canonicalCountry
}

// FormatE164Phone combines a dial code (e.g. "+234") and local subscriber digits (e.g. "8031234567")
// into a compact ITU-T E.164 string ("+2348031234567") for Keycloak attributes or SMS gateways.
func FormatE164Phone(dialCode, localNumber string) string {
	digits := extractDigitsOnly(localNumber)
	if digits == "" {
		return ""
	}
	codeDigits := extractDigitsOnly(dialCode)
	if codeDigits == "" {
		return "+" + digits
	}
	return "+" + codeDigits + digits
}

// NormalizePhoneWithCountry normalizes a raw phone string and country string into an international
// E.164-prefixed display format ("+<dial_code> <grouped_digits>") and canonical country name.
func NormalizePhoneWithCountry(rawPhone string, country string) (normalizedPhone string, canonicalCountry string) {
	dialCode, localDigits, canonicalCountry := SplitAndNormalizePhone("", rawPhone, country)
	if localDigits == "" {
		return "", canonicalCountry
	}
	formattedLocal := formatLocalPhoneDigits(localDigits)
	return dialCode + " " + formattedLocal, canonicalCountry
}

// NormalizeUserPhoneAndCountry normalizes u.PhoneCountryCode ("+234"), u.PhoneNumber ("8031234567" — clean
// local subscriber digits), u.Country ("Nigeria" — canonical name), and KycProfile.AuthorizedRepPhoneCountryCode /
// AuthorizedRepPhone in-place. Returns true if any field on u was modified.
func NormalizeUserPhoneAndCountry(u *User) bool {
	if u == nil {
		return false
	}

	rawCountry := ""
	if u.Country != nil {
		rawCountry = *u.Country
	}
	rawPhoneCode := ""
	if u.PhoneCountryCode != nil {
		rawPhoneCode = *u.PhoneCountryCode
	}
	rawPhone := ""
	if u.PhoneNumber != nil {
		rawPhone = *u.PhoneNumber
	}

	dialCode, localNumber, normCountry := SplitAndNormalizePhone(rawPhoneCode, rawPhone, rawCountry)
	changed := false

	if hasNonEmptyPtr(u.Country) && normCountry != "" && *u.Country != normCountry {
		u.Country = &normCountry
		changed = true
	}

	if hasNonEmptyPtr(u.PhoneNumber) {
		if localNumber != "" && *u.PhoneNumber != localNumber {
			u.PhoneNumber = &localNumber
			changed = true
		}
		if dialCode != "" && (u.PhoneCountryCode == nil || *u.PhoneCountryCode != dialCode) {
			u.PhoneCountryCode = &dialCode
			changed = true
		}
	} else if hasNonEmptyPtr(u.PhoneCountryCode) && dialCode != "" && *u.PhoneCountryCode != dialCode {
		u.PhoneCountryCode = &dialCode
		changed = true
	}

	if len(u.KycProfile) > 0 {
		var kycMap map[string]interface{}
		if err := json.Unmarshal(u.KycProfile, &kycMap); err == nil && len(kycMap) > 0 {
			repPhone, _ := kycMap["authorized_rep_phone"].(string)
			repCode, _ := kycMap["authorized_rep_phone_country_code"].(string)
			if strings.TrimSpace(repPhone) != "" || strings.TrimSpace(repCode) != "" {
				repCountry, _ := kycMap["country_of_incorporation"].(string)
				if strings.TrimSpace(repCountry) == "" {
					repCountry = normCountry
				}
				normRepCode, normRepPhone, _ := SplitAndNormalizePhone(repCode, repPhone, repCountry)
				kycUpdated := false
				if normRepPhone != "" && repPhone != normRepPhone {
					kycMap["authorized_rep_phone"] = normRepPhone
					kycUpdated = true
				}
				if normRepCode != "" && repCode != normRepCode {
					kycMap["authorized_rep_phone_country_code"] = normRepCode
					kycUpdated = true
				}
				if kycUpdated {
					if b, err := json.Marshal(kycMap); err == nil {
						u.KycProfile = json.RawMessage(b)
						changed = true
					}
				}
			}
		}
	}

	return changed
}

// KYCProfilePayload represents the structured KYC/KYB JSONB payload stored on a user profile.
type KYCProfilePayload struct {
	// Personal KYC fields
	DateOfBirth      string `json:"date_of_birth,omitempty"`
	Nationality      string `json:"nationality,omitempty"`
	IDDocumentType   string `json:"id_document_type,omitempty"`
	IDDocumentNumber string `json:"id_document_number,omitempty"`

	// Business KYB fields
	CompanyLegalName              string `json:"company_legal_name,omitempty"`
	RegistrationNumber            string `json:"registration_number,omitempty"`
	TaxID                         string `json:"tax_id,omitempty"`
	CountryOfIncorporation        string `json:"country_of_incorporation,omitempty"`
	Industry                      string `json:"industry,omitempty"`
	CompanySize                   string `json:"company_size,omitempty"`
	Website                       string `json:"website,omitempty"`
	AuthorizedRepName             string `json:"authorized_rep_name,omitempty"`
	AuthorizedRepRole             string `json:"authorized_rep_role,omitempty"`
	AuthorizedRepEmail            string `json:"authorized_rep_email,omitempty"`
	AuthorizedRepPhoneCountryCode string `json:"authorized_rep_phone_country_code,omitempty" db:"authorized_rep_phone_country_code"`
	AuthorizedRepPhone            string `json:"authorized_rep_phone,omitempty"`
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
	PhoneCountryCode            *string             `json:"phone_country_code,omitempty" db:"phone_country_code"`
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
// Also normalizes u.PhoneCountryCode, u.PhoneNumber, and u.Country in-place via NormalizeUserPhoneAndCountry.
func CalculateProfileCompletion(u *User) (int, []string) {
	if u == nil {
		return 0, nil
	}

	NormalizeUserPhoneAndCountry(u)

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
