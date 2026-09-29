package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageModelsSerialization(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	bucketID := uuid.New()
	projectID := uuid.New()

	t.Run("Bucket JSON serialization", func(t *testing.T) {
		bucket := Bucket{
			ID:                bucketID,
			ProjectID:         projectID,
			Name:              "test-bucket",
			CanonicalName:     "dk-" + projectID.String() + "-test-bucket",
			Region:            "us-east-1",
			Status:            BucketStatusActive,
			Policy:            BucketPolicyPrivate,
			CORSEnabled:       true,
			VersioningEnabled: false,
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		data, err := json.Marshal(bucket)
		require.NoError(t, err)

		var unmarshaled Bucket
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, bucket.ID, unmarshaled.ID)
		assert.Equal(t, bucket.ProjectID, unmarshaled.ProjectID)
		assert.Equal(t, bucket.Name, unmarshaled.Name)
		assert.Equal(t, bucket.CanonicalName, unmarshaled.CanonicalName)
		assert.Equal(t, bucket.Region, unmarshaled.Region)
		assert.Equal(t, bucket.Status, unmarshaled.Status)
		assert.Equal(t, bucket.Policy, unmarshaled.Policy)
		assert.Equal(t, bucket.CORSEnabled, unmarshaled.CORSEnabled)
		assert.Equal(t, bucket.VersioningEnabled, unmarshaled.VersioningEnabled)
	})

	t.Run("BucketQuota JSON serialization", func(t *testing.T) {
		quota := BucketQuota{
			BucketID:         bucketID,
			MaxSizeBytes:     10 * 1024 * 1024 * 1024, // 10 GB
			MaxObjects:       100000,
			CurrentSizeBytes: 1024 * 1024,
			CurrentObjects:   42,
			UpdatedAt:        now,
		}

		data, err := json.Marshal(quota)
		require.NoError(t, err)

		var unmarshaled BucketQuota
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, quota.BucketID, unmarshaled.BucketID)
		assert.Equal(t, quota.MaxSizeBytes, unmarshaled.MaxSizeBytes)
		assert.Equal(t, quota.MaxObjects, unmarshaled.MaxObjects)
		assert.Equal(t, quota.CurrentSizeBytes, unmarshaled.CurrentSizeBytes)
		assert.Equal(t, quota.CurrentObjects, unmarshaled.CurrentObjects)
	})

	t.Run("StorageServiceAccount JSON serialization", func(t *testing.T) {
		accountID := uuid.New()
		expiresAt := now.Add(24 * time.Hour)
		sa := StorageServiceAccount{
			ID:              accountID,
			ProjectID:       projectID,
			AccessKey:       "DKSA1234567890",
			VaultSecretPath: "projects/" + projectID.String() + "/storage/credentials/DKSA1234567890",
			PolicyName:      "project-storage-policy",
			Description:     "Service account for backup jobs",
			CreatedAt:       now,
			ExpiresAt:       &expiresAt,
		}

		data, err := json.Marshal(sa)
		require.NoError(t, err)

		var unmarshaled StorageServiceAccount
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, sa.ID, unmarshaled.ID)
		assert.Equal(t, sa.ProjectID, unmarshaled.ProjectID)
		assert.Equal(t, sa.AccessKey, unmarshaled.AccessKey)
		assert.Equal(t, sa.VaultSecretPath, unmarshaled.VaultSecretPath)
		assert.Equal(t, sa.PolicyName, unmarshaled.PolicyName)
		assert.Equal(t, sa.Description, unmarshaled.Description)
	})

	t.Run("Presigned URLs and ObjectItem serialization", func(t *testing.T) {
		req := PresignedURLRequest{
			Key:           "uploads/data.csv",
			Operation:     "upload",
			ExpiresIn:     900,
			ContentType:   "text/csv",
			ContentLength: 2048,
		}
		reqBytes, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaledReq PresignedURLRequest
		err = json.Unmarshal(reqBytes, &unmarshaledReq)
		require.NoError(t, err)
		assert.Equal(t, "uploads/data.csv", unmarshaledReq.Key)
		assert.Equal(t, int64(900), unmarshaledReq.ExpiresIn)

		resp := PresignedURLResponse{
			URL:       "https://s3.example.com/bucket/uploads/data.csv?signature=xxx",
			Method:    "PUT",
			Key:       "uploads/data.csv",
			ExpiresIn: 900,
		}
		respBytes, err := json.Marshal(resp)
		require.NoError(t, err)

		var unmarshaledResp PresignedURLResponse
		err = json.Unmarshal(respBytes, &unmarshaledResp)
		require.NoError(t, err)
		assert.Equal(t, resp.URL, unmarshaledResp.URL)
		assert.Equal(t, "PUT", unmarshaledResp.Method)

		item := ObjectItem{
			Key:          "uploads/data.csv",
			Size:         2048,
			LastModified: now,
			ETag:         "\"d41d8cd98f00b204e9800998ecf8427e\"",
			ContentType:  "text/csv",
			StorageClass: "STANDARD",
		}
		itemBytes, err := json.Marshal(item)
		require.NoError(t, err)

		var unmarshaledItem ObjectItem
		err = json.Unmarshal(itemBytes, &unmarshaledItem)
		require.NoError(t, err)
		assert.Equal(t, item.Key, unmarshaledItem.Key)
		assert.Equal(t, item.Size, unmarshaledItem.Size)
	})
}
