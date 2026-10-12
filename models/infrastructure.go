package models

import (
	"time"

	"github.com/google/uuid"
)

// Provider represents a physical or logical infrastructure cloud provider or colocation vendor.
type Provider struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	ProviderType string    `json:"provider_type" db:"provider_type"` // e.g. "baremetal", "colo", "cloud"
	Description  *string   `json:"description,omitempty" db:"description"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Region represents a geographical or logical datacenter region within an infrastructure provider.
type Region struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ProviderID  *uuid.UUID `json:"provider_id,omitempty" db:"provider_id"`
	Name        string     `json:"name" db:"name"`
	Code        string     `json:"code" db:"code"`
	CountryCode string     `json:"country_code" db:"country_code"`
	City        string     `json:"city" db:"city"`
	Status      string     `json:"status" db:"status"` // e.g. "active", "degraded", "maintenance"
	IsEnabled   bool       `json:"is_enabled" db:"is_enabled"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Cluster represents an operational Kubernetes cluster managed within a datacenter region.
type Cluster struct {
	ID                   uuid.UUID         `json:"id" db:"id"`
	Name                 string            `json:"name" db:"name"`
	ServerURL            string            `json:"server_url" db:"server_url"`
	RegionID             uuid.UUID         `json:"region_id" db:"region_id"`
	Labels               map[string]string `json:"labels,omitempty" db:"labels"`
	Status               string            `json:"status" db:"status"`               // e.g. "provisioning", "active", "draining", "offline"
	HealthStatus         string            `json:"health_status" db:"health_status"` // e.g. "healthy", "degraded", "unreachable"
	MaxWorkloads         int               `json:"max_workloads" db:"max_workloads"`
	CurrentWorkloadCount int               `json:"current_workload_count" db:"current_workload_count"`
	CPUCapacityM         int               `json:"cpu_capacity_m" db:"cpu_capacity_m"`
	CPUAllocatedM        int               `json:"cpu_allocated_m" db:"cpu_allocated_m"`
	MemoryCapacityMi     int               `json:"memory_capacity_mi" db:"memory_capacity_mi"`
	MemoryAllocatedMi    int               `json:"memory_allocated_mi" db:"memory_allocated_mi"`
	IsCordoned           bool              `json:"is_cordoned" db:"is_cordoned"`
	AutoscalingEnabled   bool              `json:"autoscaling_enabled" db:"autoscaling_enabled"`
	HighWatermarkPercent int               `json:"high_watermark_percent" db:"high_watermark_percent"`
	AnsibleHost          *string           `json:"ansible_host,omitempty" db:"ansible_host"`
	LastHeartbeatAt      *time.Time        `json:"last_heartbeat_at,omitempty" db:"last_heartbeat_at"`
	CreatedAt            time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at" db:"updated_at"`
}

// ClusterLoadMetrics captures calculated allocation utilization metrics and scheduling load scores.
type ClusterLoadMetrics struct {
	ClusterID          uuid.UUID `json:"cluster_id"`
	WorkloadCount      int       `json:"workload_count"`
	CPUUsagePercent    float64   `json:"cpu_usage_percent"`
	MemoryUsagePercent float64   `json:"memory_usage_percent"`
	Score              float64   `json:"score"`
}
