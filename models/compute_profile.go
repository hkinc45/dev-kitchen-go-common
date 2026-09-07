package models

import (
	"time"

	"github.com/google/uuid"
)

// ComputeProfile represents a Super Admin configurable compute capacity tier profile.
type ComputeProfile struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Slug          string    `json:"slug" db:"slug"`
	CPULimitM     int       `json:"cpu_limit_m" db:"cpu_limit_m"`
	MemoryLimitMi int       `json:"memory_limit_mi" db:"memory_limit_mi"`
	DisplayOrder  int       `json:"display_order" db:"display_order"`
	IsDefault     bool      `json:"is_default" db:"is_default"`
	IsEnabled     bool      `json:"is_enabled" db:"is_enabled"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
