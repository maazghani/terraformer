package diagnostics

import "encoding/json"

// PlanSummary aggregates resource change counts from terraform show -json output.
type PlanSummary struct {
	Create  int `json:"create"`
	Update  int `json:"update"`
	Delete  int `json:"delete"`
	Replace int `json:"replace"`
	NoOp    int `json:"no_op"`
}

// showJSONPayload mirrors the relevant subset of `terraform show -json`.
type showJSONPayload struct {
	ResourceChanges []showJSONResourceChange `json:"resource_changes"`
}

type showJSONResourceChange struct {
	Change struct {
		Actions []string `json:"actions"`
	} `json:"change"`
}

// ParsePlanSummary parses Terraform show -json output and returns a normalized
// PlanSummary aggregating resource change counts. If JSON parsing fails, it
// returns a zero-valued PlanSummary without panicking.
func ParsePlanSummary(stdout string) PlanSummary {
	var out showJSONPayload
	if err := json.Unmarshal([]byte(stdout), &out); err != nil {
		return PlanSummary{}
	}
	var s PlanSummary
	for _, rc := range out.ResourceChanges {
		s = applyAction(s, rc.Change.Actions)
	}
	return s
}

func applyAction(s PlanSummary, actions []string) PlanSummary {
	switch {
	case len(actions) == 2 && actions[0] == "delete" && actions[1] == "create":
		s.Replace++
	case len(actions) == 2 && actions[0] == "create" && actions[1] == "delete":
		s.Replace++
	case len(actions) == 1 && actions[0] == "create":
		s.Create++
	case len(actions) == 1 && actions[0] == "update":
		s.Update++
	case len(actions) == 1 && actions[0] == "delete":
		s.Delete++
	case len(actions) == 1 && actions[0] == "no-op":
		s.NoOp++
	}
	return s
}
