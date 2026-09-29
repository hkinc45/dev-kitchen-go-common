package models

import (
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
	AccessKey       string     `json:"access_key" db:"access_key"`
	VaultSecretPath string     `json:"vault_secret_path" db:"vault_secret_path"`
	PolicyName      string     `json:"policy_name" db:"policy_name"`
	Description     string     `json:"description" db:"description"`
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
