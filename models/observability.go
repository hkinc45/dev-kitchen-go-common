package models

import "time"

// Severity levels for deployment issue events.
const (
	SeverityWarning  = "Warning"
	SeverityError    = "Error"
	SeverityCritical = "Critical"
)

// DeploymentIssueSource defines the origin of the deployment issue.
const (
	SourceArgoCD     = "ArgoCD"
	SourceKubernetes = "Kubernetes"
)

// DeploymentIssueEvent represents a diagnostic warning or sync failure
// from ArgoCD or Kubernetes warning events.
type DeploymentIssueEvent struct {
	Source         string    `json:"source"`          // "ArgoCD" | "Kubernetes"
	Severity       string    `json:"severity"`        // "Warning" | "Error" | "Critical"
	Reason         string    `json:"reason"`          // e.g. "CrashLoopBackOff", "ImagePullBackOff", "ComparisonError", "SyncError", "OOMKilled"
	Message        string    `json:"message"`         // Actionable error message
	InvolvedObject string    `json:"involved_object"` // Pod or Resource name
	Timestamp      time.Time `json:"timestamp"`
	Count          int32     `json:"count"`
}

// ContainerMetricTelemetry captures real-time CPU and memory usage and limits
// measured directly from the Kubernetes Metrics Server.
type ContainerMetricTelemetry struct {
	ContainerName        string    `json:"container_name"`
	CPUUsageMillicores   int64     `json:"cpu_usage_millicores"`
	CPURequestMillicores int64     `json:"cpu_request_millicores"`
	CPULimitMillicores   int64     `json:"cpu_limit_millicores"`
	MemoryUsageBytes     int64     `json:"memory_usage_bytes"`
	MemoryUsageMiB       float64   `json:"memory_usage_mib"`
	MemoryRequestMiB     float64   `json:"memory_request_mib"`
	MemoryLimitMiB       float64   `json:"memory_limit_mib"`
	Timestamp            time.Time `json:"timestamp"`
}

// CIBuildLogResponse represents raw CI build pipeline log output from Gitea Actions.
type CIBuildLogResponse struct {
	RunID      int64      `json:"run_id"`
	Status     string     `json:"status"`
	Conclusion string     `json:"conclusion"`
	Logs       string     `json:"logs"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	Duration   string     `json:"duration,omitempty"`
}

// LogStreamRequest defines parameters for streaming Kubernetes Pod logs via SSE.
type LogStreamRequest struct {
	Container  string `form:"container" json:"container"`
	TailLines  int64  `form:"tail_lines" json:"tail_lines"`
	Follow     bool   `form:"follow" json:"follow"`
	Timestamps bool   `form:"timestamps" json:"timestamps"`
	Previous   bool   `form:"previous" json:"previous"`
}
