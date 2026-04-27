package diagnostics

import "encoding/json"

// validateJSONOutput mirrors the relevant subset of `terraform validate -json`.
type validateJSONOutput struct {
	Diagnostics []validateJSONDiagnostic `json:"diagnostics"`
}

type validateJSONDiagnostic struct {
	Severity string             `json:"severity"`
	Summary  string             `json:"summary"`
	Detail   string             `json:"detail"`
	Range    *validateJSONRange `json:"range,omitempty"`
}

type validateJSONRange struct {
	Filename string             `json:"filename"`
	Start    validateJSONOffset `json:"start"`
}

type validateJSONOffset struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// ParseValidateJSON parses Terraform validate -json output and returns normalized
// Diagnostic structures. If JSON parsing fails, it returns a best-effort fallback
// diagnostic containing the stderr or stdout content.
func ParseValidateJSON(stdout, stderr string) []Diagnostic {
	var out validateJSONOutput
	if err := json.Unmarshal([]byte(stdout), &out); err != nil {
		// Return a best-effort fallback diagnostic so callers always get
		// actionable feedback even when structured JSON is unavailable.
		// Severity is always "error" here because a JSON parse failure means
		// we cannot present structured diagnostics — the raw output may itself
		// contain errors that Terraform failed to encode as JSON.
		fallback := Diagnostic{Severity: "error"}
		if stderr != "" {
			fallback.Summary = stderr
		} else {
			fallback.Summary = stdout
		}
		return []Diagnostic{fallback}
	}
	diags := make([]Diagnostic, 0, len(out.Diagnostics))
	for _, d := range out.Diagnostics {
		nd := Diagnostic{
			Severity: d.Severity,
			Summary:  d.Summary,
			Detail:   d.Detail,
		}
		if d.Range != nil {
			nd.File = d.Range.Filename
			nd.Line = d.Range.Start.Line
			nd.Column = d.Range.Start.Column
		}
		diags = append(diags, nd)
	}
	return diags
}
