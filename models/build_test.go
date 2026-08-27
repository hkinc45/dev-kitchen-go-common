package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBuildInfoJSON(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	tests := []struct {
		name     string
		info     BuildInfo
		expected string
	}{
		{
			name: "full struct serialization",
			info: BuildInfo{
				WorkflowRunID: 101,
				CommitSHA:     "abc123def456",
				Status:        BuildStatusSuccess,
				Conclusion:    "success",
				Duration:      "1m 30s",
				HTMLURL:       "http://gitea.local/owner/repo/actions/runs/101",
				CreatedAt:     &now,
			},
		},
		{
			name: "empty struct serialization",
			info: BuildInfo{
				WorkflowRunID: 0,
				Status:        BuildStatusUnknown,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.info)
			if err != nil {
				t.Fatalf("unexpected marshal error: %v", err)
			}

			var unmarshaled BuildInfo
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}

			if unmarshaled.WorkflowRunID != tt.info.WorkflowRunID {
				t.Errorf("expected WorkflowRunID %d, got %d", tt.info.WorkflowRunID, unmarshaled.WorkflowRunID)
			}
			if unmarshaled.Status != tt.info.Status {
				t.Errorf("expected Status %s, got %s", tt.info.Status, unmarshaled.Status)
			}
		})
	}
}
