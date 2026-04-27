package desiredstate

import (
	"testing"
)

func TestDesiredState_Validation(t *testing.T) {
	tests := []struct {
		name    string
		ds      DesiredState
		wantErr bool
	}{
		{
			name: "empty desired state is valid",
			ds: DesiredState{
				Resources: []ResourceExpectation{},
			},
			wantErr: false,
		},
		{
			name: "valid resource expectation",
			ds: DesiredState{
				Resources: []ResourceExpectation{
					{
						Address: "local_file.example",
						Actions: []string{"create"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple resources with different actions",
			ds: DesiredState{
				Resources: []ResourceExpectation{
					{
						Address: "local_file.example",
						Actions: []string{"create"},
					},
					{
						Address: "local_file.other",
						Actions: []string{"update"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "resource with multiple actions",
			ds: DesiredState{
				Resources: []ResourceExpectation{
					{
						Address: "local_file.example",
						Actions: []string{"delete", "create"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "forbidden actions list",
			ds: DesiredState{
				Resources:        []ResourceExpectation{},
				ForbiddenActions: []string{"delete"},
			},
			wantErr: false,
		},
		{
			name: "resource without address",
			ds: DesiredState{
				Resources: []ResourceExpectation{
					{
						Address: "",
						Actions: []string{"create"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "resource without actions",
			ds: DesiredState{
				Resources: []ResourceExpectation{
					{
						Address: "local_file.example",
						Actions: []string{},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ds.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DesiredState.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComparisonResult_Basic(t *testing.T) {
	tests := []struct {
		name   string
		result ComparisonResult
		want   string
	}{
		{
			name: "not_implemented status",
			result: ComparisonResult{
				OK:         true,
				Status:     "not_implemented",
				Matched:    false,
				Mismatches: []Mismatch{},
				Warnings:   []string{"Desired-state comparison is stubbed in this version."},
			},
			want: "not_implemented",
		},
		{
			name: "matched status",
			result: ComparisonResult{
				OK:         true,
				Status:     "matched",
				Matched:    true,
				Mismatches: []Mismatch{},
				Warnings:   []string{},
			},
			want: "matched",
		},
		{
			name: "mismatched status",
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
			want: "mismatched",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.Status != tt.want {
				t.Errorf("ComparisonResult.Status = %v, want %v", tt.result.Status, tt.want)
			}
		})
	}
}
