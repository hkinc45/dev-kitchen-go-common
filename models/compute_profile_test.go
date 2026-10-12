package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferWorkloadProfileSKU(t *testing.T) {
	tests := []struct {
		name        string
		cpu         int64
		ramBytes    int64
		expectedSKU string
	}{
		{
			name:        "Unspecified 0/0 defaults to 1cpu-1gb",
			cpu:         0,
			ramBytes:    0,
			expectedSKU: "recipe.profile.1cpu-1gb",
		},
		{
			name:        "Negative values default to 1cpu-1gb",
			cpu:         -10,
			ramBytes:    -500,
			expectedSKU: "recipe.profile.1cpu-1gb",
		},
		{
			name:        "0.5cpu 512mb exact",
			cpu:         500,
			ramBytes:    512 * 1024 * 1024,
			expectedSKU: "recipe.profile.0.5cpu-512mb",
		},
		{
			name:        "0.25cpu 256mb fits in 0.5cpu-512mb",
			cpu:         250,
			ramBytes:    256 * 1024 * 1024,
			expectedSKU: "recipe.profile.0.5cpu-512mb",
		},
		{
			name:        "1cpu 1gb exact",
			cpu:         1000,
			ramBytes:    1024 * 1024 * 1024,
			expectedSKU: "recipe.profile.1cpu-1gb",
		},
		{
			name:        "1cpu 2gb exact",
			cpu:         1000,
			ramBytes:    2048 * 1024 * 1024,
			expectedSKU: "recipe.profile.1cpu-2gb",
		},
		{
			name:        "2cpu 4gb exact",
			cpu:         2000,
			ramBytes:    4096 * 1024 * 1024,
			expectedSKU: "recipe.profile.2cpu-4gb",
		},
		{
			name:        "4cpu 8gb exact",
			cpu:         4000,
			ramBytes:    8192 * 1024 * 1024,
			expectedSKU: "recipe.profile.4cpu-8gb",
		},
		{
			name:        "Above 4cpu 8gb clamps to 4cpu-8gb",
			cpu:         8000,
			ramBytes:    16384 * 1024 * 1024,
			expectedSKU: "recipe.profile.4cpu-8gb",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sku := InferWorkloadProfileSKU(tc.cpu, tc.ramBytes)
			assert.Equal(t, tc.expectedSKU, sku)
		})
	}
}

func TestResolveAllocationMultiplier(t *testing.T) {
	tests := []struct {
		allocType string
		expected  float64
	}{
		{"lightweight", 0.5},
		{"Lightweight", 0.5},
		{"LIGHTWEIGHT", 0.5},
		{"performance", 2.0},
		{"Performance", 2.0},
		{"general_purpose", 1.0},
		{"standard", 1.0},
		{"", 1.0},
	}

	for _, tc := range tests {
		t.Run(tc.allocType, func(t *testing.T) {
			m := ResolveAllocationMultiplier(tc.allocType)
			assert.Equal(t, tc.expected, m)
		})
	}
}
