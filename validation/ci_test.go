package validation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsValidJobTimestamp(t *testing.T) {
	now := time.Now()
	assert.True(t, IsValidJobTimestamp(&now))

	valid2024 := time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC)
	assert.True(t, IsValidJobTimestamp(&valid2024))

	assert.False(t, IsValidJobTimestamp(nil))

	zeroTime := time.Time{}
	assert.False(t, IsValidJobTimestamp(&zeroTime))

	epochTime := time.Unix(0, 0)
	assert.False(t, IsValidJobTimestamp(&epochTime))

	oldTime := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.False(t, IsValidJobTimestamp(&oldTime))
}

func TestIsValidJobTimestampUnix(t *testing.T) {
	assert.True(t, IsValidJobTimestampUnix(time.Now().Unix()))
	assert.True(t, IsValidJobTimestampUnix(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix()))
	assert.False(t, IsValidJobTimestampUnix(0))
	assert.False(t, IsValidJobTimestampUnix(-100))
	assert.False(t, IsValidJobTimestampUnix(time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC).Unix()))
}

func TestClampCIDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected int64
	}{
		{"normal duration", 120, 120},
		{"minimum duration 1", 1, 1},
		{"zero duration clamped to 1", 0, 1},
		{"negative duration clamped to 1", -50, 1},
		{"exactly maximum 7200", 7200, 7200},
		{"excessive duration 10000 capped to 7200", 10000, 7200},
		{"huge duration 86400 capped to 7200", 86400, 7200},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, ClampCIDuration(tc.input))
		})
	}
}
