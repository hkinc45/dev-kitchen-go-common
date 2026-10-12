package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInfrastructureModels(t *testing.T) {
	t.Run("Provider Model Fields", func(t *testing.T) {
		id := uuid.New()
		desc := "MainOne Colocation Facility"
		now := time.Now()

		provider := Provider{
			ID:           id,
			Name:         "MainOne MDXi",
			Slug:         "mainone-colo",
			ProviderType: "colo",
			Description:  &desc,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		assert.Equal(t, id, provider.ID)
		assert.Equal(t, "MainOne MDXi", provider.Name)
		assert.Equal(t, "mainone-colo", provider.Slug)
		assert.Equal(t, "colo", provider.ProviderType)
		assert.True(t, provider.IsActive)
		assert.Equal(t, desc, *provider.Description)
	})

	t.Run("Region Model Fields", func(t *testing.T) {
		id := uuid.New()
		provID := uuid.New()
		now := time.Now()

		region := Region{
			ID:          id,
			ProviderID:  &provID,
			Name:        "Lagos MainOne MDXi",
			Code:        "LOS-01",
			CountryCode: "NG",
			City:        "Lagos",
			Status:      "active",
			IsEnabled:   true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		assert.Equal(t, id, region.ID)
		assert.Equal(t, provID, *region.ProviderID)
		assert.Equal(t, "LOS-01", region.Code)
		assert.Equal(t, "NG", region.CountryCode)
		assert.True(t, region.IsEnabled)
	})

	t.Run("Cluster Model and Load Metrics", func(t *testing.T) {
		clusterID := uuid.New()
		regionID := uuid.New()
		now := time.Now()
		host := "192.168.1.100"

		cluster := Cluster{
			ID:                   clusterID,
			Name:                 "k8s-mainone-01",
			ServerURL:            "https://192.168.1.100:6443",
			RegionID:             regionID,
			Labels:               map[string]string{"env": "production"},
			Status:               "active",
			HealthStatus:         "healthy",
			MaxWorkloads:         100,
			CurrentWorkloadCount: 25,
			CPUCapacityM:         32000,
			CPUAllocatedM:        8000,
			MemoryCapacityMi:     131072,
			MemoryAllocatedMi:    32768,
			IsCordoned:           false,
			AutoscalingEnabled:   true,
			HighWatermarkPercent: 80,
			AnsibleHost:          &host,
			LastHeartbeatAt:      &now,
			CreatedAt:            now,
			UpdatedAt:            now,
		}

		assert.Equal(t, clusterID, cluster.ID)
		assert.Equal(t, "k8s-mainone-01", cluster.Name)
		assert.Equal(t, 25, cluster.CurrentWorkloadCount)
		assert.False(t, cluster.IsCordoned)
		assert.Equal(t, 80, cluster.HighWatermarkPercent)

		metrics := ClusterLoadMetrics{
			ClusterID:          clusterID,
			WorkloadCount:      cluster.CurrentWorkloadCount,
			CPUUsagePercent:    25.0,
			MemoryUsagePercent: 25.0,
			Score:              0.25,
		}

		assert.Equal(t, clusterID, metrics.ClusterID)
		assert.Equal(t, 0.25, metrics.Score)
	})
}
