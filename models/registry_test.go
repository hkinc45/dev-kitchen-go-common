package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainerRegistry_JSONSerialization(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	regID := uuid.New()
	projID := uuid.New()
	regPath := "projects/test-proj/registry"
	regionID := uuid.New()

	tests := []struct {
		name     string
		registry ContainerRegistry
	}{
		{
			name: "full global registry",
			registry: ContainerRegistry{
				ID:        regID,
				Name:      "Harbor Global",
				Type:      "harbor",
				URL:       "https://harbor.example.com",
				IsDefault: true,
				Scope:     "global",
				Status:    "active",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "project custom registry with vault and region",
			registry: ContainerRegistry{
				ID:              regID,
				Name:            "Custom ECR",
				Type:            "ecr",
				URL:             "123456789.dkr.ecr.us-east-1.amazonaws.com",
				IsDefault:       false,
				Scope:           "project",
				ProjectID:       &projID,
				VaultSecretPath: &regPath,
				RegionID:        &regionID,
				Status:          "active",
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.registry)
			require.NoError(t, err)

			var deserialized ContainerRegistry
			err = json.Unmarshal(data, &deserialized)
			require.NoError(t, err)

			assert.Equal(t, tt.registry.ID, deserialized.ID)
			assert.Equal(t, tt.registry.Name, deserialized.Name)
			assert.Equal(t, tt.registry.Type, deserialized.Type)
			assert.Equal(t, tt.registry.URL, deserialized.URL)
			assert.Equal(t, tt.registry.IsDefault, deserialized.IsDefault)
			assert.Equal(t, tt.registry.Scope, deserialized.Scope)
			assert.Equal(t, tt.registry.Status, deserialized.Status)
			assert.True(t, tt.registry.CreatedAt.Equal(deserialized.CreatedAt))
		})
	}
}

func TestProjectRegistryBinding_JSONSerialization(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	binding := ProjectRegistryBinding{
		ID:         uuid.New(),
		ProjectID:  uuid.New(),
		RegistryID: uuid.New(),
		IsPrimary:  true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	data, err := json.Marshal(binding)
	require.NoError(t, err)

	var deserialized ProjectRegistryBinding
	err = json.Unmarshal(data, &deserialized)
	require.NoError(t, err)

	assert.Equal(t, binding.ID, deserialized.ID)
	assert.Equal(t, binding.ProjectID, deserialized.ProjectID)
	assert.Equal(t, binding.RegistryID, deserialized.RegistryID)
	assert.Equal(t, binding.IsPrimary, deserialized.IsPrimary)
}

func TestRobotAccount_JSONSerialization(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	exp := now.Add(30 * 24 * time.Hour)

	tests := []struct {
		name  string
		robot RobotAccount
	}{
		{
			name: "pull only robot account",
			robot: RobotAccount{
				ID:          101,
				Name:        "robot$org-123+ci-pull",
				Description: "CI Pull token",
				ProjectID:   "org-123",
				Scope:       "pull",
				Permissions: "pull",
				CreatedAt:   now,
				ExpiresAt:   &exp,
				Disabled:    false,
			},
		},
		{
			name: "push_pull disabled robot account",
			robot: RobotAccount{
				ID:          102,
				Name:        "robot$org-123+dev-push",
				Description: "Developer CLI token",
				ProjectID:   "org-123",
				Scope:       "push_pull",
				Permissions: "push_pull",
				CreatedAt:   now,
				Disabled:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.robot)
			require.NoError(t, err)

			var deserialized RobotAccount
			err = json.Unmarshal(data, &deserialized)
			require.NoError(t, err)

			assert.Equal(t, tt.robot.ID, deserialized.ID)
			assert.Equal(t, tt.robot.Name, deserialized.Name)
			assert.Equal(t, tt.robot.ProjectID, deserialized.ProjectID)
			assert.Equal(t, tt.robot.Permissions, deserialized.Permissions)
			assert.Equal(t, tt.robot.Disabled, deserialized.Disabled)
		})
	}
}

func TestRegistryArtifact_JSONSerialization(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	recipeID := uuid.New()
	recipeName := "web-service"
	commitHash := "abcd1234efgh"
	commitMsg := "feat: update container build"
	author := "Developer <dev@example.com>"

	artifact := RegistryArtifact{
		Digest:        "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Tags:          []string{"v1.0.0", "latest"},
		Size:          104857600,
		PushTime:      now,
		PullTime:      &now,
		Status:        "active",
		RecipeID:      &recipeID,
		RecipeName:    &recipeName,
		CommitHash:    &commitHash,
		CommitMessage: &commitMsg,
		CommitAuthor:  &author,
		Vulnerabilities: &VulnerabilitySummary{
			Total:     12,
			Critical:  1,
			High:      3,
			Medium:    5,
			Low:       3,
			Fixable:   4,
			ScannedAt: &now,
		},
	}

	data, err := json.Marshal(artifact)
	require.NoError(t, err)

	var deserialized RegistryArtifact
	err = json.Unmarshal(data, &deserialized)
	require.NoError(t, err)

	assert.Equal(t, artifact.Digest, deserialized.Digest)
	assert.Equal(t, artifact.Tags, deserialized.Tags)
	assert.Equal(t, artifact.Size, deserialized.Size)
	assert.Equal(t, artifact.Status, deserialized.Status)
	assert.Equal(t, *artifact.RecipeID, *deserialized.RecipeID)
	assert.Equal(t, *artifact.RecipeName, *deserialized.RecipeName)
	assert.Equal(t, artifact.Vulnerabilities.Critical, deserialized.Vulnerabilities.Critical)
	assert.Equal(t, artifact.Vulnerabilities.Fixable, deserialized.Vulnerabilities.Fixable)
}
