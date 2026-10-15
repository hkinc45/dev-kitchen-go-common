package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Constants for object storage
const (
	ComponentTypeObjectStorage = "object_storage"

	BucketPolicyPrivate    = "private"
	BucketPolicyPublicRead = "public-read"

	BucketStatusActive   = "active"
	BucketStatusCreating = "creating"
	BucketStatusDeleting = "deleting"
	BucketStatusError    = "error"
)

// Standard storage billing SKUs
const (
	StorageSKUBlockStandard   = "storage.block.gb"
	StorageSKUBlockReplicated = "storage.block.replicated.gb"
	StorageSKUNFS             = "storage.nfs.gb"
	StorageSKUEphemeral       = "storage.ephemeral.gb"
	StorageSKUBackup          = "storage.backup.gb"
	StorageSKUPVCDefault      = "storage.pvc.gb_month"
)

// Bucket represents an S3-compliant object storage bucket.
type Bucket struct {
	ID                uuid.UUID `json:"id" db:"id"`
	ProjectID         uuid.UUID `json:"project_id" db:"project_id"`
	Name              string    `json:"name" db:"name"`
	CanonicalName     string    `json:"canonical_name" db:"canonical_name"`
	Region            string    `json:"region" db:"region"`
	Status            string    `json:"status" db:"status"`
	Policy            string    `json:"policy" db:"policy"`
	CORSEnabled       bool      `json:"cors_enabled" db:"cors_enabled"`
	VersioningEnabled bool      `json:"versioning_enabled" db:"versioning_enabled"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// BucketQuota defines storage capacity and object count limits for a bucket.
type BucketQuota struct {
	BucketID         uuid.UUID `json:"bucket_id" db:"bucket_id"`
	MaxSizeBytes     int64     `json:"max_size_bytes" db:"max_size_bytes"`
	MaxObjects       int64     `json:"max_objects" db:"max_objects"`
	CurrentSizeBytes int64     `json:"current_size_bytes" db:"current_size_bytes"`
	CurrentObjects   int64     `json:"current_objects" db:"current_objects"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// StorageServiceAccount represents a project-scoped S3 credential set.
type StorageServiceAccount struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	ProjectID       uuid.UUID  `json:"project_id" db:"project_id"`
	Name            string     `json:"name" db:"name"`
	AccessKey       string     `json:"access_key" db:"access_key"`
	VaultSecretPath string     `json:"vault_secret_path" db:"vault_secret_path"`
	PolicyName      string     `json:"policy_name" db:"policy_name"`
	Description     string     `json:"description" db:"description"`
	AllowedBuckets  []string   `json:"allowed_buckets" db:"allowed_buckets"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty" db:"expires_at"`
}

// PresignedURLRequest encapsulates parameters needed to generate a presigned URL.
type PresignedURLRequest struct {
	Key           string `json:"key"`
	Operation     string `json:"operation,omitempty"` // "upload" / "download" or "PUT" / "GET"
	ExpiresIn     int64  `json:"expires_in,omitempty"` // seconds
	ContentType   string `json:"content_type,omitempty"`
	ContentLength int64  `json:"content_length,omitempty"`
}

// PresignedURLResponse contains the generated presigned URL and metadata.
type PresignedURLResponse struct {
	URL       string `json:"url"`
	Method    string `json:"method"`
	Key       string `json:"key"`
	ExpiresIn int64  `json:"expires_in"`
}

// ObjectItem represents an object stored in a bucket for web file explorer and listing.
type ObjectItem struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
	ContentType  string    `json:"content_type"`
	StorageClass string    `json:"storage_class"`
}

// StorageDriverType identifies the underlying Kubernetes CSI / block storage provider.
type StorageDriverType string

const (
	StorageDriverLonghorn    StorageDriverType = "longhorn"
	StorageDriverCSISnapshot StorageDriverType = "csi-snapshot"
	StorageDriverStandard    StorageDriverType = "standard"
)

// BackupScheduleOption describes an available recurring backup/snapshot schedule tier.
type BackupScheduleOption struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Cron   string `json:"cron"`
	Retain int    `json:"retain"`
	Task   string `json:"task"`
}

// VolumeDataSource represents the CSI dataSource used to populate a PVC.
type VolumeDataSource struct {
	APIGroup string `json:"api_group,omitempty"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
}

// VolumeRecord represents a normalized view of a project PVC and its backup schedule.
type VolumeRecord struct {
	Name                  string            `json:"name"`
	Namespace             string            `json:"namespace"`
	ProjectID             string            `json:"project_id"`
	RecipeID              string            `json:"recipe_id,omitempty"`
	RecipeName            string            `json:"recipe_name,omitempty"`
	StorageClass          string            `json:"storage_class"`
	StorageSKU            string            `json:"storage_sku,omitempty"`
	ReplicationFactor     int               `json:"replication_factor,omitempty"`
	DriverType            StorageDriverType `json:"driver_type"`
	CapacityGB            int               `json:"capacity_gb"`
	SizeGiB               int               `json:"size_gib"`
	Capacity              string            `json:"capacity"`
	Status                string            `json:"status"`
	Phase                 string            `json:"phase"`
	VolumeName            string            `json:"volume_name,omitempty"`
	AccessModes           []string          `json:"access_modes"`
	BackupSchedule        string            `json:"backup_schedule"`
	SupportsSnapshots     bool              `json:"supports_snapshots"`
	SupportsRemoteBackups bool              `json:"supports_remote_backups"`
	FromBackup            string            `json:"from_backup,omitempty"`
	DataSource            *VolumeDataSource `json:"data_source,omitempty"`
	CreatedAt             time.Time         `json:"created_at"`
}

// SnapshotRecord represents a normalized CSI VolumeSnapshot or Longhorn Backup.
type SnapshotRecord struct {
	ID                  string            `json:"id"`
	Name                string            `json:"name"`
	Namespace           string            `json:"namespace,omitempty"`
	VolumeName          string            `json:"volume_name,omitempty"`
	PVCName             string            `json:"pvc_name"`
	ProjectID           string            `json:"project_id"`
	RecipeID            string            `json:"recipe_id,omitempty"`
	RecipeName          string            `json:"recipe_name,omitempty"`
	Type                string            `json:"type"` // "snapshot" or "backup"
	DriverType          StorageDriverType `json:"driver_type"`
	VolumeSnapshotClass string            `json:"volume_snapshot_class,omitempty"`
	SizeBytes           int64             `json:"size_bytes"`
	SizeGiB             int               `json:"size_gib"`
	RestoreSize         string            `json:"restore_size,omitempty"`
	Status              string            `json:"status"` // "ready", "pending", "error"
	ReadyToUse          bool              `json:"ready_to_use"`
	BackupURL           string            `json:"backup_url,omitempty"`
	CreatedAt           time.Time         `json:"created_at"`
}

// UpdateBackupPolicyRequest defines the payload to update a PVC's recurring backup schedule.
type UpdateBackupPolicyRequest struct {
	Schedule string `json:"schedule"` // "none", "hourly", "daily", "weekly"
	Retain   int    `json:"retain,omitempty"`
	Task     string    `json:"task,omitempty"`
}

// CreateSnapshotRequest defines the payload to create a CSI VolumeSnapshot or backup.
type CreateSnapshotRequest struct {
	PVCName             string            `json:"pvc_name"`
	Name                string            `json:"name,omitempty"`
	Type                string            `json:"type,omitempty"` // "snapshot" (default) or "backup"
	VolumeSnapshotClass string            `json:"volume_snapshot_class,omitempty"`
	RecipeName          string            `json:"recipe_name,omitempty"`
	Labels              map[string]string `json:"labels,omitempty"`
}

// CreateVolumeFromSnapshotRequest defines the payload to provision a new PVC from a snapshot, clone, or remote backup.
type CreateVolumeFromSnapshotRequest struct {
	Name               string `json:"name"`
	SourceSnapshotName string `json:"source_snapshot_name,omitempty"`
	SnapshotName       string `json:"snapshot_name,omitempty"`
	SourcePVCName      string `json:"source_pvc_name,omitempty"`
	SourceKind         string `json:"source_kind,omitempty"`
	FromBackupURL      string `json:"from_backup_url,omitempty"`
	StorageClass       string `json:"storage_class,omitempty"`
	StorageSKU         string `json:"storage_sku,omitempty"`
	ReplicationFactor  int    `json:"replication_factor,omitempty"`
	CapacityGB         int    `json:"capacity_gb,omitempty"`
	SizeGiB            int    `json:"size_gib,omitempty"`
	AccessMode         string `json:"access_mode,omitempty"`
	BackupSchedule     string `json:"backup_schedule,omitempty"`
	RecipeID           string `json:"recipe_id,omitempty"`
	RecipeName         string `json:"recipe_name,omitempty"`
}

// ResolveStorageReplicationFactor determines the effective replica count (1..5) for a storage volume
// from an explicit replica count or by inspecting the Kubernetes StorageClass name.
func ResolveStorageReplicationFactor(storageClass string, explicitReplicas int) int {
	if explicitReplicas >= 1 && explicitReplicas <= 5 {
		return explicitReplicas
	}
	sc := strings.ToLower(strings.TrimSpace(storageClass))
	switch {
	case strings.Contains(sc, "3x") || strings.Contains(sc, "ha-3") || strings.Contains(sc, "triple"):
		return 3
	case strings.Contains(sc, "2x") || strings.Contains(sc, "ha-2") || strings.Contains(sc, "replicated") || strings.Contains(sc, "rwx"):
		return 2
	default:
		return 1
	}
}

// ResolveStorageTierSKU determines the canonical billing SKU for a volume based on an explicit SKU,
// the storage class name, and access modes.
func ResolveStorageTierSKU(storageClass string, explicitSKU string, accessModes []string) string {
	sku := strings.TrimSpace(explicitSKU)
	if sku != "" {
		return sku
	}
	sc := strings.ToLower(strings.TrimSpace(storageClass))
	if sc == "emptydir" || strings.Contains(sc, "ephemeral") || strings.Contains(sc, "scratch") {
		return StorageSKUEphemeral
	}
	if strings.Contains(sc, "rwx") || strings.Contains(sc, "nfs") {
		return StorageSKUNFS
	}
	for _, m := range accessModes {
		if strings.EqualFold(strings.TrimSpace(m), "ReadWriteMany") || strings.EqualFold(strings.TrimSpace(m), "RWX") {
			return StorageSKUNFS
		}
	}
	if strings.Contains(sc, "replicated") || strings.Contains(sc, "ha-2x") || strings.Contains(sc, "ha-3x") || strings.Contains(sc, "2x") || strings.Contains(sc, "3x") {
		return StorageSKUBlockReplicated
	}
	return StorageSKUBlockStandard
}

// ResolveStorageClassForReplication maps a base storage class and desired replica count
// to the concrete cluster StorageClass name.
func ResolveStorageClassForReplication(baseClass string, replicas int, replicaClassMap map[string]string) string {
	base := strings.TrimSpace(baseClass)
	if base == "" {
		base = "longhorn"
	}
	if replicas < 1 {
		replicas = 1
	}
	if len(replicaClassMap) > 0 {
		if mapped, ok := replicaClassMap[strconv.Itoa(replicas)]; ok && strings.TrimSpace(mapped) != "" {
			return strings.TrimSpace(mapped)
		}
		if mapped, ok := replicaClassMap[fmt.Sprintf("%dx", replicas)]; ok && strings.TrimSpace(mapped) != "" {
			return strings.TrimSpace(mapped)
		}
	}
	lower := strings.ToLower(base)
	if strings.Contains(lower, "longhorn") {
		if strings.Contains(lower, "rwx") {
			return "longhorn-rwx"
		}
		switch replicas {
		case 3:
			return "longhorn-ha-3x"
		case 2:
			return "longhorn-replicated"
		case 1:
			if lower == "longhorn-single" {
				return "longhorn-single"
			}
			return "longhorn"
		}
	}
	return base
}


