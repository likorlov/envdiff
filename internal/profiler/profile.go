package profiler

import (
	"fmt"
	"sort"
	"strings"
)

// Profile represents a named environment profile with required and optional keys.
type Profile struct {
	Name     string
	Required []string
	Optional []string
}

// Violation describes a profile conformance issue.
type Violation struct {
	Key     string
	Message string
}

// Check validates an env map against a Profile, returning any violations.
func Check(env map[string]string, p Profile) []Violation {
	var violations []Violation

	for _, key := range p.Required {
		val, ok := env[key]
		if !ok {
			violations = append(violations, Violation{
				Key:     key,
				Message: "required key is missing",
			})
		} else if strings.TrimSpace(val) == "" {
			violations = append(violations, Violation{
				Key:     key,
				Message: "required key has empty value",
			})
		}
	}

	return violations
}

// Coverage returns the percentage of optional keys present in env (0–100).
func Coverage(env map[string]string, p Profile) float64 {
	if len(p.Optional) == 0 {
		return 100.0
	}
	present := 0
	for _, key := range p.Optional {
		if _, ok := env[key]; ok {
			present++
		}
	}
	return float64(present) / float64(len(p.Optional)) * 100.0
}

// FormatViolations returns a human-readable summary of profile violations.
func FormatViolations(violations []Violation) string {
	if len(violations) == 0 {
		return "no profile violations"
	}
	sorted := make([]Violation, len(violations))
	copy(sorted, violations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	var sb strings.Builder
	for _, v := range sorted {
		fmt.Fprintf(&sb, "  [%s] %s\n", v.Key, v.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}
