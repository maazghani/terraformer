package desiredstate

import (
	"fmt"
)

// DesiredState represents the expected state of resources in a Terraform plan.
type DesiredState struct {
	Resources        []ResourceExpectation `json:"resources"`
	ForbiddenActions []string              `json:"forbidden_actions,omitempty"`
}

// ResourceExpectation defines the expected actions for a specific resource.
type ResourceExpectation struct {
	Address string   `json:"address"`
	Actions []string `json:"actions"`
}

// ComparisonResult represents the outcome of comparing a plan against desired state.
type ComparisonResult struct {
	OK         bool       `json:"ok"`
	Status     string     `json:"status"`
	Matched    bool       `json:"matched"`
	Mismatches []Mismatch `json:"mismatches"`
	Warnings   []string   `json:"warnings"`
}

// Mismatch represents a single discrepancy between desired state and actual plan.
type Mismatch struct {
	Address string `json:"address"`
	Reason  string `json:"reason"`
}

// Validate checks that the DesiredState is well-formed.
func (ds *DesiredState) Validate() error {
	for i, res := range ds.Resources {
		if res.Address == "" {
			return fmt.Errorf("resource at index %d has empty address", i)
		}
		if len(res.Actions) == 0 {
			return fmt.Errorf("resource %q has no actions specified", res.Address)
		}
	}
	return nil
}
