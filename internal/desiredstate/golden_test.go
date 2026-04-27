package desiredstate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGolden_ComparisonResults(t *testing.T) {
	goldenDir := filepath.Join("..", "..", "testdata", "golden", "desiredstate")

	tests := []struct {
		name       string
		result     ComparisonResult
		goldenFile string
	}{
		{
			name: "not_implemented",
			result: ComparisonResult{
				OK:         true,
				Status:     "not_implemented",
				Matched:    false,
				Mismatches: []Mismatch{},
				Warnings: []string{
					"Desired-state comparison is stubbed in this version.",
				},
			},
			goldenFile: "not_implemented.json",
		},
		{
			name: "matched",
			result: ComparisonResult{
				OK:         true,
				Status:     "matched",
				Matched:    true,
				Mismatches: []Mismatch{},
				Warnings:   []string{},
			},
			goldenFile: "matched.json",
		},
		{
			name: "mismatched",
			result: ComparisonResult{
				OK:      true,
				Status:  "mismatched",
				Matched: false,
				Mismatches: []Mismatch{
					{
						Address: "local_file.example",
						Reason:  "Expected create but plan contains delete.",
					},
				},
				Warnings: []string{},
			},
			goldenFile: "mismatched.json",
		},
		{
			name: "not_checked",
			result: ComparisonResult{
				OK:         true,
				Status:     "not_checked",
				Matched:    false,
				Mismatches: []Mismatch{},
				Warnings: []string{
					"Desired-state was not checked for this operation.",
				},
			},
			goldenFile: "not_checked.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the result to JSON.
			got, err := json.MarshalIndent(tt.result, "", "  ")
			if err != nil {
				t.Fatalf("failed to marshal result: %v", err)
			}

			goldenPath := filepath.Join(goldenDir, tt.goldenFile)

			// Check if UPDATE_GOLDEN env var is set.
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				if err := os.MkdirAll(goldenDir, 0755); err != nil {
					t.Fatalf("failed to create golden dir: %v", err)
				}
				if err := os.WriteFile(goldenPath, got, 0644); err != nil {
					t.Fatalf("failed to write golden file: %v", err)
				}
				t.Logf("Updated golden file: %s", goldenPath)
				return
			}

			// Read the golden file.
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("failed to read golden file %s: %v (run with UPDATE_GOLDEN=1 to create it)", goldenPath, err)
			}

			// Compare.
			if string(got) != string(want) {
				t.Errorf("result does not match golden file %s\nGot:\n%s\nWant:\n%s", tt.goldenFile, string(got), string(want))
			}
		})
	}
}
