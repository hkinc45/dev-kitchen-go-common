package metadata

// Canonical Dev Kitchen Kubernetes metadata labels and annotations.
const (
	// Standard domain-scoped labels
	LabelOwnerID     = "dev.kitchen/owner-id"
	LabelProjectName = "dev.kitchen/project-name"
	LabelRecipeID    = "dev.kitchen/recipe-id"
	LabelRecipeName  = "dev.kitchen/recipe-name"

	// Legacy un-prefixed labels (for backward compatibility / fallback resolution during migration)
	LegacyLabelOwnerID     = "owner-id"
	LegacyLabelProjectName = "project-name"
	LegacyLabelRecipeID    = "recipe-id"
	LegacyLabelRecipeName  = "recipe-name"

	// Standard annotations
	AnnotationOwnerID     = "dev.kitchen/owner-id"
	AnnotationProjectName = "dev.kitchen/project-name"
	AnnotationRecipeID    = "dev.kitchen/recipe-id"
	AnnotationRecipeName  = "dev.kitchen/recipe-name"
)

// PlatformMetadata encapsulates the core identity attributes of a Dev Kitchen workload.
type PlatformMetadata struct {
	OwnerID     string
	ProjectName string
	RecipeID    string
	RecipeName  string
}

// StandardLabels returns the canonical dev.kitchen labels (no double-labeling).
func (m PlatformMetadata) StandardLabels() map[string]string {
	labels := make(map[string]string)
	if m.OwnerID != "" {
		labels[LabelOwnerID] = m.OwnerID
	}
	if m.ProjectName != "" {
		labels[LabelProjectName] = m.ProjectName
	}
	if m.RecipeID != "" {
		labels[LabelRecipeID] = m.RecipeID
	}
	if m.RecipeName != "" {
		labels[LabelRecipeName] = m.RecipeName
	}
	return labels
}

// StandardAnnotations returns the canonical dev.kitchen annotations.
func (m PlatformMetadata) StandardAnnotations() map[string]string {
	annotations := make(map[string]string)
	if m.OwnerID != "" {
		annotations[AnnotationOwnerID] = m.OwnerID
	}
	if m.ProjectName != "" {
		annotations[AnnotationProjectName] = m.ProjectName
	}
	if m.RecipeID != "" {
		annotations[AnnotationRecipeID] = m.RecipeID
	}
	if m.RecipeName != "" {
		annotations[AnnotationRecipeName] = m.RecipeName
	}
	return annotations
}

// IsLegacyLabel returns true if the given key is one of the deprecated un-prefixed labels.
func IsLegacyLabel(key string) bool {
	return key == LegacyLabelOwnerID || key == LegacyLabelProjectName || key == LegacyLabelRecipeID || key == LegacyLabelRecipeName
}

// InjectPlatformMetadata injects standard labels and annotations into valuesMap,
// stripping deprecated flat keys and preserving any custom user-defined labels/annotations.
func InjectPlatformMetadata(valuesMap map[string]interface{}, meta PlatformMetadata, existingValues map[string]interface{}, isMarketplace bool) {
	if valuesMap == nil {
		return
	}

	labels := meta.StandardLabels()
	annotations := meta.StandardAnnotations()

	// 1. commonLabels
	clMap := make(map[string]interface{})
	if existingValues != nil {
		if m, ok := existingValues["commonLabels"].(map[string]interface{}); ok {
			for k, v := range m {
				if !IsLegacyLabel(k) {
					clMap[k] = v
				}
			}
		}
	}
	if m, ok := valuesMap["commonLabels"].(map[string]interface{}); ok {
		for k, v := range m {
			if !IsLegacyLabel(k) {
				clMap[k] = v
			}
		}
	}
	for k, v := range labels {
		clMap[k] = v
	}
	valuesMap["commonLabels"] = clMap

	// 2. podLabels
	plMap := make(map[string]interface{})
	if existingValues != nil {
		if m, ok := existingValues["podLabels"].(map[string]interface{}); ok {
			for k, v := range m {
				if !IsLegacyLabel(k) {
					plMap[k] = v
				}
			}
		}
	}
	if m, ok := valuesMap["podLabels"].(map[string]interface{}); ok {
		for k, v := range m {
			if !IsLegacyLabel(k) {
				plMap[k] = v
			}
		}
	}
	for k, v := range labels {
		plMap[k] = v
	}
	valuesMap["podLabels"] = plMap

	// 3. commonAnnotations
	caMap := make(map[string]interface{})
	if existingValues != nil {
		if m, ok := existingValues["commonAnnotations"].(map[string]interface{}); ok {
			for k, v := range m {
				caMap[k] = v
			}
		}
	}
	if m, ok := valuesMap["commonAnnotations"].(map[string]interface{}); ok {
		for k, v := range m {
			caMap[k] = v
		}
	}
	for k, v := range annotations {
		caMap[k] = v
	}
	valuesMap["commonAnnotations"] = caMap

	// 4. podAnnotations
	paMap := make(map[string]interface{})
	if existingValues != nil {
		if m, ok := existingValues["podAnnotations"].(map[string]interface{}); ok {
			for k, v := range m {
				paMap[k] = v
			}
		}
	}
	if m, ok := valuesMap["podAnnotations"].(map[string]interface{}); ok {
		for k, v := range m {
			paMap[k] = v
		}
	}
	for k, v := range annotations {
		paMap[k] = v
	}
	valuesMap["podAnnotations"] = paMap

	// 5. For marketplace apps or sub-charts, strip any legacy labels from nested chart blocks if present
	for _, v := range valuesMap {
		if subMap, ok := v.(map[string]interface{}); ok {
			for _, key := range []string{"commonLabels", "podLabels"} {
				if lbls, exists := subMap[key].(map[string]interface{}); exists {
					for lk := range lbls {
						if IsLegacyLabel(lk) {
							delete(lbls, lk)
						}
					}
					for k, v := range labels {
						lbls[k] = v
					}
				}
			}
			for _, key := range []string{"commonAnnotations", "podAnnotations"} {
				if annos, exists := subMap[key].(map[string]interface{}); exists {
					for k, v := range annotations {
						annos[k] = v
					}
				}
			}
		}
	}

	// 6. Non-marketplace recipes wrap under DKAHC
	if !isMarketplace {
		var dkahcMap map[string]interface{}
		if dVal, ok := valuesMap["DKAHC"].(map[string]interface{}); ok {
			dkahcMap = dVal
		} else {
			dkahcMap = make(map[string]interface{})
		}
		dkahcMap["podLabels"] = plMap
		dkahcMap["commonLabels"] = clMap
		dkahcMap["podAnnotations"] = paMap
		dkahcMap["commonAnnotations"] = caMap
		valuesMap["DKAHC"] = dkahcMap
	}
}
