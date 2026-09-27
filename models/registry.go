package models

import (
	"time"

	"github.com/google/uuid"
)

// ContainerRegistry represents an enterprise container registry configuration.
type ContainerRegistry struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Name            string     `json:"name" db:"name"`
	Type            string     `json:"type" db:"type"` // harbor, zot, ecr, dockerhub, ghcr
	URL             string     `json:"url" db:"url"`
	IsDefault       bool       `json:"is_default" db:"is_default"`
	Scope           string     `json:"scope" db:"scope"` // global, project
	ProjectID       *uuid.UUID `json:"project_id,omitempty" db:"project_id"`
	VaultSecretPath *string    `json:"vault_secret_path,omitempty" db:"vault_secret_path"`
	RegionID        *uuid.UUID `json:"region_id,omitempty" db:"region_id"`
	Status          string     `json:"status" db:"status"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// ProjectRegistryBinding maps a project to a specific container registry.
type ProjectRegistryBinding struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ProjectID  uuid.UUID `json:"project_id" db:"project_id"`
	RegistryID uuid.UUID `json:"registry_id" db:"registry_id"`
	IsPrimary  bool      `json:"is_primary" db:"is_primary"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// RobotAccount represents a project-scoped machine robot account.
type RobotAccount struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	ProjectID   string     `json:"project_id"`
	Scope       string     `json:"scope"`
	Permissions string     `json:"permissions"` // pull, push_pull
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Disabled    bool       `json:"disabled"`
}

// CreateRobotAccountRequest payload for provisioning a project robot account.
type CreateRobotAccountRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	DurationDays int    `json:"duration_days"`
	Scope        string `json:"scope"` // pull or push_pull
}

// RegistryTokenResponse returned upon creating a robot account, showing secret once.
type RegistryTokenResponse struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
	Secret             string     `json:"secret"`
	DockerLoginCommand string     `json:"docker_login_command"`
	ExpiresAt          *time.Time `json:"expires_at,omitempty"`
}

// RegistryRepository represents a container repository under a project namespace.
type RegistryRepository struct {
	Name          string    `json:"name"`
	ProjectID     string    `json:"project_id"`
	ArtifactCount int       `json:"artifact_count"`
	PullCount     int       `json:"pull_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RegistryArtifact represents an image artifact digest with tags, metadata, and status.
type RegistryArtifact struct {
	Digest          string                `json:"digest"`
	Tags            []string              `json:"tags"`
	Size            int64                 `json:"size"`
	PushTime        time.Time             `json:"push_time"`
	PullTime        *time.Time            `json:"pull_time,omitempty"`
	Status          string                `json:"status"` // active, past, dangling
	RecipeID        *uuid.UUID            `json:"recipe_id,omitempty"`
	RecipeName      *string               `json:"recipe_name,omitempty"`
	CommitHash      *string               `json:"commit_hash,omitempty"`
	CommitMessage   *string               `json:"commit_message,omitempty"`
	CommitAuthor    *string               `json:"commit_author,omitempty"`
	Vulnerabilities *VulnerabilitySummary `json:"vulnerabilities,omitempty"`
}

// VulnerabilitySummary summarizes security scan metrics.
type VulnerabilitySummary struct {
	Total     int        `json:"total"`
	Critical  int        `json:"critical"`
	High      int        `json:"high"`
	Medium    int        `json:"medium"`
	Low       int        `json:"low"`
	Fixable   int        `json:"fixable"`
	ScannedAt *time.Time `json:"scanned_at,omitempty"`
}
