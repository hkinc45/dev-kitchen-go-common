package metadata

import (
	"testing"
)

func TestPlatformMetadata_StandardLabels(t *testing.T) {
	meta := PlatformMetadata{
		OwnerID:     "user-123",
		ProjectName: "my-project",
		RecipeID:    "recipe-456",
		RecipeName:  "my-app",
	}

	labels := meta.StandardLabels()
	if labels[LabelOwnerID] != "user-123" {
		t.Errorf("expected %s to be user-123, got %s", LabelOwnerID, labels[LabelOwnerID])
	}
	if labels[LabelProjectName] != "my-project" {
		t.Errorf("expected %s to be my-project, got %s", LabelProjectName, labels[LabelProjectName])
	}
	if labels[LabelRecipeID] != "recipe-456" {
		t.Errorf("expected %s to be recipe-456, got %s", LabelRecipeID, labels[LabelRecipeID])
	}
	if labels[LabelRecipeName] != "my-app" {
		t.Errorf("expected %s to be my-app, got %s", LabelRecipeName, labels[LabelRecipeName])
	}

	// Ensure no legacy double-labels are generated
	if _, ok := labels["owner-id"]; ok {
		t.Errorf("unexpected legacy label owner-id present")
	}
	if _, ok := labels["project-name"]; ok {
		t.Errorf("unexpected legacy label project-name present")
	}
	if _, ok := labels["recipe-id"]; ok {
		t.Errorf("unexpected legacy label recipe-id present")
	}
	if _, ok := labels["recipe-name"]; ok {
		t.Errorf("unexpected legacy label recipe-name present")
	}
}

func TestPlatformMetadata_StandardAnnotations(t *testing.T) {
	meta := PlatformMetadata{
		OwnerID:     "user-123",
		ProjectName: "my-project",
		RecipeID:    "recipe-456",
		RecipeName:  "my-app",
	}

	annotations := meta.StandardAnnotations()
	if annotations[AnnotationOwnerID] != "user-123" {
		t.Errorf("expected %s to be user-123, got %s", AnnotationOwnerID, annotations[AnnotationOwnerID])
	}
	if annotations[AnnotationProjectName] != "my-project" {
		t.Errorf("expected %s to be my-project, got %s", AnnotationProjectName, annotations[AnnotationProjectName])
	}
	if annotations[AnnotationRecipeID] != "recipe-456" {
		t.Errorf("expected %s to be recipe-456, got %s", AnnotationRecipeID, annotations[AnnotationRecipeID])
	}
	if annotations[AnnotationRecipeName] != "my-app" {
		t.Errorf("expected %s to be my-app, got %s", AnnotationRecipeName, annotations[AnnotationRecipeName])
	}
}

func TestIsLegacyLabel(t *testing.T) {
	legacy := []string{"owner-id", "project-name", "recipe-id", "recipe-name"}
	for _, l := range legacy {
		if !IsLegacyLabel(l) {
			t.Errorf("expected %s to be recognized as legacy label", l)
		}
	}

	canonical := []string{LabelOwnerID, LabelProjectName, LabelRecipeID, LabelRecipeName, "custom-label"}
	for _, c := range canonical {
		if IsLegacyLabel(c) {
			t.Errorf("expected %s to NOT be recognized as legacy label", c)
		}
	}
}

func TestInjectPlatformMetadata_Marketplace(t *testing.T) {
	meta := PlatformMetadata{
		OwnerID:     "user-123",
		ProjectName: "my-project",
		RecipeID:    "rec-456",
		RecipeName:  "my-recipe",
	}

	existingValues := map[string]interface{}{
		"commonLabels": map[string]interface{}{
			"owner-id":     "old-user",
			"project-name": "old-proj",
			"recipe-id":    "old-rec",
			"recipe-name":  "old-name",
			"custom-label": "custom-value",
		},
		"podLabels": map[string]interface{}{
			"owner-id": "old-user",
		},
	}

	valuesMap := map[string]interface{}{
		"mysql": map[string]interface{}{
			"commonLabels": map[string]interface{}{
				"owner-id":     "legacy-in-mysql",
				"another-user": "keep-me",
			},
		},
	}

	InjectPlatformMetadata(valuesMap, meta, existingValues, true)

	// Check top-level commonLabels
	cl, ok := valuesMap["commonLabels"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected commonLabels to be a map")
	}
	if cl[LabelOwnerID] != "user-123" || cl[LabelProjectName] != "my-project" || cl[LabelRecipeID] != "rec-456" || cl[LabelRecipeName] != "my-recipe" {
		t.Errorf("expected canonical labels in commonLabels, got %+v", cl)
	}
	if cl["custom-label"] != "custom-value" {
		t.Errorf("expected custom-label preserved")
	}
	if cl["owner-id"] != nil || cl["project-name"] != nil || cl["recipe-id"] != nil || cl["recipe-name"] != nil {
		t.Errorf("expected legacy labels stripped from commonLabels, got %+v", cl)
	}

	// Check nested mysql map
	mysqlMap := valuesMap["mysql"].(map[string]interface{})
	mysqlCL := mysqlMap["commonLabels"].(map[string]interface{})
	if mysqlCL[LabelOwnerID] != "user-123" || mysqlCL["another-user"] != "keep-me" {
		t.Errorf("expected canonical labels in nested mysql commonLabels, got %+v", mysqlCL)
	}
	if mysqlCL["owner-id"] != nil {
		t.Errorf("expected legacy label owner-id stripped from nested mysql map")
	}
}

func TestInjectPlatformMetadata_NonMarketplace(t *testing.T) {
	meta := PlatformMetadata{
		OwnerID:     "user-123",
		ProjectName: "my-project",
		RecipeID:    "rec-456",
		RecipeName:  "my-recipe",
	}

	valuesMap := make(map[string]interface{})
	InjectPlatformMetadata(valuesMap, meta, nil, false)

	// DKAHC should wrap
	dkahc, ok := valuesMap["DKAHC"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected DKAHC to be a map")
	}
	cl, ok := dkahc["commonLabels"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected dkahc commonLabels to be a map")
	}
	if cl[LabelOwnerID] != "user-123" || cl["owner-id"] != nil {
		t.Errorf("expected canonical labels and stripped legacy in DKAHC, got %+v", cl)
	}
}
