package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeploymentIssueEvent_JSON(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name     string
		event    DeploymentIssueEvent
		expected string
	}{
		{
			name: "CrashLoopBackOff Event",
			event: DeploymentIssueEvent{
				Source:         SourceKubernetes,
				Severity:       SeverityCritical,
				Reason:         "CrashLoopBackOff",
				Message:        "Back-off 5m0s restarting failed container web",
				InvolvedObject: "pod-123",
				Timestamp:      now,
				Count:          5,
			},
		},
		{
			name: "ArgoCD Sync Failure",
			event: DeploymentIssueEvent{
				Source:         SourceArgoCD,
				Severity:       SeverityError,
				Reason:         "SyncError",
				Message:        "Application destination server invalid",
				InvolvedObject: "app-abc",
				Timestamp:      now,
				Count:          1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.event)
			require.NoError(t, err)

			var unmarshaled DeploymentIssueEvent
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, tt.event.Source, unmarshaled.Source)
			assert.Equal(t, tt.event.Severity, unmarshaled.Severity)
			assert.Equal(t, tt.event.Reason, unmarshaled.Reason)
			assert.Equal(t, tt.event.Message, unmarshaled.Message)
			assert.Equal(t, tt.event.InvolvedObject, unmarshaled.InvolvedObject)
			assert.Equal(t, tt.event.Count, unmarshaled.Count)
			assert.True(t, tt.event.Timestamp.Equal(unmarshaled.Timestamp))
		})
	}
}

func TestContainerMetricTelemetry_JSON(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	telemetry := ContainerMetricTelemetry{
		ContainerName:        "web-service",
		CPUUsageMillicores:   15,
		CPURequestMillicores: 100,
		CPULimitMillicores:   500,
		MemoryUsageBytes:     104857600,
		MemoryUsageMiB:       100.0,
		MemoryRequestMiB:     256.0,
		MemoryLimitMiB:       512.0,
		Timestamp:            now,
	}

	data, err := json.Marshal(telemetry)
	require.NoError(t, err)

	var unmarshaled ContainerMetricTelemetry
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, telemetry.ContainerName, unmarshaled.ContainerName)
	assert.Equal(t, telemetry.CPUUsageMillicores, unmarshaled.CPUUsageMillicores)
	assert.Equal(t, telemetry.CPURequestMillicores, unmarshaled.CPURequestMillicores)
	assert.Equal(t, telemetry.CPULimitMillicores, unmarshaled.CPULimitMillicores)
	assert.Equal(t, telemetry.MemoryUsageBytes, unmarshaled.MemoryUsageBytes)
	assert.Equal(t, telemetry.MemoryUsageMiB, unmarshaled.MemoryUsageMiB)
	assert.Equal(t, telemetry.MemoryRequestMiB, unmarshaled.MemoryRequestMiB)
	assert.Equal(t, telemetry.MemoryLimitMiB, unmarshaled.MemoryLimitMiB)
	assert.True(t, telemetry.Timestamp.Equal(unmarshaled.Timestamp))
}

func TestCIBuildLogResponse_JSON(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name string
		resp CIBuildLogResponse
	}{
		{
			name: "Completed build run with logs",
			resp: CIBuildLogResponse{
				RunID:      42,
				Status:     "completed",
				Conclusion: "success",
				Logs:       "Step 1: Docker build\nSuccessfully built image",
				CreatedAt:  &now,
				Duration:   "45s",
			},
		},
		{
			name: "Pending build without timestamps",
			resp: CIBuildLogResponse{
				RunID:      43,
				Status:     "pending",
				Conclusion: "",
				Logs:       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.resp)
			require.NoError(t, err)

			var unmarshaled CIBuildLogResponse
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, tt.resp.RunID, unmarshaled.RunID)
			assert.Equal(t, tt.resp.Status, unmarshaled.Status)
			assert.Equal(t, tt.resp.Conclusion, unmarshaled.Conclusion)
			assert.Equal(t, tt.resp.Logs, unmarshaled.Logs)
			assert.Equal(t, tt.resp.Duration, unmarshaled.Duration)
			if tt.resp.CreatedAt != nil {
				require.NotNil(t, unmarshaled.CreatedAt)
				assert.True(t, tt.resp.CreatedAt.Equal(*unmarshaled.CreatedAt))
			} else {
				assert.Nil(t, unmarshaled.CreatedAt)
			}
		})
	}
}

func TestLogStreamRequest_JSON(t *testing.T) {
	req := LogStreamRequest{
		Container:  "web",
		TailLines:  500,
		Follow:     true,
		Timestamps: true,
		Previous:   false,
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)

	var unmarshaled LogStreamRequest
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, req, unmarshaled)
}
