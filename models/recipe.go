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
	RegionName            *string         `json:"region,omitempty" db:"region_name"`
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

// DisplayRecipe represents an enriched recipe struct with live telemetry.
type DisplayRecipe struct {
	Recipe
	SyncStatus      string          `json:"sync_status"`
	HealthStatus    string          `json:"health_status"`
	BuildStatus     string          `json:"build_status"`
	BuildInfo       *BuildInfo      `json:"build_info,omitempty"`
	ProvisionStatus string          `json:"provision_status"`
	CompositeStatus string          `json:"composite_status"`
	ServingDishes   []ServingDish   `json:"serving_dishes"`
	ServingCounters []ServingCounter `json:"serving_counters"`
	Pantries        []Pantry        `json:"pantries"`
	RecipeCards     []RecipeCard    `json:"recipe_cards"`
	SecretSpices    []SecretSpice   `json:"secret_spices"`
	CustomDomains   []CustomDomain  `json:"custom_domains"`
	ResourceTree    *ResourceTree   `json:"resource_tree,omitempty"`
}

type ServingDish struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	CPU      string `json:"cpu"`
	Memory   string `json:"memory"`
	Restarts int    `json:"restarts"`
}

type ServingCounter struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Ports  string `json:"ports"`
	Public bool   `json:"public"`
}

type Pantry struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Size   string `json:"size"`
	Status string `json:"status"`
}

type RecipeCard struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
	Keys     int    `json:"keys"`
}

type SecretSpice struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Keys int    `json:"keys"`
}

type CustomDomain struct {
	ID        string `json:"id"`
	Hostname  string `json:"hostname"`
	SSLStatus string `json:"sslStatus"`
}

type ResourceTree struct {
	Nodes []ResourceNode `json:"nodes"`
}

type ResourceNode struct {
	Group      string      `json:"group"`
	Version    string      `json:"version"`
	Kind       string      `json:"kind"`
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	CreatedAt  string      `json:"createdAt,omitempty"`
	Health     *HealthInfo `json:"health,omitempty"`
	ParentRefs []ParentRef `json:"parentRefs,omitempty"`
	Info       []InfoItem  `json:"info,omitempty"`
}

type HealthInfo struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type ParentRef struct {
	Group     string `json:"group"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type InfoItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NetworkPolicyRule struct {
	Name            string   `json:"name"`
	SecurityProfile string   `json:"security_profile"`
	IngressPolicy   string   `json:"ingress_policy"`
	EgressPolicy    string   `json:"egress_policy"`
	AllowedPorts    []string `json:"allowed_ports"`
}
