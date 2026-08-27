package models

import "time"

// Standardized build status constants.
const (
	BuildStatusPending  = "pending"
	BuildStatusBuilding = "building"
	BuildStatusSuccess  = "success"
	BuildStatusFailed   = "failed"
	BuildStatusUnknown  = "unknown"
)

// Standardized provision status constants.
const (
	ProvisionStatusPending      = "pending"
	ProvisionStatusProvisioning = "provisioning"
	ProvisionStatusCompleted    = "completed"
	ProvisionStatusFailed       = "failed"
)

// Standardized composite recipe lifecycle status constants.
const (
	CompositeStatusFailed       = "failed"
	CompositeStatusProvisioning = "provisioning"
	CompositeStatusBuilding     = "building"
	CompositeStatusSyncing      = "syncing"
	CompositeStatusActive       = "active"
)

// BuildInfo contains Gitea Actions CI pipeline workflow run metadata.
type BuildInfo struct {
	WorkflowRunID int64      `json:"workflow_run_id"`
	CommitSHA     string     `json:"commit_sha"`
	Status        string     `json:"status"`
	Conclusion    string     `json:"conclusion"`
	Duration      string     `json:"duration"`
	HTMLURL       string     `json:"html_url"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}
