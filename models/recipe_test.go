package models

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestResourceAttachmentJSON(t *testing.T) {
	targetRecipeID := uuid.New()

	tests := []struct {
		name       string
		attachment ResourceAttachment
	}{
		{
			name: "full database attachment serialization",
			attachment: ResourceAttachment{
				ID:         "att-" + uuid.New().String(),
				RecipeID:   targetRecipeID,
				RecipeName: "my-postgres-db",
				Type:       "database",
				Engine:     "postgresql",
				Aliases: map[string]string{
					"host":     "DATABASE_HOST",
					"port":     "DATABASE_PORT",
					"user":     "DATABASE_USER",
					"password": "DATABASE_PASSWORD",
					"database": "DATABASE_NAME",
					"uri":      "DATABASE_URL",
				},
				EnvironmentPrefix: "DB_",
			},
		},
		{
			name: "cache attachment with empty prefix",
			attachment: ResourceAttachment{
				ID:         "att-" + uuid.New().String(),
				RecipeID:   targetRecipeID,
				RecipeName: "my-redis-cache",
				Type:       "cache",
				Engine:     "redis",
				Aliases: map[string]string{
					"host": "REDIS_HOST",
					"port": "REDIS_PORT",
					"uri":  "REDIS_URL",
				},
			},
		},
		{
			name: "object storage attachment with bucket id and name",
			attachment: ResourceAttachment{
				ID:         "att-" + uuid.New().String(),
				RecipeID:   uuid.Nil,
				RecipeName: "",
				Type:       "object_storage",
				Engine:     "minio",
				BucketID:   &targetRecipeID,
				BucketName: "dk-test-bucket",
				Aliases: map[string]string{
					"endpoint": "S3_ENDPOINT",
					"bucket":   "S3_BUCKET",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.attachment)
			if err != nil {
				t.Fatalf("unexpected marshal error: %v", err)
			}

			var unmarshaled ResourceAttachment
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}

			if unmarshaled.ID != tt.attachment.ID {
				t.Errorf("expected ID %s, got %s", tt.attachment.ID, unmarshaled.ID)
			}
			if unmarshaled.RecipeID != tt.attachment.RecipeID {
				t.Errorf("expected RecipeID %s, got %s", tt.attachment.RecipeID, unmarshaled.RecipeID)
			}
			if unmarshaled.RecipeName != tt.attachment.RecipeName {
				t.Errorf("expected RecipeName %s, got %s", tt.attachment.RecipeName, unmarshaled.RecipeName)
			}
			if unmarshaled.Type != tt.attachment.Type {
				t.Errorf("expected Type %s, got %s", tt.attachment.Type, unmarshaled.Type)
			}
			if unmarshaled.Engine != tt.attachment.Engine {
				t.Errorf("expected Engine %s, got %s", tt.attachment.Engine, unmarshaled.Engine)
			}
			if unmarshaled.EnvironmentPrefix != tt.attachment.EnvironmentPrefix {
				t.Errorf("expected EnvironmentPrefix %s, got %s", tt.attachment.EnvironmentPrefix, unmarshaled.EnvironmentPrefix)
			}
			if len(unmarshaled.Aliases) != len(tt.attachment.Aliases) {
				t.Errorf("expected %d aliases, got %d", len(tt.attachment.Aliases), len(unmarshaled.Aliases))
			}
			if (tt.attachment.BucketID == nil && unmarshaled.BucketID != nil) || (tt.attachment.BucketID != nil && (unmarshaled.BucketID == nil || *unmarshaled.BucketID != *tt.attachment.BucketID)) {
				t.Errorf("expected BucketID %v, got %v", tt.attachment.BucketID, unmarshaled.BucketID)
			}
			if unmarshaled.BucketName != tt.attachment.BucketName {
				t.Errorf("expected BucketName %s, got %s", tt.attachment.BucketName, unmarshaled.BucketName)
			}
		})
	}
}

func TestRecipeWithResourceAttachmentsJSON(t *testing.T) {
	recipeID := uuid.New()
	targetID := uuid.New()

	recipe := Recipe{
		ID:   recipeID,
		Name: "frontend-app",
		ResourceAttachments: []ResourceAttachment{
			{
				ID:         "att-1",
				RecipeID:   targetID,
				RecipeName: "backend-db",
				Type:       "database",
				Engine:     "postgresql",
				Aliases: map[string]string{
					"uri": "DATABASE_URL",
				},
			},
		},
	}

	data, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("failed to marshal recipe: %v", err)
	}

	var unmarshaled Recipe
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal recipe: %v", err)
	}

	if len(unmarshaled.ResourceAttachments) != 1 {
		t.Fatalf("expected 1 resource attachment, got %d", len(unmarshaled.ResourceAttachments))
	}
	if unmarshaled.ResourceAttachments[0].RecipeName != "backend-db" {
		t.Errorf("expected attachment name 'backend-db', got '%s'", unmarshaled.ResourceAttachments[0].RecipeName)
	}

	// Also verify on DisplayRecipe
	display := DisplayRecipe{
		Recipe:       recipe,
		SyncStatus:   "Synced",
		HealthStatus: "Healthy",
	}

	dispData, err := json.Marshal(display)
	if err != nil {
		t.Fatalf("failed to marshal display recipe: %v", err)
	}

	var unmarshaledDisp DisplayRecipe
	if err := json.Unmarshal(dispData, &unmarshaledDisp); err != nil {
		t.Fatalf("failed to unmarshal display recipe: %v", err)
	}

	if len(unmarshaledDisp.ResourceAttachments) != 1 {
		t.Fatalf("expected 1 resource attachment on display recipe, got %d", len(unmarshaledDisp.ResourceAttachments))
	}
	if unmarshaledDisp.ResourceAttachments[0].RecipeID != targetID {
		t.Errorf("expected target ID %s, got %s", targetID, unmarshaledDisp.ResourceAttachments[0].RecipeID)
	}
}
