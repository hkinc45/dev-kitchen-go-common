package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateProjectNameFormat(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"valid standard", "my-project", false},
		{"valid alphanumeric", "project-123", false},
		{"valid minimum length 3", "abc", false},
		{"valid maximum length 35", "a123456789012345678901234567890123b", false},
		{"empty string", "", true},
		{"too short", "ab", true},
		{"too long (> 35)", "a12345678901234567890123456789012345c", true},
		{"contains underscore", "my_project", true},
		{"uppercase characters", "MyProject", true},
		{"consecutive hyphens", "my--project", true},
		{"leading hyphen", "-project", true},
		{"trailing hyphen", "project-", true},
		{"special symbols", "project@name", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateProjectNameFormat(tc.input)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRecipeNameFormat(t *testing.T) {
	assert.NoError(t, ValidateRecipeNameFormat("web-app-1"))
	assert.Error(t, ValidateRecipeNameFormat("a1234567890123456789012345678901")) // 31 chars > 30
	assert.Error(t, ValidateRecipeNameFormat("web_app"))
}

func TestSanitizeRFC1123(t *testing.T) {
	assert.Equal(t, "user-project", SanitizeRFC1123("User_Project!", 50, "default"))
	assert.Equal(t, "my-app", SanitizeRFC1123("---My---App---", 50, "default"))
	assert.Equal(t, "default", SanitizeRFC1123("---!@#---", 50, "default"))
	assert.Equal(t, "testuser", SanitizeRFC1123("testuser_very_long_name", 8, "default"))
}

func TestSanitizeDefaultProjectName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"john_doe", "john-doe"},
		{"John DOE!!", "john-doe"},
		{"a_very_long_username_that_exceeds_limits", "a-very-long-username-that-exc"},
		{"---", "default-project"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			res := SanitizeDefaultProjectName(tc.input)
			assert.Equal(t, tc.expected, res)
			assert.LessOrEqual(t, len(res), 29)
			assert.NoError(t, ValidateRFC1123(res, 3, 35, "project"))
		})
	}
}

func TestGenerateDefaultProjectName(t *testing.T) {
	tests := []struct {
		userID   string
		expected string
	}{
		{"11b31c16-bb1a-445f-a8f0-4da700e5d352", "project-11b31c16"},
		{"550e8400-e29b-41d4-a716-446655440000", "project-550e8400"},
		{"a1b2c3d4", "project-a1b2c3d4"},
		{"abc", "project-abc"},
		{"", "project-default"},
		{"---", "project-default"},
	}

	for _, tc := range tests {
		t.Run(tc.userID, func(t *testing.T) {
			res := GenerateDefaultProjectName(tc.userID)
			assert.Equal(t, tc.expected, res)
			assert.LessOrEqual(t, len(res), 29)
			assert.NoError(t, ValidateProjectNameFormat(res))
		})
	}
}

