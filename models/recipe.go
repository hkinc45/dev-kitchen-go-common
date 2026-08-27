package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Recipe represents a recipe struct in dev-kitchen.
type Recipe struct {
	ID                    uuid.UUID       `json:"id" db:"id"`
	Name                  string          `json:"name" db:"name"`
	ProjectID             uuid.UUID       `json:"project_id" db:"project_id"`
	SourceRepo            string          `json:"source_repo" db:"source_repo"`
	SourceRepoFullname    *string         `json:"source_repo_fullname,omitempty" db:"source_repo_fullname"`
	InternalGiteaRepoURL  *string         `json:"internal_gitea_repo_url,omitempty" db:"internal_gitea_repo_url"`
	ManifestRepoURL       *string         `json:"manifest_repo_url,omitempty" db:"manifest_repo_url"`
	Status                string          `json:"status" db:"status"`
	RegionID              *uuid.UUID      `json:"region_id,omitempty" db:"region_id"`
	WorkloadType          *string         `json:"workload_type,omitempty" db:"workload_type"`
	SourceType            string          `json:"source_type" db:"source_type"`
	Configuration         json.RawMessage `json:"configuration,omitempty" db:"configuration"`
	VCSConnectionID       *uuid.UUID      `json:"vcs_connection_id,omitempty" db:"vcs_connection_id"`
	DockerfilePath        *string         `json:"dockerfile_path,omitempty" db:"dockerfile_path"`
	BuildContext          *string         `json:"build_context,omitempty" db:"build_context"`
	PathTriggersOverride  *string         `json:"path_triggers_override,omitempty" db:"path_triggers_override"`
	ContainerPort         *int            `json:"container_port,omitempty" db:"container_port"`
	AppCommand            []string        `json:"app_command,omitempty" db:"app_command"`
	AppArgs               []string        `json:"app_args,omitempty" db:"app_args"`
	Branch                *string         `json:"branch,omitempty" db:"branch"`
	SourceCommitHash      *string         `json:"source_commit_hash,omitempty" db:"source_commit_hash"`
	SourceCommitMessage   *string         `json:"source_commit_message,omitempty" db:"source_commit_message"`
	SourceCommitAuthor    *string         `json:"source_commit_author,omitempty" db:"source_commit_author"`
	SourceCommitTimestamp *time.Time      `json:"source_commit_timestamp,omitempty" db:"source_commit_timestamp"`
	GiteaCommitHash       *string         `json:"gitea_commit_hash,omitempty" db:"gitea_commit_hash"`
	GiteaCommitMessage    *string         `json:"gitea_commit_message,omitempty" db:"gitea_commit_message"`
	GiteaCommitAuthor     *string         `json:"gitea_commit_author,omitempty" db:"gitea_commit_author"`
	GiteaCommitTimestamp  *time.Time      `json:"gitea_commit_timestamp,omitempty" db:"gitea_commit_timestamp"`
	LatestCIRunID         int64           `json:"latest_ci_run_id,omitempty" db:"latest_ci_run_id"`
	LatestBuildImage      string          `json:"latest_build_image,omitempty" db:"latest_build_image"`
	LatestBuildStatus     string          `json:"latest_build_status,omitempty" db:"latest_build_status"`
	LatestBuildDuration   string          `json:"latest_build_duration,omitempty" db:"latest_build_duration"`
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at" db:"updated_at"`
}
