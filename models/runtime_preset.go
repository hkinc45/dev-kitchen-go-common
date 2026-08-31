package models

import (
	"time"

	"github.com/google/uuid"
)

// RuntimePreset represents a configurable language runtime version & framework preset for Super Admin management.
type RuntimePreset struct {
	ID                     uuid.UUID `json:"id" db:"id"`
	Language               string    `json:"language" db:"language"`
	Version                string    `json:"version" db:"version"`
	IsEnabled              bool      `json:"is_enabled" db:"is_enabled"`
	IsDefault              bool      `json:"is_default" db:"is_default"`
	DefaultPackageManagers []string  `json:"default_package_managers,omitempty" db:"default_package_managers"`
	DefaultInstallCommand  *string   `json:"default_install_command,omitempty" db:"default_install_command"`
	DefaultBuildCommand    *string   `json:"default_build_command,omitempty" db:"default_build_command"`
	DefaultStartCommand    *string   `json:"default_start_command,omitempty" db:"default_start_command"`
	DefaultContainerPort   *int      `json:"default_container_port,omitempty" db:"default_container_port"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}
