package validation

import (
	"time"
)

const (
	// MaxCIDurationSeconds defines the hard upper bound (2 hours) for any single CI build job in Dev Kitchen.
	MaxCIDurationSeconds int64 = 7200
)

// IsValidJobTimestamp validates that a job timestamp is non-nil, non-zero, and post-2020.
func IsValidJobTimestamp(t *time.Time) bool {
	return t != nil && !t.IsZero() && t.Unix() > 0 && t.Year() >= 2020
}

// IsValidJobTimestampUnix validates that a job timestamp (in unix epoch seconds) is post-2020.
func IsValidJobTimestampUnix(unixSec int64) bool {
	if unixSec <= 0 {
		return false
	}
	t := time.Unix(unixSec, 0)
	return t.Year() >= 2020
}

// ClampCIDuration ensures a raw duration in seconds is bounded between 1 and MaxCIDurationSeconds (7200).
func ClampCIDuration(durSec int64) int64 {
	if durSec > MaxCIDurationSeconds {
		return MaxCIDurationSeconds
	}
	if durSec < 1 {
		return 1
	}
	return durSec
}
