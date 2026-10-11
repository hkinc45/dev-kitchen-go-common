package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// ComputeProfile represents a Super Admin configurable compute capacity tier profile.
type ComputeProfile struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Slug          string    `json:"slug" db:"slug"`
	SKU           string    `json:"sku" db:"sku"`
	CPULimitM     int       `json:"cpu_limit_m" db:"cpu_limit_m"`
	MemoryLimitMi int       `json:"memory_limit_mi" db:"memory_limit_mi"`
	DisplayOrder  int       `json:"display_order" db:"display_order"`
	IsDefault     bool      `json:"is_default" db:"is_default"`
	IsEnabled     bool      `json:"is_enabled" db:"is_enabled"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// InferWorkloadProfileSKU extracts or infers the canonical recipe profile SKU
// from single pod CPU millicores and RAM bytes.
func InferWorkloadProfileSKU(singlePodCPU, singlePodRAMBytes int64) string {
	// If resources are completely unspecified (0/0) or non-positive, default to the standard 1cpu-1gb general-purpose tier
	if singlePodCPU <= 0 && singlePodRAMBytes <= 0 {
		return "recipe.profile.1cpu-1gb"
	}

	memMiB := singlePodRAMBytes / (1024 * 1024)
	if singlePodCPU <= 500 && memMiB <= 512 {
		return "recipe.profile.0.5cpu-512mb"
	} else if singlePodCPU <= 1000 && memMiB <= 1024 {
		return "recipe.profile.1cpu-1gb"
	} else if singlePodCPU <= 1000 && memMiB <= 2048 {
		return "recipe.profile.1cpu-2gb"
	} else if singlePodCPU <= 2000 && memMiB <= 4096 {
		return "recipe.profile.2cpu-4gb"
	}
	return "recipe.profile.4cpu-8gb"
}

// ResolveAllocationMultiplier returns the billable multiplier for a given allocation type string.
func ResolveAllocationMultiplier(allocType string) float64 {
	switch strings.ToLower(strings.TrimSpace(allocType)) {
	case "lightweight":
		return 0.5
	case "performance":
		return 2.0
	default:
		return 1.0
	}
}
