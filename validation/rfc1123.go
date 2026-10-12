package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// RFC1123Regex validates standard Kubernetes RFC 1123 DNS-1123 label format:
	// - lowercase alphanumeric characters or '-'
	// - must start and end with an alphanumeric character
	RFC1123Regex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
)

// ValidateRFC1123 validates that a given string strictly adheres to RFC 1123 DNS label constraints.
func ValidateRFC1123(name string, minLen, maxLen int, resourceKind string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("%s name is required", resourceKind)
	}
	if len(name) < minLen {
		return fmt.Errorf("%s name must be at least %d characters long", resourceKind, minLen)
	}
	if len(name) > maxLen {
		return fmt.Errorf("%s name cannot exceed %d characters", resourceKind, maxLen)
	}
	if strings.Contains(name, "_") {
		return fmt.Errorf("%s name cannot contain underscores; use hyphens instead (e.g. 'my-%s')", resourceKind, strings.ToLower(resourceKind))
	}
	if strings.ToLower(name) != name {
		return fmt.Errorf("%s name must be lowercase; uppercase letters are not allowed", resourceKind)
	}
	if strings.Contains(name, "--") {
		return fmt.Errorf("%s name cannot contain consecutive hyphens", resourceKind)
	}
	if !RFC1123Regex.MatchString(name) {
		return fmt.Errorf("%s name must start and end with an alphanumeric character and contain only lowercase letters, numbers, and hyphens", resourceKind)
	}
	return nil
}

// ValidateProjectNameFormat validates that a project name conforms to Dev Kitchen DNS-1035/1123 rules (3-35 chars).
func ValidateProjectNameFormat(name string) error {
	return ValidateRFC1123(name, 3, 35, "project")
}

// ValidateRecipeNameFormat validates that a recipe name conforms to Dev Kitchen DNS-1035/1123 rules (3-30 chars).
func ValidateRecipeNameFormat(name string) error {
	return ValidateRFC1123(name, 3, 30, "recipe")
}

// SanitizeRFC1123 cleans any raw string into a valid RFC 1123 string with a given maximum length ceiling.
func SanitizeRFC1123(rawName string, maxLen int, fallback string) string {
	s := strings.ToLower(strings.TrimSpace(rawName))
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	s = b.String()
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	if maxLen > 0 && len(s) > maxLen {
		s = strings.TrimRight(s[:maxLen], "-")
	}
	if s == "" {
		return fallback
	}
	return s
}

// SanitizeDefaultProjectName cleans an onboarding username into a valid RFC 1123 project prefix
// capped at 29 characters to leave room for conflict retry suffixes.
// Preserves compatibility with legacy "<username>-default-project" naming patterns.
func SanitizeDefaultProjectName(rawName string) string {
	return SanitizeRFC1123(rawName, 29, "default-project")
}

// GenerateDefaultProjectName produces a collision-free RFC 1123 compliant default project name
// following the pattern "project-<uuid[:8]>" (e.g. "project-11b31c16", exactly 16 characters).
// If rawUserID is empty or lacks alphanumeric characters, it safely defaults to "project-default".
func GenerateDefaultProjectName(rawUserID string) string {
	clean := strings.ReplaceAll(strings.TrimSpace(rawUserID), "-", "")
	clean = strings.ToLower(clean)
	if len(clean) >= 8 {
		clean = clean[:8]
	}
	clean = SanitizeRFC1123(clean, 8, "default")
	return fmt.Sprintf("project-%s", clean)
}
